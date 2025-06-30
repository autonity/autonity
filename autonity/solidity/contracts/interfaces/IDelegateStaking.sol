// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

interface IDelegateStaking {

    /**
     * @notice Returns the remaining number of NTN (LNTN) that `staker` will be
     * allowed to bond (unbond) on behalf of `owner` through `bondFrom` (`unbondFrom`).
     * This is zero by default.
     */
    function stakeAllowance(address owner, address staker) external view returns (uint256);

    /**
     * @notice Sets `amount` as the stake-allowance of `staker` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits an {StakeApproval} event.
     */
    function approveStake(address staker, uint256 amount) external returns (bool);

    /**
     * @notice Emitted when the stake-allowance of a `staker` for an `owner` is set by
     * a call to `approveStake`. `value` is the new stake-allowance.
     */
    event StakeApproval(address indexed owner, address indexed staker, uint256 value);

}