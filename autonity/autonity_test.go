package autonity

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
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
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

func BenchmarkComputeCommittee(b *testing.B) {

	validatorCount := 100000
	validators, err := randomValidators(validatorCount, 30)
	require.NoError(b, err)
	contractAbi := &generated.AutonityAbi
	deployer := params.DeployerAddress
	committeeSize := 100

	b.Run("computeCommittee", func(b *testing.B) {
		stateDB, evmContract, contractAddress, err := deployAutonity(committeeSize, validators, deployer)
		require.NoError(b, err)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "finalizeInitialization")
		require.NoError(b, err)
		packedArgs, err := contractAbi.Pack("computeCommittee")
		require.NoError(b, err)
		benchmarkWithGas(b, evmContract, stateDB, header, contractAddress, packedArgs)
	})
}

func TestComputeCommittee(t *testing.T) {

	t.Run("computeCommittee validators < committee", func(t *testing.T) {
		committeeSize := 100
		validatorCount := 10
		testComputeCommittee(committeeSize, validatorCount, t)
	})

	t.Run("computeCommittee validators > committee", func(t *testing.T) {
		committeeSize := 100
		validatorCount := 1000
		testComputeCommittee(committeeSize, validatorCount, t)
	})
}

func TestElectProposer(t *testing.T) {
	samePowers := []int{100, 100, 100, 100}
	linearPowers := []int{100, 200, 400, 800}
	var ac = &AutonityContract{}
	t.Run("Proposer election should be deterministic", func(t *testing.T) {
		committee := generateCommittee(samePowers)
		for h := uint64(0); h < uint64(100); h++ {
			for r := int64(0); r <= int64(3); r++ {
				proposer1 := ac.electProposer(committee, h, r)
				proposer2 := ac.electProposer(committee, h, r)
				require.Equal(t, proposer1, proposer2)
			}
		}
	})

	t.Run("Proposer selection, print and compare the scheduling rate with same stake", func(t *testing.T) {
		committee := generateCommittee(samePowers)
		maxHeight := uint64(10000)
		maxRound := int64(4)
		//expectedRatioDelta := float64(0.01)
		counterMap := make(map[common.Address]int)
		counterMap[common.Address{}] = 1
		for h := uint64(0); h < maxHeight; h++ {
			for round := int64(0); round < maxRound; round++ {
				proposer := ac.electProposer(committee, h, round)
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

		for i, c := range committee.Members {
			stake := samePowers[i]
			scheduled := counterMap[c.Address]
			log.Print("electing ", "proposer: ", c.Address.String(), " stake: ", stake, " scheduled: ", scheduled)
		}
	})

	t.Run("Proposer selection, print and compare the scheduling rate with liner increasing stake", func(t *testing.T) {
		committee := generateCommittee(linearPowers)
		maxHeight := uint64(1000000)
		maxRound := int64(4)
		//expectedRatioDelta := float64(0.01)
		counterMap := make(map[common.Address]int)
		counterMap[common.Address{}] = 1
		for h := uint64(0); h < maxHeight; h++ {
			for round := int64(0); round < maxRound; round++ {
				proposer := ac.electProposer(committee, h, round)
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

		for _, c := range committee.Members {
			stake := c.VotingPower.Uint64()
			scheduled := counterMap[c.Address]
			log.Print("electing ", "proposer: ", c.Address.String(), " stake: ", stake, " scheduled: ", scheduled)
		}
	})
}

func generateCommittee(powers []int) *types.Committee {
	c := new(types.Committee)
	for _, p := range powers {
		privateKey, _ := crypto.GenerateKey()
		c.Members = append(c.Members, &types.CommitteeMember{
			Address:     crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower: new(big.Int).SetInt64(int64(p)),
		})
	}

	c.Sort()
	return c
}

func deployAutonity(
	committeeSize int, validators []params.Validator, deployer common.Address,
) (*state.StateDB, *EVMContract, common.Address, error) {
	abi := &generated.AutonityAbi
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

// Packs the args and then calls the function and returns result
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

func randomValidators(count int, randomPercentage int) ([]params.Validator, error) {
	if count == 0 {
		return []params.Validator{}, nil
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
			return []params.Validator{}, errors.New("Not sorted")
		}
	}

	validatorList := make([]params.Validator, count)
	// ecdsaSecretKeyList := make([]*ecdsa.PrivateKey, count)
	// blsSecretKeyList := make([]*blst.SecretKey, count)
	for i := 0; i < count; i++ {
		var privateKey *ecdsa.PrivateKey
		var secretKey blst.SecretKey
		var err error
		for {
			privateKey, err = crypto.GenerateKey()
			if err == nil {
				break
			}
		}
		for {
			secretKey, err = blst.RandKey()
			if err == nil {
				break
			}
		}
		// ecdsaSecretKeyList[i] = privateKey
		// blsSecretKeyList[i] = &secretKey
		consensusKey := secretKey.PublicKey()
		publicKey := privateKey.PublicKey
		enode := "enode://" + string(crypto.PubECDSAToHex(&publicKey)[2:]) + "@3.209.45.79:30303"
		address := crypto.PubkeyToAddress(publicKey)
		validatorList[i] = params.Validator{
			Treasury:      address,
			Enode:         enode,
			BondedStake:   big.NewInt(bondedStake[i]),
			OracleAddress: address,
			ConsensusKey:  consensusKey.Marshal(),
		}
		err = validatorList[i].Validate()
		if err != nil {
			return []params.Validator{}, err
		}
	}

	if randomPercentage == 0 || randomPercentage == 100 {
		return validatorList, nil
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
	return validatorList, nil
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
			AccountabilityContract: params.AccountabilityContractAddress,
			OracleContract:         params.OracleContractAddress,
			AcuContract:            params.ACUContractAddress,
			SupplyControlContract:  params.SupplyControlContractAddress,
			StabilizationContract:  params.StabilizationContractAddress,
			UpgradeManagerContract: params.UpgradeManagerContractAddress,
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

func testEVMProvider() func(header *types.Header, origin common.Address, stateDB vm.StateDB) *vm.EVM {
	return func(header *types.Header, origin common.Address, stateDB vm.StateDB) *vm.EVM {
		vmBlockContext := vm.BlockContext{
			Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
			CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
			BlockNumber: common.Big0,
		}
		txContext := vm.TxContext{
			Origin:   common.Address{},
			GasPrice: common.Big0,
		}
		evm := vm.NewEVM(vmBlockContext, txContext, stateDB, params.TestChainConfig, vm.Config{})
		return evm
	}
}

// to properly benchmark a contract call, it is expected that the state is same everytime the contract function is run
func benchmarkWithGas(
	b *testing.B, evmContract *EVMContract, stateDB vm.StateDB, header *types.Header,
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

func isVotersSorted(voters []common.Address, committeeMembers []types.CommitteeMember, validators []params.Validator, enodes []string, totalStake *big.Int) error {
	if len(voters) > len(validators) {
		return fmt.Errorf("More voters than validators")
	}
	if len(voters) != len(committeeMembers) {
		return fmt.Errorf("Committee size not equal to voter size")
	}
	if len(enodes) != len(committeeMembers) {
		return fmt.Errorf("not enough enodes")
	}
	positions := make(map[common.Address]int)
	for i, validator := range validators {
		if _, ok := positions[*validator.NodeAddress]; ok {
			return fmt.Errorf("duplicate validator")
		}
		positions[*validator.NodeAddress] = i
	}
	lastStake := big.NewInt(0)
	totalStakeCalculated := big.NewInt(0)
	for i, member := range committeeMembers {
		idx, ok := positions[member.Address]
		if !ok {
			return fmt.Errorf("committee member not found")
		}
		if i > 0 && lastStake.Cmp(validators[idx].BondedStake) < 0 {
			return fmt.Errorf("not sorted")
		}
		lastStake = validators[idx].BondedStake
		totalStakeCalculated = totalStakeCalculated.Add(totalStakeCalculated, lastStake)

		if !bytes.Equal(member.Address.Bytes(), validators[idx].NodeAddress.Bytes()) {
			return fmt.Errorf("Committee member address mismatch")
		}
		if member.VotingPower.Cmp(validators[idx].BondedStake) != 0 {
			return fmt.Errorf("Committee member stake mismatch")
		}
		if !bytes.Equal(member.ConsensusKey, validators[idx].ConsensusKey) {
			return fmt.Errorf("Committee member consensus key mismatch")
		}
		if enodes[i] != validators[idx].Enode {
			return fmt.Errorf("Committee member enode mismatch")
		}
		if voters[i] != validators[idx].OracleAddress {
			return fmt.Errorf("Oracle address mismatch")
		}
		if *validators[idx].State != uint8(0) {
			return fmt.Errorf("Committee member not active")
		}
		delete(positions, member.Address)
	}
	nonVoter := 0
	for _, validator := range validators {
		if _, ok := positions[*validator.NodeAddress]; ok {
			nonVoter++
			if validator.BondedStake.Cmp(lastStake) == 1 {
				return fmt.Errorf("Non voter has more stake than voter")
			}
		}
	}
	if nonVoter+len(voters) != len(validators) {
		return fmt.Errorf("not all are accounted")
	}
	if totalStakeCalculated.Cmp(totalStake) != 0 {
		return fmt.Errorf("epochTotalStake mismatch")
	}
	return nil
}

func testComputeCommittee(committeeSize int, validatorCount int, t *testing.T) {
	contractAbi := &generated.AutonityTestAbi
	deployer := params.DeployerAddress
	validators, err := randomValidators(validatorCount, 30)
	require.NoError(t, err)
	stateDB, evmContract, contractAddress, err := deployAutonityTest(committeeSize, validators, deployer)
	require.NoError(t, err)
	var header *types.Header
	_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
	require.NoError(t, err)
	res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee")
	require.NoError(t, err)
	voters := make([]common.Address, committeeSize)
	err = contractAbi.UnpackIntoInterface(&voters, "computeCommittee", res)
	require.NoError(t, err)
	res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommittee")
	require.NoError(t, err)
	members := make([]types.CommitteeMember, committeeSize)
	err = contractAbi.UnpackIntoInterface(&members, "getCommittee", res)
	require.NoError(t, err)
	res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getMaxCommitteeSize")
	require.NoError(t, err)
	var maxCommitteeSize *big.Int
	err = contractAbi.UnpackIntoInterface(&maxCommitteeSize, "getMaxCommitteeSize", res)
	require.NoError(t, err)
	require.Equal(t, true, maxCommitteeSize.Cmp(big.NewInt(int64(len(members)))) >= 0, "committee size exceeds MaxCommitteeSize")
	res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getCommitteeEnodes")
	require.NoError(t, err)
	enodes := make([]string, committeeSize)
	err = contractAbi.UnpackIntoInterface(&enodes, "getCommitteeEnodes", res)
	require.NoError(t, err)
	res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getEpochTotalBondedStake")
	require.NoError(t, err)
	totalStake := big.NewInt(0)
	err = contractAbi.UnpackIntoInterface(&totalStake, "getEpochTotalBondedStake", res)
	require.NoError(t, err)
	err = isVotersSorted(voters, members, validators, enodes, totalStake)
	require.NoError(t, err)
}
