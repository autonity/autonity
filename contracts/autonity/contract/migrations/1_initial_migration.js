const Migrations = artifacts.require("Migrations");

const Web3 = require('web3');
const TruffleConfig = require('../truffle-config');

module.exports = function(deployer, network, accounts) {
  const config = TruffleConfig.networks[network];
  const web3 = new Web3(new Web3.providers.HttpProvider('http://' + config.host + ':' + config.port));

  console.log('>> Unlocking account ' + accounts);
  web3.eth.personal.unlockAccount(accounts[0], "test", 36000);
  web3.eth.personal.unlockAccount(accounts[1], "test", 36000);
  web3.eth.personal.unlockAccount(accounts[2], "test", 36000);
  web3.eth.personal.unlockAccount(accounts[3], "test", 36000);
  web3.eth.personal.unlockAccount(accounts[4], "test", 36000);
  web3.eth.personal.unlockAccount(accounts[5], "test", 36000);
  web3.eth.personal.unlockAccount(accounts[6], "test", 36000);
  web3.eth.personal.unlockAccount(accounts[7], "test", 36000);
  web3.eth.personal.unlockAccount(accounts[8], "test", 36000);

  // Sleep for 2 seconds
  await new Promise(r => setTimeout(r, 2000));
  console.log('>> Deploying migration');
  deployer.deploy(Migrations);
};
