package main

import (
	"configuration"
	"feed"
	"fmt"
)

func main() {
	var config, err = configuration.ReadFromFile("config.json")
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
}
