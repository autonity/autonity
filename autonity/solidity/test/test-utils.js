const Autonity = artifacts.require("Autonity.sol");


const deployContract = async (accounts, enodes, userTypes, stakes, sysOperator, minGasPrice, committeeSize, version, msgSender) => {
    return Autonity.new(accounts, enodes, userTypes, stakes, sysOperator, minGasPrice, committeeSize, version, msgSender);
};


module.exports.deployContract = deployContract;
