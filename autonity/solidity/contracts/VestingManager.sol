// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./Autonity.sol";
import "./Liquid.sol";
import "./LiquidRewardManager.sol";
import "./VestingCalculator.sol";

contract VestingManager is LiquidRewardManager, VestingCalculator {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release schedule
    // UNLOCKED are tokens that got released but not yet transferred
    uint256 public contractVersion = 1;

    address private operator;

    struct Schedule {
        uint256 totalAmount;
        uint256 start;
        uint256 cliff;
        uint256 end; // or duration?
        uint256 vestingID;
        bool stackable;
        bool canceled;
    }

    mapping(uint256 => mapping(address => uint256)) private liquidVestingIDs;

    // stores the unique ids of schedules assigned to a beneficiary, but beneficiary does not need to know the id
    // beneficiary will number his schedules as: 0 for first schedule, 1 for 2nd and so on
    // we can get the unique schedule id from addressSchedules as follows
    // addressSchedules[beneficiary][0] is the unique id of his first schedule
    // addressSchedules[beneficiary][1] is the unique id of his 2nd schedule and so on
    mapping(address => uint256[]) private addressSchedules;

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

    // bondingToSchedule[bondingID] stores the unique schedule id which requested the bonding
    mapping(uint256 => uint256) private bondingToSchedule;

    // unbondingToSchedule[unbondingID] stores the unique schedule id which requested the unbonding
    mapping(uint256 => uint256) private unbondingToSchedule;

    mapping(uint256 => address) private cancelRecipient;

    constructor(address payable _autonity, address _operator) LiquidRewardManager(_autonity) {
        operator = _operator;
    }

    function newSchedule(
        address _beneficiary,
        uint256 _amount,
        uint256 _startBlock,
        uint256 _cliffBlock,
        uint256 _endBlock,
        bool _stackable
    ) virtual onlyOperator public {
        require(_cliffBlock >= _startBlock, "cliff must be greater to start");
        require(_endBlock > _cliffBlock, "end must be greater than cliff");

        bool transfered = autonity.transferFrom(operator, address(this), _amount);
        require(transfered, "amount not approved");

        uint256 scheduleID = schedules.length;
        uint256 vestingID = _newVesting(_amount, _cliffBlock, _endBlock);
        schedules.push(Schedule(_amount, _startBlock, _cliffBlock, _endBlock, vestingID, _stackable, false));
        uint256[] storage addressScheduless = addressSchedules[_beneficiary];
        addressScheduless.push(scheduleID);
    }


    // used by beneficiary to transfer unlocked NTN and LNTN
    function releaseFunds(uint256 _id) virtual public onlyActive(_id) {
        releaseNTN(_id);
        releaseLNTN(_id);
    }

    // used by beneficiary to transfer unlocked LNTN
    function releaseLNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        address[] memory validators = _bondedValidators(scheduleID);
        for (uint256 i = 0; i < validators.length; i++) {
            uint256 amount = _withdrawAll(liquidVestingIDs[scheduleID][validators[i]]);
            if (amount > 0) {
                _transferLNTN(scheduleID, msg.sender, amount, validators[i]);
            }
        }
    }

    // used by beneficiary to transfer unlocked NTN
    function releaseNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        Schedule storage schedule = schedules[scheduleID];
        require(schedule.cliff < block.number, "not reached cliff period yet");
        uint256 amount = _withdrawAll(schedule.vestingID);
        if (amount > 0) {
            _transferNTN(scheduleID, msg.sender, amount);
        }
    }

    function releaseNTN(uint256 _id, uint256 _amount) public onlyActive(_id) {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        Schedule storage schedule = schedules[scheduleID];
        require(schedule.cliff < block.number, "not reached cliff period yet");
        _withdraw(schedule.vestingID, _amount);
        _transferNTN(scheduleID, msg.sender, _amount);
    }

    function releaseLNTN(uint256 _id, address _validator, uint256 _amount) public onlyActive(_id) {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        _withdraw(liquidVestingIDs[scheduleID][_validator], _amount);
        _transferLNTN(scheduleID, msg.sender, _amount, _validator);
    }

    // force release of all funds, NTN and LNTN, and return them to the _recipient account
    // effectively cancelling a vesting schedule
    // rewards (AUT) which have been entitled to a schedule due to bonding are not returned to _recipient
    function cancelSchedule(address _beneficiary, uint256 _id, address _recipient) virtual public onlyOperator {
        uint256 scheduleID = _getScheduleID(_beneficiary, _id);
        Schedule storage item = schedules[scheduleID];
        _transferNTN(scheduleID, _recipient, item.totalAmount);
        _removeVesting(item.vestingID);
        address[] memory validators = _bondedValidators(scheduleID);
        for (uint256 i = 0; i < validators.length; i++) {
            uint256 amount = _unlockedLiquidBalanceOf(scheduleID, validators[i]);
            if (amount > 0) {
                _transferLNTN(scheduleID, _recipient, amount, validators[i]);
            }
            _removeVesting(liquidVestingIDs[scheduleID][validators[i]]);
            delete liquidVestingIDs[scheduleID][validators[i]];
        }
        item.canceled = true;
        cancelRecipient[scheduleID] = _recipient;
    }

    function _transferNTN(uint256 _scheduleID, address _to, uint256 _amount) private {
        bool sent = autonity.transfer(_to, _amount);
        require(sent, "NTN not transfered");
        Schedule storage schedule = schedules[_scheduleID];
        schedule.totalAmount -= _amount;
    }

    function _transferLNTN(uint256 _scheduleID, address _to, uint256 _amount, address _validator) private {
        Liquid liquidContract = autonity.getValidator(_validator).liquidContract;
        bool sent = liquidContract.transfer(_to, _amount);
        require(sent, "LNTN transfer failed");
        _decreaseLiquid(_scheduleID, _validator, _amount);
    }

    // ONLY APPLY WITH STACKABLE SCHEDULE
    // all bondings are delegated, as vesting manager cannot own a validator
    function bond(uint256 _id, address _validator, uint256 _amount) virtual public onlyActive(_id) {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        Schedule storage schedule = schedules[scheduleID];
        require(schedule.stackable, "not stackable");
        require(schedule.totalAmount >= _amount, "not enough tokens");

        uint256 bondingID = autonity.getHeadBondingID();
        autonity.bond(_validator, _amount);
        bondingToSchedule[bondingID] = scheduleID+1;
        pendingBondingRequest[bondingID] = PendingBondingRequest(_amount, _validator);
        schedule.totalAmount -= _amount;
        pendingBondingVesting[bondingID] = _splitVesting(schedule.vestingID, _amount);
    }

    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public onlyActive(_id) {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        require(
            _unlockedLiquidBalanceOf(scheduleID, _validator) >= _amount,
            "not enough unlocked liquid tokens"
        );
        uint256 unbondingID = autonity.getHeadUnbondingID();
        autonity.unbond(_validator, _amount);
        pendingUnbondingRequest[unbondingID] = PendingUnbondingRequest(_amount, _validator);
        unbondingToSchedule[unbondingID] = scheduleID+1;
        _lock(scheduleID, _validator, _amount);
    }

    function claimAllRewards() virtual external {
        uint256[] storage scheduleIDs = addressSchedules[msg.sender];
        uint256 totalFees = 0;
        for (uint256 i = 0; i < scheduleIDs.length; i++) {
            totalFees += _rewards(scheduleIDs[i]);
        }
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool sent, ) = msg.sender.call{value: totalFees}("");
        require(sent, "Failed to send AUT");
    }

    function claimRewards(uint256 _id) virtual external {
        uint256 totalFees = _rewards(_getScheduleID(msg.sender, _id));
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool sent, ) = msg.sender.call{value: totalFees}("");
        require(sent, "Failed to send AUT");
    }

    // callback function for autonity when bonding is applied
    function bondingApplied(uint256 _bondingID, uint256 _liquid, bool _rejected) public onlyAutonity {
        require(bondingToSchedule[_bondingID] > 0, "invalid bonding id");
        uint256 scheduleID = bondingToSchedule[_bondingID]-1;
        Schedule storage schedule = schedules[scheduleID];
        if (_rejected) {
            uint256 amount = pendingBondingRequest[_bondingID].amount;
            schedule.totalAmount += amount;
            if (schedule.canceled) {
                _transferNTN(scheduleID, cancelRecipient[scheduleID], amount);
                _removeVesting(pendingBondingVesting[_bondingID]);
            }
            else {
                schedule.vestingID = _mergeVesting(schedule.vestingID, pendingBondingVesting[_bondingID]);
            }
        }
        else {
            address validator = pendingBondingRequest[_bondingID].validator;
            _increaseLiquid(scheduleID, validator, _liquid);
            if (schedule.canceled) {
                _transferLNTN(scheduleID, cancelRecipient[scheduleID], _liquid, validator);
                _removeVesting(pendingBondingVesting[_bondingID]);
            }
            else {
                _updateVesting(pendingBondingVesting[_bondingID], _liquid);
                liquidVestingIDs[scheduleID][validator]
                    = _mergeVesting(liquidVestingIDs[scheduleID][validator], pendingBondingVesting[_bondingID]);
            }
        }
        delete pendingBondingVesting[_bondingID];
        delete pendingBondingRequest[_bondingID];
        delete bondingToSchedule[_bondingID];
    }

    // callback function for autonity when unbonding is applied
    function unbondingApplied(uint256 _unbondingID) public onlyAutonity {
        require(unbondingToSchedule[_unbondingID] > 0, "invalid unbonding id");
        uint256 scheduleID = unbondingToSchedule[_unbondingID]-1;
        PendingUnbondingRequest memory unbondingRequst = pendingUnbondingRequest[_unbondingID];
        _unlock(scheduleID, unbondingRequst.validator, unbondingRequst.amount);
        _decreaseLiquid(scheduleID, unbondingRequst.validator, unbondingRequst.amount);
        if (schedules[scheduleID].canceled == false) {
            pendingUnbondingVesting[_unbondingID]
                = _splitVesting(liquidVestingIDs[scheduleID][unbondingRequst.validator], unbondingRequst.amount);
        }
        delete pendingUnbondingRequest[_unbondingID];
    }

    // callback function for autonity when unbonding is released
    function unbondingReleased(uint256 _unbondingID, uint256 _amount) public onlyAutonity {
        require(unbondingToSchedule[_unbondingID] > 0, "invalid unbonding id");
        uint256 scheduleID = unbondingToSchedule[_unbondingID]-1;
        Schedule storage item = schedules[scheduleID];
        item.totalAmount += _amount;
        if (item.canceled && _amount > 0) {
            _transferNTN(scheduleID, cancelRecipient[scheduleID], _amount);
        }
        if (item.canceled) {
            _removeVesting(pendingUnbondingVesting[_unbondingID]);
        }
        else {
            _updateVesting(pendingUnbondingVesting[_unbondingID], _amount);
            item.vestingID = _mergeVesting(item.vestingID, pendingUnbondingVesting[_unbondingID]);
        }
        delete unbondingToSchedule[_unbondingID];
        delete pendingUnbondingVesting[_unbondingID];
    }

    /**
     * @dev returns a unique id for each schedule
     * @param _beneficiary address of the schedule holder
     * @param _id id of the schedule assigned to beneficiary numbered from 0 to (n-1) where n = total schedules assigned to beneficiary
     */
    function _getScheduleID(address _beneficiary, uint256 _id) private view returns (uint256) {
        require(addressSchedules[_beneficiary].length > _id, "invalid schedule id");
        return addressSchedules[_beneficiary][_id];
    }

    /*
    ============================================================
        Getters
    ============================================================
    */

   function totalSchedules() public view returns (uint256) {
        return addressSchedules[msg.sender].length;
    }

    // retrieve list of current schedules assigned to a beneficiary
    function getSchedules(address _beneficiary) virtual public view returns (Schedule[] memory) {
        uint256[] storage scheduleIDs = addressSchedules[_beneficiary];
        Schedule[] memory res = new Schedule[](scheduleIDs.length);
        for (uint256 i = 0; i < res.length; i++) {
            res[i] = schedules[scheduleIDs[i]];
        }
        return res;
    }

    // unclaimed rewards for all the schedules assigned to _account
    function unclaimedRewards(address _account) virtual external view returns (uint256) {
        uint256 totalFee = 0;
        uint256[] storage ids = addressSchedules[_account];
        for (uint256 i = 0; i < ids.length; i++) {
            totalFee += _unclaimedRewards(ids[i]);
        }
        return totalFee;
    }

    function unclaimedRewards(address _account, uint256 _id) virtual public view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return _unclaimedRewards(scheduleID);
    }

    function liquidBalanceOf(address _account, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return _liquidBalanceOf(scheduleID, _validator);
    }

    function lockedLiquidBalanceOf(address _account, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return _lockedLiquidBalanceOf(scheduleID, _validator);
    }

    function unlockedLiquidBalanceOf(address _account, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return _unlockedLiquidBalanceOf(scheduleID, _validator);
    }

    // returns the list of validator addresses wich are bonded to schedule _id assigned to _account
    function getBondedValidators(address _account, uint256 _id) external view returns (address[] memory) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return _bondedValidators(scheduleID);
    }

    // amount of NTN released from schedule _id assigned to _account but not yet withdrawn by _account
    function releasedNTN(address _account, uint256 _id) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return releasedFunds(schedules[scheduleID].vestingID);
    }

    // amount of LNTN released from schedule _id assigned to _account but not yet withdrawn by _account
    function releasedLNTN(address _account, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return releasedFunds(liquidVestingIDs[scheduleID][_validator]);
    }

    /*
    ============================================================

        Modifiers

    ============================================================
    */

    /**
    * @dev Modifier that checks if the caller is the governance operator account.
    * This should be abstracted by a separate smart-contract.
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
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        require(schedules[scheduleID].canceled == false, "schedule canceled");
        _;
    }


}
