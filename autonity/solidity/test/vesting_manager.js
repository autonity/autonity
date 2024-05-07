'use strict'
const assert = require('assert')
const truffleAssert = require('truffle-assertions')
const utils = require('./utils.js')
const Autonity = artifacts.require("Autonity")
const AccountabilityTest = artifacts.require("AccountabilityTest")
const VestingManager = artifacts.require("VestingManager")
const VestingManagerTest = artifacts.require("VestingManagerTest")
const TestContract = artifacts.require("TestContract")
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

    xdescribe('set hardcoded value', async function () {
        beforeEach(async function () {
            autonity = await utils.deployAutonityTestContract(validators, autonityConfig, accountabilityConfig, deployer, operator);
            vestingManager = await VestingManagerTest.new(autonity.address, operator, {from: deployer})
        })

        it('measure gas for finalize()', async function () {
            let amount = 100000000000
            let bond = amount / 4
            console.log("starts")
            let tx = toBN(await vestingManager.unclaimedRewards.call(anyAccount)).toString()
            console.log(tx)
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
            console.log(await web3.eth.getBalance(autonity.address))
            let committee = await autonity.getCommittee()
            console.log(committee)
            console.log(validator)
            // assert.notEqual(toBN(await vestingManager.unclaimedRewards.call(anyAccount)), 0, "no reward")
            // bond again and apply that bonding
            // let bondingID = (await autonity.getHeadBondingID()).toNumber()
            await vestingManager.bond(0, validator, bond, {from: anyAccount})
            await vestingManager.bond(0, validator, bond, {from: anyAccount})
            await vestingManager.bond(0, validator, bond, {from: anyAccount})
            // await vestingManager.cancelSchedule(anyAccount, 0, operator, {from: operator})
            let appliedBonding = []
            appliedBonding.push({validator: validator, liquidAmount: bond*3})
            tx = await vestingManager.callFinalize(appliedBonding, [], [])
            console.log(tx.receipt.gasUsed)
            // tx = await vestingManager.applyBonding(bondingID, bond, false, false)
            // bondingID++
            // console.log(tx.receipt.gasUsed)
            // tx = await vestingManager.applyBonding(bondingID, bond, false, false)
            // bondingID++
            // console.log(tx.receipt.gasUsed)
            // tx = await vestingManager.applyBonding(bondingID, bond, false, false)
            // console.log(tx.receipt.gasUsed)
            // amount = 0
            assert.equal((await vestingManager.liquidBalanceOf(anyAccount, 0, validator)).toNumber(), amount, "bonding not applied")
            console.log(toBN(await vestingManager.unclaimedRewards.call(anyAccount)).toString())
            // tx = await vestingManager.testRequire.sendTransaction()
            // console.log(tx.receipt.gasUsed)
            // tx = await vestingManager.testRequireValue.sendTransaction(1)
            // console.log(tx.receipt.gasUsed)
        })

        xit('measure gas for _applyBonding(_rejected = true)', async function () {

        })

        xit('measure gas for _applyUnbonding', async function () {
            
        })

        xit('measure gas for _releaseUnbonding()', async function () {
            
        })

        // TODO: check for both case of schedule canceled or not
    })

    describe('test', async function () {
        let testContract
        beforeEach(async function () {
            autonity = await utils.deployAutonityTestContract(validators, autonityConfig, accountabilityConfig, deployer, operator);
            testContract = await TestContract.new(autonity.address, {from: deployer})
        })

        it('test epoch', async function () {
            let blockNumber = await web3.eth.getBlockNumber()
            console.log(blockNumber)
            console.log((await autonity.blockNumber()).toNumber())
            console.log((await autonity.blockNumber.call()).toNumber())
            // let block = await web3.eth.getBlock(blockNumber)
            // console.log(block)
            let epochID = (await autonity.epochID()).toNumber()
            console.log(epochID)
            let epochPeriod = (await autonity.getEpochPeriod()).toNumber()
            console.log(epochPeriod)
            console.log((await autonity.lastEpochBlock()).toNumber())
            let lastBlock = blockNumber
            while (true) {
                console.log("loop\n")
                await autonity.finalize({from: deployer})
                let newBlockNumber = await web3.eth.getBlockNumber()
                console.log(newBlockNumber)
                console.log((await autonity.blockNumber()).toNumber())
                console.log((await autonity.blockNumber.call()).toNumber())
                let newEpochID = (await autonity.epochID()).toNumber()
                console.log(newEpochID)
                console.log((await autonity.getEpochPeriod()).toNumber())
                console.log((await autonity.lastEpochBlock()).toNumber())
                assert.equal(newBlockNumber, lastBlock+1)
                lastBlock = newBlockNumber
                if (newEpochID == epochID+1) {
                    break
                }
            }
        })

        xit('test gas', async function () {
            await testContract.setValue(1, 1)

            let tx = await testContract.getValue.sendTransaction(1)
            console.log(tx.receipt.gasUsed)

            tx = await testContract.getSameValue.sendTransaction(1)
            console.log(tx.receipt.gasUsed)

            tx = await testContract.getSameValueAgain.sendTransaction(1)
            console.log(tx.receipt.gasUsed)

            tx = await testContract.getValueAndUpdate.sendTransaction(1)
            console.log(tx.receipt.gasUsed)

            await testContract.setValue(1, 1)

            let a = (await autonity.testGasGetValue.call(1, testContract.address)).toNumber()
            console.log(a)

            a = (await autonity.testGasGetSameValue.call(1, testContract.address)).toNumber()
            console.log(a)

            a = (await autonity.testGasGetSameValueAgain.call(1, testContract.address)).toNumber()
            console.log(a)

            a = (await autonity.testGasGetValueAndUpdate.call(1, testContract.address)).toNumber()
            console.log(a)
        })

        xit('test callback gas', async function () {
            let a = (await autonity.doCallback.call(testContract.address)).toNumber()
            // let a = tx
            console.log(a)
            console.log((await autonity.doCallback(testContract.address)).toNumber())
            let b = (await autonity.doCallbackAgain.call(testContract.address, 1, 2)).toNumber()
            // let b = tx
            console.log(b)
            console.log((await autonity.doCallbackAgain(testContract.address, 1, 2)).toNumber())
            let c = (await autonity.doCallbackMath.call(testContract.address, 1, 2)).toNumber()
            // let c = tx
            console.log(c)
            console.log((await autonity.doCallbackMath(testContract.address, 1, 2)).toNumber())
            let tx = await testContract.nothing.sendTransaction()
            let d = tx.receipt.gasUsed
            console.log(d)
            console.log(a - d + 20000)
            tx = await testContract.nothing.sendTransaction(1, 2)
            let e = tx.receipt.gasUsed
            console.log(e)
            console.log(b - e + 20000)
            tx = await testContract.add.sendTransaction(1, 2)
            let f = tx.receipt.gasUsed
            console.log(f)
            console.log(c - f + 20000)
            let g = (await autonity.doAll.call(testContract.address, 1, 2)).toNumber()
            // let g = tx
            console.log(g)
            console.log((await autonity.doAll(testContract.address, 1, 2)).toNumber())
            console.log(a+b+c)
            console.log(a+b+c - g)
            let h = (await autonity.doAllFast.call(testContract.address, 1, 2)).toNumber()
            // let h = tx
            console.log(h)
            console.log((await autonity.doAllFast(testContract.address, 1, 2)).toNumber())
        })

        xit('test gas', async function () {
            let bond = 100
            await autonity.mint(testContract.address, bond, {from: operator})
            await autonity.mint(operator, bond, {from: operator})
            let validator = validators[0].nodeAddress
            let tx = await autonity.bond(validator, bond, {from: operator})
            console.log(tx.receipt.gasUsed)
            let {0: a, 1: b, 2: c, 3: d} = await testContract.bond.call(validator, bond, bond)
            a = a.toNumber()
            b = b.toNumber()
            c = c.toNumber()
            console.log(a)
            console.log(`${b} : ${a-b}`)
            console.log(`${c} : ${b-c}`)
            console.log(d)
        })

        xit('test persistence', async function () {
            let bond = 100
            await autonity.mint(testContract.address, bond, {from: operator})
            let validator = validators[0].nodeAddress
            let info = await autonity.getValidator(validator)
            console.log(info.liquidSupply)
            let bondingID = (await autonity.getHeadBondingID()).toNumber()
            await testContract.bond(validator, bond, info.liquidSupply)
            await utils.endEpoch(autonity, operator, deployer)
            assert.equal(await testContract.isApplied(bondingID), true, "bonding not applied")
        })

        xit('compare', async function () {
            await testContract.testCompare();
        })

        xit('storage', async function () {
            let id = 4
            let oldAmount = 10
            let newAmount = 133
            await testContract.createBondingRequest(id, oldAmount, anyAccount)
            await testContract.testStorage(id, oldAmount, newAmount)
        })
    })

})