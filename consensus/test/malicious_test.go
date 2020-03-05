package test

import (
	"fmt"
	"testing"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus"
	tendermintCore "github.com/clearmatics/autonity/consensus/tendermint/core"
	"github.com/clearmatics/autonity/log"
)

func TestTendermintOneMalicious(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	addedValidatorsBlocks := make(map[common.Hash]uint64)
	removedValidatorsBlocks := make(map[common.Hash]uint64)
	changedValidators := tendermintCore.NewChanges()

	cases := []*testCase{
		{
			name:          "replace a valid validator with invalid one",
			numValidators: 5,
			numBlocks:     10,
			txPerPeer:     1,
			maliciousPeers: map[string]injectors{
				"VE": {
					cons: func(basic consensus.Engine) consensus.Engine {
						return tendermintCore.NewReplaceValidatorCore(basic, changedValidators)
					},
					backs: func(basic tendermintCore.Backend) tendermintCore.Backend {
						return tendermintCore.NewChangeValidatorBackend(t, basic, changedValidators, removedValidatorsBlocks, addedValidatorsBlocks)
					},
				},
			},
			addedValidatorsBlocks:   addedValidatorsBlocks,
			removedValidatorsBlocks: removedValidatorsBlocks,
			changedValidators:       changedValidators,
		},
		{
			name:          "add a validator",
			numValidators: 5,
			numBlocks:     10,
			txPerPeer:     1,
			maliciousPeers: map[string]injectors{
				"VE": {
					cons: func(basic consensus.Engine) consensus.Engine {
						return tendermintCore.NewAddValidatorCore(basic, changedValidators)
					},
					backs: func(basic tendermintCore.Backend) tendermintCore.Backend {
						return tendermintCore.NewChangeValidatorBackend(t, basic, changedValidators, removedValidatorsBlocks, addedValidatorsBlocks)
					},
				},
			},
			addedValidatorsBlocks:   addedValidatorsBlocks,
			removedValidatorsBlocks: removedValidatorsBlocks,
			changedValidators:       changedValidators,
		},
		{
			name:          "remove a validator",
			numValidators: 5,
			numBlocks:     10,
			txPerPeer:     1,
			maliciousPeers: map[string]injectors{
				"VE": {
					cons: func(basic consensus.Engine) consensus.Engine {
						return tendermintCore.NewRemoveValidatorCore(basic, changedValidators)
					},
					backs: func(basic tendermintCore.Backend) tendermintCore.Backend {
						return tendermintCore.NewChangeValidatorBackend(t, basic, changedValidators, removedValidatorsBlocks, addedValidatorsBlocks)
					},
				},
			},
			addedValidatorsBlocks:   addedValidatorsBlocks,
			removedValidatorsBlocks: removedValidatorsBlocks,
			changedValidators:       changedValidators,
		},
	}

	for _, testCase := range cases {
		testCase := testCase
		t.Run(fmt.Sprintf("test case %s", testCase.name), func(t *testing.T) {
			runTest(t, testCase, log.LvlError)
		})
	}
}
