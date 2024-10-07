// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;
import "./LiquidState.sol";

contract LiquidStateTest is LiquidState {
    constructor(
        address _validator,
        address payable _treasury,
        uint256 _commissionRate,
        string memory _index,
        address _liquidLogicAddress
    ) LiquidState(
        _validator,
        _treasury,
        _commissionRate,
        _index,
        _liquidLogicAddress
    ) {}

    function liquidLogicContract() public view returns (address) {
        return _liquidLogicContract();
    }
}