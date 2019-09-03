package configuration

import (
	"errors"
	"fmt"
	cui "github.com/jroimartin/gocui"
	"log"
	"math/rand"
	"strconv"
	"strings"
)

const ViewsFeedSources = "feedSources"
const ViewsFeedDetails = "feedDetails"
const ViewsFeedResults = "feedResults"
const ViewHelp = "help"
const ErrorsFailedAddingView = "Failed adding view"

var err error
var config *Config
var editedFeed *FeedSource
var ErrSave = errors.New("save")
var allViews = make(map[string]*ListView)
var currentInputPrompt *InputView
var viewToFallBackTo = ViewsFeedSources

type CliView struct {
	gui  *cui.Gui
	view *cui.View
}

// LIST VIEW
type ListView struct {
	CliView
	items     []string
	DrawItems func()
}

func (listView *ListView) Init(gui *cui.Gui, name string, items []string, title string) {
	listView.items = items
	listView.gui = gui

	listView.view, _ = listView.gui.SetView(name, 0, 0, 1, 1)
	listView.view.Clear()
	listView.view.Title = " " + title + " "

	allViews[name] = listView
}

func (listView *ListView) Draw() {
	x0, y0, x1, y1 := getViewDimensions(listView.gui, listView.view.Name())
	listView.view, err = listView.gui.SetView(listView.view.Name(), x0, y0, x1, y1)
	if err != nil {
		log.Fatal("Cannot update sites view", err)
	}
	listView.view.Clear()

	if listView.DrawItems != nil {
		listView.DrawItems()
	} else {
		for _, item := range listView.items {
			_, _ = fmt.Fprintln(listView.view, item)
		}
	}
}

func (listView *CliView) Focus() error {
	listView.view.Highlight = true
	if _, err = listView.gui.SetCurrentView(listView.view.Name()); err != nil {
		panic(err)
		return nil
	}
	return nil
}

// INPUT VIEW
type InputView struct {
	CliView
	content          string
	onSubmitCallback func(content string) error
}

func (input *InputView) Init(gui *cui.Gui, content string, title string, onSubmitCallback func(content string) error) {
	input.content = content
	input.onSubmitCallback = onSubmitCallback
	input.gui = gui
	x0, y0, x1, y1 := getCenteredViewDimensions(gui, 1)
	input.view, _ = input.gui.SetView("input_"+strconv.Itoa(rand.Int()), x0, y0, x1, y1)
	input.view.Title = title
	input.view.Editable = true

	_, _ = gui.SetCurrentView(input.view.Name())
	_, _ = fmt.Fprint(input.view, content)

	if err = gui.SetKeybinding(input.view.Name(), cui.KeyEnter, cui.ModNone, inputViewApply); err != nil {
		log.Fatal("Failed to set keybindings")
	}

	if err = gui.SetKeybinding(input.view.Name(), cui.KeyEsc, cui.ModNone, inputExit); err != nil {
		log.Fatal("Failed to set keybindings")
	}

	gui.Cursor = true
	currentInputPrompt = input
	_ = input.Focus()
}

func (input *InputView) Close() error {
	currentInputPrompt = nil
	input.gui.DeleteKeybindings(input.view.Name())
	return deleteNamedView(input.gui, input.view)
}

func inputExit(gui *cui.Gui, view *cui.View) error {
	return currentInputPrompt.Close()
}

func inputViewApply(gui *cui.Gui, view *cui.View) error {
	return currentInputPrompt.onSubmitCallback(strings.Trim(view.ViewBuffer(), "\n\t\r"))
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
	if err = gui.SetKeybinding(ViewsFeedDetails, cui.KeyEnter, cui.ModNone, onEnterInDetails); err != nil {
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

func onEnterInDetails(gui *cui.Gui, view *cui.View) error {
	list := allViews[view.Name()]
	_, y := view.Cursor()
	itemCount := len(list.items)
	viewToFallBackTo = view.Name()
	if y < itemCount {
		item := editedFeed.SearchPhrases[y]
		namePrompt := new(InputView)
		namePrompt.Init(gui, item, "Edit (Esc to cancel)", func(content string) error {
			editedFeed.SearchPhrases[y] = content
			allViews[ViewsFeedDetails].Draw()
			_ = deleteNamedView(currentInputPrompt.gui, currentInputPrompt.view)
			return nil
		})
	} else if y == (itemCount + 1) {
		v := new(InputView)
		v.Init(gui, editedFeed.PostProcess, "Edit (Esc to cancel)", func(content string) error {
			editedFeed.PostProcess = content
			_ = deleteNamedView(currentInputPrompt.gui, currentInputPrompt.view)
			return nil
		})
	}

	return nil
}

func closeHelp(gui *cui.Gui, view *cui.View) error {
	return deleteNamedView(gui, allViews[ViewHelp].view)
}

func deleteNamedView(gui *cui.Gui, view *cui.View) error {
	_ = allViews[viewToFallBackTo].Focus()
	gui.Cursor = false
	return gui.DeleteView(view.Name())
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
	allViews[ViewsFeedDetails].view.Clear()
	return allViews[ViewsFeedSources].Focus()
}

func editCurrentSource(gui *cui.Gui, view *cui.View) error {
	_, selectedItem := view.Cursor()
	editedFeed = config.GetFeedAt(selectedItem)
	v := new(ListView)
	v.Init(gui, ViewsFeedDetails, editedFeed.SearchPhrases, editedFeed.Url)
	v.DrawItems = func() {
		for _, item := range v.items {
			_, _ = fmt.Fprintln(v.view, item)
		}
		_, _ = fmt.Fprintln(v.view, "")
		_, _ = fmt.Fprintln(v.view, "Edit Post-processing")
	}
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
	columnWidth := 2 * viewWidth / 3
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
