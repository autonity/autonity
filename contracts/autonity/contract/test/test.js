'use strict';


const Autonity = artifacts.require('Autonity.sol');
const validatorsList = [
    '0x627306090abaB3A6e1400e9345bC60c78a8BEf57',
    '0xf17f52151EbEF6C7334FAD080c5704D77216b732',
    '0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef',
    '0x821aEa9a577a9b44299B9c15c88cf3087F3b5544',
    '0x0d1d4e623D10F9FBA5Db95830F7d3839406C6AF2'
];

const whiteList = [
    "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
    "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
    "enode://e766ac390e2d99b559aef773c3656fa8d50df2310496ac26ca6c3fc84e21dabb8a0162cc8e34f938d45e0a8ed04955f8ddf1c380182f8ef17a3f08885064505f@172.25.0.13:30303",
    "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
    "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
];

contract('Autonity', function(accounts) {

    it('test validator can get validator list', async function() {
        const token = await Autonity.deployed();

        var getValidatorsResult = await token.GetValidators({from:accounts[0]});
        assert.deepEqual(getValidatorsResult, validatorsList);
    });

    it('test non validator can get validator list', async function() {
        const token = await Autonity.deployed();
        var getValidatorsResult = await token.GetValidators({from:accounts[7]});

        assert.deepEqual(getValidatorsResult, validatorsList)
    });


    it('test non validator cant add validator', async function() {
        const token = await Autonity.deployed();

        try {
            var r =await token.AddValidator(accounts[7], {from:accounts[6]})

        } catch (e) {
            //skip error
        }

        var getValidatorsResult = await token.GetValidators({from:accounts[0]});
        assert.deepEqual(getValidatorsResult, validatorsList)
    });

    it('test add validator and check that it is in get validator list', async function() {
        const token = await Autonity.deployed();

        token.AddValidator(accounts[7], {from:accounts[0]})

        var getValidatorsResult = await token.GetValidators({from:accounts[0]});

        var expected = validatorsList;
        expected.push(accounts[7]);
        assert.deepEqual(getValidatorsResult, expected)

        token.RemoveValidator(accounts[7], {from:accounts[0]})
        getValidatorsResult = await token.GetValidators({from:accounts[0]});
        assert.deepEqual(getValidatorsResult, validatorsList)
    });




    it('test non Governance operator cant add to whitelist', async function() {
        const token = await Autonity.deployed();
        var enode = whiteList[0];
        try {
            var r =await token.RemoveEnode(enode, {from:accounts[6]})

        } catch (e) {
            //skip error
        }

        var getWhitelistResult = await token.getWhitelist({from:accounts[0]});
        assert.deepEqual(getWhitelistResult, whiteList)
    });

    it('test non Governance operator cant remove from whitelist', async function() {
        const token = await Autonity.deployed();
        var enode = "enode://testenode";
        try {
            var r =await token.AddEnode(enode, {from:accounts[6]})

        } catch (e) {
            //skip error
        }

        var getWhitelistResult = await token.getWhitelist({from:accounts[0]});
        assert.deepEqual(getWhitelistResult, whiteList)
    });


    it('test Governance operator can add/remove to whitelist', async function() {
        const token = await Autonity.deployed();
        var enode = "enode://testenode";
        await token.AddEnode(enode, {from:accounts[0]});

        var getValidatorsResult = await token.getWhitelist({from:accounts[0]});
        var expected = whiteList;
        expected.push(enode);
        assert.deepEqual(getValidatorsResult, expected);

        await token.RemoveEnode(enode,{from:accounts[0]});
        var getValidatorsResult = await token.getWhitelist({from:accounts[1]});
        assert.deepEqual(getValidatorsResult, whiteList);
    });
});
