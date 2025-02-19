// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../ProtocolConstants.sol";

library StakingPoolMath {
    function computeRewardsFromFeeFactor(
        uint256 _lastFeeFactor,
        uint256 _feeFactor,
        uint256 _balance
    ) internal pure returns (uint256) {
        // assuming valid inputs, so `_lastFeeFactor >= _feeFactor`
        return (_lastFeeFactor - _feeFactor) * _balance / FEE_FACTOR_UNIT_RECIP;
    }

    function computeRewardFraction(
        uint256 _rewards,
        uint256 _share,
        uint256 _totalShare
    ) internal pure returns (uint256) {
        if (_share == 0) {
            return 0;
        }
        // assuming valid inputs, so `_share <= _totalShare`
        return (_rewards * _share) / _totalShare;
    }

    /**
     * @dev Calculates minted liquid amount from the amount of newton bonded.
     * As liquid is minted, in case of `_totalDelegation == 0`, we mint in 1:1 ratio.
     * @param _newtonBonded new bonded stake
     * @param _totalDelegation total delegated stake in existence or in a pool
     * @param _totalLiquid total liquid in supply or in a pool
     */
    function liquidFromStake(
        uint256 _newtonBonded,
        uint256 _totalDelegation,
        uint256 _totalLiquid
    ) internal pure returns (uint256) {
        if (_totalDelegation == 0) {
            return _newtonBonded;
        }
        return (_totalLiquid * _newtonBonded) / _totalDelegation;
    }

    /**
     * @dev Calculates amount of unbonding stake from the amount of liquid. As nothing is minted here,
     * we can assume `_liquidBurning <= _totalLiquid`.
     * @param _liquidBurning amount of liquid being burnt
     * @param _totalLiquid total liquid in supply or in a pool
     * @param _totalDelegation total delegated stake in existence or in a pool 
     */
    function unbondingStakeFromLiquid(
        uint256 _liquidBurning,
        uint256 _totalLiquid,
        uint256 _totalDelegation
    ) internal pure returns (uint256) {
        if (_liquidBurning == 0) {
            return 0;
        }
        // assuming valid inputs `_liquidBurning <= _totalLiquid`
        return (_liquidBurning * _totalDelegation) / _totalLiquid;
    }

    /**
     * @dev Calculates amount of unbonding share from unbonding stake. Unbonding share is minted,
     * so in case of `_totalUnbondingStake == 0`, we mint in 1:1 ratio.
     * @param _unbondingStake amount of stake under unbonding
     * @param _totalUnbondingStake total unbonding stake in existence or in a pool
     * @param _totalUnbondingShare total unbonding share in existence or in a pool
     */
    function unbondingShareFromStake(
        uint256 _unbondingStake,
        uint256 _totalUnbondingStake,
        uint256 _totalUnbondingShare
    ) internal pure returns (uint256) {
        if (_totalUnbondingStake == 0) {
            return _unbondingStake;
        }
        return (_unbondingStake * _totalUnbondingShare) / _totalUnbondingStake;
    }

    /**
     * @dev Calculates amount of unbonding share from requested self unbonding stake or liquid amount.
     * As unbonding share is not minted here, we can assume `_requestAmount <= _totalUnbondingAmount`.
     * @param _requestAmount amount of liquid (delegation) or newton (self-delegation) in unbonding request
     * @param _totalUnbondingAmount total self-unbonding stake or liquid in existence or in a pool
     * @param _totalUnbondingShare total self-unbonding share or unbonding share in existence or in a pool
     */
    function unbondingShareFromRequestAmount(
        uint256 _requestAmount,
        uint256 _totalUnbondingAmount,
        uint256 _totalUnbondingShare
    ) internal pure returns (uint256) {
        if (_requestAmount == 0) {
            return 0;
        }
        // assuming valid inputs `_requestAmount <= _totalUnbondingAmount`
        return (_requestAmount * _totalUnbondingShare) / _totalUnbondingAmount;
    }

    /**
     * @dev Calculates amount of stakes released from unbonding shares. As nothing is minted here,
     * we can assume `_unbondingShare <= _totalUnbondingStake`.
     * @param _unbondingShare unbonding share
     * @param _totalUnbondingShare total unbonding share in existence or in a pool
     * @param _totalUnbondingStake total unbonding stake in existence or in a pool
     */
    function releasedStakeFromShare(
        uint256 _unbondingShare,
        uint256 _totalUnbondingShare,
        uint256 _totalUnbondingStake
    ) internal pure returns (uint256) {
        if (_unbondingShare == 0) {
            return 0;
        }
        // assuming valid inputs `_unbondingShare <= _totalUnbondingStake`
        return (_unbondingShare * _totalUnbondingStake) / _totalUnbondingShare;
    }
}