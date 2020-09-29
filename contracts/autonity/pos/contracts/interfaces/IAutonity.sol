// SPDX-License-Identifier: MIT

pragma solidity ^0.7.0;
pragma experimental ABIEncoderV2;

/**
* @dev Protocol interface of the Autonity Contract
*/

interface IAutonity {

    struct CommitteeMember {
        address payable addr;
        uint256 votingPower;
    }

    enum UserType { Participant, Stakeholder, Validator}

    /**
    * @return `bytecode` the new contract bytecode.
    * @return `contractAbi` the new contract ABI.
    * @dev Implementation of {IAutonity retrieveContract}.
    */
    function retrieveContract() external view returns(string memory, string memory);

    function getMinimumGasPrice() external view returns(uint256);

    function getProposer(uint256 height, uint256 round) external view returns(address);

    function getCommittee() external view returns (CommitteeMember[] memory committee);

    function getWhitelist() external view returns (string[] memory);

    /** @dev finalize is the block state finalisation function. It is called
    * each block after processing every transactions within it. It must be restricted to the
    * protocol only.
    *
    * @param amount The amount of transaction fees collected for this block.
    * @return upgrade Set to true if an autonity contract upgrade is available.
    * @return committee The next block consensus committee.
    */
    function finalize(uint256 amount) external returns(bool upgrade, CommitteeMember[] memory committee);


    function retrieveState() external view returns(
        address[] memory _addr,
        string[] memory _enode,
        uint256[] memory _userType,
        uint256[] memory _stake,
        address _operatorAccount,
        uint256 _minGasPrice,
        uint256 _committeeSize,
        string memory _contractVersion);

    /**
     * @dev Emitted when the Minimum Gas Price was updated and set to `gasPrice`.
     *
     * Note that `gasPrice` may be zero.
     */
    event MinimumGasPriceUpdated(uint256 gasPrice);

    /**
     * @dev Emitted when the Autonity Contract was upgrae to a new version (`version`).
     */
    event ContractUpgraded(string version);
}
