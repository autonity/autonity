// SPDX-License-Identifier: LGPL-3.0-only

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

    function applyStakingOperations() public {
       _stakingOperations();
   }

     function getEpochTotalBondedStake() public view returns (uint256) {
          return epochTotalBondedStake;
     }

   function getBondingRequest(uint256 _id) public view returns (BondingRequest memory) {
        return bondingMap[_id];
   }

   function getUnbondingRequest(uint256 _id) public view returns (UnbondingRequest memory) {
        return unbondingMap[_id];
   }

   function getTailBondingID() public view returns (uint256) {
     return tailBondingID;
   }

   function getHeadBondingID() public view returns (uint256) {
     return headBondingID;
   }

   function getLastUnlockedUnbonding() public view returns (uint256) {
     return lastUnlockedUnbonding;     
   }

   function getHeadUnbondingID() public view returns (uint256) {
     return headUnbondingID;
   }

   function testComputeCommittee() public {
      _stakingOperations();
      (address[] memory voters, address[] memory reporters) = computeCommittee();
      address[] memory addresses = new address[](voters.length);
      uint256 totalStake = 0;
      uint256 lastStake = 0;
      require(committee.length <= config.protocol.committeeSize, "committee size exceeds MaxCommitteeSize");
      for (uint256 i = 0; i < committee.length; i++) {
        address memberAddress = committee[i].addr;
        require(memberAddress != address(0), "invalid address");
        addresses[i] = memberAddress;
        uint256 stake = committee[i].votingPower;
        require(stake > 0, "0 stake in committee");
        totalStake += stake;
        if (i > 0) {
          require(lastStake >= stake, "committee members not sorted");
        }
        lastStake = stake;
        Validator storage validator = validators[memberAddress];
        require(validator.nodeAddress == memberAddress, "validator does not exist");
        require(validator.bondedStake == stake, "stake mismatch");
        require(validator.oracleAddress == voters[i], "oracle address mismatch");
        require(committee[i].addr == reporters[i], "accountability reporter address mismatch");
        require(validator.state == ValidatorState.active, "validator not active");
        require(keccak256(abi.encodePacked(validator.enode)) == keccak256(abi.encodePacked(committeeNodes[i])), "enode mismatch");
        require(keccak256(abi.encodePacked(validator.consensusKey)) == keccak256(abi.encodePacked(committee[i].consensusKey)), "consensus key mismatch");
      }
      require(totalStake == epochTotalBondedStake, "total stake mismatch");

      for (uint i = 0; i < validatorList.length; i++) {
        bool foundMatch = false;
        for (uint j = 0; j < addresses.length; j++) {
          if (validatorList[i] == addresses[j]) {
            foundMatch = true;
            break;
          }
        }

        Validator storage validator = validators[validatorList[i]];
        if (foundMatch == false) {
          require(validator.bondedStake <= lastStake, "high stake for non-committee member");
        }
      }
   }

}
