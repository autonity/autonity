package faultdetector

import (
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/common"
	tendermintBackend "github.com/clearmatics/autonity/consensus/tendermint/backend"
	"github.com/clearmatics/autonity/consensus/tendermint/bft"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

var (
	// todo: refine the window and buffer range in contract which can be tuned during run time.
	randomDelayWindow   = 1000 * 10 // (0, 10] seconds random time window
	msgBufferInHeight   = 60        // buffer such range of msgs in height at msg store.
	errFutureMsg        = errors.New("future height msg")
	errGarbageMsg       = errors.New("garbage msg")
	errNotCommitteeMsg  = errors.New("msg from none committee member")
	errProposer         = errors.New("proposal is not from proposer")
	errProposal         = errors.New("proposal have invalid values")
	errEquivocation     = errors.New("equivocation happens")
	errUnknownMsg       = errors.New("unknown consensus msg")
	errInvalidChallenge = errors.New("invalid challenge")
)

// Fault detector, it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	// use below 3 members to send proof via transaction issuing.
	wg      sync.WaitGroup
	afdFeed event.Feed
	scope   event.SubscriptionScope

	// use below 2 members to forward consensus msg from protocol manager to faultdetector.
	tendermintMsgSub *event.TypeMuxSubscription
	tendermintMsgMux *event.TypeMuxSilent

	// below 2 members subscribe block event to trigger execution
	// of rule engine and make proof of innocent.
	blockChan chan core.ChainEvent
	blockSub  event.Subscription

	// chain context to validate consensus msgs.
	blockchain *core.BlockChain

	// node address
	address common.Address

	// msg store
	msgStore *MsgStore

	// future height msg buffer
	futureMsgs map[uint64][]*tendermintCore.Message

	// buffer quorum for blocks.
	totalPowers map[uint64]uint64
	logger      log.Logger
}

// call by ethereum object to create fd instance.
func NewFaultDetector(chain *core.BlockChain, nodeAddress common.Address) *FaultDetector {
	logger := log.New("faultdetector", nodeAddress)
	fd := &FaultDetector{
		address:          nodeAddress,
		blockChan:        make(chan core.ChainEvent, 300),
		blockchain:       chain,
		msgStore:         newMsgStore(),
		logger:           logger,
		tendermintMsgMux: event.NewTypeMuxSilent(logger),
		futureMsgs:       make(map[uint64][]*tendermintCore.Message),
		totalPowers:      make(map[uint64]uint64),
	}

	// register faultdetector contracts on evm's precompiled contract set.
	registerAFDContracts(chain)

	// subscribe tendermint msg
	s := fd.tendermintMsgMux.Subscribe(events.MessageEvent{})
	fd.tendermintMsgSub = s
	return fd
}

func (fd *FaultDetector) quorum(h uint64) uint64 {
	power, ok := fd.totalPowers[h]
	if ok {
		return bft.Quorum(power)
	}

	return bft.Quorum(fd.blockchain.GetHeaderByNumber(h).TotalVotingPower())
}

func (fd *FaultDetector) savePower(h uint64, power uint64) {
	fd.totalPowers[h] = power
}

func (fd *FaultDetector) deletePower(h uint64) {
	delete(fd.totalPowers, h)
}

// listen for new block events from block-chain, do the tasks like take challenge and provide proof for innocent, the
// AFD rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)

	for {
		select {
		// chain event update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev := <-fd.blockChan:
			fd.savePower(ev.Block.Number().Uint64(), ev.Block.Header().TotalVotingPower())

			// take my accusations from latest state DB, and provide innocent proof if there are any.
			err := fd.handleAccusations(ev.Block, ev.Block.Root())
			if err != nil {
				fd.logger.Warn("handle challenge", "faultdetector", err)
			}

			// before run rule engine over msg store, check to process any buffered msg.
			fd.processBufferedMsgs(ev.Block.NumberU64())

			// run rule engine over msg store on each height update.
			quorum := fd.quorum(ev.Block.NumberU64() - 1)
			fd.runRuleEngine(ev.Block.NumberU64(), quorum)

			// msg store delete msgs out of buffering window.
			fd.msgStore.DeleteMsgsAtHeight(ev.Block.NumberU64() - uint64(msgBufferInHeight))

			// delete power out of buffering window.
			fd.deletePower(ev.Block.NumberU64() - uint64(msgBufferInHeight))
		// to handle consensus msg from p2p layer.
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			if !ok {
				return
			}
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				msg := new(tendermintCore.Message)
				if err := msg.FromPayload(e.Payload); err != nil {
					fd.logger.Error("invalid payload", "faultdetector", err)
					continue
				}
				if err := fd.processMsg(msg); err != nil {
					fd.logger.Error("process consensus msg", "faultdetector", err)
					continue
				}
			}

		case <-fd.blockSub.Err():
			return
		}
	}
}

// HandleMsg is called by p2p protocol manager to deliver the consensus msg to faultdetector.
func (fd *FaultDetector) HandleMsg(addr common.Address, msg p2p.Msg) {
	if msg.Code != tendermintBackend.TendermintMsg {
		return
	}

	var data []byte
	if err := msg.Decode(&data); err != nil {
		log.Error("cannot decode consensus msg", "from", addr)
		return
	}

	// post consensus event to event loop.
	fd.tendermintMsgMux.Post(events.MessageEvent{Payload: data})
}

func (fd *FaultDetector) Stop() {
	fd.scope.Close()
	fd.blockSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.wg.Wait()
	unRegisterAFDContracts()
}

// call by ethereum object to subscribe proofs Events.
func (fd *FaultDetector) SubscribeAFDEvents(ch chan<- AccountabilityEvent) event.Subscription {
	return fd.scope.Track(fd.afdFeed.Subscribe(ch))
}

// get accusations from chain via autonityContract calls, and provide innocent proofs if there were any challenge on node.
func (fd *FaultDetector) handleAccusations(block *types.Block, hash common.Hash) error {
	var innocentProofs []autonity.OnChainProof
	state, err := fd.blockchain.StateAt(hash)
	if err != nil {
		return err
	}

	accusations := fd.blockchain.GetAutonityContract().GetAccusations(block.Header(), state)
	for i := 0; i < len(accusations); i++ {
		if accusations[i].Sender == fd.address {
			c, err := decodeProof(accusations[i].Rawproof)
			if err != nil {
				continue
			}

			p, err := fd.getInnocentProof(c)
			if err != nil {
				continue
			}
			innocentProofs = append(innocentProofs, p)
		}
	}

	// send proofs via standard transaction.
	fd.sendProofs(false, innocentProofs)
	return nil
}

func (fd *FaultDetector) randomDelay() {
	// wait for random milliseconds (under the range of 10 seconds) to check if need to rise challenge.
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(randomDelayWindow)
	time.Sleep(time.Duration(n) * time.Millisecond)
}

// send proofs via event which will handled by ethereum object to signed the TX to send proof.
func (fd *FaultDetector) sendProofs(withDelay bool, proofs []autonity.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		if !withDelay {
			fd.afdFeed.Send(AccountabilityEvent{Proofs: proofs})
		}

		if withDelay {
			fd.randomDelay()
			unPresented := fd.filterPresentedOnes(&proofs)
			if len(unPresented) != 0 {
				fd.afdFeed.Send(AccountabilityEvent{Proofs: unPresented})
			}
		}
	}()
}

func (fd *FaultDetector) filterPresentedOnes(proofs *[]autonity.OnChainProof) []autonity.OnChainProof {
	// get latest chain state.
	var result []autonity.OnChainProof
	state, err := fd.blockchain.State()
	if err != nil {
		return nil
	}
	header := fd.blockchain.CurrentBlock().Header()

	presentedAccusation := fd.blockchain.GetAutonityContract().GetAccusations(header, state)
	presentedMisbehavior := fd.blockchain.GetAutonityContract().GetChallenges(header, state)

	for i := 0; i < len(*proofs); i++ {
		present := false
		for j := 0; j < len(presentedAccusation); j++ {
		    if (*proofs)[i].Msghash == presentedAccusation[j].Msghash &&
		        (*proofs)[i].Type.Cmp(new(big.Int).SetUint64(uint64(Accusation))) == 0 {
				present = true
			}
		}

		for j := 0; j < len(presentedMisbehavior); j++ {
			if (*proofs)[i].Msghash == presentedMisbehavior[j].Msghash &&
				(*proofs)[i].Type.Cmp(new(big.Int).SetUint64(uint64(Misbehaviour))) == 0 {
				present = true
			}
		}

		if !present {
			result = append(result, (*proofs)[i])
		}
	}

	return result
}
