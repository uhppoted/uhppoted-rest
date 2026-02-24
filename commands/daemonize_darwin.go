package commands

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/uhppoted/uhppote-core/types"
	"github.com/uhppoted/uhppoted-lib/config"
	xpath "github.com/uhppoted/uhppoted-lib/encoding/plist"
)

type info struct {
	Label            string
	Executable       string
	ConfigDirectory  string
	WorkDir          string
	BindAddress      *types.BindAddr
	BroadcastAddress *types.BroadcastAddr
}

type plist struct {
	Label             string
	Program           string
	WorkingDirectory  string
	ProgramArguments  []string
	KeepAlive         bool
	RunAtLoad         bool
	StandardOutPath   string
	StandardErrorPath string
}

const newsyslog = `#logfilename                                       [owner:group]  mode  count  size   when  flags [/pid_file]  [sig_num]
{{range .}}{{.LogFile}}  :              644   30     10000  @T00  J     {{.PID}}
{{end}}`

type Daemonize struct {
}

func NewDaemonize() *Daemonize {
	return &Daemonize{}
}

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("daemonize", flag.ExitOnError)
}

func (cmd *Daemonize) Description() string {
	return "Daemonizes uhppoted as a service/daemon"
}

func (cmd *Daemonize) Usage() string {
	return ""
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Daemonizes %s as a service/daemon that runs on startup\n", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Daemonize) Execute(args ...any) error {
	fmt.Println("   ... daemonizing")
	executable, err := os.Executable()
	if err != nil {
		return err
	}

	bind, broadcast, _ := config.DefaultIpAddresses()

	d := info{
		Label:            "com.github.uhppoted-rest",
		Executable:       executable,
		ConfigDirectory:  "/usr/local/etc/com.github.uhppoted",
		WorkDir:          "/usr/local/var/com.github.uhppoted",
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
	}

	if err := cmd.launchd(&d); err != nil {
		return err
	}

	if err := cmd.mkdirs(); err != nil {
		return err
	}

	if err := cmd.logrotate(); err != nil {
		return err
	}

	if err := cmd.firewall(); err != nil {
		return err
	}

	if err := cmd.conf(&d); err != nil {
		return err
	}

	fmt.Println("   ... com.github.uhppoted-rest registered as a LaunchDaemon")
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Println("   sudo launchctl load /Library/LaunchDaemons/com.github.uhppoted-rest.plist")
	fmt.Println()

	return nil
}

func (cmd *Daemonize) launchd(d *info) error {
	path := filepath.Join("/Library/LaunchDaemons", "com.github.uhppoted-rest.plist")
	_, err := os.Stat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	pl := plist{
		Label:             d.Label,
		Program:           d.Executable,
		WorkingDirectory:  "/usr/local/var/com.github.uhppoted",
		ProgramArguments:  []string{},
		KeepAlive:         true,
		RunAtLoad:         true,
		StandardOutPath:   "/usr/local/var/log/com.github.uhppoted-rest.log",
		StandardErrorPath: "/usr/local/var/log/com.github.uhppoted-rest.log",
	}

	if !os.IsNotExist(err) {
		current, err := cmd.parse(path)
		if err != nil {
			return err
		}

		pl.WorkingDirectory = current.WorkingDirectory
		pl.ProgramArguments = current.ProgramArguments
		pl.KeepAlive = current.KeepAlive
		pl.RunAtLoad = current.RunAtLoad
		pl.StandardOutPath = current.StandardOutPath
		pl.StandardErrorPath = current.StandardErrorPath
	}

	return cmd.daemonize(path, pl)
}

func (cmd *Daemonize) parse(path string) (*plist, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	p := plist{}
	decoder := xpath.NewDecoder(f)
	err = decoder.Decode(&p)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (cmd *Daemonize) daemonize(path string, p any) error {
	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	encoder := xpath.NewEncoder(f)
	if err = encoder.Encode(p); err != nil {
		return err
	}

	return nil
}

func (cmd *Daemonize) mkdirs() error {
	dir := "/usr/local/var/com.github.uhppoted"

	fmt.Printf("   ... creating '%s'\n", dir)

	return os.MkdirAll(dir, 0644)
}

func (cmd *Daemonize) conf(d *info) error {
	path := filepath.Join(d.ConfigDirectory, "uhppoted.conf")

	fmt.Printf("   ... creating '%s'\n", path)

	// initialise config from existing uhppoted.conf
	cfg := config.NewConfig()
	if f, err := os.Open(path); err != nil && !os.IsNotExist(err) {
		return err
	} else if err == nil {
		err := cfg.Read(f)
		f.Close()
		if err != nil {
			return err
		}
	}

	// write back config with any updated information
	if f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644); err != nil {
		return err
	} else {
		defer f.Close()

		return cfg.Write(f)
	}
}

func (cmd *Daemonize) logrotate() error {
	dir := "/usr/local/var/log"
	pid := "/usr/local/var/com.github.uhppoted/uhppoted-rest.pid"
	logfiles := []struct {
		LogFile string
		PID     string
	}{
		{
			LogFile: filepath.Join(dir, "com.github.uhppoted-rest.log"),
			PID:     pid,
		},
		{
			LogFile: filepath.Join(dir, "com.github.uhppoted-rest.err"),
			PID:     pid,
		},
	}

	t := template.Must(template.New("logrotate.conf").Parse(newsyslog))
	path := filepath.Join("/etc/newsyslog.d", "uhppoted-rest.conf")

	fmt.Printf("   ... creating '%s'\n", path)

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, logfiles)
}

func (cmd *Daemonize) firewall() error {
	fmt.Println()
	fmt.Println("   ***")
	fmt.Println("   *** WARNING: adding 'uhppoted-rest' to the application firewall and unblocking incoming connections")
	fmt.Println("   ***")
	fmt.Println()

	path, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get path to executable: %v", err)
	}

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

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--add", path)
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("failed to add 'uhppoted-rest' to the application firewall (%v)", err)
		}

		command = exec.Command("/usr/libexec/ApplicationFirewall/socketfilterfw", "--unblockapp", path)
		out, err = command.CombinedOutput()
		fmt.Printf("   > %s", out)
		if err != nil {
			return fmt.Errorf("failed to unblock 'uhppoted-rest' on the application firewall (%v)", err)
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
