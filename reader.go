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
		var feedRss = feed.GetRSSFeed("https://scnlog.me/" + config.Settings[settingID].Key + "/feed/")
		for i := 0; i < len(feedRss.Channel.Items); i++ {
			fmt.Println(feedRss.Channel.Items[i].Identify())
		}
	}
}
