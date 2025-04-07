// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./lib/Precompiled.sol";
import {IConfigEvents} from "./interfaces/IConfigEvents.sol";
import {ReentrancyGuard} from "./ReentrancyGuard.sol";

contract UpgradeManager is IConfigEvents, ReentrancyGuard {
    address internal autonity;
    address internal operator;

    constructor(address _autonity, address _operator){
        autonity = _autonity;
        operator = _operator;
    }
    event UpgradeResult(address indexed contractAddress, bool success);

    /** @dev Call the in-protocol EVM replace mechanism. Requires specific tool to interact.
    * Restricted to the operator account.
    *  @param _target is the target contract address to be updated.
    *  @param _data is the contract creation code.
    */
    function upgrade(address _target, string memory _data) external virtual nonReentrant onlyOperator {
        address precompile = Precompiled.UPGRADER_CONTRACT;
        bytes memory _input = abi.encodePacked(_target, _data);

        bool success;
        uint256 returnSize;
        bytes memory returnData;

        assembly {
            let result := delegatecall(gas(), precompile, add(_input, 32), mload(_input), 0, 0)
            success := result
            returnSize := returndatasize()
            //load free memory pointer
            returnData := mload(0x40)
            // update free memory pointer to point to the end of the allocated space for returnData
            mstore(0x40, add(returnData, add(returnSize, 32)))
            // write the size of the returnData first
            mstore(returnData, returnSize)
            // copy actual return data after the length
            returndatacopy(add(returnData, 32), 0, returnSize)
        }
        emit UpgradeResult(_target, success);
        assembly {
            if iszero(success) {
                revert(add(returnData, 32), returnSize)
            }
            // return raw data skipping the length
            return(add(returnData, 32), returnSize)
        }
    }

    /// @return returns the autonity contract address
    function getAutonity() external virtual view nonReentrantView returns (address) {
        return autonity;
    }

    /// @return returns the operator address
    function getOperator() external virtual view nonReentrantView returns (address) {
        return operator;
    }

    /*
    * @notice Set the Operator account. Restricted to the Operator account.
    * @param _account the new operator account.
    */
    function setOperator(address _account) external virtual nonReentrant onlyAutonity {
        emit IConfigEvents.ConfigUpdateAddress("operator", operator, _account);
        operator = _account;
    }

    /**
    * @dev Modifier that checks if the caller is the operator contract.
    */
    modifier onlyOperator  {
        require(operator == msg.sender, "caller is not the operator");
        _;
    }

    /**
    * @dev Modifier that checks if the caller is the Autonity Contract
    */
    modifier onlyAutonity  {
        require(autonity == msg.sender, "caller is not the Autonity contract");
        _;
    }
}
