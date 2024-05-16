// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../interfaces/INonStakableVestingVault.sol";
import "./ScheduleBase.sol";

contract NonStakableVesting is INonStakableVestingVault, ScheduleBase {

    struct ScheduleClass {
        uint256 start;
        uint256 cliff;
        uint256 end;
        uint256 totalAmount;
        uint256 totalUnlocked;
        uint256 lastUnlockTime;
    }

    // stores all the schedule classes, there should not be too many of them, for the sake of efficiency
    // of unlockTokens() function
    ScheduleClass[] private scheduleClasses;

    // class id that some schedule is subscribed to
    mapping(uint256 => uint256) private classID;

    constructor(
        address payable _autonity, address _operator
    ) ScheduleBase(_autonity, _operator) {}

    /**
     * @notice creates a new class of schedule, restricted to operator
     * the class has totalAmount = 0 initially. As new schedules are subscribed to the class, its totalAmount increases
     * At any point, totalAmount of class is the sum of totalValue all the schedules that are subscribed to the class.
     * totalValue of a schedule can be calculated via _calculateTotalValue function
     * @param _startTime start time
     * @param _cliffTime cliff time. cliff period = _cliffTime - _startTime
     * @param _endTime end time, total duration of the schedule = _endTime - _startTime
     */
    function createScheduleClass(
        uint256 _startTime,
        uint256 _cliffTime,
        uint256 _endTime
    ) virtual public onlyOperator {
        scheduleClasses.push(ScheduleClass(_startTime, _cliffTime, _endTime, 0, 0, 0));
    }

    /**
     * @notice creates a new non-stakable schedule, restricted to only operator
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
        _classData.totalAmount += _amount;
        classID[_scheduleID] = _scheduleClass;
    }

    /**
     * @notice used by beneficiary to transfer all unlocked NTN of some schedule to his own address
     * @param _id id of the schedule numbered from 0 to (n-1) where n = total schedules entitled to
     * the beneficiary (excluding canceled ones). So any beneficiary can number their schedules
     * from 0 to (n-1). Beneficiary does not need to know the unique global schedule id which can
     * be retrieved via _getUniqueScheduleID function
     */
    function releaseAllFunds(uint256 _id) virtual external { // onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _releaseNTN(_scheduleID, _unlockedFunds(_scheduleID));
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice used by beneficiary to transfer some amount of unlocked NTN of some schedule to his own address
     * @param _amount amount of NTN to release
     */
    function releaseFund(uint256 _id, uint256 _amount) virtual external { // onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        require(_amount <= _unlockedFunds(_scheduleID), "not enough unlocked funds");
        _releaseNTN(_scheduleID, _amount);
    }

    /**
     * @notice changes the beneficiary of some schedule to the _recipient address. _recipient can release tokens from the schedule
     * only operator is able to call the function
     * @param _beneficiary beneficiary address whose schedule will be canceled
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (excluding canceled ones)
     * @param _recipient whome the schedule is transferred to
     */
    function cancelSchedule(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator {
        _cancelSchedule(_beneficiary, _id, _recipient);
    }

    /**
     * @notice Unlock tokens of all schedules upto current time, restricted to autonity only.
     * @dev It calculates the newly unlocked tokens upto current time and also updates the amount
     * of total unlocked tokens and the time of unlock for each class of schedule
     * Autonity must mint _totalNewUnlocked tokens, because this contract knows that for each _class,
     * _class.totalUnlocked tokens are now unlocked and available to release
     */
    function unlockTokens() external onlyAutonity returns (uint256 _totalNewUnlocked) {
        uint256 _currentTime = block.timestamp;
        for (uint256 i = 0; i < scheduleClasses.length; i++) {
            ScheduleClass storage _class = scheduleClasses[i];
            if (_class.cliff > _currentTime || _class.totalAmount == _class.totalUnlocked) {
                continue;
            }
            _class.lastUnlockTime = _currentTime;
            uint256 _unlocked = _calculateUnlockedFunds(_class.start, _class.end, _currentTime, _class.totalAmount);
            _totalNewUnlocked += _unlocked - _class.totalUnlocked;
            _class.totalUnlocked = _unlocked;
        }
        return _totalNewUnlocked;
    }

    /**
     * @dev calculates the total value of the schedule, which is constant for non stakable schedules
     * @param _scheduleID unique global id of the schedule
     */
    function _calculateTotalValue(uint256 _scheduleID) private view returns (uint256) {
        Schedule storage _schedule = schedules[_scheduleID];
        return _schedule.currentNTNAmount + _schedule.withdrawnValue;
    }

    /**
     * @dev calculates the amount of funds that are unlocked but not released yet. calculates upto _class.lastUnlockTime
     * where _class = schedule class subsribed by the schedule.
     * The unlock mechanism is epoch based, but instead of taking time from autonity.lastEpochBlock(), we take the time
     * from _class.lastUnlockTime. Because the locked tokens are not minted from genesis. This way it is ensured that
     * the unlocked tokens are minted by calling the function unlockTokens()
     */
    function _unlockedFunds(uint256 _scheduleID) private view returns (uint256) {
        return _calculateAvailableUnlockedFunds(
            _scheduleID,
            _calculateTotalValue(_scheduleID),
            scheduleClasses[classID[_scheduleID]].lastUnlockTime
        );
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice returns the amount of unlocked but not yet released funds in NTN for some schedule
     */
    function unlockedFunds(
        address _beneficiary, uint256 _id
    ) virtual external view returns (uint256) {
        return _unlockedFunds(_getUniqueScheduleID(_beneficiary, _id));
    }
}