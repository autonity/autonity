package autonity

import (
	_ "embed"
	"fmt"
	"math/big"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

// This file contains the list of all protocol contract upgrades
// By default all protocol upgrades are deployed at genesis.
// Exclusions can be made using the `ExclusionList` field. This
// field should be set for network where the upgrade already happened
// via operator transaction.

//go:embed solidity/contracts/Oracle.sol
var oracleCode string

//go:embed solidity/contracts/upgrades/0/Oracle0.sol
var oracle0Code string

//go:embed solidity/contracts/UpgradeManager.sol
var upgradeManagerCode string

//go:embed solidity/contracts/upgrades/1/UpgradeManager1.sol
var upgradeManager1Code string

var upgradeManagerABI = generated.UpgradeManager1Abi

type ProtocolUpgrade struct {
	Upgrades      []ContractUpgrade
	Name          string
	Description   string
	ExclusionList []*big.Int
}

type ContractUpgrade struct {
	Target        params.ProtocolContract
	Abi           *abi.ABI
	Bytecode      []byte
	Hash          common.Hash // codehash
	Args          []interface{}
	VersionString string

	// used for showing code diff with upcheck
	BaseCode     string
	UpgradedCode string
}

// exclusion list containing bakerloo and mainnet
func standardExclusionList() []*big.Int {
	return []*big.Int{
		new(big.Int).Set(params.AutMainnetChainConfig.ChainID),
		new(big.Int).Set(params.BakerlooChainConfig.ChainID),
	}
}

// assembles the deployment bytecode of the upgrade
func (u ContractUpgrade) DeploymentBytecode() ([]byte, error) {
	// if no args are specified, `constructorArgs` will be == []
	constructorArgs, err := u.Abi.Pack("", u.Args...)
	if err != nil {
		return nil, fmt.Errorf("cannot pack upgrade args: %w", err)
	}
	return append(u.Bytecode, constructorArgs...), nil
}

// assembles the calldata to be inserted in an operator tx to do the upgrade
func (u ContractUpgrade) Calldata() ([]byte, error) {
	payload, err := u.DeploymentBytecode()
	if err != nil {
		return nil, fmt.Errorf("cannot build deployment bytecode for %s: %w", u.Target.String(), err)
	}
	var calldata []byte
	if u.VersionString == "" {
		calldata, err = upgradeManagerABI.Pack("upgrade", u.Target, string(payload))
	} else {
		calldata, err = upgradeManagerABI.Pack("upgrade0", u.Target, string(payload), u.VersionString)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot build calldata for %s: %w", u.Target.String(), err)
	}
	return calldata, nil
}

func (p ProtocolUpgrade) Calldata() ([]byte, error) {
	if len(p.Upgrades) == 1 {
		calldata, err := p.Upgrades[0].Calldata()
		if err != nil {
			return nil, fmt.Errorf("cannot build upgrade calldata: %w", err)
		}
		return calldata, nil
	}

	// > 1 contract to upgrade atomically
	addresses := make([]common.Address, 0, len(p.Upgrades))
	bytecodes := make([]string, 0, len(p.Upgrades))
	versionStrings := make([]string, 0, len(p.Upgrades))
	allVersionsEmpty := true
	for _, upgrade := range p.Upgrades {
		deploymentBytecode, err := upgrade.DeploymentBytecode()
		if err != nil {
			return nil, fmt.Errorf("cannot build deployment bytecode for %s: %w", upgrade.Target.String(), err)
		}
		addresses = append(addresses, upgrade.Target.Address())
		bytecodes = append(bytecodes, string(deploymentBytecode))
		versionStrings = append(versionStrings, upgrade.VersionString)
		if upgrade.VersionString != "" {
			allVersionsEmpty = false
		}
	}

	var calldata []byte
	var err error
	if allVersionsEmpty {
		calldata, err = upgradeManagerABI.Pack("upgradeMultiple", addresses, bytecodes)
	} else {
		calldata, err = upgradeManagerABI.Pack("upgradeMultiple0", addresses, bytecodes, versionStrings)
	}
	if err != nil {
		return nil, fmt.Errorf("cannot build calldata for %s: %w", p.Name, err)
	}
	return calldata, nil
}

var (
	// used for all protocol contracts at genesis
	genesisVersion = bindings.UpgradeManager1version{
		Number: "1.0.0",
		Block:  new(big.Int),
	}

	Upgrades = []ProtocolUpgrade{
		// protocol upgrade 0
		{
			Upgrades: []ContractUpgrade{
				// oracle contract bugfix upgrade
				{
					Target:        params.ProtocolContract(params.OracleContractAddress),
					Abi:           &generated.Oracle0Abi,
					Bytecode:      generated.Oracle0Bytecode,
					Hash:          generated.Oracle0CodeHash,
					Args:          nil,
					VersionString: "", // on-chain versioning not implemented yet
					BaseCode:      oracleCode,
					UpgradedCode:  oracle0Code,
				},
			},
			Name: "Oracle Contract slashing bugfix",
			Description: "The `slash()` function was incorrectly called, passing the oracle " +
				"address instead of the node address of the offender. This was causing a " +
				"revert in the `getValidator()` function in the autonity contract when " +
				"trying to do the actual slashing.",
			ExclusionList: standardExclusionList(),
		},
		// protocol upgrade 1
		{
			Upgrades: []ContractUpgrade{
				// upgrade manager with version history upgrade
				{
					Target:   params.ProtocolContract(params.UpgradeManagerContractAddress),
					Abi:      &generated.UpgradeManager1Abi,
					Bytecode: generated.UpgradeManager1Bytecode,
					Hash:     generated.UpgradeManager1CodeHash,
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
							generated.Oracle0CodeHash,
							// upgrade manager version 1.1.0 (this upgrade)
							generated.UpgradeManager1CodeHash,
						},
						[]bindings.UpgradeManager1version{
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
					VersionString: "", // will be set in constructor of the updated code
					BaseCode:      upgradeManagerCode,
					UpgradedCode:  upgradeManager1Code,
				},
			},
			Name: "UpgradeManager contract new features",
			Description: "New features for the upgrade manager: " +
				"1. versioning. Every contract codehash is associated a semver version number (x.x.x). " +
				"2. atomic upgrade of multiple contract at once is now possible with the new function upgradeMultiple.",
			ExclusionList: standardExclusionList(),
		},
	}
)
