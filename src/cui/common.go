package cui

import cui "github.com/jroimartin/gocui"

func ExitRightNow(gui *cui.Gui, view *cui.View) error {
	return cui.ErrQuit
}