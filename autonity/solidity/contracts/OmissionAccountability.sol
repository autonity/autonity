// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./Autonity.sol";

contract OmissionAccountability is IOmissionAccountability {
    // Used for fixed-point arithmetic during computation of inactivity score
    uint256 public constant SCALE_FACTOR = 10_000;

    struct Config {
        uint256 inactivityThreshold;        // threshold to determine if a validator is an offender at the end of epoch.
        uint256 lookbackWindow;
        uint256 pastPerformanceWeight;
        uint256 initialJailingPeriod;       // initial number of epoch an offender will be jailed for
        uint256 initialProbationPeriod;     // initial number of epoch an offender will be set under probation for
        uint256 initialSlashingRate;

        // TODO(lorenzo) this parameter should be updated by accountability when it changes there
        // OR the change is triggered by the autonity contract, which then updates both accountability and omission
        uint256 slashingRatePrecision;      // should be the same one used in Accountability.sol

        // number of blocks to wait before generating activity proof.
        // e.g. activity proof of block x is for block x - delta
        uint256 delta;
    }

    // shadow copies of variables in Autonity.sol, updated once a epoch
    Autonity.CommitteeMember[] private committee;
    address[] private treasuries; // treasuries of the committee members
    uint256 private lastEpochBlock;

    uint256 private newLookbackWindow; // applied at epoch end
    uint256 private newDelta; // applied at epoch end
    address private operator;

    mapping(uint256 => bool) public faultyProposers;                         // marks height where proposer is faulty
    mapping(uint256 => mapping(address=>bool)) public inactiveValidators;    // inactive validators for each height
    // counter of inactive blocks for each validator (considering lookback window). It is reset at the end of the epoch.
    mapping(address => uint256) public inactivityCounter;

    // net (total - quorum) proposer effort included in the activity proof. Reset at epoch end.
    uint256 public totalEffort;

    mapping(address => uint256) public proposerEffort;

    // epoch inactivity score for each committee member. Updated at every epoch.
    mapping(address => uint256) public inactivityScores;

    mapping(address => uint256) public probationPeriods; // in epochs
    mapping(address => uint256) public repeatedOffences; // reset as soon as an entire probation period is completed without offences.

    uint256[] public epochCollusionDegree; // maps epoch number to the collusion degree

    Config public config;
    Autonity internal autonity; // for access control in setters function.

    event InactivitySlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound);
    event DeltaUpdated(uint256 delta);

    constructor(address payable _autonity, address _operator, address[] memory _treasuries, Config memory _config) {
        autonity = Autonity(_autonity);
        operator = _operator;
        config = _config;
        Autonity.CommitteeMember[] memory _committee = autonity.getCommittee();
        for(uint256 i=0;i<_committee.length;i++){
            committee.push(_committee[i]);
        }
        treasuries = _treasuries;

        newLookbackWindow = config.lookbackWindow;
        newDelta = config.delta;
    }

    /**
    * @notice called by the Autonity Contract at block finalization.
    * @param _epochEnded, true if this is the last block of the epoch
    */
    function finalize(bool _epochEnded) external virtual onlyAutonity {
        // if we are at the first delta blocks of the epoch, the activity proof should be empty
        bool _mustBeEmpty = block.number <= lastEpochBlock + config.delta;

        uint256[1] memory _committeeSlot; // declare it as array to easily access from assembly
        assembly{
            mstore(_committeeSlot,committee.slot)
        }

        (bool _isProposerOmissionFaulty, uint256 _proposerEffort, address[] memory _absentees) = Precompiled.computeAbsentees(_mustBeEmpty, config.delta, _committeeSlot[0]);

        // short-circuit function if the proof has to be empty
        if(_mustBeEmpty){
            return;
        }

        uint256 targetHeight = block.number - config.delta;

        if (_isProposerOmissionFaulty) {
            faultyProposers[targetHeight] = true;
            inactivityCounter[block.coinbase]++;
        }else{
            faultyProposers[targetHeight] = false;
            proposerEffort[block.coinbase] += _proposerEffort;
            totalEffort += _proposerEffort;

            _recordAbsentees(_absentees, targetHeight);
        }

        if(_epochEnded){
            uint256 collusionDegree = _computeInactivityScoresAndCollusionDegree();
            _punishInactiveValidators(collusionDegree);

            // reset inactivity counters
            for(uint256 i=0;i<committee.length;i++) {
                inactivityCounter[committee[i].addr] = 0;
            }

            // store collusion degree in state. This is useful for slashed validators to verify their slashing rate
            epochCollusionDegree.push(collusionDegree);

            // update lookback window and delta if changed
            config.lookbackWindow = newLookbackWindow;
            if(config.delta != newDelta) {
                emit DeltaUpdated(newDelta);
            }
            config.delta = newDelta;

        }
    }

    function _recordAbsentees(address[] memory _absentees, uint256 targetHeight) internal virtual {
        for(uint256 i=0; i < _absentees.length; i++) {
            inactiveValidators[targetHeight][_absentees[i]] = true;
        }

        if(targetHeight < lastEpochBlock + config.lookbackWindow) {
            return;
        }

        // for each absent of target height, check the lookback window to see if he was online at some point
        // if online even once in the lookback window, consider him online for this block
        // NOTE: the current block is included in the window, (h - delta - lookback, h - delta]
        for(uint256 i=0; i < _absentees.length; i++) {
            bool confirmedAbsent = true;
            uint256 initialLookBackWindow = config.lookbackWindow;
            for(uint256 h = targetHeight-1; h >targetHeight-initialLookBackWindow; h--) {
                if(faultyProposers[h]) {
                    // we do not have data for h, extend the lookback window if possible
                    if(targetHeight-lastEpochBlock <= initialLookBackWindow) {
                        // we do not have enough blocks to extend the window. let's consider the validator not absent.
                        confirmedAbsent=false;
                        break;
                    }
                    // we can extend the window
                    initialLookBackWindow++;
                    continue;
                }

                // if the validator is not found in even only one of the inactive lists, it is not considered offline
                if(!inactiveValidators[h][_absentees[i]]){
                    confirmedAbsent = false;
                    break;
                }
            }
            // if the absentee was absent for the entirety of the lookback period, increment his inactivity counter
            if (confirmedAbsent) {
                inactivityCounter[_absentees[i]]++;
            }
        }
    }

    // returns collusion degree
    function _computeInactivityScoresAndCollusionDegree() internal virtual returns (uint256) {
        uint256 epochPeriod = autonity.getEpochPeriod();
        uint256 collusionDegree = 0;

        // compute aggregated scores + collusion degree
        for(uint256 i=0;i<committee.length;i++){
            address nodeAddress = committee[i].addr;

            // first config.lookbackWindow-1 blocks of the epoch are accountable, but we do not have enough info to determine if a validator was offline/online
            // last delta blocks of the epoch are not accountable due to committee change
            uint256 inactivityScore = (inactivityCounter[nodeAddress]*SCALE_FACTOR / (epochPeriod-config.lookbackWindow+1-config.delta));

            // there is an edge case where inactivityScore can be > 100%. We cap it at 100%.
            // this can happen for example if we have a network with a single validator, that is never including any activity proof,
            // thus always being considered a faulty proposer and getting his inactivityCounter increased even when we do not have lookback blocks yet
            if(inactivityScore > SCALE_FACTOR){
                inactivityScore = SCALE_FACTOR;
            }

            uint256 aggregatedInactivityScore = ((inactivityScore*(SCALE_FACTOR-config.pastPerformanceWeight)) + (inactivityScores[nodeAddress] * config.pastPerformanceWeight))/SCALE_FACTOR;
            if(aggregatedInactivityScore > config.inactivityThreshold){
                collusionDegree++;
            }
            inactivityScores[nodeAddress] = aggregatedInactivityScore;
        }
        return collusionDegree;
    }

    function _punishInactiveValidators(uint256 collusionDegree) internal virtual {
        // reduce probation periods + dish out punishment
        for(uint256 i=0;i<committee.length;i++){
            address nodeAddress = committee[i].addr;
            Autonity.Validator memory _val = autonity.getValidator(nodeAddress);

            // if the validator has already been slashed by accountability in this epoch,
            // do not punish him for omission too. It would be unfair since peer ignore msgs from jailed vals.
            // However, do not decrease his probation since he was not fully honest
            // NOTE: validator already jailed by accountability are nonetheless taken into account into the collusion degree of omission
            if(_val.state == ValidatorState.jailed || _val.state == ValidatorState.jailbound){
                continue;
            }

            // here validator is either active or has been paused in the current epoch (but still participated to consensus)

            if(inactivityScores[nodeAddress] <= config.inactivityThreshold){
                // NOTE: probation period of a validator gets decreased only if he is part of the committee
                if(probationPeriods[nodeAddress] > 0){
                    probationPeriods[nodeAddress]--;
                    // if decreased to zero, then zero out also the offences counter
                    if(probationPeriods[nodeAddress] == 0){
                        repeatedOffences[nodeAddress] = 0;
                    }
                }
            }else{
                // punish validator if his inactivity is greater than threshold
                repeatedOffences[nodeAddress]++;
                uint256 offenceSquared = repeatedOffences[nodeAddress]*repeatedOffences[nodeAddress];
                uint256 jailingPeriod = config.initialJailingPeriod * offenceSquared;
                uint256 probationPeriod = config.initialProbationPeriod * offenceSquared;

                _val.jailReleaseBlock = block.number + jailingPeriod;
                _val.state = ValidatorState.jailed;

                // if already on probation, slash
                if(probationPeriods[nodeAddress] > 0){
                    _slash(_val, config.initialSlashingRate*offenceSquared*collusionDegree);
                }else{
                    autonity.updateValidatorAndTransferSlashedFunds(_val);
                }

                // whether slashed or not, update the probation period (cumulatively)
                probationPeriods[nodeAddress] += probationPeriod;
            }
        }
    }

    // similar logic as Accountability.sol _slash function, with a few tweaks.
    // If updating this func, probably makes sense to update the one in Accountability.sol as well.
    function _slash(Autonity.Validator memory _val, uint256 _slashingRate) internal virtual {
        if(_slashingRate > config.slashingRatePrecision) {
            _slashingRate = config.slashingRatePrecision;
        }

        uint256 _availableFunds = _val.bondedStake + _val.unbondingStake + _val.selfUnbondingStake;
        uint256 _slashingAmount =  (_slashingRate * _availableFunds)/config.slashingRatePrecision;

        // in case of 100% slash, we jailbound the validator
        if (_slashingAmount > 0 && _slashingAmount == _availableFunds) {
            _val.bondedStake = 0;
            _val.selfBondedStake = 0;
            _val.selfUnbondingStake = 0;
            _val.unbondingStake = 0;
            _val.totalSlashed += _slashingAmount;
            _val.state = ValidatorState.jailbound;
            _val.jailReleaseBlock = 0;
            autonity.updateValidatorAndTransferSlashedFunds(_val);
            emit InactivitySlashingEvent(_val.nodeAddress, _slashingAmount, 0, true);
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
            // as we cannot store fraction here, we are taking floor for both unbondingSlash and delegatedSlash
            // In case both variable unbondingStake and bondedStake are positive, this modification
            // will ensure that no variable reaches 0 too fast where the other one is too big. In this case both variables
            // will reach 0 only when slashed 100%.
            // That means the fairness issue: https://github.com/autonity/autonity/issues/819 will only be triggered
            // if 100% stake is slashed
            uint256 _unbondingSlash = (_remaining * _val.unbondingStake) /
                (_val.unbondingStake + _val.bondedStake);
            uint256 _delegatedSlash = (_remaining * _val.bondedStake) /
                (_val.unbondingStake + _val.bondedStake);
            _val.unbondingStake -= _unbondingSlash;
            _val.bondedStake -= _delegatedSlash;
            _remaining -= _unbondingSlash + _delegatedSlash;
        }

        // if positive amount remains
        _slashingAmount -= _remaining;
        _val.totalSlashed += _slashingAmount;

        autonity.updateValidatorAndTransferSlashedFunds(_val);

        emit InactivitySlashingEvent(_val.nodeAddress, _slashingAmount, _val.jailReleaseBlock, false);
    }

    /*
    * @notice called by the Autonity contract at epoch finalization, to redistribute the proposer rewards based on the effort
    * @param _ntnRewards, amount of NTN reserved for proposer rewards
    */
    function distributeProposerRewards(uint256 _ntnReward) external payable virtual onlyAutonity {
        uint256 atnReward = msg.value;

        for(uint256 i=0; i < committee.length; i++) {
            address nodeAddress = committee[i].addr;
            if(proposerEffort[nodeAddress] > 0){
               uint256 atnProposerReward = (proposerEffort[nodeAddress] * atnReward) / totalEffort;
               uint256 ntnProposerReward = (proposerEffort[nodeAddress] * _ntnReward) / totalEffort;

               // if for some reasons, funds can't be transferred to the treasury (sneaky contract)
               (bool ok, ) = treasuries[i].call{value: atnProposerReward, gas: 2300}("");
               // well, too bad, it goes to the autonity global treasury.
               if(!ok) {
                   // TODO(lorenzo) check return value?
                   autonity.getTreasuryAccount().call{value:atnProposerReward}("");
               }

               autonity.transfer(treasuries[i],ntnProposerReward);

               // reset after usage
               proposerEffort[nodeAddress] = 0;
           }
        }

        totalEffort = 0;
    }

    /*
    * @notice get the inactivity score of a validator for the last finalized epoch
    * @param node address for the validator
    * @return the inactivity score of the validator in the last finalized epoch
    */
    function getInactivityScore(address _validator) external view virtual returns (uint256) {
        return inactivityScores[_validator];
    }

    /*
    * @notice gets the scale factor used for fixed point computations in this contract
    * @return the scale factor used for fixed point computations
    */
    function getScaleFactor() external pure virtual returns (uint256) {
        return SCALE_FACTOR;
    }

    /*
    * @notice gets the delta used to determine how many block to wait before generating the activity proof
    * @return the delta number of blocks to wait before generating the activity proof
    */
    function getDelta() external view virtual returns (uint256) {
        return newDelta;
    }

    /*
    * @notice retrieves the lookback window value and whether an update of it is in progress.
    * @return the lookback window current value
    */
    function getLookbackWindow() external view virtual returns (uint256) {
        return newLookbackWindow;
    }

    /*
    * @notice gets the total proposer effort accumulated up to this block
    * @return the total proposer effort accumulated up to this block
    */
    function getTotalEffort() external view virtual returns (uint256) {
        return totalEffort;
    }

    /*
    * @notice sets committee node addresses and treasuries
    * @dev restricted to the Autonity contract. It is used to mirror this information in the omission contract when the autonity contract changes
    * @param _committee, committee members
    * @param _treasuries, treasuries of the new committee
    */
    function setCommittee(Autonity.CommitteeMember[] memory _committee, address[] memory _treasuries) external virtual onlyAutonity{
        delete committee;
        for(uint256 i=0;i<_committee.length;i++){
            committee.push(_committee[i]);
        }
        treasuries = _treasuries;
    }

    /* @notice sets the lastEpochBlock in the omission contract
    * @dev restricted to the Autonity contract. It is used to mirror this information when it is updated at epoch finalize.
    * @param _lastEpochBlock, last block of the past epoch
    */
    function setLastEpochBlock(uint256 _lastEpochBlock) external virtual onlyAutonity {
        lastEpochBlock = _lastEpochBlock;
    }

    /* @notice sets the operator in the omission contract
    * @dev restricted to the Autonity contract. It is used to mirror the operator account.
    * @param _operator, the new operator account
    */
    function setOperator(address _operator) external virtual onlyAutonity {
        operator = _operator;
    }

    // config update methods

    /* @notice sets the inactivity threshold
    * @dev restricted to the operator
    * @param _inactivityThreshold, the new value for inactivity threshold
    */
    function setInactivityThreshold(uint256 _inactivityThreshold) external virtual onlyOperator {
        require(_inactivityThreshold <= SCALE_FACTOR, "cannot exceed scale factor");
        config.inactivityThreshold = _inactivityThreshold;
    }

    /* @notice sets the past performance weight
    * @dev restricted to the operator
    * @param _pastPerformanceWeight, the new value for the past performance weight
    */
    function setPastPerformanceWeight(uint256 _pastPerformanceWeight) external virtual onlyOperator{
        require(_pastPerformanceWeight <= SCALE_FACTOR, "cannot exceed scale factor");
        require(_pastPerformanceWeight <= config.inactivityThreshold, "pastPerformanceWeight cannot be greater than inactivityThreshold");
        config.pastPerformanceWeight = _pastPerformanceWeight;
    }

    /* @notice sets the initial jailing period
    * @dev restricted to the operator
    * @param _initialJailingPeriod, the new value for the initial jailing period
    */
    function setInitialJailingPeriod(uint256 _initialJailingPeriod) external virtual onlyOperator {
        config.initialJailingPeriod = _initialJailingPeriod;
    }

    /* @notice sets the initial probation period
    * @dev restricted to the operator
    * @param _initialProbationPeriod, the new value for the initial probation period
    */
    function setInitialProbationPeriod(uint256 _initialProbationPeriod) external virtual onlyOperator {
        config.initialProbationPeriod = _initialProbationPeriod;
    }

    /* @notice sets the initial slashing rate
    * @dev restricted to the operator
    * @param _initialSlashingRate, the new value for the initial slashing rate
    */
    function setInitialSlashingRate(uint256 _initialSlashingRate) external virtual onlyOperator {
        require(_initialSlashingRate <= SCALE_FACTOR, "cannot exceed scale factor");
        config.initialSlashingRate = _initialSlashingRate;
    }

    /* @notice sets the lookback window. It will get updated at epoch end
    * @dev restricted to the operator
    * @param _lookbackWindow, the new value for the lookbackWindow
    */
    function setLookbackWindow(uint256 _lookbackWindow) external virtual onlyOperator {
        require(_lookbackWindow >= 1, "lookbackWindow cannot be 0");
        uint256 epochPeriod = autonity.getEpochPeriod();

        // utilize newDelta for comparison, so that if delta is also being changed in this epoch we take the new value
        require(epochPeriod > newDelta+_lookbackWindow-1,"epoch period needs to be greater than delta+lookbackWindow-1");
        newLookbackWindow = _lookbackWindow;
    }

    /* @notice sets delta. It will get updated at epoch end
    * @dev restricted to the operator
    * @param _delta, the new value for delta
    */
    function setDelta(uint256 _delta) external virtual onlyOperator {
        require(_delta >= 1, "delta cannot be 0");
        uint256 epochPeriod = autonity.getEpochPeriod();

        // utilize newLookbackWindow for comparison, so that if delta is also being changed in this epoch we take the new value
        require(epochPeriod > _delta+newLookbackWindow-1,"epoch period needs to be greater than delta+lookbackWindow-1");
        newDelta = _delta;
    }

    /**
    * @dev Modifier that checks if the caller is the autonity contract.
    */
    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to the Autonity Contract");
        _;
    }


    /**
    * @dev Modifier that checks if the caller is the operator.
    */
    modifier onlyOperator {
        require(operator == msg.sender, "restricted to operator");
        _;
    }
}
