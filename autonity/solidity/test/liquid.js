const truffleAssert = require('truffle-assertions');
const ValidatorLNEW = artifacts.require("Liquid")

const toWei = web3.utils.toWei;
const toBN = web3.utils.toBN;

// Note, tokens are denominted in Wei in this example, but it may not
// be possible to use such small units, depending on how the division
// is implemented in the LiquidNewtonPullFees contract.

contract("Liquid", accounts => {

  // Accounts.
  let rewardSource = accounts[0];
  let treasury = accounts[1];
  let validator = accounts[2];
  let delegatorA = accounts[3];
  let delegatorB = accounts[4];
  let delegatorC = accounts[5];

  // Contract - deployed for each test here
  async function deployLNEW(commissionPercent = 0) {
    // Cannot extract this from the ABI, so have to hard-code it.
    let FEE_FACTOR_UNIT_RECIP = toBN("10000");
    let commission =
      FEE_FACTOR_UNIT_RECIP.mul(toBN(commissionPercent)).div(toBN("100"));
    let lnew = await ValidatorLNEW.new(validator, treasury, commission);
    await lnew.mint(validator, toWei("10000", "ether"));
    return lnew;
  };

  // let lnew;
  // beforeEach(async () => {
  //     let commissionRate = LiquidNewtonPullFees.FEE_FACTOR_UNIT_RECIP.div("2");
  //     lnew = await ValidatorLNEW.new(validator, treasury, commissionRate);
  //     await lnew.mint(validator, toWei("10000", "ether"));
  // });

  async function withdrawAndCheck(lnew, address, expectFees) {
    const origBalance = toBN(await web3.eth.getBalance(address));
    assert.equal(expectFees, await lnew.unclaimedRewards(address));

    // Withdraw
    const txret = await lnew.claimRewards.sendTransaction({from: address});
    const txid = txret.tx;
    const receipt = txret.receipt;
    const tx = await web3.eth.getTransaction(txid);
    const gasCost = toBN(tx.gasPrice).mul(toBN(txret.receipt.gasUsed));
    // expectBalance = origBalance - gasCost + expectFees
    const expectBalance = origBalance.sub(gasCost).add(toBN(expectFees));

    // Balance should have increased by expectFees, and remaining
    // unclaimed fees should be 0
    assert.equal(await lnew.unclaimedRewards(address), "0");
    assert.equal(await web3.eth.getBalance(address), expectBalance);
  };

  it("reward single validator", async () => {
    let lnew = await deployLNEW();

    // Initial state
    assert.equal(await lnew.totalSupply(), toWei("10000", "ether"));
    assert.equal(await lnew.balanceOf(validator), toWei("10000", "ether"));
    assert.equal(await lnew.unclaimedRewards(validator), "0");
    [delegatorA, delegatorB].forEach(async user => {
      assert.equal(await lnew.balanceOf(user), "0");
      assert.equal(await lnew.unclaimedRewards(user), "0");
    });

    // Send 10 AUT as a reward.  Perform a call first (not a tx)
    // in order to check the returned value.
    let distributed = toBN(await lnew.redistribute.call(
      {from: rewardSource, value: toWei("10", "ether")}));
    assert.isTrue(distributed.lte(toBN(toWei("10", "ether"))));
    assert.isTrue(distributed.gt(toBN(toWei("9.9999", "ether"))));
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("10", "ether")});

    // Check distribution (only validator should hold this)
    assert.equal(await lnew.totalSupply(), toWei("10000", "ether"));
    assert.equal(await lnew.balanceOf(validator), toWei("10000", "ether"));
    assert.equal(await lnew.unclaimedRewards(validator), toWei("10", "ether"));
    [delegatorA, delegatorB].forEach(async user => {
      assert.equal(await lnew.balanceOf(user), "0");
      assert.equal(await lnew.unclaimedRewards(user), "0");
    });
  });

  it("reward multiple validators", async () => {
    let lnew = await deployLNEW();

    // delegatorA bonds 8000 NEW
    // delegatorB bonds 2000 NEW
    await lnew.mint(delegatorA, toWei("8000", "ether"));
    await lnew.mint(delegatorB, toWei("2000", "ether"));
    assert.equal(await lnew.totalSupply(), toWei("20000", "ether"));
    assert.equal(await lnew.balanceOf(validator), toWei("10000", "ether"));
    assert.equal(await lnew.balanceOf(delegatorA), toWei("8000", "ether"));
    assert.equal(await lnew.balanceOf(delegatorB), toWei("2000", "ether"));

    // Send 20 AUT as a reward and check distribution
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("20", "ether")});
    assert.equal(await lnew.unclaimedRewards(validator), toWei("10", "ether"));
    assert.equal(await lnew.unclaimedRewards(delegatorA), toWei("8", "ether"));
    assert.equal(await lnew.unclaimedRewards(delegatorB), toWei("2", "ether"));
  });

  it("transfer LNEW", async () => {
    let lnew = await deployLNEW();

    // delegatorA bonds 8000 NEW
    // delegatorB bonds 2000 NEW
    // 20 AUT reward
    await lnew.mint(delegatorA, toWei("8000", "ether"));
    await lnew.mint(delegatorB, toWei("2000", "ether"));
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("20", "ether")});

    // delegatorA gives delegatorC 3000 LNEW
    await lnew.transfer.sendTransaction(
      delegatorC, toWei("3000", "ether"), {from: delegatorA})
    assert.equal(await lnew.totalSupply(), toWei("20000", "ether"));
    assert.equal(await lnew.balanceOf(validator), toWei("10000", "ether"));
    assert.equal(await lnew.balanceOf(delegatorA), toWei("5000", "ether"));
    assert.equal(await lnew.balanceOf(delegatorB), toWei("2000", "ether"));
    assert.equal(await lnew.balanceOf(delegatorC), toWei("3000", "ether"));

    // Another 20 AUT reward.  Check distribution.
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("20", "ether")});
    // validator has 10 + 10
    assert.equal(await lnew.unclaimedRewards(validator), toWei("20", "ether"));
    // delegatorA has 8 + 5
    assert.equal(await lnew.unclaimedRewards(delegatorA), toWei("13", "ether"));
    // delegatorB has 2 + 2
    assert.equal(await lnew.unclaimedRewards(delegatorB), toWei("4", "ether"));
    // delegatorC has 3
    assert.equal(await lnew.unclaimedRewards(delegatorC), toWei("3", "ether"));
  });

  it("burn LNEW", async () => {
    let lnew = await deployLNEW();

    // delegatorA bonds 8000 NEW and burns 3000 LNEW
    await lnew.mint(delegatorA, toWei("8000", "ether"));
    await lnew.burn(delegatorA, toWei("3000", "ether"));
    assert.equal(await lnew.totalSupply(), toWei("15000", "ether"));
    assert.equal(await lnew.balanceOf(validator), toWei("10000", "ether"));
    assert.equal(await lnew.balanceOf(delegatorA), toWei("5000", "ether"));

    // Send 15 AUT as a reward and check distribution
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("15", "ether")});
    assert.equal(await lnew.unclaimedRewards(validator), toWei("10", "ether"));
    assert.equal(await lnew.unclaimedRewards(delegatorA), toWei("5", "ether"));
  });

  it("claiming rewards", async () => {
    let lnew = await deployLNEW();

    // delegatorA bonds 10000 NEW
    await lnew.mint(delegatorA, toWei("10000", "ether"));

    // Send 20 AUT as a reward (validator and delegatorA each
    // earn 10). Withdraw and check balance.
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("20", "ether")});
    await withdrawAndCheck(lnew, delegatorA, toWei("10", "ether"));

    // Send 40 AUT as a reward (validator and delegatorA each
    // earn 20). Withdraw and check balance.
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("40", "ether")});
    await withdrawAndCheck(lnew, delegatorA, toWei("20", "ether"));
  });

  it("accumulating rewards", async () => {
    let lnew = await deployLNEW();

    // delegatorA bonds 10000 NEW (total 20000 delegated)
    await lnew.mint(delegatorA, toWei("10000", "ether"));

    // Send 20 AUT as a reward (delegatorA earns 10)
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("20", "ether")});

    // Other delegators bond 20000 NEW (total of 40000 NEW bonded)
    await lnew.mint(delegatorB, toWei("12000", "ether"));
    await lnew.mint(delegatorC, toWei("8000", "ether"));

    // Send 20 AUT as a reward (delegatorA earns 5)
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("20", "ether")});

    // Other delegators bond 10000 NEW (total of 50000 NEW bonded)
    await lnew.mint(validator, toWei("2000", "ether"));
    await lnew.mint(delegatorC, toWei("8000", "ether"));

    // Send 50 AUT as a reward (delegatorA earns 10)
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("50", "ether")});

    // Check delegatorA's total fees were 10 + 5 + 10 = 25
    assert.equal(
      await lnew.unclaimedRewards(delegatorA),
      toWei("25", "ether"));
  });

  it("commission", async () => {
    // use 50% commission for simplcity
    const lnew = await deployLNEW(50);
    const treasuryBalance = toBN(await web3.eth.getBalance(treasury));

    // delegatorA bonds 10000 NEW (total 20000 delegated)
    await lnew.mint(delegatorA, toWei("10000", "ether"));

    // Send 40 AUT as a reward (treasury earns 20, delegatorA earns 10)
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("40", "ether")});

    // Other delegators bond 20000 NEW (total of 40000 NEW bonded)
    await lnew.mint(delegatorB, toWei("12000", "ether"));
    await lnew.mint(delegatorC, toWei("8000", "ether"));

    // Send 40 AUT as a reward (treasury earns 20 delegatorA earns 5)
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("40", "ether")});

    // Other delegators bond 10000 NEW (total of 50000 NEW bonded)
    await lnew.mint(validator, toWei("2000", "ether"));
    await lnew.mint(delegatorC, toWei("8000", "ether"));

    // Send 100 AUT as a reward (treasury earns 50, delegatorA earns 10)
    await lnew.redistribute.sendTransaction(
      {from: rewardSource, value: toWei("100", "ether")});

    // Check treasury balance increased by: 20 + 20 + 50 = 90
    assert.equal(
      toBN(await web3.eth.getBalance(treasury)).sub(treasuryBalance),
      toWei("90", "ether"));

    // Check delegatorA's total fees: 10 + 5 + 10 = 25
    assert.equal(
      await lnew.unclaimedRewards(delegatorA),
      toWei("25", "ether"));
  });

  it("allowances", async () => {
    const lnew = await deployLNEW();

    // delegatorA bonds 10000 NEW
    await lnew.mint(delegatorA, toWei("10000", "ether"));

    // delegatorC should not be able to transfer on A's behalf
    assert.equal(await lnew.allowance(delegatorA, delegatorC), "0");
    await truffleAssert.fails(lnew.transferFrom.sendTransaction(
      delegatorA, delegatorB, toWei("1000", "ether"), {from: delegatorC}));

    // A grants C permission to spend 5000.
    await lnew.approve.sendTransaction(
      delegatorC, toWei("5000", "ether"), {from: delegatorA});
    assert.equal(
      await lnew.allowance(delegatorA, delegatorC),
      toWei("5000", "ether"));

    // C sends 1000 of A's LNEW to B
    await lnew.transferFrom.sendTransaction(
      delegatorA, delegatorB, toWei("1000", "ether"), {from: delegatorC});

    // Check balances and allowances
    assert.equal(await lnew.balanceOf(delegatorA), toWei("9000", "ether"));
    assert.equal(await lnew.balanceOf(delegatorB), toWei("1000", "ether"));
    assert.equal(await lnew.balanceOf(delegatorC), "0");
    assert.equal(
      await lnew.allowance(delegatorA, delegatorC),
      toWei("4000", "ether"));

    // Sending 4001 should fail.
    await truffleAssert.fails(lnew.transferFrom.sendTransaction(
      delegatorA, delegatorB, toWei("4001", "ether"), {from: delegatorC}));
  });
});
