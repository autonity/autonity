package core

import (
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/p2p/enode"
	"reflect"
)

const (
	//participant: Authorized to operate a full node, is able to join the network, is not authorized to own stake.
	UserParticipant UserType = "participant"
	//member: Authorized to operate a full node, is able to join the network, authorized to own stake.
	UserStakeHolder UserType = "stakeholder"
	//validator: Authorized to operate a full node, is able to join the network, authorized to own stake, participate in consensus.
	UserValidator UserType = "validator"
)

type UserType string

func (ut UserType) IsValid() bool {
	if ut == UserStakeHolder || ut == UserValidator || ut == UserParticipant {
		return true
	}
	return false
}

// Autonity contract config. It'is used for deployment.
type AutonityContract struct {
	// Address of the validator who deploys contract stored in bytecode
	Deployer common.Address `json:"deployer" toml:",omitempty"`
	// Bytecode of validators contract // would like this type to be []byte but the unmarshalling is not working
	Bytecode string `json:"bytecode" toml:",omitempty"`
	// Json ABI of the contract
	ABI         string         `json:"abi "toml:",omitempty"`
	MinGasPrice uint64         `json:"minGasPrice" toml:",omitempty"`
	Operator    common.Address `json:"operator" toml:",omitempty"`
	Users       []User         `json:"users" "toml:",omitempty"`
}

func (ac *AutonityContract) Validate() error {
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
			return err
		}
	}
	return nil
}

//User - is used to put predefined accounts to genesis
type User struct {
	Address common.Address
	Enode   string
	Type    UserType
	Stake   uint64
}

func (u *User) Validate() error {
	if !u.Type.IsValid() {
		return errors.New("incorrect user type")
	}

	if reflect.DeepEqual(u.Address, common.Address{}) && len(u.Enode) == 0 {
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
		if err != nil {
			return fmt.Errorf("fail to parse enode for account %v, error:%v", u.Address, err)
		}

		addrFromEnode := EnodeToAddress(n)
		if reflect.DeepEqual(u.Address, common.Address{}) {
			u.Address = addrFromEnode
		} else if !reflect.DeepEqual(u.Address, addrFromEnode) {
			return errors.New("if both user.enode and user.address are defined, then the derived address from user.enode must be equal to user.address")
		}
	}

	return nil
}

//GetParticipantUsers - returns list of participants
func (ac *AutonityContract) GetValidatorUsers() []User {
	var users []User
	for i := range ac.Users {
		if ac.Users[i].Type == UserValidator {
			users = append(users, ac.Users[i])
		}
	}
	return users
}

func EnodeToAddress(n *enode.Node) common.Address {
	addrByte := n.ID().Bytes()[12:]
	return common.BytesToAddress(addrByte)
}
