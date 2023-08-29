/* TODO(tariq) tests for accountability.sol
 *
 * high priority:
 * 2. verify that when having multiple slashing events in the slashing queue, the offenceCount goes up (and consequently the slashing rate). _performSlashingTasks function
 * 3. verify that a validator with an history of past offences gets slashed more than a clean one (and exactly the amount more we expect).
 * 4. Accusation flow tests (test canAccuse/canSlash when appriopriate)
 *    0. issue multiple accusations on different blocks --> check that only the one who expired get converted to misbehavior
      1. Validator is accused and submit proof of innocence before the window is expired --> no slashing
      2. Validator is accused and does not submit proof of innocence --> accusation is promoted to misbehavior and validator gets slashed
      3. Validator is accused and submits proof of innocence **after** the window is expired --> accusation is promoted to misbehavior and validator gets slashed
      3. Validator is accused while already under accusation --> 2nd accusation reverts
        - canAccuse should return a deadline for when we can submit the 2nd accusation
      4. Validators is under accusation and someone sends proof of misbehavior against him (for the same epoch of the accusation). 
         The accused validator does not publish proof of innocence for the accusation. Outcome:
           - if misbehaviour severity > accusation severity --> only misbehaviour slashing takes effect
           - if misbehaviour severity < accusation severity --> both offences are slashed
  * 5. cannot submit misbehaviour for validator already slashed for the offence epoch with a higher severity than the submitted misb
  *       require(slashingHistory[_offender][_epoch] < _severity, "already slashed at the proof's epoch");
  *     - also canSlash should return false
  * 6. same thing for accusation
  *     - canAccuse should return false.
  * 7. edge scenario. validator is sentenced for 2 misbehaviour with 1st misb severity < 2nd misb severity in the same epoch. He should be slashed for both
  * 8. validator already slashed for an epoch, but accusation with higher severity is issued against him --> accusation is valid and should lead to slashing if not addressed
  * 9. other edge cases?

low priority:
 * 3. verify rule --> severity mapping
 * 4. verify severity --> slashign rate mapping
 * 5. test chunked event processing (handleEvent function)
 *      - test also case where multiple validators are sending interleaved chunks
 * 6. test _handle* functions edge cases (e.g. invalid proof, block in future, etc.) --> tx should revert
 * 7. whitebox testing (better to leave for when the implementation will be less prone to changes)
 *    - verify that the accusation queue, the slashing queue update and the other internal structures are updated as we expect
 */

'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const Autonity = artifacts.require("Autonity");
const Accountability = artifacts.require("Accountability");
const AccountabilityTest = artifacts.require("AccountabilityTest");
const toBN = web3.utils.toBN

function ruleToRate(accountabilityConfig,rule){
  //TODO(lorenzo) create mapping rule to rate once finalized in autonity.sol. bypass severity conversion?
  return accountabilityConfig.baseSlashingRateMid
}

async function slashAndVerify(autonity,accountability,accountabilityConfig,event,epochOffenceCount){
  let offender = await autonity.getValidator(event.offender)

  let baseRate = ruleToRate(accountabilityConfig,event.rule)

  let slashingRate = toBN(baseRate).add(toBN(epochOffenceCount).mul(toBN(accountabilityConfig.collusionFactor))).add(toBN(offender.provableFaultCount).mul(toBN(accountabilityConfig.historyFactor)));  
  // cannot slash more than 100%
  if(slashingRate.gt(toBN(accountabilityConfig.slashingRatePrecision))) {
    slashingRate = toBN(accountabilityConfig.slashingRatePrecision)
  }

  let availableFunds = toBN(offender.bondedStake).add(toBN(offender.unbondingStake)).add(toBN(offender.selfUnbondingStake))
  let slashingAmount = (slashingRate.mul(availableFunds).div(toBN(accountabilityConfig.slashingRatePrecision))).toNumber() 
  let originalSlashingAmount = slashingAmount

  let autonityTreasury = await autonity.getTreasuryAccount()
  let autonityTreasuryBalance = await autonity.balanceOf(autonityTreasury)
 
  await accountability.slash(event,epochOffenceCount)
  let slashingBlock = await web3.eth.getBlockNumber()
  let offenderSlashed = await autonity.getValidator(offender.nodeAddress);
  
  // first unbonding self stake is slashed (PAS)
  let expectedSelfUnbondingStake = (slashingAmount > parseInt(offender.selfUnbondingStake)) ? 0 : parseInt(offender.selfUnbondingStake) - slashingAmount;
  assert.equal(parseInt(offenderSlashed.selfUnbondingStake), expectedSelfUnbondingStake)
  slashingAmount = (expectedSelfUnbondingStake == 0) ? slashingAmount - parseInt(offender.selfUnbondingStake) : 0
  if(slashingAmount == 0)
    return

  // then self stake is slashed (PAS)
  let expectedSelfBondedStake = (slashingAmount > parseInt(offender.selfBondedStake)) ? 0 : parseInt(offender.selfBondedStake) - slashingAmount;
  assert.equal(parseInt(offenderSlashed.selfBondedStake), expectedSelfBondedStake)
  slashingAmount = (expectedSelfBondedStake == 0) ? slashingAmount - parseInt(offender.selfBondedStake) : 0
  if(slashingAmount == 0)
    return

  // then remaining slash is distributed equally between delegated stake and delegated unbonding stake
  let delegatedStake = parseInt(offender.bondedStake) - parseInt(offender.selfBondedStake)
  let delegatedSlash = (toBN(slashingAmount).mul(toBN(delegatedStake)).div(toBN(delegatedStake).add(toBN(offender.unbondingStake)))).toNumber()
  let unbondingDelegatedSlash = (toBN(slashingAmount).mul(toBN(offender.unbondingStake)).div(toBN(delegatedStake).add(toBN(offender.unbondingStake)))).toNumber()
  assert.equal(parseInt(offenderSlashed.bondedStake) - parseInt(offenderSlashed.selfBondedStake), delegatedStake - delegatedSlash)
  assert.equal(parseInt(offenderSlashed.unbondingStake), parseInt(offender.unbondingStake) - unbondingDelegatedSlash)

  // check total slashed
  assert.equal(parseInt(offenderSlashed.totalSlashed), parseInt(offender.totalSlashed) + originalSlashingAmount)

  // check provable fault count increases
  assert.equal(parseInt(offenderSlashed.provableFaultCount), parseInt(offender.provableFaultCount) + 1)

  // check that validator is jailed for correct amount of time
  // state: 0 --> active, 1 --> paused, 2 --> jailed
  let currentEpochPeriod = await autonity.getEpochPeriod();
  let jailSentence = toBN(offenderSlashed.provableFaultCount).mul(toBN(accountabilityConfig.jailFactor)).mul(currentEpochPeriod)
  assert.equal(parseInt(offenderSlashed.state), 2)
  assert.equal(parseInt(offenderSlashed.jailReleaseBlock),slashingBlock + jailSentence.toNumber())
  
  // check that slashed amount goes to the autonity treasury
  let autonityTreasuryBalanceAfterSlash = await autonity.balanceOf(autonityTreasury)
  assert.equal(autonityTreasuryBalanceAfterSlash.toString(), autonityTreasuryBalance.add(toBN(originalSlashingAmount)).toString())
}

contract('Accountability', function (accounts) {
  before(async function () {
    console.log("\tAttempting to mock enode verifier precompile. Will (rightfully) fail if running against Autonity network")
    await utils.mockEnodePrecompile()
  });
  for (let i = 0; i < accounts.length; i++) {
    console.log("account: ", i, accounts[i]);
  }

  const operator = accounts[5];
  const deployer = accounts[6];
  const anyAccount = accounts[7];
  const name = "Newton";
  const symbol = "NTN";
  const minBaseFee = 5000;
  const committeeSize = 1000;
  const epochPeriod = 30;
  const delegationRate = 100;
  const unBondingPeriod = 60;
  const treasuryAccount = accounts[8];
  const treasuryFee = "10000000000000000";
  const minimumEpochPeriod = 30;
  const version = 0;
  const zeroAddress = "0x0000000000000000000000000000000000000000";

  const autonityConfig = {
    "policy": {
      "treasuryFee": treasuryFee,
      "minBaseFee": minBaseFee,
      "delegationRate": delegationRate,
      "unbondingPeriod" : unBondingPeriod,
      "treasuryAccount": treasuryAccount,
    },
    "contracts": {
      "oracleContract" : zeroAddress, // gets updated in deployContracts() 
      "accountabilityContract": zeroAddress, // gets updated in deployContracts()
      "acuContract" :zeroAddress,
      "supplyControlContract" :zeroAddress,
      "stabilizationContract" :zeroAddress,
    },
    "protocol": {
      "operatorAccount": operator,
      "epochPeriod": epochPeriod,
      "blockPeriod": minimumEpochPeriod,
      "committeeSize": committeeSize,
    },
    "contractVersion":version,
  };

  const accountabilityConfig = {
    "innocenceProofSubmissionWindow": 30,
    "latestAccountabilityEventsRange": 256,
    "baseSlashingRateLow": 500,
    "baseSlashingRateMid": 1000,
    "collusionFactor": 550,
    "historyFactor": 750,
    "jailFactor": 60,
    "slashingRatePrecision": 10000
  }

  const genesisEnodes = [
    "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
    "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
    "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
    "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
  ]

  // precomputed using aut validator compute-address
  // TODO(lorenzo) derive them from enodes or privatekeys
  const genesisNodeAddresses = [
    "0x850C1Eb8D190e05845ad7F84ac95a318C8AaB07f",
    "0x4AD219b58a5b46A1D9662BeAa6a70DB9F570deA5",
    "0xc443C6c6AE98F5110702921138D840e77dA67702",
    "0x09428E8674496e2D1E965402F33A9520c5fCBbE2",
  ]

  const genesisPrivateKeys = [
   "a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91", 
   "aa4b77b1305f8f265e81599587c623d8950624f3e1bd9c121ef2461a7a1e7527",
   "4ec99383dc50aa3f3117fcbfba7b69188ba60d3418185fb353c9a69d066e55d9",
   "0c8698f456533170fe07c6dcb753d47bef8bedd46443efa57a859c989887b56b",
  ]
  
  // enodes with no validator registered at genesis
  const freeEnodes = [
    "enode://a7ecd2c1b8c0c7d7ab9cc12e620605a762865d381eb1bc5417dcf07599571f84ce5725f404f66d3e254d590ae04e4e8f18fe9e23cd29087d095a0c37d0443252@3.209.45.79:30303",
  ];

  // TODO(lorenzo) derive them from enodes or privatekeys
  const freeAddresses = [
    "0xDE03B7806f885Ae79d2aa56568b77caDB0de073E",
  ]

  const freePrivateKeys = [
    "e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b",
  ]

  const baseValidator = {
    "selfBondedStake": 0,
    "totalSlashed": 0,
    "jailReleaseBlock": 0,
    "provableFaultCount" :0,
    "liquidSupply": 0,
    "registrationBlock": 0,
    "state": 0,
    "liquidContract" : zeroAddress,
    "selfUnbondingStake" : 0,
    "selfUnbondingShares" : 0,
    "unbondingStake" : 0,
    "unbondingShares" : 0,
  }
  
  // accounts[2] is skipped because it is used as a genesis validator when running against autonity
  // this can cause interference in reward distribution tests
  const validators = [
    { ...baseValidator,
      "treasury": accounts[0],
      "nodeAddress": genesisNodeAddresses[0],
      "oracleAddress": accounts[0],
      "enode": genesisEnodes[0],
      "commissionRate": 100,
      "bondedStake": 100,
    },
    { ...baseValidator,
      "treasury": accounts[1],
      "nodeAddress": genesisNodeAddresses[1],
      "oracleAddress": accounts[1],
      "enode": genesisEnodes[1],
      "commissionRate": 100,
      "bondedStake": 90,
    },
    { ...baseValidator,
      "treasury": accounts[3],
      "nodeAddress": genesisNodeAddresses[2],
      "oracleAddress": accounts[3],
      "enode": genesisEnodes[2],
      "commissionRate": 100,
      "bondedStake": 110,
    },
    { ...baseValidator,
      "treasury": accounts[4],
      "nodeAddress": genesisNodeAddresses[3],
      "oracleAddress": accounts[4],
      "enode": genesisEnodes[3],
      "commissionRate": 100,
      "bondedStake": 120,
    },
  ];


  // initial validators ordered by bonded stake
  const orderedValidatorsList = [
    genesisNodeAddresses[0],
    genesisNodeAddresses[1],
    genesisNodeAddresses[2],
    genesisNodeAddresses[3],
  ];

  let autonity;
  let accountability;
  describe.skip('Contract initial state', function () {
    before(async function () {
      autonity = await Autonity.new(validators, autonityConfig, {from: deployer});
      await autonity.finalizeInitialization({from: deployer});
      accountability = await Accountability.new(autonity.address, accountabilityConfig, {from: deployer});
    });
    //TODO(tariq) low priority.
    // test that config gets set properly at contract deploy 
  });
  describe.skip('Contract permissioning', function () {
    before(async function () {
      autonity = await Autonity.new(validators, autonityConfig, {from: deployer});
      await autonity.finalizeInitialization({from: deployer});
      accountability = await Accountability.new(autonity.address, accountabilityConfig, {from: deployer});
    });
    //TODO(tariq) modifiers (low priority)
    // only registered validators can submit accountability events (handleEvent)
    // only autonity can call finalize(), setEpochPeriod() and distributeRewards()
  });
  describe('Slashing', function () {
    beforeEach(async function () {
      autonity = await Autonity.new(validators, autonityConfig, {from: deployer});
      await autonity.finalizeInitialization({from: deployer});
      accountability = await AccountabilityTest.new(autonity.address, accountabilityConfig, {from: deployer});
      await autonity.setAccountabilityContract(accountability.address, {from:operator});
    });
    it("test stake slashing priority (PAS first)", async function() { 
      let offender = await autonity.getValidator(validators[0].nodeAddress)
      let reporter = validators[1].treasury
      let delegator = anyAccount

      // only selfbonded stake after genesis
      let genesisStake = 100
      assert.equal(offender.bondedStake,genesisStake)
      assert.equal(offender.selfBondedStake,genesisStake)
      assert.equal(offender.bondedStake - offender.selfBondedStake,0) // delegatedStake
      assert.equal(offender.selfUnbondingStake,0)
      assert.equal(offender.unbondingStake,0) // unbonding delegated stake

      // add some delegated stake
      let delegatedStake = 100
      await autonity.mint(delegator, delegatedStake, {from: operator});
      await autonity.bond(offender.nodeAddress, delegatedStake, {from: delegator});
      await autonity.finalizeInitialization({from: deployer}) // I use finalizeInitialization as a way to trigger the staking operations
      offender = await autonity.getValidator(validators[0].nodeAddress)
      assert.equal(offender.bondedStake,genesisStake + delegatedStake)
      assert.equal(offender.selfBondedStake,genesisStake)
      assert.equal(offender.bondedStake - offender.selfBondedStake,delegatedStake) 
      assert.equal(offender.selfUnbondingStake,0)
      assert.equal(offender.unbondingStake,0) 

      // unbond some self-bonded stake and some delegated stake
      let unbondSelf = 10
      let unbondDelegated = 50
      await autonity.unbond(offender.nodeAddress, unbondSelf, {from: offender.treasury});
      await autonity.unbond(offender.nodeAddress, unbondDelegated, {from: delegator});
      await autonity.finalizeInitialization({from: deployer}) 
      offender = await autonity.getValidator(validators[0].nodeAddress)
      assert.equal(offender.bondedStake,genesisStake + delegatedStake - unbondSelf - unbondDelegated)
      assert.equal(offender.selfBondedStake,genesisStake - unbondSelf)
      assert.equal(offender.bondedStake - offender.selfBondedStake,delegatedStake - unbondDelegated) 
      assert.equal(offender.selfUnbondingStake,unbondSelf)
      assert.equal(offender.unbondingStake,unbondDelegated) 

      // trigger slash event
      const event = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": 0, // PN rule --> severity mid
        "reporter": reporter,
        "offender": offender.nodeAddress,
        "rawProof": [], // not checked by the _slash function
        "block": 1,
        "epoch": 0,
        "reportingBlock": 2,
        "messageHash": 0, // not checked by the _slash function
      }

      let epochOffenceCount = 1
      await slashAndVerify(autonity,accountability,accountabilityConfig,event,epochOffenceCount);
      
      // let's slash another time to make sure to slash also the non-pas stake this time
      epochOffenceCount = 8
      await slashAndVerify(autonity,accountability,accountabilityConfig,event,epochOffenceCount);
    });
  });
});

