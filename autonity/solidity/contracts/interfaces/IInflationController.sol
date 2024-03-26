// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 ;
/**
 * @dev Interface of the Inflation Controller Contract
 */
interface IInflationController {
    /**
    * @notice Main function. Calculate NTN current supply delta.
    */
    function calculateSupplyDelta(uint256 _currentSupply, uint256 _lastEpochBlock, uint256 _currentBlock) external view returns (uint256);
}
