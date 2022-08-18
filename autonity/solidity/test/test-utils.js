const Autonity = artifacts.require("Autonity.sol");
const AutonityTest = artifacts.require("AutonityTest");


const deployContract = async (validators, config, msgSender) => {
    return Autonity.new(validators, config, msgSender);
};

const deployTestContract = async (validators, config, msgSender) => {
    return AutonityTest.new(validators, config, msgSender);
};


module.exports.deployContract = deployContract;
module.exports.deployTestContract = deployTestContract;