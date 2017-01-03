package papertrail_usage

import (
	"io/ioutil"
	"log"
	"strings"
)

type Options struct {
	ApiTokenPath string
	ApiToken     string
	MetricsPort  int
}

func Usage() string {
	return `{ (optional)
	   "ApiTokenPath": STRING   text file containing Papertrail API token
     "MetricsPort":  INT      port to serve metrics on, e.g. 9102
	}`
}

func validateOptions(options *Options) {
	if options.ApiTokenPath == "" {
		log.Fatalf("Missing memory_usage.ApiTokenPath")
	}
	if options.MetricsPort == 0 {
		log.Fatalf("Missing memory_usage.MetricsPort")
	}

	var bytes []byte
	var err error
	if bytes, err = ioutil.ReadFile(options.ApiTokenPath); err != nil {
		log.Fatalf("Error from ReadFile: %v\n", err)
	}
	options.ApiToken = strings.TrimSpace(string(bytes))
}
