package feed

import (
	"encoding/xml"
	"fmt"
	"log"
	"rssReader/src/cli"
	"rssReader/src/configuration"
	"strings"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

func Read(config *configuration.Config) {
	for confID := range *config.GetFeeds() {
		configFeed := config.GetFeedAt(confID)

		if cli.IsVerboseDebug() {
			fmt.Println(fmt.Sprintf("Reading channel: `%s`", configFeed.Url))
			fmt.Println(fmt.Sprintf("Last checked item ID: %d", configFeed.MaxChecked))
		}
		allFeed := getRSSFeed(configFeed.Url)
		allFeed.filter(configFeed)
		if cli.IsVerboseInfo() {
			fmt.Println(fmt.Sprintf("Found %d new entries for channel `%s`", len(allFeed.Channel.Items), configFeed.Url))
		}
		allFeed.ListAll()
	}
}

func getRSSFeed(channelUrl string) *Rss {
	if cli.IsVerboseDebug() {
		fmt.Println("Reading URL " + channelUrl)
	}

	xmlBytes, err := GetRssReader(*cli.Downloader).GetXML(channelUrl)
	if err != nil {
		log.Fatalln(fmt.Sprintf("Failed to get XML at %v: %v", channelUrl, err.Error()))
	}

	var feed Rss
	err2 := xml.Unmarshal(xmlBytes, &feed)

	if err2 != nil {
		log.Fatalln(fmt.Sprintf("Error parsing: %v", err2))
	}

	return &feed
}

func (rss *Rss) filter(feedSource *configuration.FeedSource) {
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
