package main

import (
	"encoding/json"
	"fmt"
	"github.com/danielstutzman/prometheus-custom-metrics/cloudfront_logs"
	"github.com/danielstutzman/prometheus-custom-metrics/memory_usage"
	"github.com/danielstutzman/prometheus-custom-metrics/piwik_exporter"
	"github.com/danielstutzman/prometheus-custom-metrics/security_updates"
	"github.com/danielstutzman/prometheus-custom-metrics/url_to_ping"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Options struct {
	PortNum         int
	CloudfrontLogs  *cloudfront_logs.Options
	MemoryUsage     bool
	PiwikExporter   bool
	SecurityUpdates bool
	UrlToPing       *url_to_ping.Options
}

func usagef(format string, args ...interface{}) {
	log.Printf(`Usage: %s '{"PortNum":INT,  Port number to run web server on
  	"CloudfrontLogs": %s,
		"MemoryUsage": %s,
		"PiwikExporter": %s,
		"SecurityUpdates": %s,
		"UrlToPing": %s
	}`, os.Args[0], cloudfront_logs.Usage(), memory_usage.Usage(),
		piwik_exporter.Usage(), security_updates.Usage(), url_to_ping.Usage())
	log.Fatalf(format, args...)
}

func serveMetrics(portNum int) {
	http.Handle("/metrics", prometheus.Handler())
	err := http.ListenAndServe(fmt.Sprintf(":%d", portNum), nil)
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

	go serveMetrics(options.PortNum)

	if options.MemoryUsage {
		memory_usage.Main()
	}
	if options.PiwikExporter {
		piwik_exporter.Main()
	}
	if options.SecurityUpdates {
		security_updates.Main()
	}
	if options.UrlToPing != nil {
		url_to_ping.Main(options.UrlToPing)
	}
	if options.CloudfrontLogs != nil { // Run last since it's slow
		cloudfront_logs.Main(options.CloudfrontLogs)
	}

	runtime.Goexit() // don't exit main; keep running web server
}
