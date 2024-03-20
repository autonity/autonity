package protocol

import (
	"fmt"
	"math/big"
	"time"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"
)

// HandlerFunc is a callback to invoke from an outside runner after the boilerplate
// exchanges have passed.
type HandlerFunc func(peer *Peer) error

// Backend defines the data retrieval methods to serve remote requests and the
// callback methods to invoke on remote deliveries.
type Backend interface {
	// Chain retrieves the blockchain object to serve data.
	Chain() *core.BlockChain

	// RunPeer is invoked when a peer joins on the `consensus` protocol. The ACN
	// should do any peer maintenance work, handshakes and validations. If all
	// is passed, control should be given back to the `ACN` to process the
	// inbound messages going forward.
	RunPeer(peer *Peer, handler HandlerFunc) error

	// PeerInfo retrieves all known `acn` information about a peer.
	PeerInfo(id enode.ID) interface{}
}

// NodeInfo represents a short summary of the `ACN` protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network    uint64              `json:"network"`    // Ethereum network ID (1=Frontier, 2=Morden, Ropsten=3, Rinkeby=4)
	Difficulty *big.Int            `json:"difficulty"` // Total difficulty of the host's blockchain
	Genesis    common.Hash         `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash         `json:"head"`       // Hex hash of the host's best owned block
}

// nodeInfo retrieves some `acn` protocol metadata about the running host node.
func nodeInfo(chain *core.BlockChain, network uint64) *NodeInfo {
	head := chain.CurrentBlock()
	return &NodeInfo{
		Network:    network,
		Difficulty: chain.GetTd(head.Hash(), head.NumberU64()),
		Genesis:    chain.Genesis().Hash(),
		Config:     chain.Config(),
		Head:       head.Hash(),
	}
}

func MakeProtocols(backend Backend, network uint64) []p2p.Protocol {
	protocols := make([]p2p.Protocol, len(ProtocolVersions))
	for i, version := range ProtocolVersions {
		version := version // Closure

		protocols[i] = p2p.Protocol{
			Name:    ProtocolName,
			Version: version,
			Length:  protocolLengths[version],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := NewPeer(version, p, rw)
				defer peer.Close()

				return backend.RunPeer(peer, func(peer *Peer) error {
					return Handle(backend, peer)
				})
			},
			NodeInfo: func() interface{} {
				return nodeInfo(backend.Chain(), network)
			},
			PeerInfo: func(id enode.ID) interface{} {
				return backend.PeerInfo(id)
			},
			Attributes:     nil,
			DialCandidates: nil,
		}
	}
	return protocols
}

// Handle is invoked whenever an `consensus` connection is made that successfully passes
// the protocol handshake. This method will keep processing messages until the
// connection is torn down.
func Handle(backend Backend, peer *Peer) error {
	errCh := make(chan error, 1)
	for {
		if err := handleMessage(backend, peer, errCh); err != nil {
			peer.Log().Debug("Message handling failed in `acn`", "err", err)
			return err
		}
		select {
		case err := <-errCh:
			peer.Log().Error("Message handling failed in consensus core", "err", err)
			return err
		default:
			// do nothing
		}
	}
}

// handleMessage is invoked whenever an inbound message is received from a remote
// peer. The remote connection is torn down upon returning any error.
func handleMessage(backend Backend, peer *Peer, errCh chan<- error) error {
	// Read the next message from the remote peer, and ensure it's fully consumed
	msg, err := peer.rw.ReadMsg()
	if err != nil {
		return err
	}
	if msg.Size > MaxMessageSize {
		return fmt.Errorf("%w: %v > %v", errMsgTooLarge, msg.Size, MaxMessageSize)
	}
	defer msg.Discard()

	// Track the amount of time it takes to serve the request and run the ACN
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%d/%#02x", p2p.HandleHistName, ProtocolName, peer.Version(), msg.Code)
		defer func(start time.Time) {
			metrics.GetOrRegisterBufferedGauge(h, nil).Add(time.Since(start).Nanoseconds())
		}(time.Now())
	}
	if handler, ok := backend.Chain().Engine().(consensus.Handler); ok {
		if handled, err := handler.HandleMsg(peer.address, msg, errCh); handled {
			return err
		}
	}
	return fmt.Errorf("%w: %v", errInvalidMsgCode, msg.Code)
}
