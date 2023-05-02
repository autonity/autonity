package params

import (
	"errors"
	"github.com/autonity/autonity/common/ocdefault"
	"github.com/autonity/autonity/log"
)

const votePeriod = 60

// OracleContractGenesis Autonity contract config. It'is used for deployment.
type OracleContractGenesis struct {
	// Bytecode of validators contract
	// would like this type to be []byte but the unmarshalling is not working
	Bytecode string `json:"bytecode,omitempty" toml:",omitempty"`
	// Json ABI of the contract
	ABI        string   `json:"abi,omitempty" toml:",omitempty"`
	Symbols    []string `json:"symbols"`
	VotePeriod uint64   `json:"votePeriod"`
}

// Prepare prepares the AutonityContractGenesis by filling in missing fields.
// It returns an error if the configuration is invalid.
func (ocg *OracleContractGenesis) Prepare() error {
	if len(ocg.Bytecode) == 0 && len(ocg.ABI) > 0 ||
		len(ocg.Bytecode) > 0 && len(ocg.ABI) == 0 {
		return errors.New("it is an error to set only of oracle contract abi or bytecode")
	}

	if len(ocg.Bytecode) == 0 && len(ocg.ABI) == 0 {
		log.Info("Setting up Oracle default contract")
		ocg.ABI = ocdefault.ABI()
		ocg.Bytecode = ocdefault.Bytecode()
	} else {
		log.Info("Setting up custom Oracle protocol contract")
	}

	if len(ocg.Symbols) == 0 {
		ocg.Symbols = []string{"NTNUSD", "NTNAUD", "NTNCAD", "NTNEUR", "NTNGBP", "NTNJPY", "NTNSEK"}
	}

	if ocg.VotePeriod == 0 {
		ocg.VotePeriod = votePeriod
	}
	return nil
}
