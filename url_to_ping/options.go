package url_to_ping

import (
	"github.com/danielstutzman/prometheus-custom-metrics/json_value"
)

type Options struct {
	Pop3CredsJson     string
	EmailMaxAgeInMins int
	EmailSubject      string
	SuccessUrl        string
}

func Usage() string {
	return `{ (optional)
  	"Pop3CredsJson":     STRING,  path to file with {"Username":, "Password":}
		"EmailMaxAgeInMins": INT,     e.g. 60 to expect an email every hour
		"EmailSubject":      STRING   Specific subject line to look for
		"SuccessUrl":        STRING   e.g. "https://nosnch.in/abcdef"
	}`
}

func HandleOptions(section map[string]interface{}, path string,
	usagef func(string, ...interface{})) *Options {

	options := Options{}
	for key, value := range section {
		switch key {
		case "Pop3CredsJson":
			options.Pop3CredsJson = json_value.ToString(value, path+".Pop3CredsJson", usagef)
		case "EmailMaxAgeInMins":
			options.EmailMaxAgeInMins =
				json_value.ToInt(value, path+".EmailMaxAgeInMins", usagef)
		case "EmailSubject":
			options.EmailSubject = json_value.ToString(value, path+".EmailSubject", usagef)
		case "SuccessUrl":
			options.SuccessUrl = json_value.ToString(value, path+".SuccessUrl", usagef)
		default:
			usagef("Unknown key %s.%s", path, key)
		}
	}

	if options.Pop3CredsJson == "" {
		usagef("Missing %s.Pop3CredsJson", path)
	} else if options.EmailMaxAgeInMins == 0 {
		usagef("Need nonzero %s.EmailMaxAgeInMins", path)
	} else if options.EmailSubject == "" {
		usagef("Missing %s.EmailSubject", path)
	} else if options.SuccessUrl == "" {
		usagef("Missing %s.SuccessUrl", path)
	}

	return &options
}
