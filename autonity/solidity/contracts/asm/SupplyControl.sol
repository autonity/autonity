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

import {ISupplyControl} from "./ISupplyControl.sol";

/// @title ASM Supply Control Contract Implementation
/// @notice Controls the supply of Auton on the network.
/// @dev Intended to be deployed by the protocol at genesis. The stabilizer is
/// expected to be the Stabilization Contract.
contract SupplyControl is ISupplyControl {
    /// The account that is authorized to mint and burn.
    address public stabilizer;

    /// The total supply of Auton under management.
    uint256 public totalSupply;

    /// The Autonity Contract address.
    address private _autonity;

    /// The Governance Operator account address.
    address private _operator;

    error InvalidAmount();
    error InvalidRecipient();
    error Unauthorized();
    error ZeroValue();

    modifier nonZeroValue() {
        if (msg.value == 0) revert ZeroValue();
        _;
    }

    modifier onlyAutonity() {
        if (msg.sender != _autonity) revert Unauthorized();
        _;
    }

    modifier onlyOperator() {
        if (msg.sender != _operator) revert Unauthorized();
        _;
    }

    modifier onlyStabilizer() {
        if (msg.sender != stabilizer) revert Unauthorized();
        _;
    }

    /// Deploy the contract and fund it with Auton supply.
    /// @param autonity Address of the Autonity Contract
    /// @param operator Address of the Governance Operator
    /// @param stabilizer_ The address that is authorized to mint and burn
    /// @dev The message value is the Auton supply to seed.
    constructor(
        address autonity,
        address operator,
        address stabilizer_
    ) payable nonZeroValue {
        _autonity = autonity;
        _operator = operator;
        stabilizer = stabilizer_;
        totalSupply = msg.value;
    }

    /// Mint Auton and send it to the recipient.
    /// @param recipient Recipient of the Auton
    /// @param amount Amount of Auton to mint (non-zero)
    /// @dev Only the stabilizer is authorized to mint Auton. The recipient
    /// cannot be the stabilizer or the zero address.
    function mint(address recipient, uint amount) external onlyStabilizer {
        if (recipient == address(0) || recipient == stabilizer)
            revert InvalidRecipient();
        if (amount == 0 || amount > address(this).balance)
            revert InvalidAmount();
        payable(recipient).transfer(amount);
        emit Mint(recipient, amount);
    }

    /// Burn Auton by taking it out of circulation.
    /// @dev Only the stabilizer is authorized to burn Auton.
    function burn() external payable nonZeroValue onlyStabilizer {
        emit Burn(msg.value);
    }

    /// Set the Governance Operator account address.
    /// @param operator Address of the new Governance Operator
    /// @dev Only the Autonity Contract is authorized to set the Governance
    /// Operator account address.
    function setOperator(address operator) external onlyAutonity {
        _operator = operator;
    }

    /// Update the stabilizer that is authorized to mint and burn.
    /// @param stabilizer_ The new stabilizer account
    /// @dev Only the operator can update the stabilizer address.
    function setStabilizer(address stabilizer_) external onlyOperator {
        stabilizer = stabilizer_;
    }

    /// The supply of Auton available for minting.
    function availableSupply() external view returns (uint) {
        return address(this).balance;
    }
}
