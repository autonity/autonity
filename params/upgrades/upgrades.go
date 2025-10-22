package upgrades

import (
	"math/big"

	"github.com/autonity/autonity/accounts/abi"
	bindings1 "github.com/autonity/autonity/autonity/bindings/1"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
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

// TODO: move somewhere else?
func concatHashes(hashes ...common.Hash) []byte {
	var b []byte
	for _, hash := range hashes {
		b = append(b, hash.Bytes()...)
	}
	return b
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
					Abi:      &generated0.OracleAbi,
					Bytecode: generated0.OracleBytecode,
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
					Abi:      &generated1.UpgradeManagerAbi,
					Bytecode: generated1.UpgradeManagerBytecode,
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
							// protocol group version 1.0.0
							crypto.Keccak256Hash(
								concatHashes(
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
								),
							),
							// ASM group version 1.0.0
							crypto.Keccak256Hash(
								concatHashes(
									generated.ACUCodeHash,
									generated.SupplyControlCodeHash,
									generated.StabilizationCodeHash,
									generated.InflationControllerCodeHash,
									generated.AuctioneerCodeHash,
								),
							),
							// oracle version 1.0.1
							generated0.OracleCodeHash,
							// protocol group version 1.0.1
							crypto.Keccak256Hash(
								concatHashes(
									generated.AutonityCodeHash,
									generated.AccountabilityCodeHash,
									generated0.OracleCodeHash, // oracle 1.0.1
									generated.ACUCodeHash,
									generated.SupplyControlCodeHash,
									generated.StabilizationCodeHash,
									generated.UpgradeManagerCodeHash,
									generated.InflationControllerCodeHash,
									generated.OmissionAccountabilityCodeHash,
									generated.AuctioneerCodeHash,
								),
							),
							// upgrade manager version 1.1.0 (this upgrade)
							generated1.UpgradeManagerCodeHash,
							// protocol group version 1.1.0
							crypto.Keccak256Hash(
								concatHashes(
									generated.AutonityCodeHash,
									generated.AccountabilityCodeHash,
									generated0.OracleCodeHash, // oracle 1.0.1
									generated.ACUCodeHash,
									generated.SupplyControlCodeHash,
									generated.StabilizationCodeHash,
									generated1.UpgradeManagerCodeHash, // upgrade manager 1.1.0
									generated.InflationControllerCodeHash,
									generated.OmissionAccountabilityCodeHash,
									generated.AuctioneerCodeHash,
								),
							),
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
							genesisVersion,
							genesisVersion,
							// oracle contract upgrade
							{
								Number: "1.0.1",
								Block:  new(big.Int),
							},
							{
								Number: "1.0.1",
								Block:  new(big.Int),
							},
							// upgrade manager contract upgrade
							{
								Number: "1.1.0",
								Block:  new(big.Int),
							},
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
