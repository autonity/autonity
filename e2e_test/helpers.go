package e2e

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/helpers"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/rlp"
	"github.com/stretchr/testify/require"
)

var AutonityContractAddr = crypto.CreateAddress(common.Address{}, 0)
var NonNilValue = common.Hash{0x1}

func NewProposeMsg(address common.Address, block *types.Block, h uint64, r int64, vr int64, signer func([]byte) ([]byte, error)) *message.Message {
	proposal := message.NewProposal(r, new(big.Int).SetUint64(h), vr, block, signer)
	v, err := rlp.EncodeToBytes(proposal)
	if err != nil {
		return nil
	}
	return &message.Message{
		Code:          message.MsgProposal,
		Payload:       v,
		Address:       address,
		CommittedSeal: []byte{},
		ConsensusMsg:  message.ConsensusMsg(proposal),
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

func NewVoteMsg(code uint8, h uint64, r int64, v common.Hash, c *core.Core) *message.Message {
	vote := &message.Vote{
		Round:             r,
		Height:            new(big.Int).SetUint64(h),
		ProposedBlockHash: v,
	}
	encodedVote, _ := rlp.EncodeToBytes(vote)
	msg := &message.Message{
		Code:          code,
		Payload:       encodedVote,
		Address:       c.Address(),
		CommittedSeal: []byte{},
		ConsensusMsg:  message.ConsensusMsg(vote),
	}
	if code == message.MsgPrecommit {
		seal := helpers.PrepareCommittedSeal(v, r, new(big.Int).SetUint64(h))
		msg.CommittedSeal, _ = c.Backend().Sign(seal)
	}
	return msg
}

// DefaultSignAndBroadcast just do the msg gossiping without any simulation.
func DefaultSignAndBroadcast(ctx context.Context, c *core.Core, m *message.Message) {
	payload, err := c.SignMessage(m)
	if err != nil {
		return
	}
	_ = c.Backend().Broadcast(ctx, c.CommitteeSet().Committee(), message.TendermintMessageCode(m), payload)
}

func NextProposeRound(currentRound int64, c *core.Core) int64 {
	for r := currentRound + 1; ; r++ {
		p := c.CommitteeSet().GetProposer(r)
		if p.Address == c.Address() {
			return r
		}
	}
}

func AccountabilityEventDetected(t *testing.T, faultyValidator common.Address, eventType autonity.AccountabilityEventType,
	rule autonity.Rule, network Network) bool {

	n := network[1]
	autonityContract, _ := autonity.NewAccountability(autonity.AccountabilityContractAddress, n.WsClient)
	var events []autonity.AccountabilityEvent
	var err error
	if eventType == autonity.Misbehaviour {
		events, err = autonityContract.GetValidatorFaults(nil, faultyValidator)
	} else {
		var event autonity.AccountabilityEvent
		event, err = autonityContract.GetValidatorAccusation(nil, faultyValidator)
		events = []autonity.AccountabilityEvent{event}
	}
	require.NoError(t, err)
	found := false
	for _, e := range events {
		if e.Offender == faultyValidator && e.Rule == uint8(rule) {
			found = true
		}
	}

	// Go through every block receipt and look for log emitted by the autonity contract

	return found
}
