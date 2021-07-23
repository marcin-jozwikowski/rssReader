package publishing

import (
	"strconv"
)

type Release struct {
	Size           int
	Url            string
	Subtitle       string
	InternalResult string
	Piece          *Piece
}

func (r Release) ToString() string {
	line := strconv.Itoa(r.Size) + " MB | " + r.Piece.Title + r.Subtitle + " | "

	if r.InternalResult == "" {
		line += r.Piece.Publishing.Name
	} else {
		line += r.InternalResult
	}

	return line
}

type ByReleaseSize []*Release

func (a ByReleaseSize) Len() int           { return len(a) }
func (a ByReleaseSize) Less(i, j int) bool { return a[i].Size < a[j].Size }
func (a ByReleaseSize) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
