package middleware

import (
	"net/http"
	"time"
	"spyal/pkg/utils/metrics"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

const (
	successCode    = 200
	errorCodesStart = 400
)

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
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
		m.InFlightRequests.Inc()
		defer m.InFlightRequests.Dec()

		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
			statusCode:     successCode,
		}

		next.ServeHTTP(rw, r)

		duration := time.Since(start).Seconds()

		path := r.URL.Path
		method := r.Method
		statusCode := rw.statusCode
		status := http.StatusText(statusCode)

		m.TotalRequests.WithLabelValues(method, path, status).Inc()
		m.RequestDuration.WithLabelValues(method, path, status).Observe(duration)
		m.ResponseSize.WithLabelValues(method, path, status).Observe(float64(rw.responseSize))

		if statusCode >= errorCodesStart {
			m.ErrorRequests.WithLabelValues(method, path).Inc()
		}
	})
}
