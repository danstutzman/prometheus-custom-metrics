package billing_gcloud

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"log"
)

type Options struct {
	Bigquery    bigquery.Options
	MetricsPort int
}

func Usage() string {
	return `{ (optional)
	    "Bigquery": ` + bigquery.Usage() +
		`, MetricsPort":     INT      port to serve metrics on, e.g. 9102
    }`
}

func validateOptions(options *Options) {
	bigquery.ValidateOptions(&options.Bigquery)

	if options.MetricsPort == 0 {
		log.Fatalf("Missing billing_gcloud.MetricsPort")
	}
}
