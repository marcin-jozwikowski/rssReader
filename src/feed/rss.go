package feed

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"log"
	"regexp"
	"rssReader/src/cli"
	"strings"
	"sync"
	"github.com/PuerkitoBio/goquery"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type ResultItem struct {
	Item       Item
	FeedSource *FeedSource
}

func (item *ResultItem) Identify() string {
	if len(item.FeedSource.PostProcess) > 0 {
		item.ApplyPostProcess()
	}
	return item.Item.Identify()
}

func (item *ResultItem) ApplyPostProcess() {
	r := regexp.MustCompile(item.FeedSource.PostProcess)
	item.Item.ApplyPostProcessRegex(r)
}

func Read(config *Config) {
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
		fmt.Println(returnItem.Identify())
	}
}

func readOneFeed(configFeed *FeedSource, waitGroup *sync.WaitGroup, results chan ResultItem) {
	if cli.IsVerboseDebug() {
		fmt.Println(fmt.Sprintf("Reading channel: `%s`", configFeed.Url))
		fmt.Println(fmt.Sprintf("Last checked item ID: %d", configFeed.MaxChecked))
	}
	var allFeed *Rss
	if configFeed.IsHTML {
		allFeed = getHTMLFeed(configFeed)
	} else {
		allFeed = getRSSFeed(configFeed)
	}
	if allFeed != nil {
		allFeed.filterOut(configFeed)
		if cli.IsVerboseInfo() {
			fmt.Println(fmt.Sprintf("Found %d new entries for channel `%s`", len(allFeed.Channel.Items), configFeed.Url))
		}
		allFeed.WriteAllItemsToChannel(results, configFeed)
	}

	waitGroup.Done()
}

func getHTMLFeed(configFeed *FeedSource) *Rss {
	var feed = Rss{}
	var urlReader = GetURLReader()
	readItems := 0
	page := 1
	for {
		var xmlBytes []byte
		var err error

		if configFeed.IsPaginated {
			xmlBytes, err = urlReader.GetContentPaginated(configFeed, page)
		} else {
			xmlBytes, err = urlReader.GetContent(configFeed)
		}
		if err == nil {
			if readerPage, documentError := goquery.NewDocumentFromReader(bytes.NewReader(xmlBytes)); documentError == nil {
				readerPage.Find(".post").Each(func(id int, selection *goquery.Selection) {
					link := selection.Find("div.title > h1 > a")
					class, _ := selection.Find("div.title").Attr("class")
					href, _ := link.Attr("href")
					created, _ := selection.Find("div.title > small > span.localtime").Attr("data-lttime")
					item := Item{
						Title:   link.Text(),
						Link:    href,
						Guid:    class,
						Created: created,
					}
					feed.Channel.Items = append(feed.Channel.Items, item)
					readItems++
				})
			} else {
				log.Println(fmt.Sprintf("Failed to document at %v: %v", configFeed.Url, documentError.Error()))
				break
			}
		} else {
			log.Println(fmt.Sprintf("Failed to get HTML at %v: %v", configFeed.Url, err.Error()))
			break
		}

		if readItems > 100 || page > 20 {
			break
		}
		page++
	}

	return &feed
}

func getRSSFeed(configFeed *FeedSource) *Rss {
	if cli.IsVerboseDebug() {
		fmt.Println("Reading URL " + configFeed.Url)
	}

	var feed Rss
	if xmlBytes, err := GetURLReader().GetContent(configFeed); err == nil {
		if err2 := xml.Unmarshal(xmlBytes, &feed); err2 == nil {
			return &feed
		} else {
			log.Println(fmt.Sprintf("Error parsing URL %v: %v", configFeed.Url, err2))
		}
	} else {
		log.Println(fmt.Sprintf("Failed to get XML at %v: %v", configFeed.Url, err.Error()))
	}

	return nil
}

func (rss *Rss) filterOut(feedSource *FeedSource) {
	if cli.IsVerboseDebug() {
		fmt.Println("Checking against: " + strings.Join(feedSource.SearchPhrases, " | "))
	}
	maxChecked := feedSource.MaxChecked
	for testItemPosition := 0; testItemPosition < rss.Channel.GetItemsCount(); {
		// for each found RSS item:
		testItem := rss.Channel.GetItemAt(testItemPosition)
		testItemID, _ := testItem.GetID()
		if testItemID > feedSource.MaxChecked {
			if testItemID > maxChecked {
				maxChecked = testItemID
			}
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
		} else {
			if cli.IsVerboseDebug() {
				fmt.Println(fmt.Sprintf("Item ID of %d is not newer than max %d", testItemID, feedSource.MaxChecked))
			}
			// if this item has been already checked only those preceeding it can be used
			rss.Channel.DropItemAtPosition(testItemPosition)
			break
		}
	}
	feedSource.SetMaxChecked(maxChecked)
}

func (rss *Rss) WriteAllItemsToChannel(results chan ResultItem, configFeed *FeedSource) {
	for _, item := range rss.Channel.GetAllItems() {
		results <- ResultItem{Item: item, FeedSource: configFeed}
	}
}
