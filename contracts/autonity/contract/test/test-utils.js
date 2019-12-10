const Autonity = artifacts.require("Autonity.sol");


const deployContract = async (accounts, enodes, userTypes, stakes, sysOperator, minGasPrice, msgSender) => {

    return Autonity.new(accounts, enodes, userTypes, stakes, sysOperator, minGasPrice, msgSender);
};


module.exports.deployContract= deployContract;
