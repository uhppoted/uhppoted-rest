package commands

import (
	"flag"
	"fmt"
	xml "github.com/uhppoted/uhppoted-lib/encoding/plist"
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
	fmt.Printf("    Deregisters %s from launchd as a service/daemon\n", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Undaemonize) Execute(args ...interface{}) error {
	fmt.Println("   ... undaemonizing")

	path := filepath.Join("/Library/LaunchDaemons", "com.github.uhppoted-rest.plist")
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if os.IsNotExist(err) {
		fmt.Printf("   ... %s does not exist - nothing to do for launchd\n", path)
		return nil
	}

	p, err := cmd.parse(path)
	if err != nil {
		return err
	}

	if err := cmd.launchd(path, *p); err != nil {
		return err
	}

	if err := cmd.logrotate(); err != nil {
		return err
	}

	if err := cmd.rmdirs(*p); err != nil {
		return err
	}

	if err := cmd.firewall(*p); err != nil {
		return err
	}

	fmt.Println("   ... com.github.uhppoted.uhppoted-rest unregistered as a LaunchDaemon")
	fmt.Println()
	fmt.Println("   Any uhppoted-rest log files can still be found in directory /usr/local/var/log and should be removed manually")
	fmt.Println()

	return nil
}

func (cmd *Undaemonize) parse(path string) (*info, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	p := info{}
	decoder := xml.NewDecoder(f)
	err = decoder.Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (cmd *Undaemonize) launchd(path string, d info) error {
	fmt.Printf("   ... unloading LaunchDaemon\n")
	command := exec.Command("launchctl", "unload", path)
	out, err := command.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("failed to unload '%s' (%v)", d.Label, err)
	}

	fmt.Printf("   ... removing '%s'\n", path)
	err = os.Remove(path)
	if err != nil {
		return err
	}

	return nil
}

func (cmd *Undaemonize) logrotate() error {
	path := filepath.Join("/etc/newsyslog.d", "uhppoted-rest.conf")

	fmt.Printf("   ... removing '%s'\n", path)

	return os.Remove(path)
}

func (cmd *Undaemonize) rmdirs(d info) error {
	dir := d.WorkDir

	fmt.Printf("   ... removing '%s'\n", dir)

	return os.RemoveAll(dir)
}

func (cmd *Undaemonize) firewall(d info) error {
	fmt.Println()
	fmt.Println("   ***")
	fmt.Printf("   *** WARNING: removing '%s' from the application firewall", SERVICE)
	fmt.Println("   ***")
	fmt.Println()

	path := d.Executable
	command := exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	out, err := command.CombinedOutput()
	fmt.Printf("   > %s", out)
	if err != nil {
		return fmt.Errorf("failed to retrieve application firewall global state (%v)", err)
	}

	if strings.Contains(string(out), "State = 1") {
		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "off")
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("failed to disable the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--remove", path)
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("failed to remove 'uhppoted-rest' from the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--setglobalstate", "on")
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("failed to re-enable the application firewall (%v)", err)
		}

		fmt.Println()
	}

	return nil
}
