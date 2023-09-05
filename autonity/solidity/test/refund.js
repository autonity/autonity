const Oracle = artifacts.require("Oracle")
const Autonity = artifacts.require("Autonity")
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
      autonity = await Autonity.at("0xbd770416a3345f91e4b34576cb804a576fa48eb1")

      // register a new validator that we will use for voting in the refund test.
      // using the genesis validator makes fees computations harder because it receives the tip when a block is added to the chain
      let treasury = accounts[8];
      let oracleAddr = accounts[8]
      let nodeAddr = "0xDE03B7806f885Ae79d2aa56568b77caDB0de073E"
      let enode = "enode://a7ecd2c1b8c0c7d7ab9cc12e620605a762865d381eb1bc5417dcf07599571f84ce5725f404f66d3e254d590ae04e4e8f18fe9e23cd29087d095a0c37d0443252@3.209.45.79:30303"
      let nodeKey = "e59be7e486afab41ec6ef6f23746d78e5dbf9e3f9b0ac699b5566e4f675e976b"
      let treasuryProof = web3.eth.accounts.sign(treasury, nodeKey);
      let oracleProof = await web3.eth.sign(treasury, oracleAddr);
      let multisig = treasuryProof.signature + oracleProof.substring(2)
      await autonity.registerValidator(enode, oracleAddr, multisig, {from: treasury});

      // bond to it
      await autonity.bond(nodeAddr, 10, {from: accounts[8]});

      // wait for epoch to end so that accounts[8] becomes a committee member
      let currentEpoch = (await autonity.epochID()).toNumber()
      for (;;){
        await utils.timeout(5000)
        let epoch = (await autonity.epochID()).toNumber()
        if(epoch > currentEpoch){
          break;
        }
      }
      
      // wait for an additional oracle round so that he is a "fully valid" oracle voter
      let currentRound = (await oracle.getRound()).toNumber()
      for (;;){
        await utils.timeout(5000)
        let round = (await oracle.getRound()).toNumber()
        if(round > currentRound){
          break;
        }
      }

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
    it('fee is refunded for valid vote', async function () {
      const proposer = accounts[2];
      const origBalance = toBN(await web3.eth.getBalance(accounts[8]));
      const proposerInitBalance = toBN(await web3.eth.getBalance(proposer));
      const autonityInitBalance = toBN(await web3.eth.getBalance(autonity.address));

      await oracle.vote(0, [], 0, {from:accounts[8]});
     
      // check that voter balance did not change (refund was successfull)
      const updatedBalance = toBN(await web3.eth.getBalance(accounts[8]))
      assert.equal(updatedBalance.toString(), origBalance.toString());

      /*
       * normally the baseFee gets sent to the Autonity Contract for redistribution and the tip to the block proposer (see core/state_transition.go TransitionDb())
       * since for the oracle vote we are refunding both the baseFee and the tip, the balance of the AC and the block proposer should not change.
       * add asserts for these two conditions to ensure that we are not duplicating money.
       */
      assert.equal(await web3.eth.getBalance(proposer), proposerInitBalance.toString(), "proposer balance changed");
      assert.equal(await web3.eth.getBalance(autonity.address), autonityInitBalance.toString(), "autonity balance changed");
    });
    it('double vote, only first is refunded', async function () {
      let currentEpoch = (await autonity.epochID()).toNumber()
      const proposer = accounts[2];
      let proposerInitBalance = toBN(await web3.eth.getBalance(proposer));
      let autonityInitBalance = toBN(await web3.eth.getBalance(autonity.address));
      
      // first vote gets refunded
      let origBalance = toBN(await web3.eth.getBalance(accounts[8]));
      let round = await oracle.getRound()
      await oracle.vote(0, [], 0, {from:accounts[8]});
      let updatedBalance = toBN(await web3.eth.getBalance(accounts[8]))
      assert.equal(updatedBalance.toString(), origBalance.toString());
      
      /*
       * normally the baseFee gets sent to the Autonity Contract for redistribution and the tip to the block proposer (see core/state_transition.go TransitionDb())
       * since for the oracle vote we are refunding both the baseFee and the tip, the balance of the AC and the block proposer should not change.
       * add asserts for these two conditions to ensure that we are not duplicating money.
       */
      assert.equal(await web3.eth.getBalance(proposer), proposerInitBalance.toString(), "proposer balance changed");
      assert.equal(await web3.eth.getBalance(autonity.address), autonityInitBalance.toString(), "autonity balance changed");

      // make sure we are still in the same round
      round2 = await oracle.getRound()
      assert.equal(round.toString(),round2.toString())
      
      // second vote should fail, with !=0 gas expense
      await truffleAssert.fails(
        oracle.vote(0,[],0,{from:accounts[8]}),
        truffleAssert.ErrorType.REVERT,
        "already voted"
      );
      let failedTxHash = (await web3.eth.getBlock("latest")).transactions[0]
      const tx = await web3.eth.getTransaction(failedTxHash);
      const receipt = await web3.eth.getTransactionReceipt(failedTxHash);
      assert.equal(receipt.status,false)
      
      // compute total gasCost, baseFee and effectiveTip
      txBlock = await web3.eth.getBlock(tx.blockNumber)
      baseFee = toBN(txBlock.baseFeePerGas)
      const gasCost = toBN(tx.gasPrice).mul(toBN(receipt.gasUsed));
      const effectiveTip = BN.min(toBN(tx.maxPriorityFeePerGas),toBN(tx.maxFeePerGas).sub(baseFee))
      const tip = effectiveTip.mul(toBN(receipt.gasUsed));
      const baseCost = toBN(baseFee).mul(toBN(receipt.gasUsed));

      // gasCost = baseCost + tip
      assert.equal(gasCost.toString(),baseCost.add(tip).toString())
    
      // make sure that we are still in the same epoch --> no fee redistribution has happened
      let epoch = (await autonity.epochID()).toNumber()
      assert.equal(epoch,currentEpoch)
      
      // gasCost should have been spent
      updatedBalance2 = toBN(await web3.eth.getBalance(accounts[8]))
      assert.equal(updatedBalance2.toString(), updatedBalance.sub(gasCost).toString());

      /*
       * No refund in case of failed vote.
       * check that the basefee has been sent to the autonity contract and the tip to the proposer
       * the proposer is always accounts[2] since we are running on a 1-node autonity test network
       */
      assert.equal(await web3.eth.getBalance(proposer), proposerInitBalance.add(tip).toString(), "proposer did not receive tip");
      assert.equal(await web3.eth.getBalance(autonity.address), autonityInitBalance.add(baseCost).toString(), "autonity did not receive basefee");
    });
  });
});
