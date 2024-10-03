// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../AccessAutonity.sol";

abstract contract ValidatorManager is AccessAutonity {

    /** @dev Stores the array of validators bonded to the contract. */
    address[] internal bondedValidators;

    /** 
     * @dev `validatorIndex[validator]` stores the `index+1` of validator in `bondedValidators` array.
     */
    mapping(address => uint256) private validatorIndex;

    struct LinkedValidator {
        ILiquidLogic liquidStateContract;
        // both of the following are offset by 1
        uint256 lastBondingEpoch;
        uint256 lastUnbondingEpoch;
        uint256 lastUnbondingID;
    }

    mapping(address => LinkedValidator) private validators;

    /*
    ============================================================
         Internals
    ============================================================
     */

    function _unbondingRequestExpired(address _validator, uint256 _unbondingID) internal {
        if (validators[_validator].lastUnbondingID == _unbondingID+1) {
            validators[_validator].lastUnbondingID = 0;
            validators[_validator].lastUnbondingEpoch = 0;
        }
    }

    function _newUnbondingRequested(address _validator, uint256 _epochID, uint256 _unbondingID) internal {
        validators[_validator].lastUnbondingID = _unbondingID+1;
        validators[_validator].lastUnbondingEpoch = _epochID+1;
    }

    function _bondingRequestExpired(address _validator, uint256 _epochID) internal {
        if (validators[_validator].lastBondingEpoch == _epochID+1) {
            validators[_validator].lastBondingEpoch = 0;
        }
    }

    /**
     * @dev Adds the validator in the list and inform that new bonding is requested.
     * @param _validator validator address
     */
    function _newBondingRequested(address _validator, uint256 _epochID) internal {
        _addValidator(_validator);
        validators[_validator].lastBondingEpoch = _epochID+1;
    }

    function _initiateValidator(address _validator) private {
        validators[_validator].liquidStateContract = autonity.getValidator(_validator).liquidStateContract;
    }

    /**
     * @dev Adds validator in `bondedValidators` array.
     */
    function _addValidator(address _validator) internal {
        if (validatorIndex[_validator] > 0) return;
        bondedValidators.push(_validator);
        // offset by 1 to handle empty value
        validatorIndex[_validator] = bondedValidators.length;
        if (address(validators[_validator].liquidStateContract) == address(0)) {
            _initiateValidator(_validator);
        }
    }

    /**
     * @dev Removes all the validators that are not needed for some contract anymore, i.e. any validator
     * that has 0 liquid for that contract and all rewards from the validator are claimed.
     */
    function _clearValidators() internal {
        address _myAddress = address(this);
        uint256 _atn;
        uint256 _ntn;
        LinkedValidator storage _validator;
        uint256 _epochID = _getEpochID();
        for (uint256 _idx = 0; _idx < bondedValidators.length ; _idx++) {
            // if both liquid balance and unclaimed rewards are 0 and no new bonding is requested
            // then the validator is not needed anymore
            _validator = validators[bondedValidators[_idx]];
            while (true) {
                if (_validator.lastBondingEpoch == _epochID+1 || _validator.lastUnbondingEpoch == _epochID+1) {
                    break;
                }
                if (_validator.lastUnbondingID > 0) {
                    if (autonity.isUnbondingReleased(_validator.lastUnbondingID) == false) {
                        break;
                    }
                }
                if (_validator.liquidStateContract.balanceOf(_myAddress) > 0) {
                    break;
                }
                (_atn, _ntn) = _validator.liquidStateContract.unclaimedRewards(_myAddress);
                if (_atn > 0 || _ntn > 0) {
                    break;
                }
                _removeValidator(bondedValidators[_idx]);
                if (_idx >= bondedValidators.length) {
                    break;
                }
                _validator = validators[bondedValidators[_idx]];
            }
        }
    }

    function _removeValidator(address _validator) private {
        // index is offset by 1
        uint256 _idxDelete = validatorIndex[_validator];
        if (_idxDelete == bondedValidators.length) {
            // it is the last validator in the array
            bondedValidators.pop();
            delete validatorIndex[_validator];
            delete validators[_validator];
            return;
        }
        
        // the `_validator` to be deleted sits in the middle of the array
        address _lastValidator = bondedValidators[bondedValidators.length-1];
        // replacing the `_validator` in `_idxDelete-1` with `_lastValidator`, effectively deleting `_validator`
        bondedValidators[_idxDelete-1] = _lastValidator;
        validatorIndex[_lastValidator] = _idxDelete;
        // deleting the last one
        bondedValidators.pop();
        delete validators[_validator];
        delete validatorIndex[_validator];
    }

    function _getEpochID() internal view returns (uint256) {
        return autonity.epochID();
    }

    function _unclaimedRewards(address _validator) internal view returns (uint256, uint256) {
        return validators[_validator].liquidStateContract.unclaimedRewards(address(this));
    }

    function _liquidStateContract(address _validator) internal returns (ILiquidLogic) {
        if (address(validators[_validator].liquidStateContract) == address(0)) {
            _initiateValidator(_validator);
        }
        return validators[_validator].liquidStateContract;
    }

    function _liquidBalance(ILiquidLogic _liquidContract) internal view returns (uint256) {
        return _liquidContract.balanceOf(address(this));
    }

    function _lockedLiquidBalance(ILiquidLogic _liquidContract) internal view returns (uint256) {
        return _liquidContract.lockedBalanceOf(address(this));
    }

    function _unlockedLiquidBalance(ILiquidLogic _liquidContract) internal view returns (uint256) {
        return _liquidContract.unlockedBalanceOf(address(this));
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns the list of validators bonded some contract.
     */
    function getLinkedValidators() virtual external view returns (address[] memory) {
        return bondedValidators;
    }

    /**
     * @notice Returns the amount of LNTN for some contract.
     * @param _validator validator address
     */
    function liquidBalance(address _validator) public virtual view returns (uint256) {
        ILiquidLogic _liquidContract = validators[_validator].liquidStateContract;
        if (address(_liquidContract) == address(0)) {
            _liquidContract = autonity.getValidator(_validator).liquidStateContract;
        }
        return _liquidBalance(_liquidContract);
    }

    /**
     * @notice Returns the amount of unlocked LNTN for some contract.
     * @param _validator validator address
     */
    function unlockedLiquidBalance(address _validator) public virtual view returns (uint256) {
        ILiquidLogic _liquidContract = validators[_validator].liquidStateContract;
        if (address(_liquidContract) == address(0)) {
            _liquidContract = autonity.getValidator(_validator).liquidStateContract;
        }
        return _unlockedLiquidBalance(_liquidContract);
    }

    /**
     * @notice Returns the amount of locked LNTN for some contract.
     * @param _validator validator address
     */
    function lockedLiquidBalance(address _validator) public virtual view returns (uint256) {
        ILiquidLogic _liquidContract = validators[_validator].liquidStateContract;
        if (address(_liquidContract) == address(0)) {
            _liquidContract = autonity.getValidator(_validator).liquidStateContract;
        }
        return _lockedLiquidBalance(_liquidContract);
    }

}