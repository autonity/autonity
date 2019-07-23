package core

import (
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/p2p/enode"
	"reflect"
)

const (
	UserMember             UserType = "member"
	UserValidator          UserType = "validator"
	UserGovernanceOperator UserType = "GovernanceOperator"
)

type UserType string

func (ut UserType) IsValid() bool {
	if ut == UserMember || ut == UserGovernanceOperator || ut == UserValidator {
		return true
	}
	return false
}

// Autonity contract config. It'is used for deployment.
type AutonityContract struct {
	// Address of the validator who deploys contract stored in bytecode
	Deployer common.Address `json:"deployer" toml:",omitempty"`
	// Bytecode of validators contract // would like this type to be []byte but the unmarshalling is not working
	Bytecode    string `json:"bytecode" toml:",omitempty"`
	MinGasPrice uint64 `json:"minGasPrice" toml:",omitempty"`
	// Json ABI of the contract
	ABI   string `json:"abi "toml:",omitempty"`
	Users []User `json:"users" "toml:",omitempty"`
}

func (ac *AutonityContract) Validate() error {
	if len(ac.Bytecode) == 0 || len(ac.ABI) == 0 {
		return errors.New("autonity contract is empty")
	}
	if reflect.DeepEqual(ac.Deployer, common.Address{}) {
		return errors.New("deployer is empty")
	}
	numOfGovernanceOperator := 0
	for i := range ac.Users {
		if err := ac.Users[i].Validate(); err != nil {
			return err
		}
		if ac.Users[i].Type == UserGovernanceOperator {
			numOfGovernanceOperator++
		}
	}
	if numOfGovernanceOperator != 1 {
		return errors.New("incorrect number of GovernanceOperators")
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

	n, err := enode.ParseV4(u.Enode)
	if err != nil {
		return fmt.Errorf("Fail to parse enode for account %v, error:%v", u.Address, err)
	}

	if reflect.DeepEqual(u.Address, common.Address{}) {
		u.Address = EnodeToAddress(n)
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

//GetGovernanceOperator - returns GovernanceOperator
func (ac *AutonityContract) GetGovernanceOperator() User {
	for i := range ac.Users {
		if ac.Users[i].Type == UserGovernanceOperator {
			return ac.Users[i]
		}
	}
	return User{}
}

func EnodeToAddress(n *enode.Node) common.Address {
	addrByte := n.ID().Bytes()[12:]
	return common.BytesToAddress(addrByte)
}
