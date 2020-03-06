package cli

import (
	"flag"
)

var ConfigFileName *string
var RunEditor *bool
var Verbose *int
var ResetChecked *bool
var PageReadLimit *int

func init() {
	ConfigFileName = flag.String("configFile", "config.json", "Config file location")
	RunEditor = flag.Bool("editConfig", false, "Run configuration editor")
	Verbose = flag.Int("verbose", DefaultVerbose, "Verbose level: 0-None ... 3-All")
	ResetChecked = flag.Bool("resetChecked", false, "Reset last checked counters")
	PageReadLimit = flag.Int("pageReadLimit", 20, "Maximum pages to read")
	flag.Parse()

	SetVerbose(*Verbose)
}
