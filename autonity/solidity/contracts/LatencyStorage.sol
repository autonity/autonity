// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract LatencyStorage {

    struct Latency {
        address receiver;
        uint256 latency;
    }

    mapping(address => mapping(address => uint256)) private latency;

    constructor() {}

    function updateLatency(Latency[] memory _latency) public {
        address _sender = msg.sender;
        for (uint256 i = 0; i < _latency.length; i++) {
            latency[_sender][_latency[i].receiver] = _latency[i].latency;
        }
    }

    function getLatency(address _sender, address _receiver) public view returns (uint256) {
        return latency[_sender][_receiver];
    }

    function getTotalLatency(address _sender, address _receiver) public view returns (uint256) {
        uint256 _latency = latency[_sender][_receiver];
        uint256 _reverseLatency = latency[_receiver][_sender];
        if (_latency > 0 && _reverseLatency > 0) {
            return _latency + _reverseLatency;
        }
        if (_latency > 0) {
            return _latency + _latency;
        }
        return _reverseLatency + _reverseLatency;
    }

    function getMultipleLatency(address _sender, address[] memory _receiver) public view returns (Latency[] memory) {
        Latency[] memory _latency = new Latency[](_receiver.length);
        for (uint256 i = 0; i < _latency.length; i++) {
            _latency[i] = Latency(_receiver[i], latency[_sender][_receiver[i]]);
        }
        return _latency;
    }

    function getMultipleTotalLatency(address _sender, address[] memory _receiver) public view returns (Latency[] memory) {
        Latency[] memory _latency = new Latency[](_receiver.length);
        for (uint256 i = 0; i < _latency.length; i++) {
            _latency[i] = Latency(_receiver[i], getTotalLatency(_sender, _receiver[i]));
        }
        return _latency;
    }

}