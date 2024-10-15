// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./Autonity.sol";

contract AccessAutonity {
    
    Autonity internal autonity;

    constructor(address payable _autonity) {
        autonity = Autonity(_autonity);
    }

    /*
    ============================================================

        Modifiers

    ============================================================
     */

    /**
     * @dev Modifier that checks if the caller is the governance operator account.
     */
    modifier onlyOperator {
        require(autonity.getOperator() == msg.sender, "caller is not the operator");
        _;
    }

    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to Autonity contract");
        _;
    }
}