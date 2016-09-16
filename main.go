package main

import (
	"flag"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bigquery "google.golang.org/api/bigquery/v2"
	"io/ioutil"
	"log"
)

func main() {
	pemPath := flag.String("pem_path", "", "path to Google credentials in JSON format")
	googleCloudProjectId := flag.String("google_cloud_project_id", "",
		"Project number or project ID")
	flag.Parse()
	if *pemPath == "" {
		log.Fatal("Missing --pem_path")
	}
	if *googleCloudProjectId == "" {
		log.Fatal("Missing --google_cloud_project_id")
	}

	pemKeyBytes, err := ioutil.ReadFile(*pemPath)
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
	response, err := service.Jobs.Query(*googleCloudProjectId, &bigquery.QueryRequest{
		Query: sql,
	}).Do()
	if err != nil {
		panic(err)
	}

	for _, row := range response.Rows {
		log.Printf("Row: %v", row.F[0].V.(string))
	}
}
