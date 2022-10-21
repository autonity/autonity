'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./test-utils');
const liquidContract = artifacts.require("Liquid")
//todo: move gas analysis to separate js file

contract('Autonity', function (accounts) {

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
  const treasuryFee = 1;
  const minimumEpochPeriod = 30;
  const version = 1;

  const config = {
    "operatorAccount": operator,
    "treasuryAccount": treasuryAccount,
    "treasuryFee": treasuryFee,
    "minBaseFee": minBaseFee,
    "delegationRate": delegationRate,
    "epochPeriod": epochPeriod,
    "unbondingPeriod": unBondingPeriod,
    "committeeSize": committeeSize,
    "contractVersion": version,
    "blockPeriod": minimumEpochPeriod,
  };

  const validators = [
    {
      "treasury": accounts[0],
      "addr": accounts[0],
      "enode": "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
      "commissionRate": 100,
      "bondedStake": 100,
      "totalSlashed": 0,
      "liquidContract": accounts[0],
      "liquidSupply": 0,
      "registrationBlock": 0,
      "state": 0,
    },
    {
      "treasury": accounts[1],
      "addr": accounts[1],
      "enode": "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
      "commissionRate": 100,
      "bondedStake": 90,
      "totalSlashed": 0,
      "liquidContract": accounts[1],
      "liquidSupply": 0,
      "registrationBlock": 0,
      "state": 0,
    },
    {
      "treasury": accounts[3],
      "addr": accounts[3],
      "enode": "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
      "commissionRate": 100,
      "bondedStake": 110,
      "totalSlashed": 0,
      "liquidContract": accounts[3],
      "liquidSupply": 0,
      "registrationBlock": 0,
      "state": 0,
    },
    {
      "treasury": accounts[4],
      "addr": accounts[4],
      "enode": "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303",
      "commissionRate": 100,
      "bondedStake": 120,
      "totalSlashed": 0,
      "liquidContract": accounts[4],
      "liquidSupply": 0,
      "registrationBlock": 0,
      "state": 0,
    },
  ];

  // initial validators ordered by bonded stake
  const orderedValidatorsList = [
    accounts[0],
    accounts[1],
    accounts[3],
    accounts[4],
  ];

  const freeEnodes = [
    "enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303",
    "enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303",
    "enode://ca6de62fce278f96aea6ec5a2daadb877e51651247cb96ee310a318def462913b653963c155a0ef6c7d50048bba6e6cea881130857413d9f50a621546b590758@34.255.23.113:30303",
    "enode://279944d8dcd428dffaa7436f25ca0ca43ae19e7bcf94a8fb7d1641651f92d121e972ac2e8f381414b80cc8e5555811c2ec6e1a99bb009b3f53c4c69923e11bd8@35.158.244.151:30303",
    "enode://8499da03c47d637b20eee24eec3c356c9a2e6148d6fe25ca195c7949ab8ec2c03e3556126b0d7ed644675e78c4318b08691b7b57de10e5f0d40d05b09238fa0a@52.187.207.27:30303"
  ];

  let autonity;

  describe('Contract initial state', function () {
    beforeEach(async function () {
      autonity = await utils.deployContract(validators, config, {from: deployer});
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
      for (let i = 0; i < committeeEnodes.length; i++) {
        let present = false;
        for (let j = 0; j < validators.length; j++) {
          if (committeeEnodes[i] === validators[j].enode) {
            present = true;
            break;
          }
        }
        assert(present === true, "cannot find committee enode from validator pool");
      }
    });

    it('test getValidator, balanceOf, and totalSupply after contract construction', async function () {
      let total = 0;
      for (let i = 0; i < validators.length; i++) {
        total += validators[i].bondedStake;
        let b = await autonity.balanceOf(validators[i].addr, {from: anyAccount});
        // since all stake token are bonded by default, those validators have no Newton token in the account.
        assert.equal(b.toNumber(), 0, "initial balance of validator is not expected");

        let v = await autonity.getValidator(validators[i].addr, {from: anyAccount});
        assert.equal(v.treasury.toString(), validators[i].treasury.toString(), "treasury addr is not expected at contract construction");
        assert.equal(v.addr.toString(), validators[i].addr.toString(), "validator addr is not expected at contract construction");
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
      autonity = await utils.deployTestContract(validators, config, {from: deployer});
    });

    it("should revert with bad input", async () => {
      await truffleAssert.fails(
        autonity.changeCommissionRate(accounts[1], 1337, {from:accounts[3]}),
        truffleAssert.ErrorType.REVERT,
        "require caller to be validator admin account"
      );

      await truffleAssert.fails(
        autonity.changeCommissionRate(accounts[5], 1337, {from:accounts[3]}),
        truffleAssert.ErrorType.REVERT,
        "validator must be registered"
      );

      await truffleAssert.fails(
        autonity.changeCommissionRate(accounts[3], 13370, {from:accounts[3]}),
        truffleAssert.ErrorType.REVERT,
        "require correct commission rate"
      );

    });

    it("should change a validator commission rate with correct inputs", async () => {
      const txChangeRate = await autonity.changeCommissionRate(accounts[1], 1337, {from:accounts[1]});
      truffleAssert.eventEmitted(txChangeRate, 'CommissionRateChange', (ev) => {
        return ev.validator === accounts[1] && ev.rate.toString() == "1337";
      }, 'should emit correct event');

      await autonity.changeCommissionRate(accounts[3], 1339, {from:accounts[3]});
      await autonity.changeCommissionRate(accounts[1], 1338, {from:accounts[1]});

      const txApplyCommChange = await autonity.applyNewCommissionRates({from:deployer});
      const v1 = await autonity.getValidator(accounts[1]);
      assert.equal(v1.commissionRate,1338);

      const v3 = await autonity.getValidator(accounts[3]);
      assert.equal(v3.commissionRate,1339);

    })

    it("should change a validator commission rate only after unbonding period", async () => {
      await autonity.setUnbondingPeriod(5, {from:operator});
      await autonity.changeCommissionRate(accounts[1], 1338, {from:accounts[1]});
      await autonity.applyNewCommissionRates({from:deployer});
      let v1 = await autonity.getValidator(accounts[1]);
      assert.equal(v1.commissionRate,100);
      await new Promise((resolve, reject) => {
        let wait = setTimeout(() => {
          clearTimeout(wait);
          resolve();
        }, 10000)
      })
      await autonity.applyNewCommissionRates({from:deployer});
      v1 = await autonity.getValidator(accounts[1]);
      assert.equal(v1.commissionRate,1338);
    });
  })

  describe('Set protocol parameters only by operator account', function () {
    beforeEach(async function () {
      autonity = await utils.deployContract(validators, config, {from: deployer});
    });

    it('test set min base fee by operator', async function () {
      await autonity.setMinimumBaseFee(50000, {from: operator});
      let mGP = await autonity.getMinimumBaseFee({from: operator});
      assert(50000 == mGP, "min gas price is not expected");
    });

    it('test regular validator cannot set min base fee', async function () {
      let initMGP = await autonity.getMinimumBaseFee({from: operator});

      try {
        let r = await autonity.setMinimumBaseFee(50000, {from: accounts[9]});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        let minGP = await autonity.getMinimumBaseFee({from: operator});
        assert.deepEqual(initMGP, minGP);
      }
    });

    it('test set committee size by operator', async function () {
      await autonity.setCommitteeSize(500, {from: operator});
      let cS = await autonity.getMaxCommitteeSize({from: operator});
      assert(500 == cS, "committee size is not expected");
    });

    it('test regular validator cannot set committee size', async function () {
      let initCommitteeSize = await autonity.getMaxCommitteeSize({from: operator});

      try {
        let r = await autonity.setCommitteeSize(500, {from: accounts[9]});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        let cS = await autonity.getMaxCommitteeSize({from: operator});
        assert.deepEqual(initCommitteeSize, cS);
      }
    });

    it('test set un-bonding period by operator', async function () {
      await autonity.setUnbondingPeriod(120, {from: operator});
    });

    it('test regular validator cannot set un-bonding period', async function () {
      try {
        let r = await autonity.setUnbondingPeriod(120, {from: accounts[9]});
        assert.fail('Expected throw not received', r);
      } catch (e) {
      }
    });

    it('test extend epoch period by operator', async function () {
      await autonity.setEpochPeriod(90, {from: operator});
    });

    it('test set operator account by operator', async function () {
      let newOperator = accounts[9];
      await autonity.setOperatorAccount(newOperator, {from: operator});
      let nOP = await autonity.getOperator({from: operator});
      assert.deepEqual(newOperator, nOP);
    });

    it('test regular validator cannot set operator account', async function () {
      let initOperator = await autonity.getOperator({from: operator});

      try {
        let r = await autonity.setOperatorAccount(accounts[1], {from: accounts[9]});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        let op = await autonity.getOperator({from: operator});
        assert.deepEqual(initOperator, op);
      }
    });

    it('test set treasury account by operator', async function () {
      let newTreasury = accounts[1];
      await autonity.setTreasuryAccount(newTreasury, {from: operator});
    });

    it('test regular validator cannot set treasury account', async function () {
      try {
        let r = await autonity.setTreasuryAccount(accounts[9], {from: accounts[9]});
        assert.fail('Expected throw not received', r);
      } catch (e) {
      }
    });

    it('test set treasury fee by operator', async function () {
      let newFee = treasuryFee + 1;
      await autonity.setTreasuryFee(newFee, {from: operator});
    });

    it('test set treasury fee with invalid value by operator', async function () {
      // treasury fee should never exceed 1e9.
      let newFee = 10000000000;
      try {
        let r = await autonity.setTreasuryFee(newFee, {from: operator});
        assert.fail('Expected throw not received', r);
      } catch (e) {
      }
    });

    it('test regular validator cannot set treasury fee', async function () {
      try {
        let newFee = treasuryFee + 1;
        let r = await autonity.setTreasuryAccount(newFee, {from: accounts[9]});
        assert.fail('Expected throw not received', r);
      } catch (e) {
      }
    });
  });

  describe('Test cases for ERC-20 token management', function () {
    beforeEach(async function () {
      autonity = await utils.deployContract(validators, config, {from: deployer});
    });

    it('test mint Newton by operator', async function () {
      let account = accounts[7];
      let tokenMint = 20;
      let initSupply = await autonity.totalSupply();
      await autonity.mint(account, tokenMint, {from: operator});
      let balance = await autonity.balanceOf(account);
      let newSupply = await autonity.totalSupply();
      assert(balance == tokenMint, "account balance is not expected");
      assert.equal(newSupply.toNumber(), initSupply.toNumber() + tokenMint, "total supply is not expected");
    });

    it('test regular validator cannot mint Newton', async function () {
      let initBalance = await autonity.balanceOf(accounts[1]);
      let tokenMint = 20;
      try {
        let r = await autonity.mint(accounts[1], tokenMint, {from: anyAccount});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        let balance = await autonity.balanceOf(accounts[1]);
        assert.deepEqual(initBalance, balance);
      }
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
      try {
        let r = await autonity.burn(accounts[1], tokenBurn, {from: anyAccount});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        let balance = await autonity.balanceOf(accounts[1]);
        assert.deepEqual(initBalance, balance);
      }
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

      try {
        let r = await autonity.transfer(accounts[1], amount, {from: accounts[3]});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        let bA = await autonity.balanceOf(accounts[1]);
        let bB = await autonity.balanceOf(accounts[3]);
        assert.equal(initBalanceB.toNumber(), bB.toNumber(), "sender balance is not expected");
        assert.equal(initBalanceA.toNumber(), bA.toNumber(), "receiver balance is not expected");
      }
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
      autonity = await utils.deployContract(validators, config, {from: deployer});
    });

    it('Add validator with already registered address', async function () {
      let newValidator = accounts[0];
      let enode = "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303";
      let privateKey = 'a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91';
      let proof = web3.eth.accounts.sign(newValidator, privateKey);

      try {
        let r = await autonity.registerValidator(enode, proof.signature, {from: newValidator});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        let vals = await autonity.getValidators();
        assert.equal(vals.length, validators.length, "validator pool is not expected");
      }
    });

    it('Add a validator with invalid enode address', async function () {
      let newValidator = accounts[8];
      let enode = "enode://invalidEnodeAddress@172.25.0.11:30303";
      let privateKey = 'e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b';
      let proof = web3.eth.accounts.sign(newValidator, privateKey);

      try {
        let r = await autonity.registerValidator(enode, proof.signature, {from: newValidator});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        let vals = await autonity.getValidators();
        assert.equal(vals.length, validators.length, "validator pool is not expected");
      }
    });

    it('Add a validator with valid meta data', async function () {
      let issuerAccount = accounts[8];
      let newValAddr ='0xDE03B7806f885Ae79d2aa56568b77caDB0de073E';
      let enode = "enode://a7ecd2c1b8c0c7d7ab9cc12e620605a762865d381eb1bc5417dcf07599571f84ce5725f404f66d3e254d590ae04e4e8f18fe9e23cd29087d095a0c37d0443252@3.209.45.79:30303";
      let privateKey = 'e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b';
      let proof = web3.eth.accounts.sign(issuerAccount, privateKey);

      await autonity.registerValidator(enode, proof.signature, {from: issuerAccount});
      let vals = await autonity.getValidators();
      assert.equal(vals.length, validators.length + 1, "validator pool is not expected");

      let v = await autonity.getValidator(newValAddr, {from: issuerAccount});

      const liquidABI = liquidContract["abi"]
      const liquid = new web3.eth.Contract(liquidABI, v.liquidContract);
      assert.equal(await liquid.methods.name().call(),"LNTN-"+(vals.length-1))
      assert.equal(await liquid.methods.symbol().call(),"LNTN-"+(vals.length-1))
      assert.equal(v.treasury.toString(), issuerAccount.toString(), "treasury addr is not expected");
      assert.equal(v.addr.toString(), newValAddr.toString(), "validator addr is not expected");
      assert.equal(v.enode.toString(), enode.toString(), "validator enode is not expected");
      assert(v.bondedStake == 0, "validator bonded stake is not expected");
      assert(v.totalSlashed == 0, "validator total slash counter is not expected");
      assert(v.state == 0, "validator state is not expected");
    });

    it('Pause a validator', async function () {
      let validator ='0xDE03B7806f885Ae79d2aa56568b77caDB0de073E';
      let issuerAccount = accounts[8];
      let enode = "enode://a7ecd2c1b8c0c7d7ab9cc12e620605a762865d381eb1bc5417dcf07599571f84ce5725f404f66d3e254d590ae04e4e8f18fe9e23cd29087d095a0c37d0443252@3.209.45.79:30303";
      let privateKey = 'e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b';

      let proof = web3.eth.accounts.sign(issuerAccount, privateKey);
      /* try disabling it with msg.sender not the treasury account, it should fails */
      try {
        let r = await autonity.pauseValidator(validator, {from: issuerAccount});
        assert.fail('Expected throw not received', r);
      } catch (e) {}

      /* disabling a non registered validator should fail */
      try {
        let r = await autonity.pauseValidator(validator, {from: accounts[7]});
        assert.fail('Expected throw not received', r);
      } catch (e) {}

      await autonity.registerValidator(enode, proof.signature, {from: issuerAccount});
      await autonity.pauseValidator(validator, {from: issuerAccount});
      let v = await autonity.getValidator(validator, {from: issuerAccount});
      assert(v.state == 1, "validator state is not expected");

      /* try disabling it again, it should fail */
      try {
        let r = await autonity.pauseValidator(validator, {from: issuerAccount});
        assert.fail('Expected throw not received', r);
      } catch (e) {}
    });

    it("Re-active a paused validator", async function () {
      let issuerAccount = accounts[8];

      let validator ='0xDE03B7806f885Ae79d2aa56568b77caDB0de073E';
      let enode = "enode://a7ecd2c1b8c0c7d7ab9cc12e620605a762865d381eb1bc5417dcf07599571f84ce5725f404f66d3e254d590ae04e4e8f18fe9e23cd29087d095a0c37d0443252@3.209.45.79:30303";
      let privateKey = 'e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b';
      /* activating a non-existing validator should fail */
      try {
        let r = await autonity.activateValidator(validator, {from: issuerAccount});
        assert.fail('Expected throw not received', r);
      } catch (e) {}

      let proof = web3.eth.accounts.sign(issuerAccount, privateKey);
      /* activating an already active validator should fail */
      await autonity.registerValidator(enode, proof.signature, {from: issuerAccount});
      try {
        let r = await autonity.activateValidator(validator, {from: issuerAccount});
        assert.fail('Expected throw not received', r);
      } catch (e) {}
      await autonity.pauseValidator(validator, {from: issuerAccount});
      let v = await autonity.getValidator(validator, {from: issuerAccount});
      assert(v.state == 1, "validator state is not expected");
      await autonity.activateValidator(validator, {from: issuerAccount});
      v = await autonity.getValidator(validator, {from: issuerAccount});
      assert(v.state == 0, "validator state is not expected");
    })
  });

  describe('Proposer election base on stake weighted sampling', function () {
    beforeEach(async function () {
      autonity = await utils.deployContract(validators, config, {from: deployer});
    });

    it('the election should be deterministic on same height and round.', async function () {
      let height;
      for (height = 0; height < 10; height++) {
        let round;
        for (round = 0; round < 3; round++) {
          let proposer1 = await autonity.getProposer(height, round);
          let proposer2 = await autonity.getProposer(height, round);
          assert(proposer1 === proposer2, "proposer election should be deterministic on same height and round");
        }
      }
    });
  });

  describe('Proposer selection, print and compare the scheduling rate with same stake.', function () {
    let stakes = [100, 100, 100, 100];
    beforeEach(async function () {
      let copyValidators = validators;
      for (let i = 0; i < copyValidators.length; i++) {
        copyValidators[i].bondedStake = stakes[i];
      }
      autonity = await utils.deployContract(copyValidators, config, {from: deployer});
    });

    it('get proposer, print and compare the scheduling rate with same stake.', async function () {
      let height;
      let maxHeight = 10000;
      let maxRound = 4;
      let expectedRatioDelta = 0.01;
      let counterMap = new Map();
      for (height = 0; height < maxHeight; height++) {
        let round;
        for (round = 0; round < maxRound; round++) {
          let proposer = await autonity.getProposer(height, round);
          if (counterMap.has(proposer) === true) {
            counterMap.set(proposer, counterMap.get(proposer) + 1)
          } else {
            counterMap.set(proposer, 1)
          }
        }
      }

      let totalStake = 0;
      stakes.forEach(function (v, index) {
        totalStake += v
      });

      validators.forEach(function (val, index) {
        let stake = stakes[index];
        let expectedRatio = stake / totalStake;
        let scheduled = counterMap.get(val.addr);
        let actualRatio = scheduled / (maxHeight * maxRound);
        let delta = Math.abs(expectedRatio - actualRatio);
        console.log("\t proposer: " + val.addr + " stake: " + stake + " was scheduled: " + scheduled + " times from " + maxHeight * maxRound + " times scheduling"
          + " expectedRatio: " + expectedRatio + " actualRatio: " + actualRatio + " delta: " + delta);

        if (delta > expectedRatioDelta) {
          assert.fail("Unexpected proposer scheduling rate delta.")
        }
      });
    });
  });

  describe('Proposer selection, print and compare the scheduling rate with liner increasing stake.', function () {
    let stakes = [100, 200, 400, 800];
    beforeEach(async function () {
      let copyValidators = validators;
      for (let i = 0; i < copyValidators.length; i++) {
        copyValidators[i].bondedStake = stakes[i];
      }
      autonity = await utils.deployContract(copyValidators, config, {from: deployer});
    });

    it('get proposer, print and compare the scheduling rate with same stake.', async function () {
      let maxHeight = 10000;
      let maxRound = 4;
      let expectedRatioDelta = 0.01;
      let counterMap = new Map();
      for (let height = 0; height < maxHeight; height++) {
        for (let round = 0; round < maxRound; round++) {
          let proposer = await autonity.getProposer(height, round);
          if (counterMap.has(proposer) === true) {
            counterMap.set(proposer, counterMap.get(proposer) + 1)
          } else {
            counterMap.set(proposer, 1)
          }
        }
      }

      let totalStake = 0;
      stakes.forEach(function (v, index) {
        totalStake += v
      });

      validators.forEach(function (val, index) {
        let stake = stakes[index];
        let expectedRatio = stake / totalStake;
        let scheduled = counterMap.get(val.addr);
        let actualRatio = scheduled / (maxHeight * maxRound);
        let delta = Math.abs(expectedRatio - actualRatio);
        console.log("\t proposer: " + val.addr + " stake: " + stake + " was scheduled: " + scheduled + " times from " + maxHeight * maxRound + " times scheduling"
          + " expectedRatio: " + expectedRatio + " actualRatio: " + actualRatio + " delta: " + delta);

        if (delta > expectedRatioDelta) {
          assert.fail("Unexpected proposer scheduling rate delta.")
        }

      });
    });
  });

  describe('Bonding and unbonding requests', function () {
    beforeEach(async function () {
      autonity = await utils.deployContract(validators, config, {from: deployer});
    });

    it('Bond to a valid validator', async function () {
      let newAccount = accounts[8];
      let tokenMint = 200;
      await autonity.mint(newAccount, tokenMint, {from: operator});
      // bond new minted Newton to a registered validator.
      await autonity.bond(validators[0].addr, tokenMint, {from: newAccount});
      // num of stakings from contract construction equals: length of validators and the latest bond.
      let numOfStakings = validators.length + 1;
      let stakings = await autonity.getBondingReq(0, numOfStakings);
      assert.equal(stakings[numOfStakings - 1].amount, tokenMint, "stake bonding amount is not expected");
      assert.equal(stakings[numOfStakings - 1].delegator, newAccount, "delegator addr is not expected");
      assert.equal(stakings[numOfStakings - 1].delegatee, validators[0].addr, "delegatee addr is not expected");
    });

    it('does not bond on a non-registered validator', async function () {
      // mint Newton for a new account.
      let newAccount = accounts[8];
      let tokenMint = 200;
      await autonity.mint(newAccount, tokenMint, {from: operator});
      // bond new minted Newton to a not registered validator.
      try {
        let r = await autonity.bond(anyAccount, tokenMint, {from: newAccount});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        // bonding should be failed, then the staking slot should not equal to the bonding meta data.
        let numOfStakings = validators.length + 1;
        let stakings = await autonity.getBondingReq(0, numOfStakings);
        assert.notEqual(stakings[numOfStakings - 1].amount, tokenMint, "stake bonding amount is not expected");
        assert.notEqual(stakings[numOfStakings - 1].delegator, newAccount, "delegator addr is not expected");
        assert.notEqual(stakings[numOfStakings - 1].delegatee, validators[0].addr, "delegatee addr is not expected");
      }
    });

    it('un-bond from a valid validator', async function () {
      let tokenUnBond = 10;
      let from = validators[0].addr;
      // unBond from self, a registered validator.
      await autonity.unbond(from, tokenUnBond, {from: from});
      let numOfUnBonding = 1;
      let unStakings = await autonity.getUnbondingReq(0, numOfUnBonding);
      assert.equal(unStakings[numOfUnBonding - 1].amount, tokenUnBond, "stake bonding amount is not expected");
      assert.equal(unStakings[numOfUnBonding - 1].delegator, from, "delegator addr is not expected");
      assert.equal(unStakings[numOfUnBonding - 1].delegatee, from, "delegatee addr is not expected");
    });

    it('does not unbond from not registered validator', async function () {
      let unRegisteredVal = anyAccount;
      let tokenUnBond = 10;
      try {
        let r = await autonity.unbond(unRegisteredVal, tokenUnBond, {from: validators[0].addr});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        // un-bonding should be failed, then the un-staking slot should not equal to the bonding meta data.
        let numOfUnStaking = 1;
        let unStakings = await autonity.getUnbondingReq(0, numOfUnStaking);
        assert.notEqual(unStakings[numOfUnStaking - 1].amount, tokenUnBond, "stake bonding amount is not expected");
        assert.notEqual(unStakings[numOfUnStaking - 1].delegator, validators[0].addr, "delegator addr is not expected");
        assert.notEqual(unStakings[numOfUnStaking - 1].delegatee, unRegisteredVal, "delegatee addr is not expected");
      }
    });

    it("can't unbond from  avalidator with the amount exceeding the available balance", async function () {
      let tokenUnBond = 99999;
      let from = validators[0].addr;
      try {
        let r = await autonity.unbond(from, tokenUnBond, {from: from});
        assert.fail('Expected throw not received', r);
      } catch (e) {
        // un-bonding should be failed, then the un-staking slot should not equal to the bonding meta data.
        let numOfUnStaking = 1;
        let unStakings = await autonity.getUnbondingReq(0, numOfUnStaking);
        assert.notEqual(unStakings[numOfUnStaking - 1].amount, tokenUnBond, "stake bonding amount is not expected");
        assert.notEqual(unStakings[numOfUnStaking - 1].delegator, from, "delegator addr is not expected");
        assert.notEqual(unStakings[numOfUnStaking - 1].delegatee, from, "delegatee addr is not expected");
      }
    });

    it("can't bond to a paused validator", async function () {
      await autonity.pauseValidator(validators[0].addr, {from: validators[0].addr});
      try {
        let r = await autonity.bond(validators[0].addr, 100, {from: validators[0].addr});
        assert.fail('Expected throw not received', r);
      } catch (e) {}
    });
  });


  // todo: fix below testcases.
  /*
  describe('Test apply bonding and un-bonding with epoch finalize()', function () {
      let vals = [
          {
              "treasury": accounts[0],
              "addr": accounts[0],
              "enode": "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
              "commissionRate": 10000,
              "bondedStake": 100,
              "selfBondedStake": 100,
              "totalSlashed": 0,
              "liquidContract": accounts[0],
              "liquidSupply": 0,
              "registrationBlock": 0,
              "state": 0,
          },
      ];
      let copyParams = config;
      beforeEach(async function () {
          // before deploy the contract, customized a shorten epoch period and the last Epoch block base on current
          // test chain's context.
          let currentHeight = await web3.eth.getBlockNumber();
          let customizedEpochPeriod = 20;
          copyParams.lastEpochBlock = currentHeight;
          copyParams.epochPeriod = customizedEpochPeriod;

          token = await utils.deployContract(vals, copyParams, {from: deployer});
          console.log("contract deployed with LastEpochBlock: ", currentHeight, "EpochPeriod: ", customizedEpochPeriod, "at address: ", token.address);
      });

      it('test bond stake token to new added validators, new validators should be elected as committee member', async function() {
          // register 2 new validators.
          let commissionRate = 10;
          let validator1 = accounts[1];
          let enodeVal1 = "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303";
          let validator2 = accounts[3];
          let enodeVal2 = "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303";

          await token.registerValidator(enodeVal1, {from: validator1});
          await token.registerValidator(enodeVal2, {from: validator2});

          // system operator mint Newton for user.
          let user = accounts[7];
          let tokenMint = 100;
          await token.mint(user, tokenMint, {from: operator});

          // user bond Newton to validator 2.
          await token.bond(validator2, tokenMint, {from: user});

          // 1st finalize with reward to be counted for the epoch.
          await token.finalize(0, {from: deployer});
          // during the period, try to call finalize with zero reward till the end block of epoch.
          for (;;) {
              // if last block of epoch updated, then the finalize() must be finished.
              let lastEpochBlock = await token.getLastEpochBlock();
              //let currentHeight = await web3.eth.getBlockNumber();
              //console.log("last epoch block: ", lastEpochBlock.toNumber(), "current height: ", currentHeight);
              if (lastEpochBlock > copyParams.lastEpochBlock) {
                  let committee = await token.getCommittee();
                  let presented = false;
                  for (let j=0; j<committee.length; j++) {
                      if (validator2 == committee[j].addr) {
                          presented = true;
                      }
                  }
                  assert.equal(presented, true);
                  break;
              } else {
                  token.finalize(0, {from: deployer});
              }
          }
      });
      it('test un-bond stake from validator, zero bonded validator should not be elected as committee member', async function() {
          // todo: write this case.
      });
  });

  describe('Test DPoS reward distribution with weighted staking', function () {
      // to make the test run faster, we use a shorten epoch period and customized lastEpochBlock only for test.
      let copyParams = config;
      beforeEach(async function () {
          // before deploy the contract, customized a shorten epoch period and the last Epoch block base on current
          // test chain's context.
          let currentHeight = await web3.eth.getBlockNumber();
          let customizedEpochPeriod = 20;
          copyParams.lastEpochBlock = currentHeight;
          copyParams.epochPeriod = customizedEpochPeriod;

          token = await utils.deployContract(validators, copyParams, {from: deployer});
          console.log("contract deployed with LastEpochBlock: ", currentHeight, "EpochPeriod: ", customizedEpochPeriod, "at address: ", token.address);
      });

      it('test finalize with none deployer account, exception should rise.', async function () {
          let blockReward = 100000;
          let omitted = [];
          try {
              let r = await token.finalize(blockReward, {from: anyAccount});
              assert.fail('Expected throw not received', r);
          } catch (e) {
          }
      });

      it('test finalize with no fund at autonity contract account, exception should rise.', async function () {
          let reward = 100000;
          // contract account should have no fund.
          let contractBalance = await web3.eth.getBalance(token.address);
          assert.equal(contractBalance, 0, "contract account have unexpected fund");

          // since we customized the epoch period and lastEpochBlock on the creation of this test target, a brand
          // new autonity contract was deployed over an brand new address at beforeEach(), the epoch rotation should
          // comes within 10 blocks, at height: lastEpochBlock + EpochPeriod. During the period, we try to call the
          // finalize function to count the reward with 100000 only once, then try to call the finalize function with
          // no reward over heights till the end block of epoch, and the last finalize function should rise an
          // exception that the contract has not sufficient fund to distribute the rewards.

          // 1st finalize with reward 100000 to be counted for the epoch.
          await token.finalize(reward, {from: deployer});
          let lastEpochBlock = copyParams.lastEpochBlock;
          let epochPeriod = copyParams.epochPeriod;
          let endBlock = lastEpochBlock + epochPeriod;
          try {
              // during the period, try to call finalize with zero reward till the end block of epoch, since there is
              // no fund in autnoity contract account, it fails finally.
              for (;;) {
                  let currentHeight = await web3.eth.getBlockNumber();
                  if (currentHeight <= endBlock) {
                      await token.finalize(0, {from: deployer});
                  } else {
                      assert.fail('Expected throw not received');
                  }
              }
          } catch (e) {
          }
      });

      it('test finalize with reward distribution over no delegations on the network', async function () {
          let reward = 1000000000000000;
          // contract account should have no fund.
          let initFund = await web3.eth.getBalance(token.address);
          assert.equal(initFund, 0, "contract account have unexpected fund");

          // load reward/fund to contract account, to get it distributed latter on.
          await web3.eth.sendTransaction({from: anyAccount, to: token.address, value: reward});
          let loadedBalance = await web3.eth.getBalance(token.address);
          assert.equal(loadedBalance, reward, "contract account have unexpected balance");

          let initBalanceV0 = await web3.eth.getBalance(validators[0].addr);
          console.log("v0, addr: ", validators[0].addr, "init balance: ", initBalanceV0);
          let initBalanceV1 = await web3.eth.getBalance(validators[1].addr);
          console.log("v1, addr: ", validators[1].addr, "init balance: ", initBalanceV1);
          let initBalanceV2 = await web3.eth.getBalance(validators[2].addr);
          console.log("v2, addr: ", validators[2].addr, "init balance: ", initBalanceV2);
          let initBalanceV3 = await web3.eth.getBalance(validators[3].addr);
          console.log("v3, addr: ", validators[3].addr, "init balance: ", initBalanceV3);

          let initBalanceTreasury = await web3.eth.getBalance(treasuryAccount);
          console.log("tr, addr: ", treasuryAccount, "init balance: ", initBalanceTreasury);

          // since we customized the epoch period and lastEpochBlock on the creation of this test target, a brand
          // new autonity contract was deployed over an brand new address at beforeEach(), the epoch rotation should
          // comes within 10 blocks, at height: lastEpochBlock + EpochPeriod. During the period, we try to call the
          // finalize function to count the none zero reward with only once, then try to call the finalize function with
          // no reward over heights till the end block of epoch, and the last finalize function of end block should
          // to distribute the rewards.

          // 1st finalize with reward to be counted for the epoch.
          await token.finalize(reward, {from: deployer});

          // during the period, try to call finalize with zero reward till the end block of epoch.
          for (;;) {
              // if last block of epoch updated, then the distribution must be finished.
              let lastEpochBlock = await token.getLastEpochBlock();
              let currentHeight = await web3.eth.getBlockNumber();
              console.log("last epoch block: ", lastEpochBlock.toNumber(), "current height: ", currentHeight);
              if (lastEpochBlock > copyParams.lastEpochBlock) {
                  // funds locked in autonity contract should be cleared.
                  let leftFund = await web3.eth.getBalance(token.address);
                  console.log("autonity contract init balance: ", initFund, "left balance: ", leftFund);
                  assert.equal(leftFund, initFund, "left fund in autonity contract is not expected");

                  let afterBalanceV0 = await web3.eth.getBalance(validators[0].addr);
                  console.log("v0, addr: ", validators[0].addr, "after balance: ", afterBalanceV0);
                  let afterBalanceV1 = await web3.eth.getBalance(validators[1].addr);
                  console.log("v1, addr: ", validators[1].addr, "after balance: ", afterBalanceV1);
                  let afterBalanceV2 = await web3.eth.getBalance(validators[2].addr);
                  console.log("v2, addr: ", validators[2].addr, "after balance: ", afterBalanceV2);
                  let afterBalanceV3 = await web3.eth.getBalance(validators[3].addr);
                  console.log("v3, addr: ", validators[3].addr, "after balance: ", afterBalanceV3);

                  let afterBalanceTreasury = await web3.eth.getBalance(treasuryAccount);
                  console.log("tr, addr: ", treasuryAccount, "after balance: ", afterBalanceTreasury);

                  let rewardV0 = web3.utils.toBN(afterBalanceV0).sub(web3.utils.toBN(initBalanceV0));
                  let rewardV1 = web3.utils.toBN(afterBalanceV1).sub(web3.utils.toBN(initBalanceV1));
                  let rewardV2 = web3.utils.toBN(afterBalanceV2).sub(web3.utils.toBN(initBalanceV2));
                  let rewardV3 = web3.utils.toBN(afterBalanceV3).sub(web3.utils.toBN(initBalanceV3));
                  let rewardTreasury = web3.utils.toBN(afterBalanceTreasury).sub(web3.utils.toBN(initBalanceTreasury));

                  let actualReward = rewardV0.add(rewardV1).add(rewardV2).add(rewardV3).add(rewardTreasury);

                  assert.equal(actualReward.toNumber(), reward, "total distributed reward not expected");
                  console.log("actual distributed reward:   ", actualReward.toNumber());
                  console.log("expected distributed reward: ", reward);
                  break;
              } else {
                  token.finalize(0, {from: deployer});
              }
          }
      });
      it('test finalize with reward distribution with user stake bonding', async function() {
          let reward = 1000000000000000;
          // contract account should have no fund.
          let initFund = await web3.eth.getBalance(token.address);
          assert.equal(initFund, 0, "contract account have unexpected fund");

          // load reward/fund to contract account, to get it distributed latter on.
          await web3.eth.sendTransaction({from: anyAccount, to: token.address, value: reward});
          let loadedBalance = await web3.eth.getBalance(token.address);
          assert.equal(loadedBalance, reward, "contract account have unexpected balance");

          // mint Newton for external users.
          let alice = accounts[7];
          let bob = accounts[9];
          await token.mint(alice, 400, {from: operator});
          await token.mint(bob, 400, {from: operator});

          // bond Newton in different validators.
          await token.bond(validators[0].addr, 100, {from: alice});
          await token.bond(validators[1].addr, 100, {from: bob});
          await token.bond(validators[2].addr, 100, {from: alice});
          await token.bond(validators[3].addr, 100, {from: bob});

          // apply the bonding at current epoch, and wait until the epoch is finalized.
          for (;;) {
              // if last block of epoch updated, then the distribution must be finished.
              let lastEpochBlock = await token.getLastEpochBlock();
              let currentHeight = await web3.eth.getBlockNumber();
              console.log("last epoch block: ", lastEpochBlock.toNumber(), "current height: ", currentHeight);
              if (lastEpochBlock > copyParams.lastEpochBlock) {
                  // check the bonded stake should grows according to the new bonding by Alice and Bob.
                  let val0 = await token.getValidator(validators[0].addr);
                  assert(val0.bondedStake == validators[0].bondedStake + 100, "bonded stake is not expected after bonding");
                  let val1 = await token.getValidator(validators[1].addr);
                  assert(val1.bondedStake == validators[1].bondedStake + 100, "bonded stake is not expected after bonding");
                  break;
              } else {
                  token.finalize(0, {from: deployer});
              }
          }

          // apply reward distribution at 2nd epoch, and wait until the 2nd epoch is finalized, and check the reward distribution.
          let initBalanceTreasury = await web3.eth.getBalance(treasuryAccount);
          let initBalanceV0 = await web3.eth.getBalance(validators[0].addr);
          let initBalanceV1 = await web3.eth.getBalance(validators[1].addr);
          let initBalanceV2 = await web3.eth.getBalance(validators[2].addr);
          let initBalanceV3 = await web3.eth.getBalance(validators[3].addr);
          let initBalanceAlice = await web3.eth.getBalance(alice);
          let initBalanceBob = await web3.eth.getBalance(bob);

          await token.finalize(reward, {from: deployer});
          for (;;) {
              // if last block of epoch updated, then the distribution must be finished.
              let lastEpochBlock = await token.getLastEpochBlock();
              let currentHeight = await web3.eth.getBlockNumber();
              console.log("last epoch block: ", lastEpochBlock.toNumber(), "current height: ", currentHeight);
              if (lastEpochBlock > (copyParams.lastEpochBlock+20)) {
                  // 2nd epoch is finalized, now check the correctness of reward distribution.
                  // the reward should go to the treasury accounts, and validator: 0, 1, 2, 3 and user: alice and bob.
                  let postBalanceTreasury = await web3.eth.getBalance(treasuryAccount);
                  let postBalanceV0 = await web3.eth.getBalance(validators[0].addr);
                  let postBalanceV1 = await web3.eth.getBalance(validators[1].addr);
                  let postBalanceV2 = await web3.eth.getBalance(validators[2].addr);
                  let postBalanceV3 = await web3.eth.getBalance(validators[3].addr);
                  let postBalanceAlice = await web3.eth.getBalance(alice);
                  let postBalanceBob = await web3.eth.getBalance(bob);

                  let rewardTreasury = web3.utils.toBN(postBalanceTreasury).sub(web3.utils.toBN(initBalanceTreasury));
                  let rewardV0 = web3.utils.toBN(postBalanceV0).sub(web3.utils.toBN(initBalanceV0));
                  let rewardV1 = web3.utils.toBN(postBalanceV1).sub(web3.utils.toBN(initBalanceV1));
                  let rewardV2 = web3.utils.toBN(postBalanceV2).sub(web3.utils.toBN(initBalanceV2));
                  let rewardV3 = web3.utils.toBN(postBalanceV3).sub(web3.utils.toBN(initBalanceV3));
                  let rewardAlice = web3.utils.toBN(postBalanceAlice).sub(web3.utils.toBN(initBalanceAlice));
                  let rewardBob = web3.utils.toBN(postBalanceBob).sub(web3.utils.toBN(initBalanceBob));
                  console.log("treasury reward: ", rewardTreasury.toNumber(), "v0 reward: ", rewardV0.toNumber(), "v1 reward: ", rewardV1.toNumber(),
                      "v2 reward: ", rewardV2.toNumber(), "v3 reward: ", rewardV3.toNumber(), "Alice reward: ", rewardAlice.toNumber(), "Bob reward: ", rewardBob.toNumber());

                  let actualReward = rewardV0.add(rewardV1).add(rewardV2).add(rewardV3).add(rewardTreasury).add(rewardAlice).add(rewardBob);

                  assert.equal(actualReward.toNumber(), reward, "total distributed reward not expected");
                  break;
              } else {
                  token.finalize(0, {from: deployer});
              }
          }
      });
  });*/
});