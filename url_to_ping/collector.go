package url_to_ping

import (
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
	"net/http"
)

type UrlToPingCollector struct {
	url         string
	desc        *prometheus.Desc
	numRequests int
}

func (collector *UrlToPingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *UrlToPingCollector) Collect(ch chan<- prometheus.Metric) {
	resp, err := http.Get(collector.url)
	if err != nil {
		log.Printf("Error from http.Get of %s: %s", collector.url, err)
	}
	defer resp.Body.Close()
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error from ioutil.ReadAll of %s: %s", collector.url, err)
	}

	collector.numRequests += 1
	ch <- prometheus.MustNewConstMetric(
		collector.desc,
		prometheus.CounterValue,
		float64(collector.numRequests),
	)
}

func NewUrlToPingCollector(url string) *UrlToPingCollector {
	return &UrlToPingCollector{
		url: url,
		desc: prometheus.NewDesc(
			"url_to_ping_requests",
			"Number of times collector has hit url_to_ping",
			[]string{},
			prometheus.Labels{},
		),
		numRequests: 0,
	}
}
