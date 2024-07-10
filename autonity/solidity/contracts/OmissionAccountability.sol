// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./Autonity.sol";
contract OmissionAccountability is IOmissionAccountability {
    // TODO(lorenzo) should this be in the config?
    uint256 public constant SCALE_FACTOR = 10_000;

    struct Config {
        uint256 negligibleThreshold;        // threshold to determine if a validator is an offender at the end of epoch.
        uint256 omissionLookBackWindow;
        uint256 activityProofRewardRate; //TODO(lorenzo) this is actually needed only in  autonity
        uint256 maxCommitteeSize; // TODO(lorenzo) this as welll
        uint256 pastPerformanceWeight;
        uint256 initialJailingPeriod;       // initial number of epoch an offender will be jailed for
        uint256 initialProbationPeriod;     // initial number of epoch an offender will be set under probation for
        uint256 initialSlashingRate;
    }

    // shadow copies of variables in Autonity.sol
    address[] private committee;
    uint256 private lastEpochBlock;

    mapping(uint256 => bool) public faultyProposers;            // marks height where proposer is faulty
    mapping(uint256 => address[]) public inactiveValidators;    // list of inactive validators for each height
    mapping(address => uint256) public inactivityCounter;       // counter of inactive blocks for each validator (considering lookback window)

    // net (total - quorum) proposer effort included in the activity proof
    uint256 public totalAccumulatedEffort;
    mapping(address => uint256) public proposerEffort;

    mapping(address => uint256) public pastEpochInactivityScores;
    mapping(address => uint256) public currentEpochInactivityScores;

    mapping(address => uint256) public probationPeriods;
    mapping(address => uint256) public repeatedOffences; // reset as soon as an entire probation period is completed without offences.

    // TODO(lorenzo) implement this
    //mapping(address => uint256) public activityPercentage;
    //mapping(address => uint256) public overallFaultyBlocks;

    Config public config;
    Autonity internal autonity; // for access control in setters function.

    // TODO(Lorenzo) Memory or storage as data location for committee?
    constructor(address payable _autonity, address[] memory _committee, Config memory _config) {
        autonity = Autonity(_autonity);
        config = _config;
        committee = _committee;
    }

    //TODO(lorenzo): update comments and interface. Restore at symbol in front of notice and param
    /**
    * notice called by the Autonity Contract at block finalization, it receives activity report.
    * param isProposerOmissionFaulty is true when the proposer provides invalid activity proof of current height.
    * param ids stores faulty proposer's ID when isProposerOmissionFaulty is true, otherwise it carries current height
    * activity proof, which are the signers of precommit of current height - delta.
    */
    function finalize(address[] memory _absentees, address _proposer, uint256 _proposerEffort, bool _isProposerOmissionFaulty, bool _epochEnded) external onlyAutonity {
        uint256 targetHeight = block.number - DELTA;

        if (_isProposerOmissionFaulty) {
            faultyProposers[targetHeight] = true;
            inactivityCounter[_proposer]++;
        }else{
            faultyProposers[targetHeight] = false;
            proposerEffort[_proposer] += _proposerEffort;
            totalAccumulatedEffort += _proposerEffort;

            inactiveValidators[targetHeight] = _absentees;
            //TODO(lorenzo) check for off by one error
            if(targetHeight > lastEpochBlock + config.omissionLookBackWindow) {
                // for each absent of target height, check the lookback window to see if he was online at some point
                for(uint256 i=0; i < _absentees.length; i++) {
                    bool confirmedAbsent = true;
                    uint256 initialLookBackWindow = config.omissionLookBackWindow;
                    for(uint256 j=targetHeight-1;j>=targetHeight-initialLookBackWindow;j--) {
                        if(faultyProposers[j]) {
                            // if we do not have data for a certain height, extend the window
                            initialLookBackWindow++;
                            continue;
                        }
                        if(j == lastEpochBlock){
                            // if we end up here it means that we extended the lookback window too much and arrive at the start at the epoch
                            // we do not have enough information, so let's just consider the validator as not absent
                            confirmedAbsent=false;
                            break;
                        }

                        // if the validator is active even just once in the lookback window, then we consider him as not absent
                        // TODO(lorenzo) maybe too much complexity here
                        bool found = false;
                        for (uint256 k=0; k<inactiveValidators[j].length;k++){
                            if(_absentees[i] == inactiveValidators[j][k]){
                                found = true;
                            }
                        }

                        // if the validator is not found in even only one of the inactive lists, it is not considered offline
                        if(!found){
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
        }

        if(_epochEnded){
            uint256 epochPeriod = autonity.getEpochPeriod(); //TODO(lorenzo) not sure if better to shadow copy it + deal with changes in it
            uint256 collusionDegree = 0;

            // compute aggregated scores + collusion degree
            for(uint256 i=0;i<committee.length;i++){
                //TODO(lorenzo) use some kind of SCALE FACTOR?
                uint256 inactivityScore = (inactivityCounter[committee[i]]*SCALE_FACTOR / epochPeriod);
                // TODO(lorenzo) properly scale the past performance weight to avoid underflow
                //uint256 aggregatedInactivityScore =  inactivityScore*(1-config.pastPerformanceWeight) + currentEpochInactivityScores[committee[i]] * config.pastPerformanceWeight;
                uint256 aggregatedInactivityScore = inactivityScore;
                pastEpochInactivityScores[committee[i]] = currentEpochInactivityScores[committee[i]];
                currentEpochInactivityScores[committee[i]] =  aggregatedInactivityScore;
                if(aggregatedInactivityScore > config.negligibleThreshold){
                    collusionDegree++;
                }
            }

            // reduce probation periods + dish out punishment
            for(uint256 i=0;i<committee.length;i++){
                Autonity.Validator memory _val = autonity.getValidator(committee[i]);

                // if the validator has already been slashed by accountability in this epoch,
                // do not punish him for omission too. It would be unfair since peer ignore msgs from jailed vals.
                if(_val.state == ValidatorState.jailed || _val.state == ValidatorState.jailbound){
                    continue;
                }

                // TODO(lorenzo) this way probation period decreases only if the val is part of the committee. Is it what we want?
                if(_val.state == ValidatorState.active && probationPeriods[committee[i]] > 0){
                    probationPeriods[committee[i]]--;
                    // if decreased to zero, then zero out also the offences counter
                    if(probationPeriods[committee[i]] == 0){
                       repeatedOffences[committee[i]] = 0;
                    }
                }

                // punish validator if his inactivity is greater than threshold
                if(currentEpochInactivityScores[committee[i]] > config.negligibleThreshold){
                    repeatedOffences[committee[i]]++;
                    uint256 offenceSquared = repeatedOffences[committee[i]]*repeatedOffences[committee[i]];
                    uint256 jailingPeriod = config.initialJailingPeriod * offenceSquared;
                    uint256 probationPeriod = config.initialProbationPeriod * offenceSquared;

                    _val.jailReleaseBlock = block.number + jailingPeriod;
                    _val.state = ValidatorState.jailed;

                    if(probationPeriods[committee[i]] > 0){
                        _slash(_val,config.initialSlashingRate*offenceSquared*collusionDegree);
                    }else{
                        // TODO(lorenzo) should be fine even if we are not slashing,but double check
                        autonity.updateValidatorAndTransferSlashedFunds(_val);
                    }

                    // whether slashed or not, update the probation period
                    probationPeriods[committee[i]] += probationPeriod;
                }
            }

        }
    }

    // TODO(Lorenzo) the code is copied from the accountability contract, would be good to make it a func.
    // OR defer the slashing to the accountability contract
    function _slash(Autonity.Validator memory _val, uint256 _slashingRate) internal{
        if(_slashingRate > SCALE_FACTOR) {
            _slashingRate = SCALE_FACTOR;
        }

        uint256 _availableFunds = _val.bondedStake + _val.unbondingStake + _val.selfUnbondingStake;
        uint256 _slashingAmount =  (_slashingRate * _availableFunds)/SCALE_FACTOR;

        // in case of 100% slash, we jailbound the validator
        if (_slashingAmount > 0 && _slashingAmount == _availableFunds) {
            _val.bondedStake = 0;
            _val.selfBondedStake = 0;
            _val.selfUnbondingStake = 0;
            _val.unbondingStake = 0;
            _val.totalSlashed += _slashingAmount;
            //_val.provableFaultCount += 1;
            _val.state = ValidatorState.jailbound;
            _val.jailReleaseBlock = 0;
            autonity.updateValidatorAndTransferSlashedFunds(_val);
            //emit SlashingEvent(_val.nodeAddress, _slashingAmount, 0, true, _event.id);
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
        //_val.provableFaultCount += 1;
        //_val.jailReleaseBlock = block.number + config.jailFactor * _val.provableFaultCount * epochPeriod;
        //_val.state = ValidatorState.jailed; // jailed validators can't participate in consensus

        autonity.updateValidatorAndTransferSlashedFunds(_val);

        //emit SlashingEvent(_val.nodeAddress, _slashingAmount, _val.jailReleaseBlock, false, _event.id);
    }

    function distributeProposerRewards() external payable onlyAutonity {
        uint256 _atnReward = msg.value;

        for(uint256 i=0; i < committee.length; i++) {
           if(proposerEffort[committee[i]] > 0){
               //TODO(lorenzo) send it to the treasury?
               // TODO(lorenzo) needs scaling?
               uint256 proposerReward = (proposerEffort[committee[i]] * _atnReward) / totalAccumulatedEffort;
               committee[i].call{value: proposerReward, gas: 2300}("");
           }
        }
    }

    function getInactivityScore(address _validator) external view returns (uint256) {
        return currentEpochInactivityScores[_validator];
    }

    function getScaleFactor() external pure returns (uint256) {
        return SCALE_FACTOR;
    }

    function getProposerRewardsRate() external view returns (uint256) {
        return config.activityProofRewardRate;
    }

    function setCommittee(address[] memory _committee) external onlyAutonity {
        committee = _committee;
    }

    function setLastEpochBlock(uint256 _lastEpochBlock) external onlyAutonity {
        lastEpochBlock = _lastEpochBlock;
    }

    /**
    * @dev Modifier that checks if the caller is the autonity contract.
    */
    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to the Autonity Contract");
        _;
    }
}
