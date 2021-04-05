package main

import (
	"fmt"
	"os"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/command"
	"github.com/uhppoted/uhppoted-rest/commands"
)

var cli = []uhppoted.Command{
	&commands.RUN,
	&commands.CONFIG,
	commands.NewDaemonize(),
	commands.NewUndaemonize(),
	&uhppoted.Version{
		Application: commands.SERVICE,
		Version:     uhppote.VERSION,
	},
}

var help = uhppoted.NewHelp(commands.SERVICE, cli, &commands.RUN)

func main() {
	cmd, err := uhppoted.Parse(cli, &commands.RUN, help)
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Printf("\nERROR: %v\n\n", err)
		os.Exit(1)
	}
}
