// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import {Autonity} from "./Autonity.sol";

// Used for testing contract upgrade mechanism.
contract TestBase {
    string public Foo;
    constructor(string memory _foo){
        Foo = _foo;
    }
}

contract TestUpgraded is TestBase {
    string public Bar;
    constructor(string memory _bar, string memory _foo) TestBase(_foo) {
        Bar = _bar;
    }
    function FooBar(string memory _foo) public {
        Foo = _foo;
    }
}
