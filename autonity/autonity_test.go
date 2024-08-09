package autonity

import (
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net"
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
	"github.com/autonity/autonity/p2p/enode"
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

func TestUpdateEnode(t *testing.T) {

	contractAbi := &generated.AutonityTestAbi
	deployer := params.DeployerAddress
	validators, err := randomValidators(10, 100)
	require.NoError(t, err)

	t.Run("Success: update node ip of enode", func(t *testing.T) {
		stateDB, evmContract, contractAddress, _ := deployAutonity(5, validators, deployer)
		var header *types.Header
		val := validators[0]
		node, _ := enode.ParseV4(val.Enode)
		tempNode := enode.NewV4(node.Pubkey(), net.IP{127, 0, 0, 1}, 1234, 0)
		_, err = callContractFunctionAs(evmContract, contractAddress,
			stateDB, header, contractAbi, val.Treasury,
			"updateEnode", val.NodeAddress, tempNode.String())
		require.NoError(t, err)

		// make sure validator enode has taken new value
		res, err := callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getValidator", val.NodeAddress)
		require.NoError(t, err)
		av, err := contractAbi.Unpack("getValidator", res)
		require.NoError(t, err)
		out := abi.ConvertType(av[0], new(AutonityValidator)).(*AutonityValidator)
		require.Equal(t, tempNode.String(), out.Enode)
	})

	t.Run("Failure: update to invalid enode format", func(t *testing.T) {
		stateDB, evmContract, contractAddress, _ := deployAutonity(5, validators, deployer)
		var header *types.Header
		val := validators[0]
		testKey, _ := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")
		tempNode := enode.NewV4(&testKey.PublicKey, net.IP{127, 0, 0, 1}, 1234, 0)
		res, err := callContractFunctionAs(evmContract, contractAddress,
			stateDB, header, contractAbi, val.Treasury,
			"updateEnode", val.NodeAddress, tempNode.String()+"junk")
		require.Error(t, err)

		revertReason, _ := abi.UnpackRevert(res)
		require.Equal(t, "enode error", revertReason)
	})

	t.Run("Failure: update node address of enode", func(t *testing.T) {
		stateDB, evmContract, contractAddress, _ := deployAutonity(5, validators, deployer)
		var header *types.Header
		val := validators[0]
		testKey, _ := crypto.HexToECDSA("45a915e4d060149eb4365960e6a7a45f334393093061116b197e3240065ff2d8")
		tempNode := enode.NewV4(&testKey.PublicKey, net.IP{127, 0, 0, 1}, 1234, 0)
		res, err := callContractFunctionAs(evmContract, contractAddress,
			stateDB, header, contractAbi, val.Treasury,
			"updateEnode", val.NodeAddress, tempNode.String())
		require.Error(t, err)

		revertReason, _ := abi.UnpackRevert(res)
		require.Equal(t, "validator node address can't be updated", revertReason)
	})

	t.Run("Failure: update ip of enode, unknown msg.sender", func(t *testing.T) {
		stateDB, evmContract, contractAddress, _ := deployAutonity(5, validators, deployer)
		var header *types.Header
		val := validators[0]
		node, _ := enode.ParseV4(val.Enode)
		res, err := callContractFunctionAs(evmContract, contractAddress,
			stateDB, header, contractAbi, validators[1].Treasury,
			"updateEnode", val.NodeAddress, node.String())
		require.Error(t, err)

		revertReason, _ := abi.UnpackRevert(res)
		require.Equal(t, "require caller to be validator treasury account", revertReason)
	})

	t.Run("Failure: update ip of enode, validator not registered", func(t *testing.T) {
		stateDB, evmContract, contractAddress, _ := deployAutonity(5, validators, deployer)
		var header *types.Header
		randVals, _ := randomValidators(1, 100)
		val := randVals[0]
		node, _ := enode.ParseV4(val.Enode)
		res, err := callContractFunctionAs(evmContract, contractAddress,
			stateDB, header, contractAbi, val.Treasury,
			"updateEnode", val.NodeAddress, node.String())
		require.Error(t, err)

		revertReason, _ := abi.UnpackRevert(res)
		require.Equal(t, "validator not registered", revertReason)
	})

	t.Run("Failure: update ip of enode, validator In Committee", func(t *testing.T) {
		// update committee size to 10, so all the validators will be in committeee
		stateDB, evmContract, contractAddress, _ := deployAutonity(10, validators, deployer)
		var header *types.Header
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "applyStakingOperations")
		require.NoError(t, err)
		_, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "computeCommittee")
		require.NoError(t, err)

		val := validators[0]
		node, _ := enode.ParseV4(val.Enode)
		tempNode := enode.NewV4(node.Pubkey(), net.IP{127, 0, 0, 1}, 1234, 0)
		res, err := callContractFunctionAs(evmContract, contractAddress,
			stateDB, header, contractAbi, val.Treasury,
			"updateEnode", val.NodeAddress, tempNode.String())
		require.Error(t, err)
		revertReason, _ := abi.UnpackRevert(res)
		require.Equal(t, "validator must not be in committee", revertReason)

		// make sure validator enode remains unchanged
		res, err = callContractFunction(evmContract, contractAddress, stateDB, header, contractAbi, "getValidator", val.NodeAddress)
		require.NoError(t, err)
		av, err := contractAbi.Unpack("getValidator", res)
		require.NoError(t, err)
		out := abi.ConvertType(av[0], new(AutonityValidator)).(*AutonityValidator)
		require.Equal(t, node.String(), out.Enode)
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

// Packs the args and then calls the function and returns result
func callContractFunctionAs(
	evmContract *EVMContract, contractAddress common.Address, stateDB *state.StateDB, header *types.Header, abi *abi.ABI,
	origin common.Address, methodName string, args ...interface{}, //nolint:unparam
) ([]byte, error) {
	argsPacked, err := abi.Pack(methodName, args...)
	if err != nil {
		return make([]byte, 0), err
	}
	return evmContract.CallContractFuncAs(stateDB, header, contractAddress, origin, argsPacked)
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
			TreasuryFee:             new(big.Int).SetUint64(params.TestAutonityContractConfig.TreasuryFee),
			MinBaseFee:              new(big.Int).SetUint64(params.TestAutonityContractConfig.MinBaseFee),
			DelegationRate:          new(big.Int).SetUint64(params.TestAutonityContractConfig.DelegationRate),
			UnbondingPeriod:         new(big.Int).SetUint64(params.TestAutonityContractConfig.UnbondingPeriod),
			InitialInflationReserve: (*big.Int)(params.TestAutonityContractConfig.InitialInflationReserve),
			WithholdingThreshold:    new(big.Int).SetUint64(params.TestAutonityContractConfig.WithholdingThreshold),
			ProposerRewardRate:      new(big.Int).SetUint64(params.TestAutonityContractConfig.ProposerRewardRate),
			WithheldRewardsPool:     params.TestAutonityContractConfig.Operator,
			TreasuryAccount:         params.TestAutonityContractConfig.Operator,
		},
		Contracts: AutonityContracts{
			AccountabilityContract:         params.AccountabilityContractAddress,
			OmissionAccountabilityContract: params.OmissionAccountabilityContractAddress,
			OracleContract:                 params.OracleContractAddress,
			AcuContract:                    params.ACUContractAddress,
			SupplyControlContract:          params.SupplyControlContractAddress,
			StabilizationContract:          params.StabilizationContractAddress,
			UpgradeManagerContract:         params.UpgradeManagerContractAddress,
			InflationControllerContract:    params.InflationControllerContractAddress,
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

func isSorted(voters []common.Address, nodeAddresses []common.Address, treasuries []common.Address, committeeMembers []types.CommitteeMember, validators []params.Validator, enodes []string, totalStake *big.Int) error {
	if len(voters) != len(nodeAddresses) || len(voters) != len(treasuries) {
		return fmt.Errorf("Mismatch between addresses length")
	}
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
		if !bytes.Equal(member.ConsensusKeyBytes, validators[idx].ConsensusKey) {
			return fmt.Errorf("Committee member consensus key mismatch")
		}
		if enodes[i] != validators[idx].Enode {
			return fmt.Errorf("Committee member enode mismatch")
		}
		if voters[i] != validators[idx].OracleAddress {
			return fmt.Errorf("Oracle address mismatch")
		}
		if nodeAddresses[i] != *validators[idx].NodeAddress {
			return fmt.Errorf("Node address mismatch")
		}
		if treasuries[i] != validators[idx].Treasury {
			return fmt.Errorf("Treasury account mismatch")
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
	addresses, err := contractAbi.Unpack("computeCommittee", res)
	voters := addresses[0].([]common.Address)
	nodeAddresses := addresses[1].([]common.Address)
	treasuries := addresses[2].([]common.Address)
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
	err = isSorted(voters, nodeAddresses, treasuries, members, validators, enodes, totalStake)
	require.NoError(t, err)
}
