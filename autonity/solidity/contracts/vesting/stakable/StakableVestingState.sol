// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../DelegateCaller.sol";
import "./StakableVestingStorage.sol";

contract StakableVestingState is StakableVestingStorage {

    constructor(address payable _autonity) AccessAutonity(_autonity) {
        managerContract = StakableVestingManager(payable(msg.sender));
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_stakableVestingLogicContract()`. Will run if no other
     * function in the contract matches the call data.
     */
    fallback() payable external {
        DelegateCaller._delegate(
            _stakableVestingLogicContract()
        );
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_stakableVestingLogicContract()`. Will run if call data
     * is empty.
     */
    receive() payable external {
        DelegateCaller._delegate(
            _stakableVestingLogicContract()
        );
    }

    /**
     ============================================================

        Internals

     ============================================================
     */

    /**
     * @dev Fetch liquid logic contract address from autonity
     */
    function _stakableVestingLogicContract() internal view returns (address) {
        address _address = managerContract.stakableVestingLogicContract();
        require(_address != address(0), "stakable vesting logic contract not set");
        return _address;
    }

}
