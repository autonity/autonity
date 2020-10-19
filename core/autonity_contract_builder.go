package core

import (
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/rawdb"
	"github.com/clearmatics/autonity/ethdb"
	"github.com/clearmatics/autonity/params"
)

func NewAutonityContractFromConfig(db ethdb.Database, hg HeaderGetter, evmP *defaultEVMProvider, autonityConfig *params.AutonityContractGenesis) (*autonity.Contract, error) {
	var JSONString = autonityConfig.ABI
	bytes, err := rawdb.GetKeyValue(db, []byte(autonity.ABISPEC))

	if err != nil {
		return nil, err
	}
	if bytes != nil {
		JSONString = string(bytes)
	}
	return autonity.NewAutonityContract(
		db,
		autonityConfig.Operator,
		autonityConfig.MinGasPrice,
		JSONString,
		evmP,
	)
}
