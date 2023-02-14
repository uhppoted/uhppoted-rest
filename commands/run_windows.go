package commands

import (
	"flag"
	"fmt"
	syslog "log"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-lib/config"
	filelogger "github.com/uhppoted/uhppoted-lib/eventlog"
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

type service struct {
	name string
	conf *config.Config
	cmd  *Run
}

type EventLog struct {
	log *eventlog.Log
}

var RUN = Run{
	configuration: filepath.Join(workdir(), "uhppoted.conf"),
	dir:           workdir(),
	pidFile:       filepath.Join(workdir(), "uhppoted-rest.pid"),
	logFile:       filepath.Join(workdir(), "logs", "uhppoted-rest.log"),
	logFileSize:   10,
	console:       false,
	debug:         false,
}

func (cmd *Run) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("", flag.ExitOnError)

	flagset.StringVar(&cmd.configuration, "config", cmd.configuration, "Sets the configuration file path")
	flagset.StringVar(&cmd.dir, "dir", cmd.dir, "Work directory")
	flagset.StringVar(&cmd.pidFile, "pid", cmd.pidFile, "Sets the service PID file path")
	flagset.StringVar(&cmd.logLevel, "log-level", cmd.logLevel, "Sets the logging level (none/debug/info/warn/error)")
	flagset.StringVar(&cmd.logFile, "logfile", cmd.logFile, "Sets the log file path")
	flagset.IntVar(&cmd.logFileSize, "logfilesize", cmd.logFileSize, "Sets the log file size before forcing a log rotate")
	flagset.BoolVar(&cmd.console, "console", cmd.console, "Run as command-line application")
	flagset.BoolVar(&cmd.debug, "debug", cmd.debug, "Displays internal information for diagnosing errors")

	return flagset
}

func (r *Run) Execute(args ...interface{}) error {
	log.Infof("", "uhppoted-rest daemon %s - %s (PID %d)\n", uhppote.VERSION, "Microsoft Windows", os.Getpid())

	f := func(c *config.Config) error {
		return r.start(c)
	}

	return r.execute(f)
}

func (cmd *Run) start(c *config.Config) error {
	var logger *syslog.Logger

	if cmd.console {
		logger = syslog.New(os.Stdout, "", syslog.LstdFlags)
	} else if eventlogger, err := eventlog.Open(SERVICE); err == nil {
		defer eventlogger.Close()

		events := EventLog{eventlogger}
		logger = syslog.New(&events, SERVICE, syslog.Ldate|syslog.Ltime|syslog.LUTC)
	} else {
		events := filelogger.Ticker{Filename: cmd.logFile, MaxSize: cmd.logFileSize}
		logger = syslog.New(&events, "", syslog.Ldate|syslog.Ltime|syslog.LUTC)
	}

	log.SetLogger(logger)
	log.SetLevel(cmd.logLevel)
	log.Infof("", "uhppoted-rest service - start\n")

	if cmd.console {
		cmd.run(c)
		return nil
	}

	uhppoted := service{
		name: "uhppoted-rest",
		conf: c,
		cmd:  cmd,
	}

	logger.Printf("uhppoted-rest service - starting\n")

	if err := svc.Run("uhppoted-rest", &uhppoted); err != nil {
		fmt.Printf("   Unable to execute ServiceManager.Run request (%v)\n", err)
		fmt.Println()
		fmt.Println("   To run uhppoted-rest as a command line application, type:")
		fmt.Println()
		fmt.Println("     > uhppoted-rest --console")
		fmt.Println()

		logger.Fatalf("Error executing ServiceManager.Run request: %v", err)
		return err
	}

	logger.Printf("uhppoted-rest daemon - started\n")

	return nil
}

func (s *service) Execute(args []string, r <-chan svc.ChangeRequest, status chan<- svc.Status) (ssec bool, errno uint32) {
	log.Debugf("", "uhppoted-rest service - Execute")

	const commands = svc.AcceptStop | svc.AcceptShutdown

	status <- svc.Status{State: svc.StartPending}

	interrupt := make(chan os.Signal, 1)
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			err := s.cmd.listen(s.conf, interrupt)

			if err != nil {
				log.Errorf("", "%v", err)
				continue
			}

			log.Infof("", "exit")
			break
		}
	}()

	status <- svc.Status{State: svc.Running, Accepts: commands}

loop:
	for c := range r {
		log.Debugf("", "uhppoted-rest service - select: %v  %v\n", c.Cmd, c.CurrentStatus)
		switch c.Cmd {
		case svc.Interrogate:
			log.Debugf("", "uhppoted-rest service - svc.Interrogate %v\n", c.CurrentStatus)
			status <- c.CurrentStatus

		case svc.Stop:
			interrupt <- syscall.SIGINT
			log.Infof("", "uhppoted-rest service- svc.Stop\n")
			break loop

		case svc.Shutdown:
			interrupt <- syscall.SIGTERM
			log.Infof("", "uhppoted-rest service - svc.Shutdown\n")
			break loop

		default:
			log.Debugf("", "uhppoted-rest service - svc.????? (%v)\n", c.Cmd)
		}
	}

	log.Infof("", "uhppoted-rest service - stopping")
	status <- svc.Status{State: svc.StopPending}
	wg.Wait()
	status <- svc.Status{State: svc.Stopped}
	log.Infof("", "uhppoted-rest service - stopped\n")

	return false, 0
}

func (e *EventLog) Write(p []byte) (int, error) {
	err := e.log.Info(1, string(p))
	if err != nil {
		return 0, err
	}

	return len(p), nil
}
