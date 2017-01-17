package collectors

import (
	"fmt"
	"github.com/danielstutzman/prometheus-custom-metrics/collectors/billing_gcloud"
	"github.com/danielstutzman/prometheus-custom-metrics/collectors/cloudfront_logs"
	"github.com/danielstutzman/prometheus-custom-metrics/collectors/memory_usage"
	"github.com/danielstutzman/prometheus-custom-metrics/collectors/papertrail_usage"
	"github.com/danielstutzman/prometheus-custom-metrics/collectors/piwik_exporter"
	"github.com/danielstutzman/prometheus-custom-metrics/collectors/security_updates"
	"github.com/danielstutzman/prometheus-custom-metrics/collectors/url_to_ping"
	"github.com/prometheus/client_golang/prometheus"
)

type Options struct {
	BillingGcloud   *billing_gcloud.Options
	CloudfrontLogs  *cloudfront_logs.Options
	MemoryUsage     *memory_usage.Options
	PapertrailUsage *papertrail_usage.Options
	PiwikExporter   *piwik_exporter.Options
	SecurityUpdates *security_updates.Options
	UrlToPing       *url_to_ping.Options
}

type CollectorsByPort map[int][]prometheus.Collector

func NewCollectorsByPort() CollectorsByPort {
	return map[int][]prometheus.Collector{}
}

func (collectorsByPort CollectorsByPort) addCollector(collector prometheus.Collector,
	metricsPort int) {
	_, ok := collectorsByPort[metricsPort]
	if !ok {
		collectorsByPort[metricsPort] = []prometheus.Collector{}
	}

	collectorsByPort[metricsPort] = append(collectorsByPort[metricsPort], collector)
}

func Usage() string {
	return fmt.Sprintf(`{
    "BillingGcloud": %s, "CloudfrontLogs": %s, "MemoryUsage": %s, "PapertrailUsage": %s, "PiwikExporter": %s, "SecurityUpdates": %s, "UrlToPing": %s
  }`,
		billing_gcloud.Usage(),
		cloudfront_logs.Usage(),
		memory_usage.Usage(),
		papertrail_usage.Usage(),
		piwik_exporter.Usage(),
		security_updates.Usage(),
		url_to_ping.Usage(),
	)
}

func Setup(opts *Options) CollectorsByPort {
	collectorsByPort := NewCollectorsByPort()
	add := collectorsByPort.addCollector
	if opts.BillingGcloud != nil {
		collector := billing_gcloud.MakeCollector(opts.BillingGcloud)
		add(collector, opts.BillingGcloud.MetricsPort)
	}
	if opts.CloudfrontLogs != nil {
		collector := cloudfront_logs.MakeCollector(opts.CloudfrontLogs)
		add(collector, opts.CloudfrontLogs.MetricsPort)
		go collector.InitFromBigqueryAndS3() // run in the background since it's slow
	}
	if opts.MemoryUsage != nil {
		add(memory_usage.MakeCollector(opts.MemoryUsage),
			opts.MemoryUsage.MetricsPort)
	}
	if opts.PapertrailUsage != nil {
		add(papertrail_usage.MakeCollector(opts.PapertrailUsage),
			opts.PapertrailUsage.MetricsPort)
	}
	if opts.PiwikExporter != nil {
		add(piwik_exporter.MakeCollector(opts.PiwikExporter),
			opts.PiwikExporter.MetricsPort)
	}
	if opts.SecurityUpdates != nil {
		add(security_updates.MakeCollector(opts.SecurityUpdates),
			opts.SecurityUpdates.MetricsPort)
	}
	if opts.UrlToPing != nil {
		add(url_to_ping.MakeCollector(opts.UrlToPing),
			opts.UrlToPing.MetricsPort)
	}
	return collectorsByPort
}
