package tm

import (
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/p2p"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/p2p/enr"
	"github.com/autonity/autonity/params"
	"math/big"
	"time"
)

// Handler is a callback to invoke from an outside runner after the boilerplate
// exchanges have passed.
type Handler func(peer *Peer) error

// TODO: update the comments
// Backend defines the data retrieval methods to serve remote requests and the
// callback methods to invoke on remote deliveries.
type Backend interface {
	// Chain retrieves the blockchain object to serve data.
	Chain() *core.BlockChain

	// TxPool retrieves the transaction pool object to serve data.
	TxPool() TxPool

	// RunPeer is invoked when a peer joins on the `consensus` protocol. The handler
	// should do any peer maintenance work, handshakes and validations. If all
	// is passed, control should be given back to the `handler` to process the
	// inbound messages going forward.
	RunPeer(peer *Peer, handler Handler) error

	// PeerInfo retrieves all known `consensus` information about a peer.
	PeerInfo(id enode.ID) interface{}

	// Handle is a callback to be invoked when a data packet is received from
	// the remote peer. Only packets not consumed by the protocol handler will
	// be forwarded to the backend.
	Handle(peer *Peer, packet Packet) error
}

// TxPool defines the methods needed by the protocol handler to serve transactions.
type TxPool interface {
	// Get retrieves the transaction from the local txpool with the given hash.
	Get(hash common.Hash) *types.Transaction
}

// NodeInfo represents a short summary of the `eth` sub-protocol metadata
// known about the host peer.
type NodeInfo struct {
	Network    uint64              `json:"network"`    // Ethereum network ID (1=Frontier, 2=Morden, Ropsten=3, Rinkeby=4)
	Difficulty *big.Int            `json:"difficulty"` // Total difficulty of the host's blockchain
	Genesis    common.Hash         `json:"genesis"`    // SHA3 hash of the host's genesis block
	Config     *params.ChainConfig `json:"config"`     // Chain configuration for the fork rules
	Head       common.Hash         `json:"head"`       // Hex hash of the host's best owned block
}

// nodeInfo retrieves some `eth` protocol metadata about the running host node.
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

// TODO: we need to disconnect the consensus node and not redial it if
// it is not in the consensus committee, can be handled in the consensus layer though
func MakeProtocols(backend Backend, network uint64) []p2p.Protocol {
	protocols := make([]p2p.Protocol, len(ProtocolVersions))
	for i, version := range ProtocolVersions {
		version := version // Closure

		protocols[i] = p2p.Protocol{
			Name:    ProtocolName,
			Version: version,
			Length:  protocolLengths[version],
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				peer := NewPeer(version, p, rw, backend.TxPool())
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
			Attributes:     []enr.Entry{currentENREntry(backend.Chain())},
			DialCandidates: nil,
		}
	}
	return protocols
}

// Handle is invoked whenever an `eth` connection is made that successfully passes
// the protocol handshake. This method will keep processing messages until the
// connection is torn down.
func Handle(backend Backend, peer *Peer) error {
	errCh := make(chan error, 1)
	for {
		if err := handleMessage(backend, peer, errCh); err != nil {
			peer.Log().Debug("Message handling failed in `eth`", "err", err)
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

	// Track the amount of time it takes to serve the request and run the handler
	if metrics.Enabled {
		h := fmt.Sprintf("%s/%s/%d/%#02x", p2p.HandleHistName, ProtocolName, peer.Version(), msg.Code)
		defer func(start time.Time) {
			sampler := func() metrics.Sample {
				return metrics.ResettingSample(
					metrics.NewExpDecaySample(1028, 0.015),
				)
			}
			metrics.GetOrRegisterHistogramLazy(h, nil, sampler).Update(time.Since(start).Microseconds())
		}(time.Now())
	}
	if handler, ok := backend.Chain().Engine().(consensus.Handler); ok {
		if handled, err := handler.HandleMsg(peer.address, msg, errCh); handled {
			return err
		}
	}
	/*
		//TODO: verify if this to be done in consensus handler or not
		var handlers = eth66
		//if peer.Version() >= ETH67 { // Left in as a sample when new protocol is added
		//	handlers = eth67
		//}
		if handler := handlers[msg.Code]; handler != nil {
			return handler(backend, msg, peer)
		}
	*/
	return fmt.Errorf("%w: %v", errInvalidMsgCode, msg.Code)
}
