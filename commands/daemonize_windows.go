package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/config"
)

type Daemonize struct {
	name        string
	description string
}

type info struct {
	Name             string
	Description      string
	Executable       string
	WorkDir          string
	BindAddress      *types.BindAddr
	BroadcastAddress *types.BroadcastAddr
}

func NewDaemonize() *Daemonize {
	return &Daemonize{
		name:        "uhppoted-rest",
		description: "UHPPOTE UTO311-L0x access card controllers service/daemon",
	}
}

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("daemonize", flag.ExitOnError)
}

func (cmd *Daemonize) Description() string {
	return fmt.Sprintf("Registers %s as a Windows service", SERVICE)
}

func (cmd *Daemonize) Usage() string {
	return ""
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Registers %s as a windows Service that runs on startup\n", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Daemonize) Execute(args ...interface{}) error {
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	bind, broadcast, _ := config.DefaultIpAddresses()

	d := info{
		Name:             cmd.name,
		Description:      cmd.description,
		Executable:       executable,
		WorkDir:          workdir(),
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
	}

	if err := cmd.register(&d); err != nil {
		return err
	}

	if err := cmd.mkdirs(&d); err != nil {
		return err
	}

	if err := cmd.conf(&d); err != nil {
		return err
	}

	fmt.Printf("   ... %s registered as a Windows system service\n", SERVICE)
	fmt.Println()
	fmt.Println("   The service will start automatically on the next system restart. Start it manually from the")
	fmt.Println("   'Services' application or from the command line by executing the following command:")
	fmt.Println()
	fmt.Printf("     > net start %s", SERVICE)
	fmt.Println()

	return nil
}

func (cmd *Daemonize) register(d *info) error {
	config := mgr.Config{
		DisplayName:      d.Name,
		Description:      d.Description,
		StartType:        mgr.StartAutomatic,
		DelayedAutoStart: true,
	}

	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(d.Name)
	if err == nil {
		s.Close()
		return fmt.Errorf("service %s already exists", d.Name)
	}

	s, err = m.CreateService(d.Name, d.Executable, config, "is", "auto-started")
	if err != nil {
		return err
	}

	defer s.Close()

	err = eventlog.InstallAsEventCreate(d.Name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		s.Delete()
		return fmt.Errorf("InstallAsEventCreate() failed: %v", err)
	}

	return nil
}

func (cmd *Daemonize) mkdirs(d *info) error {
	directories := []string{
		d.WorkDir,
		filepath.Join(d.WorkDir, "rest"),
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}
	}

	return nil
}

func (cmd *Daemonize) conf(d *info) error {
	path := filepath.Join(d.WorkDir, "uhppoted.conf")

	fmt.Printf("   ... creating '%s'\n", path)

	// initialise config from existing uhppoted.conf
	cfg := config.NewConfig()
	if f, err := os.Open(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		err := cfg.Read(f)
		f.Close()
		if err != nil {
			return err
		}
	}

	// replace line endings and write back config with any updated information
	var b strings.Builder
	var replacer = strings.NewReplacer(
		"\r\n", "\r\n",
		"\r", "\r\n",
		"\n", "\r\n",
	)

	if err := cfg.Write(&b); err != nil {
		return err
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.Write([]byte(replacer.Replace(b.String()))); err != nil {
		return err
	}

	return nil
}
