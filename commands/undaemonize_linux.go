package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Undaemonize struct {
}

func NewUndaemonize() *Undaemonize {
	return &Undaemonize{}
}

func (cmd *Undaemonize) Name() string {
	return "undaemonize"
}

func (cmd *Undaemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("undaemonize", flag.ExitOnError)
}

func (cmd *Undaemonize) Description() string {
	return fmt.Sprintf("Deregisters %s as a service/daemon", SERVICE)
}

func (cmd *Undaemonize) Usage() string {
	return ""
}

func (cmd *Undaemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s undaemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Deregisters %s as a systed service/daemon\n", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Undaemonize) Execute(args ...interface{}) error {
	fmt.Println("   ... undaemonizing")

	path := filepath.Join("/etc/systemd/system", "uhppoted-rest.service")
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		fmt.Printf("   ... %s does not exist - nothing to do for systemd\n", path)
		return nil
	}

	if err := cmd.systemd(path); err != nil {
		return err
	}

	if err := cmd.logrotate(); err != nil {
		return err
	}

	if err := cmd.rmdirs(); err != nil {
		return err
	}

	fmt.Println("   ... uhppoted-rest unregistered as a systemd service")
	fmt.Println()
	fmt.Println("   Log files in directory /var/log/uhppoted and configuration files in /etc/uhppoted should be removed manually")
	fmt.Println()

	return nil
}

func (cmd *Undaemonize) systemd(path string) error {
	fmt.Printf("   ... stopping uhppoted-rest service\n")
	command := exec.Command("systemctl", "stop", "uhppoted-rest")
	out, err := command.CombinedOutput()
	if strings.TrimSpace(string(out)) != "" {
		fmt.Printf("   > %s\n", out)
	}
	if err != nil {
		return fmt.Errorf("failed to stop '%s' (%v)", "uhppoted-rest", err)
	}

	fmt.Printf("   ... removing '%s'\n", path)
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Undaemonize) logrotate() error {
	path := filepath.Join("/etc/logrotate.d", "uhppoted-rest")

	fmt.Printf("   ... removing '%s'\n", path)

	err := os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Undaemonize) rmdirs() error {
	dir := "/var/uhppoted/rest"

	fmt.Printf("   ... removing '%s'\n", dir)

	return os.RemoveAll(dir)
}
