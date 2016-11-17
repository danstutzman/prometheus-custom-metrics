package url_to_ping

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
