// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";
import "../interfaces/IStakeProxy.sol";
import "../Liquid.sol";
import "./LiquidRewardManager.sol";

contract VestingManager is IStakeProxy, LiquidRewardManager {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release schedule
    // UNLOCKED are tokens that got released but not yet transferred
    uint256 public contractVersion = 1;
    // TODO: review: is UPDATE_FACTOR_UNIT_RECIP value high enough or too high
    uint256 public constant UPDATE_FACTOR_UNIT_RECIP = 1_000_000_000;
    // TODO: review: a way to measure requiredGasBond and requiredGasUnbond realtime?
    uint256 public requiredGasBond = 100_000;
    uint256 public requiredGasUnbond = 100_000;

    address private operator;

    struct Schedule {
        uint256 totalValue;
        uint256 withdrawnValue;
        uint256 currentNTNAmount;
        uint256 start;
        uint256 cliff;
        uint256 end;
        bool canStake;
        bool canceled;
    }

    // mapping(uint256 => mapping(address => uint256)) private liquidVestingIDs;

    // stores the unique ids of schedules assigned to a beneficiary, but beneficiary does not need to know the id
    // beneficiary will number his schedules as: 0 for first schedule, 1 for 2nd and so on
    // we can get the unique schedule id from beneficiarySchedules as follows
    // beneficiarySchedules[beneficiary][0] is the unique id of his first schedule
    // beneficiarySchedules[beneficiary][1] is the unique id of his 2nd schedule and so on
    mapping(address => uint256[]) private beneficiarySchedules;

    // list of all schedules
    Schedule[] private schedules;

    struct PendingBondingRequest {
        uint256 amount;
        uint256 epochID;
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
    // mapping(uint256 => mapping(address => uint256)) private pendingStake;
    // mapping(address => uint256) private totalPendingStake;
    // mapping(address => uint256[]) private validatorToSchedule;
    // uint256[] private pendingBondingIDs;
    uint256[] private unbondingSchedule;
    mapping(uint256 => address[]) private scheduleUnbondingValidators;

    // bondingToSchedule[bondingID] stores the unique schedule (id+1) which requested the bonding
    mapping(uint256 => uint256) private bondingToSchedule;
    mapping(uint256 => uint256[]) private scheduleToBonding;

    // unbondingToSchedule[unbondingID] stores the unique schedule (id+1) which requested the unbonding
    mapping(uint256 => uint256) private unbondingToSchedule;
    mapping(uint256 => mapping(uint256 => uint256)) private scheduleToUnbonding;
    uint256 private tailPendingUnbondingID;
    uint256 private headPendingUnbondingID;

    mapping(uint256 => address) private cancelRecipient;

    struct Ratio {
        uint256 liquidSupply;
        uint256 delegatedStake;
        int256 lastTotalValueUpdateFactor;
    }

    mapping(uint256 => mapping(address => int256)) totalValueUpdateFactor;

    // Stores ratio for each pair of (id,v) where id = unique schedule id and v = validator address.
    // There should be only one ratio for one validator, but if we keep only one object of ratio for a validator then
    // each time the ratio is updated, all schedules bonded to that validator will have to be updated.
    // To make it efficient, we store object of ratio for each pair of (id,v), if schedule id has bonded to validator v.
    // So if the ratio of (id,v) is updated then only the schedule of id needs to be updated.
    mapping(address => Ratio) private liquidRatio;

    constructor(address payable _autonity, address _operator) LiquidRewardManager(_autonity) {
        operator = _operator;
    }

    /**
     * @notice creates a new schedule
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _startBlock start block of the vesting
     * @param _cliffBlock cliff block of the vesting
     * @param _endBlock end block of the vesting
     * @param _canStake if the NTN can be staked or not
     */
    function newSchedule(
        address _beneficiary,
        uint256 _amount,
        uint256 _startBlock,
        uint256 _cliffBlock,
        uint256 _endBlock,
        bool _canStake
    ) virtual onlyOperator public {
        require(_cliffBlock >= _startBlock, "cliff must be greater to start");
        require(_endBlock > _cliffBlock, "end must be greater than cliff");

        bool _transferred = autonity.transferFrom(operator, address(this), _amount);
        require(_transferred, "amount not approved");

        uint256 _scheduleID = schedules.length;
        // uint256 _vestingID = _newVesting(_amount, _startBlock, _cliffBlock, _endBlock);
        schedules.push(Schedule(_amount, 0, _amount, _startBlock, _cliffBlock, _endBlock, _canStake, false));
        beneficiarySchedules[_beneficiary].push(_scheduleID);
    }


    /**
     * @notice used by beneficiary to transfer all unlocked NTN and LNTN of some schedule to his own address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     */
    function releaseFunds(uint256 _id) virtual external onlyActive(_id) {
        // first NTN is released
        releaseAllNTN(_id);
        // if there still remains some released amount, i.e. not enough NTN, then LNTN is released
        releaseAllLNTN(_id);
    }

    /**
     * @notice used by beneficiary to transfer all unlocked NTN of some schedule to his own address
     */
    function releaseAllNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.cliff <= block.number, "not reached cliff period yet");
        uint256 _unlocked = _unlockedFunds(_scheduleID);
        if (_unlocked > _schedule.currentNTNAmount) {
            _updateAndTransferNTN(_scheduleID, msg.sender, _schedule.currentNTNAmount);
        }
        else if (_unlocked > 0) {
            _updateAndTransferNTN(_scheduleID, msg.sender, _unlocked);
        }
    }

    /**
     * @notice used by beneficiary to transfer all unlocked LNTN of some schedule to his own address
     */
    function releaseAllLNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.cliff <= block.number, "not reached cliff period yet");
        uint256 _unlocked = _unlockedFunds(_scheduleID);
        uint256 _remaining = _unlocked;
        address[] memory _validators = _bondedValidators(_scheduleID);
        for (uint256 i = 0; i < _validators.length && _remaining > 0; i++) {
            uint256 _balance = _unlockedLiquidBalanceOf(_scheduleID, _validators[i]);
            if (_balance == 0) {
                continue;
            }
            Ratio storage _ratio = liquidRatio[_validators[i]];
            uint256 _value = _calculateLNTNValue(_ratio.liquidSupply, _ratio.delegatedStake, _balance);
            if (_remaining >= _value) {
                _remaining -= _value;
                _updateAndTransferLNTN(_scheduleID, msg.sender, _balance, _validators[i]);
            }
            else {
                uint256 _liquid = _getLiquidFromNTN(_ratio.delegatedStake, _ratio.liquidSupply, _remaining);
                require(_liquid <= _balance, "conversion not working");
                _remaining = 0;
                _updateAndTransferLNTN(_scheduleID, msg.sender, _liquid, _validators[i]);
            }
        }
        _schedule.withdrawnValue += _unlocked - _remaining;
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice used by beneficiary to transfer some amount of unlocked NTN of some schedule to his own address
     * @param _amount amount to transfer
     */
    function releaseNTN(uint256 _id, uint256 _amount) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.cliff <= block.number, "not reached cliff period yet");
        uint256 _unlocked = _unlockedFunds(_scheduleID);
        if (_amount < _unlocked) {
            _updateAndTransferNTN(_scheduleID, msg.sender, _amount);
        }
        else {
            _updateAndTransferNTN(_scheduleID, msg.sender, _unlocked);
        }
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice used by beneficiary to transfer some amount of unlocked LNTN of some schedule to his own address
     * @param _validator address of the validator
     * @param _amount amount of LNTN to transfer
     */
    function releaseLNTN(uint256 _id, address _validator, uint256 _amount) virtual external onlyActive(_id) {
        require(_amount > 0, "require positive amount to transfer");
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.cliff <= block.number, "not reached cliff period yet");
        require(_unlockedLiquidBalanceOf(_scheduleID, _validator) >= _amount, "not enough unlocked LNTN");
        uint256 _unlocked = _unlockedFunds(_scheduleID);
        Ratio storage _ratio = liquidRatio[_validator];
        uint256 _value = _calculateLNTNValue(_ratio.liquidSupply, _ratio.delegatedStake, _amount);
        if (_value > _unlocked) {
            // not enough released, reduce liquid accordingly
            _value = _unlocked;
            _amount = _getLiquidFromNTN(_ratio.delegatedStake, _ratio.liquidSupply, _unlocked);
            require(_amount <= _unlockedLiquidBalanceOf(_scheduleID, _validator), "conversion not working");
        }
        _schedule.withdrawnValue += _value;
        _updateAndTransferLNTN(_scheduleID, msg.sender, _amount, _validator);
        
    }

    /**
     * @notice force release of all funds, NTN and LNTN, of some schedule and return them to the _recipient account
     * effectively cancelling a vesting schedule. only operator is able to call the function
     * rewards (AUT) which have been entitled to a schedule due to bonding are not returned to _recipient
     * @param _beneficiary beneficiary address whose schedule will be canceled
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     * @param _recipient to whome the all funds will be transferred
     */
    function cancelSchedule(address _beneficiary, uint256 _id, address _recipient) virtual external onlyOperator {
        // we don't care about updating the schedule as the schedule is canceled
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        Schedule storage _item = schedules[_scheduleID];
        require(_item.canceled == false, "schedule already canceled");
        address[] memory _validators = _bondedValidators(_scheduleID);
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _amount = _unlockedLiquidBalanceOf(_scheduleID, _validators[i]);
            if (_amount > 0) {
                _updateAndTransferLNTN(_scheduleID, _recipient, _amount, _validators[i]);
            }
        }
        _item.canceled = true;
        cancelRecipient[_scheduleID] = _recipient;
        _updateAndTransferNTN(_scheduleID, _recipient, _item.currentNTNAmount);
    }

    function setRequiredGasBond(uint256 _gas) external onlyOperator {
        requiredGasBond = _gas;
    }

    function setRequiredGasUnbond(uint256 _gas) external onlyOperator {
        requiredGasUnbond = _gas;
    }

    /**
     * @notice Used by beneficiary to bond some NTN of some schedule. only schedules with canStake = true can be staked. Use canStake function
     * to check if can be staked or not. Bonded from the vesting manager contract. All bondings are delegated, as vesting manager cannot own a validator
     * @param _id id of the schedule numbered from 0 to (n-1) where n = total schedules entitled to the beneficiary (including canceled ones)
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(uint256 _id, address _validator, uint256 _amount) virtual public payable onlyActive(_id) returns (uint256) {
        require(msg.value >= requiredGasBond * autonity.stakingGasPrice(), "not enough gas given for notification on bonding");
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.canStake, "not able to stake");
        require(_schedule.currentNTNAmount >= _amount, "not enough tokens");

        uint256 _bondingID = autonity.bond{value: msg.value}(_validator, _amount);
        _schedule.currentNTNAmount -= _amount;
        // offset by 1 to handle empty value
        scheduleToBonding[_scheduleID].push(_bondingID);
        bondingToSchedule[_bondingID] = _scheduleID+1;
        pendingBondingRequest[_bondingID] = PendingBondingRequest(_amount, _epochID());
        _addValidator(_scheduleID, _validator);
        // if (pendingStake[_scheduleID][_validator] == 0) {
        //     validatorToSchedule[_validator].push(_scheduleID);
        // }
        // pendingStake[_scheduleID][_validator] += _amount;
        // totalPendingStake[_validator] += _amount;
        // pendingBondingIDs.push(_bondingID);
        return _bondingID;
    }

    /**
     * @notice Used by beneficiary to unbond some LNTN of some schedule.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public payable onlyActive(_id) returns (uint256) {
        require(msg.value >= requiredGasUnbond * autonity.stakingGasPrice(), "not enough gas given for notification on unbonding");
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _cleanup(_scheduleID);
        _lock(_scheduleID, _validator, _amount);
        uint256 _unbondingID = autonity.unbond{value: msg.value}(_validator, _amount);
        // offset by 1 to handle empty value
        unbondingToSchedule[_unbondingID] = _scheduleID+1;
        pendingUnbondingRequest[_unbondingID] = PendingUnbondingRequest(_amount, _epochID(), _validator, false, false);
        scheduleToUnbonding[_scheduleID][headPendingUnbondingID] = _unbondingID;
        headPendingUnbondingID++;
        return _unbondingID;
    }

    /**
     * @notice used by beneficiary to claim all rewards (AUT) which is entitled due to bonding
     */
    function claimRewards() virtual external {
        uint256[] storage _scheduleIDs = beneficiarySchedules[msg.sender];
        uint256 _totalFees = 0;
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            _cleanup(_scheduleIDs[i]);
            _totalFees += _rewards(_scheduleIDs[i]);
        }
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool _sent, ) = msg.sender.call{value: _totalFees}("");
        require(_sent, "Failed to send AUT");
    }

    /**
     * @notice can be used to send AUT to the contract
     */
    function receiveAut() external payable {
        // do nothing
    }

    function rewardsDistributed(
        address[] memory _validators,
        uint256[] memory _delegatedStake,
        uint256[] memory _liquidSupply
    ) external onlyAutonity {
        _updateValidatorRewardAndRatio(_validators, _delegatedStake, _liquidSupply);
    }

    /**
     * @dev implements IStakeProxy.bondingApplied(), a callback function for autonity when bonding is applied
     * @param _bondingID bonding id from Autonity when bonding is requested
     * @param _liquid amount of LNTN after bonding applied successfully
     * @param _selfDelegation true if self bonded, false for delegated bonding
     * @param _rejected true if bonding request was rejected
     */
    function bondingApplied(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) external onlyAutonity {
        _applyBonding(_bondingID, _validator, _liquid, _selfDelegation, _rejected);
    }

    // /**@dev implements IStakeProxy.unbondingApplied(). callback function for autonity when unbonding is applied
    //  * @param _unbondingID unbonding id from Autonity when unbonding is requested
    //  */
    function unbondingApplied(uint256 _unbondingID, address _validator, bool _rejected) external onlyAutonity {
        _applyUnbonding(_unbondingID, _validator, _rejected);
    }

    // /**@dev implements IStakeProxy.unbondingReleased(). callback function for autonity when unbonding is released
    //  * @param _unbondingID unbonding id from Autonity when unbonding is requested
    //  * @param _amount amount of NTN released
    //  */
    function unbondingReleased(uint256 _unbondingID, address _validator, uint256 _amount, bool _rejected) external onlyAutonity {
        _releaseUnbonding(_unbondingID, _validator, _amount, _rejected);
    }

    /**
     * @dev returns equivalent amount of NTN using the ratio.
     * takes floor to match the calculation with Autonity
     * @param _liquidSupply amount of LNTN in the ratio
     * @param _delegatedStake amount of NTN in the ratio
     * @param _amount amount of LNTN to be converted
     */
    function _calculateLNTNValue(uint256 _liquidSupply, uint256 _delegatedStake, uint256 _amount) private pure returns (uint256) {
        // return math.floor(_amount * _delegatedStake / _liquidSupply)
        return _amount * _delegatedStake / _liquidSupply;
    }

    /**
     * @dev returns equivalent amount of NTN using the ratio.
     * takes ceil to match the conversion from liquid to NTN in _calculateLNTNValue
     * @param _delegatedStake amount of NTN in the ratio
     * @param _liquidSupply amount of LNTN in the ratio
     * @param _amount amount of NTN to be converted
     */
    function _getLiquidFromNTN(uint256 _delegatedStake, uint256 _liquidSupply, uint256 _amount) private pure returns (uint256) {
        // return math.ceil(_amount * _liquidSupply / _delegatedStake)
        return (_amount * _liquidSupply + _delegatedStake - 1) / _delegatedStake;
    }

    function _updateSchedule(uint256 _scheduleID) private {
        address[] memory _validators = _bondedValidators(_scheduleID);
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _balance = _liquidBalanceOf(_scheduleID, _validators[i]);
            if (_balance == 0) {
                continue;
            }
            // view only call
            Autonity.Validator memory _validatorInfo = autonity.getValidator(_validators[i]);
            _updateRatio(_validators[i], _validatorInfo.bondedStake - _validatorInfo.selfBondedStake, _validatorInfo.liquidSupply);
            _updateScheduleTotalValue(_scheduleID, _validators[i]);
        }
    }

    function _updateRatio(address _validator, uint256 _delegatedStake, uint256 _liquidSupply) private {
        Ratio storage _ratio = liquidRatio[_validator];
        uint256 _oldLiquidSupply = _ratio.liquidSupply;
        uint256 _oldDelegatedStake = _ratio.delegatedStake;
        // the first time ratio is set, _oldLiquidSupply = 0 and lastTotalValueUpdateFactor has no need to update
        if (_oldLiquidSupply > 0 && _delegatedStake * _oldLiquidSupply != _oldDelegatedStake * _liquidSupply) {
            _ratio.lastTotalValueUpdateFactor += int(UPDATE_FACTOR_UNIT_RECIP * _delegatedStake / _liquidSupply)
                                                - int(UPDATE_FACTOR_UNIT_RECIP * _oldDelegatedStake / _oldLiquidSupply);
        }
        _ratio.delegatedStake = _delegatedStake;
        _ratio.liquidSupply = _liquidSupply;
    }

    function _updateScheduleTotalValue(uint256 _scheduleID, address _validator) private {
        Ratio storage _ratio = liquidRatio[_validator];
        Schedule storage _schedule = schedules[_scheduleID];
        int256 _lastTotalValueUpdateFactor = _ratio.lastTotalValueUpdateFactor;
        _schedule.totalValue +=_liquidBalanceOf(_scheduleID, _validator)
                                * uint(_lastTotalValueUpdateFactor - totalValueUpdateFactor[_scheduleID][_validator]);
        totalValueUpdateFactor[_scheduleID][_validator] = _lastTotalValueUpdateFactor;
    }

    function _unlockedFunds(uint256 _scheduleID) private returns (uint256) {
        Schedule storage _schedule = schedules[_scheduleID];
        if (block.number < _schedule.cliff) return 0;
        // update the ratio of LNTN:NTN for each bonding of this schedule
        // also update the initial amount of the schedule according to the value of LNTN
        _updateSchedule(_scheduleID);
        uint256 _unlocked = 0;
        if (block.number >= _schedule.end) {
            _unlocked = _schedule.totalValue;
        }
        else {
            _unlocked = _schedule.totalValue * (block.number - _schedule.start) / (_schedule.end - _schedule.start);
        }

        if (_unlocked > _schedule.withdrawnValue) {
            return _unlocked - _schedule.withdrawnValue;
        }
        return 0;
    }

    /**
     * @dev returns a unique id for each schedule
     * @param _beneficiary address of the schedule holder
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     */
    function _getUniqueScheduleID(address _beneficiary, uint256 _id) private view returns (uint256) {
        require(beneficiarySchedules[_beneficiary].length > _id, "invalid schedule id");
        return beneficiarySchedules[_beneficiary][_id];
    }

    function _updateAndTransferNTN(uint256 _scheduleID, address _to, uint256 _amount) private {
        schedules[_scheduleID].currentNTNAmount -= _amount;
        schedules[_scheduleID].withdrawnValue += _amount;
        _transferNTN(_to, _amount);
    }

    function _updateAndTransferLNTN(uint256 _scheduleID, address _to, uint256 _amount, address _validator) private {
        _decreaseLiquid(_scheduleID, _validator, _amount);
        _transferLNTN(_to, _amount, _validator);
    }

    function _transferNTN(address _to, uint256 _amount) private {
        bool _sent = autonity.transfer(_to, _amount);
        require(_sent, "NTN not transferred");
    }

    function _transferLNTN(address _to, uint256 _amount, address _validator) private {
        bool _sent = _liquidContract(_validator).transfer(_to, _amount);
        require(_sent, "LNTN transfer failed");
    }

    function _updateValidatorRewardAndRatio(
        address[] memory _validators,
        uint256[] memory _delegatedStake,
        uint256[] memory _liquidSupply
    ) internal {
        for (uint256 i = 0; i < _validators.length; i++) {
            _updateUnclaimedReward(_validators[i]);
            if (_liquidSupply[i] > 0) {
                _updateRatio(_validators[i], _delegatedStake[i], _liquidSupply[i]);
            }
        }
    }

    function _applyBonding(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) internal {
        require(_selfDelegation == false, "bonding should be delegated");
        uint256 _scheduleID = bondingToSchedule[_bondingID] - 1;
        delete bondingToSchedule[_bondingID];

        Schedule storage _schedule = schedules[_scheduleID];
        if (_schedule.canceled) {
            if (_rejected) {
                _transferNTN(cancelRecipient[_scheduleID], pendingBondingRequest[_bondingID].amount);
            }
            else {
                _transferLNTN(cancelRecipient[_scheduleID], _liquid, _validator);
            }
            delete pendingBondingRequest[_bondingID];
            return;
        }

        if (_rejected) {
            _schedule.currentNTNAmount += pendingBondingRequest[_bondingID].amount;
        }
        else {
            // total value needs to be updated before increasing liquid
            _updateScheduleTotalValue(_scheduleID, _validator);
            _increaseLiquid(_scheduleID, _validator, _liquid);
        }
        delete pendingBondingRequest[_bondingID];
    }

    function _applyUnbonding(uint256 _unbondingID, address _validator, bool _rejected) internal {
        uint256 _scheduleID = unbondingToSchedule[_unbondingID] - 1;
        Schedule storage _schedule = schedules[_scheduleID];
        PendingUnbondingRequest storage _unbondingRequest = pendingUnbondingRequest[_unbondingID];
        uint256 _liquid = _unbondingRequest.liquidAmount;
        _unlock(_scheduleID, _validator, _liquid);
        if (_rejected) {
            _unbondingRequest.rejected = true;
            if (_schedule.canceled) {
                _updateAndTransferLNTN(_scheduleID, cancelRecipient[_scheduleID], _liquid, _validator);
            }
            return;
        }
        _unbondingRequest.applied = true;
        // total value needs to be updated before decreasing liquid
        _updateScheduleTotalValue(_scheduleID, _validator);
        Ratio storage _ratio = liquidRatio[_validator];
        _schedule.totalValue -= _calculateLNTNValue(_ratio.liquidSupply, _ratio.delegatedStake, _liquid);
        _decreaseLiquid(_scheduleID, _validator, _liquid);
    }

    function _releaseUnbonding(uint256 _unbondingID, address _validator, uint256 _amount, bool _rejected) internal {
        uint256 _scheduleID = unbondingToSchedule[_unbondingID] - 1;
        delete unbondingToSchedule[_unbondingID];

        Schedule storage _schedule = schedules[_scheduleID];
        if (_rejected) {
            // if released of unbonding is rejected, then it is assumed that the applying of unbonding was also rejected
            // or reverted at Autonity because Autonity could not notify us (_applyUnbonding reverted)
            // if applying of unbonding was successful, then releasing of unbonding cannot be rejected
            PendingUnbondingRequest storage _unbondingRequest = pendingUnbondingRequest[_unbondingID];
            if (!_unbondingRequest.rejected) {
                require(_unbondingRequest.applied == false, "unbonding applied successfully but release is rejected");
                uint256 _liquid = _unbondingRequest.liquidAmount;
                // _applyUnbonding failed for some reason
                _unlock(_scheduleID, _validator, _liquid);
                _unbondingRequest.rejected = true;
                if (_schedule.canceled) {
                    _updateAndTransferLNTN(
                        _scheduleID, cancelRecipient[_scheduleID], _liquid, _validator
                    );
                }
            }
            delete pendingUnbondingRequest[_unbondingID];
            return;
        }
        delete pendingUnbondingRequest[_unbondingID];
        if (_schedule.canceled) {
            _transferNTN(cancelRecipient[_scheduleID], _amount);
        }
        else {
            _schedule.currentNTNAmount += _amount;
            _schedule.totalValue += _amount;
        }
    }

    /**
     * @dev The schedule needs to be cleaned before bonding, unbonding or claiming rewards.
     * _cleanup removes any unnecessary validator from the list, removes pending bonding or unbonding request
     * that has been rejected or reverted but vesting manager could not be notified. If the clean up is not
     * done, then liquid balance could be incorrect due to not handling the bonding or unbonding request.
     * @param _scheduleID unique global schedule id
     */
    function _cleanup(uint256 _scheduleID) private {
        _revertPendingBondingRequest(_scheduleID);
        _revertPendingUnbondingRequest(_scheduleID);
        _clearValidators(_scheduleID);
    }

    function _revertPendingBondingRequest(uint256 _scheduleID) private {
        // We can keep the request in array because ALL the bonding requests are processed at epoch end
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
            // if the request is not processed successfully, then we need to revert it
            // in case the request is successfully processed, we delete it from pendingBondingRequest after processing
            if (_oldBondingRequest.amount > 0) {
                _totalAmount += _oldBondingRequest.amount;
                delete pendingBondingRequest[_oldBondingID];
                delete bondingToSchedule[_oldBondingID];
            }
        }
        delete scheduleToBonding[_scheduleID];
        if (_totalAmount == 0) {
            return;
        }
        Schedule storage _schedule = schedules[_scheduleID];
        if (_schedule.canceled) {
            _transferNTN(cancelRecipient[_scheduleID], _totalAmount);
        }
        else {
            _schedule.currentNTNAmount += _totalAmount;
        }
    }

    function _revertPendingUnbondingRequest(uint256 _scheduleID) private {
        uint256 _processedID = tailPendingUnbondingID;
        uint256 _unbondingID;
        PendingUnbondingRequest storage _unbondingRequest;
        Schedule storage _schedule = schedules[_scheduleID];
        mapping(uint256 => uint256) storage _unbondingIDs = scheduleToUnbonding[_scheduleID];
        uint256 _lastID = headPendingUnbondingID;
        for(; _processedID < _lastID; _processedID++) {
            _unbondingID = _unbondingIDs[_processedID];
            Autonity.UnbondingReleaseState _releaseState = autonity.getUnbondingReleaseState(_unbondingID);
            if (_releaseState == Autonity.UnbondingReleaseState.notReleased) {
                // all the rest have not been released yet
                break;
            }
            delete _unbondingIDs[_processedID];
            
            if (_releaseState == Autonity.UnbondingReleaseState.released) {
                // unbonding was released successfully
                continue;
            }

            delete unbondingToSchedule[_unbondingID];
            _unbondingRequest = pendingUnbondingRequest[_unbondingID];
            address _validator = _unbondingRequest.validator;

            if (_releaseState == Autonity.UnbondingReleaseState.reverted) {
                // it means the unbonding was released, but later reverted due to failing to notify us
                // that means unbonding was applied successfully
                require(_unbondingRequest.applied, "unbonding was released but not applied succesfully");
                _increaseLiquid(_scheduleID, _validator, autonity.getRevertingAmount(_unbondingID));
            }

            if (_releaseState == Autonity.UnbondingReleaseState.rejected) {
                require(_unbondingRequest.applied == false, "unbonding was applied successfully but release was rejected");
                // _unbondingRequest.rejected = true means we already rejected it
                if (_unbondingRequest.rejected == false) {
                    _unlock(_scheduleID, _validator, _unbondingRequest.liquidAmount);
                }
            }
            
            delete pendingUnbondingRequest[_unbondingID];
            if (_schedule.canceled) {
                _updateAndTransferLNTN(
                    _scheduleID, cancelRecipient[_scheduleID], _liquidBalanceOf(_scheduleID, _validator), _validator
                );
            }
        }
        tailPendingUnbondingID = _processedID;
    }

    function _epochID() internal view returns (uint256) {
        return autonity.epochID();
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice returns a schedule entitled to _beneficiary
     * @param _beneficiary beneficiary address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     */
    function getSchedule(address _beneficiary, uint256 _id) virtual external view returns (Schedule memory) {
        return schedules[_getUniqueScheduleID(_beneficiary, _id)];
    }

    /**
     * @notice returns if beneficiary can stake from his schedule
     * @param _beneficiary beneficiary address
     */
    function canStake(address _beneficiary, uint256 _id) virtual external view returns (bool) {
        return schedules[_getUniqueScheduleID(_beneficiary, _id)].canStake;
    }

    /**
     * @notice returns the number of schudeled entitled to some beneficiary
     * @param _beneficiary address of the beneficiary
     */
    function totalSchedules(address _beneficiary) virtual external view returns (uint256) {
        return beneficiarySchedules[_beneficiary].length;
    }

    /**
     * @notice returns the list of current schedules assigned to a beneficiary
     * @param _beneficiary address of the beneficiary
     */
    function getSchedules(address _beneficiary) virtual external view returns (Schedule[] memory) {
        uint256[] storage _scheduleIDs = beneficiarySchedules[_beneficiary];
        Schedule[] memory _res = new Schedule[](_scheduleIDs.length);
        for (uint256 i = 0; i < _res.length; i++) {
            _res[i] = schedules[_scheduleIDs[i]];
        }
        return _res;
    }

    /**
     * @notice returns the amount of all unclaimed rewards (AUT) due to all the bonding from schedules entitled to beneficiary
     * @param _beneficiary beneficiary address
     */
    function unclaimedRewards(address _beneficiary) virtual external returns (uint256) {
        uint256 _totalFee = 0;
        uint256[] storage _scheduleIDs = beneficiarySchedules[_beneficiary];
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            _cleanup(_scheduleIDs[i]);
            _totalFee += _unclaimedRewards(_scheduleIDs[i]);
        }
        return _totalFee;
    }

    /**
     * @notice returns the amount of LNTN for some schedule
     * @param _beneficiary beneficiary address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
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
    function unlockedFunds(address _beneficiary, uint256 _id) virtual external returns (uint256) {
        return _unlockedFunds(_getUniqueScheduleID(_beneficiary, _id));
    }

    /*
    ============================================================

        Modifiers

    ============================================================
     */

    /**
     * @dev Modifier that checks if the caller is the governance operator account.
     */
    modifier onlyOperator {
        require(operator == msg.sender, "caller is not the operator");
        _;
    }

    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to Autonity contract");
        _;
    }

    modifier onlyActive(uint256 _id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        require(schedules[_scheduleID].canceled == false, "schedule canceled");
        _;
    }


}
