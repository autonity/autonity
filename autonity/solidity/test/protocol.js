'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const liquidContract = artifacts.require("Liquid")
const AccountabilityTest = artifacts.require("AccountabilityTest")
const config = require("./config");

// testing protocol contracts interactions.

const ValidatorState = {
  active : 0,
  paused : 1,
  jailed : 2
}


async function modifiedSlashingFeeAccountability(autonity, accountabilityConfig, operator, deployer) {
  let config = JSON.parse(JSON.stringify(accountabilityConfig));
  // so that we don't encounter error due to fraction and we don't do 100% slashing
  config.collusionFactor = 0,
  config.historyFactor = 0;
  let accountability = await AccountabilityTest.new(autonity.address, config, {from: deployer});
  await autonity.setAccountabilityContract(accountability.address, {from:operator});
  return accountability;
}

async function slash(config, accountability, epochOffenceCount, offender, reporter) {
  const event = {
    "chunks": 1, 
    "chunkId": 1,
    "eventType": 0,
    "rule": 0, // PN rule --> severity mid
    "reporter": reporter,
    "offender": offender,
    "rawProof": [], 
    "block": 1,
    "epoch": 0,
    "reportingBlock": 2,
    "messageHash": 0, 
  }
  let tx = await accountability.slash(event, epochOffenceCount);
  let txEvent;
  truffleAssert.eventEmitted(tx, 'SlashingEvent', (ev) => {
    txEvent = ev;
    return ev.amount.toNumber() > 0;
  });
  let slashingRate = utils.ruleToRate(config, event.rule) / config.slashingRatePrecision;
  return {txEvent, slashingRate};
}

// only self-unbond
async function selfUnbondAndSlash(config, autonity, accountability, delegator, validator, tokenUnbond, count, operator, deployer) {
  const initBalance = (await autonity.balanceOf(delegator)).toNumber();
  let valInfo = await autonity.getValidator(validator);
  let tokenUnbondArray = [];
  let requestID = (await autonity.getHeadUnbondingID()).toNumber();
  for (let i = 0; i < count; i++) {
    tokenUnbondArray.push((1+i)*tokenUnbond);
    await autonity.unbond(validator, tokenUnbondArray[i], {from: delegator});
  }
  // let unbonding apply
  await utils.endEpoch(autonity, operator, deployer);

  valInfo = await autonity.getValidator(validator);
  let totalCurrentStake = Number(valInfo.bondedStake) + Number(valInfo.selfUnbondingStake);
  let selfUnbondingShares = Number(valInfo.selfUnbondingShares);
  let selfUnbondingStake = Number(valInfo.selfUnbondingStake);
  let requests = [];
  for (let i = 0; i < count; i++) {
    let request = await autonity.getUnbondingRequest(requestID);
    requests.push(request);
    let share = tokenUnbondArray[i] * selfUnbondingShares / selfUnbondingStake;
    checkUnbondingRequest(request, tokenUnbondArray[i], share, true);
    requestID++;
  }

  // slash
  let {txEvent, slashingRate} = await slash(config, accountability, 1, validator, validator);
  valInfo = await autonity.getValidator(validator);
  assert.equal(
    Number(valInfo.bondedStake) + Number(valInfo.selfUnbondingStake),
    totalCurrentStake - txEvent.amount.toNumber(),
    "slashing amount does not match"
  );
  assert.equal(txEvent.amount.toNumber(), totalCurrentStake * slashingRate, "unexpected slashing");
  let selfUnbondingStakeAfterSlash = Number(valInfo.selfUnbondingStake);
  assert(selfUnbondingStakeAfterSlash > 0, "slashing all selfUnbondingStake does not work well in this case");
  await utils.mineTillUnbondingRelease(autonity, operator, deployer, false);
  // release NTN
  await utils.endEpoch(autonity, operator, deployer);
  let balanceIncrease = 0;
  let factor = selfUnbondingStakeAfterSlash / selfUnbondingStake;
  for (let i = 0; i < count; i++) {
    let balanceIncreaseFraction = selfUnbondingStakeAfterSlash * Number(requests[i].unbondingShare) / selfUnbondingShares;
    let expectedIncrease = selfUnbondingStake * Number(requests[i].unbondingShare) / selfUnbondingShares;
    assert.equal(balanceIncreaseFraction, expectedIncrease * factor, "unexpected slashing");
    balanceIncrease += balanceIncreaseFraction;
  }
  assert.equal((await autonity.balanceOf(delegator)).toNumber(), initBalance + balanceIncrease, "incorrect balance");
}

// only non-self-unbond
// make sure selfBondedStake = 0 so that slashing can be applied to unbondingStake and delegatedStake
async function unbondAndSlash(config, autonity, accountability, delegators, validator, tokenUnbond, operator, deployer, slashCount) {
  // not applicable for 100% slash
  let balances = [];
  let tokenUnbondArray = [];
  for (let i = 0; i < delegators.length; i++) {
    balances.push((await autonity.balanceOf(delegators[i])).toNumber());
    tokenUnbondArray.push((i+1)*tokenUnbond);
  }
  let valInfo = await autonity.getValidator(validator);
  let liquidSupply = Number(valInfo.liquidSupply);
  let delegatedStakes = Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake);
  let totalCurrentStake = delegatedStakes;
  let delegatee = [];
  delegatee.push(validator);
  let requestID = (await autonity.getHeadUnbondingID()).toNumber();
  await utils.bulkUnbondingRequest(autonity, delegators, delegatee, tokenUnbondArray);
  // let unbonding apply
  await utils.endEpoch(autonity, operator, deployer);
  valInfo = await autonity.getValidator(validator);
  let unbondingStakes = Number(valInfo.unbondingStake);
  let unbondingShares = Number(valInfo.unbondingShares);
  let requests = [];
  for (let i = 0; i < delegators.length; i++) {
    let request = await autonity.getUnbondingRequest(requestID);
    requests.push(request);
    let newtonAmount = tokenUnbondArray[i] * delegatedStakes / liquidSupply;
    let share = newtonAmount * unbondingShares / unbondingStakes;
    checkUnbondingRequest(request, tokenUnbondArray[i], share, false);
    requestID++;
  }
  let {txEvent, slashingRate} = await slash(config, accountability, 1, validator, validator);
  slashCount++;
  valInfo = await autonity.getValidator(validator);
  assert.equal(
    Number(valInfo.bondedStake) + Number(valInfo.unbondingStake),
    totalCurrentStake - txEvent.amount.toNumber(),
    "slashing amount does not match"
  );
  // conversion ratio chaned due to slashing
  liquidSupply = Number(valInfo.liquidSupply);
  delegatedStakes = Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake);
  unbondingStakes = Number(valInfo.unbondingStake);
  unbondingShares = Number(valInfo.unbondingShares);
  assert(delegatedStakes > 0, "100% slashing");
  await utils.mineTillUnbondingRelease(autonity, operator, deployer, false);
  // release NTN
  await utils.endEpoch(autonity, operator, deployer);
  let factor = 1;
  // previous slashing affects the delegated stake and LNTN:NTN ratio
  // so we need to take all slashing into account to compare NTN release with expected amount (NTN without slashing)
  while (slashCount > 0) {
    factor = factor * (1 - slashingRate);
    slashCount--;
  }
  for (let i = 0; i < delegators.length; i++) {
    let balanceIncrease = unbondingStakes * Number(requests[i].unbondingShare) / unbondingShares;
    assert.equal(balanceIncrease, tokenUnbondArray[i] * factor, "unexpected slashing");
    assert.equal((await autonity.balanceOf(delegators[i])).toNumber(), balances[i] + balanceIncrease, "unexpected balance");
  }

}


// only non-self-unbond
// make sure selfBondedStake = 0 so that slashing can be applied to delegatedStake
async function bondSlashUnbond(config, autonity, accountability, delegators, validator, tokenBond, tokenUnbond, operator, deployer) {
  // not applicable for 100% slash
  let valInfo = await autonity.getValidator(validator);
  const valLiquid = await liquidContract.at(valInfo.liquidContract);
  let balances = [];
  let tokenBondArray = [];
  let tokenUnbondArray = [];
  for (let i = 0; i < delegators.length; i++) {
    balances.push((await valLiquid.balanceOf(delegators[i])).toNumber());
    tokenBondArray.push((i+1)*tokenBond);
    tokenUnbondArray.push((i+1)*tokenUnbond);
  }
  let delegatee = [];
  delegatee.push(validator);
  await utils.bulkBondingRequest(autonity, operator, delegators, delegatee, tokenBondArray);
  // let bonding apply
  await utils.endEpoch(autonity, operator, deployer);
  valInfo = await autonity.getValidator(validator);
  let liquidSupply = Number(valInfo.liquidSupply);
  let delegatedStakes = Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake);
  // LNTN minted
  for (let i = 0; i < delegators.length; i++) {
    let liquidAmount = tokenBondArray[i] * liquidSupply / delegatedStakes;
    assert.equal((await valLiquid.balanceOf(delegators[i])).toNumber(), balances[i] + liquidAmount, "unexpected LNTN balance");
  }
  let totalCurrentStake = delegatedStakes + Number(valInfo.unbondingStake);

  // to compare with expected NTN without slashing, need to store old ratio
  const oldDelegatedStakes = delegatedStakes;
  let {txEvent, slashingRate} = await slash(config, accountability, 1, validator, validator);
  valInfo = await autonity.getValidator(validator);
  assert.equal(
    Number(valInfo.bondedStake) + Number(valInfo.unbondingStake),
    totalCurrentStake - txEvent.amount.toNumber(),
    "slashing amount does not match"
  );
  // conversion ratio chaned due to slashing
  delegatedStakes = Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake);
  assert(delegatedStakes > 0, "100% slashing");
  let requestID = (await autonity.getHeadUnbondingID()).toNumber();
  await utils.bulkUnbondingRequest(autonity, delegators, delegatee, tokenUnbondArray);
  // let unbonding apply
  await utils.endEpoch(autonity, operator, deployer);
  valInfo = await autonity.getValidator(validator);
  let unbondingStakes = Number(valInfo.unbondingStake);
  let unbondingShares = Number(valInfo.unbondingShares);
  for (let i = 0; i < delegators.length; i++) {
    let request = await autonity.getUnbondingRequest(requestID);
    let newtonAmount = tokenUnbondArray[i] * delegatedStakes / liquidSupply;
    let expectedNewton = tokenUnbondArray[i] * oldDelegatedStakes / liquidSupply;
    assert.equal(newtonAmount, expectedNewton * (1 - slashingRate), "unexpected NTN conversion");
    let share = newtonAmount * unbondingShares / unbondingStakes;
    checkUnbondingRequest(request, tokenUnbondArray[i], share, false);
    requestID++;
  }
  // release NTN so slashing affects only LNTN:NTN
  await utils.mineTillUnbondingRelease(autonity, operator, deployer);
  return txEvent.releaseBlock.toNumber();
}


function checkUnbondingRequest(request, newtonAmount, share, selfDelegation) {
  assert.equal(request.amount, newtonAmount, "unexpected unbonding amount");
  assert.equal(request.unbondingShare, share, "unexpected unbonding share");
  assert.equal(request.selfDelegation, selfDelegation, "unexpected self delegation");
  assert.equal(request.unlocked, true, "unexpected unbondingRequest.unlocked");
}

contract('Protocol', function (accounts) {
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
  const treasuryAccount = accounts[8];

  let autonityConfig = config.autonityConfig(operator, treasuryAccount)
  const accountabilityConfig = {
    "innocenceProofSubmissionWindow": 30,
    "latestAccountabilityEventsRange": 256,
    "baseSlashingRateLow": 500,
    "baseSlashingRateMid": 1000,
    "collusionFactor": 550,
    "historyFactor": 750,
    "jailFactor": 1,
    "slashingRatePrecision": 10000
  }

  const genesisNodeAddresses = config.GENESIS_NODE_ADDRESSES

  // accounts[2] is skipped because it is used as a genesis validator when running against autonity
  // this can cause interference in reward distribution tests
  const validators = config.validators(accounts)

  let autonity;
  let accountability;
  describe('After effects of slashing', function () {
    beforeEach(async function () {
      autonity = await utils.deployAutonityTestContract(validators, autonityConfig, accountabilityConfig,  deployer, operator);
      accountability = await AccountabilityTest.new(autonity.address, accountabilityConfig, {from: deployer});
      await autonity.setAccountabilityContract(accountability.address, {from:operator});
    });

    it('unbondingShares:unbondingStake conversion ratio', async function () {
      // issue multiple unbonding requests (non-selfBonded) in different epochs, interleaved with slashing events
      // and check that the unbonding shares related fields change accordingly


      accountability = await modifiedSlashingFeeAccountability(autonity, accountabilityConfig, operator, deployer);
      
      let delegatee = [];
      let delegators = [];
      let tokenBondArray = []
      const maxCount = 3;
      const tokenBond = 100000000;
      const tokenUnbond = 1000;

      for (let i = 0; i < Math.min(validators.length, maxCount); i++) {
        delegators.push(validators[i].treasury);
        tokenBondArray.push(tokenBond);
      }
      let validator = validators[maxCount].nodeAddress;
      let treasury = validators[maxCount].treasury;
      delegatee.push(validator);
      let valInfo = await autonity.getValidator(validator);
      // so slashing applies to delegated stake and unbonding stake
      await autonity.unbond(validator, Number(valInfo.selfBondedStake), {from: treasury});
      await utils.bulkBondingRequest(autonity, operator, delegators, delegatee, tokenBondArray);

      // mine blocks until unbonding period is reached
      await utils.mineTillUnbondingRelease(autonity, operator, deployer);
      // requests will be processed at epoch end
      await utils.endEpoch(autonity, operator, deployer);
      // request unbonding and slash
      await unbondAndSlash(accountabilityConfig, autonity, accountability, delegators, validator, tokenUnbond, operator, deployer, 0);
      // repeat
      await unbondAndSlash(accountabilityConfig, autonity, accountability, delegators, validator, tokenUnbond, operator, deployer, 1);
    });


    it('selfUnbondingShares:selfUnbondingStake conversion ratio', async function () {
      // issue multiple unbonding requests (selfBonded) in different epochs, interleaved with slashing events
      // and check that the unbonding shares related fields change accordingly


      accountability = await modifiedSlashingFeeAccountability(autonity, accountabilityConfig, operator, deployer);
      
      const validator = validators[0].nodeAddress;
      const delegator = validators[0].treasury;
      const tokenBond = 100000000 - validators[0].bondedStake;
      const maxCount = 4;

      await autonity.mint(delegator, tokenBond, {from: operator});
      await autonity.bond(validator, tokenBond, {from: delegator});
      await utils.endEpoch(autonity, operator, deployer);
      let valInfo = await autonity.getValidator(validator);
      assert.equal(Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake), 0, "delegated stake exists");
      let totalSelfBonded = Number(valInfo.bondedStake);
      let tokenUnbond = totalSelfBonded * 2 / 100;

      await selfUnbondAndSlash(accountabilityConfig, autonity, accountability, delegator, validator, tokenUnbond, maxCount, operator, deployer);
      
      // repeat
      valInfo = await autonity.getValidator(validator);
      totalSelfBonded = Number(valInfo.bondedStake);
      tokenUnbond = totalSelfBonded * 2 / 100;
      await selfUnbondAndSlash(accountabilityConfig, autonity, accountability, delegator, validator, tokenUnbond, maxCount, operator, deployer);

    });

    it.skip('unbondingShares:unbondingStake 100% slash edge case', async function () {
      // see https://github.com/autonity/autonity/issues/819
      // unbonding period needs to be increased for this test to work
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;
      const delegatorA = accounts[8];
      const delegatorB = accounts[9];
      const tokenMint = 100;
      
      let delegatorBalance_A = (await autonity.balanceOf(delegatorA)).toNumber();
      let delegatorBalance_B = (await autonity.balanceOf(delegatorB)).toNumber();
      await autonity.mint(delegatorA, tokenMint, {from: operator});
      await autonity.mint(delegatorB, tokenMint, {from: operator});

      // bond from delegatorA
      await autonity.bond(validator, tokenMint, {from: delegatorA});
      // let bonding apply
      await utils.endEpoch(autonity, operator, deployer);

      await autonity.unbond(validator, tokenMint, {from: delegatorA});
      let requestIDA = (await autonity.getHeadUnbondingID()).toNumber() - 1;
      await utils.endEpoch(autonity, operator, deployer);
      let unbondingRequestA = await autonity.getUnbondingRequest(requestIDA);
      assert.equal(unbondingRequestA.unbondingShare, tokenMint, "unexpected unbondingShare");

      // so that 100% slashing occurs
      // this is only to treasury and delegatorA
      let epochOffenceCount = accountabilityConfig.slashingRatePrecision;
      let valInfo = await autonity.getValidator(validator);
      let totalSlash = Number(valInfo.bondedStake) + Number(valInfo.selfUnbondingStake) + Number(valInfo.unbondingStake);
      let {txEvent, slashingRate} = await slash(accountabilityConfig, accountability, epochOffenceCount, validator, treasury);
      let releaseBlock = txEvent.releaseBlock.toNumber();
      assert.equal(txEvent.amount.toNumber(), totalSlash, "100% slash did not happen");
      valInfo = await autonity.getValidator(validator);
      let totalStake = Number(valInfo.bondedStake) + Number(valInfo.selfUnbondingStake) + Number(valInfo.unbondingStake);
      assert.equal(totalStake, 0, "stake remaining after 100% slash");

      while (await web3.eth.getBlockNumber() < releaseBlock) {
        utils.mineEmptyBlock();
      }
      await autonity.activateValidator(validator, {from: treasury});
      // let delegatorB bond
      // no slashing will occur from now, so delegatoB should get full NTN returned when unbonded
      await autonity.bond(validator, tokenMint, {from: delegatorB});

      // let bonding apply
      await utils.endEpoch(autonity, operator, deployer);
      const valLiquidContract = await liquidContract.at(valInfo.liquidContract);
      assert.equal((await valLiquidContract.balanceOf(delegatorB)).toNumber(), tokenMint);

      // delegatorB should get full NTN while delegatorA should get 0 NTN
      await autonity.unbond(validator, tokenMint, {from: delegatorB});
      let requestIDB = (await autonity.getHeadUnbondingID()).toNumber() - 1;
      // let unbonding apply
      await utils.endEpoch(autonity, operator, deployer);
      let unbondingRequestB = await autonity.getUnbondingRequest(requestIDB);
      assert.equal(unbondingRequestB.unbondingShare, tokenMint, "unexpected unbondingShare");
      // unbonding request of delegatorB
      releaseBlock = (await autonity.getUnbondingPeriod()).toNumber() + Number(unbondingRequestB.requestBlock);
      while (await web3.eth.getBlockNumber() < releaseBlock) {
        utils.mineEmptyBlock();
      }
      // unbonding is released
      valInfo = await autonity.getValidator(validator);
      await utils.endEpoch(autonity, operator, deployer);
      assert.equal((await autonity.balanceOf(delegatorB)).toNumber(), delegatorBalance_B + tokenMint, "unexpected balance");
      assert.equal((await autonity.balanceOf(delegatorA)).toNumber(), delegatorBalance_A, "unexpected balance");
    });

    it('LNTN:NTN conversion ratio', async function () {
      // issue multiple bond and unbond request with interleaved slashing events, and check that the NTN:LNTN ratio is always what we expect
      
      accountability = await modifiedSlashingFeeAccountability(autonity, accountabilityConfig, operator, deployer);

      let delegatee = [];
      let delegators = [];
      const maxCount = 3;

      for (let i = 0; i < Math.min(validators.length, maxCount); i++) {
        delegators.push(validators[i].treasury);
      }
      let validator = validators[maxCount].nodeAddress;
      let treasury = validators[maxCount].treasury;
      delegatee.push(validator);
      let valInfo = await autonity.getValidator(validator);
      // so slashing applies to delegated stake and unbonding stake
      await autonity.unbond(validator, Number(valInfo.selfBondedStake), {from: treasury});

      // mine blocks until unbonding period is reached
      await utils.mineTillUnbondingRelease(autonity, operator, deployer);
      const tokenBond = 1000;
      const tokenUnbond = 100;
      let roundingFactor = 10000;
      // request bonding and slash
      let releaseBlock = await bondSlashUnbond(accountabilityConfig, autonity, accountability, delegators, validator, tokenBond * roundingFactor, tokenUnbond, operator, deployer);
      while (await web3.eth.getBlockNumber() < releaseBlock) {
        await utils.mineEmptyBlock();
      }
      await autonity.activateValidator(validator, {from: treasury});
      // repeat
      let slashingRate = utils.ruleToRate(accountabilityConfig, 0) / accountabilityConfig.slashingRatePrecision; // rule 0 --> severity mid
      // multiplying tokenBond with roundingFactor so we don't get fraction ratio of LNTN:NTN or unbondingShare:unbondingStake
      roundingFactor = roundingFactor * (1 - slashingRate);
      await bondSlashUnbond(accountabilityConfig, autonity, accountability, delegators, validator, tokenBond * roundingFactor, tokenUnbond, operator, deployer);
    });

    it.skip('LNTN:NTN 100% slash edge case', async function () {
      // see https://github.com/autonity/autonity/issues/819
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;
      const delegatorA = accounts[8];
      const delegatorB = accounts[9];
      const tokenMint = 100;
      
      let delegatorBalance_A = (await autonity.balanceOf(delegatorA)).toNumber();
      let delegatorBalance_B = (await autonity.balanceOf(delegatorB)).toNumber();
      await autonity.mint(delegatorA, tokenMint, {from: operator});
      await autonity.mint(delegatorB, tokenMint, {from: operator});

      // bond from delegatorA
      await autonity.bond(validator, tokenMint, {from: delegatorA});
      // let bonding apply
      await utils.endEpoch(autonity, operator, deployer);

      // so that 100% slashing occurs
      // this is only to treasury and delegatorA
      let epochOffenceCount = accountabilityConfig.slashingRatePrecision;
      let valInfo = await autonity.getValidator(validator);
      let totalSlash = Number(valInfo.bondedStake) + Number(valInfo.selfUnbondingStake) + Number(valInfo.unbondingStake);
      let {txEvent, slashingRate} = await slash(accountabilityConfig, accountability, epochOffenceCount, validator, treasury);
      let releaseBlock = txEvent.releaseBlock.toNumber();
      assert.equal(txEvent.amount.toNumber(), totalSlash, "100% slash did not happen");
      valInfo = await autonity.getValidator(validator);
      let totalStake = Number(valInfo.bondedStake) + Number(valInfo.selfUnbondingStake) + Number(valInfo.unbondingStake);
      assert.equal(totalStake, 0, "stake remaining after 100% slash");

      while (await web3.eth.getBlockNumber() < releaseBlock) {
        utils.mineEmptyBlock();
      }
      await autonity.activateValidator(validator, {from: treasury});
      // let delegatorB bond
      // no slashing will occur from now, so delegatoB should get full NTN returned when unbonded
      await autonity.bond(validator, tokenMint, {from: delegatorB});

      // let bonding apply
      await utils.endEpoch(autonity, operator, deployer);
      const valLiquidContract = await liquidContract.at(valInfo.liquidContract);
      assert.equal((await valLiquidContract.balanceOf(delegatorB)).toNumber(), tokenMint);

      // depending on the soln of the problem, the following portion of code might change
      // delegatorB should get full NTN while delegatorA should get 0 NTN
      await autonity.unbond(validator, tokenMint, {from: delegatorA});
      let requestIDA = (await autonity.getHeadUnbondingID()).toNumber() - 1;
      await autonity.unbond(validator, tokenMint, {from: delegatorB});
      let requestIDB = (await autonity.getHeadUnbondingID()).toNumber() - 1;
      // let unbonding apply
      await utils.endEpoch(autonity, operator, deployer);
      let unbondingRequestA = await autonity.getUnbondingRequest(requestIDA);
      assert.equal(unbondingRequestA.unbondingShare, 0, "unexpected unbondingShare");
      let unbondingRequestB = await autonity.getUnbondingRequest(requestIDB);
      assert.equal(unbondingRequestB.unbondingShare, tokenMint, "unexpected unbondingShare");

    });

    it('jailed validator rewards go to proof reporter', async function () {
      // verify that if a validators is jailed (it has been slashed for a misbehaviour)
      // his share of rewards goes to the proof reporter.

      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;
      const reporter = validators[1].nodeAddress;
      const reporterTreasury = validators[1].treasury;

      let committee = await autonity.getCommittee();
      let totalBondedStake = 0;
      let validatorReward;
      let reporterReward;
      for (let i = 0; i < committee.length; i++) {
        totalBondedStake += Number(committee[i].votingPower);
        if (committee[i].addr == validator) {
          validatorReward = Number(committee[i].votingPower);
        } else if (committee[i].addr == reporter) {
          reporterReward = Number(committee[i].votingPower);
        }
      }

      let reward = totalBondedStake;
      // send funds to contract account, to get them distributed later on.
      await web3.eth.sendTransaction({from: anyAccount, to: autonity.address, value: reward});

      let epochOffenceCount = 1;
      await slash(accountabilityConfig, accountability, epochOffenceCount, validator, reporter)
      let treasuryBalance = await web3.eth.getBalance(treasury);
      let reporterTreasuryBalance = Number(await web3.eth.getBalance(reporterTreasury));
      await utils.endEpoch(autonity, operator, deployer);
      assert.equal(await web3.eth.getBalance(treasury), treasuryBalance, "jailed validator got reward");
      assert.equal(
        await web3.eth.getBalance(reporterTreasury),
        validatorReward + reporterReward + reporterTreasuryBalance,
        "reporter did not get reward from jailed validator"
      );
    });

    it('jailed validator cannot be activated', async function () {
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;

      let epochOffenceCount = 1;
      let {txEvent, slashingRate} = await slash(accountabilityConfig, accountability, epochOffenceCount, validator, treasury);
      let releaseBlock = txEvent.releaseBlock.toNumber();

      let validatorInfo = await autonity.getValidator(validator);
      assert.equal(validatorInfo.state, ValidatorState.jailed, "validator not jailed");
      await truffleAssert.fails(
        autonity.activateValidator(validator, {from: treasury}),
        truffleAssert.ErrorType.REVERT,
        "validator still in jail"
      );

      while (await web3.eth.getBlockNumber() < releaseBlock) {
        utils.mineEmptyBlock();
      }
      let tx = await autonity.activateValidator(validator, {from: treasury});
      truffleAssert.eventEmitted(tx, 'ActivatedValidator', (ev) => {
        return ev.treasury === treasury && ev.addr === validator;
      });

    });
  });
});
