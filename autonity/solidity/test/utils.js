const assert = require('assert');
const util = require('util');
const config = require('./config');
const exec = util.promisify(require('child_process').exec);
const Autonity = artifacts.require("Autonity");
const Accountability = artifacts.require("Accountability");
const UpgradeManager = artifacts.require("UpgradeManager");
const Oracle = artifacts.require("Oracle")
const Acu = artifacts.require("ACU")
const SupplyControl = artifacts.require("SupplyControl")
const Stabilization = artifacts.require("Stabilization")
const InflationController = artifacts.require("InflationController")
const NonStakableVesting = artifacts.require("NonStakableVesting")
const AutonityTest = artifacts.require("AutonityTest");
const mockEnodeVerifier = artifacts.require("MockEnodeVerifier")
const mockCommitteeSelector = artifacts.require("MockCommitteeSelector")
const EC = require('elliptic').ec;
const ec = new EC('secp256k1');
const keccak256 = require('keccak256');
const ethers = require('ethers');
const truffleAssert = require('truffle-assertions');

// Validator Status in Autonity Contract
const ValidatorState = {
  active : 0,
  paused : 1,
  jailed : 2,
  jailbound : 3
}

// todo: remove this function?
// end epoch so the LastEpochBlock is closer
// then set epoch period 
async function shortenEpochPeriod(autonity, epochPeriod, operator, deployer) {
  await endEpoch(autonity, operator, deployer);
  await autonity.setEpochPeriod(epochPeriod, {from: operator});

  let currentEpoch = (await autonity.epochID()).toNumber();
  let lastEpochBlock = (await autonity.getLastEpochBlock()).toNumber();
  let oldEpochPeriod = (await autonity.getEpochPeriod()).toNumber();
  let nextEpochBlock = lastEpochBlock+oldEpochPeriod;
  let currentHeight = await web3.eth.getBlockNumber();

  // close epoch to take the shorten epoch into active state.
  console.log("currentHeight: ", currentHeight, "lastEpochBlock: ",
      lastEpochBlock, "oldEPeriod: ", oldEpochPeriod, "nextEpochBlock: ", nextEpochBlock);
  if (currentHeight > nextEpochBlock) {
    console.log("current height is higher than the next epoch block, finalize epoch at once");
    await autonity.finalize({from: deployer})
  } else {
    console.log("current height is lower than the next epoch block, try to finalize epoch");
    for (let i=currentHeight;i<=nextEpochBlock;i++) {
      let height = await web3.eth.getBlockNumber()
      console.log("try to finalize epoch", "height: ", height, "next epoch block: ", nextEpochBlock);
      autonity.finalize({from: deployer})
      let epochID = (await autonity.epochID()).toNumber()
      if (epochID === currentEpoch+1) {
        break;
      }
      await waitForNewBlock(height);
    }
  }
}

async function endEpoch(contract,operator,deployer){
    let lastEpochBlock = (await contract.getLastEpochBlock()).toNumber();
    let oldEpochPeriod = (await contract.getEpochPeriod()).toNumber();
    let nextEpochBlock = lastEpochBlock+oldEpochPeriod;
    let currentHeight = await web3.eth.getBlockNumber();
    let currentEpoch = (await contract.epochID()).toNumber();

    for (let i=currentHeight;i<=nextEpochBlock;i++) {
      let height = await web3.eth.getBlockNumber()
      console.log("try to finalize epoch", "height: ", height, "next epoch block: ", nextEpochBlock);
      contract.finalize({from: deployer})
      let epochID = (await contract.epochID()).toNumber()
      if (epochID === currentEpoch+1) {
        break;
      }
      await waitForNewBlock(height);
    }

    let newEpoch = (await contract.epochID()).toNumber()
    assert.equal(newEpoch, currentEpoch+1)
}

async function validatorState(autonity, validatorAddresses) {
  let expectedValInfo = [];
  for (let i = 0; i < validatorAddresses.length; i++) {
    expectedValInfo.push(await autonity.getValidator(validatorAddresses[i]));
  }
  return expectedValInfo;
}

async function bulkBondingRequest(autonity, operator, delegators, delegatee, tokenMint) {

  let bondingCount = 0;
  for (let i = 0; i < delegators.length; i++) {
    let totalMint = tokenMint[i] * delegatee.length;
    await autonity.mint(delegators[i], totalMint, {from: operator});
    for (let j = 0; j < delegatee.length; j++) {
      await autonity.bond(delegatee[j], tokenMint[i], {from: delegators[i]});
      bondingCount++;
    }
  }
  return bondingCount;

}

async function bulkUnbondingRequest(autonity, delegators, delegatee, tokenUnbond) {
  let unbondingCount = 0;
  for (let i = 0; i < delegators.length; i++) {
    for (let j = 0; j < delegatee.length; j++) {
      await autonity.unbond(delegatee[j], tokenUnbond[i], {from: delegators[i]});
      unbondingCount++;
    }
  }
  return unbondingCount;
}

async function mineTillUnbondingRelease(autonity, operator, deployer, maybeReleasedAlready = true) {
  let requestID = (await autonity.getHeadUnbondingID()).toNumber() - 1;
  let request = await autonity.getUnbondingRequest(requestID);
  let currentUnbondingPeriod = (await autonity.getUnbondingPeriod()).toNumber();
  let unbondingReleaseHeight = Number(request.requestBlock) + currentUnbondingPeriod;
  let lastEpochBlock = (await autonity.getLastEpochBlock()).toNumber();
  if (!maybeReleasedAlready) {
    // the following needs to be true in case unbonding not released already:
    // UnbondingRequestBlock + UnbondingPeriod > LastEpochBlock
    assert(
      unbondingReleaseHeight > lastEpochBlock,
      `unbonding period too short for testing, request-block: ${Number(request.requestBlock)}, unbonding-period: ${currentUnbondingPeriod}, `
      + `last-epoch-block: ${lastEpochBlock}`
    );
  }
  // mine blocks until unbonding period is reached
  while (await web3.eth.getBlockNumber() < unbondingReleaseHeight) {
    await mineEmptyBlock();
  }
}

// nodejs sleep
function timeout(ms) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

// set solidity bytecode at arbitrary address address
async function setCode(addr, code, contractName) {
  return new Promise((resolve, reject) => {
    web3.currentProvider.send({
      method: "evm_setAccountCode",
      params: [addr, code]
    }, (err, res) => {
      if (res?.result) { resolve(`\tSuccessfully mocked ${contractName} precompile.`); }
      else { reject(`\tError while mocking ${contractName} precompile.`); }
    });
  });
}

async function mockPrecompile() {
  await mockEnodePrecompile();
  await mockCommitteeSelectorPrecompile();
}

// mock verify enode precompiled contract
async function mockEnodePrecompile() {
      console.log("\tAttempting to mock enode verifier precompile. Will (rightfully) fail if running against Autonity network")
      const instance = await mockEnodeVerifier.new();
      console.log("enode verifier mocker address: ", instance.address)
      const code = await web3.eth.getCode(instance.address);
      const verifyEnodeAddr = "0x00000000000000000000000000000000000000ff";
      await setCode(verifyEnodeAddr, code, "enode verifier").then(
        (result) => {
            console.log(result); 
        },
        (error) => {
            console.log(error); 
    });
}

// mock committee selector precompiled contract
async function mockCommitteeSelectorPrecompile() {
  console.log("\tAttempting to mock committee selector precompile. Will (rightfully) fail if running against Autonity network")
  const instance = await mockCommitteeSelector.new();
  console.log("committee selector mocker address: ", instance.address)
  const code = await web3.eth.getCode(instance.address);
  const contractAddress = "0x00000000000000000000000000000000000000fa";
  await setCode(contractAddress, code, "committee selector").then(
    (result) => {
        console.log(result); 
    },
    (error) => {
        console.log(error); 
    });
}

// mine an empty block.
// If we are on an autonity network the rpc request will fail.
// In that case we just wait for an empty block to be mined
async function mineEmptyBlock() {
  let height = await web3.eth.getBlockNumber()
  let evmMineSuccess;
  await _mineEmptyBlock().then(
    (result) => {
      evmMineSuccess = true
    },
    (error) => {
      evmMineSuccess = false
    })
  if(!evmMineSuccess){
    await waitForNewBlock(height)
  }
}

async function waitForNewBlock(height){
  for(;;){
    let newHeight = await web3.eth.getBlockNumber()
    if (newHeight > height){
      break
    }
    timeout(50)
  }
}

// request ganache to mine empty block
async function _mineEmptyBlock() {
  return new Promise((resolve, reject) => {
    web3.currentProvider.send({
      method: "evm_mine",
    }, (err, res) => {
      if (res?.result) { resolve(); }
      else { 
        reject();
      }
    });
  });
}

const createAutonityContract = async (validators, autonityConfig, deployer) => {
    return Autonity.new(validators, autonityConfig, deployer);
}

const createAutonityTestContract = async (validators, autonityConfig, deployer) => {
  return AutonityTest.new(validators, autonityConfig, deployer);
}

async function initialize(autonity, autonityConfig, validators, accountabilityConfig, deployer, operator) {
  await autonity.finalizeInitialization({from: deployer});

  // accountability contract
  const accountability = await Accountability.new(autonity.address, accountabilityConfig, {from: deployer});
  
  // oracle contract
  let voters = validators.map((item, index) => (item.oracleAddress));
  const oracle = await Oracle.new(voters, autonity.address, operator, [], 30, {from: deployer});

  // acu contract (temporary empty basket and scale = 2)
  const acu = await Acu.new([], [], 2, autonity.address, operator, oracle.address, {from: deployer});
  
  // supply control contract. we will set the stabilizer address later
  const supplyControl = await SupplyControl.new(autonity.address,operator,"0x0000000000000000000000000000000000000000",{from:deployer,value:1})

  // stabilization contract, random temporary config and zeroAddress as collateral token

  const stabilization = await Stabilization.new(config.STABILIZATION_CONFIG,autonity.address,operator,oracle.address,supplyControl.address,"0x0000000000000000000000000000000000000000",{from:deployer})
  const upgradeManager = await UpgradeManager.new(autonity.address,operator,{from:deployer})

  await supplyControl.setStabilizer(stabilization.address,{from:operator});
  
  // non stakable contract
  const nonStakableVesting = await NonStakableVesting.new(autonity.address, operator, {from: deployer})
  
  await autonity.setAccountabilityContract(accountability.address, {from:operator});
  await autonity.setAcuContract(acu.address, {from: operator});
  await autonity.setSupplyControlContract(acu.address, {from: operator});
  await autonity.setStabilizationContract(acu.address, {from: operator});
  await autonity.setOracleContract(oracle.address, {from:operator});
  await autonity.setUpgradeManagerContract(upgradeManager.address, {from:operator});
  await autonity.setNonStakableVestingContract(nonStakableVesting.address, {from: operator})
}

// deploys protocol contracts
// set shortenEpoch = false if no need to call utils.endEpoch
const deployContracts = async (validators, autonityConfig, accountabilityConfig, deployer, operator, shortenEpoch = true) => {
    // we deploy first the inflation controller contract because it requires a genesis timestamp
    // greater than the one of the autonity contract. This is obviously not going to happen for a real network but
    // we can't really simulate a proper genesis sequence with truffle. As consequence all calculations
    // regarding the inflation rate will be wrong here which should be tested using the native go framework.
    const inflationController = await InflationController.new(config.INFLATION_CONTROLLER_CONFIG ,{from:deployer})
    // autonity contract
    // As the chain height might exceed the lastEpochBlock(0)+EpochPeriod(30) of the newly deployed AC with a lots of
    // blocks, it makes the AC impossible to finalize an epoch, thus we resolved a correct boundary of the 1st epoch
    // at this point, to make the deployed contract have a chance to finalize the 1st epoch, and then apply the default
    // 30 blocks epoch period for the testing.
    let currentHeight = await web3.eth.getBlockNumber();
    let firstEpochEndBlock = currentHeight+10;
    let copyAutonityConfig = autonityConfig;
    copyAutonityConfig.protocol.epochPeriod = firstEpochEndBlock;

    const autonity = await createAutonityContract(validators, copyAutonityConfig, {from: deployer});
    // now apply the correct epoch period, and wait it to be applied at the end of the 1st epoch finalization.
    await autonity.setEpochPeriod(autonityConfig.protocol.epochPeriod, {from: operator})

    // wait for the firstEpochEndBlock, and try to finalize it until the epoch rotation happens.
    for (let i=currentHeight;i<=firstEpochEndBlock;i++) {
      let height = await web3.eth.getBlockNumber()
      console.log("try to finalize 1st epoch after AC deployment", "height: ", height, "next epoch block: ", firstEpochEndBlock);
      autonity.finalize({from: deployer})
      let epochID = (await autonity.epochID()).toNumber()
      // if the epoch rotates from 0 to 1, then we have done the setup.
      if (epochID === 1) {
        break;
      }
      await waitForNewBlock(height);
    }

    await autonity.setInflationControllerContract(inflationController.address, {from:operator});
    await initialize(autonity, autonityConfig, validators, accountabilityConfig, deployer, operator);

    return autonity;
};

// deploys AutonityTest, a contract inheriting Autonity and exposing the "_applyNewCommissionRates" function
// set shortenEpoch = false if no need to call utils.endEpoch
const deployAutonityTestContract = async (validators, autonityConfig, accountabilityConfig, deployer, operator, shortenEpoch = true) => {
    const inflationController = await InflationController.new(config.INFLATION_CONTROLLER_CONFIG,{from:deployer})

    // As the chain height might exceed the lastEpochBlock(0)+EpochPeriod(30) of the newly deployed AC with a lots of
    // blocks, it makes the AC impossible to finalize an epoch, thus we resolved a correct boundary of the 1st epoch
    // at this point, to make the deployed contract have a chance to finalize the 1st epoch, and then apply the default
    // 30 blocks epoch period for the testing.
    let currentHeight = await web3.eth.getBlockNumber();
    let firstEpochEndBlock = currentHeight+10;
    let copyAutonityConfig = autonityConfig;
    copyAutonityConfig.protocol.epochPeriod = firstEpochEndBlock;

    const autonityTest = await createAutonityTestContract(validators, copyAutonityConfig, {from: deployer});
    // now apply the correct epoch period, and wait it to be applied at the end of the 1st epoch finalization.
    await autonityTest.setEpochPeriod(autonityConfig.protocol.epochPeriod, {from: operator})

    // wait for the firstEpochEndBlock, and try to finalize it until the epoch rotation happens.
    for (let i=currentHeight;i<=firstEpochEndBlock;i++) {
      let height = await web3.eth.getBlockNumber()
      console.log("try to finalize 1st epoch after testAC deployment", "height: ", height, "next epoch block: ", firstEpochEndBlock);
      autonityTest.finalize({from: deployer})
      let epochID = (await autonityTest.epochID()).toNumber()
      // if the epoch rotates from 0 to 1, then we have done the setup.
      if (epochID === 1) {
        break;
      }
      await waitForNewBlock(height);
    }

    await autonityTest.setInflationControllerContract(inflationController.address, {from:operator});
    await initialize(autonityTest, autonityConfig, validators, accountabilityConfig, deployer, operator);
    return autonityTest;
};

function ruleToRate(accountabilityConfig,rule){
  //TODO(lorenzo) create mapping rule to rate once finalized in autonity.sol. bypass severity conversion?
  return accountabilityConfig.baseSlashingRateMid
}

async function signTransaction(from, to, privateKey, methodRequest = null) {
  let data = "0x";
  let gasLimit = 1000000000;
  if (methodRequest != null) {
    data = methodRequest.data;
    gasLimit = methodRequest.gas;
  }
  let tx = {
    from: from,
    to: to,
    gas: gasLimit,
    data: data
  }
  return await web3.eth.accounts.signTransaction(tx, privateKey);
}

async function signAndSendTransaction(from, to, privateKey, methodRequest = null) {
  let signedTx = await signTransaction(from, to, privateKey, methodRequest);
  return await web3.eth.sendSignedTransaction(signedTx.rawTransaction);
}

function bytesToHex(bytes) {
  let hex = "0x";
  for (let i = 0; i < bytes.length; i++) {
    hex += (bytes[i] > 15) ? bytes[i].toString(16) : "0" + bytes[i].toString(16);
  }
  return hex;
}

function randomInt() {
  const MAX = 1e10;
  return Math.floor(Math.random() * MAX);
}

function randomPrivateKey() {
  let key = [];
  for (let i = 0; i < 32; i++) {
    key.push(randomInt() % 256);
  }
  return bytesToHex(key).substring(2);
}

function privateKeyToEnode(privateKey) {
  let key = publicKey(privateKey);
  key = key.substring(key.length - 128);
  return publicKeyToEnode(key);
}

function publicKeyToEnode(publicKey) {
  return "enode://" + publicKey + "@3.209.45.79:30303";
}

function publicKeyObject(privateKey) {
  return ec.keyFromPrivate(privateKey).getPublic();
}

function publicKeyCompressed(privateKey, hex = true) {
  let publicKey = publicKeyObject(privateKey);
  return (hex == true) ? publicKey.encodeCompressed("hex") : new Uint8Array(publicKey.encodeCompressed());
}

function publicKey(privateKey, hex = true) {
  let publicKey = publicKeyObject(privateKey);
  return (hex == true) ? publicKey.encode("hex") : new Uint8Array(publicKey.encode());
}

function address(publicKeyUncompressedBytes) {
  return ethers.utils.getAddress("0x" + keccakHash(publicKeyUncompressedBytes.subarray(1)).substring(24));
}

function generateMultiSig(nodekey, oraclekey, treasuryAddr) {
  let treasuryProof = web3.eth.accounts.sign(treasuryAddr, nodekey);
  let oracleProof = web3.eth.accounts.sign(treasuryAddr, oraclekey);
  let multisig = treasuryProof.signature + oracleProof.signature.substring(2)
  return multisig
}

async function generateAutonityPOP(autonityKeysFile, oracleKeyHex, treasuryAddress) {
  const command = `../../../build/bin/autonity genOwnershipProof --autonitykeys ${autonityKeysFile} --oraclekeyhex ${oracleKeyHex} ${treasuryAddress}`;
  try {
    const { stdout, stderr } = await exec(command);
    if (stderr) {
      throw new Error(stderr);
    }
    const outputLines = stdout.split('\n');
    const signatures = outputLines[0].trim();
    return { signatures };
  } catch (error) {
    return { error: error.message };
  }
}

async function generateAutonityKeys(filePath) {
  try {
    const command = `../../../build/bin/autonity genAutonityKeys --writeaddress ${filePath}`;
    const { stdout, stderr } = await exec(command);
    if (stderr) {
      throw new Error(stderr);
    }
    const nodeAddress = stdout.match(/Node address: (0x[0-9a-fA-F]+)/)[1];
    const nodePublicKey = stdout.match(/Node public key: (0x[0-9a-fA-F]+)/)[1];
    const nodeConsensusKey = stdout.match(/Consensus public key: (0x[0-9a-fA-F]+)/)[1];
    return { nodeAddress, nodePublicKey, nodeConsensusKey };
  } catch (error) {
    throw new Error(`Failed to execute command: ${error.message}`);
  }
}

function keccakHash(input) {
  return keccak256(Buffer.from(input)).toString('hex');
}

async function slash(config, accountability, epochOffenceCount, offender, reporter) {
  const event = {
    "chunks": 1,
    "chunkId": 1,
    "eventType": 0,
    "rule": 0, // PN rule --> severity mid
    "reporter": reporter,
    "offender": offender,
    "rawProof": [],
    "id": 0,
    "block": 1,
    "epoch": 0,
    "reportingBlock": 2,
    "messageHash": 0,
  }
  let tx = await accountability.slash(event, epochOffenceCount);
  let txEvent;
  truffleAssert.eventEmitted(tx, 'SlashingEvent', (ev) => {
    txEvent = ev;
    return ev.amount.toNumber() > 0;
  });
  let slashingRate = ruleToRate(config, event.rule) / config.slashingRatePrecision;
  return {txEvent, slashingRate};
}

module.exports.deployContracts = deployContracts;
module.exports.deployAutonityTestContract = deployAutonityTestContract;
module.exports.mineEmptyBlock = mineEmptyBlock;
module.exports.setCode = setCode;
module.exports.mockPrecompile = mockPrecompile;
module.exports.mockCommitteeSelectorPrecompile = mockCommitteeSelectorPrecompile;
module.exports.timeout = timeout;
module.exports.waitForNewBlock = waitForNewBlock;
module.exports.endEpoch = endEpoch;
module.exports.validatorState = validatorState;
module.exports.bulkBondingRequest = bulkBondingRequest;
module.exports.bulkUnbondingRequest = bulkUnbondingRequest;
module.exports.mineTillUnbondingRelease = mineTillUnbondingRelease;
module.exports.ruleToRate = ruleToRate;
module.exports.signTransaction = signTransaction;
module.exports.signAndSendTransaction = signAndSendTransaction;
module.exports.bytesToHex = bytesToHex;
module.exports.randomPrivateKey = randomPrivateKey;
module.exports.generateMultiSig = generateMultiSig;
module.exports.ValidatorState = ValidatorState;
module.exports.generateAutonityPOP = generateAutonityPOP;
module.exports.generateAutonityKeys = generateAutonityKeys;
module.exports.publicKeyToEnode = publicKeyToEnode;
module.exports.privateKeyToEnode = privateKeyToEnode;
module.exports.publicKeyCompressed = publicKeyCompressed;
module.exports.publicKey = publicKey;
module.exports.address = address;
module.exports.slash = slash;
