const Autonity = artifacts.require("Autonity.sol");

module.exports = function(deployer, network, accounts) {
     const vals = [
    {
      "treasury": accounts[0],
      "nodeAddress": accounts[0],
      "oracleAddress": accounts[0],
      "enode": "enode://d73b857969c86415c0c000371bcebd9ed3cca6c376032b3f65e58e9e2b79276fbc6f59eb1e22fcd6356ab95f42a666f70afd4985933bd8f3e05beb1a2bf8fdde@172.25.0.11:30303",
      "commissionRate": 100,
      "bondedStake": 100,
      "selfBondedStake": 0,
      "totalSlashed": 0,
      "jailReleaseBlock": 0,
      "totalSlashed":0,
      "provableFaultCount" :0,
      "liquidContract": accounts[0],
      "liquidSupply": 0,
      "registrationBlock": 0,
      "state": 0,
    },
    {
      "treasury": accounts[1],
      "nodeAddress": accounts[1],
      "oracleAddress": accounts[1],
      "enode": "enode://1f207dfb3bcbbd338fbc991ec13e40d204b58fe7275cea48cfeb53c2c24e1071e1b4ef2959325fe48a5893de8ff37c73a24a412f367e505e5dec832813da546a@172.25.0.12:30303",
      "commissionRate": 100,
      "bondedStake": 90,
      "selfBondedStake": 0,
      "totalSlashed": 0,
      "jailReleaseBlock": 0,
      "totalSlashed":0,
      "provableFaultCount" :0,
      "liquidContract": accounts[1],
      "liquidSupply": 0,
      "registrationBlock": 0,
      "state": 0,
    },
    {
      "treasury": accounts[3],
      "nodeAddress": accounts[3],
      "oracleAddress": accounts[3],
      "enode": "enode://438a5c2cd8fdc2ecbc508bf7362e41c0f0c3754ba1d3267127a3756324caf45e6546b02140e2144b205aeb372c96c5df9641485f721dc7c5b27eb9e35f5d887b@172.25.0.14:30303",
      "commissionRate": 100,
      "bondedStake": 110,
      "selfBondedStake": 0,
      "totalSlashed": 0,
      "jailReleaseBlock": 0,
      "totalSlashed":0,
      "provableFaultCount" :0,
      "liquidContract": accounts[3],
      "liquidSupply": 0,
      "registrationBlock": 0,
      "state": 0,
    },
    {
      "treasury": accounts[4],
      "nodeAddress": accounts[4],
      "oracleAddress": accounts[4],
      "enode": "enode://3ce6c053cb563bfd94f4e0e248510a07ccee1bc836c9784da1816dba4b10564e7be1ba42e0bd8d73c8f6274f8e9878dc13814adb381c823264265c06048b4b59@172.25.0.15:30303",
      "commissionRate": 100,
      "bondedStake": 120,
      "selfBondedStake": 0,
      "totalSlashed": 0,
      "jailReleaseBlock": 0,
      "totalSlashed":0,
      "provableFaultCount" :0,
      "liquidContract": accounts[4],
      "liquidSupply": 0,
      "registrationBlock": 0,
      "state": 0,
    },
  ];



  const operator = accounts[5];
  const deployerAcc = accounts[6];
  const anyAccount = accounts[7];
  const name = "Newton";
  const symbol = "NTN";
  const minBaseFee = 5000;
  const committeeSize = 1000;
  const epochPeriod = 30;
  const delegationRate = 100;
  const unBondingPeriod = 60;
  const treasuryAccount = accounts[8];
  const treasuryFee = 1;
  const minimumEpochPeriod = 1;
  const version = 1;

  const config = {
    "operatorAccount": operator,
    "treasuryAccount": treasuryAccount,
    "treasuryFee": treasuryFee,
    "minBaseFee": minBaseFee,
    "delegationRate": delegationRate,
    "epochPeriod": epochPeriod,
    "unbondingPeriod": unBondingPeriod,
    "committeeSize": committeeSize,
    "contractVersion": version,
    "blockPeriod": minimumEpochPeriod,
    "oracleContract" : treasuryAccount, //temporary
    "accountabilityContract": treasuryAccount,
  };

    if (network != "test") {
        deployer.deploy(Autonity, vals, config, { from:deployerAcc} );
    } else {
        console.log("Skip migration for Autonity")
    }
};
