package bigquery

import (
	"log"
)

type Options struct {
	GcloudPemPath   string
	GcloudProjectId string
}

func Usage() string {
	return `{ (optional)
      "GcloudPemPath":     STRING,  path to Google JSON creds for BigQuery,
                                      e.g. "./Speech-ba6281533dc8.json"
      "GcloudProjectId":   STRING   Project number or project ID for BigQuery
    }`
}

func validateOptions(options *Options) {
	if options.GcloudPemPath == "" {
		log.Fatalf("Missing cloudfront_logs.GcloudPemPath")
	}
	if options.GcloudProjectId == "" {
		log.Fatalf("Missing cloudfront_logs.GcloudProjectId")
	}
}

func Setup(opts *Options) *BigqueryConnection {
	return NewBigqueryConnection(opts.GcloudPemPath, opts.GcloudProjectId)
}
