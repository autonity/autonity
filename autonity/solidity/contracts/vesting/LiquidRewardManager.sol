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
        bool newBondingRequested;
    }

    // stores total liquid and lastUnrealisedFeeFactor for each validator
    // lastUnrealisedFeeFactor is used to calculate unrealised rewards for schedules with the same logic as done in Liquid.sol
    mapping(address => LiquidInfo) private liquidInfo;

    // stores the array of validators bonded to a schedule
    mapping(uint256 => address[]) private bondedValidators;
    // validatorIdx[_id][_validator] stores the (index+1) of _validator in bondedValidators[_id] array
    mapping(uint256 => mapping(address => uint256)) private validatorIdx;

    mapping(uint256 => mapping(address => Account)) private accounts;

    constructor(address payable _autonity) {
        autonity = Autonity(_autonity);
    }

    function _bondingRequestExpired(uint256 _id, address _validator) internal {
        accounts[_id][_validator].newBondingRequested = false;
    }

    /**
     * @dev initiate a pair (_id, _validator) only once for lifetime
     * initiating the pair helps to reduces maximum allowed gas usage for notification of staking operations
     * @param _id schedule id
     * @param _validator validator address
     */
    function _initiate(uint256 _id, address _validator) internal {
        _addValidator(_id, _validator);
        
        Account storage _account = accounts[_id][_validator];
        _account.newBondingRequested = true;
        if (_account.initiated) {
            return;
        }
        _account.realisedFee = OFFSET;
        _account.unrealisedFeeFactor = OFFSET;
        _account.liquidBalance = OFFSET;
        _account.lockedLiquidBalance = OFFSET;
        _account.initiated = true;
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

    /**
     * @dev _decreaseLiquid, _increaseLiquid, _realiseFees follow the same logic as done in Liquid.sol
     * the only difference is we update unclaimed rewards from the validator first, because this contract does not know
     * when the epoch ends, and so cannot claim rewards at each epoch end.
     * _updateUnclaimedReward or _claimRewards must be called before calling _decreaseLiquid, _increaseLiquid functions.
     * In the process of notification from autonity for staking operation, at first vesting contract is notified about the validators
     * in rewardsDistributed function, where we call _updateUnclaimedReward. Then vesting contract is notified about staking operations
     * in bondingApplied and unbondingApplied where we call _decreaseLiquid and _increaseLiquid
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

    /**
     * @dev _updateUnclaimedReward or _claimRewards must be called before this function
     */
    function _realiseFees(uint256 _id, address _validator) private returns (uint256 _realisedFees) {
        uint256 _lastUnrealisedFeeFactor = liquidInfo[_validator].lastUnrealisedFeeFactor;
        Account storage _account = accounts[_id][_validator];
        _realisedFees = _account.realisedFee
                        + _computeUnrealisedFees(
                            _account.liquidBalance-OFFSET, _account.unrealisedFeeFactor, _lastUnrealisedFeeFactor
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

    /**
     * @dev claims all rewards from the liquid contract of the _validator
     * @param _validator validator address
     */
    function _claimRewards(address _validator) private {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        _updateUnclaimedReward(_validator);
        // because all the rewards are being claimed
        // offset by 1
        _liquidInfo.unclaimedRewards = OFFSET;
        _liquidInfo.liquidContract.claimRewards();
    }

    /**
     * @dev call _claimRewards(_id) only when rewards are claimed
     * calculates total rewards for a schedule and resets realisedFees[id][validator] as reward is claimed
     */ 
    function _claimRewards(uint256 _id) internal returns (uint256) {
        address[] memory _validators = bondedValidators[_id];
        uint256 _totalFees;
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

    /**
     * @dev adds _validator in bondedValidators[_id] array
     */
    function _addValidator(uint256 _id, address _validator) private {
        if (validatorIdx[_id][_validator] > 0) return;
        address[] storage _validators = bondedValidators[_id];
        _validators.push(_validator);
        // offset by 1 to handle empty value
        validatorIdx[_id][_validator] = _validators.length;
        if (liquidInfo[_validator].liquidContract == Liquid(address(0))) {
            _initiateValidator(_validator);
        }
    }

    /**
     * @dev removes all the validators that are not needed for _id anymore, i.e. any validator
     * that has 0 liquid for _id and all rewards from the validator are claimed
     * @param _id schedule id
     */
    function _clearValidators(uint256 _id) internal {
        // must take a storage pointer
        // otherwise _validators will not be updated when _removeValidator is called
        address[] storage _validators = bondedValidators[_id];
        Account storage _account;
        for (uint256 i = 0; i < _validators.length ; i++) {
            _account = accounts[_id][_validators[i]];
            // actual balance is (balance - OFFSET)
            // if liquidBalance is 0, then unrealisedFee is 0
            // if both liquidBalance and realisedFee are 0 and no new bonding is requested then the validator is not needed anymore
            while (
                _account.newBondingRequested == false && _account.liquidBalance == OFFSET && _account.realisedFee == OFFSET
            ) {
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
        uint256 _idxDelete = validatorIdx[_id][_validator]-1;
        if (_idxDelete+1 == _validators.length) {
            // it is the last validator in the array
            _validators.pop();
            delete validatorIdx[_id][_validator];
            return;
        }
        
        // the _validator to be deleted sits in the middle of the array
        address _lastValidator = _validators[_validators.length-1];
        // replacing the _validator in _idxDelete with _lastValidator, effectively deleting it
        _validators[_idxDelete] = _lastValidator;
        validatorIdx[_id][_lastValidator] = _idxDelete+1;
        // deleting the last one
        _validators.pop();
        delete validatorIdx[_id][_validator];
    }

    /**
     * @dev updates the unclaimed AUT from _validator and also updates lastUnrealisedFeeFactor which is used
     * to compute unrealised fees for accounts
     */
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

    /**
     * @dev fetches the unclaimedRewards from liquid contract and calculates the changes in lastUnrealisedFeeFactor
     * but does not update any state variable. This function exists only to support _unclaimedRewards to act as view only function
     */
    function _unfetchedFeeFactor(address _validator) private view returns (uint256) {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        Liquid _contract = _liquidInfo.liquidContract;
        uint256 _totalLiquid = _contract.balanceOf(address(this));
        if (_totalLiquid == 0) {
            return 0;
        }
        // _contract.unclaimedRewards is increased by the same OFFSET of _liquidInfo.unclaimedRewards
        return (_contract.unclaimedRewards(address(this))+OFFSET-_liquidInfo.unclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
    }

    /**
     * @dev calculates the rewards yet to claim for _id
     */
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
                            _account.liquidBalance-OFFSET, _account.unrealisedFeeFactor, _lastUnrealisedFeeFactor
                        );
        
        }
        return _totalFee;
    }

    function _isNewBondingRequested(uint256 _id, address _validator) internal view returns (bool) {
        return accounts[_id][_validator].newBondingRequested;
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

    /**
     * @dev returns the list of validator addresses wich are bonded to schedule _id assigned to _account
     */
    function _bondedValidators(uint256 _id) internal view returns (address[] memory) {
        return bondedValidators[_id];
    }

    function _liquidContract(address _validator) internal view returns (Liquid) {
        return liquidInfo[_validator].liquidContract;
    }

}