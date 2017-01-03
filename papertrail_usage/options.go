package papertrail_usage

import (
	"log"
)

type Options struct {
	ApiToken    string
	MetricsPort int
}

func Usage() string {
	return `{ (optional)
	   "ApiToken":     STRING   token from https://papertrailapp.com/user/edit
     "MetricsPort":  INT      port to serve metrics on, e.g. 9102
	}`
}

func validateOptions(options *Options) {
	if options.ApiToken == "" {
		log.Fatalf("Missing memory_usage.ApiToken")
	}
	if options.MetricsPort == 0 {
		log.Fatalf("Missing memory_usage.MetricsPort")
	}
}
