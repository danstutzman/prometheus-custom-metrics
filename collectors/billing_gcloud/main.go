package billing_gcloud

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
)

func MakeCollector(options *Options) *BillingGcloudCollector {
	validateOptions(options)
	bigqueryConn := bigquery.NewBigqueryConnection(&options.Bigquery)
	return NewBillingGcloudCollector(options, bigqueryConn)
}
