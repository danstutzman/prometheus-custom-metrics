package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

type CloudfrontCollector struct {
	bigquery                  *BigqueryConnection
	s3                        *S3Connection
	siteNameStatusToNumVisits map[SiteNameStatus]int
	desc                      *prometheus.Desc
}

func (collector *CloudfrontCollector) InitFromBigqueryAndS3() {
	collector.siteNameStatusToNumVisits =
		collector.bigquery.QuerySiteNameStatusToNumVisits()
	collector.syncNewCloudfrontLogsToBigquery()
}

func (collector *CloudfrontCollector) syncNewCloudfrontLogsToBigquery() {
	for _, s3Path := range collector.s3.ListPaths() {
		visits := collector.s3.DownloadVisitsForPath(s3Path)
		for _, visit := range visits {
			siteNameStatus := SiteNameStatus{
				visit["x-host-header"],
				RollUpExactStatus(atoi(visit["sc-status"])),
			}
			collector.siteNameStatusToNumVisits[siteNameStatus] += 1
		}
		collector.bigquery.UploadVisits(s3Path, visits)
		collector.s3.DeletePath(s3Path)
	}
}

func (collector *CloudfrontCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *CloudfrontCollector) Collect(ch chan<- prometheus.Metric) {
	collector.syncNewCloudfrontLogsToBigquery()
	for siteNameStatus, numVisits := range collector.siteNameStatusToNumVisits {
		ch <- prometheus.MustNewConstMetric(
			collector.desc,
			prometheus.CounterValue,
			float64(numVisits),
			siteNameStatus.SiteName,
			siteNameStatus.Status,
		)
	}
}

func NewCloudfrontCollector(s3 *S3Connection,
	bigquery *BigqueryConnection) *CloudfrontCollector {

	return &CloudfrontCollector{
		s3:       s3,
		bigquery: bigquery,
		desc: prometheus.NewDesc(
			"cloudfront_visits",
			"Number of visits in CloudFront S3 logs.",
			[]string{"site_name", "status"},
			prometheus.Labels{},
		),
	}
}
