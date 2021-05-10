package cui

import cui "github.com/jroimartin/gocui"

type CliView struct {
	gui  *cui.Gui
	view *cui.View
}

func (view *CliView) GetView() *cui.View {
	return view.view
}

type ViewDimensions struct {
	top    int
	left   int
	right  int
	bottom int
}

func NewViewDimensions(top int, left int, right int, bottom int) ViewDimensions {
	return ViewDimensions{
		top:    top,
		left:   left,
		right:  right,
		bottom: bottom,
	}
}
