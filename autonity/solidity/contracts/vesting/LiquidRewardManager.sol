// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";
import "../Liquid.sol";

contract LiquidRewardManager {

    uint256 public constant FEE_FACTOR_UNIT_RECIP = 1_000_000_000;

    // uint256 private epochID;
    // uint256 private epochFetchedBlock;
    address private operator;
    Autonity internal autonity;

    struct LiquidInfo {
        uint256 lastUnrealisedFeeFactor;
        uint256 unclaimedRewards;
        Liquid liquidContract;
    }

    // stores total liquid and lastUnrealisedFeeFactor for each validator
    // lastUnrealisedFeeFactor is used to calculate unrealised rewards for schedules with the same logic as done in Liquid.sol
    mapping(address => LiquidInfo) private liquidInfo;

    // stores the array of validators bonded to a schedule
    mapping(uint256 => address[]) private bondedValidators;
    // stores the (index+1) of validator in bondedValidators[id] array
    mapping(uint256 => mapping(address => uint256)) private validatorIdx;

    mapping(uint256 => mapping(address => uint256)) private liquidBalances;
    mapping(uint256 => mapping(address => uint256)) private lockedLiquidBalances;
    mapping(uint256 => mapping(address => uint256)) private withdrawnLiquid;

    // realisedFees[id][validator] stores the realised reward entitled to a schedule for a validator
    // unrealisedFeeFactors[id][validator] is used to calculate unrealised rewards. it only updates
    // when the liquid balance of the schedule is updated following the same logic done in Liquid.sol
    mapping(uint256 => mapping(address => uint256)) private realisedFees;
    mapping(uint256 => mapping(address => uint256)) private unrealisedFeeFactors;

    constructor(address payable _autonity) {
        autonity = Autonity(_autonity);
    }

    function _unlockAllAndBurn(uint256 _id, address _validator) internal {
        _decreaseLiquid(_id, _validator, lockedLiquidBalances[_id][_validator]);
        delete lockedLiquidBalances[_id][_validator];
    }

    function _unlockAll(uint256 _id, address _validator) internal {
        delete lockedLiquidBalances[_id][_validator];
    }

    function _unlock(uint256 _id, address _validator, uint256 _amount) internal {
        lockedLiquidBalances[_id][_validator] -= _amount;
    }

    function _lock(uint256 _id, address _validator, uint256 _amount) internal {
        require(liquidBalances[_id][_validator] - lockedLiquidBalances[_id][_validator] >= _amount, "not enough unlocked liquid tokens");
        lockedLiquidBalances[_id][_validator] += _amount;
    }

    /**
     * @dev _decreaseLiquid, _increaseLiquid, _realiseFees, _computeUnrealisedFees follows the same logic as done in Liquid.sol
     * the only difference is in _realiseFees, we claim rewards from the validator first, because this contract does not know
     * when the epoch ends, and so cannot claim rewards at each epoch end.
     * _claimRewards claim rewards from a validator at most once per epoch, so spamming _claimRewards is not a problem
     */
    function _decreaseLiquid(uint256 _id, address _validator, uint256 _amount) internal {
        _realiseFees(_id, _validator);
        liquidBalances[_id][_validator] -= _amount;
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
    }

    function _increaseLiquid(uint256 _id, address _validator, uint256 _amount) internal {
        _realiseFees(_id, _validator);
        liquidBalances[_id][_validator] += _amount;
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
    }

    function _realiseFees(uint256 _id, address _validator) private returns (uint256 _realisedFees) {
        uint256 _lastUnrealisedFeeFactor = liquidInfo[_validator].lastUnrealisedFeeFactor;
        _realisedFees = realisedFees[_id][_validator] + _computeUnrealisedFees(_id, _validator, _lastUnrealisedFeeFactor);
        realisedFees[_id][_validator] = _realisedFees;
        unrealisedFeeFactors[_id][_validator] = _lastUnrealisedFeeFactor;
    }

    function _computeUnrealisedFees(uint256 _id, address _validator, uint256 _lastUnrealisedFeeFactor) private view returns (uint256) {
        return (_lastUnrealisedFeeFactor-unrealisedFeeFactors[_id][_validator]) * liquidBalances[_id][_validator] / FEE_FACTOR_UNIT_RECIP;
    }

    function _claimRewards(address _validator) internal {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        _updateUnclaimedReward(_validator);
        // because all the rewards are being claimed
        _liquidInfo.unclaimedRewards = 0;
        _liquidInfo.liquidContract.claimRewards();
    }

    /**
     * @dev call _rewards(_id) only when rewards are claimed
     * calculates total rewards for a schedule and deletes realisedFees[id][validator] as reward is claimed
     */ 
    function _rewards(uint256 _id) internal returns (uint256) {
        address[] memory _validators = bondedValidators[_id];
        uint256 _totalFees = 0;
        for (uint256 i = 0; i < _validators.length; i++) {
            _claimRewards(_validators[i]);
            _totalFees += _realiseFees(_id, _validators[i]);
            delete realisedFees[_id][_validators[i]];
        }
        return _totalFees;
    }

    function _addValidator(uint256 _id, address _validator) internal {
        if (validatorIdx[_id][_validator] > 0) return;
        address[] storage _validators = bondedValidators[_id];
        _validators.push(_validator);
        validatorIdx[_id][_validator] = _validators.length;
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        if (_liquidInfo.liquidContract == Liquid(address(0))) {
            _liquidInfo.liquidContract = autonity.getValidator(_validator).liquidContract;
        }
    }

    function _clearValidators(uint256 _id) internal {
        address[] storage _validators = bondedValidators[_id];
        for (uint256 i = 0; i < _validators.length ; i++) {
            while (liquidBalances[_id][_validators[i]] == 0 && realisedFees[_id][_validators[i]] == 0) {
                _removeValidator(_id, _validators[i]);
            }
        }
    }

    function _removeValidator(uint256 _id, address _validator) private {
        address[] storage _validators = bondedValidators[_id];
        uint256 _idx = validatorIdx[_id][_validator]-1;
        // removing _validator by replacing it with last one
        _validators[_idx] = _validators[_validators.length-1];
        // deleting the last one
        _validators.pop();
        // update validatorIdx
        delete validatorIdx[_id][_validator];

        if (_idx < _validators.length) {
            validatorIdx[_id][_validators[_idx]] = _idx+1;
        }
    }

    function _updateUnclaimedReward(address _validator) internal {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        Liquid _contract = _liquidInfo.liquidContract;
        if (_contract == Liquid(address(0))) {
            _contract = autonity.getValidator(_validator).liquidContract;
            _liquidInfo.liquidContract = _contract;
        }
        uint256 _totalLiquid = _contract.balanceOf(address(this));
        if (_totalLiquid == 0) {
            return;
        }
        uint256 _reward = _contract.unclaimedRewards(address(this));
        _liquidInfo.lastUnrealisedFeeFactor += (_reward-_liquidInfo.unclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        _liquidInfo.unclaimedRewards = _reward;
    }

    function _unclaimedRewards(uint256 _id) internal returns (uint256) {
        uint256 _totalFee = 0;
        address[] memory _validators = bondedValidators[_id];
        for (uint256 i = 0; i < _validators.length; i++) {
            _updateUnclaimedReward(_validators[i]);
            _totalFee += _realiseFees(_id, _validators[i]);
        }
        return _totalFee;
    }

    function _liquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        return liquidBalances[_id][_validator];
    }

    function _unlockedLiquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        return liquidBalances[_id][_validator] - lockedLiquidBalances[_id][_validator];
    }

    function _lockedLiquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        return lockedLiquidBalances[_id][_validator];
    }

    // returns the list of validator addresses wich are bonded to schedule _id assigned to _account
    function _bondedValidators(uint256 _id) internal view returns (address[] memory) {
        return bondedValidators[_id];
    }

    function _liquidContract(address _validator) internal view returns (Liquid) {
        return liquidInfo[_validator].liquidContract;
    }

}