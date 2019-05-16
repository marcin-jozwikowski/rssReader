package main

import (
	"configuration"
	"feed"
	"flag"
	"fmt"
	"strconv"
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

	for channelName, filterValues := range config {
		var channelMaxID int
		if cache[channelName] != nil {
			channelMaxID, _ = strconv.Atoi(cache[channelName][0])
		} else {
			channelMaxID = 0
		}
		allFeed := feed.GetRSSFeed(channelName)
		matching, channelMaxID := allFeed.Filter(filterValues, channelMaxID)
		for matchID := 0; matchID < len(matching); matchID++ {
			fmt.Println(matching[matchID].Identify())
		}
		cache[channelName] = []string{strconv.Itoa(channelMaxID)}
	}

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
