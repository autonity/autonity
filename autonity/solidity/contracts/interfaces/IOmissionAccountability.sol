// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.0;

interface IOmissionAccountability {
    /**
    * @notice called by the Autonity Contract at block finalization, it receives activity report.
    * @param isProposerOmissionFaulty is true when the proposer provides invalid activity proof of current height.
    * @param ids stores faulty proposer's ID when isProposerOmissionFaulty is true, otherwise it carries current height
    * activity proof which is the signers of precommit of current height - dela.
    */
    function finalize(bool isProposerOmissionFaulty, uint256[] memory ids) external;
}