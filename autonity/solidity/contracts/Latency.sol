// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.19;

import {AccessAutonity} from "./AccessAutonity.sol";
import {ILatency} from "./interfaces/ILatency.sol";

contract Latency is ILatency, AccessAutonity {
    // Below this threshold, the committee is too small for clustering
    // to be efficient, so we limit the latency calculations
    uint256 public constant SCALE_THRESHOLD_FOR_CLUSTERING = 1;

    // Frequent cluster reorganizations can lead to proposal relaying failures and reduce the robustness of the
    // semi synchronous messaging channel, which underpins our message relaying rules. This improvement aims to simplify
    // the current clustering view synchronization mechanism. Instead of triggering a cluster reorganization with
    // every measurement, we will establish a default clustering view at the start of each epoch. Following quorum
    // reports on latency measurements, an optimized clustering view will be generated using latency matrices.
    // Given that we operate in a semi-synchronized system, we assume that the emission of the KMOptimization
    // event can reach most nodes within 5 blocks. Consequently, the new clustering view will be utilized for messaging
    // at height: block.number + KM_OPTIMIZATION_DELTA.
    uint256 public constant KM_OPTIMIZATION_DELTA = 5;

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

    // latency tracks the most recent latency report of each validator with respect to
    // every other validator in the committee
    mapping(address => mapping(address => uint8)) public latency;

    // lastReportedEpoch tracks the epoch in which each validator last reported its latency
    mapping(address => uint256) public lastReportedEpoch;

    // committee is the current committee of validators
    address[] public committee;

    // epoch is used to prevent the same validator to report multiple times in the same epoch
    // it is incremented on every committee change, and does not neccessarily have to be the same as the epoch in the
    // autonity contract.
    uint256 public epoch;

    // the counter counts the reported measurements of an epoch.
    uint256 public reports;

    // the block number that the optimized clusters will be activated for current epoch.
    uint256 public kmOptimizedHeight;

    constructor(address payable _autonity, address[] memory initialCommittee) AccessAutonity(_autonity) {
        committee = initialCommittee;
        epoch = 1;
    }

    /*
    ┌────────────┐
    │ Modifiers  │
    └────────────┘
    */
    modifier onlyCommittee(address address_) {
        bool isCommittee = false;
        for (uint256 i = 0; i < committee.length; i++) {
            if (committee[i] == address_) {
                isCommittee = true;
                break;
            }
        }
        require(isCommittee, "Latency: not committee");
        _;
    }


    modifier onlyOncePerEpoch() {
        require(lastReportedEpoch[msg.sender] < epoch, "Latency: already reported in this epoch");
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
    function report(uint8[] memory _latency) external onlyCommittee(msg.sender) onlyOncePerEpoch {
        require(_latency.length == committee.length, "Latency: invalid length");
        require(committee.length > SCALE_THRESHOLD_FOR_CLUSTERING, "Latency: committee too small");
        require(kmOptimizedHeight == 0, "Latency: clustering optimization already done for current epoch");

        for (uint256 i = 0; i < _latency.length; i++) {
            latency[msg.sender][committee[i]] = _latency[i];
        }
        lastReportedEpoch[msg.sender] = epoch;
        reports++;
        emit Reported(msg.sender, _latency.length);
        if (reports >= committee.length*2/3) {
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
        committee = _committee;
        epoch++;
        reports = 0;
        kmOptimizedHeight = 0;
    }

    /*
    ┌────────────────────┐
    │ View Functions     │
    └────────────────────┘
    */

    /// @notice Read the latency matrix
    /// @return The latency matrix for the current committee
    function read() external view returns (uint8[][] memory) {
        uint8[][] memory result = new uint8[][](committee.length);
        for (uint256 i = 0; i < committee.length; i++) {
            result[i] = new uint8[](committee.length);
            for (uint256 j = 0; j < committee.length; j++) {
                result[i][j] = latency[committee[i]][committee[j]];
            }
        }
        return result;
    }

    /// @notice Read the latency report of a specific reporter
    /// @param reporter The address of the reporter
    /// @return The latency report of the reporter
    function readReport(address reporter) external onlyCommittee(reporter) view returns (uint8[] memory) {
        uint8[] memory result = new uint8[](committee.length);
        for (uint256 i = 0; i < committee.length; i++) {
            result[i] = latency[reporter][committee[i]];
        }
        return result;
    }

    /// @notice Get the current committee
    /// @return The current committee node addresses
    function getCommittee() external view returns (address[] memory) {
        return committee;
    }

    /// @notice Get latency metrics status for current epoch.
    /// @return A tuple which contains the caller's last report epoch, current epoch, and the KM optimization height of current epoch.
    function getMetricsStatus(address reporter) external view returns (uint256, uint256, uint256) {
        return (lastReportedEpoch[reporter], epoch, kmOptimizedHeight);
    }
}