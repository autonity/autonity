// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {ILiquid} from "./ILiquid.sol";
import {ValidatorState} from "../Autonity.sol";

interface IStakingPool {
    /* Used for epoched staking */
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
        bool released;
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
        ILiquid _liquidContract,
        uint256 _amount,
        address payable _recipient,
        uint256 _epochID,
        bool _selfBond
    ) external returns (uint256);

    function collectRewards(address[] memory _validators, ILiquid[] memory _liquidContracts) external;
    function applyBonding(uint256 _epochID) external;
    function applyUnbonding(uint256 _epochID) external;
    function releaseUnbondingStake(uint256 _epochID) external;

    function getBondingRequest(uint256 _id) external view returns (BondingRequest memory);
    function getUnbondingRequest(uint256 _id) external view returns (UnbondingRequest memory);
    function getBondingArrayLength() external view returns (uint256);
    function getUnbondingArrayLength() external view returns (uint256);
    function rejectedBonding() external view returns (uint256);
    function releasedStakes() external view returns (uint256);
    function isUnbondingReleased(uint256 _id) external view returns (bool);
    function getUnbondingShare(uint256 _id) external view returns (uint256);

    event BondingPoolRejected(address indexed validator, uint256 selfBondingAmount, uint256 delegatingAmount, ValidatorState state);
}