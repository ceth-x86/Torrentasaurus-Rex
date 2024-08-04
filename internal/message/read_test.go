package message

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadKeepAliveMessage(t *testing.T) {
	data := []byte{0, 0, 0, 0} // keep-alive message
	r := bytes.NewReader(data)

	message, err := Read(r)
	assert.NoError(t, err)
	assert.Nil(t, message)
}

func TestReadMessage(t *testing.T) {
	data := []byte{0, 0, 0, 5, 1, 2, 3, 4, 5} // valid message
	r := bytes.NewReader(data)

	message, err := Read(r)
	require.NoError(t, err)
	require.NotNil(t, message)
	assert.Equal(t, messageID(1), message.ID)
	assert.Equal(t, []byte{2, 3, 4, 5}, message.Payload)
}

func TestReadMessageWithShortLength(t *testing.T) {
	data := []byte{0, 0, 0, 5, 1, 2, 3} // short length
	r := bytes.NewReader(data)

	message, err := Read(r)
	assert.Error(t, err)
	assert.Nil(t, message)
	assert.Contains(t, err.Error(), "failed to read message body")
}
