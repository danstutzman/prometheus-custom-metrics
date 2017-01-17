package main

import (
	"encoding/json"
	"fmt"
	"github.com/danielstutzman/prometheus-custom-metrics/collectors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

type Options struct {
	Collectors *collectors.Options
}

func usagef(format string, args ...interface{}) {
	log.Printf(`Usage: %s '{
  "Collectors": %s
}`, os.Args[0], collectors.Usage())
	log.Fatalf(format, args...)
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

	if options.Collectors == nil {
		usagef("Missing options.Collectors")
	}
	collectorsByPort := collectors.Setup(options.Collectors)
	if len(collectorsByPort) == 0 {
		log.Fatalf("No collectors were set up")
	}

	for portNum, collectors := range collectorsByPort {
		go serveMetrics(collectors, portNum)
	}

	runtime.Goexit() // don't exit main; keep running web server
}
