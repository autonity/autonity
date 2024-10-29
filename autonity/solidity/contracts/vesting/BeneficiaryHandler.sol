// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../AccessAutonity.sol";

/** @title Manager of the beneficiary list */
abstract contract BeneficiaryHandler is AccessAutonity {

    /**
     * @dev Stores the unique ids of contracts assigned to a beneficiary, but beneficiary does not need to know the id
     * beneficiary will number his contracts as: 0 for first contract, 1 for 2nd and so on.
     * We can get the unique contract id from beneficiaryContracts as follows:
     * `beneficiaryContracts[beneficiary][0]` is the unique id of his first contract
     * `beneficiaryContracts[beneficiary][1]` is the unique id of his 2nd contract and so on
     */
    mapping(address => uint256[]) internal beneficiaryContracts;
    uint256 private totalContractsCreated;

    event BeneficiaryChanged(address indexed newBeneficiary, address indexed oldBeneficiary, uint256 contractID);

    /*
    ============================================================
         Internals
    ============================================================
     */

    function _newContractCreated(address _beneficiary) internal returns (uint256) {
        uint256 _contractID = totalContractsCreated;
        beneficiaryContracts[_beneficiary].push(_contractID);
        totalContractsCreated++;
        return _contractID;
    }

    /**
     * @notice Changes the beneficiary of some contract to the recipient address. The recipient address can release tokens from the contract as it unlocks.
     * @param _beneficiary beneficiary address whose contract will be canceled
     * @param _contractID unique global id of the contract
     * @param _recipient whome the contract is transferred to
     */
    function _changeContractBeneficiary(
        address _beneficiary, uint256 _contractID, address _recipient
    ) internal {
        uint256[] storage _contractIDs = beneficiaryContracts[_beneficiary];
        uint256[] memory _newContractIDs = new uint256[] (_contractIDs.length - 1);
        uint256 j = 0;
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            if (_contractIDs[i] == _contractID) {
                continue;
            }
            _newContractIDs[j++] = _contractIDs[i];
        }
        beneficiaryContracts[_beneficiary] = _newContractIDs;
        beneficiaryContracts[_recipient].push(_contractID);
        emit BeneficiaryChanged(_recipient, _beneficiary, _contractID);
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns the number of contracts entitled to some beneficiary.
     * @param _beneficiary address of the beneficiary
     */
    function totalContracts(address _beneficiary) virtual external view returns (uint256) {
        return beneficiaryContracts[_beneficiary].length;
    }

    /**
     * @notice Returns a unique id for each contract.
     * @param _beneficiary address of the contract holder
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones)
     */
    function getUniqueContractID(address _beneficiary, uint256 _id) public view returns (uint256) {
        require(beneficiaryContracts[_beneficiary].length > _id, "invalid contract id");
        return beneficiaryContracts[_beneficiary][_id];
    }
}