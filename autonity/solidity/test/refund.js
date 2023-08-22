const Oracle = artifacts.require("Oracle")
const truffleAssert = require('truffle-assertions');
const assert = require('assert')
const utils = require('./utils.js');
const toBN = web3.utils.toBN;
const BN = require('bn.js');

// this vote refund tests cannot be run on ganache, since it does not have the refund mechanism

contract("Oracle", accounts => {
  let oracle;

  describe('Oracle vote refund', function() {
    before(async() => {
      // for testing the refund we need to interact with the oracle contract deployed at genesis.
      // the refund logic checks if the vote is sent to this specific oracle contract.
      oracle = await Oracle.at("0x47e9Fbef8C83A1714F1951F142132E6e90F5fa5D")
    })
    afterEach(async() => {
      // after each test we wait for a round change, to have a clean contract state again
      let round =  await oracle.getRound();
      let curRound = +round;
      while (+curRound == +round) {
        await utils.timeout(1000);
        curRound = await oracle.getRound();
      }
    })

    //TODO(lorenzo) these tests can probably be merged/simplified, as they came from when we were refunding only the basefee
    it('fee is refunded for valid vote (no tip)', async function () {
      const origBalance = toBN(await web3.eth.getBalance(accounts[2]));

      txret = await oracle.vote(0, [], 0, {from:accounts[2], maxPriorityFeePerGas:0});
      
      // fetch base fee that got applied to the vote tx
      //txBlock = await web3.eth.getBlock(txret.receipt.blockNumber)
      //baseFee = toBN(txBlock.baseFeePerGas)
     
      // compute total gasCost, refund (cost from basefee) + tip
      const tx = await web3.eth.getTransaction(txret.tx);
      const gasCost = toBN(tx.gasPrice).mul(toBN(txret.receipt.gasUsed));
      //const effectiveTip = BN.min(toBN(tx.maxPriorityFeePerGas),toBN(tx.maxFeePerGas).sub(baseFee))
      //const tip = effectiveTip.mul(toBN(txret.receipt.gasUsed));
      //const refund = toBN(baseFee).mul(toBN(txret.receipt.gasUsed));

      // gasCost = refund + tip
      //assert.equal(gasCost.toString(),refund.add(tip).toString())

      //const expectedBalance = origBalance.sub(gasCost).add(refund);
      const updatedBalance = toBN(await web3.eth.getBalance(accounts[2]))

      assert.equal(updatedBalance.toString(), origBalance.toString());
    });
    it('fee is refunded for valid vote (with tip)', async function () {
      const origBalance = toBN(await web3.eth.getBalance(accounts[2]));

      txret = await oracle.vote(0, [], 0, {from:accounts[2]});
      
      // fetch base fee that got applied to the vote tx
      //txBlock = await web3.eth.getBlock(txret.receipt.blockNumber)
      //baseFee = toBN(txBlock.baseFeePerGas)
     
      // compute total gasCost, refund (cost from basefee) + tip
      const tx = await web3.eth.getTransaction(txret.tx);
      const gasCost = toBN(tx.gasPrice).mul(toBN(txret.receipt.gasUsed));
      //const effectiveTip = BN.min(toBN(tx.maxPriorityFeePerGas),toBN(tx.maxFeePerGas).sub(baseFee))
      //const tip = effectiveTip.mul(toBN(txret.receipt.gasUsed));
      //const refund = toBN(baseFee).mul(toBN(txret.receipt.gasUsed));

      // gasCost = refund + tip
      //assert.equal(gasCost.toString(),refund.add(tip).toString())

      //expectedBalance = origBalance - gasCost + refund (we pay only the tip)
      //const expectedBalance = origBalance.sub(gasCost).add(refund);
      const updatedBalance = toBN(await web3.eth.getBalance(accounts[2]))

      assert.equal(updatedBalance.toString(), origBalance.toString());
    });
    it('double vote, only first is refunded', async function () {
      let origBalance = toBN(await web3.eth.getBalance(accounts[2]));

      let round = await oracle.getRound()

      let txret = await oracle.vote(0, [], 0, {from:accounts[2]});
      
      // fetch base fee that got applied to the vote tx
      //let txBlock = await web3.eth.getBlock(txret.receipt.blockNumber)
      //let baseFee = toBN(txBlock.baseFeePerGas)
     
      // compute total gasCost, refund (cost from basefee) + tip
      let tx = await web3.eth.getTransaction(txret.tx);
      let gasCost = toBN(tx.gasPrice).mul(toBN(txret.receipt.gasUsed));
      //let effectiveTip = BN.min(toBN(tx.maxPriorityFeePerGas),toBN(tx.maxFeePerGas).sub(baseFee))
      //let tip = effectiveTip.mul(toBN(txret.receipt.gasUsed));
      //let refund = toBN(baseFee).mul(toBN(txret.receipt.gasUsed));

      // gasCost = refund + tip
      //assert.equal(gasCost.toString(),refund.add(tip).toString())

      //expectedBalance = origBalance - gasCost + refund (we pay only the tip)
      //let expectedBalance = origBalance.sub(gasCost).add(refund);
      let updatedBalance = toBN(await web3.eth.getBalance(accounts[2]))

      assert.equal(updatedBalance.toString(), origBalance.toString());

      // make sure we are still in the same round
      round2 = await oracle.getRound()
      assert.equal(round.toString(),round2.toString())
      
      // second vote with 0 tip, base fee should not be refunded
      await truffleAssert.fails(
        oracle.vote(0, [], 0, {from:accounts[2],maxPriorityFeePerGas:0}),
        truffleAssert.ErrorType.REVERT,
        "already voted"
      );
      
      updatedBalance2 = toBN(await web3.eth.getBalance(accounts[2]))

      // some gas should have been spent
      assert.notEqual(updatedBalance.toString(), updatedBalance2.toString());
    });

  });
});
