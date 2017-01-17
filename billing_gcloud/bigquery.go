package billing_gcloud

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
	"strings"
)

type BigqueryConnection struct {
	projectId string
	datasetId string
	service   *bigquery.Service
}

func parseFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
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

func RollUpProduct(googleProduct string, resourceType string) string {
	if googleProduct == "BigQuery" {
		return "bigquery"
	} else if googleProduct == "Cloud Storage" {
		return "cloudstorage"
	} else if googleProduct == "Compute Engine" {
		if strings.Contains(resourceType, "Storage") {
			return "computeengine_storage"
		} else if strings.Contains(resourceType, "CPU running in") {
			return "computeengine_instance"
		} else {
			return "computeengine_other"
		}
	} else {
		return "other"
	}
}

func (conn *BigqueryConnection) QueryProductToSumCost() map[string]float64 {
	sql := `SELECT 
		  product,
  		resource_type,
		SUM(cost) AS sum_cost
		FROM ` + fmt.Sprintf("`%s.gcp_billing_export_*`", conn.datasetId) + `
		WHERE currency = 'USD'
		GROUP BY product, resource_type
		HAVING sum_cost >= 0.01`

	var response *bigquery.QueryResponse
	var err error
	err = backoff.Retry(func() error {
		log.Printf("Querying product to sum cost (USD)...")
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

	productToSumCost := map[string]float64{}
	for _, row := range response.Rows {
		googleProduct := row.F[0].V.(string)
		resourceType := row.F[1].V.(string)
		sumCost := parseFloat64(row.F[2].V.(string))
		product := RollUpProduct(googleProduct, resourceType)
		productToSumCost[product] = sumCost
	}
	return productToSumCost
}
