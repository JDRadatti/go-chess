package chess

type Player int8

const (
	WHITE Player = 0
	BLACK Player = 1
)

const (
	WIDTH       int = 8
	HEIGHT      int = 8
	NUM_SQUARES int = 64
)

type Board struct {
	squares   [NUM_SQUARES]Square
	turn      Player // 0 for white, 1 for black
	whiteKing *Square
	blackKing *Square
	gameOver  bool
}

func NewBoardClassic() Board {
	board := Board{
		squares: InitSquaresClassic(),
		turn:    WHITE,
	}

	board.whiteKing = &board.squares[60]
	board.blackKing = &board.squares[4]
	return board
}

// NewBoardFrom for testing purposes
func NewBoardFrom(b []byte) Board {
	board := Board{
		squares: InitSquaresFrom(b),
		turn:    WHITE,
	}
	for i, s := range board.squares {
		if s.empty() {
			continue
		}
		switch s.piece.symbol {
		case 'k':
			board.whiteKing = &board.squares[i]
		case 'K':
			board.blackKing = &board.squares[i]
		}
	}
	return board
}

// Move executes a move from start to dest, if valid and updates
// necessary state.
func (b *Board) Move(start Square, dest Square) bool {
	if !b.validMove(start, dest) {
		return false
	}

	// Update state
	start.piece, dest.piece = nil, start.piece
	switch start.piece.symbol {
	case 'k':
		b.whiteKing = &dest
	case 'K':
		b.blackKing = &dest
	}
	return true
}

// validMove checks if a move from start to destination
// is valid. A move if valid if it satisfies the following criteria:
// 1. game is not over
// 2. start contains a piece belonging to player of current turn
// 3. dest is a valid move direction of the piece on start
// 4. there are no pieces between current start and dest
// 5. end is empty or contains piece belonging to opponent
// 6. if in check, can only move king to safety
// edge cases:
// 1. Knights can jump pieces
// 2. King cannot move to square attacked by enemy(i.e. cannot move to its own death)
// 3. Casling
//   - The King is not currently in check prior to castling, the Rook can be attacked prior to castling, but not the King
//   - The King is not in check on the square the King would be on after castling
//   - The King is not in check on any of the squares the King passes through while castling
//   - The King and the Rook involved have not moved yet during the game
//   - All of the squares in between the King and the Rook are unoccupied by another piece
//
// 4. Pawns can move two squares on their first move
// 5. en passant - pawn can capture diagonally iff:
//   - The capturing pawn must have advanced exactly three ranks to perform this move.
//   - The captured pawn must have moved two squares in one move, landing right next to the capturing pawn.
//   - The en passant capture must be performed on the turn immediately after the pawn being captured moves. If the player does not capture en passant on that turn, they no longer can do it later.
//
// 6. Pawns can upgrade when reaching the other side
func (b *Board) validMove(start Square, dest Square) bool {

	if b.gameOver || start.empty() || start.piece.player != b.turn {
		return false
	}
	// check if in check
	// if in check, check checkmate

	// check if castle
	if !b.clearMove(start, dest) {
		return false
	}

	// check edge cases
	switch start.piece.symbol {
	case 'k':
		fallthrough
	case 'K':
		if b.attacked(dest) {
			return false
		}
	case 'p':
		if dest.empty() && start.index-dest.index != WIDTH {
			return false
		}
	case 'P':
		if dest.empty() && start.index-dest.index != -WIDTH {
			return false
		}
	}

	start.piece, dest.piece = nil, start.piece
	return true
}

// clearMove move checks if a moves from start to dest is:
// 1. possible based on the piece in start square
// 2. there are no pieces between start and dest
// 3. dest is either empty or contains an opponent piece
func (b *Board) clearMove(start Square, dest Square) bool {
	dir, steps := move(start, dest)
	if dir == INVALID || steps > start.piece.maxSteps {
		return false
	}
	if !start.piece.validDirection(dir) {
		return false
	}
	for i := steps; i > 0; i-- {
		currIndex := start.index + i*dir
		if currIndex < 0 || currIndex >= NUM_SQUARES {
			return false
		}
		currSquare := b.squares[currIndex]
		if !currSquare.empty() && steps > 1 ||
			currSquare.samePlayer(&start) && steps == 1 {
			return false
		}
	}
	return true
}

// attacked checks if the given square is attacked by
// an opponent piece. Uses current player's turn to check for opponent.
// A square is attacked by an opponent piece if the opponent piece
// can CAPTURE a piece on that square if it were their turn.
func (b *Board) attacked(square Square) bool {
	for _, s := range b.squares {
		if s.index == square.index || s.empty() {
			continue
		}
		// Kings cannot attack
		if s.piece.symbol == 'K' || s.piece.symbol == 'k' {
			continue
		}
		if b.turn != s.piece.player { // opponent piece
			if b.clearMove(s, square) {
				return true
			}
		}
	}
	return false
}
