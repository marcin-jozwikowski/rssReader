package configuration

import (
	"errors"
	"fmt"
	cui "github.com/jroimartin/gocui"
	"log"
)

const ViewsFeedSources = "feedSources"
const ViewsFeedDetails = "feedDetails"
const ViewsFeedResults = "feedResults"
const ViewHelp = "help"
const ErrorsFailedAddingView = "Failed adding view"

var err error
var config *Config
var ErrSave = errors.New("save")
var allViews = make(map[string]*ListView)
var viewToFallBackTo = ViewsFeedSources

type CliView struct {
	gui      *cui.Gui
	viewName string
}

type ListView struct {
	CliView
	items []string
}

type InputView struct {
	CliView
	content string
}

func (listView *ListView) Init(gui *cui.Gui, name string, items []string, title string) {
	listView.viewName = name
	listView.items = items
	listView.gui = gui

	view, _ := listView.gui.SetView(name, 0, 0, 1, 1)
	view.Clear()
	view.Title = " " + title + " "

	for _, item := range listView.items {
		_, _ = fmt.Fprintln(view, item)
	}

	allViews[name] = listView
}

func (listView *ListView) Draw() {
	x0, y0, x1, y1 := getViewDimensions(listView.gui, listView.viewName)
	if _, err := listView.gui.SetView(listView.viewName, x0, y0, x1, y1); err != nil {
		log.Fatal("Cannot update sites view", err)
	}
}

func (listView *ListView) Focus() error {
	v, err := listView.gui.View(listView.viewName)
	if err != nil {
		panic(err)
		return err
	}
	v.Highlight = true
	if _, err = listView.gui.SetCurrentView(listView.viewName); err != nil {
		panic(err)
		return nil
	}
	return nil
}

func (listView *ListView) ClearContent() {
	v, _ := listView.gui.View(listView.viewName)
	v.Clear()
}

func (configuration *Config) Edit() bool {
	config = configuration
	return createCUI()
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

	v := new(ListView)
	v.Init(gui, ViewsFeedSources, config.AsList(), "Sources")

	v = new(ListView)
	v.Init(gui, ViewsFeedDetails, []string{}, "Select source to view details")

	_ = allViews[ViewsFeedSources].Focus()

	initAllKeyBindings(gui)

	finalError := gui.MainLoop()
	switch finalError {
	case nil:
	case cui.ErrQuit:
		return false
	}

	return true
}

func initAllKeyBindings(gui *cui.Gui) {
	// ViewsFeedSources
	if err = gui.SetKeybinding(ViewsFeedSources, cui.KeyArrowDown, cui.ModNone, moveCursorDown); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding(ViewsFeedSources, cui.KeyArrowUp, cui.ModNone, moveCursorUp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding(ViewsFeedSources, cui.KeyEnter, cui.ModNone, editCurrentSource); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding(ViewsFeedSources, cui.KeyArrowRight, cui.ModNone, editCurrentSource); err != nil {
		log.Fatal("Failed to set keybindings")
	}

	// ViewsFeedDetails
	if err = gui.SetKeybinding(ViewsFeedDetails, cui.KeyArrowDown, cui.ModNone, moveCursorDown); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding(ViewsFeedDetails, cui.KeyArrowUp, cui.ModNone, moveCursorUp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding(ViewsFeedDetails, cui.KeyArrowLeft, cui.ModNone, focusOnSources); err != nil {
		log.Fatal("Failed to set keybindings")
	}

	// help
	if err = gui.SetKeybinding(ViewHelp, cui.KeyCtrlH, cui.ModNone, closeHelp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding(ViewsFeedSources, cui.KeyCtrlH, cui.ModNone, showHelp); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding(ViewsFeedDetails, cui.KeyCtrlH, cui.ModNone, showHelp); err != nil {
		log.Fatal("Failed to set keybindings")
	}

	// general
	if err = gui.SetKeybinding("", cui.KeyCtrlS, cui.ModNone, saveAndExit); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding("", cui.KeyCtrlQ, cui.ModNone, exitRightNow); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding("", cui.KeyCtrlU, cui.ModNone, doSomething); err != nil {
		log.Fatal("Failed to set keybindings")
	}
}

func closeHelp(gui *cui.Gui, view *cui.View) error {
	_ = allViews[viewToFallBackTo].Focus()
	return gui.DeleteView(ViewHelp)
}

func showHelp(gui *cui.Gui, view *cui.View) error {
	viewToFallBackTo = view.Name()
	x0, y0, x1, y1 := getCenteredViewDimensions(gui, 2)
	v, _ := gui.SetView(ViewHelp, x0, y0, x1, y1)
	_, _ = gui.SetCurrentView(ViewHelp)
	v.Clear()
	v.Title = " Help (press Ctrl-H again to quit)"
	_, _ = fmt.Fprintln(v, "Show help here")
	_, _ = gui.SetCurrentView(ViewHelp)

	return nil
}

func focusOnSources(gui *cui.Gui, view *cui.View) error {
	allViews[ViewsFeedDetails].ClearContent()
	return allViews[ViewsFeedSources].Focus()
}

func editCurrentSource(gui *cui.Gui, view *cui.View) error {
	_, selectedItem := view.Cursor()
	selectedFeed := config.GetFeedAt(selectedItem)
	v := new(ListView)
	v.Init(gui, ViewsFeedDetails, selectedFeed.SearchPhrases, selectedFeed.Url)
	_ = v.Focus()

	return nil
}

func moveCursorUp(i *cui.Gui, view *cui.View) error {
	view.MoveCursor(0, -1, false)
	return nil
}

func moveCursorDown(i *cui.Gui, view *cui.View) error {
	view.MoveCursor(0, 1, false)
	return nil
}

func doSomething(gui *cui.Gui, view *cui.View) error {
	return allViews[ViewsFeedDetails].Focus()
}

func exitRightNow(gui *cui.Gui, view *cui.View) error {
	return cui.ErrQuit
}

func saveAndExit(i *cui.Gui, view *cui.View) error {
	return ErrSave
}

func getViewDimensions(gui *cui.Gui, viewName string) (int, int, int, int) {
	viewWidth, viewHeight := gui.Size()
	columnWidth := viewWidth / 3
	columnHeight := viewHeight / 2

	switch viewName {
	case ViewsFeedSources:
		return 0, 0, columnWidth - 1, columnHeight - 1
	case ViewsFeedResults:
		return 0, columnHeight, viewWidth - 1, viewHeight - 1
	case ViewsFeedDetails:
		return columnWidth, 0, viewWidth - 1, columnHeight - 1
	}
	return 0, 0, 0, 0
}

func getCenteredViewDimensions(gui *cui.Gui, lines int) (int, int, int, int) {
	viewWidth, viewHeight := gui.Size()
	top := (viewHeight / 2) - (lines / 2)
	return 10, top, viewWidth - 10, top + lines + 1
}

func layoutManager(gui *cui.Gui) error {
	for _, listView := range allViews {
		listView.Draw()
	}
	return nil
}
