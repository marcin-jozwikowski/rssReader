package main

import (
	"cli"
	"configuration"
	"feed"
	"flag"
	"fmt"
)

func main() {
	configFileName := flag.String("configFile", "config.json", "Config file location")
	cacheFileName := flag.String("cacheFile", "cache.json", "Cache file location")
	runEditor := flag.Bool("editConfig", false, "Run configuration editor")
	isVerbose := flag.Int("verbose", cli.DefaultVerbose, "Verbose level: 0-None ... 3-All")
	flag.Parse()

	cli.SetVerbose(*isVerbose)
	config, configErr := configuration.ReadFromFile(*configFileName)
	cache, cacheErr := configuration.ReadFromFile(*cacheFileName)
	if configErr != nil {
		if cli.IsVerbose() {
			fmt.Println(configErr.Error())
		}
		*runEditor = true // enforce config editor
	}
	if cacheErr != nil {
		if cli.IsVerbose() {
			fmt.Println(cacheErr.Error())
		}
	}

	if *runEditor {
		config.Edit()
		_ = config.WriteToFile(*configFileName)
		return
	}

	cache = feed.Read(config, cache)
	_ = cache.WriteToFile(*cacheFileName)
}
