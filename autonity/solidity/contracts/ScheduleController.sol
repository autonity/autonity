// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract ScheduleController {

    struct Schedule {
        uint256 totalAmount;
        uint256 unlockedAmount;
        uint256 start;
        uint256 totalDuration;
        uint256 lastUnlockTime;
    }

    mapping(address => Schedule[]) internal vaultSchedules;
    address[] internal vaults;

    constructor() {}

    /**
     * @notice Creates a new schedule.
     * @param _scheduleVault address of the vault which holds the token for this schedule
     * @param _amount total amount of the schedule
     * @param _startTime start time
     * @param _totalDuration total duration of the schedule
     */
    function _createSchedule(
        address _scheduleVault,
        uint256 _amount,
        uint256 _startTime,
        uint256 _totalDuration
    ) internal {
        Schedule[] storage _schedules = vaultSchedules[_scheduleVault];
        if (_schedules.length == 0) {
            vaults.push(_scheduleVault);
        }
        _schedules.push(Schedule(_amount, 0, _startTime, _totalDuration, 0));
    }

    function _unlockSchedules(uint256 _unlockTime) internal returns (uint256 _newUnlocked) {
        for (uint256 _vaultIndex = 0; _vaultIndex < vaults.length; _vaultIndex++) {
            Schedule[] storage _schedules = vaultSchedules[vaults[_vaultIndex]];
            for (uint256 _scheduleIndex = 0; _scheduleIndex < _schedules.length; _scheduleIndex++) {
                Schedule storage _schedule = _schedules[_scheduleIndex];
                require(_unlockTime >= _schedule.lastUnlockTime, "schedule already unlocked for given time");
                if (_unlockTime <= _schedule.start) {
                    continue;
                }
                _schedule.lastUnlockTime = _unlockTime;
                uint256 _unlocked;
                if (_unlockTime - _schedule.start >= _schedule.totalDuration) {
                    _unlocked = _schedule.totalAmount;
                }
                else {
                    _unlocked = (_unlockTime - _schedule.start) * _schedule.totalAmount / _schedule.totalDuration;
                }
                _newUnlocked += _unlocked - _schedule.unlockedAmount;
                _schedule.unlockedAmount = _unlocked;
            }
        }
    }

    /**
     * @notice Returns the schedule at index = `_id` in the `vaultSchedules[_vault]` array.
     * @param _vault address of the vault for the schedule
     * @param _id index of the schedule
     */
    function getSchedule(address _vault, uint256 _id) public view returns (Schedule memory) {
        Schedule[] storage _schedules = vaultSchedules[_vault];
        require(_schedules.length > _id, "schedule does not exist");
        return _schedules[_id];
    }

    /**
     * Returns total number of schedules for the vault at address `_vault`.
     * @param _vault address of the vault for the schedules
     */
    function getTotalSchedules(address _vault) public view returns (uint256) {
        return vaultSchedules[_vault].length;
    }
}