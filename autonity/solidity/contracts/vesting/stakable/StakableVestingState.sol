// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../lib/DelegateCaller.sol";
import "./StakableVestingStorage.sol";

contract StakableVestingState is StakableVestingStorage {
    using DelegateCaller for address;

    constructor(address payable _autonity) AccessAutonity(_autonity) {
        managerContract = IStakableVestingManager(payable(msg.sender));
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_stakableVestingLogicContract()`. Will run if no other
     * function in the contract matches the call data.
     */
    fallback() payable external {
        _stakableVestingLogicContract().delegate();
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_stakableVestingLogicContract()`. Will run if call data
     * is empty.
     */
    receive() payable external {
        _stakableVestingLogicContract().delegate();
    }

    /**
     ============================================================

        Internals

     ============================================================
     */

    /**
     * @dev Fetch stakable vesting logic contract address from autonity
     */
    function _stakableVestingLogicContract() internal view returns (address) {
        address _address = managerContract.stakableVestingLogicContract();
        require(_address != address(0), "stakable vesting logic contract not set");
        return _address;
    }

}
