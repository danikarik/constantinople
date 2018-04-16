package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/cors"

	"github.com/danikarik/constantinople/pkg/auth"
	"github.com/danikarik/constantinople/pkg/metric"
	"github.com/danikarik/constantinople/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"
	"github.com/golang/glog"
	servertiming "github.com/mitchellh/go-server-timing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
}

// App stands for application container.
type App struct {
	addr string
	mux  *chi.Mux
}

func setup(addr string, options Options) (*auth.Auth, *cors.Cors, error) {
	var (
		defaultOrigins = []string{"*"}
		defaultDebug   = true
	)
	if !glog.V(3) {
		defaultOrigins = options.Origins
		defaultDebug = false
	}
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
		Debug:      defaultDebug,
	})
	if err != nil {
		return nil, nil, err
	}
	crs := cors.New(cors.Options{
		AllowedOrigins:   defaultOrigins,
		AllowedMethods:   []string{"GET", "POST", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Cookie"},
		AllowCredentials: true,
		Debug:            defaultDebug,
	})
	return auth, crs, nil
}

// New is a contructor for App container.
func New(addr string, options Options) (*App, error) {
	auth, crs, err := setup(addr, options)
	if err != nil {
		return nil, err
	}
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
	r.Mount(auth.Router("/auth"))
	r.Mount("/debug", middleware.Profiler())
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
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
