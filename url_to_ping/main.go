package url_to_ping

import (
	"encoding/json"
	"github.com/prometheus/client_golang/prometheus"
	"io/ioutil"
	"log"
)

type Pop3Creds struct {
	Username string
	Password string
}

func Main(options *Options) {
	var jsonBytes []byte
	var err error
	if jsonBytes, err = ioutil.ReadFile(options.Pop3CredsJson); err != nil {
		log.Fatalf("Error from ReadFile: %v\n", err)
	}
	var pop3Creds Pop3Creds
	json.Unmarshal(jsonBytes, &pop3Creds)
	if pop3Creds.Username == "" {
		log.Fatalf("Missing Username in %s", options.Pop3CredsJson)
	}
	if pop3Creds.Password == "" {
		log.Fatalf("Missing Password in %s", options.Pop3CredsJson)
	}

	collector := NewUrlToPingCollector(options, &pop3Creds)
	prometheus.MustRegister(collector)
}
