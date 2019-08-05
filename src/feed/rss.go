package feed

import (
	"encoding/xml"
	"fmt"
	"log"
	"rssReader/src/cli"
	"rssReader/src/configuration"
	"strings"
	"sync"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type ResultItem struct {
	Item       Item
	FeedSource *configuration.FeedSource
}

func Read(config *configuration.Config) {
	var waitGroup sync.WaitGroup
	results := make(chan ResultItem, config.CountFeeds()*20)

	for confID := range *config.GetFeeds() {
		waitGroup.Add(1)
		go readOneFeed(config.GetFeedAt(confID), &waitGroup, results)
	}

	waitGroup.Wait()
	close(results)

	for {
		returnItem, hasMore := <-results
		if !hasMore {
			break
		}
		fmt.Println(returnItem.Item.Identify())
	}
}

func readOneFeed(configFeed *configuration.FeedSource, waitGroup *sync.WaitGroup, results chan ResultItem) {
	if cli.IsVerboseDebug() {
		fmt.Println(fmt.Sprintf("Reading channel: `%s`", configFeed.Url))
		fmt.Println(fmt.Sprintf("Last checked item ID: %d", configFeed.MaxChecked))
	}
	if allFeed := getRSSFeed(configFeed.Url); allFeed != nil {
		allFeed.filterOut(configFeed)
		if cli.IsVerboseInfo() {
			fmt.Println(fmt.Sprintf("Found %d new entries for channel `%s`", len(allFeed.Channel.Items), configFeed.Url))
		}
		allFeed.WriteAllItemsToChannel(results, configFeed)
	}

	waitGroup.Done()
}

func getRSSFeed(channelUrl string) *Rss {
	if cli.IsVerboseDebug() {
		fmt.Println("Reading URL " + channelUrl)
	}

	var feed Rss
	if xmlBytes, err := GetURLReader().GetContent(channelUrl); err == nil {
		if err2 := xml.Unmarshal(xmlBytes, &feed); err2 == nil {
			return &feed
		} else {
			log.Println(fmt.Sprintf("Error parsing URL %v: %v", channelUrl, err2))
		}
	} else {
		log.Println(fmt.Sprintf("Failed to get XML at %v: %v", channelUrl, err.Error()))
	}

	return nil
}

func (rss *Rss) filterOut(feedSource *configuration.FeedSource) {
	if cli.IsVerboseDebug() {
		fmt.Println("Checking against: " + strings.Join(feedSource.SearchPhrases, " | "))
	}
	maxChecked := feedSource.MaxChecked
	for testItemPosition := 0; testItemPosition < rss.Channel.GetItemsCount(); {
		// for each found RSS item:
		testItem := rss.Channel.GetItemAt(testItemPosition)
		testItemID, _ := testItem.GetID()
		if maxChecked < testItemID {
			// if current itemID is greater than last checked
			if cli.IsVerboseDebug() {
				fmt.Println("Checking item: " + testItem.Title)
			}
			if testItem.HasMatch(&feedSource.SearchPhrases) {
				if cli.IsVerboseInfo() {
					fmt.Println(fmt.Sprintf("Item %s has a math", testItem.Title))
				}
				testItemPosition++
				continue
			} else {
				rss.Channel.RemoveItem(*testItem)
			}
			feedSource.SetMaxChecked(testItemID)
		} else {
			if cli.IsVerboseDebug() {
				fmt.Println(fmt.Sprintf("Item ID of %d is not newer than max %d", testItemID, feedSource.MaxChecked))
			}
			// if this item has been already checked only those preceeding it can be used
			rss.Channel.DropItemAtPosition(testItemPosition)
			break
		}
	}
}

func (rss *Rss) WriteAllItemsToChannel(results chan ResultItem, configFeed *configuration.FeedSource) {
	for _, item := range rss.Channel.GetAllItems() {
		results <- ResultItem{Item: item, FeedSource: configFeed}
	}
}
