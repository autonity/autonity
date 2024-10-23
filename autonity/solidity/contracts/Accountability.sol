// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "./interfaces/IAccountability.sol";
import "./Autonity.sol";
import {SLASHING_RATE_PRECISION} from "./ProtocolConstants.sol";
import {AccessAutonity} from "./AccessAutonity.sol";

contract Accountability is IAccountability, AccessAutonity {

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

    //Todo(youssef): consider another structure purely for internal events

    Event[] public events;
    Config public config;

    // slashing rewards beneficiaries: validator => reporter
    mapping(address => address) public beneficiaries;

    mapping(address => uint256[]) private validatorFaults;

    // number of times a validator has been slashed in the past
    mapping(address => uint256) public history;

    // validatorAccusation maps a validator with an accusation
    // the id is incremented by one to handle the special case id = 0.
    mapping(address => uint256) private validatorAccusation;

    address[] internal lastCommittee;
    address[] internal curCommittee;
    mapping(address => bool) private allowedReporters;

    // mapping address => epoch => severity
    mapping (address =>  mapping(uint256 => uint256)) public slashingHistory;

    // pending slashing and accusations tasks for this epoch
    uint256[] private slashingQueue;
    uint256[] private accusationsQueue;
    uint256 internal accusationsQueueFirst = 0;

    constructor(address payable _autonity, Config memory _config) AccessAutonity(_autonity) {
        _ratesSanityCheck(_config.baseSlashingRates);
        _factorsSanityCheck(_config.factors);

        Autonity.CommitteeMember[] memory committee = autonity.getCommittee();
        for (uint256 i=0; i < committee.length; i++) {
            curCommittee.push(committee[i].addr);
            allowedReporters[committee[i].addr] = true;
        }
        config = _config;
    }

    /**
    * @notice called by the Autonity Contract at block finalization, before
    * processing reward redistribution.
    * @param _epochEnd whether or not the current block is the last one from the epoch.
    */
    function finalize(bool _epochEnd) external virtual onlyAutonity {
        // on each block, try to promote accusations without proof of innocence into misconducts.
        _promoteGuiltyAccusations();
        if (_epochEnd) {
            _performSlashingTasks();
        }
    }

   /**
    * @notice called by the Autonity Contract at block finalization, to reward the reporter of
    * a valid proof.
    * @param _offender validator account which got slashed.
    * @param _ntnReward total amount of ntn to be transferred to the repoter. MUST BE AVAILABLE
    * in the accountability contract balance.
    */
    function distributeRewards(address _offender, uint256 _ntnReward) payable external virtual onlyAutonity {
        // There is an edge-case scenario where slashing events for the
        // same accused validator are created during the same epoch.
        // In this case we only reward the last reporter.
        Autonity.Validator memory _reporter = autonity.getValidator(beneficiaries[_offender]);

        autonity.autobond(_reporter.nodeAddress, _ntnReward, 0);

        // if for some reasons, funds can't be transferred to the reporter treasury (sneaky contract)
        (bool ok, ) = _reporter.treasury.call{value:msg.value, gas: 2300}("");
        // well, too bad, it goes to the autonity global treasury.
        if(!ok) {
            autonity.getTreasuryAccount().call{value:msg.value}("");
        }
        delete beneficiaries[_offender];
    }

    /**
    * @notice Handle a misbehaviour event. Need to be called by a registered validator account
    * as the treasury-linked account will be used in case of a successful slashing event.
    */
    function handleMisbehaviour(Event memory _event) public virtual onlyAFDReporter {
        require(_event.reporter == msg.sender, "event reporter must be caller");
        require(_event.eventType == EventType.FaultProof, "wrong event type for misbehaviour");
        _handleFaultProof(_event);
    }

    /**
    * @notice Handle an accusation event. Need to be called by a registered validator account
    * as the treasury-linked account will be used in case of a successful slashing event.
    */
    function handleAccusation(Event memory _event) public virtual onlyAFDReporter {
        require(_event.reporter == msg.sender, "event reporter must be caller");
        require(_event.eventType == EventType.Accusation, "wrong event type for accusation");
        _handleAccusation(_event);
    }

    /**
    * @notice Handle an innocence proof. Need to be called by a registered validator account
    * as the treasury-linked account will be used in case of a successful slashing event.
    */
    function handleInnocenceProof(Event memory _event) public virtual onlyAFDReporter {
        require(_event.reporter == msg.sender, "event reporter must be caller");
        require(_event.eventType == EventType.InnocenceProof, "wrong event type for innocence proof");
        _handleInnocenceProof(_event);
    }

    // @dev return true if sending the event can lead to slashing
    function canSlash(address _offender, Rule _rule, uint256 _block) public virtual view returns (bool) {
        require(_rule >= Rule.PN && _rule <= Rule.Equivocation, "rule id must be valid");
        uint256 _severity = _ruleSeverity(_rule);
        uint256 _epoch = autonity.getEpochFromBlock(_block);

        return slashingHistory[_offender][_epoch] < _severity;
    }

    // @dev return true sender can accuse, can cover the cost for accusation
    function canAccuse(address _offender, Rule _rule, uint256 _block) public virtual view
    returns (bool _result, uint256 _deadline) {
        require(_rule >= Rule.PN && _rule <= Rule.Equivocation, "rule id must be valid");
        uint256 _severity = _ruleSeverity(_rule);
        uint256 _epoch = autonity.getEpochFromBlock(_block);
        if (slashingHistory[_offender][_epoch] >= _severity){
            _result = false;
            _deadline = 0;
        } else if (validatorAccusation[_offender] != 0){
            Event storage _accusation =  events[validatorAccusation[_offender] - 1];
            _result = false;
            _deadline = _accusation.block + config.innocenceProofSubmissionWindow;
        } else {
            _result = true;
            _deadline = 0;
        }
    }

    function getValidatorAccusation(address _val) public virtual view returns (Event memory){
        require(validatorAccusation[_val] > 0 , "no accusation");
        return events[validatorAccusation[_val] - 1];
    }

    function getValidatorFaults(address _val) public virtual view returns (Event[] memory){
        Event[] memory _events = new Event[](validatorFaults[_val].length);
        for(uint256 i = 0; i < validatorFaults[_val].length; i++) {
            _events[i] = events[validatorFaults[_val][i]];
        }
        return _events;
    }

    function _handleFaultProof(Event memory _ev) internal virtual {
        // Validate the misbehaviour proof
        (bool _success, address _offender, uint256 _ruleId, uint256 _block, uint256 _messageHash) =
            Precompiled.verifyAccountabilityEvent(Precompiled.MISBEHAVIOUR_CONTRACT, _ev.rawProof);

        require(_success, "failed proof verification");
        require(_offender == _ev.offender, "offender mismatch");
        require(_ruleId == uint256(_ev.rule), "rule id mismatch");
        require(_block < block.number, "can't be in the future");
        require(_block > 0, "can't be at genesis");
        
        uint256 _epoch = autonity.getEpochFromBlock(_block);
        
        _ev.block = _block;
        _ev.epoch = _epoch;
        _ev.reportingBlock = block.number;
        _ev.messageHash = _messageHash;

        _handleValidFaultProof(_ev);
    }

    function _handleValidFaultProof(Event memory _ev) internal virtual {
        uint256 _severity = _ruleSeverity(_ev.rule);
        require(slashingHistory[_ev.offender][_ev.epoch] < _severity, "already slashed at the proof's epoch");

        _ev.id = events.length;
        events.push(_ev);
        
        validatorFaults[_ev.offender].push(_ev.id);
        slashingQueue.push(_ev.id);
        slashingHistory[_ev.offender][_ev.epoch] = _severity;

        emit NewFaultProof(_ev.offender, _severity, _ev.id, autonity.epochID());
    }

    function _handleAccusation(Event memory _ev) internal virtual {
        // Validate the accusation proof. It also does height related checks 
        (bool _success, address _offender, uint256 _ruleId, uint256 _block, uint256 _messageHash) =
            Precompiled.verifyAccountabilityEvent(Precompiled.ACCUSATION_CONTRACT, _ev.rawProof);
        require(_success, "failed accusation verification");
        require(_offender == _ev.offender, "offender mismatch");
        require(_ruleId == uint256(_ev.rule), "rule id mismatch");

        uint256 _epoch = autonity.getEpochFromBlock(_block);

        _ev.block = _block;
        _ev.epoch = _epoch;
        _ev.reportingBlock = block.number;
        _ev.messageHash = _messageHash;

        _handleValidAccusation(_ev);
    }
    
    function _handleValidAccusation(Event memory _ev) internal virtual {
        require(validatorAccusation[_ev.offender] == 0, "already processing an accusation");
        uint256 _severity = _ruleSeverity(_ev.rule);
        require(slashingHistory[_ev.offender][_ev.epoch] < _severity, "already slashed at the proof's epoch");

        _ev.id = events.length;
        events.push(_ev);

        // off-by-one adjustement to hande special case id = 0
        validatorAccusation[_ev.offender] = _ev.id + 1;
        accusationsQueue.push(_ev.id + 1);

        emit NewAccusation(_ev.offender, _severity, _ev.id);
    }

    function _handleInnocenceProof(Event memory _ev) internal virtual {
        (bool _success, address _offender, uint256 _ruleId, uint256 _block, uint256 _messageHash) =
                Precompiled.verifyAccountabilityEvent(Precompiled.INNOCENCE_CONTRACT, _ev.rawProof);

        require(_success, "failed innocence verification");
        require(_offender == _ev.offender, "offender mismatch");
        require(_ruleId == uint256(_ev.rule), "rule id mismatch");
        require(_block < block.number, "can't be in the future");
        
        _ev.block = _block;
        _ev.messageHash = _messageHash;
        _ev.reportingBlock = block.number;
        _handleValidInnocenceProof(_ev);
    }

    function _handleValidInnocenceProof(Event memory _ev) internal virtual {
        uint256 _accusation = validatorAccusation[_ev.offender];
        require(_accusation != 0, "no associated accusation");

        require(events[_accusation - 1].rule == _ev.rule, "unmatching proof and accusation rule id");
        require(events[_accusation - 1].block == _ev.block, "unmatching proof and accusation block");
        require(events[_accusation - 1].messageHash == _ev.messageHash, "unmatching proof and accusation hash");

        // innocence proof is valid, remove accusation.
        for(uint256 i = accusationsQueueFirst;
                    i < accusationsQueue.length; i++){
            if(accusationsQueue[i] == _accusation ){
                accusationsQueue[i] = 0;
                break;
            }
        }
        validatorAccusation[_ev.offender] = 0;

        emit InnocenceProven(_ev.offender, 0);
    }

    /**
    * @notice Take funds away from faulty node account.
    * @dev Emit a {SlashingEvent} event for the fined account
    */
    function _slash(Event memory _event, uint256 _epochOffencesCount, uint256 _epochPeriod) internal virtual {
        address _offender = _event.offender;

        // last reporter is the beneficiary
        beneficiaries[_offender] = _event.reporter;

        // if already jailbound, validator has 0 stake
        if (autonity.getValidatorState(_offender) == ValidatorState.jailbound) {
            return;
        }

        uint256 _baseRate = _baseSlashingRate(_ruleSeverity(_event.rule));

        uint256 _slashingRate = _baseRate +
            (_epochOffencesCount * config.factors.collusion) +
            (history[_offender] * config.factors.history);

        history[_offender] += 1;
        uint256 _jailtime = config.factors.jail * history[_offender] * _epochPeriod;

        (uint256 _slashingAmount, uint256 _jailReleaseBlock, bool _isJailbound) = autonity.slashAndJail(
            _offender,
            _slashingRate,
            _jailtime,
            ValidatorState.jailed,
            ValidatorState.jailbound
        );
        emit SlashingEvent(_offender, _slashingAmount, _jailReleaseBlock, _isJailbound, _event.id);
    }

    /**
    * @notice perform slashing over faulty validators at the end of epoch. The fine in stake token are moved from
    * validator account to autonity contract account, and the corresponding slash counter as a reputation for validator
    * increase too.
    * @dev Emit a {NodeSlashed} event for every account that are slashed.
    */
    function _performSlashingTasks() internal virtual {
        // Find the total number of offences submitted during the current epoch
        // as the slashing rate depends on it.
        uint256 _offensesCount;
        uint256 _currentEpoch = autonity.epochID();
        for (uint256 i = 0; i < slashingQueue.length; i++) {
            if(events[slashingQueue[i]].epoch == _currentEpoch){
                _offensesCount += 1;
            }
        }

        uint256 _epochPeriod = autonity.getCurrentEpochPeriod();

        for (uint256 i = 0; i < slashingQueue.length; i++) {
            _slash(events[slashingQueue[i]], _offensesCount, _epochPeriod);
        }
        // reset pending slashing task queue for next epoch.
        delete slashingQueue;
    }


    /**
    * @notice promote accusations without innocence proof in the proof submission into misbehaviour.
    */
    function _promoteGuiltyAccusations() internal virtual {
        uint256 i = accusationsQueueFirst;
        uint256 _epochID = autonity.epochID();
        for(; i < accusationsQueue.length; i++){
            uint256 _id = accusationsQueue[i];
            if (_id == 0) {
                continue;
            }
            _id -= 1; // shift by one to handle event id = 0
            Event memory _ev = events[_id];
            if(_ev.reportingBlock + config.innocenceProofSubmissionWindow > block.number) {
                // The queue is ordered by time of submission so we can break here.
                break;
            }
            delete validatorAccusation[_ev.offender];
            uint256 _severity = _ruleSeverity(_ev.rule);
             if(slashingHistory[_ev.offender][_ev.epoch] >= _severity){
                // we skip this accusation as a fault proof has been reported during the submission window.
                continue;
            }
            slashingHistory[_ev.offender][_ev.epoch] = _severity;
            validatorFaults[_ev.offender].push(_id);
            slashingQueue.push(_id);

            emit NewFaultProof(_ev.offender, _severity, _id, _epochID);
        }
        accusationsQueueFirst = i;
    }

    function _ruleSeverity(Rule _rule) internal virtual pure returns (uint256) {
        if (_rule == Rule.Equivocation) {
            return uint256(Severity.Low); // can easily be done by mistake by starting two clients with the same key
        }
        if (_rule == Rule.PN) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.PO) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.PVN) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.PVO) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.PVO12) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.C) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.C1) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.InvalidProposal) {
            return uint256(Severity.High); // can only be committed if purposely modifying the client
        }
        if (_rule == Rule.InvalidProposer) {
            return uint256(Severity.High); // can only be committed if purposely modifying the client
        }
        revert("unknown rule");
    }

    function _baseSlashingRate(uint256 _severity) internal virtual view returns (uint256) {
        if (_severity == uint256(Severity.Low)) {
            return config.baseSlashingRates.low;
        }
        if (_severity == uint256(Severity.Mid)) {
            return config.baseSlashingRates.mid;
        }
        if (_severity == uint256(Severity.High)) {
            return config.baseSlashingRates.high;
        }
        revert("unknown severity");
    }

    function _grantReportAccess(address[] memory _members) internal virtual onlyAutonity {
        for (uint256 i=0; i < _members.length; i++) {
            if (allowedReporters[_members[i]] == true) {
                continue;
            }
            allowedReporters[_members[i]] = true;
        }
    }

    function _revokeReportAccess(address[] memory _members) internal virtual onlyAutonity {
        for (uint256 i=0; i < _members.length; i++) {
            if (allowedReporters[_members[i]] == false) {
                continue;
            }
            delete allowedReporters[_members[i]];
        }
    }

    function _ratesSanityCheck(BaseSlashingRates memory _rates) internal virtual {
        require(_rates.low <= _rates.mid, "slashing rate: low cannot exceed mid");
        require(_rates.mid <= _rates.high, "slashing rate: mid cannot exceed high");
        require(_rates.high <= SLASHING_RATE_PRECISION, "slashing rate: high cannot exceed precision");
    }

    function _factorsSanityCheck(Factors memory _factors) internal virtual {
        require(_factors.collusion <= SLASHING_RATE_PRECISION, "collusion factor cannot exceed slashing rate");
        require(_factors.history <= SLASHING_RATE_PRECISION, "history factor cannot exceed slashing rate");
    }


    /*
    ============================================================
        Setters
    ============================================================
    */

    /**
    * @notice setCommittee, called by the AC at epoch change, it removes stale
    * committee from the reporter set, then replace the last committee with current
    * committee, and set the current committee with the input new committee.
    * @dev restricted to the autonity contract
    * @param _newCommittee, the committee for the new epoch
    */
    function setCommittee(address[] memory _newCommittee) external virtual onlyAutonity {
        // revoke report access for stale committee.
        _revokeReportAccess(lastCommittee);

        // as last revoke might withdraw access for overlapped committee, thus
        // make sure they are still in the reporter set, then replace last committee
        // with current committee.
        _grantReportAccess(curCommittee);
        lastCommittee = curCommittee;

        // grant access for new committee and replace current committee with new one.
        _grantReportAccess(_newCommittee);
        curCommittee = _newCommittee;
    }

    /*
    * @notice sets the innocence proof submission window
    * @dev restricted to the operator
    * @param _window, the new value for the window (in blocks)
    */
    function setInnocenceProofSubmissionWindow(uint256 _window) external virtual onlyOperator {
        config.innocenceProofSubmissionWindow = _window;
    }

    /*
    * @notice sets the base slashing rates
    * @dev restricted to the operator
    * @param _rates, the new rates
    */
    function setBaseSlashingRates(BaseSlashingRates memory _rates) external virtual onlyOperator {
        _ratesSanityCheck(_rates);
        config.baseSlashingRates = _rates;
    }

    /*
    * @notice sets the new punishment factors
    * @dev restricted to the operator
    * @param _factor, the new factor
    */
    function setFactors(Factors memory _factors) external virtual onlyOperator {
        _factorsSanityCheck(_factors);
        config.factors = _factors;
    }

    /*
    ============================================================
        Modifiers
    ============================================================
    */

    /**
    * @dev Modifier that checks if the caller is an AFD reporter (a committee member of last or current epoch).
    */
    modifier onlyAFDReporter {
        require(allowedReporters[msg.sender] == true, "function restricted to a committee member");
        _;
    }
}
