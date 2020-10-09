'use strict';
const assert = require('assert');
const utils = require('./test-utils');
//todo: move gas analysis to separate js file

contract('Autonity', function (accounts) {
    const roleValidator = 2;
    const roleStakeHolder = 1;
    const roleParticipant = 0;
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
    const freeEnodes = [
        "enode://d860a01f9722d78051619d1e2351aba3f43f943f6f00718d1b9baa4101932a1f5011f16bb2b1bb35db20d6fe28fa0bf09636d26a87d31de9ec6203eeedb1f666@18.138.108.67:30303",
        "enode://22a8232c3abc76a16ae9d6c3b164f98775fe226f0917b0ca871128a74a8e9630b458460865bab457221f1d448dd9791d24c4e5d88786180ac185df813a68d4de@3.209.45.79:30303",
        "enode://ca6de62fce278f96aea6ec5a2daadb877e51651247cb96ee310a318def462913b653963c155a0ef6c7d50048bba6e6cea881130857413d9f50a621546b590758@34.255.23.113:30303",
        "enode://279944d8dcd428dffaa7436f25ca0ca43ae19e7bcf94a8fb7d1641651f92d121e972ac2e8f381414b80cc8e5555811c2ec6e1a99bb009b3f53c4c69923e11bd8@35.158.244.151:30303",
        "enode://8499da03c47d637b20eee24eec3c356c9a2e6148d6fe25ca195c7949ab8ec2c03e3556126b0d7ed644675e78c4318b08691b7b57de10e5f0d40d05b09238fa0a@52.187.207.27:30303"
    ];
    const userTypes = [2, 2, 2, 2, 2];
    const stakes = [100, 90, 80, 110, 120];
    const minGasPrice = 0;
    const committeeSize = 1000;
    const operator = accounts[0];
    const deployer = accounts[8];
    const version = "v1.0.0";
    let token;

    describe('Metrics', function() { // test failing

        beforeEach(async function() {
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version, {from: accounts[8]} );
        });

        it('test dump network Economic metric data.', async function () {
            let data = await token.dumpEconomicMetrics({from: operator});
            let minGasPrice = await token.getMinimumGasPrice({from: operator});
            let sum = 0;
            for (let i = 0; i < data.accounts.length; i++) {
                let stake = await token.balanceOf(data.accounts[i], {from: operator});
                assert.deepEqual(Number(data.stakes[i]), Number(stake));
                sum += Number(data.stakes[i])
            }

            assert.deepEqual(data.accounts, validatorsList);
            assert.deepEqual(Number(data.mingasprice), Number(minGasPrice))
            assert.deepEqual(Number(data.stakesupply), sum)
        });
    });

    describe('Initial state', function() {
        beforeEach(async function(){
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version, { from:accounts[8]} );
        });

        it('test validator can get initial validator list', async function () {
            let getValidatorsResult = await token.getValidators({from: operator});
            assert.deepEqual(getValidatorsResult, validatorsList);
        });

        it('test non validator account can get initial validator list', async function () {
            let getValidatorsResult = await token.getValidators({from: accounts[7]});
            assert.deepEqual(getValidatorsResult, validatorsList)
        });

        it('test default minimum gas price equals to 0', async function () {
            let minGasPrice = await token.getMinimumGasPrice({from: operator});
            assert(0 == minGasPrice, "default min gas price different of zero");
        });

        it('test validator can get the initial enodes whitelist', async function () {
            let enodesWhitelist = await token.getWhitelist({from: operator});
            assert.deepEqual(enodesWhitelist, whiteList);
        });
    });

    describe('Fee Distribution', function() {

        beforeEach(async function(){
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version, { from:accounts[0]} );
        });

        it('test redistribution fails with empty balance', async function () {
            try {
                await token.finalize(10000, {from: deployer});
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
                await token.finalize(10000, {from: accounts[0]});
                assert.fail('Expected throw not received', r);
            } catch (e) {

            }
        });

        it('test redistribution', async function () {
            let st = await token.getStakeholders({from: operator});

            for (let i = 0; i < st.length; i++) {
                await web3.eth.sendTransaction({from: st[i], to: token.address, value: 10000});
            }

            let balances = [];
            for (let i = 0; i < st.length; i++) {
                balances[i] = await web3.eth.getBalance(st[i]);
            }

            let performAmount = 10000;
            let totalStake= stakes.reduce((a,b) => a + b);
            let stakeholdersPart = stakes.map(element => element * performAmount / totalStake);

            await token.finalize(performAmount, {from: operator});

            let balancesAfter = [];
            for (let i = 0; i < st.length; i++) {
                balancesAfter[i] = await web3.eth.getBalance(st[i]);
            }

            for (let i = 0; i < st.length; i++) {
                let check = web3.utils.toBN(balancesAfter[i])
                    .sub(web3.utils.toBN(balances[i]))
                    .eq(web3.utils.toBN(stakeholdersPart[i]));

                assert(check, "not equal")
            }
        });
    });

    describe('Governance - System Operator', function() {

        beforeEach(async function(){
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version, { from:accounts[8]} );
        });

        it('test Governance operator can add/remove to whitelist', async function () {
            let enode = freeEnodes[0];
            let tx = await token.addUser(accounts[8], 20, enode, roleValidator, {from: operator});
            // console.log("\tGas used to add validator to whitelist = " + tx.receipt.gasUsed.toString() + " gas");
            let getValidatorsResult = await token.getWhitelist({from: operator});
            let expected = whiteList.slice();
            expected.push(enode);
            assert.deepEqual(getValidatorsResult, expected);

            tx = await token.removeUser(accounts[8], {from: operator});
            // console.log("\tGas used to remove val from whitelist = " + tx.receipt.gasUsed.toString() + " gas");
            getValidatorsResult = await token.getWhitelist({from: accounts[1]});
            assert.deepEqual(getValidatorsResult, whiteList);
        });

        it('test add validator and check that it is in get validator list', async function () {
            let expected = validatorsList.slice();
            expected.push(accounts[7]);
            let tx = await token.addUser(accounts[7], 100, freeEnodes[0], roleValidator, {from: operator});
            // console.log("\tGas used to add new validator = " + tx.receipt.gasUsed.toString() + " gas");
            let getValidatorsResult = await token.getValidators({from: operator});
            assert.deepEqual(expected, getValidatorsResult);

            tx = await token.removeUser(accounts[7], {from: operator});
            // console.log("\tGas used to remove validator = " + tx.receipt.gasUsed.toString() + " gas");
            getValidatorsResult = await token.getValidators({from: operator});
            assert.deepEqual(validatorsList, getValidatorsResult)
        });


        it('test validator cannot call change user type function', async function() {
            // Upgrades
            // test that a validator can't call the changeUserType function
            try {
              await token.addUser(accounts[6], 0, freeEnodes[0], roleParticipant, {from: operator});
              await token.addUser(accounts[7], 100, freeEnodes[0], roleValidator, {from: operator});
              await token.changeUserType(accounts[6], 1, {from: accounts[6]});
              assert.fail('Expected throw not received');
            } catch (e) {
              await token.removeUser(accounts[6], {from: operator});
              await token.removeUser(accounts[7], {from: operator});
            }
          });

          it('test upgrades to userType', async function() {
            // participant -> stakeholder (0 -> 1)
            await token.addUser(accounts[6], 0, freeEnodes[0], roleParticipant, {from: operator});
            await token.changeUserType(accounts[6], roleStakeHolder, {from: operator});
            let thisUserType = (await token.getUser(accounts[6],{from: accounts[6]})).userType;
            assert (thisUserType == roleStakeHolder, "wrong user type");
            await token.removeUser(accounts[6], {from: operator});

            // participant -> stakeholder -> validator (0 -> 1 -> 2)
            try {
                await token.addUser(accounts[6], 0, freeEnodes[0], roleParticipant, {from: operator});
                await token.changeUserType(accounts[6], roleValidator, {from: operator});
                assert.fail('Expected throw not received');
            } catch (e) {
                // user type shouldn't change since 0 stake of validator is not permitted.
                thisUserType = (await token.getUser(accounts[6], {from: accounts[6]})).userType;
                assert (thisUserType == roleParticipant, "wrong user type");
            }
            // upgrade to stakeholder and mint stake for it.
            await token.changeUserType(accounts[6], roleStakeHolder, {from: operator});
            await token.mint(accounts[6], 20, {from: operator});
            // upgrade to validator
            await token.changeUserType(accounts[6], roleValidator, {from: operator});
            thisUserType = (await token.getUser(accounts[6], {from: accounts[6]})).userType;
            assert (thisUserType == roleValidator, "wrong user type");
            let thisUserStake = await token.balanceOf(accounts[6], {from: accounts[6]});
            assert (thisUserStake == 20);
            await token.removeUser(accounts[6], {from: operator});

            // stakeholder -> validator (1 -> 2)
            await token.addUser(accounts[6], 100, freeEnodes[0], roleStakeHolder, {from: operator});
            await token.changeUserType(accounts[6], roleValidator, {from: operator});
            thisUserType = (await token.getUser(accounts[6], {from: accounts[6]})).userType;
            assert (thisUserType == roleValidator, "wrong user type");
            thisUserStake = await token.balanceOf(accounts[6], {from: accounts[6]});
            assert (thisUserStake == 100);
            await token.removeUser(accounts[6], {from: operator});
          });

          it('test downgrades to userType', async function() {
            // valiator -> stakeholder (2 -> 1)
            await token.addUser(accounts[6], 100, freeEnodes[0], roleValidator, {from: operator});
            await token.changeUserType(accounts[6], roleStakeHolder, {from: operator});
            let thisUserType = (await token.getUser(accounts[6], {from: accounts[6]})).userType;
            assert (thisUserType == roleStakeHolder, "wrong user type");
            let thisUserStake = await token.balanceOf(accounts[6], {from: accounts[6]});
            assert (thisUserStake == 100);
            await token.removeUser(accounts[6], {from: operator});

            // validator -> participant (2 -> 0)
            try {
              // ensure that a validator with stake cannot be downgraded
              await token.addUser(accounts[6], 100, freeEnodes[0], roleValidator, {from: operator});
              await token.changeUserType(accounts[6], roleParticipant, {from: operator});
              assert.fail('Expected throw not received');
            } catch (e) {
              await token.removeUser(accounts[6], {from: operator});
              await token.addUser(accounts[6], 10, freeEnodes[0], roleValidator, {from: operator});
              await token.changeUserType(accounts[6], roleStakeHolder, {from: operator});
              await token.redeem(accounts[6], 10, {from: operator});
              await token.changeUserType(accounts[6], roleParticipant, {from: operator});
              thisUserType = (await token.getUser(accounts[6], {from: accounts[6]})).userType;
              assert (thisUserType == roleParticipant, "wrong user type");
              await token.removeUser(accounts[6], {from: operator});
            }

            // stakeholder -> participant (1 -> 0)
            try {
              // ensure that a participant with stake cannot be downgraded
              await token.addUser(accounts[6], 100, freeEnodes[0], roleStakeHolder, {from: operator});
              await token.changeUserType(accounts[6], roleParticipant, {from: operator});
              assert.fail('Expected throw not received');
            } catch (e) {
              await token.removeUser(accounts[6], {from: operator});
              await token.addUser(accounts[6], 0, freeEnodes[0], roleStakeHolder, {from: operator});
              await token.changeUserType(accounts[6], roleParticipant, {from: operator});
              thisUserType = (await token.getUser(accounts[6], {from: accounts[6]})).userType;
              assert (thisUserType == roleParticipant, "wrong user type");
              await token.removeUser(accounts[6], {from: operator});
            }

        });

        it('test create participant account check it and remove it', async function () {
            let tx = await token.addUser(accounts[9], 0, freeEnodes[0], roleParticipant, {from: operator});
            //console.log("\tGas used to add participant = " + tx.receipt.gasUsed.toString() + " gas");
            let addMemberResult = (await token.getUser(accounts[9]));

            assert(accounts[9] === addMemberResult.addr);

            tx = await token.removeUser(accounts[9], {from: operator});
            // console.log("\tGas used to remove participant = " + tx.receipt.gasUsed.toString() + " gas");
            let removeMemberResult = (await token.getUser(accounts[9]));

            assert("0x0000000000000000000000000000000000000000" === removeMemberResult.addr);
        });

        it('test non validator cannot add validator', async function () {
            let enode = freeEnodes[0];
            try {
                let r = await token.addUser(accounts[7], 20, enode, roleValidator, {from: accounts[6]})

            } catch (e) {
                let getValidatorsResult = await token.getValidators({from: operator});
                assert.deepEqual(getValidatorsResult, validatorsList);
                return
            }

            assert.fail('Expected throw not received');
        });

        it('test that _createUser() does not allow duplicates', async function () {
            try {
              await token._createUser(accounts[6], freeEnodes[0], roleValidator, 100, 0, {from: operator});
              // the duplicate
              await token._createUser(accounts[6], freeEnodes[0], roleValidator, 100, 0, {from: operator});
              assert.fail('Expected throw not received');
            } catch (e) {
              return
            }
        });

        it('test non Governance operator cannot add validator', async function () {
            let enode = freeEnodes[0];
            try {
                let r = await token.addUser(accounts[6], 20, enode, roleValidator, {from: accounts[6]});
                assert.fail('Expected throw not received');

            } catch (e) {
                //skip error
                let getWhitelistResult = await token.getWhitelist({from: operator});
                assert.deepEqual(getWhitelistResult, whiteList);
            }
        });

        it('test non Governance operator cannot remove user', async function () {
            let enode = whiteList[0];
            try {
                let r = await token.removeUser(accounts[2], {from: accounts[6]});
                assert.fail('Expected throw not received', r);

            } catch (e) {
                let getWhitelistResult = await token.getWhitelist({from: operator});
                assert.deepEqual(getWhitelistResult, whiteList);
            }
        });

        it('test create/remove participants by non governance operator', async function () {
            let errorOnAddNewMember = false;
            let errorOnRemoveMember = false;

            try {
                await token.addUser(accounts[8], 0, freeEnodes[0], roleParticipant, {from: accounts[7]});
            } catch (e) {
                errorOnAddNewMember = true
            }
            let addMemberResult = await token.getUser(accounts[8]);
            // we check for the 0x0 address to see if the user is no longer saved in the contract.
            // in the future, we might want to return an error.
            assert("0x0000000000000000000000000000000000000000" === addMemberResult.addr);

            await token.addUser(accounts[8], 0, freeEnodes[0], roleParticipant, {from: operator});

            addMemberResult = await token.getUser(accounts[8]);
            assert(accounts[8] === addMemberResult.addr);


            try {
                await token.removeUser(accounts[8], {from: accounts[7]});
            } catch (e) {
                errorOnRemoveMember = true
            }

            let removeMemberResult = await token.getUser(accounts[8]);
            assert(accounts[8] === removeMemberResult.addr);
            await token.removeUser(accounts[8], {from: operator});

            removeMemberResult = await token.getUser(accounts[8]);
            assert("0x0000000000000000000000000000000000000000" === removeMemberResult.addr);

            assert(true === errorOnAddNewMember);
            assert(true === errorOnRemoveMember);

        });
    });

    describe('Committee Selection', function () {

        beforeEach(async function(){
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version,  { from:accounts[8]} );
        });

        it('test set max committee size by operator account', async function () {

            await token.setCommitteeSize(4, {from: operator});
            let maxCommitteeSize = await token.getMaxCommitteeSize();
            assert(4 == maxCommitteeSize, "maxCommittee size was not set correctly");

        });

        it('test regular validator cannot set the max committee size', async function() {

            let initMaxCommitteeSize = await token.getMaxCommitteeSize();

            try {
                let r = await token.setCommitteeSize(50, {from: accounts[6]});
                assert.fail('Expected throw not received', r);

            } catch (e) {
                let maxCommitteeSize = await token.getMaxCommitteeSize();
                assert.deepEqual(initMaxCommitteeSize,  maxCommitteeSize);
            }
        });

        it('test set committee when committee size is equal than the number of validators', async function() {

            let validators = await token.getValidators();
            let committeeSize = validators.length;

            await token.setCommitteeSize(committeeSize, {from: operator});
            let maxCommitteeSize = await token.getMaxCommitteeSize();
            assert(committeeSize == maxCommitteeSize, "maxCommittee size was not set correctly");

            try {
                let r = await token.setCommittee({from: deployer});
                assert.fail('Expected throw not received', r);
            } catch (e) {

            }
            await token.computeCommittee({from: deployer});
            let committeeResult = await token.getCommittee();
            let committeeValidators = [];

            for (let i = 0; i < committeeResult.length; i++) {
                committeeValidators.push(committeeResult[i][0])
            }

            assert.deepEqual(committeeValidators.slice().sort(), validators.slice().sort(), "Committee should be equal than validator set");

        });

        it('test set committee when committee size is greater than the number of validators', async function() {
            let validators = await token.getValidators();
            let committeeSize = validators.length + 5;

            await token.setCommitteeSize(committeeSize, {from: operator});
            let maxCommitteeSize = await token.getMaxCommitteeSize();
            assert(committeeSize == maxCommitteeSize, "maxCommittee size was not set correctly");

            try {
                let r = await token.setCommittee({from: deployer});
                assert.fail('Expected throw not received', r);
            } catch (e) {

            }
            await token.computeCommittee({from: deployer});
            let committeeResult = await token.getCommittee();
            let committeeValidators = [];

            for (let i = 0; i < committeeResult.length; i++) {
                committeeValidators.push(committeeResult[i][0])
            }

            assert.deepEqual(committeeValidators.slice().sort(), validators.slice().sort(), "Committee should be equal than validator set");

        });

        it('test set committee when committee size is smaller than the number of validators', async function() {

            try {
                let validators = await token.getValidators();
                let committeeSize = validators.length - 2;

                await token.setCommitteeSize(committeeSize, {from: operator});
                let maxCommitteeSize = await token.getMaxCommitteeSize();
                assert(committeeSize == maxCommitteeSize, "maxCommittee size was not set correctly");
                await token.computeCommittee({from: deployer});

                let r = await token.setCommittee({from: deployer});
                assert.fail('Expected throw not received', r);

                let committeeResult = await token.getCommittee();
                let committeeValidators = [];

                for (let i = 0; i < r.length; i++) {
                    committeeValidators.push(r[i][0])
                }

                // Mock committee selection
                let indexesToBeRemoved = [1,2]
                while(indexesToBeRemoved.length) {
                    validators.splice(indexesToBeRemoved.pop(), 1);
                }
                assert.deepEqual(committeeValidators.slice().sort(), validators.slice().sort(), "Error while creating new committee");
            }catch (e) {

            }

        });
    });

    describe('Stake Token', function() {

        beforeEach(async function(){
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version,  { from:accounts[8]} );
        });

        it('test user type downgraded when all stake redeemed', async function f() {
            let initStake = 100;
            await token.addUser(accounts[8], initStake, freeEnodes[0], roleValidator, {from: operator});
            let initialUserType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert(initialUserType == roleValidator, "wrong user type");

            await token.redeem(accounts[8], initStake, {from: operator});
            let newUserType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert (newUserType == roleStakeHolder, "wrong user type");
        });

        it('test user type not downgraded when not all stake redeemed', async function f() {
            let initStake = 100;
            await token.addUser(accounts[8], initStake, freeEnodes[0], roleValidator, {from: operator});
            let initialUserType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert(initialUserType == roleValidator, "wrong user type");
            await token.redeem(accounts[8], 50, {from: operator});

            let userType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert (userType == roleValidator, "wrong user type");
        });

        it('test user type downgraded when all stake transferred', async function f() {
            let initStake = 100;
            await token.addUser(accounts[8], initStake, freeEnodes[0], roleValidator, {from: operator});
            let initialUserType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert(initialUserType == roleValidator, "wrong user type");
            await token.transfer(accounts[2], initStake, {from: accounts[8]});
            let userType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert (userType == roleStakeHolder, "wrong user type");
        });

        it('test user type not downgraded when not all stake transferred', async function f() {
            let initStake = 100;
            await token.addUser(accounts[8], initStake, freeEnodes[0], roleValidator, {from: operator});
            let initialUserType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert(initialUserType == roleValidator, "wrong user type");
            await token.transfer(accounts[2], 50, {from: accounts[8]});
            let userType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert (userType == roleValidator, "wrong user type");
        });

        it('test user type downgrade of last validator in the network', async function f() {
            let initStake = 100;
            await token.addUser(accounts[8], initStake, freeEnodes[0], roleValidator, {from: operator});
            let initialUserType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
            assert(initialUserType == roleValidator, "wrong user type");

            let balance = await token.balanceOf(accounts[1], {from: accounts[1]});
            await token.redeem(accounts[1], balance, {from: operator});

            balance = await token.balanceOf(accounts[2], {from: accounts[2]});
            await token.redeem(accounts[2], balance, {from: operator});

            balance = await token.balanceOf(accounts[3], {from: accounts[3]});
            await token.redeem(accounts[3], balance, {from: operator});

            balance = await token.balanceOf(accounts[4], {from: accounts[4]});
            await token.redeem(accounts[4], balance, {from: operator});

            balance = await token.balanceOf(accounts[5], {from: accounts[5]});
            await token.redeem(accounts[5], balance, {from: operator});

            let getValidatorsResult = await token.getValidators({from: operator});
            assert(getValidatorsResult.length == 1, "wrong number of validators");

            try {
                await token.redeem(accounts[8], initStake, {from: operator});
                assert.fail('Expected throw not received');
            } catch (e) {
                let userType = (await token.getUser(accounts[8], {from: accounts[8]})).userType;
                assert(userType == roleValidator, "wrong user type");
                let getValidatorsResult = await token.getValidators({from: operator});
                assert(getValidatorsResult.length == 1, "wrong number of validators");
            }
        });

        it('test create account, add stake, check that it is added, remove stake', async function () {
            await token.addUser(accounts[7], 0, freeEnodes[0], roleStakeHolder, {from: operator});
            let getStakeResult = await token.balanceOf(accounts[7], {from: accounts[7]});
            assert(0 == getStakeResult, "unexpected tokens");

            let tx = await token.mint(accounts[7], 100, {from: operator});
            // console.log("\tGas used to mint stake = " + tx.receipt.gasUsed.toString() + " gas");
            getStakeResult = await token.balanceOf(accounts[7], {from: accounts[7]});
            assert(100 == getStakeResult, "tokens are not minted");
            tx = await token.redeem(accounts[7], 100, {from: operator});

            // console.log("\tGas used to redeem stake = " + tx.receipt.gasUsed.toString() + " gas");

            getStakeResult = await token.balanceOf(accounts[7], {from: accounts[7]});
            assert(0 == getStakeResult, "unexpected tokens");

            await token.removeUser(accounts[7], {from: operator});
        });

        it('test create account, get error when redeem empty stake', async function () {

            await token.addUser(accounts[7], 0, freeEnodes[0], roleStakeHolder, {from: operator});
            let getStakeResult = await token.balanceOf(accounts[7], {from: accounts[7]});
            assert(0 == getStakeResult, "unexpected tokens not minted");

            try {
                await token.redeem(accounts[7], 100, {from: operator});
                assert.fail('Expected throw not received');
            } catch (e) {
                getStakeResult = await token.balanceOf(accounts[7], {from: accounts[7]});
                assert(0 == getStakeResult, "unexpected tokens");
                await token.removeUser(accounts[7], {from: operator});
            }
        });

        it('test transfer stake', async function () {
            let getStakeResult = await token.balanceOf(validatorsList[2], {from: validatorsList[2]});
            assert(stakes[2] == getStakeResult, "unexpected tokens");

            await token.redeem(accounts[2], 50, {from: operator});
            let balance_after = await token.balanceOf(validatorsList[2], {from: validatorsList[2]});
            assert(stakes[2] - 50, balance_after);

            await token.removeUser(accounts[4], {from: operator});
            await token.removeUser(accounts[5], {from: operator});
        });

        it("should give accounts[2] authority to spend account[1]'s token", async function() {
            await token.addStakeholder(accounts[6], freeEnodes[0], 100, {from: operator});
            await token.addStakeholder(accounts[7], freeEnodes[1], 0, {from: operator});
            await token.addStakeholder(accounts[8], freeEnodes[2], 0, {from: operator});

            await token.approve(accounts[7], 70,  {from: accounts[6]});
            let allowance = await token.allowance(accounts[6], accounts[7]);
            assert.equal(allowance, 70, 'allowance is wrong');
            await token.transferFrom(accounts[6], accounts[8], 40, {from: accounts[7]});
            let recipientBalance = await token.balanceOf(accounts[8]);
            assert.equal(recipientBalance, 40, 'accounts[8] balance is wrong');
            let senderBalance = await token.balanceOf(accounts[6]);
            assert.equal(senderBalance, 60, 'accounts[6] balance is wrong');
            let spenderBalance = await token.balanceOf(accounts[7]);
            assert.equal(spenderBalance, 0, 'accounts[7] balance is wrong');
            let remainingAllowance = await token.allowance(accounts[6], accounts[7]);
            assert.equal(remainingAllowance, 30, 'final allowance is wrong');

            await token.removeUser(accounts[6], {from: operator});
            await token.removeUser(accounts[7], {from: operator});
            await token.removeUser(accounts[8], {from: operator});
        });

        it('test validator can get users list', async function () {
            var getValidatorsResult = await token.retrieveState({from: operator});
            let addresses = getValidatorsResult[0];
            let types = getValidatorsResult[1];
            let stake = getValidatorsResult[2];
            let enodes = getValidatorsResult[3];
            let commisionRate = getValidatorsResult[4];
            assert.deepEqual(getValidatorsResult[0], validatorsList);

            await token.addUser(accounts[7], 0, freeEnodes[0], roleStakeHolder, {from: operator});
            let getStakeResult = await token.balanceOf(accounts[7], {from: accounts[7]});
            assert(0 == getStakeResult, "unexpected tokens");

            await token.addUser(accounts[6], 0, freeEnodes[0], roleStakeHolder, {from: operator});
            getStakeResult = await token.balanceOf(accounts[6], {from: accounts[6]});
            assert(0 == getStakeResult, "unexpected tokens");

            await token.mint(accounts[7], 100, {from: operator});

            getStakeResult = await token.balanceOf(accounts[7], {from: accounts[7]});
            assert(100 == getStakeResult, "tokens are not minted");

            let tx = await token.transfer(accounts[6], 50, {from: accounts[7]});
            // console.log("\tGas used to send state token = " + tx.receipt.gasUsed.toString() + " gas");

            getStakeResult = await token.balanceOf(accounts[7], {from: accounts[7]});
            assert(50 == getStakeResult, "unexpected tokens");

            getStakeResult = await token.balanceOf(accounts[6],{from: accounts[6]});
            assert(50 == getStakeResult, "unexpected tokens");

            await token.redeem(accounts[7], 50, {from: operator});
            await token.redeem(accounts[5], 50, {from: operator});
            await token.removeUser(accounts[7], {from: operator});
            await token.removeUser(accounts[5], {from: operator});
        });
    });

    describe('Proposer selection, Normal case.', function() {

        beforeEach(async function(){
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version,  { from:accounts[8]} );
        });

        it('get proposer, proposer should be determinated across nodes on same height and round.', async function () {
            await token.computeCommittee({from: deployer});
            let height;
            for (height = 0; height < 10; height++) {
                let round;
                for (round = 0; round < 3; round ++){
                    let proposer1 = await token.getProposer(height, round);
                    let proposer2 = await token.getProposer(height, round);
                    assert(proposer1 === proposer2, "proposer should be determinated on same height and round")
                }
            }
        });

    });

    describe('Proposer selection, print and compare the scheduling rate with same stake.', function() {
        let stakes = [100, 100, 100, 100, 100];
        beforeEach(async function(){
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version,  { from:accounts[8]} );
        });

        it('get proposer, print and compare the scheduling rate with same stake.', async function () {
            await token.computeCommittee({from: deployer});
            let height;
            let maxHeight = 10000;
            let maxRound = 4;
            let expectedRatioDelta = 0.01;
            let counterMap = new Map();
            for (height = 0; height < maxHeight; height++) {
                let round;
                for (round = 0; round < maxRound; round ++){
                    let proposer = await token.getProposer(height, round);
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

            validatorsList.forEach(function (addr, index) {
                let stake = stakes[index];
                let expectedRatio = stake / totalStake;
                let scheduled = counterMap.get(addr);
                let actualRatio = scheduled / (maxHeight * maxRound);
                let delta = Math.abs(expectedRatio - actualRatio);
                console.log("\t proposer: " + addr + " stake: " + stake + " was scheduled: " + scheduled + " times of " + maxHeight*maxRound + " times scheduling"
                 + " expectedRatio: " + expectedRatio + " actualRatio: " + actualRatio + " delta: " + delta);

                if (delta > expectedRatioDelta) {
                    assert.fail("Unexpected proposer scheduling rate delta.")
                }

            });
        });
    });

    describe('Proposer selection, print and compare the scheduling rate with different stake.', function() {
        let stakes = [100, 200, 400, 800, 1600];
        beforeEach(async function(){
            token = await utils.deployContract(validatorsList, whiteList,
                userTypes, stakes, operator, minGasPrice, committeeSize, version,  { from:accounts[8]} );
        });

        it('get proposer, print and compare the scheduling rate with same stake.', async function () {
            await token.computeCommittee({from: deployer});
            let height;
            let maxHeight = 10000;
            let maxRound = 4;
            let expectedRatioDelta = 0.01;
            let counterMap = new Map();
            for (height = 0; height < maxHeight; height++) {
                let round;
                for (round = 0; round < maxRound; round ++){
                    let proposer = await token.getProposer(height, round);
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

            validatorsList.forEach(function (addr, index) {
                let stake = stakes[index];
                let expectedRatio = stake / totalStake;
                let scheduled = counterMap.get(addr);
                let actualRatio = scheduled / (maxHeight * maxRound);
                let delta = Math.abs(expectedRatio - actualRatio);
                console.log("\t proposer: " + addr + " stake: " + stake + " was scheduled: " + scheduled + " times of " + maxHeight*maxRound + " times scheduling"
                    + " expectedRatio: " + expectedRatio + " actualRatio: " + actualRatio + " delta: " + delta);

                if (delta > expectedRatioDelta) {
                    assert.fail("Unexpected proposer scheduling rate delta.")
                }

            });
        });
    });

});
