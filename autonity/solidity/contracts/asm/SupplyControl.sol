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
import {IAutonity} from "../interfaces/IAutonity.sol";

/// @title ASM Supply Control Contract Implementation
/// @notice Controls the supply of Auton on the network.
/// @dev Intended to be deployed by the protocol at genesis. The operator is
/// expected to be the Stabilization Contract.
contract SupplyControl is ISupplyControl {
    /// The Autonity Contract Address
    IAutonity private autonity;

    /// The account that is authorized to mint and burn.
    address public stabilizer;

    /// The total supply of Auton under management.
    uint256 public totalSupply;

    error InvalidAmount();
    error InvalidRecipient();
    error Unauthorized();
    error ZeroValue();

    /// Deploy the contract and fund it with Auton supply.
    /// @param _autonity The Autonity Contract address
    /// @dev The message value is the Auton supply to seed.
    constructor(address payable _autonity, address _stabilizer) payable nonZeroValue {
        autonity = IAutonity(_autonity);
        stabilizer = _stabilizer;
        totalSupply = msg.value;
    }

    /// Mint Auton and send it to the recipient.
    /// @param _recipient Recipient of the Auton
    /// @param _amount Amount of Auton to mint (non-zero)
    /// @dev Only the Stabilizer is authorized to mint Auton. The recipient
    /// cannot be the operator or the zero address.
    function mint(address _recipient, uint _amount) external onlyStabilizer {
        if (_recipient == address(0) || _recipient == stabilizer)
            revert InvalidRecipient();
        if (_amount == 0 || _amount > address(this).balance)
            revert InvalidAmount();
        payable(_recipient).transfer(_amount);
        emit Mint(_recipient, _amount);
    }

    /// Burn Auton by taking it out of circulation.
    /// @dev Only the stabilizer is authorized to burn Auton.
    function burn() external payable nonZeroValue onlyStabilizer {
        emit Burn(msg.value);
    }

    /// Update the stabilizer that is authorized to mint and burn.
    /// @param _stabilizer The new operator account
    /// @dev Only the Autonity governance operator can update the stabilizer account.
    function setStabilizer(address _stabilizer) external onlyOperator {
        stabilizer = _stabilizer;
    }

    /// The supply of Auton available for minting.
    function availableSupply() external view returns (uint) {
        return address(this).balance;
    }

    modifier nonZeroValue() {
        if (msg.value == 0) revert ZeroValue();
        _;
    }

    modifier onlyStabilizer() {
        if (msg.sender != stabilizer) revert Unauthorized();
        _;
    }

    modifier onlyOperator() {
        if (msg.sender != autonity.getOperator()) revert Unauthorized();
        _;
    }

}
