package commands

import (
	"flag"
	syslog "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/eventlog"
	"github.com/uhppoted/uhppoted-rest/log"
)

type Run struct {
	configuration string
	dir           string
	pidFile       string
	logFile       string
	logFileSize   int
	console       bool
	debug         bool
}

var RUN = Run{
	configuration: "/etc/uhppoted/uhppoted.conf",
	dir:           "/var/uhppoted",
	pidFile:       "/var/uhppoted/uhppoted-rest.pid",
	logFile:       "/var/log/uhppoted/uhppoted-rest.log",
	logFileSize:   10,
	console:       false,
	debug:         false,
}

func (r *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("run", flag.ExitOnError)

	flagset.StringVar(&r.configuration, "config", r.configuration, "Sets the configuration file path")
	flagset.StringVar(&r.dir, "dir", r.dir, "Work directory")
	flagset.StringVar(&r.pidFile, "pid", r.pidFile, "Sets the service PID file path")
	flagset.StringVar(&r.logFile, "logfile", r.logFile, "Sets the log file path")
	flagset.IntVar(&r.logFileSize, "logfilesize", r.logFileSize, "Sets the log file size before forcing a log rotate")
	flagset.BoolVar(&r.console, "console", r.console, "Writes log entries to stdout")
	flagset.BoolVar(&r.debug, "debug", r.debug, "Displays internal information for diagnosing errors")

	return flagset
}

func (r *Run) Execute(args ...interface{}) error {
	log.Infof("", "uhppoted-rest daemon %s - %s (PID %d)\n", uhppote.VERSION, "Linux", os.Getpid())

	f := func(c *config.Config) error {
		return r.exec(c)
	}

	return r.execute(f)
}

func (r *Run) exec(c *config.Config) error {
	logger := syslog.New(os.Stdout, "", syslog.LstdFlags)

	if !r.console {
		events := eventlog.Ticker{Filename: r.logFile, MaxSize: r.logFileSize}
		logger = syslog.New(&events, "", syslog.Ldate|syslog.Ltime|syslog.LUTC)
		rotate := make(chan os.Signal, 1)

		signal.Notify(rotate, syscall.SIGHUP)

		go func() {
			for {
				<-rotate
				log.Infof("", "rotating uhppoted-rest log file '%s'\n", r.logFile)
				events.Rotate()
			}
		}()
	}

	log.SetLogger(logger)
	r.run(c)

	return nil
}
