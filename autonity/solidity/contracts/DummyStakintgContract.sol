// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./Autonity.sol";
import "./interfaces/IStakeProxy.sol";

contract DummyStakintgContract is IStakeProxy {

    enum Fail {
        Never,
        AfterRewardsDitribution,
        AfterUnbondingApplied
    }

    Autonity public autonity;
    bool public revertStaking;
    Fail public revertStep;

    uint256 public bondingCost;
    uint256 public unbondingCost;

    uint256[] public requestedBondings;
    uint256[] public requestedUnbondings;

    struct BondingApplied {
        address validator;
        uint256 liquid;
        bool selfDelegation;
        bool rejected;
        bool applied;
    }

    mapping(uint256 => BondingApplied) public notifiedBondings;

    struct UnbondingApplied {
        address validator;
        bool rejected;
        bool applied;
    }

    mapping(uint256 => UnbondingApplied) public notifiedUnbonding;

    struct UnbondingReleased {
        uint256 amount;
        bool rejeced;
        bool applied;
    }

    mapping(uint256 => UnbondingReleased) public notifiedRelease;

    uint256 constant public Validator_Staked = 1;
    uint256 constant public Validator_Rewarded = 2;
    
    // mapping(uint256 => address[]) public notifiedRewardsDistribution;
    // uint256[] public epochIDs;
    mapping(uint256 => mapping(address => uint256)) public validatorsState;

    constructor(address payable _autonity) {
        autonity = Autonity(_autonity);
        revertStep = Fail.Never;
        revertStaking = false;
        bondingCost = autonity.maxBondAppliedGas() + autonity.maxRewardsDistributionGas();
        bondingCost *= autonity.stakingGasPrice();
        unbondingCost = autonity.maxUnbondAppliedGas() + autonity.maxUnbondReleasedGas() + autonity.maxRewardsDistributionGas();
        unbondingCost *= autonity.stakingGasPrice();
    }

    function revertStakingOperations() public {
        revertStaking = true;
    }

    function processStakingOperations() public {
        revertStaking = false;
    }

    function failAfterRewardsDistribution() public {
        revertStep = Fail.AfterRewardsDitribution;
    }

    function failAfterUnbondingApplied() public {
        revertStep = Fail.AfterUnbondingApplied;
    }

    function failNever() public {
        revertStep = Fail.Never;
        revertStaking = false;
    }

    function removeRequestedBondingIDs() public {
        delete requestedBondings;
    }

    function removeRequestedUnbondingIDs() public {
        delete requestedUnbondings;
    }

    function bond(address _validator, uint256 _amount) public payable returns (uint256) {
        uint256 _bondingID = autonity.bond{value: bondingCost}(_validator, _amount);
        requestedBondings.push(_bondingID);
        validatorsState[autonity.epochID()][_validator] = Validator_Staked;
        notifiedBondings[_bondingID] = BondingApplied(_validator, _amount, false, false, false);
        return _bondingID;
    }

    function unbond(address _validator, uint256 _amount) public payable returns (uint256) {
        uint256 _unbondingID = autonity.unbond{value: unbondingCost}(_validator, _amount);
        requestedUnbondings.push(_unbondingID);
        validatorsState[autonity.epochID()][_validator] = Validator_Staked;
        notifiedUnbonding[_unbondingID] = UnbondingApplied(_validator, false, false);
        notifiedRelease[_unbondingID] = UnbondingReleased(_amount, false, false);
        return _unbondingID;
    }

    function bondingApplied(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) external onlyAutonity {
        require(revertStaking == false, "bondingApplied reverts");
        notifiedBondings[_bondingID] = BondingApplied(_validator, _liquid, _selfDelegation, _rejected, true);
    }

    function unbondingApplied(uint256 _unbondingID, address _validator, bool _rejected) external onlyAutonity {
        require(revertStaking == false, "unbondingApplied reverts");
        notifiedUnbonding[_unbondingID] = UnbondingApplied(_validator, _rejected, true);

        if (revertStep == Fail.AfterUnbondingApplied) {
            revertStaking = true;
        }
    }

    function unbondingReleased(uint256 _unbondingID, uint256 _amount, bool _rejected) external onlyAutonity {
        require(revertStaking == false, "unbondingReleased reverts");
        notifiedRelease[_unbondingID] = UnbondingReleased(_amount, _rejected, true);
    }

    function rewardsDistributed(address[] memory _validators) external onlyAutonity {
        require(revertStaking == false, "rewardsDistributed reverts");

        mapping(address => uint256) storage _map = validatorsState[autonity.epochID()];
        for (uint256 i = 0; i < _validators.length; i++) {
            _map[_validators[i]] = Validator_Rewarded;
        }

        // in the next notification call, it will fail
        if (revertStep == Fail.AfterRewardsDitribution) {
            revertStaking = true;
        }
    }

    function receiveATN() external payable {}

    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to Autonity contract");
        _;
    }
}