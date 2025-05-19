// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.19;

import {StabilizationMath} from "../asm/lib/StabilizationMath.sol";

contract StabilizationMathTest {
    using StabilizationMath for uint256;

    function sqrtIncreaseAuctionAmount(
        uint256 startTimestamp,
        uint256 currentTimestamp,
        uint256 maximumOffer,
        uint256 minimumOffer,
        uint256 duration
    ) external pure returns (uint256) {
        return StabilizationMath.sqrtIncreaseAuctionAmount(
            startTimestamp,
            currentTimestamp,
            maximumOffer,
            minimumOffer,
            duration
        );
    }

    function linearDecreaseAuctionAmount(
        uint256 startTimestamp,
        uint256 currentTimestamp,
        uint256 minimumOffer,
        uint256 initialOffer,
        uint256 duration
    ) external pure returns (uint256) {
        return StabilizationMath.linearDecreaseAuctionAmount(
            startTimestamp,
            currentTimestamp,
            minimumOffer,
            initialOffer,
            duration
        );
    }

    function SCALE_FACTOR() external pure returns (uint256) {
        return StabilizationMath.SCALE_FACTOR;
    }

    function SECONDS_IN_YEAR() external pure returns (uint256) {
        return StabilizationMath.SECONDS_IN_YEAR;
    }

    function NTN_SYMBOL() external pure returns (string memory) {
        return StabilizationMath.NTN_SYMBOL;
    }

    function SCALE() external pure returns (uint256) {
        return StabilizationMath.SCALE;
    }
}