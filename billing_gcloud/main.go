package billing_gcloud

func MakeCollector(options *Options) *BillingGcloudCollector {
	validateOptions(options)
	return NewBillingGcloudCollector(options)
}
