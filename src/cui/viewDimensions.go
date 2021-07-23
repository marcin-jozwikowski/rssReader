package cui

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
