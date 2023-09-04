'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const liquidContract = artifacts.require("Liquid")
const AccountabilityTest = artifacts.require("AccountabilityTest")
const toBN = web3.utils.toBN;

// testing protocol contracts interactions.

const ValidatorState = {
  active : 0,
  paused : 1,
  jailed : 2
}

async function slash(accountability, epochOffenceCount, offender, reporter) {
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
  return txEvent;
}

// only non-self-unbond
// make sure selfBondedStake = 0 so that slashing can be applied to unbondingStake and delegatedStake
async function slashAndUnbond(autonity, accountability, delegators, validator, tokenUnbond, operator, deployer) {
  // not applicable for 100% slash
  let balances = [];
  for (let i = 0; i < delegators.length; i++) {
    balances.push((await autonity.balanceOf(delegators[i])).toNumber());
  }
  let valInfo = await autonity.getValidator(validator);
  let liquidSupply = Number(valInfo.liquidSupply);
  let delegatedStakes = Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake);
  let totalCurrentStake = delegatedStakes;
  let delegatee = [];
  delegatee.push(validator);
  let requestID = (await autonity.getHeadUnbondingID()).toNumber();
  let count = await utils.bulkUnbondingRequest(autonity, delegators, delegatee, tokenUnbond);
  // let unbonding apply
  await utils.endEpoch(autonity, operator, deployer);
  valInfo = await autonity.getValidator(validator);
  let unbondingStakes = Number(valInfo.unbondingStake);
  let unbondingShares = Number(valInfo.unbondingShares);
  let requests = [];
  while (count > 0) {
    let request = await autonity.getUnbondingRequest(requestID);
    requests.push(request);
    let newtonAmount = tokenUnbond * delegatedStakes / liquidSupply;
    let share = newtonAmount * unbondingShares / unbondingStakes;
    assert.equal(request.amount, tokenUnbond, "unexpected unbonding amount");
    assert.equal(request.unbondingShare, share, "unexpected unbonding share");
    assert.equal(request.selfDelegation, false, "unexpected self delegation");
    assert.equal(request.unlocked, true, "unexpected unbondingRequest.unlocked");
    requestID++;
    count--;
  }
  let txEvent = await slash(accountability, 1, validator, validator);
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
  let lastEpochBlock = (await autonity.getLastEpochBlock()).toNumber();
  let lastRequest = requests[requests.length - 1];
  let currentUnbondingPeriod = (await autonity.getUnbondingPeriod()).toNumber();
  let unbondingReleaseHeight = Number(lastRequest.requestBlock) + currentUnbondingPeriod;
  // the following needs to be true:
  // UnbondingRequestBlock + UnbondingPeriod > LastEpochBlock
  assert(
    unbondingReleaseHeight > lastEpochBlock,
    `unbonding period too short for testing, request-block: ${Number(lastRequest.requestBlock)}, unbonding-period: ${currentUnbondingPeriod}, `
    + `last-epoch-block: ${lastEpochBlock}`
  );
  // mine blocks until unbonding period is reached
  while (await web3.eth.getBlockNumber() < unbondingReleaseHeight) {
    await utils.mineEmptyBlock();
  }
  // release NTN
  await utils.endEpoch(autonity, operator, deployer);
  for (let i = 0; i < delegators.length; i++) {
    let balanceIncrease = unbondingStakes * Number(requests[i].unbondingShare) / unbondingShares;
    assert.equal((await autonity.balanceOf(delegators[i])).toNumber(), balances[i] + balanceIncrease, "unexpected balance");
  }

}


// only non-self-unbond
// make sure selfBondedStake = 0 so that slashing can be applied to delegatedStake
async function slashAndBond(autonity, accountability, delegators, validator, tokenBond, tokenUnbond, operator, deployer) {
  // not applicable for 100% slash
  let valInfo = await autonity.getValidator(validator);
  const valLiquid = await liquidContract.at(valInfo.liquidContract);
  let balances = [];
  for (let i = 0; i < delegators.length; i++) {
    balances.push((await valLiquid.balanceOf(delegators[i])).toNumber());
  }
  let delegatee = [];
  delegatee.push(validator);
  valInfo = await autonity.getValidator(validator);
  let liquidSupply = Number(valInfo.liquidSupply);
  let delegatedStakes = Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake);
  await utils.bulkBondingRequest(autonity, operator, delegators, delegatee, tokenBond);
  // let bonding apply
  await utils.endEpoch(autonity, operator, deployer);
  // LNTN minted
  for (let i = 0; i < delegators.length; i++) {
    let liquidAmount;
    if (delegatedStakes == 0) {
      liquidAmount = tokenBond;
    } else {
      liquidAmount = Math.floor(tokenBond * liquidSupply / delegatedStakes);
    }
    delegatedStakes += tokenBond;
    liquidSupply += liquidAmount;
    assert.equal((await valLiquid.balanceOf(delegators[i])).toNumber(), balances[i] + liquidAmount, "unexpected LNTN balance");
  }
  valInfo = await autonity.getValidator(validator);
  assert.equal(Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake), delegatedStakes, "delegated stake mismatch");
  assert.equal(Number(valInfo.liquidSupply), liquidSupply, "liquid supply mismatch");
  let totalCurrentStake = delegatedStakes + Number(valInfo.unbondingStake);

  let txEvent = await slash(accountability, 1, validator, validator);
  valInfo = await autonity.getValidator(validator);
  assert.equal(
    Number(valInfo.bondedStake) + Number(valInfo.unbondingStake),
    totalCurrentStake - txEvent.amount.toNumber(),
    "slashing amount does not match"
  );
  // conversion ratio chaned due to slashing
  liquidSupply = Number(valInfo.liquidSupply);
  delegatedStakes = Number(valInfo.bondedStake) - Number(valInfo.selfBondedStake);
  let unbondingStakes = Number(valInfo.unbondingStake);
  let unbondingShares = Number(valInfo.unbondingShares);
  assert(delegatedStakes > 0, "100% slashing");
  let requestID = (await autonity.getHeadUnbondingID()).toNumber();
  let count = await utils.bulkUnbondingRequest(autonity, delegators, delegatee, tokenUnbond);
  // let unbonding apply
  await utils.endEpoch(autonity, operator, deployer);
  while (count > 0) {
    let request = await autonity.getUnbondingRequest(requestID);
    let newtonAmount = Math.floor(tokenUnbond * delegatedStakes / liquidSupply);
    liquidSupply -= tokenUnbond;
    delegatedStakes -= newtonAmount;
    let share;
    if (unbondingStakes == 0) {
      share = newtonAmount;
    } else {
      share = Math.floor(newtonAmount * unbondingShares / unbondingStakes);
    }
    unbondingShares += share;
    unbondingStakes += newtonAmount;
    assert.equal(request.amount, tokenUnbond, "unexpected unbonding amount");
    assert.equal(request.unbondingShare, share, "unexpected unbonding share");
    assert.equal(request.selfDelegation, false, "unexpected self delegation");
    assert.equal(request.unlocked, true, "unexpected unbondingRequest.unlocked");
    requestID++;
    count--;
  }
  return txEvent.releaseBlock.toNumber();
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
      // issue multiple unbonding requests (both selfBonded and not) in different epochs, interleaved with slashing events
      // and check that the unbonding shares related fields change accordingly. Relevant fields to check:


      let config = JSON.parse(JSON.stringify(accountabilityConfig));
      // so that we don't encounter error due to fraction and we don't do 100% slashing
      config.collusionFactor = 0,
      config.historyFactor = 0;
      accountability = await AccountabilityTest.new(autonity.address, config, {from: deployer});
      await autonity.setAccountabilityContract(accountability.address, {from:operator});
      
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
      const tokenBond = 100000000;
      const tokenUnbond = 1000000;
      await utils.bulkBondingRequest(autonity, operator, delegators, delegatee, tokenBond);

      let requestID = (await autonity.getHeadUnbondingID()).toNumber() - 1;
      let request = await autonity.getUnbondingRequest(requestID);
      let currentUnbondingPeriod = (await autonity.getUnbondingPeriod()).toNumber();
      let unbondingReleaseHeight = Number(request.requestBlock) + currentUnbondingPeriod;
      // mine blocks until unbonding period is reached
      while (await web3.eth.getBlockNumber() < unbondingReleaseHeight) {
        await utils.mineEmptyBlock();
      }
      // requests will be processed at epoch end
      await utils.endEpoch(autonity, operator, deployer);
      // request unbonding and slash
      await slashAndUnbond(autonity, accountability, delegators, validator, tokenUnbond, operator, deployer);
      // repeat
      await slashAndUnbond(autonity, accountability, delegators, validator, tokenUnbond, operator, deployer);
    });

    it.skip('unbondingShares:unbondingStake 100% slash edge case', async function () {
      // see https://github.com/autonity/autonity/issues/819
      // low priority for now as we already now that the problem is there, however we need a test for the future fix
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
      let txEvent = await slash(accountability, epochOffenceCount, validator, treasury);
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

      let requestID = (await autonity.getHeadUnbondingID()).toNumber() - 1;
      let request = await autonity.getUnbondingRequest(requestID);
      let currentUnbondingPeriod = (await autonity.getUnbondingPeriod()).toNumber();
      let unbondingReleaseHeight = Number(request.requestBlock) + currentUnbondingPeriod;
      // mine blocks until unbonding period is reached
      while (await web3.eth.getBlockNumber() < unbondingReleaseHeight) {
        await utils.mineEmptyBlock();
      }
      // requests will be processed at epoch end
      await utils.endEpoch(autonity, operator, deployer);
      const tokenBond = 100000000;
      const factor = 1000;
      // request bonding and slash
      let releaseBlock = await slashAndBond(autonity, accountability, delegators, validator, tokenBond, factor, operator, deployer);
      while (await web3.eth.getBlockNumber() < releaseBlock) {
        await utils.mineEmptyBlock();
      }
      await autonity.activateValidator(validator, {from: treasury});
      // repeat
      await slashAndBond(autonity, accountability, delegators, validator, tokenBond, factor, operator, deployer);
    });

    it.skip('LNTN:NTN 100% slash edge case', async function () {
      // see https://github.com/autonity/autonity/issues/819
      // low priority for now as we already now that the problem is there, however we need a test for the future fix
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
      let txEvent = await slash(accountability, epochOffenceCount, validator, treasury);
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
      await slash(accountability, epochOffenceCount, validator, reporter)
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
      let txEvent = await slash(accountability, epochOffenceCount, validator, treasury);
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
      let height = await web3.eth.getBlockNumber();
      truffleAssert.eventEmitted(tx, 'ActivatedValidator', (ev) => {
        return ev.treasury === treasury && ev.addr === validator;
      });

    });
  });
});
