package protocol

import (
	"time"

	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/metrics"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/p2p"
)

func getIntPointer(val int) *int {
	return &val
}

var (
	ProposalWriteBg  = metrics.NewRegisteredBufferedGauge("acn/proposal/write", nil, getIntPointer(1000))  // time between round start and proposal sent
	PrevoteWriteBg   = metrics.NewRegisteredBufferedGauge("acn/prevote/write", nil, getIntPointer(5000))   // time between round start and proposal received
	PrecommitWriteBg = metrics.NewRegisteredBufferedGauge("acn/precommit/write", nil, getIntPointer(5000)) // time to verify proposal
	DefaultWriteBg   = metrics.NewRegisteredBufferedGauge("acn/any/write", nil, nil)                       // time to verify proposal
)

// Peer is a collection of relevant information we have about a `acn` peer.
type Peer struct {
	id      string // Unique ID for the peer, cached
	address common.Address

	*p2p.Peer                   // The embedded P2P package peer
	rw        p2p.MsgReadWriter // Input/output streams for snap
	version   uint              // Protocol version negotiated

}

// peerInfo represents a short summary of the `acn` protocol metadata known
// about a connected peer.
type peerInfo struct {
	Version uint `json:"version"` // Acn protocol version negotiated
}

// NewPeer create a wrapper for a network connection and negotiated  protocol
// version.
func NewPeer(version uint, p *p2p.Peer, rw p2p.MsgReadWriter) *Peer {
	peer := &Peer{
		id:      p.ID().String(),
		address: crypto.PubkeyToAddress(*p.Node().Pubkey()),
		Peer:    p,
		rw:      rw,
		version: version,
	}
	return peer
}

// Close can be used to do peer related clean up, nothing for now
func (p *Peer) Close() {
	// nothing to do
}

// ID retrieves the peer's unique identifier.
func (p *Peer) ID() string {

	return p.id
}

func (p *Peer) Address() common.Address {
	return p.address
}

func (p *Peer) Send(msgcode uint64, data interface{}) error {
	if metrics.Enabled {
		defer func(start time.Time) {
			getWriteMetric(msgcode).Add(time.Since(start).Nanoseconds())
		}(time.Now())
	}
	return p2p.Send(p.rw, msgcode, data)
}

func (p *Peer) SendRaw(msgcode uint64, data []byte) error {
	if metrics.Enabled {
		defer func(start time.Time) {
			getWriteMetric(msgcode).Add(time.Since(start).Nanoseconds())
		}(time.Now())
	}
	return p2p.SendRaw(p.rw, msgcode, data)
}

func getWriteMetric(msgCode uint64) metrics.BufferedGauge {
	switch msgCode {
	case 0x11:
		return ProposalWriteBg
	case 0x12:
		return PrevoteWriteBg
	case 0x13:
		return PrecommitWriteBg
	}
	return DefaultWriteBg
}

// Version retrieves the peer's negoatiated `acn` protocol version.
func (p *Peer) Version() uint {
	return p.version
}

// ConsensusPeerInfo gathers and returns some `acn` protocol metadata known about a peer.
func (p *Peer) ConsensusPeerInfo() *peerInfo {
	return &peerInfo{
		Version: p.Version(),
	}
}
