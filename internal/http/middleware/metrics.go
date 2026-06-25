package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	requestsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total number of HTTP requests.",
	}, []string{"method", "path", "status"})

	errorsTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_errors_total",
		Help: "Total number of HTTP requests that resulted in a 4xx or 5xx response.",
	}, []string{"method", "path", "status"})

	responseDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_duration_seconds",
		Help:    "HTTP response latency in seconds.",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "path"})
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(rw.status)
		method := r.Method
		path := r.Pattern

		requestsTotal.WithLabelValues(method, path, status).Inc()
		responseDuration.WithLabelValues(method, path).Observe(duration)

		if rw.status >= 400 {
			errorsTotal.WithLabelValues(method, path, status).Inc()
		}
	})
}
