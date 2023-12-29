// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blake2b"
	"github.com/autonity/autonity/crypto/bls12381"
	"github.com/autonity/autonity/crypto/bn256"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/autonity/autonity/params"

	// lint:ignore SA1019 Needed for precompile
	"golang.org/x/crypto/ripemd160"
)

// PrecompiledContractRWMutex to fix the race condition of precompiled contracts loading on the start-up phase.
var PrecompiledContractRWMutex = sync.RWMutex{}

// PrecompiledContract is the basic interface for native Go contracts. The implementation
// requires a deterministic gas count based on the input size of the Run method of the
// contract.
type PrecompiledContract interface {
	RequiredGas(input []byte) uint64                                                              // RequiredPrice calculates the contract gas use
	Run(input []byte, blockNumber uint64, stateDB StateDB, caller common.Address) ([]byte, error) // Run runs the precompiled contract
}

// PrecompiledContractsHomestead contains the default set of pre-compiled Ethereum
// contracts used in the Frontier and Homestead releases.
var PrecompiledContractsHomestead = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}): &ecrecover{},
	common.BytesToAddress([]byte{2}): &sha256hash{},
	common.BytesToAddress([]byte{3}): &ripemd160hash{},
	common.BytesToAddress([]byte{4}): &dataCopy{},

	common.BytesToAddress([]byte{255}): &checkEnode{},
	common.BytesToAddress([]byte{251}): &QuickSort{},
	common.BytesToAddress([]byte{250}): &QuickSortFast{},
	common.BytesToAddress([]byte{249}): &structTester{},
	common.BytesToAddress([]byte{248}): &ComputeCommitteeReadOnly{},
	common.BytesToAddress([]byte{247}): &QuickSortIterate{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{246}): &QuickSortIterateFast{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{245}): &SortLibrarySliceStable{},
	common.BytesToAddress([]byte{244}): &SortLibrarySort{},
	common.BytesToAddress([]byte{243}): &TestStructLocation{},
	common.BytesToAddress([]byte{242}): &ComputeCommitteeReadAndWrite{},
	common.BytesToAddress([]byte{241}): &ComputeCommitteeReadAndWriteReturnVoters{},
}

// PrecompiledContractsByzantium contains the default set of pre-compiled Ethereum
// contracts used in the Byzantium release.
var PrecompiledContractsByzantium = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}): &ecrecover{},
	common.BytesToAddress([]byte{2}): &sha256hash{},
	common.BytesToAddress([]byte{3}): &ripemd160hash{},
	common.BytesToAddress([]byte{4}): &dataCopy{},
	common.BytesToAddress([]byte{5}): &bigModExp{eip2565: false},
	common.BytesToAddress([]byte{6}): &bn256AddByzantium{},
	common.BytesToAddress([]byte{7}): &bn256ScalarMulByzantium{},
	common.BytesToAddress([]byte{8}): &bn256PairingByzantium{},

	common.BytesToAddress([]byte{255}): &checkEnode{},
	common.BytesToAddress([]byte{251}): &QuickSort{},
	common.BytesToAddress([]byte{250}): &QuickSortFast{},
	common.BytesToAddress([]byte{249}): &structTester{},
	common.BytesToAddress([]byte{248}): &ComputeCommitteeReadOnly{},
	common.BytesToAddress([]byte{247}): &QuickSortIterate{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{246}): &QuickSortIterateFast{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{245}): &SortLibrarySliceStable{},
	common.BytesToAddress([]byte{244}): &SortLibrarySort{},
	common.BytesToAddress([]byte{243}): &TestStructLocation{},
	common.BytesToAddress([]byte{242}): &ComputeCommitteeReadAndWrite{},
	common.BytesToAddress([]byte{241}): &ComputeCommitteeReadAndWriteReturnVoters{},
}

// PrecompiledContractsIstanbul contains the default set of pre-compiled Ethereum
// contracts used in the Istanbul release.
var PrecompiledContractsIstanbul = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}): &ecrecover{},
	common.BytesToAddress([]byte{2}): &sha256hash{},
	common.BytesToAddress([]byte{3}): &ripemd160hash{},
	common.BytesToAddress([]byte{4}): &dataCopy{},
	common.BytesToAddress([]byte{5}): &bigModExp{eip2565: false},
	common.BytesToAddress([]byte{6}): &bn256AddIstanbul{},
	common.BytesToAddress([]byte{7}): &bn256ScalarMulIstanbul{},
	common.BytesToAddress([]byte{8}): &bn256PairingIstanbul{},
	common.BytesToAddress([]byte{9}): &blake2F{},

	common.BytesToAddress([]byte{255}): &checkEnode{},
	common.BytesToAddress([]byte{251}): &QuickSort{},
	common.BytesToAddress([]byte{250}): &QuickSortFast{},
	common.BytesToAddress([]byte{249}): &structTester{},
	common.BytesToAddress([]byte{248}): &ComputeCommitteeReadOnly{},
	common.BytesToAddress([]byte{247}): &QuickSortIterate{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{246}): &QuickSortIterateFast{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{245}): &SortLibrarySliceStable{},
	common.BytesToAddress([]byte{244}): &SortLibrarySort{},
	common.BytesToAddress([]byte{243}): &TestStructLocation{},
	common.BytesToAddress([]byte{242}): &ComputeCommitteeReadAndWrite{},
	common.BytesToAddress([]byte{241}): &ComputeCommitteeReadAndWriteReturnVoters{},
}

// PrecompiledContractsBerlin contains the default set of pre-compiled Ethereum
// contracts used in the Berlin release.
var PrecompiledContractsBerlin = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{1}): &ecrecover{},
	common.BytesToAddress([]byte{2}): &sha256hash{},
	common.BytesToAddress([]byte{3}): &ripemd160hash{},
	common.BytesToAddress([]byte{4}): &dataCopy{},
	common.BytesToAddress([]byte{5}): &bigModExp{eip2565: true},
	common.BytesToAddress([]byte{6}): &bn256AddIstanbul{},
	common.BytesToAddress([]byte{7}): &bn256ScalarMulIstanbul{},
	common.BytesToAddress([]byte{8}): &bn256PairingIstanbul{},
	common.BytesToAddress([]byte{9}): &blake2F{},

	common.BytesToAddress([]byte{255}): &checkEnode{},
	common.BytesToAddress([]byte{251}): &QuickSort{},
	common.BytesToAddress([]byte{250}): &QuickSortFast{},
	common.BytesToAddress([]byte{249}): &structTester{},
	common.BytesToAddress([]byte{248}): &ComputeCommitteeReadOnly{},
	common.BytesToAddress([]byte{247}): &QuickSortIterate{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{246}): &QuickSortIterateFast{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{245}): &SortLibrarySliceStable{},
	common.BytesToAddress([]byte{244}): &SortLibrarySort{},
	common.BytesToAddress([]byte{243}): &TestStructLocation{},
	common.BytesToAddress([]byte{242}): &ComputeCommitteeReadAndWrite{},
	common.BytesToAddress([]byte{241}): &ComputeCommitteeReadAndWriteReturnVoters{},
}

// PrecompiledContractsBLS contains the set of pre-compiled Ethereum
// contracts specified in EIP-2537. These are exported for testing purposes.
var PrecompiledContractsBLS = map[common.Address]PrecompiledContract{
	common.BytesToAddress([]byte{10}): &bls12381G1Add{},
	common.BytesToAddress([]byte{11}): &bls12381G1Mul{},
	common.BytesToAddress([]byte{12}): &bls12381G1MultiExp{},
	common.BytesToAddress([]byte{13}): &bls12381G2Add{},
	common.BytesToAddress([]byte{14}): &bls12381G2Mul{},
	common.BytesToAddress([]byte{15}): &bls12381G2MultiExp{},
	common.BytesToAddress([]byte{16}): &bls12381Pairing{},
	common.BytesToAddress([]byte{17}): &bls12381MapG1{},
	common.BytesToAddress([]byte{18}): &bls12381MapG2{},

	common.BytesToAddress([]byte{255}): &checkEnode{},
	common.BytesToAddress([]byte{251}): &QuickSort{},
	common.BytesToAddress([]byte{250}): &QuickSortFast{},
	common.BytesToAddress([]byte{249}): &structTester{},
	common.BytesToAddress([]byte{248}): &ComputeCommitteeReadOnly{},
	common.BytesToAddress([]byte{247}): &QuickSortIterate{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{246}): &QuickSortIterateFast{queue: make([]*boundary, 0)},
	common.BytesToAddress([]byte{245}): &SortLibrarySliceStable{},
	common.BytesToAddress([]byte{244}): &SortLibrarySort{},
	common.BytesToAddress([]byte{243}): &TestStructLocation{},
	common.BytesToAddress([]byte{242}): &ComputeCommitteeReadAndWrite{},
	common.BytesToAddress([]byte{241}): &ComputeCommitteeReadAndWriteReturnVoters{},
}

var (
	PrecompiledAddressesBerlin    []common.Address
	PrecompiledAddressesIstanbul  []common.Address
	PrecompiledAddressesByzantium []common.Address
	PrecompiledAddressesHomestead []common.Address
)

func init() {
	for k := range PrecompiledContractsHomestead {
		PrecompiledAddressesHomestead = append(PrecompiledAddressesHomestead, k)
	}
	for k := range PrecompiledContractsByzantium {
		PrecompiledAddressesByzantium = append(PrecompiledAddressesByzantium, k)
	}
	for k := range PrecompiledContractsIstanbul {
		PrecompiledAddressesIstanbul = append(PrecompiledAddressesIstanbul, k)
	}
	for k := range PrecompiledContractsBerlin {
		PrecompiledAddressesBerlin = append(PrecompiledAddressesBerlin, k)
	}
}

// ActivePrecompiles returns the precompiles enabled with the current configuration.
func ActivePrecompiles(rules params.Rules) []common.Address {
	switch {
	case rules.IsBerlin:
		return PrecompiledAddressesBerlin
	case rules.IsIstanbul:
		return PrecompiledAddressesIstanbul
	case rules.IsByzantium:
		return PrecompiledAddressesByzantium
	default:
		return PrecompiledAddressesHomestead
	}
}

// RunPrecompiledContract runs and evaluates the output of a precompiled contract.
// It returns
// - the returned bytes,
// - the _remaining_ gas,
// - any error that occurred
func RunPrecompiledContract(
	p PrecompiledContract, input []byte, suppliedGas uint64, blockNumber uint64, stateDB StateDB, caller common.Address,
) (ret []byte, remainingGas uint64, err error) {
	gasCost := p.RequiredGas(input)
	if suppliedGas < gasCost {
		return nil, 0, ErrOutOfGas
	}
	suppliedGas -= gasCost
	output, err := p.Run(input, blockNumber, stateDB, caller)
	return output, suppliedGas, err
}

const ThreadLimit = 10000000

type structTester struct{}

func (a *structTester) RequiredGas(_ []byte) uint64 {
	return 1
}

func (a *structTester) Run(input []byte, _ uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	res := make([]byte, 32)
	for i := 0; i < len(input); i += 32 {
		_, err := fmt.Printf("%v, %v : %v \n", i>>5, &input[i], input[i:i+32])
		if err != nil {
			return res, err
		}
	}
	res[31] = 1
	return res, nil
}

func formatInput(input []byte) []*StakeWithID {
	length := len(input) / 64
	validators := make([]*StakeWithID, length)
	for i := 32; i < len(input); i += 64 {
		idx := uint32(i >> 6)
		item := StakeWithID{
			ValidatorID: idx,
			Stake:       big.NewInt(0).SetBytes(input[i : i+32]),
		}
		validators[idx] = &item
	}
	return validators
}

func formatOutput(validators []*StakeWithID, input []byte) []byte {
	length := len(validators)
	result := make([]byte, length*32+32)
	result[31] = 1
	j := 32
	for i := 0; i < length; i++ {
		idx := validators[i].ValidatorID
		// the address of validator 'idx' is at the slice [idx*64 : idx*64 + 32]
		copy(result[j:j+32], input[(idx<<6):((idx<<6)|32)])
		j += 32
	}
	return result
}

type SortLibrarySliceStable struct{}

func (a *SortLibrarySliceStable) RequiredGas(_ []byte) uint64 {
	return 1
}

func (a *SortLibrarySliceStable) Run(input []byte, _ uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	validators := formatInput(input)

	if len(validators) > 1 {
		sort.SliceStable(validators, func(i, j int) bool {
			return validators[i].Stake.Cmp(validators[j].Stake) == 1
		})
	}

	return formatOutput(validators, input), nil
}

type SortLibrarySort struct{}

type validatorSorter struct {
	validators []*StakeWithID
}

func (sorter *validatorSorter) Len() int {
	return len(sorter.validators)
}

func (sorter *validatorSorter) Less(i, j int) bool {
	return sorter.validators[i].Stake.Cmp(sorter.validators[j].Stake) == 1
}

func (sorter *validatorSorter) Swap(i, j int) {
	sorter.validators[i], sorter.validators[j] = sorter.validators[j], sorter.validators[i]
}

func (a *SortLibrarySort) RequiredGas(_ []byte) uint64 {
	return 1
}

func (a *SortLibrarySort) Run(input []byte, _ uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	validators := formatInput(input)

	sorter := &validatorSorter{validators: validators}
	if len(validators) > 1 {
		sort.Sort(sorter)
	}

	return formatOutput(sorter.validators, input), nil
}

type TestStructLocation struct{}

func (a *TestStructLocation) RequiredGas(_ []byte) uint64 {
	return 1
}

func (a *TestStructLocation) Run(input []byte, _ uint64, stateDB StateDB, caller common.Address) ([]byte, error) {
	fmt.Println(len(input))
	for i := 0; i < len(input); i += 32 {
		fmt.Println(input[i : i+32])
	}
	// step 1: Retrieve all validators from storage

	// question 1: how do we retrieve the contract address?
	//solutions :
	//	- Send it as an input, so we access it by using "input" here
	//	- Retrieve the contract from autonity config
	//  - Add new argument in Run function to pass Caller Contract Reference "ContractRef"
	fmt.Println(caller)
	slot := input[0:32]
	baseOffset := big.NewInt(0).SetBytes(slot)
	for i := 0; i < 10; i++ {
		item := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Bytes()
		fmt.Println(item)
		baseOffset.Add(baseOffset, big.NewInt(1))
	}
	// baseOffset := crypto.Keccak256Hash(validatorListSlot).Big()
	// addresses := make([]common.Address, validatorListSize)

	// // optimisation possible here: introduce concurrency
	// for i := 0; i < validatorListSize; i++ {
	// 	addresses[i] = common.BytesToAddress(stateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes())
	// 	fmt.Println(addresses[i])
	// 	baseOffset.Add(baseOffset, big.NewInt(1))
	// }

	// // We need reference of validator mapping + relative offset of bondedStake + relavtive offset of state

	// validatorsSlot := input[32:64]
	// for _, address := range addresses {
	// 	key := make([]byte, 32)
	// 	copy(key[12:32], address.Bytes())
	// 	baseOffset := crypto.Keccak256Hash(append(key, validatorsSlot...)).Big()
	// 	data := stateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes()
	// 	fmt.Println(common.BytesToAddress(data))
	// 	baseOffset.Add(baseOffset, big.NewInt(1))
	// 	data = stateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes()
	// 	fmt.Println(common.BytesToAddress(data))
	// 	baseOffset.Add(baseOffset, big.NewInt(4))
	// 	data = stateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes()
	// 	fmt.Println(big.NewInt(0).SetBytes(data))
	// 	baseOffset.Add(baseOffset, big.NewInt(3))
	// 	data = stateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes()
	// 	fmt.Println(big.NewInt(0).SetBytes(data))

	// 	// compute storage location of Validator mapping
	// 	// compute absolute location of bondedStake
	// 	// compute absolute location of state
	// 	// retrieve the data
	// }
	// // run quick sort

	// // // Save that directly into storage of validator
	// // // we can use stateDB.SetState()

	// // // SOLIDITY : [SIZE(uint256)]
	// // //
	// // // step 2: Run committee selection

	// // // step3: Save data directly into storage from here

	// // // if success : return true
	res := make([]byte, 32)
	res[31] = 1
	return res, nil
}

// type ComputeCommitteeReadOnlyTest struct{}

// func (a *ComputeCommitteeReadOnlyTest) RequiredGas(_ []byte) uint64 {
// 	return 1
// }

// func (a *ComputeCommitteeReadOnlyTest) Run(input []byte, _ uint64, evm *EVM, caller common.Address) ([]byte, error) {
// 	// for i := 0; i < len(input); i += 32 {
// 	// 	fmt.Println(input[i : i+32])
// 	// }
// 	// step 1: Retrieve all validators from storage

// 	// question 1: how do we retrieve the contract address?
// 	//solutions :
// 	//	- Send it as an input, so we access it by using "input" here
// 	//	- Retrieve the contract from autonity config
// 	//  - Add new argument in Run function to pass Caller Contract Reference "ContractRef"
// 	// fmt.Println(caller)
// 	validatorListSlot := input[0:32]
// 	// fmt.Println(validatorListSlot)
// 	validatorListSize := int(evm.StateDB.GetState(caller, common.BytesToHash(validatorListSlot)).Big().Uint64())
// 	// fmt.Println(validatorListSize)
// 	baseOffset := crypto.Keccak256Hash(validatorListSlot).Big()
// 	addresses := make([]common.Address, validatorListSize)

// 	// optimisation possible here: introduce concurrency
// 	for i := 0; i < validatorListSize; i++ {
// 		addresses[i] = common.BytesToAddress(evm.StateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes())
// 		// fmt.Println(addresses[i])
// 		baseOffset.Add(baseOffset, big.NewInt(1))
// 	}

// 	// We need reference of validator mapping + relative offset of bondedStake + relavtive offset of state

// 	validatorsSlot := input[32:64]
// 	validators := make([]*types.CommitteeMember, 0)
// 	threshold := big.NewInt(0)
// 	for _, address := range addresses {
// 		key := make([]byte, 32)
// 		copy(key[12:32], address.Bytes())
// 		baseOffset := crypto.Keccak256Hash(append(key, validatorsSlot...)).Big()
// 		// data := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Bytes()
// 		// fmt.Println(common.BytesToAddress(data))
// 		// bondedStake is at slot 5
// 		baseOffset.Add(baseOffset, big.NewInt(5))
// 		bondedStake := evm.StateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Big()
// 		// fmt.Println(bondedStake)
// 		if bondedStake.Cmp(threshold) == 1 {
// 			validators = append(validators, &types.CommitteeMember{
// 				Address:     address,
// 				VotingPower: bondedStake,
// 			})
// 		}
// 	}
// 	sort.SliceStable(validators, func(i, j int) bool {
// 		return validators[i].VotingPower.Cmp(validators[j].VotingPower) == 1
// 	})

// 	result := make([]byte, 64+len(validators)*32)
// 	result[31] = 1
// 	binary.BigEndian.PutUint32(result[60:64], uint32(len(validators)))
// 	for i := 64; i < len(result); i += 32 {
// 		copy(result[i:i+32], validators[(i>>5)-2].Address.Bytes())
// 	}
// 	return result, nil
// }

type ComputeCommitteeReadOnly struct{}

func (a *ComputeCommitteeReadOnly) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

func (a *ComputeCommitteeReadOnly) Run(input []byte, _ uint64, stateDB StateDB, caller common.Address) ([]byte, error) {
	// for i := 0; i < len(input); i += 32 {
	// 	fmt.Println(input[i : i+32])
	// }
	// step 1: Retrieve all validators from storage

	// question 1: how do we retrieve the contract address?
	//solutions :
	//	- Send it as an input, so we access it by using "input" here
	//	- Retrieve the contract from autonity config
	//  - Add new argument in Run function to pass Caller Contract Reference "ContractRef"
	// fmt.Println(caller)
	validatorListSlot := input[0:32]
	// fmt.Println(validatorListSlot)
	validatorListSize := int(stateDB.GetState(caller, common.BytesToHash(validatorListSlot)).Big().Uint64())
	// fmt.Println(validatorListSize)
	baseOffset := crypto.Keccak256Hash(validatorListSlot).Big()
	addresses := make([]common.Address, validatorListSize)

	// optimisation possible here: introduce concurrency
	for i := 0; i < validatorListSize; i++ {
		addresses[i] = common.BytesToAddress(stateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes())
		// fmt.Println(addresses[i])
		baseOffset.Add(baseOffset, big.NewInt(1))
	}

	// We need reference of validator mapping + relative offset of bondedStake + relavtive offset of state

	validatorsSlot := input[32:64]
	validators := make([]*types.CommitteeMember, 0)
	threshold := big.NewInt(0)
	for _, address := range addresses {
		key := make([]byte, 32)
		copy(key[12:32], address.Bytes())
		baseOffset := crypto.Keccak256Hash(append(key, validatorsSlot...)).Big()
		// data := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Bytes()
		// fmt.Println(common.BytesToAddress(data))
		// bondedStake is at slot 5
		baseOffset.Add(baseOffset, big.NewInt(5))
		bondedStake := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Big()
		// fmt.Println(bondedStake)
		if bondedStake.Cmp(threshold) == 1 {
			validators = append(validators, &types.CommitteeMember{
				Address:     address,
				VotingPower: bondedStake,
			})
		}
	}
	sort.SliceStable(validators, func(i, j int) bool {
		return validators[i].VotingPower.Cmp(validators[j].VotingPower) == 1
	})

	result := make([]byte, 64+len(validators)*32)
	result[31] = 1
	binary.BigEndian.PutUint32(result[60:64], uint32(len(validators)))
	for i := 64; i < len(result); i += 32 {
		copy(result[i+12:i+32], validators[(i>>5)-2].Address.Bytes())
	}
	return result, nil
}

type ComputeCommitteeReadAndWrite struct{}

func (a *ComputeCommitteeReadAndWrite) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

func (a *ComputeCommitteeReadAndWrite) Run(input []byte, _ uint64, stateDB StateDB, caller common.Address) ([]byte, error) {
	// for i := 0; i < len(input); i += 32 {
	// 	fmt.Println(input[i : i+32])
	// }
	// step 1: Retrieve all validators from storage

	// question 1: how do we retrieve the contract address?
	//solutions :
	//	- Send it as an input, so we access it by using "input" here
	//	- Retrieve the contract from autonity config
	//  - Add new argument in Run function to pass Caller Contract Reference "ContractRef"
	// fmt.Println(caller)
	validatorListSlot := input[0:32]
	validatorsSlot := input[32:64]
	committeeSlot := input[64:96]
	committeeLenConfig := binary.BigEndian.Uint32(input[124:128])
	// fmt.Println(validatorListSlot)
	validatorListSize := int(stateDB.GetState(caller, common.BytesToHash(validatorListSlot)).Big().Uint64())
	// fmt.Println(validatorListSize)
	baseOffset := crypto.Keccak256Hash(validatorListSlot).Big()
	addresses := make([]common.Address, validatorListSize)

	// optimisation possible here: introduce concurrency
	for i := 0; i < validatorListSize; i++ {
		addresses[i] = common.BytesToAddress(stateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes())
		// fmt.Println(addresses[i])
		baseOffset.Add(baseOffset, big.NewInt(1))
	}

	// We need reference of validator mapping + relative offset of bondedStake + relavtive offset of state

	validators := make([]*types.CommitteeMember, 0)
	threshold := big.NewInt(0)
	for _, address := range addresses {
		key := make([]byte, 32)
		copy(key[12:32], address.Bytes())
		baseOffset := crypto.Keccak256Hash(append(key, validatorsSlot...)).Big()
		// data := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Bytes()
		// fmt.Println(common.BytesToAddress(data))
		// bondedStake is at slot 5
		baseOffset.Add(baseOffset, big.NewInt(5))
		bondedStake := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Big()
		// fmt.Println(bondedStake)
		if bondedStake.Cmp(threshold) == 1 {
			validators = append(validators, &types.CommitteeMember{
				Address:     address,
				VotingPower: bondedStake,
			})
		}
	}
	sort.SliceStable(validators, func(i, j int) bool {
		return validators[i].VotingPower.Cmp(validators[j].VotingPower) == 1
	})

	committeeSize := int(stateDB.GetState(caller, common.BytesToHash(committeeSlot)).Big().Uint64())
	// fmt.Println(validatorListSize)
	baseOffset = crypto.Keccak256Hash(committeeSlot).Big()

	// delete old committee members : type CommitteeMember
	for i := 0; i < committeeSize; i++ {
		// delete address
		stateDB.SetState(caller, common.Hash(baseOffset.Bytes()), common.Hash{})
		baseOffset.Add(baseOffset, big.NewInt(1))
		// delete voting power
		stateDB.SetState(caller, common.Hash(baseOffset.Bytes()), common.Hash{})
		baseOffset.Add(baseOffset, big.NewInt(1))
	}
	if committeeLenConfig > uint32(len(validators)) {
		committeeSize = len(validators)
	} else {
		committeeSize = int(committeeLenConfig)
	}

	// put new committee members : type CommitteeMember
	// voters := make([]byte, committeeSize*32)
	// 4 for uint32
	committeeLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(committeeLenBytes, uint32(committeeSize))
	// save committeeSize in committeeSlot
	stateDB.SetState(caller, common.BytesToHash(committeeSlot), common.BytesToHash(committeeLenBytes))
	baseOffset = crypto.Keccak256Hash(committeeSlot).Big()
	for i := 0; i < committeeSize; i++ {
		// save address
		stateDB.SetState(caller, common.Hash(baseOffset.Bytes()), validators[i].Address.Hash())
		baseOffset.Add(baseOffset, big.NewInt(1))
		// save voting power
		stateDB.SetState(caller, common.Hash(baseOffset.Bytes()), common.BytesToHash(validators[i].VotingPower.Bytes()))
		baseOffset.Add(baseOffset, big.NewInt(1))
	}

	result := make([]byte, 32)
	result[31] = 1
	return result, nil
}

type ComputeCommitteeReadAndWriteReturnVoters struct{}

func (a *ComputeCommitteeReadAndWriteReturnVoters) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

func (a *ComputeCommitteeReadAndWriteReturnVoters) Run(input []byte, _ uint64, stateDB StateDB, caller common.Address) ([]byte, error) {
	// for i := 0; i < len(input); i += 32 {
	// 	fmt.Println(input[i : i+32])
	// }
	// step 1: Retrieve all validators from storage

	// question 1: how do we retrieve the contract address?
	//solutions :
	//	- Send it as an input, so we access it by using "input" here
	//	- Retrieve the contract from autonity config
	//  - Add new argument in Run function to pass Caller Contract Reference "ContractRef"
	// fmt.Println(caller)
	validatorListSlot := input[0:32]
	validatorsSlot := input[32:64]
	committeeSlot := input[64:96]
	epochBondedSlot := input[96:128]
	committeeLenConfig := binary.BigEndian.Uint32(input[156:160])
	// fmt.Println(validatorListSlot)
	validatorListSize := int(stateDB.GetState(caller, common.BytesToHash(validatorListSlot)).Big().Uint64())
	// fmt.Println(validatorListSize)
	baseOffset := crypto.Keccak256Hash(validatorListSlot).Big()
	addresses := make([]common.Address, validatorListSize)

	// optimisation possible here: introduce concurrency
	for i := 0; i < validatorListSize; i++ {
		addresses[i] = common.BytesToAddress(stateDB.GetState(caller, common.Hash(baseOffset.Bytes())).Bytes())
		// fmt.Println(addresses[i])
		baseOffset.Add(baseOffset, big.NewInt(1))
	}

	// We need reference of validator mapping + relative offset of bondedStake + relavtive offset of state

	validators := make([]*types.CommitteeMember, 0)
	threshold := big.NewInt(0)
	for _, address := range addresses {
		key := make([]byte, 32)
		copy(key[12:32], address.Bytes())
		baseOffset := crypto.Keccak256Hash(append(key, validatorsSlot...)).Big()
		// data := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Bytes()
		// fmt.Println(common.BytesToAddress(data))
		// bondedStake is at slot 5
		baseOffset.Add(baseOffset, big.NewInt(5))
		bondedStake := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Big()
		// fmt.Println(bondedStake)
		if bondedStake.Cmp(threshold) == 1 {
			validators = append(validators, &types.CommitteeMember{
				Address:     address,
				VotingPower: bondedStake,
			})
		}
	}
	sort.SliceStable(validators, func(i, j int) bool {
		return validators[i].VotingPower.Cmp(validators[j].VotingPower) == 1
	})

	committeeSize := int(stateDB.GetState(caller, common.BytesToHash(committeeSlot)).Big().Uint64())
	// fmt.Println(validatorListSize)
	baseOffset = crypto.Keccak256Hash(committeeSlot).Big()

	// delete old committee members : type CommitteeMember
	for i := 0; i < committeeSize; i++ {
		// delete address
		stateDB.SetState(caller, common.Hash(baseOffset.Bytes()), common.Hash{})
		baseOffset.Add(baseOffset, big.NewInt(1))
		// delete voting power
		stateDB.SetState(caller, common.Hash(baseOffset.Bytes()), common.Hash{})
		baseOffset.Add(baseOffset, big.NewInt(1))
	}
	if committeeLenConfig > uint32(len(validators)) {
		committeeSize = len(validators)
	} else {
		committeeSize = int(committeeLenConfig)
	}

	// put new committee members : type CommitteeMember
	voters := make([]byte, committeeSize*32)
	// 4 for uint32
	committeeLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(committeeLenBytes, uint32(committeeSize))
	// save committeeSize in committeeSlot
	stateDB.SetState(caller, common.BytesToHash(committeeSlot), common.BytesToHash(committeeLenBytes))
	baseOffset = crypto.Keccak256Hash(committeeSlot).Big()
	totalStake := big.NewInt(0)
	for i := 0; i < committeeSize; i++ {
		// save address
		stateDB.SetState(caller, common.Hash(baseOffset.Bytes()), validators[i].Address.Hash())
		baseOffset.Add(baseOffset, big.NewInt(1))
		// save voting power
		stateDB.SetState(caller, common.Hash(baseOffset.Bytes()), common.BytesToHash(validators[i].VotingPower.Bytes()))
		baseOffset.Add(baseOffset, big.NewInt(1))

		totalStake = totalStake.Add(totalStake, validators[i].VotingPower)

		// get oracleAddress
		key := make([]byte, 32)
		copy(key[12:32], validators[i].Address.Bytes())
		mapItemOffset := crypto.Keccak256Hash(append(key, validatorsSlot...)).Big()
		// oracleAddress is at slot 2
		mapItemOffset.Add(mapItemOffset, big.NewInt(2))
		oracleAddress := stateDB.GetState(caller, common.BytesToHash(mapItemOffset.Bytes())).Bytes()
		// voters[i*32:i*32+32] will store oracleAddress of i'th member
		copy(voters[(i<<5):(i<<5)+32], oracleAddress)
	}
	// write epochTotalBondedStake
	stateDB.SetState(caller, common.BytesToHash(epochBondedSlot), common.BytesToHash(totalStake.Bytes()))

	result := make([]byte, 64)
	result[31] = 1
	binary.BigEndian.PutUint32(result[60:64], uint32(committeeSize))
	result = append(result, voters...)
	return result, nil
}

type TestCommitteeRead struct{}

func (a *TestCommitteeRead) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

func (a *TestCommitteeRead) Run(input []byte, _ uint64, stateDB StateDB, caller common.Address) ([]byte, error) {
	// for i := 0; i < len(input); i += 32 {
	// 	fmt.Println(input[i : i+32])
	// }
	// step 1: Retrieve all validators from storage

	// question 1: how do we retrieve the contract address?
	//solutions :
	//	- Send it as an input, so we access it by using "input" here
	//	- Retrieve the contract from autonity config
	//  - Add new argument in Run function to pass Caller Contract Reference "ContractRef"
	// fmt.Println(caller)
	committeeSlot := input[0:32]
	validatorsSlot := input[32:64]
	fmt.Println(big.NewInt(0).SetBytes(committeeSlot))
	committeeSize := int(stateDB.GetState(caller, common.BytesToHash(committeeSlot)).Big().Uint64())
	fmt.Println(committeeSize)
	baseOffset := crypto.Keccak256Hash(committeeSlot).Big()

	zeroAddress := common.BytesToAddress(make([]byte, 0))
	for i := 0; i < committeeSize; i++ {
		fmt.Println(baseOffset)
		address := common.BytesToAddress(stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Bytes())
		fmt.Printf("addres %v\t", address)
		baseOffset.Add(baseOffset, big.NewInt(1))
		votingPower := stateDB.GetState(caller, common.BytesToHash(baseOffset.Bytes())).Big()
		fmt.Printf("votingPower %v\n", votingPower)
		baseOffset.Add(baseOffset, big.NewInt(1))
		if bytes.Equal(address.Bytes(), zeroAddress.Bytes()) {
			return make([]byte, 0), fmt.Errorf("zero address")
		}
		if votingPower.Cmp(big.NewInt(0)) <= 0 {
			return make([]byte, 0), fmt.Errorf("stake not positive")
		}

		key := make([]byte, 32)
		copy(key[12:32], address.Bytes())
		pairOffset := crypto.Keccak256Hash(append(key, validatorsSlot...)).Big()
		pairOffset.Add(pairOffset, big.NewInt(5))
		bondedStake := stateDB.GetState(caller, common.BytesToHash(pairOffset.Bytes())).Big()
		if bondedStake.Cmp(votingPower) != 0 {
			return make([]byte, 0), fmt.Errorf("stake mismatch")
		}
	}
	result := make([]byte, 32)
	result[31] = 1
	return result, nil
}

const KB = 1024

type QuickSort struct{}

type StakeWithID struct {
	ValidatorID uint32
	Stake       *big.Int
}

// RequiredGas the gas cost to sort validator list
func (a *QuickSort) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the validator list and sort it according to bonded stake in descending order
// and then returns the addresses only to reduce the memory usage
func (a *QuickSort) Run(input []byte, _ uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	length := len(input) / 64
	validators := formatInput(input)

	if length > 1 {
		structQuickSort(validators, 0, int32(length)-1)
	}

	return formatOutput(validators, input), nil
}

func structQuickSort(validators []*StakeWithID, low int32, high int32) {
	// Set the pivot element in its right sorted index in the array
	pivot := validators[(high+low)/2].Stake
	// isLeftSorted stores if the left subarray with indexes [left, i-1] sorted or not
	isLeftSorted := true
	// isRightSorted stores if the right subarray with indexes [j+1, right] sorted or not
	isRightSorted := true
	i := low
	j := high
	for i <= j {
		for validators[i].Stake.Cmp(pivot) == 1 {
			i++
			// check if elements at (i-1) and (i-2) are sorted or not
			if isLeftSorted && i-1 > low {
				isLeftSorted = validators[i-2].Stake.Cmp(validators[i-1].Stake) >= 0
			}
		}
		for pivot.Cmp(validators[j].Stake) == 1 {
			j--
			// check if elements at (j+1) and (j+2) are sorted or not
			if isRightSorted && j+1 < high {
				isRightSorted = validators[j+1].Stake.Cmp(validators[j+2].Stake) >= 0
			}
		}
		if i <= j {
			validators[i], validators[j] = validators[j], validators[i]
			i++
			j--
			if isLeftSorted && i-1 > low {
				isLeftSorted = validators[i-2].Stake.Cmp(validators[i-1].Stake) >= 0
			}
			if isRightSorted && j+1 < high {
				isRightSorted = validators[j+1].Stake.Cmp(validators[j+2].Stake) >= 0
			}
		}
	}
	// Recursion call in the left partition of the array
	if !isLeftSorted && low < j {
		structQuickSort(validators, low, j)
	}
	// Recursion call in the right partition
	if !isRightSorted && i < high {
		structQuickSort(validators, i, high)
	}
}

type QuickSortFast struct{}

// RequiredGas the gas cost to sort validator list
func (a *QuickSortFast) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the validator list and sort it according to bonded stake in descending order
// and then returns the addresses only to reduce the memory usage
func (a *QuickSortFast) Run(input []byte, _ uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	length := len(input) / 64
	validators := formatInput(input)

	if length > 1 {
		task := sync.WaitGroup{}
		threadUsed := 0
		structQuickSortFast(validators, 0, int32(length)-1, &task, &threadUsed)
		task.Wait()
	}

	return formatOutput(validators, input), nil
}

func structQuickSortFast(validators []*StakeWithID, low int32, high int32, task *sync.WaitGroup, threadUsed *int) {
	// Set the pivot element in its right sorted index in the array
	pivot := validators[(high+low)/2].Stake
	// isLeftSorted stores if the left subarray with indexes [left, i-1] sorted or not
	isLeftSorted := true
	// isRightSorted stores if the right subarray with indexes [j+1, right] sorted or not
	isRightSorted := true
	i := low
	j := high
	for i <= j {
		for validators[i].Stake.Cmp(pivot) == 1 {
			i++
			// check if elements at (i-1) and (i-2) are sorted or not
			if isLeftSorted && i-1 > low {
				isLeftSorted = validators[i-2].Stake.Cmp(validators[i-1].Stake) >= 0
			}
		}
		for pivot.Cmp(validators[j].Stake) == 1 {
			j--
			// check if elements at (j+1) and (j+2) are sorted or not
			if isRightSorted && j+1 < high {
				isRightSorted = validators[j+1].Stake.Cmp(validators[j+2].Stake) >= 0
			}
		}
		if i <= j {
			validators[i], validators[j] = validators[j], validators[i]
			i++
			j--
			if isLeftSorted && i-1 > low {
				isLeftSorted = validators[i-2].Stake.Cmp(validators[i-1].Stake) >= 0
			}
			if isRightSorted && j+1 < high {
				isRightSorted = validators[j+1].Stake.Cmp(validators[j+2].Stake) >= 0
			}
		}
	}

	// task := sync.WaitGroup{}

	if !isLeftSorted && low < j {
		if *threadUsed < ThreadLimit {
			task.Add(1)
			(*threadUsed)++
			go func() {
				// Recursion call in the left partition of the array
				structQuickSortFast(validators, low, j, task, threadUsed)
				task.Done()
				(*threadUsed)--
			}()
		} else {
			// Recursion call in the left partition of the array
			structQuickSortFast(validators, low, j, task, threadUsed)
		}
	}
	if !isRightSorted && i < high {
		if *threadUsed < ThreadLimit {
			task.Add(1)
			(*threadUsed)++
			go func() {
				// Recursion call in the right partition
				structQuickSortFast(validators, i, high, task, threadUsed)
				task.Done()
				(*threadUsed)--
			}()
		} else {
			// Recursion call in the right partition
			structQuickSortFast(validators, i, high, task, threadUsed)
		}
	}
}

type QuickSortIterate struct {
	queue []*boundary
}

func NewQuickSortIterate() *QuickSortIterate {
	return &QuickSortIterate{queue: make([]*boundary, 0)}
}

type boundary struct {
	low  int32
	high int32
}

// RequiredGas the gas cost to sort validator list
func (a *QuickSortIterate) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the validator list and sort it according to bonded stake in descending order
// and then returns the addresses only to reduce the memory usage
func (sort *QuickSortIterate) Run(input []byte, _ uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	validators := formatInput(input)

	if len(validators) > 1 {
		sort.structQuickSortIterate(validators)
	}

	return formatOutput(validators, input), nil
}

func (sort *QuickSortIterate) partition(validators []*StakeWithID, low int32, high int32) {
	// isLeftSorted stores if the left subarray with indexes [left, i-1] sorted or not
	isLeftSorted := true
	// isRightSorted stores if the right subarray with indexes [j+1, right] sorted or not
	isRightSorted := true
	i := low
	j := high
	// Set the pivot element in its right sorted index in the array
	pivot := validators[(high+low)/2].Stake
	for i <= j {
		for validators[i].Stake.Cmp(pivot) == 1 {
			i++
			// check if elements at (i-1) and (i-2) are sorted or not
			if isLeftSorted && i-1 > low {
				isLeftSorted = validators[i-2].Stake.Cmp(validators[i-1].Stake) >= 0
			}
		}
		for pivot.Cmp(validators[j].Stake) == 1 {
			j--
			// check if elements at (j+1) and (j+2) are sorted or not
			if isRightSorted && j+1 < high {
				isRightSorted = validators[j+1].Stake.Cmp(validators[j+2].Stake) >= 0
			}
		}
		if i <= j {
			validators[i], validators[j] = validators[j], validators[i]
			i++
			j--
			if isLeftSorted && i-1 > low {
				isLeftSorted = validators[i-2].Stake.Cmp(validators[i-1].Stake) >= 0
			}
			if isRightSorted && j+1 < high {
				isRightSorted = validators[j+1].Stake.Cmp(validators[j+2].Stake) >= 0
			}
		}
	}

	if !isLeftSorted && low < j {
		// need to sort the left portion
		sort.queue = append(sort.queue, &boundary{low: low, high: j})
	}
	if !isRightSorted && i < high {
		// need to sort the right portion
		sort.queue = append(sort.queue, &boundary{low: i, high: high})
	}
}

func (sort *QuickSortIterate) structQuickSortIterate(validators []*StakeWithID) {
	sort.queue = append(sort.queue, &boundary{low: 0, high: int32(len(validators)) - 1})
	for len(sort.queue) > 0 {
		// pop the first item
		topItem := sort.queue[0]
		sort.queue = sort.queue[1:]
		sort.partition(validators, topItem.low, topItem.high)
	}
}

type QuickSortIterateFast struct {
	queueLock sync.RWMutex
	queue     []*boundary
}

func NewQuickSortIterateFast() *QuickSortIterateFast {
	return &QuickSortIterateFast{queue: make([]*boundary, 0)}
}

// RequiredGas the gas cost to sort validator list
func (a *QuickSortIterateFast) RequiredGas(input []byte) uint64 {
	times := uint64(len(input)/KB + 1)
	return params.AutonityAFDContractGasPerKB * times
}

// Run take the validator list and sort it according to bonded stake in descending order
// and then returns the addresses only to reduce the memory usage
func (sort *QuickSortIterateFast) Run(input []byte, _ uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	validators := formatInput(input)

	if len(validators) > 1 {
		sort.structQuickSortIterateFast(validators)
	}

	return formatOutput(validators, input), nil
}

func (sort *QuickSortIterateFast) partition(validators []*StakeWithID, low int32, high int32) {
	// isLeftSorted stores if the left subarray with indexes [left, i-1] sorted or not
	isLeftSorted := true
	// isRightSorted stores if the right subarray with indexes [j+1, right] sorted or not
	isRightSorted := true
	i := low
	j := high
	// Set the pivot element in its right sorted index in the array
	pivot := validators[(high+low)/2].Stake
	for i <= j {
		for validators[i].Stake.Cmp(pivot) == 1 {
			i++
			// check if elements at (i-1) and (i-2) are sorted or not
			if isLeftSorted && i-1 > low {
				isLeftSorted = validators[i-2].Stake.Cmp(validators[i-1].Stake) >= 0
			}
		}
		for pivot.Cmp(validators[j].Stake) == 1 {
			j--
			// check if elements at (j+1) and (j+2) are sorted or not
			if isRightSorted && j+1 < high {
				isRightSorted = validators[j+1].Stake.Cmp(validators[j+2].Stake) >= 0
			}
		}
		if i <= j {
			validators[i], validators[j] = validators[j], validators[i]
			i++
			j--
			if isLeftSorted && i-1 > low {
				isLeftSorted = validators[i-2].Stake.Cmp(validators[i-1].Stake) >= 0
			}
			if isRightSorted && j+1 < high {
				isRightSorted = validators[j+1].Stake.Cmp(validators[j+2].Stake) >= 0
			}
		}
	}

	if !isLeftSorted && low < j {
		// need to sort the left portion
		sort.queueLock.Lock()
		sort.queue = append(sort.queue, &boundary{low: low, high: j})
		sort.queueLock.Unlock()
	}
	if !isRightSorted && i < high {
		// need to sort the right portion
		sort.queueLock.Lock()
		sort.queue = append(sort.queue, &boundary{low: i, high: high})
		sort.queueLock.Unlock()
	}
}

func (sort *QuickSortIterateFast) structQuickSortIterateFast(validators []*StakeWithID) {
	sort.queue = append(sort.queue, &boundary{low: 0, high: int32(len(validators)) - 1})
	task := sync.WaitGroup{}
	for {
		threadUsed := 0
		for len(sort.queue) > 0 && threadUsed < ThreadLimit {
			// pop the first item
			sort.queueLock.Lock()
			topItem := sort.queue[0]
			sort.queue = sort.queue[1:]
			sort.queueLock.Unlock()

			task.Add(1)
			threadUsed++
			go func() {
				sort.partition(validators, topItem.low, topItem.high)
				task.Done()
				threadUsed--
			}()
		}
		task.Wait()
		if len(sort.queue) == 0 {
			break
		}
	}
}

// ECRECOVER implemented as a native contract.
type ecrecover struct{}

func (c *ecrecover) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (c *ecrecover) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	const ecRecoverInputLength = 128

	input = common.RightPadBytes(input, ecRecoverInputLength)
	// "input" is (hash, v, r, s), each 32 bytes
	// but for ecrecover we want (r, s, v)

	r := new(big.Int).SetBytes(input[64:96])
	s := new(big.Int).SetBytes(input[96:128])
	v := input[63] - 27

	// tighter sig s values input homestead only apply to tx sigs
	if !allZero(input[32:63]) || !crypto.ValidateSignatureValues(v, r, s, false) {
		return nil, nil
	}
	// We must make sure not to modify the 'input', so placing the 'v' along with
	// the signature needs to be done on a new allocation
	sig := make([]byte, 65)
	copy(sig, input[64:128])
	sig[64] = v
	// v needs to be at the end for libsecp256k1
	pubKey, err := crypto.Ecrecover(input[:32], sig)
	// make sure the public key is a valid one
	if err != nil {
		return nil, nil
	}

	// the first byte of pubkey is bitcoin heritage
	return common.LeftPadBytes(crypto.Keccak256(pubKey[1:])[12:], 32), nil
}

// SHA256 implemented as a native contract.
type sha256hash struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *sha256hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*params.Sha256PerWordGas + params.Sha256BaseGas
}
func (c *sha256hash) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	h := sha256.Sum256(input)
	return h[:], nil
}

// RIPEMD160 implemented as a native contract.
type ripemd160hash struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *ripemd160hash) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*params.Ripemd160PerWordGas + params.Ripemd160BaseGas
}
func (c *ripemd160hash) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	ripemd := ripemd160.New()
	ripemd.Write(input)
	return common.LeftPadBytes(ripemd.Sum(nil), 32), nil
}

// data copy implemented as a native contract.
type dataCopy struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
//
// This method does not require any overflow checking as the input size gas costs
// required for anything significant is so high it's impossible to pay for.
func (c *dataCopy) RequiredGas(input []byte) uint64 {
	return uint64(len(input)+31)/32*params.IdentityPerWordGas + params.IdentityBaseGas
}
func (c *dataCopy) Run(in []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	return in, nil
}

// bigModExp implements a native big integer exponential modular operation.
type bigModExp struct {
	eip2565 bool
}

var (
	big0      = big.NewInt(0)
	big1      = big.NewInt(1)
	big3      = big.NewInt(3)
	big4      = big.NewInt(4)
	big7      = big.NewInt(7)
	big8      = big.NewInt(8)
	big16     = big.NewInt(16)
	big20     = big.NewInt(20)
	big32     = big.NewInt(32)
	big64     = big.NewInt(64)
	big96     = big.NewInt(96)
	big480    = big.NewInt(480)
	big1024   = big.NewInt(1024)
	big3072   = big.NewInt(3072)
	big199680 = big.NewInt(199680)
)

// modexpMultComplexity implements bigModexp multComplexity formula, as defined in EIP-198
//
// def mult_complexity(x):
//
//	if x <= 64: return x ** 2
//	elif x <= 1024: return x ** 2 // 4 + 96 * x - 3072
//	else: return x ** 2 // 16 + 480 * x - 199680
//
// where is x is max(length_of_MODULUS, length_of_BASE)
func modexpMultComplexity(x *big.Int) *big.Int {
	switch {
	case x.Cmp(big64) <= 0:
		x.Mul(x, x) // x ** 2
	case x.Cmp(big1024) <= 0:
		// (x ** 2 // 4 ) + ( 96 * x - 3072)
		x = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(x, x), big4),
			new(big.Int).Sub(new(big.Int).Mul(big96, x), big3072),
		)
	default:
		// (x ** 2 // 16) + (480 * x - 199680)
		x = new(big.Int).Add(
			new(big.Int).Div(new(big.Int).Mul(x, x), big16),
			new(big.Int).Sub(new(big.Int).Mul(big480, x), big199680),
		)
	}
	return x
}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bigModExp) RequiredGas(input []byte) uint64 {
	var (
		baseLen = new(big.Int).SetBytes(getData(input, 0, 32))
		expLen  = new(big.Int).SetBytes(getData(input, 32, 32))
		modLen  = new(big.Int).SetBytes(getData(input, 64, 32))
	)
	if len(input) > 96 {
		input = input[96:]
	} else {
		input = input[:0]
	}
	// Retrieve the head 32 bytes of exp for the adjusted exponent length
	var expHead *big.Int
	if big.NewInt(int64(len(input))).Cmp(baseLen) <= 0 {
		expHead = new(big.Int)
	} else {
		if expLen.Cmp(big32) > 0 {
			expHead = new(big.Int).SetBytes(getData(input, baseLen.Uint64(), 32))
		} else {
			expHead = new(big.Int).SetBytes(getData(input, baseLen.Uint64(), expLen.Uint64()))
		}
	}
	// Calculate the adjusted exponent length
	var msb int
	if bitlen := expHead.BitLen(); bitlen > 0 {
		msb = bitlen - 1
	}
	adjExpLen := new(big.Int)
	if expLen.Cmp(big32) > 0 {
		adjExpLen.Sub(expLen, big32)
		adjExpLen.Mul(big8, adjExpLen)
	}
	adjExpLen.Add(adjExpLen, big.NewInt(int64(msb)))
	// Calculate the gas cost of the operation
	gas := new(big.Int).Set(math.BigMax(modLen, baseLen))
	if c.eip2565 {
		// EIP-2565 has three changes
		// 1. Different multComplexity (inlined here)
		// in EIP-2565 (https://eips.ethereum.org/EIPS/eip-2565):
		//
		// def mult_complexity(x):
		//    ceiling(x/8)^2
		//
		//where is x is max(length_of_MODULUS, length_of_BASE)
		gas = gas.Add(gas, big7)
		gas = gas.Div(gas, big8)
		gas.Mul(gas, gas)

		gas.Mul(gas, math.BigMax(adjExpLen, big1))
		// 2. Different divisor (`GQUADDIVISOR`) (3)
		gas.Div(gas, big3)
		if gas.BitLen() > 64 {
			return math.MaxUint64
		}
		// 3. Minimum price of 200 gas
		if gas.Uint64() < 200 {
			return 200
		}
		return gas.Uint64()
	}
	gas = modexpMultComplexity(gas)
	gas.Mul(gas, math.BigMax(adjExpLen, big1))
	gas.Div(gas, big20)

	if gas.BitLen() > 64 {
		return math.MaxUint64
	}
	return gas.Uint64()
}

func (c *bigModExp) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	var (
		baseLen = new(big.Int).SetBytes(getData(input, 0, 32)).Uint64()
		expLen  = new(big.Int).SetBytes(getData(input, 32, 32)).Uint64()
		modLen  = new(big.Int).SetBytes(getData(input, 64, 32)).Uint64()
	)
	if len(input) > 96 {
		input = input[96:]
	} else {
		input = input[:0]
	}
	// Handle a special case when both the base and mod length is zero
	if baseLen == 0 && modLen == 0 {
		return []byte{}, nil
	}
	// Retrieve the operands and execute the exponentiation
	var (
		base = new(big.Int).SetBytes(getData(input, 0, baseLen))
		exp  = new(big.Int).SetBytes(getData(input, baseLen, expLen))
		mod  = new(big.Int).SetBytes(getData(input, baseLen+expLen, modLen))
	)
	if mod.BitLen() == 0 {
		// Modulo 0 is undefined, return zero
		return common.LeftPadBytes([]byte{}, int(modLen)), nil
	}
	return common.LeftPadBytes(base.Exp(base, exp, mod).Bytes(), int(modLen)), nil
}

// newCurvePoint unmarshals a binary blob into a bn256 elliptic curve point,
// returning it, or an error if the point is invalid.
func newCurvePoint(blob []byte) (*bn256.G1, error) {
	p := new(bn256.G1)
	if _, err := p.Unmarshal(blob); err != nil {
		return nil, err
	}
	return p, nil
}

// newTwistPoint unmarshals a binary blob into a bn256 elliptic curve point,
// returning it, or an error if the point is invalid.
func newTwistPoint(blob []byte) (*bn256.G2, error) {
	p := new(bn256.G2)
	if _, err := p.Unmarshal(blob); err != nil {
		return nil, err
	}
	return p, nil
}

// runBn256Add implements the Bn256Add precompile, referenced by both
// Byzantium and Istanbul operations.
func runBn256Add(input []byte) ([]byte, error) {
	x, err := newCurvePoint(getData(input, 0, 64))
	if err != nil {
		return nil, err
	}
	y, err := newCurvePoint(getData(input, 64, 64))
	if err != nil {
		return nil, err
	}
	res := new(bn256.G1)
	res.Add(x, y)
	return res.Marshal(), nil
}

// bn256Add implements a native elliptic curve point addition conforming to
// Istanbul consensus rules.
type bn256AddIstanbul struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256AddIstanbul) RequiredGas(input []byte) uint64 {
	return params.Bn256AddGasIstanbul
}

func (c *bn256AddIstanbul) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	return runBn256Add(input)
}

// bn256AddByzantium implements a native elliptic curve point addition
// conforming to Byzantium consensus rules.
type bn256AddByzantium struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256AddByzantium) RequiredGas(input []byte) uint64 {
	return params.Bn256AddGasByzantium
}

func (c *bn256AddByzantium) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	return runBn256Add(input)
}

// runBn256ScalarMul implements the Bn256ScalarMul precompile, referenced by
// both Byzantium and Istanbul operations.
func runBn256ScalarMul(input []byte) ([]byte, error) {
	p, err := newCurvePoint(getData(input, 0, 64))
	if err != nil {
		return nil, err
	}
	res := new(bn256.G1)
	res.ScalarMult(p, new(big.Int).SetBytes(getData(input, 64, 32)))
	return res.Marshal(), nil
}

// bn256ScalarMulIstanbul implements a native elliptic curve scalar
// multiplication conforming to Istanbul consensus rules.
type bn256ScalarMulIstanbul struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256ScalarMulIstanbul) RequiredGas(input []byte) uint64 {
	return params.Bn256ScalarMulGasIstanbul
}

func (c *bn256ScalarMulIstanbul) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	return runBn256ScalarMul(input)
}

// bn256ScalarMulByzantium implements a native elliptic curve scalar
// multiplication conforming to Byzantium consensus rules.
type bn256ScalarMulByzantium struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256ScalarMulByzantium) RequiredGas(input []byte) uint64 {
	return params.Bn256ScalarMulGasByzantium
}

func (c *bn256ScalarMulByzantium) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	return runBn256ScalarMul(input)
}

var (
	// true32Byte is returned if the bn256 pairing check succeeds.
	true32Byte = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}

	// false32Byte is returned if the bn256 pairing check fails.
	false32Byte = make([]byte, 32)

	// errBadPairingInput is returned if the bn256 pairing input is invalid.
	errBadPairingInput = errors.New("bad elliptic curve pairing size")
)

// runBn256Pairing implements the Bn256Pairing precompile, referenced by both
// Byzantium and Istanbul operations.
func runBn256Pairing(input []byte) ([]byte, error) {
	// Handle some corner cases cheaply
	if len(input)%192 > 0 {
		return nil, errBadPairingInput
	}
	// Convert the input into a set of coordinates
	var (
		cs []*bn256.G1
		ts []*bn256.G2
	)
	for i := 0; i < len(input); i += 192 {
		c, err := newCurvePoint(input[i : i+64])
		if err != nil {
			return nil, err
		}
		t, err := newTwistPoint(input[i+64 : i+192])
		if err != nil {
			return nil, err
		}
		cs = append(cs, c)
		ts = append(ts, t)
	}
	// Execute the pairing checks and return the results
	if bn256.PairingCheck(cs, ts) {
		return true32Byte, nil
	}
	return false32Byte, nil
}

// bn256PairingIstanbul implements a pairing pre-compile for the bn256 curve
// conforming to Istanbul consensus rules.
type bn256PairingIstanbul struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256PairingIstanbul) RequiredGas(input []byte) uint64 {
	return params.Bn256PairingBaseGasIstanbul + uint64(len(input)/192)*params.Bn256PairingPerPointGasIstanbul
}

func (c *bn256PairingIstanbul) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	return runBn256Pairing(input)
}

// bn256PairingByzantium implements a pairing pre-compile for the bn256 curve
// conforming to Byzantium consensus rules.
type bn256PairingByzantium struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bn256PairingByzantium) RequiredGas(input []byte) uint64 {
	return params.Bn256PairingBaseGasByzantium + uint64(len(input)/192)*params.Bn256PairingPerPointGasByzantium
}

func (c *bn256PairingByzantium) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	return runBn256Pairing(input)
}

type blake2F struct{}

func (c *blake2F) RequiredGas(input []byte) uint64 {
	// If the input is malformed, we can't calculate the gas, return 0 and let the
	// actual call choke and fault.
	if len(input) != blake2FInputLength {
		return 0
	}
	return uint64(binary.BigEndian.Uint32(input[0:4]))
}

const (
	blake2FInputLength        = 213
	blake2FFinalBlockBytes    = byte(1)
	blake2FNonFinalBlockBytes = byte(0)
)

var (
	errBlake2FInvalidInputLength = errors.New("invalid input length")
	errBlake2FInvalidFinalFlag   = errors.New("invalid final flag")
)

func (c *blake2F) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Make sure the input is valid (correct length and final flag)
	if len(input) != blake2FInputLength {
		return nil, errBlake2FInvalidInputLength
	}
	if input[212] != blake2FNonFinalBlockBytes && input[212] != blake2FFinalBlockBytes {
		return nil, errBlake2FInvalidFinalFlag
	}
	// Parse the input into the Blake2b call parameters
	var (
		rounds = binary.BigEndian.Uint32(input[0:4])
		final  = (input[212] == blake2FFinalBlockBytes)

		h [8]uint64
		m [16]uint64
		t [2]uint64
	)
	for i := 0; i < 8; i++ {
		offset := 4 + i*8
		h[i] = binary.LittleEndian.Uint64(input[offset : offset+8])
	}
	for i := 0; i < 16; i++ {
		offset := 68 + i*8
		m[i] = binary.LittleEndian.Uint64(input[offset : offset+8])
	}
	t[0] = binary.LittleEndian.Uint64(input[196:204])
	t[1] = binary.LittleEndian.Uint64(input[204:212])

	// Execute the compression function, extract and return the result
	blake2b.F(&h, m, t, final, rounds)

	output := make([]byte, 64)
	for i := 0; i < 8; i++ {
		offset := i * 8
		binary.LittleEndian.PutUint64(output[offset:offset+8], h[i])
	}
	return output, nil
}

var (
	errBLS12381InvalidInputLength          = errors.New("invalid input length")
	errBLS12381InvalidFieldElementTopBytes = errors.New("invalid field element top bytes")
	errBLS12381G1PointSubgroup             = errors.New("g1 point is not on correct subgroup")
	errBLS12381G2PointSubgroup             = errors.New("g2 point is not on correct subgroup")
)

// bls12381G1Add implements EIP-2537 G1Add precompile.
type bls12381G1Add struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381G1Add) RequiredGas(input []byte) uint64 {
	return params.Bls12381G1AddGas
}

func (c *bls12381G1Add) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 G1Add precompile.
	// > G1 addition call expects `256` bytes as an input that is interpreted as byte concatenation of two G1 points (`128` bytes each).
	// > Output is an encoding of addition operation result - single G1 point (`128` bytes).
	if len(input) != 256 {
		return nil, errBLS12381InvalidInputLength
	}
	var err error
	var p0, p1 *bls12381.PointG1

	// Initialize G1
	g := bls12381.NewG1()

	// Decode G1 point p_0
	if p0, err = g.DecodePoint(input[:128]); err != nil {
		return nil, err
	}
	// Decode G1 point p_1
	if p1, err = g.DecodePoint(input[128:]); err != nil {
		return nil, err
	}

	// Compute r = p_0 + p_1
	r := g.New()
	g.Add(r, p0, p1)

	// Encode the G1 point result into 128 bytes
	return g.EncodePoint(r), nil
}

// bls12381G1Mul implements EIP-2537 G1Mul precompile.
type bls12381G1Mul struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381G1Mul) RequiredGas(input []byte) uint64 {
	return params.Bls12381G1MulGas
}

func (c *bls12381G1Mul) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 G1Mul precompile.
	// > G1 multiplication call expects `160` bytes as an input that is interpreted as byte concatenation of encoding of G1 point (`128` bytes) and encoding of a scalar value (`32` bytes).
	// > Output is an encoding of multiplication operation result - single G1 point (`128` bytes).
	if len(input) != 160 {
		return nil, errBLS12381InvalidInputLength
	}
	var err error
	var p0 *bls12381.PointG1

	// Initialize G1
	g := bls12381.NewG1()

	// Decode G1 point
	if p0, err = g.DecodePoint(input[:128]); err != nil {
		return nil, err
	}
	// Decode scalar value
	e := new(big.Int).SetBytes(input[128:])

	// Compute r = e * p_0
	r := g.New()
	g.MulScalar(r, p0, e)

	// Encode the G1 point into 128 bytes
	return g.EncodePoint(r), nil
}

// bls12381G1MultiExp implements EIP-2537 G1MultiExp precompile.
type bls12381G1MultiExp struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381G1MultiExp) RequiredGas(input []byte) uint64 {
	// Calculate G1 point, scalar value pair length
	k := len(input) / 160
	if k == 0 {
		// Return 0 gas for small input length
		return 0
	}
	// Lookup discount value for G1 point, scalar value pair length
	var discount uint64
	if dLen := len(params.Bls12381MultiExpDiscountTable); k < dLen {
		discount = params.Bls12381MultiExpDiscountTable[k-1]
	} else {
		discount = params.Bls12381MultiExpDiscountTable[dLen-1]
	}
	// Calculate gas and return the result
	return (uint64(k) * params.Bls12381G1MulGas * discount) / 1000
}

func (c *bls12381G1MultiExp) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 G1MultiExp precompile.
	// G1 multiplication call expects `160*k` bytes as an input that is interpreted as byte concatenation of `k` slices each of them being a byte concatenation of encoding of G1 point (`128` bytes) and encoding of a scalar value (`32` bytes).
	// Output is an encoding of multiexponentiation operation result - single G1 point (`128` bytes).
	k := len(input) / 160
	if len(input) == 0 || len(input)%160 != 0 {
		return nil, errBLS12381InvalidInputLength
	}
	var err error
	points := make([]*bls12381.PointG1, k)
	scalars := make([]*big.Int, k)

	// Initialize G1
	g := bls12381.NewG1()

	// Decode point scalar pairs
	for i := 0; i < k; i++ {
		off := 160 * i
		t0, t1, t2 := off, off+128, off+160
		// Decode G1 point
		if points[i], err = g.DecodePoint(input[t0:t1]); err != nil {
			return nil, err
		}
		// Decode scalar value
		scalars[i] = new(big.Int).SetBytes(input[t1:t2])
	}

	// Compute r = e_0 * p_0 + e_1 * p_1 + ... + e_(k-1) * p_(k-1)
	r := g.New()
	g.MultiExp(r, points, scalars)

	// Encode the G1 point to 128 bytes
	return g.EncodePoint(r), nil
}

// bls12381G2Add implements EIP-2537 G2Add precompile.
type bls12381G2Add struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381G2Add) RequiredGas(input []byte) uint64 {
	return params.Bls12381G2AddGas
}

func (c *bls12381G2Add) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 G2Add precompile.
	// > G2 addition call expects `512` bytes as an input that is interpreted as byte concatenation of two G2 points (`256` bytes each).
	// > Output is an encoding of addition operation result - single G2 point (`256` bytes).
	if len(input) != 512 {
		return nil, errBLS12381InvalidInputLength
	}
	var err error
	var p0, p1 *bls12381.PointG2

	// Initialize G2
	g := bls12381.NewG2()
	r := g.New()

	// Decode G2 point p_0
	if p0, err = g.DecodePoint(input[:256]); err != nil {
		return nil, err
	}
	// Decode G2 point p_1
	if p1, err = g.DecodePoint(input[256:]); err != nil {
		return nil, err
	}

	// Compute r = p_0 + p_1
	g.Add(r, p0, p1)

	// Encode the G2 point into 256 bytes
	return g.EncodePoint(r), nil
}

// bls12381G2Mul implements EIP-2537 G2Mul precompile.
type bls12381G2Mul struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381G2Mul) RequiredGas(input []byte) uint64 {
	return params.Bls12381G2MulGas
}

func (c *bls12381G2Mul) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 G2MUL precompile logic.
	// > G2 multiplication call expects `288` bytes as an input that is interpreted as byte concatenation of encoding of G2 point (`256` bytes) and encoding of a scalar value (`32` bytes).
	// > Output is an encoding of multiplication operation result - single G2 point (`256` bytes).
	if len(input) != 288 {
		return nil, errBLS12381InvalidInputLength
	}
	var err error
	var p0 *bls12381.PointG2

	// Initialize G2
	g := bls12381.NewG2()

	// Decode G2 point
	if p0, err = g.DecodePoint(input[:256]); err != nil {
		return nil, err
	}
	// Decode scalar value
	e := new(big.Int).SetBytes(input[256:])

	// Compute r = e * p_0
	r := g.New()
	g.MulScalar(r, p0, e)

	// Encode the G2 point into 256 bytes
	return g.EncodePoint(r), nil
}

// bls12381G2MultiExp implements EIP-2537 G2MultiExp precompile.
type bls12381G2MultiExp struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381G2MultiExp) RequiredGas(input []byte) uint64 {
	// Calculate G2 point, scalar value pair length
	k := len(input) / 288
	if k == 0 {
		// Return 0 gas for small input length
		return 0
	}
	// Lookup discount value for G2 point, scalar value pair length
	var discount uint64
	if dLen := len(params.Bls12381MultiExpDiscountTable); k < dLen {
		discount = params.Bls12381MultiExpDiscountTable[k-1]
	} else {
		discount = params.Bls12381MultiExpDiscountTable[dLen-1]
	}
	// Calculate gas and return the result
	return (uint64(k) * params.Bls12381G2MulGas * discount) / 1000
}

func (c *bls12381G2MultiExp) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 G2MultiExp precompile logic
	// > G2 multiplication call expects `288*k` bytes as an input that is interpreted as byte concatenation of `k` slices each of them being a byte concatenation of encoding of G2 point (`256` bytes) and encoding of a scalar value (`32` bytes).
	// > Output is an encoding of multiexponentiation operation result - single G2 point (`256` bytes).
	k := len(input) / 288
	if len(input) == 0 || len(input)%288 != 0 {
		return nil, errBLS12381InvalidInputLength
	}
	var err error
	points := make([]*bls12381.PointG2, k)
	scalars := make([]*big.Int, k)

	// Initialize G2
	g := bls12381.NewG2()

	// Decode point scalar pairs
	for i := 0; i < k; i++ {
		off := 288 * i
		t0, t1, t2 := off, off+256, off+288
		// Decode G1 point
		if points[i], err = g.DecodePoint(input[t0:t1]); err != nil {
			return nil, err
		}
		// Decode scalar value
		scalars[i] = new(big.Int).SetBytes(input[t1:t2])
	}

	// Compute r = e_0 * p_0 + e_1 * p_1 + ... + e_(k-1) * p_(k-1)
	r := g.New()
	g.MultiExp(r, points, scalars)

	// Encode the G2 point to 256 bytes.
	return g.EncodePoint(r), nil
}

// bls12381Pairing implements EIP-2537 Pairing precompile.
type bls12381Pairing struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381Pairing) RequiredGas(input []byte) uint64 {
	return params.Bls12381PairingBaseGas + uint64(len(input)/384)*params.Bls12381PairingPerPairGas
}

func (c *bls12381Pairing) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 Pairing precompile logic.
	// > Pairing call expects `384*k` bytes as an inputs that is interpreted as byte concatenation of `k` slices. Each slice has the following structure:
	// > - `128` bytes of G1 point encoding
	// > - `256` bytes of G2 point encoding
	// > Output is a `32` bytes where last single byte is `0x01` if pairing result is equal to multiplicative identity in a pairing target field and `0x00` otherwise
	// > (which is equivalent of Big Endian encoding of Solidity values `uint256(1)` and `uin256(0)` respectively).
	k := len(input) / 384
	if len(input) == 0 || len(input)%384 != 0 {
		return nil, errBLS12381InvalidInputLength
	}

	// Initialize BLS12-381 pairing engine
	e := bls12381.NewPairingEngine()
	g1, g2 := e.G1, e.G2

	// Decode pairs
	for i := 0; i < k; i++ {
		off := 384 * i
		t0, t1, t2 := off, off+128, off+384

		// Decode G1 point
		p1, err := g1.DecodePoint(input[t0:t1])
		if err != nil {
			return nil, err
		}
		// Decode G2 point
		p2, err := g2.DecodePoint(input[t1:t2])
		if err != nil {
			return nil, err
		}

		// 'point is on curve' check already done,
		// Here we need to apply subgroup checks.
		if !g1.InCorrectSubgroup(p1) {
			return nil, errBLS12381G1PointSubgroup
		}
		if !g2.InCorrectSubgroup(p2) {
			return nil, errBLS12381G2PointSubgroup
		}

		// Update pairing engine with G1 and G2 ponits
		e.AddPair(p1, p2)
	}
	// Prepare 32 byte output
	out := make([]byte, 32)

	// Compute pairing and set the result
	if e.Check() {
		out[31] = 1
	}
	return out, nil
}

// decodeBLS12381FieldElement decodes BLS12-381 elliptic curve field element.
// Removes top 16 bytes of 64 byte input.
func decodeBLS12381FieldElement(in []byte) ([]byte, error) {
	if len(in) != 64 {
		return nil, errors.New("invalid field element length")
	}
	// check top bytes
	for i := 0; i < 16; i++ {
		if in[i] != byte(0x00) {
			return nil, errBLS12381InvalidFieldElementTopBytes
		}
	}
	out := make([]byte, 48)
	copy(out[:], in[16:])
	return out, nil
}

// bls12381MapG1 implements EIP-2537 MapG1 precompile.
type bls12381MapG1 struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381MapG1) RequiredGas(input []byte) uint64 {
	return params.Bls12381MapG1Gas
}

func (c *bls12381MapG1) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 Map_To_G1 precompile.
	// > Field-to-curve call expects `64` bytes an an input that is interpreted as a an element of the base field.
	// > Output of this call is `128` bytes and is G1 point following respective encoding rules.
	if len(input) != 64 {
		return nil, errBLS12381InvalidInputLength
	}

	// Decode input field element
	fe, err := decodeBLS12381FieldElement(input)
	if err != nil {
		return nil, err
	}

	// Initialize G1
	g := bls12381.NewG1()

	// Compute mapping
	r, err := g.MapToCurve(fe)
	if err != nil {
		return nil, err
	}

	// Encode the G1 point to 128 bytes
	return g.EncodePoint(r), nil
}

// bls12381MapG2 implements EIP-2537 MapG2 precompile.
type bls12381MapG2 struct{}

// RequiredGas returns the gas required to execute the pre-compiled contract.
func (c *bls12381MapG2) RequiredGas(input []byte) uint64 {
	return params.Bls12381MapG2Gas
}

func (c *bls12381MapG2) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	// Implements EIP-2537 Map_FP2_TO_G2 precompile logic.
	// > Field-to-curve call expects `128` bytes an an input that is interpreted as a an element of the quadratic extension field.
	// > Output of this call is `256` bytes and is G2 point following respective encoding rules.
	if len(input) != 128 {
		return nil, errBLS12381InvalidInputLength
	}

	// Decode input field element
	fe := make([]byte, 96)
	c0, err := decodeBLS12381FieldElement(input[:64])
	if err != nil {
		return nil, err
	}
	copy(fe[48:], c0)
	c1, err := decodeBLS12381FieldElement(input[64:])
	if err != nil {
		return nil, err
	}
	copy(fe[:48], c1)

	// Initialize G2
	g := bls12381.NewG2()

	// Compute mapping
	r, err := g.MapToCurve(fe)
	if err != nil {
		return nil, err
	}

	// Encode the G2 point to 256 bytes
	return g.EncodePoint(r), nil
}

// checkEnode implemented as a native contract.
type checkEnode struct{}

func (c checkEnode) RequiredGas(_ []byte) uint64 {
	return params.AutonityEnodeCheckGas
}
func (c checkEnode) Run(input []byte, blockNumber uint64, stateDB StateDB, _ common.Address) ([]byte, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("invalid enode - empty")
	}
	out := make([]byte, 64)
	nodeStr := string(input)

	node, err := enode.ParseV4NoResolve(nodeStr)
	if err != nil {
		copy(out[32:], true32Byte)
		return out, nil
	}

	address := crypto.PubkeyToAddress(*node.Pubkey())
	copy(out, address.Bytes())
	copy(out[32:], false32Byte)
	return out, nil
}
