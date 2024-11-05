// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;
import "../Autonity.sol";

contract LiquidStorage {
    mapping(address => uint256) internal balances;
    mapping(address => uint256) internal lockedBalances;

    mapping(address => mapping (address => uint256)) internal allowances;
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

    Autonity internal autonityContract; //not hardcoded for testing purposes
}