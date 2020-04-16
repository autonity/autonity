package params

import (
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/acdefault"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
	"reflect"
)

const (
	//participant: Authorized to operate a full node, is able to join the network, is not authorized to own stake.
	UserParticipant = "participant"
	//member: Authorized to operate a full node, is able to join the network, authorized to own stake.
	UserStakeHolder = "stakeholder"
	//validator: Authorized to operate a full node, is able to join the network, authorized to own stake, participate in consensus.
	UserValidator = "validator"
)

var userTypeID = map[UserType]int{
	UserParticipant: 0,
	UserStakeHolder: 1,
	UserValidator:   2,
}

type UserType string

func (ut UserType) IsValid() bool {
	if ut == UserStakeHolder || ut == UserValidator || ut == UserParticipant {
		return true
	}
	return false
}

func (ut UserType) GetID() int {
	return userTypeID[ut]
}

// Autonity contract config. It'is used for deployment.
type AutonityContractGenesis struct {
	// Address of the validator who deploys contract stored in bytecode
	Deployer common.Address `json:"deployer" toml:",omitempty"`
	// Bytecode of validators contract
	// would like this type to be []byte but the unmarshalling is not working
	Bytecode string `json:"bytecode" toml:",omitempty"`
	// Json ABI of the contract
	ABI         string         `json:"abi" toml:",omitempty"`
	MinGasPrice uint64         `json:"minGasPrice" toml:",omitempty"`
	Operator    common.Address `json:"operator" toml:",omitempty"`
	Users       []User         `json:"users" toml:",omitempty"`
}

func (ac *AutonityContractGenesis) AddDefault() *AutonityContractGenesis {
	if len(ac.Bytecode) == 0 || len(ac.ABI) == 0 {
		log.Info("Default Validator smart contract set")
		ac.ABI = acdefault.ABI()
		ac.Bytecode = acdefault.Bytecode()
	} else {
		log.Info("User specified Validator smart contract set")
	}
	if reflect.DeepEqual(ac.Deployer, common.Address{}) {
		ac.Deployer = acdefault.Deployer()
	}
	if reflect.DeepEqual(ac.Operator, common.Address{}) {
		ac.Operator = acdefault.Governance()
	}

	for i := range ac.Users {
		if reflect.DeepEqual(ac.Users[i].Address, common.Address{}) {
			if n, err := enode.ParseV4WithResolve(ac.Users[i].Enode); n != nil {
				ac.Users[i].Address = EnodeToAddress(n)
			} else {
				log.Error("Error parsing enode", "enode", ac.Users[i].Enode, "err", err)
			}
		}

		//TODO: Do we need to check if enodeToAddress corresponds to the address provided (in case is not {})
	}
	return ac
}

func (ac *AutonityContractGenesis) Validate() error {
	if len(ac.Bytecode) == 0 || len(ac.ABI) == 0 {
		return errors.New("autonity contract is empty")
	}

	if reflect.DeepEqual(ac.Deployer, common.Address{}) {
		return errors.New("deployer is empty")
	}

	if reflect.DeepEqual(ac.Operator, common.Address{}) {
		return errors.New("governance operator is empty")
	}

	for i := range ac.Users {
		if err := ac.Users[i].Validate(); err != nil {
			return fmt.Errorf("user %v got validation error %v", ac.Users[i], err)
		}
	}

	if len(ac.GetValidatorUsers()) == 0 {
		return errors.New("validators list is empty")
	}

	return nil
}

func (ac *AutonityContractGenesis) GetContractAddress() (common.Address, error) {
	if reflect.DeepEqual(ac.Deployer, common.Address{}) {
		return common.Address{}, errors.New("deployer must be not nil")
	}
	return crypto.CreateAddress(ac.Deployer, 0), nil
}

//User - is used to put predefined accounts to genesis
type User struct {
	Address common.Address `json:"address"`
	Enode   string         `json:"enode"`
	Type    UserType       `json:"type"`
	Stake   uint64         `json:"stake"`
}

func (u *User) Validate() error {
	if !u.Type.IsValid() {
		return errors.New("incorrect user type")
	}

	if reflect.DeepEqual(u.Address, common.Address{}) && len(u.Enode) == 0 { //TODO: Check if && is correct here
		return errors.New("user.enode or user.address must be defined")
	}

	if u.Type == UserParticipant && u.Stake > 0 {
		return errors.New("user.stake must be nil or equal to 0 for users of type participant")
	}

	if u.Type == UserValidator && len(u.Enode) == 0 {
		return errors.New("if user.type is validator then user.enode must be defined")
	}
	if len(u.Enode) > 0 {
		n, err := enode.ParseV4WithResolve(u.Enode)
		if n == nil {
			return fmt.Errorf("fail to parse enode for account %v, error:%v", u.Address, err)
		}

		//todo do we need this check? --> I think we do!
		//if reflect.DeepEqual(u.Address, common.Address{}) {
		//	return errors.New("if both user.enode and user.address are defined, then the derived address from user.enode must be equal to user.address")
		//}
	}

	return nil
}

//GetValidatorUsers - returns list of validators
func (ac *AutonityContractGenesis) GetValidatorUsers() []User {
	var users []User
	for i := range ac.Users {
		if ac.Users[i].Type == UserValidator {
			users = append(users, ac.Users[i])
		}
	}
	return users
}

//GetStakeHolderUsers - returns list of stakeholders
func (ac *AutonityContractGenesis) GetStakeHolderUsers() []User {
	var users []User
	for i := range ac.Users {
		if ac.Users[i].Type == UserStakeHolder { //TODO: Validators are also stakeholders -> Fix this!
			users = append(users, ac.Users[i])
		}
	}
	return users
}

func EnodeToAddress(n *enode.Node) common.Address {
	addrByte := n.ID().Bytes()[12:]
	return common.BytesToAddress(addrByte)
}
