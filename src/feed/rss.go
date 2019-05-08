package feed

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Rss struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

func getXML(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return []byte{}, fmt.Errorf("GET error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return []byte{}, fmt.Errorf("Status error: %v", resp.StatusCode)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, fmt.Errorf("Read body: %v", err)
	}

	return data, nil
}

func GetRSSFeed(categoryName string) Rss {
	url := fmt.Sprintf("https://scnlog.me/%v/feed/", categoryName)
	xmlBytes, err := getXML(url)
	if err != nil {
		log.Printf("Failed to get XML at %v: %v", url, err)
	}

	var feed Rss
	err2 := xml.Unmarshal(xmlBytes, &feed)

	if err2 != nil {
		log.Printf("Error parsing: %v", err2)
	}

	return feed
}

func (rss *Rss) Filter(values []string) []Item {
	var result []Item
	for itemID := 0; itemID < len(rss.Channel.Items); itemID++ {
		for testValueID := 0; testValueID < len(values); testValueID++ {
			testItem := rss.Channel.Items[itemID]
			if testItem.Matches(values[testValueID]) {
				result = append(result, testItem)
			}
		}
	}

	return result
}


