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

import {IACU} from "./IACU.sol";
import {IOracle} from "../interfaces/IOracle.sol";

/// @title ASM ACU Contract
/// @notice Computes the value of the ACU, an optimal currency basket of
/// 7 free-floating fiat currencies.
/// @dev Intended to be deployed by the protocol at genesis.
contract ACU is IACU {
    /// The Oracle round of the current ACU value.
    uint256 public round;
    /// The decimal places used to represent the ACU as a fixed-point integer.
    /// It is also the scale used to represent the basket quantities.
    uint256 public scale;
    /// The multiplier for scaling numbers to the ACU scaled representation.
    uint256 public scaleFactor;

    string[] private _symbols;
    uint256[] private _quantities;
    int256 private _value;
    address private _autonity;
    address private _operator;
    IOracle private _oracle;
    bytes32 private constant SYMBOL_USD =
        keccak256(abi.encodePacked("USD/USD"));

    /// The ACU value was updated.
    event Updated(uint height, uint timestamp, uint256 round, int256 value);
    /// The ACU symbols, quantites, or scale were modified.
    event BasketModified(string[] symbols, uint256[] quantities, uint256 scale);

    error InvalidBasket();
    error NoACUValue();
    error Unauthorized();

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
    function update() external onlyAutonity returns (bool status) {
        uint256 latestRound = _oracle.getRound() - 1;
        if (round >= latestRound) return false;
        int256 sumProduct = 0;
        int256 oraclePrecision = int256(_oracle.getPrecision());
        for (uint i = 0; i < _symbols.length; i++) {
            int256 price;
            if (keccak256(abi.encodePacked(_symbols[i])) == SYMBOL_USD) {
                price = oraclePrecision;
            } else {
                IOracle.RoundData memory roundData = _oracle.getRoundData(
                    latestRound,
                    _symbols[i]
                );
                if (roundData.status != 0) return false;
                price = roundData.price;
            }
            sumProduct += (price * int256(_quantities[i]));
        }

        _value = sumProduct / oraclePrecision;
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
        _operator = operator;
    }

    /// Set the Oracle Contract address that is used to retrieve prices.
    /// @param oracle Address of the new Oracle Contract
    /// @dev Only the Autonity Contract is authorized to set the Oracle
    /// Contract address.
    function setOracle(address oracle) external onlyAutonity {
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
        scaleFactor = 10 ** scale;
        emit BasketModified(symbols_, quantities_, scale_);
    }

    /*
    ┌────────────────┐
    │ View Functions │
    └────────────────┘
    */

    /// The latest ACU value that was computed.
    /// @return ACU value in fixed-point integer representation
    function value() external view returns (int256) {
        if (round == 0) revert NoACUValue();
        return _value;
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
}
