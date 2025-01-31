// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import "../asm/interfaces/IAuctioneer.sol";

contract AuctioneerMock is IAuctioneer {
    address public operator;

    function paidInterest() external payable override {
        // do nothing
    }

    function setOperator(address _operator) external override {
        operator = _operator;
    }
}