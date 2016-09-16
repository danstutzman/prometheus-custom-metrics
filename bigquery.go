package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bigquery "google.golang.org/api/bigquery/v2"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
)

var S3_PATH_REGEXP = regexp.MustCompile(
	`^([A-Z0-9]+)\.([0-9]{4})-([0-9]{2})-([0-9]{2})-([0-9]{2}).([0-9a-f]{8}).gz$`)

type BigqueryConnection struct {
	projectId string
	datasetId string
	service   *bigquery.Service
}

func NewBigqueryConnection(pemPath, projectId, datasetId string) *BigqueryConnection {
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

	return &BigqueryConnection{
		projectId: projectId,
		datasetId: datasetId,
		service:   service,
	}
}

func (conn *BigqueryConnection) createSitesTable() {
	log.Printf("Creating sites table first...")
	response, err := conn.service.Tables.Insert(conn.projectId, conn.datasetId,
		&bigquery.Table{
			Schema: &bigquery.TableSchema{
				Fields: []*bigquery.TableFieldSchema{
					{Name: "site_domain_name", Type: "STRING", Mode: "REQUIRED"},
					{Name: "s3_path", Type: "STRING", Mode: "REQUIRED"},
				},
			},
			TableReference: &bigquery.TableReference{
				DatasetId: conn.datasetId,
				ProjectId: conn.projectId,
				TableId:   "sites",
			},
		}).Do()
	if err != nil {
		panic(err)
	}
	_ = response
}

func (conn *BigqueryConnection) QueryLastS3Paths() []string {
	log.Printf("Querying last S3 paths, query 1/2...")
	sql1 := fmt.Sprintf(`SELECT site_domain_name, max(s3_path) AS last_s3_path
		FROM %s.sites
		GROUP BY site_domain_name`, conn.datasetId)
	response1, err := conn.service.Jobs.Query(conn.projectId,
		&bigquery.QueryRequest{Query: sql1}).Do()
	if err != nil {
		if err.Error() == "googleapi: Error 404: Not found: Table speech-danstutzman:cloudfront_logs.sites, notFound" {
			conn.createSitesTable()
			response1, err = conn.service.Jobs.Query(conn.projectId,
				&bigquery.QueryRequest{Query: sql1}).Do()
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}

	lastS3Paths := []string{}
	log.Printf("Querying last S3 paths, query 2/2...")
	if response1.TotalRows > 0 {
		likes := []string{}
		for _, row := range response1.Rows {
			lastS3Path := row.F[1].V.(string)
			groups := S3_PATH_REGEXP.FindStringSubmatch(lastS3Path)
			if groups == nil {
				panic(fmt.Errorf("s3_path value of '%s' didn't match regexp '%s'",
					lastS3Path, S3_PATH_REGEXP))
			}
			like := fmt.Sprintf("'%s.%s-%s-%s-%s.%%'",
				groups[1], groups[2], groups[3], groups[4])
			likes = append(likes, like)
		}
		sql2 := fmt.Sprintf(`SELECT s3_path
			FROM %s.sites
			WHERE last_s3_path LIKE %s`, conn.datasetId, strings.Join(likes, " OR "))
		response2, err := conn.service.Jobs.Query(conn.projectId,
			&bigquery.QueryRequest{Query: sql2}).Do()
		if err != nil {
			panic(err)
		}
		for _, row := range response2.Rows {
			s3Path := row.F[1].V.(string)
			lastS3Paths = append(lastS3Paths, s3Path)
		}
	}
	return lastS3Paths
}

func (conn *BigqueryConnection) TestQuery() {
	log.Printf("Querying BigQuery...")
	sql := `SELECT COUNT(*)
	  FROM cloudfront_logs.cloudfront_logs
		LIMIT 1000`
	response, err := conn.service.Jobs.Query(conn.projectId,
		&bigquery.QueryRequest{Query: sql}).Do()
	if err != nil {
		panic(err)
	}

	for _, row := range response.Rows {
		log.Printf("Row: %v", row.F[0].V.(string))
	}
}
