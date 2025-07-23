// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity 0.8.30;

import "../interfaces/IAutonity.sol";
import {ReentrancyGuard} from "../ReentrancyGuard.sol";

contract LiquidStorage is ReentrancyGuard {
    mapping(address => uint256) internal balances;
    mapping(address => uint256) internal lockedBalances;

    mapping(address => mapping (address => uint256)) internal allowances;
    mapping(address => mapping (address => uint256)) internal unbondingAllowances;
    uint256 internal supply;

    mapping(address => uint256) internal atnRealisedFees;
    mapping(address => uint256) internal atnUnrealisedFeeFactors;
    uint256 internal atnLastUnrealisedFeeFactor;

    string internal liquidName;
    string internal liquidSymbol;

    address internal validator;
    address payable internal treasury;
    uint256 internal commissionRate;

    uint256 internal treasuryUnclaimedATN;

    IAutonity internal autonityContract; //not hardcoded for testing purposes
}