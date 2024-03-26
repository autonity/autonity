// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";
import "../interfaces/IStakeProxy.sol";
import "../Liquid.sol";
import "./LiquidRewardManager.sol";
import "./VestingCalculator.sol";

contract VestingManager is IStakeProxy, LiquidRewardManager, VestingCalculator {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release schedule
    // UNLOCKED are tokens that got released but not yet transferred
    uint256 public contractVersion = 1;

    address private operator;

    struct Schedule {
        uint256 totalAmount;
        uint256 start;
        uint256 cliff;
        uint256 end;
        uint256 vestingID;
        bool canStake;
        bool canceled;
    }

    mapping(uint256 => mapping(address => uint256)) private liquidVestingIDs;

    // stores the unique ids of schedules assigned to a beneficiary, but beneficiary does not need to know the id
    // beneficiary will number his schedules as: 0 for first schedule, 1 for 2nd and so on
    // we can get the unique schedule id from beneficiarySchedules as follows
    // beneficiarySchedules[beneficiary][0] is the unique id of his first schedule
    // beneficiarySchedules[beneficiary][1] is the unique id of his 2nd schedule and so on
    mapping(address => uint256[]) private beneficiarySchedules;

    // list of all schedules
    Schedule[] private schedules;

    struct PendingBondingRequest {
        uint256 amount;
        address validator;
    }

    struct PendingUnbondingRequest {
        uint256 amount;
        address validator;
    }

    mapping(uint256 => PendingBondingRequest) private pendingBondingRequest;
    mapping(uint256 => uint256) private pendingBondingVesting;
    mapping(uint256 => PendingUnbondingRequest) private pendingUnbondingRequest;
    mapping(uint256 => uint256) private pendingUnbondingVesting;

    // bondingToSchedule[bondingID] stores the unique schedule (id+1) which requested the bonding
    mapping(uint256 => uint256) private bondingToSchedule;

    // unbondingToSchedule[unbondingID] stores the unique schedule (id+1) which requested the unbonding
    mapping(uint256 => uint256) private unbondingToSchedule;

    mapping(uint256 => address) private cancelRecipient;

    constructor(address payable _autonity, address _operator) LiquidRewardManager(_autonity) {
        operator = _operator;
    }

    /** @notice creates a new schedule
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _startBlock start block of the vesting
     * @param _cliffBlock cliff block of the vesting
     * @param _endBlock end block of the vesting
     * @param _canStake if the NTN can be staked or not
     */
    function newSchedule(
        address _beneficiary,
        uint256 _amount,
        uint256 _startBlock,
        uint256 _cliffBlock,
        uint256 _endBlock,
        bool _canStake
    ) virtual onlyOperator public {
        require(_cliffBlock >= _startBlock, "cliff must be greater to start");
        require(_endBlock > _cliffBlock, "end must be greater than cliff");

        bool _transferred = autonity.transferFrom(operator, address(this), _amount);
        require(_transferred, "amount not approved");

        uint256 _scheduleID = schedules.length;
        uint256 _vestingID = _newVesting(_amount, _startBlock, _cliffBlock, _endBlock);
        schedules.push(Schedule(_amount, _startBlock, _cliffBlock, _endBlock, _vestingID, _canStake, false));
        beneficiarySchedules[_beneficiary].push(_scheduleID);
    }


    /** @notice used by beneficiary to transfer all unlocked NTN and LNTN of some schedule to his own address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     */
    function releaseFunds(uint256 _id) virtual public onlyActive(_id) {
        releaseAllLNTN(_id);
        releaseAllNTN(_id);
    }

    /** @notice used by beneficiary to transfer all unlocked LNTN of some schedule to his own address
     */
    function releaseAllLNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        address[] memory _validators = _bondedValidators(_scheduleID);
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _amount = _decreaseAllUnlocked(liquidVestingIDs[_scheduleID][_validators[i]]);
            if (_amount > 0) {
                _transferLNTN(_scheduleID, msg.sender, _amount, _validators[i]);
            }
        }
    }

    /** @notice used by beneficiary to transfer all unlocked NTN of some schedule to his own address
     */
    function releaseAllNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.cliff < block.number, "not reached cliff period yet");
        uint256 _amount = _decreaseAllUnlocked(_schedule.vestingID);
        if (_amount > 0) {
            _transferNTN(_scheduleID, msg.sender, _amount);
        }
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /** @notice used by beneficiary to transfer some amount of unlocked NTN of some schedule to his own address
     * @param _amount amount to transfer
     */
    function releaseNTN(uint256 _id, uint256 _amount) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.cliff < block.number, "not reached cliff period yet");
        _decreaseUnlocked(_schedule.vestingID, _amount);
        _transferNTN(_scheduleID, msg.sender, _amount);
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /** @notice used by beneficiary to transfer some amount of unlocked LNTN of some schedule to his own address
     * @param _validator address of the validator
     * @param _amount amount to transfer
     */
    function releaseLNTN(uint256 _id, address _validator, uint256 _amount) virtual public onlyActive(_id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        _decreaseUnlocked(liquidVestingIDs[_scheduleID][_validator], _amount);
        _transferLNTN(_scheduleID, msg.sender, _amount, _validator);
    }

    /** @notice force release of all funds, NTN and LNTN, of some schedule and return them to the _recipient account
     * effectively cancelling a vesting schedule. only operator is able to call the function
     * rewards (AUT) which have been entitled to a schedule due to bonding are not returned to _recipient
     * @param _beneficiary beneficiary address whose schedule will be canceled
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     * @param _recipient to whome the all funds will be transferred
     */
    function cancelSchedule(address _beneficiary, uint256 _id, address _recipient) virtual public onlyOperator {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        Schedule storage _item = schedules[_scheduleID];
        _transferNTN(_scheduleID, _recipient, _item.totalAmount);
        _removeVesting(_item.vestingID);
        address[] memory _validators = _bondedValidators(_scheduleID);
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _amount = _unlockedLiquidBalanceOf(_scheduleID, _validators[i]);
            if (_amount > 0) {
                _transferLNTN(_scheduleID, _recipient, _amount, _validators[i]);
            }
            _removeVesting(liquidVestingIDs[_scheduleID][_validators[i]]);
            delete liquidVestingIDs[_scheduleID][_validators[i]];
        }
        _item.canceled = true;
        cancelRecipient[_scheduleID] = _recipient;
    }

    /** @notice Used by beneficiary to bond some NTN of some schedule. only schedules with canStake = true can be staked. Use canStake function
     * to check if can be staked or not. Bonded from the vesting manager contract. All bondings are delegated, as vesting manager cannot own a validator
     * @param _id id of the schedule numbered from 0 to (n-1) where n = total schedules entitled to the beneficiary (including canceled ones)
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(uint256 _id, address _validator, uint256 _amount) virtual public onlyActive(_id) returns (uint256) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.canStake, "not able to stake");
        require(_schedule.totalAmount >= _amount, "not enough tokens");

        uint256 _bondingID = autonity.bond(_validator, _amount);
        bondingToSchedule[_bondingID] = _scheduleID+1;
        pendingBondingRequest[_bondingID] = PendingBondingRequest(_amount, _validator);
        _schedule.totalAmount -= _amount;
        pendingBondingVesting[_bondingID] = _createDelegation(_schedule.vestingID, _amount);
        return _bondingID;
    }

    /** @notice Used by beneficiary to unbond some LNTN of some schedule.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public onlyActive(_id) returns (uint256) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        require(
            _unlockedLiquidBalanceOf(_scheduleID, _validator) >= _amount,
            "not enough unlocked liquid tokens"
        );
        uint256 _unbondingID = autonity.unbond(_validator, _amount);
        pendingUnbondingRequest[_unbondingID] = PendingUnbondingRequest(_amount, _validator);
        unbondingToSchedule[_unbondingID] = _scheduleID+1;
        _lock(_scheduleID, _validator, _amount);
        return _unbondingID;
    }

    /** @notice used by beneficiary to claim all rewards (AUT) which is entitled due to bonding */
    function claimRewards() virtual external {
        uint256[] storage _scheduleIDs = beneficiarySchedules[msg.sender];
        uint256 _totalFees = 0;
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            _totalFees += _rewards(_scheduleIDs[i]);
        }
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool _sent, ) = msg.sender.call{value: _totalFees}("");
        require(_sent, "Failed to send AUT");
    }

    /** @notice can be used to send AUT to the contract */
    function receiveAut() external payable {
        // do nothing
    }

    /** @dev implements IStakeProxy.bondingApplied(). a callback function for autonity when bonding is applied
     * @param _bondingID bonding id from Autonity when bonding is requested
     * @param _liquid amount of LNTN for bonding
     * @param _selfDelegation true if self bonded, false for delegated bonding
     * @param _rejected true if bonding request was rejected
     */
    function bondingApplied(uint256 _bondingID, uint256 _liquid, bool _selfDelegation, bool _rejected) external onlyAutonity {
        _applyBonding(_bondingID, _liquid, _selfDelegation, _rejected);
    }

    /** @dev implements IStakeProxy.unbondingApplied(). callback function for autonity when unbonding is applied
     * @param _unbondingID unbonding id from Autonity when unbonding is requested
     */
    function unbondingApplied(uint256 _unbondingID) external onlyAutonity {
        _applyUnbonding(_unbondingID);
    }

    /** @dev implements IStakeProxy.unbondingReleased(). callback function for autonity when unbonding is released
     * @param _unbondingID unbonding id from Autonity when unbonding is requested
     * @param _amount amount of NTN released
     */
    function unbondingReleased(uint256 _unbondingID, uint256 _amount) external onlyAutonity {
        _releaseUnbonding(_unbondingID, _amount);
    }

    /** @dev returns a unique id for each schedule
     * @param _beneficiary address of the schedule holder
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     */
    function _getUniqueScheduleID(address _beneficiary, uint256 _id) private view returns (uint256) {
        require(beneficiarySchedules[_beneficiary].length > _id, "invalid schedule id");
        return beneficiarySchedules[_beneficiary][_id];
    }

    function _transferNTN(uint256 _scheduleID, address _to, uint256 _amount) private {
        bool _sent = autonity.transfer(_to, _amount);
        require(_sent, "NTN not transferred");
        schedules[_scheduleID].totalAmount -= _amount;
    }

    function _transferLNTN(uint256 _scheduleID, address _to, uint256 _amount, address _validator) private {
        Liquid _liquidContract = autonity.getValidator(_validator).liquidContract;
        bool _sent = _liquidContract.transfer(_to, _amount);
        require(_sent, "LNTN transfer failed");
        _decreaseLiquid(_scheduleID, _validator, _amount);
    }

    function _applyBonding(uint256 _bondingID, uint256 _liquid, bool _selfDelegation, bool _rejected) internal {
        require(_selfDelegation == false, "bonding should be delegated");
        require(bondingToSchedule[_bondingID] > 0, "invalid bonding id");
        uint256 _scheduleID = bondingToSchedule[_bondingID]-1;
        Schedule storage _schedule = schedules[_scheduleID];
        if (_rejected) {
            uint256 _amount = pendingBondingRequest[_bondingID].amount;
            _schedule.totalAmount += _amount;
            _schedule.vestingID = _mergeVesting(_schedule.vestingID, pendingBondingVesting[_bondingID]);
        }
        else {
            address _validator = pendingBondingRequest[_bondingID].validator;
            _increaseLiquid(_scheduleID, _validator, _liquid);
            _updateVesting(pendingBondingVesting[_bondingID], _liquid);
            liquidVestingIDs[_scheduleID][_validator]
                = _mergeVesting(liquidVestingIDs[_scheduleID][_validator], pendingBondingVesting[_bondingID]);
        }
        if (_schedule.canceled) {
            if (_rejected) {
                _transferNTN(_scheduleID, cancelRecipient[_scheduleID], pendingBondingRequest[_bondingID].amount);
                _removeVesting(_schedule.vestingID);
            }
            else {
                address _validator = pendingBondingRequest[_bondingID].validator;
                _transferLNTN(_scheduleID, cancelRecipient[_scheduleID], _liquid, _validator);
                _removeVesting(liquidVestingIDs[_scheduleID][_validator]);
            }
        }
        delete pendingBondingVesting[_bondingID];
        delete pendingBondingRequest[_bondingID];
        delete bondingToSchedule[_bondingID];
    }

    function _applyUnbonding(uint256 _unbondingID) internal {
        require(unbondingToSchedule[_unbondingID] > 0, "invalid unbonding id");
        uint256 _scheduleID = unbondingToSchedule[_unbondingID]-1;
        PendingUnbondingRequest storage _unbondingRequst = pendingUnbondingRequest[_unbondingID];
        _unlock(_scheduleID, _unbondingRequst.validator, _unbondingRequst.amount);
        _decreaseLiquid(_scheduleID, _unbondingRequst.validator, _unbondingRequst.amount);
        if (schedules[_scheduleID].canceled == false) {
            pendingUnbondingVesting[_unbondingID]
                = _createDelegation(liquidVestingIDs[_scheduleID][_unbondingRequst.validator], _unbondingRequst.amount);
        }
        delete pendingUnbondingRequest[_unbondingID];
    }

    function _releaseUnbonding(uint256 _unbondingID, uint256 _amount) internal {
        require(unbondingToSchedule[_unbondingID] > 0, "invalid unbonding id");
        uint256 _scheduleID = unbondingToSchedule[_unbondingID]-1;
        Schedule storage _item = schedules[_scheduleID];
        _item.totalAmount += _amount;
        if (_item.canceled && _amount > 0) {
            _transferNTN(_scheduleID, cancelRecipient[_scheduleID], _amount);
        }
        if (_item.canceled) {
            _removeVesting(pendingUnbondingVesting[_unbondingID]);
        }
        else {
            _updateVesting(pendingUnbondingVesting[_unbondingID], _amount);
            _item.vestingID = _mergeVesting(_item.vestingID, pendingUnbondingVesting[_unbondingID]);
        }
        delete unbondingToSchedule[_unbondingID];
        delete pendingUnbondingVesting[_unbondingID];
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /** @notice returns a schedule entitled to _beneficiary
     * @param _beneficiary beneficiary address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     */
    function getSchedule(address _beneficiary, uint256 _id) virtual external view returns (Schedule memory) {
        return schedules[_getUniqueScheduleID(_beneficiary, _id)];
    }

    /** @notice returns if beneficiary can stake from his schedule
     * @param _beneficiary beneficiary address
     */
    function canStake(address _beneficiary, uint256 _id) virtual external view returns (bool) {
        return schedules[_getUniqueScheduleID(_beneficiary, _id)].canStake;
    }

    /** @notice returns the number of schudeled entitled to some beneficiary
     * @param _beneficiary address of the beneficiary
     */
    function totalSchedules(address _beneficiary) virtual external view returns (uint256) {
        return beneficiarySchedules[_beneficiary].length;
    }

    /** @notice returns the list of current schedules assigned to a beneficiary
     * @param _beneficiary address of the beneficiary
     */
    function getSchedules(address _beneficiary) virtual external view returns (Schedule[] memory) {
        uint256[] storage _scheduleIDs = beneficiarySchedules[_beneficiary];
        Schedule[] memory _res = new Schedule[](_scheduleIDs.length);
        for (uint256 i = 0; i < _res.length; i++) {
            _res[i] = schedules[_scheduleIDs[i]];
        }
        return _res;
    }

    /** @notice returns the amount of all unclaimed rewards (AUT) due to all the bonding from schedules entitled to beneficiary
     * @param _beneficiary beneficiary address
     */
    function unclaimedRewards(address _beneficiary) virtual external returns (uint256) {
        uint256 _totalFee = 0;
        uint256[] storage _scheduleIDs = beneficiarySchedules[_beneficiary];
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            _totalFee += _unclaimedRewards(_scheduleIDs[i]);
        }
        return _totalFee;
    }

    /** @notice returns the amount of LNTN for some schedule
     * @param _beneficiary beneficiary address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (including canceled ones)
     * @param _validator validator address
     */
    function liquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        return _liquidBalanceOf(_scheduleID, _validator);
    }

    /** @notice returns the amount of locked LNTN for some schedule */
    function lockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        return _lockedLiquidBalanceOf(_scheduleID, _validator);
    }

    /** @notice returns the amount of unlocked LNTN for some schedule */
    function unlockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        return _unlockedLiquidBalanceOf(_scheduleID, _validator);
    }

    /** @notice returns the list of validators bonded some schedule */
    function getBondedValidators(address _beneficiary, uint256 _id) virtual external view returns (address[] memory) {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        return _bondedValidators(_scheduleID);
    }

    /** @notice returns the amount of NTN released for some schedule */
    function releasedNTN(address _beneficiary, uint256 _id) virtual external view returns (uint256) {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        return releasedFunds(schedules[_scheduleID].vestingID);
    }

    /** @notice returns the amount of LNTN released for some schedule that bonded some funds */
    function releasedLNTN(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        return releasedFunds(liquidVestingIDs[_scheduleID][_validator]);
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

    modifier onlyActive(uint256 _id) {
        uint256 _scheduleID = _getUniqueScheduleID(msg.sender, _id);
        require(schedules[_scheduleID].canceled == false, "schedule canceled");
        _;
    }


}
