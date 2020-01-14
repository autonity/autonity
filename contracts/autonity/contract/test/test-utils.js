const Autonity = artifacts.require("Autonity.sol");


const deployContract = async (accounts, enodes, userTypes, stakes, sysOperator, minGasPrice, bondPeriod, committeeSize,msgSender) => {

    return Autonity.new(accounts, enodes, userTypes, stakes, sysOperator, minGasPrice, bondPeriod, committeeSize, msgSender);
};


module.exports.deployContract= deployContract;
