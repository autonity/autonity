'use strict';


const Autonity = artifacts.require('Autonity.sol');
const validatorsList = [
    '0x627306090abaB3A6e1400e9345bC60c78a8BEf57',
    '0xf17f52151EbEF6C7334FAD080c5704D77216b732',
    '0xC5fdf4076b8F3A5357c5E395ab970B5B54098Fef',
    '0x821aEa9a577a9b44299B9c15c88cf3087F3b5544',
    '0x0d1d4e623D10F9FBA5Db95830F7d3839406C6AF2'
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

    it('test add validator and check that it is in get validator list', async function() {
        const token = await Autonity.deployed();

        token.AddValidator(accounts[7], {from:accounts[0]})

        var getValidatorsResult = await token.GetValidators({from:accounts[0]});

        var expected = validatorsList;
        expected.push(accounts[7]);
        assert.deepEqual(getValidatorsResult, expected)
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
});
