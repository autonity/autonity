// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.19;

import {Precompiled} from "./lib/Precompiled.sol";
import {AccessAutonity} from "./AccessAutonity.sol";
import {ILatency} from "./interfaces/ILatency.sol";

contract Latency is ILatency, AccessAutonity {
    // Frequent cluster reorganizations can lead to proposal relaying failures and reduce the robustness of the
    // semi synchronous messaging channel, which underpins our message relaying rules. This improvement aims to simplify
    // the current clustering view synchronization mechanism. Instead of triggering a cluster reorganization with
    // every measurement, we will establish a default clustering view at the start of each epoch. Following quorum
    // reports on latency measurements, an optimized clustering view will be generated using latency matrices.
    // Given that we operate in a semi-synchronized system, we assume that the emission of the KMOptimization
    // event can reach most nodes within 10 blocks. Consequently, the new clustering view will be utilized for messaging
    // at height: block.number + KM_OPTIMIZATION_DELTA.
    uint256 public constant KM_OPTIMIZATION_DELTA = 10;

    // The default latency is the median of [0, 255) which is the value space of uint8.
    // As we simplified the view synchronization with one-time reorg base on quorum reports,
    // thus for those validators who missed the report, we applied the default latency for view building.
    uint8 public constant DEFAULT_LATENCY = 128;

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
    event Reported(address indexed reporter, uint256 length);

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

    // lastReportedEpoch tracks the epoch in which each validator last reported its latency
    mapping(address => uint256) public lastReportedEpoch;

    // epochPlusOne is used to prevent the same validator to report multiple times in the same epoch
    // it is incremented on every committee change, and does not neccessarily have to be the same as the epoch in the
    // autonity contract.
    uint256 public epochPlusOne;

    // the counter counts the reported measurements of an epoch.
    uint256 public reports;

    // the block number that the optimized clusters will be activated for current epoch.
    uint256 public kmOptimizedHeight;

    constructor(address payable _autonity, address[] memory initialCommittee) AccessAutonity(_autonity) {
        committee = initialCommittee;

        // init latencies storage from precompiled contract.
        uint256 _matrixSlot;
        assembly {
            _matrixSlot := latencies.slot
        }
        uint8[] memory empty = new uint8[](0);
        require(Precompiled.updateLatency(uint256(committee.length), _matrixSlot, uint256(0), empty) == Precompiled.SUCCESS, "cannot init latency matrix");

        epochPlusOne = 1;
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

        uint256 _matrixSlot;
        assembly {
            _matrixSlot := latencies.slot
        }
        require(Precompiled.updateLatency(uint256(committee.length), _matrixSlot, index, _latency) == Precompiled.SUCCESS, "cannot insert latency report");

        lastReportedEpoch[msg.sender] = epochPlusOne;
        reports++;

        emit Reported(msg.sender, _latency.length);
        bool twoThirds = reports * 3 >= committee.length * 2;
        // Initial trigger at exactly the 2/3 crossing
        bool initialKMEvent = (reports * 3 >= committee.length * 2) && ((reports - 1) * 3 < committee.length * 2);
        // Subsequent triggers every 5 reports after crossing the 2/3 threshold OR at reaching committee length
        bool periodicKMEvent = (twoThirds && ((reports - thresholdReports()) % 5 == 0) || reports == committee.length);

        if (initialKMEvent || periodicKMEvent) {
            kmOptimizedHeight = block.number+ KM_OPTIMIZATION_DELTA;
            emit KMOptimization(kmOptimizedHeight);
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
        // update new committee and reinit the latency matrix.
        committee = _committee;

        // reset legacy matrix.
        uint256 _matrixSlot;
        assembly {
            _matrixSlot := latencies.slot
        }
        uint8[] memory empty = new uint8[](0);
        require(Precompiled.updateLatency(uint256(committee.length), _matrixSlot, uint256(0), empty) == Precompiled.SUCCESS, "cannot reset latency matrix");

        epochPlusOne++;
        reports = 0;
        kmOptimizedHeight = 0;
    }

    /*
    ┌────────────────────┐
    │ View Functions     │
    └────────────────────┘
    */

    function thresholdReports() internal view returns (uint256) {
        return (committee.length * 2 + 2) / 3; // +2 for proper ceiling division
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
    function readReport(uint256 _index) external view returns (address[] memory, uint8[] memory) {
        require(_index < committee.length, "invalid index of reporter");
        return (committee, latencies[_index]);
    }

    /// @notice Get the current committee
    /// @return The current committee node addresses
    function getCommittee() external view returns (address[] memory) {
        return committee;
    }

    /// @notice Get latency metrics status for current epoch.
    /// @return A tuple which contains the caller's  current epoch, and the KM optimization height of current epoch.
    function getMetricsStatus() external view returns (uint256) {
        return kmOptimizedHeight ;
    }
}