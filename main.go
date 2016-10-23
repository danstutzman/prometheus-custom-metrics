package main

import (
	"encoding/json"
	"fmt"
	"github.com/danielstutzman/prometheus-custom-metrics/cloudfront_logs"
	"github.com/danielstutzman/prometheus-custom-metrics/json_value"
	"github.com/danielstutzman/prometheus-custom-metrics/piwik_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
	"runtime"
)

type Options struct {
	PortNum        int
	CloudfrontLogs *cloudfront_logs.Options
	PiwikExporter  bool
}

func usagef(format string, args ...interface{}) {
	log.Printf(`Usage: %s '{"PortNum":INT,  Port number to run web server on
  	"CloudfrontLogs": %s,
		"PiwikExporter": %s
	}`, os.Args[0], cloudfront_logs.Usage(), piwik_exporter.Usage())
	log.Fatalf(format, args...)
}

func handleOptions(optionsMap map[string]interface{}) Options {
	options := Options{}

	for key, value := range optionsMap {
		switch key {
		case "PortNum":
			options.PortNum = json_value.ToInt(value, "Options.PortNum", usagef)
		case "CloudfrontLogs":
			options.CloudfrontLogs = cloudfront_logs.HandleOptions(
				json_value.ToMap(value, "Options.CloudfrontLogs", usagef),
				"Options.CloudfrontLogs", usagef)
		case "PiwikExporter":
			options.PiwikExporter = json_value.ToBool(value, "Options.PiwikExporter", usagef)
		default:
			usagef("Unknown key \"%s\" in options", key)
		}
	}

	if options.PortNum == 0 {
		usagef("Missing Options.PortNum")
	}

	return options
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

	optionsMap := map[string]interface{}{}
	if err := json.Unmarshal([]byte(os.Args[1]), &optionsMap); err != nil {
		usagef("Error from json.Unmarshal of options: %v", err)
	}
	options := handleOptions(optionsMap)

	go serveMetrics(options.PortNum)

	if options.CloudfrontLogs != nil {
		cloudfront_logs.Main(*options.CloudfrontLogs)
	}
	if options.PiwikExporter {
		piwik_exporter.Main()
	}

	runtime.Goexit() // don't exit main; keep running web server
}
