package main

import (
	"configuration"
	"feed"
	"flag"
	"fmt"
)

func main() {
	configFileName := flag.String("config", "config.json", "Config file location")
	cacheFileName := flag.String("cacheFile", "cache.json", "Cache file location")
	runEditor := flag.Bool("runConfig", false, "Run configuration editor")
	flag.Parse()

	config, configErr := configuration.ReadFromFile(*configFileName)
	cache, cacheErr := configuration.ReadFromFile(*cacheFileName)
	if configErr != nil {
		fmt.Println(configErr.Error())
		*runEditor = true // enforce config editor
	}
	if cacheErr != nil {
		fmt.Println(cacheErr.Error())
	}

	if *runEditor {
		config.Edit()
		_ = config.WriteToFile(*configFileName)
		return
	}

	cache = feed.Read(config, cache)
	_ = cache.WriteToFile(*cacheFileName)
}
