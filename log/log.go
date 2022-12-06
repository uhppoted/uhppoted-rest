package log

import (
	"fmt"
	syslog "log"

	"github.com/uhppoted/uhppoted-lib/log"
)

type LogLevel int

const (
	none LogLevel = iota
	debug
	info
	warn
	errors
)

const f = "%-12v %v"

func SetDebug(enabled bool) {
	log.SetDebug(enabled)
}

func SetLevel(level string) {
	log.SetLevel(level)
}

func SetLogger(logger *syslog.Logger) {
	log.SetLogger(logger)
}

func AddFatalHook(f func()) {
	log.AddFatalHook(f)
}

func Debugf(tag string, format string, args ...any) {
	s := fmt.Sprintf(f, tag, format)

	log.Debugf(s, args...)
}

func Infof(tag string, format string, args ...any) {
	s := fmt.Sprintf(f, tag, format)

	log.Infof(s, args...)
}

func Warnf(tag string, format string, args ...any) {
	s := fmt.Sprintf(f, tag, format)

	log.Warnf(s, args...)
}

func Errorf(tag string, format string, args ...any) {
	s := fmt.Sprintf(f, tag, format)

	log.Errorf(s, args...)
}

func Fatalf(tag string, format string, args ...any) {
	s := fmt.Sprintf(f, tag, format)

	log.Fatalf(s, args...)
}
