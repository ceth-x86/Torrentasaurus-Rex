package message

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidatePayloadLengthLess(t *testing.T) {
	tests := []struct {
		name         string
		minLength    int
		actualLength int
		expectedErr  error
	}{
		{"ValidLength", 5, 10, nil},
		{"InvalidLength", 10, 5, fmt.Errorf("%w: %d < %d", ErrPayloadTooShort, 5, 10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePayloadLengthLess(tt.minLength, tt.actualLength)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidatePayloadLengthEqual(t *testing.T) {
	tests := []struct {
		name         string
		minLength    int
		actualLength int
		expectedErr  error
	}{
		{"ValidLength", 5, 5, nil},
		{"InvalidLength", 5, 10, fmt.Errorf("%w: %d != %d", ErrPayloadLength, 10, 5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validatePayloadLengthEqual(tt.minLength, tt.actualLength)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateMessageID(t *testing.T) {
	tests := []struct {
		name        string
		expected    messageID
		actual      messageID
		expectedErr error
	}{
		{"ValidMessageID", MsgPiece, MsgPiece, nil},
		{"InvalidMessageID", MsgPiece, MsgHave, fmt.Errorf("%w: expected %d, got %d", ErrInvalidMessageID, MsgPiece, MsgHave)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMessageID(tt.expected, tt.actual)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateIndex(t *testing.T) {
	tests := []struct {
		name        string
		expected    int
		actual      int
		expectedErr error
	}{
		{"ValidIndex", 5, 5, nil},
		{"InvalidIndex", 5, 10, fmt.Errorf("%w: expected %d, got %d", ErrIndexMismatch, 5, 10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateIndex(tt.expected, tt.actual)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateBeginOffset(t *testing.T) {
	tests := []struct {
		name        string
		begin       int
		bufLen      int
		expectedErr error
	}{
		{"ValidBeginOffset", 5, 10, nil},
		{"InvalidBeginOffset", 10, 5, fmt.Errorf("%w: %d >= %d", ErrBeginOffsetTooHigh, 10, 5)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateBeginOffset(tt.begin, tt.bufLen)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestValidateDataLength(t *testing.T) {
	tests := []struct {
		name        string
		begin       int
		dataLen     int
		bufLen      int
		expectedErr error
	}{
		{"ValidDataLength", 5, 5, 10, nil},
		{"InvalidDataLength", 5, 10, 10, fmt.Errorf("%w: %d + %d > %d", ErrDataTooLong, 5, 10, 10)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateDataLength(tt.begin, tt.dataLen, tt.bufLen)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestParsePiecePayload(t *testing.T) {
	tests := []struct {
		name          string
		payload       []byte
		expectedIndex int
		expectedBegin int
		expectedData  []byte
		expectedErr   error
	}{
		{
			"ValidPayload",
			func() []byte {
				buf := make([]byte, 12)
				binary.BigEndian.PutUint32(buf[0:4], uint32(1))
				binary.BigEndian.PutUint32(buf[4:8], uint32(2))
				copy(buf[8:], []byte("data"))
				return buf
			}(),
			1,
			2,
			[]byte("data"),
			nil,
		},
		{
			"InvalidPayloadLength",
			[]byte{1, 2, 3, 4, 5, 6},
			0,
			0,
			nil,
			fmt.Errorf("%w: %d < %d", ErrPayloadTooShort, 6, 8),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, begin, data, err := parsePiecePayload(tt.payload)
			assert.Equal(t, tt.expectedIndex, index)
			assert.Equal(t, tt.expectedBegin, begin)
			assert.Equal(t, tt.expectedData, data)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestParsePiece(t *testing.T) {
	tests := []struct {
		name        string
		index       int
		buf         []byte
		msg         *Message
		expectedLen int
		expectedErr error
	}{
		{
			"ValidPiece",
			1,
			make([]byte, 10),
			&Message{
				ID: MsgPiece,
				Payload: func() []byte {
					buf := make([]byte, 12)
					binary.BigEndian.PutUint32(buf[0:4], uint32(1))
					binary.BigEndian.PutUint32(buf[4:8], uint32(2))
					copy(buf[8:], []byte("data"))
					return buf
				}(),
			},
			4,
			nil,
		},
		{
			"InvalidMessageID",
			1,
			make([]byte, 10),
			&Message{ID: MsgHave, Payload: []byte{1, 2, 3, 4}},
			0,
			fmt.Errorf("%w: expected %d, got %d", ErrInvalidMessageID, MsgPiece, MsgHave),
		},
		{
			"InvalidPayload",
			1,
			make([]byte, 10),
			&Message{
				ID:      MsgPiece,
				Payload: []byte{1, 2, 3, 4, 5, 6},
			},
			0,
			fmt.Errorf("%w: %d < %d", ErrPayloadTooShort, 6, 8),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n, err := ParsePiece(tt.index, tt.buf, tt.msg)
			assert.Equal(t, tt.expectedLen, n)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}

func TestParseHave(t *testing.T) {
	tests := []struct {
		name          string
		msg           *Message
		expectedIndex int
		expectedErr   error
	}{
		{
			"ValidHave",
			&Message{
				ID: MsgHave,
				Payload: func() []byte {
					buf := make([]byte, 4)
					binary.BigEndian.PutUint32(buf, uint32(1))
					return buf
				}(),
			},
			1,
			nil,
		},
		{
			"InvalidMessageID",
			&Message{ID: MsgPiece, Payload: []byte{1, 2, 3, 4}},
			0,
			fmt.Errorf("%w: expected %d, got %d", ErrInvalidMessageID, MsgHave, MsgPiece),
		},
		{
			"InvalidPayloadLength",
			&Message{ID: MsgHave, Payload: []byte{1, 2, 3}},
			0,
			fmt.Errorf("%w: %d != %d", ErrPayloadLength, 3, 4),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			index, err := ParseHave(tt.msg)
			assert.Equal(t, tt.expectedIndex, index)
			assert.Equal(t, tt.expectedErr, err)
		})
	}
}
