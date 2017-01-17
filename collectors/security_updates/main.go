package security_updates

func MakeCollector(options *Options) *SecurityUpdatesCollector {
	validateOptions(options)
	return NewSecurityUpdatesCollector()
}
