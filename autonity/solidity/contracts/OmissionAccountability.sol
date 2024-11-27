// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import {Autonity, ValidatorState} from "./Autonity.sol";
import {Precompiled} from "./lib/Precompiled.sol";
import {IOmissionAccountability} from "./interfaces/IOmissionAccountability.sol";
import {SLASHING_RATE_SCALE_FACTOR} from "./ProtocolConstants.sol";

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
        // number of blocks to wait before generating activity proof.
        // e.g. activity proof of block x is for block x - delta
        uint256 delta;
    }

    // shadow copies of variables in Autonity.sol, updated once a epoch
    Autonity.CommitteeMember[] internal committee;
    address[] internal treasuries; // treasuries of the committee members
    uint256 internal epochBlock;

    uint256 internal newLookbackWindow; // applied at epoch end
    uint256 internal newDelta;          // applied at epoch end
    address internal operator;

    mapping(uint256 => bool) public faultyProposers;                         // marks height where proposer is faulty
    uint256 public faultyProposersInWindow;                                  // number of faulty proposers in the current lookback window

    mapping(uint256 => mapping(address => bool)) public inactiveValidators;  // inactive validators for each height
    address[] internal absenteesLastHeight;                                    // absentees of previous height
    // last active block in the epoch for a validator. default value of -1 means they are not on a offline blocks streak. It gets reset at epoch rotation when the new committee is set.
    mapping(address => int256) public lastActive;

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
    event InactivityJailingEvent(address validator, uint256 releaseBlock);

    constructor(
        address payable _autonity,
        address _operator,
        address[] memory _treasuries,
        Config memory _config
    ) {
        // config sanity checks
        require(_config.inactivityThreshold <= SCALE_FACTOR, "inactivity threshold cannot exceed scale factor");
        require(_config.pastPerformanceWeight <= SCALE_FACTOR, "past performance weight cannot exceed scale factor");
        require(_config.initialSlashingRate <= SLASHING_RATE_SCALE_FACTOR, "initial slashing rate cannot exceed slashing rate scale factor");

        autonity = Autonity(_autonity);

        // fetch committee and make sure that delta is set correctly in the autonity contract
        Autonity.EpochInfo memory epochInfo = autonity.getEpochInfo();
        require(epochInfo.delta == _config.delta, "mismatch between delta stored in Autonity contract and the one in Omission contract");

        operator = _operator;
        config = _config;
        for (uint256 i = 0; i < epochInfo.committee.length; i++) {
            committee.push(epochInfo.committee[i]);
            lastActive[committee[i].addr] = - 1;
        }
        treasuries = _treasuries;

        newLookbackWindow = config.lookbackWindow;
        newDelta = config.delta;
    }

    /**
    * @notice called by the Autonity Contract at block finalization.
    * @param _epochEnded, true if this is the last block of the epoch
    * @return the current delta value
    */
    function finalize(bool _epochEnded) external virtual onlyAutonity returns (uint256) {
        // if we are at the first delta blocks of the epoch, the activity proof should be empty
        bool _mustBeEmpty = block.number <= epochBlock + config.delta;

        uint256 _committeeSlot;
        assembly{
            _committeeSlot := committee.slot
        }

        (bool _isProposerOmissionFaulty, uint256 _proposerEffort, address[] memory _absentees) = Precompiled.computeAbsentees(_mustBeEmpty, config.delta, _committeeSlot);

        // short-circuit function if the proof has to be empty
        if (_mustBeEmpty) {
            return config.delta;
        }

        uint256 _targetHeight = block.number - config.delta;

        // if we have already processed a full lookback window,
        // decrease the faulty proposer count for each faulty proposer at the tail of the window
        // window: (h - delta - lookback - faultyProposers, h - delta]
        if (_targetHeight >= epochBlock + config.lookbackWindow + faultyProposersInWindow) {
            while(faultyProposersInWindow > 0 && faultyProposers[_targetHeight - config.lookbackWindow - faultyProposersInWindow + 1]) {
                faultyProposersInWindow--;
            }
        }

        if (_isProposerOmissionFaulty) {
            faultyProposers[_targetHeight] = true;
            inactivityCounter[block.coinbase]++;
            faultyProposersInWindow++;
        } else {
            faultyProposers[_targetHeight] = false;
            proposerEffort[block.coinbase] += _proposerEffort;
            totalEffort += _proposerEffort;

            _recordAbsentees(_absentees, _targetHeight);
        }

        if (_epochEnded) {
            uint256 _collusionDegree = _computeInactivityScoresAndCollusionDegree();
            _punishInactiveValidators(_collusionDegree);

            // clean up
            for (uint256 i = 0; i < committee.length; i++) {
                inactivityCounter[committee[i].addr] = 0;
            }
            faultyProposersInWindow = 0;
            delete absenteesLastHeight;

            // store collusion degree in state. This is useful for slashed validators to verify their slashing rate
            epochCollusionDegree.push(_collusionDegree);

            // update lookback window and delta if changed
            config.lookbackWindow = newLookbackWindow;
            config.delta = newDelta;
        }
        return config.delta;
    }

    function _contains(address[] memory _absentees, address _account) internal pure returns (bool) {
        for (uint256 i = 0; i < _absentees.length; i++) {
            if (_absentees[i] == _account) {
                return true;
            }
        }
        return false;
    }


    function _recordAbsentees(address[] memory _absentees, uint256 _targetHeight) internal virtual {
        // check if absentees of past height came back online
        for (uint256 i = 0; i < absenteesLastHeight.length; i++) {
            if (!_contains(_absentees, absenteesLastHeight[i])) {
                // for all addresses who were inactive last height, if they are active now,
                // we can reset their lastActive counter
                lastActive[absenteesLastHeight[i]] = - 1;
            }
        }
        absenteesLastHeight = _absentees;

        // for each absent of target height, check the lookback window to see if he was online at some point
        // if online even once in the lookback window, consider him online for this block
        // NOTE: the current block is included in the window, (h - delta - lookback, h - delta]
        // if we include the faulty proposers extension, the formula window becomes:
        // (h - delta - lookback - faultyProposers, h - delta]
        for (uint256 i = 0; i < _absentees.length; i++) {
            inactiveValidators[_targetHeight][_absentees[i]] = true;
            // if this is the first time they are inactive, we can set their lastActive counter to the last block
            if (lastActive[_absentees[i]] == - 1) {
                lastActive[_absentees[i]] = int256(_targetHeight) - 1;
            }
            // check if we have enough blocks to correctly evaluate activity status
            if (_targetHeight < epochBlock + config.lookbackWindow + faultyProposersInWindow) {
                continue;
            }
            // validator is having a streak of offline blocks
            if (uint256(lastActive[_absentees[i]]) <= _targetHeight - (config.lookbackWindow + faultyProposersInWindow)) {
                // the validator was not active at some point in the lookback window
                inactivityCounter[_absentees[i]]++;
            }
        }
    }

    // returns collusion degree
    function _computeInactivityScoresAndCollusionDegree() internal virtual returns (uint256) {
        uint256 _epochPeriod = autonity.getCurrentEpochPeriod();
        uint256 _collusionDegree = 0;

        // first config.lookbackWindow-1 blocks of the epoch are accountable, but we do not have enough info to determine if a validator was offline/online
        // last delta blocks of the epoch are not accountable due to committee change
        uint256 _qualifiedBlocks = _epochPeriod - config.lookbackWindow + 1 - config.delta;

        // weight of current epoch performance
        uint256 _currentPerformanceWeight = SCALE_FACTOR - config.pastPerformanceWeight;

        // compute aggregated scores + collusion degree
        for (uint256 i = 0; i < committee.length; i++) {
            address _nodeAddress = committee[i].addr;

            // there is an edge case where inactivityCounter could be > qualifiedBlocks. However we cap it at qualifiedBlocks to prevent having > 100% inactivity score
            // this can happen for example if we have a network with a single validator, that is never including any activity proof,
            // thus always being considered a faulty proposer and getting his inactivityCounter increased even when we do not have lookback blocks yet
            if (inactivityCounter[_nodeAddress] > _qualifiedBlocks) {
                inactivityCounter[_nodeAddress] = _qualifiedBlocks;
            }

            /* the following formula is refactored to minimize precision loss, prioritizing multiplications over divisions
            *  A more intuitive but equivalent construction:
            *  aggregatedInactivityScore = (currentInactivityScore * currentPerformanceWeight + pastInactivityScore * pastPerformanceWeight) / SCALE_FACTOR
            *  with currentInactivityScore = (currentInactivityCounter * SCALE_FACTOR) / qualifiedBlocks
            */
            uint256 _aggregatedInactivityScore =
                (
                    inactivityCounter[_nodeAddress] * SCALE_FACTOR * _currentPerformanceWeight
                    + inactivityScores[_nodeAddress] * config.pastPerformanceWeight * _qualifiedBlocks
                )
                / (SCALE_FACTOR * _qualifiedBlocks);

            if (_aggregatedInactivityScore > config.inactivityThreshold) {
                _collusionDegree++;
            }
            inactivityScores[_nodeAddress] = _aggregatedInactivityScore;
        }
        return _collusionDegree;
    }

    function _punishInactiveValidators(uint256 _collusionDegree) internal virtual {
        // reduce probation periods + dish out punishment
        for (uint256 i = 0; i < committee.length; i++) {
            address _nodeAddress = committee[i].addr;

            // if the validator has already been slashed by accountability in this epoch,
            // do not punish him for omission too. It would be unfair since peer ignore msgs from jailed vals.
            // However, do not decrease his probation since he was not fully honest
            // NOTE: validator already jailed by accountability are nonetheless taken into account into the collusion degree of omission
            ValidatorState _state = autonity.getValidatorState(_nodeAddress);
            if (_state == ValidatorState.jailed || _state == ValidatorState.jailbound) {
                continue;
            }

            // here validator is either active or has been paused in the current epoch (but still participated to consensus)

            if (inactivityScores[_nodeAddress] <= config.inactivityThreshold) {
                // NOTE: probation period of a validator gets decreased only if he is part of the committee
                if (probationPeriods[_nodeAddress] > 0) {
                    probationPeriods[_nodeAddress]--;
                    // if decreased to zero, then zero out also the offences counter
                    if (probationPeriods[_nodeAddress] == 0) {
                        repeatedOffences[_nodeAddress] = 0;
                    }
                }
            } else {
                // punish validator if his inactivity is greater than threshold
                repeatedOffences[_nodeAddress]++;
                uint256 _offenceSquared = repeatedOffences[_nodeAddress] * repeatedOffences[_nodeAddress];
                uint256 _jailingPeriod = config.initialJailingPeriod * _offenceSquared;
                uint256 _probationPeriod = config.initialProbationPeriod * _offenceSquared;

                // if already on probation, slash and jail
                if (probationPeriods[_nodeAddress] > 0) {
                    uint256 _slashingRate = config.initialSlashingRate * _offenceSquared * _collusionDegree;
                    uint256 _slashingAmount;
                    uint256 _jailReleaseBlock;
                    bool _isJailbound;
                    (_slashingAmount, _jailReleaseBlock, _isJailbound) = autonity.slashAndJail(
                        _nodeAddress,
                        _slashingRate,
                        _jailingPeriod,
                        ValidatorState.jailedForInactivity,
                        ValidatorState.jailboundForInactivity
                    );
                    emit InactivitySlashingEvent(_nodeAddress, _slashingAmount, _jailReleaseBlock, _isJailbound);
                } else {
                    // if not, only jail
                    uint256 _jailReleaseBlock = autonity.jail(_nodeAddress, _jailingPeriod, ValidatorState.jailedForInactivity);
                    emit InactivityJailingEvent(_nodeAddress, _jailReleaseBlock);
                }

                // whether slashed or not, update the probation period (cumulatively)
                probationPeriods[_nodeAddress] += _probationPeriod;
            }
        }
    }

    /*
    * @notice called by the Autonity contract at epoch finalization, to redistribute the proposer rewards based on the effort
    * @param _ntnRewards, amount of NTN reserved for proposer rewards
    */
    function distributeProposerRewards(uint256 _ntnReward) external payable virtual onlyAutonity {
        uint256 _atnReward = address(this).balance;

        for (uint256 i = 0; i < committee.length; i++) {
            address _nodeAddress = committee[i].addr;
            if (proposerEffort[_nodeAddress] > 0) {
                uint256 _atnProposerReward = (proposerEffort[_nodeAddress] * _atnReward) / totalEffort;
                uint256 _ntnProposerReward = (proposerEffort[_nodeAddress] * _ntnReward) / totalEffort;

                if (_atnProposerReward > 0) {
                    // if for some reasons, funds can't be transferred to the treasury (sneaky contract)
                    (bool _ok,) = treasuries[i].call{value: _atnProposerReward, gas: 2300}("");
                    // well, too bad, it goes to the autonity global treasury.
                    if (!_ok) {
                        autonity.getTreasuryAccount().call{value: _atnProposerReward}("");
                    }
                }

                if (_ntnProposerReward > 0) {
                    autonity.autobond(_nodeAddress, _ntnProposerReward, 0);
                }

                // reset after usage
                proposerEffort[_nodeAddress] = 0;
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
    * @notice gets the delta used to determine how many block to wait before generating the activity proof. If the delta will change at epoch end,
    * the new value will be returned
    * @return the delta number of blocks to wait before generating the activity proof
    */
    function getDelta() external view virtual returns (uint256) {
        return newDelta;
    }

    /*
    * @notice retrieves the lookback window value and whether an update of it is in progress. If the lookback window will change at epoch end,
    * the new value will be returned
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
    * @notice get the absentees of last height
    * @return the absentees of last height
    */
    function getAbsenteesLastHeight() external view virtual returns (address[] memory){
        return absenteesLastHeight;
    }

    /*
    * @notice sets committee node addresses and treasuries
    * @dev restricted to the Autonity contract. It is used to mirror this information in the omission contract when the autonity contract changes
    * @param _committee, committee members
    * @param _treasuries, treasuries of the new committee
    */
    function setCommittee(Autonity.CommitteeMember[] memory _committee, address[] memory _treasuries) external virtual onlyAutonity {
        delete committee;
        for (uint256 i = 0; i < _committee.length; i++) {
            committee.push(_committee[i]);
            lastActive[committee[i].addr] = - 1; // clean up last active block for all new committee members
        }
        treasuries = _treasuries;
    }

    /* @notice sets the current epoch block in the omission contract
    * @dev restricted to the Autonity contract. It is used to mirror this information when it is updated at epoch finalize.
    * @param _epochBlock, epoch block of the current epoch
    */
    function setEpochBlock(uint256 _epochBlock) external virtual onlyAutonity {
        epochBlock = _epochBlock;
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
        require(_inactivityThreshold >= config.pastPerformanceWeight, "inactivityThreshold needs to be greater or equal to pastPerformanceWeight");
        config.inactivityThreshold = _inactivityThreshold;
    }

    /* @notice sets the past performance weight
    * @dev restricted to the operator
    * @param _pastPerformanceWeight, the new value for the past performance weight
    */
    function setPastPerformanceWeight(uint256 _pastPerformanceWeight) external virtual onlyOperator {
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
        require(_initialSlashingRate <= SLASHING_RATE_SCALE_FACTOR, "cannot exceed slashing rate scale factor");
        config.initialSlashingRate = _initialSlashingRate;
    }

    /* @notice sets the lookback window. It will get updated at epoch end
    * @dev restricted to the operator
    * @param _lookbackWindow, the new value for the lookbackWindow
    */
    function setLookbackWindow(uint256 _lookbackWindow) external virtual onlyOperator {
        require(_lookbackWindow >= 1, "lookbackWindow cannot be 0");
        uint256 _epochPeriod = autonity.getEpochPeriod();

        // utilize newDelta for comparison, so that if delta is also being changed in this epoch we take the new value
        require(_epochPeriod > newDelta + _lookbackWindow - 1, "epoch period needs to be greater than delta+lookbackWindow-1");
        newLookbackWindow = _lookbackWindow;
    }

    /* @notice sets delta. It will get updated at epoch end
    * @dev restricted to the operator
    * @param _delta, the new value for delta
    */
    function setDelta(uint256 _delta) external virtual onlyOperator {
        require(_delta >= 2, "delta needs to be at least 2"); // cannot be 1 due to optimistic block building
        uint256 _epochPeriod = autonity.getEpochPeriod();

        // utilize newLookbackWindow for comparison, so that if delta is also being changed in this epoch we take the new value
        require(_epochPeriod > _delta + newLookbackWindow - 1, "epoch period needs to be greater than delta+lookbackWindow-1");
        newDelta = _delta;
    }

    /**
    * @dev Modifier that checks if the caller is the autonity contract.
    */
    modifier onlyAutonity {
        require(msg.sender == address(autonity), "function restricted to the Autonity Contract");
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
