package cloudfront_logs

import (
	"github.com/prometheus/client_golang/prometheus"
	"log"
)

func Main(options *Options) {
	if options.S3CredsPath == "" {
		log.Fatalf("Missing options.S3CredsPath")
	}
	if options.S3Region == "" {
		log.Fatalf("Missing options.S3Region")
	}
	if options.S3BucketName == "" {
		log.Fatalf("Missing options.S3BucketName")
	}
	if options.GcloudPemPath == "" {
		log.Fatalf("Missing options.GcloudPemPath")
	}
	if options.GcloudProjectId == "" {
		log.Fatalf("Missing options.GcloudProjectId")
	}

	s3 := NewS3Connection(options.S3CredsPath, options.S3Region, options.S3BucketName)
	bigquery := NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, "cloudfront_logs")
	collector := NewCloudfrontCollector(s3, bigquery)
	collector.InitFromBigqueryAndS3()
	prometheus.MustRegister(collector)
}
