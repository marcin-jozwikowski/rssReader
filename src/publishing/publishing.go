package publishing

import "sort"

type Publishing struct {
	Name   string
	Pieces []*Piece
}

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

func (s *Publishing) GetPieceByAt(id int) *Piece {
	return s.Pieces[id]
}

func (s *Publishing) Sort() {
	sort.Sort(ByPieceTitle(s.Pieces))
}
