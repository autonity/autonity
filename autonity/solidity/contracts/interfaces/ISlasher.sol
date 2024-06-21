// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "../Autonity.sol";

interface ISlasher {

    function jail(
        Autonity.Validator memory _val,
        uint256 _jailtime,
        ValidatorState _newJailedState
    ) external returns (
        Autonity.Validator memory
    );

    function jailbound(
        Autonity.Validator memory _val,
        ValidatorState _newJailboundState
    ) external returns (
        Autonity.Validator memory
    );

    function slash(
        Autonity.Validator memory _val,
        uint256 _slashingRate
    ) external returns (
        Autonity.Validator memory,
        uint256 // slashingAmount
    );

    function slashAndJail(
        Autonity.Validator memory _val,
        uint256 _slashingRate,
        uint256 _jailtime,
        ValidatorState _newJailedState,
        ValidatorState _newJailboundState
    ) external returns (
        Autonity.Validator memory,  // slashedVal
        uint256,                    // slashingAmount
        bool                        // isJailbound
    );
}
