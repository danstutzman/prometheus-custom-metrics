package papertrail_usage

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
	"net/http"
)

const PAPERTRAIL_USAGE_URL = "https://papertrailapp.com/api/v1/accounts.json"

type PapertrailUsageCollector struct {
	options       *Options
	jsonKeyToDesc map[string]*prometheus.Desc
}

func (collector *PapertrailUsageCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, desc := range collector.jsonKeyToDesc {
		ch <- desc
	}
}

func (collector *PapertrailUsageCollector) queryUsage() map[string]float64 {
	client := &http.Client{}
	req, err := http.NewRequest("GET", PAPERTRAIL_USAGE_URL, nil)
	req.Header.Add("X-Papertrail-Token", collector.options.ApiToken)
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error from http.Get of %s: %s", PAPERTRAIL_USAGE_URL, err)
	}
	defer resp.Body.Close()

	usageJson, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error from ioutil.ReadAll of %s: %s", PAPERTRAIL_USAGE_URL, err)
	}

	usage := map[string]float64{}
	err = json.Unmarshal([]byte(usageJson), &usage)
	if err != nil {
		log.Printf("Error from json.Unmarshal of %s: %s", usage, err)
	}

	return usage
}

func (collector *PapertrailUsageCollector) Collect(ch chan<- prometheus.Metric) {
	usage := collector.queryUsage()
	for jsonKey, desc := range collector.jsonKeyToDesc {
		ch <- prometheus.MustNewConstMetric(
			desc,
			prometheus.GaugeValue,
			usage[jsonKey],
		)
	}
}

func NewPapertrailUsageCollector(options *Options) *PapertrailUsageCollector {
	return &PapertrailUsageCollector{
		options: options,
		jsonKeyToDesc: map[string]*prometheus.Desc{
			"log_data_transfer_used": prometheus.NewDesc(
				"papertrail_log_data_transfer_used",
				"Amount of monthly storage used for Papertrail (third-party log aggregator)",
				[]string{},
				prometheus.Labels{},
			),
			"log_data_transfer_used_percent": prometheus.NewDesc(
				"papertrail_log_data_transfer_used_percent",
				"Percent of monthly storage used for Papertrail (third-party log aggregator)",
				[]string{},
				prometheus.Labels{},
			),
			"log_data_transfer_plan_limit": prometheus.NewDesc(
				"papertrail_log_data_transfer_plan_limit",
				"Available monthly storage for Papertrail (third-party log aggregator)",
				[]string{},
				prometheus.Labels{},
			),
			"log_data_transfer_hard_limit": prometheus.NewDesc(
				"papertrail_log_data_transfer_hard_limit",
				"Available monthly storage for Papertrail (third-party log aggregator)",
				[]string{},
				prometheus.Labels{},
			),
		},
	}
}
