package commands

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/config"
	"github.com/uhppoted/uhppoted-api/eventlog"
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
	configuration: config.DefaultConfig,
	dir:           "/usr/local/var/com.github.uhppoted",
	pidFile:       "/usr/local/var/com.github.uhppoted/uhppoted-rest.pid",
	logFile:       "/usr/local/var/com.github.uhppoted/logs/uhppoted-rest.log",
	logFileSize:   10,
	console:       false,
	debug:         false,
}

func (cmd *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("run", flag.ExitOnError)

	flagset.StringVar(&cmd.configuration, "config", cmd.configuration, "Sets the configuration file path")
	flagset.StringVar(&cmd.dir, "dir", cmd.dir, "Work directory")
	flagset.StringVar(&cmd.pidFile, "pid", cmd.pidFile, "Sets the service PID file path")
	flagset.StringVar(&cmd.logFile, "logfile", cmd.logFile, "Sets the log file path")
	flagset.IntVar(&cmd.logFileSize, "logfilesize", cmd.logFileSize, "Sets the log file size before forcing a log rotate")
	flagset.BoolVar(&cmd.console, "console", cmd.console, "Writes log entries to stdout")
	flagset.BoolVar(&cmd.debug, "debug", cmd.debug, "Displays internal information for diagnosing errors")

	return flagset
}

func (cmd *Run) Execute(args ...interface{}) error {
	log.Printf("uhppoted-rest daemon %s - %s (PID %d)\n", uhppote.VERSION, "MacOS", os.Getpid())

	f := func(c *config.Config) error {
		return cmd.exec(c)
	}

	return cmd.execute(f)
}

func (cmd *Run) exec(c *config.Config) error {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	if !cmd.console {
		events := eventlog.Ticker{Filename: cmd.logFile, MaxSize: cmd.logFileSize}
		logger = log.New(&events, "", log.Ldate|log.Ltime|log.LUTC)
		rotate := make(chan os.Signal, 1)

		signal.Notify(rotate, syscall.SIGHUP)

		go func() {
			for {
				<-rotate
				log.Printf("Rotating uhppoted-rest log file '%s'\n", cmd.logFile)
				events.Rotate()
			}
		}()
	}

	cmd.run(c, logger)

	return nil
}
