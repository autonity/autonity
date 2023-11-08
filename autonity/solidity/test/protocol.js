'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const liquidContract = artifacts.require("Liquid")
const AccountabilityTest = artifacts.require("AccountabilityTest")
const toBN = web3.utils.toBN;

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
  // validator needs to have non-self-bonding to be killed
  truffleAssert.eventEmitted(tx, 'ValidatorKilled', (ev) => {
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
    "jailFactor": 1,
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
    "selfUnbondingStakeLocked": 0,
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

    it('killed validator rewards go to proof reporter', async function () {
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
      assert.equal(await web3.eth.getBalance(treasury), treasuryBalance, "killed validator got reward");
      assert.equal(
        await web3.eth.getBalance(reporterTreasury),
        validatorReward + reporterReward + reporterTreasuryBalance,
        "reporter did not get reward from killed validator"
      );
    });

    it('jailed validator cannot be activated', async function () {
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;

      let epochOffenceCount = 1;
      let {txEvent, slashingRate} = await slash(accountabilityConfig, accountability, epochOffenceCount, validator, treasury);
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

    it('killed validator cannot be activated', async function () {
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
      assert.equal(validatorInfo.state, utils.ValidatorState.killed, "validator not killed");
      await truffleAssert.fails(
        autonity.activateValidator(validator, {from: treasury}),
        truffleAssert.ErrorType.REVERT,
        "validator killed permanently"
      );

      let releaseBlock = validatorInfo.jailReleaseBlock;
      while (await web3.eth.getBlockNumber() < releaseBlock) {
        utils.mineEmptyBlock();
      }
      await truffleAssert.fails(
        autonity.activateValidator(validator, {from: treasury}),
        truffleAssert.ErrorType.REVERT,
        "validator killed permanently"
      );

    });

    it('cannot bond to a killed validator', async function () {
      let validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;

      // non-self bond
      let delegator = accounts[9];
      let tokenMint = 100;
      await autonity.mint(delegator, 3*tokenMint, {from: operator});
      // 1st bond
      // without non-self bonding, validator cannot be killed
      await autonity.bond(validator, tokenMint, {from: delegator});
      await utils.endEpoch(autonity, operator, deployer);

      let balance = (await autonity.balanceOf(delegator)).toNumber();
      // 2nd bond
      await autonity.bond(validator, tokenMint, {from: delegator});
      assert.equal((await autonity.balanceOf(delegator)).toNumber(), balance - tokenMint, "balance did not decrease after bonding request");
      await killValidatorWithSlash(accountabilityConfig, accountability, validator, treasury);

      await truffleAssert.fails(
        autonity.bond(validator, tokenMint, {from: delegator}),
        truffleAssert.ErrorType.REVERT,
        "validator need to be active"
      );
      // 2nd bonding should not be applied
      await utils.endEpoch(autonity, operator, deployer);
      assert.equal((await autonity.balanceOf(delegator)).toNumber(), balance, "unexpected balance");
    });

    it('zero amount of bonded and unbonding stake after 100% slash', async function () {
      const validator = validators[0].nodeAddress;
      const treasury = validators[0].treasury;
      const delegator = accounts[8];
      const tokenMint = 100;
      await autonity.mint(delegator, tokenMint, {from: operator});
      await autonity.bond(validator, tokenMint, {from: delegator});
      // let bonding apply
      await utils.endEpoch(autonity, operator, deployer);

      await autonity.unbond(validator, tokenMint/2, {from: delegator});
      let valInfo = await autonity.getValidator(validator);
      let totalSlash = Number(valInfo.bondedStake) + Number(valInfo.selfUnbondingStake) + Number(valInfo.unbondingStake);
      let {txEvent, slashingRate} = await killValidatorWithSlash(accountabilityConfig, accountability, validator, treasury);
      assert.equal(txEvent.amount.toNumber(), totalSlash, "100% slash did not happen");
      valInfo = await autonity.getValidator(validator);
      let totalStake = Number(valInfo.bondedStake) + Number(valInfo.selfUnbondingStake) + Number(valInfo.unbondingStake);
      assert.equal(totalStake, 0, "stake remaining after 100% slash");
    });

    it('does not trigger fairness issue (unbondingStake > 0 and delegatedStake > 0)', async function () {
      // fairness issue is triggered when delegatedStake or unbondingStake becomes 0 from positive due to slashing
      // it can happen due to slashing rate = 100% or (totalStake - 1)/totalStake
      // it can also happen if both unbondingStake and delegatedStake > 0 and slashing amount >= (totalStake - 1)
      // it should not happen for slashing amount < (totalStake - 1)
      let config = JSON.parse(JSON.stringify(accountabilityConfig));
      config.collusionFactor = parseInt(config.slashingRatePrecision) - parseInt(config.baseSlashingRateMid) - 2;
      accountability = await AccountabilityTest.new(autonity.address, config, {from: deployer});
      await autonity.setAccountabilityContract(accountability.address, {from:operator});

      const validatorAddresses = [ validators[0].nodeAddress, validators[1].nodeAddress];
      const tokenUnbondFactor = [1/10, 8/10];
      const delegator = accounts[8];

      for (let iter = 0; iter < 2; iter++) {
        let validator = validatorAddresses[iter];
        let valInfo = await autonity.getValidator(validator);
        let selfBondedStake = parseInt(valInfo.bondedStake);
        // non-self bond to check fairness issue
        const tokenMint = parseInt(config.slashingRatePrecision) - selfBondedStake;
        console.log(tokenMint);
        await autonity.mint(delegator, tokenMint, {from: operator});
        await autonity.bond(validator, tokenMint, {from: delegator});
        // let bonding apply
        await utils.endEpoch(autonity, operator, deployer);
        let totalStake = selfBondedStake + tokenMint;
        let tokenUnBond = Math.ceil(tokenMint*tokenUnbondFactor[iter]);
        console.log(tokenUnBond);
        await autonity.unbond(validator, tokenUnBond, {from: delegator});
        await utils.endEpoch(autonity, operator, deployer);
        const event = {
          "chunks": 1,
          "chunkId": 1,
          "eventType": 0,
          "rule": 0, // PN rule --> severity mid
          "reporter": validator,
          "offender": validator,
          "rawProof": [], // not checked by the _slash function
          "block": 1,
          "epoch": 0,
          "reportingBlock": 2,
          "messageHash": 0, // not checked by the _slash function
        }
        let epochOffenceCount = 1;
        let tx = await accountability.slash(event, epochOffenceCount);
        // checking if highest possible slashing can be done without triggering fairness issue
        truffleAssert.eventEmitted(tx, 'SlashingEvent', (ev) => {
          return ev.amount.toNumber() > 0 && ev.amount.toNumber() == totalStake - 2;
        });
        valInfo = await autonity.getValidator(validator);
        assert.equal(valInfo.state, utils.ValidatorState.jailed, "validator not jailed");
        assert(parseInt(valInfo.bondedStake) > 0 && parseInt(valInfo.unbondingStake) > 0, "fairness issue triggered");
      }

    });


  });
});
