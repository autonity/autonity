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

	treasury1, _ := crypto.GenerateKey()
	nodeKey1, _ := crypto.GenerateKey()
	oracleKey1, _ := crypto.GenerateKey()
	nodeAddr1 := crypto.PubkeyToAddress(nodeKey1.PublicKey)
	enode1 := enode.NewV4(&nodeKey1.PublicKey, net.ParseIP("127.0.0.1"), 30303, 0)
	consensusKey1, err := blst.RandKey()
	require.NoError(t, err)

	treasury2, _ := crypto.GenerateKey()
	nodeKey2, _ := crypto.GenerateKey()
	oracleKey2, _ := crypto.GenerateKey()
	nodeAddr2 := crypto.PubkeyToAddress(nodeKey2.PublicKey)
	enode2 := enode.NewV4(&nodeKey2.PublicKey, net.ParseIP("127.0.0.1"), 30303, 0)
	consensusKey2, err := blst.RandKey()
	require.NoError(t, err)

	contractConfig := AutonityContractGenesis{
		Operator:         common.HexToAddress("0xff"),
		MaxCommitteeSize: 21,
		Validators: []*Validator{
			{
				Treasury:      crypto.PubkeyToAddress(treasury1.PublicKey),
				Enode:         enode1.String(),
				NodeAddress:   &nodeAddr1,
				OracleAddress: crypto.PubkeyToAddress(oracleKey1.PublicKey),
				BondedStake:   big.NewInt(1),
				ConsensusKey:  consensusKey1.PublicKey().Marshal(),
			},
			{
				Treasury:      crypto.PubkeyToAddress(treasury2.PublicKey),
				Enode:         enode2.String(),
				NodeAddress:   &nodeAddr2,
				OracleAddress: crypto.PubkeyToAddress(oracleKey2.PublicKey),
				BondedStake:   big.NewInt(1),
				ConsensusKey:  consensusKey2.PublicKey().Marshal(),
			},
		},
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
