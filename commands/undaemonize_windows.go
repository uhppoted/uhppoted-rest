package commands

import (
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
)

type Undaemonize struct {
	name        string
	description string
}

func NewUndaemonize() *Undaemonize {
	return &Undaemonize{
		name:        "uhppoted-rest",
		description: "uhppoted-rest Service Interface to UTO311-L0x devices",
	}
}

func (cmd *Undaemonize) Name() string {
	return "undaemonize"
}

func (cmd *Undaemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("undaemonize", flag.ExitOnError)
}

func (cmd *Undaemonize) Execute(args ...interface{}) error {
	fmt.Println("   ... undaemonizing")

	dir := workdir()
	if err := cmd.unregister(); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted-rest unregistered as a Windows system service")
	fmt.Println()
	fmt.Printf("   Log files and configuration files in directory %s should be removed manually", dir)
	fmt.Println()

	return nil
}

func (cmd *Undaemonize) unregister() error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}

	defer m.Disconnect()

	s, err := m.OpenService(cmd.name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", cmd.name)
	}

	defer s.Close()

	err = s.Delete()
	if err != nil {
		return err
	}

	err = eventlog.Remove(cmd.name)
	if err != nil {
		return fmt.Errorf("RemoveEventLogSource() failed: %s", err)
	}

	return nil
}

func (cmd *Undaemonize) Description() string {
	return fmt.Sprintf("Deregisters the %s service", SERVICE)
}

func (cmd *Undaemonize) Usage() string {
	return ""
}

func (cmd *Undaemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s undaemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Deregisters %s as a Windows service\n", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}
