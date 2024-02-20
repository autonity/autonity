'use strict';
const assert = require('assert');
const { Buffer } = require('node:buffer');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const liquidContract = artifacts.require("Liquid")
const AccountabilityTest = artifacts.require("AccountabilityTest")
const config = require('./config.js')

contract('Autonity', function (accounts) {

    for (let i = 0; i < accounts.length; i++) {
        console.log("account: ", i, accounts[i]);
    }

    const operator = accounts[5];
    const deployer = accounts[6];
    const anyAccount = accounts[7];
    const treasuryAccount = accounts[8];

    const accountabilityConfig = config.ACCOUNTABILITY_CONFIG
    const genesisEnodes = config.GENESIS_ENODES
    const genesisNodeAddresses = config.GENESIS_NODE_ADDRESSES
    const baseValidator = config.BASE_VALIDATOR
    const genesisPrivateKeys = config.GENESIS_PRIVATE_KEYS
    let autonityConfig = config.autonityConfig(operator, treasuryAccount)
    let validators = config.validators(accounts)

    let autonity;
    let accountability;

    describe('Validator management', function () {
        let enode;
        let node;
        let oracle;
        let consensusKey;
        let pop;
        beforeEach(async function () {
            autonity = await utils.deployContracts(validators, autonityConfig, accountabilityConfig, deployer, operator);

            const nodeKeyInfo = await utils.generateAutonityKeys(`./autonity/data/test0.key`)
            const oracleKey = utils.randomPrivateKey();
            const popInfo = await utils.generateAutonityPOP(`./autonity/data/test0.key`, oracleKey, accounts[8])
            enode = utils.publicKeyToEnode(nodeKeyInfo.nodePublicKey.substring(2))
            node = nodeKeyInfo.nodeAddress
            oracle = utils.address(utils.publicKey(oracleKey, false))
            consensusKey = Buffer.from(nodeKeyInfo.nodeConsensusKey.substring(2), 'hex');
            pop = Buffer.from(popInfo.signatures.substring(2), 'hex')
        });

        it('Add validator with already registered address', async function () {
            let treasury = accounts[8];
            let enode = genesisEnodes[0]
            // multisig length is checked before validator already registered (it is not verified though)
            let multisig = utils.generateMultiSig(genesisPrivateKeys[0],genesisPrivateKeys[0],treasury)

            let consensusKey = Buffer.from('845681310fe66ed10629e76cc5aa20f3ec8b853af9f3dee8a6318f3fb81c0adcaaa0a776dc066127e743bba6b0349bc0', 'hex');
            let consensusKeyProof = '0x88a19caac1d02d2efb3675ec9fe99936b1170641b03d7525674ee001446cfd204fa5ba0b5e362d71294f3ba2f758695115a17101fc70b73fe90d7eb83950c3f7ad598b6740698b8e78fb48821c47762cdf2de889deede80fe2e7c085e48562c4';
            multisig = multisig + consensusKeyProof.substring(2);

            await truffleAssert.fails(
                autonity.registerValidator(enode, genesisNodeAddresses[0], consensusKey, multisig, {from: treasury}),
                truffleAssert.ErrorType.REVERT,
                "validator already registered"
            );

            let vals = await autonity.getValidators();
            assert.equal(vals.length, validators.length, "validator pool is not expected");
        });

        it('Add a validator with invalid enode address', async function () {
            let treasury = accounts[8];
            let enode = "enode://invalidEnodeAddress@172.25.0.11:30303";
            await truffleAssert.fails(
                autonity.registerValidator(enode, oracle, consensusKey, pop, {from: treasury}),
                truffleAssert.ErrorType.REVERT,
                "enode error"
            );

            let vals = await autonity.getValidators();
            assert.equal(vals.length, validators.length, "validator pool is not expected");
        });

        it('Add a validator with invalid oracle proof', async function () {
            let treasury = accounts[8];
            // set a wrong oracle address on purpose.
            let oracleAddr = treasury
            await truffleAssert.fails(
                autonity.registerValidator(enode, oracleAddr, consensusKey, pop, {from: treasury}),

                truffleAssert.ErrorType.REVERT,
                "Invalid oracle key ownership proof provided"
            );

            let vals = await autonity.getValidators();
            assert.equal(vals.length, validators.length, "validator pool is not expected");
        });

        it('Add a validator with valid meta data', async function () {
            let treasury = accounts[8];
            await autonity.registerValidator(enode, oracle, consensusKey, pop, {from: treasury});
            let vals = await autonity.getValidators();
            assert.equal(vals.length, validators.length + 1, "validator pool is not expected");

            let v = await autonity.getValidator(node, {from: treasury});

            const liquidABI = liquidContract["abi"]
            const liquid = new web3.eth.Contract(liquidABI, v.liquidContract);
            assert.equal(await liquid.methods.name().call(),"LNTN-"+(vals.length-1))
            assert.equal(await liquid.methods.symbol().call(),"LNTN-"+(vals.length-1))

            assert.equal(v.treasury.toString(), treasury.toString(), "treasury addr is not expected");
            assert.equal(v.nodeAddress.toString(), node.toString(), "validator addr is not expected");
            assert.equal(v.enode.toString(), enode.toString(), "validator enode is not expected");
            assert(v.bondedStake == 0, "validator bonded stake is not expected");
            assert(v.totalSlashed == 0, "validator total slash counter is not expected");
            assert(v.state == 0, "validator state is not expected");
        });

        it('Pause a validator', async function () {
            let treasury = accounts[8];
            // disabling a non registered validator should fail
            await truffleAssert.fails(
                autonity.pauseValidator(node, {from: treasury}),
                truffleAssert.ErrorType.REVERT,
                "validator must be registered"
            );
            await autonity.registerValidator(enode, oracle, consensusKey, pop, {from: treasury});

            // try disabling it with msg.sender not the treasury account, it should fails
            await truffleAssert.fails(
                autonity.pauseValidator(node, {from: accounts[7]}),
                truffleAssert.ErrorType.REVERT,
                "require caller to be validator admin account"
            );

            await autonity.pauseValidator(node, {from: treasury});
            let v = await autonity.getValidator(node, {from: treasury});
            assert(v.state == 1, "validator state is not expected");

            // try disabling it again, it should fail
            await truffleAssert.fails(
                autonity.pauseValidator(node, {from: treasury}),
                truffleAssert.ErrorType.REVERT,
                "validator must be active"
            );
        });

        it("Re-active a paused validator", async function () {
            let treasury = accounts[8];
            // activating a non-existing validator should fail
            await truffleAssert.fails(
                autonity.activateValidator(node, {from: treasury}),
                truffleAssert.ErrorType.REVERT,
                "validator must be registered"
            );

            await autonity.registerValidator(enode, oracle, consensusKey, pop, {from: treasury});
            // activating from non-treasury account should fail
            await truffleAssert.fails(
                autonity.activateValidator(node, {from: accounts[7]}),
                truffleAssert.ErrorType.REVERT,
                "require caller to be validator treasury account"
            );

            // activating an already active validator should fail
            await truffleAssert.fails(
                autonity.activateValidator(node, {from: treasury}),
                truffleAssert.ErrorType.REVERT,
                "validator already active"
            );
            await autonity.pauseValidator(node, {from: treasury});
            let v = await autonity.getValidator(node, {from: treasury});
            assert(v.state == 1, "validator state is not expected");
            await autonity.activateValidator(node, {from: treasury});
            v = await autonity.getValidator(node, {from: treasury});
            assert(v.state == 0, "validator state is not expected");
        })
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
        let enode1;
        let enode2;
        let node1;
        let node2;
        let oracle1;
        let oracle2;
        let consensusKey1;
        let consensusKey2;
        let pop1;
        let pop2;
        beforeEach(async function () {
            // set short epoch period
            let customizedEpochPeriod = 20;
            copyParams.protocol.epochPeriod = customizedEpochPeriod;

            autonity = await utils.deployContracts(vals, copyParams, accountabilityConfig, deployer, operator);
            assert.equal(customizedEpochPeriod,(await autonity.getEpochPeriod()).toNumber());

            const nodeKeyInfo = await utils.generateAutonityKeys(`./autonity/data/test1.key`)
            const oracleKey = utils.randomPrivateKey();
            const popInfo = await utils.generateAutonityPOP(`./autonity/data/test1.key`, oracleKey, accounts[1])
            enode1 = utils.publicKeyToEnode(nodeKeyInfo.nodePublicKey.substring(2))
            node1 = nodeKeyInfo.nodeAddress
            oracle1 = utils.address(utils.publicKey(oracleKey, false))
            consensusKey1 = Buffer.from(nodeKeyInfo.nodeConsensusKey.substring(2), 'hex');
            pop1 = Buffer.from(popInfo.signatures.substring(2), 'hex')

            const nodeKeyInfo2 = await utils.generateAutonityKeys(`./autonity/data/test2.key`)
            const oracleKey2 = utils.randomPrivateKey();
            const popInfo2 = await utils.generateAutonityPOP(`./autonity/data/test2.key`, oracleKey2, accounts[3])
            enode2 = utils.publicKeyToEnode(nodeKeyInfo2.nodePublicKey.substring(2))
            node2 = nodeKeyInfo2.nodeAddress
            oracle2 = utils.address(utils.publicKey(oracleKey2, false))
            consensusKey2 = Buffer.from(nodeKeyInfo2.nodeConsensusKey.substring(2), 'hex');
            pop2 = Buffer.from(popInfo2.signatures.substring(2), 'hex')
        });

        it('test bond stake token to new added validators, new validators should be elected as committee member', async function() {
            // register 2 new validators.
            let treasury1 = accounts[1];
            let treasury2 = accounts[3];

            await autonity.registerValidator(enode1, oracle1, consensusKey1, pop1, {from: treasury1});
            await autonity.registerValidator(enode2, oracle2, consensusKey2, pop2, {from: treasury2});

            // system operator mint Newton for user.
            let user = accounts[7];
            let tokenMint = 100;

            await autonity.mint(user, tokenMint, {from: operator});

            // user bond Newton to node 2.
            await autonity.bond(node2, tokenMint, {from: user});

            // close epoch to ensure bonding is applied
            await utils.endEpoch(autonity,operator,deployer);

            let committee = await autonity.getCommittee();
            let presented = false;
            for (let j=0; j<committee.length; j++) {
                if (node2 == committee[j].addr) {
                    presented = true;
                }
                // we should not find the 0 bonded stake new validator
                if (node1 == committee[j].addr) {
                    assert.fail("found unexpected committee member")
                }
            }
            assert.equal(presented, true);
        });
        it('test un-bond stake from validator, zero bonded validator should not be elected as committee member', async function() {
            // register 2 new validators.
            let treasury1 = accounts[1];
            let treasury2 = accounts[3];
            await autonity.registerValidator(enode1, oracle1, consensusKey1, pop1, {from: treasury1});
            await autonity.registerValidator(enode2, oracle2, consensusKey2, pop2, {from: treasury2});

            // system operator mint Newton for user.
            let user = accounts[7];
            let tokenMint = 100;

            await autonity.mint(user, tokenMint, {from: operator});

            // bond NTN to the 2 validators
            await autonity.bond(node1, 20, {from: user});
            await autonity.bond(node2, 20, {from: user});

            // close epoch to ensure bonding is applied
            await utils.endEpoch(autonity,operator,deployer);

            let committee = await autonity.getCommittee();
            assert.equal(committee.length,3)

            await autonity.unbond(node1, 20, {from: user});

            // voting power gets reduced right away at epoch end
            await utils.endEpoch(autonity,operator,deployer);
            committee = await autonity.getCommittee();
            assert.equal(committee.length,2)

            let presented = false;
            for (let j=0; j<committee.length; j++) {
                // we should not find the validator we unbonded from
                if (node1 == committee[j].addr) {
                    assert.fail("found unexpected committee member")
                }
                // we should find the other one
                if (node2 == committee[j].addr) {
                    presented = true;
                }
            }
            assert.equal(presented, true);
        });
        it('test more than committeeSize bonded validators, the ones with less stake should remain outside of the committee', async function() {
            // re-deploy with 4 validators instead of 1
            autonity = await utils.deployContracts(validators, copyParams, accountabilityConfig, deployer, operator);

            // set committeeSize to 0, minimum stake validator should be excluded
            await autonity.setCommitteeSize(3, {from: operator});
            await autonity.computeCommittee({from:deployer})

            let minimumStakeAddress;
            let minimumStake = Number.MAX_VALUE
            for (let i=0; i<validators.length; i++) {
                if(validators[i].bondedStake < minimumStake) {
                    minimumStake = validators[i].bondedStake
                    minimumStakeAddress = validators[i].nodeAddress
                }
            }
            let committee = await autonity.getCommittee({from: anyAccount})

            assert.equal(committee.length,3)
            for (let i=0; i<committee.length; i++) {
                assert.notEqual(committee[i].addr,minimumStakeAddress)
            }

        });
    });

    describe('After effects of slashing, ', function () {
        beforeEach(async function () {
            autonity = await utils.deployAutonityTestContract(validators, autonityConfig, accountabilityConfig,  deployer, operator);
            accountability = await AccountabilityTest.new(autonity.address, accountabilityConfig, {from: deployer});
            await autonity.setAccountabilityContract(accountability.address, {from:operator});
        });
        it('does not trigger fairness issue (unbondingStake > 0 and delegatedStake > 0)', async function () {
            // fairness issue is triggered when delegatedStake or unbondingStake becomes 0 from positive due to slashing
            // it can happen due to slashing rate = 100%
            // it should not happen for slashing amount < totalStake
            let config = JSON.parse(JSON.stringify(accountabilityConfig));
            // modifying config so we get slashingAmount = totalStake - 1, the highest slash possible without triggering fairness issue
            const expectedBondedStake = parseInt(config.slashingRatePrecision);
            const expectedSlash = expectedBondedStake - 1;
            config.collusionFactor = expectedSlash - parseInt(config.baseSlashingRateMid);
            accountability = await AccountabilityTest.new(autonity.address, config, {from: deployer});
            await autonity.setAccountabilityContract(accountability.address, {from:operator});

            const tokenUnbondFactor = [1/10, 9/10, 1/100, 99/100, 1/1000, 999/1000, 1/10000000, 9999999/10000000];
            const delegator = accounts[9];
            const balance = (await autonity.balanceOf(delegator)).toNumber();
            let validatorAddresses = [];
            for (let i = 0; i < Math.min(validators.length, tokenUnbondFactor.length); i++) {
                validatorAddresses.push(validators[i].nodeAddress);
            }

            let keyGenCounter = 0;
            while (tokenUnbondFactor.length > validatorAddresses.length) {
                const treasury = accounts[8];
                const nodeKeyInfo = await utils.generateAutonityKeys(`./autonity/data/${keyGenCounter}.key`)
                const oracleKey = utils.randomPrivateKey();
                const popInfo = await utils.generateAutonityPOP(`./autonity/data/${keyGenCounter}.key`, oracleKey, treasury)
                const enode = utils.publicKeyToEnode(nodeKeyInfo.nodePublicKey.substring(2))
                const oracleAddress = utils.address(utils.publicKey(oracleKey, false))
                let consensusKey = Buffer.from(nodeKeyInfo.nodeConsensusKey.substring(2), 'hex');
                let signatures = Buffer.from(popInfo.signatures.substring(2), 'hex')
                await autonity.registerValidator(enode, oracleAddress, consensusKey, signatures, {from: treasury});
                validatorAddresses.push(nodeKeyInfo.nodeAddress);
                keyGenCounter++
            }

            let tokenMinted = []
            for (let iter = 0; iter < validatorAddresses.length; iter++) {
                let validator = validatorAddresses[iter];
                let validatorInfo = await autonity.getValidator(validator);
                let bondedStake = parseInt(validatorInfo.bondedStake);
                // non-self bond to check fairness issue
                const tokenMint = expectedBondedStake - bondedStake;
                tokenMinted.push(tokenMint);
                await autonity.mint(delegator, tokenMint, {from: operator});
                await autonity.bond(validator, tokenMint, {from: delegator});
            }
            // let bonding apply
            await utils.endEpoch(autonity, operator, deployer);

            for (let iter = 0; iter < validatorAddresses.length; iter++) {
                const validator = validatorAddresses[iter];
                let tokenUnBond = Math.max(1, Math.floor(tokenMinted[iter]*tokenUnbondFactor[iter]));
                await autonity.unbond(validator, tokenUnBond, {from: delegator});
            }
            // let unbonding apply and unbondingStake create
            await utils.endEpoch(autonity, operator, deployer);

            for (let iter = 0; iter < validatorAddresses.length; iter++) {
                const validator = validatorAddresses[iter];
                let {txEvent, _} = await utils.slash(config, accountability, 1, validator, validator);
                // checking if highest possible slashing can be done without triggering fairness issue
                // cannot slash (totalStake - 1) because both delegated and unbonding slash is floored
                assert.equal(txEvent.amount.toNumber(), expectedSlash-1, "highest slash did not happen");
                let validatorInfo = await autonity.getValidator(validator);
                assert.equal(validatorInfo.state, utils.ValidatorState.jailed, "validator not jailed");
                assert(parseInt(validatorInfo.bondedStake) > 0 && parseInt(validatorInfo.unbondingStake) > 0, "fairness issue triggered");
            }
            await utils.mineTillUnbondingRelease(autonity, operator, deployer);
            assert.equal((await autonity.balanceOf(delegator)).toNumber(), balance, "unbonding released");
        });
    });
});
