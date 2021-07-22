package cli

import (
	"flag"
	"os"
)

var ConfigFileName *string
var RunEditor *bool
var Verbose *int

func init() {
	dirname, _ := os.UserHomeDir()
	filename := "sources.json"
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		filename = dirname + "/.tvReader/" + filename
	}
	ConfigFileName = flag.String("configFile", filename, "Config file location")
	RunEditor = flag.Bool("editConfig", false, "Run configuration editor")
	Verbose = flag.Int("verbose", DefaultVerbose, "Verbose level: 0-None ... 3-All")
	flag.Parse()

	SetVerbose(*Verbose)
}
