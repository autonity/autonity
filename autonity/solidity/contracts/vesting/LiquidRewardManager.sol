// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";

contract LiquidRewardManager {

    uint256 public constant FEE_FACTOR_UNIT_RECIP = 1_000_000_000;

    address private operator;
    Autonity private autonity;

    struct RewardEvent {
        uint256 epochID;
        uint256 totalLiquid;
        uint256 stakingRequestID;
        bool bonding;
        bool eventExist;
        bool applied;
    }

    struct RewardTracker {
        uint256 atnUnclaimedRewards;
        uint256 ntnUnclaimedRewards;
        uint256 lastUpdateEpochID;
        address payable liquidStateContract;
        RewardEvent lastRewardEvent;
        RewardEvent pendingRewardEvent;
        mapping(uint256 => uint256) atnLastUnrealisedFeeFactor;
        mapping(uint256 => uint256) ntnLastUnrealisedFeeFactor;
    }

    struct Account {
        uint256 liquidBalance;
        uint256 lockedLiquidBalance;
        uint256 atnRealisedFee;
        uint256 atnUnrealisedFeeFactor;
        uint256 ntnRealisedFee;
        uint256 ntnUnrealisedFeeFactor;
        bool newBondingRequested;
    }

    // stores total liquid and lastUnrealisedFeeFactor for each validator
    // lastUnrealisedFeeFactor is used to calculate unrealised rewards for schedules with the same logic as done in Liquid.sol
    mapping(address => RewardTracker) private rewardTracker;

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
     * @dev Adds the validator in the list and inform that new bonding is requested
     * @param _id schedule id
     * @param _validator validator address
     */
    function _newBondingRequested(uint256 _id, address _validator, uint256 _bondingID) internal {
        _addValidator(_id, _validator);
        accounts[_id][_validator].newBondingRequested = true;
        _newPendingRewardEvent(_validator, _bondingID, true);
    }

    function _unlockAndBurnLiquid(uint256 _id, address _validator, uint256 _amount, uint256 _epochID) internal {
        Account storage _account = accounts[_id][_validator];
        _account.lockedLiquidBalance -= _amount;
        _burnLiquid(_id, _validator, _amount, _epochID);
    }

    function _lock(uint256 _id, address _validator, uint256 _amount) internal {
        Account storage _account = accounts[_id][_validator];
        require(_account.liquidBalance - _account.lockedLiquidBalance >= _amount, "not enough unlocked liquid tokens");
        _account.lockedLiquidBalance += _amount;
    }

    /**
     * @dev Burns some liquid tokens that represents liquid bonded to some validator from some contract.
     * The following functions: _burnLiquid, _mintLiquid, _realiseFees follow the same logic as done in Liquid.sol.
     * The only difference is that the liquid is not updated immediately. The liquid update reflects the changes after
     * the staking operations of epochID are applied.
     */
    function _burnLiquid(uint256 _id, address _validator, uint256 _amount, uint256 _epochID) internal {
        _realiseFees(_id, _validator, _epochID);
        Account storage _account = accounts[_id][_validator];
        _account.liquidBalance -= _amount;
    }

    /**
     * @dev Mints some liquid tokens that represents liquid bonded to some validator from some contract.
     */
    function _mintLiquid(uint256 _id, address _validator, uint256 _amount, uint256 _epochID) internal {
        _realiseFees(_id, _validator, _epochID);
        Account storage _account = accounts[_id][_validator];
        _account.liquidBalance += _amount;
    }

    /**
     * @dev Realise fees until epochID. Must update rewards before realising fees.
     */
    function _realiseFees(uint256 _id, address _validator, uint256 _epochID) private returns (uint256 _atnRealisedFees, uint256 _ntnRealisedFees) {
        _updateUnclaimedReward(_validator, _epochID);
        uint256 _atnLastUnrealisedFeeFactor = rewardTracker[_validator].atnLastUnrealisedFeeFactor[_epochID];
        Account storage _account = accounts[_id][_validator];
        uint256 _balance = _account.liquidBalance;
        _atnRealisedFees = _account.atnRealisedFee
                        + _computeUnrealisedFees(
                            _balance, _account.atnUnrealisedFeeFactor, _atnLastUnrealisedFeeFactor
                        );
        _account.atnRealisedFee = _atnRealisedFees;
        _account.atnUnrealisedFeeFactor = _atnLastUnrealisedFeeFactor;

        uint256 _ntnLastUnrealisedFeeFactor = rewardTracker[_validator].ntnLastUnrealisedFeeFactor[_epochID];
        _ntnRealisedFees = _account.ntnRealisedFee
                        + _computeUnrealisedFees(
                            _balance, _account.ntnUnrealisedFeeFactor, _ntnLastUnrealisedFeeFactor
                        );
        _account.ntnRealisedFee = _ntnRealisedFees;
        _account.ntnUnrealisedFeeFactor = _ntnLastUnrealisedFeeFactor;
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
     * @dev Claims all rewards from the liquid contract of the _validator
     * @param _validator validator address
     */
    function _claimRewards(address _validator) private {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        uint256 _epochID = _getEpochID();
        require(_epochID > 0, "no rewards until first epoch finalized");
        _updateUnclaimedReward(_validator, _epochID-1);
        // because all the rewards are being claimed
        // offset by 1
        _rewardTracker.atnUnclaimedRewards = 0;
        _rewardTracker.ntnUnclaimedRewards = 0;
        ILiquidLogic(_rewardTracker.liquidStateContract).claimRewards();
    }

    /**
     * @dev Calculates total rewards for a schedule and resets realisedFees[id][validator] as reward is claimed
     */
    function _claimRewards(uint256 _id) internal returns (uint256 _atnTotalFees, uint256 _ntnTotalFees) {
        address[] memory _validators = bondedValidators[_id];
        for (uint256 i = 0; i < _validators.length; i++) {
            (uint256 _atnFee, uint256 _ntnFee) = _claimRewards(_id, _validators[i]);
            _atnTotalFees += _atnFee;
            _ntnTotalFees += _ntnFee;
        }
    }

    function _claimRewards(uint256 _id, address _validator) internal returns (uint256 _atnFee, uint256 _ntnFee) {
        Account storage _account = accounts[_id][_validator];
        _claimRewards(_validator);
        (_atnFee, _ntnFee) = _realiseFees(_id, _validator, _getEpochID());
        _account.atnRealisedFee = 0;
        _account.ntnRealisedFee = 0;
    }

    function _initiateValidator(address _validator) private {
        rewardTracker[_validator].liquidStateContract = autonity.getValidator(_validator).liquidStateContract;
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
        if (rewardTracker[_validator].liquidStateContract == address(0)) {
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
        for (uint256 _idx = 0; _idx < _validators.length ; _idx++) {
            _account = accounts[_id][_validators[_idx]];
            // actual balance is (balance - OFFSET)
            // if liquidBalance is 0, then unrealisedFee is 0
            // if both liquidBalance and realisedFee are 0 and no new bonding is requested then the validator is not needed anymore
            while (
                _account.newBondingRequested == false && _account.liquidBalance == 0
                && _account.atnRealisedFee == 0 && _account.ntnRealisedFee == 0
            ) {
                _removeValidator(_id, _validators[_idx]);
                if (_idx >= _validators.length) {
                    break;
                }
                _account = accounts[_id][_validators[_idx]];
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
     * @dev If pending rewrad event exists and if the event is not from current epoch, then the pending rewrad event
     * replaces the last reward event.
     */
    function _updateLastRewardEvent(address _validator) internal {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        RewardEvent storage _pending = _rewardTracker.pendingRewardEvent;
        RewardEvent storage _last = _rewardTracker.lastRewardEvent;
        if (_pending.eventExist == true && _pending.epochID < _getEpochID()) {
            require(_last.applied == true, "reward needs to be updated from last event before replacing");
            _last = _pending;
            _pending.eventExist = false;
        }
    }

    function _newPendingRewardEvent(address _validator, uint256 _stakingRequestID, bool _bonding) internal {
        _updateLastRewardEvent(_validator);
        // if there is already a pending event, replacing it is fine
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        _rewardTracker.pendingRewardEvent = RewardEvent(
            _getEpochID(),
            _rewardTracker.liquidContract.balanceOf(address(this)),
            _stakingRequestID,
            _bonding,
            true,
            false
        );
    }

    function _updatePendingEvent(address _validator) internal {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        if (_rewardTracker.pendingRewardEvent.eventExist == true) {
            _rewardTracker.pendingRewardEvent.totalLiquid = _rewardTracker.liquidContract.balanceOf(address(this));
        }
    }

    function _updateLastUnrealisedFeeFactor(
        RewardTracker storage _rewardTracker,
        uint256 _epochID,
        uint256 _atnReward,
        uint256 _ntnReward,
        uint256 _totalLiquid
    ) private {
        mapping(uint256 => uint256) storage _lastUnrealisedFeeFactor =_rewardTracker.atnLastUnrealisedFeeFactor;
        _lastUnrealisedFeeFactor[_epochID] = _lastUnrealisedFeeFactor[_rewardTracker.lastUpdateEpochID]
                                            + (_atnReward-_rewardTracker.atnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        _rewardTracker.atnUnclaimedRewards = _atnReward;

        _lastUnrealisedFeeFactor =_rewardTracker.ntnLastUnrealisedFeeFactor;
        _lastUnrealisedFeeFactor[_epochID] = _lastUnrealisedFeeFactor[_rewardTracker.lastUpdateEpochID]
                                            + (_ntnReward-_rewardTracker.ntnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        _rewardTracker.ntnUnclaimedRewards = _ntnReward;
        _rewardTracker.lastUpdateEpochID = _epochID;
    }

    /**
     * @dev Updates the unclaimed rewards from validator and the lastUnrealisedFeeFactor which is used
     * to compute unrealised fees for accounts. Both is updated until given epoch id. The lastUnrealisedFeeFactor
     * is kept with history, so we have a mapping(epochID => value) of lastUnrealisedFeeFactor instead of a single variable.
     * The history is needed because the liquid balance of some account is not updated immediately, instead it can be updated
     * some time later, whenever the related account sends some transaction that will require the updated liquid balance.
     */
    function _updateUnclaimedReward(address _validator, uint256 _epochID) private {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        
        uint256 _atnReward;
        uint256 _ntnReward;

        // first update last event
        RewardEvent storage _lastRewardEvent = _rewardTracker.lastRewardEvent;
        uint256 _totalLiquid = _lastRewardEvent.totalLiquid;
        if (_lastRewardEvent.eventExist == true && _totalLiquid > 0) {
            uint256 _currentUdpateEpoch = _lastRewardEvent.epochID;
            require(_rewardTracker.lastUpdateEpochID <= _currentUdpateEpoch, "rewards already updated for epoch");
            if (_lastRewardEvent.bonding) {
                (_atnReward, _ntnReward) = autonity.getRewardsTillBonding(_lastRewardEvent.stakingRequestID);
            }
            else {
                (_atnReward, _ntnReward) = autonity.getRewardsTillUnbonding(_lastRewardEvent.stakingRequestID);
            }
            _lastRewardEvent.applied = true;
            _updateLastUnrealisedFeeFactor(_rewardTracker, _currentUdpateEpoch, _atnReward, _ntnReward, _totalLiquid);
            if (_currentUdpateEpoch >= _epochID) {
                return;
            }
        }

        // there is no event after last event, so we update the latest rewards
        require(_getEpochID() > _epochID, "epoch not finalized yet");
        _epochID = _getEpochID() - 1;
        require(_rewardTracker.lastUpdateEpochID <= _epochID, "rewards already updated for epoch");
        Liquid _contract = _rewardTracker.liquidContract;
        _totalLiquid = _contract.balanceOf(address(this));
        (_atnReward, _ntnReward) = _contract.unclaimedRewards(address(this));
        _updateLastUnrealisedFeeFactor(_rewardTracker, _epochID, _atnReward, _ntnReward, _totalLiquid);

        // lastRewardEvent is in the past now
        _lastRewardEvent.eventExist = false;
    }

    /**
     * @dev fetches the unclaimedRewards from liquid contract and calculates the changes in lastUnrealisedFeeFactor
     * but does not update any state variable. This function exists only to support _unclaimedRewards to act as view only function
     */
    function _unfetchedFeeFactor(address _validator) private view returns (uint256 _atnFeeFactor, uint256 _ntnFeeFactor) {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        ILiquidLogic _contract = ILiquidLogic(_rewardTracker.liquidStateContract);
        uint256 _totalLiquid = _contract.balanceOf(address(this));
        if (_totalLiquid > 0) {
            // reward is increased by  OFFSET
            (uint256 _atnReward, uint256 _ntnReward) = _contract.unclaimedRewards(address(this));
            _atnFeeFactor = (_atnReward-_rewardTracker.atnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
            _ntnFeeFactor = (_ntnReward-_rewardTracker.ntnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        }
    }

    function _pendingEventFeeFactor(address _validator) private view returns (uint256 _atnFeeFactor, uint256 _ntnFeeFactor) {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        RewardEvent storage _pendingEvent = _rewardTracker.pendingRewardEvent;
        require(_pendingEvent.eventExist == true, "there must be pending reward event");
        require(_pendingEvent.epochID < _getEpochID(), "pending event should be from past epoch");

        uint256 _totalLiquid = _pendingEvent.totalLiquid;
        if (_totalLiquid > 0) {
            uint256 _atnReward;
            uint256 _ntnReward;

            // getting rewards at the time the staking request was applied
            if (_pendingEvent.bonding == true) {
                (_atnReward, _ntnReward) = autonity.getRewardsTillBonding(_pendingEvent.stakingRequestID);
            }
            else {
                (_atnReward, _ntnReward) = autonity.getRewardsTillUnbonding(_pendingEvent.stakingRequestID);
            }

            _atnFeeFactor = (_atnReward-_rewardTracker.atnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
            _ntnFeeFactor = (_ntnReward-_rewardTracker.ntnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        }
    }

    /**
     * @dev calculates the rewards yet to claim for _id from _validator 
     */
    function _unclaimedRewards(
        uint256 _id,
        address _validator,
        int256 _balanceChange,
        uint256 _updateEpochID
    ) internal view returns (uint256 _atnReward, uint256 _ntnReward) {
        Account storage _account = accounts[_id][_validator];
        (uint256 _atnLastUnrealisedFeeFactor, uint256 _ntnLastUnrealisedFeeFactor) = _unfetchedFeeFactor(_validator);
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        uint256 _lastUpdateEpoch = _rewardTracker.lastUpdateEpochID;
        _atnLastUnrealisedFeeFactor += _rewardTracker.atnLastUnrealisedFeeFactor[_lastUpdateEpoch];
        _ntnLastUnrealisedFeeFactor += _rewardTracker.ntnLastUnrealisedFeeFactor[_lastUpdateEpoch];

        uint256 _balance = _account.liquidBalance;
        _atnReward = _account.atnRealisedFee
                    + _computeUnrealisedFees(
                        _balance, _account.atnUnrealisedFeeFactor, _atnLastUnrealisedFeeFactor
                    );

        _ntnReward = _account.ntnRealisedFee
                    + _computeUnrealisedFees(
                        _balance, _account.ntnUnrealisedFeeFactor, _ntnLastUnrealisedFeeFactor
                    );

        if (_balanceChange == 0) {
            return (_atnReward, _ntnReward);
        }

        require(_updateEpochID < _getEpochID(), "balance cannot change before epoch is finalized");
        // there is a balance change after _updateEpochID due to staking request but they are not applied yet
        // which means there has been no new staking request after _updateEpochID
        // we can also be certain that _updateEpochID < current epoch ID
        uint256 _atnPastUnrealisedFeeFactor;
        uint256 _ntnPastUnrealisedFeeFactor;
        if (_rewardTracker.lastUpdateEpochID >= _updateEpochID) {
            // the balance change takes effect after _updateEpochID
            // so only rewards from current epoch until _updateEpochID is affected by balance change
            _atnPastUnrealisedFeeFactor = _rewardTracker.atnLastUnrealisedFeeFactor[_updateEpochID];
            _ntnPastUnrealisedFeeFactor = _rewardTracker.ntnLastUnrealisedFeeFactor[_updateEpochID];
        }
        else {
            // We have not updated rewards at or after _updateEpochID, so there must be a pending reward event
            // for those staking requests. Otherwise we would already update rewards for _updateEpochID
            // or there should be no staking requests (_balanceChange == 0)
            (_atnPastUnrealisedFeeFactor, _ntnPastUnrealisedFeeFactor) = _pendingEventFeeFactor(_validator);
        }

        if (_balanceChange > 0) {
            _atnReward += _computeUnrealisedFees(
                uint256(_balanceChange), _atnPastUnrealisedFeeFactor, _atnLastUnrealisedFeeFactor
            );
            _ntnReward += _computeUnrealisedFees(
                uint256(_balanceChange), _ntnPastUnrealisedFeeFactor, _ntnLastUnrealisedFeeFactor
            );
        }
    }

    function _isNewBondingRequested(uint256 _id, address _validator) internal view returns (bool) {
        return accounts[_id][_validator].newBondingRequested;
    }

    function _liquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        return accounts[_id][_validator].liquidBalance;
    }

    function _unlockedLiquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        Account storage _account = accounts[_id][_validator];
        return _account.liquidBalance - _account.lockedLiquidBalance;
    }

    function _lockedLiquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        return accounts[_id][_validator].lockedLiquidBalance;
    }

    /**
     * @dev returns the list of validator addresses wich are bonded to schedule _id assigned to _account
     */
    function _bondedValidators(uint256 _id) internal view returns (address[] memory) {
        return bondedValidators[_id];
    }

    function _liquidStateContract(address _validator) internal view returns (address) {
        return rewardTracker[_validator].liquidStateContract;
    }

    function _getEpochID() internal view returns (uint256) {
        return autonity.epochID();
    }

}