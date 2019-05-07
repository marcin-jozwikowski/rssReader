package feed

import "encoding/xml"

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

func (channel *Channel) Identify() string {
	return channel.XMLName.Local
}
