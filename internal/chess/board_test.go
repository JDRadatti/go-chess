package chess

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"testing"
)

func checkState(t *testing.T, board *Board) {
	// make sure the kings are tracked properly
	found := 0
	for i, square := range board.squares {
		if square.empty() {
			continue
		}
		if square.piece.symbol == 'k' {
			assert.Equal(t, i, board.whiteKing.index)
			found++
		} else if square.piece.symbol == 'K' {
			assert.Equal(t, i, board.blackKing.index)
			found++
		}
	}
	assert.Equal(t, 2, found, "both kings must exist")
}

func TestValidMoves(t *testing.T) {
	inputs := []struct {
		name         string
		board        []byte
		startSquares []int // index of square in board.squares
		destSquares  []int // index of square in board.squares
		expected     []bool
		turn         []Player
	}{
		{
			name: "basic pawn movements",
			board: []byte{
				'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
				'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
				'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
			},
			startSquares: []int{48, 48, 48, 8, 8, 8, 8},
			destSquares:  []int{40, 41, 39, 16, 0, 24, 26}, // 39 tests edge case (literally)
			expected:     []bool{true, false, false, true, false, true, false},
			turn:         []Player{WHITE, WHITE, WHITE, BLACK, BLACK, BLACK, BLACK},
		},
		{
			name: "basic knight movement",
			board: []byte{
				'R', ' ', 'B', 'Q', 'K', 'B', 'N', 'R',
				'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
				'N', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'N', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'p', 'p', 'p', ' ', 'p', ' ', 'p', 'p',
				'r', 'n', 'b', 'p', 'k', 'b', 'n', 'r',
			},
			startSquares: []int{57, 57, 57, 16, 16, 16, 36, 36, 36, 36, 36, 36, 36, 36},
			destSquares:  []int{40, 41, 42, 0, 10, 1, 19, 21, 26, 30, 42, 46, 51, 53},
			expected:     []bool{true, false, true, false, false, true, true, true, true, true, true, true, true, true},
			turn:         []Player{WHITE, WHITE, WHITE, BLACK, BLACK, BLACK, BLACK, BLACK, BLACK, BLACK, BLACK, BLACK, BLACK, BLACK},
		},
		{
			name: "basic rook movement",
			board: []byte{
				'R', ' ', 'B', 'Q', 'K', 'B', 'N', 'R',
				' ', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'p', ' ', ' ', ' ', 'N', ' ', ' ', ' ',
				'P', ' ', ' ', ' ', ' ', ' ', ' ', 'r',
				'b', 'p', 'p', ' ', 'p', ' ', 'p', 'p',
				' ', 'n', 'b', 'p', 'k', 'b', 'n', ' ',
			},
			startSquares: []int{0, 0, 0, 47, 47, 47, 47, 47},
			destSquares:  []int{1, 24, 56, 56, 23, 47, 41, 0},
			expected:     []bool{true, true, false, false, true, false, true, false},
			turn:         []Player{BLACK, BLACK, BLACK, WHITE, WHITE, WHITE, WHITE, WHITE},
		},
		{
			name: "basic bishop movement",
			board: []byte{
				'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
				'P', 'P', 'P', ' ', 'P', 'P', 'P', 'P',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'b', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'p', 'p', 'p', ' ', 'p', 'p', 'p', 'p',
				'r', 'n', 'b', 'p', 'k', 'b', 'n', 'r',
			},
			startSquares: []int{2, 2, 2, 32, 32},
			destSquares:  []int{11, 9, 16, 11, 7},
			expected:     []bool{true, false, false, true, false},
			turn:         []Player{BLACK, BLACK, BLACK, WHITE, WHITE},
		},
		{
			name: "basic queen movement",
			board: []byte{
				'R', 'N', 'B', ' ', 'K', 'B', 'N', 'R',
				'P', 'P', 'P', ' ', 'P', 'P', 'P', 'P',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'q', ' ', ' ', ' ', 'Q', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'p', 'p', 'p', ' ', 'p', 'p', 'p', 'p',
				'r', 'n', 'b', ' ', 'k', 'b', 'n', 'r',
			},
			startSquares: []int{37, 37, 37, 37, 33, 33, 33},
			destSquares:  []int{34, 51, 45, 39, 19, 17, 36},
			expected:     []bool{true, true, true, true, true, true, true},
			turn:         []Player{BLACK, BLACK, BLACK, BLACK, WHITE, WHITE, WHITE},
		},
		{
			name: "basic king movement",
			board: []byte{
				' ', 'K', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'p', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'b', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', 'k', ' ', ' ',
			},
			startSquares: []int{1, 1, 1, 1, 61},
			destSquares:  []int{0, 2, 8, 10, 60},
			expected:     []bool{false, false, true, false, true},
			turn:         []Player{BLACK, BLACK, BLACK, BLACK, WHITE},
		},
	}

	for j, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			board := NewBoardFrom(input.board) // reset board every time
			assert.Equal(t, len(input.startSquares), len(input.destSquares))
			assert.Equal(t, len(input.startSquares), len(input.expected))
			assert.Equal(t, len(input.startSquares), len(input.turn))
			for i, startI := range input.startSquares {
				start := board.squares[startI]
				dest := board.squares[input.destSquares[i]]
				board.turns = int(input.turn[i])
				assert.Equal(t, input.expected[i], board.validMove(start, dest),
					fmt.Sprintf("test %d, subtest %d", j, i))
			}

		})
	}
}

// TestMoves tests a series of moves with
// Move() instead of validMove() [which is used by TestValidMove].
// The difference between these two functions is Move() updates state.
func TestMoves(t *testing.T) {
	inputs := []struct {
		name         string
		board        []byte
		startSquares []int // index of square in board.squares
		destSquares  []int // index of square in board.squares
		expected     []bool
	}{
		{
			name: "london opening. black mirrors. opposite castle. all valid",
			board: []byte{
				'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
				'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
				'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
			},
			startSquares: []int{51, 11, 58, 2, 52, 12, 61, 5, 62, 15, 60, 1, 57, 3, 48, 4},
			destSquares:  []int{35, 27, 37, 29, 44, 20, 43, 19, 47, 23, 63, 16, 42, 11, 40, 0},
			expected:     []bool{true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true},
		},
	}

	for j, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			board := NewBoardFrom(input.board)
			assert.Equal(t, len(input.startSquares), len(input.destSquares))
			assert.Equal(t, len(input.startSquares), len(input.expected))
			for i, startI := range input.startSquares {
				start := board.squares[startI]
				dest := board.squares[input.destSquares[i]]
				assert.Equal(t, input.expected[i], board.Move(start, dest),
					fmt.Sprintf("test %d, subtest %d", j, i))
				log.Println(board.String())
			}

		})

	}
}

func TestCheck(t *testing.T) {
	inputs := []struct {
		name      string
		board     []byte
		kingIndex int
		player    Player // the player to check if in check
		expected  bool   // whether the given board is a check for player
	}{
		{
			name: "basic check on white king with queen",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'Q', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			player:    WHITE,
			kingIndex: 60,
			expected:  true,
		},
		{
			name: "basic check on black king with queen",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'q', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			kingIndex: 4,
			player:    BLACK,
			expected:  true,
		},
		{
			name: "basic game position without check",
			board: []byte{
				'R', 'N', ' ', 'Q', 'K', ' ', 'N', 'R',
				'P', 'P', 'P', ' ', ' ', 'P', 'P', ' ',
				' ', ' ', ' ', 'B', 'P', ' ', ' ', 'P',
				' ', ' ', ' ', 'P', ' ', 'B', ' ', ' ',
				' ', ' ', ' ', 'p', ' ', 'b', ' ', ' ',
				' ', ' ', ' ', 'b', 'p', ' ', ' ', 'n',
				'p', 'p', 'p', ' ', ' ', 'p', 'p', 'p',
				'r', 'n', ' ', 'q', 'k', ' ', ' ', 'r',
			},
			kingIndex: 60,
			player:    WHITE,
			expected:  false,
		},
		{
			name: "check blocked opponent piece",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', 'Q', ' ',
				' ', ' ', ' ', ' ', ' ', 'N', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			kingIndex: 60,
			player:    WHITE,
			expected:  false,
		},
		{
			name: "revealed check",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'N', ' ', 'Q', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			kingIndex: 60,
			player:    WHITE,
			expected:  true,
		},
	}

	for j, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			board := NewBoardFrom(input.board)
			board.turns = int(input.player)
			square := board.squares[input.kingIndex]
			assert.Equal(t, input.expected, board.inCheck(square),
				fmt.Sprintf("test %d", j))
		})

	}
}

func TestMoveInCheck(t *testing.T) {
	inputs := []struct {
		name       string
		board      []byte
		startIndex int
		destIndex  int
		player     Player // the player of current turn
		expected   bool   // whether the given move is valid
	}{
		{
			name: "basic queen move out of check",
			board: []byte{
				'K', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'Q', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'k', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startIndex: 56,
			destIndex:  57,
			player:     WHITE,
			expected:   true,
		},
		{
			name: "basic queen capture to avoid check",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				' ', ' ', ' ', 'P', 'P', 'P', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', 'n', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			startIndex: 12,
			destIndex:  21,
			player:     BLACK,
			expected:   true,
		},
		{
			name: "queen move in double check",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				' ', ' ', ' ', 'q', 'P', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', 'n', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			startIndex: 4,
			destIndex:  5,
			player:     BLACK,
			expected:   true,
		},
		{
			name: "basic queen move out of check INVALID",
			board: []byte{
				'K', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'Q', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'k', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startIndex: 56,
			destIndex:  48,
			player:     WHITE,
			expected:   false,
		},
		{
			name: "basic queen capture to avoid check INVALID",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				' ', ' ', ' ', 'P', 'P', 'P', ' ', ' ',
				' ', ' ', 'Q', ' ', ' ', 'n', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			startIndex: 11,
			destIndex:  18,
			player:     BLACK,
			expected:   false,
		},
		{
			name: "queen move in double check INVALID",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				' ', ' ', ' ', 'q', 'P', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', 'n', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			startIndex: 4,
			destIndex:  11,
			player:     BLACK,
			expected:   false,
		},
		{
			name: "cannot castle in check",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', 'R',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'r', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			startIndex: 4,
			destIndex:  7,
			player:     BLACK,
			expected:   false,
		},
		{
			name: "random valid move",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', ' ',
				'P', ' ', ' ', ' ', 'P', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', 'n', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			startIndex: 8,
			destIndex:  16,
			player:     BLACK,
			expected:   false,
		},
		{
			name: "cannot move into a check",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', 'R',
				' ', ' ', ' ', ' ', 'R', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', 'r', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			startIndex: 4,
			destIndex:  5,
			player:     BLACK,
			expected:   false,
		},
		{
			name: "cannot move expose your king",
			board: []byte{
				' ', ' ', ' ', ' ', 'K', ' ', ' ', 'R',
				' ', ' ', ' ', ' ', 'R', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'r', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
			},
			startIndex: 12,
			destIndex:  13,
			player:     BLACK,
			expected:   false,
		},
	}

	for j, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			board := NewBoardFrom(input.board)
			board.turns = int(input.player)
			start := board.squares[input.startIndex]
			dest := board.squares[input.destIndex]
			assert.Equal(t, input.expected, board.Move(start, dest),
				fmt.Sprintf("test %d", j))
		})

	}
}

func TestMate(t *testing.T) {
	inputs := []struct {
		name          string
		board         []byte
		player        Player // the player of current turn
		expectedCheck bool   // whether the given move is valid
		expectedMate  bool   // whether the given move is valid
	}{
		{
			name: "basic queen check, no mate",
			board: []byte{
				'K', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'Q', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'k', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			player:        WHITE,
			expectedCheck: true,
			expectedMate:  false,
		},
		{
			name: "ladder mate",
			board: []byte{
				'K', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'q', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'k', 'r', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			player:        BLACK,
			expectedCheck: true,
			expectedMate:  true,
		},
		{
			name: "night and pawn mate",
			board: []byte{
				'K', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'p', ' ', 'n', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'k', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			player:        BLACK,
			expectedCheck: true,
			expectedMate:  true,
		},
		{
			name: "night and bishop no mate",
			board: []byte{
				'K', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'b', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', 'n', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'k', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			player:        BLACK,
			expectedCheck: true,
			expectedMate:  false,
		},
		{
			name: "No chec",
			board: []byte{
				'K', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'P', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', 'b', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'k', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			player:        BLACK,
			expectedCheck: false,
			expectedMate:  false,
		},
		{
			name: "Stale mate",
			board: []byte{
				'K', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'P', 'P', 'b', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				'k', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			player:        BLACK,
			expectedCheck: false,
			expectedMate:  false,
		},
	}

	for j, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			board := NewBoardFrom(input.board)
			board.turns = int(input.player)
			check, mate := board.checkOrMate(input.player)
			assert.Equal(t, input.expectedCheck, check,
				fmt.Sprintf("test %d for check", j))
			assert.Equal(t, input.expectedMate, mate,
				fmt.Sprintf("test %d for mate", j))
		})

	}
}
