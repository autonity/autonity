'use strict';

//it's kind of e2e test. They have sequence dependency
//todo get rid of it.

const Autonity = artifacts.require('Autonity.sol');
// const validatorsList = [
//     '0x627306090abaB3A6e1400e9345bC60c78a8BEf57',
//     '0xf17f52151EbEF6C7334FAD080c5704D77216b732',
//     '0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef',
//     '0x821aEa9a577a9b44299B9c15c88cf3087F3b5544',
//     '0x0d1d4e623D10F9FBA5Db95830F7d3839406C6AF2'
// ];
//
// const whiteList = [
//     "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
//     "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
//     "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
//     "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
//     "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
// ];
//
// const governanceOperatorAccount = "0x627306090abaB3A6e1400e9345bC60c78a8BEf57";
// const deployer = "0x0F4F2Ac550A1b4e2280d04c21cEa7EBD822934b5";

contract('Autonity', function (accounts) {
    const validatorsList = [
        accounts[1],
        accounts[2],
        accounts[3],
        accounts[4],
        accounts[5],
    ];

    const whiteList = [
        "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
        "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
        "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
        "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
        "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
    ];

    const governanceOperatorAccount = accounts[0];
    const deployer = accounts[8];
    let token;

    describe('Initial state', function() {

        beforeEach(async function(){
            token = await Autonity.deployed();
        })

        it('test validator can get initial validator list', async function () {
            let getValidatorsResult = await token.getValidators({from: governanceOperatorAccount});
            assert.deepEqual(getValidatorsResult, validatorsList);
        });

        it('test non validator account can get initial validator list', async function () {
            let getValidatorsResult = await token.getValidators({from: accounts[7]});
            assert.deepEqual(getValidatorsResult, validatorsList)
        });

        it('test default minimum gas price equals to 0', async function () {
            let minGasPrice = await token.getMinimumGasPrice({from: governanceOperatorAccount});
            assert(0 == minGasPrice, "default min gas price different of zero");
        });

        it('test validator can get the initial enodes whitelist', async function () {
            let enodesWhitelist = await token.getWhitelist({from: governanceOperatorAccount});
            assert.deepEqual(enodesWhitelist, whiteList);
        });

    });

    describe('Fee Distribution', function() {

        beforeEach(async function(){
            token = await Autonity.deployed();
        })

        it('test redistribution fails with empty balance', async function () {
            try {
                await token.performRedistribution(10000, {from: deployer});
                assert.fail('Expected throw not received', r);
            } catch (e) {

            }
        });

        it('test redistribution fails with not deployer', async function () {
            let st = await token.getStakeholders();

            for (let i = 0; i < st.length; i++) {
                await web3.eth.sendTransaction({from: st[i], to: token.address, value: 10000});
            }
            let balance = await web3.eth.getBalance(token.address);
            if (balance < 10000) {
                assert.fail("incorrect balance")
            }
            try {
                await token.performRedistribution(10000, {from: accounts[0]});
                assert.fail('Expected throw not received', r);
            } catch (e) {

            }
        });

        it('test redistribution', async function () {
            let st = await token.getStakeholders();
            let balances = [];
            for (let i = 0; i < st.length; i++) {
                balances[i] = await web3.eth.getBalance(st[i]);
            }

            let performAmount = 10000;
            let stackholdersPart = performAmount / st.length;

            await token.performRedistribution(performAmount, {from: deployer});

            let balancesAfter = [];
            for (let i = 0; i < st.length; i++) {
                balancesAfter[i] = await web3.eth.getBalance(st[i]);
            }

            for (let i = 0; i < st.length; i++) {
                let check = web3.utils.toBN(balancesAfter[i])
                    .sub(web3.utils.toBN(balances[i]))
                    .eq(web3.utils.toBN(stackholdersPart));

                assert(check, "not equal")
            }

        });
    });

    describe('Governance - System Operator', function() {

        beforeEach(async function(){
            token = await Autonity.deployed();
        })

        it('test Governance operator can add/remove to whitelist', async function () {
            let enode = "enode://testenode";
            let tx = await token.addValidator(accounts[8], 20, enode, {from: governanceOperatorAccount});
            console.log("\tGas used to add validator to whitelist = " + tx.receipt.gasUsed.toString() + " gas");
            let getValidatorsResult = await token.getWhitelist({from: governanceOperatorAccount});
            let expected = whiteList.slice();
            expected.push(enode);
            assert.deepEqual(getValidatorsResult, expected);

            tx = await token.removeUser(accounts[8], {from: governanceOperatorAccount});
            console.log("\tGas used to remove val from whitelist = " + tx.receipt.gasUsed.toString() + " gas");
            getValidatorsResult = await token.getWhitelist({from: accounts[1]});
            assert.deepEqual(getValidatorsResult, whiteList);
        });

        it('test add validator and check that it is in get validator list', async function () {
            let expected = validatorsList.slice();
            expected.push(accounts[7]);
            let tx = await token.addValidator(accounts[7], 100, "not nil enode", {from: governanceOperatorAccount});
            console.log("\tGas used to add new validator = " + tx.receipt.gasUsed.toString() + " gas");
            let getValidatorsResult = await token.getValidators({from: governanceOperatorAccount});
            assert.deepEqual(expected, getValidatorsResult);

            tx = await token.removeUser(accounts[7], {from: governanceOperatorAccount});
            console.log("\tGas used to remove validator = " + tx.receipt.gasUsed.toString() + " gas");
            getValidatorsResult = await token.getValidators({from: governanceOperatorAccount});
            assert.deepEqual(validatorsList, getValidatorsResult)
        });

        it('test create participant account check it and remove it', async function () {

            let tx = await token.addParticipant(accounts[9], "some enode", {from: governanceOperatorAccount});
            console.log("\tGas used to add participant = " + tx.receipt.gasUsed.toString() + " gas");

            let addMemberResult = await token.checkMember(accounts[9]);

            assert(true === addMemberResult);

            tx = await token.removeUser(accounts[9], {from: governanceOperatorAccount});
            console.log("\tGas used to remove participant = " + tx.receipt.gasUsed.toString() + " gas");
            let removeMemberResult = await token.checkMember(accounts[9]);

            assert(false === removeMemberResult);
        });

        it('test non validator cannot add validator', async function () {

            try {
                let r = await token.addValidator(accounts[7], {from: accounts[6]})

            } catch (e) {
                let getValidatorsResult = await token.getValidators({from: governanceOperatorAccount});
                assert.deepEqual(getValidatorsResult, validatorsList);
                return
            }

            assert.fail('Expected throw not received');
        });

        it('test non Governance operator cannot add validator', async function () {
            let enode = "enode://testenode";
            try {
                let r = await token.addValidator(accounts[6], 20, enode, {from: accounts[6]});
                assert.fail('Expected throw not received');

            } catch (e) {
                //skip error
                let getWhitelistResult = await token.getWhitelist({from: governanceOperatorAccount});
                assert.deepEqual(getWhitelistResult, whiteList);
            }
        });

        it('test non Governance operator cannot remove user', async function () {
            let enode = whiteList[0];
            try {
                let r = await token.removeUser(accounts[2], {from: accounts[6]});
                assert.fail('Expected throw not received', r);

            } catch (e) {
                let getWhitelistResult = await token.getWhitelist({from: governanceOperatorAccount});
                assert.deepEqual(getWhitelistResult, whiteList);
            }
        });

        it('test create/remove participants by non governance operator', async function () {
            let errorOnAddNewMember = false;
            let errorOnRemoveMember = false;

            try {
                await token.addParticipant(accounts[8], "some enode", {from: accounts[7]});
            } catch (e) {
                errorOnAddNewMember = true
            }
            let addMemberResult = await token.checkMember(accounts[8]);
            assert(false === addMemberResult);

            await token.addParticipant(accounts[8], "some enode", {from: governanceOperatorAccount});

            addMemberResult = await token.checkMember(accounts[8]);
            assert(true === addMemberResult);


            try {
                await token.removeUser(accounts[8], {from: accounts[7]});
            } catch (e) {
                errorOnRemoveMember = true
            }

            let removeMemberResult = await token.checkMember(accounts[8]);
            assert(true === removeMemberResult);
            await token.removeUser(accounts[8], {from: governanceOperatorAccount});

            removeMemberResult = await token.checkMember(accounts[8]);
            assert(false === removeMemberResult);

            assert(true === errorOnAddNewMember);
            assert(true === errorOnRemoveMember);

        });
    });

    describe('Stake Token', function() {

        beforeEach(async function(){
            token = await Autonity.deployed();
        })

        it('test create account, add stake, check that it is added, remove stake', async function () {

            await token.addStakeholder(accounts[7], "some enode", 0, {from: governanceOperatorAccount});
            let getStakeResult = await token.getStake({from: accounts[7]});
            assert(0 == getStakeResult, "unexpected tokens");

            let tx = await token.mintStake(accounts[7], 100, {from: governanceOperatorAccount});
            console.log("\tGas used to mint stake = " + tx.receipt.gasUsed.toString() + " gas");

            getStakeResult = await token.getStake({from: accounts[7]});
            assert(100 == getStakeResult, "tokens are not minted");

            tx = await token.redeemStake(accounts[7], 100, {from: governanceOperatorAccount});

            console.log("\tGas used to redeem stake = " + tx.receipt.gasUsed.toString() + " gas");

            getStakeResult = await token.getStake({from: accounts[7]});
            assert(0 == getStakeResult, "unexpected tokens");

            await token.removeUser(accounts[7], {from: governanceOperatorAccount});
        });

        it('test create account, get error when redeem empty stake', async function () {

            await token.addStakeholder(accounts[7], "some enode", 0, {from: governanceOperatorAccount});
            let getStakeResult = await token.getStake({from: accounts[7]});
            assert(0 == getStakeResult, "unexpected tokens not minted");

            try {
                await token.redeemStake(accounts[7], 100, {from: governanceOperatorAccount});
                assert.fail('Expected throw not received');
            } catch (e) {
                getStakeResult = await token.getStake({from: accounts[7]});
                assert(0 == getStakeResult, "unexpected tokens");
                await token.removeUser(accounts[7], {from: governanceOperatorAccount});
            }
        });

        it('test transfer stake', async function () {

            await token.addStakeholder(accounts[7], "some enode", 0, {from: governanceOperatorAccount});
            let getStakeResult = await token.getStake({from: accounts[7]});
            assert(0 == getStakeResult, "unexpected tokens");

            await token.addStakeholder(accounts[5], "some enode", 0, {from: governanceOperatorAccount});
            getStakeResult = await token.getStake({from: accounts[5]});
            assert(0 == getStakeResult, "unexpected tokens");

            await token.mintStake(accounts[7], 100, {from: governanceOperatorAccount});

            getStakeResult = await token.getStake({from: accounts[7]});
            assert(100 == getStakeResult, "tokens are not minted");

            let tx = await token.send(accounts[5], 50, {from: accounts[7]});
            console.log("\tGas used to send state token = " + tx.receipt.gasUsed.toString() + " gas");

            getStakeResult = await token.getStake({from: accounts[7]});
            assert(50 == getStakeResult, "unexpected tokens");

            getStakeResult = await token.getStake({from: accounts[5]});
            assert(50 == getStakeResult, "unexpected tokens");

            await token.redeemStake(accounts[7], 50, {from: governanceOperatorAccount});
            await token.redeemStake(accounts[5], 50, {from: governanceOperatorAccount});
            await token.removeUser(accounts[7], {from: governanceOperatorAccount});
            await token.removeUser(accounts[5], {from: governanceOperatorAccount});
        });
    });

    it('test dump network Economic metric data.', async function () {
        const token = await Autonity.deployed();
        let data = await token.dumpEconomicsMetricData();
        let minGasPrice = await token.getMinimumGasPrice();
        let sum = 0;
        for (let i = 0; i < data.accounts.length; i++) {
            let stake = await token.getAccountStake(data.accounts[i]);
            assert.deepEqual(Number(data.stakes[i]), Number(stake));
            sum += Number(data.stakes[i])
        }

        assert.deepEqual(data.accounts, validatorsList);
        assert.deepEqual(Number(data.mingasprice), Number(minGasPrice))
        assert.deepEqual(Number(data.stakesupply), sum)
    });

    describe('Metrics', function() {

        beforeEach(async function(){
            token = await Autonity.deployed();
        })

        it('test dump network Economic metric data.', async function () {
            let data = await token.dumpEconomicsMetricData({from: governanceOperatorAccount});
            let minGasPrice = await token.getMinimumGasPrice({from: governanceOperatorAccount});
            let sum = 0;
            for (let i = 0; i < data.accounts.length; i++) {
                let stake = await token.getAccountStake(data.accounts[i], {from: governanceOperatorAccount});
                assert.deepEqual(Number(data.stakes[i]), Number(stake));
                sum += Number(data.stakes[i])
            }

            assert.deepEqual(data.accounts, validatorsList);
            assert.deepEqual(Number(data.mingasprice), Number(minGasPrice))
            assert.deepEqual(Number(data.stakesupply), sum)
        });
    });

});
