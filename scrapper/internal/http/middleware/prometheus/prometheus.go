package prometheus

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	httpRequestTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"})

	httpErrorsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of all HTTP errors",
		},
		[]string{"method", "path", "status"})

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration",
			Help:    "Duration of HTTP requests",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"})
)

func PromMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			duration := time.Since(start).Seconds()

			var path string

			rout := chi.RouteContext(r.Context())

			if pattern := rout.RoutePattern(); pattern != "" {
				path = pattern
			} else {
				path = r.URL.Path
			}

			httpRequestTotal.WithLabelValues(r.Method, path, strconv.Itoa(ww.Status())).Inc()
			httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)

			if ww.Status() >= 400 {
				httpErrorsTotal.WithLabelValues(r.Method, path, strconv.Itoa(ww.Status())).Inc()
			}
		}

		return http.HandlerFunc(fn)
	}
}
