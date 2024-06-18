// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";
import "../Liquid.sol";

contract LiquidRewardManager {

    uint256 public constant FEE_FACTOR_UNIT_RECIP = 1_000_000_000;

    address private operator;
    Autonity private autonity;

    /**
     * @dev This structure tracks the activity that requires to update rewards. Any activity like increase or decrease of liquid balances
     * or claiming rewards require updating rewards before the activity is applied. This structure tracks those activities, when such
     * activity is requested.
     * 
     * There are two types of reward event:
     * Pending Reward Event: A reward event that cannot be applied yet. Bonding or unbonding request creates such reward event.
     * Last Reward Event: A reward event that can be applied. When a pending reward event can be applied, after the epcoh is finalized,
     * it becomes last reward event. If there already exists a last reward event, then it is applied before being replaced by a pending reward event.
     */
    struct RewardEvent {
        uint256 epochID;
        uint256 totalLiquid;
        uint256 stakingRequestID;
        bool isBonding;
        bool eventExist;
        bool applied;
    }

    /**
     * @dev Tracks rewards for each validator.
     */
    struct RewardTracker {
        uint256 atnUnclaimedRewards;
        uint256 ntnUnclaimedRewards;
        // offset by 1 to handle empty value
        uint256 lastUpdateEpochID;
        Liquid liquidContract;
        RewardEvent lastRewardEvent;
        RewardEvent pendingRewardEvent;
        mapping(uint256 => uint256) atnLastUnrealisedFeeFactor;
        mapping(uint256 => uint256) ntnLastUnrealisedFeeFactor;
        mapping(uint256 => bool) unrealisedFeeFactorUpdated;
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
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        require(_rewardTracker.unrealisedFeeFactorUpdated[_epochID] == true, "unrealised fee factor not updated");
        uint256 _atnLastUnrealisedFeeFactor = _rewardTracker.atnLastUnrealisedFeeFactor[_epochID];
        Account storage _account = accounts[_id][_validator];
        uint256 _balance = _account.liquidBalance;
        _atnRealisedFees = _account.atnRealisedFee
                        + _computeUnrealisedFees(
                            _balance, _account.atnUnrealisedFeeFactor, _atnLastUnrealisedFeeFactor
                        );
        _account.atnRealisedFee = _atnRealisedFees;
        _account.atnUnrealisedFeeFactor = _atnLastUnrealisedFeeFactor;

        uint256 _ntnLastUnrealisedFeeFactor = _rewardTracker.ntnLastUnrealisedFeeFactor[_epochID];
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
        require(_lastUnrealisedFeeFactor >= _unrealisedFeeFactor, "invalid unrealised fee factor");
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
        _rewardTracker.atnUnclaimedRewards = 0;
        _rewardTracker.ntnUnclaimedRewards = 0;
        
        // cannot apply last reward event anymore
        // because the unclaimed reward from this event is constant, not realtime
        _rewardTracker.lastRewardEvent.eventExist = false;
        _liquidContract(_validator).claimRewards();
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
        (_atnFee, _ntnFee) = _realiseFees(_id, _validator, _getEpochID()-1);
        _account.atnRealisedFee = 0;
        _account.ntnRealisedFee = 0;
    }

    function _initiateValidator(address _validator) private {
        rewardTracker[_validator].liquidContract = autonity.getValidator(_validator).liquidContract;
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
        if (address(rewardTracker[_validator].liquidContract) == address(0)) {
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

    function getPendingRewardEvent(address _validator) public view returns (RewardEvent memory) {
        return rewardTracker[_validator].pendingRewardEvent;
    }

    function getLastRewardEvent(address _validator) public view returns (RewardEvent memory) {
        return rewardTracker[_validator].lastRewardEvent;
    }

    /**
     * @dev If pending rewrad event exists and if the event is not from current epoch, then the pending rewrad event
     * replaces the last reward event.
     */
    function _updateLastRewardEvent(address _validator) private {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        RewardEvent storage _pending = _rewardTracker.pendingRewardEvent;
        if (_pending.eventExist == true && _pending.epochID < _getEpochID()) {
            RewardEvent storage _last = _rewardTracker.lastRewardEvent;
            require(_last.eventExist == false || _last.applied == true, "last event needs to be applied before being replaced");
            _rewardTracker.lastRewardEvent = _rewardTracker.pendingRewardEvent;
            _pending.eventExist = false;
        }
    }

    function _newPendingRewardEvent(address _validator, uint256 _stakingRequestID, bool _isBonding) internal {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        if (_rewardTracker.lastRewardEvent.eventExist == true && _rewardTracker.lastRewardEvent.applied == false) {
            _updateUnclaimedReward(_validator, _rewardTracker.lastRewardEvent.epochID);
        }

        _updateLastRewardEvent(_validator);

        // if there is already a pending event, replacing it is fine
        _rewardTracker.pendingRewardEvent = RewardEvent(
            _getEpochID(),
            _liquidContract(_validator).balanceOf(address(this)),
            _stakingRequestID,
            _isBonding,
            true,
            false
        );
    }

    function _updatePendingEventLiquid(address _validator) internal {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        RewardEvent storage _pending = _rewardTracker.pendingRewardEvent;
        if (_pending.eventExist == true && _pending.epochID == _getEpochID()) {
            _pending.totalLiquid = _liquidContract(_validator).balanceOf(address(this));
        }
    }

    function _updateLastUnrealisedFeeFactor(
        RewardTracker storage _rewardTracker,
        uint256 _epochID,
        uint256 _atnReward,
        uint256 _ntnReward,
        uint256 _totalLiquid
    ) private {
        mapping(uint256 => uint256) storage _lastUnrealisedFeeFactor = _rewardTracker.atnLastUnrealisedFeeFactor;

        if (_rewardTracker.lastUpdateEpochID > 0) {
            _lastUnrealisedFeeFactor[_epochID] = _lastUnrealisedFeeFactor[_rewardTracker.lastUpdateEpochID-1];
        }

        if (_totalLiquid > 0) {
            _lastUnrealisedFeeFactor[_epochID] += (_atnReward-_rewardTracker.atnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        }

        _rewardTracker.atnUnclaimedRewards = _atnReward;

        _lastUnrealisedFeeFactor = _rewardTracker.ntnLastUnrealisedFeeFactor;

        if (_rewardTracker.lastUpdateEpochID > 0) {
            _lastUnrealisedFeeFactor[_epochID] = _lastUnrealisedFeeFactor[_rewardTracker.lastUpdateEpochID-1];
        }

        if (_totalLiquid > 0) {
            _lastUnrealisedFeeFactor[_epochID] += (_ntnReward-_rewardTracker.ntnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        }

        // fee factor updated
        _rewardTracker.unrealisedFeeFactorUpdated[_epochID] = true;

        _rewardTracker.ntnUnclaimedRewards = _ntnReward;
        // offset by 1 to handle empty value
        _rewardTracker.lastUpdateEpochID = _epochID+1;
    }

    function _applyLastRewardEvent(RewardTracker storage _rewardTracker) private {

        RewardEvent storage _lastRewardEvent = _rewardTracker.lastRewardEvent;

        if (_lastRewardEvent.eventExist == true) {
            uint256 _currentUdpateEpoch = _lastRewardEvent.epochID;
            require(_rewardTracker.lastUpdateEpochID <= _currentUdpateEpoch+1, "reward update event is passed");

            uint256 _atnReward;
            uint256 _ntnReward;
            if (_lastRewardEvent.isBonding) {
                (_atnReward, _ntnReward) = autonity.getRewardsTillBonding(_lastRewardEvent.stakingRequestID);
            }
            else {
                (_atnReward, _ntnReward) = autonity.getRewardsTillUnbonding(_lastRewardEvent.stakingRequestID);
            }

            _lastRewardEvent.applied = true;
            _updateLastUnrealisedFeeFactor(_rewardTracker, _currentUdpateEpoch, _atnReward, _ntnReward, _lastRewardEvent.totalLiquid);
        }

    }

    /**
     * @dev Updates the unclaimed rewards from validator and the lastUnrealisedFeeFactor which is used
     * to compute unrealised fees for accounts. Both is updated until given epoch id. The lastUnrealisedFeeFactor
     * is kept with history, so we have a mapping(epochID => value) of lastUnrealisedFeeFactor instead of a single variable.
     * The history is needed because the liquid balance of some account is not updated immediately, instead it can be updated
     * some time later, whenever the related account sends some transaction that will require the updated liquid balance.
     * @param _validator validator address, from which we will claim rewards
     * @param _epochID the epochID untill which we need to fetch rewards
     * 
     * To update unclaimed rewards, first we need to apply the last reward event (see: struct RewardEvent). Then if there is a
     * pending reward event from some past epoch, it replaces the current last reward event. Then we apply the new last reward
     * event again. After that, if we are still behind the input epochID, then we fetch the last updated rewards.
     */
    function _updateUnclaimedReward(address _validator, uint256 _epochID) private {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        _applyLastRewardEvent(_rewardTracker);

        if (_rewardTracker.pendingRewardEvent.eventExist == true) {
            _updateLastRewardEvent(_validator);
            _applyLastRewardEvent(_rewardTracker);
        }

        if (_rewardTracker.lastUpdateEpochID >= _epochID+1) {
            return;
        }

        // there is no event after last event, so we update the latest rewards

        // for every input epochID we need to have unrealisedFeeFactor for that epochID
        // as we are fetching the last reward, input epochID needs to match last epochID
        // for which rewards are distributed
        require(_getEpochID()-1 == _epochID, "cannot update rewards for input epochID");
        Liquid _contract = _liquidContract(_validator);
        (uint256 _atnReward, uint256 _ntnReward) = _contract.unclaimedRewards(address(this));
        _updateLastUnrealisedFeeFactor(_rewardTracker, _epochID, _atnReward, _ntnReward, _contract.balanceOf(address(this)));

        // lastRewardEvent is in the past now
        // cannot be applied anymore
        _rewardTracker.lastRewardEvent.eventExist = false;
    }

    /**
     * @dev Fetches the unclaimedRewards from liquid contract and calculates the changes in lastUnrealisedFeeFactor
     * but does not update any state variable. This function helps to calculate unclaimed rewards.
     */
    function _unfetchedFeeFactor(
        Liquid _contract,
        uint256 _atnLastReward,
        uint256 _ntnLastReward
    ) private view returns (
        uint256 _atnFeeFactor,
        uint256 _ntnFeeFactor
    ) {
        require(address(_contract) != address(0), "validator not initiated");
        uint256 _totalLiquid = _contract.balanceOf(address(this));
        if (_totalLiquid > 0) {
            (uint256 _atnReward, uint256 _ntnReward) = _contract.unclaimedRewards(address(this));
            _atnFeeFactor = (_atnReward-_atnLastReward) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
            _ntnFeeFactor = (_ntnReward-_ntnLastReward) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        }
    }

    /**
     * @dev Applies the reward event without changing the state and generates rewards.
     * This function helps to calculate unclaimed rewards.
     * @param _rewardEvent Reward event to apply
     */
    function _rewardEventSimulator(
        RewardEvent storage _rewardEvent,
        uint256 _atnLastReward,
        uint256 _ntnLastReward
    ) private view returns (
        uint256 _atnNewReward,
        uint256 _ntnNewReward,
        uint256 _atnFeeFactor,
        uint256 _ntnFeeFactor
    ) {
        // getting rewards at the time the staking request was applied
        if (_rewardEvent.isBonding == true) {
            (_atnNewReward, _ntnNewReward) = autonity.getRewardsTillBonding(_rewardEvent.stakingRequestID);
        }
        else {
            (_atnNewReward, _ntnNewReward) = autonity.getRewardsTillUnbonding(_rewardEvent.stakingRequestID);
        }

        uint256 _totalLiquid = _rewardEvent.totalLiquid;
        if (_totalLiquid > 0) {
            _atnFeeFactor = (_atnNewReward-_atnLastReward) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
            _ntnFeeFactor = (_ntnNewReward-_ntnLastReward) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        }
    }

    /**
     * @dev Generates unrealised fee factor for validator till the rewards for epochID have been distributed.
     * This function helps to calculate unclaimed rewards.
     * @param _validator validator address
     * @param _epochID epochID
     */
    function _generateUnrealisedFeeFactor(
        address _validator,
        uint256 _epochID
    ) private view returns (
        uint256 _atnUnrealisedFeeFactor,
        uint256 _ntnUnrealisedFeeFactor
    ) {

        require(_epochID < _getEpochID(), "epoch not finalized yet");
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        uint256 _lastRewardUpdate = _rewardTracker.lastUpdateEpochID;
        if (_lastRewardUpdate >= _epochID+1) {
            require(_rewardTracker.unrealisedFeeFactorUpdated[_epochID] == true, "unrealised fee factor not updated");
            return (_rewardTracker.atnLastUnrealisedFeeFactor[_epochID], _rewardTracker.ntnLastUnrealisedFeeFactor[_epochID]);
        }

        if (_lastRewardUpdate > 0) {
            _atnUnrealisedFeeFactor = _rewardTracker.atnLastUnrealisedFeeFactor[_lastRewardUpdate-1];
            _ntnUnrealisedFeeFactor = _rewardTracker.ntnLastUnrealisedFeeFactor[_lastRewardUpdate-1];
        }
        
        uint256 _atnReward = _rewardTracker.atnUnclaimedRewards;
        uint256 _ntnReward = _rewardTracker.ntnUnclaimedRewards;
        uint256 _atnFeeFactorFetched;
        uint256 _ntnFeeFactorFetched;

        // apply last reward event
        RewardEvent storage _rewardEvent = _rewardTracker.lastRewardEvent;
        if (_rewardEvent.eventExist == true && _lastRewardUpdate < _rewardEvent.epochID+1) {
            _lastRewardUpdate = _rewardEvent.epochID+1;
            require(_epochID >= _rewardEvent.epochID, "cannot generate unrealised fee factor for epochID from last reward event");
            (_atnReward, _ntnReward, _atnFeeFactorFetched, _ntnFeeFactorFetched) = _rewardEventSimulator(_rewardEvent, _atnReward, _ntnReward);
            _atnUnrealisedFeeFactor += _atnFeeFactorFetched;
            _ntnUnrealisedFeeFactor += _ntnFeeFactorFetched;

            if (_epochID == _rewardEvent.epochID) {
                return (_atnUnrealisedFeeFactor, _ntnUnrealisedFeeFactor);
            }
        }

        // apply pending reward event
        _rewardEvent = _rewardTracker.pendingRewardEvent;
        if (_rewardEvent.eventExist == true && _rewardEvent.epochID < _getEpochID() && _lastRewardUpdate < _rewardEvent.epochID+1) {
            _lastRewardUpdate = _rewardEvent.epochID+1;
            require(_epochID >= _rewardEvent.epochID, "cannot generate unrealised fee factor for epochID from pending reward event");
            (_atnReward, _ntnReward, _atnFeeFactorFetched, _ntnFeeFactorFetched) = _rewardEventSimulator(_rewardEvent, _atnReward, _ntnReward);
            _atnUnrealisedFeeFactor += _atnFeeFactorFetched;
            _ntnUnrealisedFeeFactor += _ntnFeeFactorFetched;

            if (_epochID == _rewardEvent.epochID) {
                return (_atnUnrealisedFeeFactor, _ntnUnrealisedFeeFactor);
            }
        }

        // need to fetch latest reward update
        require(_epochID == _getEpochID()-1, "cannot generate unrealised fee factor");
        (_atnFeeFactorFetched, _ntnFeeFactorFetched) = _unfetchedFeeFactor(_rewardTracker.liquidContract, _atnReward, _ntnReward);
        _atnUnrealisedFeeFactor += _atnFeeFactorFetched;
        _ntnUnrealisedFeeFactor += _ntnFeeFactorFetched;
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
        if (_account.liquidBalance == 0 && _balanceChange == 0) {
            return (_account.atnRealisedFee, _account.ntnRealisedFee);
        }

        uint256 _lastEpochID = _getEpochID()-1;
        if (_updateEpochID > 0) {
            _lastEpochID = _updateEpochID-1;
        }
        uint256 _balance = _account.liquidBalance;
        (uint256 _atnMidUnrealisedFeeFactor, uint256 _ntnMidUnrealisedFeeFactor) = _generateUnrealisedFeeFactor(_validator, _lastEpochID);

        _atnReward = _account.atnRealisedFee
                    + _computeUnrealisedFees(
                        _balance, _account.atnUnrealisedFeeFactor, _atnMidUnrealisedFeeFactor
                    );

        _ntnReward = _account.ntnRealisedFee
                    + _computeUnrealisedFees(
                        _balance, _account.ntnUnrealisedFeeFactor, _ntnMidUnrealisedFeeFactor
                    );

        if (_updateEpochID == 0) {
            return (_atnReward, _ntnReward);
        }
        

        require(int256(_account.liquidBalance) + _balanceChange >= 0, "balance change is wrong");
        _balance = uint256(int256(_account.liquidBalance) + _balanceChange);

        require(_updateEpochID <= _getEpochID(), "balance cannot change before epoch is finalized");
        (uint256 _atnLastUnrealisedFeeFactor, uint256 _ntnLastUnrealisedFeeFactor) = _generateUnrealisedFeeFactor(_validator, _getEpochID()-1);

        _atnReward += _computeUnrealisedFees(
            _balance, _atnMidUnrealisedFeeFactor, _atnLastUnrealisedFeeFactor
        );

        _ntnReward += _computeUnrealisedFees(
            _balance, _ntnMidUnrealisedFeeFactor, _ntnLastUnrealisedFeeFactor
        );
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

    function _liquidContract(address _validator) internal returns (Liquid) {
        RewardTracker storage _rewardTracker = rewardTracker[_validator];
        if (address(_rewardTracker.liquidContract) == address(0)) {
            _initiateValidator(_validator);
        }
        return _rewardTracker.liquidContract;
    }

    function _getEpochID() internal view returns (uint256) {
        return autonity.epochID();
    }

}