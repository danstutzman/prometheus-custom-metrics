package security_updates

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"os"
	"os/exec"
	"regexp"
	"strconv"
)

var CMD_LINE = []string{"/usr/lib/update-notifier/apt-check", "--human-readable"}

var OUTPUT_REGEXP = regexp.MustCompile(`([0-9])+ packages can be updated.
([0-9]+) updates are security updates.
`)

type SecurityUpdatesCollector struct {
	descForSecurityUpdates *prometheus.Desc
	descForRebootRequired  *prometheus.Desc
}

func (collector *SecurityUpdatesCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.descForSecurityUpdates
	ch <- collector.descForRebootRequired
}

func (collector *SecurityUpdatesCollector) Collect(ch chan<- prometheus.Metric) {
	cmd := exec.Command(CMD_LINE[0], CMD_LINE[1:]...)
	output, err := cmd.Output()
	if err != nil {
		panic(fmt.Errorf("Error from exec.Command %s: %", CMD_LINE, err))
	}

	match := OUTPUT_REGEXP.FindStringSubmatch(string(output))
	if match == nil {
		panic(fmt.Errorf("Output '%s' from %s didn't match regex %s", output, CMD_LINE,
			OUTPUT_REGEXP))
	}
	numSecurityUpdates, err := strconv.Atoi(match[2])
	if err != nil {
		panic(fmt.Errorf("Error from Atoi of '%s': %s", match[2], err))
	}

	ch <- prometheus.MustNewConstMetric(
		collector.descForSecurityUpdates,
		prometheus.GaugeValue,
		float64(numSecurityUpdates))

	var isRebootRequired int
	_, err = os.Stat("/var/run/reboot-required")
	if os.IsNotExist(err) {
		isRebootRequired = 0
	} else {
		isRebootRequired = 1
	}

	ch <- prometheus.MustNewConstMetric(
		collector.descForRebootRequired,
		prometheus.GaugeValue,
		float64(isRebootRequired))
}

func NewSecurityUpdatesCollector() *SecurityUpdatesCollector {
	return &SecurityUpdatesCollector{
		descForSecurityUpdates: prometheus.NewDesc(
			"ubuntu_security_updates",
			"Number of Ubuntu packages that are security updates",
			[]string{},
			prometheus.Labels{},
		),
		descForRebootRequired: prometheus.NewDesc(
			"is_reboot_required",
			"Nonzero if Ubuntu applied security updates that now require a reboot",
			[]string{},
			prometheus.Labels{},
		),
	}
}
