package feed

import "encoding/xml"

type Item struct {
	XMLName xml.Name `xml:"item"`
	Title   string   `xml:"title"`
	Guid    string   `xml:"guid"`
}

func (item *Item) Identify() string {
	return item.Title + " ---> " + item.Guid
}