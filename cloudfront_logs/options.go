package cloudfront_logs

import (
	"fmt"
)

type Options struct {
	S3CredsPath     string
	S3Region        string
	S3BucketName    string
	GcloudPemPath   string
	GcloudProjectId string
}

func Usage() string {
	return `{ (optional)
		"S3CredsPath":     STRING,  path to AWS credentials file, e.g. "./s3.creds.ini"
		"S3Region":        STRING,  AWS region for S3, e.g. "us-east-1"
		"S3BucketName":    STRING,  Name of S3 bucket, e.g. "cloudfront-logs-danstutzman"
		"GcloudPemPath":   STRING,  path to Google credentials in JSON format, e.g. "./Speech-ba6281533dc8.json"
		"GcloudProjectId": STRING   Project number or project ID
	}`
}

func validateOptions(options *Options) error {
	}
	if options.S3CredsPath == "" {
		return fmt.Errorf("Missing options.S3CredsPath")
	}
	if options.S3Region == "" {
		return fmt.Errorf("Missing options.S3Region")
	}
	if options.S3BucketName == "" {
		return fmt.Errorf("Missing options.S3BucketName")
	}
	if options.GcloudPemPath == "" {
		return fmt.Errorf("Missing options.GcloudPemPath")
	}
	if options.GcloudProjectId == "" {
		return fmt.Errorf("Missing options.GcloudProjectId")
	}
	return nil
}
