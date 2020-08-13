package commands

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"github.com/uhppoted/uhppoted-api/config"
)

type usergroup string

type Daemonize struct {
	usergroup usergroup
}

type info struct {
	Description      string
	Documentation    string
	Executable       string
	PID              string
	User             string
	Group            string
	Uid              int
	Gid              int
	LogFiles         []string
	BindAddress      *net.UDPAddr
	BroadcastAddress *net.UDPAddr
}

const serviceTemplate = `[Unit]
Description={{.Description}}
Documentation={{.Documentation}}
After=syslog.target network.target

[Service]
Type=simple
ExecStart={{.Executable}}
PIDFile={{.PID}}
User={{.User}}
Group={{.Group}}

[Install]
WantedBy=multi-user.target
`

const logRotateTemplate = `{{range .LogFiles}}{{.}} {
    daily
    rotate 30
    compress
        compresscmd /bin/bzip2
        compressext .bz2
        dateext
    missingok
    notifempty
    su uhppoted uhppoted
    postrotate
       /usr/bin/killall -HUP uhppoted
    endscript
}{{end}}
`

const confTemplate = `# UDP
bind.address = {{.BindAddress}}
broadcast.address = {{.BroadcastAddress}}

# REST API
rest.http.enabled = false
rest.http.port = 8080
rest.https.enabled = true
rest.https.port = 8443
rest.tls.key = /etc/uhppoted/rest/uhppoted.key
rest.tls.certificate = /etc/uhppoted/rest/uhppoted.cert
rest.tls.ca = /etc/uhppoted/rest/ca.cert

# OPEN API
# openapi.enabled = false
# openapi.directory = {{.WorkDir}}\rest\openapi

# DEVICES
# Example configuration for UTO311-L04 with serial number 405419896
# UT0311-L0x.405419896.address = 192.168.1.100:60000
# UT0311-L0x.405419896.door.1 = Front Door
# UT0311-L0x.405419896.door.2 = Side Door
# UT0311-L0x.405419896.door.3 = Garage
# UT0311-L0x.405419896.door.4 = Workshop
`

func NewDaemonize() *Daemonize {
	return &Daemonize{
		usergroup: "uhppoted:uhppoted",
	}
}

func (cmd *Daemonize) Name() string {
	return "daemonize"
}

func (cmd *Daemonize) FlagSet() *flag.FlagSet {
	flagset := flag.NewFlagSet("daemonize", flag.ExitOnError)
	flagset.Var(&cmd.usergroup, "user", "user:group for uhppoted service")

	return flagset
}

func (cmd *Daemonize) Description() string {
	return fmt.Sprintf("Registers %s as a service/daemon", SERVICE)
}

func (cmd *Daemonize) Usage() string {
	return "daemonize [--user <user:group>]"
}

func (cmd *Daemonize) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s daemonize [--user <user:group>]\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Registers %s as a systemd service/daemon that runs on startup.", SERVICE)
	fmt.Println("      Defaults to the user:group uhppoted:uhppoted unless otherwise specified")
	fmt.Println("      with the --user option")
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Daemonize) Execute(args ...interface{}) error {
	fmt.Println("   ... daemonizing")

	executable, err := os.Executable()
	if err != nil {
		return err
	}

	uid, gid, err := getUserGroup(string(cmd.usergroup))
	if err != nil {
		return err
	}

	bind, broadcast, _ := config.DefaultIpAddresses()

	d := info{
		Description:      "UHPPOTE UTO311-L0x access card controllers service/daemon ",
		Documentation:    "https://github.com/uhppoted/uhppoted-rest",
		Executable:       executable,
		PID:              "/var/uhppoted/uhppoted-rest.pid",
		User:             "uhppoted",
		Group:            "uhppoted",
		Uid:              uid,
		Gid:              gid,
		LogFiles:         []string{"/var/log/uhppoted/uhppoted-rest.log"},
		BindAddress:      &bind,
		BroadcastAddress: &broadcast,
	}

	if err := cmd.systemd(&d); err != nil {
		return err
	}

	if err := cmd.logrotate(&d); err != nil {
		return err
	}

	if err := cmd.mkdirs(&d); err != nil {
		return err
	}

	if err := cmd.conf(&d); err != nil {
		return err
	}

	fmt.Printf("   ... %s registered as a systemd service\n", SERVICE)
	fmt.Println()
	fmt.Println("   The daemon will start automatically on the next system restart - to start it manually, execute the following command:")
	fmt.Println()
	fmt.Printf("     > sudo systemctl start  %s", SERVICE)
	fmt.Printf("     > sudo systemctl status %s", SERVICE)
	fmt.Println()
	fmt.Println("   The firewall may need additional rules to allow UDP broadcast e.g. for UFW:")
	fmt.Println()
	fmt.Printf("     > sudo ufw allow from %s to any port 60000 proto udp\n", d.BindAddress.IP)
	fmt.Println()

	return nil
}

func (cmd *Daemonize) systemd(d *info) error {
	path := filepath.Join("/etc/systemd/system", "uhppoted-rest.service")
	t := template.Must(template.New("uhppoted-rest.service").Parse(serviceTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, d)
}

func (cmd *Daemonize) logrotate(d *info) error {
	path := filepath.Join("/etc/logrotate.d", "uhppoted-rest")
	t := template.Must(template.New("uhppoted-rest.logrotate").Parse(logRotateTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	return t.Execute(f, d)
}

func (cmd *Daemonize) conf(d *info) error {
	path := filepath.Join("/etc/uhppoted", "uhppoted.conf")
	t := template.Must(template.New("uhppoted.conf").Parse(confTemplate))

	fmt.Printf("   ... creating '%s'\n", path)
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer f.Close()

	err = t.Execute(f, d)
	if err != nil {
		return err
	}

	return os.Chown(path, d.Uid, d.Gid)
}

func (cmd *Daemonize) mkdirs(d *info) error {
	directories := []string{
		"/var/uhppoted",
		"/var/log/uhppoted",
		"/etc/uhppoted",
		"/etc/uhppoted/rest",
	}

	for _, dir := range directories {
		fmt.Printf("   ... creating '%s'\n", dir)

		if err := os.MkdirAll(dir, 0770); err != nil {
			return err
		}

		if err := os.Chown(dir, d.Uid, d.Gid); err != nil {
			return err
		}
	}

	return nil
}

func getUserGroup(s string) (int, int, error) {
	match := regexp.MustCompile(`(\w+?):(\w+)`).FindStringSubmatch(s)
	if match == nil {
		return 0, 0, fmt.Errorf("Invalid user:group '%s'", s)
	}

	u, err := user.Lookup(match[1])
	if err != nil {
		return 0, 0, err
	}

	g, err := user.LookupGroup(match[2])
	if err != nil {
		return 0, 0, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return 0, 0, err
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return 0, 0, err
	}

	return uid, gid, nil
}

// usergroup::flag.Value

func (f *usergroup) String() string {
	if f == nil {
		return "uhppoted:uhppoted"
	}

	return string(*f)
}

func (f *usergroup) Set(s string) error {
	_, _, err := getUserGroup(s)
	if err != nil {
		return err
	}

	*f = usergroup(strings.TrimSpace(s))

	return nil
}
