package cloudfront_logs

import (
	"github.com/danielstutzman/prometheus-custom-metrics/util"
	"github.com/prometheus/client_golang/prometheus"
)

type CloudfrontCollector struct {
	bigquery                          *util.BigqueryConnection
	s3                                *S3Connection
	siteNameStatusToNumVisits         map[SiteNameStatus]int
	siteNameStatusToNumVisitsDesc     *prometheus.Desc
	siteNameToRequestSecondsSum       map[string]float64
	siteNameToRequestSecondsSumDesc   *prometheus.Desc
	siteNameToRequestSecondsCount     map[string]int
	siteNameToRequestSecondsCountDesc *prometheus.Desc
}

func (collector *CloudfrontCollector) InitFromBigqueryAndS3() {
	collector.siteNameStatusToNumVisits =
		QuerySiteNameStatusToNumVisits(collector.bigquery)
	collector.siteNameToRequestSecondsSum,
		collector.siteNameToRequestSecondsCount =
		QuerySiteNameToRequestSeconds(collector.bigquery)
	collector.syncNewCloudfrontLogsToBigquery()
}

func (collector *CloudfrontCollector) syncNewCloudfrontLogsToBigquery() {
	sem := make(chan bool, 5)
	for _, s3Path := range collector.s3.ListPaths() {
		sem <- true
		go func(s3Path string) {
			visits := collector.s3.DownloadVisitsForPath(s3Path)
			for _, visit := range visits {
				siteName := visit["x-host-header"]
				siteNameStatus := SiteNameStatus{
					siteName,
					rollUpExactStatus(util.Atoi(visit["sc-status"])),
				}
				collector.siteNameStatusToNumVisits[siteNameStatus] += 1

				collector.siteNameToRequestSecondsSum[siteName] +=
					util.ParseFloat64(visit["time-taken"])
				collector.siteNameToRequestSecondsCount[siteName] += 1
			}
			UploadVisits(collector.bigquery, s3Path, visits)
			collector.s3.DeletePath(s3Path)
			<-sem
		}(s3Path)
	}
	for i := 0; i < cap(sem); i++ {
		sem <- true
	}
}

func (collector *CloudfrontCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.siteNameStatusToNumVisitsDesc
	ch <- collector.siteNameToRequestSecondsSumDesc
	ch <- collector.siteNameToRequestSecondsCountDesc
}

func (collector *CloudfrontCollector) Collect(ch chan<- prometheus.Metric) {
	collector.syncNewCloudfrontLogsToBigquery()
	for siteNameStatus, numVisits := range collector.siteNameStatusToNumVisits {
		ch <- prometheus.MustNewConstMetric(
			collector.siteNameStatusToNumVisitsDesc,
			prometheus.CounterValue,
			float64(numVisits),
			siteNameStatus.SiteName,
			siteNameStatus.Status,
		)
	}
	for siteName, requestSecondsSum := range collector.siteNameToRequestSecondsSum {
		ch <- prometheus.MustNewConstMetric(
			collector.siteNameToRequestSecondsSumDesc,
			prometheus.CounterValue,
			float64(requestSecondsSum),
			siteName,
		)
		ch <- prometheus.MustNewConstMetric(
			collector.siteNameToRequestSecondsCountDesc,
			prometheus.CounterValue,
			float64(collector.siteNameToRequestSecondsCount[siteName]),
			siteName,
		)
	}
}

func NewCloudfrontCollector(s3 *S3Connection,
	bigquery *util.BigqueryConnection) *CloudfrontCollector {

	return &CloudfrontCollector{
		s3:       s3,
		bigquery: bigquery,
		siteNameStatusToNumVisitsDesc: prometheus.NewDesc(
			"cloudfront_visits",
			"Number of visits in CloudFront S3 logs.",
			[]string{"site_name", "status"},
			prometheus.Labels{},
		),
		siteNameToRequestSecondsSumDesc: prometheus.NewDesc(
			"cloudfront_request_seconds_sum",
			"Total duration of requests in CloudFront S3 logs.",
			[]string{"site_name"},
			prometheus.Labels{},
		),
		siteNameToRequestSecondsCountDesc: prometheus.NewDesc(
			"cloudfront_request_seconds_count",
			"Number of requests in CloudFront S3 logs.",
			[]string{"site_name"},
			prometheus.Labels{},
		),
	}
}
