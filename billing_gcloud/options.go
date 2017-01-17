package billing_gcloud

import (
	"log"
)

type Options struct {
	MetricsPort       int
	GcloudPemPath     string
	GcloudProjectId   string
	GcloudDatasetName string
}

func Usage() string {
	return `{ (optional)
     "MetricsPort":       INT      port to serve metrics on, e.g. 9102
     "GcloudPemPath":     STRING,  path to Google credentials in JSON format,
		                                 e.g. "./Speech-ba6281533dc8.json"
     "GcloudProjectId":   STRING   Project number or project ID
     "GcloudDatasetName": STRING   Name of dataset
	}`
}

func validateOptions(options *Options) {
	if options.MetricsPort == 0 {
		log.Fatalf("Missing memory_usage.MetricsPort")
	}
	if options.GcloudPemPath == "" {
		log.Fatalf("Missing cloudfront_logs.GcloudPemPath")
	}
	if options.GcloudProjectId == "" {
		log.Fatalf("Missing cloudfront_logs.GcloudProjectId")
	}
	if options.GcloudDatasetName == "" {
		log.Fatalf("Missing cloudfront_logs.GcloudDatasetName")
	}
}
