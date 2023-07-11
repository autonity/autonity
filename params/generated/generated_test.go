package generated

import (
	"encoding/json"
	"github.com/autonity/autonity/accounts/abi"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestJson(t *testing.T) {
	data, err := json.Marshal(&AutonityUpgradeTestAbi)
	require.NoError(t, err)
	var Abi abi.ABI
	err = json.Unmarshal(data, &Abi)
	require.NoError(t, err)
}
