package reader

type Release struct {
	Size           int
	Url            string
	Subtitle       string
	InternalResult string
	Piece          *Piece
}

type ByReleaseSize []*Release

func (a ByReleaseSize) Len() int           { return len(a) }
func (a ByReleaseSize) Less(i, j int) bool { return a[i].Size < a[j].Size }
func (a ByReleaseSize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
