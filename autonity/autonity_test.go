package autonity

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/rawdb"
	"github.com/autonity/autonity/core/state"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

func BenchmarkComputeCommittee(b *testing.B) {

	validatorCount := 100000
	validators, _, err := randomValidators(validatorCount, 30)
	require.NoError(b, err)
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 100

	// b.Run("computeCommittee", func(b *testing.B) {
	// 	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(b, err)
	// 	var header *types.Header
	// 	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	// 	require.NoError(b, err)
	// 	packedArgs, err := contractAbi.Pack("computeCommittee")
	// 	require.NoError(b, err)
	// 	_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
	// 	require.NoError(b, err)
	// 	benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	// })

	// b.Run("computeCommitteeOptimzed", func(b *testing.B) {
	// 	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(b, err)
	// 	var header *types.Header
	// 	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	// 	require.NoError(b, err)
	// 	packedArgs, err := contractAbi.Pack("computeCommitteeOptimzed")
	// 	require.NoError(b, err)
	// 	_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
	// 	require.NoError(b, err)
	// 	benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	// })

	// b.Run("computeCommitteePrecompiledSorting", func(b *testing.B) {
	// 	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(b, err)
	// 	var header *types.Header
	// 	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	// 	require.NoError(b, err)
	// 	packedArgs, err := contractAbi.Pack("computeCommitteePrecompiledSorting")
	// 	require.NoError(b, err)
	// 	_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
	// 	require.NoError(b, err)
	// 	benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	// })

	// b.Run("computeCommitteePrecompiledSortingFast", func(b *testing.B) {
	// 	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(b, err)
	// 	var header *types.Header
	// 	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	// 	require.NoError(b, err)
	// 	packedArgs, err := contractAbi.Pack("computeCommitteePrecompiledSortingFast")
	// 	require.NoError(b, err)
	// 	_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
	// 	require.NoError(b, err)
	// 	benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	// })

	// b.Run("computeCommitteePrecompiledSortingIterate", func(b *testing.B) {
	// 	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(b, err)
	// 	var header *types.Header
	// 	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	// 	require.NoError(b, err)
	// 	packedArgs, err := contractAbi.Pack("computeCommitteePrecompiledSortingIterate")
	// 	require.NoError(b, err)
	// 	_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
	// 	require.NoError(b, err)
	// 	benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	// })

	// b.Run("computeCommitteePrecompiledSortingIterateFast", func(b *testing.B) {
	// 	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(b, err)
	// 	var header *types.Header
	// 	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	// 	require.NoError(b, err)
	// 	packedArgs, err := contractAbi.Pack("computeCommitteePrecompiledSortingIterateFast")
	// 	require.NoError(b, err)
	// 	_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
	// 	require.NoError(b, err)
	// 	benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	// })

	b.Run("computeCommittee_ReadPrecompiled_WriteSolidity", func(b *testing.B) {
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(b, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(b, err)
		packedArgs, err := contractAbi.Pack("computeCommittee_ReadPrecompiled_WriteSolidity")
		require.NoError(b, err)
		_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
		require.NoError(b, err)
		benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	})

	b.Run("computeCommittee_FullPrecompiled", func(b *testing.B) {
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(b, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(b, err)
		packedArgs, err := contractAbi.Pack("computeCommittee_FullPrecompiled")
		require.NoError(b, err)
		_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
		require.NoError(b, err)
		benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	})

	b.Run("computeCommittee_FullPrecompiled_Return", func(b *testing.B) {
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(b, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(b, err)
		packedArgs, err := contractAbi.Pack("computeCommittee_FullPrecompiled_Return")
		require.NoError(b, err)
		_, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
		require.NoError(b, err)
		benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	})
}

func BenchmarkSortingAlgo(b *testing.B) {
	validatorCount := 1000
	validators, _, err := randomValidators(validatorCount, 70)
	require.NoError(b, err)
	input, err := inputToSort(validators)
	require.NoError(b, err)

	stateDB, _, _, err := initalizeEvm(nil)
	require.NoError(b, err)

	contracts := make([]vm.PrecompiledContract, 0)
	contractNames := make([]string, 0)
	contracts = append(contracts, &vm.SortLibrarySliceStable{})
	contractNames = append(contractNames, "SortLibrarySliceStable")
	contracts = append(contracts, &vm.SortLibrarySort{})
	contractNames = append(contractNames, "SortLibrarySort")
	contracts = append(contracts, vm.NewQuickSortIterateFast())
	contractNames = append(contractNames, "QuickSortIterateFast")
	contracts = append(contracts, vm.NewQuickSortIterate())
	contractNames = append(contractNames, "QuickSortIterate")
	contracts = append(contracts, &vm.QuickSortFast{})
	contractNames = append(contractNames, "QuickSortFast")
	contracts = append(contracts, &vm.QuickSort{})
	contractNames = append(contractNames, "QuickSort")

	for i, contract := range contracts {
		b.Run(contractNames[i], func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, err = contract.Run(input, 0, stateDB, common.Address{})
				require.NoError(b, err)
			}
		})
	}
}

func TestComputeCommittee(t *testing.T) {
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}

	t.Run("computeCommittee_ReadPrecompiled_WriteSolidity validators < committee", func(t *testing.T) {
		committeeSize := 100
		validatorCount := 10
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(t, err)
		res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee_ReadPrecompiled_WriteSolidity")
		// t.Log("error")
		// t.Log(err)
		// t.Log("res")
		// t.Log(res)
		// t.Log(string(res))
		require.NoError(t, err)
		voters := make([]common.Address, validatorCount)
		err = contractAbi.UnpackIntoInterface(&voters, "computeCommittee_ReadPrecompiled_WriteSolidity", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommittee")
		require.NoError(t, err)
		members := make([]types.CommitteeMember, validatorCount)
		err = contractAbi.UnpackIntoInterface(&members, "getCommittee", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getEpochTotalBondedStake")
		require.NoError(t, err)
		totalStake := big.NewInt(0)
		err = contractAbi.UnpackIntoInterface(&totalStake, "getEpochTotalBondedStake", res)
		require.NoError(t, err)
		err = isVotersSorted(voters, members, validators, totalStake)
		require.NoError(t, err)
	})

	t.Run("computeCommittee_ReadPrecompiled_WriteSolidity validators == committee", func(t *testing.T) {
		committeeSize := 100
		validatorCount := committeeSize
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(t, err)
		res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee_ReadPrecompiled_WriteSolidity")
		// t.Log("error")
		// t.Log(err)
		// t.Log("res")
		// t.Log(res)
		// t.Log(string(res))
		require.NoError(t, err)
		voters := make([]common.Address, committeeSize)
		err = contractAbi.UnpackIntoInterface(&voters, "computeCommittee_ReadPrecompiled_WriteSolidity", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommittee")
		require.NoError(t, err)
		members := make([]types.CommitteeMember, committeeSize)
		err = contractAbi.UnpackIntoInterface(&members, "getCommittee", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getEpochTotalBondedStake")
		require.NoError(t, err)
		totalStake := big.NewInt(0)
		err = contractAbi.UnpackIntoInterface(&totalStake, "getEpochTotalBondedStake", res)
		require.NoError(t, err)
		err = isVotersSorted(voters, members, validators, totalStake)
		require.NoError(t, err)
	})

	t.Run("computeCommittee_ReadPrecompiled_WriteSolidity validators > committee", func(t *testing.T) {
		committeeSize := 100
		validatorCount := 1000
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(t, err)
		res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee_ReadPrecompiled_WriteSolidity")
		// t.Log("error")
		// t.Log(err)
		// t.Log("res")
		// t.Log(res)
		// t.Log(string(res))
		require.NoError(t, err)
		voters := make([]common.Address, committeeSize)
		err = contractAbi.UnpackIntoInterface(&voters, "computeCommittee_ReadPrecompiled_WriteSolidity", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommittee")
		require.NoError(t, err)
		members := make([]types.CommitteeMember, committeeSize)
		err = contractAbi.UnpackIntoInterface(&members, "getCommittee", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getEpochTotalBondedStake")
		require.NoError(t, err)
		totalStake := big.NewInt(0)
		err = contractAbi.UnpackIntoInterface(&totalStake, "getEpochTotalBondedStake", res)
		require.NoError(t, err)
		err = isVotersSorted(voters, members, validators, totalStake)
		require.NoError(t, err)

	})

	t.Run("computeCommittee_FullPrecompiled", func(t *testing.T) {
		committeeSize := 100
		validatorCount := 1000
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(t, err)
		res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee_FullPrecompiled")
		// t.Log("error")
		// t.Log(err)
		// t.Log("res")
		// t.Log(res)
		// t.Log(string(res))
		require.NoError(t, err)
		voters := make([]common.Address, committeeSize)
		err = contractAbi.UnpackIntoInterface(&voters, "computeCommittee_FullPrecompiled", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommittee")
		require.NoError(t, err)
		members := make([]types.CommitteeMember, committeeSize)
		err = contractAbi.UnpackIntoInterface(&members, "getCommittee", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getEpochTotalBondedStake")
		require.NoError(t, err)
		totalStake := big.NewInt(0)
		err = contractAbi.UnpackIntoInterface(&totalStake, "getEpochTotalBondedStake", res)
		require.NoError(t, err)
		err = isVotersSorted(voters, members, validators, totalStake)
		require.NoError(t, err)
	})

	t.Run("computeCommittee_FullPrecompiled_Return", func(t *testing.T) {
		committeeSize := 100
		validatorCount := 1000
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(t, err)
		res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee_FullPrecompiled_Return")
		// t.Log("error")
		// t.Log(err)
		// t.Log("res")
		// t.Log(res)
		// t.Log(string(res))
		require.NoError(t, err)
		voters := make([]common.Address, committeeSize)
		err = contractAbi.UnpackIntoInterface(&voters, "computeCommittee_FullPrecompiled_Return", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommittee")
		require.NoError(t, err)
		members := make([]types.CommitteeMember, committeeSize)
		err = contractAbi.UnpackIntoInterface(&members, "getCommittee", res)
		require.NoError(t, err)
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getEpochTotalBondedStake")
		require.NoError(t, err)
		totalStake := big.NewInt(0)
		err = contractAbi.UnpackIntoInterface(&totalStake, "getEpochTotalBondedStake", res)
		require.NoError(t, err)
		err = isVotersSorted(voters, members, validators, totalStake)
		require.NoError(t, err)
	})
}

func BenchmarkComputeCommitteePrecompiled(b *testing.B) {
	committeeSize := 100
	validatorCount := 100000
	validators, _, err := randomValidators(validatorCount, 30)
	require.NoError(b, err)
	deployer := common.Address{}
	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	require.NoError(b, err)
	var header *types.Header
	contractAbi := &generated.AutonityTestAbi
	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	require.NoError(b, err)
	res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getValidatorListSlot")
	require.NoError(b, err)
	validatorListSlot := big.NewInt(0)
	err = contractAbi.UnpackIntoInterface(&validatorListSlot, "getValidatorListSlot", res)
	require.NoError(b, err)
	// b.Log(validatorListSlot)

	res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getValidatorsSlot")
	require.NoError(b, err)
	validatorsSlot := big.NewInt(0)
	err = contractAbi.UnpackIntoInterface(&validatorsSlot, "getValidatorsSlot", res)
	require.NoError(b, err)
	// b.Log(validatorsSlot)
	input := make([]byte, 64)
	validatorListSlot.FillBytes(input[0:32])
	validatorsSlot.FillBytes(input[32:64])
	precompiledContract := &vm.ComputeCommitteeReadOnly{}
	// precompiledContractTest := &vm.ComputeCommitteeTest{}
	// evm := evmContract.evmProvider(header, deployer, stateDB)

	b.Run("precompiledContract", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err = precompiledContract.Run(input, 0, stateDB, contractAddress)
			require.NoError(b, err)
		}
	})

	// b.Run("precompiledContractTest", func(b *testing.B) {
	// 	for i := 0; i < b.N; i++ {
	// 		_, err = precompiledContractTest.Run(input, 0, evm, contractAddress)
	// 		require.NoError(b, err)
	// 	}
	// })
}

func TestComputeCommitteePrecompiled(t *testing.T) {
	committeeSize := 100
	validatorCount := 1000
	validators, _, err := randomValidators(validatorCount, 30)
	require.NoError(t, err)

	// deploy contrac
	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, common.Address{})
	require.NoError(t, err)
	var header *types.Header
	contractAbi := &generated.AutonityTestAbi
	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	require.NoError(t, err)
	res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getValidatorListSlot")
	require.NoError(t, err)
	validatorListSlot := big.NewInt(0)
	err = contractAbi.UnpackIntoInterface(&validatorListSlot, "getValidatorListSlot", res)
	require.NoError(t, err)
	t.Log(validatorListSlot)

	res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getValidatorsSlot")
	require.NoError(t, err)
	validatorsSlot := big.NewInt(0)
	err = contractAbi.UnpackIntoInterface(&validatorsSlot, "getValidatorsSlot", res)
	require.NoError(t, err)
	t.Log(validatorsSlot)

	precompiledContract := &vm.ComputeCommitteeReadOnly{}
	input := make([]byte, 64)
	validatorListSlot.FillBytes(input[0:32])
	validatorsSlot.FillBytes(input[32:64])
	t.Log(input)
	output, err := precompiledContract.Run(input, 0, stateDB, contractAddress)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), big.NewInt(0).SetBytes(output[0:32]))
	t.Log(output[32:64])
	require.Equal(t, big.NewInt(int64(validatorCount)), big.NewInt(0).SetBytes(output[32:64]))
	require.Equal(t, validatorCount*32+64, len(output))
	positions := make(map[common.Address]int)
	for i, validator := range validators {
		positions[*validator.NodeAddress] = i
	}
	var lastStake *big.Int
	zeroAddress := common.BytesToAddress(make([]byte, 0))
	for i := 64; i < len(output); i += 32 {
		address := common.BytesToAddress(output[i : i+32])
		require.False(t, bytes.Equal(zeroAddress.Bytes(), address.Bytes()))
		idx, ok := positions[address]
		require.True(t, ok)
		if i > 64 {
			require.True(t, lastStake.Cmp(validators[idx].BondedStake) >= 0)
		}
		lastStake = validators[idx].BondedStake
	}
}

func TestCommitteeRead(t *testing.T) {
	committeeSize := 10
	validatorCount := 10
	validators, _, err := randomValidators(validatorCount, 30)
	require.NoError(t, err)

	// for _, validator := range validators {
	// 	t.Log(validator.NodeAddress)
	// }

	// deploy contrac
	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, common.Address{})
	require.NoError(t, err)
	var header *types.Header
	contractAbi := &generated.AutonityTestAbi
	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	require.NoError(t, err)

	res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee_ReadPrecompiled_WriteSolidity")
	// _, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee")
	t.Log("error")
	t.Log(err)
	t.Log("res")
	t.Log(res)
	t.Log(string(res))
	require.NoError(t, err)

	res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommitteeSlot")
	require.NoError(t, err)
	committeeSlot := big.NewInt(0)
	err = contractAbi.UnpackIntoInterface(&committeeSlot, "getCommitteeSlot", res)
	require.NoError(t, err)
	t.Log(committeeSlot)

	hash := make([]byte, 32)
	committeeSlot.FillBytes(hash)
	t.Log(hash)
	calculatedSlot := crypto.Keccak256Hash(hash).Big()

	for i := 0; i < committeeSize; i++ {
		idx := big.NewInt(int64(i))
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommitteeMemberSlot", idx)
		require.NoError(t, err)
		membmerSlot := big.NewInt(0)
		err = contractAbi.UnpackIntoInterface(&membmerSlot, "getCommitteeMemberSlot", res)
		require.NoError(t, err)
		t.Log(membmerSlot)
		// t.Log(calculatedSlot)
		require.Equal(t, calculatedSlot, membmerSlot)
		calculatedSlot.Add(calculatedSlot, big.NewInt(2))

		// res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommitteeMember", idx)
		// require.NoError(t, err)
		// t.Log(res)
		// membmer := types.CommitteeMember{Address: common.BytesToAddress(make([]byte, 0)), VotingPower: big.NewInt(0)}
		// err = contractAbi.UnpackIntoInterface(&membmer, "getCommitteeMember", res)
		// require.NoError(t, err)
		// t.Log(membmer)
	}

	res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getValidatorsSlot")
	require.NoError(t, err)
	validatorsSlot := big.NewInt(0)
	err = contractAbi.UnpackIntoInterface(&validatorsSlot, "getValidatorsSlot", res)
	require.NoError(t, err)
	t.Log(validatorsSlot)

	tester := &vm.TestCommitteeRead{}
	input := make([]byte, 64)
	committeeSlot.FillBytes(input[0:32])
	validatorsSlot.FillBytes(input[32:64])
	output, err := tester.Run(input, 0, stateDB, contractAddress)
	require.NoError(t, err)
	require.Equal(t, big.NewInt(1), big.NewInt(0).SetBytes(output))
}

func TestSorting(t *testing.T) {
	// Deploy contract for each test
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 100
	validatorCount := 1000

	t.Run("test sorting with 0% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 0)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "testSorting")
		require.NoError(t, err)
	})

	t.Run("test sorting with 30% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "testSorting")
		require.NoError(t, err)
	})

	t.Run("test sorting with 70% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 70)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "testSorting")
		require.NoError(t, err)
	})

	t.Run("test sorting with 100% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 100)
		require.NoError(t, err)
		stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "testSorting")
		require.NoError(t, err)
	})
}

func TestSortLibrarySliceTable(t *testing.T) {
	validatorCount := 1000
	stateDB, _, _, err := initalizeEvm(nil)
	require.NoError(t, err)

	contracts := make([]vm.PrecompiledContract, 0)
	contractNames := make([]string, 0)
	contracts = append(contracts, &vm.SortLibrarySliceStable{})
	contractNames = append(contractNames, "SortLibrarySliceStable")
	contracts = append(contracts, &vm.SortLibrarySort{})
	contractNames = append(contractNames, "SortLibrarySort")
	contracts = append(contracts, vm.NewQuickSortIterateFast())
	contractNames = append(contractNames, "QuickSortIterateFast")
	contracts = append(contracts, vm.NewQuickSortIterate())
	contractNames = append(contractNames, "QuickSortIterate")
	contracts = append(contracts, &vm.QuickSortFast{})
	contractNames = append(contractNames, "QuickSortFast")
	contracts = append(contracts, &vm.QuickSort{})
	contractNames = append(contractNames, "QuickSort")

	for i, contract := range contracts {
		t.Run(contractNames[i], func(t *testing.T) {
			validators, _, err := randomValidators(validatorCount, 0)
			require.NoError(t, err)
			input, err := inputToSort(validators)
			require.NoError(t, err)
			output, err := contract.Run(input, 0, stateDB, common.Address{})
			require.NoError(t, err)
			require.NoError(t, isOutputSorted(output, validators))
		})

		t.Run(contractNames[i], func(t *testing.T) {
			validators, _, err := randomValidators(validatorCount, 30)
			require.NoError(t, err)
			input, err := inputToSort(validators)
			require.NoError(t, err)
			output, err := contract.Run(input, 0, stateDB, common.Address{})
			require.NoError(t, err)
			require.NoError(t, isOutputSorted(output, validators))
		})

		t.Run(contractNames[i], func(t *testing.T) {
			validators, _, err := randomValidators(validatorCount, 70)
			require.NoError(t, err)
			input, err := inputToSort(validators)
			require.NoError(t, err)
			output, err := contract.Run(input, 0, stateDB, common.Address{})
			require.NoError(t, err)
			require.NoError(t, isOutputSorted(output, validators))
		})

		t.Run(contractNames[i], func(t *testing.T) {
			validators, _, err := randomValidators(validatorCount, 100)
			require.NoError(t, err)
			input, err := inputToSort(validators)
			require.NoError(t, err)
			output, err := contract.Run(input, 0, stateDB, common.Address{})
			require.NoError(t, err)
			require.NoError(t, isOutputSorted(output, validators))
		})
	}
}

func TestStruct(t *testing.T) {
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 100
	validatorCount := 100
	validators, _, err := randomValidators(validatorCount, 0)
	require.NoError(t, err)
	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	require.NoError(t, err)

	t.Run("test committee struct 1", func(t *testing.T) {
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "testCommitteeStruct", big.NewInt(1))
		require.NoError(t, err)
	})

	t.Run("test committee struct 2", func(t *testing.T) {
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "testCommitteeStruct", big.NewInt(2))
		require.NoError(t, err)
	})

}

func TestAssemblyProperArrray(t *testing.T) {
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 10
	validatorCount := 10
	validators, _, err := randomValidators(validatorCount, 0)
	require.NoError(t, err)
	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	require.NoError(t, err)

	t.Run("testAssemblyProperArrray", func(t *testing.T) {
		var header *types.Header
		argsPacked, err := contractAbi.Pack("testAssemblyProperArrray")
		require.NoError(t, err)
		res, _, err := evmContract.CallContractFunc(stateDB, header, contractAddress, argsPacked)
		// t.Log("error:")
		// t.Log(err)
		// t.Log("res")
		// t.Log(res)
		// t.Log(string(res))
		require.NoError(t, err)
		// var key common.Address
		// var mapLocation *big.Int
		// var location *big.Int
		// var calculatedLocation *big.Int
		unpacked, err := contractAbi.Unpack("testAssemblyProperArrray", res)
		require.NoError(t, err)
		t.Log(unpacked)
		t.Log(unpacked...)
		// baseOffset := crypto.Keccak256Hash(append(key, validatorsSlot...))
	})

}

func TestStructLocation(t *testing.T) {
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 10
	validatorCount := 10
	validators, _, err := randomValidators(validatorCount, 0)
	require.NoError(t, err)
	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	require.NoError(t, err)

	t.Run("testStructLocation", func(t *testing.T) {
		var header *types.Header
		argsPacked, err := contractAbi.Pack("testStructLocation", big.NewInt(1), big.NewInt(2), big.NewInt(3), big.NewInt(4), big.NewInt(5))
		require.NoError(t, err)
		res, _, err := evmContract.CallContractFunc(stateDB, header, contractAddress, argsPacked)
		t.Log("error:")
		t.Log(err)
		t.Log("res")
		t.Log(res)
		t.Log(string(res))
		unpacked, err := contractAbi.Unpack("testStructLocation", res)
		require.NoError(t, err)
		// t.Log(unpacked)
		t.Log(unpacked...)
		require.NoError(t, err)

		// get item
		argsPacked, err = contractAbi.Pack("getItem")
		require.NoError(t, err)
		res, _, err = evmContract.CallContractFunc(stateDB, header, contractAddress, argsPacked)
		t.Log("error:")
		t.Log(err)
		t.Log("res")
		t.Log(res)
		t.Log(string(res))
		unpacked, err = contractAbi.Unpack("getItem", res)
		require.NoError(t, err)
		// t.Log(unpacked)
		t.Log(unpacked...)
		require.NoError(t, err)
		// var key common.Address
		// var mapLocation *big.Int
		// var location *big.Int
		// var calculatedLocation *big.Int

		// baseOffset := crypto.Keccak256Hash(append(key, validatorsSlot...))
	})
}

func TestElectProposer(t *testing.T) {
	height := uint64(9999)
	samePowers := []int{100, 100, 100, 100}
	linearPowers := []int{100, 200, 400, 800}
	var ac = &AutonityContract{}
	t.Run("Proposer election should be deterministic", func(t *testing.T) {
		committee := generateCommittee(samePowers)
		parentHeader := newBlockHeader(height, committee)
		for h := uint64(0); h < uint64(100); h++ {
			for r := int64(0); r <= int64(3); r++ {
				proposer1 := ac.electProposer(parentHeader, h, r)
				proposer2 := ac.electProposer(parentHeader, h, r)
				require.Equal(t, proposer1, proposer2)
			}
		}
	})

	t.Run("Proposer selection, print and compare the scheduling rate with same stake", func(t *testing.T) {
		committee := generateCommittee(samePowers)
		parentHeader := newBlockHeader(height, committee)
		maxHeight := uint64(10000)
		maxRound := int64(4)
		//expectedRatioDelta := float64(0.01)
		counterMap := make(map[common.Address]int)
		counterMap[common.Address{}] = 1
		for h := uint64(0); h < maxHeight; h++ {
			for round := int64(0); round < maxRound; round++ {
				proposer := ac.electProposer(parentHeader, h, round)
				_, ok := counterMap[proposer]
				if ok {
					counterMap[proposer]++
				} else {
					counterMap[proposer] = 1
				}
			}
		}

		totalStake := 0
		for _, s := range samePowers {
			totalStake += s
		}

		for i, c := range committee {
			stake := samePowers[i]
			scheduled := counterMap[c.Address]
			log.Print("electing ", "proposer: ", c.Address.String(), " stake: ", stake, " scheduled: ", scheduled)
		}
	})

	t.Run("Proposer selection, print and compare the scheduling rate with liner increasing stake", func(t *testing.T) {
		committee := generateCommittee(linearPowers)
		parentHeader := newBlockHeader(height, committee)
		maxHeight := uint64(1000000)
		maxRound := int64(4)
		//expectedRatioDelta := float64(0.01)
		counterMap := make(map[common.Address]int)
		counterMap[common.Address{}] = 1
		for h := uint64(0); h < maxHeight; h++ {
			for round := int64(0); round < maxRound; round++ {
				proposer := ac.electProposer(parentHeader, h, round)
				_, ok := counterMap[proposer]
				if ok {
					counterMap[proposer]++
				} else {
					counterMap[proposer] = 1
				}
			}
		}

		totalStake := 0
		for _, s := range samePowers {
			totalStake += s
		}

		for _, c := range committee {
			stake := c.VotingPower.Uint64()
			scheduled := counterMap[c.Address]
			log.Print("electing ", "proposer: ", c.Address.String(), " stake: ", stake, " scheduled: ", scheduled)
		}
	})
}

func newBlockHeader(height uint64, committee types.Committee) *types.Header {
	// use random nonce to create different blocks
	var nonce types.BlockNonce
	for i := 0; i < len(nonce); i++ {
		nonce[0] = byte(rand.Intn(256)) //nolint
	}
	return &types.Header{
		Number:    new(big.Int).SetUint64(height),
		Nonce:     nonce,
		Committee: committee,
	}
}

func generateCommittee(powers []int) types.Committee {
	vals := make(types.Committee, 0)
	for _, p := range powers {
		privateKey, _ := crypto.GenerateKey()
		committeeMember := types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetInt64(int64(p)),
		}
		vals = append(vals, committeeMember)
	}
	sort.Sort(vals)
	return vals
}

func deployAutonityTest(
	committeeSize int, validators []params.Validator, deployer common.Address,
) (*state.StateDB, *EVMContract, common.Address, error) {
	abi := &generated.AutonityTestAbi
	stateDB, evm, evmContract, err := initalizeEvm(abi)
	if err != nil {
		return stateDB, evmContract, common.Address{}, err
	}
	contractConfig := autonityTestConfig()
	contractConfig.Protocol.OperatorAccount = common.Address{}
	contractConfig.Protocol.CommitteeSize = big.NewInt(int64(committeeSize))
	args, err := abi.Pack("", validators, contractConfig)
	if err != nil {
		return stateDB, evmContract, common.Address{}, err
	}
	contractAddress, err := deployContract(generated.AutonityTestBytecode, args, deployer, evm)
	return stateDB, evmContract, contractAddress, err
}

func initalizeEvm(abi *abi.ABI) (*state.StateDB, *vm.EVM, *EVMContract, error) {
	ethDb := rawdb.NewMemoryDatabase()
	db := state.NewDatabase(ethDb)
	stateDB, err := state.New(common.Hash{}, db, nil)
	if err != nil {
		return new(state.StateDB), new(vm.EVM), new(EVMContract), err
	}
	evm := createTestVM(stateDB)
	evmContract := NewEVMContract(testEVMProvider(), abi, ethDb, params.TestChainConfig)
	return stateDB, evm, evmContract, nil
}

func deployContract(byteCode []byte, args []byte, deployer common.Address, evm *vm.EVM) (common.Address, error) {
	gas := uint64(math.MaxUint64)
	value := common.Big0
	data := append(byteCode, args...)
	_, contractAddress, _, err := evm.Create(vm.AccountRef(deployer), data, gas, value)
	return contractAddress, err
}

// Packs the args and then calls the function
// can also return result if needed
func callContractFunction(
	evmContract *EVMContract, contractAddress common.Address, stateDB *state.StateDB, header *types.Header, abi *abi.ABI,
	methodName string, args ...interface{},
) ([]byte, error) {
	argsPacked, err := abi.Pack(methodName, args...)
	if err != nil {
		return make([]byte, 0), err
	}
	res, _, err := evmContract.CallContractFunc(stateDB, header, contractAddress, argsPacked)
	return res, err
}

func randomValidators(count int, randomPercentage int) ([]params.Validator, []*ecdsa.PrivateKey, error) {
	if count == 0 {
		return []params.Validator{}, []*ecdsa.PrivateKey{}, nil
	}

	bondedStake := make([]int64, count)
	for i := 0; i < count; i++ {
		bondedStake[i] = int64(rand.Uint64() >> 1)
	}
	if randomPercentage < 100 {
		sort.SliceStable(bondedStake, func(i, j int) bool {
			return bondedStake[i] > bondedStake[j]
		})
		if count > 1 && bondedStake[0] < bondedStake[1] {
			return []params.Validator{}, []*ecdsa.PrivateKey{}, errors.New("Not sorted")
		}
	}

	validatorList := make([]params.Validator, count)
	privateKeyList := make([]*ecdsa.PrivateKey, count)
	for i := 0; i < count; i++ {
		var privateKey *ecdsa.PrivateKey
		var err error
		for {
			privateKey, err = crypto.GenerateKey()
			if err == nil {
				break
			}
		}
		privateKeyList[i] = privateKey
		publicKey := privateKey.PublicKey
		enode := "enode://" + string(crypto.PubECDSAToHex(&publicKey)[2:]) + "@3.209.45.79:30303"
		address := crypto.PubkeyToAddress(publicKey)
		validatorList[i] = params.Validator{
			Treasury:      address,
			Enode:         enode,
			BondedStake:   big.NewInt(bondedStake[i]),
			OracleAddress: address,
		}
		err = validatorList[i].Validate()
		if err != nil {
			return []params.Validator{}, []*ecdsa.PrivateKey{}, err
		}
	}

	if randomPercentage == 0 || randomPercentage == 100 {
		return validatorList, privateKeyList, nil
	}

	randomValidatorCount := count * randomPercentage / 100
	randomIndex := make(map[uint32]bool)
	randomIndex[0] = true
	for i := 0; i < randomValidatorCount; i++ {
		var idx uint32
		for {
			idx = rand.Uint32() % uint32(count)
			_, ok := randomIndex[idx]
			if !ok {
				break
			}
		}

		stake := validatorList[idx-1].BondedStake
		validatorList[idx].BondedStake = new(big.Int).Add(stake, big.NewInt(int64(rand.Uint64()>>1)))
		randomIndex[idx] = true
	}
	return validatorList, privateKeyList, nil
}

func autonityTestConfig() AutonityConfig {
	config := AutonityConfig{
		Policy: AutonityPolicy{
			TreasuryFee:     new(big.Int).SetUint64(params.TestAutonityContractConfig.TreasuryFee),
			MinBaseFee:      new(big.Int).SetUint64(params.TestAutonityContractConfig.MinBaseFee),
			DelegationRate:  new(big.Int).SetUint64(params.TestAutonityContractConfig.DelegationRate),
			UnbondingPeriod: new(big.Int).SetUint64(params.TestAutonityContractConfig.UnbondingPeriod),
			TreasuryAccount: params.TestAutonityContractConfig.Operator,
		},
		Contracts: AutonityContracts{
			AccountabilityContract: AccountabilityContractAddress,
			OracleContract:         OracleContractAddress,
			AcuContract:            ACUContractAddress,
			SupplyControlContract:  SupplyControlContractAddress,
			StabilizationContract:  StabilizationContractAddress,
		},
		Protocol: AutonityProtocol{
			OperatorAccount: params.TestAutonityContractConfig.Operator,
			EpochPeriod:     new(big.Int).SetUint64(params.TestAutonityContractConfig.EpochPeriod),
			BlockPeriod:     new(big.Int).SetUint64(params.TestAutonityContractConfig.BlockPeriod),
			CommitteeSize:   new(big.Int).SetUint64(params.TestAutonityContractConfig.MaxCommitteeSize),
		},
		ContractVersion: big.NewInt(1),
	}
	return config
}

func createTestVM(state vm.StateDB) *vm.EVM {
	vmBlockContext := vm.BlockContext{
		Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
		CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
		BlockNumber: common.Big0,
	}

	txContext := vm.TxContext{
		Origin:   common.Address{},
		GasPrice: common.Big0,
	}

	evm := vm.NewEVM(vmBlockContext, txContext, state, params.TestChainConfig, vm.Config{})
	return evm
}

func testEVMProvider() func(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
	return func(header *types.Header, origin common.Address, statedb *state.StateDB) *vm.EVM {
		vmBlockContext := vm.BlockContext{
			Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
			BlockNumber: common.Big0,
		}
		txContext := vm.TxContext{
			Origin:   common.Address{},
			GasPrice: common.Big0,
		}
		evm := vm.NewEVM(vmBlockContext, txContext, statedb, params.TestChainConfig, vm.Config{})
		return evm
	}
}

func benchmarkWithGas(
	b *testing.B, evmContract *EVMContract, stateDB *state.StateDB, header *types.Header,
	contractAddress common.Address, packedArgs []byte,
) {
	gas := uint64(math.MaxUint64)
	var gasUsed uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, gasLeft, err := evmContract.CallContractFunc(stateDB, header, contractAddress, packedArgs)
		require.NoError(b, err)
		gasUsed += gas - gasLeft
	}
	b.Log(1.0 * gasUsed / uint64(b.N))
}

func inputToSort(validators []params.Validator) ([]byte, error) {
	validatorCount := len(validators)
	input := make([]byte, validatorCount*64)
	for i := 0; i < validatorCount; i++ {
		copied := copy(input[i*64+12:i*64+32], validators[i].NodeAddress.Bytes())
		if copied != 20 {
			return input, fmt.Errorf("Could not copy address")
		}
		stake := validators[i].BondedStake.Bytes()
		if len(stake) > 32 {
			return input, fmt.Errorf("stake(big.int) size greater than 32")
		}
		copied = copy(input[i*64+64-len(stake):i*64+64], stake)
		if len(stake) != copied {
			return input, fmt.Errorf("Could not copy stake(big.int)")
		}
	}
	return input, nil
}

func isOutputSorted(output []byte, validators []params.Validator) error {
	if output[31] != 1 {
		return fmt.Errorf("unsuccessful call")
	}
	// only addresses in the output
	if len(validators)+1 != len(output)/32 {
		return fmt.Errorf("length mismatch")
	}
	position := make(map[common.Address]int)
	for i, validator := range validators {
		if _, ok := position[*validator.NodeAddress]; ok {
			return fmt.Errorf("duplicate validator")
		}
		position[*validator.NodeAddress] = i
	}
	lastStake := big.NewInt(0)
	for i := 32; i < len(output); i += 32 {
		address := common.BytesToAddress(output[i : i+32])
		idx, ok := position[address]
		if !ok {
			return fmt.Errorf("validator not found")
		}
		stake := validators[idx].BondedStake
		if i > 32 && lastStake.Cmp(stake) < 0 {
			return fmt.Errorf("not sorted")
		}
		lastStake = stake
	}
	return nil
}

func isVotersSorted(voters []common.Address, committeeMembers []types.CommitteeMember, validators []params.Validator, totalStake *big.Int) error {
	if len(voters) > len(validators) {
		return fmt.Errorf("More voters than validators")
	}
	if len(voters) != len(committeeMembers) {
		return fmt.Errorf("Committee size not equal to voter size")
	}
	positions := make(map[common.Address]int)
	for i, validator := range validators {
		if _, ok := positions[validator.OracleAddress]; ok {
			return fmt.Errorf("duplicate validator")
		}
		positions[validator.OracleAddress] = i
	}
	lastStake := big.NewInt(0)
	totalStakeCalculated := big.NewInt(0)
	for i, voter := range voters {
		idx, ok := positions[voter]
		if !ok {
			return fmt.Errorf("voter not found")
		}
		if i > 0 && lastStake.Cmp(validators[idx].BondedStake) < 0 {
			return fmt.Errorf("not sorted")
		}
		lastStake = validators[idx].BondedStake
		totalStakeCalculated = totalStakeCalculated.Add(totalStakeCalculated, lastStake)

		if !bytes.Equal(committeeMembers[i].Address.Bytes(), validators[idx].NodeAddress.Bytes()) {
			return fmt.Errorf("Committee member mismatch")
		}
	}
	if totalStakeCalculated.Cmp(totalStake) != 0 {
		return fmt.Errorf("epochTotalStake mismatch")
	}
	return nil
}
