// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./BeneficiaryHandler.sol";
import "./ContractBase.sol";

contract NonStakableVesting is BeneficiaryHandler, ContractBase {

    struct ScheduleTracker {
        uint256 unsubscribedAmount;
        uint256 expiredFromContract;
        uint256 withdrawnAmount; // withdrawn by treasury account
        bool initialized;
    }

    mapping(uint256 => ScheduleTracker) internal scheduleTracker;

    /** @dev List of all contracts */
    ContractBase.Contract[] internal contracts;
    uint256[] internal expiredFundsFromContract;

    /** @dev ID of schedule that some contract is subscribed to. */
    mapping(uint256 => uint256) internal subscribedTo;

    constructor(address payable _autonity) AccessAutonity(_autonity) {}

    /**
     * @notice Creates a new non-stakable contract which subscribes to some schedule.
     * If the contract is created before the start timestamp, the beneficiary is entitled to NTN as it unlocks.
     * Otherwise, the contract already has some unlocked NTN which is not entitled to beneficiary. However, NTN that will
     * be unlocked in future will be entitled to beneficiary.
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _scheduleID schedule to subscribe
     * @param _cliffDuration cliff duration of the contract
     * @custom:restricted-to operator account
     */
    function newContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _scheduleID,
        uint256 _cliffDuration
    ) virtual onlyOperator public {
        ScheduleController.Schedule memory _schedule = autonity.getSchedule(address(this), _scheduleID);
        ScheduleTracker storage _scheduleTracker = scheduleTracker[_scheduleID];

        if (!_scheduleTracker.initialized) {
            _initiateSchedule(_scheduleTracker, _schedule.totalAmount);
        }
        require(_scheduleTracker.unsubscribedAmount >= _amount, "not enough funds to create a new contract under schedule");
        uint256 _contractID = _newContractCreated(_beneficiary);
        require(contracts.length == _contractID, "invalid contract id");

        // `_expiredFunds` = the amount of funds that have been unlocked already, in case the contract was created later than the `_schedule.start`
        // the `_expiredFunds` belongs to the treasury account, not the `_beneficiary`
        uint256 _expiredFunds = _calculateUnlockedFunds(_schedule.unlockedAmount, _schedule.totalAmount, _amount);
        ContractBase.Contract memory _contract = _createContract(
            _beneficiary, _amount - _expiredFunds, _schedule.start, _cliffDuration, _schedule.totalDuration, false
        );
        contracts.push(_contract);
        expiredFundsFromContract.push(_expiredFunds);

        subscribedTo[_contractID] = _scheduleID;
        _scheduleTracker.unsubscribedAmount -= _amount;
        _scheduleTracker.expiredFromContract += _expiredFunds;
    }

    /**
     * @notice Transfers all the unsubscribed funds of the schedule to the treasury account after the schedule total duration has expired.
     * @param _scheduleID id of the schedule
     * @custom:restricted-to treasury account
     */
    function releaseFundsForTreasury(uint256 _scheduleID) virtual external onlyAutonityTreasury {
        ScheduleController.Schedule memory _schedule = autonity.getSchedule(address(this), _scheduleID);
        require(_schedule.lastUnlockTime >= _schedule.start + _schedule.totalDuration, "schedule total duration not expired yet");
        ScheduleTracker storage _scheduleTracker = scheduleTracker[_scheduleID];

        if (!_scheduleTracker.initialized) {
            _initiateSchedule(_scheduleTracker, _schedule.totalAmount);
        }
        uint256 _withdrawable = _scheduleTracker.unsubscribedAmount + _scheduleTracker.expiredFromContract - _scheduleTracker.withdrawnAmount;
        _transferNTN(msg.sender, _withdrawable);
        _scheduleTracker.withdrawnAmount += _withdrawable;
    }

    /**
     * @notice Changes the beneficiary of some contract to the recipient address. The recipient address can release and stake tokens from the contract.
     * @param _beneficiary beneficiary address whose contract will be canceled
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding already canceled ones)
     * @param _recipient whome the contract is transferred to
     * @custom:restricted-to operator account
     */
    function changeContractBeneficiary(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator {
        uint256 _contractID = getUniqueContractID(_beneficiary, _id);
        _changeContractBeneficiary(_beneficiary, _contractID, _recipient);
    }

    /**
     * @notice Used by beneficiary to transfer all unlocked NTN of some contract to his own address.
     * @param _id id of the contract numbered from 0 to (n-1) where n = total contracts entitled to
     * the beneficiary (excluding canceled ones). So any beneficiary can number their contracts
     * from 0 to (n-1). Beneficiary does not need to know the unique global contract id.
     */
    function releaseAllNTN(uint256 _id) virtual external {
        uint256 _contractID = getUniqueContractID(msg.sender, _id);
        _releaseNTN(contracts[_contractID], _unlockedFunds(_contractID));
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address.
     * @param _amount amount of NTN to release
     * @param _id id of the contract numbered from 0 to (n-1) where n = total contracts entitled to
     * the beneficiary (excluding canceled ones). So any beneficiary can number their contracts
     * from 0 to (n-1). Beneficiary does not need to know the unique global contract id.
     */
    function releaseNTN(uint256 _id, uint256 _amount) virtual external {
        uint256 _contractID = getUniqueContractID(msg.sender, _id);
        require(_amount <= _unlockedFunds(_contractID), "not enough unlocked funds");
        _releaseNTN(contracts[_contractID], _amount);
    }

    /*
    ============================================================
         Internals
    ============================================================
     */

    function _initiateSchedule(ScheduleTracker storage _scheduleTracker, uint256 _totalAmount) internal {
        _scheduleTracker.unsubscribedAmount = _totalAmount;
        _scheduleTracker.initialized = true;
    }

    /**
     * @dev Calculates the total value of the contract, which is constant for non stakable contracts.
     * @param _contractID unique global id of the contract
     */
    function _calculateTotalValue(uint256 _contractID) internal view returns (uint256) {
        Contract storage _contract = contracts[_contractID];
        return _contract.currentNTNAmount + _contract.withdrawnValue + expiredFundsFromContract[_contractID];
    }

    /**
     * @dev Calculates the amount of withdrawable funds upto `schedule.lastUnlockTime`, which is the last epoch time,
     * where schedule = schedule subsribed by the contract.
     */
    function _unlockedFunds(uint256 _contractID) internal view returns (uint256) {
        ContractBase.Contract storage _contract = contracts[_contractID];
        require(_contract.start + _contract.cliffDuration <= block.timestamp, "cliff period not reached yet");
        ScheduleController.Schedule memory _schedule = autonity.getSchedule(address(this), subscribedTo[_contractID]);
        return _calculateUnlockedFunds(
            _schedule.unlockedAmount,
            _schedule.totalAmount,
            _calculateTotalValue(_contractID)
        ) - _contract.withdrawnValue - expiredFundsFromContract[_contractID];
    }

    function _calculateUnlockedFunds(
        uint256 _scheduleUnlockAmount, uint256 _scheduleTotalAmount, uint256 _contractTotalAmount
    ) internal pure returns (uint256) {
        return (_scheduleUnlockAmount * _contractTotalAmount) / _scheduleTotalAmount;
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns the amount of withdrawable funds upto the last epoch time.
     */
    function unlockedFunds(
        address _beneficiary, uint256 _id
    ) virtual external view returns (uint256) {
        return _unlockedFunds(getUniqueContractID(_beneficiary, _id));
    }

    function getExpiredFunds(address _beneficiary, uint256 _id) virtual external view returns (uint256) {
        return expiredFundsFromContract[getUniqueContractID(_beneficiary, _id)];
    }

    function getContract(address _beneficiary, uint256 _id) virtual external view returns (ContractBase.Contract memory) {
        return contracts[getUniqueContractID(_beneficiary, _id)];
    }

    function getContracts(address _beneficiary) virtual external view returns (ContractBase.Contract[] memory) {
        uint256[] storage _contractIDs = beneficiaryContracts[_beneficiary];
        ContractBase.Contract[] memory _res = new ContractBase.Contract[] (_contractIDs.length);
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            _res[i] = contracts[_contractIDs[i]];
        }
        return _res;
    }

    /*
    ============================================================

        Modifiers

    ============================================================
     */

    modifier onlyAutonityTreasury {
        require(msg.sender == autonity.getTreasuryAccount(), "caller is not treasury account");
        _;
    }
}