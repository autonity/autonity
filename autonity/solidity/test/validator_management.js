'use strict';
const assert = require('assert');
const { Buffer } = require('node:buffer');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const liquidContract = artifacts.require("Liquid")

function generateMultiSig(nodekey,oraclekey,treasuryAddr) {
    let treasuryProof = web3.eth.accounts.sign(treasuryAddr, nodekey);
    let oracleProof = web3.eth.accounts.sign(treasuryAddr, oraclekey);
    let multisig = treasuryProof.signature + oracleProof.signature.substring(2)
    return multisig
}

contract('Autonity', function (accounts) {

    for (let i = 0; i < accounts.length; i++) {
        console.log("account: ", i, accounts[i]);
    }

    const operator = accounts[5];
    const deployer = accounts[6];
    const anyAccount = accounts[7];
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
        "validatorKey": "0x00",
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

    let autonity;

    describe('Validator management', function () {
        beforeEach(async function () {
            autonity = await utils.deployContracts(validators, autonityConfig, accountabilityConfig, deployer, operator);
        });

        it('Add validator with already registered address', async function () {
            let newValidator = accounts[0];
            let enode = genesisEnodes[0]
            // multisig length is checked before validator already registered (it is not verified though)
            let multisig = generateMultiSig(genesisPrivateKeys[0],genesisPrivateKeys[0],newValidator)
            let validatorKey = Buffer.from('845681310fe66ed10629e76cc5aa20f3ec8b853af9f3dee8a6318f3fb81c0adcaaa0a776dc066127e743bba6b0349bc0', 'hex');
            let validatorKeyProof = '0x88a19caac1d02d2efb3675ec9fe99936b1170641b03d7525674ee001446cfd204fa5ba0b5e362d71294f3ba2f758695115a17101fc70b73fe90d7eb83950c3f7ad598b6740698b8e78fb48821c47762cdf2de889deede80fe2e7c085e48562c4';
            multisig = multisig + validatorKeyProof.substring(2);

            await truffleAssert.fails(
                autonity.registerValidator(enode, genesisNodeAddresses[0], validatorKey, multisig, {from: newValidator}),
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
            let validatorKey = Buffer.from('b4c9a6216f9e39139b8ea2b36f277042bbf5e1198d8e01cff0cca816ce5cc820e219025d2fa399b133d3fc83920eeca5', 'hex');
            let validatorKeyProof = '0xa141b3c759ad5eec4def611fc4cb028f1edb0f363f9f415c692998b0b6e677acdfb7e2ac23e3e848027b5e19e56b550c15a87ccc81e6f8ebd34fa54850ec0fe192567bf4aefcddb06f6c00bee4768010013b162a91d4f7ed397568affe497532';

            let multisig = generateMultiSig(privateKey,privateKey,newValidator)
            multisig = multisig + validatorKeyProof.substring(2)

            await truffleAssert.fails(
                autonity.registerValidator(enode, newValidator, validatorKey, multisig, {from: newValidator}),
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
            let validatorKey = Buffer.from('b4c9a6216f9e39139b8ea2b36f277042bbf5e1198d8e01cff0cca816ce5cc820e219025d2fa399b133d3fc83920eeca5', 'hex');
            let validatorKeyProof = '0xa141b3c759ad5eec4def611fc4cb028f1edb0f363f9f415c692998b0b6e677acdfb7e2ac23e3e848027b5e19e56b550c15a87ccc81e6f8ebd34fa54850ec0fe192567bf4aefcddb06f6c00bee4768010013b162a91d4f7ed397568affe497532';
            multisig = multisig + validatorKeyProof.substring(2)

            await truffleAssert.fails(
                autonity.registerValidator(enode, oracleAddr, validatorKey, multisig, {from: newValidator}),
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

            // generate the validator Key and multisigs from console:
            //./autonity genOwnershipProof --nodekeyhex e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b --oraclekeyhex e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b 0xe12b43B69E57eD6ACdd8721Eb092BF7c8D41Df41
            let validatorKey = Buffer.from("b4c9a6216f9e39139b8ea2b36f277042bbf5e1198d8e01cff0cca816ce5cc820e219025d2fa399b133d3fc83920eeca5", 'hex');
            let multisig = Buffer.from("d4b63f6b5535d7255dbb5ecc5092c7eb042de1d20dff80535321dc1f8fa3cf8844a2927ad86d4e74573b5af4bb69a2a788d0e98a0d2410aed51d355985836cb701d4b63f6b5535d7255dbb5ecc5092c7eb042de1d20dff80535321dc1f8fa3cf8844a2927ad86d4e74573b5af4bb69a2a788d0e98a0d2410aed51d355985836cb70191c4492d13544d3ea23aab9b051796e11285f519dc2316cac3d96c5f3d594459474438b09f6e60a25ea22938ed6379760b573466601576a1967cb5aceabe12c4aa2e27f67666f1a3af5fbc4b7209cb83f7e76a4be4c03e1dc99d662f9ea883ec", "hex");
            let oracleAddr = newValAddr

            await autonity.registerValidator(enode, oracleAddr, validatorKey, multisig, {from: issuerAccount});
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
            let validatorKey = Buffer.from("b4c9a6216f9e39139b8ea2b36f277042bbf5e1198d8e01cff0cca816ce5cc820e219025d2fa399b133d3fc83920eeca5", "hex");
            let multisigs = Buffer.from("d4b63f6b5535d7255dbb5ecc5092c7eb042de1d20dff80535321dc1f8fa3cf8844a2927ad86d4e74573b5af4bb69a2a788d0e98a0d2410aed51d355985836cb701d4b63f6b5535d7255dbb5ecc5092c7eb042de1d20dff80535321dc1f8fa3cf8844a2927ad86d4e74573b5af4bb69a2a788d0e98a0d2410aed51d355985836cb70191c4492d13544d3ea23aab9b051796e11285f519dc2316cac3d96c5f3d594459474438b09f6e60a25ea22938ed6379760b573466601576a1967cb5aceabe12c4aa2e27f67666f1a3af5fbc4b7209cb83f7e76a4be4c03e1dc99d662f9ea883ec", "hex");
            let oracleAddr = validator

            // disabling a non registered validator should fail
            await truffleAssert.fails(
                autonity.pauseValidator(validator, {from: issuerAccount}),
                truffleAssert.ErrorType.REVERT,
                "validator must be registered"
            );

            await autonity.registerValidator(enode, oracleAddr, validatorKey, multisigs, {from: issuerAccount});

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
            // activating a non-existing validator should fail
            await truffleAssert.fails(
                autonity.activateValidator(validator, {from: issuerAccount}),
                truffleAssert.ErrorType.REVERT,
                "validator must be registered"
            );

            let validatorKey = Buffer.from("b4c9a6216f9e39139b8ea2b36f277042bbf5e1198d8e01cff0cca816ce5cc820e219025d2fa399b133d3fc83920eeca5", "hex");
            let multisigs = Buffer.from("d4b63f6b5535d7255dbb5ecc5092c7eb042de1d20dff80535321dc1f8fa3cf8844a2927ad86d4e74573b5af4bb69a2a788d0e98a0d2410aed51d355985836cb701d4b63f6b5535d7255dbb5ecc5092c7eb042de1d20dff80535321dc1f8fa3cf8844a2927ad86d4e74573b5af4bb69a2a788d0e98a0d2410aed51d355985836cb70191c4492d13544d3ea23aab9b051796e11285f519dc2316cac3d96c5f3d594459474438b09f6e60a25ea22938ed6379760b573466601576a1967cb5aceabe12c4aa2e27f67666f1a3af5fbc4b7209cb83f7e76a4be4c03e1dc99d662f9ea883ec", "hex");
            let oracleAddr = validator
            await autonity.registerValidator(enode, oracleAddr, validatorKey, multisigs, {from: issuerAccount});

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

    describe('Test committee members rotation through bonding/unbonding', function () {
        let validatorKey1 = Buffer.from("a39f5fd136836a203bfd13d8acc631199c478d9aaa67b147989bdc75676c9e084c0e3396011ff370ca4635723c335a03", "hex");
        let multisig1 = Buffer.from("b958d8998c700728340e78f5371eda293602de4e0dccde8184ddb65c87c5b21b7bf4374c8df5b32cf8b611746e21403ecc1ab4182baba1a67962d4d84b95350101b958d8998c700728340e78f5371eda293602de4e0dccde8184ddb65c87c5b21b7bf4374c8df5b32cf8b611746e21403ecc1ab4182baba1a67962d4d84b953501018e33e67e311ff80635094f1eaddca209064013850fa521dda481bd6f96702f491de7ae38537278fc807da025041fdc08023188392f68c232ce2cf0d9d03a106e572a11696356d6bdbf70ff4058040d735eebb5527b7e2c2137aed3d5532ae684", "hex");
        let validatorKey2 = Buffer.from("9271d72f26539bbb1beb011b63fa63c56a7c225e9e20933cc8a501204fdf8b302922e11e9d45015d6547dd4e117b9c5e", "hex");
        let multisig2 = Buffer.from("2bcd02051836b04282d158c70e00236ec868019563b9caa7a6e1fc35fbc648ea5526dbdca54b0f3b5448325462b202a792582ff37ce04cf1f0c166e271dfc339012bcd02051836b04282d158c70e00236ec868019563b9caa7a6e1fc35fbc648ea5526dbdca54b0f3b5448325462b202a792582ff37ce04cf1f0c166e271dfc33901817c5c5ce485cf0a26c46d4931f3a40dfe9da29f5da11710cffb1f26d8db692a9034eb59a797a8b3111d11f182ce03ad0e02820dfcdb4f3c7740b083509840003f347f0f476a6d51d5208a5e23781822c84c540089e86450c15342d062232d9e", "hex");
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
            let validator1 = accounts[1]; // treasury
            let oracle1 = genesisNodeAddresses[1] // oracle address = node address
            let enodeVal1 = genesisEnodes[1]

            let validator2 = accounts[3];
            let oracle2 = genesisNodeAddresses[2]
            let enodeVal2 = genesisEnodes[2]

            await token.registerValidator(enodeVal1, oracle1, validatorKey1, multisig1, {from: validator1});
            await token.registerValidator(enodeVal2, oracle2, validatorKey2, multisig2, {from: validator2});

            // system operator mint Newton for user.
            let user = accounts[7];
            let tokenMint = 100;
            await token.mint(user, tokenMint, {from: operator});

            // user bond Newton to validator 2.
            await token.bond(genesisNodeAddresses[2], tokenMint, {from: user});

            // close epoch to ensure bonding is applied
            await utils.endEpoch(token,operator,deployer);

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
            let validator1 = accounts[1]; // treasury
            let oracle1 = genesisNodeAddresses[1] // oracle address = node address
            let enodeVal1 = genesisEnodes[1]

            let validator2 = accounts[3];
            let oracle2 = genesisNodeAddresses[2]
            let enodeVal2 = genesisEnodes[2]

            await token.registerValidator(enodeVal1, oracle1, validatorKey1, multisig1, {from: validator1});
            await token.registerValidator(enodeVal2, oracle2, validatorKey2, multisig2, {from: validator2});

            // system operator mint Newton for user.
            let user = accounts[7];
            let tokenMint = 100;
            await token.mint(user, tokenMint, {from: operator});

            // bond NTN to the 2 validators
            await token.bond(genesisNodeAddresses[1], 20, {from: user});
            await token.bond(genesisNodeAddresses[2], 20, {from: user});

            // close epoch to ensure bonding is applied
            await utils.endEpoch(token,operator,deployer);

            let committee = await token.getCommittee();
            assert.equal(committee.length,3)

            await token.unbond(genesisNodeAddresses[1], 20, {from: user});

            // voting power gets reduced right away at epoch end
            await utils.endEpoch(token,operator,deployer);
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
});
