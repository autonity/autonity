package core

import (
	"errors"
	"fmt"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/p2p/enode"
	"reflect"
)

const (
	UserParticipant        UserType = "participant"
	UserMember             UserType = "member"
	UserGovernanceOperator UserType = "GovernanceOperator"
)

type UserType string

func (ut UserType) IsValid() bool {
	if ut == UserMember || ut == UserParticipant || ut == UserGovernanceOperator {
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
	ABI   string `json:"abi "toml:",omitempty"`
	Users []User `json:"users"`
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

func (u User) Validate() error {
	if !u.Type.IsValid() {
		return errors.New("incorrect user type")
	}
	if u.Type == UserMember && u.Stake > 0 {
		return fmt.Errorf("user %v has a stake, but shouldn't", u.Address)
	}
	if reflect.DeepEqual(u.Address, common.Address{}) {
		return errors.New("account is empty")
	}
	if _, err := enode.ParseV4(u.Enode); err != nil {
		return fmt.Errorf("Fail to parse enode for account %v, error:%v", u.Address, err)
	}
	return nil
}

//GetParticipantUsers - returns list of participants
func (ac AutonityContract) GetParticipantUsers() []User {
	var users []User
	for i := range ac.Users {
		if ac.Users[i].Type == UserParticipant {
			users = append(users, ac.Users[i])
		}
	}
	return users
}

//GetGovernanceOperator - returns GovernanceOperator
func (ac AutonityContract) GetGovernanceOperator() User {
	for i := range ac.Users {
		if ac.Users[i].Type == UserGovernanceOperator {
			return ac.Users[i]
		}
	}
	return User{}
}
