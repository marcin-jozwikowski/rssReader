package cli

import ("flag")

var ConfigFileName *string
var CacheFileName *string
var RunEditor *bool
var Verbose *int

func init()  {
	ConfigFileName = flag.String("configFile", "config.json", "Config file location")
	CacheFileName = flag.String("cacheFile", "cache.json", "Cache file location")
	RunEditor = flag.Bool("editConfig", false, "Run configuration editor")
	Verbose = flag.Int("verbose", DefaultVerbose, "Verbose level: 0-None ... 3-All")
	flag.Parse()

	SetVerbose(*Verbose)
}
