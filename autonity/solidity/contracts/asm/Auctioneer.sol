// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import {IStabilization} from "./interfaces/IStabilization.sol";
import {StabilizationMath} from "./lib/StabilizationMath.sol";
import "./lib/StabilizationErrors.sol";
import {AuctionLib} from "./lib/AuctionLib.sol";
import {IERC20} from "../interfaces/IERC20.sol";
import {IOracle} from "../interfaces/IOracle.sol";

contract Auctioneer {
    using AuctionLib for AuctionLib.AuctionSet;
    struct Config {
        uint256 liquidationAuctionDuration;
        uint256 liquidationAuctionDiscount; // value between [0,1) with SCALE_FACTOR precision

        uint256 interestAuctionDuration;
        uint256 interestAuctionDiscount; // value between [0,1) with SCALE_FACTOR precision

        uint256 interestAuctionThreshold; // in ATN
    }

    event AuctionedDebt(address indexed debtor, address indexed biddor, uint256 collateralAmount, uint256 debtAmount);
    event AuctionedInterest(address indexed biddor, uint256 interestAmount, uint256 paymentAmount);
    event NewInterestAuction(uint256 auctionId, uint256 amount, uint256 startRound);

    // Public state
    Config public config;
    IERC20 public collateralToken;

    // Internal state
    IStabilization internal _stabilization;
    IOracle internal _oracle;
    AuctionLib.AuctionSet internal auctions;
    uint256 internal _pendingAllocatedInterest;

    // Modifiers

    modifier onlyStabilization() {
        if (msg.sender != address(_stabilization)) {
            revert Unauthorized();
        }
        _;
    }

    constructor(Config memory config_, address stabilization_, address oracle_, address collateralToken_) {
        config = config_;
        _stabilization = IStabilization(stabilization_);
        _oracle = IOracle(oracle_);
        collateralToken = IERC20(collateralToken_);
    }

    /*
    ┌────────────────────────┐
    │ External Functions     │
    └────────────────────────┘
    */

    // @notice Place a bid to liquidate a CDP that is undercollateralized
    // @param debtor The address of the CDP owner
    // @param liquidatableRound The earliest round in which the CDP was liquidatable
    // @param ntnAmount The amount of NTN to receive in exchange for paying off the debt
    // @dev The caller must send the debt amount in ATN (via msg.value), and ntnAmount must be less than or equal to
    // maxLiquidationReturn for the caller to successfully execute a liquidation.
    function bidDebt(address debtor, uint256 liquidatableRound, uint256 ntnAmount) external payable {
        IStabilization.CDP memory cdp = _stabilization.cdps(debtor);
        uint256 debtAmount = _stabilization.debtAmount(debtor, block.timestamp);
        if (msg.value < debtAmount) {
            revert InvalidAmount();
        }
        IOracle.RoundData memory round = _oracle.getRoundData(liquidatableRound, StabilizationMath.NTN_SYMBOL);
        if (round.timestamp < cdp.timestamp) {
            revert InvalidRound(liquidatableRound);
        }

        // check if the CDP was liquidatable during the oracle round liquidatableRound
        if (
            !StabilizationMath.underCollateralized(
            cdp.collateral,
            round.price,
            debtAmount,
            _stabilization.config().liquidationRatio)
        ) {
            revert NotLiquidatable();
        }

        uint256 maxNtnAmount = maxLiquidationReturn(debtor, liquidatableRound);

        if (ntnAmount > maxNtnAmount) {
            revert BidTooLow(maxNtnAmount, ntnAmount);
        }

        _stabilization.liquidate{value: msg.value}(debtor, ntnAmount, msg.sender);
        emit AuctionedDebt(debtor, msg.sender, ntnAmount, debtAmount);
    }

    // @notice Place a bid on an interest auction
    // @param auction The ID of the auction
    // @param ntnAmount The amount of NTN to pay for the interest (must be greater than or equal to minInterestPayment)
    function bidInterest(uint256 auction, uint256 ntnAmount) external {
        AuctionLib.Auction storage interestAuction = auctions.get(auction);
        uint256 ntnToPay = minInterestPayment(auction);
        uint256 atnToReceive = interestAuction.amount;

        if (ntnAmount < ntnToPay) {
            revert BidTooLow(ntnToPay, ntnAmount);
        }

        if (collateralToken.allowance(msg.sender, address(this)) < ntnAmount) {
            revert InsufficientAllowance();
        }

        bool success = collateralToken.transferFrom(msg.sender, address(this), ntnAmount);
        if (!success) {
            revert TransferFailed();
        }

        auctions.remove(auction);

        // transfer ATN
        (bool ok,) = msg.sender.call{value: atnToReceive, gas: 2300}("");
        if (!ok) {
            revert TransferFailed();
        }

        emit AuctionedInterest(msg.sender, atnToReceive, ntnToPay);
        // ToDo: send the NTN somewhere
    }

    /*
    ┌────────────────────────┐
    │ Permissioned Functions │
    └────────────────────────┘
    */

    function paidInterest() external payable onlyStabilization {
        _pendingAllocatedInterest += msg.value;
        if (_pendingAllocatedInterest >= config.interestAuctionThreshold) {
            uint256 startRound = _oracle.getRound() - 1;
            uint256 auction = auctions.push(_pendingAllocatedInterest, startRound, block.timestamp);
            emit NewInterestAuction(auction, _pendingAllocatedInterest, block.timestamp);
            _pendingAllocatedInterest = 0;
        }
    }

    /*
    ┌────────────────┐
    │ View Functions │
    └────────────────┘
    */

    function openAuctions() external view returns (AuctionLib.Auction[] memory) {
        return auctions.values();
    }

    function getAuction(uint256 auction) external view returns (AuctionLib.Auction memory) {
        return auctions.get(auction);
    }

    function maxLiquidationReturn(address debtor, uint256 liquidatableRound) public view returns (uint256) {
        IOracle.RoundData memory round = _oracle.getRoundData(liquidatableRound, StabilizationMath.NTN_SYMBOL);
        IStabilization.CDP memory cdp = _stabilization.cdps(debtor);
        return StabilizationMath.sqrtIncreaseAuctionAmount(
            round.timestamp,
            block.timestamp,
            cdp.collateral,
            _calculateInitialReturn(
                cdp.collateral
            ),
            config.liquidationAuctionDuration
        );
    }

    function minInterestPayment(uint256 auction) public view returns (uint256) {
        AuctionLib.Auction storage interestAuction = auctions.get(auction);
        if (interestAuction.startTimestamp == 0) {
            revert InvalidAuctionId();
        }
        IOracle.RoundData memory round = _oracle.getRoundData(interestAuction.startRound, StabilizationMath.NTN_SYMBOL);
        return StabilizationMath.linearDecreaseAuctionAmount(
            interestAuction.startTimestamp,
            block.timestamp,
            0, // TODO: what is a logical minimum ?
            _calculateInitialCost(interestAuction.amount, round.price),
            config.interestAuctionDuration
        );
    }

    /*
    ┌────────────────────┐
    │ Internal Functions │
    └────────────────────┘
    */

    function _calculateInitialReturn(
        uint256 collateral
    ) internal view returns (uint256) {
        uint256 L = _stabilization.config().liquidationRatio;
        return collateral * StabilizationMath.SCALE_FACTOR / L;
    }

    function _calculateInitialCost(uint256 interestAmount, uint256 collateralPrice) internal view returns (uint256) {
        // TODO(scott): double check this calculation
        uint256 oracleScaleFactor = 10 ** _oracle.getDecimals();
        uint256 priceDiscounted = collateralPrice - (collateralPrice * config.interestAuctionDiscount) / StabilizationMath.SCALE_FACTOR;
        return (interestAmount * oracleScaleFactor) / priceDiscounted;
    }
}








