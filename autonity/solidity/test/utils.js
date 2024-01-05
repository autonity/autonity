const assert = require('assert');
const util = require('util');
const exec = util.promisify(require('child_process').exec);
const Autonity = artifacts.require("Autonity");
const Accountability = artifacts.require("Accountability");
const Oracle = artifacts.require("Oracle")
const Acu = artifacts.require("ACU")
const SupplyControl = artifacts.require("SupplyControl")
const Stabilization = artifacts.require("Stabilization")
const AutonityTest = artifacts.require("AutonityTest");
const mockEnodeVerifier = artifacts.require("MockEnodeVerifier")
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


// end epoch so the LastEpochBlock is closer
// then set epoch period 
async function shortenEpochPeriod(autonity, epochPeriod, operator, deployer) {
  await endEpoch(autonity, operator, deployer);
  await autonity.setEpochPeriod(epochPeriod, {from: operator});
}

// while testing we might ran into situations were currentHeight > lastEpochBlock + epochPeriod
// in this case in order to be able to finalize we need to setEpochPeriod to a bigger value
// also we need to take into account that if we are running against autonity, the network will keep mining as we do these operations
async function endEpoch(contract,operator,deployer){
  let lastEpochBlock = (await contract.getLastEpochBlock()).toNumber();
  let currentHeight = await web3.eth.getBlockNumber();
  let currentEpoch = (await contract.epochID()).toNumber()
  let delta = currentHeight - lastEpochBlock
  let epochPeriod = delta + 5

  await contract.setEpochPeriod(epochPeriod,{from: operator})

  assert.equal(epochPeriod,(await contract.getEpochPeriod()).toNumber())

  // close epoch
  for (let i=0;i<(lastEpochBlock + epochPeriod) - currentHeight;i++) {
    let height = await web3.eth.getBlockNumber()
    contract.finalize({from: deployer})
    await waitForNewBlock(height);
  }
  let newEpoch = (await contract.epochID()).toNumber()
  assert.equal(currentEpoch+1,newEpoch)
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
async function setCode(addr, code) {
  return new Promise((resolve, reject) => {
    web3.currentProvider.send({
      method: "evm_setAccountCode",
      params: [addr, code]
    }, (err, res) => {
      if (res?.result) { resolve("\tSuccessfully mocked verifier precompile."); }
      else { reject("\tError while mocking verifier precompile."); }
    });
  });
}

// mock verify enode precompiled contract
async function mockEnodePrecompile() {
      const instance = await mockEnodeVerifier.new();
      console.log("enode verifier mocker address: ", instance.address)
      const code = await web3.eth.getCode(instance.address);
      const verifyEnodeAddr = "0x00000000000000000000000000000000000000ff";
      await setCode(verifyEnodeAddr, code).then(
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
    timeout(100)
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
  config = { 
    "borrowInterestRate" : 0,
    "liquidationRatio" : 1,
    "minCollateralizationRatio" : 2,
    "minDebtRequirement" : 0,
    "targetPrice" : 0,
  }
  const stabilization = await Stabilization.new(config,autonity.address,operator,oracle.address,supplyControl.address,"0x0000000000000000000000000000000000000000",{from:deployer})

  // setters
  await supplyControl.setStabilizer(stabilization.address,{from:operator});
  
  await autonity.setAccountabilityContract(accountability.address, {from:operator});
  await autonity.setAcuContract(acu.address, {from: operator});
  await autonity.setSupplyControlContract(acu.address, {from: operator});
  await autonity.setStabilizationContract(acu.address, {from: operator});
  await autonity.setOracleContract(oracle.address, {from:operator});
  await shortenEpochPeriod(autonity, autonityConfig.protocol.epochPeriod, operator, deployer);
}

// deploys protocol contracts
const deployContracts = async (validators, autonityConfig, accountabilityConfig, deployer, operator) => {
    // autonity contract
    const autonity = await createAutonityContract(validators, autonityConfig, {from: deployer});
    await initialize(autonity, autonityConfig, validators, accountabilityConfig, deployer, operator);
    return autonity;
};

// deploys AutonityTest, a contract inheriting Autonity and exposing the "_applyNewCommissionRates" function
const deployAutonityTestContract = async (validators, autonityConfig, accountabilityConfig, deployer, operator) => {
    const autonityTest = await createAutonityTestContract(validators, autonityConfig, {from: deployer});
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

function publicKeyToEnode(publicKey) {
  return "enode://" + publicKey + "@3.209.45.79:30303";
}

function publicKeyObject(privateKey) {
  return ec.keyFromPrivate(privateKey).getPublic();
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

async function generateAutonityPOP(nodeKeyFile, oracleKeyHex, treasuryAddress) {
  const command = `../../../build/bin/autonity genOwnershipProof --nodekey ${nodeKeyFile} --oraclekeyhex ${oracleKeyHex} ${treasuryAddress}`;
  try {
    const { stdout, stderr } = await exec(command);
    if (stderr) {
      throw new Error(stderr);
    }
    const outputLines = stdout.split('\n');
    const nodeConsensusKey = outputLines.find(line => line.startsWith('Node consensus key:')).split(':')[1].trim();
    const signatures = outputLines.find(line => line.startsWith('Signatures:')).split(':')[1].trim();
    return { nodeConsensusKey, signatures };
  } catch (error) {
    return { error: error.message };
  }
}

async function generateNodeKey(filePath) {
  try {
    const command = `../../../build/bin/autonity genNodeKey --writeaddress ${filePath}`;
    const { stdout, stderr } = await exec(command);
    const nodeAddress = stdout.match(/Node address: (0x[0-9a-fA-F]+)/)[1];
    const nodePublicKey = stdout.match(/Node public key: (0x[0-9a-fA-F]+)/)[1];
    const nodeConsensusKey = stdout.match(/Node consensus key: (0x[0-9a-fA-F]+)/)[1];
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
module.exports.mockEnodePrecompile = mockEnodePrecompile;
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
module.exports.generateNodeKey = generateNodeKey;
module.exports.publicKeyToEnode = publicKeyToEnode;
module.exports.publicKey = publicKey;
module.exports.address = address;
module.exports.slash = slash;
