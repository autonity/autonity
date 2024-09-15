package collusion

import (
	"math/big"
	"math/rand"
	"sync"

	"github.com/autonity/autonity/log"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/cmd/gengen/gengen"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/bft"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	e2e "github.com/autonity/autonity/e2e_test"
)

/**
 * Collusion test, it simulates 1/3 faulty validators working together to manipulate voting of consensus.
 * The collusion framework is base on a faulty party which has 1/3 faulty nodes in it, the framework select a leader out
 * from the faulty nodes, as a leader of faulty nodes, it needs to coordinate the misbehaviour of all faulty followers
 * by according to the accountability rules in the test context, by putting the messages in the queues of each steps of
 * each faulty followers, the followers pick up them from the queue of per step then broadcast it to the committee. Thus,
 * we have such collusion simulated to verify if the accountability module can detect such kind of misbehaviour and to
 * verify the correctness and live-ness of consensus.
 */

type collusionPlaner interface {
	setupRoles(leader *gengen.Validator, followers []*gengen.Validator)
}

type faultyBroadcaster interface {
	SetupCollusionContext()
	BroadcastAll(msg message.Msg)
	Height() *big.Int
	Backend() interfaces.Backend
	Address() common.Address
}

// configurations for different testcases
var collusions = map[autonity.Rule]*collusion{}
var collusionsLock = sync.RWMutex{}

type collusion struct {
	leader       *gengen.Validator
	followers    []*gengen.Validator
	h            uint64
	r            int64
	rule         autonity.Rule
	invalidValue *types.Block
	lock         sync.RWMutex
}

func (c *collusion) context() (uint64, int64, *types.Block) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.h, c.r, c.invalidValue
}

func (c *collusion) setupContext(h uint64, r int64, invalidValue *types.Block) {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.h = h
	c.r = r
	c.invalidValue = invalidValue
}

func (c *collusion) contextReady() bool {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.invalidValue != nil
}

func getCollusion(rule autonity.Rule) *collusion {
	collusionsLock.RLock()
	defer collusionsLock.RUnlock()
	return collusions[rule]
}

func newBlockHeader(height uint64) *types.Header {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[i] = byte(rand.Intn(256)) //nolint
	}
	return &types.Header{
		Number: new(big.Int).SetUint64(height),
		Nonce:  nonce,
	}
}

func initCollusion(vals []*gengen.Validator, rule autonity.Rule, planer collusionPlaner) {
	// it assumes the voting power of validator for collusion test equals to 1 for each.
	f := bft.F(new(big.Int).SetUint64(uint64(len(vals))))
	if f.Uint64() < 2 {
		panic("collusion test requires at least two faulty nodes")
	}
	var faultyMembers []*gengen.Validator
	for i := uint64(0); i < f.Uint64(); i++ {
		faultyMembers = append(faultyMembers, vals[i])
	}

	b := &collusion{
		rule:      rule,
		leader:    faultyMembers[0],
		followers: faultyMembers[1:],
	}
	planer.setupRoles(b.leader, b.followers)

	collusionsLock.Lock()
	defer collusionsLock.Unlock()
	collusions[rule] = b
}

func validProposer(address common.Address, h uint64, r int64, contract *autonity.ProtocolContracts, committee *types.Committee) bool {
	return address == contract.Proposer(committee, nil, h-1, r)
}

func sendPrevote(c *core.Core, rule autonity.Rule) {
	h, r, v := getCollusion(rule).context()
	// if the leader haven't set up the context, skip.
	if v == nil || c.Height().Uint64() < h {
		return
	}

	// send prevote for the planned invalid proposal.
	committee, err := c.Backend().BlockChain().CommitteeByHeight(h)
	if err != nil {
		panic(err)
	}

	if rule == autonity.PVO && h == c.Height().Uint64() {
		log.Debug("prevote collusion simulated", "rule", rule, "h", c.Height(), "r", r, "v", v.Hash(), "node", c.Address())
		// send prevote for the planned invalid proposal for PVO.
		vote := message.NewPrevote(r, h, v.Hash(), c.Backend().Sign, committee.MemberByAddress(c.Address()), committee.Len())
		c.SetSentPrevote(true)
		c.BroadcastAll(vote)
		return
	}

	// send prevote for the planned invalid proposal for PVN
	log.Debug("prevote collusion simulated", "rule", rule, "h", c.Height(), "r", r, "v", v.Hash(), "node", c.Address())
	vote := message.NewPrevote(r, h, v.Hash(), c.Backend().Sign, committee.MemberByAddress(c.Address()), committee.Len())
	c.SetSentPrevote(true)
	c.BroadcastAll(vote)
}

func sendProposal(c faultyBroadcaster, rule autonity.Rule, msg message.Msg) {
	ctx := getCollusion(rule)
	if !ctx.contextReady() {
		c.SetupCollusionContext()
		c.BroadcastAll(msg)
		return
	}

	h, r, v := ctx.context()
	if rule == autonity.PVN && h < c.Height().Uint64() {
		c.SetupCollusionContext()
		c.BroadcastAll(msg)
		return
	}
	if h != c.Height().Uint64() {
		c.BroadcastAll(msg)
		return
	}

	vr := r - 1
	if rule == autonity.PVN {
		vr = -1
	}

	// send invalid proposal with the planed data.
	committee, err := c.Backend().BlockChain().CommitteeByHeight(h)
	if err != nil {
		panic(err)
	}
	p := message.NewPropose(r, h, vr, v, c.Backend().Sign, committee.MemberByAddress(c.Address()))
	c.BroadcastAll(p)
}

func setupCollusionContext(c faultyBroadcaster, rule autonity.Rule) {
	leader := c.Address()
	futureHeight := c.Height().Uint64() + 5
	round := int64(0)
	epoch, _ := c.Backend().BlockChain().LatestEpoch()

	contract := c.Backend().BlockChain().ProtocolContracts()
	for ; ; round++ {
		// select a none proposer to propose faulty value in PVN context.
		if rule == autonity.PVN && !validProposer(leader, futureHeight, round, contract, epoch.Committee) {
			break
		}

		// select a proposer to propose faulty value in PVO and C1 context
		if round != 0 && rule != autonity.PVN && validProposer(leader, futureHeight, round, contract, epoch.Committee) {
			break
		}
	}

	b := types.NewBlockWithHeader(newBlockHeader(futureHeight))
	if rule == autonity.PVN {
		e2e.FuzBlock(b, new(big.Int).SetUint64(futureHeight))
	}

	getCollusion(rule).setupContext(futureHeight, round, b)
	log.Debug("setup collusion context done for", "rule", rule)
}
