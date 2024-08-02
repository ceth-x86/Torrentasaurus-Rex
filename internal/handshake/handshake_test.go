package handshake

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestSerialize(t *testing.T) {
	infoHash := [InfoHashSize]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	peerID := [PeerIDSize]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}

	handshake := &Handshake{
		Pstr:     ProtocolName,
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	expected := []byte{
		byte(len(ProtocolName)),                                                                       // Pstrlen
		'B', 'i', 't', 'T', 'o', 'r', 'r', 'e', 'n', 't', ' ', 'p', 'r', 'o', 't', 'o', 'c', 'o', 'l', // Pstr
		0, 0, 0, 0, 0, 0, 0, 0, // Reserved bytes
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, // InfoHash
		21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, // PeerID
	}

	result := handshake.serialize()

	if !bytes.Equal(result, expected) {
		t.Errorf("Expected %x, but got %x", expected, result)
	}
}

func TestRead(t *testing.T) {
	infoHash := [InfoHashSize]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}
	peerID := [PeerIDSize]byte{21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40}

	handshake := &Handshake{
		Pstr:     ProtocolName,
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	serialized := handshake.serialize()
	reader := bytes.NewReader(serialized)

	result, err := read(reader)
	if err != nil {
		t.Fatalf("Expected no error, but got %v", err)
	}

	if result.Pstr != ProtocolName {
		t.Errorf("Expected Pstr %s, but got %s", ProtocolName, result.Pstr)
	}
	if !bytes.Equal(result.InfoHash[:], infoHash[:]) {
		t.Errorf("Expected InfoHash %x, but got %x", infoHash, result.InfoHash)
	}
	if !bytes.Equal(result.PeerID[:], peerID[:]) {
		t.Errorf("Expected PeerID %x, but got %x", peerID, result.PeerID)
	}
}

func createClientAndServer(t *testing.T) (clientConn, serverConn net.Conn) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	require.Nil(t, err)

	// net.Dial does not block, so we need this signalling channel to make sure
	// we don't return before serverConn is ready
	done := make(chan struct{})
	go func() {
		defer ln.Close()
		serverConn, err = ln.Accept()
		require.Nil(t, err)
		done <- struct{}{}
	}()
	clientConn, err = net.Dial("tcp", ln.Addr().String())
	<-done

	return clientConn, serverConn
}

func TestCompleteHandshake(t *testing.T) {
	tests := map[string]struct {
		clientInfohash  [20]byte
		clientPeerID    [20]byte
		serverHandshake []byte
		output          *Handshake
		fails           bool
	}{
		"successful handshake": {
			clientInfohash:  [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
			clientPeerID:    [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			serverHandshake: []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0, 134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116, 45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
			output: &Handshake{
				Pstr:     ProtocolName,
				InfoHash: [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
				PeerID:   [20]byte{45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
			},
			fails: false,
		},
		"wrong infohash": {
			clientInfohash:  [20]byte{134, 212, 200, 0, 36, 164, 105, 190, 76, 80, 188, 90, 16, 44, 247, 23, 128, 49, 0, 116},
			clientPeerID:    [20]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20},
			serverHandshake: []byte{19, 66, 105, 116, 84, 111, 114, 114, 101, 110, 116, 32, 112, 114, 111, 116, 111, 99, 111, 108, 0, 0, 0, 0, 0, 0, 0, 0, 0xde, 0xe8, 0x6a, 0x7f, 0xa6, 0xf2, 0x86, 0xa9, 0xd7, 0x4c, 0x36, 0x20, 0x14, 0x61, 0x6a, 0x0f, 0xf5, 0xe4, 0x84, 0x3d, 45, 83, 89, 48, 48, 49, 48, 45, 192, 125, 147, 203, 136, 32, 59, 180, 253, 168, 193, 19},
			output:          nil,
			fails:           true,
		},
	}

	for _, test := range tests {
		clientConn, serverConn := createClientAndServer(t)
		serverConn.Write(test.serverHandshake)

		h, err := CompleteHandshake(clientConn, test.clientInfohash, test.clientPeerID)

		if test.fails {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, h, test.output)
		}
	}
}
