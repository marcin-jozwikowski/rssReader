package main

import (
	"configuration"
	"feed"
	"flag"
	"fmt"
)

const configType_main = "main"
const configType_cache = "cache"

func main() {
	configFileName := flag.String("config", "config.json", "Config file location")
	cacheFileName := flag.String("cacheFile", "cache.json", "Cache file location")
	runEditor := flag.Bool("runConfig", false, "Run configuration editor")
	flag.Parse()

	config, configErr := configuration.ReadFromFile(*configFileName, getDefaultConfig(configType_main))
	cache, cacheErr := configuration.ReadFromFile(*cacheFileName, getDefaultConfig(configType_cache))
	if configErr != nil || cacheErr != nil {
		if configErr != nil {
			fmt.Println(configErr.Error())
		}
		if cacheErr != nil {
			fmt.Println(cacheErr.Error())
		}
		return
	}

	if *runEditor {
		config.Edit()
		_ = config.WriteToFile(*configFileName)
		return
	}

	cache = feed.Read(config, cache)
	_ = cache.WriteToFile(*cacheFileName)
}

func getDefaultConfig(configType string) configuration.Config {
	var conf = make(configuration.Config,1)
	switch configType {
	case configType_main:
		conf["tv-shows"] = []string{"ALF", "That 70's Show"}
		break
	}
	return conf
}
