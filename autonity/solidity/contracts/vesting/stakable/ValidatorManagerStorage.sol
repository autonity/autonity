// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../AccessAutonity.sol";

abstract contract ValidatorManagerStorage is AccessAutonity {

    /** @dev Stores the array of validators bonded to the contract. */
    address[] internal bondedValidators;

    /** 
     * @dev `validatorIndex[validator]` stores the `index+1` of validator in `bondedValidators` array.
     */
    mapping(address => uint256) internal validatorIndex;

    struct LinkedValidator {
        ILiquidLogic liquidStateContract;
        // both of the following are offset by 1
        uint256 lastBondingEpoch;
        uint256 lastUnbondingEpoch;
        uint256 lastUnbondingID;
    }

    mapping(address => LinkedValidator) internal validators;
}