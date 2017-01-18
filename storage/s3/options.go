package s3

import (
	"log"
)

type Options struct {
	S3CredsPath  string
	S3Region     string
	S3BucketName string
}

func Usage() string {
	return `
	    "S3CredsPath":     STRING,  path to AWS credentials file, e.g. "./s3.creds.ini"
      "S3Region":        STRING,  AWS region for S3, e.g. "us-east-1"
      "S3BucketName":    STRING,  Name of S3 bucket, e.g. "cloudfront-logs-danstutzman"`
}

func ValidateOptions(options *Options) {
	if options.S3CredsPath == "" {
		log.Fatalf("Missing cloudfront_logs.S3CredsPath")
	}
	if options.S3Region == "" {
		log.Fatalf("Missing cloudfront_logs.S3Region")
	}
	if options.S3BucketName == "" {
		log.Fatalf("Missing cloudfront_logs.S3BucketName")
	}
}
