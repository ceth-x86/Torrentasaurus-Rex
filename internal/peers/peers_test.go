package peers

import (
	"encoding/binary"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPeerString(t *testing.T) {
	peer := Peer{
		IP:   net.IPv4(192, 168, 0, 1),
		Port: 6881,
	}
	expected := "192.168.0.1:6881"
	assert.Equal(t, expected, peer.String())
}

func TestUnmarshal(t *testing.T) {
	// Test with valid input
	peersBin := []byte{
		192, 168, 0, 1, 0x1A, 0xE1, // 192.168.0.1:6881
		10, 0, 0, 2, 0x1A, 0xE1, // 10.0.0.2:6881
	}

	expected := []Peer{
		{IP: net.IPv4(192, 168, 0, 1).To4(), Port: 6881},
		{IP: net.IPv4(10, 0, 0, 2).To4(), Port: 6881},
	}

	result, err := Unmarshal(peersBin)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	// Test with malformed input (length not a multiple of peerSize)
	peersBinMalformed := []byte{
		192, 168, 0, 1, 0x1A, // Malformed data
	}

	result, err = Unmarshal(peersBinMalformed)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestUnmarshalEmptyInput(t *testing.T) {
	// Test with empty input
	peersBin := []byte{}

	expected := []Peer{}

	result, err := Unmarshal(peersBin)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUnmarshalSinglePeer(t *testing.T) {
	// Test with single peer
	peersBin := make([]byte, 6)
	copy(peersBin[0:4], net.IPv4(127, 0, 0, 1).To4())
	binary.BigEndian.PutUint16(peersBin[4:6], 6881)

	expected := []Peer{
		{IP: net.IPv4(127, 0, 0, 1).To4(), Port: 6881},
	}

	result, err := Unmarshal(peersBin)
	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}
