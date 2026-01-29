package fba

import (
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/vm/pdk/storage"
)

// FBAContract holds the state variables using the SDK
type FBAContract struct {
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

// NewFBAContract initializes the contract state structure
func NewFBAContract() *FBAContract {
	return &FBAContract{}
}
