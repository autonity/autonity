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

func TestMaliciousBehaviourPN(t *testing.T) {
	genesisHookPN := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.PN)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RulePNChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.PN, autonity.Misbehaviour, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourPN",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookPN,
		finalAssert:   RulePNChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestMaliciousBehaviourPO(t *testing.T) {
	genesisHookPO := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.PO)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RulePOChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.PO, autonity.Misbehaviour, validators["VA"])
	}
	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourPO",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookPO,
		finalAssert:   RulePOChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestMaliciousBehaviourPVN(t *testing.T) {
	genesisHookPVN := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.PVN)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RulePVNChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.PVN, autonity.Misbehaviour, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourPVN",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookPVN,
		finalAssert:   RulePVNChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestMaliciousBehaviourPVO1(t *testing.T) {
	genesisHookPVO1 := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.PVO1)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RulePVO1Checker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.PVO1, autonity.Misbehaviour, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourPVO1",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookPVO1,
		finalAssert:   RulePVO1Checker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestMaliciousBehaviourPVO2(t *testing.T) {
	genesisHookPVO2 := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.PVO2)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RulePVO2Checker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.PVO2, autonity.Misbehaviour, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourPVO2",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookPVO2,
		finalAssert:   RulePVO2Checker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestMaliciousBehaviourC(t *testing.T) {
	genesisHookC := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.C)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleCChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.C, autonity.Misbehaviour, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourC",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookC,
		finalAssert:   RuleCChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestMaliciousBehaviourInvalidProposal(t *testing.T) {
	genesisHookInvalidProposal := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.InvalidProposal)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleInvalidProposalChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.InvalidProposal, autonity.Misbehaviour, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourInvalidProposal",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookInvalidProposal,
		finalAssert:   RuleInvalidProposalChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestMaliciousBehaviourInvalidProposer(t *testing.T) {
	genesisHookInvalidProposer := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.InvalidProposer)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleInvalidProposerChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.InvalidProposer, autonity.Misbehaviour, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourInvalidProposer",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookInvalidProposer,
		finalAssert:   RuleInvalidProposerChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestMaliciousBehaviourEquivocation(t *testing.T) {
	genesisHookEquivocation := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.MisbehaviourRuleID = new(uint8)
		*c.MisbehaviourRuleID = uint8(faultdetector.Equivocation)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleEquivocationChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.Equivocation, autonity.Misbehaviour, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorMaliciousBehaviourInvalidProposer",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookEquivocation,
		finalAssert:   RuleEquivocationChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestAccusationRulePO(t *testing.T) {
	genesisHookAccusationPO := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.AccusationRuleID = new(uint8)
		*c.AccusationRuleID = uint8(faultdetector.PO)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleAccusationPOChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.PO, autonity.Accusation, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorAccusationPO",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookAccusationPO,
		finalAssert:   RuleAccusationPOChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestAccusationRulePVN(t *testing.T) {
	genesisHookAccusationPVN := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.AccusationRuleID = new(uint8)
		*c.AccusationRuleID = uint8(faultdetector.PVN)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleAccusationPVNChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.PVN, autonity.Accusation, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorAccusationPVN",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookAccusationPVN,
		finalAssert:   RuleAccusationPVNChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestAccusationRulePVO(t *testing.T) {
	genesisHookAccusationPVO := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.AccusationRuleID = new(uint8)
		*c.AccusationRuleID = uint8(faultdetector.PVO)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleAccusationPVOChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.PVO, autonity.Accusation, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorAccusationPVO",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookAccusationPVO,
		finalAssert:   RuleAccusationPVOChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestAccusationRuleC(t *testing.T) {
	genesisHookAccusationC := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.AccusationRuleID = new(uint8)
		*c.AccusationRuleID = uint8(faultdetector.C)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleAccusationCChecker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.C, autonity.Accusation, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorAccusationC",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookAccusationC,
		finalAssert:   RuleAccusationCChecker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func TestAccusationRuleC1(t *testing.T) {
	genesisHookAccusationC1 := func(g *core.Genesis) *core.Genesis {
		var c config.MisbehaviourConfig
		c.AccusationRuleID = new(uint8)
		*c.AccusationRuleID = uint8(faultdetector.C1)
		g.Config.Tendermint.BlockPeriod = 1
		g.Config.Tendermint.MisbehaviourConfig = &c
		return g
	}
	RuleAccusationC1Checker := func(t *testing.T, validators map[string]*testNode) {
		requireOnChainProof(t, faultdetector.C1, autonity.Accusation, validators["VA"])
	}

	tc := testCase{
		name:          "TestFaultDetectorAccusationC1",
		numValidators: 6,
		numBlocks:     60,
		genesisHook:   genesisHookAccusationC1,
		finalAssert:   RuleAccusationC1Checker,
	}
	t.Run(fmt.Sprintf("test case %s", tc.name), func(t *testing.T) {
		runTest(t, &tc)
	})
}

func requireOnChainProof(t *testing.T, rule faultdetector.Rule, proofType uint8, node *testNode) {
	stateDB, err := node.service.BlockChain().State()
	if err != nil {
		t.Fatal(err)
	}

	var proofs []autonity.OnChainProof

	if proofType == autonity.Misbehaviour {
		proofs = node.service.BlockChain().GetAutonityContract().GetMisBehaviours(
			node.service.BlockChain().CurrentHeader(),
			stateDB,
		)
	} else {
		proofs = node.service.BlockChain().GetAutonityContract().GetAccusations(
			node.service.BlockChain().CurrentHeader(),
			stateDB,
		)
	}

	presented := false
	for _, ms := range proofs {
		if ms.Type == proofType {
			p := new(faultdetector.Proof)
			err := rlp.DecodeBytes(ms.Rawproof, p)
			if err != nil {
				continue
			}

			if p.Rule == rule {
				presented = true
			}
		}
	}

	if !presented {
		t.Fatal("The on-chain proof of rule is not expected on chain")
	}
}
