// Package gas provides gas calculation and configuration for PDK contracts.
package gas

import (
	"fmt"
	"math/big"

	"github.com/autonity/autonity/accounts/abi"
)

// Calculator is a function that calculates dynamic gas cost based on input parameters.
type Calculator func(params []byte) uint64

// MethodGas defines the gas cost configuration for a specific method.
type MethodGas struct {
	Base       uint64
	Calculator Calculator
}

// Config holds the gas configuration for all methods of a contract.
type Config struct {
	methods    map[string]MethodGas
	selectors  map[[4]byte]MethodGas
	defaultGas uint64
}

// NewConfig creates a new gas configuration with a default gas cost.
func NewConfig(defaultGas uint64) *Config {
	return &Config{
		methods:    make(map[string]MethodGas),
		selectors:  make(map[[4]byte]MethodGas),
		defaultGas: defaultGas,
	}
}

// SetMethodGas sets the gas configuration for a specific method by name.
func (c *Config) SetMethodGas(methodName string, gas MethodGas) {
	c.methods[methodName] = gas
}

// Finalize resolves method names to selectors using the contract ABI.
func (c *Config) Finalize(contractABI abi.ABI) error {
	for methodName, methodGas := range c.methods {
		method, ok := contractABI.Methods[methodName]
		if !ok {
			return fmt.Errorf("gas config: method %s not found in ABI", methodName)
		}

		var selector [4]byte
		copy(selector[:], method.ID[:4])
		c.selectors[selector] = methodGas
	}
	return nil
}

// GetGas calculates the gas cost for a given input.
func (c *Config) GetGas(input []byte) uint64 {
	if len(input) < 4 {
		return c.defaultGas
	}

	var selector [4]byte
	copy(selector[:], input[:4])

	methodGas, ok := c.selectors[selector]
	if !ok {
		return c.defaultGas
	}

	totalGas := methodGas.Base

	if methodGas.Calculator != nil {
		params := input[4:]
		dynamicGas := methodGas.Calculator(params)
		totalGas += dynamicGas
	}

	return totalGas
}

// GetDefaultGas returns the default gas cost.
func (c *Config) GetDefaultGas() uint64 {
	return c.defaultGas
}

// ArrayLengthCalculator charges based on array length from ABI encoding
func ArrayLengthCalculator(costPerItem uint64) Calculator {
	return func(params []byte) uint64 {
		if len(params) < 32 {
			return 0
		}
		length := new(big.Int).SetBytes(params[:32]).Uint64()
		return length * costPerItem
	}
}

// FixedCalculator returns a constant additional gas cost
func FixedCalculator(additionalGas uint64) Calculator {
	return func(_ []byte) uint64 {
		return additionalGas
	}
}
