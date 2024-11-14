// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import {Autonity, ValidatorState} from "./Autonity.sol";
import {SLASHING_RATE_SCALE_FACTOR} from "./ProtocolConstants.sol";

contract Slasher {
    address private autonity;

    constructor(address _autonity){
        autonity = _autonity; // we could use msg.sender but it would make deploying a new upgraded slasher contract harder
    }

    /**
    * @dev modifies the passed validator structure to enable jailing
    * @param _val, the validator to be jailed
    * @param _jailtime, the jailing time to be assigned to the validator
    * @param _newJailedState, the validator state to be applied
    * @return the modified validator
    */
    function jail(
        Autonity.Validator memory _val,
        uint256 _jailtime,
        ValidatorState _newJailedState
    ) external virtual onlyAutonity returns (
        Autonity.Validator memory
    ){
        _jail(_val, _jailtime, _newJailedState);
        return _val;
    }

    function _jail(
        Autonity.Validator memory _val,
        uint256 _jailtime,
        ValidatorState _newJailedState
    ) internal virtual {
        _val.jailReleaseBlock = block.number + _jailtime;
        _val.state = _newJailedState;
    }

    /**
    * @dev modifies the passed validator structure to enable jailbounding (perpetual jail)
    * @param _val, the validator to be jailbound
    * @param _newJailboundState, the validator state to be applied
    * @return the modified validator
    */
    function jailbound(
        Autonity.Validator memory _val,
        ValidatorState _newJailboundState
    ) external virtual onlyAutonity returns (
        Autonity.Validator memory
    ){
        _jailbound(_val, _newJailboundState);
        return _val;
    }

    function _jailbound(
        Autonity.Validator memory _val,
        ValidatorState _newJailboundState
    ) internal virtual {
        _val.jailReleaseBlock = 0;
        _val.state = _newJailboundState;
    }

    /**
    * @dev applies slashing to the passed validator struct
    * @dev NOTE: 100% slash is not allowed and if attempted will cause a revert
    * @param _val, the validator to be slashed
    * @param _slashingRate, the rate for the slash
    * @return the modified validator
    * @return the slashing amount
    */
    function slash(
        Autonity.Validator memory _val,
        uint256 _slashingRate
    ) external virtual onlyAutonity returns (
        Autonity.Validator memory,
        uint256
    ){
        uint256 _slashingAmount = _slash(_val, _slashingRate);
        return (_val, _slashingAmount);
    }

    function _slash(
        Autonity.Validator memory _val,
        uint256 _slashingRate
    ) internal virtual returns (
        uint256 // slashingAmount
    ){
        require(_slashingRate < SLASHING_RATE_SCALE_FACTOR, "cannot slash 100% without jailbounding");
        uint256 _availableFunds = _val.bondedStake + _val.unbondingStake + _val.selfUnbondingStake;

        // here 0 <= slashingAmount < availableFunds
        uint256 _slashingAmount = (_slashingRate * _availableFunds) / SLASHING_RATE_SCALE_FACTOR;

        uint256 _remaining = _slashingAmount;
        // -------------------------------------------
        // Implementation of Penalty Absorbing Stake
        // -------------------------------------------
        // Self-unbonding stake gets slashed in priority.
        if (_val.selfUnbondingStake >= _remaining) {
            _val.selfUnbondingStake -= _remaining;
            _remaining = 0;
        } else {
            _remaining -= _val.selfUnbondingStake;
            _val.selfUnbondingStake = 0;
        }
        // Then self-bonded stake
        if (_remaining > 0) {
            if (_val.selfBondedStake >= _remaining) {
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
        return _slashingAmount;
    }

    /**
      * @dev slashes and jails the specified validator
      * @param _val, the validator to be slashed
      * @param _slashingRate, the rate to be used
      * @param _jailtime, the jailing time to be assigned to the validator
      * @param _newJailedState, the validator state to be applied for jailing
      * @param _newJailboundState, the validator state to be applied in case of 100% slashing
      * @return the modified validator
      * @return the amount slashed in NTN
      * @return a flag that signals if the validator has been permanently jailed
      */
    function slashAndJail(
        Autonity.Validator memory _val,
        uint256 _slashingRate,
        uint256 _jailtime,
        ValidatorState _newJailedState,
        ValidatorState _newJailboundState
    ) external virtual onlyAutonity returns (
        Autonity.Validator memory,  // slashedVal
        uint256,                    // slashingAmount
        bool                        // isJailbound
    ){

        // in case of >= 100% slash, slash all funds and jailbound validator
        if (_slashingRate >= SLASHING_RATE_SCALE_FACTOR) {
            uint256 _availableFunds = _val.bondedStake + _val.unbondingStake + _val.selfUnbondingStake;
            _val.bondedStake = 0;
            _val.selfBondedStake = 0;
            _val.selfUnbondingStake = 0;
            _val.unbondingStake = 0;
            _val.totalSlashed += _availableFunds;
            _jailbound(_val, _newJailboundState);
            return (_val, _availableFunds, true);
        }

        uint256 _slashingAmount = _slash(_val, _slashingRate);
        _jail(_val, _jailtime, _newJailedState);
        return (_val, _slashingAmount, false);
    }

    /**
      * @notice returns the scale factor used for slashing
      * @return slashing scale factor
      */
    function getSlashingScaleFactor() external view returns (uint256) {
        return SLASHING_RATE_SCALE_FACTOR;
    }

    modifier onlyAutonity {
        require(
            msg.sender == autonity,
            "Call restricted to the Autonity Contract");
        _;
    }
}
