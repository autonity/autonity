// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;

import "../Oracle.sol";

// this contract mocks the autonity contract for use in Oracle tests
contract OracleAutonityMockTest {
    uint256 public epochPeriod;
    Oracle public oracle;

    constructor(uint256 _epochPeriod) {
        epochPeriod = _epochPeriod;
    }

    function setOracle(address payable _oracle) public {
        oracle = Oracle(_oracle);
    }

    function finalize() public {
        bool newRound = oracle.finalize();

        if (newRound) {
            oracle.updateVoters();
        }
    }

    function getEpochPeriod() public view returns (uint256) {
        return epochPeriod;
    }

    function getCurrentEpochPeriod() public view returns (uint256) {
        return epochPeriod;
    }

    function autobond(address _validator, uint256 _selfBonded, uint256 _bondedStake) public {
        // do nothing
    }

    function slash(address _validator, uint256 _rate) public {
        // do nothing
    }

    function setVoters(address[] memory _newVoters, address[] memory _treasury, address[] memory _validator) public {
        oracle.setVoters(_newVoters, _treasury, _validator);
    }
}