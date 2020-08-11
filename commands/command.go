package commands

import (
	"flag"
	"fmt"
)

const SERVICE = `uhppoted-rest`

func helpOptions(flagset *flag.FlagSet) {
	flags := 0
	count := 0

	flag.VisitAll(func(f *flag.Flag) {
		count++
	})

	flagset.VisitAll(func(f *flag.Flag) {
		flags++
		fmt.Printf("    --%-13s %s\n", f.Name, f.Usage)
	})

	if count > 0 {
		fmt.Println()
		fmt.Println("  Options:")
		flag.VisitAll(func(f *flag.Flag) {
			fmt.Printf("    --%-13s %s\n", f.Name, f.Usage)
		})
	}

	if flags > 0 {
		fmt.Println()
	}
}
