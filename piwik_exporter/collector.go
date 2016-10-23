package piwik_exporter

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ziutek/mymysql/mysql"
	"io"
	"strconv"
	"time"
)

type PiwikCollector struct {
	db                 mysql.Conn
	idsiteToSiteName   map[int]string
	descForPiwikVisits *prometheus.Desc
	descForQueryTime   *prometheus.Desc
}

func atoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(fmt.Errorf("Error from Atoi of '%s': %s", s, err))
	}
	return i
}

func queryIdsiteToNumVisits(db mysql.Conn) map[int]int {
	idsiteToNumVisits := map[int]int{}

	sql := `SELECT idsite,
		      COUNT(*) as num_visits
		    FROM piwik_log_visit
		    GROUP BY idsite;`
	result, err := db.Start(sql)
	if err != nil {
		panic(fmt.Errorf("Error from db.Start with sql %s: %s", sql, err))
	}

	row := result.MakeRow()
	for {
		err := result.ScanRow(row)
		if err == io.EOF {
			break // No more rows
		}
		if err != nil {
			panic(fmt.Errorf("Error from ScanRow: %s", err))
		}

		idsite := atoi(string(row[0].([]byte)))
		numVisits := atoi(string(row[1].([]byte)))
		idsiteToNumVisits[idsite] = numVisits
	}

	return idsiteToNumVisits
}

func queryIdsiteToSiteName(db mysql.Conn) map[int]string {
	idsiteToSiteName := map[int]string{}

	sql := `SELECT idsite, name FROM piwik_site;`
	result, err := db.Start(sql)
	if err != nil {
		panic(fmt.Errorf("Error from db.Start with sql %s: %s", sql, err))
	}

	row := result.MakeRow()
	for {
		err := result.ScanRow(row)
		if err == io.EOF {
			break // No more rows
		}
		if err != nil {
			panic(fmt.Errorf("Error from ScanRow: %s", err))
		}

		idsite := atoi(string(row[0].([]byte)))
		siteName := string(row[1].([]byte))
		idsiteToSiteName[idsite] = siteName
	}

	return idsiteToSiteName
}

func (collector *PiwikCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.descForPiwikVisits
}

func (collector *PiwikCollector) Collect(ch chan<- prometheus.Metric) {
	timeBeforeNumVisitsQuery := time.Now()
	idsiteToNumVisits := queryIdsiteToNumVisits(collector.db)
	timeOfNumVisitsQuery := time.Since(timeBeforeNumVisitsQuery)

	areAllSiteNamesKnown := true
	for idsite, _ := range idsiteToNumVisits {
		if _, known := collector.idsiteToSiteName[idsite]; !known {
			areAllSiteNamesKnown = false
		}
	}
	if !areAllSiteNamesKnown {
		collector.idsiteToSiteName = queryIdsiteToSiteName(collector.db)
	}

	for idsite, numVisits := range idsiteToNumVisits {
		ch <- prometheus.MustNewConstMetric(
			collector.descForPiwikVisits,
			prometheus.CounterValue,
			float64(numVisits),
			collector.idsiteToSiteName[idsite],
		)
	}
	ch <- prometheus.MustNewConstMetric(
		collector.descForQueryTime,
		prometheus.GaugeValue,
		// report as integer for better compression
		timeOfNumVisitsQuery.Seconds(),
	)
}

func NewPiwikCollector(db mysql.Conn) *PiwikCollector {
	return &PiwikCollector{
		db:               db,
		idsiteToSiteName: map[int]string{},
		descForPiwikVisits: prometheus.NewDesc(
			"piwik_visits",
			"Number of site visits in Piwik database.",
			[]string{"site_name"},
			prometheus.Labels{},
		),
		descForQueryTime: prometheus.NewDesc(
			"piwik_query_seconds",
			"Duration of Piwik database query in seconds",
			[]string{},
			prometheus.Labels{},
		),
	}
}
