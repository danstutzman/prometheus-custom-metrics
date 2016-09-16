package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bigquery "google.golang.org/api/bigquery/v2"
	"io/ioutil"
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

func testS3(credsPath, region, bucketName string) {
	log.Printf("Creating AWS session...")
	session, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Credentials: credentials.NewSharedCredentials(credsPath, ""),
			Region:      aws.String(region),
		},
	})
	if err != nil {
		panic(fmt.Errorf("Couldn't create AWS session: %s", err))
	}
	s3Service := s3.New(session)

	log.Printf("Listing objects in s3://%s...", bucketName)
	resp, err := s3Service.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		//		ContinuationToken: aws.String("Token"),
		//		Delimiter:         aws.String("Delimiter"),
		//		EncodingType:      aws.String("EncodingType"),
		//		FetchOwner:        aws.Bool(true),
		//		MaxKeys:           aws.Int64(1),
		//		Prefix:            aws.String("Prefix"),
		//		StartAfter:        aws.String("StartAfter"),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("Couldn't ListObjectsV2: %s", err))
	}
	fmt.Println(resp)
}

func testGcloud(pemPath, projectId string) {
	pemKeyBytes, err := ioutil.ReadFile(pemPath)
	if err != nil {
		panic(err)
	}

	log.Printf("Obtaining OAuth2 token...")
	token, err := google.JWTConfigFromJSON(pemKeyBytes, bigquery.BigqueryScope)
	client := token.Client(oauth2.NoContext)

	service, err := bigquery.New(client)
	if err != nil {
		panic(err)
	}

	log.Printf("Querying BigQuery...")
	sql := `SELECT COUNT(*)
	  FROM cloudfront_logs.cloudfront_logs
		LIMIT 1000`
	response, err := service.Jobs.Query(projectId,
		&bigquery.QueryRequest{
			Query: sql,
		}).Do()
	if err != nil {
		panic(err)
	}

	for _, row := range response.Rows {
		log.Printf("Row: %v", row.F[0].V.(string))
	}
}

func main() {
	args := parseArgsOrFatal()

	testS3(args.s3CredsPath, args.s3Region, args.s3BucketName)
	testGcloud(args.gcloudPemPath, args.gcloudProjectId)
}
