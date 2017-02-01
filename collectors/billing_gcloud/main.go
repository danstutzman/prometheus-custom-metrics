package billing_gcloud

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"github.com/sirupsen/logrus"
)

func MakeCollector(options *Options, log *logrus.Logger) *BillingGcloudCollector {
	validateOptions(options)
	bigqueryConn := bigquery.NewBigqueryConnection(&options.Bigquery, log)
	return NewBillingGcloudCollector(options, bigqueryConn)
}
