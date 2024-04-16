// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

interface IStakeProxy {

    struct BondingApplied {
        address validator;
        uint256 liquidAmount;
    }

    struct UnbondingReleased {
        uint256 unbondingID;
        uint256 releasedAmount;
    }

    function bondingApplied(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) external;
    function unbondingApplied(uint256 _unbondingID, address _validator, bool _rejected) external;
    function unbondingReleased(uint256 _unbondingID, address _validator, uint256 _amount, bool _rejected) external;
    // function finalize(BondingApplied[] memory _appliedBonding, uint256[] memory _rejectedBonding, UnbondingReleased[] memory _releasedUnbonding) external;
    function distributeRewards(address[] memory _validators) external;
    function receiveAut() external payable;
}