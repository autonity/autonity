package protocol

import (
	"math/big"
	"sync"

	"github.com/autonity/autonity/crypto"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/p2p"
)

// Peer is a collection of relevant information we have about a `acn` peer.
type Peer struct {
	id      string // Unique ID for the peer, cached
	address common.Address

	*p2p.Peer                   // The embedded P2P package peer
	rw        p2p.MsgReadWriter // Input/output streams for snap
	version   uint              // Protocol version negotiated

	head common.Hash // Latest advertised head block hash
	td   *big.Int    // Latest advertised head block total difficulty

	term chan struct{} // Termination channel to stop the broadcasters
	lock sync.RWMutex  // Mutex protecting the internal fields
}

// peerInfo represents a short summary of the `acn` protocol metadata known
// about a connected peer.
type peerInfo struct {
	Version    uint     `json:"version"`    // Acn protocol version negotiated
	Difficulty *big.Int `json:"difficulty"` // Total difficulty of the peer's blockchain
	Head       string   `json:"head"`       // Hex hash of the peer's best owned block
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
		term:    make(chan struct{}),
	}
	// Start up all the broadcasters
	// no block broadcasting for consensus Peers
	//go peer.broadcastBlocks()
	return peer
}

// Close signals the broadcast goroutine to terminate. Only ever call this if
// you created the peer yourself via NewPeer, Otherwise let whoever created it
// clean it up!
func (p *Peer) Close() {
	close(p.term)
}

// ID retrieves the peer's unique identifier.
func (p *Peer) ID() string {

	return p.id
}

func (p *Peer) Address() common.Address {
	return p.address
}

func (p *Peer) Send(msgcode uint64, data interface{}) error {
	return p2p.Send(p.rw, msgcode, data)
}

func (p *Peer) SendRaw(msgcode uint64, data []byte) error {
	return p2p.SendRaw(p.rw, msgcode, data)
}

// Version retrieves the peer's negoatiated `acn` protocol version.
func (p *Peer) Version() uint {
	return p.version
}

// Head retrieves the current head hash and total difficulty of the peer.
func (p *Peer) Head() (hash common.Hash, td *big.Int) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	copy(hash[:], p.head[:])
	return hash, new(big.Int).Set(p.td)
}

// SetHead updates the head hash and total difficulty of the peer.
func (p *Peer) SetHead(hash common.Hash, td *big.Int) {
	p.lock.Lock()
	defer p.lock.Unlock()

	copy(p.head[:], hash[:])
	p.td.Set(td)
}

// info gathers and returns some `eth` protocol metadata known about a peer.
func (p *Peer) info() *peerInfo {
	hash, td := p.Head()

	return &peerInfo{
		Version:    p.Version(),
		Difficulty: td,
		Head:       hash.Hex(),
	}
}
