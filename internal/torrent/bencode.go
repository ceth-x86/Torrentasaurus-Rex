package torrent

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"github.com/jackpal/bencode-go"
)

// bencodeInfo represents the information contained in the "info" section of the bencode torrent file.
type bencodeInfo struct {
	Pieces      string `bencode:"pieces"`
	PieceLength int    `bencode:"piece length"`
	Length      int    `bencode:"length"`
	Name        string `bencode:"name"`
}

type BencodeTorrentFile struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

// hash computes the SHA-1 hash for BencodeInfo.
func (i *bencodeInfo) hash() ([20]byte, error) {
	var buf bytes.Buffer
	if err := bencode.Marshal(&buf, *i); err != nil {
		return [20]byte{}, fmt.Errorf("failed to marshal bencode info: %w", err)
	}
	return sha1.Sum(buf.Bytes()), nil
}

// splitPieceHashes splits the pieces string into individual SHA-1 hashes.
func (i *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	const hashLen = 20
	buf := []byte(i.Pieces)
	if len(buf)%hashLen != 0 {
		return nil, fmt.Errorf("incorrect data of length %d", len(buf))
	}

	numHashes := len(buf) / hashLen
	hashes := make([][20]byte, numHashes)
	for i := 0; i < numHashes; i++ {
		copy(hashes[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return hashes, nil
}
