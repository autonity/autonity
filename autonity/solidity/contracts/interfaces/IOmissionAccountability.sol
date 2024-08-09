// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

interface IOmissionAccountability {
    /** //TODO(lorenzo) restore documentation tags (at symbol before notice and param)
    * notice called by the Autonity Contract at block finalization, it receives activity report.
    * param isProposerOmissionFaulty is true when the proposer provides invalid activity proof of current height.
    * param ids stores faulty proposer's ID when isProposerOmissionFaulty is true, otherwise it carries current height
    * activity proof which is the signers of precommit of current height - dela.
    * param epochEnded signals whether we are finalizing the epoch
    */
    function finalize(address[] memory _absentees, address _proposer, uint256 _proposerEffort, bool _isProposerOmissionFaulty, bool _epochEnded) external;

    function setCommittee(address[] memory _nodeAddresses, address[] memory _treasuries) external;
    function setLastEpochBlock(uint256 _lastEpochBlock) external;
    function setOperator(address _operator) external;
    function getInactivityScore(address _validator) external view returns (uint256);
    function getScaleFactor() external pure returns (uint256);
    function distributeProposerRewards(uint256 _ntnReward) external payable;
    function getLookbackWindow() external view virtual returns (uint256,bool);
    function getTotalEffort() external view virtual returns (uint256);


    /**
    * @dev Event emitted after a successful slashing.
    */
    event InactivitySlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound);
}