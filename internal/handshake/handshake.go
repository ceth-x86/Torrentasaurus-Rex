// TODO: refactore
// TODO: tests (don't forget about completeHandshake)

package handshake

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"time"
)

const (
	ProtocolName      = "BitTorrent protocol"
	ReservedBytesSize = 8
	InfoHashSize      = 20
	PeerIDSize        = 20
	FixedHeaderSize   = ReservedBytesSize + InfoHashSize + PeerIDSize
)

// A Handshake is a special message that a peer uses to identify itself
type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

// CompleteHandshake performs the handshake process with the peer
func CompleteHandshake(conn net.Conn, infohash, peerID [InfoHashSize]byte) (*Handshake, error) {
	conn.SetDeadline(time.Now().Add(3 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	req := &Handshake{
		Pstr:     ProtocolName,
		InfoHash: infohash,
		PeerID:   peerID,
	}
	_, err := conn.Write(req.serialize())
	if err != nil {
		return nil, fmt.Errorf("failed to write handshake: %w", err)
	}

	res, err := read(conn)
	if err != nil {
		return nil, fmt.Errorf("failed to read handshake: %w", err)
	}
	if !bytes.Equal(res.InfoHash[:], infohash[:]) {
		return nil, fmt.Errorf("expected infohash %x but got %x", res.InfoHash, infohash)
	}
	return res, nil
}

// serialize serializes the handshake to a buffer
func (h *Handshake) serialize() []byte {
	buf := make([]byte, len(h.Pstr)+1+FixedHeaderSize)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], make([]byte, ReservedBytesSize))
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

// read parses a handshake from a stream
func read(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, fmt.Errorf("failed to read pstrlen: %w", err)
	}

	pstrlen := int(lengthBuf[0])
	if pstrlen == 0 {
		return nil, fmt.Errorf("pstrlen cannot be 0")
	}

	handshakeBuf := make([]byte, FixedHeaderSize+pstrlen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		return nil, err
	}

	var infoHash, peerID [InfoHashSize]byte

	copy(infoHash[:], handshakeBuf[pstrlen+ReservedBytesSize:pstrlen+ReservedBytesSize+InfoHashSize])
	copy(peerID[:], handshakeBuf[pstrlen+ReservedBytesSize+InfoHashSize:])

	h := Handshake{
		Pstr:     string(handshakeBuf[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return &h, nil
}
