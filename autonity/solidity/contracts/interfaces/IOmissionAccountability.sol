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
    function finalize(address[] memory absentees, address proposer, uint256 proposerEffort, bool isProposerOmissionFaulty, uint256 lastEpochBlock, uint256 delta, bool epochEnded) external;
    function setCommittee(address[] memory _committee) external;
}