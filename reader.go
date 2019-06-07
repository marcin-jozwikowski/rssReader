package main

import (
	"cli"
	"configuration"
	"feed"
	"fmt"
)

func main() {
	config, configErr := configuration.ReadFromFile(*cli.ConfigFileName)
	cache, cacheErr := configuration.ReadFromFile(*cli.CacheFileName)
	if configErr != nil {
		if cli.IsVerbose() {
			fmt.Println(configErr.Error())
		}
		*cli.RunEditor = true // enforce config editor
	}
	if cacheErr != nil {
		if cli.IsVerbose() {
			fmt.Println(cacheErr.Error())
		}
	}

	if *cli.RunEditor {
		config.Edit()
		_ = config.WriteToFile(*cli.ConfigFileName)
		return
	}

	cache = feed.Read(config, cache)
	_ = cache.WriteToFile(*cli.CacheFileName)
}
