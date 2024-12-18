// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;

import "../Oracle.sol";

contract OracleTest is Oracle {
    constructor(
        address[] memory _voters,
        address[] memory _nodeAddresses,
        address[] memory _treasuries,
        string[] memory _symbols,
        Config memory _config
    ) Oracle(
        _voters,
        _nodeAddresses,
        _treasuries,
        _symbols,
        _config
    ) {}

    function makeVoter(address _voter) public {
        voterInfo[_voter].isVoter = true;
    }

    function setVotersNow(address[] memory _voters) public {
        voters = _voters;
    }

    function priceLen() public view returns (uint256) {
        return prices.length;
    }

    function aggregateReports(uint _index) public {
        require(_index < symbols.length, "index out of range");
        _aggregateReports(_index);
    }

    function makeReportAvailable(address _voter) public {
        // require(_index < symbols.length, "index out of range");
        voterInfo[_voter].reportAvailable = true;
    }
}