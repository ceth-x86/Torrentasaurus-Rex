package peers

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}

func Unmarshal(peersBin []byte) ([]Peer, error) {
	const (
		peerSize = 6 // 4 bytes for IP, 2 bytes for port
		ipSize   = 4
	)

	if len(peersBin)%peerSize != 0 {
		return nil, errors.New("received malformed peers")
	}

	numPeers := len(peersBin) / peerSize
	peers := make([]Peer, numPeers)

	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		ipBytes := peersBin[offset : offset+ipSize]
		portBytes := peersBin[offset+ipSize : offset+peerSize]

		peers[i] = Peer{
			IP:   ipBytes,
			Port: binary.BigEndian.Uint16(portBytes),
		}
	}
	return peers, nil
}

func GeneratePeerID() ([20]byte, error) {
	var peerID [20]byte

	if _, err := rand.Read(peerID[:]); err != nil {
		return [20]byte{}, errors.New("failed to generate peer ID: " + err.Error())
	}
	return peerID, nil
}
