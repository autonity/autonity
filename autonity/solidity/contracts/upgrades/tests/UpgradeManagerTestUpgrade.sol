// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.10;

import {UpgradeManager1} from "../1/UpgradeManager1.sol";

contract UpgradeManagerTestUpgrade is UpgradeManager1{
    constructor() UpgradeManager1(new bytes32[](0), new version[](0)) {}

    function getVersion(bytes32) external pure override returns (version memory) {
        revert("getVersion is bricked");
    }

}
