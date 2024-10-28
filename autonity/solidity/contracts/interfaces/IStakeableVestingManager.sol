// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2;

interface IStakeableVestingManager {
    function stakeableVestingLogicContract() external view returns (address);
}