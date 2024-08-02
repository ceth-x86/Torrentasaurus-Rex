package bitfields

import (
	"Torrentasaurus_Rex/internal/message"
	"fmt"
	"net"
	"time"
)

// A Bitfield represents the pieces that a peer has
type Bitfield []byte

// HasPiece tells if a bitfields has a particular index set
func (bf Bitfield) HasPiece(index int) bool {
	if !bf.isValidIndex(index) {
		return false
	}
	byteIndex, offset := bf.byteAndOffset(index)
	return bf[byteIndex]>>uint(7-offset)&1 != 0
}

// SetPiece sets a bit in the bitfields
func (bf Bitfield) SetPiece(index int) {
	if !bf.isValidIndex(index) {
		return
	}
	byteIndex, offset := bf.byteAndOffset(index)
	bf[byteIndex] |= 1 << uint(7-offset)
}

// RecvBitfield receives a bitfield from a connection
func RecvBitfield(conn net.Conn) (Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	msg, err := message.Read(conn)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, fmt.Errorf("expected bitfields but got %v", msg)
	}
	if msg.ID != message.MsgBitfield {
		return nil, fmt.Errorf("expected bitfields but got ID %d", msg.ID)
	}

	return Bitfield(msg.Payload), nil
}

// isValidIndex checks if the index is within the bounds of the bitfield
func (bf Bitfield) isValidIndex(index int) bool {
	byteIndex := index / 8
	return byteIndex >= 0 && byteIndex < len(bf)
}

// byteAndOffset returns the byte index and bit offset for a given piece index
func (bf Bitfield) byteAndOffset(index int) (int, int) {
	return index / 8, index % 8
}
