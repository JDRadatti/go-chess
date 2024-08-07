// Piece represents a piece on a chess board
// Piece must not contain logic pertaining to a board because
// each game uses a pointer to the same pieces to avoid unnecessary allocs.
package chess

const (
	KW    byte = 'k'
	KB    byte = 'K'
	QW    byte = 'q'
	QB    byte = 'Q'
	RW    byte = 'r'
	RB    byte = 'R'
	BW    byte = 'b'
	BB    byte = 'B'
	NW    byte = 'n'
	NB    byte = 'N'
	PW    byte = 'p'
	PB    byte = 'P'
	EMPTY byte = ' '
)

var Pieces map[byte]*Piece = map[byte]*Piece{
	KW: KingW,
	KB: KingB,
	QW: QueenW,
	QB: QueenB,
	RW: RookW,
	RB: RookB,
	BW: BishopW,
	BB: BishopB,
	NW: KnightW,
	NB: KnightB,
	PW: PawnW,
	PB: PawnB,
}

type Piece struct {
	symbol     byte   // byte representation of this piece
	directions []int  // directions allowed for this piece
	maxSteps   int    // number of steps in each direction allowed
	player     Player // owner of this piece
}

func (p *Piece) validDirection(dir int) bool {
	for _, d := range p.directions {
		if dir == d {
			return true
		}
	}
	return false
}

func (p *Piece) king() bool {
	return p.symbol == KW || p.symbol == KB
}

func (p *Piece) rook() bool {
	return p.symbol == RW || p.symbol == RB
}

func (p *Piece) bishop() bool {
	return p.symbol == BW || p.symbol == BB
}

func (p *Piece) pawn() bool {
	return p.symbol == PW || p.symbol == PB
}

func (p *Piece) knight() bool {
	return p.symbol == NW || p.symbol == NB
}

func (p *Piece) queen() bool {
	return p.symbol == QW || p.symbol == QB
}


var KingW = &Piece{
	symbol:     KW,
	directions: []int{NORTH, EAST, SOUTH, WEST, NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   1,
	player:     WHITE,
}
var KingB = &Piece{
	symbol:     KB,
	directions: []int{NORTH, EAST, SOUTH, WEST, NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   1,
	player:     BLACK,
}

var QueenW = &Piece{
	symbol:     QW,
	directions: []int{NORTH, EAST, SOUTH, WEST, NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   WIDTH,
	player:     WHITE,
}

var QueenB = &Piece{
	symbol:     QB,
	directions: []int{NORTH, EAST, SOUTH, WEST, NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   WIDTH,
	player:     BLACK,
}

var BishopW = &Piece{
	symbol:     BW,
	directions: []int{NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   WIDTH,
	player:     WHITE,
}

var BishopB = &Piece{
	symbol:     BB,
	directions: []int{NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   WIDTH,
	player:     BLACK,
}

var RookW = &Piece{
	symbol:     RW,
	directions: []int{NORTH, EAST, SOUTH, WEST},
	maxSteps:   WIDTH,
	player:     WHITE,
}

var RookB = &Piece{
	symbol:     RB,
	directions: []int{NORTH, EAST, SOUTH, WEST},
	maxSteps:   WIDTH,
	player:     BLACK,
}

var KnightW = &Piece{
	symbol:     NW,
	directions: []int{KNIGHT0, KNIGHT1, KNIGHT2, KNIGHT3, KNIGHT4, KNIGHT5, KNIGHT6, KNIGHT7},
	maxSteps:   1,
	player:     WHITE,
}

var KnightB = &Piece{
	symbol:     NB,
	directions: []int{KNIGHT0, KNIGHT1, KNIGHT2, KNIGHT3, KNIGHT4, KNIGHT5, KNIGHT6, KNIGHT7},
	maxSteps:   1,
	player:     BLACK,
}

var PawnW = &Piece{
	symbol:     PW,
	directions: []int{NORTH, NORTHWEST, NORTHEAST, NORTH2},
	maxSteps:   2,
	player:     WHITE,
}

var PawnB = &Piece{
	symbol:     PB,
	directions: []int{SOUTH, SOUTHWEST, SOUTHEAST, SOUTH2},
	maxSteps:   2,
	player:     BLACK,
}
