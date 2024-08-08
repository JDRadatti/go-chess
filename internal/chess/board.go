package chess

type Player int8

const (
	WHITE Player = 0
	BLACK Player = 1
)

const (
	WIDTH         int = 8
	HEIGHT        int = 8
	NUM_SQUARES   int = 64
	CASTLE_OFFSET int = 2
)

// Board represents the state of a chess game.
// Board uses a 1D array with indicies as shown below.
//
//	8   0   1   2   3   4   5   6   7
//	7   8   9  10  11  12  13  14  15
//	6  16  17  18  19  20  21  22  23
//	5  24  25  26  27  28  29  30  31
//	4  32  33  34  35  36  37  38  39
//	3  40  41  42  43  44  45  46  47
//	2  48  49  50  51  52  53  54  55
//	1  56  57  58  59  60  61  62  63
//	   a   b   c   d   e   f   g   h
type Board struct {
	squares   [NUM_SQUARES]*Square
	turns     int
	whiteKing *Square
	blackKing *Square
	gameOver  bool
}

func NewBoardClassic() Board {
	board := Board{
		squares: InitSquaresClassic(),
	}

	board.whiteKing = board.squares[60]
	board.blackKing = board.squares[4]
	return board
}

// NewBoardFrom for testing purposes
func NewBoardFrom(b []byte) Board {
	board := Board{
		squares: InitSquaresFrom(b),
	}
	for i, s := range board.squares {
		if s.empty() {
			continue
		}
		switch s.piece.symbol {
		case KW:
			board.whiteKing = board.squares[i]
		case KB:
			board.blackKing = board.squares[i]
		}
	}
	return board
}

// Move executes a move from start to dest, if valid and updates
// necessary state.
func (b *Board) Move(start *Square, dest *Square) bool {

	if b.gameOver || start.empty() || start.piece.player != b.turn() {
		return false
	}

	if b.castle(start, dest) {
		return true
	} else if !b.validMove(start, dest) {
		return false
	}

	move := NewMove(start, dest, start.piece, dest.piece)
	b.makeMove(move)

	if b.inCheck(b.currentKing()) { // cannot move into a check
		b.undoMove(move)
		return false
	}
	b.turns++

	// note: this must be after incrementing turns
	check, mate, stale := b.checkOrMateOrStale()
	move.check = check
	move.mate = mate
	b.gameOver = move.mate || stale

	return true
}

func (b *Board) undoMove(move *Move) {
	move.startSquare.piece, move.destSquare.piece = move.startPiece, move.destPiece

	// Update state
	switch move.startPiece.symbol {
	case KW:
		b.whiteKing = move.startSquare
	case KB:
		b.blackKing = move.startSquare
	}

	move.startSquare.markUnmoved()
}

func (b *Board) makeMove(move *Move) {
	move.startSquare.piece, move.destSquare.piece = nil, move.startPiece

	// Update state
	switch move.destSquare.piece.symbol {
	case KW:
		b.whiteKing = move.destSquare
	case KB:
		b.blackKing = move.destSquare
	}

	move.startSquare.markMoved()
}

// validMove checks if a move from start to destination
// is valid. A move if valid if it satisfies the following criteria:
// NOTE: this function is stateless, so it does not check things like
// if the game is over and if its the correct players turn
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
// 4. Pawns can move two squares on their first move
// 5. en passant - pawn can capture diagonally iff:
//   - The capturing pawn must have advanced exactly three ranks to perform this move.
//   - The captured pawn must have moved two squares in one move, landing right next to the capturing pawn.
//   - The en passant capture must be performed on the turn immediately after the pawn being captured moves. If the player does not capture en passant on that turn, they no longer can do it later.
//
// 6. Pawns can upgrade when reaching the other side
func (b *Board) validMove(start *Square, dest *Square) bool {

	if !b.clearMove(start, dest) {
		return false
	}

	// check edge cases
	switch start.piece.symbol {
	case KW:
		fallthrough
	case KB:
		if b.attacked(dest) {
			return false
		}
	case PW:
		fallthrough
	case PB:
		indexDiff := start.index - dest.index
		if indexDiff < 0 {
			indexDiff = -indexDiff
		}
		if dest.empty() && indexDiff%WIDTH != 0 {
			return false // moved diagonal without capture
		} else if start.hasMoved() && indexDiff > WIDTH {
			return false // moved two forward on non-first move
		}
	}
	return true
}

// clearMove move checks if a moves from start to dest is:
// 1. possible based on the piece in start square
// 2. there are no pieces between start and dest
// 3. dest is either empty or contains an opponent piece
func (b *Board) clearMove(start *Square, dest *Square) bool {
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
		if !currSquare.empty() {
			if !currSquare.samePlayer(start) && i == steps {
				continue
			}
			return false
		}
	}
	return true
}

// attacked checks if the given square is attacked by
// an opponent piece. Uses current player's turn to check for opponent.
// A square is attacked by an opponent piece if the opponent piece
// can CAPTURE a piece on that square if it were their turn.
func (b *Board) attacked(square *Square) bool {
	for _, s := range b.squares {
		if s.index == square.index || s.empty() {
			continue
		}
		// Kings cannot attack
		if s.piece.king() {
			continue
		}
		if b.turn() != s.piece.player { // opponent piece
			if b.clearMove(s, square) {
				return true
			}
		}
	}
	return false
}

func (b *Board) turn() Player {
	return Player(b.turns % 2)
}

// inCheck returns true iff the current player is in check.
func (b *Board) inCheck(s *Square) bool {
	if s.empty() || !s.piece.king() {
		return false
	}
	return b.attacked(s)
}

// checkOrMateOrStale checks if the current player is either in
// check, checkmate, or stalemate
// This should not be called to get just check (instead, use inCheck)
func (b *Board) checkOrMateOrStale() (bool, bool, bool) {
	hasValidMoves := b.hasValidMoves(b.turn())
	check := b.inCheck(b.currentKing())
	return check, check && !hasValidMoves, !check && !hasValidMoves
}

// inCheck returns if the given player has valid moves.
// checkmate if checked and has no moves
// stalemate if not checked and has no moves
func (b *Board) hasValidMoves(player Player) bool {

	// Iterate through all valid moves for the given player
	// and check for a valid move
	for _, square := range b.squares {
		if square.empty() || square.piece.player != player {
			continue
		}

		for _, dir := range square.piece.directions {
			currIndex := square.index + dir
			if currIndex < 0 || currIndex >= NUM_SQUARES {
				continue
			}
			nextSquare := b.squares[currIndex]
			if b.validMove(square, nextSquare) {
				return check, false
			}
		}
	}
	return check, true
}

// castle returns true iff the move from start to dest is a valid castle move
// castle UPDATES STATE
// Castling rules:
//   - The King is not currently in check prior to castling, the Rook can be attacked prior to castling, but not the King
//   - The King is not in check on the square the King would be on after castling
//   - The King is not in check on any of the squares the King passes through while castling
//   - The King and the Rook involved have not moved yet during the game
//   - All of the squares in between the King and the Rook are unoccupied by another piece
//
// note: does not check for correct turn
func (b *Board) castle(start *Square, dest *Square) bool {
	if start.empty() || dest.empty() || !start.piece.king() || !dest.piece.rook() {
		return false
	}

	if start.hasMoved() || dest.hasMoved() {
		return false
	}

	if b.inCheck(start) {
		return false // cannot castle if in check
	}

	var left, right int
	var kingI, rookI int
	if start.index < dest.index { // kingside castle
		left, right = start.index, start.index+CASTLE_OFFSET
		kingI, rookI = right, right-1
	} else { // queenside castle
		left, right = dest.index+CASTLE_OFFSET, start.index
		kingI, rookI = left, left+1
		if !b.squares[left-1].empty() {
			return false
		}
	}
	for i := left; i <= right; i++ {
		if !b.squares[i].empty() && i != start.index ||
			b.attacked(b.squares[i]) {
			return false
		}
	}

	// clear rook
	start.piece, b.squares[kingI].piece = nil, start.piece
	dest.piece, b.squares[rookI].piece = nil, dest.piece
	b.turns++
	return true
}

// currentKing returns the king belonging to the player
// whose turn it is
func (b *Board) currentKing() *Square {
	switch b.turn() {
	case WHITE:
		return b.whiteKing
	case BLACK:
		return b.blackKing
	}
	return nil
}

func (b *Board) String() string {
	builder := ""
	counter := 0
	for _, square := range b.squares {
		if counter%8 == 0 {
			builder += "\n"
		}

		if square.empty() {
			builder += "_ "
		} else {
			builder += string(square.piece.symbol) + " "
		}
		counter++
	}
	return builder
}
