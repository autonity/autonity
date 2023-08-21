const OracleContract = artifacts.require("Oracle")
const truffleAssert = require('truffle-assertions');
const assert = require('assert')
const utils = require('./utils.js');
const VOTE_PERIOD = 3;
var Web3 = require("web3");

contract("Oracle", accounts => {
  before(async() => {
    await deployOracleContract();
  })

  const operator = accounts[0];
  const autonity = accounts[6];
  const deployer = accounts[5];
  // Accounts.
  const voterAccounts = [
      accounts[0],
      accounts[1],
      accounts[3],
      accounts[4],
  ]
  let oracle;
  //let w3ContractObj;
  let symbols = ["NTN/USD","NTN/AUD","NTN/CAD","NTN/EUR","NTN/GBP","NTN/JPY","NTN/SEK"]

  function getSymbolIndex(sym) {
    for (let i = 0; i < symbols.length; i++){
      const symbol = symbols[i];
      if (sym == symbol) {
        return i;
      }
    }
  }

  async function waitForNRounds(numRounds) {
    let round =  await oracle.getRound();
    let curRound = +round;
    while (+curRound != (+round + numRounds)) {
      await utils.timeout(100);
      await oracle.finalize({from:autonity});
      curRound = await oracle.getRound();
    }
  }


  async function deployOracleContract() {
    oracle = await OracleContract.new(voterAccounts, autonity, operator, symbols, VOTE_PERIOD);
    await web3.eth.sendTransaction({
      from: deployer,
      to: oracle.address,
      value: web3.utils.toWei("999854114594577", "wei")
    });
    //console.log("oracle address:"+oracle.address);
    // using websocket provider to access apis
    /*
    let w3       = new Web3(new Web3.providers.WebsocketProvider(web3.eth.currentProvider.host));
    w3ContractObj = new w3.eth.Contract(oracle.abi, oracle.address);
    w3ContractObj.events.NewRound({fromBlock:0})
        .on('data', (event) => {
         console.log("New Round Event:" + event.returnValues.round);
        })
        .on('error', console.error);
    w3ContractObj.events.NewSymbols({fromBlock:0})
        .on('data', (event) => {
         console.log("New Symbol Event:" + event.returnValues);
        })
        .on('error', console.error);
    w3ContractObj.events.Debug({fromBlock:0})
        .on('data', (event) => {
   //       console.log("Debug Event::" + event.returnValues.msg);
        })
        .on('error', console.error);
    w3ContractObj.events.DebugAS({fromBlock:0})
        .on('data', (event) => {
    //      console.log("DebugAS Event::" + event.returnValues.msg + " num: ", event.returnValues.num);
        })
        .on('error', console.error);

     */
  };

  const median = arr => {
    const mid = Math.floor(arr.length / 2),
        nums = [...arr].sort((a, b) => a - b);
    return arr.length % 2 !== 0 ? nums[mid] : Math.floor((nums[mid - 1] + nums[mid]) / 2);
  }

  let rounds = [];
  function updateRoundData(symChangeRound, syms, voterList) {
    for (let rIdx = symChangeRound; rIdx < rounds.length; rIdx++) {
      const voters = [];
      const expPrice = [];
      let pricesWithCommit = [];
      for (let vIdx = 0; vIdx < voterList.length; vIdx++) {
        let prices = [];
        prices =  Array.from({length: syms.length}, () => Math.floor(Math.random() * 100) + 1);
        pricesWithCommit = prices.map((item, index) => ({type:"int256", v: item}));
        const salt = rIdx + rIdx*10 +rIdx* 100;
        pricesWithCommit.push({type: "uint256", v: salt});
        pricesWithCommit.push({type: "address", v: voterList[vIdx]});
        const voter = {prices: prices, pricesWithCommit};
        voters.push(voter);
      }
      // expected price for each symbol in a round
      for (let sIdx = 0; sIdx < syms.length; sIdx++) {
        const symbolPrice = []
        for (let vIdx = 0; vIdx < voterList.length; vIdx++) {
          symbolPrice.push(voters[vIdx].prices[sIdx]);
        }
        expPrice.push(median(symbolPrice));
      }
      rounds[rIdx].expPrice = expPrice;
      rounds[rIdx].voters = voters
    }
  }
  // Sample Round data
  // { // round object
  //   "voters": [ { // list of all voter
  //     "prices": [ 95, 100, 6, 82, 70, 96, 10 ], //last round prices
  //     "pricesWithCommit": [ { "type": "int256", "v": 95 }, // current round prices, with salt and msg.sender value
  //       { "type": "int256", "v": 100 },
  //       { "type": "int256", "v": 6 },
  //       { "type": "int256", "v": 82 },
  //       { "type": "int256", "v": 70 },
  //       { "type": "int256", "v": 96 },
  //       { "type": "int256", "v": 10 },
  //       { "type": "uint256", "v": 111 },  // Salt value
  //       { "type": "address", "v": "0x627306090abaB3A6e1400e9345bC60c78a8BEf57" } // msg.sender
  //     ]
  //   },
  //   "expPrice": [ 12, 51, 24, 21, 65, 58, 38 ] // average price for each symbol
  // }

  function generateRoundData(numRounds, syms) {
    rounds = [];
    // Push first value as empty, since rounds start at 1
    for (let rIdx = 0; rIdx < numRounds; rIdx++) {
      let prices = [];
      let pricesWithCommit = [];
      let voter = {prices: prices, pricesWithCommit};
      const voters = [];
      const expPrice = [];
      let round = {};
      if (rIdx == 0) {
        // pushing empty data set to define valid keys
        voters.push(voter);
        round = {voters, expPrice };
        rounds.push(round);
        continue;
      }
      for (let vIdx = 0; vIdx < voterAccounts.length; vIdx++) {
        prices =  Array.from({length: syms.length}, () => Math.floor(Math.random() * 100) + 1);
        pricesWithCommit = prices.map((item, index) => ({type:"int256", v: item}));
        const salt = rIdx + rIdx*10 +rIdx* 100;
        pricesWithCommit.push({type: "uint256", v: salt});
        pricesWithCommit.push({type: "address", v: voterAccounts[vIdx]});
        voter = {prices: prices, pricesWithCommit};
        voters.push(voter);
      }

      // expected price for each symbol in a round
      for (let sIdx = 0; sIdx < syms.length; sIdx++) {
        const symbolPrice = []
        for (let vIdx = 0; vIdx < voterAccounts.length; vIdx++) {
          symbolPrice.push(voters[vIdx].prices[sIdx]);
        }
        expPrice.push(median(symbolPrice));
      }
      round = {voters, expPrice };
      rounds.push(round);
    }
  }

  describe('Contract initial state', function() {
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

  describe('Contract set api test', function() {
    it('Test sorting of voters', async function () {
      let newVoters = [
        "0xfF00000000000000000000000000000000000000",
        "0xaa00000000000000000000000000000000000000",
        "0x1100000000000000000000000000000000000000",
        "0x6600000000000000000000000000000000000000",
        "0xd228247B4f57587F6d2A479669e277699643135B",
        "0xF7cA6855Df4B0f725aC0dA6B54DD5CDF7E4c21d8",
      ]
      await oracle.setVoters(newVoters, {from:autonity});
      let updatedVoters = await oracle.getVoters();
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
      await oracle.setVoters(newVoters, {from:autonity});
      let updatedVoters = await oracle.getVoters();
      assert.deepEqual(
        newVoters.slice().sort(function (a, b) {
          return a.toLowerCase().localeCompare(b.toLowerCase());
        }),
        updatedVoters, "voters are not as expected"
      );
    });
    it('Test update voters - empty voter list', async function () {
      let newVoters = [];
      await truffleAssert.fails(
        oracle.setVoters(newVoters, {from:autonity}),
        truffleAssert.ErrorType.REVERT,
        "Voters can't be empty"
      );
    });

    it('Test update symbols', async function () {
      let newSymbols = ["NTN/USD","NTN/AUD","NTN/CAD","NTN/EUR","NTN/GBP","NTN/JPY"]
      await oracle.setSymbols(newSymbols, {from:operator});
      await waitForNRounds(2);
      let syms = await oracle.getSymbols();
      assert.deepEqual(syms, newSymbols, "symbols are not as expected");
    });

    it('Test update empty symbol list', async function () {
      let newSymbols = [];
      await truffleAssert.fails(
      oracle.setSymbols(newSymbols, {from:operator}),
      truffleAssert.ErrorType.REVERT,
      "symbols can't be empty"
    );
    });

    it('Test round update', async function () {
      const curRound = await oracle.getRound();
      await waitForNRounds(1)
      const newRound = await oracle.getRound();
      assert(+curRound+1 == +newRound, "round is not updated");
    });
  });


  describe('Contract running state', function() {
    beforeEach(async() => {
      await deployOracleContract()
    })
    it('Test update Symbols in same round ', async function () {
      let newSymbols = ["NTN/USD","NTN/AUD","NTN/CAD","NTN/EUR","NTN/GBP","NTN/JPY"];
      await oracle.setSymbols(newSymbols, {from:operator});
      newSymbols = ["NTN/USD","NTN/AUD","NTN/CAD","NTN/EUR","NTN/GBP"];
      await truffleAssert.fails(
        oracle.setSymbols(newSymbols, {from:operator}),
        truffleAssert.ErrorType.REVERT,
        "can't be updated in this round"
      );
    });

    it('Test update Symbols in subsequent round ', async function () {
      let newSymbols = ["NTN/USD","NTN/AUD","NTN/CAD","NTN/EUR","NTN/GBP","NTN/JPY"]
      await oracle.setSymbols(newSymbols, {from:operator});
      await waitForNRounds(1)
      await truffleAssert.fails(
        oracle.setSymbols(newSymbols, {from:operator}),
        truffleAssert.ErrorType.REVERT,
        "can't be updated in this round"
      );
    });

    it('Test vote - multiple votes in same round', async function () {
      generateRoundData(2, symbols);
      // round starts with 1
      let commit = web3.utils.soliditySha3(...rounds[1].voters[0].pricesWithCommit);
      // balance before vote
      await oracle.vote(commit, [], 0, {from:voterAccounts[0]});
      // second vote should revert
      await truffleAssert.fails(
        oracle.vote(commit, [], 0, {from:voterAccounts[0]}),
        truffleAssert.ErrorType.REVERT,
        "already voted",
      );
    });

    it('Test vote - empty report for existing voter', async function () {
      generateRoundData(3, symbols);
      // start round with 1
      for (let rId = 1; rId < rounds.length; rId++){
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++){
          const voter = round.voters[i];
          const commit = web3.utils.soliditySha3(...voter.pricesWithCommit);
          if (rId == 1) {
            await oracle.vote(commit, [], 0, {from:voterAccounts[i]});
          } else {
            const pricesWithcommit = rounds[rId-1].voters[i].pricesWithCommit;
            const salt = pricesWithcommit[pricesWithcommit.length -2].v;
            // should return because price report lenth doesn't match symbol length
            //TODO: should be verified by slashing event
            await oracle.vote(commit, [], salt, {from:voterAccounts[i]});
          }
        }
        await waitForNRounds(1)
      }
    });

    it('Test vote - retrieve price successfully for latest round', async function () {
      generateRoundData(10, symbols);
      for (let rId = 1; rId < rounds.length; rId++){
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++){
          const voter = round.voters[i];
          const commit = web3.utils.soliditySha3(...voter.pricesWithCommit);
          if (rId == 1) {
            await oracle.vote(commit, [], 0, {from:voterAccounts[i]});
          } else {
            const pricesWithcommit = rounds[rId-1].voters[i].pricesWithCommit;
            const salt = pricesWithcommit[pricesWithcommit.length -2].v;
            await oracle.vote(commit, rounds[rId-1].voters[i].prices, salt, {from:voterAccounts[i]});
          }
        }
        await waitForNRounds(1)
      }
      let roundData = await oracle.latestRoundData("NTN/GBP");
      //console.log("Price data received for round latest:"+ roundData[0] +
          //" price:"+roundData[1]+ " timestamp: "+roundData[2] + " status:"+roundData[3]);
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN/GBP");
      //console.log("exp Price:"+ rounds[roundID-1].expPrice[symIndex], " rec price:"+ +roundData[1]);
      assert(+roundData[1] == +rounds[roundID-1].expPrice[symIndex], "price is not as expected");
    });

    it('Test vote - skip voting round', async function () {
      generateRoundData(6, symbols);
      for (let rId = 1; rId < rounds.length; rId++){
        if (rId == 3) {
          //skipping round 3 - no voting
          await waitForNRounds(1)
          continue;
        }
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++){
          const voter = round.voters[i];
          const commit = web3.utils.soliditySha3(...voter.pricesWithCommit);
          if (rId == 1) {
            await oracle.vote(commit, [], 0, {from:voterAccounts[i]});
          } else {
            const pricesWithcommit = rounds[rId-1].voters[i].pricesWithCommit;
            const salt = pricesWithcommit[pricesWithcommit.length -2].v;
            await oracle.vote(commit, rounds[rId-1].voters[i].prices, salt, {from:voterAccounts[i]});
          }
        }
        await waitForNRounds(1)
      }
      // since we skipped round 3 - did not send commit and report both,
      // price will not be calculated in round 3 and 4 both
      let roundData = await oracle.getRoundData(3, "NTN/GBP");
      assert(roundData.status != 0, "status should not be success");
      roundData = await oracle.getRoundData(4, "NTN/GBP");
      assert(roundData.status != 0, "status should not be success");
    });

    it('Test vote - commit mismatch', async function () {
      generateRoundData(3, symbols);
      for (let rId = 1; rId < rounds.length; rId++){
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++){
          const voter = round.voters[i];

          let commit = web3.utils.soliditySha3(...voter.pricesWithCommit);
          if (rId == 1) {
              commit = 243432; // wrong commit values for all voters
              await oracle.vote(commit, [], 0, {from:voterAccounts[i]});
          } else {
            const pricesWithcommit = rounds[rId-1].voters[i].pricesWithCommit;
            const salt = pricesWithcommit[pricesWithcommit.length -2].v;
            await oracle.vote(commit, rounds[rId-1].voters[i].prices, salt, {from:voterAccounts[i]});
          }
        }
        await waitForNRounds(1)
      }
      let roundData = await oracle.latestRoundData("NTN/GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN/GBP");
      assert(roundData.status != 0, "status should not be success");
      assert(+roundData[1] != +rounds[roundID-1].expPrice[symIndex], "price should not be as expected");
    });

    it('Test vote committe change - Add one new member and remove one', async function () {
      generateRoundData(5, symbols);
      let voterCopy =  voterAccounts.slice();
      let rID;
      let commit;
      for (rId = 1; rId < rounds.length; rId++) {
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++){
          const voter = round.voters[i];
          commit = web3.utils.soliditySha3(...voter.pricesWithCommit);
          switch (rId) {
            case 1:
              await oracle.vote(commit, [], 0, {from:voterCopy[i]});
              break;
            default:
              const pricesWithcommit = rounds[rId-1].voters[i].pricesWithCommit;
              const salt = pricesWithcommit[pricesWithcommit.length -2].v;
              await oracle.vote(commit, rounds[rId-1].voters[i].prices, salt, {from:voterCopy[i]});
          }
        }
        switch (rId) {
          case 2:
            // update voters in round 2
            voterCopy = [
              accounts[0],
              accounts[1],
              accounts[3],
              accounts[5],
            ]
            await oracle.setVoters(voterCopy, {from:autonity});
            const syms =  await oracle.getSymbols();
            updateRoundData(rId, syms, voterCopy);
            break;
          case 3:
            // round-3 specific extra vote from the removed voter
            // use the data from the last voter in voter dataset, since we are voting
            // on behalf of him
            const prices = rounds[rId-1].voters[rounds[rId-1].voters.length-1].prices;
            const pricesWithCommit = rounds[rId-1].voters[rounds[rId-1].voters.length-1].pricesWithCommit;
            const salt = pricesWithCommit[pricesWithCommit.length -2].v;
            await oracle.vote(commit, prices, salt, {from:voterAccounts[3]});
            break;
        }
        await waitForNRounds(1)
      }
      let roundData = await oracle.latestRoundData("NTN/GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN/GBP");
      //console.log("exp Price:"+ rounds[roundID-1].expPrice[symIndex], " rec price:"+ +roundData[1]);
      assert(+roundData[1] == +rounds[roundID-1].expPrice[symIndex], "price is not as expected");
    });

    it('Test vote update symbols ', async function () {
      generateRoundData(5, symbols);
      let voterCopy =  voterAccounts.slice();
      let rId; // roundID
      let commit;
      for (rId = 1; rId < rounds.length; rId++) {
        switch (rId) {
          case 2:
            // update symbols in 2nd round
            const newSymbols = ["NTN/USD","NTN/AUD","NTN/CAD","NTN/EUR","NTN/GBP","NTN/JPY"]
            await oracle.setSymbols(newSymbols, {from:operator});
            break;
          case 3:
            const syms =  await oracle.getSymbols();
            //Imp Note: oracle.getVoters() can return voters in different order
            // to update roundData use local voterCopy/voterAccounts
            updateRoundData(rId, syms, voterCopy);
            break;
        }
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++){
          commit = web3.utils.soliditySha3(...round.voters[i].pricesWithCommit);
          switch (rId) {
            case 1:
              await oracle.vote(commit, [], 0, {from:voterCopy[i]});
              break;
            default:
              const pricesWithcommit = rounds[rId-1].voters[i].pricesWithCommit;
              const salt = pricesWithcommit[pricesWithcommit.length -2].v;
              await oracle.vote(commit, rounds[rId-1].voters[i].prices, salt, {from:voterCopy[i]});
          }
        }
        await waitForNRounds(1)
      }
      let roundData = await oracle.latestRoundData("NTN/GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN/GBP");
      //console.log("exp Price:"+ rounds[roundID-1].expPrice[symIndex], " rec price:"+ +roundData[1]);
      assert(+roundData[1] == +rounds[roundID-1].expPrice[symIndex], "price is not as expected");
    });

    it('Test vote - committe change and update symbols ', async function () {
      generateRoundData(10, symbols);
      let voterCopy =  voterAccounts.slice();
      let rId;
      let commit;
      let syms;
      for (rId = 1; rId < rounds.length; rId++) {
        const round = rounds[rId];
        for (let i = 0; i < round.voters.length; i++){
          const voter = round.voters[i];
          commit = web3.utils.soliditySha3(...voter.pricesWithCommit);
          switch (rId) {
            case 1:
              await oracle.vote(commit, [], 0, {from:voterCopy[i]});
              break;
            default:
              const pricesWithcommit = rounds[rId-1].voters[i].pricesWithCommit;
              const salt = pricesWithcommit[pricesWithcommit.length -2].v;
              await oracle.vote(commit, rounds[rId-1].voters[i].prices, salt, {from:voterCopy[i]});
          }
        }
        switch (rId) {
          case 2:
            // update voters in round 2
            voterCopy = [
              accounts[0],
              accounts[1],
              accounts[3],
              accounts[5],
            ]
            await oracle.setVoters(voterCopy, {from:autonity});
            syms =  await oracle.getSymbols();
            updateRoundData(rId, syms, voterCopy);
            // update symbols in 2nd round
            let newSymbols = ["NTN/USD","NTN/AUD","NTN/CAD","NTN/EUR","NTN/GBP","NTN/JPY"]
            await oracle.setSymbols(newSymbols, {from:operator});
            break;
          case 3:
            // round-3 specific extra vote from the removed voter
            // use the data from the last voter in voter dataset, since we are voting
            // on behalf of him
            syms =  await oracle.getSymbols();
            updateRoundData(rId, syms, voterCopy);
            // recalculate the commit since symbols are changed now
            commit = web3.utils.soliditySha3(...rounds[rId].voters[rId].pricesWithCommit);
            const prices = rounds[rId-1].voters[rounds[rId-1].voters.length-1].prices;
            const pricesWithCommit = rounds[rId-1].voters[rounds[rId-1].voters.length-1].pricesWithCommit;
            const salt = pricesWithCommit[pricesWithCommit.length -2].v;
            await oracle.vote(commit, prices, salt, {from:voterAccounts[3]});
            break;
        }
        await waitForNRounds(1)
      }
      let roundData = await oracle.latestRoundData("NTN/GBP");
      const roundID = roundData[0];
      const symIndex = getSymbolIndex("NTN/GBP");
      //console.log("exp Price:"+ rounds[roundID-1].expPrice[symIndex], " rec price:"+ +roundData[1]);
      assert(+roundData[1] == +rounds[roundID-1].expPrice[symIndex], "price is not as expected");
    });
  });
});
