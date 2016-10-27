package memory_usage

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var PS_CMD_LINE = []string{"/bin/ps", "axwwo", "rss,args"}

type RegexpAndGroup struct {
	regexp *regexp.Regexp
	group  string
}

var REGEXP_TO_GROUP_MAPPINGS = []RegexpAndGroup{
	{regexp.MustCompile(`^[^ ]*blackbox_exporter`), "blackbox_exporter"},
	{regexp.MustCompile(`^[^ ]*grafana-server`), "grafana-server"},
	{regexp.MustCompile(`^nginx`), "nginx"},
	{regexp.MustCompile(`^[^ ]*node_exporter`), "node_exporter"},
	{regexp.MustCompile(`^[^ ]*apache2`), "piwik"},
	{regexp.MustCompile(`^[^ ]*mysqld`), "piwik"},
	{regexp.MustCompile(`^php-fpm`), "piwik"},
	{regexp.MustCompile(`^[^ ]*postgres`), "postgres"},
	{regexp.MustCompile(`^[^ ]*postgres_exporter`), "postgres_exporter"},
	{regexp.MustCompile(`^[^ ]*prometheus`), "prometheus"},
	{regexp.MustCompile(`^[^ ]*prometheus-custom-metrics`), "prometheus-custom-metrics"},
	{regexp.MustCompile(`^[^ ]*remote_syslog`), "remote_syslog"},
	{regexp.MustCompile(`^[^ ]*rsyslogd`), "rsyslogd"},
	{regexp.MustCompile(`^[^ ]*unicorn`), "unicorn"},
	{regexp.MustCompile(`/var/www/vocabincontext/golang/backend`), "vocabincontext"},
}

var OTHER_GROUP = "other"

type MemoryUsageCollector struct {
	desc *prometheus.Desc
}

func (collector *MemoryUsageCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func grepAndGroupPsOutput() map[string]float64 {
	cmd := exec.Command(PS_CMD_LINE[0], PS_CMD_LINE[1:]...)
	output, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("Error from exec.Command %s: %", PS_CMD_LINE, err))
	}

	groupToMegabytes := map[string]float64{}
	for _, line := range strings.Split(string(output), "\n") {
		values := strings.Fields(line)
		if len(values) >= 2 {
			kilobytes := values[0]
			binary := values[1]

			group := OTHER_GROUP
			for _, regexpToGroup := range REGEXP_TO_GROUP_MAPPINGS {
				if regexpToGroup.regexp.Match([]byte(binary)) {
					group = regexpToGroup.group
					break
				}
			}

			kilobytesInt, err := strconv.Atoi(kilobytes)
			if err == nil {
				megabytes := float64(kilobytesInt) / 1024.0
				groupToMegabytes[group] += megabytes
			}
		}
	}
	return groupToMegabytes
}

func (collector *MemoryUsageCollector) Collect(ch chan<- prometheus.Metric) {
	groupToMegabytes := grepAndGroupPsOutput()
	for group, megabytes := range groupToMegabytes {
		ch <- prometheus.MustNewConstMetric(
			collector.desc,
			prometheus.GaugeValue,
			megabytes,
			group,
		)
	}
}

func NewMemoryUsageCollector() *MemoryUsageCollector {
	return &MemoryUsageCollector{
		desc: prometheus.NewDesc(
			"memory_usage_in_mb",
			"Memory usage by related processes in megabytes",
			[]string{"process_group"},
			prometheus.Labels{},
		),
	}
}
