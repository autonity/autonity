// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../AccessAutonity.sol";

abstract contract ValidatorManager is AccessAutonity {

    /** @dev Stores the array of validators bonded to the contract. */
    address[] internal bondedValidators;

    /** 
     * @dev `validatorIdx[validator]` stores the `index+1` of validator in `bondedValidators` array.
     */
    mapping(address => uint256) private validatorIdx;

    struct BondedValidator {
        ILiquidLogic liquidStateContract;
        bool newBondingRequested;
    }

    mapping(address => BondedValidator) private validators;

    /*
    ============================================================
         Internals
    ============================================================
     */

    function _bondingRequestExpired(address _validator) internal {
        validators[_validator].newBondingRequested = false;
    }

    /**
     * @dev Adds the validator in the list and inform that new bonding is requested.
     * @param _validator validator address
     */
    function _newBondingRequested(address _validator) internal {
        _addValidator(_validator);
        validators[_validator].newBondingRequested = true;
    }

    function _initiateValidator(address _validator) private {
        validators[_validator].liquidStateContract = autonity.getValidator(_validator).liquidStateContract;
    }

    /**
     * @dev Adds validator in `bondedValidators` array.
     */
    function _addValidator(address _validator) internal {
        if (validatorIdx[_validator] > 0) return;
        bondedValidators.push(_validator);
        // offset by 1 to handle empty value
        validatorIdx[_validator] = bondedValidators.length;
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
        address _validator;
        ILiquidLogic _stateContract;
        uint256 _atn;
        uint256 _ntn;
        for (uint256 _idx = 0; _idx < bondedValidators.length ; _idx++) {
            // if both liquid balance and unclaimed rewards are 0 and no new bonding is requested
            // then the validator is not needed anymore
            _validator = bondedValidators[_idx];
            _stateContract = _liquidStateContract(_validator);
            while (true) {
                if (validators[_validator].newBondingRequested || _stateContract.balanceOf(_myAddress) > 0) {
                    break;
                }
                (_atn, _ntn) = _stateContract.unclaimedRewards(_myAddress);
                if (_atn > 0 || _ntn > 0) {
                    break;
                }
                _removeValidator(_validator);
                if (_idx >= bondedValidators.length) {
                    break;
                }
                _validator = bondedValidators[_idx];
                _stateContract = _liquidStateContract(_validator);
            }
        }
    }

    function _removeValidator(address _validator) private {
        // index is offset by 1
        uint256 _idxDelete = validatorIdx[_validator];
        if (_idxDelete == bondedValidators.length) {
            // it is the last validator in the array
            bondedValidators.pop();
            delete validatorIdx[_validator];
            validators[_validator].liquidStateContract = ILiquidLogic(address(0));
            return;
        }
        
        // the `_validator` to be deleted sits in the middle of the array
        address _lastValidator = bondedValidators[bondedValidators.length-1];
        // replacing the `_validator` in `_idxDelete-1` with `_lastValidator`, effectively deleting `_validator`
        bondedValidators[_idxDelete-1] = _lastValidator;
        validatorIdx[_lastValidator] = _idxDelete;
        // deleting the last one
        bondedValidators.pop();
        validators[_validator].liquidStateContract = ILiquidLogic(address(0));
        delete validatorIdx[_validator];
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

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns the list of validators bonded some contract.
     */
    function getBondedValidators() virtual external view returns (address[] memory) {
        return bondedValidators;
    }

    /**
     * @notice Returns the amount of LNTN for some contract.
     * @param _validator validator address
     */
    function liquidBalanceOf(address _validator) public virtual view returns (uint256) {
        return validators[_validator].liquidStateContract.balanceOf(address(this));
    }

    /**
     * @notice Returns the amount of unlocked LNTN for some contract.
     * @param _validator validator address
     */
    function unlockedLiquidBalanceOf(address _validator) public virtual view returns (uint256) {
        return validators[_validator].liquidStateContract.unlockedBalanceOf(address(this));
    }

    /**
     * @notice Returns the amount of locked LNTN for some contract.
     * @param _validator validator address
     */
    function lockedLiquidBalanceOf(address _validator) public virtual view returns (uint256) {
        return validators[_validator].liquidStateContract.lockedBalanceOf(address(this));
    }

}