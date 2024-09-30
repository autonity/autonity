// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./BeneficiaryHandler.sol";
import "./StakableVesting.sol";

contract StakableVestingManager is BeneficiaryHandler {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release contract
    // UNLOCKED are tokens that can be withdrawn
    uint256 public contractVersion = 1;

    /**
     * @notice Sum of total amount of contracts that can be created.
     * Each time a new contract is created, `totalNominal` is decreased.
     * Address(this) should have `totalNominal` amount of NTN availabe at genesis,
     * otherwise withdrawing or bonding from a contract is not possible.
     */
    uint256 public totalNominal;

    StakableVesting[] private contracts;

    constructor(address payable _autonity) AccessAutonity(_autonity) {}

    /**
     * @notice Creates a new stakable contract.
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _startTime start time of the vesting
     * @param _cliffDuration cliff period
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
        require(_startTime + _cliffDuration >= autonity.lastEpochTime(), "contract cliff duration is past");
        require(totalNominal >= _amount, "not enough stake reserved to create a new contract");

        uint256 _contractID = _newContractCreated(_beneficiary);
        require(_contractID == contracts.length, "invalid contract id");
        contracts.push(
            new StakableVesting(
                payable(autonity),
                _beneficiary,
                _amount,
                _startTime,
                _cliffDuration,
                _totalDuration
            )
        );
        totalNominal -= _amount;
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
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        contracts[_contractID].changeContractBeneficiary(_recipient);
        _changeContractBeneficiary(_beneficiary, _contractID, _recipient);
    }

    /**
     * @notice Set the value of totalNominal.
     * In case totalNominal is increased, the increased amount should be minted
     * and transferred to the address of this contract, otherwise newly created vesting
     * contracts will not have funds to withdraw or bond. See `ewContract()`.
     * @custom:restricted-to operator account
     */
    function setTotalNominal(uint256 _newTotalNominal) virtual external onlyOperator {
        totalNominal = _newTotalNominal;
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

    function getContractAccount(address _beneficiary, uint256 _id) external virtual view returns (StakableVesting) {
        return contracts[_getUniqueContractID(_beneficiary, _id)];
    }

    function getContracts(address _beneficiary) external virtual view returns (ContractBase.Contract[] memory) {
        uint256[] storage _contractIDs = beneficiaryContracts[_beneficiary];
        ContractBase.Contract[] memory _res = new ContractBase.Contract[] (_contractIDs.length);
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            _res[i] = contracts[_contractIDs[i]].getContract();
        }
        return _res;
    }

}
