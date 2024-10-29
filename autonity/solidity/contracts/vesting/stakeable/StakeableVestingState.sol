// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../lib/DelegateCaller.sol";
import "./StakeableVestingStorage.sol";

/**
 * @title State of the Stakeable Vesting Smart Contract for vesting and staking funds
 * @notice Each stakeable vesting contract has its own separate smart contract.
 */
contract StakeableVestingState is StakeableVestingStorage {
    using DelegateCaller for address;

    constructor(address payable _autonity) AccessAutonity(_autonity) {
        managerContract = IStakeableVestingManager(payable(msg.sender));
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_stakeableVestingLogicContract()`. Will run if no other
     * function in the contract matches the call data.
     */
    fallback() payable external {
        _stakeableVestingLogicContract().delegate();
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_stakeableVestingLogicContract()`. Will run if call data
     * is empty.
     */
    receive() payable external {
        _stakeableVestingLogicContract().delegate();
    }

    /**
     ============================================================

        Internals

     ============================================================
     */

    /**
     * @dev Fetch stakeable vesting logic contract address from autonity
     */
    function _stakeableVestingLogicContract() internal view returns (address) {
        address _address = managerContract.stakeableVestingLogicContract();
        require(_address != address(0), "stakeable vesting logic contract not set");
        return _address;
    }

}
