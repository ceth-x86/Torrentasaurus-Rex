package exchange

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateBoundsForPiece(t *testing.T) {
	tests := []struct {
		name          string
		pieceLength   int
		length        int
		index         int
		expectedBegin int
		expectedEnd   int
	}{
		{"Normal case", 100, 1000, 3, 300, 400},
		{"Last piece", 100, 1000, 9, 900, 1000},
		{"Partial piece", 100, 950, 9, 900, 950},
		{"First piece", 100, 1000, 0, 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exchange{
				PieceLength: tt.pieceLength,
				Length:      tt.length,
			}
			begin, end := e.calculateBoundsForPiece(tt.index)
			assert.Equal(t, tt.expectedBegin, begin)
			assert.Equal(t, tt.expectedEnd, end)
		})
	}
}

func TestCalculatePieceSize(t *testing.T) {
	tests := []struct {
		name         string
		pieceLength  int
		length       int
		index        int
		expectedSize int
	}{
		{"Normal case", 100, 1000, 3, 100},
		{"Last piece", 100, 1000, 9, 100},
		{"Partial piece", 100, 950, 9, 50},
		{"First piece", 100, 1000, 0, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Exchange{
				PieceLength: tt.pieceLength,
				Length:      tt.length,
			}
			size := e.calculatePieceSize(tt.index)
			assert.Equal(t, tt.expectedSize, size)
		})
	}
}
