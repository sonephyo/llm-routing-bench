package metrics

import "github.com/prometheus/client_golang/prometheus"

var RequestCount = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "lb_requests_total",
		Help: "Total requests per backend",
	},
	[]string{"backend"},
)

var RequestLatency = prometheus.NewHistogramVec(
	prometheus.HistogramOpts{
		Name: "lb_request_duration_seconds",
		Help: "Request latency per backend",
	},
	[]string{"backend"},
)
