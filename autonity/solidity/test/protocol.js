'use strict';
const assert = require('assert');
const truffleAssert = require('truffle-assertions');
const utils = require('./utils.js');
const liquidContract = artifacts.require("Liquid")
const toBN = web3.utils.toBN;

// testing protocol contracts interactions.

contract('Protocol', function (accounts) {
    before(async function () {
      console.log("\tAttempting to mock enode verifier precompile. Will (rightfully) fail if running against Autonity network")
      await utils.mockEnodePrecompile()
    });

  for (let i = 0; i < accounts.length; i++) {
    console.log("account: ", i, accounts[i]);
  }

  const operator = accounts[5];
  const deployer = accounts[6];
  const anyAccount = accounts[7];
  const name = "Newton";
  const symbol = "NTN";
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


  // initial validators ordered by bonded stake
  const orderedValidatorsList = [
    genesisNodeAddresses[0],
    genesisNodeAddresses[1],
    genesisNodeAddresses[2],
    genesisNodeAddresses[3],
  ];

  let autonity;
  describe('After effects of slashing', function () {
      beforeEach(async function () {
        autonity = await utils.deployContracts(validators, autonityConfig, accountabilityConfig,  deployer, operator);
      });
    it('unbondingShares:unbondingStake conversion ratio', async function () {
        /* TODO(tariq) issue multiple unbonding requests (both selfBonded and not) in different epochs, interleaved with slashing events
         * and check that the unbonding shares related fields change accordingly. Relevant fields to check:
         * struct Validator {
              uint256 bondedStake;
              uint256 unbondingStake;
              uint256 unbondingShares;
              uint256 selfBondedStake;
              uint256 selfUnbondingStake;
              uint256 selfUnbondingShares;
              uint256 liquidSupply;
            }
            struct UnbondingRequest {
              address payable delegator;
              address delegatee;
              uint256 amount;
              uint256 unbondingShare;
              bool selfDelegation;
            }
          * When a slashing event happens, the ratio of exchange between "unbonding shares" <--> "unbonding stake" changes and it is not 1:1 anymore
          * Example:
          * unbondingShares = 100
          * unbondingStake = 100
          * exchange ratio 1 share : 1 NTN
          *
          * Then slash of 50 NTNs. Now:
          *
          * unbondingShares = 100
          * unbondingStake = 50
          * exchange ratio 2 share : 1 NTN
          *
          * The value of the share decreased, as you now need 2 shares to redeem 1 NTN
          */
    });
    it.skip('unbondingShares:unbondingStake 100% slash edge case', async function () {
      // see https://github.com/autonity/autonity/issues/819
      // low priority for now as we already now that the problem is there, however we need a test for the future fix
    });
    it('LNTN:NTN conversion ratio', async function () {
      // TODO(tariq) the same logic used for converting unbondingShares <--> unbondingStake is used for the NTN:LNTN conversion
      // The same example applies when NTN is slashed --> LNTN loses value --> conversion rate is now not 1:1 anymore
      // issue multiple bond and unbond request with interleaved slashing events, and check that the NTN:LNTN ratio is always what we expect
    });
    it.skip('LNTN:NTN 100% slash edge case', async function () {
      // see https://github.com/autonity/autonity/issues/819
      // low priority for now as we already now that the problem is there, however we need a test for the future fix
    });
    it('jailed validator rewards go to proof reporter', async function () {
      /* TODO(tariq) verify that if a validators is jailed (it has been slashed for a misbehaviour)
       * his share of rewards goes to the proof reporter. relevant code in _performRedistribution() function of Autonity.sol
       *          if (_val.state == ValidatorState.jailed) {
       *             config.contracts.accountabilityContract.distributeRewards{value: _reward}(committee[i].addr);
       *             continue;
       *         }
       */         

    });
  });
});
