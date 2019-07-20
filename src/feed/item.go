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
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Guid    string   `xml:"guid"`
	Created string   `xml:"pubDate"`
}

func (item *Item) Identify() string {
	created, _ := time.Parse(time.RFC1123Z, item.Created)
	return fmt.Sprintf("[%s] %s ---> %s", created.Format("2006-01-02 15:04:05"), item.Guid, item.Title)
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
