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
    uint256 public constant UPDATE_FACTOR_UNIT_RECIP = 1_000_000_000;

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
    mapping(uint256 => uint256[]) private scheduleToUnbonding;

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

    /** @notice creates a new schedule
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


    /** @notice used by beneficiary to transfer all unlocked NTN and LNTN of some schedule to his own address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     */
    function releaseFunds(uint256 _id) virtual external onlyActive(_id) {
        // first NTN is released
        releaseAllNTN(_id);
        // if there still remains some released amount, i.e. not enough NTN, then LNTN is released
        releaseAllLNTN(_id);
    }

    /** @notice used by beneficiary to transfer all unlocked NTN of some schedule to his own address */
    function releaseAllNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
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

    /** @notice used by beneficiary to transfer all unlocked LNTN of some schedule to his own address */
    function releaseAllLNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
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
    /** @notice used by beneficiary to transfer some amount of unlocked NTN of some schedule to his own address
     * @param _amount amount to transfer
     */
    function releaseNTN(uint256 _id, uint256 _amount) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
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
    /** @notice used by beneficiary to transfer some amount of unlocked LNTN of some schedule to his own address
     * @param _validator address of the validator
     * @param _amount amount of LNTN to transfer
     */
    function releaseLNTN(uint256 _id, address _validator, uint256 _amount) virtual external onlyActive(_id) {
        require(_amount > 0, "require positive amount to transfer");
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
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

    /** @notice force release of all funds, NTN and LNTN, of some schedule and return them to the _recipient account
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

    function revertPendingBondingRequest(address _beneficiary, uint256 _id) public {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
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
            if (_oldBondingRequest.amount > 0) {
                _totalAmount += _oldBondingRequest.amount;
                delete pendingBondingRequest[_oldBondingID];
                delete bondingToSchedule[_oldBondingID];
            }
        }
        Schedule storage _schedule = schedules[_scheduleID];
        if (_schedule.canceled) {
            _transferNTN(cancelRecipient[_scheduleID], _totalAmount);
        }
        else {
            _schedule.currentNTNAmount += _totalAmount;
        }
        delete scheduleToBonding[_scheduleID];
    }

    function revertPendingUnbondingRequest(address _beneficiary, uint256 _id) public {
        // TODO: the following is incorrect. need to include unbondingPeriod from Autonity. But a better solution may be possible
        // without including unbondingPeriod as it can changed by operator
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        uint256[] storage _unbondingIDs = scheduleToUnbonding[_scheduleID];
        uint256 _length = _unbondingIDs.length;
        if (_length > 0) {
            uint256 _oldUnbondingID = _unbondingIDs[0];
            PendingUnbondingRequest storage _oldUnbondingRequest = pendingUnbondingRequest[_oldUnbondingID];
            // will revert request from some previous epoch, request from current epoch will not be reverted
            if (_oldUnbondingRequest.liquidAmount > 0 && _oldUnbondingRequest.epochID < _epochID()) {
                Schedule storage _schedule = schedules[_scheduleID];
                for (uint256 i = 0; i < _length; i++) {
                    _oldUnbondingID = _unbondingIDs[i];
                    _oldUnbondingRequest = pendingUnbondingRequest[_oldUnbondingID];
                    address _validator = _oldUnbondingRequest.validator;
                    if (_lockedLiquidBalanceOf(_scheduleID, _validator) >= _oldUnbondingRequest.liquidAmount) {
                        // unbonding was not applied, just unlock the locked liquid
                        _unlock(_scheduleID, _validator, _oldUnbondingRequest.liquidAmount);
                    }
                    else {
                        // unbonding was applied, but released was reverted. mint necessary amount
                        uint256 _liquidToMint = _oldUnbondingRequest.liquidAmount - _lockedLiquidBalanceOf(_scheduleID, _validator);
                        _unlockAll(_scheduleID, _validator);
                        _increaseLiquid(_scheduleID, _validator, _liquidToMint);
                    }
                    
                    if (_schedule.canceled) {
                        _updateAndTransferLNTN(_scheduleID, cancelRecipient[_scheduleID], _oldUnbondingRequest.liquidAmount, _validator);
                    }
                    delete pendingUnbondingRequest[_oldUnbondingID];
                    delete unbondingToSchedule[_oldUnbondingID];
                }
                
                delete scheduleToUnbonding[_scheduleID];
            }
        }
    }

    /** @notice Used by beneficiary to bond some NTN of some schedule. only schedules with canStake = true can be staked. Use canStake function
     * to check if can be staked or not. Bonded from the vesting manager contract. All bondings are delegated, as vesting manager cannot own a validator
     * @param _id id of the schedule numbered from 0 to (n-1) where n = total schedules entitled to the beneficiary (including canceled ones)
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(uint256 _id, address _validator, uint256 _amount) virtual public payable onlyActive(_id) returns (uint256) {
        // in case bonding request for some previous epoch was reverted due to failing to notify us
        revertPendingBondingRequest(msg.sender, _id);

        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
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

    /** @notice Used by beneficiary to unbond some LNTN of some schedule.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public payable onlyActive(_id) returns (uint256) {
        // in case unbonding request for some previous epoch was reverted due to failing to notify us
        revertPendingUnbondingRequest(msg.sender, _id);

        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _lock(_scheduleID, _validator, _amount);
        uint256 _unbondingID = autonity.unbond{value: msg.value}(_validator, _amount);
        // offset by 1 to handle empty value
        unbondingToSchedule[_unbondingID] = _scheduleID+1;
        scheduleToUnbonding[_scheduleID].push(_unbondingID);
        pendingUnbondingRequest[_unbondingID] = PendingUnbondingRequest(_amount, _epochID(), _validator, false);
        _clearValidators(_scheduleID);
        return _unbondingID;
    }

    /** @notice used by beneficiary to claim all rewards (AUT) which is entitled due to bonding */
    function claimRewards() virtual external {
        uint256[] storage _scheduleIDs = beneficiarySchedules[msg.sender];
        uint256 _totalFees = 0;
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            _totalFees += _rewards(_scheduleIDs[i]);
        }
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool _sent, ) = msg.sender.call{value: _totalFees}("");
        require(_sent, "Failed to send AUT");
    }

    /** @notice callback function for autonity to finalize epoch */
    function finalize(
        BondingApplied[] memory _bonding,
        uint256[] memory _rejectedBonding,
        UnbondingReleased[] memory _releasedUnbonding
    ) external onlyAutonity {
        // _finalize(_bonding, _rejectedBonding, _releasedUnbonding);
    }

    /** @notice can be used to send AUT to the contract */
    function receiveAut() external payable {
        // do nothing
    }

    function distributeRewards(address[] memory _validators) external {}

    // /** @dev implements IStakeProxy.bondingApplied(). a callback function for autonity when bonding is applied
    //  * @param _bondingID bonding id from Autonity when bonding is requested
    //  * @param _liquid amount of LNTN for bonding
    //  * @param _selfDelegation true if self bonded, false for delegated bonding
    //  * @param _rejected true if bonding request was rejected
    //  */
    function bondingApplied(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) external onlyAutonity {
        _applyBonding(_bondingID, _validator, _liquid, _selfDelegation, _rejected);
    }

    // /** @dev implements IStakeProxy.unbondingApplied(). callback function for autonity when unbonding is applied
    //  * @param _unbondingID unbonding id from Autonity when unbonding is requested
    //  */
    function unbondingApplied(uint256 _unbondingID, address _validator, bool _rejected) external onlyAutonity {
        _applyUnbonding(_unbondingID, _validator, _rejected);
    }

    // /** @dev implements IStakeProxy.unbondingReleased(). callback function for autonity when unbonding is released
    //  * @param _unbondingID unbonding id from Autonity when unbonding is requested
    //  * @param _amount amount of NTN released
    //  */
    function unbondingReleased(uint256 _unbondingID, address _validator, uint256 _amount, bool _rejected) external onlyAutonity {
        _releaseUnbonding(_unbondingID, _validator, _amount, _rejected);
    }

    /** @dev returns equivalent amount of NTN using the ratio.
     * takes floor to match the calculation with Autonity
     * @param _liquidSupply amount of LNTN in the ratio
     * @param _delegatedStake amount of NTN in the ratio
     * @param _amount amount of LNTN to be converted
     */
    function _calculateLNTNValue(uint256 _liquidSupply, uint256 _delegatedStake, uint256 _amount) private pure returns (uint256) {
        // return math.floor(_amount * _delegatedStake / _liquidSupply)
        return _amount * _delegatedStake / _liquidSupply;
    }

    /** @dev returns equivalent amount of NTN using the ratio.
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
        _clearValidators(_scheduleID);
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
        // the first time ratio is set, _oldLiquidSupply = 0 and lastTotalValueUpdateFactor has no need to update
        if (_oldLiquidSupply > 0) {
            _ratio.lastTotalValueUpdateFactor += int(UPDATE_FACTOR_UNIT_RECIP * _delegatedStake / _liquidSupply)
                                                - int(UPDATE_FACTOR_UNIT_RECIP * _ratio.delegatedStake / _oldLiquidSupply);
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

    /** @dev returns a unique id for each schedule
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
        if (_rejected) {
            _unlock(_scheduleID, _validator, _liquid);
            if (_schedule.canceled) {
                _updateAndTransferLNTN(_scheduleID, cancelRecipient[_scheduleID], _liquid, _validator);
            }
            _unbondingRequest.rejected = true;
            return;
        }
        _updateScheduleTotalValue(_scheduleID, _validator);
        Ratio storage _ratio = liquidRatio[_validator];
        _schedule.totalValue -= _calculateLNTNValue(_ratio.liquidSupply, _ratio.delegatedStake, _liquid);
        _unlock(_scheduleID, _validator, _liquid);
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
                uint256 _liquid = _unbondingRequest.liquidAmount;
                // _applyUnbonding failed for some reason
                _unlock(_scheduleID, _validator, _liquid);
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

    // /** @dev Callback function for autonity at finalize function when epoch is ended. Internal for testing purpose
    //  * This function handles the bonding and unbonding related mechanism.
    //  * Follow _stakingOperations() function in Autonity.sol to apply the operations correctly.
    //  * Note: Autonity does not know about vesting manager and will use the callback for any contract
    //  * that has bonded or unbonded in autonity contract.
    //  * @param _appliedBonding list of bonding requests that are applied
    //  * @param _rejectedBonding list of bonding id that are rejected
    //  * @param _releasedUnbonding list of unbonding requests that are applied
    //  */
    // function _finalize(
    //     BondingApplied[] memory _appliedBonding, // 20 * 100
    //     uint256[] memory _rejectedBonding,
    //     UnbondingReleased[] memory _releasedUnbonding
    // ) internal {
    //     // all bonding requests are processed (applied or rejected) at finalize
    //     // handle rejected bonding requests first
    //     for (uint256 i = 0; i < _rejectedBonding.length; i++) {
    //         uint256 _scheduleID = bondingToSchedule[_rejectedBonding[i]];
    //         delete bondingToSchedule[_rejectedBonding[i]];
    //         PendingBondingRequest storage _request = pendingBondingRequest[_rejectedBonding[i]];
    //         address _validator = _request.validator;
    //         delete pendingStake[_scheduleID][_validator];
    //         delete totalPendingStake[_validator];
    //         delete validatorToSchedule[_validator];
    //         schedules[_scheduleID].currentNTNAmount += _request.amount;
    //         delete pendingBondingRequest[_rejectedBonding[i]];
    //     }
        
    //     // we assume all bondings are delegation, so liquid amount is positive
    //     uint256 _length;
    //     for (uint256 i = 0; i < _appliedBonding.length; i++) {
    //         _updateUnclaimedReward(_appliedBonding[i].validator);
    //         uint256 _totalDelegatedStake = totalPendingStake[_appliedBonding[i].validator];
    //         delete totalPendingStake[_appliedBonding[i].validator];
    //         _updateRatio(_appliedBonding[i].validator, _totalDelegatedStake, _appliedBonding[i].liquidAmount);
    //         uint256[] storage _scheduleIDs = validatorToSchedule[_appliedBonding[i].validator];
    //         _length = _scheduleIDs.length;
    //         for (uint256 j = 0; j < _length; j++) {
    //             uint256 _scheduleID = _scheduleIDs[j];
    //             _updateScheduleTotalValue(_scheduleID, _appliedBonding[i].validator);
    //             uint256 _delegatedStake = pendingStake[_scheduleID][_appliedBonding[i].validator];
    //             delete pendingStake[_scheduleID][_appliedBonding[i].validator];
    //             uint256 _liquidAmount = _getLiquidFromNTN(_totalDelegatedStake, _appliedBonding[i].liquidAmount, _delegatedStake);
    //             _totalDelegatedStake -= _delegatedStake;
    //             _appliedBonding[i].liquidAmount -= _liquidAmount;
    //             _increaseLiquid(_scheduleID, _appliedBonding[i].validator, _liquidAmount);
    //         }
    //         delete validatorToSchedule[_appliedBonding[i].validator];
    //     }

    //     _length = pendingBondingIDs.length;
    //     for (uint256 i = 0; i < _length; i++) {
    //         delete bondingToSchedule[pendingBondingIDs[i]];
    //         delete pendingBondingRequest[pendingBondingIDs[i]];
    //     }
    //     delete pendingBondingIDs;

    //     // all unbonding requests are applied at finalize
    //     _length = unbondingSchedule.length;
    //     for (uint256 i = 0; i < _length; i++) {
    //         uint256 _scheduleID = unbondingSchedule[i];
    //         uint256 _totalValueDecrease;
    //         address[] storage _validators = scheduleUnbondingValidators[_scheduleID];
    //         uint256 _validatorCount = _validators.length;
    //         for (uint256 j = 0; j < _validatorCount; j++) {
    //             address _validator = _validators[j];
    //             _updateScheduleTotalValue(_scheduleID, _validator);
    //             Ratio storage _ratio = liquidRatio[_validator];
    //             _totalValueDecrease += _calculateLNTNValue(
    //                 _ratio.liquidSupply, _ratio.delegatedStake, _lockedLiquidBalanceOf(_scheduleID, _validator)
    //             );
    //             // We lock LNTN only if unbonding is requested.
    //             // Because all unbondings are applied, we unlock all and burn them
    //             _unlockAllAndBurn(_scheduleID, _validator);
    //         }
    //         delete scheduleUnbondingValidators[_scheduleID];
    //         Schedule storage _schedule = schedules[_scheduleID];
    //         _schedule.totalValue -= _totalValueDecrease;
    //     }
    //     delete unbondingSchedule;

    //     // handle released unbonding
    //     for (uint256 i = 0; i < _releasedUnbonding.length; i++) {
    //         uint256 _scheduleID = unbondingToSchedule[_releasedUnbonding[i].unbondingID];
    //         delete unbondingToSchedule[_releasedUnbonding[i].unbondingID];
    //         Schedule storage _schedule = schedules[_scheduleID];
    //         _schedule.totalValue += _releasedUnbonding[i].releasedAmount;
    //         _schedule.currentNTNAmount += _releasedUnbonding[i].releasedAmount;
    //     }
    // }

    function _epochID() internal view returns (uint256) {
        return autonity.epochID();
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /** @notice returns a schedule entitled to _beneficiary
     * @param _beneficiary beneficiary address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     */
    function getSchedule(address _beneficiary, uint256 _id) virtual external view returns (Schedule memory) {
        return schedules[_getUniqueScheduleID(_beneficiary, _id)];
    }

    /** @notice returns if beneficiary can stake from his schedule
     * @param _beneficiary beneficiary address
     */
    function canStake(address _beneficiary, uint256 _id) virtual external view returns (bool) {
        return schedules[_getUniqueScheduleID(_beneficiary, _id)].canStake;
    }

    /** @notice returns the number of schudeled entitled to some beneficiary
     * @param _beneficiary address of the beneficiary
     */
    function totalSchedules(address _beneficiary) virtual external view returns (uint256) {
        return beneficiarySchedules[_beneficiary].length;
    }

    /** @notice returns the list of current schedules assigned to a beneficiary
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

    /** @notice returns the amount of all unclaimed rewards (AUT) due to all the bonding from schedules entitled to beneficiary
     * @param _beneficiary beneficiary address
     */
    function unclaimedRewards(address _beneficiary) virtual external returns (uint256) {
        uint256 _totalFee = 0;
        uint256[] storage _scheduleIDs = beneficiarySchedules[_beneficiary];
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            _totalFee += _unclaimedRewards(_scheduleIDs[i]);
        }
        return _totalFee;
    }

    /** @notice returns the amount of LNTN for some schedule
     * @param _beneficiary beneficiary address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     * @param _validator validator address
     */
    function liquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _liquidBalanceOf(_getUniqueScheduleID(_beneficiary, _id), _validator);
    }

    /** @notice returns the amount of locked LNTN for some schedule */
    function lockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _lockedLiquidBalanceOf(_getUniqueScheduleID(_beneficiary, _id), _validator);
    }

    /** @notice returns the amount of unlocked LNTN for some schedule */
    function unlockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _unlockedLiquidBalanceOf(_getUniqueScheduleID(_beneficiary, _id), _validator);
    }

    /** @notice returns the list of validators bonded some schedule */
    function getBondedValidators(address _beneficiary, uint256 _id) virtual external view returns (address[] memory) {
        return _bondedValidators(_getUniqueScheduleID(_beneficiary, _id));
    }

    /** @notice returns the amount of released funds in NTN for some schedule */
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
