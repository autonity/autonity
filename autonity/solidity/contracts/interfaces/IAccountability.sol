// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

interface IAccountability {

    struct BaseSlashingRates {
        uint256 low;
        uint256 mid;
        uint256 high;
    }

    struct Factors {
        uint256 collusion;
        uint256 history;
        uint256 jail;
    }

    struct Config {
        uint256 innocenceProofSubmissionWindow;
        uint256 delta;
        uint256 range;
        BaseSlashingRates baseSlashingRates;
        Factors factors;
    }

    enum EventType {
        FaultProof,
        Accusation,
        InnocenceProof
    }

    // Must match autonity/types.go
    enum Rule {
        PN,
        PO,
        PVN,
        PVO,
        PVO12,
        C,
        C1,

        InvalidProposal, // The value proposed by proposer cannot pass the blockchain's validation.
        InvalidProposer, // A proposal sent from none proposer nodes of the committee.
        Equivocation    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.
    }

    enum Severity {
        Reserved, // artificially bump the starting severity value to 1
        Low,
        Mid,
        High
    }

    struct Event {
        EventType eventType; // Accountability event types: Misbehaviour, Accusation, Innocence.
        Rule rule;           // Rule ID defined in AFD rule engine.
        address reporter;    // The node address of the validator who report this event, for incentive protocol.
        address offender;    // The corresponding node address of this accountability event.
        bytes rawProof;      // rlp encoded bytes of Proof object.

        uint256 id;             // index of the event in the Events array. Will be populated internally.
        uint256 block;          // block when the event occurred. Will be populated internally.
        uint256 epoch;          // epoch when the event occurred. Will be populated internally.
        uint256 reportingBlock; // block when the event got reported. Will be populated internally.
        uint256 messageHash;    // hash of the main evidence. Will be populated internally.
    }

    /**
    * @return config, the config of the accountability contract
    */
    function getConfig() external view returns (Config memory);

    /**
    * @return gracePeriod, the current grace period in accountability
    */
    function getGracePeriod() external view returns (uint256);

    /**
    * @notice called by the Autonity Contract at block finalization, before
    * processing reward redistribution.
    * @param _epochEnd whether or not the current block is the last one from the epoch.
    * @return range, the height range for the provable fault detector
    * @return delta, the delta for the provable fault detector
    * @return gracePeriod, the current gracePeriod value
    */
    function finalize(bool _epochEnd) external returns (uint256,uint256,uint256);

    /**
    * @notice distribute slashing rewards to reporters.
    * @param _validator the address of the validator node being slashed.
    * @param _ntnReward ntn rewards to be distributed
    */
    function distributeRewards(address _validator, uint256 _ntnReward) external payable;

    /**
    * @notice called by the Autonity Contract when the committee is updated.
    * @param _committee the new committee member addresses
    */
    function setCommittee(address[] memory _committee) external;

    /**
    * @notice Event emitted when a fault proof has been submitted. The reported validator
    * will be silenced and slashed at the end of the current epoch.
    */
    event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id, uint256 _epoch);

    /**
    * @notice Event emitted when a reporter is rewarded for submitting a valid proof
    */
    event ReporterRewarded(address _reporter, address indexed _offender, uint256 _ntnReward, uint256 _atnReward);

    /**
    * @notice Event emitted after receiving an accusation, the reported validator has
    * a certain amount of time to submit a proof-of-innocence, otherwise, he gets slashed.
    */
    event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id);

    /**
    * @notice Event emitted after receiving a proof-of-innocence cancelling an accusation.
    */
    event InnocenceProven(address indexed _offender, uint256 _id);

    /**
    * @notice Event emitted after a successful slashing.
    */
    event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId);

    /**
    * @notice Event emitted after accountability factors (collusion, history and jail) are updated
    */
    event AccountabilityFactorsUpdate(Factors oldFactors, Factors newFactors);

    /**
    * @notice Event emitted after base slashing rates are updated
    */
    event BaseSlashingRateUpdate(BaseSlashingRates oldRates, BaseSlashingRates newRates);

}
