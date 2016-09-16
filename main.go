package main

import (
	"flag"
	"log"
)

type Args struct {
	s3CredsPath     string
	s3Region        string
	s3BucketName    string
	gcloudPemPath   string
	gcloudProjectId string
}

func parseArgsOrFatal() Args {
	s3CredsPath := flag.String("s3_creds_path", "",
		"path to AWS credentials file, e.g. ./s3.creds.ini")
	s3Region := flag.String("s3_region", "",
		"AWS region for S3, e.g. us-east-1")
	s3BucketName := flag.String("s3_bucket_name", "",
		"Name of S3 bucket, e.g. cloudfront-logs-danstutzman")
	gcloudPemPath := flag.String("gcloud_pem_path", "",
		"path to Google credentials in JSON format, e.g. ./Speech-ba6281533dc8.json")
	gcloudProjectId := flag.String("gcloud_project_id", "",
		"Project number or project ID")
	flag.Parse()

	if *s3CredsPath == "" {
		log.Fatal("Missing --s3_creds_path")
	}
	if *s3Region == "" {
		log.Fatal("Missing --s3_region")
	}
	if *s3BucketName == "" {
		log.Fatal("Missing --s3_bucket_name")
	}
	if *gcloudPemPath == "" {
		log.Fatal("Missing --gcloud_pem_path")
	}
	if *gcloudProjectId == "" {
		log.Fatal("Missing --gcloud_project_id")
	}

	return Args{
		s3CredsPath:     *s3CredsPath,
		s3Region:        *s3Region,
		s3BucketName:    *s3BucketName,
		gcloudPemPath:   *gcloudPemPath,
		gcloudProjectId: *gcloudProjectId,
	}
}

func main() {
	args := parseArgsOrFatal()

	testS3(args.s3CredsPath, args.s3Region, args.s3BucketName)
	testGcloud(args.gcloudPemPath, args.gcloudProjectId)
}
