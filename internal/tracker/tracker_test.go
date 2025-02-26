package tracker

import (
	"Torrentasaurus_Rex/internal/torrent"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBuildTrackerURL(t *testing.T) {
	tf := &torrent.TorrentFile{
		Announce: "http://example.com/announce",
		InfoHash: [20]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78},
		Length:   12345,
	}
	peerID := [20]byte{0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01}

	expectedURL := "http://example.com/announce?compact=1&downloaded=0&info_hash=%124Vx%9A%BC%DE%F0%124Vx%9A%BC%DE%F0%124Vx&left=12345&peer_id=%AB%CD%EF%01%23Eg%89%AB%CD%EF%01%23Eg%89%AB%CD%EF%01&port=6881&uploaded=0"

	resultURL, err := BuildTrackerURL(tf, peerID)
	require.NoError(t, err)
	assert.Equal(t, expectedURL, resultURL)
}

func TestBuildTrackerURLError(t *testing.T) {
	tf := &torrent.TorrentFile{
		Announce: "://bad_url",
		InfoHash: [20]byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12, 0x34, 0x56, 0x78},
		Length:   12345,
	}
	peerID := [20]byte{0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0x01}

	resultURL, err := BuildTrackerURL(tf, peerID)
	require.Error(t, err)
	assert.Empty(t, resultURL)
}
