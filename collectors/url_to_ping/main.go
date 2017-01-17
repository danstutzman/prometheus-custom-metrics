package url_to_ping

import (
	"github.com/prometheus/client_golang/prometheus"
)

func MakeCollector(options *Options) prometheus.Collector {
	validateOptions(options)
	return NewUrlToPingCollector(options)
}
