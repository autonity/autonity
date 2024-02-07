// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./Precompiled.sol";

contract UpgradeManager {
    address public autonity;
    address public operator;

    constructor(address _autonity, address _operator){
        autonity = _autonity;
        operator = _operator;
    }

    /** @dev Call the in-protocol EVM replace mechanism. Requires specific tool to interact.
    * Restricted to the operator account.
    *  @param _target is the target contract address to be updated.
    *  @param _data is the contract creation code.
    */
    function upgrade(address _target, string memory _data) external onlyOperator {
        address precompile = Precompiled.UPGRADER_CONTRACT;
        bytes memory _input = abi.encodePacked(_target, _data);
        assembly {
            let result := delegatecall(gas(), precompile, add(_input,32), mload(_input), 0, 0)
            returndatacopy(0, 0, returndatasize())
            switch result
            case 0 { revert(0, returndatasize()) }
            default { return(0, returndatasize()) }
        }
    }

    /*
    * @notice Set the Operator account. Restricted to the Operator account.
    * @param _account the new operator account.
    */
    function setOperator(address _account) external onlyAutonity {
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
