package autonity

import (
	"log"
	"math/big"
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

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

func BenchmarkContractFunction(b *testing.B) {
	// Deploy contract
	db := state.NewDatabase(rawdb.NewMemoryDatabase())
	state, err := state.New(common.Hash{}, db, nil)
	require.NoError(b, err)
	vmctx := vm.BlockContext{
		Transfer:    func(vm.StateDB, common.Address, common.Address, *big.Int) {},
		CanTransfer: func(vm.StateDB, common.Address, *big.Int) bool { return true },
		BlockNumber: common.Big0,
	}

	txContext := vm.TxContext{
		Origin:   common.Address{},
		GasPrice: common.Big0,
	}

	evm := vm.NewEVM(vmctx, txContext, state, params.TestChainConfig, vm.Config{})
	gas := uint64(math.MaxUint64)
	value := common.Big0

	contractConfig := AutonityConfig{
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

	args, err := generated.AutonityAbi.Pack("", []params.Validator{}, contractConfig)
	require.NoError(b, err)
	data := append(generated.AutonityBytecode, args...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		evm.Create(vm.AccountRef(common.Address{}), data, gas, value)
		//require.NoError(b, err)
		//b.Log(autonity)
		//b.Log("gas used:", gas-leftOverGas)
	}
	//genesisContracts := NewGenesisEVMContract()
	//DeployAutonityContract()
	// get bindings via :
	//NewProtocolContracts
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
