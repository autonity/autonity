const Autonity = artifacts.require("Autonity");
const Accountability = artifacts.require("Accountability");
const Oracle = artifacts.require("Oracle")
const Acu = artifacts.require("ACU")
const SupplyControl = artifacts.require("SupplyControl")
const Stabilization = artifacts.require("Stabilization")
const AutonityTest = artifacts.require("AutonityTest");
const mockEnodeVerifier = artifacts.require("MockEnodeVerifier")

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
      if (res?.result) { resolve("\tSuccessfully mocked enode verifier precompile."); }
      else { reject("\tError while mocking enode verifier precompile."); }
    });
  });
}

// mock verify enode precompiled contract
async function mockEnodePrecompile() {
      const instance = await mockEnodeVerifier.new();
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

const createAutonityTestContract = async (validators, autonityConfig, unbondingPeriod, deployer) => {
  return AutonityTest.new(validators, autonityConfig, unbondingPeriod, deployer);
}

async function initialize(autonity, validators, accountabilityConfig, deployer, operator) {
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
}

// deploys protocol contracts
const deployContracts = async (validators, autonityConfig, accountabilityConfig, deployer, operator) => {
    // autonity contract
    const autonity = await createAutonityContract(validators, autonityConfig, {from: deployer});
    await initialize(autonity, validators, accountabilityConfig, deployer, operator);
    return autonity;
};

// deploys AutonityTest, a contract inheriting Autonity and exposing the "_applyNewCommissionRates" function
const deployAutonityTestContract = async (validators, autonityConfig, accountabilityConfig, deployer, operator, unbondingPeriod = 0) => {
    const autonityTest = await createAutonityTestContract(validators, autonityConfig, unbondingPeriod, {from: deployer});
    await initialize(autonityTest, validators, accountabilityConfig, deployer, operator);
    return autonityTest;
};


module.exports.deployContracts = deployContracts;
module.exports.deployAutonityTestContract = deployAutonityTestContract;
module.exports.mineEmptyBlock = mineEmptyBlock;
module.exports.setCode = setCode;
module.exports.mockEnodePrecompile = mockEnodePrecompile;
module.exports.timeout = timeout;
module.exports.waitForNewBlock = waitForNewBlock;
