package params

import (
	"github.com/autonity/autonity/crypto/blst"
	"math/big"
	"net"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/p2p/enode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareAutonityContract(t *testing.T) {
	committeeSize := 21
	contractConfig := AutonityContractGenesis{
		Operator:         common.HexToAddress("0xff"),
		MaxCommitteeSize: uint64(committeeSize),
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
	assert.NoError(t, contractConfig.Prepare())
}

func TestPrepareAutonityContract_ParticipantHaveStake_Fail(t *testing.T) {
	contractConfig := AutonityContractGenesis{
		Operator:         common.HexToAddress("0xff"),
		MaxCommitteeSize: 21,
		Validators: []*Validator{
			{
				Treasury: common.Address{},
				Enode:    "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
			},
		},
	}
	assert.Error(t, contractConfig.Prepare(), "Expecting Prepare to return error")
}

func TestPrepareAutonityContract_InvalidAddrOrEnode_Fail(t *testing.T) {
	t.Skip("Do we need it?")
	address := common.HexToAddress("0x123")
	contractConfig := AutonityContractGenesis{
		Operator:         common.HexToAddress("0xff"),
		MaxCommitteeSize: 21,
		Validators: []*Validator{
			{
				NodeAddress: &address,
				Enode:       "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				BondedStake: big.NewInt(10),
				Treasury:    common.Address{},
			},
		},
	}
	assert.Error(t, contractConfig.Prepare(), "Expecting Prepare to return error")
}

func TestPrepareAutonityContract_GovernanceOperatorNotExisted_Fail(t *testing.T) {
	contractConfig := AutonityContractGenesis{
		MaxCommitteeSize: 21,
		Validators:       []*Validator{},
	}
	assert.Error(t, contractConfig.Prepare(), "Expecting Prepare to return error")
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
	require.NoError(t, contractConfig.Prepare())
	assert.NotNil(t, contractConfig.Validators[0].NodeAddress, "Failed to add user address")
}

func TestPrepareAutonityContract_CommitteSizeNotProvided_Fail(t *testing.T) {
	contractConfig := AutonityContractGenesis{
		Operator:   common.HexToAddress("0xff"),
		Validators: []*Validator{},
	}
	assert.Error(t, contractConfig.Prepare(), "Expecting Prepare to return error")
}
