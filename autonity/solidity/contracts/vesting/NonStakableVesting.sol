// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./ContractBase.sol";

contract NonStakableVesting is ContractBase {

    struct ScheduleTracker {
        uint256 unsubscribedAmount;
        uint256 expiredFromContract;
        uint256 withdrawnAmount; // withdrawn by treasury account
        bool initiated;
    }

    mapping(uint256 => ScheduleTracker) internal scheduleTracker;

    /** @dev ID of schedule that some contract is subscribed to. */
    mapping(uint256 => uint256) internal subscribedTo;

    constructor(
        address payable _autonity, address _operator
    ) ContractBase(_autonity, _operator) {}

    /**
     * @notice Creates a new non-stakable contract which subscribes to some schedule.
     * @dev If the contract is created before cliff period has passed, the beneficiary is entitled to all the NTN after cliff period.
     * Otherwise, the contract already has some unlocked NTN which is not entitled to beneficiary. However, NTN that will
     * be unlocked in future will be entitled to beneficiary.
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _scheduleID schedule to subscribe
     * @param _cliffDuration cliff duration in seconds for the contract
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

        if (!_scheduleTracker.initiated) {
            _scheduleTracker.unsubscribedAmount = _schedule.totalAmount;
            _scheduleTracker.initiated = true;
        }
        require(_scheduleTracker.unsubscribedAmount >= _amount, "not enough funds to create a new contract under schedule");

        uint256 _expiredFunds = _calculateUnlockedFunds(_schedule.unlockedAmount, _schedule.totalAmount, _amount);
        uint256 _contractID = _createContract(
            _beneficiary, _amount, _expiredFunds, _schedule.start, _cliffDuration, _schedule.totalDuration, false
        );

        subscribedTo[_contractID] = _scheduleID;
        _scheduleTracker.unsubscribedAmount -= _amount;
        _scheduleTracker.expiredFromContract += _expiredFunds;
    }

    function releaseFundsForTreasury(uint256 _scheduleID) external virtual {
        require(msg.sender == autonity.getTreasuryAccount(), "caller is not treasury account");
        ScheduleController.Schedule memory _schedule = autonity.getSchedule(address(this), _scheduleID);
        ScheduleTracker storage _scheduleTracker = scheduleTracker[_scheduleID];

        if (!_scheduleTracker.initiated) {
            _scheduleTracker.unsubscribedAmount = _schedule.totalAmount;
            _scheduleTracker.initiated = true;
        }
        uint256 _unlocked = _calculateUnlockedFunds(_schedule.unlockedAmount, _schedule.totalAmount, _scheduleTracker.unsubscribedAmount)
                            + _scheduleTracker.expiredFromContract - _scheduleTracker.withdrawnAmount;
        bool _sent = autonity.transfer(msg.sender, _unlocked);
        require(_sent, "transfer failed");
        _scheduleTracker.withdrawnAmount += _unlocked;
    }

    /**
     * @notice Used by beneficiary to transfer all unlocked NTN of some contract to his own address.
     * @param _id id of the contract numbered from 0 to (n-1) where n = total contracts entitled to
     * the beneficiary (excluding canceled ones). So any beneficiary can number their contracts
     * from 0 to (n-1). Beneficiary does not need to know the unique global contract id.
     */
    function releaseAllFunds(uint256 _id) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _releaseNTN(_contractID, _unlockedFunds(_contractID));
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address.
     * @param _amount amount of NTN to release
     */
    function releaseFund(uint256 _id, uint256 _amount) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        require(_amount <= _unlockedFunds(_contractID), "not enough unlocked funds");
        _releaseNTN(_contractID, _amount);
    }

    /**
     * @notice Changes the beneficiary of some contract to the recipient address. The recipient address can release tokens from the contract as it unlocks.
     * @param _beneficiary beneficiary address whose contract will be canceled
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones)
     * @param _recipient whome the contract is transferred to
     * @custom:restricted-to operator account
     */
    function changeContractBeneficiary(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator {
        _changeContractBeneficiary(
            _getUniqueContractID(_beneficiary, _id),
            _beneficiary,
            _recipient
        );
    }

    /**
     * @dev Calculates the total value of the contract, which is constant for non stakable contracts.
     * @param _contractID unique global id of the contract
     */
    function _calculateTotalValue(uint256 _contractID) internal view returns (uint256) {
        Contract storage _contract = contracts[_contractID];
        return _contract.currentNTNAmount + _contract.withdrawnValue + _contract.expiredFunds;
    }

    /**
     * @dev Calculates the amount of withdrawable funds upto `schedule.lastUnlockTime`, which is the last epoch time,
     * where schedule = schedule subsribed by the contract.
     */
    function _unlockedFunds(uint256 _contractID) internal view returns (uint256) {
        ScheduleController.Schedule memory _schedule = autonity.getSchedule(address(this), subscribedTo[_contractID]);
        Contract storage _contract = contracts[_contractID];
        return _calculateUnlockedFunds(
            _schedule.unlockedAmount,
            _schedule.totalAmount,
            _calculateTotalValue(_contractID)
        ) - _contract.withdrawnValue - _contract.expiredFunds;
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
        return _unlockedFunds(_getUniqueContractID(_beneficiary, _id));
    }
}