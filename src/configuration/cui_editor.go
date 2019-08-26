package configuration

import (
	"errors"
	cui "github.com/jroimartin/gocui"
	"log"
)

const ViewsFeedSources = "feedSources"
const ViewsFeedDetails = "feedDetails"
const ViewsFeedResults = "feedResults"
const ErrorsFailedAddingView = "Failed adding view"

var gui *cui.Gui
var err error
var config *Config
var ErrSave = errors.New("save")

func (configuration *Config) Edit() bool {
	config = configuration
	return createCUI()
}

func createCUI() bool {
	gui, err = cui.NewGui(cui.OutputNormal)
	if err != nil {
		log.Fatal("Failed to build CUI", err)
	}
	defer gui.Close()

	gui.SetManagerFunc(layoutManager)
	err = initAllViews()
	if err != nil {
		log.Fatal(ErrorsFailedAddingView, err)
	}
	initAllKeyBindings()

	finalError := gui.MainLoop()
	switch finalError {
	case nil:
	case cui.ErrQuit:
		return false
	}

	return true
}

func initAllKeyBindings() {
	if err = gui.SetKeybinding("", cui.KeyCtrlS, cui.ModNone, saveAndExit); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding("", cui.KeyCtrlD, cui.ModNone, exitRightNow); err != nil {
		log.Fatal("Failed to set keybindings")
	}
	if err = gui.SetKeybinding("", cui.KeyCtrlU, cui.ModNone, doSomething); err != nil {
		log.Fatal("Failed to set keybindings")
	}
}

func initAllViews() error {
	gui.SelFgColor = cui.ColorGreen | cui.AttrBold
	gui.BgColor = cui.ColorDefault
	gui.Highlight = true

	_, _ = gui.SetView(ViewsFeedResults, 0, 0, 1, 1)
	_, _ = gui.SetView(ViewsFeedDetails, 0, 0, 1, 1)
	_, _ = gui.SetView(ViewsFeedSources, 0, 0, 1, 1)

	return focusView(ViewsFeedSources)
}

func focusView(viewName string) error {
	v, err := gui.View(viewName)
	if err != nil {
		return err
	}
	v.Highlight = true
	if _, err = gui.SetCurrentView(v.Name()); err != nil {
		return err
	}
	return nil
}

func doSomething(gui *cui.Gui, view *cui.View) error {
	return focusView(ViewsFeedResults)
}

func exitRightNow(gui *cui.Gui, view *cui.View) error {
	return cui.ErrQuit
}

func saveAndExit(i *cui.Gui, view *cui.View) error {
	return ErrSave
}

func layoutManager(gui *cui.Gui) error {
	viewWidth, viewHeight := gui.Size()
	columnWidth := viewWidth / 3
	columnHeight := viewHeight / 2

	if _, err := gui.SetView(ViewsFeedSources, 0, 0, columnWidth-1, columnHeight-1); err != nil {
		log.Fatal("Cannot update sites view", err)
	}

	if _, err := gui.SetView(ViewsFeedDetails, columnWidth, 0, viewWidth-1, columnHeight-1); err != nil {
		log.Fatal("Cannot update sites view", err)
	}

	if _, err := gui.SetView(ViewsFeedResults, 0, columnHeight, viewWidth-1, viewHeight-1); err != nil {
		log.Fatal("Cannot update sites view", err)
	}

	return nil
}
