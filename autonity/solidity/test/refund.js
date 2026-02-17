const Oracle = artifacts.require("Oracle")
const Autonity = artifacts.require("Autonity")
const truffleAssert = require('truffle-assertions');
const assert = require('assert')
const utils = require('./utils.js');
const toBN = web3.utils.toBN;
const BN = require('bn.js');
const { Buffer } = require('node:buffer');
const BLOCK_REWARD = toBN(web3.utils.toWei("2", "ether")); // Constantinople+ ethash block reward

// this vote refund tests cannot be run on ganache, since it does not have the refund mechanism

contract("Oracle", accounts => {
  let oracle;
  // Use the genesis oracle voter (see genesis-tendermint.json oracleAddress).
  const voter = accounts[9];

  describe('Oracle vote refund', function() {
    before(async() => {
      // for testing the refund we need to interact with the oracle contract deployed at genesis.
      // the refund logic checks if the vote is sent to this specific oracle contract.
      oracle = await Oracle.at("0x47e9Fbef8C83A1714F1951F142132E6e90F5fa5D")
      autonity = await Autonity.at("0xbd770416a3345f91e4b34576cb804a576fa48eb1")
    })
    // NOTE: Oracle rounds are advanced by the protocol calling Autonity.finalize(), which then
    // calls Oracle.finalize(). On our 1-node FullFaker truffle testnet this protocol finalization
    // isn’t automatically executed every block, so tests must not rely on rounds changing “by waiting”.
    it('valid vote is refunded; double vote is not refunded', async function () {
      const currentEpoch = (await autonity.getEpochID()).toNumber()
      const round = await oracle.getRound()

      // ---- First vote (successful, reimbursable) ----
      const voteTx = await oracle.vote(0, [], 0, 0, {from: voter});
      const voteBlockNum = voteTx.receipt.blockNumber;
      const voteBlockHex = web3.utils.toHex(voteBlockNum);
      const votePrevBlockHex = web3.utils.toHex(voteBlockNum - 1);
      const voteBlock = await web3.eth.getBlock(voteBlockNum);
      const voteMiner = voteBlock.miner;

      // check that voter balance did not change (refund was successful)
      const voterBalBefore = toBN(await web3.eth.getBalance(voter, votePrevBlockHex));
      const voterBalAfter = toBN(await web3.eth.getBalance(voter, voteBlockHex));
      assert.equal(voterBalAfter.toString(), voterBalBefore.toString());

      /*
       * normally the baseFee gets sent to the Autonity Contract for redistribution and the tip to the block proposer (see core/state_transition.go TransitionDb())
       * since for the oracle vote we are refunding both the baseFee and the tip, the balance of the AC and the block proposer should not change.
       */
      const proposerBefore = toBN(await web3.eth.getBalance(voteMiner, votePrevBlockHex));
      const proposerAfter = toBN(await web3.eth.getBalance(voteMiner, voteBlockHex));
      // Ethash still mints the static block reward in this testnet setup. Ensure the vote didn't add any tx fees.
      assert.equal(proposerAfter.sub(proposerBefore).toString(), BLOCK_REWARD.toString(), "proposer received tx fees on reimbursed vote");

      const autonityBefore = toBN(await web3.eth.getBalance(autonity.address, votePrevBlockHex));
      const autonityAfter = toBN(await web3.eth.getBalance(autonity.address, voteBlockHex));
      assert.equal(autonityAfter.toString(), autonityBefore.toString(), "autonity balance changed");

      // make sure we are still in the same round
      const round2 = await oracle.getRound()
      assert.equal(round.toString(), round2.toString())

      // ---- Second vote (reverted, not reimbursable) ----
      const beforeFailBlock = await web3.eth.getBlockNumber()
      await utils.failsRevert(
        oracle.vote(0, [], 0, 0, {from: voter}),
        "already voted"
      );
      const afterFailBlock = await web3.eth.getBlockNumber()
      assert.ok(afterFailBlock > beforeFailBlock, "block number did not advance after reverted vote tx")

      const failedBlock = await web3.eth.getBlock(afterFailBlock)
      assert.ok(failedBlock.transactions.length > 0, "no txs found in block containing reverted vote")
      const failedTxHash = failedBlock.transactions[0]

      const tx = await web3.eth.getTransaction(failedTxHash);
      const receipt = await web3.eth.getTransactionReceipt(failedTxHash);
      assert.equal(receipt.status, false)

      // compute total gasCost, baseFee and effectiveTip
      const txBlock = await web3.eth.getBlock(tx.blockNumber)
      const baseFee = toBN(txBlock.baseFeePerGas || "0x0")
      const effectiveGasPriceRaw = receipt.effectiveGasPrice || tx.gasPrice
      assert.ok(effectiveGasPriceRaw != null, "missing effective gas price on reverted tx")
      const effectiveGasPrice = toBN(effectiveGasPriceRaw)
      const gasUsed = toBN(receipt.gasUsed)
      const gasCost = effectiveGasPrice.mul(gasUsed)
      const baseCost = baseFee.mul(gasUsed)
      const tip = effectiveGasPrice.sub(baseFee).mul(gasUsed)

      // gasCost = baseCost + tip
      assert.equal(gasCost.toString(), baseCost.add(tip).toString())

      // make sure that we are still in the same epoch --> no fee redistribution has happened
      const epoch = (await autonity.getEpochID()).toNumber()
      assert.equal(epoch, currentEpoch)

      // gasCost should have been spent
      const failBlockHex = web3.utils.toHex(afterFailBlock);
      const failPrevBlockHex = web3.utils.toHex(afterFailBlock - 1);
      const voterBalBeforeFail = toBN(await web3.eth.getBalance(voter, failPrevBlockHex));
      const voterBalAfterFail = toBN(await web3.eth.getBalance(voter, failBlockHex));
      assert.equal(voterBalAfterFail.toString(), voterBalBeforeFail.sub(gasCost).toString());

      /*
       * No refund in case of failed vote.
       * check that the basefee has been sent to the autonity contract and the tip to the proposer
       */
      const proposerBeforeFail = toBN(await web3.eth.getBalance(failedBlock.miner, failPrevBlockHex));
      const proposerAfterFail = toBN(await web3.eth.getBalance(failedBlock.miner, failBlockHex));
      assert.equal(
        proposerAfterFail.toString(),
        proposerBeforeFail.add(BLOCK_REWARD).add(tip).toString(),
        "proposer did not receive block reward + tip",
      );

      const autonityBeforeFail = toBN(await web3.eth.getBalance(autonity.address, failPrevBlockHex));
      const autonityAfterFail = toBN(await web3.eth.getBalance(autonity.address, failBlockHex));
      assert.equal(autonityAfterFail.toString(), autonityBeforeFail.add(baseCost).toString(), "autonity did not receive basefee");
    });
  });
});
