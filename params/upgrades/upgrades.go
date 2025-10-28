package upgrades

import (
	"math/big"

	"github.com/autonity/autonity/accounts/abi"
	bindings1 "github.com/autonity/autonity/autonity/bindings/1"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
	generated0 "github.com/autonity/autonity/params/upgrades/generated/0"
	generated1 "github.com/autonity/autonity/params/upgrades/generated/1"
)

// This file serves the purpose of defining exclusion lists for all ... TODO

type ProtocolUpgrade struct {
	Upgrades      []ContractUpgrade
	Description   string
	ExclusionList []*big.Int
}

type ContractUpgrade struct {
	Target   params.ProtocolContract
	Abi      *abi.ABI
	Bytecode []byte
	Args     []interface{}
}

var (
	// used for all protocol contracts at genesis
	genesisVersion = bindings1.UpgradeManager1version{
		Number: "1.0.0",
		Block:  new(big.Int),
	}

	Upgrades = []ProtocolUpgrade{
		// protocol upgrade 0
		{
			Upgrades: []ContractUpgrade{
				// oracle contract bugfix upgrade
				{
					Target:   params.ProtocolContract(params.OracleContractAddress),
					Abi:      &generated0.Oracle0Abi,
					Bytecode: generated0.Oracle0Bytecode,
					Args:     nil,
				},
			},
			Description: "The `slash()` function was incorrectly called, passing the oracle" +
				"\naddress instead of the node address of the offender. This was causing a" +
				"\nrevert in the `getValidator()` function in the autonity contract when" +
				"\ntrying to do the actual slashing.",
			ExclusionList: []*big.Int{
				new(big.Int).SetUint64(params.AutMainnetNetworkID),
				new(big.Int).SetUint64(65010004), // bakerloo
			},
		},
		// protocol upgrade 1
		{
			Upgrades: []ContractUpgrade{
				// upgrade manager with version history upgrade
				{
					Target:   params.ProtocolContract(params.UpgradeManagerContractAddress),
					Abi:      &generated1.UpgradeManager1Abi,
					Bytecode: generated1.UpgradeManager1Bytecode,
					Args: []interface{}{
						[]common.Hash{
							// version 1.0.0 contract hashes
							generated.AutonityCodeHash,
							generated.AccountabilityCodeHash,
							generated.OracleCodeHash,
							generated.ACUCodeHash,
							generated.SupplyControlCodeHash,
							generated.StabilizationCodeHash,
							generated.UpgradeManagerCodeHash,
							generated.InflationControllerCodeHash,
							generated.OmissionAccountabilityCodeHash,
							generated.AuctioneerCodeHash,
							// oracle version 1.0.1
							generated0.Oracle0CodeHash,
							// upgrade manager version 1.1.0 (this upgrade)
							generated1.UpgradeManager1CodeHash,
						},
						[]bindings1.UpgradeManager1version{
							genesisVersion,
							genesisVersion,
							genesisVersion,
							genesisVersion,
							genesisVersion,
							genesisVersion,
							genesisVersion,
							genesisVersion,
							genesisVersion,
							genesisVersion,
							// oracle contract upgrade
							{
								Number: "1.0.1",
								Block:  new(big.Int),
							},
							// upgrade manager contract upgrade
							{
								Number: "1.1.0",
								Block:  new(big.Int),
							},
						},
					},
				},
			},
			Description: "TODO",
			ExclusionList: []*big.Int{
				new(big.Int).SetUint64(params.AutMainnetNetworkID),
				new(big.Int).SetUint64(65010004), // bakerloo
			},
		},
	}
)
