package chess

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromAlgebraicValid(t *testing.T) {
	board := NewBoardFrom([]byte{
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	})
	width, height := 7, 8
	index := 0
	destSquare, destIndex := "a8", 0
	for j := width; j >= 0; j-- {
		rank := '1' + j
		for i := 0; i < height; i++ {
			file := 'a' + i
			square := fmt.Sprintf("%c", file) + fmt.Sprintf("%c", rank) + destSquare
			s1, s2 := board.fromAlgebraic(square)
			assert.Equal(t, index, s1.index)
			assert.Equal(t, destIndex, s2.index)
			index++
		}
	}

	startSquare, startIndex := "a8", 0
	index = 0
	for j := width; j >= 0; j-- {
		rank := '1' + j
		for i := 0; i < height; i++ {
			file := 'a' + i
			square := startSquare + fmt.Sprintf("%c", file) + fmt.Sprintf("%c", rank)
			s1, s2 := board.fromAlgebraic(square)
			assert.Equal(t, index, s2.index)
			assert.Equal(t, startIndex, s1.index)
			index++
		}
	}
}

func TestFromAlgebraicInvalid(t *testing.T) {
	board := NewBoardFrom([]byte{
		'R', 'N', 'B', 'Q', 'K', 'B', 'N', 'R',
		'P', 'P', 'P', 'P', 'P', 'P', 'P', 'P',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		' ', ' ', ' ', ' ', ' ', ' ', ' ', ' ',
		'p', 'p', 'p', 'p', 'p', 'p', 'p', 'p',
		'r', 'n', 'b', 'q', 'k', 'b', 'n', 'r',
	})
	inputs := []string{
		"a9b2",
		"abcd",
		"1234",
		"ijkl",
		"i9j0",
		"a8",
		"",
		"*!@#",
		"--",
        "a1b2c3",
	}
	for _, input := range inputs {
		s1, s2 := board.fromAlgebraic(input)
        var expected *Square
		assert.Equal(t, expected, s1)
		assert.Equal(t, expected, s2)
	}
}

func TestToAlgebraic(t *testing.T) {
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
