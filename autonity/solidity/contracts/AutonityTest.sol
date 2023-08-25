// SPDX-License-Identifier: MIT

pragma solidity ^0.8.3;

import "./interfaces/IERC20.sol";
import "./Liquid.sol";
import "./Upgradeable.sol";
import "./Precompiled.sol";
import "./Autonity.sol";

/** @title Proof-of-Stake Autonity Contract */

contract AutonityTest is Autonity {

    constructor(Validator[] memory _validators,
                Config memory _config,
                uint256 _unbodingPeriod) Autonity(_validators, _config) {

        config.policy.unbondingPeriod = _unbodingPeriod;
    }

   function applyNewCommissionRates() public onlyProtocol {
        Autonity._applyNewCommissionRates();
   }

   function getBondingRequest(uint256 _id) public view returns (BondingRequest memory) {
        return bondingMap[_id];
   }

   function getUnbondingRequest(uint256 _id) public view returns (UnbondingRequest memory) {
        return unbondingMap[_id];
   }

   function getFirstPendingBondingRequest() public view returns (uint256) {
     require(tailBondingID < headBondingID, "No pending bonding requests");
     return tailBondingID;
   }

   function getLastRequestedBondingRequest() public view returns (uint256) {
     require(headBondingID > 0, "No bonding is requested");
     return headBondingID - 1;     
   }

   function getFirstPendingUnbondingRequest() public view returns (uint256) {
     require(lastUnlockedUnbonding < headUnbondingID, "No pending unbonding request");
     return lastUnlockedUnbonding;     
   }

   function getLastRequestedUnbondingRequest() public view returns (uint256) {
     require(headUnbondingID > 0, "No unbonding is requested");
     return headUnbondingID - 1;
   }

}
