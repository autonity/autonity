package protocol

import (
	"time"

	"github.com/autonity/autonity/common/fixsizecache"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/p2p/enode"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/p2p"
)

var (
	ProposalWriteBg  = metrics.NewRegisteredBufferedGauge("acn/proposal/write", nil, metrics.GetIntPointer(1000))  // time to write proposal to wire
	PrevoteWriteBg   = metrics.NewRegisteredBufferedGauge("acn/prevote/write", nil, metrics.GetIntPointer(5000))   // time to write prevote to wire
	PrecommitWriteBg = metrics.NewRegisteredBufferedGauge("acn/precommit/write", nil, metrics.GetIntPointer(5000)) // time to write precommit to wire
	DefaultWriteBg   = metrics.NewRegisteredBufferedGauge("acn/any/write", nil, nil)
)

const (
	buckets = 199
	entries = 10
)

// Peer is a collection of relevant information we have about a `acn` peer.
type Peer struct {
	id      enode.ID // Unique ID for the peer, cached
	address common.Address

	*p2p.Peer                   // The embedded P2P package peer
	rw        p2p.MsgReadWriter // Input/output streams for snap
	version   uint              // Protocol version negotiated
	cache     *fixsizecache.Cache[common.Hash, bool]
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
		id:      p.ID(),
		address: crypto.PubkeyToAddress(*p.Node().Pubkey()),
		Peer:    p,
		rw:      rw,
		version: version,
		cache:   fixsizecache.New[common.Hash, bool](buckets, entries, 0, fixsizecache.HashKey[common.Hash]),
	}
	return peer
}

func (p *Peer) Cache() *fixsizecache.Cache[common.Hash, bool] {
	return p.cache
}

// Close can be used to do peer related clean up, nothing for now
func (p *Peer) Close() {
	// nothing to do
}

// ID retrieves the peer's unique identifier.
func (p *Peer) ID() enode.ID {

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
