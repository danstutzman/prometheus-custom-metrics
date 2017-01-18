package cloudfront_logs

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/s3"
	"log"
)

type Options struct {
	S3          s3.Options
	Bigquery    bigquery.Options
	MetricsPort int
}

func Usage() string {
	return `{ (optional)
    "S3": ` + s3.Usage() + `, "BigQuery": ` + bigquery.Usage() +
		`, "MetricsPort":     INT,     port to serve metrics on, e.g. 9102
  }`
}

func validateOptions(options *Options) {
	s3.ValidateOptions(&options.S3)
	bigquery.ValidateOptions(&options.Bigquery)

	if options.MetricsPort == 0 {
		log.Fatalf("Missing cloudfront_logs.MetricsPort")
	}
}
