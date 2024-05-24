// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";

contract ContractBase {

    struct Contract {
        uint256 currentNTNAmount;
        uint256 withdrawnValue;
        uint256 start;
        uint256 cliff;
        uint256 end;
        bool canStake;
    }

    // stores the unique ids of contracts assigned to a beneficiary, but beneficiary does not need to know the id
    // beneficiary will number his contracts as: 0 for first contract, 1 for 2nd and so on
    // we can get the unique contract id from beneficiaryContracts as follows
    // beneficiaryContracts[beneficiary][0] is the unique id of his first contract
    // beneficiaryContracts[beneficiary][1] is the unique id of his 2nd contract and so on
    mapping(address => uint256[]) internal beneficiaryContracts;

    // list of all contracts
    Contract[] internal contracts;

    Autonity internal autonity;
    address private operator;

    constructor(address payable _autonity, address _operator) {
        autonity = Autonity(_autonity);
        operator = _operator;
    }

    function _createContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffTime,
        uint256 _endTime,
        bool _canStake
    ) internal returns (uint256) {
        require(_startTime >= block.timestamp, "contract cannot start before creating it");
        require(_cliffTime >= _startTime, "cliff must be greater than or equal to start");
        require(_endTime > _cliffTime, "end must be greater than cliff");

        uint256 _contractID = contracts.length;
        contracts.push(
            Contract(
                _amount, 0, _startTime, _cliffTime, _endTime, _canStake
            )
        );
        beneficiaryContracts[_beneficiary].push(_contractID);
        return _contractID;
    }

    function _releaseNTN(
        uint256 _contractID, uint256 _amount
    ) internal returns (uint256 _remaining) {
        Contract storage _contract = contracts[_contractID];
        if (_amount > _contract.currentNTNAmount) {
            _remaining = _amount - _contract.currentNTNAmount;
            _updateAndTransferNTN(_contractID, msg.sender, _contract.currentNTNAmount);
        }
        else if (_amount > 0) {
            _updateAndTransferNTN(_contractID, msg.sender, _amount);
        }
    }

    function _calculateAvailableUnlockedFunds(
        uint256 _contractID, uint256 _totalValue, uint256 _time
    ) internal view returns (uint256) {
        Contract storage _contract = contracts[_contractID];
        require(_time >= _contract.cliff, "cliff period not reached yet");

        uint256 _unlocked = _calculateUnlockedFunds(_contract.start, _contract.end, _time, _totalValue);
        if (_unlocked > _contract.withdrawnValue) {
            return _unlocked - _contract.withdrawnValue;
        }
        return 0;
    }

    function _calculateUnlockedFunds(
        uint256 _start, uint256 _end, uint256 _time, uint256 _totalAmount
    ) internal pure returns (uint256) {
        if (_time >= _end) {
            return _totalAmount;
        }
        return _totalAmount * (_time - _start) / (_end - _start);
    }

    function _cancelContract(
        address _beneficiary, uint256 _id, address _recipient
    ) internal {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        _changeContractBeneficiary(_contractID, _beneficiary, _recipient);
    }

    function _changeContractBeneficiary(
        uint256 _contractID, address _oldBeneficiary, address _newBeneficiary
    ) private {
        uint256[] storage _contractIDs = beneficiaryContracts[_oldBeneficiary];
        uint256[] memory _newContractIDs = new uint256[] (_contractIDs.length - 1);
        uint256 j = 0;
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            if (_contractIDs[i] == _contractID) {
                continue;
            }
            _newContractIDs[j++] = _contractIDs[i];
        }
        beneficiaryContracts[_oldBeneficiary] = _newContractIDs;
        beneficiaryContracts[_newBeneficiary].push(_contractID);
    }

    /**
     * @dev returns a unique id for each contract
     * @param _beneficiary address of the contract holder
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones)
     */
    function _getUniqueContractID(address _beneficiary, uint256 _id) internal view returns (uint256) {
        require(beneficiaryContracts[_beneficiary].length > _id, "invalid contract id");
        return beneficiaryContracts[_beneficiary][_id];
    }

    function _updateAndTransferNTN(uint256 _contractID, address _to, uint256 _amount) internal {
        Contract storage _contract = contracts[_contractID];
        _contract.currentNTNAmount -= _amount;
        _contract.withdrawnValue += _amount;
        _transferNTN(_to, _amount);
    }

    function _transferNTN(address _to, uint256 _amount) internal {
        bool _sent = autonity.transfer(_to, _amount);
        require(_sent, "NTN not transferred");
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice returns a contract entitled to _beneficiary
     * @param _beneficiary beneficiary address
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones)
     */
    function getContract(address _beneficiary, uint256 _id) virtual external view returns (Contract memory) {
        return contracts[_getUniqueContractID(_beneficiary, _id)];
    }

    /**
     * @notice returns the list of current contracts assigned to a beneficiary
     * @param _beneficiary address of the beneficiary
     */
    function getContracts(address _beneficiary) virtual external view returns (Contract[] memory) {
        uint256[] storage _contractIDs = beneficiaryContracts[_beneficiary];
        Contract[] memory _res = new Contract[](_contractIDs.length);
        for (uint256 i = 0; i < _res.length; i++) {
            _res[i] = contracts[_contractIDs[i]];
        }
        return _res;
    }

    /**
     * @notice returns if beneficiary can stake from his contract
     * @param _beneficiary beneficiary address
     */
    function canStake(address _beneficiary, uint256 _id) virtual external view returns (bool) {
        return contracts[_getUniqueContractID(_beneficiary, _id)].canStake;
    }

    /**
     * @notice returns the number of schudeled entitled to some beneficiary
     * @param _beneficiary address of the beneficiary
     */
    function totalContracts(address _beneficiary) virtual external view returns (uint256) {
        return beneficiaryContracts[_beneficiary].length;
    }


    /*
    ============================================================

        Modifiers

    ============================================================
     */

    /**
     * @dev Modifier that checks if the caller is the governance operator account.
     */
    modifier onlyOperator {
        require(operator == msg.sender, "caller is not the operator");
        _;
    }

    modifier onlyAutonity {
        require(msg.sender == address(autonity) , "function restricted to Autonity contract");
        _;
    }
}