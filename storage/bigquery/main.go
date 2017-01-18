package bigquery

import (
	"fmt"
	"github.com/cenkalti/backoff"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	bigquery "google.golang.org/api/bigquery/v2"
	"google.golang.org/api/googleapi"
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

type BigqueryConnection struct {
	projectId string
	service   *bigquery.Service
}

func Atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func ParseFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func NewBigqueryConnection(pemPath, projectId string) *BigqueryConnection {
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
		service:   service,
	}
}

func (conn *BigqueryConnection) Query(sql, description string) []*bigquery.TableRow {
	var response *bigquery.QueryResponse
	var err error
	err = backoff.Retry(func() error {
		log.Printf("Querying %s...", description)
		response, err = conn.service.Jobs.Query(conn.projectId, &bigquery.QueryRequest{
			Query:        sql,
			UseLegacySql: googleapi.Bool(false),
		}).Do()
		if err != nil {
			err2, isGoogleApiError := err.(*googleapi.Error)
			if isGoogleApiError && (err2.Code == 500 || err2.Code == 503) {
				// then let backoff retry the query
			} else {
				log.Fatalf("Error %s; query was %s", err, sql)
			}
		}
		return err
	}, backoff.NewExponentialBackOff())
	if err != nil {
		panic(err)
	}

	return response.Rows
}

func (conn *BigqueryConnection) CreateTable(dataset string, tableName string,
	fields []*bigquery.TableFieldSchema) {

	log.Printf("Creating %s table...", tableName)
	_, err := conn.service.Tables.Insert(conn.projectId, dataset,
		&bigquery.Table{
			Schema: &bigquery.TableSchema{Fields: fields},
			TableReference: &bigquery.TableReference{
				DatasetId: dataset,
				ProjectId: conn.projectId,
				TableId:   tableName,
			},
		}).Do()
	if err != nil {
		panic(err)
	}

	log.Printf("Waiting 30 seconds for BigQuery to catch up...")
	time.Sleep(30 * time.Second)
}

func (conn *BigqueryConnection) InsertRows(dataset string, tableName string,
	createTable func(), rows []*bigquery.TableDataInsertAllRequestRows) {

	var err error
	err = backoff.Retry(func() error {
		log.Printf("Inserting rows to %s...", tableName)
		_, err := conn.service.Tabledata.InsertAll(conn.projectId, dataset,
			tableName, &bigquery.TableDataInsertAllRequest{Rows: rows}).Do()
		if err != nil {
			err2, isGoogleApiError := err.(*googleapi.Error)
			if isGoogleApiError && (err2.Code == 500 || err2.Code == 503) {
				// then let backoff retry the query
			} else {
				log.Fatalf("Error %s inserting rows", err)
			}
		}
		return err
	}, backoff.NewExponentialBackOff())
	if err != nil {
		panic(err)
	}

	if err != nil {
		log.Println(err)
		if err.Error() == fmt.Sprintf(
			"googleapi: Error 404: Not found: Table %s:%s.%s, notFound",
			conn.projectId, dataset, tableName) {

			createTable()

			// Now retry the insert
			err = backoff.Retry(func() error {
				_, err := conn.service.Tabledata.InsertAll(conn.projectId, dataset,
					tableName, &bigquery.TableDataInsertAllRequest{Rows: rows}).Do()
				if err != nil {
					err2, isGoogleApiError := err.(*googleapi.Error)
					if isGoogleApiError && (err2.Code == 500 || err2.Code == 503) {
						// then let backoff retry the query
					} else {
						log.Fatalf("Error %s inserting rows", err)
					}
				}
				return err
			}, backoff.NewExponentialBackOff())
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
}
