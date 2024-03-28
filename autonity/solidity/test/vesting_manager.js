'use strict'
const assert = require('assert')
const truffleAssert = require('truffle-assertions')
const utils = require('./utils.js')
const Autonity = artifacts.require("Autonity")
const AccountabilityTest = artifacts.require("AccountabilityTest")
const VestingManager = artifacts.require("VestingManager")
const VestingManagerTest = artifacts.require("VestingManagerTest")
const toBN = web3.utils.toBN
const config = require("./config")

function checkEventFails(tx, eventName, name) {
    truffleAssert.eventEmitted(
        tx, `${eventName}`, (ev) => {
            return true
        },
        `should emit correct event ${name}`
    )
}

contract('VestingManager', function (accounts) {
    before(async function () {
        await utils.mockPrecompile()
    });
    for (let i = 0; i < accounts.length; i++) {
        console.log("account: ", i, accounts[i])
    }
  
    const operator = accounts[5]
    const deployer = accounts[6]
    const anyAccount = accounts[7]
    const treasuryAccount = accounts[8]
    const zeroAddress = "0x0000000000000000000000000000000000000000"
    const accountabilityConfig = config.ACCOUNTABILITY_CONFIG
    const autonityConfig = config.autonityConfig(operator, treasuryAccount)
    const genesisPrivateKeys = config.GENESIS_PRIVATE_KEYS
  
    // accounts[2] is skipped because it is used as a genesis validator when running against autonity
    // this can cause interference in reward distribution tests
    let validators = config.validators(accounts)
  
    let autonity
    let accountability
    let vestingManager

    xdescribe('Reward flow', async function () {
        beforeEach(async function () {
            autonity = await Autonity.new(validators, autonityConfig, {from: deployer})
            await autonity.finalizeInitialization({from: deployer})
            vestingManager = await VestingManager.new(autonity.address, operator, {from: deployer})
        })
        // add tests
    })

    xdescribe('Vesting flow', async function () {
        beforeEach(async function () {
            autonity = await Autonity.new(validators, autonityConfig, {from: deployer})
            await autonity.finalizeInitialization({from: deployer})
            accountability = await AccountabilityTest.new(autonity.address, accountabilityConfig, {from: deployer})
            await autonity.setAccountabilityContract(accountability.address, {from:operator})
            vestingManager = await VestingManager.new(autonity.address, operator, {from: deployer})
        })
        // add tests
    })

    xdescribe('Bonding and Unbonding flow', async function () {
        beforeEach(async function () {
            autonity = await utils.deployAutonityTestContract(validators, autonityConfig, accountabilityConfig, deployer, operator);
            vestingManager = await VestingManager.new(autonity.address, operator, {from: deployer})
        })

        it('can bond before release', async function () {
            let amount = 100
            let validator = validators[0].nodeAddress
            await autonity.mint(operator, amount, {from: operator})
            await autonity.approve(vestingManager.address, amount, {from: operator})
            await vestingManager.newSchedule(anyAccount, amount, 0, 0, 100000000000, true, {from: operator})
            assert.notEqual((await vestingManager.releasedNTN(anyAccount, 0)).toNumber(), amount, "not released")
            let id = (await autonity.getHeadBondingID()).toNumber()
            await vestingManager.bond(0, validator, amount, {from: anyAccount})
            let request = await autonity.getBondingRequest(id)
            assert.equal(request.amount, amount)
            await utils.endEpoch(autonity, operator, deployer)
            let schedule = await vestingManager.getSchedule(anyAccount, 0);
            assert.equal(schedule.totalAmount, 0, "not bonded")
            assert.equal((await vestingManager.liquidBalanceOf(anyAccount, 0, validator)).toNumber(), amount, "no LNTN")
        })

        // add tests
    })

    xdescribe('Access restriction', async function () {
        beforeEach(async function () {
            autonity = await Autonity.new(validators, autonityConfig, {from: deployer})
            await autonity.finalizeInitialization({from: deployer})
            vestingManager = await VestingManager.new(autonity.address, operator, {from: deployer})
        })

        it('cannot apply bonding operations', async function () {
            await truffleAssert.fails(
                vestingManager.bondingApplied(0, 0, false, false),
                truffleAssert.ErrorType.REVERT,
                "function restricted to Autonity contract"
            )

            await truffleAssert.fails(
                vestingManager.unbondingApplied(0),
                truffleAssert.ErrorType.REVERT,
                "function restricted to Autonity contract"
            )

            await truffleAssert.fails(
                vestingManager.unbondingReleased(0, 0),
                truffleAssert.ErrorType.REVERT,
                "function restricted to Autonity contract"
            )
        })

        it('cannot create or cancel schedule', async function () {
            await truffleAssert.fails(
                vestingManager.newSchedule(zeroAddress, 0, 0, 0, 0, false),
                truffleAssert.ErrorType.REVERT,
                "caller is not the operator"
            )

            await truffleAssert.fails(
                vestingManager.cancelSchedule(zeroAddress, 0, zeroAddress),
                truffleAssert.ErrorType.REVERT,
                "caller is not the operator"
            )
        })
    })

    describe('set hardcoded value', async function () {
        beforeEach(async function () {
            autonity = await utils.deployAutonityTestContract(validators, autonityConfig, accountabilityConfig, deployer, operator);
            vestingManager = await VestingManagerTest.new(autonity.address, operator, {from: deployer})
        })

        it('measure gas for _applyBonding', async function () {
            let amount = 100000000000
            let bond = amount / 2
            let validator = validators[0].nodeAddress
            await autonity.mint(operator, amount, {from: operator})
            await autonity.approve(vestingManager.address, amount, {from: operator})
            await vestingManager.newSchedule(anyAccount, amount, 0, 0, 100000000000, true, {from: operator})
            await vestingManager.bond(0, validator, bond, {from: anyAccount})
            await utils.endEpoch(autonity, operator, deployer)
            console.log((await vestingManager.liquidBalanceOf(anyAccount, 0, validator)).toNumber())
            assert.equal((await vestingManager.liquidBalanceOf(anyAccount, 0, validator)).toNumber(), bond, "no LNTN")
            let info = await autonity.getValidator(validator)
            console.log(info.liquidSupply)
            // vesting manager should get some reward at epoch end
            let initFunds = await web3.eth.getBalance(autonity.address)
            let reward = web3.utils.toWei("1", "ether")
            console.log(await web3.eth.getBalance(autonity.address))
            console.log(await web3.eth.getBalance(anyAccount))
            await web3.eth.sendTransaction({from: anyAccount, to: autonity.address, value: reward})
            console.log(await web3.eth.getBalance(autonity.address))
            console.log(await web3.eth.getBalance(anyAccount))
            // let newBalance =  toBN(reward).add(toBN(initFunds))
            // console.log(newBalance)
            // assert.equal(toBN(await web3.eth.getBalance(autonity.address)), newBalance, "autonity balance not changed")
            await utils.endEpoch(autonity, operator, deployer)
            let committee = await autonity.getCommittee()
            console.log(committee)
            console.log(validator)
            console.log((await vestingManager.unclaimedRewards(anyAccount, 0)).toNumber())
            assert.notEqual((await vestingManager.unclaimedRewards(anyAccount, 0)).toNumber(), 0, "no reward")
            // bond again and apply that bonding
            let bondingID = (await autonity.getHeadBondingID()).toNumber()
            await vestingManager.bond(0, validator, bond, {from: anyAccount})
            let tx = await vestingManager.applyBonding(bondingID, bond, false, false)
            assert.equal((await vestingManager.liquidBalanceOf(anyAccount, 0, validator)).toNumber(), bond+bond, "bonding not applied")
            console.log(tx.receipt.gasUsed)
        })
    })
})