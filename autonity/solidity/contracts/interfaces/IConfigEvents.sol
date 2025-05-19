// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

interface IConfigEvents {
    /**
    * @notice Emitted after updating config parameter of type uint
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    */
    event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue);

    /**
    * @notice Emitted after updating config parameter of type int
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    */
    event ConfigUpdateInt(string name, int256 oldValue, int256 newValue);

    /**
    * @notice Emitted after updating config parameter of type address
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    */
    event ConfigUpdateAddress(string name, address oldValue, address newValue);
}