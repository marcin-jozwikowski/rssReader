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

func Read(config *configuration.Config) {
	var waitGroup sync.WaitGroup

	for confID := range *config.GetFeeds() {
		waitGroup.Add(1)
		go readOneFeed(config.GetFeedAt(confID), &waitGroup)
	}

	waitGroup.Wait()
}

func readOneFeed(configFeed *configuration.FeedSource, waitGroup *sync.WaitGroup)  {
	if cli.IsVerboseDebug() {
		fmt.Println(fmt.Sprintf("Reading channel: `%s`", configFeed.Url))
		fmt.Println(fmt.Sprintf("Last checked item ID: %d", configFeed.MaxChecked))
	}
	if allFeed := getRSSFeed(configFeed.Url); allFeed != nil {
		allFeed.filterOut(configFeed)
		if cli.IsVerboseInfo() {
			fmt.Println(fmt.Sprintf("Found %d new entries for channel `%s`", len(allFeed.Channel.Items), configFeed.Url))
		}
		allFeed.ListAll()
	}

	waitGroup.Done()
}

func getRSSFeed(channelUrl string) *Rss {
	if cli.IsVerboseDebug() {
		fmt.Println("Reading URL " + channelUrl)
	}

	var feed Rss
	if xmlBytes, err := GetRssReader(*cli.Downloader).GetXML(channelUrl); err == nil {
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
	for testItemPosition := 0; testItemPosition < len(rss.Channel.Items); {
		// for each found RSS item:
		testItem := &rss.Channel.Items[testItemPosition]
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
			rss.Channel.Items = rss.Channel.Items[:testItemPosition]
			break
		}
	}
}

func (rss *Rss) ListAll() {
	for _, match := range rss.Channel.Items {
		if cli.IsVerbose() {
			fmt.Println(match.Identify())
		}
	}
}
