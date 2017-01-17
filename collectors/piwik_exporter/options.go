package piwik_exporter

import (
	"log"
)

type Options struct {
	MetricsPort int
}

func Usage() string {
	return `{ (optional)
      "MetricsPort":     INT     port to serve metrics on, e.g. 9102
    }`
}

func validateOptions(options *Options) {
	if options.MetricsPort == 0 {
		log.Fatalf("Missing piwik_exporter.MetricsPort")
	}
}
