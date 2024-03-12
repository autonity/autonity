// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

contract VestingManager {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release schedule
    // UNLOCKED are tokens that got released but not yet transferred
    uint256 public contractVersion = 1;

    struct Schedule {
        uint256 amount;
        uint256 start;
        uint256 cliff;
        uint256 end; // or duration?
        uint256 stackable;
    }

    mapping(address => Schedule[]) internal schedules;

    constructor(address _autonity, address _operator){
        // save autonity and operator account  - with standard modifers
    }

    function newSchedule(
        address _beneficiary,
        uint256 _amount,
        uint256 _startBlock,
        uint256 _cliffBlock,
        uint256 _endBlock,
        bool _stackable) virtual onlyOperator public {
        // need to be AUTHORIZED to transfer AMOUNT
        require(_cliffBlock >= _startBlock, "cliff must be greater to start");
        // use transferFrom(ERC20) here - meaning operator has to authorize transfer from VestingManager
        // check if funds were transferered if so create schedule
    }

    // retrieve list of current schedules assigned to a beneficiary
    function getSchedules(address _beneficiary) virtual public {

    }

    // used by beneficiary to transfer unlocked ntn
    function releaseFunds(uint256 _id) virtual public {
        // not only unlocked token but unlocked LNTN too !!!
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
    function bond(id, validator, amount) virtual public {

    }

    function unbond(id, validator, amount) virtual public {

    }
    // add std modifiers here

    // - Commission rate - one or two?
    // - staking rewards, who gets them are they locked
    //          - NTN and ATN immediate withdraw
    //          -  at our discretion
    //
    // - can you stake unlocked but not transferred yet tokens

    //  -- -


    // one address - one contract?? that would simplify accounting for sure !

    // we don't know how much we are going to get from autonity "bond" function.
}
