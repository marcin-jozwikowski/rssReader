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
	flag.Parse()

	config, err := configuration.ReadFromFile(*configFileName, getDefaultConfig(configType_main))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	cache, err := configuration.ReadFromFile(*cacheFileName, getDefaultConfig(configType_cache))
	if err != nil {
		fmt.Println(err.Error())
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
