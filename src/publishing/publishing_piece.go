package publishing

import (
	"rssReader/src/cui"
	"sort"
	"strconv"
)

type Piece struct {
	Title      string
	Releases   []*Release
	Publishing *Publishing
}

func (e Piece) ToString() string {
	line := e.Title + " | " + strconv.Itoa(len(e.Releases)) + " Releases"
	return line
}

func (e Piece) GetListViewItems() *[]cui.ListViewItem {
	var r []cui.ListViewItem
	for _, c := range e.Releases {
		r = append(r, c)
	}
	return &r
}

type ByPieceTitle []*Piece

func (a ByPieceTitle) Len() int           { return len(a) }
func (a ByPieceTitle) Less(i, j int) bool { return a[i].Title > a[j].Title } //desc
func (a ByPieceTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (e *Piece) addRelease(size int, url string, subtitle string) {
	e.Releases = append(e.Releases, &Release{Url: url, Size: size, Subtitle: subtitle, Piece: e})
	sort.Sort(ByReleaseSize(e.Releases))
}

func (e *Piece) GetReleaseAt(id int) *Release {
	return e.Releases[id]
}
