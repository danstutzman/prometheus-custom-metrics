package security_updates

import (
	"github.com/prometheus/client_golang/prometheus"
)

func Usage() string {
	return `BOOL (e.g. true)`
}

func Main() {
	collector := NewSecurityUpdatesCollector()
	prometheus.MustRegister(collector)
}
