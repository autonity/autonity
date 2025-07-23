// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.30;

import "./IAutonity.sol";

interface IOmissionAccountability {
    function finalize(bool _epochEnded) external returns (uint256);
    function setCommittee(IAutonity.CommitteeMember[] memory _committee, address[] memory _treasuries) external;
    function setEpochBlock(uint256 _epochBlock) external;
    function setOperator(address _operator) external;
    function getInactivityScore(address _validator) external view returns (uint256);
    function getScaleFactor() external pure returns (uint256);
    function distributeProposerRewards(uint256 _ntnReward) external payable;
    function getTotalEffort() external view returns (uint256);
    function getLookbackWindow() external view returns (uint256);
    function getDelta() external view returns (uint256);
}