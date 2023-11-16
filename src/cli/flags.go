package cli

import (
	"flag"
)

var ConfigFileName *string
var RunEditor *bool
var Verbose *int
var OneTimeUrl *string

func init() {
	ConfigFileName = flag.String("configFile", GetConfigFileLocation("sources.json"), "Config file location")
	RunEditor = flag.Bool("editConfig", false, "Run configuration editor")
	Verbose = flag.Int("verbose", DefaultVerbose, "Verbose level: 0-None ... 3-All")
	OneTimeUrl = flag.String("url", "", "Url for the one-time config")
	flag.Parse()

	SetVerbose(*Verbose)
}
