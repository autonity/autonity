package autonitytests

import (
	"math/big"
	"testing"

	"github.com/autonity/autonity/autonity/tests"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/require"
)

func TestCirculatingSupply(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("total supply and circulating supply are same", setup, func(r *tests.Runner) {
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, totalSupply, circulatingSupply)
	})

	tests.RunWithSetup("minting increases both circulating and total supply", setup, func(r *tests.Runner) {
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)

		amount := big.NewInt(100)
		r.NoError(
			r.Autonity.Mint(r.Operator, tests.User, amount),
		)
		newTotalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		newCirculatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(totalSupply, amount), newTotalSupply)
		require.Equal(r.T, new(big.Int).Add(circulatingSupply, amount), newCirculatingSupply)
	})

	tests.RunWithSetup("burning decreases both circulating and total supply", setup, func(r *tests.Runner) {
		mintAmount := big.NewInt(100)
		r.NoError(
			r.Autonity.Mint(r.Operator, tests.User, mintAmount),
		)
		totalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		circulatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		burnAmount := big.NewInt(60)
		r.NoError(
			r.Autonity.Burn(r.Operator, tests.User, burnAmount),
		)
		newTotalSupply, _, err := r.Autonity.TotalSupply(nil)
		require.NoError(r.T, err)
		newCirculatingSupply, _, err := r.Autonity.CirculatingSupply(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Sub(totalSupply, burnAmount), newTotalSupply)
		require.Equal(r.T, new(big.Int).Sub(circulatingSupply, burnAmount), newCirculatingSupply)
	})

}

func TestDuplicateOracleAddress(t *testing.T) {
	setup := func() *tests.Runner {
		return tests.Setup(t, nil)
	}

	tests.RunWithSetup("no duplicate oracle address in autonity deployment", setup, func(r *tests.Runner) {
		config, _, err := r.Autonity.Config(nil)
		require.NoError(r.T, err)
		validators := r.Committee.Validators
		for i := 0; i < len(validators); i++ {
			// there is a check for empty address in Prepare(), so there is no check for this on the contract side
			validators[i].OracleAddress = common.Address{}
		}
		config.ContractVersion = common.Big0
		_, _, _, err = r.DeployAutonity(
			nil,
			validators,
			config,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: oracle server exists", err.Error())
	})

	tests.RunWithSetup("no duplicate oracle address in validator registration", setup, func(r *tests.Runner) {
		// register validator with genesis validator oracle address
		validator, signature, nodeKey, blsKey, err := tests.RandomValidator()
		require.NoError(r.T, err)
		duplicateOracleKey, err := crypto.HexToECDSA(params.TestNodeKeys[0])
		require.NoError(r.T, err)
		pop, err := crypto.AutonityPOPProof(nodeKey, duplicateOracleKey, validator.Treasury.Hex(), blsKey)
		require.NoError(r.T, err)
		validator.OracleAddress = crypto.PubkeyToAddress(duplicateOracleKey.PublicKey)
		_, err = r.Autonity.RegisterValidator(
			tests.FromSender(validator.Treasury, nil),
			validator.Enode,
			validator.OracleAddress,
			validator.ConsensusKey,
			pop,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: oracle server exists", err.Error())

		validator.OracleAddress = validator.Treasury
		r.NoError(
			r.Autonity.RegisterValidator(
				tests.FromSender(validator.Treasury, nil),
				validator.Enode,
				validator.OracleAddress,
				validator.ConsensusKey,
				signature,
			),
		)

		// register validator with existing validator oracle address
		newValidator, signature, newNodeKey, blsKey, err := tests.RandomValidator()
		require.NoError(r.T, err)
		require.NotEqual(r.T, nodeKey, newNodeKey)
		pop, err = crypto.AutonityPOPProof(newNodeKey, nodeKey, newValidator.Treasury.Hex(), blsKey)
		require.NoError(r.T, err)
		newValidator.OracleAddress = validator.OracleAddress
		_, err = r.Autonity.RegisterValidator(
			tests.FromSender(newValidator.Treasury, nil),
			newValidator.Enode,
			newValidator.OracleAddress,
			newValidator.ConsensusKey,
			pop,
		)
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: oracle server exists", err.Error())

		newValidator.OracleAddress = newValidator.Treasury
		r.NoError(
			r.Autonity.RegisterValidator(
				tests.FromSender(newValidator.Treasury, nil),
				newValidator.Enode,
				newValidator.OracleAddress,
				newValidator.ConsensusKey,
				signature,
			),
		)
	})
}
