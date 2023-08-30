// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
/**
 * @dev Interface of the Autonity Contract.
 * Import this over Autonity.sol.
 */

interface IAutonity {

    /**
    * @notice Returns the current operator account.
    */
    function getOperator() external view returns (address);

    /**
    * @notice Returns the current Oracle account.
    */
    function getOracle() external view returns (address);
}
