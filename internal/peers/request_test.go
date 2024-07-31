package peers

import (
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackpal/bencode-go"
	"github.com/stretchr/testify/assert"
)

// Mock response structure
type mockTrackerResponse struct {
	Peers string `bencode:"peers"`
}

func TestRequestSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := mockTrackerResponse{
			Peers: string([]byte{192, 168, 0, 1, 0x1A, 0xE1, 10, 0, 0, 2, 0x1A, 0xE1}),
		}
		w.WriteHeader(http.StatusOK)
		err := bencode.Marshal(w, response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	expectedPeers := []Peer{
		{IP: net.IPv4(192, 168, 0, 1).To4(), Port: 6881},
		{IP: net.IPv4(10, 0, 0, 2).To4(), Port: 6881},
	}

	peers, err := Request(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, expectedPeers, peers)
}

func TestRequestFailure(t *testing.T) {
	// Create a mock server that returns an error status
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	}))
	defer server.Close()

	peers, err := Request(server.URL)
	assert.Error(t, err)
	assert.Nil(t, peers)
}

func TestRequestMalformedResponse(t *testing.T) {
	// Create a mock server that returns malformed data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("malformed data"))
	}))
	defer server.Close()

	peers, err := Request(server.URL)
	assert.Error(t, err)
	assert.Nil(t, peers)
}
