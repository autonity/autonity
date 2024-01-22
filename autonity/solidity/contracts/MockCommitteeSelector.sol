// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./Autonity.sol";

contract MockCommitteeSelector {

    constructor() {}

    fallback(bytes calldata input) external returns (bytes memory) {
        bytes memory ret = new bytes(32);
        if (input.length != 5*32) {
            return ret;
        }
        // read slots
        // detailed specifications on how to read state variable from slot: https://docs.soliditylang.org/en/latest/internals/layout_in_storage.html#storage-inplace-encoding
        bytes memory validatorListSlot = new bytes(32);
        bytes memory validatorsSlot = new bytes(32);
        bytes memory committeeSlot = new bytes(32);
        bytes memory epochTotalBondedStakeSlot = new bytes(32);
        bytes memory committeeSizeBytes = new bytes(32);
        assembly {
            calldatacopy(add(validatorListSlot, 0x20), input.offset, 0x20)
            calldatacopy(add(validatorsSlot, 0x20), add(input.offset, 0x20), 0x20)
            calldatacopy(add(committeeSlot, 0x20), add(input.offset, 0x40), 0x20)
            calldatacopy(add(epochTotalBondedStakeSlot, 0x20), add(input.offset, 0x60), 0x20)
            calldatacopy(add(committeeSizeBytes, 0x20), add(input.offset, 0x80), 0x20)
        }
        uint256 committeeSize = uint256(bytes32(committeeSizeBytes));
        uint256 validatorCount;
        assembly {
            validatorCount := sload(mload(add(validatorListSlot, 0x20)))
        }
        uint256 validatorListSlotHash = uint256(keccak256(validatorListSlot));
        uint256 threshold = 0;
        uint256 count = 0;

        // read validators info from storage
        {
            for (uint256 i = 0; i < validatorCount; i++) {
                bytes memory key = new bytes(64);
                uint256 slot = validatorListSlotHash + i;
                assembly {
                    mstore(add(key, 0x20), sload(slot))
                    mstore(add(key, 0x40), mload(add(validatorsSlot, 0x20)))
                }
                uint256 validatorStakeSlot = uint256(keccak256(key))+5;
                uint256 validatorStateSlot = uint256(keccak256(key))+18;
                uint256 bondedStake;
                ValidatorState state;
                assembly {
                    bondedStake := sload(validatorStakeSlot)
                    state := sload(validatorStateSlot)
                }

                if (bondedStake > threshold && state == ValidatorState.active) {
                    count++;
                }
            }
        }

        Autonity.CommitteeMember[] memory validators = new Autonity.CommitteeMember[](count);
        {
            uint256 j = 0;
            for (uint256 i = 0; i < validatorCount; i++) {
                bytes memory key = new bytes(64);
                uint256 slot = validatorListSlotHash + i;
                assembly {
                    mstore(add(key, 0x20), sload(slot))
                    mstore(add(key, 0x40), mload(add(validatorsSlot, 0x20)))
                }
                uint256 validatorStakeSlot = uint256(keccak256(key))+5;
                uint256 validatorStateSlot = uint256(keccak256(key))+18;
                uint256 bondedStake;
                ValidatorState state;
                assembly {
                    bondedStake := sload(validatorStakeSlot)
                    state := sload(validatorStateSlot)
                }

                if (bondedStake > threshold && state == ValidatorState.active) {
                    bytes memory addressBytes = new bytes(32);
                    assembly {
                        mstore(add(addressBytes, 0x20), mload(add(key, 0x20)))
                    }
                    address nodeAddress = address(uint160(uint256(bytes32(addressBytes))));
                    validators[j] = Autonity.CommitteeMember(nodeAddress, bondedStake);
                    j++;
                }
            }
        }

        if (committeeSize > validators.length) {
            committeeSize = validators.length;
        }
        _sortByStake(validators);

        // write committee nodes in persistent storage
        {
            assembly {
                sstore(mload(add(committeeSlot, 0x20)), committeeSize)
            }
            uint256 committeeSlotBase = uint256(keccak256(committeeSlot));
            uint256 totalStake = 0;
            for (uint256 i = 0 ; i < committeeSize; i++) {
                uint256 nodeAddress = uint256(uint160(validators[i].addr));
                uint256 bondedStake = validators[i].votingPower;
                totalStake += bondedStake;
                assembly {
                    sstore(committeeSlotBase, nodeAddress)
                }
                committeeSlotBase++;
                assembly {
                    sstore(committeeSlotBase, bondedStake)
                }
                committeeSlotBase++;
            }
            assembly {
                sstore(mload(add(epochTotalBondedStakeSlot, 0x20)), totalStake)
            }
        }

        // success
        ret[31] = bytes1(uint8(1));
        return ret;
    }

    function _sortByStake(Autonity.CommitteeMember[] memory _validators) internal pure {
        _structQuickSort(_validators, int(0), int(_validators.length - 1));
    }

    /**
    * @dev QuickSort algorithm sorting in ascending order by stake.
    */
    function _structQuickSort(Autonity.CommitteeMember[] memory _users, int _low, int _high) internal pure {

        int _i = _low;
        int _j = _high;
        if (_i == _j) return;
        uint _pivot = _users[uint(_low + (_high - _low) / 2)].votingPower;
        // Set the pivot element in its right sorted index in the array
        while (_i <= _j) {
            while (_users[uint(_i)].votingPower > _pivot) _i++;
            while (_pivot > _users[uint(_j)].votingPower) _j--;
            if (_i <= _j) {
                (_users[uint(_i)], _users[uint(_j)]) = (_users[uint(_j)], _users[uint(_i)]);
                _i++;
                _j--;
            }
        }
        // Recursion call in the left partition of the array
        if (_low < _j) {
            _structQuickSort(_users, _low, _j);
        }
        // Recursion call in the right partition
        if (_i < _high) {
            _structQuickSort(_users, _i, _high);
        }
    }

}
