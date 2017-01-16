package billing_gcloud

import (
	"github.com/prometheus/client_golang/prometheus"
)

type BillingGcloudCollector struct {
	options *Options
	desc    *prometheus.Desc
}

func (collector *BillingGcloudCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *BillingGcloudCollector) Collect(ch chan<- prometheus.Metric) {
	//usage := collector.queryUsage()
	ch <- prometheus.MustNewConstMetric(
		collector.desc,
		prometheus.CounterValue,
		0.00,
		"test",
	)
}

func NewBillingGcloudCollector(options *Options) *BillingGcloudCollector {
	return &BillingGcloudCollector{
		options: options,
		desc: prometheus.NewDesc(
			"billing_gcloud_sum_cost_usd",
			"Total spent on Google Cloud since enabled export",
			[]string{"product"},
			prometheus.Labels{},
		),
	}
}
