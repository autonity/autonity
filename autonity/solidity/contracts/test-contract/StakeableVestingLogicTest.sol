// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../vesting/stakeable/StakeableVestingLogic.sol";

contract StakeableVestingLogicTest is StakeableVestingLogic {
    constructor(address payable _autonity) StakeableVestingLogic(_autonity) {}

    function clearValidators() public {
        _clearValidators();
    }
}