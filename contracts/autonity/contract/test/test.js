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

    it('test validator can get validator list', async function () {
        const token = await Autonity.deployed();

        var getValidatorsResult = await token.getValidators({from: governanceOperatorAccount});
        assert.deepEqual(getValidatorsResult, validatorsList);
    });

    it('test validator can get users list', async function () {
        const token = await Autonity.deployed();

        var getValidatorsResult = await token.getUsers({from: governanceOperatorAccount});
        let addresses = getValidatorsResult[0];
        let types = getValidatorsResult[1];
        let stake = getValidatorsResult[2];
        let enodes = getValidatorsResult[3];
        // assert.deepEqual(getValidatorsResult, validatorsList);
        var a = Autonity.new(addresses, enodes, types, stake,accounts[0], 0, { from:accounts[8]});
        let b = await a;
        console.log(b);

    });


    it('test redistribution fails with empty balance', async function () {
        const token = await Autonity.deployed();

        try {
            await token.performRedistribution(10000, {from: deployer});
            assert.fail('Expected throw not received', r);
        } catch (e) {

        }
    });

    it('test redistribution fails with not deployer', async function () {
        const token = await Autonity.deployed();
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
        const token = await Autonity.deployed();
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


    it('test non validator can get validator list', async function () {
        const token = await Autonity.deployed();
        var getValidatorsResult = await token.getValidators({from: accounts[7]});

        assert.deepEqual(getValidatorsResult, validatorsList)
    });


    it('test non validator cant add validator', async function () {
        const token = await Autonity.deployed();

        try {
            var r = await token.addValidator(accounts[7], {from: accounts[6]})

        } catch (e) {
            var getValidatorsResult = await token.getValidators({from: governanceOperatorAccount});
            assert.deepEqual(getValidatorsResult, validatorsList);
            return
        }

        assert.fail('Expected throw not received');
    });


    it('test add validator and check that it is in get validator list', async function () {
        const token = await Autonity.deployed();
        let expected = validatorsList.slice();
        expected.push(accounts[7]);
        await token.addValidator(accounts[7], 100, "not nil enode", {from: governanceOperatorAccount});

        var getValidatorsResult = await token.getValidators({from: governanceOperatorAccount});
        assert.deepEqual(expected, getValidatorsResult);

        await token.removeUser(accounts[7], {from: governanceOperatorAccount});
        getValidatorsResult = await token.getValidators({from: governanceOperatorAccount});
        assert.deepEqual(validatorsList, getValidatorsResult)
    });


    it('test non Governance operator cant remove user', async function () {
        const token = await Autonity.deployed();
        var enode = whiteList[0];
        try {
            let r = await token.removeUser(accounts[2], {from: accounts[6]});
            assert.fail('Expected throw not received', r);

        } catch (e) {
            var getWhitelistResult = await token.getWhitelist({from: governanceOperatorAccount});
            assert.deepEqual(getWhitelistResult, whiteList);
        }
    });

    it('test non Governance operator cant add validator', async function () {
        const token = await Autonity.deployed();
        var enode = "enode://testenode";
        try {
            var r = await token.addValidator(accounts[6], 20, enode, {from: accounts[6]});
            assert.fail('Expected throw not received');

        } catch (e) {
            //skip error
            var getWhitelistResult = await token.getWhitelist({from: governanceOperatorAccount});
            assert.deepEqual(getWhitelistResult, whiteList);
        }
    });


    it('test Governance operator can add/remove to whitelist', async function () {
        const token = await Autonity.deployed();
        var enode = "enode://testenode";
        await token.addValidator(accounts[8], 20, enode, {from: governanceOperatorAccount});

        var getValidatorsResult = await token.getWhitelist({from: governanceOperatorAccount});
        let expected = whiteList.slice();
        expected.push(enode);
        assert.deepEqual(getValidatorsResult, expected);

        await token.removeUser(accounts[8], {from: governanceOperatorAccount});
        getValidatorsResult = await token.getWhitelist({from: accounts[1]});
        assert.deepEqual(getValidatorsResult, whiteList);
    });


    it('test create participant account check it and remove it', async function () {
        const token = await Autonity.deployed();

        await token.addParticipant(accounts[9], "some enode", {from: governanceOperatorAccount});
        var addMemberResult = await token.checkMember(accounts[9]);

        assert(true == addMemberResult);

        await token.removeUser(accounts[9], {from: governanceOperatorAccount});
        var removeMemberResult = await token.checkMember(accounts[9]);

        assert(false == removeMemberResult);
    });


    it('test create account, add stake, check that it is added, remove stake', async function () {
        const token = await Autonity.deployed();

        await token.addStakeholder(accounts[7], "some enode", 0, {from: governanceOperatorAccount});
        var getStakeResult = await token.getStake({from: accounts[7]});
        assert(0 == getStakeResult, "unexpected tokens");

        await token.mintStake(accounts[7], 100, {from: governanceOperatorAccount});

        getStakeResult = await token.getStake({from: accounts[7]});
        assert(100 == getStakeResult, "tokens are not minted");

        await token.redeemStake(accounts[7], 100, {from: governanceOperatorAccount});

        getStakeResult = await token.getStake({from: accounts[7]});
        assert(0 == getStakeResult, "unexpected tokens");

        await token.removeUser(accounts[7], {from: governanceOperatorAccount});
    });


    it('test create account, get error when redeem empty stake', async function () {
        const token = await Autonity.deployed();

        await token.addStakeholder(accounts[7], "some enode", 0, {from: governanceOperatorAccount});
        var getStakeResult = await token.getStake({from: accounts[7]});
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

    it('test create/remove participants by non governance operator', async function () {
        const token = await Autonity.deployed();
        var errorOnAddNewMember = false;
        var errorOnRemoveMember = false;

        try {
            await token.addParticipant(accounts[8], "some enode", {from: accounts[7]});
        } catch (e) {
            errorOnAddNewMember = true
        }
        var addMemberResult = await token.checkMember(accounts[8]);
        assert(false == addMemberResult);

        await token.addParticipant(accounts[8], "some enode", {from: governanceOperatorAccount});

        addMemberResult = await token.checkMember(accounts[8]);
        assert(true == addMemberResult);


        try {
            await token.removeUser(accounts[8], {from: accounts[7]});
        } catch (e) {
            errorOnRemoveMember = true
        }

        var removeMemberResult = await token.checkMember(accounts[8]);
        assert(true == removeMemberResult);
        await token.removeUser(accounts[8], {from: governanceOperatorAccount});

        removeMemberResult = await token.checkMember(accounts[8]);
        assert(false == removeMemberResult);

        assert(true == errorOnAddNewMember);
        assert(true == errorOnRemoveMember);

    });

    it('test transfer stake', async function () {
        const token = await Autonity.deployed();

        await token.addStakeholder(accounts[7], "some enode", 0, {from: governanceOperatorAccount});
        var getStakeResult = await token.getStake({from: accounts[7]});
        assert(0 == getStakeResult, "unexpected tokens");

        await token.addStakeholder(accounts[5], "some enode", 0, {from: governanceOperatorAccount});
        var getStakeResult = await token.getStake({from: accounts[5]});
        assert(0 == getStakeResult, "unexpected tokens");

        await token.mintStake(accounts[7], 100, {from: governanceOperatorAccount});

        getStakeResult = await token.getStake({from: accounts[7]});
        assert(100 == getStakeResult, "tokens are not minted");

        await token.send(accounts[5], 50, {from: accounts[7]});

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
