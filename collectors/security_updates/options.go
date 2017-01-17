package security_updates

import (
	"log"
)

type Options struct {
	MetricsPort int
}

func Usage() string {
	return `{ (optional)
      "MetricsPort":       INT     port to run web server on, e.g. 9102
    }`
}

func validateOptions(options *Options) {
	if options.MetricsPort == 0 {
		log.Fatalf("Missing security_updates.MetricsPort")
	}
}
