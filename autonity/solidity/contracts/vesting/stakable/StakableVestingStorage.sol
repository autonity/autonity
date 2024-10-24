// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../interfaces/IStakableVestingManager.sol";
import "../ContractBase.sol";
import {StakingRequestQueue} from "./QueueLib.sol";
import "./ValidatorManagerStorage.sol";

abstract contract StakableVestingStorage is ValidatorManagerStorage {
    constructor() {}
    address internal beneficiary;
    IStakableVestingManager internal managerContract;
    ContractBase.Contract internal stakableContract;

    struct ContractValuation {
        uint256 totalShare;
        uint256 withdrawnShare;
    }

    ContractValuation internal contractValuation;
    StakingRequestQueue internal bondingQueue;
    StakingRequestQueue internal unbondingQueue;
}