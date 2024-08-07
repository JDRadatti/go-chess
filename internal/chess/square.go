package chess

type Square struct {
	index int
	moved int8
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
	return s.moved == 1
}

func (s *Square) markMoved() {
	s.moved = 1
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
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	}
	return InitSquaresFrom(board)
}
