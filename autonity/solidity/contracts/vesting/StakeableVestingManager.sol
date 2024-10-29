// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../interfaces/IStakeableVestingManager.sol";
import "./BeneficiaryHandler.sol";
import "./stakeable/StakeableVestingLogic.sol";
import "./stakeable/StakeableVestingState.sol";

/**
 * @title Manager Contract for stakeable vesting smart contracts
 * @notice It only creates new stakeable vesting contract which deploys new smart contract for each new stakeable vesting contract.
 * It can also manage the beneficiary of an existing stakeable vesting contract.
 */
contract StakeableVestingManager is BeneficiaryHandler, IStakeableVestingManager {
    uint256 public contractVersion = 1;

    address public stakeableVestingLogicContract;

    IStakeableVesting[] private contracts;

    constructor(address payable _autonity) AccessAutonity(_autonity) {
        stakeableVestingLogicContract = address(new StakeableVestingLogic(_autonity));
    }

    function setStakeableVestingLogicContract(address _contract) virtual external onlyOperator {
        require(_contract != address(0), "invalid contract address");
        stakeableVestingLogicContract = _contract;
    }

    /**
     * @notice Creates a new stakeable contract. the operation is invalid if the cliff duration is already past.
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _startTime start time of the contract
     * @param _cliffDuration cliff duration of the contract
     * @param _totalDuration total duration of the contract
     * @custom:restricted-to operator account
     */
    function newContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration
    ) virtual onlyOperator public {
        require(_startTime >= block.timestamp, "contract cannot start before creation");
        require(autonity.balanceOf(address(this)) >= _amount, "not enough stake reserved to create a new contract");

        uint256 _contractID = _newContractCreated(_beneficiary);
        require(_contractID == contracts.length, "invalid contract id");
        IStakeableVesting _stakeableVestingContract = IStakeableVesting(
            address(new StakeableVestingState(payable(autonity)))
        );
        _stakeableVestingContract.createContract(
            _beneficiary,
            _amount,
            _startTime,
            _cliffDuration,
            _totalDuration
        );
        contracts.push(_stakeableVestingContract);
        bool _sent = autonity.transfer(address(_stakeableVestingContract), _amount);
        require(_sent, "failed to transfer NTN");
    }

    /**
     * @notice Changes the beneficiary of some contract to the recipient address. The recipient address can release and stake tokens from the contract.
     * Rewards which have been entitled to the beneficiary due to bonding from this contract are not transferred to recipient, but transferred to the old beneficiary.
     * @param _beneficiary beneficiary address whose contract will be canceled
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding already canceled ones)
     * @param _recipient whome the contract is transferred to
     * @custom:restricted-to operator account
     */
    function changeContractBeneficiary(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator {
        uint256 _contractID = getUniqueContractID(_beneficiary, _id);
        contracts[_contractID].changeContractBeneficiary(_recipient);
        _changeContractBeneficiary(_beneficiary, _contractID, _recipient);
    }

    /**
     * @dev Receive Auton function https://solidity.readthedocs.io/en/v0.7.2/contracts.html#receive-ether-function
     */
    receive() external payable {}

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns the smart contract account that holds the corresponding stake-able vesting contract.
     * @param _uniqueContractID unique id of the contract
     */
    function getContractAccount(uint256 _uniqueContractID) external virtual view returns (IStakeableVesting) {
        require(_uniqueContractID < contracts.length, "invalid contract id");
        return contracts[_uniqueContractID];
    }

    /**
     * @notice Returns the smart contract account that holds the corresponding stake-able vesting contract.
     * @param _beneficiary address of the beneficiary of the contract
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding already canceled ones)
     */
    function getContractAccount(address _beneficiary, uint256 _id) external virtual view returns (IStakeableVesting) {
        return contracts[getUniqueContractID(_beneficiary, _id)];
    }

    /**
     * @notice Returns all the smart contract accounts that holds the corresponding stake-able vesting contract.
     * @param _beneficiary address of the beneficiary of the contract
     */
    function getContractAccounts(address _beneficiary) external virtual view returns (IStakeableVesting[] memory) {
        uint256[] storage _contractIDs = beneficiaryContracts[_beneficiary];
        IStakeableVesting[] memory _contracts = new IStakeableVesting[] (_contractIDs.length);
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            _contracts[i] = contracts[_contractIDs[i]];
        }
        return _contracts;
    }

    /**
     * @notice Returns all the contracts entitled to `_beneficiary`.
     * @param _beneficiary address of the beneficiary of the contract
     */
    function getContracts(address _beneficiary) external virtual view returns (ContractBase.Contract[] memory) {
        uint256[] storage _contractIDs = beneficiaryContracts[_beneficiary];
        ContractBase.Contract[] memory _res = new ContractBase.Contract[] (_contractIDs.length);
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            _res[i] = contracts[_contractIDs[i]].getContract();
        }
        return _res;
    }

}
