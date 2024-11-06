// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

interface IAccountability {
    /**
    * @notice called by the Autonity Contract at block finalization, before
    * processing reward redistribution.
    * @param _epochEnd whether or not the current block is the last one from the epoch.
    */
    function finalize(bool _epochEnd) external;

    /**
    * @notice distribute slashing rewards to reporters.
    * @param _validator the address of the validator node being slashed.
    */
    function distributeRewards(address _validator, uint256 _ntnReward) external payable;

    /**
    * @notice called by the Autonity Contract when the committee is updated.
    * @param _committee the new committee member addresses;
    */
    function setCommittee(address[] memory _committee) external;

    /**
    * @dev Event emitted when a fault proof has been submitted. The reported validator
    * will be silencied and slashed at the end of the current epoch.
    */
    event NewFaultProof(address indexed _offender, uint256 _severity, uint256 _id);

    /**
    * @dev Event emitted after receiving an accusation, the reported validator has
    * a certain amount of time to submit a proof-of-innocence, otherwise, he gets slashed.
    */
    event NewAccusation(address indexed _offender, uint256 _severity, uint256 _id);

    /**
    * @dev Event emitted after receiving a proof-of-innocence cancelling an accusation.
    */
    event InnocenceProven(address indexed _offender, uint256 _id);

    /**
    * @dev Event emitted after a successful slashing.
    */
    event SlashingEvent(address validator, uint256 amount, uint256 releaseBlock, bool isJailbound, uint256 eventId);
}
