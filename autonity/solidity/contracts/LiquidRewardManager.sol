// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./Autonity.sol";
import "./Liquid.sol";

contract LiquidRewardManager {

    uint256 public constant FEE_FACTOR_UNIT_RECIP = 1_000_000_000;

    uint256 private epochID;
    uint256 private epochFetchedBlock;
    address private operator;
    Autonity internal autonity;

    struct LiquidInfo {
        uint256 totalLiquid;
        uint256 lastUnrealisedFeeFactor;
    }

    // stores total liquid and lastUnrealisedFeeFactor for each validator
    // lastUnrealisedFeeFactor is used to calculate unrealised rewards for schedules with the same logic as done in Liquid.sol
    mapping(address => LiquidInfo) private validatorLiquids;

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

    // rewardsClaimedEpoch[validator] stores the last epoch where the rewards from validator were claimed
    mapping(address => uint256) private rewardsClaimedEpoch;

    constructor(address payable _autonity) {
        autonity = Autonity(_autonity);
    }

    function _unlock(uint256 _id, address _validator, uint256 _amount) internal {
        require(lockedLiquidBalances[_id][_validator] >= _amount, "not enough locked balance");
        lockedLiquidBalances[_id][_validator] -= _amount;
    }

    function _lock(uint256 _id, address _validator, uint256 _amount) internal {
        require(_unlockedLiquidBalanceOf(_id, _validator) >= _amount, "not enough unlocked balance");
        lockedLiquidBalances[_id][_validator] += _amount;
    }

    /**
     * @dev _decreaseLiquid, _increaseLiquid, _realiseFees, _computeUnrealisedFees follows the same logic as done in Liquid.sol
     * the only difference is in _realiseFees, we claim rewards from the validator first, because this contract does not know
     * when the epoch ends, and so cannot claim rewards at each epoch end.
     * _claimRewards claim rewards from a validator at most once per epoch, so spamming _claimRewards is not a problem
     */
    function _decreaseLiquid(uint256 _id, address _validator, uint256 _amount) internal {
        require(
            liquidBalances[_id][_validator] - lockedLiquidBalances[_id][_validator] >= _amount,
            "not enough unlocked liquid tokens"
        );

        _realiseFees(_id, _validator);
        liquidBalances[_id][_validator] -= _amount;
        validatorLiquids[_validator].totalLiquid -= _amount;
        if (liquidBalances[_id][_validator] == 0) {
            _removeValidator(_id, _validator);
            delete unrealisedFeeFactors[_id][_validator];
        }
    }

    function _increaseLiquid(uint256 _id, address _validator, uint256 _amount) internal {
        _realiseFees(_id, _validator);
        if (liquidBalances[_id][_validator] == 0) {
            _addValidator(_id, _validator);
        }
        liquidBalances[_id][_validator] += _amount;
        validatorLiquids[_validator].totalLiquid += _amount;
    }

    function _realiseFees(uint256 _id, address _validator) private returns (uint256 _realisedFees) {
        _claimRewards(_validator);
        uint256 _unrealisedFees = _computeUnrealisedFees(_id, _validator);
        _realisedFees = realisedFees[_id][_validator] + _unrealisedFees;
        realisedFees[_id][_validator] = _realisedFees;
        unrealisedFeeFactors[_id][_validator] = validatorLiquids[_validator].lastUnrealisedFeeFactor;
    }

    function _computeUnrealisedFees(uint256 _id, address _validator) private view returns (uint256) {
        uint256 balance = liquidBalances[_id][_validator];
        if (balance == 0) {
            return 0;
        }
        uint256 _unrealisedFeeFactor =
            validatorLiquids[_validator].lastUnrealisedFeeFactor - unrealisedFeeFactors[_id][_validator];
        uint256 _unrealisedFee = (_unrealisedFeeFactor * balance) / FEE_FACTOR_UNIT_RECIP;
        return _unrealisedFee;
    }

    function _claimRewards(address _validator) private {
        if (rewardsClaimedEpoch[_validator] == _epochID()) {
            return;
        }
        Liquid liquidContract = autonity.getValidator(_validator).liquidContract;
        uint256 reward = address(this).balance;
        liquidContract.claimRewards();
        reward = address(this).balance - reward;
        if (reward > 0) {
            LiquidInfo storage liquidInfo = validatorLiquids[_validator];
            require(liquidInfo.totalLiquid > 0, "got reward from validator with no liquid supply"); // this shouldn't happen
            liquidInfo.lastUnrealisedFeeFactor += (reward * FEE_FACTOR_UNIT_RECIP) / liquidInfo.totalLiquid;
        }
        rewardsClaimedEpoch[_validator] = _epochID();
    }

    /**
     * @dev call _rewards(_id) only when rewards are claimed
     * calculates total rewards for a schedule and deletes realisedFees[id][validator] as reward is claimed
     */ 
    function _rewards(uint256 _id) internal returns (uint256) {
        address[] storage validators = bondedValidators[_id];
        uint256 totalFees = 0;
        for (uint256 i = 0; i < validators.length; i++) {
            totalFees += _realiseFees(_id, validators[i]);
            delete realisedFees[_id][validators[i]];
        }
        return totalFees;
    }

    function _addValidator(uint256 _id, address _validator) private {
        address[] storage validators = bondedValidators[_id];
        validators.push(_validator);
        validatorIdx[_id][_validator] = validators.length;
    }

    function _removeValidator(uint256 _id, address _validator) private {
        address[] storage validators = bondedValidators[_id];
        uint256 idx = validatorIdx[_id][_validator]-1;
        // removing _validator by replacing it with last one and then deleting the last one
        validators[idx] = validators[validators.length-1];
        validators.pop();
        delete validatorIdx[_id][_validator];

        if (idx < validators.length) {
            validatorIdx[_id][validators[idx]] = idx+1;
        }
    }

    function _epochID() private returns (uint256) {
        if (epochFetchedBlock < block.number) {
            epochFetchedBlock = block.number;
            epochID = autonity.epochID();
        }
        return epochID;
    }

    function _unclaimedRewards(uint256 _id) internal view returns (uint256) {
        uint256 totalFee = 0;
        address[] storage validators = bondedValidators[_id];
        for (uint256 i = 0; i < validators.length; i++) {
            address validator = validators[i];
            totalFee += realisedFees[_id][validator] + _computeUnrealisedFees(_id, validator);
        }
        return totalFee;
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

}