package commands

import (
	"flag"
	"fmt"
	"path/filepath"
)

var CONFIG = Dump{
	config: filepath.Join(workdir(), "uhppoted.conf"),
}

type Dump struct {
	config string
}

func (cmd *Dump) Name() string {
	return "config"
}

func (cmd *Dump) FlagSet() *flag.FlagSet {
	return flag.NewFlagSet("config", flag.ExitOnError)
}

func (cmd *Dump) Description() string {
	return fmt.Sprintf("Displays all the configuration information for %s", SERVICE)
}

func (cmd *Dump) Usage() string {
	return ""
}

func (cmd *Dump) Help() {
	fmt.Println()
	fmt.Printf("  Usage: %s config\n", SERVICE)
	fmt.Println()
	fmt.Printf("    Displays all the configuration information for %s\n", SERVICE)
	fmt.Println()

	helpOptions(cmd.FlagSet())
}

func (cmd *Dump) Execute(args ...interface{}) error {
	if err := dump(cmd.config); err != nil {
		return err
	}

	return nil
}
