package billing_gcloud

import (
	"log"
)

func MakeCollector(options *Options) *BillingGcloudCollector {
	validateOptions(options)
	bigquery := NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, options.GcloudDatasetName)
	log.Printf("product to sum cost: %v", bigquery.QueryProductToSumCost())
	return NewBillingGcloudCollector(options)
}
