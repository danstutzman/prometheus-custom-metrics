package main

import (
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
)

type Args struct {
	s3CredsPath     string
	s3Region        string
	s3BucketName    string
	gcloudPemPath   string
	gcloudProjectId string
	portNum         int
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
	portNum := flag.Int("port_num", 0, "Port number to run web server on")
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
	if *portNum == 0 {
		log.Fatal("Missing --port_num")
	}

	return Args{
		s3CredsPath:     *s3CredsPath,
		s3Region:        *s3Region,
		s3BucketName:    *s3BucketName,
		gcloudPemPath:   *gcloudPemPath,
		gcloudProjectId: *gcloudProjectId,
		portNum:         *portNum,
	}
}

func main() {
	args := parseArgsOrFatal()
	s3 := NewS3Connection(args.s3CredsPath, args.s3Region, args.s3BucketName)
	bigquery := NewBigqueryConnection(args.gcloudPemPath,
		args.gcloudProjectId, "cloudfront_logs")
	collector := NewCloudfrontCollector(s3, bigquery)
	collector.InitFromBigqueryAndS3()

	prometheus.MustRegister(collector)
	http.Handle("/metrics", prometheus.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", args.portNum), nil)
	if err != nil {
		panic(err)
	}
}
