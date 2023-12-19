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

   function testSorting() public {
       // testing  _sortByStakeOptimized
       // apply staking operations first, so everyone has positive stake
       _stakingOperations();
       CommitteeMember[] memory _validatorList = new CommitteeMember[](validatorList.length);
       for (uint256 i = 0; i < validatorList.length; i++) {
            Validator storage _user = validators[validatorList[i]];
            CommitteeMember memory _item = CommitteeMember(_user.nodeAddress, _user.bondedStake);
            _validatorList[i] = _item;
       }
       _sortByStakeOptimized(_validatorList);
       
       // check if sorted
       for (uint256 i = 1; i < _validatorList.length; i++) {
           require(_validatorList[i-1].votingPower >= _validatorList[i].votingPower, "not sorted");
       }
       require(_validatorList[0].votingPower > 0, "no positive stake");
   }

     function testSortingPrecompiled() public {
          // testing  _sortByStakeOptimized
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          CommitteeMember[] memory _validatorList = new CommitteeMember[](validatorList.length);
          for (uint256 i = 0; i < validatorList.length; i++) {
               Validator storage _user = validators[validatorList[i]];
               CommitteeMember memory _item = CommitteeMember(_user.nodeAddress, _user.bondedStake);
               _validatorList[i] = _item;
          }
          address[] memory addresses = _sortByStakePrecompiled(_validatorList, validatorList.length);

          // check if sorted
          uint256 lastStake = 0;
          for (uint256 i = 0; i < addresses.length; i++) {
               require(addresses[i] != address(0), "invalid address");
               Validator storage validator = validators[addresses[i]];
               require(validator.nodeAddress == addresses[i], "validator does not exit");

               if (i > 0) {
                    require(validator.bondedStake <= lastStake, "not sorted");
               }
               lastStake = validator.bondedStake;
          }
     }

     function testSortingPrecompiledFast() public {
          // testing  _sortByStakeOptimized
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          CommitteeMember[] memory _validatorList = new CommitteeMember[](validatorList.length);
          for (uint256 i = 0; i < validatorList.length; i++) {
               Validator storage _user = validators[validatorList[i]];
               CommitteeMember memory _item = CommitteeMember(_user.nodeAddress, _user.bondedStake);
               _validatorList[i] = _item;
          }
          address[] memory addresses = _sortByStakePrecompiledFast(_validatorList, validatorList.length);

          // check if sorted
          uint256 lastStake = 0;
          for (uint256 i = 0; i < addresses.length; i++) {
               require(addresses[i] != address(0), "invalid address");
               Validator storage validator = validators[addresses[i]];
               require(validator.nodeAddress == addresses[i], "validator does not exit");

               if (i > 0) {
                    require(validator.bondedStake <= lastStake, "not sorted");
               }
               lastStake = validator.bondedStake;
          }
     }

     function testCommitteeStruct(uint256 count) public {
          if (count > validatorList.length) {
               count = validatorList.length;
          }
          _stakingOperations();
          CommitteeMember[] memory _validatorList = new CommitteeMember[](count);
          for (uint256 i = 0; i < count; i++) {
               Validator storage _user = validators[validatorList[i]];
               CommitteeMember memory _item = CommitteeMember(_user.nodeAddress, _user.bondedStake);
               _validatorList[i] = _item;
          }

          uint _length = 1000;
          uint _returnDataLength = 32;
          uint256[1] memory _returnData;
          address to = address(0xf9);
          assembly {
               //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
               if iszero(staticcall(gas(), to, _validatorList, _length, _returnData, _returnDataLength)) {
                    revert(0, 0)
               }
          }

          require(_returnData[0] == 1, "unsuccessful call");
     }

     function testValidatorStruct(uint256 count) public {
          if (count > validatorList.length) {
               count = validatorList.length;
          }
          _stakingOperations();
          Validator[] memory _validatorList = new Validator[](count);
          for (uint256 i = 0; i < count; i++) {
               Validator memory _user = validators[validatorList[i]];
               _validatorList[i] = _user;
          }

          uint _length = 5000;
          uint _returnDataLength = 32;
          uint256[1] memory _returnData;
          address to = address(0xf9);
          assembly {
               //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
               if iszero(staticcall(gas(), to, _validatorList, _length, _returnData, _returnDataLength)) {
                    revert(0, 0)
               }
          }

          require(_returnData[0] == 1, "unsuccessful call");
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

}
