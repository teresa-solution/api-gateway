package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

var (
	HttpRequestsTotal = prometheus.NewCounterVec( // Changed to HttpRequestsTotal
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)
)

func InitMetrics() {
	err := prometheus.Register(HttpRequestsTotal)
	if err != nil {
		log.Error().Err(err).Msg("Failed to register HttpRequestsTotal metric")
	}
}
