package billing_aws

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/s3"
)

func MakeCollector(options *Options) *BillingAwsCollector {
	validateOptions(options)

	bigqueryConn := bigquery.NewBigqueryConnection(&options.Bigquery)
	s3Conn := s3.NewS3Connection(&options.S3)

	return NewBillingAwsCollector(options, bigqueryConn, s3Conn)
}
