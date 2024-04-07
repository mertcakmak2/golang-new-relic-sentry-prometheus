package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	HttpRequestCountWithPath = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total_with_path",
			Help: "Number of HTTP requests by path.",
		},
		[]string{"url"},
	)

	// PROMQL => rate(http_request_duration_seconds_sum{}[5m]) / rate(http_request_duration_seconds_count{}[5m])
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "http_request_duration_seconds",
			Help: "Average response time of HTTP requests.",
		},
		[]string{"path"},
	)
)

func init() {
	prometheus.MustRegister(HttpRequestCountWithPath)
	prometheus.MustRegister(HttpRequestDuration)
}
