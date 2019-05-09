package main

import (
	"configuration"
	"feed"
	"flag"
	"fmt"
)

func main() {
	configFileName := flag.String("config", "config.json", "Config file location")
	flag.Parse()

	config, err := configuration.ReadFromFile(*configFileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for settingID := 0; settingID < len(config.Settings); settingID++ {
		setting := config.Settings[settingID]
		allFeed := feed.GetRSSFeed(setting.Key)
		matching := allFeed.Filter(setting.Values)
		for matchID := 0; matchID < len(matching); matchID++ {
			fmt.Println(matching[matchID].Identify())
		}
	}

	err = config.WriteToFile(*configFileName)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
