// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import {Autonity} from "./Autonity.sol";
import {Precompiled} from "./Precompiled.sol";
import {Slasher} from "./Slasher.sol";
import {IOmissionAccountability} from "./interfaces/IOmissionAccountability.sol";

contract OmissionAccountability is IOmissionAccountability, Slasher {
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
    Autonity.CommitteeMember[] private committee;
    address[] private treasuries; // treasuries of the committee members
    uint256 private epochBlock;

    uint256 private newLookbackWindow; // applied at epoch end
    uint256 private newDelta;          // applied at epoch end
    address private operator;

    mapping(uint256 => bool) public faultyProposers;                         // marks height where proposer is faulty
    uint256 public faultyProposersInWindow;                                  // number of faulty proposers in the current lookback window

    mapping(uint256 => mapping(address => bool)) public inactiveValidators;  // inactive validators for each height
    mapping(address => uint256) public lastActive; // counter of active blocks for each validator. It is reset at the end of the epoch. A default value of 0 means they have been active in the last window.
    address[] public absenteesLastHeight;

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

    constructor(
        address payable _autonity,
        address _operator,
        address[] memory _treasuries,
        Config memory _config
    ) Slasher(_autonity) {
        autonity = Autonity(_autonity);

        // fetch committee and make sure that delta is set correctly in the autonity contract
        (Autonity.CommitteeMember[] memory committee,,,,uint256 delta) = autonity.getEpochInfo();
        require(delta == config.delta, "mismatch between delta stored in Autonity contract and the one in Omission contract");

        operator = _operator;
        config = _config;
        for (uint256 i = 0; i < committee.length; i++) {
            committee.push(committee[i]);
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

        uint256[1] memory _committeeSlot; // declare it as array to easily access from assembly
        assembly{
            mstore(_committeeSlot, committee.slot)
        }

        (bool _isProposerOmissionFaulty, uint256 _proposerEffort, address[] memory _absentees) = Precompiled.computeAbsentees(_mustBeEmpty, config.delta, _committeeSlot[0]);

        // short-circuit function if the proof has to be empty
        if (_mustBeEmpty) {
            return config.delta;
        }

        uint256 targetHeight = block.number - config.delta;

        if (_isProposerOmissionFaulty) {
            faultyProposers[targetHeight] = true;
            inactivityCounter[block.coinbase]++;
            faultyProposersInWindow++;
        } else {
            faultyProposers[targetHeight] = false;
            proposerEffort[block.coinbase] += _proposerEffort;
            totalEffort += _proposerEffort;
            lastActive[block.coinbase] = block.number; // why?

            _recordAbsentees(_absentees, targetHeight);
        }

        if (
            (targetHeight >= epochBlock + config.lookbackWindow) &&
            (faultyProposers[targetHeight - config.lookbackWindow])
        ) {
            faultyProposersInWindow--;
        }


        if (_epochEnded) {
            uint256 collusionDegree = _computeInactivityScoresAndCollusionDegree();
            _punishInactiveValidators(collusionDegree);

            // reset inactivity counters
            for (uint256 i = 0; i < committee.length; i++) {
                inactivityCounter[committee[i].addr] = 0;
                lastActive[committee[i].addr] = 0;
            }

            // store collusion degree in state. This is useful for slashed validators to verify their slashing rate
            epochCollusionDegree.push(collusionDegree);

            // update lookback window and delta if changed
            config.lookbackWindow = newLookbackWindow;
            config.delta = newDelta;
        }
        return config.delta;
    }

    function _contains(address[] memory _absentees, address _account) internal pure returns (bool) {
        for (uint i = 0; i < _absentees.length; i++) {
            if (_absentees[i] == _account) {
                return true;
            }
        }
        return false;
    }


    function _recordAbsentees(address[] memory _absentees, uint256 targetHeight) internal virtual {
        for (uint i = 0; i < absenteesLastHeight.length; i++) {
            if (!_contains(_absentees, absenteesLastHeight[i])) {
                // for all addresses who were inactive last height, if they are active now,
                // we can reset their lastActive counter
                lastActive[absenteesLastHeight[i]] = 0;
            }
        }

        if (targetHeight < epochBlock + config.lookbackWindow) {
            return;
        }

        // for each absent of target height, check the lookback window to see if he was online at some point
        // if online even once in the lookback window, consider him online for this block
        // NOTE: the current block is included in the window, (h - delta - lookback, h - delta]
        for (uint256 i = 0; i < _absentees.length; i++) {
            inactiveValidators[targetHeight][_absentees[i]] = true;
            // if this is the first time they are inactive, we can set their lastActive counter to the last block
            if (lastActive[_absentees[i]] == 0) {
                lastActive[_absentees[i]] = targetHeight - 1;
                continue;
            }
            if (lastActive[_absentees[i]] > targetHeight - (config.lookbackWindow + faultyProposersInWindow)) {
                // the validator was not active at some point in the lookback window
                inactivityCounter[_absentees[i]]++;
            }
        }
    }

    // returns collusion degree
    function _computeInactivityScoresAndCollusionDegree() internal virtual returns (uint256) {
        uint256 epochPeriod = autonity.getCurrentEpochPeriod();
        uint256 collusionDegree = 0;

        // first config.lookbackWindow-1 blocks of the epoch are accountable, but we do not have enough info to determine if a validator was offline/online
        // last delta blocks of the epoch are not accountable due to committee change
        uint256 qualifiedBlocks = epochPeriod - config.lookbackWindow + 1 - config.delta;

        // weight of current epoch performance
        uint256 currentPerformanceWeight = SCALE_FACTOR - config.pastPerformanceWeight;

        // compute aggregated scores + collusion degree
        for (uint256 i = 0; i < committee.length; i++) {
            address nodeAddress = committee[i].addr;

            // there is an edge case where inactivityCounter could be > qualifiedBlocks. However we cap it at qualifiedBlocks to prevent having > 100% inactivity score
            // this can happen for example if we have a network with a single validator, that is never including any activity proof,
            // thus always being considered a faulty proposer and getting his inactivityCounter increased even when we do not have lookback blocks yet
            if (inactivityCounter[nodeAddress] > qualifiedBlocks) {
                inactivityCounter[nodeAddress] = qualifiedBlocks;
            }

            /* the following formula is refactored to minimize precision loss, prioritizing multiplications over divisions
            *  A more intuitive but equivalent construction:
            *  aggregatedInactivityScore = (currentInactivityScore * currentPerformanceWeight + pastInactivityScore * pastPerformanceWeight) / SCALE_FACTOR
            *  with currentInactivityScore = (currentInactivityCounter * SCALE_FACTOR) / qualifiedBlocks
            */
            uint256 aggregatedInactivityScore =
                (
                    inactivityCounter[nodeAddress] * SCALE_FACTOR * currentPerformanceWeight
                    + inactivityScores[nodeAddress] * config.pastPerformanceWeight * qualifiedBlocks
                )
                / (SCALE_FACTOR * qualifiedBlocks);

            if (aggregatedInactivityScore > config.inactivityThreshold) {
                collusionDegree++;
            }
            inactivityScores[nodeAddress] = aggregatedInactivityScore;
        }
        return collusionDegree;
    }

    function _punishInactiveValidators(uint256 collusionDegree) internal virtual {
        // reduce probation periods + dish out punishment
        for (uint256 i = 0; i < committee.length; i++) {
            address nodeAddress = committee[i].addr;
            Autonity.Validator memory _val = autonity.getValidator(nodeAddress);

            // if the validator has already been slashed by accountability in this epoch,
            // do not punish him for omission too. It would be unfair since peer ignore msgs from jailed vals.
            // However, do not decrease his probation since he was not fully honest
            // NOTE: validator already jailed by accountability are nonetheless taken into account into the collusion degree of omission
            if (_val.state == ValidatorState.jailed || _val.state == ValidatorState.jailbound) {
                continue;
            }

            // here validator is either active or has been paused in the current epoch (but still participated to consensus)

            if (inactivityScores[nodeAddress] <= config.inactivityThreshold) {
                // NOTE: probation period of a validator gets decreased only if he is part of the committee
                if (probationPeriods[nodeAddress] > 0) {
                    probationPeriods[nodeAddress]--;
                    // if decreased to zero, then zero out also the offences counter
                    if (probationPeriods[nodeAddress] == 0) {
                        repeatedOffences[nodeAddress] = 0;
                    }
                }
            } else {
                // punish validator if his inactivity is greater than threshold
                repeatedOffences[nodeAddress]++;
                uint256 offenceSquared = repeatedOffences[nodeAddress] * repeatedOffences[nodeAddress];
                uint256 jailingPeriod = config.initialJailingPeriod * offenceSquared;
                uint256 probationPeriod = config.initialProbationPeriod * offenceSquared;

                // if already on probation, slash and jail
                if (probationPeriods[nodeAddress] > 0) {
                    uint256 slashingRate = config.initialSlashingRate * offenceSquared * collusionDegree;
                    uint256 slashingAmount;
                    uint256 jailReleaseBlock;
                    bool isJailbound;
                    (slashingAmount, jailReleaseBlock, isJailbound) = _slashAtRate(_val, slashingRate, jailingPeriod, ValidatorState.jailedForInactivity, ValidatorState.jailboundForInactivity);
                    emit InactivitySlashingEvent(_val.nodeAddress, slashingAmount, jailReleaseBlock, isJailbound);
                } else {
                    // if not, only jail
                    _val.jailReleaseBlock = block.number + jailingPeriod;
                    _val.state = ValidatorState.jailedForInactivity;
                    autonity.updateValidatorAndTransferSlashedFunds(_val);
                }

                // whether slashed or not, update the probation period (cumulatively)
                probationPeriods[nodeAddress] += probationPeriod;
            }
        }
    }

    /*
    * @notice called by the Autonity contract at epoch finalization, to redistribute the proposer rewards based on the effort
    * @param _ntnRewards, amount of NTN reserved for proposer rewards
    */
    function distributeProposerRewards(uint256 _ntnReward) external payable virtual onlyAutonity {
        uint256 atnReward = address(this).balance;

        for (uint256 i = 0; i < committee.length; i++) {
            address nodeAddress = committee[i].addr;
            if (proposerEffort[nodeAddress] > 0) {
                uint256 atnProposerReward = (proposerEffort[nodeAddress] * atnReward) / totalEffort;
                uint256 ntnProposerReward = (proposerEffort[nodeAddress] * _ntnReward) / totalEffort;

                if (atnProposerReward > 0) {
                    // if for some reasons, funds can't be transferred to the treasury (sneaky contract)
                    (bool ok,) = treasuries[i].call{value: atnProposerReward, gas: 2300}("");
                    // well, too bad, it goes to the autonity global treasury.
                    if (!ok) {
                        autonity.getTreasuryAccount().call{value: atnProposerReward}("");
                    }
                }

                if (ntnProposerReward > 0) {
                    autonity.transfer(treasuries[i], ntnProposerReward);
                }

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
    * @notice sets committee node addresses and treasuries
    * @dev restricted to the Autonity contract. It is used to mirror this information in the omission contract when the autonity contract changes
    * @param _committee, committee members
    * @param _treasuries, treasuries of the new committee
    */
    function setCommittee(Autonity.CommitteeMember[] memory _committee, address[] memory _treasuries) external virtual onlyAutonity {
        delete committee;
        for (uint256 i = 0; i < _committee.length; i++) {
            committee.push(_committee[i]);
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
        require(epochPeriod > newDelta + _lookbackWindow - 1, "epoch period needs to be greater than delta+lookbackWindow-1");
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
        require(epochPeriod > _delta + newLookbackWindow - 1, "epoch period needs to be greater than delta+lookbackWindow-1");
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
