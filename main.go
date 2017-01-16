package main

import (
	"encoding/json"
	"fmt"
	"github.com/danielstutzman/prometheus-custom-metrics/billing_gcloud"
	"github.com/danielstutzman/prometheus-custom-metrics/cloudfront_logs"
	"github.com/danielstutzman/prometheus-custom-metrics/memory_usage"
	"github.com/danielstutzman/prometheus-custom-metrics/papertrail_usage"
	"github.com/danielstutzman/prometheus-custom-metrics/piwik_exporter"
	"github.com/danielstutzman/prometheus-custom-metrics/security_updates"
	"github.com/danielstutzman/prometheus-custom-metrics/url_to_ping"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
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

func usagef(format string, args ...interface{}) {
	log.Printf(`Usage: %s '{"PortNum":INT,  Port number to run web server on
  	"BillingGcloud": %s,
  	"CloudfrontLogs": %s,
		"MemoryUsage": %s,
		"PapertrailUsage": %s,
		"PiwikExporter": %s,
		"SecurityUpdates": %s,
		"UrlToPing": %s
	}`, os.Args[0], billing_gcloud.Usage(), cloudfront_logs.Usage(),
		memory_usage.Usage(), papertrail_usage.Usage(), piwik_exporter.Usage(),
		security_updates.Usage(), url_to_ping.Usage())
	log.Fatalf(format, args...)
}

var collectorsByPort map[int][]prometheus.Collector

func addCollector(collector prometheus.Collector, metricsPort int) {
	if collectorsByPort == nil {
		collectorsByPort = map[int][]prometheus.Collector{}
	}

	_, ok := collectorsByPort[metricsPort]
	if !ok {
		collectorsByPort[metricsPort] = []prometheus.Collector{}
	}

	collectorsByPort[metricsPort] = append(collectorsByPort[metricsPort], collector)
}

func serveMetrics(collectors []prometheus.Collector, portNum int) {
	collectorNames := []string{}
	for _, collector := range collectors {
		collectorNames = append(collectorNames, fmt.Sprintf("%T", collector))
	}
	log.Printf("Starting web server on port %d for %s",
		portNum, strings.Join(collectorNames, ", "))

	registry := prometheus.NewPedanticRegistry()
	for _, collector := range collectors {
		registry.Register(collector)
	}

	serveMux := http.NewServeMux()
	serveMux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	err := http.ListenAndServe(fmt.Sprintf(":%d", portNum), serveMux)
	if err != nil {
		panic(err)
	}
}

func main() {
	if len(os.Args) == 1 {
		usagef("You must supply a command line argument")
	}
	if len(os.Args) > 2 {
		usagef("You must supply only one command line argument")
	}

	options := Options{}
	if err := json.Unmarshal([]byte(os.Args[1]), &options); err != nil {
		usagef("Error from json.Unmarshal of options: %v", err)
	}

	if options.BillingGcloud != nil {
		collector := billing_gcloud.MakeCollector(options.BillingGcloud)
		addCollector(collector, options.BillingGcloud.MetricsPort)
	}
	if options.CloudfrontLogs != nil {
		collector := cloudfront_logs.MakeCollector(options.CloudfrontLogs)
		addCollector(collector, options.CloudfrontLogs.MetricsPort)
		go collector.InitFromBigqueryAndS3() // run in the background since it's slow
	}
	if options.MemoryUsage != nil {
		addCollector(
			memory_usage.MakeCollector(options.MemoryUsage),
			options.MemoryUsage.MetricsPort)
	}
	if options.PapertrailUsage != nil {
		addCollector(
			papertrail_usage.MakeCollector(options.PapertrailUsage),
			options.PapertrailUsage.MetricsPort)
	}
	if options.PiwikExporter != nil {
		addCollector(
			piwik_exporter.MakeCollector(options.PiwikExporter),
			options.PiwikExporter.MetricsPort)
	}
	if options.SecurityUpdates != nil {
		addCollector(
			security_updates.MakeCollector(options.SecurityUpdates),
			options.SecurityUpdates.MetricsPort)
	}
	if options.UrlToPing != nil {
		addCollector(
			url_to_ping.MakeCollector(options.UrlToPing),
			options.UrlToPing.MetricsPort)
	}

	for portNum, collectors := range collectorsByPort {
		go serveMetrics(collectors, portNum)
	}

	runtime.Goexit() // don't exit main; keep running web server
}
