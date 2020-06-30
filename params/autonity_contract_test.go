package params

import (
	"net"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrepareAutonityContract(t *testing.T) {
	key1, _ := crypto.GenerateKey()
	addr1 := crypto.PubkeyToAddress(key1.PublicKey)
	node1 := enode.NewV4(&key1.PublicKey, net.ParseIP("127.0.0.1"), 30303, 0)

	key2, _ := crypto.GenerateKey()
	addr2 := crypto.PubkeyToAddress(key2.PublicKey)
	node2 := enode.NewV4(&key2.PublicKey, net.ParseIP("127.0.0.1"), 30303, 0)

	contractConfig := AutonityContractGenesis{
		Operator: common.HexToAddress("0xff"),
		Bytecode: "some code",
		ABI:      "some abi",
		Users: []User{
			{
				Enode:   node1.String(),
				Type:    UserStakeHolder,
				Address: &addr1,
			},
			{
				Enode:   node2.String(),
				Type:    UserValidator,
				Address: &addr2,
			},
		},
	}
	assert.NoError(t, contractConfig.Prepare())
}

func TestPrepareAutonityContract_ParticipantHaveStake_Fail(t *testing.T) {
	contractConfig := AutonityContractGenesis{
		Bytecode: "some code",
		ABI:      "some abi",
		Operator: common.HexToAddress("0xff"),
		Users: []User{
			{
				Enode: "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				Type:  UserParticipant,
				Stake: 10,
			},
		},
	}
	assert.Error(t, contractConfig.Prepare(), "Expecting Prepare to return error")
}

func TestPrepareAutonityContract_ByteCodeMissed_Fail(t *testing.T) {
	address := common.HexToAddress("0x123")
	contractConfig := AutonityContractGenesis{
		ABI:      "some abi",
		Operator: common.HexToAddress("0xff"),
		Users: []User{
			{
				Address: &address,
				Enode:   "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				Type:    UserParticipant,
			},
		},
	}

	assert.Error(t, contractConfig.Prepare(), "Expecting Prepare to return error")
}

func TestPrepareAutonityContract_InvalidAddrOrEnode_Fail(t *testing.T) {
	t.Skip("Do we need it?")
	address := common.HexToAddress("0x123")
	contractConfig := AutonityContractGenesis{
		Bytecode: "some code",
		ABI:      "some abi",
		Operator: common.HexToAddress("0xff"),

		Users: []User{
			{
				Address: &address,
				Enode:   "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				Type:    UserParticipant,
			},
		},
	}
	assert.Error(t, contractConfig.Prepare(), "Expecting Prepare to return error")
}

func TestPrepareAutonityContract_GovernanceOperatorNotExisted_Fail(t *testing.T) {
	contractConfig := AutonityContractGenesis{
		Bytecode: "some code",
		ABI:      "some abi",
		Users:    []User{},
	}
	assert.Error(t, contractConfig.Prepare(), "Expecting Prepare to return error")
}
func TestPrepareAutonityContract_AddsUserAddress(t *testing.T) {
	contractConfig := &AutonityContractGenesis{
		Bytecode: "some code",
		ABI:      "some abi",
		Users: []User{
			{
				Enode: "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				Type:  UserValidator,
				Stake: 1,
			},
		},
	}
	require.NoError(t, contractConfig.Prepare())
	assert.NotNil(t, contractConfig.Users[0].Address, "Failed to add user address")
}
