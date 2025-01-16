// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.19;

import {AccessAutonity} from "./AccessAutonity.sol";

contract Latency is AccessAutonity {
    event Reported(address indexed reporter, uint256 length);
    uint256 public constant SCALE_THRESHOLD_FOR_CLUSTERING = 32;
    mapping(address => mapping(address => uint8)) public latency;
    address[] public committee;

    // to limit the unnecessary clustering reorg, we limit the report for 1 report per validator per epoch.
    uint256 public roundID;
    mapping(address => uint256) public reportedRound;

    constructor(address payable _autonity, address[] memory initialCommittee) AccessAutonity(_autonity) {
        committee = initialCommittee;
        roundID = 1;
    }

    function report(uint8[] memory _latency) external {
        require(_latency.length == committee.length, "Latency: invalid length");
        require(reportedRound[msg.sender] != roundID, "Already reported at current epoch");
        require(committee.length > SCALE_THRESHOLD_FOR_CLUSTERING, "Don't do clustering for small scale network");
        for (uint256 i = 0; i < _latency.length; i++) {
            latency[msg.sender][committee[i]] = _latency[i];
        }
        reportedRound[msg.sender] = roundID;
        emit Reported(msg.sender, _latency.length);
    }

    function setCommittee(address[] memory _committee) external onlyAutonity {
        committee = _committee;
        roundID +=1;
    }

    function read() external view returns (uint8[][] memory) {
        uint8[][] memory result = new uint8[][](committee.length);
        for (uint256 i = 0; i < committee.length; i++) {
            result[i] = new uint8[](committee.length);
            for (uint256 j = 0; j < committee.length; j++) {
                result[i][j] = latency[committee[i]][committee[j]];
            }
        }
        return result;
    }

    function getCommittee() external view returns (address[] memory) {
        return committee;
    }
}