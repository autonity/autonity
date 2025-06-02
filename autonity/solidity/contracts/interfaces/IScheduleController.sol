// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IScheduleController {
    struct Schedule {
        uint256 totalAmount;
        uint256 unlockedAmount;
        uint256 start;
        uint256 totalDuration;
        uint256 lastUnlockTime;
    }

    /**
     * @notice Returns the schedule at index = `_id` in the `vaultSchedules[_vault]` array.
     * @param _vault address of the vault for the schedule
     * @param _id index of the schedule
     */
    function getSchedule(address _vault, uint256 _id) external view returns (Schedule memory);

    /**
     * Returns total number of schedules for the vault at address `_vault`.
     * @param _vault address of the vault for the schedules
     */
    function getTotalSchedules(address _vault) external view returns (uint256);
}