// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

interface IOmissionAccountability {
    function finalize(address[] memory _absentees, address _proposer, uint256 _proposerEffort, bool _isProposerOmissionFaulty, bool _epochEnded) external;
    function setCommittee(address[] memory _nodeAddresses, address[] memory _treasuries) external;
    function setLastEpochBlock(uint256 _lastEpochBlock) external;
    function setOperator(address _operator) external;
    function getInactivityScore(address _validator) external view returns (uint256);
    function getScaleFactor() external pure returns (uint256);
    function distributeProposerRewards(uint256 _ntnReward) external payable;
    function getLookbackWindow() external view virtual returns (uint256,bool);
    function getTotalEffort() external view virtual returns (uint256);

    event InactivitySlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound);
}