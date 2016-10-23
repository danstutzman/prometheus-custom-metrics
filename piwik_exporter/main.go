package piwik_exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"log"
)

func Usage() string {
	return "true"
}

func Main() {
	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "", "piwik")
	err := db.Connect()
	if err != nil {
		log.Fatalf("Couldn't db.Connect to MySQL database: %s", err)
	}

	prometheus.MustRegister(NewPiwikCollector(db))
}
