const Migrations = artifacts.require("Migrations");

const Web3 = require('web3');
const TruffleConfig = require('../truffle-config');

module.exports = function(deployer, network, accounts) {
  console.log('>> Deploying migration');
  deployer.deploy(Migrations);
};
