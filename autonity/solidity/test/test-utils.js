const Autonity = artifacts.require("Autonity.sol");
const Accountability = artifacts.require("Accountability.sol");
const AutonityTest = artifacts.require("AutonityTest");


const deployContract = async (validators, autonityConfig, accountabilityConfig, msgSender) => {
    const autonity = await Autonity.new(validators, autonityConfig, msgSender);
    const accountability = await Accountability.new(autonity.address, accountabilityConfig);
    await autonity.setAccountabilityContract(accountability.address, msgSender);
    await autonity.finalizeInitialization(msgSender)
    return autonity
};

const deployTestContract = async (validators, config, msgSender) => {
    return AutonityTest.new(validators, config, msgSender);
};


module.exports.deployContract = deployContract;
module.exports.deployTestContract = deployTestContract;