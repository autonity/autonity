// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IStakeProxy {

    function bondingApplied(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) external;
    function unbondingApplied(uint256 _unbondingID, address _validator, bool _rejected) external;
    function unbondingReleased(uint256 _unbondingID, address _validator, uint256 _amount, bool _rejected) external;
    function rewardsDistributed(address[] memory _validators, uint256[] memory _delegatedStake, uint256[] memory _liquidSupply) external;
    function receiveAut() external payable;

}