package billing_aws

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/s3"
	"github.com/prometheus/client_golang/prometheus"
)

type BillingAwsCollector struct {
	options  *Options
	bigquery *bigquery.BigqueryConnection
	s3       *s3.S3Connection
	desc     *prometheus.Desc
}

func (collector *BillingAwsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *BillingAwsCollector) Collect(ch chan<- prometheus.Metric) {
	bucketNameToNumBytes := collector.downloadBucketNameToNumBytes()
	for bucketName, numBytes := range bucketNameToNumBytes {
		ch <- prometheus.MustNewConstMetric(
			collector.desc,
			prometheus.GaugeValue,
			numBytes,
			bucketName,
		)
	}
}

func NewBillingAwsCollector(options *Options,
	bigqueryConn *bigquery.BigqueryConnection,
	s3Conn *s3.S3Connection) *BillingAwsCollector {

	return &BillingAwsCollector{
		options:  options,
		bigquery: bigqueryConn,
		s3:       s3Conn,
		desc: prometheus.NewDesc(
			"billing_aws_s3_bucket_usage_bytes",
			"Amazon Web Services S3 storage usage, by bucket, in bytes",
			[]string{"bucket"},
			prometheus.Labels{},
		),
	}
}
