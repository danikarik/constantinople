package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
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

// RequestsResponseTime -
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

func main() {
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
		Debug:          true,
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Content-Type", "Cookie"},
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler)
	r.Use(RequestsResponseTime())

	r.Mount("/debug", middleware.Profiler())
	r.Get("/metrics", promhttp.Handler().ServeHTTP)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		timing := servertiming.FromContext(r.Context())

		m := timing.NewMetric("sql").WithDesc("SQL query").Start()
		time.Sleep(random(20, 50))
		m.Stop()

		w.WriteHeader(200)
		_, err := w.Write([]byte("Done. Check your browser inspector timing details."))
		if err != nil {
			log.Printf("Can't write http response: %s", err)
		}
	})

	h := servertiming.Middleware(r, nil)

	http.ListenAndServe(":3000", h)
}

func random(min, max int) time.Duration {
	return (time.Duration(rand.Intn(max-min) + min)) * time.Millisecond
}
