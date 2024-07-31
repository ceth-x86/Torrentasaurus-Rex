package peers

import (
	"fmt"
	"net/http"
	"time"

	"Torrentasaurus_Rex/internal/tracker"

	"github.com/jackpal/bencode-go"
)

func Request(url string) ([]Peer, error) {
	client := &http.Client{Timeout: 15 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch URL: %s, status code: %d", url, resp.StatusCode)
	}

	trackerResponse := tracker.BencodeTrackerResponse{}
	err = bencode.Unmarshal(resp.Body, &trackerResponse)
	if err != nil {
		return nil, err
	}

	return Unmarshal([]byte(trackerResponse.Peers))
}
