package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/kangaechu/richecker"
	"os"
)

type CmdOpts struct {
	Days int `short:"d" long:"days_before_expiration" description:"Print if expiration date is less than this days. (default:3)"`
}

func main() {
	var opts CmdOpts
	parser := flags.NewParser(&opts, flags.Default)
	_, err := flags.Parse(&opts)
	if err != nil {
		parser.WriteHelp(os.Stdout)
		os.Exit(1)
	}
	richecker.Check(opts.Days)
}
