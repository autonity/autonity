// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./Autonity.sol";
contract OmissionAccountability is IOmissionAccountability {

    struct Config {
        uint256 negligibleThreshold;        // threshold to determine if a validator is an offender at the end of epoch.
        uint256 omissionLookBackWindow;
        uint256 activityProofRewardRate;
        uint256 maxCommitteeSize;
        uint256 pastPerformanceWeight;
        uint256 initialJailingPeriod;       // initial number of epoch an offender will be jailed for
        uint256 initialProbationPeriod;     // initial number of epoch an offender will be set under probation for
        uint256 initialSlashingRate;
    }

    // shadow copies of variables in Autonity.sol
    address[] private committee;
    // TODO(Lorenzo) shadow copy also this one?
    //uint256 private lastEpochBlock;


    mapping(uint256 => bool) public faultyProposers;            // marks height where proposer is faulty
    mapping(uint256 => address[]) public inactiveValidators;    // list of inactive validators for each height
    mapping(address => uint256) public inactivityCounter;          // counter of inactive blocks for each validator (considering lookback window)

    // todo: to get the voting power from AC committee, or we save a shadow of it in this contract.
        uint256 public totalAccumulatedEffort;
        // TODO: count only > quorum ??
        // TODO: can we compute who was the proposer? or should we save it?
        mapping(address => uint256) public proverEfforts;

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

    // TODO(Lorenzo) add commiteee in deployment code. Also memory or storage as data locaiton?
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
    // TODO(Lorenzo) make delta a protocol param instead of passing it
    function finalize(address[] memory absentees, address proposer, uint256 proposerEffort, bool isProposerOmissionFaulty, uint256 lastEpochBlock, uint256 delta, bool epochEnded) external onlyAutonity {
        uint256 targetHeight = block.number - delta;

        if (isProposerOmissionFaulty) {
            faultyProposers[targetHeight] = true;
            inactivityCounter[proposer]++;
        }else{
            faultyProposers[targetHeight] = false;
            inactiveValidators[targetHeight] = absentees;
            //TODO(lorenzo) check for off by one error
            if(targetHeight > lastEpochBlock + config.omissionLookBackWindow) {
                // for each absent of target height, check the lookback window to see if he was online at some point
                for(uint256 i=0; i < absentees.length; i++) {
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
                            if(absentees[i] == inactiveValidators[j][k]){
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
                        inactivityCounter[absentees[i]]++;
                    }
                }
            }
            proverEfforts[proposer] += proposerEffort;
            totalAccumulatedEffort += proposerEffort;
        }

        if(epochEnded){
            /*
            // compute the inactivity score for the epoch which just ended, for all validators
            mapping(address => uint256) inactivityScore

            // TODO(lorenzo) actually better to compute the number of blocks for each validator "on the fly"
            for i:=0;i<committee.length;i++{

            }*/
        }
    }

    function setCommittee(address[] memory _committee) external onlyAutonity {
        committee = _committee;
    }

    /**
    * @dev Modifier that checks if the caller is the autonity contract.
    */
    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to the Autonity Contract");
        _;
    }
}
