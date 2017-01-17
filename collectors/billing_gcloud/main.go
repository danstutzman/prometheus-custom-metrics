package billing_gcloud

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"log"
)

func MakeCollector(options *Options) *BillingGcloudCollector {
	validateOptions(options)
	bigquery := bigquery.NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, options.GcloudDatasetName)
	log.Printf("Created bigquery connection")
	return NewBillingGcloudCollector(options, bigquery)
}
