package message

import (
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestSerializeKeepAlive tests serialization of a keep-alive message
func TestSerializeKeepAlive(t *testing.T) {
	var m *Message // nil message
	expected := createKeepAliveMessage()
	actual := m.Serialize()
	assert.Equal(t, expected, actual, "Serialized keep-alive message should be equal to expected keep-alive message")
}

// TestCreateKeepAliveMessage tests creation of a keep-alive message
func TestCreateKeepAliveMessage(t *testing.T) {
	expected := make([]byte, 4)
	actual := createKeepAliveMessage()
	assert.Equal(t, expected, actual, "Keep-alive message should be a 4-byte slice")
}

// TestSerializeMessage tests serialization of a non-nil message
func TestSerializeMessage(t *testing.T) {
	message := &Message{
		ID:      1,
		Payload: []byte{0x01, 0x02, 0x03},
	}

	expectedLength := uint32(4 + 1 + len(message.Payload)) // 4 bytes length prefix, 1 byte ID, payload length
	expected := make([]byte, 4+1+len(message.Payload))
	binary.BigEndian.PutUint32(expected[0:4], expectedLength-4)
	expected[4] = byte(message.ID)
	copy(expected[5:], message.Payload)

	actual := message.Serialize()
	assert.Equal(t, expected, actual, "Serialized message should be equal to expected serialized message")
}

// TestSerializeEmptyPayload tests serialization of a message with an empty payload
func TestSerializeEmptyPayload(t *testing.T) {
	message := &Message{
		ID:      1,
		Payload: []byte{},
	}

	expectedLength := uint32(5) // 4 bytes length prefix, 1 byte ID
	expected := make([]byte, 5)
	binary.BigEndian.PutUint32(expected[0:4], expectedLength-4)
	expected[4] = byte(message.ID)

	actual := message.Serialize()
	assert.Equal(t, expected, actual, "Serialized message with empty payload should be equal to expected serialized message")
}

// TestSerializeMessageWithNilPayload tests serialization of a message with a nil payload
func TestSerializeMessageWithNilPayload(t *testing.T) {
	message := &Message{
		ID:      1,
		Payload: nil,
	}

	expectedLength := uint32(5) // 4 bytes length prefix, 1 byte ID
	expected := make([]byte, 5)
	binary.BigEndian.PutUint32(expected[0:4], expectedLength-4)
	expected[4] = byte(message.ID)

	actual := message.Serialize()
	assert.Equal(t, expected, actual, "Serialized message with nil payload should be equal to expected serialized message")
}
