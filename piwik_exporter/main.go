package piwik_exporter

import (
	"github.com/ziutek/mymysql/mysql"
	_ "github.com/ziutek/mymysql/native"
	"log"
)

func MakeCollector(options *Options) *PiwikCollector {
	validateOptions(options)

	db := mysql.New("tcp", "", "127.0.0.1:3306", "root", "", "piwik")
	err := db.Connect()
	if err != nil {
		log.Fatalf("Couldn't db.Connect to MySQL database: %s", err)
	}

	return NewPiwikCollector(db)
}
