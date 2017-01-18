package billing_gcloud

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"log"
)

func MakeCollector(options *Options,
	bigqueryConn *bigquery.BigqueryConnection) *BillingGcloudCollector {

	validateOptions(options)
	if bigqueryConn == nil {
		log.Fatalf("Missing Bigquery configuration")
	}

	return NewBillingGcloudCollector(options, bigqueryConn)
}
