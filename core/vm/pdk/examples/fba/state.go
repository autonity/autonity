package fba

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

// Contract holds the state variables using the SDK
type Contract struct {
	// Access Control Roles (Simple single-admin for demo, can be expanded)
	Admin storage.Var[common.Address]

	// Authorized Submitters (Mapping address -> bool)
	Submitters storage.Map[common.Address, storage.Var[bool]]

	// Configuration Addresses
	ClearingAddr      storage.Var[common.Address]
	MarginAccountAddr storage.Var[common.Address]
	ProductRegistry   storage.Var[common.Address]
	CollateralToken   storage.Var[common.Address]
}

// NewContract initializes the contract state structure
func NewContract() *Contract {
	return &Contract{}
}
