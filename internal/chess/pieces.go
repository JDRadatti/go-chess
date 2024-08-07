// Piece represents a piece on a chess board
// Piece must not contain logic pertaining to a board because
// each game uses a pointer to the same pieces to avoid unnecessary allocs.
package chess

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

const KW byte = 'k'
const KB byte = 'K'
const QW byte = 'q'
const QB byte = 'Q'
const RW byte = 'r'
const RB byte = 'R'
const BW byte = 'b'
const BB byte = 'B'
const NW byte = 'n'
const NB byte = 'N'
const PW byte = 'p'
const PB byte = 'P'
const EMPTY byte = ' '

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

var KingW = &Piece{
	symbol:     'k',
	directions: []int{NORTH, EAST, SOUTH, WEST, NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   1,
	player:     WHITE,
}
var KingB = &Piece{
	symbol:     'K',
	directions: []int{NORTH, EAST, SOUTH, WEST, NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   1,
	player:     BLACK,
}

var QueenW = &Piece{
	symbol:     'q',
	directions: []int{NORTH, EAST, SOUTH, WEST, NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   WIDTH,
	player:     WHITE,
}

var QueenB = &Piece{
	symbol:     'Q',
	directions: []int{NORTH, EAST, SOUTH, WEST, NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   WIDTH,
	player:     BLACK,
}

var BishopW = &Piece{
	symbol:     'b',
	directions: []int{NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   WIDTH,
	player:     WHITE,
}

var BishopB = &Piece{
	symbol:     'B',
	directions: []int{NORTHWEST, NORTHEAST, SOUTHWEST, SOUTHEAST},
	maxSteps:   WIDTH,
	player:     BLACK,
}

var RookW = &Piece{
	symbol:     'r',
	directions: []int{NORTH, EAST, SOUTH, WEST},
	maxSteps:   WIDTH,
	player:     WHITE,
}

var RookB = &Piece{
	symbol:     'R',
	directions: []int{NORTH, EAST, SOUTH, WEST},
	maxSteps:   WIDTH,
	player:     BLACK,
}

var KnightW = &Piece{
	symbol:     'n',
	directions: []int{KNIGHT0, KNIGHT1, KNIGHT2, KNIGHT3, KNIGHT4, KNIGHT5, KNIGHT6, KNIGHT7},
	maxSteps:   1,
	player:     WHITE,
}

var KnightB = &Piece{
	symbol:     'N',
	directions: []int{KNIGHT0, KNIGHT1, KNIGHT2, KNIGHT3, KNIGHT4, KNIGHT5, KNIGHT6, KNIGHT7},
	maxSteps:   1,
	player:     BLACK,
}

var PawnW = &Piece{
	symbol:     'p',
	directions: []int{NORTH, NORTHWEST, NORTHEAST},
	maxSteps:   1,
	player:     WHITE,
}

var PawnB = &Piece{
	symbol:     'P',
	directions: []int{SOUTH, SOUTHWEST, SOUTHEAST},
	maxSteps:   1,
	player:     BLACK,
}
