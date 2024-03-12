// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./Autonity.sol";
import "./Liquid.sol";

contract VestingManager {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release schedule
    // UNLOCKED are tokens that got released but not yet transferred
    uint256 public contractVersion = 1;

    uint256 public constant FEE_FACTOR_UNIT_RECIP = 1_000_000_000;

    struct Schedule {
        address beneficiary;
        uint256 totalAmount;
        uint256 withdrawnAmount;
        uint256 start;
        uint256 cliff;
        uint256 end; // or duration?
        bool stackable;
        bool canceled;
    }

    // stores the unique ids of schedules assigned to a beneficiary, but beneficiary does not need to know the id
    // beneficiary will number his schedules as: 0 for first schedule, 1 for 2nd and so on
    // we can get the unique schedule id from addressSchedules as follows
    // addressSchedules[beneficiary][0] is the unique id of his first schedule
    // addressSchedules[beneficiary][1] is the unique id of his 2nd schedule and so on
    mapping(address => uint256[]) internal addressSchedules;

    // list of all schedules
    Schedule[] internal schedules;

    struct LiquidInfo {
        uint256 totalLiquid;
        uint256 lastUnrealisedFeeFactor;
    }

    // stores total liquid and lastUnrealisedFeeFactor for each validator
    // lastUnrealisedFeeFactor is used to calculate unrealised rewards for schedules with the same logic as done in Liquid.sol
    mapping(address => LiquidInfo) private validatorLiquids;

    // stores the array of validators bonded to a schedule
    mapping(uint256 => address[]) private bondedValidators;
    // stores the (index+1) of validator in bondedValidators[id] array
    mapping(uint256 => mapping(address => uint256)) validatorIdx;

    mapping(uint256 => mapping(address => uint256)) private liquidBalances;
    mapping(uint256 => mapping(address => uint256)) private lockedLiquidBalances;
    mapping(uint256 => mapping(address => uint256)) private withdrawnLiquid;

    // realisedFees[id][validator] stores the realised reward entitled to a schedule for a validator
    // unrealisedFeeFactors[id][validator] is used to calculate unrealised rewards. it only updates
    // when the liquid balance of the schedule is updated following the same logic done in Liquid.sol
    mapping(uint256 => mapping(address => uint256)) private realisedFees;
    mapping(uint256 => mapping(address => uint256)) private unrealisedFeeFactors;

    struct PendingBondingRequest {
        uint256 amount;
        address validator;
    }

    struct PendingUnbondingRequest {
        uint256 amount;
        address validator;
    }

    mapping(uint256 => PendingBondingRequest) private pendingBondingRequest;
    mapping(uint256 => PendingUnbondingRequest) private pendingUnbondingRequest;

    // bondingToSchedule[bondingID] stores the unique schedule id which requested the bonding
    mapping(uint256 => uint256) private bondingToSchedule;

    // unbondingToSchedule[unbondingID] stores the unique schedule id which requested the unbonding
    mapping(uint256 => uint256) private unbondingToSchedule;

    mapping(uint256 => address) private cancelRecipient;

    uint256 private epochID;
    uint256 private epochFetchedBlock;

    // rewardsClaimedEpoch[validator] stores the last epoch where the rewards from validator were claimed
    mapping(address => uint256) rewardsClaimedEpoch;

    Autonity internal autonity;
    address operator;

    constructor(address payable _autonity, address _operator){
        autonity = Autonity(_autonity);
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
        schedules.push(Schedule(_beneficiary, _amount, 0, _startBlock, _cliffBlock, _endBlock, _stackable, false));
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
        address[] storage validators = bondedValidators[scheduleID];
        for (uint256 i = 0; i < validators.length; i++) {
            address validator = validators[i];
            uint256 amount = _releasedLNTN(scheduleID, validator);
            if (amount > 0) {
                _transferLNTN(scheduleID, msg.sender, amount, validator);
                withdrawnLiquid[scheduleID][validator] += amount;
            }
        }
    }

    // used by beneficiary to transfer unlocked NTN
    function releaseNTN(uint256 _id) virtual public onlyActive(_id) {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        Schedule storage schedule = schedules[scheduleID];
        require(schedule.cliff < block.number, "not reached cliff period yet");
        uint256 amount = _releasedNTN(scheduleID);
        if (amount > 0) {
            _transferNTN(scheduleID, msg.sender, amount);
            schedule.withdrawnAmount += amount;
        }
    }

    // force release of all funds, NTN and LNTN, and return them to the _recipient account
    // effectively cancelling a vesting schedule
    // rewards (AUT) which have been entitled to a schedule due to bonding are not returned to _recipient
    function cancelSchedule(address _beneficiary, uint256 _id, address _recipient) virtual public onlyOperator {
        uint256 scheduleID = _getScheduleID(_beneficiary, _id);
        Schedule storage item = schedules[scheduleID];
        _transferNTN(scheduleID, _recipient, item.totalAmount);
        address[] storage validators = bondedValidators[scheduleID];
        for (uint256 i = 0; i < validators.length; i++) {
            address validator = validators[i];
            uint256 amount = liquidBalances[scheduleID][validator] - lockedLiquidBalances[scheduleID][validator];
            if (amount > 0) {
                _transferLNTN(scheduleID, _recipient, amount, validator);
            }
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
    }

    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public onlyActive(_id) {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        require(
            liquidBalances[scheduleID][_validator] - lockedLiquidBalances[scheduleID][_validator] >= _amount,
            "not enough unlocked liquid tokens"
        );
        uint256 unbondingID = autonity.getHeadUnbondingID();
        autonity.unbond(_validator, _amount);
        pendingUnbondingRequest[unbondingID] = PendingUnbondingRequest(_amount, _validator);
        unbondingToSchedule[unbondingID] = scheduleID+1;
        lockedLiquidBalances[scheduleID][_validator] += _amount;
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
            }
        }
        else {
            _increaseLiquid(scheduleID, pendingBondingRequest[_bondingID].validator, _liquid);
            if (schedule.canceled) {
                _transferLNTN(scheduleID, cancelRecipient[scheduleID], _liquid, pendingBondingRequest[_bondingID].validator);
            }
        }
        delete pendingBondingRequest[_bondingID];
        delete bondingToSchedule[_bondingID];
    }

    // callback function for autonity when unbonding is applied
    function unbondingApplied(uint256 _unbondingID) public onlyAutonity {
        require(unbondingToSchedule[_unbondingID] > 0, "invalid unbonding id");
        uint256 scheduleID = unbondingToSchedule[_unbondingID]-1;
        PendingUnbondingRequest memory unbondingRequst = pendingUnbondingRequest[_unbondingID];
        lockedLiquidBalances[scheduleID][unbondingRequst.validator] -= unbondingRequst.amount;
        _decreaseLiquid(scheduleID, unbondingRequst.validator, unbondingRequst.amount);
        delete pendingUnbondingRequest[_unbondingID];
    }

    // callback function for autonity when unbonding is released
    function unbondingReleased(uint256 _unbondingID, uint256 _amount) public onlyAutonity {
        require(unbondingToSchedule[_unbondingID] > 0, "invalid unbonding id");
        uint256 scheduleID = unbondingToSchedule[_unbondingID]-1;
        schedules[scheduleID].totalAmount += _amount;
        if (schedules[scheduleID].canceled && _amount > 0) {
            _transferNTN(scheduleID, cancelRecipient[scheduleID], _amount);
        }
        delete unbondingToSchedule[_unbondingID];
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

    /**
     * @dev amount of released NTN and LNTN follows a linear function from cliff period to end period
     * @param _scheduleID unique id of a schedule, for internal purpose only, beneficiary has no need of this id
     */
    function _releasedNTN(uint256 _scheduleID) private view returns (uint256) {
        Schedule storage item = schedules[_scheduleID];
        if (item.end <= block.number) {
            return item.totalAmount;
        }
        uint256 unlocked = (item.totalAmount+item.withdrawnAmount) * (block.number-item.cliff) / (item.end-item.cliff);
        if (unlocked > item.withdrawnAmount) {
            return unlocked - item.withdrawnAmount;
        }
        return 0;
    }

    function _releasedLNTN(uint256 _scheduleID, address _validator) private view returns (uint256) {
        Schedule storage item = schedules[_scheduleID];
        if (item.end <= block.number) {
            return liquidBalances[_scheduleID][_validator] - lockedLiquidBalances[_scheduleID][_validator];
        }
        uint256 withdrawn = withdrawnLiquid[_scheduleID][_validator];
        uint256 unlocked = (liquidBalances[_scheduleID][_validator]+withdrawn) * (block.number-item.cliff) / (item.end-item.cliff);
        if (unlocked > withdrawn) {
            uint256 available = liquidBalances[_scheduleID][_validator] - lockedLiquidBalances[_scheduleID][_validator];
            if (available > unlocked - withdrawn) {
                return unlocked - withdrawn;
            }
            return available;
        }
        return 0;
    }

    /**
     * @dev _decreaseLiquid, _increaseLiquid, _realiseFees, _computeUnrealisedFees follows the same logic as done in Liquid.sol
     * the only difference is in _realiseFees, we claim rewards from the validator first, because vesting manager does not know
     * when the epoch ends, and so cannot claim rewards at each epoch end.
     * _claimRewards claim rewards from a validator only once per epoch, so spamming _claimRewards is not a problem
     */
    function _decreaseLiquid(uint256 _scheduleID, address _validator, uint256 _amount) private {
        require(
            liquidBalances[_scheduleID][_validator] - lockedLiquidBalances[_scheduleID][_validator] >= _amount,
            "not enough unlocked liquid tokens"
        );

        _realiseFees(_scheduleID, _validator);
        liquidBalances[_scheduleID][_validator] -= _amount;
        validatorLiquids[_validator].totalLiquid -= _amount;
        if (liquidBalances[_scheduleID][_validator] == 0) {
            _removeValidator(_scheduleID, _validator);
            delete unrealisedFeeFactors[_scheduleID][_validator];
        }
    }

    function _increaseLiquid(uint256 _scheduleID, address _validator, uint256 _amount) private {
        _realiseFees(_scheduleID, _validator);
        if (liquidBalances[_scheduleID][_validator] == 0) {
            _addValidator(_scheduleID, _validator);
        }
        liquidBalances[_scheduleID][_validator] += _amount;
        validatorLiquids[_validator].totalLiquid += _amount;
    }

    function _realiseFees(uint256 _scheduleID, address _validator) private returns (uint256 _realisedFees) {
        _claimRewards(_validator);
        uint256 _unrealisedFees = _computeUnrealisedFees(_scheduleID, _validator);
        _realisedFees = realisedFees[_scheduleID][_validator] + _unrealisedFees;
        realisedFees[_scheduleID][_validator] = _realisedFees;
        unrealisedFeeFactors[_scheduleID][_validator] = validatorLiquids[_validator].lastUnrealisedFeeFactor;
    }

    function _computeUnrealisedFees(uint256 _scheduleID, address _validator) private view returns (uint256) {
        uint256 balance = liquidBalances[_scheduleID][_validator];
        if (balance == 0) {
            return 0;
        }
        uint256 _unrealisedFeeFactor =
            validatorLiquids[_validator].lastUnrealisedFeeFactor - unrealisedFeeFactors[_scheduleID][_validator];
        uint256 _unrealisedFee = (_unrealisedFeeFactor * balance) / FEE_FACTOR_UNIT_RECIP;
        return _unrealisedFee;
    }

    function _claimRewards(address _validator) private {
        if (rewardsClaimedEpoch[_validator] == _epochID()) {
            return;
        }
        Liquid liquidContract = autonity.getValidator(_validator).liquidContract;
        uint256 reward = address(this).balance;
        liquidContract.claimRewards();
        reward = address(this).balance - reward;
        if (reward > 0) {
            LiquidInfo storage liquidInfo = validatorLiquids[_validator];
            require(liquidInfo.totalLiquid > 0, "got reward from validator with no liquid supply"); // this shouldn't happen
            liquidInfo.lastUnrealisedFeeFactor += (reward * FEE_FACTOR_UNIT_RECIP) / liquidInfo.totalLiquid;
        }
        rewardsClaimedEpoch[_validator] = _epochID();
    }

    /**
     * @dev call _rewards(_scheduleID) only when rewards are claimed
     * calculates total rewards for a schedule and deletes realisedFees[id][validator] as reward is claimed
     */ 
    function _rewards(uint256 _scheduleID) private returns (uint256) {
        address[] storage validators = bondedValidators[_scheduleID];
        uint256 totalFees = 0;
        for (uint256 i = 0; i < validators.length; i++) {
            totalFees += _realiseFees(_scheduleID, validators[i]);
            delete realisedFees[_scheduleID][validators[i]];
        }
        return totalFees;
    }

    function _addValidator(uint256 _scheduleID, address _validator) private {
        address[] storage validators = bondedValidators[_scheduleID];
        validators.push(_validator);
        validatorIdx[_scheduleID][_validator] = validators.length;
    }

    function _removeValidator(uint256 _scheduleID, address _validator) private {
        address[] storage validators = bondedValidators[_scheduleID];
        uint256 idx = validatorIdx[_scheduleID][_validator]-1;
        // removing _validator by replacing it with last one and then deleting the last one
        validators[idx] = validators[validators.length-1];
        validators.pop();
        delete validatorIdx[_scheduleID][_validator];

        if (idx < validators.length) {
            validatorIdx[_scheduleID][validators[idx]] = idx+1;
        }
    }

    function _epochID() private returns (uint256) {
        if (epochFetchedBlock < block.number) {
            epochFetchedBlock = block.number;
            epochID = autonity.epochID();
        }
        return epochID;
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
        uint256 length = addressSchedules[_account].length;
        for (uint256 i = 0; i < length; i++) {
            totalFee += unclaimedRewards(_account, i);
        }
        return totalFee;
    }

    function unclaimedRewards(address _account, uint256 _id) virtual public view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        uint256 totalFee = 0;
        address[] storage validators = bondedValidators[scheduleID];
        for (uint256 i = 0; i < validators.length; i++) {
            address validator = validators[i];
            totalFee += realisedFees[scheduleID][validator] + _computeUnrealisedFees(scheduleID, validator);
        }
        return totalFee;
    }

    function liquidBalanceOf(address _account, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return liquidBalances[scheduleID][_validator];
    }

    function unlockedLiquidBalanceOf(address _account, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return liquidBalances[scheduleID][_validator] - lockedLiquidBalances[scheduleID][_validator];
    }

    // returns the list of validator addresses wich are bonded to schedule _id assigned to _account
    function getBondedValidators(address _account, uint256 _id) external view returns (address[] memory) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return bondedValidators[scheduleID];
    }

    // amount of NTN released from schedule _id assigned to _account but not yet withdrawn by _account
    function releasedNTN(address _account, uint256 _id) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return _releasedNTN(scheduleID);
    }

    // amount of LNTN released from schedule _id assigned to _account but not yet withdrawn by _account
    function releasedLNTN(address _account, uint256 _id, address _validator) virtual external view returns (uint256) {
        uint256 scheduleID = _getScheduleID(_account, _id);
        return _releasedLNTN(scheduleID, _validator);
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
