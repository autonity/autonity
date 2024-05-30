// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface INonStakableVestingVault {

    function unlockTokens() external returns (uint256 _totalNewUnlocked, uint256 _newUnlockedUnsubscribed);

}