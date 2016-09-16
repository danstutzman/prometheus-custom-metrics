package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bigquery "google.golang.org/api/bigquery/v2"
	"io/ioutil"
	"log"
)

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
