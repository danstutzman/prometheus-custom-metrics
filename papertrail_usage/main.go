package papertrail_usage

func MakeCollector(options *Options) *PapertrailUsageCollector {
	validateOptions(options)
	return NewPapertrailUsageCollector(options)
}
