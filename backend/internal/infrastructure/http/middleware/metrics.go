package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal    *prometheus.CounterVec
	httpRequestDuration  *prometheus.HistogramVec
	metricsOnce          sync.Once
)

func initMetrics() {
	metricsOnce.Do(func() {
		httpRequestsTotal = prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests by method, route, and status code.",
			},
			[]string{"method", "route", "status"},
		)
		httpRequestDuration = prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "HTTP request duration in seconds by method and route.",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "route"},
		)
		prometheus.MustRegister(httpRequestsTotal, httpRequestDuration)
	})
}

// PrometheusMiddleware records per-route request counts and durations.
// It uses Chi's route pattern (e.g. /articles/{slug}) rather than the raw
// URL path to keep label cardinality bounded.
func PrometheusMiddleware(next http.Handler) http.Handler {
	initMetrics()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		route := chi.RouteContext(r.Context()).RoutePattern()
		if route == "" {
			route = r.URL.Path
		}

		status := strconv.Itoa(ww.Status())
		httpRequestsTotal.WithLabelValues(r.Method, route, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, route).Observe(time.Since(start).Seconds())
	})
}
