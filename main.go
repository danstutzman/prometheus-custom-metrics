package main

import (
	"encoding/json"
	"fmt"
	"github.com/danielstutzman/prometheus-cloudfront-logs-exporter/cloudfront_logs"
	"github.com/danielstutzman/prometheus-cloudfront-logs-exporter/json_value"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net/http"
	"os"
)

type Options struct {
	PortNum        int
	CloudfrontLogs *cloudfront_logs.Options
}

func usagef(format string, args ...interface{}) {
	log.Printf(`Usage: %s '{"PortNum":INT,  Port number to run web server on
  	"CloudfrontLogs": %s
	}`, os.Args[0], cloudfront_logs.Usage())
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
}
