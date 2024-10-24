// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../vesting/stakable/StakableVestingLogic.sol";

contract StakableVestingLogicTest is StakableVestingLogic {
    constructor(address payable _autonity) StakableVestingLogic(_autonity) {}

    function clearValidators() public {
        _clearValidators();
    }
}