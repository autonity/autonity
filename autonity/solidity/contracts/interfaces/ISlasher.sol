// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "../Autonity.sol";

interface ISlasher {
    function slashAtRate(
        Autonity.Validator memory _val,
        uint256 _slashingRate,
        uint256 _jailtime,
        ValidatorState _newJailedState,
        ValidatorState _newJailboundState
    ) external returns (
        uint256 slashingAmount,
        uint256 jailReleaseBlock,
        bool isJailbound
    );
}
