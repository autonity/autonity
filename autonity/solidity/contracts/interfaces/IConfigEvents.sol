// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.30;

interface IConfigEvents {
    /**
    * @notice Emitted after updating config parameter of type uint
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    * @param appliesAtHeight block at which the change will apply
    */
    event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue, uint256 appliesAtHeight);

    /**
    * @notice Emitted after updating config parameter of type int
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    * @param appliesAtHeight block at which the change will apply
    */
    event ConfigUpdateInt(string name, int256 oldValue, int256 newValue, uint256 appliesAtHeight);

    /**
    * @notice Emitted after updating config parameter of type address
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    * @param appliesAtHeight block at which the change will apply
    */
    event ConfigUpdateAddress(string name, address oldValue, address newValue, uint256 appliesAtHeight);

    /**
    * @notice Emitted after updating config parameter of type boolean
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    * @param appliesAtHeight block at which the change will apply
    */
    event ConfigUpdateBool(string name, bool oldValue, bool newValue, uint256 appliesAtHeight);
}
