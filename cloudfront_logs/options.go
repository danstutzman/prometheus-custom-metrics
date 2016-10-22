package cloudfront_logs

import (
	"github.com/danielstutzman/prometheus-cloudfront-logs-exporter/json_value"
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

func HandleOptions(section map[string]interface{}, path string,
	usagef func(string, ...interface{})) *Options {

	options := Options{}
	for key, value := range section {
		switch key {
		case "S3CredsPath":
			options.S3CredsPath = json_value.ToString(value, path+".S3CredsPath", usagef)
		case "S3Region":
			options.S3Region = json_value.ToString(value, path+".S3Region", usagef)
		case "S3BucketName":
			options.S3BucketName = json_value.ToString(value, path+".S3BucketName", usagef)
		case "GcloudPemPath":
			options.GcloudPemPath = json_value.ToString(value, path+".GcloudPemPath", usagef)
		case "GcloudProjectId":
			options.GcloudProjectId = json_value.ToString(value, path+".GcloudProjectId",
				usagef)
		default:
			usagef("Unknown key %s.%s", path, key)
		}
	}

	if options.S3CredsPath == "" {
		usagef("Missing %s.S3CredsPath", path)
	} else if options.S3Region == "" {
		usagef("Missing %s.S3Region", path)
	} else if options.S3BucketName == "" {
		usagef("Missing %s.S3BucketName", path)
	} else if options.GcloudPemPath == "" {
		usagef("Missing %s.GcloudPemPath", path)
	} else if options.GcloudProjectId == "" {
		usagef("Missing %s.GcloudProjectId", path)
	}

	return &options
}
