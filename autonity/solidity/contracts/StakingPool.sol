// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./AccessAutonity.sol";
import "./interfaces/IStakingPool.sol";
import {EnumerableSet} from "./utils/AddressSet.sol";
import {UintQueue, UintQueueLib} from "./lib/UintQueueLib.sol";

contract StakingPool is AccessAutonity, IStakingPool {
    using EnumerableSet for EnumerableSet.AddressSet;
    using UintQueueLib for UintQueue;

    /* Staking Pool of the Validator for Future Processing */
    /**
     * @dev The fields are updated as new bonding or unboning requests appear.
     * At epoch end, validators are updated with the information from the `ValidatorPool`.
     */
    struct ValidatorPool {
        uint256 selfBondingStake;
        uint256 delegatingStake;
        uint256 selfUnbondingStake;
        uint256 burningLiquid;
        uint256 selfUnbondingShare;
        uint256 unbondingShare;
    }

    /* Staking Pool of the Delegator for Future Processing */
    /**
     * @dev The fields are updated at epoch end.
     * On external calls, delegators take their share from the pool.
     */
    struct DelegatorPool {
        uint256 liquidMinted;
        uint256 selfUnbondingShare;
        uint256 unbondingShare;
        uint256 releasedSelfStake;
        uint256 releasedStake;
        uint256 atnUnrealisedFeeFactor;
    }

    struct StakingHistory {
        ValidatorPool validatorPool;
        DelegatorPool delegatorPool;
        bool notActive;
    }

    /** @dev Stores historical data for staking operations for some epoch. */
    mapping(uint256 => mapping(address => StakingHistory)) internal epochStakes;
    mapping(uint256 => EnumerableSet.AddressSet) internal validatorBonded;
    mapping(uint256 => EnumerableSet.AddressSet) internal validatorUnbonded;
    /** @dev Tracks amount of rejected bonding newton which is necessary for passive balance update in Autonity. */
    uint256 public rejectedBonding;
    /** @dev Tracks amount of released stakes from unbonding which is necessary for passive balance update in Autonity. */
    uint256 public releasedStakes;
    /** @dev The epoch id for which unbondings will be released */
    uint256 internal unbondingEpoch;

    struct UnbondingQueue {
        UintQueue queue;
        uint256 unlockingIndex;
    }
    mapping(address => UintQueue) internal pendingBondingQueue;
    mapping(address => UnbondingQueue) internal pendingUnbondingQueue;


    BondingRequest[] internal bondingMap;
    uint256 internal headBondingID;

    UnbondingRequest[] internal unbondingMap;
    uint256 internal headUnbondingID;

    address internal operator;

    constructor(address payable _autonity, address _operator) AccessAutonity(_autonity) {
        operator = _operator;
    }

    function setOperator(address _operator) external onlyAutonity {
        operator = _operator;
    }

    function bond(
        address _validator,
        uint256 _amount,
        address payable _recipient,
        uint256 _epochID,
        bool _selfBond
    ) external onlyAutonity returns (uint256) {
        uint256 _id = bondingMap.length;
        bondingMap.push(
            BondingRequest(_recipient, _validator, _amount, block.number, _epochID, _selfBond)
        );
        pendingBondingQueue[_recipient].enqueue(_id);
        validatorBonded[_epochID].add(_validator);
        // requested boning amount goes to the validators pool which will be processed at epoch end
        ValidatorPool storage _pool = epochStakes[_epochID][_validator].validatorPool;
        if (_selfBond) {
            _pool.selfBondingStake += _amount;
        }
        else {
            _pool.delegatingStake += _amount;
        }
        return _id;
    }

    function unbond(
        address _validator,
        uint256 _amount,
        address payable _recipient,
        uint256 _epochID,
        bool _selfBond
    ) external onlyAutonity returns (uint256) {
        validatorUnbonded[_epochID].add(_validator);
        // requested unboning amount goes to the validators pool which will be processed at epoch end
        ValidatorPool storage _pool = epochStakes[_epochID][_validator].validatorPool;
        if (_selfBond) {
            _pool.selfUnbondingStake += _amount;
        }
        else {
            _pool.burningLiquid += _amount;
        }
        uint256 _id = unbondingMap.length;
        unbondingMap.push(
            UnbondingRequest(
                _recipient, _validator, _amount, 0, block.number, _epochID, false, false, _selfBond
            )
        );
        pendingUnbondingQueue[_recipient].queue.enqueue(_id);
        return _id;
    }

    /**
     * @notice Function restricted to autonity contract to apply the bonding requests at epoch end.
     * @dev Apply the bonding requests and update the validator stakes from the validators pool.
     * Also update the delegators pool accordingly from which delegators can take their shares.
     * @param _epochID Epoch id of the running epoch
     */
    function applyBonding(uint256 _epochID) external onlyAutonity {
        // we remove these validators after processing
        EnumerableSet.AddressSet storage _validators = validatorBonded[_epochID];
        uint256 _length = _validators.length();
        while (_length > 0) {
            address _validator = _validators.at(0);
            StakingHistory storage _stakeInfo = epochStakes[_epochID][_validator];
            Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
            uint256 _delegation = _stakeInfo.validatorPool.delegatingStake;
            uint256 _selfBonding = _stakeInfo.validatorPool.selfBondingStake;

            if (_validatorInfo.state == ValidatorState.active) {
                // `delegatingStake` from `validatorPool` is converted to `_liquidMinted`
                uint256 _liquidMinted = _liquidFromStake(
                    _delegation,
                    _validatorInfo.bondedStake - _validatorInfo.selfBondedStake,
                    _validatorInfo.liquidSupply
                );
                // `_liquidMinted` goes into `delegatorPool`
                _stakeInfo.delegatorPool.liquidMinted = _liquidMinted;
                // Liquid balances are not updated immediately.
                // Instead they are updated via some external call.
                // To track rewards when the balances are updated,
                // we need to track `atnLastUnrealisedFeeFactor`
                _stakeInfo.delegatorPool.atnUnrealisedFeeFactor = _validatorInfo.liquidStateContract.getUnrealisedFeeFactor();

                // update the stakes for the validator
                autonity.applyBonding(
                    _validator,
                    _delegation,
                    _selfBonding,
                    _liquidMinted
                );
            }
            else {
                _stakeInfo.notActive = true;
                // track rejected bonding amount for passive update of balances
                rejectedBonding += _selfBonding + _delegation;
                emit BondingPoolRejected(_validator, _selfBonding, _delegation, _validatorInfo.state);
            }

            require(_validators.remove(_validator), "validator not removed");
            _length--;
        }
    }

    /**
     * @notice Function restricted to autonity contract to apply the unbonding requests at epoch end.
     * @dev Apply the unbonding requests and update the validator stakes from the validators pool.
     * Also update the delegators pool accordingly from which delegators can take their shares.
     * @param _epochID Epoch id of the running epoch
     */
    function applyUnbonding(uint256 _epochID) external onlyAutonity {
        // we cannot remove these validators as the unbonding requests will be released later
        EnumerableSet.AddressSet storage _validators = validatorUnbonded[_epochID];
        uint256 _length = _validators.length();
        for (uint256 i = 0; i < _length; i++) {
            address _validator = _validators.at(i);
            StakingHistory storage _stakeInfo = epochStakes[_epochID][_validator];
            Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);

            // for delegation
            uint256 _unbondingStake;
            uint256 _unbondingShare;
            if (_stakeInfo.validatorPool.burningLiquid > 0) {
                // `burningLiquid` from `validatorPool` is converted to `_unbondingStake`
                // and then converted to `_unbondingShare`
                _unbondingStake = _unbondingStakeFromLiquid(
                    _stakeInfo.validatorPool.burningLiquid,
                    _validatorInfo.liquidSupply,
                    _validatorInfo.bondedStake - _validatorInfo.selfBondedStake
                );
                _unbondingShare = _unbondingShareFromStake(
                    _unbondingStake,
                    _validatorInfo.unbondingStake,
                    _validatorInfo.unbondingShares
                );
                // `_unbondingShare` goes into both validators pool and delegators pool
                // validators pool use this information to calculate released stake when unbonding period is finished
                // delegators will take their unbonding share from the delegators pool
                _stakeInfo.delegatorPool.unbondingShare = _unbondingShare;
                _stakeInfo.validatorPool.unbondingShare = _unbondingShare;
                // Liquid balances are not updated immediately.
                // Instead they are updated via some external call.
                // To track rewards when the balances are updated,
                // we need to track `atnLastUnrealisedFeeFactor`
                _stakeInfo.delegatorPool.atnUnrealisedFeeFactor = _validatorInfo.liquidStateContract.getUnrealisedFeeFactor();
            }

            // for self-delegation
            uint256 _selfUnbondingStake =
                (_validatorInfo.selfBondedStake >= _stakeInfo.validatorPool.selfUnbondingStake) ?
                _stakeInfo.validatorPool.selfUnbondingStake : _validatorInfo.selfBondedStake;
            // `selfUnbondingStake` from `validatorPool` is converted to `_selfUnbondingShare`
            uint256 _selfUnbondingShare = _unbondingShareFromStake(
                _selfUnbondingStake,
                _validatorInfo.selfUnbondingStake,
                _validatorInfo.selfUnbondingShares
            );
            // `_selfUnbondingShare` goes into both validators pool and delegators pool
            // validators pool use this information to calculate released stake when unbonding period is finished
            // delegator will take his self unbonding share per request from the delegators pool
            _stakeInfo.delegatorPool.selfUnbondingShare = _selfUnbondingShare;
            _stakeInfo.validatorPool.selfUnbondingShare = _selfUnbondingShare;

            autonity.applyUnbonding(
                _validator,
                _stakeInfo.validatorPool.burningLiquid,
                _stakeInfo.validatorPool.selfUnbondingStake,
                _selfUnbondingStake,
                _unbondingStake,
                _selfUnbondingShare,
                _unbondingShare
            );
        }
    }

    /**
     * @notice Function restricted to autonity contract to release the unbonding requests after the unbonding period has passed.
     * @dev Release the unbonding requests and update the delegators pool from the validators pool. Validators pool were
     * updated when the unbonding requests were applied. Unboning requests are always applied before they can be released.
     * @param _epochID Epoch id of the epoch till which unbonding requests can be released.
     */
    function releaseUnbondingStake(uint256 _epochID) external onlyAutonity {
        uint256 _processingEpoch = unbondingEpoch;
        while (_processingEpoch <= _epochID) {
            // we remove these validators after processing
            EnumerableSet.AddressSet storage _validators = validatorUnbonded[_processingEpoch];
            uint256 _length = _validators.length();
            while (_length > 0) {
                address _validator = _validators.at(0);
                StakingHistory storage _stakeInfo = epochStakes[_epochID][_validator];
                Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);

                // `unbondingShare` from `validatorPool` is converted to `_releasedStake`
                uint256 _releasedStake = _releasedStakeFromShare(
                    _stakeInfo.validatorPool.unbondingShare,
                    _validatorInfo.unbondingShares,
                    _validatorInfo.unbondingStake
                );
                // `_releasedStake` goes to delegators pool
                _stakeInfo.delegatorPool.releasedStake = _releasedStake;

                // `selfUnbondingShare` from `validatorPool` is converted to `_releasedSelfStake`
                uint256 _releasedSelfStake = _releasedStakeFromShare(
                    _stakeInfo.validatorPool.selfUnbondingShare,
                    _validatorInfo.selfUnbondingShares,
                    _validatorInfo.selfUnbondingStake
                );
                // `_releasedSelfStake` goes to delegators pool
                _stakeInfo.delegatorPool.releasedSelfStake = _releasedSelfStake;
                // track released stakes for passive update of balances
                releasedStakes += _releasedStake + _releasedSelfStake;

                autonity.releaseUnbondingStake(
                    _validator,
                    _stakeInfo.validatorPool.selfUnbondingShare,
                    _stakeInfo.validatorPool.unbondingShare,
                    _releasedSelfStake,
                    _releasedStake
                );

                require(_validators.remove(_validator), "validator not removed");
                _length--;
            }
            _processingEpoch++;
        }
        unbondingEpoch = _processingEpoch;
    }

    /**
     * @notice Persist the information of the delegator in the state from the delegators pool.
     * Restricted to autonity contract.
     * @param _delegator Delegator account
     * @param _epochID Epoch id of the running epoch
     */
    function updateDelegatorPool(address _delegator, uint256 _epochID) external onlyTokenHolder(_delegator) {
        _applyBondingRequest(_delegator, _epochID);
        _applyUnbondingRequest(_delegator, _epochID);
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    function getBondingRequest(uint256 _id) external view returns (BondingRequest memory) {
        require(_id < bondingMap.length, "request doesn't exist");
        return bondingMap[_id];
    }

    function getUnbondingRequest(uint256 _id) external view returns (UnbondingRequest memory) {
        require(_id < unbondingMap.length, "request doesn't exist");
        return unbondingMap[_id];
    }

    /*
    ============================================================

        Modifiers

    ============================================================
     */

    modifier onlyTokenHolder(address _sender) {
        require(
            _sender == msg.sender || address(autonity) == msg.sender,
            "Call restricted to account holder or autonity contract"
        );
        _;
    }

    modifier onlyLiquidHolder(address _validator, address _delegator) {
        require(
            _delegator == msg.sender || address(autonity.getValidator(_validator).liquidStateContract) == msg.sender,
            "Call restricted to account holder or liquid contract"
        );
        _;
    }

    /*
    ============================================================

        Internals

    ============================================================
     */

    function _applyBondingRequest(address _delegator, uint256 _epochID) internal {
        UintQueue storage _queue = pendingBondingQueue[_delegator];
        uint256[] storage _array = _queue.array;
        uint256 _length = _array.length;
        uint256 _topIndex = _queue.topIndex;
        BondingRequest storage _request;
        StakingHistory storage _pool;
        uint256 _rejectedBonding;

        while (_topIndex < _length) {
            _request = bondingMap[_array[_topIndex]];
            if (_request.epochID > _epochID) {
                break;
            }

            // apply the bonding request
            _pool = epochStakes[_request.epochID][_request.delegatee];
            if (_pool.notActive) {
                // validator was inactive at the time
                _rejectedBonding += _request.amount;
            }
            else if (!_request.selfDelegation && _pool.delegatorPool.liquidMinted > 0) {
                // calculate liquid amount
                uint256 _earnedLiquid = _liquidFromStake(
                    _request.amount,
                    _pool.validatorPool.delegatingStake,
                    _pool.delegatorPool.liquidMinted
                );
                // update the liquid balance
                autonity.getValidator(_request.delegatee).liquidStateContract.transferLiquidFromPool(
                    _request.delegator,
                    _earnedLiquid,
                    _pool.delegatorPool.atnUnrealisedFeeFactor
                );
                // remove the share of the delegator from the delegators pool
                _pool.delegatorPool.liquidMinted -= _earnedLiquid;

                // clean some storage
                if (_pool.delegatorPool.liquidMinted == 0 && _pool.validatorPool.burningLiquid == 0) {
                    _pool.delegatorPool.atnUnrealisedFeeFactor = 0;
                }
            }
            
            // remove the requested bonding amount from the validators pool
            if (_request.selfDelegation) {
                // remove the share of the delegator from the pool
                _pool.validatorPool.selfBondingStake -= _request.amount;
            }
            else {
                _pool.validatorPool.delegatingStake -= _request.amount;
            }

            _topIndex++;
        }
        
        // TODO (tariq): consider deleting bonding request from `bondingMap` as they are applied
        _queue.dequeue(_topIndex - _queue.topIndex);
        if (_rejectedBonding > 0) {
            rejectedBonding -= _rejectedBonding;
            autonity.updateWithRejectedBondingAmount(_delegator, _rejectedBonding);
        }
    }

    function _applyUnbondingRequest(address _delegator, uint256 _epochID) internal {
        UnbondingQueue storage _queue = pendingUnbondingQueue[_delegator];
        uint256[] storage _array = _queue.queue.array;
        uint256 _length = _array.length;
        uint256 _processingIndex = _queue.unlockingIndex;
        UnbondingRequest storage _request;
        StakingHistory storage _pool;

        while (_processingIndex < _length) {
            _request = unbondingMap[_array[_processingIndex]];
            if (_request.epochID > _epochID) {
                break;
            }

            // apply the unbonding request
            _request.unlocked = true;
            _pool = epochStakes[_request.epochID][_request.delegatee];
            if (_request.selfDelegation) {
                // calculate unbonding share
                uint256 _share = _unbondingShareFromRequestAmount(
                    _request.amount,
                    _pool.validatorPool.selfUnbondingStake,
                    _pool.delegatorPool.selfUnbondingShare
                );
                _request.unbondingShare = _share;

                // remove the share of the delegator from the delegators pool
                _pool.delegatorPool.selfUnbondingShare -= _share;
                // remove the requested amount from validators pool
                _pool.validatorPool.selfUnbondingStake -= _request.amount;
            }
            else {
                // unlock the liquid balance
                autonity.getValidator(_request.delegatee).liquidStateContract.unlockAndBurnLiquid(
                    _request.delegator,
                    _request.amount,
                    _pool.delegatorPool.atnUnrealisedFeeFactor
                );
                // calculate unbonding share
                uint256 _share = _unbondingShareFromRequestAmount(
                    _request.amount,
                    _pool.validatorPool.burningLiquid,
                    _pool.delegatorPool.unbondingShare
                );
                _request.unbondingShare = _share;
                
                // remove the share of the delegator from the delegators pool
                _pool.delegatorPool.unbondingShare -= _share;
                // remove the requested amount from validators pool
                _pool.validatorPool.burningLiquid -= _request.amount;

                // clean some storage
                if (_pool.delegatorPool.liquidMinted == 0 && _pool.validatorPool.burningLiquid == 0) {
                    _pool.delegatorPool.atnUnrealisedFeeFactor = 0;
                }
            }

            _processingIndex++;
        }
        _queue.unlockingIndex = _processingIndex;
    }

    function _releaseUnbondingStake(address _delegator, uint256 _epochID) internal {
        UintQueue storage _queue = pendingUnbondingQueue[_delegator].queue;
        uint256[] storage _array = _queue.array;
        uint256 _length = _array.length;
        uint256 _topIndex = _queue.topIndex;
        UnbondingRequest storage _request;
        StakingHistory storage _pool;

        uint256 _releasedStakes;
        while (_topIndex < _length) {
            _request = unbondingMap[_array[_topIndex]];
            if (_request.epochID > _epochID) {
                break;
            }

            // release the stakes
            _request.released = true;
            _pool = epochStakes[_request.epochID][_request.delegatee];
            if (_request.selfDelegation) {
                _releasedStakes += _releasedStakeFromShare(
                    _request.unbondingShare,
                    _pool.validatorPool.selfUnbondingShare,
                    _pool.delegatorPool.releasedSelfStake
                );
            }
            else {
                _releasedStakes += _releasedStakeFromShare(
                    _request.unbondingShare,
                    _pool.validatorPool.unbondingShare,
                    _pool.delegatorPool.releasedStake
                );
            }
            _topIndex++;
        }

        // TODO (tariq): consider deleting unbonding request from `unbondingMap` as they are released
        _queue.dequeue(_topIndex - _queue.topIndex);
        if (_releasedStakes > 0) {
            releasedStakes -= _releasedStakes;
            autonity.updateWithReleasedStake(_delegator, _releasedStakes);
        }
    }

    /**
     * @dev Calculates minted liquid amount from the amount of newton bonded.
     * As liquid is minted, in case of `_totalDelegation == 0`, we mint in 1:1 ratio.
     * @param _newtonBonded new bonded stake
     * @param _totalDelegation total delegated stake in existence or in a pool
     * @param _totalLiquid total liquid in supply or in a pool
     */
    function _liquidFromStake(
        uint256 _newtonBonded,
        uint256 _totalDelegation,
        uint256 _totalLiquid
    ) internal pure returns (uint256) {
        if (_totalDelegation == 0) {
            return _newtonBonded;
        }
        return (_totalLiquid * _newtonBonded) / _totalDelegation;
    }

    /**
     * @dev Calculates amount of unbonding stake from the amount of liquid. As nothing is minted here,
     * we can assume `_liquidBurning <= _totalLiquid`.
     * @param _liquidBurning amount of liquid being burnt
     * @param _totalLiquid total liquid in supply or in a pool
     * @param _totalDelegation total delegated stake in existence or in a pool 
     */
    function _unbondingStakeFromLiquid(
        uint256 _liquidBurning,
        uint256 _totalLiquid,
        uint256 _totalDelegation
    ) internal pure returns (uint256) {
        if (_liquidBurning == 0) {
            return 0;
        }
        // assuming valid inputs `_liquidBurning <= _totalLiquid`
        return _convertFromRatio(
            _liquidBurning,
            _totalDelegation,
            _totalLiquid
        );
    }

    /**
     * @dev Calculates amount of unbonding share from unbonding stake. Unbonding share is minted,
     * so in case of `_totalUnbondingStake == 0`, we mint in 1:1 ratio.
     * @param _unbondingStake amount of stake under unbonding
     * @param _totalUnbondingStake total unbonding stake in existence or in a pool
     * @param _totalUnbondingShare total unbonding share in existence or in a pool
     */
    function _unbondingShareFromStake(
        uint256 _unbondingStake,
        uint256 _totalUnbondingStake,
        uint256 _totalUnbondingShare
    ) internal pure returns (uint256) {
        if (_totalUnbondingStake == 0) {
            return _unbondingStake;
        }
        return _convertFromRatio(
            _unbondingStake,
            _totalUnbondingShare,
            _totalUnbondingStake
        );
    }

    /**
     * @dev Calculates amount of unbonding share from requested self unbonding stake or liquid amount.
     * As unbonding share is not minted here, we can assume `_requestAmount <= _totalUnbondingAmount`.
     * @param _requestAmount amount of liquid or self unbonding stake
     * @param _totalUnbondingAmount total unbonding stake or liquid in existence or in a pool
     * @param _totalUnbondingShare total unbonding share in existence or in a pool
     */
    function _unbondingShareFromRequestAmount(
        uint256 _requestAmount,
        uint256 _totalUnbondingAmount,
        uint256 _totalUnbondingShare
    ) internal pure returns (uint256) {
        if (_requestAmount == 0) {
            return 0;
        }
        // assuming valid inputs `_requestAmount <= _totalUnbondingAmount`
        return _convertFromRatio(
            _requestAmount,
            _totalUnbondingShare,
            _totalUnbondingAmount
        );
    }

    /**
     * @dev Calculates amount of stakes released from unbonding shares. As nothing is minted here,
     * we can assume `_unbondingShare <= _totalUnbondingStake`.
     * @param _unbondingShare unbonding share
     * @param _totalUnbondingShare total unbonding share in existence or in a pool
     * @param _totalUnbondingStake total unbonding stake in existence or in a pool
     */
    function _releasedStakeFromShare(
        uint256 _unbondingShare,
        uint256 _totalUnbondingShare,
        uint256 _totalUnbondingStake
    ) internal pure returns (uint256) {
        if (_unbondingShare == 0) {
            return 0;
        }
        // assuming valid inputs `_unbondingShare <= _totalUnbondingStake`
        return _convertFromRatio(
            _unbondingShare,
            _totalUnbondingStake,
            _totalUnbondingShare
        );
    }

    function _convertFromRatio(
        uint256 _share,
        uint256 _ratioNumerator,
        uint256 _ratioDenominator
    ) internal pure returns (uint256) {
        return (_share * _ratioNumerator) / _ratioDenominator;
    }
}