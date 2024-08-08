package chess

type Square struct {
	index int
	moved bool
	piece *Piece
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
	return s.moved
}

func (s *Square) markMoved() {
	s.moved = true
}

func (s *Square) markUnmoved() {
	s.moved = false
}

func (s *Square) String() string {
	return FILES[s.file()] + RANKS[s.rank()]
}

// samePlayer returns true iff both squares have a piece that are
// owned by the same player
func (s *Square) samePlayer(o *Square) bool {
	return !s.empty() && !o.empty() && s.piece.player == o.piece.player
}

func valid(col byte, row byte) bool {
	return 'a' <= col && col <= 'h' && '1' <= row && row <= '8'
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
