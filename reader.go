package main

import (
	"feed"
	"fmt"
)

func main() {
	var feedRss = feed.GetRSSFeed("https://scnlog.me/tv-shows/sdtv/feed/")
	for i := 0; i < len(feedRss.Channel.Items); i++ {
		fmt.Println(feedRss.Channel.Items[i].Identify())
	}
}