package tests

import (
	"reflect"

	"github.com/autonity/autonity/core/types"
)

type isEvent interface {
	SetRaw(types.Log)
	GetRaw() types.Log
}

func emitsEvent[T isEvent](logs []*types.Log, parser func(log types.Log) (T, error), expected T) bool {
	for _, log := range logs {
		actual, err := parser(*log)
		if err != nil {
			continue
		}
		// we want to ignore the Raw field in the comparison
		expected.SetRaw(*log)
		if reflect.DeepEqual(actual, expected) {
			return true
		}
	}
	return false
}
