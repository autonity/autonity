// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./VestingManager.sol";

contract VestingManagerTest is VestingManager {
    constructor(address payable _autonity, address _operator) VestingManager(_autonity, _operator) {}

    function applyBonding(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) public {
        _applyBonding(_bondingID, _validator, _liquid, _selfDelegation, _rejected);
    }

    function applyUnbonding(uint256 _unbondingID, address _validator, bool _rejected) public {
        _applyUnbonding(_unbondingID, _validator, _rejected);
    }

    function releaseUnbonding(uint256 _unbondingID, uint256 _amount, bool _rejected) public {
        _releaseUnbonding(_unbondingID, _amount, _rejected);
    }

    function updateValidatorRewardAndRatio(address[] memory _validators) public {
        _updateValidatorRewardAndRatio(_validators);
    }
}