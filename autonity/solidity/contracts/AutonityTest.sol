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

     struct TestStruct {
          uint128 a;
          uint128 b;
          uint256 c;
          uint128 d;
          uint128 e;
     }
     TestStruct item;

   function applyNewCommissionRates() public onlyProtocol {
        Autonity._applyNewCommissionRates();
   }

   function applyStakingOperations() public {
       _stakingOperations();
   }

     function validatorsToSort() internal view returns (CommitteeMember[] memory) {
          CommitteeMember[] memory _validatorList = new CommitteeMember[](validatorList.length);
          for (uint256 i = 0; i < validatorList.length; i++) {
               Validator storage _user = validators[validatorList[i]];
               CommitteeMember memory _item = CommitteeMember(_user.nodeAddress, _user.bondedStake);
               _validatorList[i] = _item;
          }
          return _validatorList;
     }

     function isSorted(CommitteeMember[] memory _validatorList) internal pure {
          for (uint256 i = 1; i < _validatorList.length; i++) {
               require(_validatorList[i-1].votingPower >= _validatorList[i].votingPower, "not sorted");
          }
          require(_validatorList[0].votingPower > 0, "no positive stake");
     }

     function isSorted(address[] memory addresses) internal view {
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

   function testSorting() public {
       // testing  _sortByStakeOptimized
       // apply staking operations first, so everyone has positive stake
       _stakingOperations();
       CommitteeMember[] memory _validatorList = validatorsToSort();
       _sortByStakeOptimized(_validatorList);
       isSorted(_validatorList);
   }

     function testSortingPrecompiled() public {
          // testing  _sortByStakeOptimized
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          require(validatorList.length <= 100, "limited return data");
          CommitteeMember[] memory _validatorList = validatorsToSort();
          address[] memory addresses = _sortByStakePrecompiled(_validatorList, validatorList.length);
          isSorted(addresses);
     }

     function testSortLibrarySort() public {
          // testing  _sortByStakeOptimized
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          require(validatorList.length <= 100, "limited return data");
          CommitteeMember[] memory _validatorList = validatorsToSort();
          address[] memory addresses = _sortByStakeSortLibrarySort(_validatorList, validatorList.length);
          isSorted(addresses);
     }

     function testSortLibrarySliceTable() public {
          // testing  _sortByStakeOptimized
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          require(validatorList.length <= 100, "limited return data");
          CommitteeMember[] memory _validatorList = validatorsToSort();
          address[] memory addresses = _sortByStakeSortLibrarySliceTable(_validatorList, validatorList.length);
          isSorted(addresses);
     }

     function testSortingPrecompiledFast() public {
          // testing  _sortByStakeOptimized
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          require(validatorList.length <= 100, "limited return data");
          CommitteeMember[] memory _validatorList = validatorsToSort();
          address[] memory addresses = _sortByStakePrecompiledFast(_validatorList, validatorList.length);
          isSorted(addresses);
     }

     function testAssemblyProperArrray() public view returns (uint256, uint256, address, uint256) {
          require(validatorList.length > 0, "no validators");
          address key = validatorList[0];
          Validator storage validator = validators[key];
          // uint256 storage bondedStake = validator.bondedStake;
          uint256[2] memory location;
          uint256 mapLocation;
          assembly {
               mapLocation := validator.offset
               // location := validator.slot
               mstore(add(location, 0x20), validator.slot)
               mstore(location, validators.slot)
          }
          uint256 calculatedLocation = uint256(keccak256(abi.encode(key, location[0])));
          require(calculatedLocation == location[1], "location mismatch");
          return (location[1], calculatedLocation, key, location[0]);
     }

     function sort() public returns (address[] memory) {
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          require(validatorList.length <= 100, "limited return data");
          CommitteeMember[] memory _validatorList = validatorsToSort();
          address[] memory addresses = _sortByStakePrecompiledIterate(_validatorList, validatorList.length);
          require(addresses.length == _validatorList.length, "return data wrong");
          return addresses;
     }

     function testSortingPrecompiledIterate() public {
          // testing  _sortByStakeOptimized
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          require(validatorList.length <= 100, "limited return data");
          CommitteeMember[] memory _validatorList = validatorsToSort();
          address[] memory addresses = _sortByStakePrecompiledIterate(_validatorList, validatorList.length);
          isSorted(addresses);
     }

     function testSortingPrecompiledIterateFast() public {
          // testing  _sortByStakeOptimized
          // apply staking operations first, so everyone has positive stake
          _stakingOperations();
          require(validatorList.length <= 100, "limited return data");
          CommitteeMember[] memory _validatorList = validatorsToSort();
          address[] memory addresses = _sortByStakePrecompiledIterateFast(_validatorList, validatorList.length);
          isSorted(addresses);
     }

     function testStructLocation(uint128 a, uint128 b, uint256 c, uint128 d, uint128 e) public returns (TestStruct memory) {
          // while (item.length < 2) {
          //      item.push(TestStruct(1,1,1,1,1));
          // }
          item = TestStruct(
               a, b, c, d, e
          );
          // item[1] = TestStruct(
          //      a, b, c, d, e
          // );
          uint256[1] memory input;
          address to = Precompiled.TEST_LOCATION_CONTRACT;
          uint256[1] memory _returnData;
          assembly {
               mstore(input, item.slot)
               //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
               if iszero(staticcall(gas(), to, input, 32, _returnData, 32)) {
                    revert(0, 0)
               }
          }
          require(_returnData[0] == 1, "unsuccessful call");
          return item;
     }

     function getItem() public view returns (TestStruct memory) {
          return item;
     }

     function getValidatorListSlot() public pure returns (uint256) {
          uint256 slot;
          assembly {
               slot := validatorList.slot
          }
          return slot;
     }

     function getValidatorsSlot() public pure returns (uint256) {
          uint256 slot;
          assembly {
               slot := validators.slot
          }
          return slot;
     }

     function getCommitteeSlot() public pure returns (uint256) {
          uint256 slot;
          assembly {
               slot := committee.slot
          }
          return slot;
     }

     function getCommitteeMemberSlot(uint256 _idx) public view returns (uint256) {
          require(committee.length > _idx, "no member");
          uint256 slot;
          CommitteeMember storage member = committee[_idx];
          assembly {
               slot := member.slot
          }
          return slot;
     }

     function getCommitteeMember(uint256 _idx) public view returns (CommitteeMember memory) {
          require(committee.length > _idx, "no member");
          CommitteeMember memory member = committee[_idx];
          return member;
     }

     function getEpochTotalBondedStake() public view returns (uint256) {
          return epochTotalBondedStake;
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
