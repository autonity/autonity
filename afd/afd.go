package afd

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/p2p"
	"math/rand"
	"sync"
	"time"
)

// Fault detector, it subscribe chain event to trigger rule engine to apply patterns over
// msg store, it send proof of challenge if it detects any potential misbehavior, either it
// read state db on each new height to get latest challenges from autonity contract's view,
// and to prove its innocent if there were any challenges on the suspicious node.
type FaultDetector struct {
	wg sync.WaitGroup
	afdFeed event.Feed
	scope event.SubscriptionScope

	blockChan chan core.ChainEvent
	blockSub event.Subscription
	blockchain *core.BlockChain
	address common.Address
	msgStore *MsgStore
	ruleEngine *RuleEngine
}

var(
	// todo: to be configured at genesis.
	randomDelayWindow = 10000
)

// call by ethereum object to create fd instance.
func NewFaultDetector(chain *core.BlockChain, nodeAddress common.Address) *FaultDetector {
	fd := &FaultDetector{
		address: nodeAddress,
		blockChan:  make(chan core.ChainEvent, 300),
		blockchain: chain,
		msgStore: new(MsgStore),
		ruleEngine: new(RuleEngine),
	}

	// init accountability precompiled contracts.
	initAccountabilityContracts(chain)
	return fd
}

// listen for new block events from block-chain, do the tasks like take challenge and provide proof for innocent, the
// AFD rule engine could also triggered from here to scan those msgs of msg store by applying rules.
func (fd *FaultDetector) FaultDetectorEventLoop() {
	fd.blockSub = fd.blockchain.SubscribeChainEvent(fd.blockChan)
	for {
		select {
		case ev := <-fd.blockChan:
			// take my challenge from latest state DB, and provide innocent proof if there are any.
			err := fd.handleMyChallenges(ev.Block, ev.Hash)
			if err != nil {
				// prints something.
			}

			// todo: tell rule engine to run patterns over msg store on each new height.
			fd.ruleEngine.run()

		case <-fd.blockSub.Err():
			return
		}
	}
}

// get challenges from blockchain via autonityContract calls.
func (fd *FaultDetector) handleMyChallenges(block *types.Block, hash common.Hash) error {
	var innocentProofs []types.OnChainProof
	state, err := fd.blockchain.StateAt(hash)
	if err != nil {
		return err
	}

	challenges := fd.blockchain.GetAutonityContract().GetChallenges(block.Header(), state)
	for i:=0; i < len(challenges); i++ {
		if challenges[i].Sender == fd.address {
			p, err := fd.proveInnocent(challenges[i])
			if err != nil {
				continue
			}
			innocentProofs = append(innocentProofs, p)
		}
	}

	// send proofs via standard transaction.
	fd.SendProofs(types.InnocentProof, innocentProofs)
	return nil
}

// get proof of innocent over msg store.
func (fd *FaultDetector) proveInnocent(challenge types.OnChainProof) (types.OnChainProof, error) {
	// todo: get proof from msg store over the rule.
	var proof types.OnChainProof
	return proof, nil
}

func (fd *FaultDetector) Stop() {
	fd.scope.Close()
	fd.blockSub.Unsubscribe()
	fd.wg.Wait()
	cleanContracts()
}

// call by ethereum object to subscribe proofs Events.
func (fd *FaultDetector) SubscribeAFDEvents(ch chan<- types.SubmitProofEvent) event.Subscription {
	return fd.scope.Track(fd.afdFeed.Subscribe(ch))
}

// send proofs via event which will handled by ethereum object to signed the TX to send proof.
func (fd *FaultDetector) SendProofs(t types.ProofType,  proofs[]types.OnChainProof) {
	fd.wg.Add(1)
	go func() {
		defer fd.wg.Done()
		if t == types.InnocentProof {
			fd.afdFeed.Send(types.SubmitProofEvent{Proofs:proofs, Type:t})
		}

		if t == types.ChallengeProof {
			// wait for random milliseconds (under the range of 10 seconds) to check if need to rise challenge.
			rand.Seed(time.Now().UnixNano())
			n := rand.Intn(randomDelayWindow)
			time.Sleep(time.Duration(n)*time.Millisecond)
			unPresented := fd.filterUnPresentedChallenges(&proofs)
			if len(unPresented) != 0 {
				fd.afdFeed.Send(types.SubmitProofEvent{Proofs:unPresented, Type:t})
			}
		}
	}()
}

func (fd *FaultDetector) filterUnPresentedChallenges(proofs *[]types.OnChainProof) []types.OnChainProof {
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
			if (*proofs)[i].Sender == challenges[i].Sender && (*proofs)[i].Height == challenges[i].Height &&
				(*proofs)[i].Round == challenges[i].Round && (*proofs)[i].Rule == challenges[i].Rule &&
				(*proofs)[i].MsgType == challenges[i].MsgType {
				present = true
			}
		}
		if !present {
			result = append(result, (*proofs)[i])
		}
	}

	return result
}

// HandleConsensusMsg is called by p2p protocol manager to deliver the consensus msg to afd.
func (fd *FaultDetector) HandleConsensusMsg(addr common.Address, msg p2p.Msg) error {
	if msg.Code != types.TendermintMsg {
		return nil
	}

	// post msg to afd event loop at once.
	return nil
}
