package memory_usage

import (
	"github.com/prometheus/client_golang/prometheus"
)

func Usage() string {
	return `BOOL (e.g. true)`
}

func Main() {
	collector := NewMemoryUsageCollector()
	prometheus.MustRegister(collector)
}
