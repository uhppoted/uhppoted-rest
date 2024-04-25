package commands

import (
	"errors"
	"fmt"
	"math"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/uhppoted/uhppoted-lib/config"
	"github.com/uhppoted/uhppoted-lib/locales"
	"github.com/uhppoted/uhppoted-lib/lockfile"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-rest/log"
	"github.com/uhppoted/uhppoted-rest/rest"
)

type status struct {
	touched time.Time
	status  types.Status
}

type alerts struct {
	missing      bool
	unexpected   bool
	touched      bool
	synchronized bool
}

type state struct {
	started time.Time

	healthcheck struct {
		touched *time.Time
		alerted bool
	}

	devices struct {
		status sync.Map
		errors sync.Map
	}
}

const (
	IDLE   = time.Duration(60 * time.Second)
	IGNORE = time.Duration(5 * time.Minute)
	DELTA  = 60
	DELAY  = 30
)

func (cmd *Run) Name() string {
	return "run"
}

func (cmd *Run) Description() string {
	return fmt.Sprintf("Runs the %s daemon/service until terminated by the system service manager", SERVICE)
}

func (cmd *Run) Usage() string {
	return fmt.Sprintf("%s [run] [--console] [--config <file>] [--dir <workdir>] [--pid <file>] [--logfile <file>] [--logfilesize <bytes>] [--debug]", SERVICE)
}

func (cmd *Run) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s [run] [--console] [--config <file>] [--dir <workdir>] [--pid <file>] [--logfile <file>] [--logfilesize <bytes>] [--debug]", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Run) execute(f func(*config.Config) error) error {
	conf := config.NewConfig()
	if err := conf.Load(cmd.configuration); err != nil {
		log.Warnf("RUN", "Could not load configuration (%v)", err)
	}

	if err := os.MkdirAll(cmd.dir, os.ModeDir|os.ModePerm); err != nil {
		return fmt.Errorf("unable to create working directory '%v': %v", cmd.dir, err)
	}

	// ... create lockfile
	lockFile := config.Lockfile{
		File:   cmd.pidFile,
		Remove: conf.LockfileRemove,
	}

	if lockFile.File == "" {
		lockFile.File = filepath.Join(os.TempDir(), fmt.Sprintf("%s.pid", SERVICE))
	}

	if kraken, err := lockfile.MakeLockFile(lockFile); err != nil {
		return err
	} else {
		defer func() {
			kraken.Release()
		}()

		log.AddFatalHook(func() {
			kraken.Release()
		})
	}

	// ... run
	return f(conf)
}

func (cmd *Run) run(c *config.Config) {
	log.Infof("", "START")

	// ... set (optional) locale
	if c.REST.Locale != "" {
		folder := filepath.Dir(c.REST.Locale)
		file := filepath.Base(c.REST.Locale)
		fs := os.DirFS(folder)
		if err := locales.Load(fs, file); err != nil {
			log.Warnf("", "%v", err)
		} else {
			log.Infof("", "using translations from %v", c.REST.Locale)
		}
	}

	// ... syscall SIG handlers

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	// ... listen forever

	for {
		err := cmd.listen(c, interrupt)
		if err != nil {
			log.Errorf("RUN", "%v", err)
			continue
		}

		log.Infof("RUN", "exit")
		break
	}

	log.Infof("RUN", "STOP")
}

func (cmd *Run) listen(c *config.Config, interrupt chan os.Signal) error {
	s := state{
		started: time.Now(),

		healthcheck: struct {
			touched *time.Time
			alerted bool
		}{
			touched: nil,
			alerted: false,
		},
		devices: struct {
			status sync.Map
			errors sync.Map
		}{
			status: sync.Map{},
			errors: sync.Map{},
		},
	}

	bind, broadcast, listen := config.DefaultIpAddresses()

	if c.BindAddress != nil {
		bind = *c.BindAddress
	}

	if c.BroadcastAddress != nil {
		broadcast = *c.BroadcastAddress
	}

	if c.ListenAddress != nil {
		listen = *c.ListenAddress
	}

	devices := []uhppote.Device{}
	for id, d := range c.Devices {
		if device := uhppote.NewDevice(d.Name, id, d.Address, d.Protocol, d.Doors); device != nil {
			devices = append(devices, *device)
		}
	}

	u := uhppote.NewUHPPOTE(bind, broadcast, listen, c.Timeout, devices, cmd.debug)

	// ... REST task

	restd := rest.RESTD{
		HTTPEnabled:               c.REST.HttpEnabled,
		HTTPPort:                  c.REST.HttpPort,
		HTTPSEnabled:              c.REST.HttpsEnabled,
		HTTPSPort:                 c.REST.HttpsPort,
		TLSKeyFile:                c.REST.TLSKeyFile,
		TLSCertificateFile:        c.REST.TLSCertificateFile,
		CACertificateFile:         c.REST.CACertificateFile,
		RequireClientCertificates: c.REST.RequireClientCertificates,
		CORSEnabled:               c.REST.CORSEnabled,
		AuthEnabled:               c.REST.AuthEnabled,
		AuthUsers:                 c.REST.Users,
		AuthGroups:                c.REST.Groups,
		HOTPWindow:                c.REST.HOTP.Range,
		HOTPCounters:              c.REST.HOTP.Counters,
		Protocol:                  c.REST.Protocol,
		OpenAPI: rest.OpenAPI{
			Enabled:   c.OpenAPI.Enabled,
			Directory: c.OpenAPI.Directory,
		},
	}

	go func() {
		restd.Run(u, devices)
	}()

	defer rest.Close()

	// ... health-check task
	k := time.NewTicker(15 * time.Second)

	defer k.Stop()

	go func() {
		for {
			<-k.C
			healthcheck(u, &s)
		}
	}()

	// ... wait until interrupted/closed

	closed := make(chan struct{})
	w := time.NewTicker(5 * time.Second)

	defer w.Stop()

	for {
		select {
		case <-w.C:
			if err := watchdog(u, &s); err != nil {
				return err
			}

		case <-interrupt:
			log.Infof("", "... interrupt")
			return nil

		case <-closed:
			log.Infof("", "... closed")

			return errors.New("server error")
		}
	}
}

func healthcheck(u uhppote.IUHPPOTE, st *state) {
	log.Debugf("health-check", "run")

	now := time.Now()
	devices := make(map[uint32]bool)

	if found, err := u.GetDevices(); err != nil {
		log.Warnf("health-check", "keep-alive error: %v", err)
	} else {
		for _, id := range found {
			devices[uint32(id.SerialNumber)] = true
		}
	}

	for id := range u.DeviceList() {
		devices[id] = true
	}

	for id := range devices {
		s, err := u.GetStatus(id)
		if err == nil {
			st.devices.status.Store(id, status{
				touched: now,
				status:  *s,
			})
		}
	}

	st.healthcheck.touched = &now
}

func watchdog(u uhppote.IUHPPOTE, st *state) error {
	warnings := 0
	errors := 0
	healthCheckRunning := false
	now := time.Now()
	seconds, _ := time.ParseDuration("1s")

	// Verify health-check

	dt := time.Since(st.started).Round(seconds)
	if st.healthcheck.touched != nil {
		dt = time.Since(*st.healthcheck.touched)
		if int64(math.Abs(dt.Seconds())) < DELAY {
			healthCheckRunning = true
		}
	}

	if int64(math.Abs(dt.Seconds())) > DELAY {
		errors += 1
		if !st.healthcheck.alerted {
			log.Errorf("watchdog", "subsystem has not run since %v (%v)", types.DateTime(st.started), dt)
			st.healthcheck.alerted = true
		}
	} else {
		if st.healthcheck.alerted {
			log.Infof("watchdog", "subsystem is running")
			st.healthcheck.alerted = false
		}
	}

	// Verify configured devices

	if healthCheckRunning {
		for id := range u.DeviceList() {
			alerted := alerts{
				missing:      false,
				unexpected:   false,
				touched:      false,
				synchronized: false,
			}

			if v, found := st.devices.errors.Load(id); found {
				alerted.missing = v.(alerts).missing
				alerted.unexpected = v.(alerts).unexpected
				alerted.touched = v.(alerts).touched
				alerted.synchronized = v.(alerts).synchronized
			}

			if _, found := st.devices.status.Load(id); !found {
				errors += 1
				if !alerted.missing {
					log.Errorf("health-check", "UTC0311-L0x %s device not found", types.SerialNumber(id))
					alerted.missing = true
				}
			}

			if v, found := st.devices.status.Load(id); found {
				touched := v.(status).touched
				t := time.Time(v.(status).status.SystemDateTime)
				dt := time.Since(t).Round(seconds)
				dtt := int64(math.Abs(time.Since(touched).Seconds()))

				if alerted.missing {
					log.Errorf("health-check", "UTC0311-L0x %s present", types.SerialNumber(id))
					alerted.missing = false
				}

				if now.After(touched.Add(IDLE)) {
					errors += 1
					if !alerted.touched {
						log.Errorf("health-check", "UTC0311-L0x %s no response for %s", types.SerialNumber(id), time.Since(touched).Round(seconds))
						alerted.touched = true
						alerted.synchronized = false
					}
				} else {
					if alerted.touched {
						log.Infof("health-check", "UTC0311-L0x %s connected", types.SerialNumber(id))
						alerted.touched = false
					}
				}

				if dtt < DELTA/2 {
					if int64(math.Abs(dt.Seconds())) > DELTA {
						errors += 1
						if !alerted.synchronized {
							log.Errorf("health-check", "UTC0311-L0x %v system time not synchronized: %v (%v)", types.SerialNumber(id), types.DateTime(t), dt)
							alerted.synchronized = true
						}
					} else {
						if alerted.synchronized {
							log.Infof("health-check", "UTC0311-L0x %v system time synchronized: %v (%v)", types.SerialNumber(id), types.DateTime(t), dt)
							alerted.synchronized = false
						}
					}
				}
			}

			st.devices.errors.Store(id, alerted)
		}
	}

	// Any unexpected devices?

	st.devices.status.Range(func(key, value interface{}) bool {
		alerted := alerts{
			missing:      false,
			unexpected:   false,
			touched:      false,
			synchronized: false,
		}

		if v, found := st.devices.errors.Load(key); found {
			alerted.missing = v.(alerts).missing
			alerted.unexpected = v.(alerts).unexpected
			alerted.touched = v.(alerts).touched
			alerted.synchronized = v.(alerts).synchronized
		}

		for id := range u.DeviceList() {
			if id == key {
				if alerted.unexpected {
					log.Errorf("health-check", "UTC0311-L0x %s added to configuration", types.SerialNumber(key.(uint32)))
					alerted.unexpected = false
					st.devices.errors.Store(id, alerted)
				}

				return true
			}
		}

		touched := value.(status).touched
		t := time.Time(value.(status).status.SystemDateTime)
		dt := time.Since(t).Round(seconds)
		dtt := int64(math.Abs(time.Since(touched).Seconds()))

		if now.After(touched.Add(IGNORE)) {
			st.devices.status.Delete(key)
			st.devices.errors.Delete(key)

			if alerted.unexpected {
				log.Warnf("health-check", "UTC0311-L0x %s disappeared", types.SerialNumber(key.(uint32)))
			}
		} else {
			warnings += 1
			if !alerted.unexpected {
				log.Warnf("health-check", "UTC0311-L0x %s unexpected device", types.SerialNumber(key.(uint32)))
				alerted.unexpected = true
			}

			if now.After(touched.Add(IDLE)) {
				warnings += 1
				if !alerted.touched {
					log.Warnf("health-check", "UTC0311-L0x %s no response for %s", types.SerialNumber(key.(uint32)), time.Since(touched).Round(seconds))
					alerted.touched = true
					alerted.synchronized = false
				}
			} else {
				if alerted.touched {
					log.Infof("health-check", "UTC0311-L0x %s connected", types.SerialNumber(key.(uint32)))
					alerted.touched = false
				}
			}

			if dtt < DELTA/2 {
				if int64(math.Abs(dt.Seconds())) > DELTA {
					warnings += 1
					if !alerted.synchronized {
						log.Warnf("health-check", "UTC0311-L0x %v system time not synchronized: %v (%v)", types.SerialNumber(key.(uint32)), types.DateTime(t), dt)
						alerted.synchronized = true
					}
				} else {
					if alerted.synchronized {
						log.Infof("health-check", "UTC0311-L0x %v system time synchronized: %v (%v)", types.SerialNumber(key.(uint32)), types.DateTime(t), dt)
						alerted.synchronized = false
					}
				}
			}

			st.devices.errors.Store(key, alerted)
		}

		return true
	})

	// 'k, done

	if errors > 0 {
		log.Errorf("watchdog", "%v errors", errors)
	} else if warnings > 0 {
		log.Warnf("watchdog", "%v warnings", warnings)
	} else {
		log.Infof("watchdog", "OK")
	}

	return nil
}
