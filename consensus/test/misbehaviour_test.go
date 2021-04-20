package test

import (
	"fmt"
	"github.com/clearmatics/autonity/autonity"
	"github.com/clearmatics/autonity/consensus/tendermint/config"
	"github.com/clearmatics/autonity/consensus/tendermint/faultdetector"
	"github.com/clearmatics/autonity/core"
	"github.com/clearmatics/autonity/rlp"
	"testing"
)

func TestMaliciousBehaviours(t *testing.T) {
	requireMisbehaviourProof := func(t *testing.T, rule faultdetector.Rule, node *testNode) {
		stateDB, err := node.service.BlockChain().State()
		if err != nil {
			t.Fatal(err)
		}
		misBehaviours := node.service.BlockChain().GetAutonityContract().GetMisBehaviours(
			node.service.BlockChain().CurrentHeader(),
			stateDB,
		)

		presented := false
		for _, ms := range misBehaviours {
			if ms.Type == autonity.Misbehaviour {
				p := new(faultdetector.Proof)
				err := rlp.DecodeBytes(ms.Rawproof, p)
				if err != nil {
					continue
				} else {
					if p.Rule == rule {
						presented = true
					}
				}
			}
		}

		if !presented {
			t.Fatal("The misbehaviour of rule is not expected on chain")
		}
	}

	genesisHookPN := func(g *core.Genesis) *core.Genesis {
		g.Config.Tendermint.MisbehaviourConfig = &config.MisbehaviourConfig{
			MisbehaviourRuleID: uint8(faultdetector.PN),
		}
		return g
	}
	RulePNChecker := func(t *testing.T, validators map[string]*testNode) {
		requireMisbehaviourProof(t, faultdetector.PN, validators["VA"])
	}

	genesisHookPO := func(g *core.Genesis) *core.Genesis {
		g.Config.Tendermint.MisbehaviourConfig = &config.MisbehaviourConfig{
			MisbehaviourRuleID: uint8(faultdetector.PO),
		}
		return g
	}
	RulePOChecker := func(t *testing.T, validators map[string]*testNode) {
		requireMisbehaviourProof(t, faultdetector.PO, validators["VA"])
	}

	genesisHookPVN := func(g *core.Genesis) *core.Genesis {
		g.Config.Tendermint.MisbehaviourConfig = &config.MisbehaviourConfig{
			MisbehaviourRuleID: uint8(faultdetector.PVN),
		}
		return g
	}
	RulePVNChecker := func(t *testing.T, validators map[string]*testNode) {
		requireMisbehaviourProof(t, faultdetector.PVN, validators["VA"])
	}

	genesisHookC := func(g *core.Genesis) *core.Genesis {
		g.Config.Tendermint.MisbehaviourConfig = &config.MisbehaviourConfig{
			MisbehaviourRuleID: uint8(faultdetector.C),
		}
		return g
	}
	RuleCChecker := func(t *testing.T, validators map[string]*testNode) {
		requireMisbehaviourProof(t, faultdetector.C, validators["VA"])
	}

	genesisHookInvalidProposal := func(g *core.Genesis) *core.Genesis {
		g.Config.Tendermint.MisbehaviourConfig = &config.MisbehaviourConfig{
			MisbehaviourRuleID: uint8(faultdetector.InvalidProposal),
		}
		return g
	}
	RuleInvalidProposalChecker := func(t *testing.T, validators map[string]*testNode) {
		requireMisbehaviourProof(t, faultdetector.InvalidProposal, validators["VA"])
	}

	genesisHookInvalidProposer := func(g *core.Genesis) *core.Genesis {
		g.Config.Tendermint.MisbehaviourConfig = &config.MisbehaviourConfig{
			MisbehaviourRuleID: uint8(faultdetector.InvalidProposer),
		}
		return g
	}
	RuleInvalidProposerChecker := func(t *testing.T, validators map[string]*testNode) {
		// check misbehaviour of PO is presented.
		requireMisbehaviourProof(t, faultdetector.InvalidProposer, validators["VA"])
	}

	genesisHookEquivocation := func(g *core.Genesis) *core.Genesis {
		g.Config.Tendermint.MisbehaviourConfig = &config.MisbehaviourConfig{
			MisbehaviourRuleID: uint8(faultdetector.Equivocation),
		}
		return g
	}
	RuleEquivocationChecker := func(t *testing.T, validators map[string]*testNode) {
		requireMisbehaviourProof(t, faultdetector.Equivocation, validators["VA"])
	}

	cases := []*testCase{
		{
			name: "TestFaultDetectorMaliciousBehaviourPN",
			numValidators: 6,
			numBlocks: 60,
			genesisHook: genesisHookPN,
			finalAssert: RulePNChecker,
		},
		{
			name: "TestFaultDetectorMaliciousBehaviourPO",
			numValidators: 6,
			numBlocks: 60,
			genesisHook: genesisHookPO,
			finalAssert: RulePOChecker,
		},
		{
			name: "TestFaultDetectorMaliciousBehaviourPVN",
			numValidators: 6,
			numBlocks: 60,
			genesisHook: genesisHookPVN,
			finalAssert: RulePVNChecker,
		},
		{
			name: "TestFaultDetectorMaliciousBehaviourC",
			numValidators: 6,
			numBlocks: 60,
			genesisHook: genesisHookC,
			finalAssert: RuleCChecker,
		},
		{
			name: "TestFaultDetectorMaliciousBehaviourInvalidProposal",
			numValidators: 6,
			numBlocks: 60,
			genesisHook: genesisHookInvalidProposal,
			finalAssert: RuleInvalidProposalChecker,
		},
		{
			name: "TestFaultDetectorMaliciousBehaviourInvalidProposer",
			numValidators: 6,
			numBlocks: 60,
			genesisHook: genesisHookInvalidProposer,
			finalAssert: RuleInvalidProposerChecker,
		},
		{
			name: "TestFaultDetectorMaliciousBehaviourEquivocation",
			numValidators: 6,
			numBlocks: 60,
			genesisHook: genesisHookEquivocation,
			finalAssert: RuleEquivocationChecker,
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
			runTest(t, testCase)
		})
	}
}