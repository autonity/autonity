// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "./StabilizationErrors.sol";
import {UD60x18, ud} from "../../lib/prb-math-4.0.1/UD60x18.sol";

library StabilizationMath {

    string internal constant NTN_SYMBOL = "NTN-ATN";
    /// The decimal places in fixed-point integer representation.
    uint256 internal constant SCALE = 18; // Match UD60x18
    /// The multiplier for scaling numbers to the required scale.
    uint256 internal constant SCALE_FACTOR = 10 ** SCALE;
    /// A year is assumed to have 365 days for interest rate calculations.
    uint256 internal constant SECONDS_IN_YEAR = 365 days;

    /*
    ┌──────────────────────┐
    │ Auction Calculations │
    └──────────────────────┘
    */
    function linearIncreaseAuctionAmount(
        uint256 startTimestamp,
        uint256 currentTimestamp,
        uint256 maximumOffer,
        uint256 initialOffer,
        uint256 duration
    ) internal pure returns (uint256){
        if (currentTimestamp <= startTimestamp) {
            return 0;
        }
        uint256 timeDelta = currentTimestamp - startTimestamp;
        // if the auction has been running for longer than the auction duration
        // the full collateral amount is receivable
        if (timeDelta >= duration) {
            return maximumOffer;
        }

        // if the auction has been running for less than the auction duration
        // the collateral receiva
        return initialOffer + ((maximumOffer - initialOffer) * timeDelta) / duration;
    }

    function linearDecreaseAuctionAmount(
        uint256 startTimestamp,
        uint256 currentTimestamp,
        uint256 minimumOffer,
        uint256 initialOffer,
        uint256 duration
    ) internal pure returns (uint256){
        if (currentTimestamp <= startTimestamp) {
            return 0;
        }
        uint256 timeDelta = currentTimestamp - startTimestamp;
        // if the auction has been running for longer than the auction duration
        // the full collateral amount is receivable
        if (timeDelta >= duration) {
            return minimumOffer;
        }

        // if the auction has been running for less than the auction duration
        // the collateral receiva
        return initialOffer - ((initialOffer - minimumOffer) * timeDelta) / duration;
    }

    /*
    ┌──────────────────┐
    │ CDP Calculations │
    └──────────────────┘
    */

    /// Calculate the maximum amount of Amount that can be borrowed for the
    /// given amount of Collateral Token.
    /// @param collateral Amount of Collateral Token backing the debt
    /// @param price The price of Collateral Token in Auton
    /// @param targetPrice The ACU value of 1 unit of debt
    /// @param mcr The minimum collateralization ratio
    /// @return The maximum Auton that can be borrowed
    function borrowLimit(
        uint256 collateral,
        uint256 price,
        uint256 targetPrice,
        uint256 mcr
    ) internal pure returns (uint256) {
        if (price == 0 || mcr == 0) revert InvalidParameter();
        return (collateral * price * targetPrice) / (mcr * SCALE_FACTOR);
    }

    /// Calculate the minimum amount of Collateral Token that must be deposited
    /// in the CDP in order to borrow the given amount of Autons.
    /// @param principal Auton amount to borrow
    /// @param price The price of Collateral Token in Auton
    /// @param mcr The minimum collateralization ratio
    /// @return The minimum Collateral Token amount required
    function minimumCollateral(
        uint256 principal,
        uint256 price,
        uint256 mcr
    ) internal pure returns (uint256) {
        if (price == 0 || mcr == 0) revert InvalidParameter();
        return (principal * mcr) / price;
    }

    /// Calculate the interest due for a given amount of debt.
    /// @param debt The debt amount
    /// @param rate The borrow interest rate
    /// @param timeBorrow The borrow time
    /// @param timeDue The time the interest is due
    /// @return
    /// @dev Makes use of the prb-math library for natural exponentiation.
    function interestDue(
        uint256 debt,
        uint256 rate,
        uint timeBorrow,
        uint timeDue
    ) internal pure returns (uint256) {
        if (timeBorrow > timeDue) revert InvalidParameter();
        UD60x18 d = ud(debt);
        UD60x18 r = ud(rate);
        UD60x18 t = ud(timeDue - timeBorrow).div(ud(SECONDS_IN_YEAR));
        UD60x18 exp = r.mul(t).exp();
        UD60x18 interest = d.mul(exp.sub(ud(SCALE_FACTOR)));
        return interest.intoUint256();
    }

    /// Determine if a debt position is undercollateralized.
    /// @param collateral The collateral amount
    /// @param price The price of Collateral Token in Auton
    /// @param debt The debt amount
    /// @param liquidationRatio The liquidation ratio
    /// @return Whether the position is liquidatable
    function underCollateralized(
        uint256 collateral,
        uint256 price,
        uint256 debt,
        uint256 liquidationRatio
    ) internal pure returns (bool) {
        if (debt == 0) return false;
        return (collateral * price) / debt < liquidationRatio;
    }
}
