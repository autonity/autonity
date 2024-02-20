/* TODO(tariq) tests for accountability.sol (low priority)
 * 1. edge case scenario: validator is sentenced for 2 misbehavior in the same epoch. For this to happen the 1st submitted misb needs to have severity < 2nd misb severity. The offender should be slashed for both misb. Instead if the 1st submitted misb has severity >= 2nd misb, only the first misb should lead to slashing. Currently this test case cannot be implemented since we use only severity mid in the autonity contract. 
 * 2. verify rule --> severity mapping
 * 3. verify severity --> slashign rate mapping
 * 4. test _handle* functions edge cases (e.g. invalid proof, block in future, etc.) --> tx should revert
 * 5. whitebox testing (better to leave for when the implementation will be less prone to changes)
 *    - verify that the accusation queue, the slashing queue update and the other internal structures are updated as we expect
 *
 * There might be additional edge cases for slashing, misbehavior and accusation flow to test.
 */

'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const Autonity = artifacts.require("Autonity");
const Accountability = artifacts.require("Accountability");
const AccountabilityTest = artifacts.require("AccountabilityTest");
const toBN = web3.utils.toBN;
const config = require("./config");


function checkEvent(event, offender, reporter, chunkId, rawProof) {
  assert.equal(event.offender, offender, "event offender mismatch");
  assert.equal(event.reporter, reporter, "event reporter mismatch");
  assert.equal(event.chunkId, chunkId, "event chunkId mismatch");
  assert.equal(event.rawProof, rawProof, "event rawProof mismatch")
}

async function slashAndVerify(autonity,accountability,accountabilityConfig,event,epochOffenceCount){
  let offender = await autonity.getValidator(event.offender)

  let baseRate = utils.ruleToRate(accountabilityConfig,event.rule)

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
 
  let tx = await accountability.slash(event,epochOffenceCount)
  let slashingBlock = tx.receipt.blockNumber
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
    await utils.mockPrecompile()
  });
  for (let i = 0; i < accounts.length; i++) {
    console.log("account: ", i, accounts[i]);
  }

  const operator = accounts[5];
  const deployer = accounts[6];
  const anyAccount = accounts[7];
  const treasuryAccount = accounts[8];
  const zeroAddress = "0x0000000000000000000000000000000000000000";

  let autonityConfig = config.autonityConfig(operator, treasuryAccount)

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
  const genesisPrivateKeys = config.GENESIS_PRIVATE_KEYS

  // accounts[2] is skipped because it is used as a genesis validator when running against autonity
  // this can cause interference in reward distribution tests
  let validators = config.validators(accounts);

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
    it("multiple slashing events in the same epoch should lead to increased penalties (collusion)",async function() {
      // insert multiple slashing events for the same epoch with different validators as offender
      const offenderAddresses = [validators[0].nodeAddress, validators[1].nodeAddress, validators[2].nodeAddress]
      const epochOffenceCount = offenderAddresses.length
      const reporter = validators[3].treasury
      const event = {
        "chunks": 1, 
        "chunkId": 1,
        "eventType": 0,
        "rule": 0, // PN rule --> severity mid
        "reporter": reporter,
        "offender":"",
        "rawProof": [], 
        "block": 1,
        "epoch": 0,
        "reportingBlock": 2,
        "messageHash": 0, 
      }
      let offenders = [];
      for (const offenderAddress of offenderAddresses) {
        event.offender = offenderAddress
        let offender = await autonity.getValidator(offenderAddress)
        // they should have only selfbonded stake
        assert.equal(offender.bondedStake,offender.selfBondedStake)
        offenders.push(offender)
        await accountability.handleValidFaultProof(event)
      }

      await accountability.performSlashingTasks()

      for (const offender of offenders) {
        let offenderSlashed = await autonity.getValidator(offender.nodeAddress);

        let baseRate = utils.ruleToRate(accountabilityConfig,event.rule);

        let slashingRate = toBN(baseRate).add(toBN(epochOffenceCount).mul(toBN(accountabilityConfig.collusionFactor))).add(toBN(offender.provableFaultCount).mul(toBN(accountabilityConfig.historyFactor)));  
        // cannot slash more than 100%
        if(slashingRate.gt(toBN(accountabilityConfig.slashingRatePrecision))) {
          slashingRate = toBN(accountabilityConfig.slashingRatePrecision)
        }

        let availableFunds = toBN(offender.bondedStake).add(toBN(offender.unbondingStake)).add(toBN(offender.selfUnbondingStake))
        let slashingAmount = (slashingRate.mul(availableFunds).div(toBN(accountabilityConfig.slashingRatePrecision))).toNumber() 

        assert.equal(parseInt(offenderSlashed.bondedStake),parseInt(offender.bondedStake) - slashingAmount)
      }
    });
    it("a validator with a history of misbehavior should get slashed more",async function() {
      let currentEpochPeriod = (await autonity.getEpochPeriod()).toNumber()
      let reporter = validators[0]
      let offender = validators[1]
      const event = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": 0, // PN rule --> severity mid
        "reporter": reporter.treasury,
        "offender": offender.nodeAddress,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0, 
      }

      // insert 3 past offences in different epochs
      await accountability.handleValidFaultProof(event)
      event.block += currentEpochPeriod
      event.epoch += 1
      event.reportingBlock = event.block + 1
      await accountability.handleValidFaultProof(event)
      event.block += currentEpochPeriod
      event.epoch += 1
      event.reportingBlock = event.block + 1
      await accountability.handleValidFaultProof(event)

      await accountability.performSlashingTasks()
      let offenderValidator = await autonity.getValidator(offender.nodeAddress)
      assert.equal(offenderValidator.provableFaultCount,'3')

      // check slashing rate on fourth offence
      let epochOffenceCount = 0
      event.block += currentEpochPeriod
      event.epoch += 1
      event.reportingBlock = event.block + 1
      await slashAndVerify(autonity,accountability,accountabilityConfig,event,epochOffenceCount);
    });
    /*  Validator is under accusation and someone sends proof of misbehavior against him (for the same epoch of the accusation). 
    *   The accused validator does not publish proof of innocence for the accusation. Outcome:
    *     - if misbehavior severity >= accusation severity --> only misbehavior slashing takes effect
    *     - if misbehavior severity < accusation severity --> both offences are slashed 
    */
    it('edge case: concurrent accusation and misbehavior submission (misb severity >= accusation severity)',async function() {
      let reporter = validators[0]
      let offender = validators[1]
      let offenderInfo = await autonity.getValidator(offender.nodeAddress)
      let PNrule = 0
      let currentBlock = await web3.eth.getBlockNumber()
      const event = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter.treasury,
        "offender": offender.nodeAddress,
        "rawProof": [],
        "block": currentBlock - 1,
        "epoch": 0,
        "reportingBlock": currentBlock,
        "messageHash": 0, 
      }
      
      // only selfbondedstake
      assert.equal(offenderInfo.bondedStake,offenderInfo.selfBondedStake)
      
      // insert accusation
      let canAccuse = await accountability.canAccuse(offender.nodeAddress,PNrule,event.block);
      assert.strictEqual(canAccuse._result,true);
      assert.strictEqual(canAccuse._deadline.toString(),'0');
      await accountability.handleValidAccusation(event);

      // insert misbheavior with same severity
      assert.strictEqual(await accountability.canSlash(offender.nodeAddress,PNrule,event.block),true);
      await accountability.handleValidFaultProof(event);
      
      // wait for accusation to expire
      for (let i = 0; i < accountabilityConfig.innocenceProofSubmissionWindow; i++) { await utils.mineEmptyBlock() }

      // promote guilty accusations (no accusation should be promoted since misb severity == accusation severity)
      let tx = await accountability.promoteGuiltyAccusations();
      truffleAssert.eventNotEmitted(tx, 'NewFaultProof')

      // misb should lead to slashing
      tx = await accountability.performSlashingTasks()
      truffleAssert.eventEmitted(tx,'SlashingEvent')
        
      let offenderSlashed = await autonity.getValidator(offender.nodeAddress);

      let epochOffenceCount = 1;
      let baseRate = utils.ruleToRate(accountabilityConfig,event.rule);

      let slashingRate = toBN(baseRate).add(toBN(epochOffenceCount).mul(toBN(accountabilityConfig.collusionFactor))).add(toBN(offender.provableFaultCount).mul(toBN(accountabilityConfig.historyFactor)));  
      // cannot slash more than 100%
      if(slashingRate.gt(toBN(accountabilityConfig.slashingRatePrecision))) {
        slashingRate = toBN(accountabilityConfig.slashingRatePrecision)
      }

      let availableFunds = toBN(offender.bondedStake).add(toBN(offender.unbondingStake)).add(toBN(offender.selfUnbondingStake))
      let slashingAmount = (slashingRate.mul(availableFunds).div(toBN(accountabilityConfig.slashingRatePrecision))).toNumber() 

      assert.equal(parseInt(offenderSlashed.bondedStake),parseInt(offender.bondedStake) - slashingAmount)
    }); 
    it.skip('edge case: concurrent accusation and misbehavior submission (misb severity < accusation severity)',async function() {
      //TODO(tariq) implement this test case. currently not implementable since we use only severity mid in autonity contract.
    }); 
  });
  describe('misbehavior flow', function () {
    beforeEach(async function () {
      autonity = await Autonity.new(validators, autonityConfig, {from: deployer});
      await autonity.finalizeInitialization({from: deployer});
      accountability = await AccountabilityTest.new(autonity.address, accountabilityConfig, {from: deployer});
      await autonity.setAccountabilityContract(accountability.address, {from:operator});
    });
    it("cannot submit misbehavior with severity X for validator already slashed for the offence epoch with severity Y >= X", async function() {
      let reporter = validators[0]
      let offender = validators[1]
      let PNrule = 0
      const event = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter.treasury,
        "offender": offender.nodeAddress,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0, 
      }
      
      assert.strictEqual(await accountability.canSlash(offender.nodeAddress,PNrule,event.block),true);

      await accountability.handleValidFaultProof(event);

      assert.strictEqual(await accountability.canSlash(offender.nodeAddress,PNrule,event.block + 1),false);

      await truffleAssert.fails(
        accountability.handleValidFaultProof(event),
        truffleAssert.ErrorType.REVERT,
        "already slashed at the proof's epoch"
      );
      // TODO(lorenzo) once implemented in contract
      // add canSlash and handleValidFaultProof asserts when submitting a proof of higher severity (slashing is possible in that case)
    });
  });
  describe('accusation flow', function () {
    beforeEach(async function () {
      autonity = await Autonity.new(validators, autonityConfig, {from: deployer});
      await autonity.finalizeInitialization({from: deployer});
      accountability = await AccountabilityTest.new(autonity.address, accountabilityConfig, {from: deployer});
      await autonity.setAccountabilityContract(accountability.address, {from:operator});
    });
    it("cannot submit accusation with severity X for validator already slashed for the offence epoch with severity Y >= X", async function() {
      let reporter = validators[0]
      let offender = validators[1]
      let PNrule = 0
      const event = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter.treasury,
        "offender": offender.nodeAddress,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0, 
      }
      
      let canAccuse = await accountability.canAccuse(offender.nodeAddress,PNrule,event.block);
      assert.strictEqual(canAccuse._result,true);
      assert.strictEqual(canAccuse._deadline.toString(),'0');

      await accountability.handleValidFaultProof(event);
      
      canAccuse = await accountability.canAccuse(offender.nodeAddress,PNrule,event.block+1);
      assert.strictEqual(canAccuse._result,false);
      assert.strictEqual(canAccuse._deadline.toString(),'0');

      await truffleAssert.fails(
        accountability.handleValidAccusation(event),
        truffleAssert.ErrorType.REVERT,
        "already slashed at the proof's epoch"
      );
      // TODO(lorenzo) once implemented in contract
      // add canAccuse and handleValidAccusation asserts when submitting an accusation of higher severity (slashing is possible in that case)

    });
    it("Cannot accuse validator already under accusation", async function() {
      let reporter = validators[0]
      let offender = validators[1]
      let PNrule = 0
      const event = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter.treasury,
        "offender": offender.nodeAddress,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0, 
      }
      let canAccuse = await accountability.canAccuse(offender.nodeAddress,PNrule,event.block);
      assert.strictEqual(canAccuse._result,true);
      assert.strictEqual(canAccuse._deadline.toString(),'0');

      await accountability.handleValidAccusation(event);
      
      canAccuse = await accountability.canAccuse(offender.nodeAddress,PNrule,event.block+1);
      assert.strictEqual(canAccuse._result,false);
      assert.strictEqual(canAccuse._deadline.toNumber(),event.block + accountabilityConfig.innocenceProofSubmissionWindow);

      await truffleAssert.fails(
        accountability.handleValidAccusation(event),
        truffleAssert.ErrorType.REVERT,
        "already processing an accusation"
      );
    });
    it("Only expired unadressed accusations are promoted to misbehavior and lead to slashing", async function() {
      let reporter = validators[0]
      let offender1 = validators[1] // will not post an innocence proof before accusation promotion --> slashed
      let offender2 = validators[2] // will post an innocence proof before accusation promotion --> no slashing
      let offender3 = validators[3] // will be accused later than offender1 and offender2, thus his accusation will not be expired when accusation are promoted
      let PNrule = 0
      const event = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter.treasury,
        "offender": "", // tofill
        "rawProof": [],
        "block": 0, // tofill
        "epoch": 0,
        "reportingBlock": 0, //tofill
        "messageHash": 0, 
      }
      // accuse offender1
      let currentBlock = await web3.eth.getBlockNumber()
      event.offender = offender1.nodeAddress
      event.block = currentBlock - 1
      let offender1Block = event.block
      event.reportingBlock = currentBlock
      let canAccuse = await accountability.canAccuse(event.offender,event.rule,event.block);
      assert.strictEqual(canAccuse._result,true);
      assert.strictEqual(canAccuse._deadline.toString(),'0');
      await accountability.handleValidAccusation(event);
      
      // accuse offender2
      currentBlock = await web3.eth.getBlockNumber()
      event.offender = offender2.nodeAddress
      event.block = currentBlock - 1
      let offender2Block = event.block
      event.reportingBlock = currentBlock
      canAccuse = await accountability.canAccuse(event.offender,event.rule,event.block);
      assert.strictEqual(canAccuse._result,true);
      assert.strictEqual(canAccuse._deadline.toString(),'0');
      await accountability.handleValidAccusation(event);
      
      // accuse offender3 with reportingBlock in the future
      currentBlock = await web3.eth.getBlockNumber()
      event.offender = offender3.nodeAddress
      event.block = currentBlock - 1
      event.reportingBlock = currentBlock + 500
      canAccuse = await accountability.canAccuse(event.offender,event.rule,event.block);
      assert.strictEqual(canAccuse._result,true);
      assert.strictEqual(canAccuse._deadline.toString(),'0');
      await accountability.handleValidAccusation(event);

      // submit valid proof of innocence for offender2
      const proof = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter.treasury,
        "offender": offender2.nodeAddress,
        "rawProof": [],
        "block": offender2Block,
        "epoch": 0,
        "reportingBlock": 0, // does not matter
        "messageHash": 0, // must match accusation's one
      }
      await accountability.handleValidInnocenceProof(proof);

      // wait for accusations to expire
      for (let i = 0; i < accountabilityConfig.innocenceProofSubmissionWindow; i++) { await utils.mineEmptyBlock() }

      // promote accusations. only offender1's accusation should be promoted to misbehavior
      let tx = await accountability.promoteGuiltyAccusations();
      // severity mid == 2
      truffleAssert.eventEmitted(tx, 'NewFaultProof', (ev) => {
        return ev._offender === offender1.nodeAddress && ev._severity == 2 && ev._id == 0
      });

      // canSlash should return false only for offender1
      assert.strictEqual(await accountability.canSlash(offender1.nodeAddress,PNrule,currentBlock),false);
      assert.strictEqual(await accountability.canSlash(offender2.nodeAddress,PNrule,currentBlock),true);
      assert.strictEqual(await accountability.canSlash(offender3.nodeAddress,PNrule,currentBlock),true);

      // offender1 should fail to submit proof of innocence, he is too late
      const proof2 = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter.treasury,
        "offender": offender1.nodeAddress,
        "rawProof": [],
        "block": offender1Block,
        "epoch": 0,
        "reportingBlock": 0, // does not matter
        "messageHash": 0, // must match accusation's one
      }
      await truffleAssert.fails(
        accountability.handleValidInnocenceProof(proof2),
        truffleAssert.ErrorType.REVERT,
        "no associated accusation",
      );
    });
  });

  describe('events', function () {
    beforeEach(async function () {
      autonity = await Autonity.new(validators, autonityConfig, {from: deployer});
      await autonity.finalizeInitialization({from: deployer});
      accountability = await AccountabilityTest.new(autonity.address, accountabilityConfig, {from: deployer});
      await autonity.setAccountabilityContract(accountability.address, {from:operator});
    });

    it("non-validator cannot submit event", async function () {
      let reporter = anyAccount;
      let offender = validators[1].nodeAddress;
      let PNrule = 0
      const event = {
        "chunks": 1,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter,
        "offender": offender,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0, 
      }
      await truffleAssert.fails(
        accountability.handleEvent(event, {from: reporter}),
        truffleAssert.ErrorType.REVERT,
        "validator not registered"
      );
    });

    it("handles chunked events", async function () {
      let reporter = validators[1].nodeAddress;
      let offender = validators[0].nodeAddress;
      let reporterPrivateKey = genesisPrivateKeys[1];
      let PNrule = 0;
      let event = {
        "chunks": 4,
        "chunkId": 2,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter,
        "offender": offender,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0,
      };
      let balance = web3.utils.toWei("10", "ether");
      await web3.eth.sendTransaction({from: validators[0].treasury, to: reporter, value: balance});
      
      // cannot submit transaction from reporter because the address is not unlocked and will require signing
      // however sendSignedTransaction method returns general error message instead of detailed error message
      // using call is similar to sending transaction but it will always revert, so does not require signing
      await truffleAssert.fails(
        accountability.handleEvent.call(event, {from: reporter}),
        truffleAssert.ErrorType.REVERT,
        "chunks must be contiguous"
      );
      
      let eventCount = event.chunks - 1;
      let currentProof = "0x";
      for (let i = 0; i < eventCount; i++) {
        let rawProof = [];
        rawProof.push(i);
        event.chunkId = i;
        event.rawProof = rawProof;
        let request = (await accountability.handleEvent.request(event, {from: reporter}));
        let receipt = await utils.signAndSendTransaction(reporter, accountability.address, reporterPrivateKey, request);
        assert.equal(receipt.status, true, "transaction failed");
        let currentEvent = await accountability.getReporterChunksMap({from: reporter});
        let hexNumber = (i > 15) ? i.toString(16) : "0" + i.toString(16);
        currentProof = currentProof + hexNumber;
        checkEvent(currentEvent, offender, reporter, i, currentProof);
      }

      // the error prooves that it is the last call and is ready to process
      event.chunkId = eventCount;
      await truffleAssert.fails(
        accountability.handleEvent.call(event, {from: reporter}),
        truffleAssert.ErrorType.REVERT,
        "failed proof verification"
      );
    });

    it("cannot submit event for another reporter", async function () {
      let reporter = validators[0].nodeAddress;
      let offender = validators[1].nodeAddress;
      let PNrule = 0;
      let event = {
        "chunks": 4,
        "chunkId": 2,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter,
        "offender": offender,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0,
      };

      // cannot submit transaction from reporter because the address is not unlocked and will require signing
      // however sendSignedTransaction method returns general error message instead of detailed error message
      // using call is similar to sending transaction but it will always revert, so does not require signing
      await truffleAssert.fails(
        accountability.handleEvent.call(event, {from: offender}),
        truffleAssert.ErrorType.REVERT,
        "event reporter must be caller"
      );
    });

    it("can reset event", async function () {
      let reporter = validators[0].nodeAddress;
      let offender = validators[1].nodeAddress;
      let reporterPrivateKey = genesisPrivateKeys[0];
      let balance = web3.utils.toWei("10", "ether");
      await web3.eth.sendTransaction({from: validators[0].treasury, to: reporter, value: balance});
      let PNrule = 0;
      let event = {
        "chunks": 4,
        "chunkId": 0,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter,
        "offender": offender,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0,
      };
      let rawProof = [];
      rawProof.push(20);
      event.rawProof = rawProof;

      let request = (await accountability.handleEvent.request(event, {from: reporter}));
      let receipt = await utils.signAndSendTransaction(reporter, accountability.address, reporterPrivateKey, request);
      assert.equal(receipt.status, true, "transaction failed");
      
      event.chunkId = 1;
      request = (await accountability.handleEvent.request(event, {from: reporter}));
      receipt = await utils.signAndSendTransaction(reporter, accountability.address, reporterPrivateKey, request);
      assert.equal(receipt.status, true, "transaction failed");

      let currentEvent = await accountability.getReporterChunksMap({from: reporter});
      let hexProof = "0x" + rawProof[0].toString(16) + rawProof[0].toString(16);
      checkEvent(currentEvent, offender, reporter, event.chunkId, hexProof);

      // reset
      event.chunkId = 0;
      request = (await accountability.handleEvent.request(event, {from: reporter}));
      receipt = await utils.signAndSendTransaction(reporter, accountability.address, reporterPrivateKey, request);
      assert.equal(receipt.status, true, "transaction failed");

      currentEvent = await accountability.getReporterChunksMap({from: reporter});
      hexProof = "0x" + rawProof[0].toString(16);
      checkEvent(currentEvent, offender, reporter, event.chunkId, hexProof);
    });

    it("sends chunked events from multiple validator", async function () {
      let reporter = [];
      let reporterPrivateKey = [];
      let rawProof = [];
      reporter.push(validators[0].nodeAddress);
      reporter.push(validators[1].nodeAddress);
      reporterPrivateKey.push(genesisPrivateKeys[0]);
      reporterPrivateKey.push(genesisPrivateKeys[1]);
      for (let i = 0; i < reporter.length; i++) {
        rawProof.push([]);
      }

      let offender = reporter;

      let balance = web3.utils.toWei("10", "ether");
      for (let i = 0; i < reporter.length; i++) {
        await web3.eth.sendTransaction({from: validators[0].treasury, to: reporter[i], value: balance});
      }

      let PNrule = 0;
      let event = {
        "chunks": 3,
        "chunkId": 0,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter,
        "offender": offender,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0,
      };

      let currentProof = "0x";
      let eventCount = event.chunks - 1;
      let count = 0;
      for (let chunkId = 0; chunkId < eventCount; chunkId++) {
        event.chunkId = chunkId;
        // reporter[i] sends event for offender[i];
        for (let i = 0; i < reporter.length; i++) {
          let sender = reporter[i];
          event.rawProof = [];
          event.rawProof.push(count);
          event.reporter = sender;
          event.offender = offender[i];
          let request = (await accountability.handleEvent.request(event, {from: sender}));
          let receipt = await utils.signAndSendTransaction(sender, accountability.address, reporterPrivateKey[i], request);
          assert.equal(receipt.status, true, "transaction failed");
          let currentEvent = await accountability.getReporterChunksMap({from: sender});
          rawProof[i].push(count);
          checkEvent(currentEvent, offender[i], sender, chunkId, utils.bytesToHex(rawProof[i]));
          count++;
        }
      }

      event.chunkId = eventCount;
      for (let i = 0; i < reporter.length; i++) {
        event.reporter = reporter[i];
        event.offender = offender[i];
        await truffleAssert.fails(
          accountability.handleEvent.call(event, {from: reporter[i]}),
          truffleAssert.ErrorType.REVERT,
          "failed proof verification"
        );
      }

    });

    it.skip("Can send event chunks with chunkId = 1", async function () {
      // right now it fails in case chunked events start with chunkId = 1
      // see https://github.com/autonity/autonity/issues/840
      let reporter = validators[0].nodeAddress;
      let offender = validators[1].nodeAddress;
      let reporterPrivateKey = genesisPrivateKeys[0];
      let balance = web3.utils.toWei("10", "ether");
      await web3.eth.sendTransaction({from: validators[0].treasury, to: reporter, value: balance});
      let PNrule = 0;
      let event = {
        "chunks": 4,
        "chunkId": 1,
        "eventType": 0,
        "rule": PNrule,
        "reporter": reporter,
        "offender": offender,
        "rawProof": [],
        "block": 10,
        "epoch": 0,
        "reportingBlock": 11,
        "messageHash": 0,
      };

      let request = (await accountability.handleEvent.request(event, {from: reporter}));
      let receipt = await utils.signAndSendTransaction(reporter, accountability.address, reporterPrivateKey, request);
      assert.equal(receipt.status, true, "transaction failed");
      let currentEvent = await accountability.getReporterChunksMap({from: reporter});
      // checkEvent(currentEvent, offender, reporter, event.chunkId, "0x");
      checkEvent(currentEvent, zeroAddress, zeroAddress, event.chunkId, "0x");
    });
  });

});