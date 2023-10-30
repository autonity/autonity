package e2e

import (
	"crypto/rand"
	"reflect"
	"testing"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/crypto"
	"github.com/stretchr/testify/require"
)

var AutonityContractAddr = crypto.CreateAddress(common.Address{}, 0)
var NonNilValue = common.Hash{0x1}

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
