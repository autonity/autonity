package tests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
	"github.com/stretchr/testify/require"
)

func makeCommit(salt *big.Int, sender common.Address, reports []*big.Int) *big.Int {
	buffer := make([]byte, 0)
	for i := range reports {
		buffer = append(buffer, common.LeftPadBytes(reports[i].Bytes(), 32)...)
	}
	buffer = append(buffer, common.LeftPadBytes(salt.Bytes(), 32)...)
	buffer = append(buffer, sender.Bytes()...)
	return new(big.Int).SetBytes(crypto.Keccak256(buffer))
}

/*
func TestSimpleVote(t *testing.T) {
	r := setup(t, nil)
	symbols, _, _ := r.oracle.GetSymbols(nil)
	votePeriod, _, _ := r.oracle.GetVotePeriod(nil)
	tests := []struct {
		votes    [][]*big.Int
		expected []*big.Int
	}{
		{
			votes: [][]*big.Int{
				{big.NewInt(1), big.NewInt(90009), big.NewInt(90), big.NewInt(100)},
				{big.NewInt(10093), big.NewInt(988), big.NewInt(90), big.NewInt(99188129399)},
				{big.NewInt(457645765), big.NewInt(237492837498), big.NewInt(18), big.NewInt(100)},
				{big.NewInt(1), big.NewInt(90009), big.NewInt(90), big.NewInt(100)},
			},
			expected: []*big.Int{big.NewInt(95)},
		},
	}

	test := func(r *runner, n int) {
		var (
			committedVotes = make([][]*big.Int, len(r.committee.validators))
			currentVotes   = make([][]*big.Int, len(r.committee.validators))
			rounds         = len(tests[n].votes)
		)
		for round := 0; round <= rounds; round++ {
			for i, validator := range r.committee.validators {
				if currentVotes[i] == nil {
					currentVotes[i] = make([]*big.Int, len(symbols))
				}
				if committedVotes[i] == nil {
					committedVotes[i] = make([]*big.Int, len(symbols))
				}
				for s := range symbols {
					if round == 0 {
						committedVotes[i][s] = common.Big0
					}
					if round == rounds {
						currentVotes[i][s] = common.Big0
					} else {
						currentVotes[i][s] = tests[n].votes[round][i]
					}
				}
				_, err := r.oracle.Vote(
					&runOptions{origin: validator.OracleAddress},
					makeCommit(common.Big1, validator.OracleAddress, currentVotes[i]),
					committedVotes[i],
					common.Big1,
				)
				require.NoError(r.t, err)
			}
			r.waitNBlocks(int(votePeriod.Int64()))
			data, _, err := r.oracle.LatestRoundData(nil, symbols[0])
			require.NoError(t, err)
			fmt.Println(data)
			committedVotes, currentVotes = currentVotes, committedVotes
		}
	}

	for n := range tests {
		r.run(fmt.Sprintf("test vote - %d", n), func(r *runner) {
			test(r, n)
		})
	}
}

*/

func TestReportPacking(t *testing.T) {
	// Define your Solidity-like function arguments
	addressType, _ := abi.NewType("address", "", nil)
	saltType, _ := abi.NewType("uint256", "", nil)
	reportType, _ := abi.NewType("tuple[]", "struct Overloader.F", []abi.ArgumentMarshaling{
		{Name: "price", Type: "uint120"},
		{Name: "confidence", Type: "uint8"}})
	args := abi.Arguments{
		{Type: reportType},  // Solidity uint256
		{Type: saltType},    // Solidity address
		{Type: addressType}, // Solidity bool
	}

	// Prepare values to encode
	reports := []IOracleReport{
		{
			Price:      big.NewInt(121212),
			Confidence: 8,
		},
		{
			Price:      big.NewInt(88282828),
			Confidence: 34,
		},
	}
	// Example uint256
	addressVal := common.HexToAddress("0x11") // Example address
	saltVal := big.NewInt(99999)              // Example bool

	// Pack values to encode them
	packed, err := args.Pack(reports, saltVal, addressVal)
	if err != nil {
		t.Fatalf("Failed to pack values: %v", err)
	}
	fmt.Printf("Encoded data: %x\n", packed)
	res, _ := args.Unpack(packed)
	fmt.Println(res)
}

func TestVotingPeriodUpdate(t *testing.T) {
	setup := func() *Runner {
		r := Setup(t, nil)
		// set big values for testing
		r.NoError(
			r.Autonity.SetEpochPeriod(r.Operator, big.NewInt(100)),
		)
		r.WaitNextEpoch()
		r.NoError(
			r.Oracle.SetVotePeriod(r.Operator, big.NewInt(30)),
		)
		return r
	}

	RunWithSetup("voting period cannot be too big", setup, func(r *Runner) {
		epochPeriod, _, err := r.Autonity.GetEpochPeriod(nil)
		require.NoError(r.T, err)
		require.True(r.T, epochPeriod.Cmp(common.Big2) >= 0, "cannot test")
		maxVotingPeriod := new(big.Int).Div(epochPeriod, big.NewInt(2))
		_, err = r.Oracle.SetVotePeriod(r.Operator, new(big.Int).Add(maxVotingPeriod, common.Big1))
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: vote period is too big", err.Error())
		r.NoError(
			r.Oracle.SetVotePeriod(r.Operator, maxVotingPeriod),
		)

		// change the parity
		epochPeriod = new(big.Int).Add(epochPeriod, common.Big1)
		r.NoError(
			r.Autonity.SetEpochPeriod(r.Operator, epochPeriod),
		)
		r.WaitNextEpoch()
		newEpochPeriod, _, err := r.Autonity.GetEpochPeriod(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, epochPeriod, newEpochPeriod)
		maxVotingPeriod = new(big.Int).Div(epochPeriod, big.NewInt(2))
		_, err = r.Oracle.SetVotePeriod(r.Operator, new(big.Int).Add(maxVotingPeriod, common.Big1))
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: vote period is too big", err.Error())
		r.NoError(
			r.Oracle.SetVotePeriod(r.Operator, maxVotingPeriod),
		)
	})

	RunWithSetup("epoch period cannot be too small", setup, func(r *Runner) {
		votingPeriod, _, err := r.Oracle.GetVotePeriod(nil)
		require.NoError(r.T, err)
		require.True(r.T, votingPeriod.Cmp(common.Big1) >= 0, "cannot test")
		minEpochPeriod := new(big.Int).Mul(votingPeriod, big.NewInt(2))
		_, err = r.Autonity.SetEpochPeriod(r.Operator, new(big.Int).Sub(minEpochPeriod, common.Big1))
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: epoch period is too small", err.Error())
		r.NoError(
			r.Autonity.SetEpochPeriod(r.Operator, minEpochPeriod),
		)
	})

	RunWithSetup("voting period respects both current and new epoch period", setup, func(r *Runner) {
		epochPeriod, _, err := r.Autonity.GetEpochPeriod(nil)
		require.NoError(r.T, err)
		require.True(r.T, epochPeriod.Cmp(common.Big2) >= 0, "cannot test")

		// set new epoch period bigger
		newEpochPeriod := new(big.Int).Add(epochPeriod, big.NewInt(10))
		maxVotingPeriod := new(big.Int).Div(epochPeriod, big.NewInt(2))
		r.NoError(
			r.Autonity.SetEpochPeriod(r.Operator, newEpochPeriod),
		)
		currentEpochPeriod, _, err := r.Autonity.GetCurrentEpochPeriod(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, epochPeriod, currentEpochPeriod)
		_, err = r.Oracle.SetVotePeriod(r.Operator, new(big.Int).Add(maxVotingPeriod, common.Big1))
		require.Equal(r.T, "execution reverted: vote period is too big", err.Error())
		r.NoError(
			r.Oracle.SetVotePeriod(r.Operator, maxVotingPeriod),
		)

		r.WaitNextEpoch()
		epochPeriod = newEpochPeriod
		currentEpochPeriod, _, err = r.Autonity.GetEpochPeriod(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, epochPeriod, currentEpochPeriod)

		// set new epoch period smaller
		newEpochPeriod = new(big.Int).Sub(epochPeriod, big.NewInt(10))
		maxVotingPeriod = new(big.Int).Div(newEpochPeriod, big.NewInt(2))
		r.NoError(
			r.Autonity.SetEpochPeriod(r.Operator, newEpochPeriod),
		)
		currentEpochPeriod, _, err = r.Autonity.GetCurrentEpochPeriod(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, epochPeriod, currentEpochPeriod)
		_, err = r.Oracle.SetVotePeriod(r.Operator, new(big.Int).Add(maxVotingPeriod, common.Big1))
		require.Error(r.T, err)
		require.Equal(r.T, "execution reverted: vote period is too big", err.Error())
		r.NoError(
			r.Oracle.SetVotePeriod(r.Operator, maxVotingPeriod),
		)
	})
}

func TestVotersUpdate(t *testing.T) {

	newVoterCheck := func(r *Runner, voters map[common.Address]struct{}, isVoter bool) {
		newVoters, _, err := r.Oracle.GetNewVoters(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, len(r.Committee.Validators), len(newVoters))

		newVoterSet := make(map[common.Address]struct{})
		for _, v := range newVoters {
			newVoterSet[v] = struct{}{}
		}

		for _, v := range r.Committee.Validators {
			_, ok := newVoterSet[v.OracleAddress]
			require.True(r.T, ok)
			voterInfo, _, err := r.Oracle.VoterInfo(nil, v.OracleAddress)
			require.NoError(r.T, err)
			require.Equal(r.T, v.Treasury, voterInfo.Treasury)
			require.Equal(r.T, v.NodeAddress, voterInfo.Validator)
			if _, ok := voters[v.OracleAddress]; ok {
				require.Equal(r.T, true, voterInfo.IsVoter)
			} else {
				require.Equal(r.T, isVoter, voterInfo.IsVoter)
			}
		}
	}

	voterCheck := func(r *Runner, expectedVoters map[common.Address]struct{}) {
		voters, _, err := r.Oracle.GetVoters(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, len(expectedVoters), len(voters))
		for _, v := range voters {
			_, ok := expectedVoters[v]
			require.True(r.T, ok)
			voterInfo, _, err := r.Oracle.VoterInfo(nil, v)
			require.NoError(r.T, err)
			validator, _, err := r.Autonity.GetValidator(nil, voterInfo.Validator)
			require.NoError(r.T, err)
			require.Equal(r.T, validator.Treasury, voterInfo.Treasury)
			require.Equal(r.T, true, voterInfo.IsVoter)
		}
	}

	getVoters := func(r *Runner) map[common.Address]struct{} {
		voters := make(map[common.Address]struct{})
		for _, v := range r.Committee.Validators {
			voters[v.OracleAddress] = struct{}{}
		}
		return voters
	}

	setup := func() *Runner {
		r := Setup(t, SetInflationReserveZero)
		require.True(r.T, len(r.Committee.Validators) >= 4, "cannot test")
		// all validators
		allValidators := r.Committee.Validators
		voterCount := 2
		for i := voterCount; i < len(allValidators); i++ {
			require.Equal(r.T, allValidators[i].SelfBondedStake, allValidators[i].BondedStake)
			r.NoError(
				r.Autonity.Unbond(
					FromSender(allValidators[i].Treasury, nil),
					allValidators[i].NodeAddress,
					allValidators[i].SelfBondedStake,
				),
			)
		}
		// so that voting period is not a factor of epoch period
		r.NoError(
			r.Autonity.SetEpochPeriod(r.Operator, big.NewInt(127)),
		)
		r.NoError(
			r.Oracle.SetVotePeriod(r.Operator, big.NewInt(11)),
		)
		r.WaitNextEpoch()
		require.Equal(r.T, voterCount, len(r.Committee.Validators))

		// check voter info update
		r.WaitNextEpoch()
		voters := getVoters(r)
		voterCheck(r, voters)
		newVoterCheck(r, voters, true)

		for {
			released, _, err := r.Autonity.IsUnbondingReleased(nil, common.Big0)
			require.NoError(r.T, err)
			if released {
				break
			}
			r.WaitNextEpoch()
		}
		return r
	}

	progressRound := func(r *Runner, round *big.Int) {
		for {
			r.WaitNBlocks(1)
			newRound, _, err := r.Oracle.GetRound(nil)
			require.NoError(r.T, err)
			if newRound.Cmp(round) == 1 {
				require.Equal(r.T, new(big.Int).Add(round, common.Big1), newRound, "cannot test") // newRound == round+1
				break
			}
		}
	}

	checkVoterUpdate := func(r *Runner, oldVoters map[common.Address]struct{}) {
		voterCheck(r, oldVoters)
		newVoterCheck(r, oldVoters, false)

		// progress round
		round, _, err := r.Oracle.GetRound(nil)
		require.NoError(r.T, err)
		progressRound(r, round)
		// new voters should get access for voting, but voters array not updated yet
		round = new(big.Int).Add(round, common.Big1)
		voterCheck(r, oldVoters)
		newVoterCheck(r, oldVoters, true)

		// progress round
		progressRound(r, round)
		// voters array should be updated
		voters := getVoters(r)
		newVoterCheck(r, voters, true)
		voterCheck(r, voters)
	}

	getCommitteeSet := func(r *Runner) map[common.Address]struct{} {
		committeeSet := make(map[common.Address]struct{})
		for _, c := range r.Committee.Validators {
			committeeSet[c.NodeAddress] = struct{}{}
		}
		return committeeSet
	}

	getAllValidators := func(r *Runner) []common.Address {
		allValidators, _, err := r.Autonity.GetValidators(nil)
		require.NoError(r.T, err)
		return allValidators
	}

	checkCommittee := func(r *Runner, expectedCommittee map[common.Address]struct{}) {
		committeeSet := getCommitteeSet(r)
		require.Equal(r.T, len(expectedCommittee), len(committeeSet))
		for c := range committeeSet {
			_, ok := expectedCommittee[c]
			require.True(r.T, ok)
		}
	}

	addToCommittee := func(r *Runner, newVoter int) map[common.Address]struct{} {
		committeeSet := getCommitteeSet(r)
		allValidators := getAllValidators(r)
		newCommitteeSet := make(map[common.Address]struct{})
		for i := 0; newVoter > 0 && i < len(allValidators); i++ {
			if _, ok := committeeSet[allValidators[i]]; ok {
				continue
			}
			newVoter--
			validator, _, err := r.Autonity.GetValidator(nil, allValidators[i])
			require.NoError(r.T, err)
			newtonBalance := r.GetNewtonBalanceOf(validator.Treasury)
			r.NoError(
				r.Autonity.Bond(
					FromSender(validator.Treasury, nil),
					validator.NodeAddress,
					newtonBalance,
				),
			)
			newCommitteeSet[validator.NodeAddress] = struct{}{}
		}
		require.True(r.T, newVoter == 0, "cannot test")
		return newCommitteeSet
	}

	removeFromCommittee := func(r *Runner, removeVoter int) map[common.Address]struct{} {
		require.True(r.T, removeVoter <= len(r.Committee.Validators))
		removed := make(map[common.Address]struct{})
		for i := 0; i < removeVoter; i++ {
			validator := r.Committee.Validators[i]
			r.NoError(
				r.Autonity.Unbond(
					FromSender(validator.Treasury, nil),
					validator.NodeAddress,
					validator.SelfBondedStake,
				),
			)
			removed[validator.NodeAddress] = struct{}{}
		}
		return removed
	}

	RunWithSetup("new voters are updated properly (new voters and old voters have empty intersection set)", setup, func(r *Runner) {
		newCommitteeSet := addToCommittee(r, 2)
		removeFromCommittee(r, len(r.Committee.Validators))
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkCommittee(r, newCommitteeSet)
		checkVoterUpdate(r, oldVoters)
	})

	RunWithSetup("new voters are updated properly (new voters and old voters are same)", setup, func(r *Runner) {
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkVoterUpdate(r, oldVoters)
	})

	RunWithSetup("new voters are updated properly (new voters and old voters have non-empty intersection set)", setup, func(r *Runner) {
		require.True(r.T, len(r.Committee.Validators) > 1, "cannot test")
		newCommitteeSet := addToCommittee(r, 1)
		removed := removeFromCommittee(r, 1)
		for _, c := range r.Committee.Validators {
			if _, ok := removed[c.NodeAddress]; !ok {
				newCommitteeSet[c.NodeAddress] = struct{}{}
			}
		}
		require.True(r.T, len(newCommitteeSet) > 1, "cannot test")
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkCommittee(r, newCommitteeSet)
		checkVoterUpdate(r, oldVoters)
	})

	RunWithSetup("new voters are updated properly (new voters set is a subset of old voters set)", setup, func(r *Runner) {
		require.True(r.T, len(r.Committee.Validators) > 1, "cannot test")
		removed := removeFromCommittee(r, 1)
		newCommitteeSet := make(map[common.Address]struct{})
		for _, c := range r.Committee.Validators {
			if _, ok := removed[c.NodeAddress]; !ok {
				newCommitteeSet[c.NodeAddress] = struct{}{}
			}
		}
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkCommittee(r, newCommitteeSet)
		checkVoterUpdate(r, oldVoters)
	})

	RunWithSetup("new voters are updated properly (new voters set is a superset of old voters set)", setup, func(r *Runner) {
		newCommitteeSet := addToCommittee(r, 1)
		for _, c := range r.Committee.Validators {
			newCommitteeSet[c.NodeAddress] = struct{}{}
		}
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkCommittee(r, newCommitteeSet)
		checkVoterUpdate(r, oldVoters)
	})
}
