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
	logLevel      string
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

func (cmd *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("run", flag.ExitOnError)

	flagset.StringVar(&cmd.configuration, "config", cmd.configuration, "Sets the configuration file path")
	flagset.StringVar(&cmd.dir, "dir", cmd.dir, "Work directory")
	flagset.StringVar(&cmd.pidFile, "pid", cmd.pidFile, "Sets the service PID file path")
	flagset.StringVar(&cmd.logLevel, "log-level", cmd.logLevel, "Sets the logging level (none/debug/info/warn/error)")
	flagset.StringVar(&cmd.logFile, "logfile", cmd.logFile, "Sets the log file path")
	flagset.IntVar(&cmd.logFileSize, "logfilesize", cmd.logFileSize, "Sets the log file size before forcing a log rotate")
	flagset.BoolVar(&cmd.console, "console", cmd.console, "Writes log entries to stdout")
	flagset.BoolVar(&cmd.debug, "debug", cmd.debug, "Displays internal information for diagnosing errors")

	return flagset
}

func (r *Run) Execute(args ...interface{}) error {
	log.Infof("", "uhppoted-rest daemon %s - %s (PID %d)\n", uhppote.VERSION, "Linux", os.Getpid())

	f := func(c *config.Config) error {
		return r.exec(c)
	}

	return r.execute(f)
}

func (cmd *Run) exec(c *config.Config) error {
	logger := syslog.New(os.Stdout, "", syslog.LstdFlags)

	if !cmd.console {
		events := eventlog.Ticker{Filename: cmd.logFile, MaxSize: cmd.logFileSize}
		logger = syslog.New(&events, "", syslog.Ldate|syslog.Ltime|syslog.LUTC)
		rotate := make(chan os.Signal, 1)

		signal.Notify(rotate, syscall.SIGHUP)

		go func() {
			for {
				<-rotate
				log.Infof("", "rotating uhppoted-rest log file '%s'\n", cmd.logFile)
				events.Rotate()
			}
		}()
	}

	log.SetLogger(logger)
	log.SetLevel(cmd.logLevel)
	cmd.run(c)

	return nil
}
