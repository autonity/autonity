package tests

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/crypto"
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
				{big.NewInt(1), big.NewInt(237492837498), big.NewInt(18), big.NewInt(100)},
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
