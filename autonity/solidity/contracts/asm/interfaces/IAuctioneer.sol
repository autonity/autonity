// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

interface IAuctioneer {
    function paidInterest() external payable;
    function setOperator(address operator) external;
}