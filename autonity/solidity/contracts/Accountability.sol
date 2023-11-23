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
        uint256 slashingRatePrecision;
    }

    uint256 public epochPeriod;

    enum EventType {
        FaultProof,
        Accusation,
        InnocenceProof
    }

    // Must match autonity/accountability_types.go
    enum Rule {
        PN,
        PO,
        PVN,
        PVO,
        PVO12,
        PVO3,
        C,
        C1,

        InvalidProposal, // The value proposed by proposer cannot pass the blockchain's validation.
        InvalidProposer, // A proposal sent from none proposer nodes of the committee.
        Equivocation,    // Multiple distinguish votes(proposal, prevote, precommit) sent by validator.

        InvalidRound,          // message contains invalid round number or step.
        WrongValidRound, // message was signed by sender, but it cannot be decoded.
        GarbageMessage // message sender is not the member of current committee.
    }

    enum Severity {
        Minor,
        Low,
        Mid,
        High,
        Critical
    }

    struct Event {
        uint8 chunks;        // Counter of number of chunks for oversize accountability event
        uint8 chunkId;       // Chunk index to construct the oversize accountability event
        EventType eventType; // Accountability event types: Misbehaviour, Accusation, Innocence.
        Rule rule;           // Rule ID defined in AFD rule engine.
        address reporter;    // The node address of the validator who report this event, for incentive protocol.
        address offender;    // The corresponding node address of this accountability event.
        bytes rawProof;      // rlp encoded bytes of Proof object.

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

    // validatorAccusation maps a validator with an accusation
    // the id is incremented by one to handle the special case id = 0.
    mapping(address => uint256) private validatorAccusation;

    mapping(address => Event) internal reporterChunksMap;

    // mapping address => epoch => severity
    mapping (address =>  mapping(uint256 => uint256)) public slashingHistory;

    // pending slashing and accusations tasks for this epoch
    uint256[] private slashingQueue;
    uint256[] private accusationsQueue;
    uint256 internal accusationsQueueFirst = 0;

    constructor(address payable _autonity, Config memory _config){
        autonity = Autonity(_autonity);
        epochPeriod = autonity.getEpochPeriod();

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

    function distributeRewards(address _validator) payable external onlyAutonity {
        // There is an edge-case scenario where slashing events for the
        // same accused validator are created during the same epoch.
        // In this case we only reward the last reporter.
        address _reporterTreasury = autonity.getValidator(beneficiaries[_validator]).treasury;
        // if for some reasons, funds can't be transferred to the reporter treasury (sneaky contract)
        (bool ok, ) = _reporterTreasury.call{value:msg.value, gas: 2300}("");
        // well, too bad, it goes to the autonity global treasury.
        if(!ok) {
            autonity.getTreasuryAccount().call{value:msg.value}("");
        }
        delete beneficiaries[_validator];
    }

    /**
    * @notice Handle an accountability event. Need to be called by a registered validator account
    * as the treasury-linked account will be used in case of a successful slashing event.
    * todo(youssef): rethink modifiers here, consider splitting this into multiple functions.
    */
    function handleEvent(Event memory _event) public onlyValidator {
        require(_event.reporter == msg.sender, "event reporter must be caller");
        // if the event is a chunked event, store it.
        if (_event.chunks > 1) {
            bool _readyToProcess = _storeChunk(_event);
            // for the last chunk we should process directly the event from here.
            if (!_readyToProcess) {
                return;
            }
            _event = reporterChunksMap[msg.sender];
        }

        if (_event.eventType == EventType.FaultProof) {
            _handleFaultProof(_event);
            return;
        }
        if (_event.eventType == EventType.Accusation) {
            _handleAccusation(_event);
            return;
        }
        if (_event.eventType == EventType.InnocenceProof) {
            _handleInnocenceProof(_event);
            return;
        }
        // todo(youssef): consider reverting here to help with refund
        return;
    }

    // @dev return true if sending the event can lead to slashing
    function canSlash(address _offender, Rule _rule, uint256 _block) public view returns (bool) {
        uint256 _severity = _ruleSeverity(_rule);
        uint256 _epoch = autonity.getEpochFromBlock(_block);

        return slashingHistory[_offender][_epoch] < _severity ;
    }

    // @dev return true sender can accuse, can cover the cost for accusation
    function canAccuse(address _offender, Rule _rule, uint256 _block) public view
    returns (bool _result, uint256 _deadline) {
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

        // in the case of a garbage message, where the inner consensus message can't be decoded
        // we assign the last block number as the block height as the epoch id is not yet available
        if(_block == 0){
            _block = block.number - 1;
        }
        
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

        events.push(_ev);
        uint256 _eventId = events.length - 1;
        validatorFaults[_ev.offender].push(_eventId);
        slashingQueue.push(_eventId);
        slashingHistory[_ev.offender][_ev.epoch] = _severity;

        emit NewFaultProof(_ev.offender, _severity, _eventId);
    }

    function _handleAccusation(Event memory _ev) internal {
        // Validate the accusation proof
        (bool _success, address _offender, uint256 _ruleId, uint256 _block, uint256 _messageHash) =
            Precompiled.verifyAccountabilityEvent(Precompiled.ACCUSATION_CONTRACT, _ev.rawProof);
        require(_success, "failed accusation verification");
        require(_offender == _ev.offender, "offender mismatch");
        require(_ruleId == uint256(_ev.rule), "rule id mismatch");
        require(_block < block.number, "can't be in the future");

        uint256 _epoch = autonity.getEpochFromBlock(_block);

        _ev.block = _block;
        _ev.epoch = _epoch;
        _ev.reportingBlock = block.number;
        _ev.messageHash = _messageHash;

        _handleValidAccusation(_ev);
    }
    
    function _handleValidAccusation(Event memory _ev) internal {
        require(validatorAccusation[_ev.offender] == 0, "already processing an accusation");
        uint256 _severity = _ruleSeverity(_ev.rule);
        require(slashingHistory[_ev.offender][_ev.epoch] < _severity, "already slashed at the proof's epoch");

        events.push(_ev);
        uint256 _eventId = events.length - 1;

        // off-by-one adjustement to hande special case id = 0
        validatorAccusation[_ev.offender] = _eventId + 1;
        accusationsQueue.push(_eventId + 1);

        emit NewAccusation(_ev.offender, _severity, _eventId);
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

    // @dev only supporting one chunked event per sender
    function _storeChunk(Event memory _event) internal returns (bool){
        // saving a chunk with id 0 will reset the local store
        if (_event.chunkId == 0) {
            reporterChunksMap[msg.sender] = _event;
            return false;
        }
        require(reporterChunksMap[msg.sender].chunkId + 1 == _event.chunkId, "chunks must be contiguous");
        BytesLib.concatStorage(reporterChunksMap[msg.sender].rawProof, _event.rawProof);
        reporterChunksMap[msg.sender].chunkId += 1;
        // return true if it's the final chunk, to be processed immediately
        return _event.chunkId + 1 == _event.chunks;
    }

    /**
    * @notice Take funds away from faulty node account.
    * @dev Emit a {SlashingEvent} event for the fined account or {ValidatorJailbound} event for being jailed permanently
    */
    function _slash(Event memory _event, uint256 _epochOffencesCount) internal {
        // The assumption here is that the node hasn't been slashed yet for the proof's epoch.
        //_val must be returned - no error check
        Autonity.Validator memory _val = autonity.getValidator(_event.offender);
        // last reporter is the beneficiary
        beneficiaries[_event.offender] = _event.reporter;

        uint256 _baseRate = _baseSlashingRate(_ruleSeverity(_event.rule));
        uint256 _history = _val.provableFaultCount;

        uint256 _slashingRate = _baseRate +
            (_epochOffencesCount * config.collusionFactor) +
            ( _history * config.historyFactor);

        if(_slashingRate > config.slashingRatePrecision) {
            _slashingRate = config.slashingRatePrecision;
        }

        uint256 _availableFunds = _val.bondedStake + _val.unbondingStake + _val.selfUnbondingStake;
        uint256 _slashingAmount =  (_slashingRate * _availableFunds)/config.slashingRatePrecision;

        // in case of (_slashingAmount = _availableFunds - 1) or 100% slash, we slash all stakes and jailbound the validator
        if (_slashingAmount > 0 && _slashingAmount >= _availableFunds - 1) {
            _val.bondedStake = 0;
            _val.selfBondedStake = 0;
            _val.selfUnbondingStake = 0;
            _val.unbondingStake = 0;
            _val.totalSlashed += _availableFunds;
            _val.provableFaultCount += 1;
            _val.state = ValidatorState.jailbound;
            autonity.updateValidatorAndTransferSlashedFunds(_val);
            emit ValidatorJailbound(_val.nodeAddress, _availableFunds);
            return;
        }
        uint256 _remaining = _slashingAmount;
        // -------------------------------------------
        // Implementation of Penalty Absorbing Stake
        // -------------------------------------------
        // Self-unbonding stake gets slashed in priority.
        if(_val.selfUnbondingStake >= _remaining){
            _val.selfUnbondingStake -= _remaining;
            _remaining = 0;
        } else {
            _remaining -= _val.selfUnbondingStake;
            _val.selfUnbondingStake = 0;
        }
        // Then self-bonded stake
        if (_remaining > 0){
            if(_val.selfBondedStake >= _remaining) {
                _val.selfBondedStake -= _remaining;
                _val.bondedStake -= _remaining;
                _remaining = 0;
            } else {
                _remaining -= _val.selfBondedStake;
                _val.bondedStake -= _val.selfBondedStake;
                _val.selfBondedStake = 0;
            }
        }
        // --------------------------------------------
        // Remaining stake to be slashed is split equally between the delegated
        // stake pool and the non-self unbonding stake pool.
        // As a reminder, the delegated stake pool is bondedStake - selfBondedStake.
        // if _remaining > 0 then bondedStake = delegated stake, because all selfBondedStake is slashed
        if (_remaining > 0 && (_val.unbondingStake + _val.bondedStake > 0)) {
            uint256 _unbondingSlash = 0;
            uint256 _delegatedSlash = 0;
            // as we cannot store fraction here, we are taking floor for the smaller one between unbondingStake and bondedStake
            // and ceil for the larger one. In case both variable unbondingStake and bondedStake are positive, this modification
            // will ensure that no variable reaches 0 too fast where the other one is too big. In this case the bigger one
            // will reach 0 first, and the smaller one will be 0 or 1.
            // That means the fairness issue: https://github.com/autonity/autonity/issues/819 will only be triggered if 100% stake
            // is slashed or (slashingAmount = totalStake - 1)
            if (_val.unbondingStake <= _val.bondedStake) {
                _unbondingSlash = (_remaining * _val.unbondingStake) /
                                        (_val.unbondingStake + _val.bondedStake);
                _delegatedSlash = _remaining - _unbondingSlash;
            } else {
                _delegatedSlash = (_remaining * _val.bondedStake) /
                                        (_val.unbondingStake + _val.bondedStake);
                _unbondingSlash = _remaining - _delegatedSlash;
            }
            _val.unbondingStake -= _unbondingSlash;
            _val.bondedStake -= _delegatedSlash;

        }

        _val.totalSlashed += _slashingAmount;
        _val.provableFaultCount += 1;
        _val.jailReleaseBlock = block.number + config.jailFactor * _val.provableFaultCount * epochPeriod;
        _val.state = ValidatorState.jailed; // jailed validators can't participate in consensus

        autonity.updateValidatorAndTransferSlashedFunds(_val);

        emit SlashingEvent(_val.nodeAddress, _slashingAmount, _val.jailReleaseBlock);
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
        //
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
    * @dev Modifier that checks if the caller is the slashing contract.
    */
    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to the validator");
        _;
    }

    /**
    * @dev Modifier that checks if the caller is a registered validator.
    */
    modifier onlyValidator {
        Autonity.Validator memory _val = autonity.getValidator(msg.sender);
        require(_val.nodeAddress == msg.sender, "function restricted to a registered validator");
        _;
    }

}
