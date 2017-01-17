package url_to_ping

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Options struct {
	MetricsPort       int
	Pop3CredsJson     string
	Pop3Creds         *Pop3Creds
	EmailMaxAgeInMins int
	EmailSubject      string
	SuccessUrl        string
}

func Usage() string {
	return `{ (optional)
      "MetricsPort":       INT,     port to run web server on, e.g. 9102
      "Pop3CredsJson":     STRING,  path to file with {"Username":, "Password":}
      "EmailMaxAgeInMins": INT,     e.g. 60 to expect an email every hour
      "EmailSubject":      STRING   Specific subject line to look for
      "SuccessUrl":        STRING   e.g. "https://nosnch.in/abcdef"
    }`
}

type Pop3Creds struct {
	Username string
	Password string
}

func validateOptions(options *Options) {
	if options.MetricsPort == 0 {
		log.Fatalf("Missing url_to_ping.MetricsPort")
	}
	if options.Pop3CredsJson == "" {
		log.Fatalf("Missing url_to_ping.Pop3CredsJson")
	}
	if options.EmailMaxAgeInMins == 0 {
		log.Fatalf("Missing url_to_ping.EmailMaxAgeInMins")
	}
	if options.EmailSubject == "" {
		log.Fatalf("Missing url_to_ping.EmailSubject")
	}
	if options.SuccessUrl == "" {
		log.Fatalf("Missing url_to_ping.SuccessUrl")
	}

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
	options.Pop3Creds = &pop3Creds
}
