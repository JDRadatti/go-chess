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

var Pieces map[byte]*Piece = map[byte]*Piece{
	'k': KingW,
	'K': KingB,
	'q': QueenW,
	'Q': QueenB,
	'r': RookW,
	'R': RookB,
	'b': BishopW,
	'B': BishopB,
	'n': KnightW,
	'N': KnightB,
	'p': PawnW,
	'P': PawnB,
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
