// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../interfaces/INonStakableVestingVault.sol";
import "./ContractBase.sol";

contract NonStakableVesting is INonStakableVestingVault, ContractBase {

    struct Schedule {
        uint256 start;
        uint256 cliff;
        uint256 end;
        uint256 totalAmount;
        uint256 totalUnlocked;
        uint256 lastUnlockTime;
    }

    // stores all the schedules, there should not be too many of them, for the sake of efficiency
    // of unlockTokens() function
    Schedule[] private schedules;

    // id of schedule that some contract is subscribed to
    mapping(uint256 => uint256) private subscribedTo;

    constructor(
        address payable _autonity, address _operator
    ) ContractBase(_autonity, _operator) {}

    /**
     * @notice creates a new schedule, restricted to operator
     * the schedule has totalAmount = 0 initially. As new contracts are subscribed to the schedule, its totalAmount increases
     * At any point, totalAmount of schedule is the sum of totalValue all the contracts that are subscribed to the schedule.
     * totalValue of a contract can be calculated via _calculateTotalValue function
     * @param _startTime start time
     * @param _cliffTime cliff time. cliff period = _cliffTime - _startTime
     * @param _endTime end time, total duration of the schedule = _endTime - _startTime
     */
    function createSchedule(
        uint256 _startTime,
        uint256 _cliffTime,
        uint256 _endTime
    ) virtual public onlyOperator {
        schedules.push(Schedule(_startTime, _cliffTime, _endTime, 0, 0, 0));
    }

    /**
     * @notice creates a new non-stakable contract, restricted to only operator
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _scheduleID schedule to subscribe
     */
    function newContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _scheduleID
    ) virtual onlyOperator public {
        require(_scheduleID < schedules.length, "invalid schedule ID");
        Schedule storage _schedule = schedules[_scheduleID];
        uint256 _contractID = _createContract(
            _beneficiary, _amount, _schedule.start, _schedule.cliff, _schedule.end, false
        );
        _schedule.totalAmount += _amount;
        subscribedTo[_contractID] = _scheduleID;
    }

    /**
     * @notice used by beneficiary to transfer all unlocked NTN of some contract to his own address
     * @param _id id of the contract numbered from 0 to (n-1) where n = total contracts entitled to
     * the beneficiary (excluding canceled ones). So any beneficiary can number their contracts
     * from 0 to (n-1). Beneficiary does not need to know the unique global contract id which can
     * be retrieved via _getUniqueContractID function
     */
    function releaseAllFunds(uint256 _id) virtual external { // onlyActive(_id) {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _releaseNTN(_contractID, _unlockedFunds(_contractID));
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address
     * @param _amount amount of NTN to release
     */
    function releaseFund(uint256 _id, uint256 _amount) virtual external { // onlyActive(_id) {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        require(_amount <= _unlockedFunds(_contractID), "not enough unlocked funds");
        _releaseNTN(_contractID, _amount);
    }

    /**
     * @notice changes the beneficiary of some contract to the _recipient address. _recipient can release tokens from the contract
     * only operator is able to call the function
     * @param _beneficiary beneficiary address whose contract will be canceled
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones)
     * @param _recipient whome the contract is transferred to
     */
    function cancelContract(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator {
        _cancelContract(_beneficiary, _id, _recipient);
    }

    /**
     * @notice Unlock tokens of all schedules upto current time, restricted to autonity only.
     * @dev It calculates the newly unlocked tokens upto current time and also updates the amount
     * of total unlocked tokens and the time of unlock for each schedule
     * Autonity must mint _totalNewUnlocked tokens, because this contract knows that for each _schedule,
     * _schedule.totalUnlocked tokens are now unlocked and available to release
     */
    function unlockTokens() external onlyAutonity returns (uint256 _totalNewUnlocked) {
        uint256 _currentTime = block.timestamp;
        for (uint256 i = 0; i < schedules.length; i++) {
            Schedule storage _schedule = schedules[i];
            if (_schedule.cliff > _currentTime || _schedule.totalAmount == _schedule.totalUnlocked) {
                continue;
            }
            _schedule.lastUnlockTime = _currentTime;
            uint256 _unlocked = _calculateUnlockedFunds(_schedule.start, _schedule.end, _currentTime, _schedule.totalAmount);
            _totalNewUnlocked += _unlocked - _schedule.totalUnlocked;
            _schedule.totalUnlocked = _unlocked;
        }
        return _totalNewUnlocked;
    }

    /**
     * @dev calculates the total value of the contract, which is constant for non stakable contracts
     * @param _contractID unique global id of the contract
     */
    function _calculateTotalValue(uint256 _contractID) private view returns (uint256) {
        Contract storage _contract = contracts[_contractID];
        return _contract.currentNTNAmount + _contract.withdrawnValue;
    }

    /**
     * @dev calculates the amount of funds that are unlocked but not released yet. calculates upto _schedule.lastUnlockTime
     * where _schedule = schedule subsribed by the contract.
     * The unlock mechanism is epoch based, but instead of taking time from autonity.lastEpochBlock(), we take the time
     * from _schedule.lastUnlockTime. Because the locked tokens are not minted from genesis. This way it is ensured that
     * the unlocked tokens are minted by calling the function unlockTokens()
     */
    function _unlockedFunds(uint256 _contractID) private view returns (uint256) {
        return _calculateAvailableUnlockedFunds(
            _contractID,
            _calculateTotalValue(_contractID),
            schedules[subscribedTo[_contractID]].lastUnlockTime
        );
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice returns the amount of unlocked but not yet released funds in NTN for some contract
     */
    function unlockedFunds(
        address _beneficiary, uint256 _id
    ) virtual external view returns (uint256) {
        return _unlockedFunds(_getUniqueContractID(_beneficiary, _id));
    }
}