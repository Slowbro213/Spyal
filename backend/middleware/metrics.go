package middleware

import (
	"net/http"
	"spyal/core"
	"spyal/pkg/utils/metrics"
	"time"
)

type responseWriter struct {
	core.Middleware
	responseSize int
}

const (
	successCode     = 200
	errorCodesStart = 400
)

func (rw *responseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.responseSize += size
	return size, err
}

// TrackMetrics wraps an http.Handler and observes Prometheus metrics.
func TrackMetrics(m *metrics.Metrics, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if core.IsWebSocketRequest(r) {
			next.ServeHTTP(w, r)
			return
		}

		m.InFlightRequests.Inc()
		defer m.InFlightRequests.Dec()

		start := time.Now()

		rw := &responseWriter{
			Middleware: core.Middleware{
				ResponseWriter: w,
				StatusCode:     successCode,
			},
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		path := r.URL.Path
		method := r.Method
		statusCode := rw.StatusCode
		status := http.StatusText(statusCode)

		m.TotalRequests.WithLabelValues(method, path, status).Inc()
		m.RequestDuration.WithLabelValues(method, path, status).Observe(duration)
		m.ResponseSize.WithLabelValues(method, path, status).Observe(float64(rw.responseSize))

		if statusCode >= errorCodesStart {
			m.ErrorRequests.WithLabelValues(method, path).Inc()
		}
	})
}
