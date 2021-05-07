package reader

import (
	"fmt"
	cui "github.com/jroimartin/gocui"
	"log"
	listCui "rssReader/src/cui"
	"strconv"
)

const ViewsSources = "viewSources"

var allViews = make(map[string]*listCui.ListView)
var runtimeConfig *RuntimeConfig

func RunCUI(config *RuntimeConfig) {
	runtimeConfig = config
	createCUI()
}

func runReader() {
	for rcId := range runtimeConfig.Sources {
		go func(source *DataSource) {
			source.AddResultingShow(runForDataSource(source))
			allViews[ViewsSources].Update()
		}(&runtimeConfig.Sources[rcId])
	}
}

func createCUI() bool {
	gui, err := cui.NewGui(cui.OutputNormal)
	if err != nil {
		log.Fatal("Failed to build CUI", err)
	}
	defer gui.Close()

	gui.SelFgColor = cui.ColorGreen | cui.AttrBold
	gui.BgColor = cui.ColorDefault
	gui.Highlight = true
	gui.SetManagerFunc(layoutManager)

	v := new(listCui.ListView)
	v.Init(gui, ViewsSources, runtimeConfig.SourcesAsList(), "Sources", getViewDimensions(gui, ViewsSources))
	allViews[ViewsSources] = v
	v.DrawItems = viewSourcesDrawItems

	//v = new(ListView)
	//v.Init(gui, ViewsFeedDetails, []string{}, "Select source to view details")
	//v.DrawItems = viewFeedDetailsDrawItems
	//
	//v = new(ListView)
	//v.Init(gui, ViewsFeedResults, []string{}, "Results")
	//v.view.Wrap = true

	_ = allViews[ViewsSources].Focus()

	initAllKeyBindings(gui)

	runReader()

	finalError := gui.MainLoop()
	switch finalError {
	case nil:
	case cui.ErrQuit:
		return false
	}

	return true
}

func viewSourcesDrawItems() {
	for _, source := range runtimeConfig.Sources {
		if source.show == nil {
			_, _ = fmt.Fprintln(allViews[ViewsSources].GetView(), source.Name)
		} else {
			_, _ = fmt.Fprintln(allViews[ViewsSources].GetView(), source.Name+" | "+strconv.Itoa(len(source.show.Episodes))+" Episodes")
		}
	}
}

func initAllKeyBindings(gui *cui.Gui) {
	if err := gui.SetKeybinding("", cui.KeyCtrlQ, cui.ModNone, listCui.ExitRightNow); err != nil {
		log.Fatal("Failed to set keybindings")
	}
}

func layoutManager(gui *cui.Gui) error {
	for _, listView := range allViews {
		listView.Draw()
	}
	return nil
}

func getViewDimensions(gui *cui.Gui, viewName string) listCui.ViewDimensions {
	viewWidth, viewHeight := gui.Size()
	columnWidth := 2 * viewWidth / 3
	columnHeight := viewHeight / 2

	switch viewName {
	case ViewsSources:
		return listCui.NewViewDimensions(0, 0, columnWidth-1, columnHeight-1)
		//case ViewsFeedResults:
		//	return 0, columnHeight, viewWidth - 1, viewHeight - 1
		//case ViewsFeedDetails:
		//	return columnWidth, 0, viewWidth - 1, columnHeight - 1
	}
	return listCui.NewViewDimensions(0, 0, 0, 0)
}
