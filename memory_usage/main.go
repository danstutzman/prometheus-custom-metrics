package memory_usage

func MakeCollector(options *Options) *MemoryUsageCollector {
	validateOptions(options)
	return NewMemoryUsageCollector()
}
