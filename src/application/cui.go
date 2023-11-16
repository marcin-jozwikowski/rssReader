package application

import (
	"log"
	listCui "rssReader/src/cui"
	"rssReader/src/publishing"

	cui "github.com/jroimartin/gocui"
)

const ViewsSources = "viewSources"
const ViewsPieces = "viewsPieces"
const ViewsReleases = "viewReleases"

var allViews = make(map[string]*listCui.ListView)
var runtimeConfig *RuntimeConfig
var currentSourceId = -1
var currentEntryId = -1

func RunCUI(config *RuntimeConfig) {
	runtimeConfig = config
	createCUI()
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
	v.Init(gui, ViewsSources, runtimeConfig.GetListViewItems(), "Sources", getViewDimensions(gui, ViewsSources))
	allViews[ViewsSources] = v

	v = new(listCui.ListView)
	v.Init(gui, ViewsPieces, &[]listCui.ListViewItem{}, "Select source to view its pieces", getViewDimensions(gui, ViewsPieces))
	allViews[ViewsPieces] = v

	v = new(listCui.ListView)
	v.Init(gui, ViewsReleases, &[]listCui.ListViewItem{}, "Select a piece to view its releases", getViewDimensions(gui, ViewsReleases))
	allViews[ViewsReleases] = v

	_ = allViews[ViewsSources].Focus()

	initAllKeyBindings(gui)

	finalError := gui.MainLoop()
	switch finalError {
	case nil:
	case cui.ErrQuit:
		return false
	}

	return true
}

func viewSourcesSelectEntry(_ *cui.Gui, view *cui.View) error {
	_, selectedSource := view.Cursor()
	_, offset := view.Origin()
	currentSourceId = selectedSource + offset
	if runtimeConfig.GetSourceAt(currentSourceId).GetResultingPublishing() == nil {
		// run the reader
		go func(source *DataSource) {
			source.AddResultingPublishing(RunForDataSource(source))
			allViews[ViewsSources].Update()
		}(runtimeConfig.GetSourceAt(currentSourceId))
		allViews[ViewsSources].Update()
	} else {
		// display the items
		allViews[ViewsPieces].GetView().Clear()
		allViews[ViewsPieces].SetItems(runtimeConfig.GetSourceAt(currentSourceId).GetListViewItems())
		_ = allViews[ViewsPieces].Focus()
	}
	return nil
}

func viewEntriesSelectEntry(_ *cui.Gui, view *cui.View) error {
	_, selectedEntry := view.Cursor()
	_, offset := view.Origin()
	selectedEntry += offset
	currentEntryId = selectedEntry
	allViews[ViewsReleases].GetView().Clear()
	allViews[ViewsReleases].SetItems(runtimeConfig.GetSourceAt(currentSourceId).GetResultingPublishing().GetPieceByAt(currentEntryId).GetListViewItems())
	_ = allViews[ViewsReleases].Focus()
	return nil
}

func initAllKeyBindings(gui *cui.Gui) {
	if err := gui.SetKeybinding("", cui.KeyCtrlQ, cui.ModNone, listCui.ExitRightNow); err != nil {
		log.Fatal("Failed to set keybindings")
	}

	if err := gui.SetKeybinding(ViewsSources, cui.KeyArrowDown, cui.ModNone, moveCursorDown); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsSources, cui.KeyArrowUp, cui.ModNone, moveCursorUp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsSources, cui.KeyEnter, cui.ModNone, viewSourcesSelectEntry); err != nil {
		log.Fatal("Failed to set keybindings")
	}

	if err := gui.SetKeybinding(ViewsPieces, cui.KeyArrowDown, cui.ModNone, moveCursorDown); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsPieces, cui.KeyArrowUp, cui.ModNone, moveCursorUp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsPieces, cui.KeyEnter, cui.ModNone, viewEntriesSelectEntry); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsPieces, cui.KeyBackspace, cui.ModNone, viewEntriesClose); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsPieces, cui.KeyBackspace2, cui.ModNone, viewEntriesClose); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsPieces, cui.KeyEsc, cui.ModNone, viewEntriesClose); err != nil {
		log.Fatal("Failed to set keybindings")
	}

	if err := gui.SetKeybinding(ViewsReleases, cui.KeyArrowDown, cui.ModNone, moveCursorDown); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsReleases, cui.KeyArrowUp, cui.ModNone, moveCursorUp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsReleases, cui.KeyBackspace, cui.ModNone, viewReleasesClose); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsReleases, cui.KeyBackspace2, cui.ModNone, viewReleasesClose); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsReleases, cui.KeyEsc, cui.ModNone, viewReleasesClose); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err := gui.SetKeybinding(ViewsReleases, cui.KeyEnter, cui.ModNone, viewReleasesSelectRelease); err != nil {
		log.Fatal("Failed to set keybindings")
	}
}

func viewReleasesSelectRelease(_ *cui.Gui, view *cui.View) error {
	_, selectedRelease := view.Cursor()
	_, offset := view.Origin()
	selectedRelease += offset

	rel := runtimeConfig.GetSourceAt(currentSourceId).GetResultingPublishing().GetPieceByAt(currentEntryId).GetReleaseAt(selectedRelease)

	go func(release *publishing.Release) {
		RunForRelease(runtimeConfig.GetSourceAt(currentSourceId), release)
		allViews[ViewsReleases].Update()
	}(rel)

	return nil
}

func viewEntriesClose(_ *cui.Gui, _ *cui.View) error {
	currentSourceId = -1
	_ = allViews[ViewsPieces].GetView().SetCursor(0, 0)
	allViews[ViewsPieces].GetView().Clear()
	_ = allViews[ViewsSources].Focus()
	return nil
}

func viewReleasesClose(_ *cui.Gui, _ *cui.View) error {
	currentEntryId = -1
	_ = allViews[ViewsReleases].GetView().SetCursor(0, 0)
	allViews[ViewsReleases].GetView().Clear()
	_ = allViews[ViewsPieces].Focus()
	return nil
}

func moveCursorUp(_ *cui.Gui, view *cui.View) error {
	view.MoveCursor(0, -1, false)
	return nil
}

func moveCursorDown(_ *cui.Gui, view *cui.View) error {
	_, y := view.Cursor()
	if str, _ := view.Line(y + 1); str != "" {
		view.MoveCursor(0, 1, false)
	}
	return nil
}

func layoutManager(_ *cui.Gui) error {
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
	case ViewsPieces:
		return listCui.NewViewDimensions(0, columnWidth, viewWidth-1, viewHeight-1)
	case ViewsReleases:
		return listCui.NewViewDimensions(columnHeight, 0, columnWidth-1, viewHeight-1)
	}
	return listCui.NewViewDimensions(0, 0, 0, 0)
}
