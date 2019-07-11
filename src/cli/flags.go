package cli

import (
	"flag"
)

var ConfigFileName *string
var CacheFileName *string
var RunEditor *bool
var Verbose *int
var Downloader *string
var DownloaderParams *string

func init() {
	ConfigFileName = flag.String("configFile", "config.json", "Config file location")
	CacheFileName = flag.String("cacheFile", "cache.json", "Cache file location")
	RunEditor = flag.Bool("editConfig", false, "Run configuration editor")
	Verbose = flag.Int("verbose", DefaultVerbose, "Verbose level: 0-None ... 3-All")
	Downloader = flag.String("downloader", "surf", "Download command to run. See online documentation for details.")
	DownloaderParams = flag.String("downloaderParams", "", "Download command params... See online documentation for details.")
	flag.Parse()

	SetVerbose(*Verbose)
}
