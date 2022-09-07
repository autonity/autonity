package test

import (
	"crypto/rand"
	"fmt"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core/messageutils"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"math/big"
	"reflect"
)

var autonityContractAddr = crypto.CreateAddress(common.Address{}, 0)

func MsgPropose(address common.Address, block *types.Block, h uint64, r int64, vr int64) *messageutils.Message {
	proposal := messageutils.NewProposal(r, new(big.Int).SetUint64(h), vr, block)
	v, err := messageutils.Encode(proposal)
	if err != nil {
		return nil
	}
	return &messageutils.Message{
		Code:          messageutils.MsgProposal,
		Msg:           v,
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
