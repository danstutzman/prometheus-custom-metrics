package billing_aws

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/s3"
	"log"
)

type Options struct {
	Bigquery    bigquery.Options
	S3          s3.Options
	MetricsPort int
}

func Usage() string {
	return `{ (optional)
    "Bigquery": ` + bigquery.Usage() + `, "S3": ` + s3.Usage() +
		`, MetricsPort":     INT      port to serve metrics on, e.g. 9102
  }`
}

func validateOptions(options *Options) {
	bigquery.ValidateOptions(&options.Bigquery)
	s3.ValidateOptions(&options.S3)

	if options.MetricsPort == 0 {
		log.Fatalf("Missing billing_aws.MetricsPort")
	}
}
