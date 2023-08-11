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

import {IOracle} from "../interfaces/IOracle.sol";
import {IAutonity} from "../interfaces/IAutonity.sol";
import {IACU} from "./IACU.sol";

/// @title ASM ACU Contract
/// @notice Computes the value of the ACU, an optimal currency basket of
/// 7 free-floating fiat currencies.
/// @dev Intended to be deployed by the protocol at genesis.
contract ACU is IACU {
    bytes32 private constant SYMBOL_USD = keccak256(abi.encodePacked("USD/USD"));
    /// The Oracle round of the current ACU value.
    uint256 public round;
    /// The decimal places used to represent the ACU as a fixed-point integer.
    /// It is also the scale used to represent the basket quantities.
    uint256 public scale;
    /// The multiplier for scaling numbers to the ACU scaled representation.
    uint256 public scaleFactor;

    string[] public symbols;
    uint256[] public quantities;

    /// The resulting computed ACU value
    /// The initial value will be set to 0
    int256 public value;

    IOracle private oracle;
    IAutonity private autonity;

    /// The ACU value was updated.
    event Updated(uint height, uint timestamp, uint256 round, int256 value);
    /// The ACU symbols, quantites, or scale were modified.
    event BasketModified(string[] symbols, uint256[] quantities, uint256 scale);

    error InvalidBasket();
    error Unauthorized();


    /// Create and deploy the ASM ACU Contract.
    /// @param _symbols The symbols used to retrieve prices
    /// @param _quantities The basket quantity corresponding to each symbol
    /// @param _scale The scale for quantities and the ACU value
    /// @param _autonity Address of the Autonity Contract
    constructor(
        string[] memory _symbols,
        uint256[] memory _quantities,
        uint256 _scale,
        address payable _autonity
    ) validBasket(_symbols, _quantities) {
        symbols = _symbols;
        quantities = _quantities;
        scale = _scale;
        scaleFactor = 10 ** _scale;
        autonity = IAutonity(_autonity);
        oracle = IOracle(autonity.getOracle());
    }

    /*
    ┌────────────────────┐
    │ Operator Functions │
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
    /// @dev Only the operator is authorized to trigger the computation of ACU.
    function update() external onlyAutonity returns (bool status) {
        uint256 _latestRound = oracle.getRound() - 1;
        if (round >= _latestRound) return false;
        int256 _sumProduct = 0;
        int256 _oraclePrecision = int256(oracle.getPrecision());
        for (uint i = 0; i < symbols.length; i++) {
            int256 _price;
            if (keccak256(abi.encodePacked(symbols[i])) == SYMBOL_USD) {
                _price = _oraclePrecision;
            } else {
                IOracle.RoundData memory roundData = oracle.getRoundData(
                    _latestRound,
                    symbols[i]
                );
                if (roundData.status != 0) return false;
                _price = roundData.price;
            }
            _sumProduct += (_price * int256(quantities[i]));
        }

        value = _sumProduct / _oraclePrecision;
        round = _latestRound;

        // solhint-disable-next-line not-rely-on-time
        emit Updated(block.number, block.timestamp, round, value);
        return true;
    }

    /// Modify the ACU symbols, quantites, or scale.
    /// @param _symbols The symbols used to retrieve prices
    /// @param _quantities The basket quantity corresponding to each symbol
    /// @param _scale The scale for quantities and the ACU value
    /// @dev Only the operator is authorized to modify the basket.
    function modifyBasket(
        string[] memory _symbols,
        uint256[] memory _quantities,
        uint256 _scale
    ) external validBasket(_symbols, _quantities) onlyOperator {
        symbols = _symbols;
        quantities = _quantities;
        scale = _scale;
        scaleFactor = 10 ** scale;
        emit BasketModified(_symbols, _quantities, _scale);
    }

    /// Set the Oracle Contract address that is used to retrieve prices.
    /// @param _oracle Address of the new Oracle Contract
    /// @dev Only the autonity contract is authorized to set the Oracle Contract address.
    function setOracle(address payable _oracle) external onlyAutonity {
        oracle = IOracle(_oracle);
    }

    /*
    ┌────────────────┐
    │ Modifiers       │
    └────────────────┘
    */


    modifier onlyAutonity() {
        if (msg.sender != address(autonity)) revert Unauthorized();
        _;
    }

    modifier onlyOperator() {
        if (msg.sender != autonity.getOperator()) revert Unauthorized();
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
}
