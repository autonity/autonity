// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;

// a dummy bindings contract. It imports all other contracts so the bindings can be generated for all of them, since abigen works on only one .sol file.
// I also considered to simply call abigen for every contract separately, but it compiles contract and it's dependencies,
// so in order to have separate Go bindings file for each, we would need to manually provide a list of exclusion (abigen doesn't
// let you name a contract you want to generate, only a list of excluded types).

import "./Autonity.sol";
import "./Accountability.sol";
import "./Liquid.sol";
import "./Oracle.sol";
import "./asm/ACU.sol";
import "./asm/Stabilization.sol";
import "./asm/SupplyControl.sol";
import "./AutonityUpgradeTest.sol";
import "./Tests.sol";
import "./UpgradeManager.sol";
