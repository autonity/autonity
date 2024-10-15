package e2e

import (
	"crypto/rand"
	"fmt"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/params"
	fuzz "github.com/google/gofuzz"
	"math/big"
	"reflect"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
)

var NonNilValue = common.Hash{0x1}

var ErrAccountabilityEventMissing = fmt.Errorf("required accountability event is missing")

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

func NextProposeRound(currentRound int64, c *core.Core) int64 {
	for r := currentRound + 1; ; r++ {
		p := c.CommitteeSet().GetProposer(r)
		if p.Address == c.Address() {
			return r
		}
	}
}

func AccountabilityEventDetected(t *testing.T, faultyValidator common.Address, eventType autonity.AccountabilityEventType,
	rule autonity.Rule, network Network) error {

	var lastEpochID int64 = -1
	for _, n := range network {
		header := n.Eth.BlockChain().CurrentHeader()
		db, err := n.Eth.BlockChain().StateAt(header.Root)
		require.NoError(t, err)
		epochID, err := n.Eth.BlockChain().ProtocolContracts().AutonityContract.EpochID(header, db)
		require.NoError(t, err)
		if !epochID.IsInt64() {
			require.Fail(t, "fatal error: epoch id does not fit in int64")
		}
		if lastEpochID == -1 {
			lastEpochID = epochID.Int64()
		} else {
			require.Equal(t, lastEpochID, epochID.Int64(), "epoch id does not match for nodes")
		}
	}

	n := network[1]
	accountabilityContract, _ := autonity.NewAccountability(params.AccountabilityContractAddress, n.WsClient)
	var events []autonity.AccountabilityEvent
	if eventType == autonity.Misbehaviour {
		faults, err := accountabilityContract.GetValidatorFaults(nil, faultyValidator)
		require.NoError(t, err)
		events = append(events, faults...)
	} else {
		iter, err := accountabilityContract.FilterNewAccusation(nil, []common.Address{faultyValidator})
		require.NoError(t, err)
		for iter.Next() {
			event, err := accountabilityContract.Events(nil, iter.Event.Id)
			require.NoError(t, err)
			events = append(events, event)
		}
	}

	found := false
	for _, e := range events {
		if e.Offender == faultyValidator && e.Rule == uint8(rule) {
			found = true
		}
	}

	if !found {
		return ErrAccountabilityEventMissing
	}

	// check if the reporter of accountability events is reimbursed.
	for _, e := range events {
		err := network.CheckReimbursement(e.ReportingBlock.Uint64(), e.Reporter)
		if err != nil {
			return err
		}
	}
	return nil
}

func FuzBlock(p *types.Block, height *big.Int) {
	fakeTransactions := make([]*types.Transaction, 0)
	f := fuzz.New()
	for i := 0; i < 5; i++ {
		var fakeTransaction types.Transaction
		f.Fuzz(&fakeTransaction)
		var tx types.LegacyTx
		f.Fuzz(&tx)
		fakeTransaction.SetInner(&tx)
		fakeTransactions = append(fakeTransactions, &fakeTransaction)
	}
	p.SetTransactions(fakeTransactions)
	var hash common.Hash
	f.Fuzz(&hash)
	var atmHash atomic.Value
	atmHash.Store(hash)
	// nil hash
	p.SetHash(atmHash)
	p.SetHeaderNumber(height)
}
