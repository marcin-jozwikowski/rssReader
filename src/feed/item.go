package feed

import (
	"encoding/xml"
	"regexp"
	"strconv"
	"strings"
)

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Guid    string   `xml:"guid"`
}

func (item *Item) Identify() string {
	return item.Guid + " ---> " + item.Title
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
