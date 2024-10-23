// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "./Autonity.sol";

contract Slasher {
    Autonity private autonity;

    uint256 private constant SLASHING_RATE_PRECISION = 10_000;

    constructor(address payable _autonity){
        autonity = Autonity(_autonity);
    }

    /**
    * @notice generic slashing function
    * @param _val, the validator to be slashed
    * @param _slashingRate, the rate to be used
    * @param _jailtime, the jailing time to be assigned to the validator
    * @param _newJailedState, the validator state to be applied for jailing
    * @param _newJailboundState, the validator state to be applied in case of 100% slashing
    * @return the slashing amount, the jail release block and whether the validator got jailbound or not
    */
    function _slashAtRate(
        Autonity.Validator memory _val,
        uint256 _slashingRate,
        uint256 _jailtime,
        ValidatorState _newJailedState,
        ValidatorState _newJailboundState
    ) internal virtual returns (
        uint256 slashingAmount,
        uint256 jailReleaseBlock,
        bool isJailbound
    ){
        if(_slashingRate > SLASHING_RATE_PRECISION) {
            _slashingRate = SLASHING_RATE_PRECISION;
        }

        uint256 availableFunds = _val.bondedStake + _val.unbondingStake + _val.selfUnbondingStake;
        slashingAmount =  (_slashingRate * availableFunds) / SLASHING_RATE_PRECISION;

        // in case of 100% slash, we jailbound the validator
        if (slashingAmount > 0 && slashingAmount == availableFunds) {
            isJailbound = true;
            jailReleaseBlock = 0;

            _val.bondedStake = 0;
            _val.selfBondedStake = 0;
            _val.selfUnbondingStake = 0;
            _val.unbondingStake = 0;
            _val.totalSlashed += slashingAmount;
            _val.state = _newJailboundState;
            _val.jailReleaseBlock = jailReleaseBlock;
            autonity.updateValidatorAndTransferSlashedFunds(_val);
            return;
        }
        uint256 remaining = slashingAmount;
        // -------------------------------------------
        // Implementation of Penalty Absorbing Stake
        // -------------------------------------------
        // Self-unbonding stake gets slashed in priority.
        if(_val.selfUnbondingStake >= remaining){
            _val.selfUnbondingStake -= remaining;
            remaining = 0;
        } else {
            remaining -= _val.selfUnbondingStake;
            _val.selfUnbondingStake = 0;
        }
        // Then self-bonded stake
        if (remaining > 0){
            if(_val.selfBondedStake >= remaining) {
                _val.selfBondedStake -= remaining;
                _val.bondedStake -= remaining;
                remaining = 0;
            } else {
                remaining -= _val.selfBondedStake;
                _val.bondedStake -= _val.selfBondedStake;
                _val.selfBondedStake = 0;
            }
        }
        // --------------------------------------------
        // Remaining stake to be slashed is split equally between the delegated
        // stake pool and the non-self unbonding stake pool.
        // As a reminder, the delegated stake pool is bondedStake - selfBondedStake.
        // if _remaining > 0 then bondedStake = delegated stake, because all selfBondedStake is slashed
        if (remaining > 0 && (_val.unbondingStake + _val.bondedStake > 0)) {
            // as we cannot store fraction here, we are taking floor for both unbondingSlash and delegatedSlash
            // In case both variable unbondingStake and bondedStake are positive, this modification
            // will ensure that no variable reaches 0 too fast where the other one is too big. In this case both variables
            // will reach 0 only when slashed 100%.
            // That means the fairness issue: https://github.com/autonity/autonity/issues/819 will only be triggered
            // if 100% stake is slashed
            uint256 unbondingSlash = (remaining * _val.unbondingStake) /
                (_val.unbondingStake + _val.bondedStake);
            uint256 delegatedSlash = (remaining * _val.bondedStake) /
                (_val.unbondingStake + _val.bondedStake);
            _val.unbondingStake -= unbondingSlash;
            _val.bondedStake -= delegatedSlash;
            remaining -= unbondingSlash + delegatedSlash;
        }

        // if positive amount remains
        slashingAmount -= remaining;
        _val.totalSlashed += slashingAmount;
        isJailbound = false;
        jailReleaseBlock = block.number + _jailtime;
        _val.jailReleaseBlock = jailReleaseBlock;
        _val.state = _newJailedState;

        autonity.updateValidatorAndTransferSlashedFunds(_val);
    }
}
