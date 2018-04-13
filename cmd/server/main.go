package main

import (
	"context"
	"flag"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/danikarik/constantinople/pkg/auth"
	"github.com/danikarik/constantinople/pkg/util"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/valve"
	"github.com/golang/glog"
	"github.com/mitchellh/go-server-timing"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/cors"
)

var httpRequestsResponseTime = prometheus.NewSummary(prometheus.SummaryOpts{
	Name: "promhttp_metric_handler_requests_response_time",
	Help: "Request response times",
})

func init() {
	prometheus.MustRegister(httpRequestsResponseTime)
}

// RequestsResponseTime observes response time per request.
func RequestsResponseTime() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			defer httpRequestsResponseTime.Observe(float64(time.Since(start).Seconds()))
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func init() {
	flag.Set("logtostderr", "true")
	flag.Set("v", "2")
	flag.Parse()
}

func main() {

	valv := valve.New()
	baseCtx := valv.Context()

	r := chi.NewRouter()

	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.NoCache)

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(middleware.Compress(5))
	r.Use(cors.New(cors.Options{
		Debug:          glog.V(3) == true,
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Content-Type", "Cookie"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler)
	r.Use(RequestsResponseTime())

	auth, err := auth.New(auth.Options{
		PKIAddress: "127.0.0.1:8000",
		Hostname:   "127.0.0.1:6379",
		Password:   "daniyar",
	})
	if err != nil {
		util.Exit("[auth] %s", err.Error())
	}

	r.Mount(auth.Router("/auth"))

	r.Mount("/debug", middleware.Profiler())
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	srv := http.Server{Addr: ":3000", Handler: chi.ServerBaseContext(baseCtx, servertiming.Middleware(r, nil))}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for range c {
			// sig is a ^C, handle it
			util.Info("[server] shutting down..")

			// first valv
			valv.Shutdown(10 * time.Second)

			// create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// start http shutdown
			srv.Shutdown(ctx)

			// verify, in worst case call cancel via defer
			select {
			case <-time.After(11 * time.Second):
				util.Info("[server] not all connections done")
			case <-ctx.Done():

			}
		}
	}()
	util.Info("[server] listening on %s", srv.Addr)
	srv.ListenAndServe()
}

func random(min, max int) time.Duration {
	return (time.Duration(rand.Intn(max-min) + min)) * time.Millisecond
}
