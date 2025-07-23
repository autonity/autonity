// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.30;

interface IAuctioneer {
    // stabilization functions
    function paidInterest() external payable;

    // autonity functions
    function setOperator(address operator) external;
    function setOracle(address oracle) external;
    function setStabilization(address stabilization) external;
}