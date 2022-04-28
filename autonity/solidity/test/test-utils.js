const Autonity = artifacts.require("Autonity.sol");


const deployContract = async (validators, config, msgSender) => {
    return Autonity.new(validators, config, msgSender);
};


module.exports.deployContract = deployContract;
