package e2e

import (
	"crypto/rand"
	"fmt"
	"github.com/autonity/autonity/params"
	"reflect"
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

	n := network[1]
	autonityContract, _ := autonity.NewAccountability(params.AccountabilityContractAddress, n.WsClient)
	var events []autonity.AccountabilityEvent
	if eventType == autonity.Misbehaviour {
		faults, err := autonityContract.GetValidatorFaults(nil, faultyValidator)
		require.NoError(t, err)
		events = append(events, faults...)
	} else {
		iter, err := autonityContract.FilterNewAccusation(nil, []common.Address{faultyValidator})
		require.NoError(t, err)
		for iter.Next() {
			event, err := autonityContract.Events(nil, iter.Event.Id)
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
