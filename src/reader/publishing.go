package reader

import "sort"

type Publishing struct {
	Name   string
	Pieces []*Piece
}

type Piece struct {
	Title      string
	Releases   []*Release
	Publishing *Publishing
}

type ByPieceTitle []*Piece

func (a ByPieceTitle) Len() int           { return len(a) }
func (a ByPieceTitle) Less(i, j int) bool { return a[i].Title > a[j].Title } //desc
func (a ByPieceTitle) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

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

func (s *Publishing) AddRelease(title string, subtitle string, size int, url string) {
	piece := s.getPieceByTitle(title)
	if piece == nil {
		piece = s.addPiece(title)
	}
	piece.addRelease(size, url, subtitle)
}

func (s *Publishing) getPieceByTitle(title string) *Piece {
	for e := range s.Pieces {
		if s.Pieces[e].Title == title {
			return s.Pieces[e]
		}
	}
	return nil
}

func (s *Publishing) addPiece(title string) *Piece {
	piece := Piece{Title: title, Publishing: s}
	s.Pieces = append(s.Pieces, &piece)

	return &piece
}

func (s *Publishing) getPieceByAt(id int) *Piece {
	return s.Pieces[id]
}

func (s *Publishing) Sort() {
	sort.Sort(ByPieceTitle(s.Pieces))
}

func (e *Piece) addRelease(size int, url string, subtitle string) {
	e.Releases = append(e.Releases, &Release{Url: url, Size: size, Subtitle: subtitle, Piece: e})
	sort.Sort(ByReleaseSize(e.Releases))
}

func (e *Piece) getReleaseAt(id int) *Release {
	return e.Releases[id]
}
