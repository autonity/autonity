// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../interfaces/INonStakableVestingVault.sol";
import "./ContractBase.sol";

contract NonStakableVesting is INonStakableVestingVault, ContractBase {

    /**
     * @notice The total amount of funds to create new locked non-stakable schedules.
     * The balance is not immediately available at the vault.
     * Rather the unlocked amount of schedules is minted at epoch end.
     * The balance tells us the max size of a newly created schedule.
     * See `createSchedule()`
     */
    uint256 public totalNominal;

    /**
     * @notice The maximum duration of any schedule or contract
     */
    uint256 public maxAllowedDuration;

    struct Schedule {
        uint256 start;
        uint256 cliffDuration;
        uint256 totalDuration;
        uint256 amount;
        uint256 unsubscribedAmount;
        uint256 totalUnlocked;
        uint256 totalUnlockedUnsubscribed;
        uint256 lastUnlockTime;
    }

    /**
     * @dev Stores all the schedules, there should not be too many of them, for the sake of efficiency
     * of `unlockTokens()` function
     */
    Schedule[] private schedules;

    /** @dev ID of schedule that some contract is subscribed to */
    mapping(uint256 => uint256) private subscribedTo;

    constructor(
        address payable _autonity, address _operator
    ) ContractBase(_autonity, _operator) {}

    /**
     * @notice Creates a new schedule.
     * @dev The schedule has unsubscribedAmount = amount initially. As new contracts are subscribed to the schedule, its unsubscribedAmount decreases.
     * At any point, `subscribedAmount of schedule = amount - unsubscribedAmount`.
     * @param _amount total amount of the schedule
     * @param _startTime start time
     * @param _cliffDuration cliff period, after _cliffDuration + _startTime, the schedule will have claimables
     * @param _totalDuration total duration of the schedule
     * @custom:restricted-to operator account
     */
    function createSchedule(
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration
    ) virtual public onlyOperator {
        require(totalNominal >= _amount, "not enough funds to create a new schedule");
        require(maxAllowedDuration >= _totalDuration, "schedule total duration exceeds max allowed duration");

        schedules.push(Schedule(_startTime, _cliffDuration, _totalDuration, _amount, _amount, 0, 0, 0));
        totalNominal -= _amount;
    }

    /**
     * @notice Creates a new non-stakable contract which subscribes to some schedule.
     * @dev If the contract is created before cliff period has passed, the beneficiary is entitled to NTN as it unlocks.
     * Otherwise, the contract already has some unlocked NTN which is not entitled to beneficiary. However, NTN that will
     * be unlocked in future will be entitled to beneficiary.
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _scheduleID schedule to subscribe
     * @custom:restricted-to operator account
     */
    function newContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _scheduleID
    ) virtual onlyOperator public {
        require(_scheduleID < schedules.length, "invalid schedule ID");

        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.unsubscribedAmount >= _amount, "not enough funds to create a new contract under schedule");

        uint256 _contractID = _createContract(
            _beneficiary, _amount, _schedule.start, _schedule.cliffDuration, _schedule.totalDuration, false
        );

        subscribedTo[_contractID] = _scheduleID;

        if (_schedule.lastUnlockTime >= _schedule.start + _schedule.cliffDuration) {
            // We have created the contract, but it already have some funds uncloked and claimable
            // those unlocked funds are unlocked from unsubscribed funds of the schedule total funds
            // which have already been transferred to treasuryAccount.
            // So the beneficiary will get the funds that will be unlocked in future
            
            // calculate unlocked portion of the unsubscribeds funds from this contract
            // it is the same as calling _unlockedFunds, but we calculate it this way
            // to account for all the _schedule.totalUnlockedUnsubscribed funds
            // otherwise there could be some _schedule.totalUnlockedUnsubscribed funds remaining
            // due to integer division and precision loss
            Contract storage _contract = contracts[_contractID];
            uint256 _unlockedFromUnsubscribed = (_contract.currentNTNAmount * _schedule.totalUnlockedUnsubscribed) / _schedule.unsubscribedAmount;
            _schedule.totalUnlockedUnsubscribed -= _unlockedFromUnsubscribed;

            // the following will prevent the beneficiary to claim _unlockedFromUnsubscribed amount
            // but the contract will follow the same linear vesting function
            _contract.currentNTNAmount -= _unlockedFromUnsubscribed;
            _contract.withdrawnValue += _unlockedFromUnsubscribed;
        }
        _schedule.unsubscribedAmount -= _amount;
    }

    /**
     * @notice Sets the `totalNominal` value.
     * @custom:restricted-to operator account
     */
    function setTotalNominal(uint256 _totalNominal) virtual external onlyOperator {
        totalNominal = _totalNominal;
    }

    /**
     * @notice Sets the max allowed duration of any schedule or contract.
     * @custom:restricted-to operator account
     */
    function setMaxAllowedDuration(uint256 _newMaxDuration) virtual external onlyOperator {
        maxAllowedDuration = _newMaxDuration;
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
     * @notice Unlock tokens of all schedules upto current time.
     * @dev It calculates the newly unlocked tokens upto current time and also updates the amount
     * of total unlocked tokens and the time of unlock for each schedule
     * Autonity must mint new unlocked tokens, because this contract knows that for each schedule,
     * `schedule.totalUnlocked` tokens are now unlocked and available to release
     * @return _newUnlockedSubscribed tokens unlocked from contract subscribed to some schedule
     * @return _newUnlockedUnsubscribed tokens unlocked from schedule.unsubscribedAmount, which is not subscribed by any contract
     * @dev `newUnlockedSubscribed` goes to the balance of address(this) and `newUnlockedUnsubscribed` goes to the treasury address.
     * See `finalize()` in Autonity.sol
     * @custom:restricted-to Autonity contract
     */
    function unlockTokens() external onlyAutonity returns (uint256 _newUnlockedSubscribed, uint256 _newUnlockedUnsubscribed) {
        uint256 _currentTime = block.timestamp;
        uint256 _totalNewUnlocked;
        for (uint256 i = 0; i < schedules.length; i++) {
            Schedule storage _schedule = schedules[i];
            if (_schedule.cliffDuration + _schedule.start > _currentTime || _schedule.amount == _schedule.totalUnlocked) {
                // we did not reach cliff, or we have unlocked everything
                continue;
            }

            _schedule.lastUnlockTime = _currentTime;
            uint256 _unlocked = _calculateTotalUnlockedFunds(_schedule.start, _schedule.totalDuration, _currentTime, _schedule.amount);
            
            if (_unlocked < _schedule.totalUnlocked) {
                // if this happens, then there is something wrong and it needs immediate attention
                _unlocked = _schedule.totalUnlocked;
            }
            _totalNewUnlocked += _unlocked - _schedule.totalUnlocked;
            _schedule.totalUnlocked = _unlocked;

            _unlocked = _calculateTotalUnlockedFunds(_schedule.start, _schedule.totalDuration, _currentTime, _schedule.unsubscribedAmount);

            if (_unlocked < _schedule.totalUnlockedUnsubscribed) {
                // if this happens, then there is something wrong and it needs immediate attention
                _unlocked = _schedule.totalUnlockedUnsubscribed;
            }
            _newUnlockedUnsubscribed += _unlocked - _schedule.totalUnlockedUnsubscribed;
            _schedule.totalUnlockedUnsubscribed = _unlocked;
        }
        _newUnlockedSubscribed = _totalNewUnlocked - _newUnlockedUnsubscribed;
    }

    /**
     * @dev Calculates the total value of the contract, which is constant for non stakable contracts.
     * @param _contractID unique global id of the contract
     */
    function _calculateTotalValue(uint256 _contractID) private view returns (uint256) {
        Contract storage _contract = contracts[_contractID];
        return _contract.currentNTNAmount + _contract.withdrawnValue;
    }

    /**
     * @dev Calculates the amount of withdrawable funds upto `schedule.lastUnlockTime`, which is the last epoch time.
     * where schedule = schedule subsribed by the contract.
     * The unlock mechanism is epoch based, but instead of taking time from `autonity.lastEpochBlock()`, we take the time
     * from `schedule.lastUnlockTime`. Because the locked tokens are not minted from genesis. This way it is ensured that
     * the unlocked tokens are minted at epoch end.
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
     * @notice Returns the amount of withdrawable funds upto the last epoch time.
     */
    function unlockedFunds(
        address _beneficiary, uint256 _id
    ) virtual external view returns (uint256) {
        return _unlockedFunds(_getUniqueContractID(_beneficiary, _id));
    }

    /**
     * @notice Returns some schedule
     * @param _id id of some schedule numbered from 0 to (n-1), where n = total schedules created
     */
    function getSchedule(uint256 _id) external view returns (Schedule memory) {
        require(schedules.length > _id, "schedule does not exist");
        return schedules[_id];
    }
}