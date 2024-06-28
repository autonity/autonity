// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./Autonity.sol";
contract OmissionAccountability is IOmissionAccountability {

    // please follow the design doc: https://github.com/clearmatics/autonity-protocol/discussions/248 to understand the
    // omission fault detection protocol.
    struct Config {
        uint256 omissionLoopBackWindow;
        uint256 activityProofRewardRate;
        uint256 maxCommitteeSize;
        uint256 pastPerformanceWeight;
        uint256 initialJailingPeriod;       // is the initial number of epoch validator will be jailed for
        uint256 initialProbationPeriod;     // is the initial number of epoch validator will be set under probation for
        uint256 initialSlashingRate;
    }

    // todo: review this look back window.
    struct HeightInactiveNodes {
        // counter for the number of distinct voting power in the activity proof, set to 0 for invalid activity proof.
        uint256 proverEffort;
        mapping(address => bool) inactiveNodes;
    }

    struct LookBackWindow {
        // start height of the look back window.
        uint256 start;

        // mapping height to its corresponding inactive nodes.
        mapping(uint256 => HeightInactiveNodes) lookBackWindow;
    }

    // mapping height to its corresponding inactive nodes.
    LookBackWindow public lookBackWindow;

    // todo: review this prover effort counter of an epoch, it reset on epoch rotation.
    struct EpochProverEfforts {
        uint256 totalAccumulatedEfforts;
        mapping(address => uint256) proverEfforts;
    }

    EpochProverEfforts public epochProverEfforts;

    // todo: review this epoch inactivity score, we need to save at least two epoch as we have different weights for each.
    struct EpochInactivityScores {
        uint256 startEpochID;

        mapping(address => uint256) accumulatedInactiveBlocks;

        // mapping epochID => (address => score);
        mapping(uint256 => mapping(address => uint256)) scores;
    }

    // todo: review the collusion degree data structure
    uint256 public epochCollusionDegree; // the number of omission faulty nodes that are addressed on current epoch.

    EpochInactivityScores public epochInactivityScores;

    // todo: review the jailing and probation metrics.
    mapping(address => uint256) public repeatedOffences; // reset as soon as an entire probation period is completed without offences.

    // todo: review the average activity percentage.
    mapping(address => uint256) public activityPercentage;

    // todo: review the overall omission faulty heights
    mapping(address => uint256) public overallFaultyBlocks;

    Config public config;
    Autonity internal autonity; // for access control in setters function.

    constructor(address payable _autonity, Config memory _config) {
        autonity = Autonity(_autonity);
        config = _config;
    }

    /**
    * @notice called by the Autonity Contract at block finalization, it receives activity report.
    * @param isProposerOmissionFaulty is true when the proposer provides invalid activity proof of current height.
    * @param ids stores faulty proposer's ID when isProposerOmissionFaulty is true, otherwise it carries current height
    * activity proof which is the signers of precommit of current height - dela.
    */
    function finalize(bool isProposerOmissionFaulty, uint256[] memory ids) external onlyAutonity {
        // todo, fill the look-back window by the input data.
    }

    /**
    * @dev Modifier that checks if the caller is the slashing contract.
    */
    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to the Autonity Contract");
        _;
    }
}