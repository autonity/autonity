// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";
import "./ScheduleBase.sol";

contract NonStakableVestingManager is ScheduleBase {

    struct ScheduleClass {
        uint256 start;
        uint256 cliff;
        uint256 end;
        uint256 totalAmount;
        uint256 totalUnlocked;
        uint256 lastUnlockBlock;
    }

    ScheduleClass[] private scheduleClasses;

    mapping(uint256 => uint256) private classID;

    constructor(
        address payable _autonity, address _operator
    ) ScheduleBase(_autonity, _operator) {}

    function createScheduleClass(
        uint256 _startBlock,
        uint256 _cliffBlock,
        uint256 _endBlock
    ) virtual public onlyOperator {
        scheduleClasses.push(ScheduleClass(_startBlock, _cliffBlock, _endBlock, 0, 0, 0));
    }

    /**
     * @notice creates a new stakable schedule, restricted to only operator
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _scheduleClass schedule class to subscribe
     */
    function newSchedule(
        address _beneficiary,
        uint256 _amount,
        uint256 _scheduleClass
    ) virtual onlyOperator public {
        require(_scheduleClass < scheduleClasses.length, "invalid schedule class");
        ScheduleClass storage _classData = scheduleClasses[_scheduleClass];
        uint256 _scheduleID = _createSchedule(
            _beneficiary, _amount, _classData.start, _classData.cliff, _classData.end, false
        );
        classID[_scheduleID] = _scheduleClass;
    }

    /**
     * @notice used by beneficiary to transfer all unlocked NTN of some schedule to his own address
     */
    function releaseAllFunds(uint256 _id) virtual external onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _releaseNTN(_scheduleID, _unlockedFunds(_scheduleID));
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice used by beneficiary to transfer some amount of unlocked NTN of some schedule to his own address
     * @param _amount amount of NTN to release
     */
    function releaseFund(uint256 _id, uint256 _amount) virtual external onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        require(_amount <= _unlockedFunds(_scheduleID), "not enough unlocked funds");
        _releaseNTN(_scheduleID, _amount);
    }

    /**
     * @notice release of all unlocked NTN of some schedule and return them to the _recipient account
     * effectively cancelling a vesting schedule. only operator is able to call the function
     * @param _beneficiary beneficiary address whose schedule will be canceled
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     * @param _recipient to whome the all funds will be transferred
     */
    function cancelSchedule(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator onlyActive(_id) {
        // TODO: remove the schedule from its class
        // TODO: we can only transfer unlocked tokens, because locked tokens are not minted yet. what to do with it?
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        _removeScheduleFromClass(_scheduleID);
        Schedule storage _schedule = schedules[_scheduleID];
        _schedule.canceled = true;
        // locked tokens are not minted yet, so we can transfer only unlocked tokens
        _updateAndTransferNTN(_scheduleID, _recipient, _unlockedFunds(_scheduleID));
    }

    /**
     * @dev Unlock tokens of all schedules upto current block, restricted to autonity only
     */
    function unlockTokens() external onlyAutonity returns (uint256) {
        uint256 _currentBlock = block.number;
        uint256 _totalNewUnlocked;
        for (uint256 i = 0; i < scheduleClasses.length; i++) {
            ScheduleClass storage _class = scheduleClasses[i];
            if (_class.cliff > _currentBlock || _class.totalAmount == _class.totalUnlocked) {
                continue;
            }
            _class.lastUnlockBlock = _currentBlock;
            uint256 _unlocked = _calculateUnlockedFunds(_class.start, _class.end, _currentBlock, _class.totalAmount);
            _totalNewUnlocked += _unlocked - _class.totalUnlocked;
            _class.totalUnlocked = _unlocked;
        }
        return _totalNewUnlocked;
    }

    function _calculateTotalValue(uint256 _scheduleID) private view returns (uint256) {
        Schedule storage _schedule = schedules[_scheduleID];
        return _schedule.currentNTNAmount + _schedule.withdrawnValue;
    }

    function _unlockedFunds(uint256 _scheduleID) private view returns (uint256) {
        return _calculateUnlockedFundsAtHeight(
            _scheduleID,
            _calculateTotalValue(_scheduleID),
            scheduleClasses[classID[_scheduleID]].lastUnlockBlock
        );
    }

    function _removeScheduleFromClass(uint256 _scheduleID) private {
        ScheduleClass storage _class = scheduleClasses[classID[_scheduleID]];
        Schedule storage _schedule = schedules[_scheduleID];
        _class.totalUnlocked -= _unlockedFunds(_scheduleID) + _schedule.withdrawnValue;
        _class.totalAmount -= _calculateTotalValue(_scheduleID);
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice returns the amount of released funds in NTN for some schedule
     */
    function unlockedFunds(
        address _beneficiary, uint256 _id
    ) virtual external view onlyActive(_id) returns (uint256) {
        return _unlockedFunds(_getUniqueScheduleID(_beneficiary, _id));
    }
}