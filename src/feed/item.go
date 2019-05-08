package feed

import (
	"encoding/xml"
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

func (item *Item) Matches(name string) bool {
	return strings.Contains(item.Title, name)
}