const OracleContract = artifacts.require("Oracle");
const OracleAutonityMock = artifacts.require("OracleAutonityMockTest");
const truffleAssert = require('truffle-assertions');
const assert = require('assert')
const utils = require('./utils.js');
const VOTE_PERIOD = 3;
// var Web3 = require("web3");
const config = require("./config");

contract("Oracle", accounts => {
  for (let i = 0; i < accounts.length; i++) {
    console.log("account: ", i, accounts[i]);
  }

  const operator = accounts[5];
  const deployer = accounts[6];

  let oracle;
  let autonityContract;
  let autonity;

  before(async () => {
    await deployOracleContract();
  });
  // Accounts.
  const voterAccounts = [
    accounts[0],
    accounts[1],
    accounts[3],
    accounts[4],
  ]

  let symbols = ["NTN-USD", "NTN-AUD", "NTN-CAD", "NTN-EUR", "NTN-GBP", "NTN-JPY", "NTN-SEK"]

  function getSymbolIndex(sym) {
    for (let i = 0; i < symbols.length; i++) {
      const symbol = symbols[i];
      if (sym == symbol) {
        return i;
      }
    }
  }

  async function waitForNRounds(numRounds) {
    let round = await oracle.getRound();
    let curRound = +round;
    while (+curRound < (+round + numRounds)) {
      await utils.timeout(100);
      await autonityContract.finalize();
      curRound = await oracle.getRound();
    }
  }


  async function deployOracleContract() {
    autonityContract = await OracleAutonityMock.new(
      30, // _epochPeriod
      {from: deployer},
    );
    autonity = autonityContract.address;
    oracle = await OracleContract.new(
      voterAccounts, // voters
      voterAccounts, // node addresses
      voterAccounts, // treasuries
      symbols,
      {
        autonity,
        operator,
        votePeriod: VOTE_PERIOD,
        outlierDetectionThreshold: 100,
        outlierSlashingThreshold: 100,
        baseSlashingRate: 10,
      },
    );

    await autonityContract.setOracle(oracle.address);

    await web3.eth.sendTransaction({
      from: deployer,
      to: oracle.address,
      value: web3.utils.toWei("999854114594577", "wei")
    });
  }

  const calculatePrice = arr => {
    const num = arr.reduce((acc, val) => acc + (val.price * val.confidence), 0);
    const den = arr.reduce((acc, val) => acc + val.confidence, 0);
    return Math.floor(num / den);
  }

  let rounds = [];


  function generateCommit(prices, salt, voter) {
    return web3.utils.soliditySha3Raw(web3.eth.abi.encodeParameters([
      {
        "components": [
          {
            "internalType": "uint120",
            "name": "price",
            "type": "uint120"
          },
          {
            "internalType": "uint8",
            "name": "confidence",
            "type": "uint8"
          }
        ],
        "internalType": "struct IOracle.Report[]",
        "name": "_reports",
        "type": "tuple[]"
      },
      "uint256",
      "address"
    ], [
      prices,
      salt,
      voter
    ]))
  }

  function generateVoters(symbols, accounts) {
    const basePrices = Array.from(
      {length: symbols.length},
      () => ({price: Math.floor(Math.random() * 100) + 1, confidence: 100}),
    )
    return Array.from(Array(accounts.length), (_, vIdx) => ({
      prices: basePrices.map(p => ({...p, price: Math.floor(p.price * (1 + 0.02 * Math.random()))})),
      commit: "",
      salt: vIdx + vIdx * 10 + vIdx * 100,
      address: accounts[vIdx],
    }));
  }

  function calcExpectedPrice(symbols, round) {
    return Array.from({length: symbols.length}, (_, sIdx) => {
      const symbolPrice = round.voters.map(v => v.prices[sIdx]);
      return calculatePrice(symbolPrice);
    });
  }

  function setCommits(votingRounds) {
    votingRounds.forEach((r, rIdx) => {
      r.voters.forEach((v) => {
        const futureVote = votingRounds[(rIdx + 1) % votingRounds.length].voters.find((nv) => nv.address === v.address);
        if (!futureVote) {
          // voter is not in the next round, so we don't need to update the commit,
          // but we do need something here, so
          v.commit = generateCommit(
            v.prices,
            v.salt,
            v.address,
          );
          return v;
        }
        v.commit = generateCommit(
          futureVote.prices,
          futureVote.salt,
          v.address,
        );
      });
    })
  }

  function updateRoundData(fromRound, newSymbols, newVoters) {
    for (let rIdx = fromRound; rIdx < rounds.length; rIdx++) {
      // if this is the first round, we need votes from all voters falling out
      let voterAddresses = newVoters;
      if (rIdx === fromRound) {
        const oldVoters = rounds[rIdx].voters.map(v => v.address);
        const diff = oldVoters.filter(v => !newVoters.includes(v));
        voterAddresses = [...newVoters, ...diff];
      }

      // regenerate round with new voters
      rounds[rIdx] = {
        voters: generateVoters(newSymbols, voterAddresses),
        expPrice: []
      };

      // regenerate expected price
      rounds[rIdx].expPrice = calcExpectedPrice(newSymbols, rounds[rIdx]);
    }

    // regenerate commits from round right before
    setCommits(rounds.slice(fromRound-1));

    return rounds;
  }

  // Sample Round data
  // { // round object
  //   "voters": [ { // list of all voter
  //     "prices": [ 95, 100, 6, 82, 70, 96, 10 ], //last round prices
  //     "commit": "0x..."
  //     "salt": 2345,
  //     "address": "0x..."
  //       }... ]
  //   },
  //   "expPrice": [ 12, 51, 24, 21, 65, 58, 38 ] // average price for each symbol
  // }

  function generateRoundData(numRounds, syms) {
    // initialize array of rounds
    rounds = Array.from(Array(numRounds), (_, rIdx) => {
      return {
        voters: generateVoters(syms, voterAccounts),
        expPrice: [],
      };
    });

    // calculate expected price
    rounds.forEach(r => {
      r.expPrice = calcExpectedPrice(syms, r);
    });

    // set commits
    setCommits(rounds);

    return rounds;
  }

  describe('Contract initial state', function () {
    it('Test get symbols', async function () {
      let syms = await oracle.getSymbols();
      assert.deepEqual(syms.slice().sort(), symbols.slice().sort(), "symbols are not as expected");
    });

    it('Test get committee', async function () {
      let vs = await oracle.getVoters();
      assert.deepEqual(
        voterAccounts.slice().sort(function (a, b) {
          return a.toLowerCase().localeCompare(b.toLowerCase());
        }),
        vs, "voters are not as expected"
      );
    });

    it('Test get round', async function () {
      let round = await oracle.getRound();
      assert(round == 1, "round value must be one at initialization")
    });
  });

  describe('Contract set api test', function () {
    it('Test sorting of voters', async function () {
      let newVoters = [
        "0xfF00000000000000000000000000000000000000",
        "0xaa00000000000000000000000000000000000000",
        "0x1100000000000000000000000000000000000000",
        "0x6600000000000000000000000000000000000000",
        "0xd228247B4f57587F6d2A479669e277699643135B",
        "0xF7cA6855Df4B0f725aC0dA6B54DD5CDF7E4c21d8",
      ]
      await autonityContract.setVoters(
        newVoters, // voters
        newVoters, // treasury
        newVoters, // nodeAddresses
        {from: deployer}
      );
      let updatedVoters = await oracle.getNewVoters();
      //console.log(updatedVoters)
      assert.deepEqual(
        newVoters.slice().sort(function (a, b) {
          return a.toLowerCase().localeCompare(b.toLowerCase());
        }),
        updatedVoters, "voters are not as expected"
      );
    });
    it('Test update voters', async function () {
      let newVoters = [
        accounts[0],
        accounts[1],
        accounts[3],
        accounts[5],
      ]
      await autonityContract.setVoters(
        newVoters, // voters
        newVoters, // treasuries
        newVoters, // nodeAddresses
        {from: deployer});
      let updatedVoters = await oracle.getNewVoters();
      assert.deepEqual(
        newVoters.slice().sort(function (a, b) {
          return a.toLowerCase().localeCompare(b.toLowerCase());
        }),
        updatedVoters, "voters are not as expected"
      );

      await waitForNRounds(5);

      let vs = await oracle.getVoters();
      assert.deepEqual(
        newVoters.slice().sort(function (a, b) {
          return a.toLowerCase().localeCompare(b.toLowerCase());
        }),
        vs, "voters are not as expected"
      );

    });
    it('Test update voters - empty voter list', async function () {
      let newVoters = [];
      await truffleAssert.fails(
        autonityContract.setVoters(
          newVoters, // voters
          newVoters, // treasuries
          newVoters, // nodeAddresses
          {from: deployer}),
        truffleAssert.ErrorType.REVERT,
        "Voters can't be empty"
      );
    });

    it('Test update symbols', async function () {
      let newSymbols = ["NTN-USD", "NTN-AUD", "NTN-CAD", "NTN-EUR", "NTN-GBP", "NTN-JPY"]
      await oracle.setSymbols(newSymbols, {from: operator});
      await waitForNRounds(2);
      let syms = await oracle.getSymbols();
      assert.deepEqual(syms, newSymbols, "symbols are not as expected");
    });

    it('Test update empty symbol list', async function () {
      let newSymbols = [];
      await truffleAssert.fails(
        oracle.setSymbols(newSymbols, {from: operator}),
        truffleAssert.ErrorType.REVERT,
        "symbols can't be empty"
      );
    });

    it('Test round update', async function () {
      const curRound = await oracle.getRound();
      await waitForNRounds(1)
      const newRound = await oracle.getRound();
      assert(+curRound + 1 == +newRound, "round is not updated");
    });
    //TODO(tariq) low priority. add test that checks that only voters can actually call vote
    //TODO(tariq) low priority. add test that checks that only autonity can call finalize(), setVoters() and setOperator()
    //TODO(tariq) low priority. add test that checks that only operator can call setSymbols()
  });

  describe('Contract running state', function () {
    beforeEach(async () => {
      await deployOracleContract();
    });

    it('Test update Symbols in same round ', async function () {
      let newSymbols = ["NTN-USD", "NTN-AUD", "NTN-CAD", "NTN-EUR", "NTN-GBP", "NTN-JPY"];
      await oracle.setSymbols(newSymbols, {from: operator});
      newSymbols = ["NTN-USD", "NTN-AUD", "NTN-CAD", "NTN-EUR", "NTN-GBP"];
      await truffleAssert.fails(
        oracle.setSymbols(newSymbols, {from: operator}),
        truffleAssert.ErrorType.REVERT,
        "can't be updated in this round"
      );
    });

    it('Test update Symbols in subsequent round ', async function () {
      let newSymbols = ["NTN-USD", "NTN-AUD", "NTN-CAD", "NTN-EUR", "NTN-GBP", "NTN-JPY"]
      await oracle.setSymbols(newSymbols, {from: operator});
      await waitForNRounds(1)
      await truffleAssert.fails(
        oracle.setSymbols(newSymbols, {from: operator}),
        truffleAssert.ErrorType.REVERT,
        "can't be updated in this round"
      );
    });

    it('Test vote - multiple votes in same round', async function () {
      generateRoundData(2, symbols);
      const commit = rounds[0].voters[0].commit;
      // balance before vote
      await oracle.vote(commit, [], 0, 0, {from: voterAccounts[0]});
      // second vote should revert
      await truffleAssert.fails(
        oracle.vote(commit, [], 0, 0, {from: voterAccounts[0]}),
        truffleAssert.ErrorType.REVERT,
        "already voted",
      );
    });

    it('Test vote - empty report for existing voter', async function () {
      generateRoundData(3, symbols);
      // start round with 1
      for (let rId = 0; rId < rounds.length; rId++) {
        const {voters} = rounds[rId];
        await Promise.all(voters.map(async (voter, i) => {
          if (rId === 0) {
            await oracle.vote(voter.commit, [], 0, 0, {from: voterAccounts[i]});
          } else {
            const {commit, salt} = voter;
            // vote with empty report
            await oracle.vote(commit, [], salt, 0, {from: voterAccounts[i]});
            //TODO: should be verified by slashing event
          }
        }));
        await waitForNRounds(1)
      }
    });

    it('Test vote - retrieve price successfully for latest round', async function () {
      generateRoundData(10, symbols);
      for (let rId = 0; rId < rounds.length; rId++) {
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++) {
          const {commit, prices, salt} = round.voters[i];
          //TODO(tariq) low priority. check that vote correctly trigger the update of the `reports` structure
          if (rId === 0) {
            await oracle.vote(commit, [], 0, 0, {from: voterAccounts[i]});
          } else {
            await oracle.vote(commit, prices, salt, 0, {from: voterAccounts[i]});
          }
        }
        await waitForNRounds(1)
      }
      //TODO(tariq) low priority. check that the price is what we expect for the other symbols as well
      let roundData = await oracle.latestRoundData("NTN-GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN-GBP");
      assert(+roundData[1] === +rounds[roundID - 1].expPrice[symIndex], "price is not as expected");
    });

    it('Test vote - skip voting round', async function () {
      generateRoundData(6, symbols);
      for (let rId = 0; rId < rounds.length; rId++) {
        if (rId === 2) {
          //skipping round 3 - no voting
          await waitForNRounds(1)
          continue;
        }
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++) {
          const {commit, prices, salt} = round.voters[i];
          if (rId === 0) {
            await oracle.vote(commit, [], 0, 0, {from: voterAccounts[i]});
          } else {
            await oracle.vote(commit, prices, salt, 0, {from: voterAccounts[i]});
          }
        }
        await waitForNRounds(1)
      }

      const sIdx = getSymbolIndex("NTN-GBP");
      const round2Data = await oracle.getRoundData(2, "NTN-GBP");
      const round3Data = await oracle.getRoundData(3, "NTN-GBP");
      const round4Data = await oracle.getRoundData(4, "NTN-GBP");

      // since we skipped round 3 - did not send commit and report both,
      // price will be reused from round 2
      assert(
        +round2Data.price === +rounds[1].expPrice[sIdx],
        "price for round 2 should be as expected"
      );
      assert(round2Data.price === round3Data.price, "price for round 3 should equal to round 2");

      // round 4 is also invalid so should be equal to round 3
      assert(round4Data.price === round3Data.price, "price for round 4 should equal to round 3");
    });

    it('Test vote - commit mismatch', async function () {
      generateRoundData(2, symbols);
      for (let rId = 0; rId < rounds.length; rId++) {
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++) {
          let {commit, prices, salt} = round.voters[i];
          //TODO(tariq) low priority. check that invalid vote correctly trigger the update of the `reports` structure
          //  reports[symbols[i]][msg.sender] = INVALID_PRICE;
          if (rId === 0) {
            commit = 243432; // wrong commit values for all voters
            await oracle.vote(commit, [], 0, 0, {from: voterAccounts[i]});
          } else {
            await oracle.vote(commit, prices, salt, 0, {from: voterAccounts[i]});
          }
        }
        await waitForNRounds(1)
      }
      let roundData = await oracle.latestRoundData("NTN-GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN-GBP");
      assert(roundData.success === false, "price should not be valid");
      assert(+roundData[1] !== +rounds[roundID - 1].expPrice[symIndex], "price should not be as expected");
    });

    it('Test voter committee change - Add one new member and remove one', async function () {
      const newVoters = [
        accounts[0],
        accounts[1],
        accounts[3],
        accounts[5],
      ]
      generateRoundData(6, symbols);
      // voters change in round 3
      updateRoundData(3, symbols, newVoters);

      for (let rId = 0; rId < 5; rId++) {
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++) {
          const {commit, prices, salt, address} = round.voters[i];
          if (rId === 0) {
            await oracle.vote(commit, [], 0, 0, {from: address});
          } else {
            await oracle.vote(commit, prices, salt, 0, {from: address});
          }
        }
        if (rId === 2) {
          // update voters in round 2
          await autonityContract.setVoters(
            newVoters, // voters
            newVoters, // treasuries
            newVoters, // nodeAddresses
            {from: deployer}
          );
        }
        await waitForNRounds(1)
        console.log("finished round", rId+1);
      }
      let roundData = await oracle.latestRoundData("NTN-GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN-GBP");
      assert(+roundData[1] === +rounds[roundID - 1].expPrice[symIndex], "price is not as expected");
    });

    it('Test vote update symbols ', async function () {
      generateRoundData(6, symbols);
      const newSymbols = ["NTN-USD", "NTN-AUD", "NTN-CAD", "NTN-EUR", "NTN-GBP", "NTN-JPY"]
      // symbols change in round 4, new commits from round 3
      updateRoundData(4, newSymbols, voterAccounts);
      for (let rId = 0; rId < 5; rId++) {
        for (let i = 0; i < rounds[rId].voters.length; i++) {
          const {commit, prices, salt, address} = rounds[rId].voters[i];
          if (rId === 0) {
            await oracle.vote(commit, [], 0, 0, {from: address});
          } else {
            await oracle.vote(commit, prices, salt, 0, {from: address});
          }
        }
        if (rId === 2) {
          // update symbols in 2nd round
          await oracle.setSymbols(newSymbols, {from: operator});
        }
        await waitForNRounds(1)
      }
      let roundData = await oracle.latestRoundData("NTN-GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN-GBP");
      //console.log("exp Price:"+ rounds[roundID-1].expPrice[symIndex], " rec price:"+ +roundData[1]);
      assert(+roundData[1] === +rounds[roundID - 1].expPrice[symIndex], "price is not as expected");
    });

    it('Test vote - committee change and update symbols ', async function () {
      generateRoundData(10, symbols);

      let newVoters = [
        accounts[0],
        accounts[1],
        accounts[3],
        accounts[5],
      ];
      let newSymbols = ["NTN-USD", "NTN-AUD", "NTN-CAD", "NTN-EUR", "NTN-GBP", "NTN-JPY"];

      // voters change in round 3
      updateRoundData(3, symbols, newVoters);

      // symbols change in round 4, new commits from round 3
      updateRoundData(4, newSymbols, newVoters);

      for (let rId = 0; rId < rounds.length - 1; rId++) {
        for (let i = 0; i < rounds[rId].voters.length; i++) {
          const {commit, prices, salt, address} = rounds[rId].voters[i];
          if (rId === 0) {
            await oracle.vote(commit, [], 0, 0, {from: address});
          } else {
            await oracle.vote(commit, prices, salt, 0, {from: address});
          }
        }
        if (rId === 2) {
          // update voters in round 2
          await autonityContract.setVoters(
            newVoters, // voters
            newVoters, // treasuries
            newVoters, // nodeAddresses
            {from: deployer}
          );
          await oracle.setSymbols(newSymbols, {from: operator});
        }
        await waitForNRounds(1)
      }
      let roundData = await oracle.latestRoundData("NTN-GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN-GBP");
      //console.log("exp Price:"+ rounds[roundID-1].expPrice[symIndex], " rec price:"+ +roundData[1]);
      assert(+roundData[1] === +rounds[roundID - 1].expPrice[symIndex], "price is not as expected");
    });
  });
});
