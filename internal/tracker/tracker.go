package tracker

import (
	"Torrentasaurus_Rex/internal/torrent"
	"fmt"
	"net/url"
	"strconv"
)

// Port to listen on
const Port uint16 = 6881

type BencodeTrackerResponse struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

func BuildTrackerURL(tf *torrent.TorrentFile, peerID [20]byte) (string, error) {
	base, err := url.Parse(tf.Announce)
	if err != nil {
		return "", fmt.Errorf("failed to parse announce URL: %w", err)
	}
	params := url.Values{
		"info_hash":  {string(tf.InfoHash[:])},
		"peer_id":    {string(peerID[:])},
		"port":       {strconv.Itoa(int(Port))},
		"uploaded":   {"0"},
		"downloaded": {"0"},
		"compact":    {"1"},
		"left":       {strconv.Itoa(tf.Length)},
	}
	base.RawQuery = params.Encode()
	return base.String(), nil
}
