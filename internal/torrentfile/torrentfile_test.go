package torrentfile

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestOpen(t *testing.T) {
	torrentFile := "testdata/ubuntu-24.04-desktop-amd64.iso.torrent"
	goldenPath := "testdata/ubuntu-24.04-desktop-amd64.iso.torrent.golden.json"

	torrent, err := Open(torrentFile)
	require.Nil(t, err)

	expected := TorrentFile{}
	golden, err := os.ReadFile(goldenPath)
	require.Nil(t, err)
	err = json.Unmarshal(golden, &expected)
	require.Nil(t, err)

	assert.Equal(t, expected, torrent)
}
