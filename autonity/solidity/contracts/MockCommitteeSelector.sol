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
                uint256 validatorStateSlot = uint256(keccak256(key))+19;
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
                uint256 stakeSlot = uint256(keccak256(key))+5;
                uint256 stateSlot = uint256(keccak256(key))+19;
                uint256 consensusKeySlot = uint256(keccak256(key))+18;
                uint256 bondedStake;
                ValidatorState state;
                uint256 consensusKeyData;
                assembly {
                    bondedStake := sload(stakeSlot)
                    state := sload(stateSlot)
                    consensusKeyData := sload(consensusKeySlot)
                }

                if (bondedStake > threshold && state == ValidatorState.active) {
                    bytes memory addressBytes = new bytes(32);
                    assembly {
                        mstore(add(addressBytes, 0x20), mload(add(key, 0x20)))
                    }
                    address nodeAddress = address(uint160(uint256(bytes32(addressBytes))));
                    bytes memory consensusKey;
                    uint256 consensusKeyLen;
                    if (consensusKeyData % 2 == 0) {
                        // short bytes
                        consensusKeyLen = (consensusKeyData % 256) / 2; // 2*length is stored in last byte
                        consensusKey = new bytes(consensusKeyLen);
                        uint256 pos = 32;
                        for (uint256 k = 0; k < consensusKeyLen; k++) {
                            pos--;
                            uint256 div = 256 ** pos;
                            consensusKey[k] = bytes1(uint8((consensusKeyData / div) % 256));
                        }
                    }
                    else {
                        // long bytes
                        consensusKeyLen = (consensusKeyData - 1) / 2;
                        consensusKey = new bytes(consensusKeyLen);
                        uint256 consensusKeyOffset = uint256(keccak256(abi.encodePacked(consensusKeySlot)));
                        for (uint256 k = 0; k < consensusKeyLen; k += 32) {
                            uint256 consensusChunk;
                            assembly {
                                consensusChunk := sload(consensusKeyOffset)
                            }
                            consensusKeyOffset++;
                            uint256 pos = 32;
                            for (uint256 h = 0; h < 32 && h+k < consensusKeyLen; h++) {
                                pos--;
                                uint256 div = 256 ** pos;
                                consensusKey[h+k] = bytes1(uint8((consensusChunk / div) % 256));
                            }
                        }
                    }
                    validators[j] = Autonity.CommitteeMember(nodeAddress, bondedStake, consensusKey);
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
                bytes memory consensusKey = validators[i].consensusKey;
                totalStake += bondedStake;
                assembly {
                    sstore(committeeSlotBase, nodeAddress)
                }
                committeeSlotBase++;
                assembly {
                    sstore(committeeSlotBase, bondedStake)
                }
                committeeSlotBase++;
                bytes memory consensusKeyChunk = new bytes(32);
                if (consensusKey.length < 32) {
                    // short bytes
                    consensusKeyChunk[31] = bytes1(uint8(consensusKey.length * 2));
                    for (uint256 j = 0; j < consensusKey.length; j++) {
                        consensusKeyChunk[j] = consensusKey[j];
                    }
                    assembly {
                        sstore(committeeSlotBase, mload(add(consensusKeyChunk, 0x20)))
                    }
                }
                else {
                    // long bytes
                    consensusKeyChunk = abi.encodePacked(consensusKey.length * 2 + 1);
                    assembly {
                        sstore(committeeSlotBase, mload(add(consensusKeyChunk, 0x20)))
                    }
                    uint256 consensusKeyOffset = uint256(keccak256(abi.encodePacked(committeeSlotBase)));
                    for (uint256 j = 0; j < consensusKey.length; j += 32) {
                        for (uint256 k = 0; k < 32 && k+j < consensusKey.length; k++) {
                            consensusKeyChunk[k] = consensusKey[k+j];
                        }
                        assembly {
                            sstore(consensusKeyOffset, mload(add(consensusKeyChunk, 0x20)))
                        }
                        consensusKeyOffset++;
                    }
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
