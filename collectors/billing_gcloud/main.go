package billing_gcloud

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage"
	"log"
)

func MakeCollector(options *Options) *BillingGcloudCollector {
	validateOptions(options)
	bigquery := storage.NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, options.GcloudDatasetName)
	log.Printf("Created bigquery connection")
	return NewBillingGcloudCollector(options, bigquery)
}
