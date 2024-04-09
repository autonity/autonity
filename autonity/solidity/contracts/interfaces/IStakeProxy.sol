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

    // function bondingApplied(uint256 _bondingID, uint256 _liquid, bool _selfDelegation, bool _rejected) external;
    // function unbondingApplied(uint256 _unbondingID) external;
    // function unbondingReleased(uint256 _unbondingID, uint256 _amount) external;
    function finalize(BondingApplied[] memory _appliedBonding, uint256[] memory _rejectedBonding, UnbondingReleased[] memory _releasedUnbonding) external;
    function receiveAut() external payable;
}