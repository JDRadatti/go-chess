package chess

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRankAndFile(t *testing.T) {
	for i := 0; i < int(HEIGHT); i++ {
		for j := 0; j < int(WIDTH); j++ {
			square := Square{index: i*int(WIDTH) + j}
			assert.Equal(t, i, square.rank(), "test file")
			assert.Equal(t, j, square.file(), "test rank")
		}
	}
}
