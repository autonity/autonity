'use strict';

const MIN_BASE_FEE = 5000;
const COMMITTEE_SIZE = 1000;
const EPOCH_PERIOD = 30;
const DELEGATION_RATE = 100;
const UN_BONDING_PERIOD = 60;
const TREASURY_FEE = "10000000000000000";
const MIN_EPOCH_PERIOD = 30;
const VERSION = 0;
const ZERO_ADDRESS = "0x0000000000000000000000000000000000000000";

const ACCOUNTABILITY_CONFIG = {
        "innocenceProofSubmissionWindow": 30,
        "latestAccountabilityEventsRange": 256,
        "baseSlashingRateLow": 500,
        "baseSlashingRateMid": 1000,
        "collusionFactor": 550,
        "historyFactor": 750,
        "jailFactor": 60,
        "slashingRatePrecision": 10000
    };

const GENESIS_ENODES = [
        "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
        "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
        "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
        "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303"
    ];

// precomputed using aut validator compute-address
// TODO(lorenzo) derive them from enodes or privatekeys
const GENESIS_NODE_ADDRESSES = [
        "0x850C1Eb8D190e05845ad7F84ac95a318C8AaB07f",
        "0x4AD219b58a5b46A1D9662BeAa6a70DB9F570deA5",
        "0xc443C6c6AE98F5110702921138D840e77dA67702",
        "0x09428E8674496e2D1E965402F33A9520c5fCBbE2",
    ];

const BASE_VALIDATOR = {
        "selfBondedStake": 0,
        "selfUnbondingStakeLocked": 0,
        "totalSlashed": 0,
        "jailReleaseBlock": 0,
        "provableFaultCount": 0,
        "liquidSupply": 0,
        "registrationBlock": 0,
        "state": 0,
        "liquidContract": ZERO_ADDRESS,
        "selfUnbondingStake": 0,
        "selfUnbondingShares": 0,
        "unbondingStake": 0,
        "unbondingShares": 0,
        "key": "0x00",
    };

const GENESIS_PRIVATE_KEYS = [
    "a4b489752489e0f47e410b8e8cbb1ac1b56770d202ffd45b346ca8355c602c91",
    "aa4b77b1305f8f265e81599587c623d8950624f3e1bd9c121ef2461a7a1e7527",
    "4ec99383dc50aa3f3117fcbfba7b69188ba60d3418185fb353c9a69d066e55d9",
    "0c8698f456533170fe07c6dcb753d47bef8bedd46443efa57a859c989887b56b",
];

function autonityConfig(operator, treasuryAccount) {
    return {
        "policy": {
            "treasuryFee": TREASURY_FEE,
            "minBaseFee": MIN_BASE_FEE,
            "delegationRate": DELEGATION_RATE,
            "unbondingPeriod" : UN_BONDING_PERIOD,
            "treasuryAccount": treasuryAccount,
        },
        "contracts": {
            "oracleContract" : ZERO_ADDRESS, // gets updated in deployContracts()
            "accountabilityContract": ZERO_ADDRESS, // gets updated in deployContracts()
            "acuContract" :ZERO_ADDRESS,
            "supplyControlContract" :ZERO_ADDRESS,
            "stabilizationContract" :ZERO_ADDRESS,
        },
        "protocol": {
            "operatorAccount": operator,
            "epochPeriod": EPOCH_PERIOD,
            "blockPeriod": MIN_EPOCH_PERIOD,
            "committeeSize": COMMITTEE_SIZE,
        },
        "contractVersion": VERSION,
    };
}

// accounts[2] is skipped because it is used as a genesis validator when running against autonity
// this can cause interference in reward distribution tests
function validators(accounts) {
    return [
        { ...BASE_VALIDATOR,
            "treasury": accounts[0],
            "nodeAddress": GENESIS_NODE_ADDRESSES[0],
            "oracleAddress": accounts[0],
            "enode": GENESIS_ENODES[0],
            "commissionRate": 100,
            "bondedStake": 100,
        },
        { ...BASE_VALIDATOR,
            "treasury": accounts[1],
            "nodeAddress": GENESIS_NODE_ADDRESSES[1],
            "oracleAddress": accounts[1],
            "enode": GENESIS_ENODES[1],
            "commissionRate": 100,
            "bondedStake": 90,
        },
        { ...BASE_VALIDATOR,
            "treasury": accounts[3],
            "nodeAddress": GENESIS_NODE_ADDRESSES[2],
            "oracleAddress": accounts[3],
            "enode": GENESIS_ENODES[2],
            "commissionRate": 100,
            "bondedStake": 110,
        },
        { ...BASE_VALIDATOR,
            "treasury": accounts[4],
            "nodeAddress": GENESIS_NODE_ADDRESSES[3],
            "oracleAddress": accounts[4],
            "enode": GENESIS_ENODES[3],
            "commissionRate": 100,
            "bondedStake": 120,
        },
    ];
}

module.exports = {
    MIN_BASE_FEE: MIN_BASE_FEE,
    COMMITTEE_SIZE: COMMITTEE_SIZE,
    VERSION: VERSION,
    ACCOUNTABILITY_CONFIG: ACCOUNTABILITY_CONFIG,
    GENESIS_ENODES: GENESIS_ENODES,
    GENESIS_NODE_ADDRESSES: GENESIS_NODE_ADDRESSES,
    BASE_VALIDATOR: BASE_VALIDATOR,
    GENESIS_PRIVATE_KEYS: GENESIS_PRIVATE_KEYS,
    autonityConfig: autonityConfig,
    validators: validators,
};