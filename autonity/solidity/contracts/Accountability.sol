// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "./interfaces/IAccountability.sol";
import "./Autonity.sol";

contract Accountability is IAccountability {

    struct Config {
        uint256 innocenceProofSubmissionWindow;

        // Slashing parameters
        uint256 baseSlashingRateLow;
        uint256 baseSlashingRateMid;
        uint256 collusionFactor;
        uint256 historyFactor;
        uint256 jailFactor;
    }

    uint256 public epochPeriod;
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
        Minor,
        Low,
        Mid,
        High,
        Critical
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

    Autonity internal autonity;
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
    // counters of rising accusation by a validator of epoch.
    // mapping address => epoch => counter
    mapping(address => mapping(uint256 => uint256)) private epochAccusationCounter;

    // mapping address => epoch => severity
    mapping (address =>  mapping(uint256 => uint256)) public slashingHistory;

    // pending slashing and accusations tasks for this epoch
    uint256[] private slashingQueue;
    uint256[] private accusationsQueue;
    uint256 internal accusationsQueueFirst = 0;

    constructor(address payable _autonity, Config memory _config) {
        autonity = Autonity(_autonity);
        epochPeriod = autonity.getCurrentEpochPeriod();
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
    function finalize(bool _epochEnd) external onlyAutonity {
        // on each block, try to promote accusations without proof of innocence into misconducts.
        _promoteGuiltyAccusations();
        if (_epochEnd) {
            _performSlashingTasks();
        }
    }

   /**
    * @notice called by the Autonity Contract at block finalization, to reward the reporter of
    * a valid proof.
    * @param _validator validator account which got slashed.
    * @param _ntnReward total amount of ntn to be transferred to the repoter. MUST BE AVAILABLE
    * in the accountability contract balance.
    */
    function distributeRewards(address _validator, uint256 _ntnReward) payable external onlyAutonity {
        // There is an edge-case scenario where slashing events for the
        // same accused validator are created during the same epoch.
        // In this case we only reward the last reporter.
        address _reporterTreasury = autonity.getValidator(beneficiaries[_validator]).treasury;

        try autonity.transfer(_reporterTreasury, _ntnReward) {}
        catch {}

        // if for some reasons, funds can't be transferred to the reporter treasury (sneaky contract)
        (bool ok, ) = _reporterTreasury.call{value:msg.value, gas: 2300}("");
        // well, too bad, it goes to the autonity global treasury.
        if(!ok) {
            autonity.getTreasuryAccount().call{value:msg.value}("");
        }
        delete beneficiaries[_validator];
    }

    /**
    * @notice Handle a misbehaviour event. Need to be called by a registered validator account
    * as the treasury-linked account will be used in case of a successful slashing event.
    */
    function handleMisbehaviour(Event memory _event) public onlyAFDReporter {
        require(_event.reporter == msg.sender, "event reporter must be caller");
        require(_event.eventType == EventType.FaultProof, "wrong event type for misbehaviour");
        _handleFaultProof(_event);
    }

    /**
    * @notice Handle an accusation event. Need to be called by a registered validator account
    * as the treasury-linked account will be used in case of a successful slashing event.
    */
    function handleAccusation(Event memory _event) public onlyAFDReporter {
        require(_event.reporter == msg.sender, "event reporter must be caller");
        require(_event.eventType == EventType.Accusation, "wrong event type for accusation");
        _handleAccusation(_event);
    }

    /**
    * @notice Handle an innocence proof. Need to be called by a registered validator account
    * as the treasury-linked account will be used in case of a successful slashing event.
    */
    function handleInnocenceProof(Event memory _event) public onlyAFDReporter {
        require(_event.reporter == msg.sender, "event reporter must be caller");
        require(_event.eventType == EventType.InnocenceProof, "wrong event type for innocence proof");
        _handleInnocenceProof(_event);
    }

    // @dev return true if sending the event can lead to slashing
    function canSlash(address _offender, Rule _rule, uint256 _block) public view returns (bool) {
        require(_rule >= Rule.PN && _rule <= Rule.Equivocation, "rule id must be valid");
        uint256 _severity = _ruleSeverity(_rule);
        uint256 _epoch = autonity.getEpochFromBlock(_block);

        return slashingHistory[_offender][_epoch] < _severity;
    }

    // @dev return true sender can accuse, can cover the cost for accusation
    function canAccuse(address _offender, Rule _rule, uint256 _block) public view
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

    function getValidatorAccusation(address _val) public view returns (Event memory){
        require(validatorAccusation[_val] > 0 , "no accusation");
        return events[validatorAccusation[_val] - 1];
    }

    function getValidatorFaults(address _val) public view returns (Event[] memory){
        Event[] memory _events = new Event[](validatorFaults[_val].length);
        for(uint256 i = 0; i < validatorFaults[_val].length; i++) {
            _events[i] = events[validatorFaults[_val][i]];
        }
        return _events;
    }

    function _handleFaultProof(Event memory _ev) internal {
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

    function _handleValidFaultProof(Event memory _ev) internal{
        uint256 _severity = _ruleSeverity(_ev.rule);
        require(slashingHistory[_ev.offender][_ev.epoch] < _severity, "already slashed at the proof's epoch");

        _ev.id = events.length;
        events.push(_ev);
        
        validatorFaults[_ev.offender].push(_ev.id);
        slashingQueue.push(_ev.id);
        slashingHistory[_ev.offender][_ev.epoch] = _severity;

        emit NewFaultProof(_ev.offender, _severity, _ev.id);
    }

    function _handleAccusation(Event memory _ev) internal {
        // Validate the accusation proof. It also does height related checks 
        (bool _success, address _offender, uint256 _ruleId, uint256 _block, uint256 _messageHash) =
            Precompiled.verifyAccountabilityEvent(Precompiled.ACCUSATION_CONTRACT, _ev.rawProof);
        require(_success, "failed accusation verification");
        require(_offender == _ev.offender, "offender mismatch");
        require(_ruleId == uint256(_ev.rule), "rule id mismatch");

        uint256 _epoch = autonity.getEpochFromBlock(_block);
        uint256 _curEpoch = autonity.epochID();

        _ev.block = _block;
        _ev.epoch = _epoch;
        _ev.reportingBlock = block.number;
        _ev.messageHash = _messageHash;

        // By according to the canAccuse(), in an EPOCH, the accusations over a validator are graded by severity, a
        // validator cannot be accused more then once unless it breaks higher severity rules, thus rate limit of rising
        // accusation during an EPOCH is calculated as: (len(committee)-1) * the max levels of rule severity. As all the
        // severity is MID now, we cap it to: (len(committee)-1) * (Severity.Mid+1) of per EPOCH for a reporter. However,
        // in reality, one is not accusable if it is processing a pending accusation, thus a malicious reporter cannot
        // DoS the protocol since the abused accusations shall be reverted without reimbursement, and it has to wait
        // until the suspected one become accusable again for innocenceProofSubmissionWindow blocks.
        uint256 _cap = (curCommittee.length-1)*(uint256(Severity.Mid)+1);
        if (_ev.epoch < _curEpoch) {
            _cap = (lastCommittee.length-1)*(uint256(Severity.Mid)+1);
        }
        require(epochAccusationCounter[msg.sender][_epoch] <= _cap, "validator shouldn't abuse accusation in an epoch");
        _handleValidAccusation(_ev);
        epochAccusationCounter[msg.sender][_epoch] += 1;
    }
    
    function _handleValidAccusation(Event memory _ev) internal {
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

    function _handleInnocenceProof(Event memory _ev) internal {
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

    function _handleValidInnocenceProof(Event memory _ev) internal {
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
    function _slash(Event memory _event, uint256 _epochOffencesCount) internal {
        address _offender = _event.offender;

        // last reporter is the beneficiary
        beneficiaries[_offender] = _event.reporter;

        // if already jailbound, validator has 0 stake
        if (autonity.getValidatorState(_offender) == ValidatorState.jailbound) {
            return;
        }

        uint256 _baseRate = _baseSlashingRate(_ruleSeverity(_event.rule));

        uint256 _slashingRate = _baseRate +
            (_epochOffencesCount * config.collusionFactor) +
            (history[_offender] * config.historyFactor);

        history[_offender] += 1;
        uint256 _jailtime = config.jailFactor * history[_offender] * epochPeriod;

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
    function _performSlashingTasks() internal {
        // Find the total number of offences submitted during the current epoch
        // as the slashing rate depends on it.
        uint256 _offensesCount;
        uint256 _currentEpoch = autonity.epochID();
        for (uint256 i = 0; i < slashingQueue.length; i++) {
            if(events[slashingQueue[i]].epoch == _currentEpoch){
                _offensesCount += 1;
            }
        }

        for (uint256 i = 0; i < slashingQueue.length; i++) {
            _slash(events[slashingQueue[i]], _offensesCount);
        }
        // reset pending slashing task queue for next epoch.
        delete slashingQueue;
    }


    /**
    * @notice promote accusations without innocence proof in the proof submission into misbehaviour.
    */
    function _promoteGuiltyAccusations() internal {
        uint256 i = accusationsQueueFirst;
        for(; i < accusationsQueue.length; i++){
            uint256 _id = accusationsQueue[i];
            if (_id == 0) {
                continue;
            }
            _id -= 1; // shift by one to handle event id = 0
            Event memory _ev = events[_id];
            //todo(youssef): complete
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

            emit NewFaultProof(_ev.offender, _severity, _id);
        }
        accusationsQueueFirst = i;
    }

    function _ruleSeverity(Rule _rule) internal pure returns (uint256) {
        if (_rule == Rule.Equivocation) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.PN) {
            return uint256(Severity.Mid);
        }
        if (_rule == Rule.PO) {
            return uint256(Severity.Mid);
        }
        // todo(youssef): finish
        return uint256(Severity.Mid);
    }

    function _baseSlashingRate(uint256 _severity) internal view returns (uint256) {
        // TODO: complete
        if (_severity == uint256(Severity.Minor)) {
            return config.baseSlashingRateMid;
        }
        if (_severity == uint256(Severity.Low)) {
            return config.baseSlashingRateMid;
        }
        if (_severity == uint256(Severity.Mid)) {
            return config.baseSlashingRateMid;
        }
        if (_severity == uint256(Severity.High)) {
            return config.baseSlashingRateMid;
        }
        if (_severity == uint256(Severity.Critical)) {
            return 10000;
        }
        return 10000;
    }

    function setEpochPeriod(uint256 _newPeriod) external onlyAutonity{
        epochPeriod = _newPeriod;
    }
    /**
    * @notice setCommittee is called by autonity only on new epoch, it remove stale
    * committee from the reporter set, then replace the last committee with current
    * committee, and set the current committee with the input new committee.
    */
    function setCommittee(address[] memory _newCommittee) external onlyAutonity{
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

    function _grantReportAccess(address[] memory _members) internal onlyAutonity {
        for (uint256 i=0; i < _members.length; i++) {
            if (allowedReporters[_members[i]] == true) {
                continue;
            }
            allowedReporters[_members[i]] = true;
        }
    }

    function _revokeReportAccess(address[] memory _members) internal onlyAutonity {
        for (uint256 i=0; i < _members.length; i++) {
            if (allowedReporters[_members[i]] == false) {
                continue;
            }
            delete allowedReporters[_members[i]];
        }
    }

    /**
    * @dev Modifier that checks if the caller is the slashing contract.
    */
    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to the Autonity protocol contract");
        _;
    }

    /**
    * @dev Modifier that checks if the caller is an AFD reporter (a committee member of last or current epoch).
    */
    modifier onlyAFDReporter {
        require(allowedReporters[msg.sender] == true, "function restricted to a committee member");
        _;
    }
}
