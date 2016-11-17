package cloudfront_logs

func MakeCollector(options *Options) *CloudfrontCollector {
	validateOptions(options)
	s3 := NewS3Connection(options.S3CredsPath, options.S3Region, options.S3BucketName)
	bigquery := NewBigqueryConnection(options.GcloudPemPath,
		options.GcloudProjectId, "cloudfront_logs")
	collector := NewCloudfrontCollector(s3, bigquery)
	return collector
}
