package cloudfront_logs

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"log"
)

func MakeCollector(options *Options,
	bigqueryConn *bigquery.BigqueryConnection) *CloudfrontCollector {

	validateOptions(options)
	if bigqueryConn == nil {
		log.Fatalf("Missing Bigquery configuration")
	}

	s3 := NewS3Connection(options.S3CredsPath, options.S3Region, options.S3BucketName)
	collector := NewCloudfrontCollector(options, s3, bigqueryConn)
	return collector
}
