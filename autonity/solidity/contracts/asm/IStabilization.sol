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

/// @title Stabilization Contract Interface
/// @dev Only meant to be used by the Autonity Contract.
interface IStabilization {
    /// Set the Oracle Contract address.
    /// @dev Restricted to the Autonity Contract.
    function setOracle(address oracle) external;
}
