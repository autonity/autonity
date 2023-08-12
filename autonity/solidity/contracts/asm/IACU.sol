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

/// @title ACU Contract Interface
/// @dev Only meant to be used by the Autonity Contract.
interface IACU {
    /// Set the Governance Operator account address.
    /// @param operator Address of the new Governance Operator
    /// @dev Only the Autonity Contract is authorized to set the Governance
    /// Operator account address.
    function setOperator(address operator) external;

    /// Set the Oracle Contract address that is used to retrieve prices.
    /// @param oracle Address of the new Oracle Contract
    /// @dev Only the Autonity Contract is authorized to set the Oracle
    /// Contract address.
    function setOracle(address oracle) external;

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
    function update() external returns (bool status);
}
