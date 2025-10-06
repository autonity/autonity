package upgrades

import (
	"math/big"

	"github.com/autonity/autonity/accounts/abi"
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
						[]string{
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.0",
							"1.0.1",
							"1.0.1",
							"1.1.0",
							"1.1.0",
						},
					},
				},
			},
			ExclusionList: []*big.Int{
				new(big.Int).SetUint64(params.AutMainnetNetworkID),
				new(big.Int).SetUint64(65010004), // bakerloo
			},
		},
	}
)
