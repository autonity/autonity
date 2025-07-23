// SPDX-License-Identifier: UNLICENSED
pragma solidity 0.8.30;

// no fallback method. cannot receive gas fee
contract DummyMaliciousContract {

    receive() external payable {
        uint iter = 0;
        while (true) {
            iter++;
        }
    }
}
