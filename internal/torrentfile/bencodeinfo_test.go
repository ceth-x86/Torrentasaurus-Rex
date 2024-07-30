package torrentfile

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBencodeInfo_Hash(t *testing.T) {
	info := bencodeInfo{
		Pieces:      "12345678901234567890",
		PieceLength: 262144,
		Length:      123456,
		Name:        "testfile.txt",
	}
	hash, err := info.hash()
	require.NoError(t, err)

	expectedHash := [20]uint8{0xa4, 0x76, 0xf7, 0xce, 0xf2, 0xa4, 0x2b, 0x54, 0xf2, 0xc9, 0x7d, 0x49, 0x3a, 0x4a, 0x4b, 0xfb, 0xf7, 0xb3, 0x21, 0x19}
	assert.Equal(t, expectedHash, hash)
}

func TestBencodeInfo_SplitPieceHashes(t *testing.T) {
	info := bencodeInfo{
		Pieces: "1234567890123456789012345678901234567890",
	}
	hashes, err := info.splitPieceHashes()
	require.NoError(t, err)

	expectedHashes := [][20]byte{
		{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30},
		{0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30, 0x31, 0x32, 0x33, 0x34, 0x35, 0x36, 0x37, 0x38, 0x39, 0x30},
	}
	assert.Equal(t, expectedHashes, hashes)
}
