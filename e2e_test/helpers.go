package test

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/test"
	"github.com/stretchr/testify/require"
	"math/big"
	"reflect"
	"testing"
)

var AutonityContractAddr = crypto.CreateAddress(common.Address{}, 0)

func NewProposeMsg(address common.Address, block *types.Block, h uint64, r int64, vr int64, sig []byte) *messageutils.Message {
	proposal := messageutils.NewProposal(r, new(big.Int).SetUint64(h), vr, block)
	proposal.LiteSig = sig
	v, err := messageutils.Encode(proposal)
	if err != nil {
		return nil
	}
	return &messageutils.Message{
		Code:          consensus.MsgProposal,
		TbftMsgBytes:  v,
		Address:       address,
		CommittedSeal: []byte{},
	}
}

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GetAllFields(v interface{}) (fieldArr []string) {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Ptr {
		panic("Need pointer!")
	}
	e := reflect.ValueOf(v).Elem()
	for i := 0; i < e.NumField(); i++ {
		fieldArr = append(fieldArr, e.Type().Field(i).Name)
	}
	return fieldArr
}

func GetAllFieldCombinations(v interface{}) (allComb [][]string) {
	fieldSet := GetAllFields(v)
	length := len(fieldSet)
	// total unique combinations are 2^length-1(not including empty)
	for i := 1; i < (1 << length); i++ {
		var comb []string
		for j := 0; j < length; j++ {
			// test if ith bit is set in range of total combinations
			// if yes, we can add to the combination
			if (i>>j)&1 == 1 {
				comb = append(comb, fieldSet[j])
			}
		}
		allComb = append(allComb, comb)
	}
	return allComb
}

func PrintStructMap(oMap map[string]reflect.Value) {
	for key, element := range oMap {
		fmt.Println("Key:", key, "=>", "Element:", element)
	}
}

var NonNilValue = common.Hash{0x1}

func NewVoteMsg(code uint8, h uint64, r int64, v common.Hash, c *core.Core) *messageutils.Message {

	var preVote = messageutils.Vote{
		Round:             r,
		Height:            new(big.Int).SetUint64(h),
		ProposedBlockHash: v,
	}

	encodedVote, err := messageutils.Encode(&preVote)
	if err != nil {
		return nil
	}

	var msg = &messageutils.Message{
		Code:          code,
		TbftMsgBytes:  encodedVote,
		Address:       c.Address(),
		CommittedSeal: []byte{},
	}
	if code == consensus.MsgPrecommit {
		// add committed seal
		seal := helpers.PrepareCommittedSeal(v, r, new(big.Int).SetUint64(h))
		msg.CommittedSeal, err = c.Backend().Sign(seal)
		if err != nil {
			c.Logger().Error("Fault simulator, error while signing committed seal", "err", err)
		}
	}
	return msg
}

// DefaultBehaviour just do the msg gossiping without any simulation.
func DefaultBehaviour(ctx context.Context, c *core.Core, m *messageutils.Message) {
	payload, err := c.FinalizeMessage(m)
	if err != nil {
		return
	}

	if err = c.Backend().Broadcast(ctx, c.CommitteeSet().Committee(), payload); err != nil {
		return
	}
}

func NextProposeRound(currentRound int64, c *core.Core) int64 {
	r := currentRound + 1
	for ; ; r++ {
		p := c.CommitteeSet().GetProposer(r)
		if p.Address == c.Address() {
			break
		}
	}
	return r
}

func DecodeMsg(iMsg *messageutils.Message, c *core.Core) *messageutils.Message {
	payload, _ := c.FinalizeMessage(iMsg)
	m := new(messageutils.Message)
	if err := m.FromPayload(payload); err != nil {
		return nil
	}
	return m
}

func AccountabilityEventDetected(t *testing.T, faultNode common.Address, tp autonity.AccountabilityEventType,
	rule autonity.Rule, network test.Network) bool {

	n := network[1]
	autonityContract, _ := NewAutonity(AutonityContractAddr, n.WsClient)
	var events []AutonityAccountabilityEvent
	var err error
	if tp == autonity.Misbehaviour {
		events, err = autonityContract.GetValidatorRecentMisbehaviours(nil, faultNode)
		require.NoError(t, err)
	} else {
		events, err = autonityContract.GetValidatorRecentAccusations(nil, faultNode)
		require.NoError(t, err)
	}
	presented := false
	for _, e := range events {
		if e.Sender == faultNode && e.Rule == uint8(rule) {
			presented = true
		}
	}
	return presented
}
