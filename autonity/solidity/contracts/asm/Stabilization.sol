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

import {IERC20} from "../interfaces/IERC20.sol";
import {IOracle} from "../interfaces/IOracle.sol";
import {IStabilization} from "./interfaces/IStabilization.sol";
import {ISupplyControl} from "./interfaces/ISupplyControl.sol";
import {IACU} from "./interfaces/IACU.sol";
import {UD60x18, ud} from "../lib/prb-math-4.0.1/UD60x18.sol";
import {StabilizationMath} from "./lib/StabilizationMath.sol";
import "./lib/StabilizationErrors.sol";
import {IAuctioneer} from "./interfaces/IAuctioneer.sol";
import {IConfigEvents} from "../interfaces/IConfigEvents.sol";

/// @title ASM Stabilization Contract
/// @notice A CDP-based stabilization mechanism for the Auton.
/// @dev Intended to be deployed by the protocol at genesis. Note that all
/// rates, ratios, prices, and amounts are represented as fixed-point integers
/// with `SCALE` decimal places.
/* solhint-disable not-rely-on-time */
contract Stabilization is IStabilization {

    /// The Config object that stores Stabilization Contract parameters.
    Config internal _config;
    /// A mapping to retrieve the CDP for an account address.
    mapping(address => CDP) internal _cdps;

    address[] private _accounts;
    address private _autonity;
    address private _operator;
    address private _auctioneer;
    address private _acu;
    IERC20 private _collateralToken;
    IOracle private _oracle;
    ISupplyControl private _supplyControl;

    // Parameters for initial CDP restrictions
    bool private _restricted;
    address private _atnSupplyOperator;
    uint256 private _defaultGenesisBorrowInterestRate;

    // For borrow interest rate update mechanism
    /**
     * @dev Stores the summation of all rates multiplied by their respective time window.
     * The current rate `config.borrowInterestRate` is not included in this because
     * the time window of `config.borrowInterestRate` is not finished yet.
     */
    uint256 private _aggregatedInterestExponent;
    /** @dev Time-stamp since the `config.borrowInterestRate` is active */
    uint256 private _borrowInterestActiveTimestamp;
    /** @dev Borrow interest rate that will be applied after some time */
    uint256 private _pendingBorrowInterestRate;
    /** @dev Time-stamp since the `_pendingBorrowInterestRate` will be active */
    uint256 private _pendingRateUpdateTimestamp;
    /** @dev Announcement window in pending which will be active since `_pendingAnnouncementUpdateTimestamp` timestamp (in seconds) */
    uint256 private _pendingAnnouncementWindow;
    /** @dev Timestamp since `_pendingAnnouncementWindow` will be active */
    uint256 private _pendingAnnouncementUpdateTimestamp;

    /// Collateral Token was deposited into a CDP
    /// @param account The CDP account address
    /// @param amount Collateral Token deposited
    event Deposit(address indexed account, uint256 amount);
    /// Collateral Token was withdrawn from a CDP
    /// @param account The CDP account address
    /// @param amount Collateral Token withdrawn
    event Withdraw(address indexed account, uint256 amount);
    /// Auton was borrowed from a CDP
    /// @param account The CDP account address
    /// @param amount Auton amount borrowed
    event Borrow(address indexed account, uint256 amount);
    /// Auton debt was paid into a CDP
    /// @param account The CDP account address
    /// @param amount Auton amount repaid
    event Repay(address indexed account, uint256 amount);
    /// A CDP was liquidated
    /// @param account The CDP account address
    /// @param liquidator The liquidator address
    event Liquidate(address indexed account, address liquidator);
    /// Transition out of the initial CDP restrictions
    event CDPRestrictionsRemoved();

    /**
     * @notice It is announced that borrow interest rate is going to be updated.
     * @param newRate The new borrow interest rate
     * @param activeSince Timestamp since the new rate will be active
     */
    event InterestRateUpdateAnnounced(uint256 newRate, uint256 activeSince, bool pendingRateOverridden);

    /**
     * @notice
     */
    event AnnouncementWindowUpdateAnnounced(uint256 newAnnouncementWindow, uint256 activeSince);

    modifier goodTime(address account, uint timestamp) {
        CDP storage cdp = _cdps[account];
        if (timestamp < cdp.timestamp) revert InvalidParameter();
        _;
    }

    modifier nonZeroAmount(uint256 amount) {
        if (amount == 0) revert InvalidAmount();
        _;
    }

    modifier onlyAutonity() {
        if (msg.sender != _autonity) revert Unauthorized();
        _;
    }

    modifier onlyAuctioneer() {
        if (msg.sender != _auctioneer) revert Unauthorized();
        _;
    }

    modifier onlyOperator() {
        if (msg.sender != _operator) revert Unauthorized();
        _;
    }

    modifier positiveMCR(uint256 ratio) {
        if (ratio == 0) revert InvalidParameter();
        _;
    }

    modifier validPrice(uint256 price) {
        if (price == 0) revert InvalidPrice();
        _;
    }

    modifier validRatios(
        uint256 liquidationRatio,
        uint256 minCollateralizationRatio
    ) {
        if (liquidationRatio >= minCollateralizationRatio)
            revert InvalidParameter();
        _;
    }

    // Restricted to the atnSupplyOperator during the initial CDP restrictions
    modifier restrictedSupplyOperator() {
        if (_restricted && msg.sender != _atnSupplyOperator) revert Unauthorized();
        _;
    }

    // Completely disabled during the initial CDP restrictions
    modifier restricted() {
        if (_restricted) revert Unauthorized();
        _;
    }

    /// Create and deploy the ASM Stabilization Contract.
    /// @param config_ Stabilization configuration
    /// @param autonity Address of the Autonity Contract
    /// @param operator Address of the Governance Operator
    /// @param oracle Address of the Oracle Contract
    /// @param supplyControl Address of the SupplyControl Contract
    /// @param auctioneer Address of the Auctioneer Contract
    /// @param acu Address of the ACU Contract
    /// @param collateralToken Address of the Collateral Token contract
    constructor(
        Config memory config_,
        address autonity,
        address operator,
        address oracle,
        address supplyControl,
        address auctioneer,
        address acu,
        IERC20 collateralToken
    )
    positiveMCR(config_.minCollateralizationRatio)
    validRatios(config_.liquidationRatio, config_.minCollateralizationRatio)
    {
        if (config_.announcementWindow == 0) revert ZeroValue();
        _config = config_;
        _autonity = autonity;
        _operator = operator;
        _oracle = IOracle(oracle);
        _supplyControl = ISupplyControl(supplyControl);
        _collateralToken = collateralToken;
        _auctioneer = auctioneer;
        _acu = acu;

        _restricted = true;
        _defaultGenesisBorrowInterestRate = config_.borrowInterestRate;
        _config.borrowInterestRate = 0;
    }

    /*
    ┌─────────────────┐
    │ Owner Functions │
    └─────────────────┘
    */

    /// Deposit Collateral Token using the ERC20 allowance mechanism.
    ///
    /// Before calling this function, the CDP owner must approve the
    /// Stabilization contract to spend Collateral Token on their behalf for
    /// the full amount to be deposited.
    /// @param amount Units of Collateral Token to deposit (non-zero)
    function deposit(uint256 amount) external nonZeroAmount(amount) restrictedSupplyOperator {
        if (_collateralToken.allowance(msg.sender, address(this)) < amount)
            revert InsufficientAllowance();

        CDP storage cdp = _cdps[msg.sender];
        if (cdp.timestamp == 0) _accounts.push(msg.sender);
        cdp.timestamp = block.timestamp; // opens the CDP
        cdp.collateral += amount;

        if (!_collateralToken.transferFrom(msg.sender, address(this), amount))
            revert TransferFailed();
        emit Deposit(msg.sender, amount);
    }

    /// Request a withdrawal of Collateral Token.
    ///
    /// The CDP must not be liquidatable and the withdrawal must not reduce the
    /// remaining Collateral Token amount below the minimum collateral amount.
    /// @param amount Units of Collateral Token to withdraw
    function withdraw(uint256 amount) external nonZeroAmount(amount) restrictedSupplyOperator {
        CDP storage cdp = _cdps[msg.sender];
        if (amount > cdp.collateral) revert InvalidAmount();
        (uint256 debt, , ) = _calculateDebtAmount(
            cdp,
            block.timestamp
        );
        uint256 price = collateralPrice();
        if (
            StabilizationMath.underCollateralized(
            cdp.collateral,
            price,
            debt,
            _config.liquidationRatio
        )
        ) revert Liquidatable();
        if (
            cdp.collateral - amount <
            StabilizationMath.minimumCollateral(
                cdp.principal,
                price,
                _config.minCollateralizationRatio
            )
        ) revert InsufficientCollateral();

        cdp.collateral -= amount;

        if (!_collateralToken.transfer(msg.sender, amount))
            revert TransferFailed();
        emit Withdraw(msg.sender, amount);
    }

    /// Borrow Auton against the CDP Collateral.
    ///
    /// The CDP must not be liquidatable, the `amount` must not exceed the
    /// borrow limit, the debt after borrowing must satisfy the minimum debt
    /// requirement.
    /// @param amount Auton to borrow
    function borrow(uint256 amount) external nonZeroAmount(amount) restrictedSupplyOperator {
        CDP storage cdp = _cdps[msg.sender];
        uint256 debt = _updateDebt(
            cdp,
            block.timestamp
        );
        debt += amount;
        if (debt < _config.minDebtRequirement) revert InvalidDebtPosition();
        uint256 price = collateralPrice();
        if (
            StabilizationMath.underCollateralized(
                cdp.collateral,
                price,
                debt,
                _config.liquidationRatio
            )
        ) revert Liquidatable();

        uint256 limit = maxBorrow(cdp.collateral);
        if (cdp.principal + amount > limit) revert InsufficientCollateral();

        cdp.principal += amount;
        _supplyControl.mint(msg.sender, amount);
        emit Borrow(msg.sender, amount);
    }

    /// Make a payment towards CDP debt.
    ///
    /// The transaction value is the payment amount. The debt after payment
    /// must satisfy the minimum debt requirement. The payment first covers
    /// the outstanding interest debt before the principal debt.
    function repay() external payable restrictedSupplyOperator {
        if (msg.value == 0) revert ZeroValue();
        CDP storage cdp = _cdps[msg.sender];
        if (cdp.principal == 0) revert NoDebtPosition();
        uint256 debt = _updateDebt(
            cdp,
            block.timestamp
        );
        if (
            (msg.value < debt) && (debt - msg.value < _config.minDebtRequirement)
        ) revert InvalidDebtPosition();

        (
            uint256 interestRecv,
            uint256 principalRecv,
            uint256 surplusRecv
        ) = _allocatePayment(cdp, msg.value);
        cdp.principal -= principalRecv;
        cdp.interest -= interestRecv;

        if (interestRecv > 0) IAuctioneer(_auctioneer).paidInterest{value: interestRecv}();
        if (principalRecv > 0) _supplyControl.burn{value: principalRecv}();
        if (surplusRecv > 0) payable(msg.sender).transfer(surplusRecv);
        emit Repay(msg.sender, msg.value);
    }

    /*
    ┌──────────────────┐
    │ Keeper Functions │
    └──────────────────┘
    */

    /// Liquidate a CDP that is undercollateralized.
    ///
    /// The liquidator must pay all the CDP debt outstanding. As a reward,
    /// the liquidator will receive the collateral that is held in the CDP. The
    /// transaction value is the payment amount. After covering the CDP's debt,
    /// any surplus is refunded to the liquidator.
    /// @param account The CDP account address to liquidate
    /// @param collateralSold The amount of collateral sold by the auctioneer
    /// @param bidder The address of the bidder
    function liquidate(address account, uint256 collateralSold, address bidder) external payable restricted onlyAuctioneer {
        if (msg.value == 0) revert ZeroValue();
        CDP storage cdp = _cdps[account];
        if (cdp.principal == 0) revert NoDebtPosition();
        if (cdp.collateral < collateralSold) revert InvalidAmount();
        (uint256 debt, uint256 accrued, ) = _calculateDebtAmount(
            cdp,
            block.timestamp
        );
        if (
            !StabilizationMath.underCollateralized(
            cdp.collateral,
            collateralPrice(),
            debt,
            _config.liquidationRatio
        )
        ) revert NotLiquidatable();

        if (msg.value < debt) revert InsufficientPayment();
        _supplyControl.burn{value: cdp.principal}();
        IAuctioneer(_auctioneer).paidInterest{value: accrued + cdp.interest}();

        uint surplus = msg.value - debt;

        uint256 collateral = cdp.collateral;
        cdp.timestamp = block.timestamp;
        cdp.collateral = collateral - collateralSold;
        cdp.principal = 0;
        cdp.interest = 0;

        if (!_collateralToken.transfer(bidder, collateralSold))
            revert TransferFailed();
        if (surplus > 0) payable(bidder).transfer(surplus);
        emit Liquidate(account, bidder);
    }

    /*
    ┌────────────────────┐
    │ Operator Functions │
    └────────────────────┘
    */

    /// Set the liquidation ratio.
    ///
    /// Must be less than the minimum collateralization ratio.
    /// @param ratio The liquidation ratio
    /// @dev Restricted to the operator.
    function setLiquidationRatio(
        uint256 ratio
    )
        external
        validRatios(ratio, _config.minCollateralizationRatio)
        onlyOperator
    {
        _config.liquidationRatio = ratio;
        emit IConfigEvents.ConfigUpdateUint("liquidationRatio", config.liquidationRatio, ratio);
    }

    /// Set the minimum collateralization ratio.
    ///
    /// Must be positive and greater than the liquidation ratio.
    /// @param ratio The minimum collateralization ratio
    /// @dev Restricted to the operator.
    function setMinCollateralizationRatio(
        uint256 ratio
    ) external positiveMCR(ratio) validRatios(_config.liquidationRatio, ratio) onlyOperator {
        _config.minCollateralizationRatio = ratio;
        emit IConfigEvents.ConfigUpdateUint("minCollateralizationRatio", config.minCollateralizationRatio, ratio);
    }

    /// Set the minimum debt requirement.
    /// @param amount The minimum debt amount
    /// @dev Restricted to the operator.
    function setMinDebtRequirement(uint256 amount) external onlyOperator {
        _config.minDebtRequirement = amount;
        emit IConfigEvents.ConfigUpdateUint("minDebtRequirement", config.minDebtRequirement, amount);
    }

    /// Set the SupplyControl Contract address.
    /// @param supplyControl The SupplyControl Contract address
    /// @dev Restricted to the operator.
    function setSupplyControl(address supplyControl) external onlyOperator {
        emit IConfigEvents.ConfigUpdateAddress("supplyControl", address(_supplyControl), supplyControl);
        _supplyControl = ISupplyControl(supplyControl);
    }

    /// Set the _atnSupplyOperator address.
    /// @param atnSupplyOperator The _atnSupplyOperator address
    /// @dev Restricted to the operator.
    function setAtnSupplyOperator(address atnSupplyOperator) external onlyOperator {
        emit IConfigEvents.ConfigUpdateAddress("atnSupplyOperator", _atnSupplyOperator, atnSupplyOperator);
        _atnSupplyOperator = atnSupplyOperator;
    }

    /// Transition out of the restricted state.
    /// @dev Restricted to the operator.
    function removeCDPRestrictions() external onlyOperator {
        if (_restricted == false) revert NotRestricted();
        _restricted = false;
        _config.borrowInterestRate = _defaultGenesisBorrowInterestRate;
        _borrowInterestActiveTimestamp = block.timestamp;
        emit CDPRestrictionsRemoved();
    }

    /**
     * @notice Updates the borrow interest rate. The new rate `newInterestRate` will take affect after the `config.announcementWindow` (in seconds).
     * @param newInterestRate The new interst rate multiplied by 10**18. If it is 5% then it should be `(5/100)*(10**18) = 50_000_000_000_000_000`
     */
    function updateBorrowInterestRate(uint256 newInterestRate) external restricted onlyOperator {
        _updateAnnouncementWindow();
        _applyInterestRateUpdate();
        bool pendingRateExist = false;
        if (_pendingRateUpdateTimestamp > 0) {
            // it exists, cannot be applied and going to be overridden the new one
            pendingRateExist = true;
        }
        _pendingBorrowInterestRate = newInterestRate;
        _pendingRateUpdateTimestamp = block.timestamp + _config.announcementWindow;

        emit InterestRateUpdateAnnounced(newInterestRate, _pendingRateUpdateTimestamp, pendingRateExist);
    }

    /**
     * @notice Updates the announcement window. The new window `window` will take affect after the `config.announcementWindow` (in seconds).
     * It requires that there is no announcement window in pending.
     */
    function updateAnnouncementWindow(uint256 window) external onlyOperator {
        if (window == 0) revert ZeroValue();
        _updateAnnouncementWindow();
        if (_pendingAnnouncementWindow > 0) revert AnnouncementWindowPending();
        _pendingAnnouncementWindow = window;
        _pendingAnnouncementUpdateTimestamp = block.timestamp + _config.announcementWindow;
        emit AnnouncementWindowUpdateAnnounced(window, _pendingAnnouncementUpdateTimestamp);
    }

    /*
    ┌────────────────────┐
    │ Autonity Functions │
    └────────────────────┘
    */

    /// Set the Governance Operator account address.
    /// @param operator Address of the new Governance Operator
    /// @dev Restricted to the Autonity Contract.
    function setOperator(address operator) external onlyAutonity {
        emit IConfigEvents.ConfigUpdateAddress("operator", _operator, operator);
        _operator = operator;
    }

    /// Set the Oracle Contract address.
    /// @param oracle Address of the new Oracle Contract
    /// @dev Restricted to the Autonity Contract.
    function setOracle(address oracle) external onlyAutonity {
        emit IConfigEvents.ConfigUpdateAddress("oracle", address(_oracle), oracle);
        _oracle = IOracle(oracle);
    }

    /*
    ┌────────────────┐
    │ View Functions │
    └────────────────┘
    */

    /// Retrieve the Stabilization configuration.
    /// @return The Stabilization configuration
    function config() external view returns (Config memory) {
        return _config;
    }

    /// Retrieve the CDP for an account address.
    /// @param owner The CDP account address
    /// @return The CDP object
    function cdps(address owner) external view returns (CDP memory) {
        return _cdps[owner];
    }

    /// Retrieve all the accounts that have opened a CDP.
    /// @return Array of CDP account addresses
    function accounts() external view returns (address[] memory) {
        return _accounts;
    }

    /// Calculate the current debt amount outstanding for a CDP.
    /// @param account The CDP account address
    /// @return debt The debt amount
    function debtAmount(address account) external view returns (uint256 debt) {
        return this.debtAmountAtTime(account, block.timestamp);
    }

    /// Calculate the debt amount outstanding for a CDP at the given timestamp.
    ///
    /// The timestamp must be equal or later than the time of the CDP last
    /// borrow or repayment.
    /// @param account The CDP account address
    /// @param timestamp The timestamp to value the debt
    /// @return debt The debt amount
    function debtAmountAtTime(
        address account,
        uint timestamp
    ) external view goodTime(account, timestamp) returns (uint256 debt) {
        CDP storage cdp = _cdps[account];
        (debt, , ) = _calculateDebtAmount(
            cdp,
            timestamp
        );
    }

    /// Determine if the CDP is currently liquidatable.
    /// @param account The CDP account address
    /// @return Whether the CDP is liquidatable
    function isLiquidatable(address account) external view returns (bool) {
        CDP storage cdp = _cdps[account];
        (uint256 debt, , ) = _calculateDebtAmount(
            cdp,
            block.timestamp
        );
        return
            StabilizationMath.underCollateralized(
            cdp.collateral,
            collateralPrice(),
            debt,
            _config.liquidationRatio
        );
    }

    /// Calculate the maximum amount that can be borrowed against the collateral.
    /// Note that this takes into account the minimum collateralization ratio or
    /// the max borrow limit, whichever is lower will determine the max borrow
    /// @param collateral The amount of Collateral Token
    /// @return The maximum borrow amount
    function maxBorrow(
        uint256 collateral
    ) public view returns (uint256) {
        uint256 borrowLimit = StabilizationMath.borrowLimit(
            collateral,
            collateralPrice(),
            debtPrice(),
            _config.targetPrice,
            _config.minCollateralizationRatio
        );

        uint256 debtLimit = StabilizationMath.debtLimit(
            collateral,
            collateralPrice(),
            _config.liquidationRatio
        );
        return borrowLimit > debtLimit ? debtLimit : borrowLimit;
    }

    /// Price the Collateral Token in Auton.
    ///
    /// Retrieves the Collateral Token price from the Oracle Contract and
    /// converts it to Auton.
    /// @return price Price of Collateral Token
    /// @dev The function reverts in case the price is invalid or unavailable.
    function collateralPrice() public view returns (uint256 price) {
        IOracle.RoundData memory data = _oracle.latestRoundData(StabilizationMath.NTN_SYMBOL);
        if (!data.success) revert PriceUnavailable(StabilizationMath.NTN_SYMBOL);
        if (data.price <= 0) revert InvalidPrice();
        price = data.price;
    }

    /// Price Auton in USD
    ///
    /// Retrieves the Auton price from the Oracle Contract
    /// @return price Price of ATN in ACU
    /// @dev The function reverts in case the price is invalid or unavailable.
    function debtPrice() public view returns (uint256) {
        IOracle.RoundData memory data = _oracle.latestRoundData(StabilizationMath.ATN_SYMBOL);
        if (!data.success) revert PriceUnavailable(StabilizationMath.ATN_SYMBOL);
        if (data.price <= 0) revert InvalidPrice();
        uint256 atnUsd = data.price;
        uint256 acuUsd = acuPrice();
        return atnUsd * StabilizationMath.SCALE_FACTOR / acuUsd;
    }

    /// Price the ACU value in USD.
    ///
    /// Retrieves the ACU value from the ACU Contract and converts it to have
    /// StabilizationMath.SCALE_FACTOR precision.
    /// @return price Price of ACU value
    /// @dev The function reverts in case the price is invalid or unavailable.
    function acuPrice() public view returns (uint256 price) {
        try IACU(_acu).value() returns (int256 acuValue) {
            return StabilizationMath.toScaleFactor(
                uint256(acuValue),
                IACU(_acu).scaleFactor()
            );
        } catch {
            revert PriceUnavailable("ACU-USD");
        }
    }

    /**
     * @notice Get the pending borrow interest rate and since when it will be active.
     * @return uint256 The pending rate
     * @return uint256 The timestamp since it will be active
     */
    function getPendingInterestRateInfo() public view returns (uint256, uint256) {
        return (_pendingBorrowInterestRate, _pendingRateUpdateTimestamp);
    }

    /**
     * @notice Get aggregated interest exponent which is the summation of all interest rate multiplied by their respective time window (in years).
     */
    function getAggregatedInterestExponent() public view returns (uint256) {
        return _calculateAggregatedInterestExponent(block.timestamp);
    }

    /**
     * @notice Get the timestamp since when the current rate is active.
     */
    function getCurrentRateActiveTimestamp() public view returns (uint256) {
        if (_pendingRateUpdateTimestamp > 0 && _pendingRateUpdateTimestamp <= block.timestamp) {
            return _pendingRateUpdateTimestamp;
        }
        return _borrowInterestActiveTimestamp;
    }

    /**
     * @notice Get the active current rate.
     */
    function getCurrentRate() public view returns (uint256) {
        if (_pendingRateUpdateTimestamp > 0 && _pendingRateUpdateTimestamp <= block.timestamp) {
            return _pendingBorrowInterestRate;
        }
        return _config.borrowInterestRate;
    }

    /**
     * @notice Get the pending announcement window and since when it will be active.
     * @return uint256 The pending announcement window
     * @return uint256 The timestamp since the pending announcement window will be active
     */
    function getPendingAnnouncementWindowInfo() public view returns (uint256, uint256) {
        return (_pendingAnnouncementWindow, _pendingAnnouncementUpdateTimestamp);
    }

    /**
     * @notice Get the announcement window in seconds.
     */
    function getAnnouncementWindow() public view returns (uint256) {
        if (_pendingAnnouncementUpdateTimestamp > 0 && _pendingAnnouncementUpdateTimestamp <= block.timestamp) {
            return _pendingAnnouncementWindow;
        }
        return _config.announcementWindow;
    }

    /*
    ┌────────────────┐
    │ Pure Functions │
    └────────────────┘
    */

    // ToDo(scott): figure out the best way to avoid this redundancy
    function borrowLimit(
        uint256 collateral,
        uint256 collateralPrice,
        uint256 debtPrice,
        uint256 targetDebtPrice,
        uint256 mcr
    ) external pure returns (uint256) {
        return StabilizationMath.borrowLimit(
            collateral,
            collateralPrice,
            debtPrice,
            targetDebtPrice,
            mcr
        );
    }

    function minimumCollateral(
        uint256 principal,
        uint256 price,
        uint256 mcr
    ) external pure returns (uint256) {
        return StabilizationMath.minimumCollateral(principal, price, mcr);
    }

    function interestDue(
        uint256 debt,
        uint256 rateExponent
    ) external pure returns (uint256) {
        return StabilizationMath.interestDue(debt, rateExponent);
    }

    function underCollateralized(
        uint256 collateral,
        uint256 price,
        uint256 debt,
        uint256 liquidationRatio
    ) external pure returns (bool) {
        return StabilizationMath.underCollateralized(collateral, price, debt, liquidationRatio);
    }

    function interestExponent(
        uint256 interestRate,
        uint256 startTimestamp,
        uint256 endTimestamp
    ) external pure returns (uint256) {
        return StabilizationMath.interestExponent(interestRate, startTimestamp, endTimestamp);
    }


    /*
    ┌────────────────────┐
    │ Internal Functions │
    └────────────────────┘
    */

    function _calculateDebtAmount(
        CDP storage cdp,
        uint timestamp
    ) internal view returns (uint256 total, uint256 accrued, uint256 totalExponent) {
        if (timestamp == 0) revert InvalidParameter();
        uint256 debt = cdp.principal + cdp.interest;
        totalExponent = _calculateAggregatedInterestExponent(timestamp);
        if (debt == 0) {
            return (0, 0, totalExponent);
        }
        if (timestamp == cdp.timestamp) accrued = 0;
        else {
            accrued = StabilizationMath.interestDue(
                debt,
                totalExponent - cdp.lastAggregatedInterestExponent
            );
        }
        total = debt + accrued;
    }

    function _updateDebt(
        CDP storage cdp,
        uint timestamp
    ) internal returns (uint256 total) {
        uint256 accrued;
        (total, accrued, cdp.lastAggregatedInterestExponent) = _calculateDebtAmount(
            cdp,
            timestamp
        );
        cdp.interest += accrued;
        cdp.timestamp = timestamp;
    }

    function _allocatePayment(
        CDP storage cdp,
        uint256 amount
    )
    internal
    view
    returns (uint256 interest, uint256 principal, uint256 surplus)
    {
        uint256 debt = cdp.principal + cdp.interest;
        interest = amount < cdp.interest ? amount : cdp.interest;
        principal = amount < debt ? amount - interest : cdp.principal;
        surplus = amount > debt ? amount - debt : 0;
    }

    function _updateAnnouncementWindow() internal {
        if (_pendingAnnouncementUpdateTimestamp > 0 && _pendingAnnouncementUpdateTimestamp <= block.timestamp) {
            _config.announcementWindow = _pendingAnnouncementWindow;
            _pendingAnnouncementWindow = 0;
            _pendingAnnouncementUpdateTimestamp = 0;
        }
    }

    function _applyInterestRateUpdate() internal {
        uint256 pendingRateActiveTimestamp = _pendingRateUpdateTimestamp;
        if (pendingRateActiveTimestamp == 0 || pendingRateActiveTimestamp > block.timestamp) {
            return;
        }
        _aggregatedInterestExponent += StabilizationMath.interestExponent(
            _config.borrowInterestRate,
            _borrowInterestActiveTimestamp,
            pendingRateActiveTimestamp
        );
        _config.borrowInterestRate = _pendingBorrowInterestRate;
        _borrowInterestActiveTimestamp = pendingRateActiveTimestamp;
        // clear some storage
        _pendingRateUpdateTimestamp = 0;
        _pendingBorrowInterestRate = 0;
    }

    /**
     * @dev Calculates total aggregated interest exponent which is the summation of all interest rates multiplied
     * by their respective time window until `timestamp`.
     * @return aggregatedInterestExponent aggregated interest exponent
     */
    function _calculateAggregatedInterestExponent(uint256 timestamp) internal view returns (uint256) {
        uint256 aggregatedInterestExponent = _aggregatedInterestExponent;
        uint256 currentRate = _config.borrowInterestRate;
        uint256 currentRateActiveTimestamp = _borrowInterestActiveTimestamp;
        // the following condition is enforces because `_aggregatedInterestExponent` state
        // variable stores the aggregated interest until `_borrowInterestActiveTimestamp`
        if (timestamp < currentRateActiveTimestamp) revert InvalidParameter();

        uint256 pendingRateActiveTimestamp = _pendingRateUpdateTimestamp;
        if (pendingRateActiveTimestamp > 0 && pendingRateActiveTimestamp <= timestamp) {
            // add the `currentRate` multiplied by its time window to the aggregation
            aggregatedInterestExponent += StabilizationMath.interestExponent(
                currentRate,
                currentRateActiveTimestamp,
                pendingRateActiveTimestamp
            );
            // update the `currentRate`
            currentRate = _pendingBorrowInterestRate;
            currentRateActiveTimestamp = pendingRateActiveTimestamp;
        }

        return aggregatedInterestExponent + StabilizationMath.interestExponent(
            currentRate,
            currentRateActiveTimestamp,
            timestamp
        );
    }
}
/* solhint-enable not-rely-on-time */
