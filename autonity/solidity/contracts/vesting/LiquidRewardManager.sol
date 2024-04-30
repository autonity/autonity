// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";
import "../Liquid.sol";

contract LiquidRewardManager {

    uint256 public constant FEE_FACTOR_UNIT_RECIP = 1_000_000_000;

    address private operator;
    Autonity private autonity;

    // when initiated, all values are offset by 1
    uint256 private constant OFFSET = 1;

    struct LiquidInfo {
        uint256 lastUnrealisedFeeFactor;
        uint256 unclaimedRewards;
        Liquid liquidContract;
    }

    struct Account {
        uint256 liquidBalance;
        uint256 lockedLiquidBalance;
        uint256 realisedFee;
        uint256 unrealisedFeeFactor;
        bool initiated;
    }

    // stores total liquid and lastUnrealisedFeeFactor for each validator
    // lastUnrealisedFeeFactor is used to calculate unrealised rewards for schedules with the same logic as done in Liquid.sol
    mapping(address => LiquidInfo) private liquidInfo;

    // stores the array of validators bonded to a schedule
    mapping(uint256 => address[]) private bondedValidators;
    // stores the (index+1) of validator in bondedValidators[id] array
    mapping(uint256 => mapping(address => uint256)) private validatorIdx;

    mapping(uint256 => mapping(address => Account)) private accounts;

    constructor(address payable _autonity) {
        autonity = Autonity(_autonity);
    }

    function _unlock(uint256 _id, address _validator, uint256 _amount) internal {
        Account storage _account = accounts[_id][_validator];
        _account.lockedLiquidBalance -= _amount;
    }

    function _lock(uint256 _id, address _validator, uint256 _amount) internal {
        Account storage _account = accounts[_id][_validator];
        require(_account.liquidBalance - _account.lockedLiquidBalance >= _amount, "not enough unlocked liquid tokens");
        _account.lockedLiquidBalance += _amount;
    }

    function _initiate(uint256 _id, address _validator) internal {
        _addValidator(_id, _validator);
        
        // lockedLiquidBalances[_id][_validator] = 1;
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated) {
            return;
        }
        _account.realisedFee = OFFSET;
        _account.unrealisedFeeFactor = OFFSET;
        _account.liquidBalance = OFFSET;
        _account.lockedLiquidBalance = OFFSET;
        _account.initiated = true;
    }

    /**
     * @dev _decreaseLiquid, _increaseLiquid, _realiseFees, _computeUnrealisedFees follows the same logic as done in Liquid.sol
     * the only difference is in _realiseFees, we claim rewards from the validator first, because this contract does not know
     * when the epoch ends, and so cannot claim rewards at each epoch end.
     * _claimRewards claim rewards from a validator at most once per epoch, so spamming _claimRewards is not a problem
     */
    function _decreaseLiquid(uint256 _id, address _validator, uint256 _amount) internal {
        _realiseFees(_id, _validator);
        Account storage _account = accounts[_id][_validator];
        _account.liquidBalance -= _amount;
    }

    function _increaseLiquid(uint256 _id, address _validator, uint256 _amount) internal {
        _realiseFees(_id, _validator);
        Account storage _account = accounts[_id][_validator];
        _account.liquidBalance += _amount;
    }

    function _realiseFees(uint256 _id, address _validator) private returns (uint256 _realisedFees) {
        uint256 _lastUnrealisedFeeFactor = liquidInfo[_validator].lastUnrealisedFeeFactor;
        Account storage _account = accounts[_id][_validator];
        _realisedFees = _account.realisedFee
                        + _computeUnrealisedFees(
                            _account.liquidBalance-1, _account.unrealisedFeeFactor, _lastUnrealisedFeeFactor
                        );
        _account.realisedFee = _realisedFees;
        _account.unrealisedFeeFactor = _lastUnrealisedFeeFactor;
        // remove the offset
        _realisedFees -= OFFSET;
    }

    function _computeUnrealisedFees(
        uint256 _balance, uint256 _unrealisedFeeFactor, uint256 _lastUnrealisedFeeFactor
    ) private pure returns (uint256) {
        if (_balance == 0) {
            return 0;
        }
        return (_lastUnrealisedFeeFactor - _unrealisedFeeFactor) * _balance / FEE_FACTOR_UNIT_RECIP;
    }

    function _claimRewards(address _validator) internal {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        _updateUnclaimedReward(_validator);
        // because all the rewards are being claimed
        // offset by 1
        _liquidInfo.unclaimedRewards = OFFSET;
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
            accounts[_id][_validators[i]].realisedFee = OFFSET;
        }
        return _totalFees;
    }

    function _initiateValidator(address _validator) private {
        Liquid _contract = autonity.getValidator(_validator).liquidContract;
        // offset by 1
        liquidInfo[_validator] = LiquidInfo(OFFSET, OFFSET, _contract);
    }

    function _addValidator(uint256 _id, address _validator) private {
        if (validatorIdx[_id][_validator] > 0) return;
        address[] storage _validators = bondedValidators[_id];
        _validators.push(_validator);
        // offset by 1 to handle empty value
        validatorIdx[_id][_validator] = _validators.length;
        if (liquidInfo[_validator].liquidContract == Liquid(address(0))) {
            // _liquidInfo.liquidContract = autonity.getValidator(_validator).liquidContract;
            _initiateValidator(_validator);
        }
    }

    function _clearValidators(uint256 _id) internal {
        address[] storage _validators = bondedValidators[_id];
        Account storage _account;
        for (uint256 i = 0; i < _validators.length ; i++) {
            _account = accounts[_id][_validators[i]];
            while (_account.liquidBalance == OFFSET && _account.realisedFee == OFFSET) {
                _removeValidator(_id, _validators[i]);
                if (i >= _validators.length) {
                    break;
                }
                _account = accounts[_id][_validators[i]];
            }
        }
    }

    function _removeValidator(uint256 _id, address _validator) private {
        address[] storage _validators = bondedValidators[_id];
        uint256 _idx = validatorIdx[_id][_validator]-1;
        if (_idx+1 == _validators.length) {
            _validators.pop();
            delete validatorIdx[_id][_validator];
            return;
        }
        
        // the _validator to be deleted sits in the middle of the array
        // removing _validator by replacing it with last one
        // otherwise we will need to iterate the whole array
        address _lastValidator = _validators[_validators.length-1];
        _validators[_idx] = _lastValidator;
        validatorIdx[_id][_lastValidator] = _idx+1;
        // deleting the last one
        _validators.pop();
        delete validatorIdx[_id][_validator];
    }

    function _updateUnclaimedReward(address _validator) internal {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        Liquid _contract = _liquidInfo.liquidContract;
        uint256 _totalLiquid = _contract.balanceOf(address(this));
        if (_totalLiquid == 0) {
            return;
        }
        // _reward is increased by the same OFFSET of _liquidInfo.unclaimedRewards
        uint256 _reward = _contract.unclaimedRewards(address(this))+OFFSET;
        _liquidInfo.lastUnrealisedFeeFactor += (_reward-_liquidInfo.unclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        _liquidInfo.unclaimedRewards = _reward;
    }

    function _unfetchedFeeFactor(address _validator) internal view returns (uint256) {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        Liquid _contract = _liquidInfo.liquidContract;
        uint256 _totalLiquid = _contract.balanceOf(address(this));
        if (_totalLiquid == 0) {
            return 0;
        }
        // _contract.unclaimedRewards is increased by the same OFFSET of _liquidInfo.unclaimedRewards
        return (_contract.unclaimedRewards(address(this))+OFFSET-_liquidInfo.unclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
    }

    function _unclaimedRewards(uint256 _id) internal view returns (uint256) {
        uint256 _totalFee;
        address[] memory _validators = bondedValidators[_id];
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _lastUnrealisedFeeFactor = liquidInfo[_validators[i]].lastUnrealisedFeeFactor
                                                + _unfetchedFeeFactor(_validators[i]);

            Account storage _account = accounts[_id][_validators[i]];
            // remove offset from realisedFee
            _totalFee += _account.realisedFee - OFFSET
                        + _computeUnrealisedFees(
                            _account.liquidBalance-1, _account.unrealisedFeeFactor, _lastUnrealisedFeeFactor
                        );
        
        }
        return _totalFee;
    }

    function _liquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated) {
            return _account.liquidBalance - OFFSET;
        }
        return 0;
    }

    function _unlockedLiquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated) {
            return _account.liquidBalance - _account.lockedLiquidBalance;
        }
        return 0;
    }

    function _lockedLiquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated) {
            return _account.lockedLiquidBalance - OFFSET;
        }
        return 0;
    }

    // returns the list of validator addresses wich are bonded to schedule _id assigned to _account
    function _bondedValidators(uint256 _id) internal view returns (address[] memory) {
        return bondedValidators[_id];
    }

    function _liquidContract(address _validator) internal view returns (Liquid) {
        return liquidInfo[_validator].liquidContract;
    }

}