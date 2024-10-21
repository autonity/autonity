// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../AccessAutonity.sol";

abstract contract ValidatorManagerStorage is AccessAutonity {

    /** @dev Stores the array of validators linked to the contract. */
    address[] internal linkedValidators;

    /** 
     * @dev `validatorIndex[validator]` stores the `index+1` of validator in `linkedValidators` array.
     */
    mapping(address => uint256) internal validatorIndex;

    struct LinkedValidator {
        ILiquid liquidStateContract;
        // the following is offset by 1
        uint256 lastBondingEpoch;
    }

    mapping(address => LinkedValidator) internal validators;
}