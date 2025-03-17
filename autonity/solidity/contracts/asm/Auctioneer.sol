// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import {IStabilization} from "./interfaces/IStabilization.sol";
import {StabilizationMath} from "./lib/StabilizationMath.sol";
import "./lib/ASMErrors.sol";
import {AuctionLib} from "./lib/AuctionLib.sol";
import {IERC20} from "../interfaces/IERC20.sol";
import {IOracle} from "../interfaces/IOracle.sol";

contract Auctioneer {
    using AuctionLib for AuctionLib.AuctionSet;
    struct Config {
        uint256 liquidationAuctionDuration;

        uint256 interestAuctionDuration;
        uint256 interestAuctionDiscount; // value between [0,1) with SCALE_FACTOR precision

        uint256 interestAuctionThreshold; // in ATN
    }

    // Events
    event AuctionedDebt(address indexed debtor, address indexed biddor, uint256 collateralAmount, uint256 debtAmount);
    event AuctionedInterest(address indexed biddor, uint256 interestAmount, uint256 paymentAmount);
    event NewInterestAuction(uint256 auctionId, uint256 amount, uint256 startRound);
    event ConfigUpdated(string field);

    // Public state
    Config public config;
    IERC20 public collateralToken;
    address public proceedAddress;

    // Internal state
    IStabilization internal _stabilization;
    address internal _operator;
    address internal _autonity;
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

    modifier onlyOperator() {
        if (msg.sender != _operator) {
            revert Unauthorized();
        }
        _;
    }

    constructor(
        Config memory config_,
        address stabilization_,
        address oracle_,
        address collateralToken_,
        address autonity_,
        address operator_
    ) {
        _validateConfig(config_);
        config = config_;
        _stabilization = IStabilization(stabilization_);
        _oracle = IOracle(oracle_);
        collateralToken = IERC20(collateralToken_);
        _autonity = autonity_;
        _operator = operator_;
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
        IOracle.RoundData memory round = _oracle.getRoundData(liquidatableRound, StabilizationMath.NTN_SYMBOL);
        if (round.timestamp < cdp.timestamp) {
            revert InvalidRound(liquidatableRound);
        }

        uint256 debtAmount = _stabilization.debtAmountAtTime(debtor, block.timestamp);
        if (msg.value < debtAmount) {
            revert InvalidAmount();
        }

        // check if the CDP was liquidatable during the oracle round liquidatableRound
        if (
            !StabilizationMath.underCollateralized(
            cdp.collateral,
            round.price,
            _stabilization.debtAmountAtTime(debtor, round.timestamp),
            _stabilization.config().liquidationRatio)
        ) {
            revert InvalidRound(liquidatableRound);
        }

        uint256 maxNtnAmount = maxLiquidationReturn(debtor, liquidatableRound);

        if (ntnAmount > maxNtnAmount) {
            revert BidTooLow(maxNtnAmount, ntnAmount);
        }

        try _stabilization.liquidate{value: msg.value}(debtor, ntnAmount, msg.sender) {
            emit AuctionedDebt(debtor, msg.sender, msg.value, debtAmount);
        } catch {
            revert NotLiquidatable();
        }
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

        // if the proceeds address has not been set, the collateral will accumulate in this contract until
        // the next auction
        if (proceedAddress != address(0)) {
            if(!collateralToken.transfer(proceedAddress, collateralToken.balanceOf(address(this)))) {
                revert TransferFailed();
            }
        }
        emit AuctionedInterest(msg.sender, atnToReceive, ntnToPay);
    }

    /*
    ┌────────────────────────┐
    │ Permissioned Functions │
    └────────────────────────┘
    */

    // @notice Deposit interest payments into the contract
    // @dev This function is called by the stabilization mechanism contract
    function paidInterest() external payable onlyStabilization {
        _pendingAllocatedInterest += msg.value;
        if (_pendingAllocatedInterest >= config.interestAuctionThreshold) {
            uint256 startRound = _oracle.getRound() - 1;
            uint256 auction = auctions.push(_pendingAllocatedInterest, startRound, block.timestamp);
            emit NewInterestAuction(auction, _pendingAllocatedInterest, block.timestamp);
            _pendingAllocatedInterest = 0;
        }
    }

    // @notice Set the operator address
    // @param operator_ The address of the operator
    function setOperator(address operator_) external {
        if (msg.sender != _autonity) {
            revert Unauthorized();
        }
        _operator = operator_;
        emit ConfigUpdated("operator");
    }

    // Operator functions

    // @notice Set the oracle address
    // @param oracle_ The address of the oracle
    function setOracle(address oracle_) external onlyOperator {
        if (oracle_ == address(0)) {
            revert InvalidParameter("oracle_");
        }
        _oracle = IOracle(oracle_);
        emit ConfigUpdated("oracle");
    }

    // @notice Set the stabilization address
    // @param stabilization_ The address of the stabilization contract
    function setStabilization(address stabilization_) external onlyOperator {
        if (stabilization_ == address(0)) {
            revert InvalidParameter("stabilization_");
        }
        _stabilization = IStabilization(stabilization_);
        emit ConfigUpdated("stabilization");
    }

    // @notice Set the liquidation auction duration
    // @param duration The duration of the liquidation auction
    function setLiquidationAuctionDuration(uint256 duration) external onlyOperator {
        if (duration == 0) {
            revert InvalidParameter("duration");
        }
        config.liquidationAuctionDuration = duration;
        emit ConfigUpdated("liquidationAuctionDuration");
    }

    // @notice Set the interest auction duration
    // @param duration The duration of the interest auction
    function setInterestAuctionDuration(uint256 duration) external onlyOperator {
        if (duration == 0) {
            revert InvalidParameter("duration");
        }
        config.interestAuctionDuration = duration;
        emit ConfigUpdated("interestAuctionDuration");
    }

    // @notice Set the interest auction discount
    // @param discount The discount applied to the interest auction
    // @dev The discount is a value between [0,1) with SCALE_FACTOR precision
    function setInterestAuctionDiscount(uint256 discount) external onlyOperator {
        if (discount >= StabilizationMath.SCALE_FACTOR) {
            revert InvalidParameter("discount");
        }
        config.interestAuctionDiscount = discount;
        emit ConfigUpdated("interestAuctionDiscount");
    }

    // @notice Set the interest auction threshold
    // @param threshold The threshold for starting an interest auction
    function setInterestAuctionThreshold(uint256 threshold) external onlyOperator {
        if (threshold == 0) {
            revert InvalidParameter("threshold");
        }
        config.interestAuctionThreshold = threshold;
        emit ConfigUpdated("interestAuctionThreshold");
    }

    // @notice Set the proceeds address
    // @param proceedAddress_ The address to send proceeds to
    function setProceedAddress(address proceedAddress_) external onlyOperator {
        proceedAddress = proceedAddress_;
        emit ConfigUpdated("proceedAddress");
    }


    /*
    ┌────────────────┐
    │ View Functions │
    └────────────────┘
    */

    // @notice Get all open interest auctions
    // @return An array of all open interest auctions
    function openAuctions() external view returns (AuctionLib.Auction[] memory) {
        return auctions.values();
    }

    // @notice Get an auction by ID
    // @param auction The ID of the auction
    function getAuction(uint256 auction) external view returns (AuctionLib.Auction memory) {
        return auctions.get(auction);
    }

    // @notice Get the maximum amount of NTN that can be returned to a liquidator for a given CDP
    // @param debtor The address of the CDP owner
    // @param liquidatableRound The earliest round in which the CDP was liquidatable
    function maxLiquidationReturn(address debtor, uint256 liquidatableRound) public view returns (uint256) {
        IOracle.RoundData memory round = _oracle.getRoundData(liquidatableRound, StabilizationMath.NTN_SYMBOL);
        IStabilization.CDP memory cdp = _stabilization.cdps(debtor);
        return StabilizationMath.sqrtIncreaseAuctionAmount(
            round.timestamp,
            block.timestamp,
            cdp.collateral,
            _stabilization.config().liquidationRatio,
            config.liquidationAuctionDuration
        );
    }

    // @notice Get the minimum amount of NTN that can be paid for an interest auction
    // @param auction The ID of the auction
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
    function _calculateInitialCost(uint256 interestAmount, uint256 collateralPrice) internal view returns (uint256) {
        uint256 oracleScaleFactor = 10 ** _oracle.getDecimals();
        uint256 priceDiscounted = collateralPrice - (collateralPrice * config.interestAuctionDiscount) / StabilizationMath.SCALE_FACTOR;
        return (interestAmount * oracleScaleFactor) / priceDiscounted;
    }
    
    function _validateConfig(Config memory config_) internal pure {
        if (config_.liquidationAuctionDuration == 0) {
            revert InvalidParameter("liquidationAuctionDuration");
        }
        if (config_.interestAuctionDuration == 0) {
            revert InvalidParameter("interestAuctionDuration");
        }
        if (config_.interestAuctionDiscount >= StabilizationMath.SCALE_FACTOR) {
            revert InvalidParameter("interestAuctionDiscount");
        }
        if (config_.interestAuctionThreshold == 0) {
            revert InvalidParameter("interestAuctionThreshold");
        }
    }
}








