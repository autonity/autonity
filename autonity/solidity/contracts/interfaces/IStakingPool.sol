// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {ILiquid} from "./ILiquid.sol";
import {ValidatorState} from "../Autonity.sol";

interface IStakingPool {
    /* Used for epoched staking */
    
    /* Staking Pool of the Validator for Future Processing */
    /**
     * @dev The fields are updated as new bonding or unboning requests appear.
     * At epoch end, validators are updated with the information from the `ValidatorPool`.
     */
    struct ValidatorBondingPool {
        uint256 selfBondingStake;
        uint256 delegatingStake;
        bool notActive;
    }

    struct ValidatorUnbondingPool {
        uint256 selfUnbondingStake;
        uint256 liquidBurning;
        uint256 selfUnbondingShare;
        uint256 unbondingShare;
    }

    /* Staking Pool of the Delegator for Future Processing */
    /**
     * @dev The fields are updated at epoch end.
     * On external calls, delegators take their share from the pool.
     */
    struct DelegatorBondingPool {
        uint256 liquidMinted;
        uint256 feeFactor;
        uint256 rewardsCollected;
    }

    struct DelegatorUnbondingPool {
        uint256 selfUnbondingShare;
        uint256 unbondingShare;
        uint256 releasedSelfStake;
        uint256 releasedStake;
        uint256 feeFactor;
        uint256 rewardsCollected;
    }

    struct PoolCollection {
        ValidatorBondingPool validatorBondingPool;
        ValidatorUnbondingPool validatorUnbondingPool;
        DelegatorBondingPool delegatorBondingPool;
        DelegatorUnbondingPool delegatorUnbondingPool;
    }
    
    struct BondingRequest {
        address payable delegator;
        address validator;
        uint256 amount;
        uint256 requestBlock;
        uint256 epochID;
        bool selfDelegation;
    }

    struct UnbondingRequest {
        address payable delegator;
        address validator;
        uint256 amount; // NTN for self-delegation, LNTN otherwise
        uint256 unbondingShare;
        uint256 requestBlock;
        uint256 epochID;
        bool unlocked;
        bool selfDelegation;
    }

    function setOperator(address _operator) external;

    function bond(
        address _validator,
        uint256 _amount,
        address payable _recipient,
        uint256 _epochID,
        bool _selfBond
    ) external returns (uint256);

    function unbond(
        address _validator,
        uint256 _amount,
        address payable _recipient,
        uint256 _epochID,
        bool _selfBond
    ) external returns (uint256);

    /* Protocol calls done in `autonity.finalize()` by autonity contract */
    function collectRewards(address[] memory _validators, ILiquid[] memory _liquidContracts) external;
    function applyBonding(uint256 _epochID) external;
    function applyUnbonding(uint256 _epochID) external;
    function releaseUnbondingStake(uint256 _epochID) external;

    /* Updates the information of the `_delegator` from the delegators pool. Both function do the same thing. */
    function updateDelegatorPool(address _delegator) external;
    function updateDelegatorPool(address _delegator, address _validator) external;

    function getBondingRequest(uint256 _id) external view returns (BondingRequest memory);
    function getUnbondingRequest(uint256 _id) external view returns (UnbondingRequest memory);
    function getBondingArrayLength() external view returns (uint256);
    function getUnbondingArrayLength() external view returns (uint256);
    function rejectedBonding() external view returns (uint256);
    function releasedStakes() external view returns (uint256);
    function isUnbondingReleased(uint256 _id) external view returns (bool);
    function getUnbondingShare(uint256 _id) external view returns (uint256);
    function calculateReleasedStake(address _delegator) external view returns (uint256);
    function calculateRejectedBonding(address _delegator, uint256 _epochID) external view returns (uint256);
    function calculateLiquidMinted(address _delegator, address _validator) external view returns (uint256);
    function calculateLiquidBurning(address _delegator, address _validator) external view returns (uint256);
    function calculateRewards(address _delegator, address _validator) external view returns (uint256);

    event BondingPoolRejected(address indexed validator, uint256 selfBondingAmount, uint256 delegatingAmount, ValidatorState state);
}