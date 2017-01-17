package cloudfront_logs

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
)

func MakeCollector(options *Options) *CloudfrontCollector {
	validateOptions(options)
	s3 := NewS3Connection(options.S3CredsPath, options.S3Region, options.S3BucketName)
	bigquery := bigquery.NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, "cloudfront_logs")
	collector := NewCloudfrontCollector(s3, bigquery)
	return collector
}
