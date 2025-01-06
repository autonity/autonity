// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import {IStabilization} from "./IStabilization.sol";
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

    function bidDebt(address debtor, uint256 liquidatableRound) external payable {
        IStabilization.CDP memory cdp = _stabilization.cdps(debtor);
        if (msg.value < cdp.principal + cdp.interest) {
            revert InvalidAmount();
        }
        IOracle.RoundData memory round = _oracle.getRoundData(liquidatableRound, StabilizationMath.NTN_SYMBOL);
        if (round.timestamp < cdp.timestamp) {
            revert InvalidRound(liquidatableRound);
        }

        if (
            !StabilizationMath.underCollateralized(
            cdp.collateral,
            round.price,
            cdp.principal + cdp.interest,
            _stabilization.config().liquidationRatio)
        ) {
            revert NotLiquidatable();
        }

        uint256 collateralToReceive = StabilizationMath.linearIncreaseAuctionAmount(
            round.timestamp,
            block.timestamp,
            cdp.collateral,
            _calculateInitialAmount(cdp.principal + cdp.interest, round.price),
            config.liquidationAuctionDuration
        );

        _stabilization.liquidate{value: msg.value}(debtor, collateralToReceive, msg.sender);
        emit AuctionedDebt(debtor, msg.sender, collateralToReceive, cdp.principal + cdp.interest);
    }

    function bidInterest(uint256 auction) external {
        AuctionLib.Auction storage interestAuction = auctions.get(auction);
        IOracle.RoundData memory round = _oracle.getRoundData(interestAuction.startRound, StabilizationMath.NTN_SYMBOL);
        uint256 ntnToPay = StabilizationMath.linearDecreaseAuctionAmount(
            round.timestamp,
            block.timestamp,
            0, // TODO: what is a logical minimum ?
            _calculateInitialCost(interestAuction.amount, round.price),
            config.interestAuctionDuration
        );

        if (collateralToken.allowance(msg.sender, address(this)) < ntnToPay) {
            revert InsufficientAllowance();
        }
        bool success = collateralToken.transferFrom(msg.sender, address(this), ntnToPay);
        if (!success) {
            revert TransferFailed();
        }

        auctions.remove(auction);

        // transfer ATN
        (bool ok, ) = msg.sender.call{value: interestAuction.amount, gas: 2300}("");
        if (!ok) {
            revert TransferFailed();
        }

        emit AuctionedInterest(msg.sender, interestAuction.amount, ntnToPay);
    }

    /*
    ┌────────────────────────┐
    │ Permissioned Functions │
    └────────────────────────┘
    */

    function paidInterest() external payable onlyStabilization {
        _pendingAllocatedInterest += msg.value;
        if (_pendingAllocatedInterest >= config.interestAuctionThreshold) {
            uint256 auction = auctions.push(_pendingAllocatedInterest, _oracle.getRound());
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

    /*
    ┌────────────────────┐
    │ Internal Functions │
    └────────────────────┘
    */

    function _calculateInitialAmount(
        uint256 debtAmount,
        uint256 collateralPrice
    ) internal view returns (uint256) {
        // TODO(scott): double check this calculation
        uint256 oracleScaleFactor = 10 ** _oracle.getDecimals();
        return oracleScaleFactor * debtAmount * StabilizationMath.SCALE_FACTOR / (collateralPrice * config.liquidationAuctionDiscount);
    }

    function _calculateInitialCost(uint256 interestAmount, uint256 collateralPrice) internal view returns (uint256) {
        // TODO(scott): double check this calculation
        uint256 oracleScaleFactor = 10 ** _oracle.getDecimals();
        return (interestAmount * StabilizationMath.SCALE_FACTOR * oracleScaleFactor) / (collateralPrice * config.interestAuctionDiscount);
    }
}








