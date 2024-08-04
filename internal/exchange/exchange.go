package exchange

import "Torrentasaurus_Rex/internal/peers"

// Exchange holds data required to download a torrent from a list of peers
type Exchange struct {
	Peers       []peers.Peer
	PeerID      [20]byte
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}
