const Autonity = artifacts.require("Autonity.sol");
const Accountability = artifacts.require("Accountability.sol");
const AutonityTest = artifacts.require("AutonityTest");


const deployContract = async (validators, config, msgSender) => {
    const autonity = await Autonity.new(validators, config, msgSender);
    const accountability = await Accountability.new(autonity.address);
    await autonity.setAccountabilityContract(accountability.address, msgSender);
    return autonity
};

const deployTestContract = async (validators, config, msgSender) => {
    return AutonityTest.new(validators, config, msgSender);
};


module.exports.deployContract = deployContract;
module.exports.deployTestContract = deployTestContract;