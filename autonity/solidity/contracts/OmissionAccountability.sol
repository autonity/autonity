// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./Autonity.sol";
contract OmissionAccountability {

    struct Config {
        uint256 omissionLoopBackWindow;
        uint256 activityProofRewardRate;
        uint256 maxCommitteeSize;
        uint256 pastPerformanceWeight;
        uint256 initialJailingPeriod;
        uint256 initialProbationPeriod;
        uint256 initialSlashingRate;
    }

    Config public config;
    Autonity internal autonity; // for access control in setters function.
    constructor(address payable _autonity, Config memory _config) {
        autonity = Autonity(_autonity);
        config = _config;
    }

    // todo: implement D4's protocol contract
}