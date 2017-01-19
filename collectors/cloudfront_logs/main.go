package cloudfront_logs

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/s3"
)

func MakeCollector(options *Options) *CloudfrontCollector {
	validateOptions(options)
	bigqueryConn := bigquery.NewBigqueryConnection(&options.Bigquery)
	s3 := s3.NewS3Connection(&options.S3)
	collector := NewCloudfrontCollector(options, s3, bigqueryConn)
	return collector
}