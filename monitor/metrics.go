package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

var TotalFailures = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "telnet_failures_total",
		Help: "Number of telnet tests failed.",
	},
	[]string{"status"},
)

var TotalPasses = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "telnet_passes_total",
		Help: "Number of telnet tests passed.",
	},
	[]string{"status"},
)

func Init() {
	prometheus.Register(TotalPasses)
	prometheus.Register(TotalFailures)
}
