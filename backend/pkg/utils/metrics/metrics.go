package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	TotalRequests   *prometheus.CounterVec
	RequestDuration *prometheus.HistogramVec
	InFlightRequests prometheus.Gauge
	ErrorRequests   *prometheus.CounterVec
	ResponseSize    *prometheus.HistogramVec
}

const (
	bucketSize = 100 
	buckets = 10 
	bucketGrowth = 2
)

func New() *Metrics {
	m := &Metrics{
		TotalRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		RequestDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_request_duration_seconds",
				Help:    "Histogram of latencies for HTTP requests",
				Buckets: prometheus.DefBuckets,
			},
			[]string{"method", "path", "status"},
		),
		InFlightRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "http_in_flight_requests",
				Help: "Current number of in-flight HTTP requests",
			},
		),
		ErrorRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_error_requests_total",
				Help: "Total number of HTTP error responses (status >= 400)",
			},
			[]string{"method", "path"},
		),
		ResponseSize: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name:    "http_response_size_bytes",
				Help:    "Histogram of HTTP response sizes in bytes",
				Buckets: prometheus.ExponentialBuckets(bucketSize, bucketGrowth, buckets),
			},
			[]string{"method", "path", "status"},
		),
	}

	// Register all metrics
	prometheus.MustRegister(
		m.TotalRequests,
		m.RequestDuration,
		m.InFlightRequests,
		m.ErrorRequests,
		m.ResponseSize,
	)

	return m
}
