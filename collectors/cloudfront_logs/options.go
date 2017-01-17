package cloudfront_logs

import (
	"log"
)

type Options struct {
	MetricsPort     int
	S3CredsPath     string
	S3Region        string
	S3BucketName    string
	GcloudPemPath   string
	GcloudProjectId string
}

func Usage() string {
	return `{ (optional)
      "MetricsPort":     INT,     port to serve metrics on, e.g. 9102
      "S3CredsPath":     STRING,  path to AWS credentials file, e.g. "./s3.creds.ini"
      "S3Region":        STRING,  AWS region for S3, e.g. "us-east-1"
      "S3BucketName":    STRING,  Name of S3 bucket, e.g. "cloudfront-logs-danstutzman"
      "GcloudPemPath":   STRING,  path to Google credentials in JSON format, e.g. "./Speech-ba6281533dc8.json"
      "GcloudProjectId": STRING   Project number or project ID
    }`
}

func validateOptions(options *Options) {
	if options.MetricsPort == 0 {
		log.Fatalf("Missing cloudfront_logs.MetricsPort")
	}
	if options.S3CredsPath == "" {
		log.Fatalf("Missing cloudfront_logs.S3CredsPath")
	}
	if options.S3Region == "" {
		log.Fatalf("Missing cloudfront_logs.S3Region")
	}
	if options.S3BucketName == "" {
		log.Fatalf("Missing cloudfront_logs.S3BucketName")
	}
	if options.GcloudPemPath == "" {
		log.Fatalf("Missing cloudfront_logs.GcloudPemPath")
	}
	if options.GcloudProjectId == "" {
		log.Fatalf("Missing cloudfront_logs.GcloudProjectId")
	}
}
