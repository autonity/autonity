package params

import (
	"math/big"
	"net"
	"testing"

	"github.com/autonity/autonity/crypto/blst"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
)

func TestPrepareChainConfig(t *testing.T) {
	committeeSize := 21
	contractConfig := AutonityContractGenesis{
		Operator:         common.HexToAddress("0xff"),
		MaxCommitteeSize: uint64(committeeSize),
		EpochPeriod:      50,
	}

	for i := 0; i < committeeSize; i++ {
		treasury, _ := crypto.GenerateKey()
		nodeKey, _ := crypto.GenerateKey()
		oracleKey, _ := crypto.GenerateKey()
		nodeAddr := crypto.PubkeyToAddress(nodeKey.PublicKey)
		enode := enode.NewV4(&nodeKey.PublicKey, net.ParseIP("127.0.0.1"), 30303, 0)
		consensusKey, err := blst.RandKey()
		require.NoError(t, err)

		validator := &Validator{
			Treasury:      crypto.PubkeyToAddress(treasury.PublicKey),
			Enode:         enode.String(),
			NodeAddress:   &nodeAddr,
			OracleAddress: crypto.PubkeyToAddress(oracleKey.PublicKey),
			BondedStake:   big.NewInt(1),
			ConsensusKey:  consensusKey.PublicKey().Marshal(),
		}
		contractConfig.Validators = append(contractConfig.Validators, validator)
	}
	chainConfig := ChainConfig{AutonityContractConfig: &contractConfig, OmissionAccountabilityConfig: DefaultOmissionAccountabilityConfig}
	assert.NoError(t, chainConfig.Prepare())
}

func TestPrepareChainConfig_ParticipantHaveStake_Fail(t *testing.T) {
	contractConfig := AutonityContractGenesis{
		Operator:         common.HexToAddress("0xff"),
		MaxCommitteeSize: 21,
		EpochPeriod:      50,
		Validators: []*Validator{
			{
				Treasury: common.Address{},
				Enode:    "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
			},
		},
	}
	chainConfig := ChainConfig{AutonityContractConfig: &contractConfig, OmissionAccountabilityConfig: DefaultOmissionAccountabilityConfig}
	err := chainConfig.Prepare()
	t.Log(err)
	assert.Error(t, err, "Expecting Prepare to return error")
}

func TestPrepareChainConfig_InvalidAddrOrEnode_Fail(t *testing.T) {
	t.Skip("Do we need it?")
	address := common.HexToAddress("0x123")
	contractConfig := AutonityContractGenesis{
		Operator:         common.HexToAddress("0xff"),
		MaxCommitteeSize: 21,
		EpochPeriod:      50,
		Validators: []*Validator{
			{
				NodeAddress: &address,
				Enode:       "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				BondedStake: big.NewInt(10),
				Treasury:    common.Address{},
			},
		},
	}
	chainConfig := ChainConfig{AutonityContractConfig: &contractConfig, OmissionAccountabilityConfig: DefaultOmissionAccountabilityConfig}
	err := chainConfig.Prepare()
	t.Log(err)
	assert.Error(t, err, "Expecting Prepare to return error")
}

func TestPrepareChainConfig_GovernanceOperatorNotExisted_Fail(t *testing.T) {
	// TODO: refactor this actually fails because there are no validators
	contractConfig := AutonityContractGenesis{
		MaxCommitteeSize: 21,
		EpochPeriod:      50,
		Validators:       []*Validator{},
	}
	chainConfig := ChainConfig{AutonityContractConfig: &contractConfig, OmissionAccountabilityConfig: DefaultOmissionAccountabilityConfig}
	err := chainConfig.Prepare()
	t.Log(err)
	assert.Error(t, err, "Expecting Prepare to return error")
}

func TestPrepareChainConfig_EpochPeriod(t *testing.T) {
	committeeSize := 21
	contractConfig := AutonityContractGenesis{
		Operator:         common.HexToAddress("0xff"),
		MaxCommitteeSize: uint64(committeeSize),
		EpochPeriod:      30,
	}

	for i := 0; i < committeeSize; i++ {
		treasury, _ := crypto.GenerateKey()
		nodeKey, _ := crypto.GenerateKey()
		oracleKey, _ := crypto.GenerateKey()
		nodeAddr := crypto.PubkeyToAddress(nodeKey.PublicKey)
		enode := enode.NewV4(&nodeKey.PublicKey, net.ParseIP("127.0.0.1"), 30303, 0)
		consensusKey, err := blst.RandKey()
		require.NoError(t, err)

		validator := &Validator{
			Treasury:      crypto.PubkeyToAddress(treasury.PublicKey),
			Enode:         enode.String(),
			NodeAddress:   &nodeAddr,
			OracleAddress: crypto.PubkeyToAddress(oracleKey.PublicKey),
			BondedStake:   big.NewInt(1),
			ConsensusKey:  consensusKey.PublicKey().Marshal(),
		}
		contractConfig.Validators = append(contractConfig.Validators, validator)
	}
	chainConfig := ChainConfig{AutonityContractConfig: &contractConfig, OmissionAccountabilityConfig: &OmissionAccountabilityGenesis{
		InactivityThreshold:    1000, // 10%
		LookbackWindow:         30,   // 30 blocks
		PastPerformanceWeight:  1000, // 10%
		InitialJailingPeriod:   300,  // 300 blocks
		InitialProbationPeriod: 24,   // 24 epochs
		InitialSlashingRate:    1000, // 10%
		Delta:                  10,   // 10 blocks
	}}
	// equation epochPeriod > delta+lookback-1 needs to be respected
	// 30 > 10+30-1 --> false --> err
	err := chainConfig.Prepare()
	t.Log(err)
	assert.Error(t, err, "Expecting Prepare to return error")
	// 30 > 10+20-1 --> true --> no err
	chainConfig.OmissionAccountabilityConfig.LookbackWindow = 20
	assert.NoError(t, chainConfig.Prepare())
	// 30 > 20+20-1 --> false --> err
	chainConfig.OmissionAccountabilityConfig.Delta = 20
	err = chainConfig.Prepare()
	t.Log(err)
	assert.Error(t, err, "Expecting Prepare to return error")
	// 40 > 20+20-1 --> true --> no err
	contractConfig.EpochPeriod = 40
	assert.NoError(t, chainConfig.Prepare())
	// 400 > 20+20-1 --> true --> no err
	contractConfig.EpochPeriod = 400
	assert.NoError(t, chainConfig.Prepare())
}

func TestPrepareAutonityContract_AddsUserAddress(t *testing.T) {
	treasury, _ := crypto.GenerateKey()
	nodeKey, _ := crypto.GenerateKey()
	oracleKey, _ := crypto.GenerateKey()
	enode := enode.NewV4(&nodeKey.PublicKey, net.ParseIP("127.0.0.1"), 30303, 0)
	consensusKey, err := blst.RandKey()
	require.NoError(t, err)

	contractConfig := &AutonityContractGenesis{
		MaxCommitteeSize: 21,
		EpochPeriod:      50,
		Validators: []*Validator{
			{
				Treasury:      crypto.PubkeyToAddress(treasury.PublicKey),
				Enode:         enode.String(),
				OracleAddress: crypto.PubkeyToAddress(oracleKey.PublicKey),
				ConsensusKey:  consensusKey.PublicKey().Marshal(),
				BondedStake:   big.NewInt(1),
			},
		},
	}
	chainConfig := ChainConfig{AutonityContractConfig: contractConfig, OmissionAccountabilityConfig: DefaultOmissionAccountabilityConfig}
	require.NoError(t, chainConfig.Prepare())
	assert.NotNil(t, contractConfig.Validators[0].NodeAddress, "Failed to add user address")
}

func TestPrepareAutonityContract_CommitteSizeNotProvided_Fail(t *testing.T) {
	contractConfig := AutonityContractGenesis{
		Operator:    common.HexToAddress("0xff"),
		EpochPeriod: 50,
		Validators:  []*Validator{},
	}
	chainConfig := ChainConfig{AutonityContractConfig: &contractConfig, OmissionAccountabilityConfig: DefaultOmissionAccountabilityConfig}
	err := chainConfig.Prepare()
	t.Log(err)
	assert.Error(t, err, "Expecting Prepare to return error")
}
