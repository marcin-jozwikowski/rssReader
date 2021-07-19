package cui

import (
	"fmt"
	cui "github.com/jroimartin/gocui"
	"log"
)

type ListView struct {
	CliView
	Title      string
	items      *[]ListViewItem
	DrawItems  func()
	Dimensions ViewDimensions
}

type ListViewItem interface {
	ToString() string
}

func (listView *ListView) Init(gui *cui.Gui, name string, items *[]ListViewItem, title string, dimensions ViewDimensions) {
	listView.items = items
	listView.gui = gui
	listView.Dimensions = dimensions

	listView.view, _ = listView.gui.SetView(name, 0, 0, 1, 1)
	listView.view.Clear()
	listView.SetTitle(title)
}

func (listView *ListView) SetTitle(title string)  {
	listView.view.Title = " " + title + " "
}

func (listView *ListView) Draw() {
	var err error
	listView.view, err = listView.gui.SetView(listView.view.Name(), listView.Dimensions.left, listView.Dimensions.top, listView.Dimensions.right, listView.Dimensions.bottom)
	if err != nil {
		log.Fatal("Cannot update sites view", err)
	}
	listView.view.Clear()

	if listView.DrawItems != nil {
		listView.DrawItems()
	} else {
		for _, item := range *listView.items {
			_, _ = fmt.Fprintln(listView.view, item.ToString())
		}
	}
}

func (listView *ListView) AddItem(item ListViewItem) {
	*listView.items = append(*listView.items, item)
	listView.Draw()
}

func (listView *ListView) ResetItems() {
	*listView.items = []ListViewItem{}
	listView.Draw()
}

func (listView *ListView) Update() {
	listView.gui.Update(func(gui *cui.Gui) error {
		return nil
	})
}

func (view *CliView) Focus() error {
	view.view.Highlight = true
	if _, err := view.gui.SetCurrentView(view.view.Name()); err != nil {
		panic(err)
		return nil
	}
	return nil
}
