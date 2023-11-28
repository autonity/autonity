package tm

import (
	"github.com/autonity/autonity/crypto"
	"math/big"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/p2p"
	mapset "github.com/deckarep/golang-set"
)

const (
	// maxKnownBlocks is the maximum block hashes to keep in the known list
	// before starting to randomly evict them.
	maxKnownBlocks = 1024

	// maxQueuedBlocks is the maximum number of block propagations to queue up before
	// dropping broadcasts. There's not much point in queueing stale blocks, so a few
	// that might cover uncles should be enough.
	maxQueuedBlocks = 4

	// maxQueuedBlockAnns is the maximum number of block announcements to queue up before
	// dropping broadcasts. Similarly to block propagations, there's no point to queue
	// above some healthy uncle limit, so use that.
	maxQueuedBlockAnns = 4
)

// max is a helper function which returns the larger of the two given integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// Peer is a collection of relevant information we have about a `eth` peer.
type Peer struct {
	id      string // Unique ID for the peer, cached
	address common.Address

	*p2p.Peer                   // The embedded P2P package peer
	rw        p2p.MsgReadWriter // Input/output streams for snap
	version   uint              // Protocol version negotiated

	head common.Hash // Latest advertised head block hash
	td   *big.Int    // Latest advertised head block total difficulty

	txpool          TxPool                 // Transaction pool used by the broadcasters for liveness checks
	knownBlocks     *knownCache            // Set of block hashes known to be known by this peer
	queuedBlocks    chan *blockPropagation // Queue of blocks to broadcast to the peer
	queuedBlockAnns chan *types.Block      // Queue of blocks to announce to the peer

	term chan struct{} // Termination channel to stop the broadcasters
	lock sync.RWMutex  // Mutex protecting the internal fields
}

// TODO: check if TXPool is really required
// NewPeer create a wrapper for a network connection and negotiated  protocol
// version.
func NewPeer(version uint, p *p2p.Peer, rw p2p.MsgReadWriter, txpool TxPool) *Peer {
	peer := &Peer{
		id:              p.ID().String(),
		address:         crypto.PubkeyToAddress(*p.Node().Pubkey()),
		Peer:            p,
		rw:              rw,
		version:         version,
		knownBlocks:     newKnownCache(maxKnownBlocks),
		queuedBlocks:    make(chan *blockPropagation, maxQueuedBlocks),
		queuedBlockAnns: make(chan *types.Block, maxQueuedBlockAnns),
		txpool:          txpool,
		term:            make(chan struct{}),
	}
	// Start up all the broadcasters
	go peer.broadcastBlocks()
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

// Version retrieves the peer's negoatiated `eth` protocol version.
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

// KnownBlock returns whether peer is known to already have a block.
func (p *Peer) KnownBlock(hash common.Hash) bool {
	return p.knownBlocks.Contains(hash)
}

// markBlock marks a block as known for the peer, ensuring that the block will
// never be propagated to this particular peer.
func (p *Peer) markBlock(hash common.Hash) {
	// If we reached the memory allowance, drop a previously known block hash
	p.knownBlocks.Add(hash)
}

// SendNewBlockHashes announces the availability of a number of blocks through
// a hash notification.
func (p *Peer) SendNewBlockHashes(hashes []common.Hash, numbers []uint64) error {
	// Mark all the block hashes as known, but ensure we don't overflow our limits
	p.knownBlocks.Add(hashes...)

	request := make(NewBlockHashesPacket, len(hashes))
	for i := 0; i < len(hashes); i++ {
		request[i].Hash = hashes[i]
		request[i].Number = numbers[i]
	}
	return p2p.Send(p.rw, NewBlockHashesMsg, request)
}

// AsyncSendNewBlockHash queues the availability of a block for propagation to a
// remote peer. If the peer's broadcast queue is full, the event is silently
// dropped.
func (p *Peer) AsyncSendNewBlockHash(block *types.Block) {
	select {
	case p.queuedBlockAnns <- block:
		// Mark all the block hash as known, but ensure we don't overflow our limits
		p.knownBlocks.Add(block.Hash())
	default:
		p.Log().Debug("Dropping block announcement", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// SendNewBlock propagates an entire block to a remote peer.
func (p *Peer) SendNewBlock(block *types.Block, td *big.Int) error {
	// Mark all the block hash as known, but ensure we don't overflow our limits
	p.knownBlocks.Add(block.Hash())
	return p2p.Send(p.rw, NewBlockMsg, &NewBlockPacket{
		Block: block,
		TD:    td,
	})
}

// AsyncSendNewBlock queues an entire block for propagation to a remote peer. If
// the peer's broadcast queue is full, the event is silently dropped.
func (p *Peer) AsyncSendNewBlock(block *types.Block, td *big.Int) {
	select {
	case p.queuedBlocks <- &blockPropagation{block: block, td: td}:
		// Mark all the block hash as known, but ensure we don't overflow our limits
		p.knownBlocks.Add(block.Hash())
	default:
		p.Log().Debug("Dropping block propagation", "number", block.NumberU64(), "hash", block.Hash())
	}
}

// knownCache is a cache for known hashes.
type knownCache struct {
	hashes mapset.Set
	max    int
}

// newKnownCache creates a new knownCache with a max capacity.
func newKnownCache(max int) *knownCache {
	return &knownCache{
		max:    max,
		hashes: mapset.NewSet(),
	}
}

// Add adds a list of elements to the set.
func (k *knownCache) Add(hashes ...common.Hash) {
	for k.hashes.Cardinality() > max(0, k.max-len(hashes)) {
		k.hashes.Pop()
	}
	for _, hash := range hashes {
		k.hashes.Add(hash)
	}
}

// Contains returns whether the given item is in the set.
func (k *knownCache) Contains(hash common.Hash) bool {
	return k.hashes.Contains(hash)
}

// Cardinality returns the number of elements in the set.
func (k *knownCache) Cardinality() int {
	return k.hashes.Cardinality()
}
