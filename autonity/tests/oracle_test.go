package tests

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"
)

func TestSimpleVote(t *testing.T) {
	r := Setup(t, nil)
	symbols, _, _ := r.Oracle.GetSymbols(nil)
	votePeriod, _, _ := r.Oracle.GetVotePeriod(nil)
	tests := []struct {
		votes    [][]IOracleReport
		expected []*big.Int
	}{
		{
			votes: [][]IOracleReport{
				{
					// Happy case scenario: no outliers, 100% confidence
					{big.NewInt(225), 100},
					{big.NewInt(226), 100},
					{big.NewInt(228), 100},
					{big.NewInt(230), 100},
				},
				{
					// no outliers, different confidence
					// median = 12900
					{big.NewInt(13500), 50},
					{big.NewInt(12600), 70},
					{big.NewInt(12800), 100},
					{big.NewInt(13000), 1},
				},
				{
					// one outlier above median
					// median = 13250 | upper threshold = 15369 @ 1.16
					{big.NewInt(13500), 50},
					{big.NewInt(12600), 70},
					{big.NewInt(15370), 100}, // outlier +16%
					{big.NewInt(13000), 1},
				},
				{
					// one outlier below median
					// median = 12800 | lower threshold = 10752
					{big.NewInt(13500), 50},
					{big.NewInt(12600), 70},
					{big.NewInt(10753), 100}, // outlier -16%
					{big.NewInt(13000), 1},
				},
			},
			expected: []*big.Int{
				big.NewInt(95),
				big.NewInt(12895),
			},
		},
	}

	test := func(r *Runner, n int) {
		var (
			committedVotes = make([][]IOracleReport, len(r.Committee.Validators))
			currentVotes   = make([][]IOracleReport, len(r.Committee.Validators))
			rounds         = len(tests[n].votes)
		)
		for round := 0; round <= rounds; round++ {
			for i, validator := range r.Committee.Validators {
				if currentVotes[i] == nil {
					currentVotes[i] = make([]IOracleReport, len(symbols))
				}
				if committedVotes[i] == nil {
					committedVotes[i] = make([]IOracleReport, len(symbols))
				}
				for s := range symbols {
					if round == 0 {
						committedVotes[i][s] = IOracleReport{common.Big0, 0}

					}
					if round == rounds {
						currentVotes[i][s] = IOracleReport{common.Big0, 0}
					} else {
						currentVotes[i][s] = tests[n].votes[round][i]
					}
				}
				_, err := r.Oracle.Vote(
					&runOptions{origin: validator.OracleAddress},
					makeCommit(r.T, common.Big1, validator.OracleAddress, currentVotes[i]),
					committedVotes[i],
					common.Big1,
					87,
				)
				require.NoError(r.T, err)
			}
			r.WaitNBlocks(int(votePeriod.Int64()))
			data, _, err := r.Oracle.LatestRoundData(nil, symbols[0])
			require.NoError(t, err)
			fmt.Println(data)
			committedVotes, currentVotes = currentVotes, committedVotes
		}
	}

	for n := range tests {
		r.Run(fmt.Sprintf("test vote - %d", n), func(r *Runner) {
			test(r, n)
		})
	}
}

// abi.encode(_reports, _salt, msg.sender) follows below encoding schema of the eth ABI specification.
var ReportABIEncodeSchema = []byte("[{\"components\":[{\"internalType\":\"uint120\",\"name\":\"price\",\"type\":\"uint120\"},{\"internalType\":\"uint8\",\"name\":\"confidence\",\"type\":\"uint8\"}],\"internalType\":\"struct Report[]\",\"name\":\"_reports\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"_salt\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}]")

func genReports(n int, price ...int) []IOracleReport {
	defaultPrice := 1000
	if len(price) > 0 {
		defaultPrice = price[0]
	}
	var reports []IOracleReport
	for i := 0; i < n; i++ {
		reports = append(reports, IOracleReport{
			Price:      big.NewInt(int64(defaultPrice)),
			Confidence: 100,
		})
	}
	return reports
}

func makeCommit(t *testing.T, salt *big.Int, sender common.Address, reports []IOracleReport) *big.Int {
	var args abi.Arguments
	err := json.Unmarshal(ReportABIEncodeSchema, &args)
	require.NoError(t, err)

	var hash common.Hash
	bytes, err := args.Pack(reports, salt, sender)
	require.NoError(t, err)

	hash = crypto.Keccak256Hash(bytes)
	return new(big.Int).SetBytes(hash[:])

}

func TestRewardsDistribution(t *testing.T) {
	// we purposefully pick an epoch period that does not divide evenly into the vote period
	// this checks that the rewards are distributed correctly even when the rewards distribution does not
	// align with the vote period end (default vote period in testing is 10 blocks)
	epochPeriod := 55

	setup := func() *Runner {
		r := Setup(t, func(genesis *params.AutonityContractGenesis) *params.AutonityContractGenesis {
			genesis.ProposerRewardRate = 0
			// set oracle reward rate to 100% to simplify the test
			genesis.OracleRewardRate = 10_000
			genesis.EpochPeriod = uint64(epochPeriod)
			genesis.TreasuryFee = 0
			return genesis
		})
		return r
	}

	RunWithSetup("all correctly reporting voters get equal reward", setup, func(r *Runner) {
		// set up some rewards
		r.GiveMeSomeMoney(r.Autonity.address, big.NewInt(1_000_000_000_000))
		totalRewards := r.RewardsAfterOneEpoch()

		oracle := r.Oracle

		symbols, _, err := oracle.GetSymbols(nil)
		require.NoError(t, err)

		period, _, err := oracle.GetVotePeriod(nil)
		require.NoError(t, err)
		votePeriod := int(period.Int64())

		ntnStakes := make(map[common.Address]*big.Int)
		atnBalances := make(map[common.Address]*big.Int)
		for _, val := range r.Committee.Validators {
			ntnStakes[val.Treasury] = val.SelfBondedStake
			atnBalances[val.Treasury] = r.GetBalanceOf(val.Treasury)
		}

		voters := func() []common.Address {
			var vs []common.Address
			for _, val := range r.Committee.Validators {
				vs = append(vs, val.OracleAddress)
			}
			return vs
		}()

		// initial commit
		for _, voter := range voters {
			_, err := oracle.Vote(
				&runOptions{origin: voter},
				makeCommit(r.T, big.NewInt(0), voter, genReports(len(symbols))),
				nil,
				big.NewInt(1),
				0,
			)
			require.NoError(t, err)
		}

		r.WaitNBlocks(votePeriod)

		// now vote until the end of the epoch
		for i := 0; i < (epochPeriod/votePeriod)-1; i++ {
			for _, voter := range voters {
				_, err := oracle.Vote(
					&runOptions{origin: voter},
					makeCommit(r.T, big.NewInt(int64(i+1)), voter, genReports(len(symbols))),
					genReports(len(symbols)),
					big.NewInt(int64(i)),
					0,
				)
				require.NoError(t, err)
			}
			r.WaitNBlocks(votePeriod)
		}

		// check performance
		for _, voter := range voters {
			performance, _, err := oracle.GetRewardPeriodPerformance(nil, voter)
			require.NoError(t, err)
			// should be confidence * n_symbols * n_rounds
			require.Equal(t, big.NewInt(int64(len(symbols)*100*(epochPeriod/votePeriod-1))), performance)
		}

		// there may still be some blocks left in the epoch
		r.WaitNextEpoch()

		// performance should be reset to zero
		for _, voter := range voters {
			performance, _, err := oracle.GetRewardPeriodPerformance(nil, voter)
			require.NoError(t, err)
			require.Equal(t, uint64(0), performance.Uint64())
		}

		// rewards should be distributed according to performance / num_voters
		expectedNTNReward := new(big.Int).Div(totalRewards.RewardNTN, big.NewInt(int64(len(voters))))
		expectedATNReward := new(big.Int).Div(totalRewards.RewardATN, big.NewInt(int64(len(voters))))
		for _, val := range r.Committee.Validators {
			ntnStakeBefore := ntnStakes[val.Treasury]
			atnBalanceBefore := atnBalances[val.Treasury]

			info, _, err := r.Autonity.GetValidator(nil, val.NodeAddress)
			require.NoError(t, err)

			ntnStakeAfter := info.SelfBondedStake
			atnBalanceAfter := r.GetBalanceOf(val.Treasury)

			diffNTN := new(big.Int).Sub(ntnStakeAfter, ntnStakeBefore)
			diffATN := new(big.Int).Sub(atnBalanceAfter, atnBalanceBefore)

			require.Equal(t, expectedNTNReward, diffNTN)
			require.True(t, diffATN.Cmp(big.NewInt(0)) > 0)
			require.Equal(t, expectedATNReward, diffATN)
		}
	})

	RunWithSetup("old voters can also get rewards", setup, func(r *Runner) {
		symbols, _, err := r.Oracle.GetSymbols(nil)
		require.NoError(r.T, err)
		oldVoter := r.Committee.Validators[0]
		// remove `oldVoter` from committee by reducing committee size
		r.NoError(
			r.Autonity.SetCommitteeSize(
				r.Operator,
				big.NewInt(int64(len(r.Committee.Validators))-1),
			),
		)
		r.NoError(
			r.Autonity.Unbond(
				FromSender(oldVoter.Treasury, nil),
				oldVoter.NodeAddress,
				oldVoter.SelfBondedStake,
			),
		)
		epochInfo, _, err := r.Autonity.GetEpochInfo(nil)
		require.NoError(r.T, err)
		r.WaitNBlocks(int(epochInfo.NextEpochBlock.Int64() - r.Evm.Context.BlockNumber.Int64()))
		currentRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(r.T, err)
		currentEpochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		votePeriod, _, err := r.Oracle.GetVotePeriod(nil)
		require.NoError(r.T, err)

		// the next block ends epoch, vote before it so the next vote earns reward
		r.NoError(
			r.Oracle.Vote(
				FromSender(oldVoter.OracleAddress, nil),
				makeCommit(r.T, common.Big0, oldVoter.OracleAddress, genReports(len(symbols))),
				nil,
				common.Big0,
				0,
			),
		)
		r.WaitNBlocks(int(votePeriod.Int64()))
		newRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(currentRound, common.Big1), newRound)
		newEpochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(currentEpochID, common.Big1), newEpochID)
		// check if `oldVoter` is removed from committee
		committee, _, err := r.Autonity.GetCommittee(nil)
		require.NoError(r.T, err)
		for _, c := range committee {
			require.NotEqual(r.T, oldVoter.NodeAddress, c.Addr)
		}

		// vote again from old voter and get reward in the next epoch
		r.NoError(
			r.Oracle.Vote(
				FromSender(oldVoter.OracleAddress, nil),
				makeCommit(r.T, common.Big0, oldVoter.OracleAddress, genReports(len(symbols))),
				genReports(len(symbols)),
				common.Big0,
				0,
			),
		)
		r.GiveMeSomeMoney(r.Autonity.address, big.NewInt(1000_000_000))
		rewards := r.RewardsAfterOneEpoch()
		rewards.RewardNTN = new(big.Int).Add(rewards.RewardNTN, r.GetNewtonBalanceOf(r.Oracle.address))
		atnBalance := r.GetBalanceOf(oldVoter.Treasury)
		validatorInfo, _, err := r.Autonity.GetValidator(nil, oldVoter.NodeAddress)
		require.NoError(r.T, err)
		selfBondedStake := validatorInfo.SelfBondedStake
		r.WaitNextEpoch()
		require.Equal(r.T, new(big.Int).Add(atnBalance, rewards.RewardATN), r.GetBalanceOf(oldVoter.Treasury), "did not get atn reward")
		validatorInfo, _, err = r.Autonity.GetValidator(nil, oldVoter.NodeAddress)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Add(selfBondedStake, rewards.RewardNTN), validatorInfo.SelfBondedStake, "did not get ntn reward")
	})
}

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
			voterValidator, _, err := r.Oracle.VoterValidators(nil, v.OracleAddress)
			require.NoError(r.T, err)

			voterTreasury, _, err := r.Oracle.VoterTreasuries(nil, v.OracleAddress)
			require.NoError(r.T, err)
			require.Equal(r.T, v.Treasury, voterTreasury)
			require.Equal(r.T, v.NodeAddress, voterValidator)

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
			voterValidator, _, err := r.Oracle.VoterValidators(nil, v)
			require.NoError(r.T, err)
			validator, _, err := r.Autonity.GetValidator(nil, voterValidator)
			require.NoError(r.T, err)
			voterTreasury, _, err := r.Oracle.VoterTreasuries(nil, v)
			require.NoError(r.T, err)
			require.Equal(r.T, validator.Treasury, voterTreasury)
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

	// set `newVoterImmediateAccess = true` if voting round ended after epoch is ended
	checkVoterUpdate := func(r *Runner, newVoterImmediateAccess bool, oldVoters map[common.Address]struct{}) {
		round, _, err := r.Oracle.GetRound(nil)
		require.NoError(r.T, err)

		if !newVoterImmediateAccess {
			// voting round did not end yet, so new voters did not get access yet
			voterCheck(r, oldVoters)
			newVoterCheck(r, oldVoters, false)
			// progress round
			progressRound(r, round)
			// new voters should get access for voting, but voters array not updated yet
			round = new(big.Int).Add(round, common.Big1)
		}

		// new voters should get access but voters array not updated yet
		voterCheck(r, oldVoters)
		newVoterCheck(r, oldVoters, true)

		// progress round
		progressRound(r, round)
		// voters array should be updated
		voters := getVoters(r)
		newVoterCheck(r, voters, true)
		voterCheck(r, voters)

		// check if old voters got their access removed
		for v := range oldVoters {
			if _, ok := voters[v]; !ok {
				voterInfo, _, err := r.Oracle.VoterInfo(nil, v)
				require.NoError(r.T, err)
				require.False(r.T, voterInfo.IsVoter)
			}
		}
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
		checkVoterUpdate(r, false, oldVoters)
	})

	RunWithSetup("new voters are updated properly (new voters and old voters are same)", setup, func(r *Runner) {
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkVoterUpdate(r, false, oldVoters)
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
		checkVoterUpdate(r, false, oldVoters)
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
		checkVoterUpdate(r, false, oldVoters)
	})

	RunWithSetup("new voters are updated properly (new voters set is a superset of old voters set)", setup, func(r *Runner) {
		newCommitteeSet := addToCommittee(r, 1)
		for _, c := range r.Committee.Validators {
			newCommitteeSet[c.NodeAddress] = struct{}{}
		}
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkCommittee(r, newCommitteeSet)
		checkVoterUpdate(r, false, oldVoters)
	})

	RunWithSetup("new voters are updated properly with edge case on voting period (votingPeriod * 2 = epochPeriod)", setup, func(r *Runner) {
		epochPeriod, _, err := r.Autonity.GetEpochPeriod(nil)
		require.NoError(r.T, err)
		votingPeriod := new(big.Int).Div(epochPeriod, big.NewInt(2))
		r.NoError(
			r.Oracle.SetVotePeriod(r.Operator, votingPeriod),
		)
		if epochPeriod.Int64()%2 == 1 {
			r.NoError(
				r.Autonity.SetEpochPeriod(r.Operator, new(big.Int).Sub(epochPeriod, common.Big1)),
			)
		}
		r.WaitNextEpoch()
		newCommitteeSet := addToCommittee(r, 2)
		removeFromCommittee(r, len(r.Committee.Validators))
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkCommittee(r, newCommitteeSet)
		checkVoterUpdate(r, false, oldVoters)
	})

	RunWithSetup("voting period and epoch period ends together with edge case (votingPeriod * 2 = epochPeriod)", setup, func(r *Runner) {
		// make the vote period 1 so it ends with epoch period
		r.NoError(
			r.Oracle.SetVotePeriod(r.Operator, common.Big1),
		)
		epochPeriod, _, err := r.Autonity.GetEpochPeriod(nil)
		require.NoError(r.T, err)
		if epochPeriod.Int64()%2 == 1 {
			r.NoError(
				r.Autonity.SetEpochPeriod(
					r.Operator,
					new(big.Int).Add(epochPeriod, common.Big1),
				),
			)
		}
		r.WaitNextEpoch()
		// new voting round and epoch starts together
		// set the edge case (votingPeriod * 2 = epochPeriod)
		epochPeriod, _, err = r.Autonity.GetEpochPeriod(nil)
		require.NoError(r.T, err)
		r.NoError(
			r.Oracle.SetVotePeriod(
				r.Operator,
				new(big.Int).Div(epochPeriod, common.Big2),
			),
		)

		// test the conditions
		votePeriod, _, err := r.Oracle.GetVotePeriod(nil)
		require.NoError(r.T, err)
		require.Equal(r.T, new(big.Int).Mul(votePeriod, common.Big2), epochPeriod, "(votingPeriod * 2 = epochPeriod) not true")
		// see if the voting round and epoch coincides
		epochID, _, err := r.Autonity.EpochID(nil)
		require.NoError(r.T, err)
		targetEpochID := new(big.Int).Add(epochID, common.Big1)
		votingRound, _, err := r.Oracle.GetRound(nil)
		require.NoError(r.T, err)
		targetVotingRound := new(big.Int).Add(votingRound, common.Big2)

		for epochID.Cmp(targetEpochID) == -1 && votingRound.Cmp(targetVotingRound) == -1 {
			r.WaitNBlocks(1)
			epochID, _, err = r.Autonity.EpochID(nil)
			require.NoError(r.T, err)
			votingRound, _, err = r.Oracle.GetRound(nil)
			require.NoError(r.T, err)
		}

		require.Equal(r.T, targetEpochID, epochID, "voting round reached but epoch did not")
		require.Equal(r.T, targetVotingRound, votingRound, "epoch reached but voting round did not")

		// check voter update
		newCommitteeSet := addToCommittee(r, 2)
		removeFromCommittee(r, len(r.Committee.Validators))
		oldVoters := getVoters(r)
		r.WaitNextEpoch()
		checkCommittee(r, newCommitteeSet)
		checkVoterUpdate(r, true, oldVoters)
		// check again
		newCommitteeSet = addToCommittee(r, 2)
		removeFromCommittee(r, len(r.Committee.Validators))
		oldVoters = getVoters(r)
		r.WaitNextEpoch()
		checkCommittee(r, newCommitteeSet)
		checkVoterUpdate(r, true, oldVoters)
	})
}

func TestAllOutliersAreNotSlashed(t *testing.T) {
	setup := func() *Runner {
		r := Setup(t, nil)
		r.NoError(
			r.Oracle.SetSlashingConfig(
				r.Operator,
				big.NewInt(int64(params.DefaultGenesisOracleConfig.OutlierSlashingThreshold)),  // 10%
				big.NewInt(int64(params.DefaultGenesisOracleConfig.OutlierDetectionThreshold)), // 15%
				big.NewInt(int64(params.DefaultGenesisOracleConfig.BaseSlashingRate)),
			),
		)
		return r
	}

	testSlashing := func(
		r *Runner, voters int,
		slashed, outliers []bool,
		prices []int,
	) {
		symbols, _, err := r.Oracle.GetSymbols(nil)
		require.NoError(r.T, err)
		validators := make([]common.Address, 0, voters)
		oracles := make([]common.Address, 0, voters)
		stakes := make([]*big.Int, 0, voters)
		for i := 0; i < voters; i++ {
			validators = append(validators, r.Committee.Validators[i].NodeAddress)
			oracles = append(oracles, r.Committee.Validators[i].OracleAddress)
			stakes = append(stakes, r.Committee.Validators[i].BondedStake)
		}

		vote := func() {
			for i, v := range oracles {
				r.NoError(
					r.Oracle.Vote(
						FromSender(v, nil),
						makeCommit(r.T, common.Big0, v, genReports(len(symbols), prices[i])),
						genReports(len(symbols), prices[i]),
						common.Big0,
						0,
					),
				)
			}
		}

		nextRound := func() {
			round, _, err := r.Oracle.GetRound(nil)
			require.NoError(r.T, err)
			votePeriod, _, err := r.Oracle.GetVotePeriod(nil)
			require.NoError(r.T, err)
			r.WaitNBlocks(int(votePeriod.Int64()))
			newRound, _, err := r.Oracle.GetRound(nil)
			require.NoError(r.T, err)
			require.Equal(r.T, new(big.Int).Add(round, common.Big1), newRound)
		}

		vote()
		nextRound()
		vote()
		for _, v := range oracles {
			info, _, err := r.Oracle.VoterInfo(nil, v)
			require.NoError(r.T, err)
			require.True(r.T, info.ReportAvailable)
			require.True(r.T, info.IsVoter)
		}
		nextRound()

		for i, v := range validators {
			valInfo, _, err := r.Autonity.GetValidator(nil, v)
			require.NoError(r.T, err)
			if slashed[i] {
				require.True(r.T, valInfo.BondedStake.Cmp(stakes[i]) == -1, "did not get slashed")
			} else {
				require.True(r.T, valInfo.BondedStake.Cmp(stakes[i]) == 0, "got slashed")
			}

			voterInfo, _, err := r.Oracle.VoterInfo(nil, oracles[i])
			require.NoError(r.T, err)
			require.True(r.T, voterInfo.IsVoter)
			if outliers[i] {
				require.False(r.T, voterInfo.ReportAvailable)
			} else {
				require.True(r.T, voterInfo.ReportAvailable)
			}
		}
	}

	RunWithSetup("outliers ratio < sqrt(slashingThreshold) are not slashed", setup, func(r *Runner) {
		// change the prices if params.DefaultGenesisOracleConfig is changed
		prices := []int{85, 100, 116}
		outliers := []bool{true, false, true}
		slashed := []bool{false, false, true}
		testSlashing(r, len(prices), slashed, outliers, prices)
	})

	RunWithSetup("outliers ratio < sqrt(slashingThreshold) are not slashed", setup, func(r *Runner) {
		// change the prices if params.DefaultGenesisOracleConfig is changed
		prices := []int{84, 100, 115}
		outliers := []bool{true, false, true}
		slashed := []bool{true, false, false}
		testSlashing(r, len(prices), slashed, outliers, prices)
	})
}
