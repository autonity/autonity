// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../interfaces/IStakeableVestingManager.sol";
import "../ContractBase.sol";
import {StakingRequestQueue} from "./QueueLib.sol";
import "./ValidatorManagerStorage.sol";

/** @title Storage of the Stakeable Vesting Smart Contract for vesting and staking funds */
abstract contract StakeableVestingStorage is ValidatorManagerStorage {
    constructor() {}
    address internal beneficiary;
    IStakeableVestingManager internal managerContract;
    ContractBase.Contract internal stakeableContract;

    struct ContractValuation {
        uint256 totalShare;
        uint256 withdrawnShare;
    }

    ContractValuation internal contractValuation;
    StakingRequestQueue internal bondingQueue;
    StakingRequestQueue internal unbondingQueue;
}