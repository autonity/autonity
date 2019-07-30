package core

import (
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/p2p/enode"
	"reflect"
	"testing"
)

func TestValidateAutonityContract(t *testing.T) {
	contractConfig := AutonityContract{
		Deployer:           common.HexToAddress("0xff"),
		GovernanceOperator: common.HexToAddress("0xff"),
		Bytecode:           "some code",
		ABI:                "some abi",
		Users: []User{
			{
				Enode: "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				Type:  UserMember,
			},
			{
				Enode: "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				Type:  UserValidator,
			},
		},
	}
	err := contractConfig.Validate()
	if err != nil {
		t.Fatal(err)
	}

	for i := range contractConfig.Users {
		if reflect.DeepEqual(contractConfig.Users[i].Address, common.Address{}) {
			t.Fatal("Empty address")
		}
	}
}
func TestValidateAutonityContract_ParticipantHaveStake_Fail(t *testing.T) {
	contractConfig := AutonityContract{
		Deployer:           common.HexToAddress("0xff"),
		Bytecode:           "some code",
		ABI:                "some abi",
		GovernanceOperator: common.HexToAddress("0xff"),
		Users: []User{
			{
				Enode: "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				Type:  UserParticipant,
				Stake: 10,
			},
		},
	}
	err := contractConfig.Validate()
	if err == nil {
		t.FailNow()
	}

}
func TestValidateAutonityContract_ByteCodeMissed_Fail(t *testing.T) {
	contractConfig := AutonityContract{
		Deployer:           common.HexToAddress("0xff"),
		ABI:                "some abi",
		GovernanceOperator: common.HexToAddress("0xff"),
		Users: []User{
			{
				Address: common.HexToAddress("0x123"),
				Enode:   "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
				Type:    UserParticipant,
			},
		},
	}
	err := contractConfig.Validate()
	if err == nil {
		t.FailNow()
	}
}
func TestValidateAutonityContract_GovernanceOperatorNotExisted_Fail(t *testing.T) {
	contractConfig := AutonityContract{
		Deployer: common.HexToAddress("0xff"),
		Bytecode: "some code",
		ABI:      "some abi",
		Users:    []User{},
	}
	err := contractConfig.Validate()
	if err == nil {
		t.FailNow()
	}
}

func TestUsePartOfEnodeAsAddress(t *testing.T) {
	k, err := crypto.GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	pubkey := k.PublicKey
	expected := crypto.PubkeyToAddress(pubkey)

	en := enode.PubkeyToIDV4(&pubkey)
	addrByte := en.Bytes()[12:]
	addr := common.BytesToAddress(addrByte)
	if expected.String() != addr.String() {
		t.Fatal(expected.String() + " != " + addr.String())
	}

}
