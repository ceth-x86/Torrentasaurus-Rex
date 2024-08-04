package message

import (
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	ErrInvalidMessageID   = errors.New("invalid message ID")
	ErrPayloadTooShort    = errors.New("payload too short")
	ErrPayloadLength      = errors.New("payload length payload")
	ErrIndexMismatch      = errors.New("index mismatch")
	ErrBeginOffsetTooHigh = errors.New("begin offset too high")
	ErrDataTooLong        = errors.New("data too long")
)

func validatePayloadLengthLess(minLength, actualLength int) error {
	if actualLength < minLength {
		return fmt.Errorf("%w: %d < %d", ErrPayloadTooShort, actualLength, minLength)
	}
	return nil
}

func validatePayloadLengthEqual(minLength, actualLength int) error {
	if actualLength != minLength {
		return fmt.Errorf("%w: %d != %d", ErrPayloadLength, actualLength, minLength)
	}
	return nil
}

func validateMessageID(expected, actual messageID) error {
	if actual != expected {
		return fmt.Errorf("%w: expected %d, got %d", ErrInvalidMessageID, expected, actual)
	}
	return nil
}

func validateIndex(expected, actual int) error {
	if expected != actual {
		return fmt.Errorf("%w: expected %d, got %d", ErrIndexMismatch, expected, actual)
	}
	return nil
}

func validateBeginOffset(begin, bufLen int) error {
	if begin >= bufLen {
		return fmt.Errorf("%w: %d >= %d", ErrBeginOffsetTooHigh, begin, bufLen)
	}
	return nil
}

func validateDataLength(begin, dataLen, bufLen int) error {
	if begin+dataLen > bufLen {
		return fmt.Errorf("%w: %d + %d > %d", ErrDataTooLong, begin, dataLen, bufLen)
	}
	return nil
}

func parsePiecePayload(payload []byte) (int, int, []byte, error) {
	if err := validatePayloadLengthLess(8, len(payload)); err != nil {
		return 0, 0, nil, err
	}

	index := int(binary.BigEndian.Uint32(payload[0:4]))
	begin := int(binary.BigEndian.Uint32(payload[4:8]))
	data := payload[8:]
	return index, begin, data, nil
}

// ParsePiece parses a PIECE message and copies its payload into a buffer
func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if err := validateMessageID(MsgPiece, msg.ID); err != nil {
		return 0, err
	}

	parsedIndex, begin, data, err := parsePiecePayload(msg.Payload)
	if err != nil {
		return 0, err
	}
	if err := validateIndex(index, parsedIndex); err != nil {
		return 0, err
	}
	if err := validateBeginOffset(begin, len(buf)); err != nil {
		return 0, err
	}

	if err := validateDataLength(begin, len(data), len(buf)); err != nil {
		return 0, err
	}

	copy(buf[begin:], data)
	return len(data), nil
}

// ParseHave parses a HAVE message
func ParseHave(msg *Message) (int, error) {
	if err := validateMessageID(MsgHave, msg.ID); err != nil {
		return 0, err
	}
	if err := validatePayloadLengthEqual(4, len(msg.Payload)); err != nil {
		return 0, err
	}
	index := int(binary.BigEndian.Uint32(msg.Payload))
	return index, nil
}
