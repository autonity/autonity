// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "../asm/interfaces/IStabilization.sol";

contract StabilizationMock is IStabilization {
    // Public state retrieval functions
    function config() external view returns (Config memory) {
        return Config(0,0,0,0,0,0);
    }

    function cdps(address owner) external view returns (CDP memory) {
        return CDP(0,0,0,0,0);
    }

    // view functions
    function debtAmountAtTime(address account, uint timestamp) external view returns (uint256) {
        return 0;
    }

    function debtAmount(address account) external view returns (uint256) {
        return 0;
    }

    /// Liquidate an undercollateralized CDP.
    /// @param account The address of the CDP owner.
    /// @param collateralSold The amount of collateral to sell.
    /// @param bidder The address of the bidder to receive the collateral.
    /// @dev Restricted to the Auctioneer Contract.
    function liquidate(address account, uint256 collateralSold, address bidder) external payable {
        // do nothing
    }

    /// Set the Governance Operator account address.
    /// @param operator Address of the new Governance Operator
    /// @dev Restricted to the Autonity Contract.
    function setOperator(address operator) external {
        // do nothing
    }

    /// Set the Oracle Contract address.
    /// @param oracle Address of the new Oracle Contract
    /// @dev Restricted to the Autonity Contract.
    function setOracle(address oracle) external {
        // do nothing
    }
}