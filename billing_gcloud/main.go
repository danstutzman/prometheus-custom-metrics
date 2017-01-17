package billing_gcloud

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage"
)

func MakeCollector(options *Options) *BillingGcloudCollector {
	validateOptions(options)
	bigquery := storage.NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, options.GcloudDatasetName)
	return NewBillingGcloudCollector(options, bigquery)
}
