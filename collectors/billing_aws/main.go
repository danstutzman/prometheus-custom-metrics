package billing_aws

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"github.com/sirupsen/logrus"

	"github.com/danielstutzman/prometheus-custom-metrics/storage/s3"
)

func MakeCollector(options *Options, log *logrus.Logger) *BillingAwsCollector {
	validateOptions(options)

	bigqueryConn := bigquery.NewBigqueryConnection(&options.Bigquery, log)
	s3Conn := s3.NewS3Connection(&options.S3)

	return NewBillingAwsCollector(options, bigqueryConn, s3Conn)
}
