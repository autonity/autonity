// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

/*
      .o.        .oooooo..o ooo        ooooo
     .888.      d8P'    `Y8 `88.       .888'
    .8"888.     Y88bo.       888b     d'888
   .8' `888.     `"Y8888o.   8 Y88. .P  888
  .88ooo8888.        `"Y88b  8  `888'   888
 .8'     `888.  oo     .d8P  8    Y     888
o88o     o8888o 8""88888P'  o8o        o888o

       Auton Stabilization Mechanism
*/

/// @title Stabilization Contract Interface
/// @dev Only meant to be used by the Autonity Contract.
interface IStabilization {
    /// Stabilization Configuration.
    struct Config {
        /// The annual continuously-compounded interest rate for borrowing.
        uint256 borrowInterestRate;
        /// Announcement window (in seconds) for borrow interest rate update.
        uint256 announcementWindow;
        /// The minimum ACU value of collateral required to maintain 1 ACU
        /// value of debt.
        uint256 liquidationRatio;
        /// The minimum ACU value of collateral required to borrow 1 ACU value
        /// of debt.
        uint256 minCollateralizationRatio;
        /// The minimum amount of debt required to maintain a CDP.
        uint256 minDebtRequirement;
        /// The ACU value of 1 unit of debt.
        uint256 targetPrice;
    }

    /// Represents a Collateralized Debt Position (CDP)
    struct CDP {
        /// The timestamp of the last borrow or repayment.
        uint timestamp;
        /// The collateral deposited with the Stabilization Contract.
        uint256 collateral;
        /// The principal debt outstanding as of `timestamp`.
        uint256 principal;
        /// The interest debt that is due at the `timestamp`.
        uint256 interest;
        /// aggregated interest exponent till last update.
        uint256 lastAggregatedInterestExponent;
    }

    // Public state retrieval functions
    function config() external view returns (Config memory);

    function cdps(address owner) external view returns (CDP memory);

    // view functions
    function debtAmountAtTime(address account, uint timestamp) external view returns (uint256);

    /// Liquidate an undercollateralized CDP.
    /// @param account The address of the CDP owner.
    /// @param collateralSold The amount of collateral to sell.
    /// @param bidder The address of the bidder to receive the collateral.
    /// @dev Restricted to the Auctioneer Contract.
    function liquidate(address account, uint256 collateralSold, address bidder) external payable;

    /// Set the Governance Operator account address.
    /// @param operator Address of the new Governance Operator
    /// @dev Restricted to the Autonity Contract.
    function setOperator(address operator) external;

    /// Set the Oracle Contract address.
    /// @param oracle Address of the new Oracle Contract
    /// @dev Restricted to the Autonity Contract.
    function setOracle(address oracle) external;
}
