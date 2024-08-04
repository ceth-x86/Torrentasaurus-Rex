package message

import "encoding/binary"

// Serialize serializes a message into a buffer of the form
// <length prefix><message ID><payload>
// Interprets `nil` as a keep-alive message
func (m *Message) Serialize() []byte {
	if m == nil {
		return createKeepAliveMessage()
	}
	return serializeMessage(m)
}

// createKeepAliveMessage creates a keep-alive message
func createKeepAliveMessage() []byte {
	return make([]byte, 4)
}

// serializeMessage serializes a non-nil message
func serializeMessage(m *Message) []byte {
	const headerSize = 4
	const idSize = 1

	length := uint32(len(m.Payload) + idSize) // +1 for ID
	buf := make([]byte, headerSize+length)
	binary.BigEndian.PutUint32(buf[0:headerSize], length)
	buf[headerSize] = byte(m.ID)
	copy(buf[headerSize+idSize:], m.Payload)
	return buf
}
