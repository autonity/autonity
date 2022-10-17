package malicious

import (
	"context"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
	e2etest "github.com/autonity/autonity/e2e_test"
	"github.com/autonity/autonity/node"
	"github.com/autonity/autonity/test"
	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

type randomBytesBroadcaster struct {
	*core.Core
}

func (s *randomBytesBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	logger := s.Logger().New("step", s.Step())
	logger.Info("Broadcasting random bytes")

	for i := 0; i < 1000; i++ {
		payload, err := e2etest.GenerateRandomBytes(2048)
		if err != nil {
			logger.Error("Failed to generate random bytes ", "err", err)
			return
		}
		if err = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), payload); err != nil {
			logger.Error("Failed to broadcast message", "msg", msg, "err", err)
			return
		}
	}
}

// TestRandomBytesBroadcaster broadcasts random bytes in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestRandomBytesBroadcaster(t *testing.T) {
	users, err := test.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Broadcaster: &randomBytesBroadcaster{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 180)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type garbageMessageBroadcaster struct {
	*core.Core
}

func (s *garbageMessageBroadcaster) Broadcast(ctx context.Context, msg *messageutils.Message) {
	logger := s.Logger().New("step", s.Step())

	var fMsg messageutils.Message
	f := fuzz.New().NilChance(0.5).Funcs(
		func(cm *messageutils.ConsensusMsg, c fuzz.Continue) {
			c.Fuzz(cm)
		})
	f.Fuzz(&fMsg)

	logger.Info("Broadcasting random bytes")

	payload, err := s.FinalizeMessage(&fMsg)
	if err != nil {
		logger.Error("Failed to finalize message", "msg", fMsg, "err", err)
		return
	}
	if err = s.Backend().Broadcast(ctx, s.CommitteeSet().Committee(), payload); err != nil {
		logger.Error("Failed to broadcast message", "msg", fMsg, "err", err)
		return
	}
}

// TestGarbageMessageBroadcaster broadcasts a garbage Messages in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbageMessageBroadcaster(t *testing.T) {
	users, err := test.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Broadcaster: &garbageMessageBroadcaster{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 180)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type garbagePrecommitSender struct {
	*core.Core
	interfaces.Precommiter
}

func (c *garbagePrecommitSender) SendPrecommit(ctx context.Context, isNil bool) {
	logger := c.Logger().New("step", c.Step())
	var precommit messageutils.Vote
	precommitFieldComb := e2etest.GetAllFieldCombinations(&precommit)
	var msg messageutils.Message
	msgFieldComb := e2etest.GetAllFieldCombinations(&msg)

	proposedBlockHash := common.Hash{}
	if !isNil {
		if h := c.CurRoundMessages().GetProposalHash(); h == (common.Hash{}) {
			c.Logger().Error("Core.sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		proposedBlockHash = c.CurRoundMessages().GetProposalHash()
	}
	//Each iteration tries to fuzz a unique set of fields and skipping
	// a few as provided by fieldsArray
	for _, fieldsArray := range precommitFieldComb {
		// a valid proposal block
		f := fuzz.New().NilChance(0.5)
		f.AllowUnexportedFields(true)
		for _, fieldName := range fieldsArray {
			f.SkipFieldsWithPattern(regexp.MustCompile(fieldName))
		}

		precommit = messageutils.Vote{
			Round:  c.Round(),
			Height: c.Height(),
		}

		precommit.ProposedBlockHash = proposedBlockHash
		// fuzzing existing precommit message, skip the fields in field array
		f.Fuzz(&precommit)
		encodedVote, err := messageutils.Encode(&precommit)
		if err != nil {
			logger.Error("Failed to encode", "subject", precommit)
			return
		}
		msg := &messageutils.Message{
			Code:          messageutils.MsgPrecommit,
			Msg:           encodedVote,
			Address:       c.Address(),
			CommittedSeal: []byte{},
		}
		//Each iteration tries to fuzz a unique set of fields and skipping
		// a few as provided by fieldsArray
		for _, fArray := range msgFieldComb {
			// a valid proposal block
			f := fuzz.New().NilChance(0.5)
			f.AllowUnexportedFields(true)
			for _, fName := range fArray {
				f.SkipFieldsWithPattern(regexp.MustCompile(fName))
			}
			f.Funcs(func(dMsg *messageutils.ConsensusMsg, fc fuzz.Continue) {})
			f.Fuzz(msg)
			// Create committed seal
			seal := helpers.PrepareCommittedSeal(precommit.ProposedBlockHash, c.Round(), c.Height())
			msg.CommittedSeal, err = c.Backend().Sign(seal)
			if err != nil {
				c.Logger().Error("Core.sendPrecommit error while signing committed seal", "err", err)
			}

			c.SetSentPrecommit(true)
			c.Br().Broadcast(ctx, msg)
		}
	}
}

// TestGarbagePrecommitter broadcasts a garbage precommit message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbagePrecommitter(t *testing.T) {
	users, err := test.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Precommitter: &garbagePrecommitSender{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type garbagePrevoter struct {
	*core.Core
	interfaces.Prevoter
}

func (c *garbagePrevoter) SendPrevote(ctx context.Context, isNil bool) {
	logger := c.Logger().New("step", c.Step())

	var prevote messageutils.Vote
	prevoteFieldComb := e2etest.GetAllFieldCombinations(&prevote)
	var msg messageutils.Message
	msgFieldComb := e2etest.GetAllFieldCombinations(&msg)
	proposedBlockHash := c.CurRoundMessages().GetProposalHash()
	if !isNil {
		if h := c.CurRoundMessages().GetProposalHash(); h == (common.Hash{}) {
			c.Logger().Error("sendPrecommit Proposal is empty! It should not be empty!")
			return
		}
		proposedBlockHash = c.CurRoundMessages().GetProposalHash()
	}

	//Each iteration tries to fuzz a unique set of fields and skipping
	// a few as provided by fieldsArray
	for _, fieldsArray := range prevoteFieldComb {
		// a valid proposal block
		f := fuzz.New().NilChance(0.5)
		f.AllowUnexportedFields(true)
		for _, fieldName := range fieldsArray {
			f.SkipFieldsWithPattern(regexp.MustCompile(fieldName))
		}
		prevote = messageutils.Vote{
			Round:  c.Round(),
			Height: c.Height(),
		}
		prevote.ProposedBlockHash = proposedBlockHash
		encodedVote, err := messageutils.Encode(&prevote)
		if err != nil {
			logger.Error("Failed to encode", "subject", prevote)
			return
		}

		msg := &messageutils.Message{
			Code:          messageutils.MsgPrevote,
			Msg:           encodedVote,
			Address:       c.Address(),
			CommittedSeal: []byte{},
		}
		for _, fArray := range msgFieldComb {
			// a valid proposal block
			f := fuzz.New().NilChance(0.5)
			f.AllowUnexportedFields(true)
			for _, fName := range fArray {
				f.SkipFieldsWithPattern(regexp.MustCompile(fName))
			}
			f.Funcs(func(dMsg *messageutils.ConsensusMsg, fc fuzz.Continue) {})
			f.Fuzz(msg)
			c.Br().Broadcast(ctx, msg)
		}
	}
	c.SetSentPrevote(true)
}

// TestGarbagePrevoter broadcasts a garbage prevote message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbagePrevoter(t *testing.T) {
	users, err := test.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Prevoter: &garbagePrevoter{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

type garbageProposer struct {
	*core.Core
	interfaces.Proposer
}

/*
type structNode struct {
	fName string
	sMap  map[string]*structNode
	fList []string
}

func generateFieldMap(v interface{}) map[string]reflect.Value {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		panic("Need pointer!")
	}
	outMap := make(map[string]reflect.Value)
	e := reflect.ValueOf(v).Elem()
	for i := 0; i < e.NumField(); i++ {
		fmt.Println("handling field => ", e.Type().Field(i).Name)
		if e.Field(i).Type().Kind() == reflect.Ptr {
			fKind := e.Field(i).Type().Elem().Kind()
			if fKind == reflect.Struct {
				fmt.Println("TODO - handle recursively")
			}
		} else if e.Field(i).Type().Kind() == reflect.Struct {
			fmt.Println("TODO - handle recursively")
		}
		outMap[e.Type().Field(i).Name] = e.Field(i)
	}
	return outMap
}
*/

func (c *garbageProposer) SendProposal(ctx context.Context, p *types.Block) {
	var proposalMsg messageutils.Proposal
	allComb := e2etest.GetAllFieldCombinations(&proposalMsg)
	//Each iteration tries to fuzz a unique set of fields and skipping
	// a few as provided by fieldsArray
	for _, fieldsArray := range allComb {
		// a valid proposal block
		proposalBlock := messageutils.NewProposal(c.Round(), c.Height(), c.ValidRound(), p)
		f := fuzz.New().NilChance(0.5)
		f.AllowUnexportedFields(true)
		for _, fieldName := range fieldsArray {
			f.SkipFieldsWithPattern(regexp.MustCompile(fieldName))
		}

		f.Funcs(
			func(r *interface{}, fc fuzz.Continue) {},
			func(tr *types.TxData, fc fuzz.Continue) {
				var txData types.LegacyTx
				fc.Fuzz(&txData)
				*tr = &txData
			},
			func(tr **types.Transaction, fc fuzz.Continue) {
				var fakeTransaction types.Transaction
				fc.Fuzz(&fakeTransaction)
				*tr = &fakeTransaction
			},
		)
		f.Fuzz(proposalBlock)
		proposal, _ := messageutils.Encode(proposalBlock)
		c.SetSentProposal(true)
		c.Backend().SetProposedBlockHash(p.Hash())

		c.Br().Broadcast(ctx, &messageutils.Message{
			Code:          messageutils.MsgProposal,
			Msg:           proposal,
			Address:       c.Address(),
			CommittedSeal: []byte{},
		})
	}
}

// TestGarbagePrevoter broadcasts a garbage proposal message in the network,
// We expect other nodes to detect this misbehaviour and discard these messages
// Receiving nodes should also disconnect misbehaving nodes
func TestGarbageProposer(t *testing.T) {
	users, err := test.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].CustHandler = &node.CustomHandler{Proposer: &garbageProposer{}}
	// creates a network of 6 users and starts all the nodes in it
	network, err := test.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown()

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}
