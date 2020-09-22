package params

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/clearmatics/autonity/crypto"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/acdefault"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
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
	// Bytecode of validators contract
	// would like this type to be []byte but the unmarshalling is not working
	Bytecode string `json:"bytecode,omitempty" toml:",omitempty"`
	// Json ABI of the contract
	ABI         string         `json:"abi,omitempty" toml:",omitempty"`
	MinGasPrice uint64         `json:"minGasPrice" toml:",omitempty"`
	Operator    common.Address `json:"operator" toml:",omitempty"`
	Users       []User         `json:"users" toml:",omitempty"`
}

// Prepare prepares the AutonityContractGenesis by filling in missing fields.
// It returns an error if the configuration is invalid.
func (ac *AutonityContractGenesis) Prepare() error {

	if len(ac.Bytecode) == 0 && len(ac.ABI) > 0 ||
		len(ac.Bytecode) > 0 && len(ac.ABI) == 0 {
		return errors.New("it is an error to set only of autonity contract abi or bytecode")
	}

	if len(ac.Bytecode) == 0 && len(ac.ABI) == 0 {
		log.Info("Default Validator smart contract set")
		ac.ABI = acdefault.ABI()
		ac.Bytecode = acdefault.Bytecode()
	} else {
		log.Info("User specified Validator smart contract set")
	}
	if reflect.DeepEqual(ac.Operator, common.Address{}) {
		ac.Operator = acdefault.Governance()
	}

	for i := range ac.Users {
		if ac.Users[i].Address == nil {
			addr, err := ac.Users[i].getAddressFromEnode()
			if err != nil {
				return fmt.Errorf("error parsing enode %q, err: %v", ac.Users[i].Enode, err)
			}
			ac.Users[i].Address = &addr
			err = ac.Users[i].Validate()
			if err != nil {
				return err
			}
		}
	}

	if len(ac.GetValidatorUsers()) == 0 {
		return errors.New("validators list is empty")
	}
	return nil
}

//User - is used to put predefined accounts to genesis
type User struct {
	Address *common.Address `json:"address,omitempty"`
	Enode   string          `json:"enode"`
	Type    UserType        `json:"type"`
	Stake   uint64          `json:"stake"`
}

// getAddressFromEnode gets the account address from the user enode.
func (u *User) getAddressFromEnode() (common.Address, error) {
	n, err := enode.ParseV4(u.Enode)
	if err != nil {
		panic(err.Error())
		return common.Address{}, fmt.Errorf("failed to parse enode %q, error:%v", u.Enode, err)
	}
	return crypto.PubkeyToAddress(*n.Pubkey()), nil
}

func (u *User) Validate() error {
	switch {
	case !u.Type.IsValid():
		return errors.New("incorrect user type")

	case u.Type == UserParticipant && u.Stake > 0:
		return errors.New("user.stake must be nil or equal to 0 for users of type participant")

	case u.Type == UserValidator && len(u.Enode) == 0:
		return errors.New("if user.type is validator then user.enode must be defined")

	case u.Address == nil && len(u.Enode) == 0:
		return errors.New("user.enode or user.address must be defined")

	case len(u.Enode) > 0:
		a, err := u.getAddressFromEnode()
		if err != nil {
			return err
		}
		// If address is set check it matches the address from the enode
		if u.Address != nil && *u.Address != a {
			return fmt.Errorf("mismatching address %q and address from enode %q", u.Address.String(), a.String())
		}
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
