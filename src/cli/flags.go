package cli

import (
	"flag"
)

var ConfigFileName *string
var RunEditor *bool
var Verbose *int
var PageReadLimit *int

func init() {
	ConfigFileName = flag.String("configFile", "sources.json", "Config file location")
	RunEditor = flag.Bool("editConfig", false, "Run configuration editor")
	Verbose = flag.Int("verbose", DefaultVerbose, "Verbose level: 0-None ... 3-All")
	PageReadLimit = flag.Int("pageReadLimit", 15, "Maximum pages to read")
	flag.Parse()

	SetVerbose(*Verbose)
}
