package afd

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/events"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"math/rand"
	"sync"
	"time"
)

var (
	// todo: config the window and buffer height in genesis.
	randomDelayWindow = 1000 * 10 // (0, 10] seconds random time window
	msgBufferInHeight = 60    // buffer such range of msgs in height at msg store.
	errFutureMsg = errors.New("future height msg")
	errGarbageMsg = errors.New("garbage msg")
	errNotCommitteeMsg = errors.New("msg from none committee member")
	errProposer = errors.New("proposal is not from proposer")
	errProposal = errors.New("proposal have invalid values")
	errEquivocation = errors.New("equivocation happens")
	errUnknownMsg = errors.New("unknown consensus msg")
	errInvalidChallenge = errors.New("invalid challenge")
	errNoEvidence = errors.New("no evidence")
)

// Fault detector, it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	// use below 3 members to send proof via transaction issuing.
	wg sync.WaitGroup
	afdFeed event.Feed
	scope event.SubscriptionScope

	// use below 2 members to forward consensus msg from protocol manager to afd.
	tendermintMsgSub *event.TypeMuxSubscription
	tendermintMsgMux *event.TypeMuxSilent

	// below 2 members subscribe block event to trigger execution
	// of rule engine and make proof of innocent.
	blockChan chan core.ChainEvent
	blockSub event.Subscription

	// chain context to validate consensus msgs.
	blockchain *core.BlockChain

	// node address
	address common.Address

	// msg store
	msgStore *MsgStore

	// future height msg buffer
	futureMsgs map[uint64][]*types.ConsensusMessage

	logger log.Logger
}


// call by ethereum object to create fd instance.
func NewFaultDetector(chain *core.BlockChain, nodeAddress common.Address) *FaultDetector {
	logger := log.New("afd", nodeAddress)
	fd := &FaultDetector{
		address: nodeAddress,
		blockChan:  make(chan core.ChainEvent, 300),
		blockchain: chain,
		msgStore: new(MsgStore),
		logger:logger,
		tendermintMsgMux:  event.NewTypeMuxSilent(logger),
		futureMsgs: make(map[uint64][]*types.ConsensusMessage),
	}

	// register afd contracts on evm's precompiled contract set.
	registerAFDContracts(chain)

	// subscribe tendermint msg
	s := fd.tendermintMsgMux.Subscribe(events.MessageEvent{})
	fd.tendermintMsgSub = s
	return fd
}

// listen for new block events from block-chain, do the tasks like take challenge and provide proof for innocent, the
// AFD rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)

	for {
		select {
		// chain event update, provide proof of innocent if one is on challenge, rule engine scanning is triggered also.
		case ev := <-fd.blockChan:
			// take my challenge from latest state DB, and provide innocent proof if there are any.
			err := fd.handleChallenges(ev.Block, ev.Hash)
			if err != nil {
				fd.logger.Warn("handle challenge","afd", err)
			}

			// before run rule engine over msg store, check to process any buffered msg.
			fd.processBufferedMsgs(ev.Block.NumberU64())

			// run rule engine over msg store on each height update.

			fd.runRuleEngine(ev.Block.NumberU64())

			// msg store delete msgs out of buffering window.
			fd.msgStore.DeleteMsgsAtHeight(ev.Block.NumberU64() - uint64(msgBufferInHeight))

		// to handle consensus msg from p2p layer.	
		case ev, ok := <-fd.tendermintMsgSub.Chan():
			if !ok {
				return
			}
			switch e := ev.Data.(type) {
			case events.MessageEvent:
				msg := new(types.ConsensusMessage)
				if err := msg.FromPayload(e.Payload); err != nil {
					fd.logger.Error("invalid payload", "afd", err)
					continue
				}
				if err := fd.processMsg(msg); err != nil {
					fd.logger.Error("process consensus msg", "afd", err)
					continue
				}
			}

		case <-fd.blockSub.Err():
			return
		}
	}
}

// HandleMsg is called by p2p protocol manager to deliver the consensus msg to afd.
func (fd *FaultDetector) HandleMsg(addr common.Address, msg p2p.Msg) {
	if msg.Code != types.TendermintMsg {
		return
	}

	var data []byte
	if err := msg.Decode(&data); err != nil {
		log.Error("cannot decode consensus msg", "from", addr)
		return
	}

	// post consensus event to event loop.
	fd.tendermintMsgMux.Post(events.MessageEvent{Payload:data})
	return
}

func (fd *FaultDetector) Stop() {
	fd.scope.Close()
	fd.blockSub.Unsubscribe()
	fd.tendermintMsgSub.Unsubscribe()
	fd.wg.Wait()
	unRegisterAFDContracts()
}

// call by ethereum object to subscribe proofs Events.
func (fd *FaultDetector) SubscribeAFDEvents(ch chan<- types.SubmitProofEvent) event.Subscription {
	return fd.scope.Track(fd.afdFeed.Subscribe(ch))
}

// get challenges from chain via autonityContract calls, and provide proofs if there were any challenge of client.
func (fd *FaultDetector) handleChallenges(block *types.Block, hash common.Hash) error {
	var innocentProofs []types.OnChainProof
	state, err := fd.blockchain.StateAt(hash)
	if err != nil {
		return err
	}

	challenges := fd.blockchain.GetAutonityContract().GetChallenges(block.Header(), state)
	for i:=0; i < len(challenges); i++ {
		if challenges[i].SenderHash == types.RLPHash(fd.address) {
			c, err := decodeProof(challenges[i].RawProofBytes)
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
	fd.sendProofs(types.InnocentProof, innocentProofs)
	return nil
}

func (fd *FaultDetector) randomDelay() {
	// wait for random milliseconds (under the range of 10 seconds) to check if need to rise challenge.
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(randomDelayWindow)
	time.Sleep(time.Duration(n)*time.Millisecond)
}

// send proofs via event which will handled by ethereum object to signed the TX to send proof.
func (fd *FaultDetector) sendProofs(t types.ProofType,  proofs[]types.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		if t == types.InnocentProof {
			fd.afdFeed.Send(types.SubmitProofEvent{Proofs:proofs, Type:t})
		}

		if t == types.ChallengeProof {
			fd.randomDelay()
			unPresented := fd.filterPresentedChallenges(&proofs)
			if len(unPresented) != 0 {
				fd.afdFeed.Send(types.SubmitProofEvent{Proofs:unPresented, Type:t})
			}
		}
	}()
}

func (fd *FaultDetector) filterPresentedChallenges(proofs *[]types.OnChainProof) []types.OnChainProof {
	// get latest chain state.
	var result []types.OnChainProof
	state, err := fd.blockchain.State()
	if err != nil {
		return nil
	}

	header := fd.blockchain.CurrentBlock().Header()
	challenges := fd.blockchain.GetAutonityContract().GetChallenges(header, state)

	for i:=0; i < len(*proofs); i++ {
		present := false
		for i:=0; i < len(challenges); i++ {
			if (*proofs)[i].MsgHash == challenges[i].MsgHash {
				present = true
			}
		}
		if !present {
			result = append(result, (*proofs)[i])
		}
	}

	return result
}
