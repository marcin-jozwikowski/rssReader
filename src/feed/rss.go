package feed

import (
	"configuration"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

func Read(config configuration.Config, cache configuration.Config) configuration.Config {
	for channelName, filterValues := range config {
		if configuration.IsVerboseDebug() {
			fmt.Println(fmt.Sprintf("Reading channel: `%s`", channelName))
		}
		var channelMaxID int
		if cache[channelName] != nil {
			channelMaxID, _ = strconv.Atoi(cache[channelName][0])
		} else {
			channelMaxID = 0
		}
		if configuration.IsVerboseDebug() {
			fmt.Println(fmt.Sprintf("Last checked item ID: %d", channelMaxID))
		}
		allFeed := getRSSFeed(channelName)
		matching, channelMaxID := allFeed.filter(filterValues, channelMaxID)
		if configuration.IsVerboseInfo() {
			fmt.Println(fmt.Sprintf("Found %d new entries for channel `%s`", len(matching), channelName))
		}
		for matchID := 0; matchID < len(matching); matchID++ {
			if configuration.IsVerbose() {
				fmt.Println(matching[matchID].Identify())
			}
		}
		cache[channelName] = []string{strconv.Itoa(channelMaxID)}
	}

	return cache
}

func getXML(url string) ([]byte, error) {
	if configuration.IsVerboseDebug() {
		fmt.Println("Reading URL " + url)
	}
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err.Error())
	}

	return data, nil
}

func getRSSFeed(categoryName string) Rss {
	url := fmt.Sprintf("https://scnlog.me/%v/feed/", categoryName)
	xmlBytes, err := getXML(url)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to get XML at %v: %v", url, err.Error()))
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
	for itemID := 0; itemID < len(rss.Channel.Items); itemID++ {
		testItem := rss.Channel.Items[itemID]
		testItemID, _ := testItem.GetID()
		if maxID < testItemID {
			for testValueID := 0; testValueID < len(values); testValueID++ {
				if testItem.Matches(values[testValueID]) {
					if configuration.IsVerboseDebug() {
						fmt.Println(fmt.Sprintf("Item by ID %s mathed by %s", testItem.Title, values[testValueID]))
					}
					result = append(result, testItem)
					break
				}
			}
			if testItemID > newMaxID {
				newMaxID = testItemID
			}
		} else {
			if configuration.IsVerboseDebug() {
				fmt.Println(fmt.Sprintf("Item ID of %d is not newer than max %d", testItemID, maxID))
			}
			break
		}
	}

	return result, newMaxID
}
