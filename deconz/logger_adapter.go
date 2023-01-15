package deconz

import "github.com/PerformLine/go-stockutil/log"

type Logger struct {
}

func (l Logger) Errorf(format string, v ...interface{}) {
	log.Fatalf(format, v...)
}

func (l Logger) Warnf(format string, v ...interface{}) {
	log.Warningf(format, v...)
}

func (l Logger) Debugf(format string, v ...interface{}) {
	log.Debugf(format, v...)
}

func (l Logger) Infof(format string, v ...interface{}) {
	log.Infof(format, v...)
}
