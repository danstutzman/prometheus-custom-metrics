package cpu

import (
	"github.com/sirupsen/logrus"
)

func MakeCollector(options *Options, log *logrus.Logger) *CpuCollector {
	validateOptions(options)
	return NewCpuCollector(log)
}
