package main

import (
	"fmt"
	"os"

	"github.com/uhppoted/uhppote-core/uhppote"
	"github.com/uhppoted/uhppoted-api/command"
	"github.com/uhppoted/uhppoted-rest/commands"
)

var cli = []uhppoted.CommandV{
	&commands.RUN,
	commands.NewDaemonize(),
	commands.NewUndaemonize(),
	&uhppoted.VersionV{
		Application: commands.SERVICE,
		Version:     uhppote.VERSION,
	},
}

var help = uhppoted.NewHelpV(commands.SERVICE, cli, &commands.RUN)

func main() {
	cmd, err := uhppoted.ParseV(cli, &commands.RUN, help)
	if err != nil {
		fmt.Printf("\nError parsing command line: %v\n\n", err)
		os.Exit(1)
	}

	if err = cmd.Execute(); err != nil {
		fmt.Printf("\nERROR: %v\n\n", err)
		os.Exit(1)
	}
}
