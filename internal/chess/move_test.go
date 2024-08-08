package chess

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestToAlgabraic(t *testing.T) {
	inputs := []struct {
		name         string
		board        []byte
		startSquares int // index of square in board.squares
		destSquares  int // index of square in board.squares
		expected     string
		player       Player
		castle       bool
	}{
		{
			name: "a4",
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
			startSquares: 48,
			destSquares:  40,
			expected:     "a3",
		},
		{
			name: "d5",
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
			startSquares: 11,
			destSquares:  27,
			expected:     "d5",
		},
		{
			name: "Na6",
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
			startSquares: 1,
			destSquares:  16,
			expected:     "Na6",
		},
		{
			name: "Knight edge case",
			board: []byte{
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'n', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', 'n', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 17,
			destSquares:  27,
			expected:     "Nbd5",
		},
		{
			name: "Knight edge case",
			board: []byte{
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'n', ' ', ' ', ' ', 'n', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 17,
			destSquares:  27,
			expected:     "Nbd5",
		},
		{
			name: "Knight edge case",
			board: []byte{
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'n', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'n', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 17,
			destSquares:  27,
			expected:     "N6d5",
		},
		{
			name: "Knight edge case",
			board: []byte{
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'n', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', 'P', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', 'n', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 17,
			destSquares:  27,
			expected:     "Nbxd5",
		},
		{
			name: "Rook edge case",
			board: []byte{
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'r', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'r', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 17,
			destSquares:  25,
			expected:     "R6b5",
		},
		{
			name: "Rook edge case",
			board: []byte{
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'r', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', 'r', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 41,
			destSquares:  25,
			expected:     "R3b5",
		},
		{
			name: "Rook edge case",
			board: []byte{
				'r', ' ', ' ', ' ', 'r', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 0,
			destSquares:  1,
			expected:     "Rab8",
		},
		{
			name: "Kingside castle",
			board: []byte{
				'r', ' ', ' ', ' ', 'k', ' ', ' ', 'r',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 4,
			castle:       true,
			destSquares:  7,
			expected:     "O-O",
		},
		{
			name: "Queenside castle",
			board: []byte{
				'r', ' ', ' ', ' ', 'k', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
				' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
			},
			startSquares: 4,
			castle:       true,
			destSquares:  0,
			expected:     "O-O-O",
		},
	}

	for j, input := range inputs {
		t.Run(input.name, func(t *testing.T) {
			board := NewBoardFrom(input.board) // reset board every time
			start := board.squares[input.startSquares]
			dest := board.squares[input.destSquares]
			move := Move{
				startSquare1: start,
				destSquare1:  dest,
				piece1:       start.piece,
				startSquare2: dest,
				destSquare2:  nil,
				piece2:       dest.piece,
				castle:       input.castle,
			}
			assert.Equal(t, input.expected, move.toAlgebraic(&board),
				fmt.Sprintf("test %d", j))

		})
	}
}
