package feed

import (
	"cli"
	"configuration"
	"encoding/xml"
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

func Read(config configuration.Config, cache configuration.Config) configuration.Config {
	for channelUrl, filterValues := range config {
		if cli.IsVerboseDebug() {
			fmt.Println(fmt.Sprintf("Reading channel: `%s`", channelUrl))
		}
		var channelMaxID int
		if cache[channelUrl] != nil {
			channelMaxID, _ = strconv.Atoi(cache[channelUrl][0])
		} else {
			channelMaxID = 0
		}
		if cli.IsVerboseDebug() {
			fmt.Println(fmt.Sprintf("Last checked item ID: %d", channelMaxID))
		}
		allFeed := getRSSFeed(channelUrl)
		matching, channelMaxID := allFeed.filter(filterValues, channelMaxID)
		if cli.IsVerboseInfo() {
			fmt.Println(fmt.Sprintf("Found %d new entries for channel `%s`", len(matching), channelUrl))
		}
		for matchID := 0; matchID < len(matching); matchID++ {
			if cli.IsVerbose() {
				fmt.Println(matching[matchID].Identify())
			}
		}
		cache[channelUrl] = []string{strconv.Itoa(channelMaxID)}
	}

	return cache
}

func getXML(url string) ([]byte, error) {
	if cli.IsVerboseDebug() {
		fmt.Println("Reading URL " + url)
	}

	reader := GetRssReader(*cli.DownloadCommand)
	return reader.GetXML(url)
}

func getRSSFeed(channelUrl string) Rss {
	xmlBytes, err := getXML(channelUrl)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to get XML at %v: %v", channelUrl, err.Error()))
	}

	var feed Rss
	err2 := xml.Unmarshal(xmlBytes, &feed)

	if err2 != nil {
		log.Fatalln(fmt.Sprintf("Error parsing: %v", err2))
	}

	return feed
}

func (rss *Rss) filter(values []string, maxID int) ([]Item, int) {
	var result []Item
	newMaxID := maxID
	if cli.IsVerboseDebug() {
		fmt.Println("Checking against: " + strings.Join(values, " | "))
	}
	for itemID := 0; itemID < len(rss.Channel.Items); itemID++ {
		// for each found RSS item:
		testItem := rss.Channel.Items[itemID]
		testItemID, _ := testItem.GetID()
		if maxID < testItemID {
			// if current itemID is greater than last checked
			if cli.IsVerboseDebug() {
				fmt.Println("Checking item: " + testItem.Title)
			}
			for testValueID := 0; testValueID < len(values); testValueID++ {
			// test item against all keywords
				if testItem.Matches(values[testValueID]) {
					if cli.IsVerboseInfo() {
						fmt.Println(fmt.Sprintf("Item %s mathed by %s", testItem.Title, values[testValueID]))
					}
					result = append(result, testItem)
					break
				}
			}
			if testItemID > newMaxID {
				newMaxID = testItemID
			}
		} else {
			if cli.IsVerboseDebug() {
				fmt.Println(fmt.Sprintf("Item ID of %d is not newer than max %d", testItemID, maxID))
			}
			break
		}
	}

	return result, newMaxID
}
