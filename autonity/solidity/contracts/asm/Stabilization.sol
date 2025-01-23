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
        (uint256 debt,) = _debtAmount(cdp, block.timestamp);
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
        (uint256 debt, uint256 accrued) = _debtAmount(cdp, block.timestamp);
        debt += amount;
        if (debt < _config.minDebtRequirement) revert InvalidDebtPosition();
        uint256 price = collateralPrice();
        if (
            StabilizationMath.underCollateralized(
            cdp.collateral,
            price,
            debt,
            _config.liquidationRatio)
        ) revert Liquidatable();

        uint256 limit = maxBorrow(cdp.collateral);
        if (debt > limit) revert InsufficientCollateral();

        cdp.timestamp = block.timestamp;
        cdp.principal += amount;
        cdp.interest += accrued;

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
        (uint256 debt, uint256 accrued) = _debtAmount(cdp, block.timestamp);
        if (
            (msg.value < debt) && (debt - msg.value < _config.minDebtRequirement)
        ) revert InvalidDebtPosition();

        cdp.interest += accrued;
        cdp.timestamp = block.timestamp;
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
        (uint256 debt, uint256 accrued) = _debtAmount(cdp, block.timestamp);
        if (
            !StabilizationMath.underCollateralized(
            cdp.collateral,
            collateralPrice(),
            debt,
            _config.liquidationRatio
        )
        ) revert NotLiquidatable();

        if (msg.value < debt) revert InsufficientPayment();
        uint surplus = msg.value - debt;

        uint256 collateral = cdp.collateral;
        cdp.timestamp = block.timestamp;
        cdp.collateral = collateral - collateralSold;
        cdp.principal = 0;
        cdp.interest = 0;

        if (!_collateralToken.transfer(bidder, collateralSold))
            revert TransferFailed();
        _supplyControl.burn{value: debt - accrued}();
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
    }

    /// Set the minimum debt requirement.
    /// @param amount The minimum debt amount
    /// @dev Restricted to the operator.
    function setMinDebtRequirement(uint256 amount) external onlyOperator {
        _config.minDebtRequirement = amount;
    }

    /// Set the SupplyControl Contract address.
    /// @param supplyControl The SupplyControl Contract address
    /// @dev Restricted to the operator.
    function setSupplyControl(address supplyControl) external onlyOperator {
        _supplyControl = ISupplyControl(supplyControl);
    }

    /// Set the _atnSupplyOperator address.
    /// @param atnSupplyOperator The _atnSupplyOperator address
    /// @dev Restricted to the operator.
    function setAtnSupplyOperator(address atnSupplyOperator) external onlyOperator {
        _atnSupplyOperator = atnSupplyOperator;
    }

    /// Transition out of the restricted state.
    /// @dev Restricted to the operator.
    function removeCDPRestrictions() external onlyOperator {
        _restricted = false;
        _config.borrowInterestRate = _defaultGenesisBorrowInterestRate;
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
        _operator = operator;
    }

    /// Set the Oracle Contract address.
    /// @param oracle Address of the new Oracle Contract
    /// @dev Restricted to the Autonity Contract.
    function setOracle(address oracle) external onlyAutonity {
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
        return this.debtAmount(account, block.timestamp);
    }

    /// Calculate the debt amount outstanding for a CDP at the given timestamp.
    ///
    /// The timestamp must be equal or later than the time of the CDP last
    /// borrow or repayment.
    /// @param account The CDP account address
    /// @param timestamp The timestamp to value the debt
    /// @return debt The debt amount
    function debtAmount(
        address account,
        uint timestamp
    ) external view goodTime(account, timestamp) returns (uint256 debt) {
        CDP storage cdp = _cdps[account];
        (debt,) = _debtAmount(cdp, timestamp);
    }

    /// Determine if the CDP is currently liquidatable.
    /// @param account The CDP account address
    /// @return Whether the CDP is liquidatable
    function isLiquidatable(address account) external view returns (bool) {
        CDP storage cdp = _cdps[account];
        (uint256 debt,) = _debtAmount(cdp, block.timestamp);
        return
            StabilizationMath.underCollateralized(
            cdp.collateral,
            collateralPrice(),
            debt,
            _config.liquidationRatio
        );
    }

    function maxBorrow(
        uint256 collateral
    ) public view returns (uint256) {
        // oracle prices are all 18 decimals, but the acu value is scaled
        // independently

        uint256 borrowLimit = StabilizationMath.borrowLimit(
            collateral,
            collateralPrice(),
            debtPrice(),
            _config.targetPrice,
            acuPrice(),
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
    /// Retrieves the Auton price from the Oracle Contract and
    /// converts it to Auton.
    /// @return price Price of Auton in USD
    /// @dev The function reverts in case the price is invalid or unavailable.
    function debtPrice() public view returns (uint256 price) {
        IOracle.RoundData memory data = _oracle.latestRoundData(StabilizationMath.ATN_SYMBOL);
        if (!data.success) revert PriceUnavailable(StabilizationMath.NTN_SYMBOL);
        if (data.price <= 0) revert InvalidPrice();
        price = data.price;
    }

    function acuPrice() public view returns (uint256 price) {
        try IACU(_acu).value() returns (int256 acuValue) {
            return StabilizationMath.toScaleFactor(
                uint256(acuValue),
                IACU(_acu).scaleFactor()
            );
        } catch {
            revert PriceUnavailable("ACU");
        }
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
        uint256 acuPrice,
        uint256 mcr
    ) external pure returns (uint256) {
        return StabilizationMath.borrowLimit(
            collateral,
            collateralPrice,
            debtPrice,
            targetDebtPrice,
            acuPrice,
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
        uint256 rate,
        uint timeBorrow,
        uint timeDue
    ) external pure returns (uint256) {
        return StabilizationMath.interestDue(debt, rate, timeBorrow, timeDue);
    }

    function underCollateralized(
        uint256 collateral,
        uint256 price,
        uint256 debt,
        uint256 liquidationRatio
    ) external pure returns (bool) {
        return StabilizationMath.underCollateralized(collateral, price, debt, liquidationRatio);
    }

    /*
    ┌────────────────────┐
    │ Internal Functions │
    └────────────────────┘
    */

    function _debtAmount(
        CDP storage cdp,
        uint timestamp
    ) internal view returns (uint256 total, uint256 accrued) {
        if (timestamp == 0) revert InvalidParameter();
        uint256 debt = cdp.principal + cdp.interest;
        if (timestamp == cdp.timestamp) accrued = 0;
        else {
            accrued = StabilizationMath.interestDue(
                debt,
                _config.borrowInterestRate,
                cdp.timestamp,
                timestamp
            );
        }
        total = debt + accrued;
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
}
/* solhint-enable not-rely-on-time */
