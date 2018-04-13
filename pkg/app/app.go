package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/danikarik/constantinople/pkg/auth"
	"github.com/danikarik/constantinople/pkg/metric"
	"github.com/danikarik/constantinople/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"
	"github.com/golang/glog"
	servertiming "github.com/mitchellh/go-server-timing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

var (
	// ErrNoAddress raises when http server address is empty.
	ErrNoAddress = errors.New("port: address to listen not specified")
	// ErrRedisConn raises when cannot connect to Redis host.
	ErrRedisConn = errors.New("redis: connection failed")
)

// Options is a configuration container to setup the application.
type Options struct {
}

// App stands for application container.
type App struct {
	addr string
	mux  *chi.Mux
}

// New is a contructor for App container.
func New(addr string, options Options) (*App, error) {

	if addr == "" {
		return nil, ErrNoAddress
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

	r.Use(cors.New(cors.Options{
		Debug:          glog.V(3) == true,
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Content-Type", "Cookie"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler)

	r.Use(metric.RequestsResponseTime())

	auth, err := auth.New(auth.Options{
		PKIAddress: "127.0.0.1:8000",
		Hostname:   "127.0.0.1:6379",
		Password:   "daniyar",
	})
	if err != nil {
		return nil, ErrRedisConn
	}

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
			util.Info("[server] shutting down application...")
			valv.Shutdown(10 * time.Second)
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
			select {
			case <-time.After(11 * time.Second):
				util.Info("[server] not all connections done...")
			case <-ctx.Done():

			}
		}
	}()
	util.Info("[server] listening on %s", srv.Addr)
	return srv.ListenAndServe()
}
