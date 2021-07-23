package cui

import cui "github.com/jroimartin/gocui"

type CliView struct {
	gui  *cui.Gui
	view *cui.View
}

func (view *CliView) GetView() *cui.View {
	return view.view
}

func (view *CliView) Focus() error {
	view.view.Highlight = true
	if _, err := view.gui.SetCurrentView(view.view.Name()); err != nil {
		panic(err)
		return nil
	}
	return nil
}
