// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import {IERC20} from "../interfaces/IERC20.sol";
import {IOracle} from "../interfaces/IOracle.sol";
import {ISupplyControl} from "./ISupplyControl.sol";

contract StabilizationState {
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
    }

    /// Stabilization Configuration.
    struct Config {
        /// The annual continuously-compounded interest rate for borrowing.
        uint256 borrowInterestRate;
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

    address[] internal _accounts;
    address internal _autonity;
    address internal _operator;
    IERC20 internal _collateralToken;
    IOracle internal _oracle;
    ISupplyControl internal _supplyControl;

    /// The Config object that stores Stabilization Contract parameters.
    Config public config;
    /// A mapping to retrieve the CDP for an account address.
    mapping(address => CDP) public cdps;

    // Parameters for initial CDP restrictions
    bool internal _restricted;
    address internal _atnSupplyOperator;
    uint256 internal _defaultGenesisBorrowInterestRate;

    // auction
    address internal _auctioneer;
}