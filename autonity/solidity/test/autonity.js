'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const liquidContract = artifacts.require("Liquid")
const toBN = web3.utils.toBN;

// while testing we might ran into situations were currentHeight > lastEpochBlock + epochPeriod
// in this case in order to be able to finalize we need to setEpochPeriod to a bigger value
// also we need to take into account that if we are running against autonity, the network will keep mining as we do these operations
async function endEpoch(contract,operator,deployer){
  let lastEpochBlock = (await contract.getLastEpochBlock()).toNumber();
  let currentHeight = await web3.eth.getBlockNumber();
  let currentEpoch = (await contract.epochID()).toNumber()
  let delta = currentHeight - lastEpochBlock
  let epochPeriod = delta + 5

  await contract.setEpochPeriod(epochPeriod,{from: operator})

  assert.equal(epochPeriod,(await contract.getEpochPeriod()).toNumber())

  // close epoch
  for (let i=0;i<(lastEpochBlock + epochPeriod) - currentHeight;i++) {
    let height = await web3.eth.getBlockNumber()
    contract.finalize({from: deployer})
    await utils.waitForNewBlock(height);
  }
  let newEpoch = (await contract.epochID()).toNumber()
  assert.equal(currentEpoch+1,newEpoch)
}

function generateMultiSig(nodekey,oraclekey,treasuryAddr) {
      let treasuryProof = web3.eth.accounts.sign(treasuryAddr, nodekey);
      let oracleProof = web3.eth.accounts.sign(treasuryAddr, oraclekey);
      let multisig = treasuryProof.signature + oracleProof.signature.substring(2)
      return multisig
}

contract('Autonity', function (accounts) {
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
  const version = 1;
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
  describe('Contract initial state', function () {
    beforeEach(async function () {
      autonity = await utils.deployContracts(validators, autonityConfig, accountabilityConfig,  deployer, operator);
    });

    it('test get token name', async function () {
      let n = await autonity.name({from: anyAccount});
      assert(name === n, "token name is not expected");
    });

    it('test get token symbol', async function () {
      let s = await autonity.symbol({from: anyAccount});
      assert(symbol === s, "token symbol is not expected");
    });

    it('test get min base fee after contract construction', async function () {
      let mBaseFee = await autonity.getMinimumBaseFee({from: operator});
      assert(minBaseFee == mBaseFee, "min base fee is not expected");
    });

    it('test get contract version after contract construction', async function () {
      let v = await autonity.getVersion({from: anyAccount});
      assert(version == v, `version of contract is not expected, has ${v} want ${version}`);
    });

    it('test get committee size after contract construction', async function () {
      let cS = await autonity.getMaxCommitteeSize({from: anyAccount});
      assert(committeeSize == cS, "committee size is not expected");
    });

    it('test get operator account after contract construction', async function () {
      let ac = await autonity.getOperator({from: anyAccount});
      assert.deepEqual(operator, ac);
    });

    it('test get validators after contract construction', async function () {
      let vals = await autonity.getValidators({from: anyAccount});
      assert.deepEqual(vals.slice().sort(), orderedValidatorsList.slice().sort(), "validator set is not expected");
    });

    it('test get committee after contract construction', async function () {
      let committee = await autonity.getCommittee({from: anyAccount});
      let committeeValidators = [];
      for (let i = 0; i < committee.length; i++) {
        committeeValidators.push(committee[i].addr);
      }
      assert.deepEqual(committeeValidators.slice().sort(), orderedValidatorsList.slice().sort(), "Committee should be equal than validator set");
    });

    it('test get committee enodes after contract construction', async function () {
      let committeeEnodes = await autonity.getCommitteeEnodes({from: anyAccount});
      assert.deepEqual(committeeEnodes.slice().sort(), genesisEnodes.slice().sort(), "Committee enodes should be equal to genesis validator enodes");
    });

    it('test getValidator, balanceOf, and totalSupply after contract construction', async function () {
      let total = 0;
      for (let i = 0; i < validators.length; i++) {
        total += validators[i].bondedStake;
        let b = await autonity.balanceOf(validators[i].treasury, {from: anyAccount});
        // since all stake token are bonded by default, those validators have no Newton token in the account.
        assert.equal(b.toNumber(), 0, "initial balance of validator is not expected");

        let v = await autonity.getValidator(validators[i].nodeAddress, {from: anyAccount});
        assert.equal(v.treasury.toString(), validators[i].treasury.toString(), "treasury addr is not expected at contract construction");
        assert.equal(v.nodeAddress.toString(), validators[i].nodeAddress.toString(), "validator addr is not expected at contract construction");
        assert.equal(v.enode.toString(), validators[i].enode.toString(), "validator enode is not expected at contract construction");
        assert(v.commissionRate == validators[i].commissionRate, "validator commission rate is not expected at contract construction");
        assert(v.bondedStake == validators[i].bondedStake, "validator bonded stake is not expected at contract construction");
        assert(v.totalSlashed == validators[i].totalSlashed, "validator total slash counter is not expected at contract construction");
        assert(v.registrationBlock == validators[i].registrationBlock, "registration block is not expected at contract construction");
        assert(v.state == validators[i].state, "validator state is not expected at contract construction");
      }
      let totalSupply = await autonity.totalSupply({from: anyAccount});
      assert.equal(total, totalSupply.toNumber(), "Newton total supply is not expected at contract construction phase");
    });
  });

  describe("Validator commission rate", () => {
    beforeEach(async function () {
      // the test contract exposes the applyNewCommissionRates function
      autonity = await utils.deployAutonityTestContract(validators, autonityConfig, {from: deployer});
    });

    it("should revert with bad input", async () => {
      await truffleAssert.fails(
        autonity.changeCommissionRate(genesisNodeAddresses[1], 1337, {from:accounts[3]}),
        truffleAssert.ErrorType.REVERT,
        "require caller to be validator admin account"
      );

      await truffleAssert.fails(
        autonity.changeCommissionRate(accounts[5], 1337, {from:accounts[3]}),
        truffleAssert.ErrorType.REVERT,
        "validator must be registered"
      );

      await truffleAssert.fails(
        autonity.changeCommissionRate(genesisNodeAddresses[3], 13370, {from:accounts[4]}),
        truffleAssert.ErrorType.REVERT,
        "require correct commission rate"
      );

    });

    it("should change a validator commission rate with correct inputs", async () => {
      const txChangeRate = await autonity.changeCommissionRate(genesisNodeAddresses[1], 1337, {from:accounts[1]});
      truffleAssert.eventEmitted(txChangeRate, 'CommissionRateChange', (ev) => {
        return ev.validator === genesisNodeAddresses[1] && ev.rate.toString() == "1337";
      }, 'should emit correct event');

      await autonity.changeCommissionRate(genesisNodeAddresses[3], 1339, {from:accounts[4]});
      await autonity.changeCommissionRate(genesisNodeAddresses[1], 1338, {from:accounts[1]});

      const txApplyCommChange = await autonity.applyNewCommissionRates({from:deployer});
      const v1 = await autonity.getValidator(genesisNodeAddresses[1]);
      assert.equal(v1.commissionRate,1338);

      const v3 = await autonity.getValidator(genesisNodeAddresses[3]);
      assert.equal(v3.commissionRate,1339);

    })

    it("should change a validator commission rate only after unbonding period", async () => {
      await autonity.setUnbondingPeriod(5, {from:operator});
      await autonity.changeCommissionRate(genesisNodeAddresses[1], 1338, {from:accounts[1]});
      await autonity.applyNewCommissionRates({from:deployer});
      let v1 = await autonity.getValidator(genesisNodeAddresses[1]);
      assert.equal(v1.commissionRate,100);

      await utils.mineEmptyBlock()
      await utils.mineEmptyBlock()
      await utils.mineEmptyBlock()
      await utils.mineEmptyBlock()
      await utils.mineEmptyBlock()
      await utils.mineEmptyBlock()
      await utils.mineEmptyBlock()

      await autonity.applyNewCommissionRates({from:deployer});
      v1 = await autonity.getValidator(genesisNodeAddresses[1]);
      assert.equal(v1.commissionRate,1338);
    });
  })

  describe('Set protocol parameters only by operator account', function () {
    beforeEach(async function () {
      autonity = await utils.deployContracts(validators, autonityConfig, accountabilityConfig, deployer, operator);
    });

    it('test set min base fee by operator', async function () {
      await autonity.setMinimumBaseFee(50000, {from: operator});
      let mGP = await autonity.getMinimumBaseFee({from: operator});
      assert(50000 == mGP, "min gas price is not expected");
    });

    it('test regular validator cannot set min base fee', async function () {
      let initMGP = await autonity.getMinimumBaseFee({from: operator});
      
      await truffleAssert.fails(
        autonity.setMinimumBaseFee(50000, {from: accounts[9]}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );

      let minGP = await autonity.getMinimumBaseFee({from: operator});
      assert.deepEqual(initMGP, minGP);
    });

    it('test set committee size by operator', async function () {
      await autonity.setCommitteeSize(500, {from: operator});
      let cS = await autonity.getMaxCommitteeSize({from: operator});
      assert(500 == cS, "committee size is not expected");
    });

    it('test regular validator cannot set committee size', async function () {
      let initCommitteeSize = await autonity.getMaxCommitteeSize({from: operator});
      
      await truffleAssert.fails(
        autonity.setCommitteeSize(500, {from: accounts[9]}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );

      let cS = await autonity.getMaxCommitteeSize({from: operator});
      assert.deepEqual(initCommitteeSize, cS);
    });

    it('test set un-bonding period by operator', async function () {
      await autonity.setUnbondingPeriod(127, {from: operator});
      let uP = await autonity.getUnbondingPeriod({from: operator});
      assert.equal(127,uP)
    });

    it('test regular validator cannot set un-bonding period', async function () {
      let initUP = await autonity.getUnbondingPeriod({from: operator});
      
      await truffleAssert.fails(
        autonity.setUnbondingPeriod(127, {from: accounts[9]}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );
      let uP = await autonity.getUnbondingPeriod({from: operator});
      assert.equal(initUP.toString(),uP.toString())
    });

    it('test extend epoch period by operator', async function () {
      await autonity.setEpochPeriod(98, {from: operator});
      let eP = await autonity.getEpochPeriod({from: operator});
      assert.equal("98",eP.toString())
    });
    
    it('test regular validator cannot extend epoch period', async function () {
      let initEP = await autonity.getEpochPeriod({from: operator});
      
      await truffleAssert.fails(
        autonity.setEpochPeriod(98, {from:accounts[9]}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );
      
      let eP = await autonity.getEpochPeriod({from: accounts[9]});
      assert.equal(initEP.toString(),eP.toString())
    });

    it('test set operator account by operator', async function () {
      let newOperator = accounts[9];
      await autonity.setOperatorAccount(newOperator, {from: operator});
      let nOP = await autonity.getOperator({from: operator});
      assert.deepEqual(newOperator, nOP);
    });

    it('test regular validator cannot set operator account', async function () {
      let initOperator = await autonity.getOperator({from: operator});
      
      await truffleAssert.fails(
        autonity.setOperatorAccount(accounts[1], {from: accounts[9]}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );

      let op = await autonity.getOperator({from: operator});
      assert.deepEqual(initOperator, op);
    });

    it('test set treasury account by operator', async function () {
      let newTreasury = accounts[1];
      await autonity.setTreasuryAccount(newTreasury, {from: operator});
      
      let treasury = await autonity.getTreasuryAccount({from: operator});
      assert.deepEqual(newTreasury,treasury)
    });

    it('test regular validator cannot set treasury account', async function () {
      let initTreasury = await autonity.getTreasuryAccount({from: operator});
      
      await truffleAssert.fails(
        autonity.setTreasuryAccount(accounts[9], {from: accounts[9]}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );
      
      let treasury = await autonity.getTreasuryAccount({from: operator});
      assert.deepEqual(initTreasury,treasury)
    });

    it('test set treasury fee by operator', async function () {
      let initFee = await autonity.getTreasuryFee({from: operator});
      let newFee = initFee + 1;
      await autonity.setTreasuryFee(newFee, {from: operator});
      let treasuryFee = await autonity.getTreasuryFee({from: operator});
      assert.equal(newFee,treasuryFee)
    });

    it.skip('test set treasury fee with invalid value by operator', async function () {
      // treasury fee should never exceed 1e9.
      let newFee = 10000000000;
      await truffleAssert.fails(
        autonity.setTreasuryFee(newFee, {from: operator}),
        truffleAssert.ErrorType.REVERT,
      );
    });

    it('test regular validator cannot set treasury fee', async function () {
      let initFee = await autonity.getTreasuryFee({from: operator});
      let newFee = initFee + 1;
      await truffleAssert.fails(
        autonity.setTreasuryFee(newFee, {from: accounts[9]}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );
      let treasuryFee = await autonity.getTreasuryFee({from: operator});
      assert.equal(treasuryFee.toString(),initFee.toString())
    });
  });

  describe('Test cases for ERC-20 token management', function () {
    beforeEach(async function () {
      autonity = await utils.deployContracts(validators, autonityConfig, accountabilityConfig, deployer, operator);
    });

    it('test mint Newton by operator', async function () {
      let account = accounts[7];
      let tokenMint = toBN('999999999990000000000000000000000');
      let initSupply = await autonity.totalSupply();
      await autonity.mint(account, tokenMint, {from: operator});
      let balance = await autonity.balanceOf(account);
      let newSupply = await autonity.totalSupply();
      assert(balance.toString() == tokenMint.toString(), "account balance is not expected");
      assert.equal(newSupply.toString(), initSupply.add(tokenMint).toString(), "total supply is not expected");
    });

    it('test regular validator cannot mint Newton', async function () {
      let initBalance = await autonity.balanceOf(accounts[1]);
      let tokenMint = 20;

      await truffleAssert.fails(
        autonity.mint(accounts[1], tokenMint, {from: anyAccount}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );
      let balance = await autonity.balanceOf(accounts[1]);
      assert.deepEqual(initBalance, balance);
    });

    it('test burn Newton by operator', async function () {
      let tokenMint = 20;
      let tokenBurn = 10;
      let initSupply = await autonity.totalSupply();
      let initBalance = await autonity.balanceOf(accounts[1]);
      // since all stake token are bonded, we mint new tokens for account, then try to burn them again.
      await autonity.mint(accounts[1], tokenMint, {from: operator});
      await autonity.burn(accounts[1], tokenBurn, {from: operator});
      let newBalance = await autonity.balanceOf(accounts[1]);
      let newSupply = await autonity.totalSupply();
      assert.equal(newBalance.toNumber(), initBalance.toNumber() + tokenMint - tokenBurn, "account balance is not expected");
      assert.equal(newSupply.toNumber(), initSupply.toNumber() + tokenMint - tokenBurn, "total supply is not expected");
    });

    it('test regular validator cannot burn Newton', async function () {
      let initBalance = await autonity.balanceOf(accounts[1]);
      let tokenBurn = 10;
      
      await truffleAssert.fails(
        autonity.burn(accounts[1], tokenBurn, {from: anyAccount}),
        truffleAssert.ErrorType.REVERT,
        "caller is not the operator"
      );
      let balance = await autonity.balanceOf(accounts[1]);
      assert.deepEqual(initBalance, balance);
    });

    it('test ERC-20 token transfer', async function () {
      // since all the token are bonded, so we mint tokens before transfer.
      let amount = 10;
      await autonity.mint(accounts[3], amount, {from: operator});

      let initBalanceA = await autonity.balanceOf(accounts[1]);
      let initBalanceB = await autonity.balanceOf(accounts[3]);
      await autonity.transfer(accounts[1], amount, {from: accounts[3]});
      let newBalanceA = await autonity.balanceOf(accounts[1]);
      let newBalanceB = await autonity.balanceOf(accounts[3]);
      assert.equal(initBalanceB.toNumber(), newBalanceB.toNumber() + amount, "sender balance is not expected");
      assert.equal(initBalanceA.toNumber(), newBalanceA.toNumber() - amount, "receiver balance is not expected");
    });

    it('test ERC-20 token transfer with no sufficient fund', async function () {
      let amount = 10000000;
      let initBalanceA = await autonity.balanceOf(accounts[1]);
      let initBalanceB = await autonity.balanceOf(accounts[3]);
      
      await truffleAssert.fails(
        autonity.transfer(accounts[1], amount, {from: accounts[3]}),
        truffleAssert.ErrorType.REVERT,
        "amount exceeds balance"
      );

      let bA = await autonity.balanceOf(accounts[1]);
      let bB = await autonity.balanceOf(accounts[3]);
      assert.equal(initBalanceB.toNumber(), bB.toNumber(), "sender balance is not expected");
      assert.equal(initBalanceA.toNumber(), bA.toNumber(), "receiver balance is not expected");
    });

    it('test ERC-20 token approve', async function () {
      let amountApproved = 10;
      let spender = accounts[1];
      let owner = accounts[3];
      await autonity.approve(spender, amountApproved, {from: owner});
      let approval = await autonity.allowance(owner, spender);
      assert.equal(approval.toNumber(), amountApproved, "token approval is not expected");
    });

    it('test ERC-20 token transferFrom', async function () {
      let owner = accounts[3];
      let spender = accounts[1];
      let amountApproved = 20;
      let amountTransfer = 10;
      // since all token were bonded, so we mint new tokens for account before they can be transfer.
      await autonity.mint(owner, 1000, {from: operator});

      let initBalanceOwner = await autonity.balanceOf(owner);
      let initBalanceSpender = await autonity.balanceOf(spender);

      await autonity.approve(spender, amountApproved, {from: owner});
      await autonity.transferFrom(owner, spender, amountTransfer, {from: spender});
      let newBalanceOwner = await autonity.balanceOf(owner);
      let newBalanceSpender = await autonity.balanceOf(spender);

      let allowance = await autonity.allowance(owner, spender);

      assert.equal(initBalanceOwner.toNumber(), newBalanceOwner.toNumber() + amountTransfer, "balance of owner is not expected");
      assert.equal(initBalanceSpender.toNumber(), newBalanceSpender.toNumber() - amountTransfer, "balance of spender is not expected");
      assert.equal(allowance.toNumber(), amountApproved - amountTransfer, "allowance is not expected");
    });
  });

  describe('Validator management', function () {
    beforeEach(async function () {
      autonity = await utils.deployContracts(validators, autonityConfig, accountabilityConfig, deployer, operator);
    });

    it('Add validator with already registered address', async function () {
      let newValidator = accounts[0];
      let enode = genesisEnodes[0]
      // multisig length is checked before validator already registered (it is not verified though)
      let multisig = generateMultiSig(genesisPrivateKeys[0],genesisPrivateKeys[0],newValidator)
      
      await truffleAssert.fails(
        autonity.registerValidator(enode, genesisNodeAddresses[0], multisig, {from: newValidator}),
        truffleAssert.ErrorType.REVERT,
        "validator already registered"
      );

      let vals = await autonity.getValidators();
      assert.equal(vals.length, validators.length, "validator pool is not expected");
    });

    it('Add a validator with invalid enode address', async function () {
      let newValidator = accounts[8];
      let enode = "enode://invalidEnodeAddress@172.25.0.11:30303";
      let privateKey = genesisPrivateKeys[0] // irrelevant
      let multisig = generateMultiSig(privateKey,privateKey,newValidator)
      
      await truffleAssert.fails(
        autonity.registerValidator(enode, newValidator, multisig, {from: newValidator}),
        truffleAssert.ErrorType.REVERT,
        "enode error"
      );

      let vals = await autonity.getValidators();
      assert.equal(vals.length, validators.length, "validator pool is not expected");
    });

    it('Add a validator with invalid oracle proof', async function () {
      let newValidator = accounts[8];
      let enode = freeEnodes[0]
      let privateKey = freePrivateKeys[0]
      // generate oracle signature with nodekey instead of treasury key
      let multisig = generateMultiSig(privateKey,privateKey,newValidator)
      let oracleAddr = newValidator // treasury address
      
      await truffleAssert.fails(
        autonity.registerValidator(enode, oracleAddr, multisig, {from: newValidator}),
        truffleAssert.ErrorType.REVERT,
        "Invalid oracle key ownership proof provided"
      );

      let vals = await autonity.getValidators();
      assert.equal(vals.length, validators.length, "validator pool is not expected");
    });

    it('Add a validator with valid meta data', async function () {
      let issuerAccount = accounts[8];
      let newValAddr = freeAddresses[0]
      let enode = freeEnodes[0]
      let privateKey = freePrivateKeys[0]
      let multisig = generateMultiSig(privateKey,privateKey,issuerAccount)
      let oracleAddr = newValAddr

      await autonity.registerValidator(enode, oracleAddr, multisig, {from: issuerAccount});
      let vals = await autonity.getValidators();
      assert.equal(vals.length, validators.length + 1, "validator pool is not expected");

      let v = await autonity.getValidator(newValAddr, {from: issuerAccount});

      const liquidABI = liquidContract["abi"]
      const liquid = new web3.eth.Contract(liquidABI, v.liquidContract);
      assert.equal(await liquid.methods.name().call(),"LNTN-"+(vals.length-1))
      assert.equal(await liquid.methods.symbol().call(),"LNTN-"+(vals.length-1))
      assert.equal(v.treasury.toString(), issuerAccount.toString(), "treasury addr is not expected");
      assert.equal(v.nodeAddress.toString(), newValAddr.toString(), "validator addr is not expected");
      assert.equal(v.enode.toString(), enode.toString(), "validator enode is not expected");
      assert(v.bondedStake == 0, "validator bonded stake is not expected");
      assert(v.totalSlashed == 0, "validator total slash counter is not expected");
      assert(v.state == 0, "validator state is not expected");
    });

    it('Pause a validator', async function () {
      let validator = freeAddresses[0];
      let issuerAccount = accounts[8];
      let enode = freeEnodes[0]
      let privateKey = freePrivateKeys[0]
      let multisig = generateMultiSig(privateKey,privateKey,issuerAccount)
      let oracleAddr = validator

      // disabling a non registered validator should fail
      await truffleAssert.fails(
        autonity.pauseValidator(validator, {from: issuerAccount}),
        truffleAssert.ErrorType.REVERT,
        "validator must be registered"
      );

      await autonity.registerValidator(enode, oracleAddr, multisig, {from: issuerAccount});

      // try disabling it with msg.sender not the treasury account, it should fails 
      await truffleAssert.fails(
        autonity.pauseValidator(validator, {from: accounts[7]}),
        truffleAssert.ErrorType.REVERT,
        "require caller to be validator admin account"
      );

      await autonity.pauseValidator(validator, {from: issuerAccount});
      let v = await autonity.getValidator(validator, {from: issuerAccount});
      assert(v.state == 1, "validator state is not expected");

      // try disabling it again, it should fail 
      await truffleAssert.fails(
        autonity.pauseValidator(validator, {from: issuerAccount}),
        truffleAssert.ErrorType.REVERT,
        "validator must be active"
      );
    });

    it("Re-active a paused validator", async function () {
      let issuerAccount = accounts[8];

      let validator = freeAddresses[0]
      let enode = freeEnodes[0]
      let privateKey = freePrivateKeys[0];
      // activating a non-existing validator should fail 
      await truffleAssert.fails(
        autonity.activateValidator(validator, {from: issuerAccount}),
        truffleAssert.ErrorType.REVERT,
        "validator must be registered"
      );

      let multisig = generateMultiSig(privateKey,privateKey,issuerAccount)
      let oracleAddr = validator
      await autonity.registerValidator(enode, oracleAddr, multisig, {from: issuerAccount});
      
      // activating from non-treasury account should fail 
      await truffleAssert.fails(
        autonity.activateValidator(validator, {from: accounts[7]}),
        truffleAssert.ErrorType.REVERT,
        "require caller to be validator treasury account"
      );
      
      // activating an already active validator should fail 
      await truffleAssert.fails(
        autonity.activateValidator(validator, {from: issuerAccount}),
        truffleAssert.ErrorType.REVERT,
        "validator already active"
      );
      await autonity.pauseValidator(validator, {from: issuerAccount});
      let v = await autonity.getValidator(validator, {from: issuerAccount});
      assert(v.state == 1, "validator state is not expected");
      await autonity.activateValidator(validator, {from: issuerAccount});
      v = await autonity.getValidator(validator, {from: issuerAccount});
      assert(v.state == 0, "validator state is not expected");
    })
  });

  describe('Bonding and unbonding requests', function () {
    beforeEach(async function () {
      autonity = await utils.deployContracts(validators, autonityConfig, accountabilityConfig, deployer, operator);
    });

    //TODO(tariq) replicate this test for a selfBonded bond request --> no LNTN minting for selfBonded stake
    it('Bond to a valid validator (not selfBonded)', async function () {
      // mint Newton for a new account.
      let newAccount = accounts[8];
      let tokenMint = 200;
      await autonity.mint(newAccount, tokenMint, {from: operator});

      // bond new minted Newton to a registered validator.
      let tx = await autonity.bond(validators[0].nodeAddress, tokenMint, {from: newAccount});
      truffleAssert.eventEmitted(tx, 'NewBondingRequest', (ev) => {
        return ev.validator === validators[0].nodeAddress && ev.delegator === newAccount && ev.selfBonded === false && ev.amount.toNumber() === tokenMint
      }, 'should emit newBondingRequest event');
      //TODO(tariq) check effects of bond:
      //                                  1. bonded NTN is substracted from balance of delegator
      //                                  2. LNTN is minted to delegator at epoch end (to trigger epoch end, see endEpoch function helper)
      
      /* TODO(tariq) the internal queues for bond and unbond are not publicly accessible anymore.
       * To run these checks we need another contract that inherits Autonity and exposes the bondingMap and unbondingMap.
       * See AutonityTest.sol for the same approach applied to applyNewCommissionRates()
      // num of stakings from contract construction equals: length of validators and the latest bond.
      let numOfStakings = validators.length + 1;
      // ids start from 0
      let latestBondingReqId = numOfStakings - 1       

      assert.equal(latestBondingReqId, (await autonity.getLastRequestedBondingRequest()).toNumber())

      let stakings = await autonity.getBondingRequests(0, latestBondingReqId);
      assert.equal(stakings[latestBondingReqId].amount, tokenMint, "stake bonding amount is not expected");
      assert.equal(stakings[latestBondingReqId].delegator, newAccount, "delegator addr is not expected");
      assert.equal(stakings[latestBondingReqId].delegatee, validators[0].nodeAddress, "delegatee addr is not expected");
      */
    });

    it('does not bond on a non registered validator', async function () {
      // mint Newton for a new account.
      let newAccount = accounts[8];
      let tokenMint = 200;
      await autonity.mint(newAccount, tokenMint, {from: operator});

      // bond new minted Newton to a not registered validator.
      await truffleAssert.fails(
        autonity.bond(anyAccount, tokenMint, {from: newAccount}),
        truffleAssert.ErrorType.REVERT,
        "validator not registered"
      );
    });
    
    it("can't bond to a paused validator", async function () {
      await autonity.pauseValidator(validators[0].nodeAddress, {from: validators[0].treasury});
      
      await truffleAssert.fails(
        autonity.bond(validators[0].nodeAddress, 100, {from: validators[0].treasury}),
        truffleAssert.ErrorType.REVERT,
        "validator need to be active"
      );
    });

    //TODO(tariq) replicate this test for a non-selfBonded unbond request --> 
    //                                                                        1. LNTN locked when unbonding request issued
    //                                                                        2. LNTN burned at the end of the epoch following the unbonding request
    //                                                                        3. Unbonding shares issued at the end of the epoch
    //                                                                        4. Unbonding shares converted to NTNs and released at the end of the unbonding period
    it('un-bond from a valid validator', async function () {
      let tokenUnBond = 10;
      let from = validators[0].treasury;
      // unBond from self, a registered validator.
      let tx = await autonity.unbond(validators[0].nodeAddress, tokenUnBond, {from: from});
      
      truffleAssert.eventEmitted(tx, 'NewUnbondingRequest', (ev) => {
        return ev.validator === validators[0].nodeAddress && ev.delegator === from && ev.selfBonded === true && ev.amount.toNumber() === tokenUnBond
      }, 'should emit newUnbondingRequest event');
      //TODO(tariq) check effects of unbond (selfBonded):
      //                                  1. unbonded NTN enters "unbonding" state at epoch end and unbonding shares are issued. validator voting power (bondedStake) decreases
      //                                  3. At the end of the unbonding period the unbonding shares are converted back to NTNs and released.

      /* TODO(tariq) the internal queues for bond and unbond are not publicly accessible anymore.
       * To run these checks we need another contract that inherits Autonity and exposes the bondingMap and unbondingMap.
       * See AutonityTest.sol for the same approach applied to applyNewCommissionRates()
      let numOfUnBonding = 1;
      let latestUnbondingReqId = numOfUnBonding - 1

      assert.equal(latestUnbondingReqId, (await autonity.getLastRequestedUnbondingRequest()).toNumber())

      let unStakings = await autonity.getUnbondingRequests(0, latestUnbondingReqId);
      assert.equal(unStakings[latestUnbondingReqId].amount, tokenUnBond, "stake unbonding amount is not expected");
      assert.equal(unStakings[latestUnbondingReqId].delegator, from, "delegator addr is not expected");
      assert.equal(unStakings[latestUnbondingReqId].delegatee, validators[0].nodeAddress, "delegatee addr is not expected");
      */
    });

    it('does not unbond from not registered validator', async function () {
      let unRegisteredVal = anyAccount;
      let tokenUnBond = 10;
      
      await truffleAssert.fails(
        autonity.unbond(unRegisteredVal, tokenUnBond, {from: validators[0].treasury}),
        truffleAssert.ErrorType.REVERT,
        "validator not registered",
      );
    });

    it("can't unbond from  avalidator with the amount exceeding the available balance", async function () {
      let tokenUnBond = 99999;
      let from = validators[0].treasury;
      
      await truffleAssert.fails(
        autonity.unbond(validators[0].nodeAddress, tokenUnBond, {from: from}),
        truffleAssert.ErrorType.REVERT,
        "insufficient self bonded newton balance"
      );
    });
    
    /* TODO(tariq) the internal queues for bond and unbond are not publicly accessible anymore.
     * To run these checks we need another contract that inherits Autonity and exposes the bondingMap and unbondingMap.
     * See AutonityTest.sol for the same approach applied to applyNewCommissionRates()
     */
    it.skip('test bonding queue logic', async function () {
      // num of stakings from contract construction equals: length of validators 
      let numOfStakings = validators.length;

      // they are all processed at contract construction time, so there should be no pending requests
      await truffleAssert.fails(
        autonity.getFirstPendingBondingRequest(),
        truffleAssert.ErrorType.REVERT,
        "no pending requests"
      );
      
      // ids start from 0
      let latestBondingReqId = numOfStakings - 1       
      assert.equal(latestBondingReqId, (await autonity.getLastRequestedBondingRequest()).toNumber())
      
      // do a new bonding req
      let newAccount = accounts[8];
      let tokenMint = 200;
      await autonity.mint(newAccount, tokenMint, {from: operator});
      await autonity.bond(validators[0].nodeAddress, tokenMint, {from: newAccount});
      numOfStakings++
      
      // ids start from 0
      latestBondingReqId = numOfStakings - 1       
      assert.equal(latestBondingReqId, (await autonity.getLastRequestedBondingRequest()).toNumber())
      assert.equal(latestBondingReqId, (await autonity.getFirstPendingBondingRequest()).toNumber())

      let stakings = await autonity.getBondingRequests(latestBondingReqId, latestBondingReqId);

      assert.equal(stakings[0].amount, tokenMint, "stake bonding amount is not expected");
      assert.equal(stakings[0].delegator, newAccount, "delegator addr is not expected");
      assert.equal(stakings[0].delegatee, validators[0].nodeAddress, "delegatee addr is not expected");
    });
    /* TODO(tariq) the internal queues for bond and unbond are not publicly accessible anymore.
     * To run these checks we need another contract that inherits Autonity and exposes the bondingMap and unbondingMap.
     * See AutonityTest.sol for the same approach applied to applyNewCommissionRates()
     */
    it.skip('test unbonding queue logic', async function () {
      // no unbondings from contract construction
      await truffleAssert.fails(
        autonity.getFirstPendingUnbondingRequest(),
        truffleAssert.ErrorType.REVERT,
        "no pending requests"
      );
      await truffleAssert.fails(
        autonity.getLastRequestedUnbondingRequest(),
        truffleAssert.ErrorType.REVERT,
        "no unbonding request processed",
      );
      
      // do a new unbonding req
      let tokenUnBond = 10;
      let from = validators[0].treasury;
      await autonity.unbond(validators[0].nodeAddress, tokenUnBond, {from: from});
      
      let latestUnbondingReqId = 0
      assert.equal(latestUnbondingReqId, (await autonity.getLastRequestedUnbondingRequest()).toNumber())
      assert.equal(latestUnbondingReqId, (await autonity.getFirstPendingUnbondingRequest()).toNumber())

      let unstakings = await autonity.getUnbondingRequests(latestUnbondingReqId, latestUnbondingReqId);

      assert.equal(unstakings[0].amount, tokenUnBond, "stake unbonding amount is not expected");
      assert.equal(unstakings[0].delegator, validators[0].treasury, "delegator addr is not expected");
      assert.equal(unstakings[0].delegatee, validators[0].nodeAddress, "delegatee addr is not expected");
    });

  });

  describe('Test committee members rotation through bonding/unbonding', function () {
      let vals = [
        { ...baseValidator,
          "treasury": accounts[0],
          "nodeAddress": genesisNodeAddresses[0],
          "oracleAddress": accounts[0],
          "enode": genesisEnodes[0],
          "commissionRate": 10000,
          "bondedStake": 100,
        },
      ];
      let copyParams = autonityConfig;
      let token;
      beforeEach(async function () {
          // set short epoch period 
          let customizedEpochPeriod = 20;
          copyParams.protocol.epochPeriod = customizedEpochPeriod;

          token = await utils.deployContracts(vals, copyParams, accountabilityConfig, deployer, operator);
          assert.equal(customizedEpochPeriod,(await token.getEpochPeriod()).toNumber());
      });
      
      it('test bond stake token to new added validators, new validators should be elected as committee member', async function() {
          // register 2 new validators.
          let commissionRate = 10;
          let validator1 = accounts[1]; // treasury
          let oracle1 = genesisNodeAddresses[1] // oracle address = node address
          let enodeVal1 = genesisEnodes[1]
          let multisig1 = generateMultiSig(genesisPrivateKeys[1],genesisPrivateKeys[1],validator1)

          let validator2 = accounts[3];
          let oracle2 = genesisNodeAddresses[2] 
          let enodeVal2 = genesisEnodes[2]
          let multisig2 = generateMultiSig(genesisPrivateKeys[2],genesisPrivateKeys[2],validator2)

          await token.registerValidator(enodeVal1, oracle1, multisig1, {from: validator1});
          await token.registerValidator(enodeVal2, oracle2, multisig2, {from: validator2});

          // system operator mint Newton for user.
          let user = accounts[7];
          let tokenMint = 100;
          await token.mint(user, tokenMint, {from: operator});

          // user bond Newton to validator 2.
          await token.bond(genesisNodeAddresses[2], tokenMint, {from: user});
          
          // close epoch to ensure bonding is applied
          await endEpoch(token,operator,deployer);
                
          let committee = await token.getCommittee();
          let presented = false;
          for (let j=0; j<committee.length; j++) {
            if (genesisNodeAddresses[2] == committee[j].addr) {
              presented = true;
            }
            // we should not find the 0 bonded stake new validator
            if (genesisNodeAddresses[1] == committee[j].addr) {
              assert.fail("found unexpected committee member")
            }
          }
          assert.equal(presented, true);
      });
      it('test un-bond stake from validator, zero bonded validator should not be elected as committee member', async function() {
          // register 2 new validators.
          let commissionRate = 10;
          let validator1 = accounts[1]; // treasury
          let oracle1 = genesisNodeAddresses[1] // oracle address = node address
          let enodeVal1 = genesisEnodes[1]
          let multisig1 = generateMultiSig(genesisPrivateKeys[1],genesisPrivateKeys[1],validator1)

          let validator2 = accounts[3];
          let oracle2 = genesisNodeAddresses[2] 
          let enodeVal2 = genesisEnodes[2]
          let multisig2 = generateMultiSig(genesisPrivateKeys[2],genesisPrivateKeys[2],validator2)

          await token.registerValidator(enodeVal1, oracle1, multisig1, {from: validator1});
          await token.registerValidator(enodeVal2, oracle2, multisig2, {from: validator2});

          // system operator mint Newton for user.
          let user = accounts[7];
          let tokenMint = 100;
          await token.mint(user, tokenMint, {from: operator});

          // bond NTN to the 2 validators
          await token.bond(genesisNodeAddresses[1], 20, {from: user});
          await token.bond(genesisNodeAddresses[2], 20, {from: user});
          
          // close epoch to ensure bonding is applied
          await endEpoch(token,operator,deployer);
                
          let committee = await token.getCommittee();
          assert.equal(committee.length,3)

          await token.unbond(genesisNodeAddresses[1], 20, {from: user});  

          // voting power gets reduced right away at epoch end
          await endEpoch(token,operator,deployer);
          committee = await token.getCommittee();
          assert.equal(committee.length,2)
          
          let presented = false;
          for (let j=0; j<committee.length; j++) {
            // we should not find the validator we unbonded from 
            if (genesisNodeAddresses[1] == committee[j].addr) {
              assert.fail("found unexpected committee member")
            }
            // we should find the other one
            if (genesisNodeAddresses[2] == committee[j].addr) {
              presented = true;
            }
          }
          assert.equal(presented, true);
      });
  });

  describe('Test DPoS reward distribution', function () {
      let copyParams = autonityConfig;
      let token;
      beforeEach(async function () {
          // set short epoch period 
          let customizedEpochPeriod = 20;
          copyParams.protocol.epochPeriod = customizedEpochPeriod;

          token = await utils.deployContracts(validators, copyParams, accountabilityConfig, deployer, operator);
          assert.equal((await token.getEpochPeriod()).toNumber(),customizedEpochPeriod);
      });

      it('test finalize with not deployer account, exception should rise.', async function () {
          await truffleAssert.fails(
            token.finalize({from: anyAccount}),
            truffleAssert.ErrorType.REVERT,
            "function restricted to the protocol",
          );
      });

      it('test reward distribution with only selfBondedStake (no delegated stake)', async function () {
          let reward = 1000000000000000;
          // contract account should have no funds.
          let initFunds = await web3.eth.getBalance(token.address);
          assert.equal(initFunds,0);

          // send funds to contract account, to get them distributed later on.
          await web3.eth.sendTransaction({from: anyAccount, to: token.address, value: reward});
          let loadedBalance = await web3.eth.getBalance(token.address);
          assert.equal(loadedBalance, reward);

          // get validators and treasury initial ATN balance
          let initBalanceV0 = toBN(await web3.eth.getBalance(validators[0].treasury));
          let initBalanceV1 = toBN(await web3.eth.getBalance(validators[1].treasury));
          let initBalanceV2 = toBN(await web3.eth.getBalance(validators[2].treasury));
          let initBalanceV3 = toBN(await web3.eth.getBalance(validators[3].treasury));
          let initBalanceTreasury = toBN(await web3.eth.getBalance(treasuryAccount));

          // close epoch --> rewards are distributed
          await endEpoch(token,operator,deployer)

          // check autonity treasury reward
          let expectedTreasuryReward = toBN(copyParams.policy.treasuryFee).mul(toBN(reward)).div(toBN(10 ** 18));
          let afterBalanceTreasury = toBN(await web3.eth.getBalance(treasuryAccount));
          assert.equal(afterBalanceTreasury.sub(initBalanceTreasury).toString(),expectedTreasuryReward.toString())

          // check validators rewards
          let validatorRewards = toBN(reward).sub(expectedTreasuryReward)
          let totalStake = toBN(validators[0].bondedStake).add(toBN(validators[1].bondedStake)).add(toBN(validators[2].bondedStake)).add(toBN(validators[3].bondedStake)) 
          assert.equal(totalStake.toString(),"420")

          let afterBalanceV0 = toBN(await web3.eth.getBalance(validators[0].treasury));
          let expectedRewardV0 = toBN(validators[0].bondedStake).mul(validatorRewards).div(totalStake);
          assert.equal(afterBalanceV0.sub(initBalanceV0).toString(),expectedRewardV0.toString())
          
          let afterBalanceV1 = toBN(await web3.eth.getBalance(validators[1].treasury));
          let expectedRewardV1 = toBN(validators[1].bondedStake).mul(validatorRewards).div(totalStake);
          assert.equal(afterBalanceV1.sub(initBalanceV1).toString(),expectedRewardV1.toString())
          
          let afterBalanceV2 = toBN(await web3.eth.getBalance(validators[2].treasury));
          let expectedRewardV2 = toBN(validators[2].bondedStake).mul(validatorRewards).div(totalStake);
          assert.equal(afterBalanceV2.sub(initBalanceV2).toString(),expectedRewardV2.toString())
          
          let afterBalanceV3 = toBN(await web3.eth.getBalance(validators[3].treasury));
          let expectedRewardV3 = toBN(validators[3].bondedStake).mul(validatorRewards).div(totalStake);
          assert.equal(afterBalanceV3.sub(initBalanceV3).toString(),expectedRewardV3.toString())

          // Autonity contract should have left only dust ATN
          let leftFund = toBN(await web3.eth.getBalance(token.address));
          assert.equal(leftFund.toString(),toBN(loadedBalance).sub(expectedTreasuryReward).sub(expectedRewardV0).sub(expectedRewardV1).sub(expectedRewardV2).sub(expectedRewardV3).toString());
      });
      it('test reward distribution with delegations', async function () {
          const COMMISSION_RATE_PRECISION = 10000

          // mint Newton for external users.
          let alice = accounts[7]; // n.b. accounts[7] is also anyAccount
          let bob = accounts[9];
          await token.mint(alice, 200, {from: operator});
          await token.mint(bob, 200, {from: operator});

          // bond Newton in different validators.
          await token.bond(validators[0].nodeAddress, 120, {from: alice});
          await token.bond(validators[1].nodeAddress, 150, {from: bob});
          await token.bond(validators[2].nodeAddress, 80, {from: alice});
          await token.bond(validators[3].nodeAddress, 50, {from: bob});
          
          // close epoch --> bondings are applied
          await endEpoch(token,operator,deployer)
          
          // check the bonded stake should grows according to the new bonding by Alice and Bob.
          let val0 = await token.getValidator(validators[0].nodeAddress);
          assert.equal(val0.bondedStake,validators[0].bondedStake + 120)
          assert.equal(val0.selfBondedStake,validators[0].bondedStake)
          let val1 = await token.getValidator(validators[1].nodeAddress);
          assert.equal(val1.bondedStake,validators[1].bondedStake + 150)
          assert.equal(val1.selfBondedStake,validators[1].bondedStake)
          let val2 = await token.getValidator(validators[2].nodeAddress);
          assert.equal(val2.bondedStake,validators[2].bondedStake + 80)
          assert.equal(val2.selfBondedStake,validators[2].bondedStake)
          let val3 = await token.getValidator(validators[3].nodeAddress);
          assert.equal(val3.bondedStake,validators[3].bondedStake + 50)
          assert.equal(val3.selfBondedStake,validators[3].bondedStake)
          
          // get initial ATN balances
          let initBalanceV0 = toBN(await web3.eth.getBalance(validators[0].treasury));
          let initBalanceV1 = toBN(await web3.eth.getBalance(validators[1].treasury));
          let initBalanceV2 = toBN(await web3.eth.getBalance(validators[2].treasury));
          let initBalanceV3 = toBN(await web3.eth.getBalance(validators[3].treasury));
          let initBalanceTreasury = toBN(await web3.eth.getBalance(treasuryAccount));
          let initBalanceAlice = toBN(await web3.eth.getBalance(alice));
          let initBalanceBob = toBN(await web3.eth.getBalance(bob));
          
          // fund contract
          let reward = 1000000000000000;
          // contract account should have no funds.
          let initFunds = await web3.eth.getBalance(token.address);
          assert.equal(initFunds,0);

          // send funds to contract account, to get them distributed later on.
          await web3.eth.sendTransaction({from: operator, to: token.address, value: reward});
          let loadedBalance = await web3.eth.getBalance(token.address);
          assert.equal(loadedBalance, reward, "contract account have unexpected balance");
          
          // close epoch --> rewards are distributed
          await endEpoch(token,operator,deployer);

          let totalRewardsDistributed = toBN(0)
          
          // check autonity treasury reward
          let expectedTreasuryReward = toBN(copyParams.policy.treasuryFee).mul(toBN(reward)).div(toBN(10 ** 18));
          let afterBalanceTreasury = toBN(await web3.eth.getBalance(treasuryAccount));
          assert.equal(afterBalanceTreasury.sub(initBalanceTreasury).toString(),expectedTreasuryReward.toString())
          totalRewardsDistributed = totalRewardsDistributed.add(expectedTreasuryReward)

          // check validators rewards
          let validatorRewards = toBN(reward).sub(expectedTreasuryReward)
          let totalStake = toBN(val0.bondedStake).add(toBN(val1.bondedStake)).add(toBN(val2.bondedStake)).add(toBN(val3.bondedStake)) 
          assert.equal(totalStake.toString(),"820")

          let afterBalanceV0 = toBN(await web3.eth.getBalance(validators[0].treasury));
          let expectedRewardV0 = toBN(val0.bondedStake).mul(validatorRewards).div(totalStake);
          let selfRewardV0 = expectedRewardV0.mul(toBN(val0.selfBondedStake)).div(toBN(val0.bondedStake))
          let delegatorRewardV0 = expectedRewardV0.sub(selfRewardV0)
          let commissionIncomeV0 = delegatorRewardV0.mul(toBN(val0.commissionRate)).div(toBN(COMMISSION_RATE_PRECISION))
          assert.equal(afterBalanceV0.sub(initBalanceV0).toString(),selfRewardV0.add(commissionIncomeV0).toString())
          totalRewardsDistributed = totalRewardsDistributed.add(selfRewardV0).add(commissionIncomeV0)
          
          let afterBalanceV1 = toBN(await web3.eth.getBalance(validators[1].treasury));
          let expectedRewardV1 = toBN(val1.bondedStake).mul(validatorRewards).div(totalStake);
          let selfRewardV1 = expectedRewardV1.mul(toBN(val1.selfBondedStake)).div(toBN(val1.bondedStake))
          let delegatorRewardV1 = expectedRewardV1.sub(selfRewardV1)
          let commissionIncomeV1 = delegatorRewardV1.mul(toBN(val1.commissionRate)).div(toBN(COMMISSION_RATE_PRECISION))
          assert.equal(afterBalanceV1.sub(initBalanceV1).toString(),selfRewardV1.add(commissionIncomeV1).toString())
          totalRewardsDistributed = totalRewardsDistributed.add(selfRewardV1).add(commissionIncomeV1)

          let afterBalanceV2 = toBN(await web3.eth.getBalance(validators[2].treasury));
          let expectedRewardV2 = toBN(val2.bondedStake).mul(validatorRewards).div(totalStake);
          let selfRewardV2 = expectedRewardV2.mul(toBN(val2.selfBondedStake)).div(toBN(val2.bondedStake))
          let delegatorRewardV2 = expectedRewardV2.sub(selfRewardV2)
          let commissionIncomeV2 = delegatorRewardV2.mul(toBN(val2.commissionRate)).div(toBN(COMMISSION_RATE_PRECISION))
          assert.equal(afterBalanceV2.sub(initBalanceV2).toString(),selfRewardV2.add(commissionIncomeV2).toString())
          totalRewardsDistributed = totalRewardsDistributed.add(selfRewardV2).add(commissionIncomeV2)

          let afterBalanceV3 = toBN(await web3.eth.getBalance(validators[3].treasury));
          let expectedRewardV3 = toBN(val3.bondedStake).mul(validatorRewards).div(totalStake);
          let selfRewardV3 = expectedRewardV3.mul(toBN(val3.selfBondedStake)).div(toBN(val3.bondedStake))
          let delegatorRewardV3 = expectedRewardV3.sub(selfRewardV3)
          let commissionIncomeV3 = delegatorRewardV3.mul(toBN(val3.commissionRate)).div(toBN(COMMISSION_RATE_PRECISION))
          assert.equal(afterBalanceV3.sub(initBalanceV3).toString(),selfRewardV3.add(commissionIncomeV3).toString())
          totalRewardsDistributed = totalRewardsDistributed.add(selfRewardV3).add(commissionIncomeV3)
          
          // check delegators unclaimed reward
          const fee_factor_unit_recip = toBN(1000000000)

          let val0Liquid = await liquidContract.at(val0.liquidContract)
          let unclaimedRewardsV0 = await val0Liquid.unclaimedRewards(alice)
          // note(lorenzo) I added the .sub(toBN(1)) because the unclaimedRewards are sometimes 1 wei lower than what we expect due to rounding in Liquid.sol
          assert.equal(unclaimedRewardsV0.toString(),delegatorRewardV0.sub(commissionIncomeV0).sub(toBN(1)).toString())
          // the 1 wei was sent to the liquid contract, but the delegator cannot claim it due to rounding
          totalRewardsDistributed = totalRewardsDistributed.add(unclaimedRewardsV0).add(toBN(1)) 
          
          // check that if we mirror the computation done in Liquid.sol, we don't need the sub(toBN(1))
          let supplyV0 = toBN(await val0Liquid.totalSupply())
          let _rewardV0 = delegatorRewardV0.sub(commissionIncomeV0)
          let _unclaimedRewardsV0 = _rewardV0.mul(fee_factor_unit_recip).div(supplyV0).mul(toBN(120)).div(fee_factor_unit_recip)
          assert.equal(unclaimedRewardsV0.toString(),_unclaimedRewardsV0.toString())
          
          let val1Liquid = await liquidContract.at(val1.liquidContract)
          let unclaimedRewardsV1 = await val1Liquid.unclaimedRewards(bob)
          // note(lorenzo) I added the .sub(toBN(1)) because the unclaimedRewards are sometimes 1 wei lower than what we expect due to rounding in Liquid.sol
          assert.equal(unclaimedRewardsV1.toString(),delegatorRewardV1.sub(commissionIncomeV1).sub(toBN(1)).toString())
          // the 1 wei was sent to the liquid contract, but the delegator cannot claim it due to rounding
          totalRewardsDistributed = totalRewardsDistributed.add(unclaimedRewardsV1).add(toBN(1))
          
          // check that if we mirror the computation done in Liquid.sol, we don't need the sub(toBN(1))
          let supplyV1 = toBN(await val1Liquid.totalSupply())
          let _rewardV1 = delegatorRewardV1.sub(commissionIncomeV1)
          let _unclaimedRewardsV1 = _rewardV1.mul(fee_factor_unit_recip).div(supplyV1).mul(toBN(150)).div(fee_factor_unit_recip)
          assert.equal(unclaimedRewardsV1.toString(),_unclaimedRewardsV1.toString())

          let val2Liquid = await liquidContract.at(val2.liquidContract)
          let unclaimedRewardsV2 = await val2Liquid.unclaimedRewards(alice)
          assert.equal(unclaimedRewardsV2.toString(),delegatorRewardV2.sub(commissionIncomeV2).toString())
          totalRewardsDistributed = totalRewardsDistributed.add(unclaimedRewardsV2)
          
          // mirror computation in liquid.sol
          let supplyV2 = toBN(await val2Liquid.totalSupply())
          let _rewardV2 = delegatorRewardV2.sub(commissionIncomeV2)
          let _unclaimedRewardsV2 = _rewardV2.mul(fee_factor_unit_recip).div(supplyV2).mul(toBN(80)).div(fee_factor_unit_recip)
          assert.equal(unclaimedRewardsV2.toString(),_unclaimedRewardsV2.toString())
          
          let val3Liquid = await liquidContract.at(val3.liquidContract)
          let unclaimedRewardsV3 = await val3Liquid.unclaimedRewards(bob)
          assert.equal(unclaimedRewardsV3.toString(),delegatorRewardV3.sub(commissionIncomeV3).toString())
          totalRewardsDistributed = totalRewardsDistributed.add(unclaimedRewardsV3)
          
          // mirror computation in liquid.sol
          let supplyV3 = toBN(await val3Liquid.totalSupply())
          let _rewardV3 = delegatorRewardV3.sub(commissionIncomeV3)
          let _unclaimedRewardsV3 = _rewardV3.mul(fee_factor_unit_recip).div(supplyV3).mul(toBN(50)).div(fee_factor_unit_recip)
          assert.equal(unclaimedRewardsV3.toString(),_unclaimedRewardsV3.toString())

          // Autonity contract should have left only dust ATN
          let leftFund = toBN(await web3.eth.getBalance(token.address));
          assert.equal(leftFund.toString(),toBN(loadedBalance).sub(totalRewardsDistributed).toString())
    });
  });
});
