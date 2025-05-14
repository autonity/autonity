// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.19;

import {Precompiled} from "./lib/Precompiled.sol";
import {AccessAutonity} from "./AccessAutonity.sol";
import {ILatency} from "./interfaces/ILatency.sol";

contract Latency is ILatency, AccessAutonity {
    // The default latency is the median of [0, 255) which is the value space of uint8.
    // As we simplified the view synchronization with one-time reorg base on quorum reports,
    // thus for those validators who missed the report, we applied the default latency for view building.
    uint8 public constant DEFAULT_LATENCY = 128;
    uint256 public constant LOCK_IN_THRESHOLD_DENOMINATOR = 1000;

    /*
    ┌────────┐
    │ Events │
    └────────┘
    */

    /**
     * @dev Emitted when a validator reports its latency
     * @param reporter The address of the validator who reported the latency
     * @param length The length of the reported latency array
     */
    event Reported(address indexed reporter, uint256 length, uint totalReported);

    /**
     * @dev Emitted when there are quorum measurements metrics is ready for K Means optimization.
     * @param height the height when the optimized clusters will be activated for messaging.
     */
    event KMOptimization(uint256 height);

    /*
    ┌────────┐
    │ State  │
    └────────┘
    */

    // latency tracks the latency report of each validator of an epoch.
    // It is a N * N matrix pre-allocated by the protocol on epoch rotation, N is the length of committee.
    uint8[][] private latencies;
    // committee is the current committee of validators
    address[] public committee;

    // last latencies save last epoch's latencies, it is used for those who recovered from a disaster to build the KM
    // clusters before a new KM clusters is constructed for current epoch. It is operated by precompile contract.
    uint8[][] private lastLatencies;
    address[] public lastCommittee;

    // lastReportedEpoch tracks the epoch in which each validator last reported its latency
    mapping(address => uint256) public lastReportedEpoch;

    // epochPlusOne is used to prevent the same validator to report multiple times in the same epoch
    // it is incremented on every committee change, and does not neccessarily have to be the same as the epoch in the
    // autonity contract.
    uint256 public epochPlusOne;

    // the counter counts the reported measurements of an epoch.
    uint256 public reports;

    // the last matrix finalization block
    uint256 public matrixLockInBlock;
    uint256 public lastMatrixLockInBlock;

    // The threshold of reports to trigger the K Means optimization (base 1000).
    uint256 public lockInThreshold;
    uint256 public lockInDelay;

    constructor(address payable _autonity, address[] memory initialCommittee) AccessAutonity(_autonity) {
        committee = initialCommittee;
        epochPlusOne = 1;
        lockInThreshold = 600; // 60%
        lockInDelay = 5; // 5 blocks
    }

    modifier onlyOncePerEpoch() {
        require(lastReportedEpoch[msg.sender] < epochPlusOne, "Latency: already reported in this epoch");
        _;
    }

    /*
    ┌────────────────────┐
    │ External Functions │
    └────────────────────┘
    */

    /// @notice Report the latency to the contract
    /// @param _latency The latency array (must match committee length)
    /// @dev The length of the latency array must be equal to the length of the committee
    function report(uint256 index, uint8[] memory _latency) external onlyOncePerEpoch {
        require(_latency.length == committee.length, "Latency: invalid length");
        require(index < committee.length, "Latency: invalid reporter index");
        require(committee[index] == msg.sender, "Latency: not a valid reporter");

        if (matrixLockInBlock != 0 && block.number > matrixLockInBlock) {
            emit Reported(msg.sender, _latency.length, reports);
            lastReportedEpoch[msg.sender] = epochPlusOne;
            reports++;
            return;
        }

        uint256 _reportsSlot;
        uint256 _matrixSlot;
        assembly {
            _reportsSlot := reports.slot
            _matrixSlot := latencies.slot
        }
        require(Precompiled.updateLatency(0, _reportsSlot, _matrixSlot, index, _latency) == Precompiled.SUCCESS, "cannot insert latency report");

        lastReportedEpoch[msg.sender] = epochPlusOne;
        reports++;

        emit Reported(msg.sender, _latency.length, reports);
        bool thresholdMet = reports * LOCK_IN_THRESHOLD_DENOMINATOR >= committee.length * lockInThreshold;
        if (thresholdMet) {
            // other validators can report until the end of the block, but a new event should not be triggered
            matrixLockInBlock = block.number;
            emit KMOptimization(block.number + lockInDelay);
        }
    }

    /*
    ┌────────────────────┐
    │ Autonity Functions │
    └────────────────────┘
    */

    /// @notice Set the committee for the new epoch
    /// @param _committee The new committee
    /// @dev This function is intended to be called by the Autonity contract
    function setCommittee(address[] memory _committee) external onlyAutonity {

        // save last matrix
        uint256 _matrixSlot;
        uint256 _lastMatrixSlot;
        uint8[] memory row;
        assembly {
            _matrixSlot := latencies.slot
            _lastMatrixSlot := lastLatencies.slot
        }
        require(Precompiled.updateLatency(_lastMatrixSlot, 0, _matrixSlot, 0, row) == Precompiled.SUCCESS, "cannot store last epoch latencies");
        // save last committee
        lastCommittee = committee;

        // update new committee.
        committee = _committee;
        epochPlusOne++;
        reports = 0;
        lastMatrixLockInBlock = matrixLockInBlock;
        matrixLockInBlock = 0;
    }

    /*
    ┌────────────────────┐
    │ View Functions     │
    └────────────────────┘
    */

    function thresholdReports() internal view returns (uint256) {
        return (committee.length * 2 + 2) / 3; // +2 for proper ceiling division
    }

    function readLastEpochReport(uint256 _index) external view returns (uint8[] memory) {
        require(_index < lastCommittee.length, "invalid index of reporter");
        require(epochPlusOne > 1, "1st epoch is not over yet");
        return lastLatencies[_index];
    }

    function readLast() external view returns (address[] memory, uint8[][] memory) {
        require(epochPlusOne > 1, "1st epoch is not over yet");
        return (lastCommittee, lastLatencies);
    }

    /// @notice Read the entire latency matrix, it would be failed once the committee size scale to a large number.
    //  It would be better to use readReport(_index) which returns a single reporters report, This is one of the bottleneck
    //  of the on-chain latency data solution.
    /// @return The committee and the latency matrix for the current epoch.
    function read() external view returns (address[] memory, uint8[][] memory) {
        return (committee, latencies);
    }

    /// @notice Read the report by reporter index.
    /// @return The committee and the latency report of the reporter.
    function readReport(uint256 _index) external view returns (uint8[] memory) {
        require(_index < committee.length, "invalid index of reporter");
        return latencies[_index];
    }


    /// @notice Check if the caller already reported for current epoch.
    /// @return True if the caller reported.
    function clientReported(uint256 _index) external view returns (bool) {
        require(_index < committee.length, "invalid index of reporter");
        if (lastReportedEpoch[committee[_index]] == epochPlusOne) {
            return true;
        }
        return false;
    }

    /// @notice Get the current committee
    /// @return The current committee node addresses
    function getCommittee() external view returns (address[] memory) {
        return committee;
    }

    /// @notice Get the last committee
    /// @return The last committee node addresses
    function getLastCommittee() external view returns (address[] memory) {
        return lastCommittee;
    }

    /// @notice Get latency metrics status for current epoch.
    /// @return A tuple which contains the caller's  current epoch, and the KM optimization height of current epoch.
    function getMetricsStatus() external view returns (uint256) {
        return matrixLockInBlock;
    }
}