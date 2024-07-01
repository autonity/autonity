// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./Autonity.sol";
contract OmissionAccountability is IOmissionAccountability {

    // please follow the design doc: https://github.com/clearmatics/autonity-protocol/discussions/248 to understand the
    // omission fault detection protocol.
    struct Config {
        uint256 negligibleThreshold;        // a threshold ratio to address if a validator is an offencer at the end of epoch.
        uint256 omissionLoopBackWindow;
        uint256 activityProofRewardRate;
        uint256 maxCommitteeSize;
        uint256 pastPerformanceWeight;
        uint256 initialJailingPeriod;       // is the initial number of epoch validator will be jailed for
        uint256 initialProbationPeriod;     // is the initial number of epoch validator will be set under probation for
        uint256 initialSlashingRate;
    }

    struct HeightInactiveNodes {
        // counter for the number of distinct voting power in the activity proof, set to 0 for invalid activity proof.
        uint256 proverEffort;
        mapping(address => bool) inactiveNodes;
    }
    // todo: if we can have a ring buffer for look-back window?
    struct LookBackWindow {
        // start height of the look back window.
        uint256 start;
        // mapping height to its corresponding inactive nodes.
        mapping(uint256 => HeightInactiveNodes) lookBackWindow;
    }

    // mapping height to its corresponding inactive nodes.
    LookBackWindow public lookBackWindow;

    // todo: to get the voting power from AC committee, or we save a shadow of it in this contract.
    struct EpochProverEfforts {
        uint256 totalAccumulatedEfforts;
        mapping(address => uint256) proverEfforts;
    }

    EpochProverEfforts public epochProverEfforts;

    // EpochInactivityScores
    mapping(address => uint256) public epochInactiveBlocks;
    mapping(address => uint256) public lastEpochInactivityScores;
    mapping(address => uint256) public currentEpochInactivityScores;

    // epochCollusionDegree can be computed
    // the number of omission faulty nodes that are addressed on current epoch
    // uint256 public epochCollusionDegree;

    mapping(address => uint256) public repeatedOffences; // reset as soon as an entire probation period is completed without offences.

    mapping(address => uint256) public activityPercentage;

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