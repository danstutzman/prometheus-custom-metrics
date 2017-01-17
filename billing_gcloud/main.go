package billing_gcloud

import (
	"github.com/danielstutzman/prometheus-custom-metrics/util"
)

func MakeCollector(options *Options) *BillingGcloudCollector {
	validateOptions(options)
	bigquery := util.NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, options.GcloudDatasetName)
	return NewBillingGcloudCollector(options, bigquery)
}
