package metric

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// counter for total http requests
	HTTPRequestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests with method, path and status",
		},
		[]string{"method", "path", "status"},
	)

	// histogram for http requests duration
	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// counter for total grpc requests
	GRPCClientRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "grpc_client_requests_total",
			Help: "Total number of gRPC client requests",
		},
		[]string{"service", "method", "status"},
	)

	// histogram for grpc requests duration
	GRPCClientDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "grpc_client_duration_seconds",
			Help:    "Duration of gRPC client requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method"},
	)
)

type Metrics struct {
	HTTPRequestTotal    *prometheus.CounterVec
	HTTPRequestDuration *prometheus.HistogramVec
}

func NewMetrics(reg prometheus.Registerer) *Metrics {
	m := &Metrics{
		HTTPRequestTotal:    HTTPRequestTotal,
		HTTPRequestDuration: HTTPRequestDuration,
	}
	reg.MustRegister(m.HTTPRequestTotal, m.HTTPRequestDuration)
	return m
}
