const Autonity = artifacts.require("Autonity.sol");


const deployContract = async (accounts, enodes, userTypes, stakes, commissionRates, sysOperator, minGasPrice, bondPeriod, committeeSize, version, msgSender) => {
    return Autonity.new(accounts, enodes, userTypes, stakes, commissionRates, sysOperator, minGasPrice, bondPeriod, committeeSize, version, msgSender);
};


module.exports.deployContract = deployContract;
