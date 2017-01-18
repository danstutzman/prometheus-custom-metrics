package storage

import (
	"fmt"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
)

type Options struct {
	Bigquery *bigquery.Options
}

func Usage() string {
	return fmt.Sprintf(`{
    "Bigquery": %s
  }`,
		bigquery.Usage(),
	)
}

func Setup(opts *Options) *bigquery.BigqueryConnection {
	var bigqueryConn *bigquery.BigqueryConnection

	if opts.Bigquery != nil {
		bigqueryConn = bigquery.Setup(opts.Bigquery)
	}

	return bigqueryConn
}
