// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "./IAutonity.sol";

interface ISlasher {

    function jail(
        IAutonity.Validator memory _val,
        uint256 _blockNumber,
        uint256 _jailtime,
        IAutonity.ValidatorState _newJailedState
    ) external returns (
        IAutonity.Validator memory
    );

    function jailbound(
        IAutonity.Validator memory _val,
        IAutonity.ValidatorState _newJailboundState
    ) external returns (
        IAutonity.Validator memory
    );

    function slash(
        IAutonity.Validator memory _val,
        uint256 _slashingRate
    ) external returns (
        IAutonity.Validator memory,
        uint256 // slashingAmount
    );

    function slashAndJail(
        IAutonity.Validator memory _val,
        uint256 _slashingRate,
        uint256 _blockNumber,
        uint256 _jailtime,
        IAutonity.ValidatorState _newJailedState,
        IAutonity.ValidatorState _newJailboundState
    ) external returns (
        IAutonity.Validator memory,  // slashedVal
        uint256,                    // slashingAmount
        bool                        // isJailbound
    );
}
