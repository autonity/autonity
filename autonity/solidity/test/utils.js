const assert = require('assert');
const util = require('util');
const config = require('./config');
const exec = util.promisify(require('child_process').exec);
const Autonity = artifacts.require("Autonity");
const Accountability = artifacts.require("Accountability");
const OmissionAccountability = artifacts.require("OmissionAccountability");
const UpgradeManager = artifacts.require("UpgradeManager");
const Oracle = artifacts.require("Oracle")
const Acu = artifacts.require("ACU")
const SupplyControl = artifacts.require("SupplyControl")
const Stabilization = artifacts.require("StabilizationMock")
const InflationController = artifacts.require("InflationController")
const AuctioneerMock = artifacts.require("AuctioneerMock")
const AutonityTest = artifacts.require("AutonityTest");
const mockEnodeVerifier = artifacts.require("MockEnodeVerifier")
const mockCommitteeSelector = artifacts.require("MockCommitteeSelector")
const EC = require('elliptic').ec;
const ec = new EC('secp256k1');
const keccak256 = require('keccak256');
const ethers = require('ethers');
const truffleAssert = require('truffle-assertions');
const {SLASHING_RATE_PRECISION} = require("./config");
const path = require("path");
const fs = require("fs");

// Resolve the autonity binary relative to this test directory so the tests work
// no matter what the process cwd is (e.g. when running `truffle test` from
// `autonity/solidity`).
const AUTONITY_BIN = path.resolve(__dirname, "../../../build/bin/autonity");

let _autonityProviderPatched = false;
let _isAutonityNetworkCached;
let _cachedGasPrice;
let _cachedGasPriceAtMs = 0;

async function getCachedGasPrice() {
  const now = Date.now();
  // Keep this short to avoid stale pricing if basefee changes.
  if (_cachedGasPrice != null && (now - _cachedGasPriceAtMs) < 2000) {
    return _cachedGasPrice;
  }
  try {
    _cachedGasPrice = await web3.eth.getGasPrice();
    _cachedGasPriceAtMs = now;
    return _cachedGasPrice;
  } catch (_) {
    return undefined;
  }
}

function _isTxIndexingInProgressError(err) {
  if (!err) return false;
  const msg = (err.message || "").toString().toLowerCase();
  if (msg.includes("transaction indexing is in progress")) return true;
  // Some providers attach the JSON-RPC error under `data.originalError`.
  const orig = err?.data?.originalError;
  if (orig && typeof orig.message === "string") {
    if (orig.message.toLowerCase().includes("transaction indexing is in progress")) return true;
  }
  return false;
}

function installAutonityProviderWorkarounds() {
  if (_autonityProviderPatched) return;
  _autonityProviderPatched = true;

  // Truffle/Web3 treats JSON-RPC errors as hard failures during receipt polling.
  // Autonity can legitimately return "transaction indexing is in progress" very
  // early after startup; for polling callers it should behave like "not yet".
  const patchSend = (provider, fnName) => {
    if (!provider || typeof provider[fnName] !== "function") return;
    const orig = provider[fnName].bind(provider);
    provider[fnName] = (payload, cb) => {
      // Preserve default behavior if no callback is provided.
      if (typeof cb !== "function") return orig(payload, cb);

      return orig(payload, (err, res) => {
        const method = payload?.method;
        const isReceiptish =
          method === "eth_getTransactionReceipt" ||
          method === "eth_getTransactionByHash";

        // Handle both "err" and "res.error" shapes across providers.
        if (isReceiptish && _isTxIndexingInProgressError(err)) {
          return cb(null, { jsonrpc: "2.0", id: payload.id, result: null });
        }
        if (isReceiptish && res && res.error) {
          const emsg = (res.error.message || "").toString().toLowerCase();
          if (emsg.includes("transaction indexing is in progress")) {
            return cb(null, { jsonrpc: "2.0", id: payload.id, result: null });
          }
        }
        return cb(err, res);
      });
    };
  };

  patchSend(web3.currentProvider, "send");
  patchSend(web3.currentProvider, "sendAsync");
}

async function isAutonityNetwork() {
  if (_isAutonityNetworkCached !== undefined) return _isAutonityNetworkCached;
  try {
    const nodeInfo = await web3.eth.getNodeInfo();
    _isAutonityNetworkCached = typeof nodeInfo === "string" && nodeInfo.toLowerCase().includes("autonity/");
  } catch (_) {
    _isAutonityNetworkCached = false;
  }
  return _isAutonityNetworkCached;
}

async function failsRevert(promise, reason) {
  // On Autonity nodes, Truffle frequently surfaces reverts as a generic
  // StatusError (status 0) without the revert reason string. Keep reason checks
  // where possible, but fall back to "reverted" when the provider doesn't carry
  // the string.
  if (await isAutonityNetwork()) {
    return truffleAssert.fails(promise, truffleAssert.ErrorType.REVERT);
  }
  if (reason != null) {
    return truffleAssert.fails(promise, truffleAssert.ErrorType.REVERT, reason);
  }
  return truffleAssert.fails(promise, truffleAssert.ErrorType.REVERT);
}

// Validator Status in Autonity Contract
const ValidatorState = {
  active : 0,
  paused : 1,
  jailed : 2,
  jailbound : 3
}

async function endEpoch(contract,operator,deployer){
  const epochPeriod = (await contract.getEpochPeriod()).toNumber();
  const currentEpoch = (await contract.getEpochID()).toNumber();
  const targetEpoch = currentEpoch + 1;

  const isAutonity = await isAutonityNetwork();

  // Autonity: `finalize()` is what applies the per-block protocol transitions in
  // these tests (the contract is a deployed test instance, not the system one).
  // Ganache: same, but also waits for the new block if needed.
  for (let i = 0; i <= epochPeriod; i++) {
    await contract.finalize({from: deployer});
    const newEpochID = (await contract.getEpochID()).toNumber();
    if (newEpochID === targetEpoch) return;
    if (!isAutonity) {
      const height = await web3.eth.getBlockNumber();
      await waitForNewBlock(height);
    }
  }

  throw new Error(`endEpoch timeout: epoch did not advance from ${currentEpoch} after ${epochPeriod + 1} finalize() calls`);
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
  // These mocks are only needed for Ganache. On an Autonity node the precompiles
  // already exist and `evm_setAccountCode` is not available.
  try {
    const nodeInfo = await web3.eth.getNodeInfo();
    if (typeof nodeInfo === "string" && nodeInfo.toLowerCase().includes("autonity/")) {
      _isAutonityNetworkCached = true;
      installAutonityProviderWorkarounds();
      console.log("\tSkipping precompile mocking on Autonity network")
      return
    }
  } catch (_) {
    // If we can't detect the node, fall through and attempt mocking.
  }

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
// In that case we need to trigger a block by sending a minimal tx, because
// the local tendermint testnet config does not produce empty blocks.
async function mineEmptyBlock() {
  const height = await web3.eth.getBlockNumber()

  // Autonity nodes don't support `evm_mine`; avoid spamming failing RPC calls.
  if (await isAutonityNetwork()) {
    const accounts = await web3.eth.getAccounts()
    if (!accounts || accounts.length === 0) {
      throw new Error("mineEmptyBlock: no accounts available to trigger a block when evm_mine is unsupported")
    }
    const from = accounts[0]

    // Ensure the tx is accepted on post-London chains.
    const gasPrice = await getCachedGasPrice()
    const tx = { from, to: from, value: "0x0", gas: 21000 }
    if (gasPrice != null) {
      tx.gasPrice = gasPrice
    }
    const res = await web3.eth.sendTransaction(tx)
    // HDWalletProvider/web3 typically waits for mining and returns a receipt.
    // If we only got a tx hash, fall back to height polling.
    if (!(res && typeof res === "object" && res.blockNumber != null)) {
      await waitForNewBlock(height)
    }
    return
  }

  let evmMineSuccess = true
  await _mineEmptyBlock().catch(() => {
    evmMineSuccess = false
  })

  if(!evmMineSuccess){
    // If the chain is already producing blocks (e.g. autonity --mine), don't spam
    // the node with transactions, just wait for the next block.
    try {
      // 1.5s covers typical CI latency and the default 1s block period.
      const start = Date.now()
      while (Date.now() - start < 1500) {
        const newHeight = await web3.eth.getBlockNumber()
        if (newHeight > height) {
          return
        }
        await timeout(50)
      }
    } catch (_) {
      // fall through to tx-triggered block
    }

    const accounts = await web3.eth.getAccounts()
    if (!accounts || accounts.length === 0) {
      throw new Error("mineEmptyBlock: no accounts available to trigger a block when evm_mine is unsupported")
    }
    const from = accounts[0]

    // Ensure the tx is accepted on post-London chains.
    const gasPrice = await getCachedGasPrice()
    const tx = { from, to: from, value: "0x0", gas: 21000 }
    if (gasPrice != null) {
      tx.gasPrice = gasPrice
    }
    const res = await web3.eth.sendTransaction(tx)
    if (!(res && typeof res === "object" && res.blockNumber != null)) {
      await waitForNewBlock(height)
    }
  }
}

async function waitForNewBlock(height){
  for(;;){
    let newHeight = await web3.eth.getBlockNumber()
    if (newHeight > height){
      break
    }
    await timeout(100)
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

async function initialize(autonity, autonityConfig, validators, accountabilityConfig, omissionAccountabilityConfig, deployer, operator) {
  await autonity.finalizeInitialization(omissionAccountabilityConfig.delta,{from:deployer});

  // accountability contract
  const accountability = await Accountability.new(autonity.address, accountabilityConfig, {from: deployer});

  // oracle contract
  let voters = validators.map((item) => (item.oracleAddress));
  let treasuries = validators.map((item) => (item.treasury));
  let nodes = validators.map((item) => (item.nodeAddress));
  const oracle = await Oracle.new(voters, nodes, treasuries, [], {
    autonity: autonity.address,
    operator,
    votePeriod: 3,
    outlierDetectionThreshold: 100,
    outlierSlashingThreshold: 100,
    baseSlashingRate: 10,
    nonRevealThreshold: 3,
    revealResetInterval: 10,
    slashingRateCap: 1000,
  }, {from: deployer});

  // acu contract (temporary empty basket and scale = 2)
  const acu = await Acu.new([], [], 2, autonity.address, operator, oracle.address, {from: deployer});

  // supply control contract. we will set the stabilizer address later
  const supplyControl = await SupplyControl.new(autonity.address,operator,"0x0000000000000000000000000000000000000000",{from:deployer,value:1})
  const auctioneer = await AuctioneerMock.new({from:deployer});
  const stabilization = await Stabilization.new({from:deployer});
  const upgradeManager = await UpgradeManager.new(autonity.address,operator,{from:deployer})

  // omission accountability contract
  const omissionAccountability = await OmissionAccountability.new(autonity.address, operator, omissionAccountabilityConfig, {from:deployer})

  await autonity.setSupplyControlContract(supplyControl.address, {from: operator});
  await autonity.setAuctioneerContract(auctioneer.address, {from: operator});
  await autonity.setStabilizationContract(stabilization.address, {from: operator});
  await autonity.setAccountabilityContract(accountability.address, {from:operator});
  await autonity.setAcuContract(acu.address, {from: operator});
  await autonity.setOracleContract(oracle.address, {from:operator});
  await autonity.setUpgradeManagerContract(upgradeManager.address, {from:operator});
  await autonity.setOmissionAccountabilityContract(omissionAccountability.address, {from: operator});

}

// deploys protocol contracts
// set shortenEpoch = false if no need to call utils.endEpoch
const deployContracts = async (validators, autonityConfig, accountabilityConfig, omissionAccountabilityConfig, deployer, operator, shortenEpoch = true) => {
    // we deploy first the inflation controller contract because it requires a genesis timestamp
    // greater than the one of the autonity contract. This is obviously not going to happen for a real network but
    // we can't really simulate a proper genesis sequence with truffle. As consequence all calculations
    // regarding the inflation rate will be wrong here which should be tested using the native go framework.
    const inflationController = await InflationController.new(config.INFLATION_CONTROLLER_CONFIG ,{from:deployer})

    const autonity = await createAutonityContract(validators, autonityConfig, {from: deployer});

    // now init autonity contract with sub protocol contracts, otherwise finalize() will be reverted.
    await autonity.setInflationControllerContract(inflationController.address, {from:operator});
    await initialize(autonity, autonityConfig, validators, accountabilityConfig, omissionAccountabilityConfig, deployer, operator);
    return autonity;
};

// deploys AutonityTest, a contract inheriting Autonity and exposing the "_applyNewCommissionRates" function
// set shortenEpoch = false if no need to call utils.endEpoch
const deployAutonityTestContract = async (validators, autonityConfig, accountabilityConfig, omissionAccountabilityConfig, deployer, operator, shortenEpoch = true) => {
    const inflationController = await InflationController.new(config.INFLATION_CONTROLLER_CONFIG,{from:deployer})

    const autonityTest = await createAutonityTestContract(validators, autonityConfig, {from: deployer});

    // now init autonity contract with sub protocol contracts, otherwise finalize() will be reverted.
    await autonityTest.setInflationControllerContract(inflationController.address, {from:operator});
    await initialize(autonityTest, autonityConfig, validators, accountabilityConfig, omissionAccountabilityConfig, deployer, operator);
    return autonityTest;
};

function ruleToRate(accountabilityConfig,rule){
  if(rule == 9) { // equivocation
    return accountabilityConfig.baseSlashingRates.low
  }
  if(rule >= 0 && rule <= 6) {
    return accountabilityConfig.baseSlashingRates.mid
  }
  if(rule == 7 || rule == 8){ // invalid proposal and invalid proposer
    return accountabilityConfig.baseSlashingRates.high
  }
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
  const autonityKeysFileResolved = path.isAbsolute(autonityKeysFile)
    ? autonityKeysFile
    : path.resolve(__dirname, autonityKeysFile);
  const command = `${AUTONITY_BIN} genOwnershipProof --autonitykeys ${autonityKeysFileResolved} --oraclekeyhex ${oracleKeyHex} ${treasuryAddress}`;
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
    const filePathResolved = path.isAbsolute(filePath)
      ? filePath
      : path.resolve(__dirname, filePath);
    fs.mkdirSync(path.dirname(filePathResolved), { recursive: true });
    const command = `${AUTONITY_BIN} genAutonityKeys --writeaddress ${filePathResolved}`;
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

async function slash(config, accountability, epochOffenceCount, offender, reporter, epochPeriod) {
  const event = {
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
  let tx = await accountability.slash(event, epochOffenceCount, epochPeriod);
  let txEvent;
  truffleAssert.eventEmitted(tx, 'SlashingEvent', (ev) => {
    txEvent = ev;
    return ev.amount.toNumber() > 0;
  });
  let slashingRate = ruleToRate(config, event.rule) / SLASHING_RATE_PRECISION;
  return {txEvent, slashingRate};
}

module.exports.deployContracts = deployContracts;
module.exports.deployAutonityTestContract = deployAutonityTestContract;
module.exports.mineEmptyBlock = mineEmptyBlock;
module.exports.setCode = setCode;
module.exports.mockPrecompile = mockPrecompile;
module.exports.mockCommitteeSelectorPrecompile = mockCommitteeSelectorPrecompile;
module.exports.isAutonityNetwork = isAutonityNetwork;
module.exports.failsRevert = failsRevert;
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
