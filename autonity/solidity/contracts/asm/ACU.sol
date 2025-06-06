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

import {IACU} from "./interfaces/IACU.sol";
import {IOracle} from "../interfaces/IOracle.sol";
import "./lib/ASMErrors.sol";
import {IConfigEvents} from "../interfaces/IConfigEvents.sol";
import {ReentrancyGuard} from "../ReentrancyGuard.sol";

/// @title ASM ACU Contract
/// @notice Computes the value of the ACU, an optimal currency basket of
/// 7 free-floating fiat currencies.
/// @dev Intended to be deployed by the protocol at genesis.
contract ACU is IACU, IConfigEvents, ReentrancyGuard {
    /// The Oracle round of the current ACU value.
    uint256 internal round;
    /// The decimal places used to represent the ACU as a fixed-point integer.
    /// It is also the scale used to represent the basket quantities.
    uint256 internal scale;
    /// The multiplier for scaling numbers to the ACU scaled representation.
    uint256 internal scaleFactor;
    /// The quantity multiplier for the ACU basket.
    uint256 internal quantityMultiplier;

    string[] private _symbols;
    uint256[] private _quantities;
    uint256 private _value;
    address private _autonity;
    address private _operator;
    IOracle private _oracle;
    bytes32 private constant SYMBOL_USD = keccak256(abi.encodePacked("USD-USD"));

    /// The ACU value was updated.
    event Updated(uint height, uint timestamp, uint256 round, uint256 value);
    /// The ACU symbols, quantites, or scale were modified.
    event BasketModified(string[] symbols, uint256[] quantities, uint256 scale);
    /// The ACU quantity multiplier has been updated
    event Rescaled(uint256 newQuantityMultiplier, uint256 oldQuantityMultiplier);

    modifier onlyAutonity() {
        if (msg.sender != _autonity) revert Unauthorized();
        _;
    }

    modifier onlyOperator() {
        if (msg.sender != _operator) revert Unauthorized();
        _;
    }

    modifier validBasket(
        string[] memory symbols_,
        uint256[] memory quantities_
    ) {
        if (symbols_.length != quantities_.length) revert InvalidBasket();
        for (uint i = 0; i < quantities_.length; i++) {
            if (quantities_[i] > uint256(type(int256).max))
                revert InvalidBasket();
        }
        _;
    }

    /// Create and deploy the ASM ACU Contract.
    /// @param symbols_ The symbols used to retrieve prices
    /// @param quantities_ The basket quantity corresponding to each symbol
    /// @param scale_ The scale for quantities and the ACU value
    /// @param autonity Address of the Autonity Contract
    /// @param operator Address of the Governance Operator
    /// @param oracle Address of the Oracle Contract
    constructor(
        string[] memory symbols_,
        uint256[] memory quantities_,
        uint256 scale_,
        address autonity,
        address operator,
        address oracle
    ) validBasket(symbols_, quantities_) {
        _symbols = symbols_;
        _quantities = quantities_;
        scale = scale_;
        scaleFactor = 10 ** scale_;
        quantityMultiplier = scaleFactor;
        _autonity = autonity;
        _operator = operator;
        _oracle = IOracle(oracle);
    }

    /*
    ┌────────────────────┐
    │ Autonity Functions │
    └────────────────────┘
    */

    /// Compute the ACU value and store it.
    ///
    /// It retrieves the latest prices from the Oracle Contract. If one or
    /// more prices are unavailable from the Oracle, it will not compute the
    /// value for that round.
    ///
    /// This function is intended to be called by the protocol during block
    /// finalization, after the Oracle Contract finalization has completed.
    /// @return status Whether the ACU value was updated successfully
    /// @dev Only the Autonity Contract is authorized to trigger the
    /// computation of the ACU.
    function update() external onlyAutonity nonReentrant returns (bool status) {
        uint256 latestRound = _oracle.getRound() - 1;
        if (round >= latestRound) return false;
        uint256 sumProduct = 0;
        uint256 oracleDecimals = uint256(_oracle.getDecimals());
        for (uint i = 0; i < _symbols.length; i++) {
            uint256 price;
            if (keccak256(abi.encodePacked(_symbols[i])) == SYMBOL_USD) {
                price = 10 ** oracleDecimals;
            } else {
                IOracle.RoundData memory roundData = _oracle.getRoundData(
                    latestRound,
                    _symbols[i]
                );
                if (!roundData.success) return false;
                price = roundData.price;
            }
            sumProduct += (price * _quantities[i]);
        }

        _value = sumProduct / 10 ** oracleDecimals;
        round = latestRound;

        // solhint-disable-next-line not-rely-on-time
        emit Updated(block.number, block.timestamp, round, _value);
        return true;
    }

    /// Set the Governance Operator account address.
    /// @param operator Address of the new Governance Operator
    /// @dev Only the Autonity Contract is authorized to set the Governance
    /// Operator account address.
    function setOperator(address operator) external onlyAutonity {
        emit IConfigEvents.ConfigUpdateAddress("operator", _operator, operator, block.number);
        _operator = operator;
    }

    /// Set the Oracle Contract address that is used to retrieve prices.
    /// @param oracle Address of the new Oracle Contract
    /// @dev Only the Autonity Contract is authorized to set the Oracle
    /// Contract address.
    function setOracle(address oracle) external onlyAutonity {
        emit IConfigEvents.ConfigUpdateAddress("oracle", address(_oracle), oracle, block.number);
        _oracle = IOracle(oracle);
    }

    /*
    ┌────────────────────┐
    │ Operator Functions │
    └────────────────────┘
    */

    /// Modify the ACU symbols, quantites, or scale.
    /// @param symbols_ The symbols used to retrieve prices
    /// @param quantities_ The basket quantity corresponding to each symbol
    /// @param scale_ The scale for quantities and the ACU value
    /// @dev Only the operator is authorized to modify the basket.
    function modifyBasket(
        string[] memory symbols_,
        uint256[] memory quantities_,
        uint256 scale_
    ) external validBasket(symbols_, quantities_) onlyOperator {
        _symbols = symbols_;
        _quantities = quantities_;
        scale = scale_;
        // Rescale the quantity multiplier to the new scaleFactor
        quantityMultiplier = quantityMultiplier * (10 ** scale_) / scaleFactor;
        scaleFactor = 10 ** scale;
        emit BasketModified(symbols_, quantities_, scale_);
    }

    // Rescale the quantity multiplier.
    /// @param newQuantityMultiplier The new quantity multiplier
    /// @notice the quantity multiplier has precision of scaleFactor
    function rescale(uint256 newQuantityMultiplier) external onlyOperator {
        uint256 oldQuantityMultiplier = quantityMultiplier;
        if (newQuantityMultiplier == 0) revert ZeroValue();
        quantityMultiplier = newQuantityMultiplier;
        emit Rescaled(newQuantityMultiplier, oldQuantityMultiplier);
    }

    /*
    ┌────────────────┐
    │ View Functions │
    └────────────────┘
    */

    /// The latest ACU value that was computed.
    /// @return ACU value in fixed-point integer representation rescaled by the
    /// quantity multiplier
    function value() external view nonReentrantView returns (uint256) {
        if (round == 0) revert NoACUValue();
        return quantityMultiplier * _value / scaleFactor;
    }

    /// The symbols that are used to compute the ACU.
    /// @return Array of symbols
    function symbols() external view returns (string[] memory) {
        return _symbols;
    }

    /// The basket quantities that are used to compute the ACU.
    /// @return Array of quantities
    function quantities() external view returns (uint256[] memory) {
        return _quantities;
    }

    /// The quantity multiplier that is used to compute the ACU.
    /// @return Quantity multiplier
    /// @dev The quantity multiplier has precision of scaleFactor
    function multiplier() external view returns (uint256) {
        return quantityMultiplier;
    }

    /// The scaled quantities used to compute the ACU.
    /// @return Array of scaled quantities
    function scaledQuantities() external view returns (uint256[] memory) {
        uint256[] memory scaled = new uint256[](_quantities.length);
        for (uint i = 0; i < _quantities.length; i++) {
            scaled[i] = _quantities[i] * quantityMultiplier / scaleFactor;
        }
        return scaled;
    }

    /// @return The multiplier for scaling numbers to the ACU scaled representation.
    function getScaleFactor() external view returns (uint256) {
        return scaleFactor;
    }

    /// @return The decimal places used to represent the ACU as a fixed-point integer.
    function getScale() external view returns (uint256) {
        return scale;
    }

    /// @return The Oracle round of the current ACU value.
    function getRound() external view nonReentrantView returns (uint256) {
        return round;
    }
}
