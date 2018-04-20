package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"

	"github.com/danikarik/constantinople/pkg/auth"
	"github.com/danikarik/constantinople/pkg/metric"
	"github.com/danikarik/constantinople/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"
	"github.com/golang/glog"
	servertiming "github.com/mitchellh/go-server-timing"
)

const (
	testCars = `{"status":"SUCCESS","result":[{"manufacturer":"Toyota","model":"RAV-4","number":"151AMM01","year":1992,"pasport_number":"AS011101010","color":"Серо-зеленый металлик","color_hex":"#4D5645","vin_code":"1GNSCAE03BR377068","engine_capacity":2000,"fines":[],"taxes":[],"deregistration":null,"reregistration":null},{"manufacturer":"BMW","model":"X7","number":"A151APM","year":1992,"pasport_number":"AS011101010","color":"Сливочный белый","color_hex":"#FDF4E3","vin_code":"1GNSCAE03BR377068","engine_capacity":2000,"fines":[{"amount":"25600 KZT","info":"Превышение скоростного режима","is_paid":false}],"taxes":[{"amount":"25600 KZT","info":"Превышение скоростного режима","is_paid":false}],"deregistration":null,"reregistration":null},{"manufacturer":"Lada","model":"Granta","number":"151BCM01","year":1992,"pasport_number":"AS011101010","color":"Серый коричневый","color_hex":"#403A3A","vin_code":"1GNSCAE03BR377068","engine_capacity":2000,"fines":[{"amount":"25600 KZT","info":"Превышение скоростного режима","is_paid":false}],"taxes":[],"deregistration":null,"reregistration":null},{"manufacturer":"Lada","model":"Granta","number":"151DDD01","year":1992,"pasport_number":"AS011101010","color":"Синий кобальт","color_hex":"#1E213D","vin_code":"1GNSCAE03BR377068","engine_capacity":2000,"fines":[],"taxes":[{"amount":"25600 KZT","info":"Превышение скоростного режима","is_paid":false}],"deregistration":null,"reregistration":null}]}`
)

var (
	// ErrHTTPAddress raises when http server address is empty.
	ErrHTTPAddress = errors.New("tcp: address not specified")
	// ErrAllowedOrigins raises when origin whitelist is empty.
	ErrAllowedOrigins = errors.New("cors: origins not specified")
)

// Options is a configuration container to setup the application.
type Options struct {
	Origins     []string
	AuthService string
	RedisHost   string
	RedisPass   string
	Username    string
	Password    string
}

// App stands for application container.
type App struct {
	addr string
	mux  *chi.Mux
}

func setup(addr string, options Options) (*auth.Auth, *cors.Cors, error) {
	var debug = glog.V(3) == true
	if addr == "" {
		return nil, nil, ErrHTTPAddress
	}
	if len(options.Origins) == 0 {
		return nil, nil, ErrAllowedOrigins
	}
	auth, err := auth.New(auth.Options{
		PKIAddress: options.AuthService,
		Hostname:   options.RedisHost,
		Password:   options.RedisPass,
		Debug:      debug,
	})
	if err != nil {
		return nil, nil, err
	}
	crs := cors.New(cors.Options{
		AllowedOrigins:   options.Origins,
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Cookie"},
		AllowCredentials: true,
		Debug:            debug,
	})
	return auth, crs, nil
}

// New is a contructor for App container.
func New(addr string, options Options) (*App, error) {
	ath, crs, err := setup(addr, options)
	if err != nil {
		return nil, err
	}
	basicAuth := auth.BasicAuth(func(user, pass string) bool {
		if user == options.Username && pass == options.Password {
			return true
		}
		return false
	})
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.NoCache)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5))
	r.Use(crs.Handler)
	r.Use(metric.RequestsResponseTime())
	r.Mount(ath.Router("/auth"))
	r.With(basicAuth).Mount("/debug", middleware.Profiler())
	r.With(basicAuth).Mount("/metrics", promhttp.Handler())
	// DELETE IN FUTURE
	r.With(ath.Middleware()).Get("/cars", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(testCars))
	})
	return &App{
		addr: addr,
		mux:  r,
	}, nil
}

// Serve the application at the specified address/port and listen for OS
// interrupt and kill signals and will attempt to stop the application
// gracefully.
func (a *App) Serve() error {
	valv := valve.New()
	baseCtx := valv.Context()
	graceStop := make(chan os.Signal, 1)
	signal.Notify(graceStop, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL)
	srv := http.Server{
		Addr: a.addr,
		Handler: chi.ServerBaseContext(
			baseCtx,
			servertiming.Middleware(a.mux, nil),
		),
	}
	go func() {
		for range graceStop {
			util.Info("shutting down application...")
			valv.Shutdown(10 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
			select {
			case <-time.After(11 * time.Second):
				util.Info("not all connections done...")
			case <-ctx.Done():

			}
		}
	}()
	util.Info("listening on %s", srv.Addr)
	return srv.ListenAndServe()
}
