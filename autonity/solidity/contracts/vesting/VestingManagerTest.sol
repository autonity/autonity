// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./VestingManager.sol";

contract VestingManagerTest is VestingManager {
    constructor(address payable _autonity, address _operator) VestingManager(_autonity, _operator) {}

    // function applyBonding(uint256 _bondingID, uint256 _liquid, bool _selfDelegation, bool _rejected) public {
    //     _applyBonding(_bondingID, _liquid, _selfDelegation, _rejected);
    // }

    // function applyUnbonding(uint256 _unbondingID) public {
    //     _applyUnbonding(_unbondingID);
    // }

    // function releaseUnbonding(uint256 _unbondingID, uint256 _amount) public {
    //     _releaseUnbonding(_unbondingID, _amount);
    // }
    /** @dev Callback function for autonity at finalize function when epoch is ended.
     * This function handles the bonding and unbonding related mechanism.
     * Follow _stakingOperations() function in Autonity.sol to apply the operations correctly.
     * Note: Autonity does not know about vesting manager and will use the callback for any contract
     * that has bonded or unbonded in autonity contract.
     * @param _bonding list of bonding requests that are applied
     * @param _rejectedBonding list of bonding id that are rejected
     * @param _releasedUnbonding list of unbonding requests that are applied
     */
    function callFinalize(
        BondingApplied[] memory _bonding,
        uint256[] memory _rejectedBonding,
        UnbondingReleased[] memory _releasedUnbonding
    ) public {
        // _finalize(_bonding, _rejectedBonding, _releasedUnbonding);
    }
}