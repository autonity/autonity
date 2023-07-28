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
import {ISupplyControl} from "./ISupplyControl.sol";
import {UD60x18, ud} from "../lib/prb-math-4.0.1/UD60x18.sol";

/// Stabilization Configuration.
struct Config {
    /// The annual continuously-compounded interest rate for borrowing.
    uint256 borrowInterestRate;
    /// The minimum ACU value of collateral required to maintain 1 ACU value of
    /// debt.
    uint256 liquidationRatio;
    /// The minimum ACU value of collateral required to borrow 1 ACU value of
    /// debt.
    uint256 minCollateralizationRatio;
    /// The minimum amount of debt required to maintain a CDP.
    uint256 minDebtRequirement;
    /// The ACU value of 1 unit of debt.
    uint256 redemptionPrice;
}

/// @title ASM Stabilization Contract
/// @notice A CDP-based stabilization mechanism for the Auton.
/// @dev Intended to be deployed by the protocol at genesis. Note that all
/// rates, ratios, prices, and amounts are represented as fixed-point integers
/// with `SCALE` decimal places.
/* solhint-disable not-rely-on-time */
contract Stabilization {
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

    /// The decimal places in fixed-point integer representation.
    uint256 public constant SCALE = 18; // Match UD60x18
    /// The multiplier for scaling numbers to the required scale.
    uint256 public constant SCALE_FACTOR = 10 ** SCALE;
    /// A year is assumed to have 365 days for interest rate calculations.
    uint256 public constant SECONDS_IN_YEAR = 365 days;
    /// The Config object that stores Stabilization Contract parameters.
    Config public config;
    /// A mapping to retrieve the CDP for an account address.
    mapping(address => CDP) public cdps;

    string private constant NTN_SYMBOL = "NTN/ATN";
    address[] private _accounts;
    address private _operator;
    IERC20 private _collateralToken;
    IOracle private _oracle;
    ISupplyControl private _supplyControl;

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

    error InsufficientAllowance();
    error InsufficientPayment();
    error InsufficientCollateral();
    error InvalidDebtPosition();
    error InvalidAmount();
    error InvalidParameter();
    error InvalidPrice();
    error Liquidatable();
    error NotLiquidatable();
    error NoDebtPosition();
    error PriceUnavailable();
    error TransferFailed();
    error Unauthorized();
    error ZeroValue();

    modifier goodTime(address account, uint timestamp) {
        CDP storage cdp = cdps[account];
        if (timestamp < cdp.timestamp) revert InvalidParameter();
        _;
    }

    modifier nonZeroAmount(uint256 amount) {
        if (amount == 0) revert InvalidAmount();
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

    /// Create and deploy the ASM Stabilization Contract.
    /// @param config_ Stabilization configuration
    /// @param operator The operator is authorized to change parameters
    /// @param oracle Address of the Oracle Contract
    /// @param supplyControl Address of the SupplyControl Contract
    /// @param collateralToken Address of the Collateral Token contract
    constructor(
        Config memory config_,
        address operator,
        address oracle,
        address supplyControl,
        IERC20 collateralToken
    )
        positiveMCR(config_.minCollateralizationRatio)
        validRatios(config_.liquidationRatio, config_.minCollateralizationRatio)
    {
        config = config_;
        _collateralToken = collateralToken;
        _operator = operator;
        _oracle = IOracle(oracle);
        _supplyControl = ISupplyControl(supplyControl);
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
    function deposit(uint256 amount) external nonZeroAmount(amount) {
        if (_collateralToken.allowance(msg.sender, address(this)) < amount)
            revert InsufficientAllowance();

        CDP storage cdp = cdps[msg.sender];
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
    function withdraw(uint256 amount) external nonZeroAmount(amount) {
        CDP storage cdp = cdps[msg.sender];
        if (amount > cdp.collateral) revert InvalidAmount();
        (uint256 debt, ) = _debtAmount(cdp, block.timestamp);
        uint256 price = collateralPrice();
        if (
            underCollateralized(
                cdp.collateral,
                price,
                debt,
                config.liquidationRatio
            )
        ) revert Liquidatable();
        if (
            cdp.collateral - amount <
            minimumCollateral(
                cdp.principal,
                price,
                config.minCollateralizationRatio
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
    function borrow(uint256 amount) external nonZeroAmount(amount) {
        CDP storage cdp = cdps[msg.sender];
        (uint256 debt, uint256 accrued) = _debtAmount(cdp, block.timestamp);
        debt += amount;
        if (debt < config.minDebtRequirement) revert InvalidDebtPosition();
        uint256 price = collateralPrice();
        if (
            underCollateralized(
                cdp.collateral,
                price,
                debt,
                config.liquidationRatio
            )
        ) revert Liquidatable();
        uint256 limit = borrowLimit(
            cdp.collateral,
            price,
            config.redemptionPrice,
            config.minCollateralizationRatio
        );
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
    function repay() external payable {
        if (msg.value == 0) revert ZeroValue();
        CDP storage cdp = cdps[msg.sender];
        if (cdp.principal == 0) revert NoDebtPosition();
        (uint256 debt, uint256 accrued) = _debtAmount(cdp, block.timestamp);
        if (
            (msg.value < debt) && (debt - msg.value < config.minDebtRequirement)
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
    function liquidate(address account) external payable {
        if (msg.value == 0) revert ZeroValue();
        CDP storage cdp = cdps[account];
        if (cdp.principal == 0) revert NoDebtPosition();
        (uint256 debt, uint256 accrued) = _debtAmount(cdp, block.timestamp);
        if (
            !underCollateralized(
                cdp.collateral,
                collateralPrice(),
                debt,
                config.liquidationRatio
            )
        ) revert NotLiquidatable();
        uint surplus = msg.value - debt;
        if (surplus < 0) revert InsufficientPayment();

        uint256 collateral = cdp.collateral;
        cdp.timestamp = block.timestamp;
        cdp.collateral = 0;
        cdp.principal = 0;
        cdp.interest = 0;

        if (!_collateralToken.transfer(msg.sender, collateral))
            revert TransferFailed();
        _supplyControl.burn{value: debt - accrued}();
        if (surplus > 0) payable(msg.sender).transfer(surplus);
        emit Liquidate(account, msg.sender);
    }

    /*
    ┌────────────────────┐
    │ Operator Functions │
    └────────────────────┘
    */

    /// Set the liquidation ratio.
    ///
    /// Must be less than the minimum collateralization ratio.
    /// @dev Restricted to the operator.
    function setLiquidationRatio(
        uint256 ratio
    )
        external
        validRatios(ratio, config.minCollateralizationRatio)
        onlyOperator
    {
        config.liquidationRatio = ratio;
    }

    /// Set the minimum collateralization ratio.
    ///
    /// Must be positive and greater than the liquidation ratio.
    /// @dev Restricted to the operator.
    function setMinCollateralizationRatio(
        uint256 ratio
    )
        external
        positiveMCR(ratio)
        validRatios(config.liquidationRatio, ratio)
        onlyOperator
    {
        config.minCollateralizationRatio = ratio;
    }

    /// Set the minimum debt requirement.
    /// @dev Restricted to the operator.
    function setMinDebtRequirement(uint256 amount) external onlyOperator {
        config.minDebtRequirement = amount;
    }

    /// Set the Oracle Contract address.
    /// @dev Restricted to the operator.
    function setOracle(address oracle) external onlyOperator {
        _oracle = IOracle(oracle);
    }

    /// Set the SupplyControl Contract address.
    /// @dev Restricted to the operator.
    function setSupplyControl(address supplyControl) external onlyOperator {
        _supplyControl = ISupplyControl(supplyControl);
    }

    /*
    ┌────────────────┐
    │ View Functions │
    └────────────────┘
    */

    /// Retrieve all the accounts that have opened a CDP.
    /// @return Array of CDP account addresses
    function accounts() external view returns (address[] memory) {
        return _accounts;
    }

    /// Calculate the current debt amount outstanding for a CDP.
    /// @return debt The debt amount
    function debtAmount(address account) external view returns (uint256 debt) {
        return this.debtAmount(account, block.timestamp);
    }

    /// Calculate the debt amount outstanding for a CDP at the given timestamp.
    ///
    /// The timestamp must be equal or later than the time of the CDP last
    /// borrow or repayment.
    /// @return debt The debt amount
    function debtAmount(
        address account,
        uint timestamp
    ) external view goodTime(account, timestamp) returns (uint256 debt) {
        CDP storage cdp = cdps[account];
        (debt, ) = _debtAmount(cdp, timestamp);
    }

    /// Determine if the CDP is currently liquidatable.
    /// @return Whether the CDP is liquidatable
    function isLiquidatable(address account) external view returns (bool) {
        return this.isLiquidatable(account, block.timestamp);
    }

    /// Determine if the CDP is liquidatable at the given timestamp.
    ///
    /// The timestamp must be equal or later than the time of the CDP last
    /// borrow or repayment.
    /// @return Whether the CDP is liquidatable
    function isLiquidatable(
        address account,
        uint timestamp
    ) external view goodTime(account, timestamp) returns (bool) {
        CDP storage cdp = cdps[account];
        (uint256 debt, ) = _debtAmount(cdp, timestamp);
        return
            underCollateralized(
                cdp.collateral,
                collateralPrice(),
                debt,
                config.liquidationRatio
            );
    }

    /// Price the Collateral Token in Auton.
    ///
    /// Retrieves the Collateral Token price from the Oracle Contract and
    /// converts it to Auton.
    /// @return price Price of Collateral Token
    /// @dev The function reverts in case the price is invalid or unavailable.
    function collateralPrice() public view returns (uint256 price) {
        IOracle.RoundData memory data = _oracle.latestRoundData(NTN_SYMBOL);
        if (data.status != 0) revert PriceUnavailable();
        if (data.price <= 0) revert InvalidPrice();
        // Convert price from Oracle precision to SCALE decimals
        if (SCALE_FACTOR > _oracle.getPrecision())
            price =
                uint256(data.price) *
                (SCALE_FACTOR / _oracle.getPrecision());
        else
            price =
                uint256(data.price) /
                (_oracle.getPrecision() / SCALE_FACTOR);
    }

    /*
    ┌──────────────┐
    │ Calculations │
    └──────────────┘
    */

    /// Calculate the maximum amount of Amount that can be borrowed for the
    /// given amount of Collateral Token.
    /// @param collateral Amount of Collateral Token backing the debt
    /// @param price The price of Collateral Token in Auton
    /// @param redemptionPrice The ACU value of 1 unit of debt
    /// @param mcr The minimum collateralization ratio
    /// @return The maximum Auton that can be borrowed
    function borrowLimit(
        uint256 collateral,
        uint256 price,
        uint256 redemptionPrice,
        uint256 mcr
    ) public pure returns (uint256) {
        if (price == 0 || mcr == 0) revert InvalidParameter();
        return (collateral * price * redemptionPrice) / (mcr * SCALE_FACTOR);
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
    ) public pure validPrice(price) returns (uint256) {
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
    ) public pure returns (uint256) {
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
    ) public pure validPrice(price) returns (bool) {
        if (debt == 0) return false;
        return (collateral * price) / debt < liquidationRatio;
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
            accrued = interestDue(
                debt,
                config.borrowInterestRate,
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
