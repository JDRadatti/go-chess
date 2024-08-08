package chess

type Square struct {
	index     int
	moveCount int
	piece     *Piece
}

func (s *Square) file() int {
	return s.index % int(WIDTH)
}

func (s *Square) rank() int {
	return s.index / int(WIDTH)
}

func (s *Square) empty() bool {
	return s.piece == nil
}

func (s *Square) hasMoved() bool {
	return s.moveCount > 0
}

func (s *Square) markMoved() {
	s.moveCount++
}

func (s *Square) markUnmoved() {
	if s.moveCount <= 0 {
		return
	}
	s.moveCount--
}

func (s *Square) String() string {
	return FILES[s.file()] + RANKS[s.rank()]
}

// samePlayer returns true iff both squares have a piece that are
// owned by the same player
func (s *Square) samePlayer(o *Square) bool {
	return !s.empty() && !o.empty() && s.piece.player == o.piece.player
}

func InitSquaresFrom(board []byte) [NUM_SQUARES]*Square {
	squares := [NUM_SQUARES]*Square{}
	for i, symbol := range board {
		squares[i] = &Square{
			index: i,
			piece: Pieces[symbol],
		}
	}
	return squares
}

func InitSquaresClassic() [NUM_SQUARES]*Square {
	board := []byte{
		RB, NB, BB, QB, KB, BB, NB, RB,
		PB, PB, PB, PB, PB, PB, PB, PB,
		EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
		EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY, EMPTY,
		PW, PW, PW, PW, PW, PW, PW, PW,
		RW, NW, BW, QW, KW, BW, NW, RW,
	}
	return InitSquaresFrom(board)
}
