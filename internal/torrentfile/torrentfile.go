package torrentfile

import (
	"fmt"
	"github.com/jackpal/bencode-go"
	"os"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

type BencodeTorrentFile struct {
	Announce string      `bencode:"announce"`
	Info     bencodeInfo `bencode:"info"`
}

func Open(path string) (TorrentFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return TorrentFile{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	btf := &BencodeTorrentFile{}
	if err := bencode.Unmarshal(file, &btf); err != nil {
		return TorrentFile{}, fmt.Errorf("failed to unmarshal bencode: %w", err)
	}

	return btf.toTorrentFile()
}

func (btf *BencodeTorrentFile) toTorrentFile() (TorrentFile, error) {
	infoHash, err := btf.Info.hash()
	if err != nil {
		return TorrentFile{}, fmt.Errorf("failed to hash info: %w", err)
	}

	pieceHashes, err := btf.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, fmt.Errorf("failed to split piece hashes: %w", err)
	}

	return TorrentFile{
		Announce:    btf.Announce,
		InfoHash:    infoHash,
		PieceHashes: pieceHashes,
		PieceLength: btf.Info.PieceLength,
		Length:      btf.Info.Length,
		Name:        btf.Info.Name,
	}, nil
}
