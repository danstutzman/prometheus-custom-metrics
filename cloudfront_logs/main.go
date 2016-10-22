package cloudfront_logs

import (
	"github.com/prometheus/client_golang/prometheus"
)

func Main(options Options) {
	s3 := NewS3Connection(options.S3CredsPath, options.S3Region, options.S3BucketName)
	bigquery := NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, "cloudfront_logs")
	collector := NewCloudfrontCollector(s3, bigquery)
	collector.InitFromBigqueryAndS3()
	prometheus.MustRegister(collector)
}
