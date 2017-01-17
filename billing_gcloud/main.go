package billing_gcloud

import ()

func MakeCollector(options *Options) *BillingGcloudCollector {
	validateOptions(options)
	bigquery := NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, options.GcloudDatasetName)
	return NewBillingGcloudCollector(options, bigquery)
}
