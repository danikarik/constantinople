package metric

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
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
