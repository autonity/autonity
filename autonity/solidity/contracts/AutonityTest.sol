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
                Config memory _config) Autonity(_validators, _config) {

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
     return tailBondingID;
   }

   function getHeadBondingID() public view returns (uint256) {
     return headBondingID;
   }

   function getFirstPendingUnbondingRequest() public view returns (uint256) {
     return lastUnlockedUnbonding;     
   }

   function getHeadUnbondingID() public view returns (uint256) {
     return headUnbondingID;
   }

}
