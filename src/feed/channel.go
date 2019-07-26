package feed

import "encoding/xml"

type Channel struct {
	XMLName xml.Name `xml:"channel"`
	Items   []Item   `xml:"item"`
}

func (channel *Channel) Identify() string {
	return channel.XMLName.Local
}

func (channel *Channel) RemoveItemAt(i int) {
	if 1 == len(channel.Items) && 0 == i {
		channel.Items = []Item{}
	} else {
		channel.Items = append(channel.Items[:i], channel.Items[i+1:]...)
	}
}

func (channel *Channel) RemoveItem(item Item) bool {
	idToRemove := -1

	for i := range channel.Items {
		if channel.Items[i] == item {
			idToRemove = i
			break
		}
	}

	if -1 != idToRemove {
		channel.RemoveItemAt(idToRemove)
		return true
	}
	return false
}

func (channel *Channel) GetAllItems() []Item {
	return channel.Items
}

func (channel *Channel) DropItemAtPosition(itemPosition int) {
	channel.Items = channel.Items[:itemPosition]
}

func (channel *Channel) GetItemsCount() int {
	return len(channel.Items)
}

func (channel *Channel) GetItemAt(testItemPosition int) *Item {
	return &channel.Items[testItemPosition]
}
