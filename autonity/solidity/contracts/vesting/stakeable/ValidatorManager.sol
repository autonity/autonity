// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./ValidatorManagerStorage.sol";

/**
 * @title Manager contract that manages validator related operations
 * @dev Can be reused by any staking contract that stakes from a single delegator.
 */
abstract contract ValidatorManager is ValidatorManagerStorage {

    /*
    ============================================================
         Internals
    ============================================================
     */

    function _bondingRequestExpired(address _validator, uint256 _epochID) internal {
        if (validators[_validator].lastBondingEpoch == _epochID+1) {
            validators[_validator].lastBondingEpoch = 0;
        }
    }

    /**
     * @dev Adds the validator in the list and tracks the `_epochID`.
     * @param _validator validator address
     * @param _epochID epoch id of the request
     */
    function _newBondingRequested(address _validator, uint256 _epochID) internal {
        _addValidator(_validator);
        validators[_validator].lastBondingEpoch = _epochID+1;
    }

    function _getLiquidStateContract(address _validator) internal view returns (ILiquid) {
        ILiquid _liquidContract = validators[_validator].liquidStateContract;
        if (address(_liquidContract) != address(0)) {
            return _liquidContract;
        }
        return autonity.getValidator(_validator).liquidStateContract;
    }

    function _initializeValidator(address _validator) private {
        validators[_validator].liquidStateContract = _getLiquidStateContract(_validator);
    }

    /**
     * @dev Adds validator in `linkedValidators` array.
     */
    function _addValidator(address _validator) internal {
        if (validatorIndex[_validator] > 0) return;
        linkedValidators.push(_validator);
        // offset by 1 to handle empty value
        validatorIndex[_validator] = linkedValidators.length;
        if (address(validators[_validator].liquidStateContract) == address(0)) {
            _initializeValidator(_validator);
        }
    }

    /**
     * @dev Removes all the validators that are not needed for the contract anymore, i.e. any validator
     * that has 0 liquid for that contract and all rewards from the validator are claimed and no pending
     * bonding or unbonding request.
     */
    function _clearValidators() internal {
        address _myAddress = address(this);
        LinkedValidator storage _validator;
        uint256 _epochID = _getEpochID();
        for (uint256 _idx = 0; _idx < linkedValidators.length ; _idx++) {
            // if both liquid balance and unclaimed rewards are 0 and no new bonding is requested
            // then the validator is not needed anymore
            _validator = validators[linkedValidators[_idx]];
            while (true) {
                if (_validator.lastBondingEpoch == _epochID+1) {
                    break;
                }
                if (_validator.liquidStateContract.balanceOf(_myAddress) > 0) {
                    break;
                }
                if (_validator.liquidStateContract.unclaimedRewards(_myAddress) > 0) {
                    break;
                }
                _removeValidator(linkedValidators[_idx]);
                if (_idx >= linkedValidators.length) {
                    break;
                }
                _validator = validators[linkedValidators[_idx]];
            }
        }
    }

    function _removeValidator(address _validator) private {
        // index is offset by 1
        uint256 _idxDelete = validatorIndex[_validator];
        if (_idxDelete == linkedValidators.length) {
            // it is the last validator in the array
            linkedValidators.pop();
            delete validatorIndex[_validator];
            delete validators[_validator];
            return;
        }
        
        // the `_validator` to be deleted sits somewhere in the middle of the array
        address _lastValidator = linkedValidators[linkedValidators.length-1];
        // replacing the `_validator` in `_idxDelete-1` with `_lastValidator`, effectively deleting `_validator`
        linkedValidators[_idxDelete-1] = _lastValidator;
        validatorIndex[_lastValidator] = _idxDelete;
        // deleting the last one
        linkedValidators.pop();
        delete validators[_validator];
        delete validatorIndex[_validator];
    }

    function _getEpochID() internal view returns (uint256) {
        return autonity.epochID();
    }

    function _unclaimedRewards(address _validator) internal view returns (uint256) {
        return _getLiquidStateContract(_validator).unclaimedRewards(address(this));
    }

    function _liquidStateContract(address _validator) internal returns (ILiquid) {
        if (address(validators[_validator].liquidStateContract) == address(0)) {
            _initializeValidator(_validator);
        }
        return validators[_validator].liquidStateContract;
    }

    function _liquidBalance(ILiquid _liquidContract) internal view returns (uint256) {
        return _liquidContract.balanceOf(address(this));
    }

    function _lockedLiquidBalance(ILiquid _liquidContract) internal view returns (uint256) {
        return _liquidContract.lockedBalanceOf(address(this));
    }

    function _unlockedLiquidBalance(ILiquid _liquidContract) internal view returns (uint256) {
        return _liquidContract.unlockedBalanceOf(address(this));
    }
}