'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const liquidContract = artifacts.require("Liquid")
const AccountabilityTest = artifacts.require("AccountabilityTest")
const config = require("./config");

// testing protocol contracts interactions.



async function modifiedSlashingFeeAccountability(autonity, accountabilityConfig, operator, deployer) {
  let config = JSON.parse(JSON.stringify(accountabilityConfig));
  // so that we don't encounter error due to fraction and we don't do 100% slashing
  config.collusionFactor = 0,
  config.historyFactor = 0;
  let accountability = await AccountabilityTest.new(autonity.address, config, {from: deployer});
  await autonity.setAccountabilityContract(accountability.address, {from:operator});
  return accountability;
}

async function killValidatorWithSlash(config, accountability, offender, reporter) {
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

  // high offence count for 100% slash
  let epochOffenceCount = config.slashingRatePrecision;
  let tx = await accountability.slash(event, epochOffenceCount);
  let txEvent;
  // validator needs to have non-self-bonding to be jailbound
  truffleAssert.eventEmitted(tx, 'SlashingEvent', (ev) => {
    txEvent = ev;
    return ev.amount.toNumber() > 0 && ev.isJailbound == true;
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
  let {txEvent, slashingRate} = await utils.slash(config, accountability, 1, validator, validator);
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
  let {txEvent, slashingRate} = await utils.slash(config, accountability, 1, validator, validator);
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
  let {txEvent, slashingRate} = await utils.slash(config, accountability, 1, validator, validator);
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
      await utils.mockPrecompile()
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

    it('unbondingShares:unbondingStake 100% slash edge case', async function () {
      // unbonding period needs to be increased for this test to work
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;
      const delegator = accounts[9];
      const tokenMint = 100;
      await autonity.mint(delegator, tokenMint, {from: operator});
      await autonity.bond(validator, tokenMint, {from: delegator});
      // let bonding apply
      await utils.endEpoch(autonity, operator, deployer);

      let balance = (await autonity.balanceOf(delegator)).toNumber();
      await autonity.unbond(validator, tokenMint, {from: delegator});
      let requestID = (await autonity.getHeadUnbondingID()).toNumber() - 1;
      // let unbonding apply
      await utils.endEpoch(autonity, operator, deployer);
      let unbondingRequest = await autonity.getUnbondingRequest(requestID);
      assert.equal(unbondingRequest.unbondingShare, tokenMint, "unexpected unbondingShare");

      await killValidatorWithSlash(accountabilityConfig, accountability, validator, treasury);
      await utils.mineTillUnbondingRelease(autonity, operator, deployer, false);
      await utils.endEpoch(autonity, operator, deployer);
      assert.equal((await autonity.balanceOf(delegator)).toNumber(), balance, "balance increased after 100% slash");
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

    it('LNTN:NTN 100% slash edge case', async function () {
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;
      const delegator = accounts[8];
      const tokenMint = 100;
      await autonity.mint(delegator, tokenMint, {from: operator});
      await autonity.bond(validator, tokenMint, {from: delegator});
      // let bonding apply
      await utils.endEpoch(autonity, operator, deployer);

      await killValidatorWithSlash(accountabilityConfig, accountability, validator, treasury);
      let balance = (await autonity.balanceOf(delegator)).toNumber();
      await autonity.unbond(validator, tokenMint, {from: delegator});
      let requestID = (await autonity.getHeadUnbondingID()).toNumber() - 1;
      await utils.mineTillUnbondingRelease(autonity, operator, deployer);
      await utils.endEpoch(autonity, operator, deployer);
      let unbondingRequest = await autonity.getUnbondingRequest(requestID);
      assert.equal(unbondingRequest.unbondingShare, 0, "unexpected unbondingShare");
      assert.equal((await autonity.balanceOf(delegator)).toNumber(), balance, "balance increased after 100% slash");

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
      await utils.slash(accountabilityConfig, accountability, epochOffenceCount, validator, reporter)
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

    it('jailbound validator rewards go to proof reporter', async function () {
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;
      const reporter = validators[1].nodeAddress;
      const reporterTreasury = validators[1].treasury;

      // non-self bond
      const delegator = accounts[9];
      const tokenMint = 100;
      await autonity.mint(delegator, tokenMint, {from: operator});
      await autonity.bond(validator, tokenMint, {from: delegator});
      await utils.endEpoch(autonity, operator, deployer);

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

      await killValidatorWithSlash(accountabilityConfig, accountability, validator, reporter)
      let treasuryBalance = await web3.eth.getBalance(treasury);
      let reporterTreasuryBalance = Number(await web3.eth.getBalance(reporterTreasury));
      await utils.endEpoch(autonity, operator, deployer);
      assert.equal(await web3.eth.getBalance(treasury), treasuryBalance, "jailbound validator got reward");
      assert.equal(
        await web3.eth.getBalance(reporterTreasury),
        validatorReward + reporterReward + reporterTreasuryBalance,
        "reporter did not get reward from jailbound validator"
      );
    });

    it('jailed validator cannot be activated', async function () {
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;

      let epochOffenceCount = 1;
      let {txEvent, slashingRate} = await utils.slash(accountabilityConfig, accountability, epochOffenceCount, validator, treasury);
      let releaseBlock = txEvent.releaseBlock.toNumber();

      let validatorInfo = await autonity.getValidator(validator);
      assert.equal(validatorInfo.state, utils.ValidatorState.jailed, "validator not jailed");
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

    it('cannot bond to a non-active validator', async function () {
      // ready task to jail, jailbound or pause a validator
      let tasks = [];
      let status = [];
      let jailTask = async function(validator, reporter) {
        let {txEvent, _} = await utils.slash(accountabilityConfig, accountability, 1, validator, reporter);
        assert.equal(txEvent.isJailbound, false, "slashed too much, validator jailbound instead of jailed");
      }
      tasks.push(jailTask);
      status.push(utils.ValidatorState.jailed);
      let jailboundTask = async function(validator, reporter) {
        await killValidatorWithSlash(accountabilityConfig, accountability, validator, reporter);
      }
      tasks.push(jailboundTask);
      status.push(utils.ValidatorState.jailbound);
      let pauseTask = async function(validator, treasury) {
        let tx = await autonity.pauseValidator(validator, {from: treasury});
        truffleAssert.eventEmitted(tx, 'PausedValidator', (ev) => {
          return ev.treasury == treasury && ev.addr == validator;
        });
      }
      tasks.push(pauseTask);
      status.push(utils.ValidatorState.paused);
      const count = tasks.length;

      const delegator = accounts[9];
      const tokenMint = 100;

      await autonity.mint(delegator, (1+count)*tokenMint, {from: operator});
      let delegatorBalance = (await autonity.balanceOf(delegator)).toNumber();
      let treasuryBalances = []
      // mint and bond
      for (let i = 0; i < count; i++) {
        let treasury = validators[i].treasury;
        let validator = validators[i].nodeAddress;
        await autonity.mint(treasury, 2*tokenMint, {from: operator});
        treasuryBalances.push((await autonity.balanceOf(treasury)).toNumber());
        await autonity.bond(validator, tokenMint, {from: delegator});
        await autonity.bond(validator, tokenMint, {from: treasury});
        assert.equal((await autonity.balanceOf(delegator)).toNumber(), delegatorBalance - tokenMint, "delegator balance did not decrease after bonding request");
        assert.equal((await autonity.balanceOf(treasury)).toNumber(), treasuryBalances[i] - tokenMint, "treasury balance did not decrease after bonding request");
        delegatorBalance -= tokenMint;
      }

      // perform tasks before bonding can be applied
      let oldValInfo = [];
      for (let i = 0; i < count; i++) {
        let treasury = validators[i].treasury;
        let validator = validators[i].nodeAddress;
        await tasks[i](validator, treasury);
        oldValInfo.push(await autonity.getValidator(validator));
        assert.equal(oldValInfo[i].state, status[i], "validator status mismatch");
        // cannot bond to jailed validator
        await truffleAssert.fails(
          autonity.bond(validator, tokenMint, {from: delegator}),
          truffleAssert.ErrorType.REVERT,
          "validator need to be active"
        );
        await truffleAssert.fails(
          autonity.bond(validator, tokenMint, {from: treasury}),
          truffleAssert.ErrorType.REVERT,
          "validator need to be active"
        );
      }

      // bonding request should be rejected
      await utils.endEpoch(autonity, operator, deployer);
      delegatorBalance += count*tokenMint;
      assert.equal((await autonity.balanceOf(delegator)).toNumber(), delegatorBalance, "unexpected delegator balance");
      for (let i = 0; i < count; i++) {
        let validator = validators[i].nodeAddress;
        let treasury = validators[i].treasury;
        assert.equal((await autonity.balanceOf(treasury)).toNumber(), treasuryBalances[i], "unexpected treasury balance");
        let newValInfo = await autonity.getValidator(validator);
        assert.equal(newValInfo.bondedStake, oldValInfo[i].bondedStake, "bondedStake changed");
        assert.equal(newValInfo.selfBondedStake, oldValInfo[i].selfBondedStake, "selfBondedStake changed");
        assert.equal(newValInfo.selfUnbondingStake, oldValInfo[i].selfUnbondingStake, "selfUnbondingStake changed");
        assert.equal(newValInfo.unbondingStake, oldValInfo[i].unbondingStake, "unbondingStake changed");
        assert.equal(newValInfo.liquidSupply, oldValInfo[i].liquidSupply, "liquidSupply changed");
        assert.equal(newValInfo.state, oldValInfo[i].state, "validator status mismatch");
      }
    });

    it('jailbound validator cannot be activated', async function () {
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;

      // non-self bond
      let newAccount = accounts[8];
      let tokenMint = 200;
      await autonity.mint(newAccount, tokenMint, {from: operator});
      await autonity.bond(validator, tokenMint, {from: newAccount});
      await utils.endEpoch(autonity, operator, deployer);

      await killValidatorWithSlash(accountabilityConfig, accountability, validator, treasury);

      let validatorInfo = await autonity.getValidator(validator);
      assert.equal(validatorInfo.state, utils.ValidatorState.jailbound, "validator not jailbound");
      await truffleAssert.fails(
        autonity.activateValidator(validator, {from: treasury}),
        truffleAssert.ErrorType.REVERT,
        "validator jailed permanently"
      );

      let releaseBlock = validatorInfo.jailReleaseBlock;
      assert.equal(releaseBlock, 0, "releaseBlock for jailbound validator");

    });

    it('kills validator for 100% slash', async function () {
      // non-self-bond and self-bond
      const validatorAddresses = [validators[0].nodeAddress, validators[1].nodeAddress];
      const delegatorAddresses = [accounts[9], validators[1].treasury];
      const tokenMint = 100;
      let balances = [];

      for (let iter = 0; iter < delegatorAddresses.length; iter++) {
        const delegator = delegatorAddresses[iter];
        const validator = validatorAddresses[iter];
        balances.push((await autonity.balanceOf(delegator)).toNumber());
        await autonity.mint(delegator, tokenMint, {from: operator});
        await autonity.bond(validator, tokenMint, {from: delegator});
      }
      // let bonding apply
      await utils.endEpoch(autonity, operator, deployer);

      for (let iter = 0; iter < delegatorAddresses.length; iter++) {
        const delegator = delegatorAddresses[iter];
        const validator = validatorAddresses[iter];
        await autonity.unbond(validator, tokenMint/2, {from: delegator});
      }
      // let unbonding apply
      await utils.endEpoch(autonity, operator, deployer);

      for (let iter = 0; iter < delegatorAddresses.length; iter++) {
        const delegator = delegatorAddresses[iter];
        const validator = validatorAddresses[iter];
        await killValidatorWithSlash(accountabilityConfig, accountability, validator, delegator)
        let valInfo = await autonity.getValidator(validator);
        assert.equal(valInfo.state, utils.ValidatorState.jailbound, "validator not jailbound");
        assert.equal(
          parseInt(valInfo.bondedStake) + parseInt(valInfo.unbondingStake) + parseInt(valInfo.selfUnbondingStake)
          , 0, "100% slash did not happen"
        );
        await utils.mineTillUnbondingRelease(autonity, operator, deployer);
        assert.equal((await autonity.balanceOf(delegator)).toNumber(), balances[iter], "unbonding released");
      }
      
    });
  });
});
