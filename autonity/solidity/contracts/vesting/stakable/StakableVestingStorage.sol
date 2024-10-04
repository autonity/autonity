// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../ContractBase.sol";
import "../StakableVestingManager.sol";
import "./ValidatorManagerStorage.sol";

abstract contract StakableVestingStorage is ValidatorManagerStorage {
    constructor() {}
    address internal beneficiary;
    StakableVestingManager internal managerContract;
    ContractBase.Contract internal stakableContract;

    struct ContractValuation {
        uint256 totalShare;
        uint256 withdrawnShare;
    }

    ContractValuation internal contractValuation;

    struct PendingBondingRequest {
        uint256 amount;
        uint256 epochID;
        address validator;
    }

    PendingBondingRequest[] internal bondingQueue;
    uint256 internal bondingQueueTopIndex;

    struct PendingUnbondingRequest {
        uint256 unbondingID;
        uint256 epochID;
        address validator;
    }

    PendingUnbondingRequest[] internal unbondingQueue;
    uint256 internal unbondingQueueTopIndex;
}