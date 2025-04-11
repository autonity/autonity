// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
/**
 * @dev Interface of the Latency Contract.
 * Used in the Autonity Contract.
 */

interface ILatency {

    /**
     * @notice Report the latency between the caller and the committee members.
     * @param _latency the latency values between the caller and the committee members.
     */
    function report(uint256 index, uint8[] memory _latency) external;

    /**
     * @notice Read the latency values between the committee members.
     * @return the latency values between the committee members.
     */
    function read() external view returns (address[] memory, uint8[][] memory);

    /**
     * @notice Set the committee members.
     * @param _committee the new committee member addresses.
     */
    function setCommittee(address[] memory _committee) external;
}