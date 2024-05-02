// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";

contract ScheduleBase {

    struct Schedule {
        uint256 currentNTNAmount;
        uint256 withdrawnValue;
        uint256 start;
        uint256 cliff;
        uint256 end;
        bool canStake;
    }

    // stores the unique ids of schedules assigned to a beneficiary, but beneficiary does not need to know the id
    // beneficiary will number his schedules as: 0 for first schedule, 1 for 2nd and so on
    // we can get the unique schedule id from beneficiarySchedules as follows
    // beneficiarySchedules[beneficiary][0] is the unique id of his first schedule
    // beneficiarySchedules[beneficiary][1] is the unique id of his 2nd schedule and so on
    mapping(address => uint256[]) internal beneficiarySchedules;

    // list of all schedules
    Schedule[] internal schedules;

    Autonity internal autonity;
    address private operator;

    constructor(address payable _autonity, address _operator) {
        autonity = Autonity(_autonity);
        operator = _operator;
    }

    function _createSchedule(
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffTime,
        uint256 _endTime,
        bool _canStake
    ) internal returns (uint256) {
        require(_cliffTime >= _startTime, "cliff must be greater than or equal to start");
        require(_endTime > _cliffTime, "end must be greater than cliff");

        uint256 _scheduleID = schedules.length;
        schedules.push(
            Schedule(
                _amount, 0, _startTime, _cliffTime, _endTime, _canStake
            )
        );
        beneficiarySchedules[_beneficiary].push(_scheduleID);
        return _scheduleID;
    }

    function _releaseNTN(
        uint256 _scheduleID, uint256 _amount
    ) internal returns (uint256 _remaining) {
        Schedule storage _schedule = schedules[_scheduleID];
        require(_schedule.cliff <= block.number, "cliff period not reached yet");
        
        if (_amount > _schedule.currentNTNAmount) {
            _remaining = _amount - _schedule.currentNTNAmount;
            _updateAndTransferNTN(_scheduleID, msg.sender, _schedule.currentNTNAmount);
        }
        else if (_amount > 0) {
            _updateAndTransferNTN(_scheduleID, msg.sender, _amount);
        }
    }

    function _calculateUnlockedFundsAtTime(
        uint256 _scheduleID, uint256 _totalValue, uint256 _time
    ) internal view returns (uint256) {
        Schedule storage _schedule = schedules[_scheduleID];
        if (_time < _schedule.cliff) return 0;

        uint256 _unlocked = _calculateUnlockedFunds(_schedule.start, _schedule.end, _time, _totalValue);
        if (_unlocked > _schedule.withdrawnValue) {
            return _unlocked - _schedule.withdrawnValue;
        }
        return 0;
    }

    function _calculateUnlockedFunds(
        uint256 _start, uint256 _end, uint256 _time, uint256 _totalAmount
    ) internal pure returns (uint256) {
        if (_time >= _end) {
            return _totalAmount;
        }
        return _totalAmount * (_time - _start) / (_end - _start);
    }

    function _cancelSchedule(
        address _beneficiary, uint256 _id, address _recipient
    ) internal {
        uint256 _scheduleID = _getUniqueScheduleID(_beneficiary, _id);
        _changeScheduleBeneficiary(_scheduleID, _beneficiary, _recipient);
    }

    function _changeScheduleBeneficiary(
        uint256 _scheduleID, address _oldBeneficiary, address _newBeneficiary
    ) private {
        uint256[] storage _scheduleIDs = beneficiarySchedules[_oldBeneficiary];
        uint256[] memory _newScheduleIDs = new uint256[] (_scheduleIDs.length - 1);
        uint256 j = 0;
        for (uint256 i = 0; i < _scheduleIDs.length; i++) {
            if (_scheduleIDs[i] == _scheduleID) {
                continue;
            }
            _newScheduleIDs[j++] = _scheduleIDs[i];
        }
        beneficiarySchedules[_oldBeneficiary] = _newScheduleIDs;
        beneficiarySchedules[_newBeneficiary].push(_scheduleID);
    }

    /**
     * @dev returns a unique id for each schedule
     * @param _beneficiary address of the schedule holder
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (excluding canceled ones)
     */
    function _getUniqueScheduleID(address _beneficiary, uint256 _id) internal view returns (uint256) {
        require(beneficiarySchedules[_beneficiary].length > _id, "invalid schedule id");
        return beneficiarySchedules[_beneficiary][_id];
    }

    function _updateAndTransferNTN(uint256 _scheduleID, address _to, uint256 _amount) internal {
        Schedule storage _schedule = schedules[_scheduleID];
        _schedule.currentNTNAmount -= _amount;
        _schedule.withdrawnValue += _amount;
        _transferNTN(_to, _amount);
    }

    function _transferNTN(address _to, uint256 _amount) internal {
        bool _sent = autonity.transfer(_to, _amount);
        require(_sent, "NTN not transferred");
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice returns a schedule entitled to _beneficiary
     * @param _beneficiary beneficiary address
     * @param _id schedule id numbered from 0 to (n-1); n = total schedules entitled to the beneficiary (excluding canceled ones)
     */
    function getSchedule(address _beneficiary, uint256 _id) virtual external view returns (Schedule memory) {
        return schedules[_getUniqueScheduleID(_beneficiary, _id)];
    }

    /**
     * @notice returns the list of current schedules assigned to a beneficiary
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

    /**
     * @notice returns if beneficiary can stake from his schedule
     * @param _beneficiary beneficiary address
     */
    function canStake(address _beneficiary, uint256 _id) virtual external view returns (bool) {
        return schedules[_getUniqueScheduleID(_beneficiary, _id)].canStake;
    }

    /**
     * @notice returns the number of schudeled entitled to some beneficiary
     * @param _beneficiary address of the beneficiary
     */
    function totalSchedules(address _beneficiary) virtual external view returns (uint256) {
        return beneficiarySchedules[_beneficiary].length;
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
}