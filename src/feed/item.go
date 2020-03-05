package feed

import (
	"encoding/xml"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	XMLName   xml.Name `xml:"item"`
	Title     string   `xml:"title"`
	Guid      string   `xml:"guid"`
	Link      string   `xml:"link"`
	Created   string   `xml:"pubDate"`
	processed string
}

func (item *Item) Identify() string {
	created, _ := time.Parse(time.RFC1123Z, item.Created)
	itemIdentifier := item.processed
	if "" == itemIdentifier {
		itemIdentifier = item.GetLink()
	}
	return fmt.Sprintf("[%s] %s ---> %s", created.Format("2006-01-02 15:04:05"), itemIdentifier, item.Title)
}

func (item *Item) HasMatch(searches *[]string) bool {
	for _, name := range *searches {
		if strings.Contains(item.Title, name) {
			return true
		}
	}
	return false
}

func (item *Item) GetID() (int, error) {
	numberPart := regexp.MustCompile(`(\d+)`).FindStringSubmatch(item.Guid)
	return strconv.Atoi(numberPart[0])
}

func (item *Item) GetLink() string {
	if item.Link == "" {
		return item.Guid
	}
	return item.Link
}

func (item *Item) ApplyPostProcessRegex(r *regexp.Regexp) {
	fakeFeed := FeedSource{
		Url:         item.Guid,
		IsProtected: false,
	}
	if fullContent, er := GetURLReader().GetContent(&fakeFeed); er == nil {
		test := string(fullContent)
		postProcessed := r.FindAllString(test, -1)
		if len(postProcessed) > 0 {
			uniq := map[string]bool{}
			for pp := range postProcessed {
				uniq[postProcessed[pp]] = true
			}
			postProcessed = []string{}
			for uniqItem, _ := range uniq {
				postProcessed = append(postProcessed, uniqItem)
			}
			item.processed = strings.Join(postProcessed, " | ")
		}
	}
}
