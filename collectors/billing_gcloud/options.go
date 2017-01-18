package billing_gcloud

import (
	"log"
)

type Options struct {
	MetricsPort     int
	BigqueryDataset string
}

func Usage() string {
	return `{ (optional)
      "MetricsPort":     INT      port to serve metrics on, e.g. 9102
      "BigqueryDataset": STRING   Name of dataset
    }`
}

func validateOptions(options *Options) {
	if options.MetricsPort == 0 {
		log.Fatalf("Missing billing_gcloud.MetricsPort")
	}
	if options.BigqueryDataset == "" {
		log.Fatalf("Missing billing_gcloud.BigqueryDataset")
	}
}
