// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../interfaces/IStakeProxy.sol";
import "./LiquidRewardManager.sol";
import "./ScheduleBase.sol";

contract StakableVesting is IStakeProxy, ScheduleBase, LiquidRewardManager {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release schedule
    // UNLOCKED are tokens that can be withdrawn
    uint256 public contractVersion = 1;
    // TODO: review: a way to measure requiredGasBond and requiredGasUnbond realtime?
    uint256 private requiredGasBond = 50_000;
    uint256 private requiredGasUnbond = 50_000;


    struct PendingBondingRequest {
        uint256 amount;
        uint256 epochID;
        address validator;
        bool processed;
    }

    struct PendingUnbondingRequest {
        uint256 liquidAmount;
        uint256 epochID;
        address validator;
        bool rejected;
        bool applied;
    }

    mapping(uint256 => PendingBondingRequest) private pendingBondingRequest;
    mapping(uint256 => PendingUnbondingRequest) private pendingUnbondingRequest;

    /**
     * @dev bondingToSchedule is needed to handle notification from autonity when bonding is applied.
     * In case the it fails to notify the vesting contract, scheduleToBonding is needed to revert the failed requests.
     * All bonding requests are applied at epoch end, so we can process all of them (failed or successful) together.
     * See bond and _revertPendingBondingRequest for more clarity
     */

    mapping(uint256 => uint256) private bondingToSchedule;
    mapping(uint256 => uint256[]) private scheduleToBonding;

    /**
     * @dev unbondingToSchedule is needed to handle notification from autonity when unbonding is applied and released.
     * In case it fails to notify the vesting contract, scheduleToUnbonding is needed to revert the failed requests.
     * Not all requests are released together at epoch end, so we cannot process all the request together.
     * tailPendingUnbondingID and headPendingUnbondingID helps to keep track of scheduleToUnbonding.
     * See unbond and _revertPendingUnbondingRequest for more clarity
     */

    mapping(uint256 => uint256) private unbondingToSchedule;
    mapping(uint256 => mapping(uint256 => uint256)) private scheduleToUnbonding;
    mapping(uint256 => uint256) private tailPendingUnbondingID;
    mapping(uint256 => uint256) private headPendingUnbondingID;

    // AUT rewards entitled to some beneficiary for bonding from some schedule before it has been cancelled
    // see cancelSchedule for more clarity.
    mapping(address => uint256) private atnRewards;
    mapping(address => uint256) private ntnRewards;

    constructor(
        address payable _autonity, address _operator
    ) LiquidRewardManager(_autonity) ScheduleBase(_autonity, _operator) {}

    /**
     * @notice creates a new stakable schedule, restricted to operator only
     * @dev _amount NTN must be minted and transferred to this contract before creating the schedule
     * otherwise the schedule cannot be released or bonded to some validator
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _startTime start time of the vesting
     * @param _cliffTime cliff time of the vesting, cliff period = _cliffTime - _startTime
     * @param _endTime end time of the vesting, total duration of the schedule = _endTime - _startTime
     */
    function newSchedule(
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffTime,
        uint256 _endTime
    ) virtual onlyOperator public {
        _createSchedule(_beneficiary, _amount, _startTime, _cliffTime, _endTime, true);
    }


    /**
     * @notice used by beneficiary to transfer all unlocked NTN and LNTN of some schedule to his own address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (excluding canceled ones)
     * So any beneficiary can number their schedules from 0 to (n-1). Beneficiary does not need to know the 
     * unique global schedule id which can be retrieved via _getUniqueScheduleID function
     */
    function releaseFunds(uint256 _id) virtual external {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        uint256 _unlocked = _unlockedFunds(_scheduleID);
        // first NTN is released
        _unlocked = _releaseNTN(_scheduleID, _unlocked);
        // if there still remains some unlocked funds, i.e. not enough NTN, then LNTN is released
        _releaseAllUnlockedLNTN(_scheduleID, _unlocked);
    }

    /**
     * @notice used by beneficiary to transfer all unlocked NTN of some schedule to his own address
     */
    function releaseAllNTN(uint256 _id) virtual external {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        _releaseNTN(_scheduleID, _unlockedFunds(_scheduleID));
    }

    /**
     * @notice used by beneficiary to transfer all unlocked LNTN of some schedule to his own address
     */
    function releaseAllLNTN(uint256 _id) virtual external {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        _releaseAllUnlockedLNTN(_scheduleID, _unlockedFunds(_scheduleID));
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice used by beneficiary to transfer some amount of unlocked NTN of some schedule to his own address
     * @param _amount amount to transfer
     */
    function releaseNTN(uint256 _id, uint256 _amount) virtual external {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        require(_amount <= _unlockedFunds(_scheduleID), "not enough unlocked funds");
        _releaseNTN(_scheduleID, _amount);
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice used by beneficiary to transfer some amount of unlocked LNTN of some schedule to his own address
     * @param _validator address of the validator
     * @param _amount amount of LNTN to transfer
     */
    function releaseLNTN(uint256 _id, address _validator, uint256 _amount) virtual external {
        require(_amount > 0, "require positive amount to transfer");
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);

        uint256 _unlockedLiquid = _unlockedLiquidBalanceOf(_scheduleID, _validator);
        require(_unlockedLiquid >= _amount, "not enough unlocked LNTN");

        uint256 _value = _calculateLNTNValue(_validator, _amount);
        require(_value <= _unlockedFunds(_scheduleID), "not enough unlocked funds");

        Schedule storage _schedule = schedules[_scheduleID];
        _schedule.withdrawnValue += _value;
        _updateAndTransferLNTN(_scheduleID, msg.sender, _amount, _validator);
    }

    /**
     * @notice changes the beneficiary of some schedule to the _recipient address. _recipient can release and stake tokens from the schedule.
     * only operator is able to call this function.
     * rewards which have been entitled to the beneficiary due to bonding from this schedule are not transferred to _recipient
     * @dev rewards earned until this point from this schedule are calculated and stored in atnRewards and ntnRewards mapping so that
     * _beneficiary can later claim them even though _beneficiary is not entitled to this schedule.
     * @param _beneficiary beneficiary address whose schedule will be canceled
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (excluding already canceled ones)
     * @param _recipient whome the schedule is transferred to
     */
    function cancelSchedule(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        (uint256 _atnReward, uint256 _ntnReward) = _claimRewards(_scheduleID);
        atnRewards[_beneficiary] += _atnReward;
        ntnRewards[_beneficiary] += _ntnReward;
        _cancelSchedule(_beneficiary, _id, _recipient);
    }

    /**
     * @notice In case some funds are missing due to some pending staking operation that failed,
     * this function updates the funds of some schedule _id entitled to _beneficiary by reverting the failed requests
     */
    function updateFunds(address _beneficiary, uint256 _id) virtual external {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        _revertPendingBondingRequest(_scheduleID);
        _revertPendingUnbondingRequest(_scheduleID);
    }

    function setRequiredGasBond(uint256 _gas) external onlyOperator {
        requiredGasBond = _gas;
    }

    function setRequiredGasUnbond(uint256 _gas) external onlyOperator {
        requiredGasUnbond = _gas;
    }

    /**
     * @notice Used by beneficiary to bond some NTN of some schedule _id.
     * All bondings are delegated, as vesting manager cannot own a validator
     * @param _id id of the schedule numbered from 0 to (n-1) where n = total schedules entitled to the beneficiary (excluding canceled ones)
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(uint256 _id, address _validator, uint256 _amount) virtual public payable returns (uint256) {
        // TODO: do we need to wait till _schedule.start before bonding??
        require(msg.value >= requiredBondingGasCost(), "not enough gas given for notification on bonding");
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.start <= block.timestamp, "schedule not started yet");
        require(_schedule.currentNTNAmount >= _amount, "not enough tokens");

        uint256 _bondingID = autonity.bond{value: msg.value}(_validator, _amount);
        _schedule.currentNTNAmount -= _amount;
        // offset by 1 to handle empty value
        scheduleToBonding[_scheduleID].push(_bondingID);
        bondingToSchedule[_bondingID] = _scheduleID+1;
        pendingBondingRequest[_bondingID] = PendingBondingRequest(_amount, _epochID(), _validator, false);
        _initiate(_scheduleID, _validator);
        return _bondingID;
    }

    /**
     * @notice Used by beneficiary to unbond some LNTN of some schedule.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public payable returns (uint256) {
        require(msg.value >= requiredUnbondingGasCost(), "not enough gas given for notification on unbonding");
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        _lock(_scheduleID, _validator, _amount);
        uint256 _unbondingID = autonity.unbond{value: msg.value}(_validator, _amount);
        // offset by 1 to handle empty value
        unbondingToSchedule[_unbondingID] = _scheduleID+1;
        pendingUnbondingRequest[_unbondingID] = PendingUnbondingRequest(_amount, _epochID(), _validator, false, false);
        uint256 _lastID = headPendingUnbondingID[_scheduleID];
        // scheduleToUnbonding[_scheduleID][_i] stores the _unbondingID of the i'th unbonding request
        scheduleToUnbonding[_scheduleID][_lastID] = _unbondingID;
        headPendingUnbondingID[_scheduleID] = _lastID+1;
        return _unbondingID;
    }

    /**
     * @notice used by beneficiary to claim all rewards which is entitled due to bonding
     * @dev Rewards from some cancelled schedules are stored in rewards mapping. All rewards from
     * schedules that are still entitled to the beneficiary need to be calculated via _claimRewards
     */
    function claimRewards() virtual external {
        uint256[] storage _scheduleIDs = beneficiarySchedules[msg.sender];
        uint256 _atnTotalFees = atnRewards[msg.sender];
        uint256 _ntnTotalFees = ntnRewards[msg.sender];
        atnRewards[msg.sender] = 0;
        ntnRewards[msg.sender] = 0;
        
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            _cleanup(_scheduleIDs[i]);
            (uint256 _atnReward, uint256 _ntnReward) = _claimRewards(_scheduleIDs[i]);
            _atnTotalFees += _atnReward;
            _ntnTotalFees += _ntnReward;
        }
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool _sent, ) = msg.sender.call{value: _atnTotalFees}("");
        require(_sent, "Failed to send AUT");

        _sent = autonity.transfer(msg.sender, _ntnTotalFees);
        require(_sent, "Failed to send NTN");
    }

    /**
     * @notice can be used to send AUT to the contract
     */
    function receiveAut() external payable {
        // do nothing
    }

    /**
     * @notice callback function restricted to autonity to notify the vesting contract to update rewards (AUT) for _validators
     * @param _validators address of the validators that have staking operations and need to update their rewards (AUT)
     */
    function rewardsDistributed(address[] memory _validators) external onlyAutonity {
        _updateValidatorReward(_validators);
    }

    /**
     * @notice implements IStakeProxy.bondingApplied(), a callback function for autonity when bonding is applied
     * @param _bondingID bonding id from Autonity when bonding was requested
     * @param _liquid amount of LNTN after bonding applied successfully
     * @param _selfDelegation true if self bonded, false for delegated bonding
     * @param _rejected true if bonding request was rejected, false if applied successfully
     */
    function bondingApplied(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) external onlyAutonity {
        _applyBonding(_bondingID, _validator, _liquid, _selfDelegation, _rejected);
    }

    /**
     * @notice implements IStakeProxy.unbondingApplied(). callback function for autonity when unbonding is applied
     * @param _unbondingID unbonding id from Autonity when unbonding was requested
     * @param _rejected true if unbonding was rejected, false if applied successfully
     */
    function unbondingApplied(uint256 _unbondingID, address _validator, bool _rejected) external onlyAutonity {
        _applyUnbonding(_unbondingID, _validator, _rejected);
    }

    /**
     * @dev implements IStakeProxy.unbondingReleased(). callback function for autonity when unbonding is released
     * @param _unbondingID unbonding id from Autonity when unbonding was requested
     * @param _amount amount of NTN released
     * @param _rejected true if unbonding was rejected, false if applied and released successfully
     */
    function unbondingReleased(uint256 _unbondingID, uint256 _amount, bool _rejected) external onlyAutonity {
        _releaseUnbonding(_unbondingID, _amount, _rejected);
    }

    /**
     * @dev returns equivalent amount of NTN using the ratio.
     * @param _validator validator address
     * @param _amount amount of LNTN to be converted
     */
    function _calculateLNTNValue(address _validator, uint256 _amount) private view returns (uint256) {
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake) / _validatorInfo.liquidSupply;
    }

    /**
     * @dev returns equivalent amount of LNTN using the ratio.
     * @param _validator validator address
     * @param _amount amount of NTN to be converted
     */
    function _getLiquidFromNTN(address _validator, uint256 _amount) private view returns (uint256) {
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * _validatorInfo.liquidSupply / (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake);
    }

    /**
     * @dev calculates the total value of the schedule, which can vary if the schedule has some LNTN
     * total value = current_NTN + withdrawn_value + (the value of LNTN converted to NTN)
     * @param _scheduleID unique global id of the schedule
     */
    function _calculateTotalValue(uint256 _scheduleID) private view returns (uint256) {
        Schedule storage _schedule = schedules[_scheduleID];
        uint256 _totalValue = _schedule.currentNTNAmount + _schedule.withdrawnValue;
        address[] memory _validators = _bondedValidators(_scheduleID);
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _balance = _liquidBalanceOf(_scheduleID, _validators[i]);
            if (_balance == 0) {
                continue;
            }
            _totalValue += _calculateLNTNValue(_validators[i], _balance);
        }
        return _totalValue;
    }

    /**
     * @dev transfers some LNTN equivalent to _availableUnlockedFunds NTN to beneficiary address.
     * In case the _scheduleID has LNTN to multiple validators, we pick one validator and try to transfer
     * as much LNTN as possible. If there still remains some more uncloked funds, then we pick another validator.
     * There is no particular order in which validator should be picked first.
     */
    function _releaseAllUnlockedLNTN(
        uint256 _scheduleID, uint256 _availableUnlockedFunds
    ) private returns (uint256 _remaining) {
        _remaining = _availableUnlockedFunds;
        address[] memory _validators = _bondedValidators(_scheduleID);
        for (uint256 i = 0; i < _validators.length && _remaining > 0; i++) {
            uint256 _balance = _unlockedLiquidBalanceOf(_scheduleID, _validators[i]);
            if (_balance == 0) {
                continue;
            }
            uint256 _value = _calculateLNTNValue(_validators[i], _balance);
            if (_remaining >= _value) {
                _remaining -= _value;
                _updateAndTransferLNTN(_scheduleID, msg.sender, _balance, _validators[i]);
            }
            else {
                uint256 _liquid = _getLiquidFromNTN(_validators[i], _remaining);
                require(_liquid <= _balance, "conversion not working");
                _remaining = 0;
                _updateAndTransferLNTN(_scheduleID, msg.sender, _liquid, _validators[i]);
            }
        }
        Schedule storage _schedule = schedules[_scheduleID];
        _schedule.withdrawnValue += _availableUnlockedFunds - _remaining;
    }

    /**
     * @dev calculates the amount of unlocked funds in NTN until last epoch block time
     */
    function _unlockedFunds(uint256 _scheduleID) private view returns (uint256) {
        return _calculateAvailableUnlockedFunds(
            _scheduleID, _calculateTotalValue(_scheduleID), autonity.lastEpochTime()
        );
    }

    function _updateAndTransferLNTN(uint256 _scheduleID, address _to, uint256 _amount, address _validator) private {
        _updateUnclaimedReward(_validator);
        _decreaseLiquid(_scheduleID, _validator, _amount);
        _transferLNTN(_to, _amount, _validator);
    }

    function _transferLNTN(address _to, uint256 _amount, address _validator) private {
        bool _sent = _liquidContract(_validator).transfer(_to, _amount);
        require(_sent, "LNTN transfer failed");
    }

    function _updateValidatorReward(address[] memory _validators) internal {
        for (uint256 i = 0; i < _validators.length; i++) {
            _updateUnclaimedReward(_validators[i]);
        }
    }

    /**
     * @dev mimic _applyBonding from Autonity.sol
     */
    function _applyBonding(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) internal {
        require(_selfDelegation == false, "bonding should be delegated");
        uint256 _scheduleID = bondingToSchedule[_bondingID] - 1;

        PendingBondingRequest storage _request = pendingBondingRequest[_bondingID];
        _request.processed = true;

        if (_rejected) {
            Schedule storage _schedule = schedules[_scheduleID];
            _schedule.currentNTNAmount += _request.amount;
        }
        else {
            _increaseLiquid(_scheduleID, _validator, _liquid);
        }
    }

    /**
     * @dev mimic _applyUnbonding from Autonity.sol
     */
    function _applyUnbonding(uint256 _unbondingID, address _validator, bool _rejected) internal {
        uint256 _scheduleID = unbondingToSchedule[_unbondingID] - 1;
        PendingUnbondingRequest storage _unbondingRequest = pendingUnbondingRequest[_unbondingID];
        uint256 _liquid = _unbondingRequest.liquidAmount;
        _unlock(_scheduleID, _validator, _liquid);

        if (_rejected) {
            _unbondingRequest.rejected = true;
            return;
        }

        _unbondingRequest.applied = true;
        _decreaseLiquid(_scheduleID, _validator, _liquid);
    }

    function _releaseUnbonding(uint256 _unbondingID, uint256 _amount, bool _rejected) internal {
        uint256 _scheduleID = unbondingToSchedule[_unbondingID] - 1;

        if (_rejected) {
            // If released of unbonding is rejected, then it is assumed that the applying of unbonding was also rejected
            // or reverted at Autonity because Autonity could not notify us (_applyUnbonding reverted).
            // If applying of unbonding was successful, then releasing of unbonding cannot be rejected.
            // Here we assume that it was rejected at vesting manager as well, otherwise if it was reverted (due to out of gas)
            // it will revert here as well. In any case, _revertPendingUnbondingRequest will handle the reverted or rejected request
            return;
        }
        Schedule storage _schedule = schedules[_scheduleID];
        _schedule.currentNTNAmount += _amount;
    }

    /**
     * @dev The schedule needs to be cleaned before bonding, unbonding or claiming rewards.
     * _cleanup removes any unnecessary validator from the list, removes pending bonding or unbonding requests
     * that have been rejected or reverted but vesting manager could not be notified. If the clean up is not
     * done, then liquid balance could be incorrect due to not handling the bonding or unbonding request.
     * @param _scheduleID unique global schedule id
     */
    function _cleanup(uint256 _scheduleID) private {
        _revertPendingBondingRequest(_scheduleID);
        _revertPendingUnbondingRequest(_scheduleID);
        _clearValidators(_scheduleID);
    }

    /**
     * @dev in case some bonding request from some previous epoch was unsuccessful and vesting contract was not notified,
     * this function handles such requests. All the requests from past epoch can be handled as the bonding requests are
     * applied at epoch end immediately. Requests from current epoch are not handled.
     * @param _scheduleID unique global id of the schedule
     */
    function _revertPendingBondingRequest(uint256 _scheduleID) private {
        uint256[] storage _bondingIDs = scheduleToBonding[_scheduleID];
        uint256 _length = _bondingIDs.length;
        uint256 _oldBondingID;
        PendingBondingRequest storage _oldBondingRequest;
        uint256 _totalAmount = 0;
        uint256 _currentEpochID = _epochID();
        for (uint256 i = 0; i < _length; i++) {
            _oldBondingID = _bondingIDs[i];
            _oldBondingRequest = pendingBondingRequest[_oldBondingID];
            // will revert request from some previous epoch, request from current epoch will not be reverted
            if (_oldBondingRequest.epochID == _currentEpochID) {
                // all the request in the array are from current epoch
                return;
            }
            _bondingRequestExpired(_scheduleID, _oldBondingRequest.validator);
            // if the request is not processed successfully, then we need to revert it
            if (_oldBondingRequest.processed == false) {
                _totalAmount += _oldBondingRequest.amount;
            }
            delete pendingBondingRequest[_oldBondingID];
            delete bondingToSchedule[_oldBondingID];
        }
        delete scheduleToBonding[_scheduleID];
        if (_totalAmount == 0) {
            return;
        }
        Schedule storage _schedule = schedules[_scheduleID];
        _schedule.currentNTNAmount += _totalAmount;
    }

    /**
     * @dev in case some unbonding request from some previous epoch was unsuccessful (not applied successfully or
     * not released successfully) and vesting contract was not notified, this function handles such requests.
     * Any request that has been processed in _releaseUnbondingStake function can be handled here.
     * Other requests need to wait.
     * @param _scheduleID unique global id of the schedule
     */
    function _revertPendingUnbondingRequest(uint256 _scheduleID) private {
        uint256 _unbondingID;
        PendingUnbondingRequest storage _unbondingRequest;
        mapping(uint256 => uint256) storage _unbondingIDs = scheduleToUnbonding[_scheduleID];
        uint256 _lastID = headPendingUnbondingID[_scheduleID];
        uint256 _processingID = tailPendingUnbondingID[_scheduleID];
        for(; _processingID < _lastID; _processingID++) {
            _unbondingID = _unbondingIDs[_processingID];
            Autonity.UnbondingReleaseState _releaseState = autonity.getUnbondingReleaseState(_unbondingID);
            if (_releaseState == Autonity.UnbondingReleaseState.notReleased) {
                // all the rest have not been released yet
                break;
            }
            delete _unbondingIDs[_processingID];
            delete unbondingToSchedule[_unbondingID];
            
            if (_releaseState == Autonity.UnbondingReleaseState.released) {
                // unbonding was released successfully
                delete pendingUnbondingRequest[_unbondingID];
                continue;
            }

            _unbondingRequest = pendingUnbondingRequest[_unbondingID];
            address _validator = _unbondingRequest.validator;

            if (_releaseState == Autonity.UnbondingReleaseState.reverted) {
                // it means the unbonding was released, but later reverted due to failing to notify us
                // that means unbonding was applied successfully
                require(_unbondingRequest.applied, "unbonding was released but not applied succesfully");
                _updateUnclaimedReward(_validator);
                _increaseLiquid(_scheduleID, _validator, autonity.getRevertingAmount(_unbondingID));
            }
            else if (_releaseState == Autonity.UnbondingReleaseState.rejected) {
                require(_unbondingRequest.applied == false, "unbonding was applied successfully but release was rejected");
                // _unbondingRequest.rejected = true means we already rejected it
                if (_unbondingRequest.rejected == false) {
                    _unlock(_scheduleID, _validator, _unbondingRequest.liquidAmount);
                }
            }
            
            delete pendingUnbondingRequest[_unbondingID];
        }
        tailPendingUnbondingID[_scheduleID] = _processingID;
    }

    function _epochID() private view returns (uint256) {
        return autonity.epochID();
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice returns the gas cost required in "wei" to notify vesting manager when the bonding is applied
     */
    function requiredBondingGasCost() public view returns (uint256) {
        return requiredGasBond * autonity.stakingGasPrice();
    }

    /**
     * @notice returns the gas cost required in "wei" to notify vesting manager when the unbonding is applied and released
     */
    function requiredUnbondingGasCost() public view returns (uint256) {
        return requiredGasUnbond * autonity.stakingGasPrice();
    }

    /**
     * @notice returns the amount of all unclaimed rewards due to all the bonding from schedules entitled to beneficiary
     * @param _beneficiary beneficiary address
     */
    function unclaimedRewards(address _beneficiary) virtual external view returns (uint256 _atnTotalFee, uint256 _ntnTotalFee) {
        _atnTotalFee = atnRewards[_beneficiary];
        _ntnTotalFee = ntnRewards[_beneficiary];
        uint256[] storage _scheduleIDs = beneficiarySchedules[_beneficiary];
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            (uint256 _atnFee, uint256 _ntnFee) = _unclaimedRewards(_scheduleIDs[i]);
            _atnTotalFee += _atnFee;
            _ntnTotalFee += _ntnFee;
        }
    }

    /**
     * @notice returns the amount of LNTN for some schedule
     * @param _beneficiary beneficiary address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (excluding canceled ones)
     * @param _validator validator address
     */
    function liquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _liquidBalanceOf(_getUniqueScheduleID(_beneficiary, _id), _validator);
    }

    /**
     * @notice returns the amount of locked LNTN for some schedule
     */
    function lockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _lockedLiquidBalanceOf(_getUniqueScheduleID(_beneficiary, _id), _validator);
    }

    /**
     * @notice returns the amount of unlocked LNTN for some schedule
     */
    function unlockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _unlockedLiquidBalanceOf(_getUniqueScheduleID(_beneficiary, _id), _validator);
    }

    /**
     * @notice returns the list of validators bonded some schedule
     */
    function getBondedValidators(address _beneficiary, uint256 _id) virtual external view returns (address[] memory) {
        return _bondedValidators(_getUniqueScheduleID(_beneficiary, _id));
    }

    /**
     * @notice returns the amount of released funds in NTN for some schedule
     */
    function unlockedFunds(address _beneficiary, uint256 _id) virtual external view returns (uint256) {
        return _unlockedFunds(_getUniqueScheduleID(_beneficiary, _id));
    }

    function scheduleTotalValue(address _beneficiary, uint256 _id) external view returns (uint256) {
        return _calculateTotalValue(_getUniqueScheduleID(_beneficiary, _id));
    }

}
