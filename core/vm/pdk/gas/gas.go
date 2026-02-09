package gas

import (
	"fmt"
	"math/big"

	"github.com/autonity/autonity/accounts/abi"
)

type Calculator func(params []byte) uint64

type MethodGas struct {
	Base       uint64
	Calculator Calculator
}

type Config struct {
	methods    map[string]MethodGas
	selectors  map[[4]byte]MethodGas
	defaultGas uint64
}

func NewConfig(defaultGas uint64) *Config {
	return &Config{
		methods:    make(map[string]MethodGas),
		selectors:  make(map[[4]byte]MethodGas),
		defaultGas: defaultGas,
	}
}

func (c *Config) SetMethodGas(methodName string, gas MethodGas) {
	c.methods[methodName] = gas
}

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
