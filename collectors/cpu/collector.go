package cpu

import (
	"bufio"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"time"
)

var USER_HZ = time.Duration(100) * time.Second

var CPU_LINE_REGEXP = regexp.MustCompile(`^cpu\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s([0-9]+)\s+([0-9]+)\s+([0-9]+)\s+([0-9]+)\s([0-9]+)`)

type ProcStatCpu struct {
	user       int
	nice       int
	system     int
	idle       int
	iowait     int
	irq        int
	softirq    int
	steal_time int
}

type CpuCollector struct {
	desc    *prometheus.Desc
	log     *logrus.Logger
	lastCpu *ProcStatCpu
}

func (collector *CpuCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *CpuCollector) atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		collector.log.Fatalf("Error from Atoi of '%s': %s", s, err)
	}
	return i
}

func (collector *CpuCollector) Collect(ch chan<- prometheus.Metric) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		collector.log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var cpuLine *ProcStatCpu
	for scanner.Scan() {
		line := scanner.Text()
		match := CPU_LINE_REGEXP.FindStringSubmatch(string(line))
		if match != nil {
			cpuLine = &ProcStatCpu{
				user:       collector.atoi(match[1]),
				nice:       collector.atoi(match[2]),
				system:     collector.atoi(match[3]),
				idle:       collector.atoi(match[4]),
				iowait:     collector.atoi(match[5]),
				irq:        collector.atoi(match[6]),
				softirq:    collector.atoi(match[7]),
				steal_time: collector.atoi(match[8]),
			}
		}
	}
	if cpuLine == nil {
		collector.log.Fatalf("No line from /proc/stat matched OUTPUT_REGEXP")
	}
	if err := scanner.Err(); err != nil {
		collector.log.Fatalf("Error from scanner.Err: %s", err)
	}

	last := collector.lastCpu
	if last != nil {
		timeSinceLast := (cpuLine.user - last.user) +
			(cpuLine.nice - last.nice) +
			(cpuLine.system - last.system) +
			(cpuLine.idle - last.idle) +
			(cpuLine.iowait - last.iowait) +
			(cpuLine.irq - last.irq) +
			(cpuLine.softirq - last.softirq) +
			(cpuLine.steal_time - last.steal_time)
		idleSinceLast := cpuLine.idle - last.idle

		ch <- prometheus.MustNewConstMetric(
			collector.desc,
			prometheus.GaugeValue,
			100.0*(1.0-(float64(idleSinceLast)/float64(timeSinceLast))),
		)
	}

	collector.lastCpu = cpuLine
}

func NewCpuCollector(log *logrus.Logger) *CpuCollector {
	return &CpuCollector{
		desc: prometheus.NewDesc(
			"cpu_usage_percent",
			"Percent cpu usage from /proc/stat (all CPUs)",
			[]string{},
			prometheus.Labels{},
		),
		log: log,
	}
}
