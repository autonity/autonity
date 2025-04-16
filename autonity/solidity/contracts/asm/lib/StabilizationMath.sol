// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "./ASMErrors.sol";
import {UD60x18, ud} from "../../lib/prb-math-4.0.1/UD60x18.sol";

library StabilizationMath {

    string internal constant NTN_SYMBOL = "NTN-ATN";
    string internal constant NTN_USD_SYMBOL = "NTN-USD";

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

    // Calculate the amount of collateral that can be received in a debt auction given the auction parameters.
    // @param startTimestamp The timestamp when the auction started
    // @param currentTimestamp The current timestamp
    // @param totalCollateral The maximum amount of collateral that can be received
    // @param liquidationRatio The liquidation ratio (must be > 1)
    // @param duration The duration of the auction
    // @return The amount of collateral to be received
    function sqrtIncreaseAuctionAmount(
        uint256 startTimestamp,
        uint256 currentTimestamp,
        uint256 totalCollateral,
        uint256 liquidationRatio,
        uint256 duration
    ) internal pure returns (uint256) {
        if (currentTimestamp < startTimestamp) {
            revert InvalidParameter("currentTimestamp > startTimestamp");
        }
        uint256 timeDelta = currentTimestamp - startTimestamp;
        // if the auction has been running for longer than the auction duration
        // the full collateral amount is receivable
        if (timeDelta >= duration) {
            return totalCollateral;
        }
        UD60x18 L = ud(liquidationRatio);
        UD60x18 C = ud(totalCollateral);
        UD60x18 t = ud(timeDelta * SCALE_FACTOR);
        UD60x18 T = ud(duration * SCALE_FACTOR);
        UD60x18 sqrtTau = t.div(T).sqrt();
        UD60x18 one = ud(SCALE_FACTOR);

        // C / (1 + (L - 1) * (1 - sqrt(t/T)))
        UD60x18 result = C.div(one.add((L.sub(one)).mul(one.sub(sqrtTau))));
        return result.intoUint256();
    }

    // Calculates the amount of collateral to be paid in an interest auction given the auction parameters.
    // @param startTimestamp The timestamp when the auction started
    // @param currentTimestamp The current timestamp
    // @param minimumOffer The minimum amount of collateral that can be paid (end of the auction)
    // @param maximumOffer The initial amount of collateral that can be paid (start of the auction)
    // @param duration The duration of the auction
    // @return The amount of collateral to be paid
    function linearDecreaseAuctionAmount(
        uint256 startTimestamp,
        uint256 currentTimestamp,
        uint256 minimumOffer,
        uint256 maximumOffer,
        uint256 duration
    ) internal pure returns (uint256){
        if (currentTimestamp < startTimestamp) {
            revert InvalidParameter("currentTimestamp > startTimestamp");
        }
        uint256 timeDelta = currentTimestamp - startTimestamp;
        // if the auction has been running for longer than the auction duration
        // the full collateral amount is receivable
        if (timeDelta >= duration) {
            return minimumOffer;
        }

        // if the auction has been running for less than the auction duration
        // the collateral received is a linear function of time
        return maximumOffer - ((maximumOffer - minimumOffer) * timeDelta) / duration;
    }

    /*
    ┌──────────────────┐
    │ CDP Calculations │
    └──────────────────┘
    */

    /// Calculate the maximum amount of ATN that can be borrowed for the
    /// given amount of Collateral Token.
    /// @param collateral Amount of Collateral Token backing the debt
    /// @param collateralPriceACU The price of Collateral Token in ACU
    /// @param targetDebtPriceACU The value of 1 unit of debt in ACU
    /// @param mcr The minimum collateralization ratio
    /// @return The maximum ATN that can be borrowed
    function borrowLimit(
        uint256 collateral,
        uint256 collateralPriceACU,
        uint256 targetDebtPriceACU,
        uint256 mcr
    ) internal pure returns (uint256) {
        if (collateralPriceACU == 0 || mcr == 0 || targetDebtPriceACU == 0) revert InvalidParameter("mcr || targetDebtPrice || collateralPrice");
        return (collateral * collateralPriceACU * SCALE_FACTOR) / (mcr * targetDebtPriceACU);
    }

    /// Calculate the minimum amount of Collateral Token that must be deposited
    /// in the CDP in order to borrow the given amount of Autons.
    /// @param principal Auton amount to borrow
    /// @param collateralPriceACU The price of Collateral Token in ACU
    /// @param targetDebtPriceACU The value of 1 unit of debt in ACU
    /// @param mcr The minimum collateralization ratio
    /// @return The minimum Collateral Token amount required
    function minimumCollateral(
        uint256 principal,
        uint256 collateralPriceACU,
        uint256 targetDebtPriceACU,
        uint256 mcr
    ) internal pure returns (uint256) {
        if (collateralPriceACU == 0 || mcr == 0) revert InvalidParameter("collateralPriceACU || mcr");
        return (principal * mcr * targetDebtPriceACU) / (collateralPriceACU * SCALE_FACTOR);
    }

    /// Calculate the interest due for a given amount of debt.
    /// @param debt The debt amount
    /// @param rateExponent The summation of the rates multiplied by their respective time window
    /// @return The interest due
    /// @dev Makes use of the prb-math library for natural exponentiation.
    function interestDue(
        uint256 debt,
        uint256 rateExponent
    ) internal pure returns (uint256) {
        UD60x18 d = ud(debt);
        UD60x18 rt = ud(rateExponent);
        UD60x18 exp = rt.exp();
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

    /// Determine the maximum amount of Auton that a user can borrow before their position becomes
    /// liquidatable
    /// @param collateral The collateral amount
    /// @param price The price of Collateral Token in ATN
    /// @param liquidationRatio The liquidation ratio, must be > 1
    /// @return The maximum amount of Auton that can be borrowed
    function debtLimit(
        uint256 collateral,
        uint256 price,
        uint256 liquidationRatio
    ) internal pure returns (uint256) {
        if (price == 0 || liquidationRatio == 0) revert InvalidParameter("price || liquidationRatio");
        return (collateral * price) / liquidationRatio;
    }

    /*
    ┌───────────┐
    │ Utilities │
    └───────────┘
    */

    /// Scale a value to SCALE_FACTOR.
    /// @param value The value to scale
    /// @param valueScaleFactor The scale factor of the value
    function toScaleFactor(uint256 value, uint256 valueScaleFactor) internal pure returns (uint256) {
        return (value * SCALE_FACTOR) / valueScaleFactor;
    }

    /**
     * @dev Calculates the interest exponent for given rate in the given time window
     * @param interestRate interest rate
     * @param startTimestamp start timestamp of the window in seconds
     * @param endTimestamp end timestamp of the window in seconds
     */
    function interestExponent(
        uint256 interestRate,
        uint256 startTimestamp,
        uint256 endTimestamp
    ) internal pure returns (uint256) {
        if (endTimestamp < startTimestamp) revert InvalidParameter("endTimestamp || startTimestamp");
        return interestRate * (endTimestamp - startTimestamp) / SECONDS_IN_YEAR;
    }
}
