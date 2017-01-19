package cloudfront_logs

import (
	"fmt"
	mybigquery "github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	bigquery "google.golang.org/api/bigquery/v2"
	"regexp"
)

var S3_PATH_REGEXP = regexp.MustCompile(
	`^([A-Z0-9]+)\.([0-9]{4})-([0-9]{2})-([0-9]{2})-([0-9]{2}).([0-9a-f]{8}).gz$`)

type SiteNameStatus struct {
	SiteName string
	Status   string
}

func maybeNull(s string) bigquery.JsonValue {
	if s == "-" {
		return nil
	} else if s == "" {
		return nil
	} else {
		return s
	}
}

func (collector *CloudfrontCollector) createVisitsTable() {
	collector.bigquery.CreateTable("visits", []*bigquery.TableFieldSchema{
		{Name: "s3_path", Type: "STRING", Mode: "REQUIRED"},
		{Name: "datetime", Type: "DATETIME", Mode: "REQUIRED"},
		{Name: "x_edge_location", Type: "STRING", Mode: "REQUIRED"},
		{Name: "c_ip", Type: "STRING", Mode: "REQUIRED"},
		{Name: "cs_method", Type: "STRING", Mode: "REQUIRED"},
		{Name: "sc_status", Type: "INTEGER", Mode: "REQUIRED"},
		{Name: "cs_referer", Type: "STRING", Mode: "NULLABLE"},
		{Name: "x_host_header", Type: "STRING", Mode: "REQUIRED"},
		{Name: "time_taken", Type: "FLOAT", Mode: "REQUIRED"},
		{Name: "x_forwarded_for", Type: "STRING", Mode: "NULLABLE"},
		{Name: "cs_protocol", Type: "STRING", Mode: "REQUIRED"},
		{Name: "cs_uri_stem", Type: "STRING", Mode: "REQUIRED"},
		{Name: "cs_user_agent", Type: "STRING", Mode: "REQUIRED"},
		{Name: "cs_uri_query", Type: "STRING", Mode: "NULLABLE"},
		{Name: "x_edge_response_result_type", Type: "STRING", Mode: "NULLABLE"},
		{Name: "sc_bytes", Type: "INTEGER", Mode: "REQUIRED"},
		{Name: "cs_host", Type: "STRING", Mode: "REQUIRED"},
		{Name: "cs_cookie", Type: "STRING", Mode: "NULLABLE"},
		{Name: "x_edge_result_type", Type: "STRING", Mode: "REQUIRED"},
		{Name: "cs_bytes", Type: "INTEGER", Mode: "REQUIRED"},
		{Name: "ssl_protocol", Type: "STRING", Mode: "NULLABLE"},
		{Name: "ssl_cipher", Type: "STRING", Mode: "NULLABLE"},
		{Name: "cs_protocol_version", Type: "STRING", Mode: "NULLABLE"},
	})
}

func (collector *CloudfrontCollector) uploadVisits(s3Path string,
	visits []map[string]string) {

	rows := make([]*bigquery.TableDataInsertAllRequestRows, 0)
	for _, visit := range visits {
		row := &bigquery.TableDataInsertAllRequestRows{
			InsertId: s3Path + "-1-" + visit["x-edge-request-id"],
			Json: map[string]bigquery.JsonValue{
				// don't include x-edge-request-id
				"s3_path":                     s3Path,
				"datetime":                    visit["date"] + "T" + visit["time"],
				"x_edge_location":             visit["x-edge-location"],
				"c_ip":                        visit["c-ip"],
				"cs_method":                   visit["cs-method"],
				"sc_status":                   visit["sc-status"],
				"cs_referer":                  maybeNull(visit["cs(Referer)"]),
				"x_host_header":               visit["x-host-header"],
				"time_taken":                  visit["time-taken"],
				"x_forwarded_for":             maybeNull(visit["x-forwarded-for"]),
				"cs_protocol":                 visit["cs-protocol"],
				"cs_uri_stem":                 visit["cs-uri-stem"],
				"cs_user_agent":               visit["cs(User-Agent)"],
				"cs_uri_query":                maybeNull(visit["cs-uri-query"]),
				"x_edge_response_result_type": maybeNull(visit["x-edge-response-result-type"]),
				"sc_bytes":                    visit["sc-bytes"],
				"cs_host":                     visit["cs(Host)"],
				"cs_cookie":                   maybeNull(visit["cs(Cookie)"]),
				"x_edge_result_type":          visit["x-edge-result-type"],
				"cs_bytes":                    visit["cs-bytes"],
				"ssl_protocol":                maybeNull(visit["ssl-protocol"]),
				"ssl_cipher":                  maybeNull(visit["ssl-cipher"]),
				"cs_protocol_version":         maybeNull(visit["cs-protocol-version"]),
			}}
		rows = append(rows, row)
	}

	collector.bigquery.InsertRows("visits", func() { collector.createVisitsTable() },
		rows)
}

func rollUpExactStatus(exactStatus int) string {
	switch {
	case exactStatus == 200:
		return "200"
	case exactStatus == 404:
		return "200"
	case exactStatus >= 500 && exactStatus <= 599:
		return "5xx"
	default:
		return "other"
	}
}

func (collector *CloudfrontCollector) querySiteNameStatusToNumVisits() map[SiteNameStatus]int {

	sql := fmt.Sprintf(`SELECT x_host_header AS site_name,
			sc_status AS exact_status,
			COUNT(*) AS num_visits
		FROM %s.visits
		GROUP BY site_name, exact_status`, collector.bigquery.DatasetName())

	rows := collector.bigquery.Query(sql, "site name * status to num visits")

	siteNameStatusToNumVisits := map[SiteNameStatus]int{}
	for _, row := range rows {
		siteName := row.F[0].V.(string)
		exactStatus := mybigquery.Atoi(row.F[1].V.(string))
		numVisits := mybigquery.Atoi(row.F[2].V.(string))
		status := rollUpExactStatus(exactStatus)
		siteNameStatusToNumVisits[SiteNameStatus{siteName, status}] += numVisits
	}
	return siteNameStatusToNumVisits
}

// Returns _sum and _count
func (collector *CloudfrontCollector) querySiteNameToRequestSeconds() (
	map[string]float64, map[string]int) {

	sql := fmt.Sprintf(`SELECT x_host_header AS site_name,
			SUM(time_taken) AS sum_time_taken,
			COUNT(*) AS num_visits
		FROM %s.visits
		GROUP BY site_name`, collector.bigquery.DatasetName())

	rows := collector.bigquery.Query(sql, "site name to request seconds")

	siteNameToRequestSecondsSum := map[string]float64{}
	siteNameToRequestSecondsCount := map[string]int{}
	for _, row := range rows {
		siteName := row.F[0].V.(string)
		sumTimeTaken := mybigquery.ParseFloat64(row.F[1].V.(string))
		numVisits := mybigquery.Atoi(row.F[2].V.(string))
		siteNameToRequestSecondsSum[siteName] = sumTimeTaken
		siteNameToRequestSecondsCount[siteName] = numVisits
	}
	return siteNameToRequestSecondsSum, siteNameToRequestSecondsCount
}