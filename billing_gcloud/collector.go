package billing_gcloud

import (
	"github.com/prometheus/client_golang/prometheus"
)

type BillingGcloudCollector struct {
	options  *Options
	bigquery *BigqueryConnection
	desc     *prometheus.Desc
}

func (collector *BillingGcloudCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *BillingGcloudCollector) Collect(ch chan<- prometheus.Metric) {
	productToSumCost := collector.bigquery.QueryProductToSumCost()
	for product, sumCost := range productToSumCost {
		ch <- prometheus.MustNewConstMetric(
			collector.desc,
			prometheus.CounterValue,
			sumCost,
			product,
		)
	}
}

func NewBillingGcloudCollector(options *Options, bigquery *BigqueryConnection) *BillingGcloudCollector {
	return &BillingGcloudCollector{
		options:  options,
		bigquery: bigquery,
		desc: prometheus.NewDesc(
			"billing_gcloud_sum_cost_usd",
			"Total spent on Google Cloud since enabled export",
			[]string{"product"},
			prometheus.Labels{},
		),
	}
}
