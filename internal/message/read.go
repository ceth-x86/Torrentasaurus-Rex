package message

import (
	"encoding/binary"
	"fmt"
	"io"
)

// Read parses a message from a stream. Returns `nil` on keep-alive message
func Read(r io.Reader) (*Message, error) {
	length, err := readMessageLength(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read message length: %w", err)
	}

	// keep-alive message
	if length == 0 {
		return nil, nil
	}

	if length < 1 {
		return nil, fmt.Errorf("invalid message length: %d", length)
	}

	message, err := readMessageBody(r, length)
	if err != nil {
		return nil, fmt.Errorf("failed to read message body: %w", err)
	}

	return message, nil
}

func readMessageLength(r io.Reader) (uint32, error) {
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lengthBuf); err != nil {
		return 0, fmt.Errorf("failed to read message length: %w", err)
	}
	return binary.BigEndian.Uint32(lengthBuf), nil
}

func readMessageBody(r io.Reader, length uint32) (*Message, error) {
	messageBuf := make([]byte, length)
	if _, err := io.ReadFull(r, messageBuf); err != nil {
		return nil, fmt.Errorf("failed to read message body: %w", err)
	}

	return &Message{
		ID:      messageID(messageBuf[0]),
		Payload: messageBuf[1:],
	}, nil
}
