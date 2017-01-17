package billing_gcloud

import (
	"fmt"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"strings"
)

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

func QueryProductToSumCost(conn *bigquery.BigqueryConnection) map[string]float64 {
	sql := `SELECT 
		  product,
  		resource_type,
		SUM(cost) AS sum_cost
		FROM ` + fmt.Sprintf("`%s.gcp_billing_export_*`", conn.DatasetId()) + `
		WHERE currency = 'USD'
		GROUP BY product, resource_type
		HAVING sum_cost >= 0.01`

	rows := conn.Query(sql, "product to sum cost")

	productToSumCost := map[string]float64{}
	for _, row := range rows {
		googleProduct := row.F[0].V.(string)
		resourceType := row.F[1].V.(string)
		sumCost := bigquery.ParseFloat64(row.F[2].V.(string))
		product := RollUpProduct(googleProduct, resourceType)
		productToSumCost[product] = sumCost
	}
	return productToSumCost
}
