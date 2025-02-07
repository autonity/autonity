package types

import (
	"math/big"
	"reflect"
	"testing"

	"github.com/autonity/autonity/common"
	"github.com/stretchr/testify/require"
)

var one = reflect.ValueOf(common.Big1)
var two = reflect.ValueOf(common.Big2)

func countFields(value reflect.Value) int {
	numFields := value.NumField()
	totalFields := 0

	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		// recursively traverse inner struct
		if field.Kind() == reflect.Struct {
			totalFields += countFields(field)
			continue
		}
		if field.Type() == reflect.TypeOf(new(big.Int)) {
			totalFields++
			continue
		}
		panic("unexpected field type: " + field.Type().String())
	}
	return totalFields
}

func setField(value reflect.Value, index int, count *int, targetValue reflect.Value) {
	numFields := value.NumField()

	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		if field.Type() == reflect.TypeOf(new(big.Int)) {
			if index == *count {
				field.Set(targetValue)
			}
			*count++
			continue
		}
		// recursively traverse inner struct
		if field.Kind() == reflect.Struct {
			setField(field, index, count, targetValue)
			continue
		}
		// crash otherwise
		panic("unexpected field type: " + field.Type().String())
	}
}

// set all fields of the config via reflection
// handles nested structs, crashes if leaves fields are not *big.Int
func setAllFields(config *ContractsConfig) {
	pointerToConfig := reflect.ValueOf(config)
	configValue := pointerToConfig.Elem()

	innerSetAllFields(configValue)

}

func innerSetAllFields(value reflect.Value) {
	numFields := value.NumField()

	for i := 0; i < numFields; i++ {
		field := value.Field(i)
		// set *big.Int to one
		if field.Type() == reflect.TypeOf(new(big.Int)) {
			field.Set(one)
			continue
		}
		// recursively traverse inner struct
		if field.Kind() == reflect.Struct {
			innerSetAllFields(field)
			continue
		}
		// crash otherwise
		panic("unexpected field type: " + field.Type().String())
	}
}

func TestConfigCopy(t *testing.T) {
	config := &ContractsConfig{}

	setAllFields(config)

	configCopy := config.Copy()

	require.True(t, reflect.DeepEqual(config, configCopy))
}

func TestIsEqual(t *testing.T) {
	t.Run("works correctly", func(t *testing.T) {
		config := &ContractsConfig{
			BlockPeriod: common.Big1,
			EpochPeriod: common.Big256,
			GasLimit:    common.Big5,
		}
		config2 := &ContractsConfig{
			BlockPeriod: new(big.Int).Set(config.BlockPeriod),
			EpochPeriod: new(big.Int).Set(config.EpochPeriod),
			GasLimit:    new(big.Int).Set(config.GasLimit),
		}

		require.True(t, config.Equal(config))
		require.True(t, config.Equal(config2))
		config2.BlockPeriod = common.Big4
		require.False(t, config.Equal(config2))
	})
	t.Run("no field was forgotten", func(t *testing.T) {
		config := &ContractsConfig{}
		config2 := &ContractsConfig{}

		pointerToConfig := reflect.ValueOf(config)
		configValue := pointerToConfig.Elem()
		numFields := countFields(configValue)

		// set both configs to all 1s
		for i := 0; i < numFields; i++ {
			count := 0
			setField(configValue, i, &count, one)
		}
		*config2 = *config
		require.True(t, config.Equal(config2))
		t.Log("initial configs")
		t.Log(config)
		t.Log(config2)

		// change one by one the fields of one config to two, to check whether a field
		// is forgotten in the equality function
		for i := 0; i < numFields; i++ {
			count := 0
			setField(configValue, i, &count, two)
			require.False(t, config.Equal(config2))
			t.Log("before reassignment")
			t.Log(config)
			t.Log(config2)
			*config2 = *config
			t.Log("after reassignment")
			t.Log(config)
			t.Log(config2)
			t.Log("===========")
			require.True(t, config.Equal(config2))
		}
	})
}
