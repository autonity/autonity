package autonity

import (
	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"math/big"
	"strings"
)

type OracleContract struct {
	evm             *vm.EVM
	autonityAddress common.Address
	genesis         *params.OracleContractGenesis
}

var oracleABI abi.ABI
var oracleContractAddress common.Address

func GetOracleABI() abi.ABI {
	return oracleABI
}

func GetOracleContractAddress() common.Address {
	return oracleContractAddress
}

func NewOracleContract(ocg *params.OracleContractGenesis, evm *vm.EVM) *OracleContract {
	return &OracleContract{
		autonityAddress: ContractAddress,
		evm:             evm,
		genesis:         ocg,
	}
}

func (oc *OracleContract) Deploy(voters []common.Address, operator common.Address) error {
	var err error
	oracleABI, err = abi.JSON(strings.NewReader(oc.genesis.ABI))
	if err != nil {
		return nil
	}
	contractBytecode := common.Hex2Bytes(oc.genesis.Bytecode)
	constructorParams, err := oracleABI.Pack("",
		voters,
		oc.autonityAddress,
		operator, // same operator as autonity
		oc.genesis.Symbols,
		new(big.Int).SetUint64(oc.genesis.VotePeriod))

	if err != nil {
		log.Error("contractABI.Pack returns err", "err", err)
		return err
	}

	data := append(contractBytecode, constructorParams...)
	gas := uint64(0xFFFFFFFF)
	value := new(big.Int).SetUint64(0x00)

	// Deploy the Autonity contract
	_, oracleContractAddress, _, err = oc.evm.Create(vm.AccountRef(Deployer), data, gas, value)
	if err != nil {
		log.Error("DeployOracleContract evm.Create", "err", err)
		return err
	}
	log.Info("Deployed Oracle Contract", "Address", oracleContractAddress.String())

	return nil

}
