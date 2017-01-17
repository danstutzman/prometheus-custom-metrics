package url_to_ping

import (
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type UrlToPingCollector struct {
	options     *Options
	desc        *prometheus.Desc
	numRequests int
}

func (collector *UrlToPingCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *UrlToPingCollector) Collect(ch chan<- prometheus.Metric) {
	since := time.Now().Add(
		time.Duration(-1*collector.options.EmailMaxAgeInMins) * time.Minute)
	if MailboxHasMailWithSubject(collector.options.Pop3Creds.Username,
		collector.options.Pop3Creds.Password, since, collector.options.EmailSubject) {

		resp, err := http.Get(collector.options.SuccessUrl)
		if err != nil {
			log.Printf("Error from http.Get of %s: %s", collector.options.SuccessUrl, err)
		}
		defer resp.Body.Close()
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error from ioutil.ReadAll of %s: %s",
				collector.options.SuccessUrl, err)
		}

		collector.numRequests += 1
	}

	ch <- prometheus.MustNewConstMetric(
		collector.desc,
		prometheus.CounterValue,
		float64(collector.numRequests),
	)
}

func NewUrlToPingCollector(options *Options) *UrlToPingCollector {
	return &UrlToPingCollector{
		options: options,
		desc: prometheus.NewDesc(
			"url_to_ping_requests",
			"Number of times collector has hit url_to_ping",
			[]string{},
			prometheus.Labels{},
		),
		numRequests: 0,
	}
}
