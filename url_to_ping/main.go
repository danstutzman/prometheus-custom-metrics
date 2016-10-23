package url_to_ping

import (
	"github.com/prometheus/client_golang/prometheus"
)

func Usage() string {
	return `STRING (URL, e.g. "https://nosnch.in/abcdef")`
}

func Main(url string) {
	collector := NewUrlToPingCollector(url)
	prometheus.MustRegister(collector)
}
