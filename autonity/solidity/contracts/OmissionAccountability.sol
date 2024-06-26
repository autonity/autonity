// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;
import "./Autonity.sol";
contract OmissionAccountability is IOmissionAccountability {

    struct Config {
        uint256 omissionLoopBackWindow;
        uint256 activityProofRewardRate;
        uint256 maxCommitteeSize;
        uint256 pastPerformanceWeight;
        uint256 initialJailingPeriod;
        uint256 initialProbationPeriod;
        uint256 initialSlashingRate;
    }

    Config public config;
    Autonity internal autonity; // for access control in setters function.
    constructor(address payable _autonity, Config memory _config) {
        autonity = Autonity(_autonity);
        config = _config;
    }

    /**
    * @notice called by the Autonity Contract at block finalization, it receives activity report.
    * @param isProposerOmissionFaulty is true when the proposer provides invalid activity proof of current height.
    * @param ids stores faulty proposer's ID when isProposerOmissionFaulty is true, otherwise it carries current height
    * activity proof which is the signers of precommit of current height - dela.
    */
    function finalize(bool isProposerOmissionFaulty, uint256[] memory ids) external onlyAutonity {
        // todo, fill the look-back window by the input data.
    }

    /**
    * @dev Modifier that checks if the caller is the slashing contract.
    */
    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to the autonity");
        _;
    }
}