package autonity

import (
	"crypto/ecdsa"
	"errors"
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
	// 	stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(b, err)
	// 	var header *types.Header
	// 	err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "applyStakingOperations")
	// 	require.NoError(b, err)
	// 	packedArgs, err := contractAbi.Pack("computeCommittee")
	// 	require.NoError(b, err)
	// 	_, _, err = evmContract.CallContractFunc(stateDb, header, contractAddress, packedArgs)
	// 	require.NoError(b, err)
	// 	benchmarkWithGas(b, evmContract, stateDb, header, contractAddress, packedArgs)
	// })

	// b.Run("computeCommitteeOptimzed", func(b *testing.B) {
	// 	stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(b, err)
	// 	var header *types.Header
	// 	err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "applyStakingOperations")
	// 	require.NoError(b, err)
	// 	packedArgs, err := contractAbi.Pack("computeCommitteeOptimzed")
	// 	require.NoError(b, err)
	// 	_, _, err = evmContract.CallContractFunc(stateDb, header, contractAddress, packedArgs)
	// 	require.NoError(b, err)
	// 	benchmarkWithGas(b, evmContract, stateDb, header, contractAddress, packedArgs)
	// })

	b.Run("computeCommitteePrecompiledSorting", func(b *testing.B) {
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(b, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "applyStakingOperations")
		require.NoError(b, err)
		packedArgs, err := contractAbi.Pack("computeCommitteePrecompiledSorting")
		require.NoError(b, err)
		_, _, err = evmContract.CallContractFunc(stateDb, header, contractAddress, packedArgs)
		require.NoError(b, err)
		benchmarkWithGas(b, evmContract, stateDb, header, contractAddress, packedArgs)
	})

	b.Run("computeCommitteePrecompiledSortingFast", func(b *testing.B) {
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(b, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "applyStakingOperations")
		require.NoError(b, err)
		packedArgs, err := contractAbi.Pack("computeCommitteePrecompiledSortingFast")
		require.NoError(b, err)
		_, _, err = evmContract.CallContractFunc(stateDb, header, contractAddress, packedArgs)
		require.NoError(b, err)
		benchmarkWithGas(b, evmContract, stateDb, header, contractAddress, packedArgs)
	})

	b.Run("computeCommitteePrecompiledSortingIterate", func(b *testing.B) {
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(b, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "applyStakingOperations")
		require.NoError(b, err)
		packedArgs, err := contractAbi.Pack("computeCommitteePrecompiledSortingIterate")
		require.NoError(b, err)
		_, _, err = evmContract.CallContractFunc(stateDb, header, contractAddress, packedArgs)
		require.NoError(b, err)
		benchmarkWithGas(b, evmContract, stateDb, header, contractAddress, packedArgs)
	})

	b.Run("computeCommitteePrecompiledSortingIterateFast", func(b *testing.B) {
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(b, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "applyStakingOperations")
		require.NoError(b, err)
		packedArgs, err := contractAbi.Pack("computeCommitteePrecompiledSortingIterateFast")
		require.NoError(b, err)
		_, _, err = evmContract.CallContractFunc(stateDb, header, contractAddress, packedArgs)
		require.NoError(b, err)
		benchmarkWithGas(b, evmContract, stateDb, header, contractAddress, packedArgs)
	})
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
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSorting")
		require.NoError(t, err)
	})

	t.Run("test sorting with 30% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSorting")
		require.NoError(t, err)
	})

	t.Run("test sorting with 70% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 70)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSorting")
		require.NoError(t, err)
	})

	t.Run("test sorting with 100% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 100)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSorting")
		require.NoError(t, err)
	})
}

func TestSortingPrecompiled(t *testing.T) {
	// Deploy contract for each test
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 100
	// in precompiled contract, we have to take array of fixed size for returnData
	// if validatorCount > committeeSize, the test will not work
	// because in the test we are expecting to get all the validators sorted and returned
	// in practice, we don't need to return all the validators, only top 100 validators
	validatorCount := committeeSize

	t.Run("test sorting with 0% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 0)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiled")
		require.NoError(t, err)
	})

	t.Run("test sorting with 30% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiled")
		require.NoError(t, err)
	})

	t.Run("test sorting with 70% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 70)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiled")
		require.NoError(t, err)
	})

	t.Run("test sorting with 100% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 100)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiled")
		require.NoError(t, err)
	})
}

func TestSortingPrecompiledFast(t *testing.T) {
	// Deploy contract for each test
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 100
	// in precompiled contract, we have to take array of fixed size for returnData
	// if validatorCount > committeeSize, the test will not work
	// because in the test we are expecting to get all the validators sorted and returned
	// in practice, we don't need to return all the validators, only top 100 validators
	validatorCount := committeeSize

	t.Run("test sorting with 0% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 0)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledFast")
		require.NoError(t, err)
	})

	t.Run("test sorting with 30% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledFast")
		require.NoError(t, err)
	})

	t.Run("test sorting with 70% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 70)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledFast")
		require.NoError(t, err)
	})

	t.Run("test sorting with 100% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 100)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledFast")
		require.NoError(t, err)
	})
}

func TestArraySlice(t *testing.T) {
	queue := make([]int, 0)
	printStuff(queue, t)
	queue = append(queue, 1)
	queue = append(queue, 2)
	printStuff(queue, t)
	queue = append(queue, 3)
	queue = append(queue, 4)
	printStuff(queue, t)
	queue = queue[1:]
	printStuff(queue, t)
	queue = queue[3:]
	printStuff(queue, t)
}

func printStuff(queue []int, t *testing.T) {
	t.Log(len(queue))
	t.Log(queue)
	for i := 0; i < len(queue); i++ {
		t.Log(&queue[i])
	}
}

func TestSort(t *testing.T) {
	// Deploy contract for each test
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 9
	// in precompiled contract, we have to take array of fixed size for returnData
	// if validatorCount > committeeSize, the test will not work
	// because in the test we are expecting to get all the validators sorted and returned
	// in practice, we don't need to return all the validators, only top 100 validators
	validatorCount := committeeSize

	// t.Run("test sorting with 0% randomness", func(t *testing.T) {
	// 	validators, _, err := randomValidators(validatorCount, 0)
	// 	require.NoError(t, err)
	// 	stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(t, err)
	// 	var header *types.Header
	// 	// err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
	// 	argsPacked, err := contractAbi.Pack("sort")
	// 	require.NoError(t, err)
	// 	res, _, err := evmContract.CallContractFunc(stateDb, header, contractAddress, argsPacked)
	// 	require.NoError(t, err)
	// 	var addresses []common.Address
	// 	err = contractAbi.UnpackIntoInterface(&addresses, "sort", res)
	// 	require.NoError(t, err)
	// 	t.Log("printing res")
	// 	t.Log(addresses)
	// })

	t.Run("test sorting with 30% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		sortedValidators := make([]params.Validator, validatorCount)
		for i := 0; i < validatorCount; i++ {
			sortedValidators[i] = validators[i]
		}
		sort.SliceStable(sortedValidators, func(i, j int) bool {
			return sortedValidators[i].BondedStake.Cmp(sortedValidators[j].BondedStake) == 1
		})
		for i := 0; i < validatorCount; i++ {
			t.Log(sortedValidators[i].NodeAddress)
			t.Log(sortedValidators[i].BondedStake)
		}
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		// err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
		argsPacked, err := contractAbi.Pack("sort")
		require.NoError(t, err)
		res, _, err := evmContract.CallContractFunc(stateDb, header, contractAddress, argsPacked)
		require.NoError(t, err)
		var addresses []common.Address
		err = contractAbi.UnpackIntoInterface(&addresses, "sort", res)
		require.NoError(t, err)
		t.Log("printing res")
		// t.Log(addresses)
		for i := 0; i < validatorCount; i++ {
			t.Log(addresses[i])
		}
	})

	// t.Run("test sorting with 70% randomness", func(t *testing.T) {
	// 	validators, _, err := randomValidators(validatorCount, 70)
	// 	require.NoError(t, err)
	// 	stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(t, err)
	// 	var header *types.Header
	// 	// err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
	// 	argsPacked, err := contractAbi.Pack("sort")
	// 	require.NoError(t, err)
	// 	res, _, err := evmContract.CallContractFunc(stateDb, header, contractAddress, argsPacked)
	// 	require.NoError(t, err)
	// 	var addresses []common.Address
	// 	err = contractAbi.UnpackIntoInterface(&addresses, "sort", res)
	// 	require.NoError(t, err)
	// 	t.Log("printing res")
	// 	t.Log(addresses)
	// })

	// t.Run("test sorting with 100% randomness", func(t *testing.T) {
	// 	validators, _, err := randomValidators(validatorCount, 100)
	// 	require.NoError(t, err)
	// 	stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	// 	require.NoError(t, err)
	// 	var header *types.Header
	// 	// err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
	// 	argsPacked, err := contractAbi.Pack("sort")
	// 	require.NoError(t, err)
	// 	res, _, err := evmContract.CallContractFunc(stateDb, header, contractAddress, argsPacked)
	// 	require.NoError(t, err)
	// 	var addresses []common.Address
	// 	err = contractAbi.UnpackIntoInterface(&addresses, "sort", res)
	// 	require.NoError(t, err)
	// 	t.Log("printing res")
	// 	t.Log(addresses)
	// })
}

func TestSortingPrecompiledIterate(t *testing.T) {
	// Deploy contract for each test
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 100
	// in precompiled contract, we have to take array of fixed size for returnData
	// if validatorCount > committeeSize, the test will not work
	// because in the test we are expecting to get all the validators sorted and returned
	// in practice, we don't need to return all the validators, only top 100 validators
	validatorCount := committeeSize

	t.Run("test sorting with 0% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 0)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
		require.NoError(t, err)
	})

	t.Run("test sorting with 30% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		// err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
		argsPacked, err := contractAbi.Pack("testSortingPrecompiledIterate")
		require.NoError(t, err)
		res, _, err := evmContract.CallContractFunc(stateDb, header, contractAddress, argsPacked)
		t.Log("printing error")
		t.Log(err)
		t.Log("printing res")
		t.Log(res)
		t.Log(string(res))
		require.NoError(t, err)
	})

	t.Run("test sorting with 70% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 70)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
		require.NoError(t, err)
	})

	t.Run("test sorting with 100% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 100)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
		require.NoError(t, err)
	})
}

func TestSortingPrecompiledIterateFast(t *testing.T) {
	// Deploy contract for each test
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 100
	// in precompiled contract, we have to take array of fixed size for returnData
	// if validatorCount > committeeSize, the test will not work
	// because in the test we are expecting to get all the validators sorted and returned
	// in practice, we don't need to return all the validators, only top 100 validators
	validatorCount := committeeSize

	t.Run("test sorting with 0% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 0)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterateFast")
		require.NoError(t, err)
	})

	t.Run("test sorting with 30% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 30)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		// err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterate")
		argsPacked, err := contractAbi.Pack("testSortingPrecompiledIterateFast")
		require.NoError(t, err)
		res, _, err := evmContract.CallContractFunc(stateDb, header, contractAddress, argsPacked)
		t.Log("printing error")
		t.Log(err)
		t.Log("printing res")
		t.Log(res)
		t.Log(string(res))
		require.NoError(t, err)
	})

	t.Run("test sorting with 70% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 70)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterateFast")
		require.NoError(t, err)
	})

	t.Run("test sorting with 100% randomness", func(t *testing.T) {
		validators, _, err := randomValidators(validatorCount, 100)
		require.NoError(t, err)
		stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
		require.NoError(t, err)
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testSortingPrecompiledIterateFast")
		require.NoError(t, err)
	})
}

func TestStruct(t *testing.T) {
	contractAbi := &generated.AutonityTestAbi
	deployer := common.Address{}
	committeeSize := 100
	validatorCount := 100
	validators, _, err := randomValidators(validatorCount, 0)
	require.NoError(t, err)
	stateDb, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	require.NoError(t, err)

	t.Run("test committee struct 1", func(t *testing.T) {
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testCommitteeStruct", big.NewInt(1))
		require.NoError(t, err)
	})

	t.Run("test committee struct 2", func(t *testing.T) {
		var header *types.Header
		err = callContractFunction(evmContract, contractAddress, stateDb, header, contractAbi, "testCommitteeStruct", big.NewInt(2))
		require.NoError(t, err)
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
	stateDb, evm, evmContract, err := initalizeEvm(abi)
	if err != nil {
		return stateDb, evmContract, common.Address{}, err
	}
	contractConfig := autonityTestConfig()
	contractConfig.Protocol.OperatorAccount = common.Address{}
	contractConfig.Protocol.CommitteeSize = big.NewInt(int64(committeeSize))
	args, err := abi.Pack("", validators, contractConfig)
	if err != nil {
		return stateDb, evmContract, common.Address{}, err
	}
	contractAddress, err := deployContract(generated.AutonityTestBytecode, args, deployer, evm)
	return stateDb, evmContract, contractAddress, err
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
	evmContract *EVMContract, contractAddress common.Address, stateDb *state.StateDB, header *types.Header, abi *abi.ABI,
	methodName string, args ...interface{},
) error {
	argsPacked, err := abi.Pack(methodName, args...)
	if err != nil {
		return err
	}
	_, _, err = evmContract.CallContractFunc(stateDb, header, contractAddress, argsPacked)
	return err
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
			Treasury:    address,
			Enode:       enode,
			BondedStake: big.NewInt(bondedStake[i]),
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
	b *testing.B, evmContract *EVMContract, stateDb *state.StateDB, header *types.Header,
	contractAddress common.Address, packedArgs []byte,
) {
	gas := uint64(math.MaxUint64)
	var gasUsed uint64

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, gasLeft, err := evmContract.CallContractFunc(stateDb, header, contractAddress, packedArgs)
		require.NoError(b, err)
		gasUsed += gas - gasLeft
	}
	b.Log(1.0 * gasUsed / uint64(b.N))
}
