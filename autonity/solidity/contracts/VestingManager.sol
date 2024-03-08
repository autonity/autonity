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
    }

    mapping(address => uint256[]) internal addressSchedules;
    Schedule[] internal schedules;

    struct LiquidInfo {
        uint256 totalLiquid;
        uint256 lastUnrealisedFeeFactor;
    }

    mapping(address => LiquidInfo) private validatorLiquids;

    // stores the array of validators bonded to a schedule
    mapping(uint256 => address[]) private bondedValidators;
    // stores the (index+1) of validator in bondedValidators[id] array
    mapping(uint256 => mapping(address => uint256)) validatorIdx;

    mapping(uint256 => mapping(address => uint256)) private liquidBalances;
    mapping(uint256 => mapping(address => uint256)) private lockedLiquidBalances;
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
    mapping(uint256 => uint256) private bondingToSchedule;
    mapping(uint256 => uint256) private unbondingToSchedule;

    uint256 private epochID;
    uint256 private epochFetchedBlock;
    mapping(address => uint256) rewardsClaimedEpoch;

    Autonity internal autonity;
    address operator;

    constructor(address payable _autonity, address _operator){
        // save autonity and operator account  - with standard modifers
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
        schedules.push(Schedule(_beneficiary, _amount, 0, _startBlock, _cliffBlock, _endBlock, _stackable));
        uint256[] storage addressScheduless = addressSchedules[_beneficiary];
        addressScheduless.push(scheduleID);
    }

    // retrieve list of current schedules assigned to a beneficiary
    function getSchedules(address _beneficiary) virtual public returns (Schedule[] memory) {
        uint256[] storage scheduleIDs = addressSchedules[_beneficiary];
        Schedule[] memory res = new Schedule[](scheduleIDs.length);
        for(uint256 i = 0; i < res.length; i++) {
            res[i] = schedules[scheduleIDs[i]];
        }
        return res;
    }

    // used by beneficiary to transfer unlocked ntn
    function releaseFunds(uint256 _id) virtual public {
        // not only unlocked token but unlocked LNTN too !!!
        // release unlocked LNTN -> unbond?
    }

    function releaseUnlockedFunds(uint256 _id) public {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        Schedule storage schedule = schedules[scheduleID];
        require(schedule.cliff < block.number, "not reached cliff period yet");
        uint256 amount =
            schedule.totalAmount * (block.number-schedule.cliff) / (schedule.end-schedule.cliff) - schedule.withdrawnAmount;
        bool sent = autonity.transfer(msg.sender, amount);
        require(sent, "NTN not transfered");
        schedule.withdrawnAmount += amount;
    }

    // force release of all funds and return them to the _recipient account
    // effectively cancelling a vesting schedule
    // - target is beneficiary
    // - flag to retrieve as well untransfered unlocked token
    function cancelSchedule(address _target, uint256 _id, address _recipient) virtual public onlyOperator {

    }

    // ONLY APPLY WITH STACKABLE SCHEDULE
    // Q : can we bond more than LOCKED with remaining UNLOCKED
    // 3 locked , 2 unlocked : can you bond with your 5?
    function bond(uint256 _id, address _validator, uint256 _amount) virtual public {
        uint256 scheduleID = _getScheduleID(msg.sender, _id);
        Schedule storage schedule = schedules[scheduleID];
        require(schedule.stackable, "not stackable");
        require(schedule.totalAmount - schedule.withdrawnAmount >= _amount, "not enough tokens");

        uint256 bondingID = autonity.getHeadBondingID();
        autonity.bond(_validator, _amount);
        bondingToSchedule[bondingID] = scheduleID+1;
        pendingBondingRequest[bondingID] = PendingBondingRequest(_amount, _validator);
        schedule.totalAmount -= _amount;
    }

    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public {
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

    function totalSchedules() public view returns (uint256) {
        return addressSchedules[msg.sender].length;
    }

    function claimAllRewards() external {
        uint256[] storage scheduleIDs = addressSchedules[msg.sender];
        uint256 totalFees = 0;
        for(uint256 i = 0; i < scheduleIDs.length; i++) {
            totalFees += _computeRewards(scheduleIDs[i]);
        }
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool sent, ) = msg.sender.call{value: totalFees}("");
        require(sent, "Failed to send AUT");
    }

    function claimRewards(uint256 _id) external {
        uint256 totalFees = _computeRewards(_getScheduleID(msg.sender, _id));
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool sent, ) = msg.sender.call{value: totalFees}("");
        require(sent, "Failed to send AUT");
    }

    function bondingApplied(uint256 _bondingID, uint256 _liquid, bool _rejected) public onlyAutonity {
        require(bondingToSchedule[_bondingID] > 0, "invalid bonding id");
        uint256 scheduleID = bondingToSchedule[_bondingID]-1;
        if (_rejected) {
            uint256 amount = pendingBondingRequest[_bondingID].amount;
            Schedule storage schedule = schedules[scheduleID];
            schedule.totalAmount += amount;
        }
        else {
            _increaseLiquid(scheduleID, pendingBondingRequest[_bondingID].validator, _liquid);
        }
        delete pendingBondingRequest[_bondingID];
        delete bondingToSchedule[_bondingID];
    }

    function unbondingApplied(uint256 _unbondingID, bool _rejected) public onlyAutonity {
        require(unbondingToSchedule[_unbondingID] > 0, "invalid unbonding id");
        uint256 scheduleID = unbondingToSchedule[_unbondingID]-1;
        PendingUnbondingRequest storage unbondingRequst = pendingUnbondingRequest[_unbondingID];
        lockedLiquidBalances[scheduleID][unbondingRequst.validator] -= unbondingRequst.amount;
        if (!_rejected) {
            _decreaseLiquid(scheduleID, unbondingRequst.validator, unbondingRequst.amount);
        }
        else {
            delete unbondingToSchedule[_unbondingID];
        }
        delete pendingUnbondingRequest[_unbondingID];
    }

    function unbondingReleased(uint256 _unbondingID, uint256 _amount) public onlyAutonity {
        require(unbondingToSchedule[_unbondingID] > 0, "invalid unbonding id");
        uint256 scheduleID = unbondingToSchedule[_unbondingID]-1;
        schedules[scheduleID].totalAmount += _amount;
        delete unbondingToSchedule[_unbondingID];
    }

    function _getScheduleID(address _beneficiary, uint256 _id) private view returns (uint256) {
        require(addressSchedules[_beneficiary].length > _id, "invalid schedule id");
        return addressSchedules[_beneficiary][_id];
    }

    function _decreaseLiquid(uint256 _scheduleID, address _validator, uint256 _amount) private {
        require(
            liquidBalances[_scheduleID][_validator] - lockedLiquidBalances[_scheduleID][_validator] >= _amount,
            "not enough unlocked liquid tokens"
        );

        _realiseFees(_scheduleID, _validator);
        liquidBalances[_scheduleID][_validator] -= _amount;
        validatorLiquids[_validator].totalLiquid -= _amount;
        if(liquidBalances[_scheduleID][_validator] == 0) {
            _removeValidator(_scheduleID, _validator);
            delete unrealisedFeeFactors[_scheduleID][_validator];
        }
    }

    function _increaseLiquid(uint256 _scheduleID, address _validator, uint256 _amount) private {
        _realiseFees(_scheduleID, _validator);
        if(liquidBalances[_scheduleID][_validator] == 0) {
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
        if(rewardsClaimedEpoch[_validator] == _epochID()) {
            return;
        }
        Autonity.Validator memory validator = autonity.getValidator(_validator);
        uint256 reward = address(this).balance;
        validator.liquidContract.claimRewards();
        reward = address(this).balance - reward;
        if(reward > 0) {
            LiquidInfo storage liquidInfo = validatorLiquids[_validator];
            require(liquidInfo.totalLiquid > 0, "got reward from validator with no liquid supply"); // this shouldn't happen
            liquidInfo.lastUnrealisedFeeFactor += (reward * FEE_FACTOR_UNIT_RECIP) / liquidInfo.totalLiquid;
        }
        rewardsClaimedEpoch[_validator] = _epochID();
    }

    function _computeRewards(uint256 _scheduleID) private returns (uint256) {
        address[] storage validators = bondedValidators[_scheduleID];
        uint256 totalFees = 0;
        for(uint256 i = 0; i < validators.length; i++) {
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

        if(idx < validators.length) {
            validatorIdx[_scheduleID][validators[idx]] = idx+1;
        }
    }

    function _epochID() private returns (uint256) {
        if(epochFetchedBlock < block.number) {
            epochFetchedBlock = block.number;
            epochID = autonity.epochID();
        }
        return epochID;
    }

    // add std modifiers here
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


}
