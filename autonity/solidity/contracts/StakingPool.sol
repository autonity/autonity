// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./AccessAutonity.sol";
import "./interfaces/IStakingPool.sol";
import {EnumerableSet} from "./utils/AddressSet.sol";
import {UintQueue, UintQueueLib} from "./lib/UintQueueLib.sol";
import {StakingPoolMath} from "./lib/StakingPoolMath.sol";

contract StakingPool is AccessAutonity, IStakingPool {
    using EnumerableSet for EnumerableSet.AddressSet;
    using UintQueueLib for UintQueue;

    /** @dev Stores historical data for staking operations for some epoch. */
    mapping(uint256 => mapping(address => PoolCollection)) internal epochStakesPool;
    mapping(uint256 => EnumerableSet.AddressSet) internal validatorBonded;
    mapping(uint256 => EnumerableSet.AddressSet) internal validatorUnbonded;
    /** @dev Tracks amount of rejected bonding newton which is necessary for passive balance update in Autonity. */
    uint256 public rejectedBonding;
    /** @dev Tracks amount of released stakes from unbonding which is necessary for passive balance update in Autonity. */
    uint256 public releasedStakes;
    /** @dev The epoch id for which unbondings will be released */
    uint256 public unbondingEpoch;
    mapping(address => uint256) internal lastFeeFactor;

    struct UnbondingQueue {
        UintQueue queue;
        uint256 unlockingIndex;
    }
    mapping(address => UintQueue) internal pendingBondingQueue;
    mapping(address => UnbondingQueue) internal pendingUnbondingQueue;


    BondingRequest[] internal bondingArray;

    UnbondingRequest[] internal unbondingArray;

    address internal operator;

    constructor(address payable _autonity, address _operator) AccessAutonity(_autonity) {
        operator = _operator;
    }

    /**
    * @dev Receive Auton function https://solidity.readthedocs.io/en/v0.7.2/contracts.html#receive-ether-function
    *
    */
    receive() external payable {}

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
        uint256 _id = bondingArray.length;
        bondingArray.push(
            BondingRequest(_recipient, _validator, _amount, block.number, _epochID, _selfBond)
        );
        pendingBondingQueue[_recipient].enqueue(_id);
        validatorBonded[_epochID].add(_validator);
        // requested boning amount goes to the validators pool which will be processed at epoch end
        ValidatorBondingPool storage _pool = epochStakesPool[_epochID][_validator].validatorBondingPool;
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
        PoolCollection storage _pool = epochStakesPool[_epochID][_validator];
        if (_selfBond) {
            _pool.validatorUnbondingPool.selfUnbondingStake += _amount;
        }
        else {
            if (_pool.validatorUnbondingPool.liquidBurning == 0) {
                _pool.delegatorUnbondingPool.feeFactor = lastFeeFactor[_validator];
            }
            _pool.validatorUnbondingPool.liquidBurning += _amount;
        }
        uint256 _id = unbondingArray.length;
        unbondingArray.push(
            UnbondingRequest(
                _recipient, _validator, _amount, 0, block.number, _epochID, false, _selfBond
            )
        );
        pendingUnbondingQueue[_recipient].queue.enqueue(_id);
        return _id;
    }

    /**
     * @notice Claims and tracks rewards for each committee member. Restricted to autonity contract.
     * @dev Must be called before doing any staking operation.
     * @param _validators committee members
     * @param _liquidContracts liquid state contract addresses of the respective committee member
     */
    function collectRewards(address[] memory _validators, ILiquid[] memory _liquidContracts) external onlyAutonity {
        uint256 _count = _validators.length;
        require(_count == _liquidContracts.length, "invalid inputs");
        uint256 _atnBalance = address(this).balance;
        for (uint256 i = 0; i < _count; i++) {
            uint256 _liquidBalance = _liquidContracts[i].balanceInContract(address(this));
            if (_liquidBalance > 0) {
                _liquidContracts[i].claimRewards();
                lastFeeFactor[_validators[i]] += (address(this).balance - _atnBalance) * FEE_FACTOR_UNIT_RECIP / _liquidBalance;
                _atnBalance = address(this).balance;
            }
        }
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
            PoolCollection storage _stakePool = epochStakesPool[_epochID][_validator];
            ValidatorBondingPool storage _validatorPool = _stakePool.validatorBondingPool;
            Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);

            if (_validatorInfo.state == ValidatorState.active) {
                // `delegatingStake` from `_validatorPool` is converted to `_liquidMinted`
                uint256 _liquidMinted = StakingPoolMath.liquidFromStake(
                    _validatorPool.delegatingStake,
                    _validatorInfo.bondedStake - _validatorInfo.selfBondedStake,
                    _validatorInfo.liquidSupply
                );

                if (_liquidMinted > 0) {
                    // `_liquidMinted` goes into `delegatorPool`
                    _stakePool.delegatorBondingPool.liquidMinted = _liquidMinted;
                    // To track rewards when delegators take their share from the pool,
                    // we need to track `lastFeeFactor`.
                    _stakePool.delegatorBondingPool.feeFactor = lastFeeFactor[_validator];
                }

                // update the stakes for the validator
                autonity.applyBonding(
                    _validator,
                    _validatorPool.selfBondingStake,
                    _validatorPool.delegatingStake,
                    _liquidMinted
                );
            }
            else {
                _validatorPool.notActive = true;
                // track rejected bonding amount for passive update of balances
                rejectedBonding += _validatorPool.selfBondingStake + _validatorPool.delegatingStake;
                emit BondingPoolRejected(
                    _validator,
                    _validatorPool.selfBondingStake,
                    _validatorPool.delegatingStake,
                    _validatorInfo.state
                );
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
            PoolCollection storage _stakePool = epochStakesPool[_epochID][_validator];
            Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);

            // for delegation
            uint256 _unbondingStake;
            uint256 _unbondingShare;
            if (_stakePool.validatorUnbondingPool.liquidBurning > 0) {
                // `liquidBurning` from `validatorUnbondingPool` is converted to `_unbondingStake`
                // and then converted to `_unbondingShare`
                _unbondingStake = StakingPoolMath.unbondingStakeFromLiquid(
                    _stakePool.validatorUnbondingPool.liquidBurning,
                    _validatorInfo.liquidSupply,
                    _validatorInfo.bondedStake - _validatorInfo.selfBondedStake
                );
                _unbondingShare = StakingPoolMath.unbondingShareFromStake(
                    _unbondingStake,
                    _validatorInfo.unbondingStake,
                    _validatorInfo.unbondingShares
                );

                // `_unbondingShare` goes into both validators pool and delegators pool
                // validators pool use this information to calculate released stake when unbonding period is finished
                // delegators will take their unbonding share from the delegators pool
                _stakePool.delegatorUnbondingPool.unbondingShare = _unbondingShare;
                _stakePool.validatorUnbondingPool.unbondingShare = _unbondingShare;

                // `_stakePool.validatorUnbondingPool.liquidBurning` is going to be burnt. Rewards accumulated due to this amount
                // needs to be stored for resdributing later on external call
                _stakePool.delegatorUnbondingPool.rewardsCollected = StakingPoolMath.computeUnrealisedRewards(
                    lastFeeFactor[_validator],
                    _stakePool.delegatorUnbondingPool.feeFactor,
                    _stakePool.validatorUnbondingPool.liquidBurning
                );
                _stakePool.delegatorUnbondingPool.feeFactor = 0;
            }

            // for self-delegation
            uint256 _selfUnbondingStake =
                (_validatorInfo.selfBondedStake >= _stakePool.validatorUnbondingPool.selfUnbondingStake) ?
                _stakePool.validatorUnbondingPool.selfUnbondingStake : _validatorInfo.selfBondedStake;
            // `selfUnbondingStake` from `validatorUnbondingPool` is converted to `_selfUnbondingShare`
            uint256 _selfUnbondingShare = StakingPoolMath.unbondingShareFromStake(
                _selfUnbondingStake,
                _validatorInfo.selfUnbondingStake,
                _validatorInfo.selfUnbondingShares
            );

            // `_selfUnbondingShare` goes into both validators pool and delegators pool
            // validators pool use this information to calculate released stake when unbonding period is finished
            // delegator will take his self unbonding share per request from the delegators pool
            _stakePool.delegatorUnbondingPool.selfUnbondingShare = _selfUnbondingShare;
            _stakePool.validatorUnbondingPool.selfUnbondingShare = _selfUnbondingShare;

            autonity.applyUnbonding(
                _validator,
                _stakePool.validatorUnbondingPool.liquidBurning,
                _stakePool.validatorUnbondingPool.selfUnbondingStake,
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
                PoolCollection storage _stakePool = epochStakesPool[_processingEpoch][_validator];
                Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);

                // `unbondingShare` from `validatorUnbondingPool` is converted to `_releasedStake`
                uint256 _releasedStake = StakingPoolMath.releasedStakeFromShare(
                    _stakePool.validatorUnbondingPool.unbondingShare,
                    _validatorInfo.unbondingShares,
                    _validatorInfo.unbondingStake
                );
                // `_releasedStake` goes to delegators pool
                _stakePool.delegatorUnbondingPool.releasedStake = _releasedStake;

                // `selfUnbondingShare` from `validatorUnbondingPool` is converted to `_releasedSelfStake`
                uint256 _releasedSelfStake = StakingPoolMath.releasedStakeFromShare(
                    _stakePool.validatorUnbondingPool.selfUnbondingShare,
                    _validatorInfo.selfUnbondingShares,
                    _validatorInfo.selfUnbondingStake
                );
                // `_releasedSelfStake` goes to delegators pool
                _stakePool.delegatorUnbondingPool.releasedSelfStake = _releasedSelfStake;
                // track released stakes for passive update of balances
                releasedStakes += _releasedStake + _releasedSelfStake;

                autonity.releaseUnbondingStake(
                    _validator,
                    _stakePool.validatorUnbondingPool.selfUnbondingShare,
                    _stakePool.validatorUnbondingPool.unbondingShare,
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
     * Restricted to the autonity contract or the delegator.
     * @param _delegator Delegator account
     */
    function updateDelegatorPool(address _delegator) external virtual onlyTokenHolder(_delegator) {
        _updateDelegatorPool(_delegator);
    }

    /**
     * @notice Persist the information of the delegator in the state from the delegators pool.
     * Restricted to the autonity contract or the delegator.
     * @param _delegator Delegator account
     * @param _validator Validator address
     */
    function updateDelegatorPool1(address _delegator, address _validator) external virtual onlyLiquidHolder(_delegator, _validator) {
        _updateDelegatorPool(_delegator);
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    function getBondingRequest(uint256 _id) external view returns (BondingRequest memory) {
        require(_id < bondingArray.length, "request doesn't exist");
        return bondingArray[_id];
    }

    function getUnbondingRequest(uint256 _id) external view returns (UnbondingRequest memory) {
        require(_id < unbondingArray.length, "request doesn't exist");
        return unbondingArray[_id];
    }

    function getBondingArrayLength() external view returns (uint256) {
        return bondingArray.length;
    }

    function getUnbondingArrayLength() external view returns (uint256) {
        return unbondingArray.length;
    }

    function isUnbondingReleased(uint256 _id) external view returns (bool) {
        require(unbondingArray.length > _id, "request doesn't exist");
        return unbondingArray[_id].epochID < unbondingEpoch;
    }

    function getUnbondingShare(uint256 _id) public view returns (uint256) {
        require(unbondingArray.length > _id, "request doesn't exist");
        UnbondingRequest storage _item = unbondingArray[_id];
        if (_item.unlocked) {
            return _item.unbondingShare;
        }

        // assuming unbonding is applied
        PoolCollection storage _pool = epochStakesPool[_item.epochID][_item.validator];
        if (_item.selfDelegation) {
            return StakingPoolMath.unbondingShareFromRequestAmount(
                _item.amount,
                _pool.validatorUnbondingPool.selfUnbondingStake,
                _pool.delegatorUnbondingPool.selfUnbondingShare
            );
        }
        return StakingPoolMath.unbondingShareFromRequestAmount(
            _item.amount,
            _pool.validatorUnbondingPool.liquidBurning,
            _pool.delegatorUnbondingPool.unbondingShare
        );
    }

    function calculateReleasedStake(address _delegator) external view returns (uint256) {
        uint256 _epochID = unbondingEpoch;
        UintQueue storage _queue = pendingUnbondingQueue[_delegator].queue;
        uint256 _length = _queue.array.length;
        uint256 _topIndex = _queue.topIndex;
        uint256 _releasedStakes;

        while (_topIndex < _length) {
            UnbondingRequest storage _request = unbondingArray[_queue.array[_topIndex]];
            if (_request.epochID >= _epochID) {
                break;
            }

            PoolCollection storage _pool = epochStakesPool[_request.epochID][_request.validator];
            // calculate released stake
            if (_request.selfDelegation) {
                _releasedStakes += StakingPoolMath.releasedStakeFromShare(
                    getUnbondingShare(_queue.array[_topIndex]),
                    _pool.validatorUnbondingPool.selfUnbondingShare,
                    _pool.delegatorUnbondingPool.releasedSelfStake
                );
            }
            else {
                _releasedStakes += StakingPoolMath.releasedStakeFromShare(
                    getUnbondingShare(_queue.array[_topIndex]),
                    _pool.validatorUnbondingPool.unbondingShare,
                    _pool.delegatorUnbondingPool.releasedStake
                );
            }
            _topIndex++;
        }
        return _releasedStakes;
    }

    function calculateRejectedBonding(address _delegator, uint256 _epochID) external view returns (uint256) {
        UintQueue storage _queue = pendingBondingQueue[_delegator];
        uint256 _length = _queue.array.length;
        uint256 _topIndex = _queue.topIndex;
        uint256 _bondingRejected;

        while (_topIndex < _length) {
            BondingRequest storage _request = bondingArray[_queue.array[_topIndex]];
            if (_request.epochID >= _epochID) {
                break;
            }

            if (epochStakesPool[_request.epochID][_request.validator].validatorBondingPool.notActive) {
                // validator was inactive at the time
                _bondingRejected += _request.amount;
            }
            _topIndex++;
        }
        return _bondingRejected;
    }

    function calculateLiquidMinted(address _delegator, address _validator) external view returns (uint256) {
        uint256 _epochID = autonity.epochID();
        UintQueue storage _queue = pendingBondingQueue[_delegator];
        uint256 _length = _queue.array.length;
        uint256 _topIndex = _queue.topIndex;
        uint256 _liquidMinted;

        while (_topIndex < _length) {
            BondingRequest storage _request = bondingArray[_queue.array[_topIndex]];
            if (_request.epochID == _epochID) {
                break;
            }
            if (_request.validator != _validator) {
                _topIndex++;
                continue;
            }
            if (_request.selfDelegation) {
                break;
            }

            PoolCollection storage _pool = epochStakesPool[_request.epochID][_validator];
            if (!_pool.validatorBondingPool.notActive) {
                // validator was active at the time and liquid is minted
                _liquidMinted += StakingPoolMath.liquidFromStake(
                    _request.amount,
                    _pool.validatorBondingPool.delegatingStake,
                    _pool.delegatorBondingPool.liquidMinted
                );
            }
            _topIndex++;
        }
        return _liquidMinted;
    }

    function calculateLiquidBurning(address _delegator, address _validator) external view returns (uint256) {
        UnbondingQueue storage _queue = pendingUnbondingQueue[_delegator];
        uint256 _length = _queue.queue.array.length;
        uint256 _processingIndex = _queue.unlockingIndex;
        uint256 _epochID = autonity.epochID();
        uint256 _liquidBurning;

        while (_processingIndex < _length) {
            UnbondingRequest storage _request = unbondingArray[_queue.queue.array[_processingIndex]];
            if (_request.validator != _validator) {
                _processingIndex++;
                continue;
            }
            if (_request.selfDelegation) {
                break;
            }
            if (_request.epochID < _epochID) {
                // these requests are already unlocked and liquid is already burnt
                _processingIndex++;
                continue;
            }

            _liquidBurning += _request.amount;
            _processingIndex++;
        }
        return _liquidBurning;
    }

    function calculateRewards(address _delegator, address _validator) external view returns (uint256) {
        uint256 _epochID = autonity.epochID();
        return _rewardsFromLiquidMinted(_delegator, _validator, _epochID) + _rewardsFromLiquidBurning(_delegator, _validator, _epochID);
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

    modifier onlyLiquidHolder(address _delegator, address _validator) {
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

    function _updateDelegatorPool(address _delegator) internal {
        uint256 _epochID = autonity.epochID();
        uint256 _rewards = _processBondingRequest(_delegator, _epochID);
        _rewards +=_processUnbondingRequest(_delegator, _epochID);
        _releaseUnbondingStake(_delegator, unbondingEpoch);

        if (_rewards > 0) {
            //   solhint-disable-next-line avoid-low-level-calls
            (bool _sent, ) = _delegator.call{value: _rewards}("");
            require(_sent, "Failed to send ATN");
        }
    }

    /**
     * @dev Apply all bonding requests from `_delegator` coming in epoch `_epochID` or before.
     */
    function _processBondingRequest(address _delegator, uint256 _epochID) internal returns (uint256) {
        UintQueue storage _queue = pendingBondingQueue[_delegator];
        uint256 _length = _queue.array.length;
        uint256 _topIndex = _queue.topIndex;
        uint256 _bondingRejected;
        uint256 _rewards;

        while (_topIndex < _length) {
            BondingRequest storage _request = bondingArray[_queue.array[_topIndex]];
            if (_request.epochID == _epochID) {
                break;
            }

            // apply the bonding request
            PoolCollection storage _pool = epochStakesPool[_request.epochID][_request.validator];
            if (_pool.validatorBondingPool.notActive) {
                // validator was inactive at the time
                _bondingRejected += _request.amount;
            }
            else if (!_request.selfDelegation && _pool.delegatorBondingPool.liquidMinted > 0) {
                // calculate liquid amount
                uint256 _totalLiquid = _pool.delegatorBondingPool.liquidMinted;
                uint256 _earnedLiquid = StakingPoolMath.liquidFromStake(
                    _request.amount,
                    _pool.validatorBondingPool.delegatingStake,
                    _totalLiquid
                );
                // update the liquid balance
                autonity.getValidator(_request.validator).liquidStateContract.transferFromPool(
                    _request.delegator,
                    _earnedLiquid
                );
                
                // calculate and send rewards due to `_earnedLiquid`
                // first calculate total rewards for `_totalLiquid` and store it
                uint256 _totalRewards = _pool.delegatorBondingPool.rewardsCollected + StakingPoolMath.computeUnrealisedRewards(
                    lastFeeFactor[_request.validator],
                    _pool.delegatorBondingPool.feeFactor,
                    _totalLiquid
                );
                _pool.delegatorBondingPool.feeFactor = lastFeeFactor[_request.validator];
                // now compute the share earned due to `_earnedLiquid`
                // this way we have minimal dust remaining in address(this).balance
                uint256 _rewardCalculated = StakingPoolMath.computeRewardFraction(
                    _totalRewards,
                    _earnedLiquid,
                    _totalLiquid
                );
                _rewards += _rewardCalculated;

                // remove the share of the delegator from the delegators pool
                _pool.delegatorBondingPool.liquidMinted = _totalLiquid - _earnedLiquid;
                _pool.delegatorBondingPool.rewardsCollected = _totalRewards - _rewardCalculated;

                // clean some storage
                if (_totalLiquid == _earnedLiquid) {
                    _pool.delegatorBondingPool.feeFactor = 0;
                }
            }
            
            // remove the requested bonding amount from the validators pool
            if (_request.selfDelegation) {
                _pool.validatorBondingPool.selfBondingStake -= _request.amount;
            }
            else {
                _pool.validatorBondingPool.delegatingStake -= _request.amount;
            }

            _topIndex++;
        }
        
        // TODO (tariq): consider deleting bonding request from `bondingArray` as they are applied
        _queue.dequeue(_topIndex - _queue.topIndex);
        if (_bondingRejected > 0) {
            rejectedBonding -= _bondingRejected;
            autonity.updateWithRejectedBondingAmount(_delegator, _bondingRejected);
        }
        return _rewards;
    }

    /**
     * @dev Apply all unbonding requests from `_delegator` coming in epoch `_epochID` or before.
     */
    function _processUnbondingRequest(address _delegator, uint256 _epochID) internal returns (uint256) {
        UnbondingQueue storage _queue = pendingUnbondingQueue[_delegator];
        uint256 _length = _queue.queue.array.length;
        uint256 _processingIndex = _queue.unlockingIndex;
        uint256 _rewards;

        while (_processingIndex < _length) {
            UnbondingRequest storage _request = unbondingArray[_queue.queue.array[_processingIndex]];
            if (_request.epochID == _epochID) {
                break;
            }

            // apply the unbonding request
            _request.unlocked = true;
            PoolCollection storage _pool = epochStakesPool[_request.epochID][_request.validator];
            if (_request.selfDelegation) {
                // calculate unbonding share
                uint256 _share = StakingPoolMath.unbondingShareFromRequestAmount(
                    _request.amount,
                    _pool.validatorUnbondingPool.selfUnbondingStake,
                    _pool.delegatorUnbondingPool.selfUnbondingShare
                );
                _request.unbondingShare = _share;

                // remove the share of the delegator from the delegators pool
                _pool.delegatorUnbondingPool.selfUnbondingShare -= _share;
                // remove the requested amount from validators pool
                _pool.validatorUnbondingPool.selfUnbondingStake -= _request.amount;
            }
            else {
                // calculate unbonding share
                uint256 _totalLiquid = _pool.validatorUnbondingPool.liquidBurning;
                uint256 _share = StakingPoolMath.unbondingShareFromRequestAmount(
                    _request.amount,
                    _totalLiquid,
                    _pool.delegatorUnbondingPool.unbondingShare
                );
                _request.unbondingShare = _share;

                // calculate and send the rewards due to `_request.amount`
                uint256 _rewardCalculated = StakingPoolMath.computeRewardFraction(
                    _pool.delegatorUnbondingPool.rewardsCollected,
                    _request.amount,
                    _totalLiquid
                );
                _rewards += _rewardCalculated;
                
                // remove the share of the delegator from the delegators pool
                _pool.delegatorUnbondingPool.unbondingShare -= _share;
                _pool.delegatorUnbondingPool.rewardsCollected -= _rewardCalculated;
                // remove the requested amount from validators pool
                _pool.validatorUnbondingPool.liquidBurning = _totalLiquid - _request.amount;
            }

            _processingIndex++;
        }
        _queue.unlockingIndex = _processingIndex;
        return _rewards;
    }

    /**
     * @dev Release all unbonding requests from `_delegator` coming in epoch `_epochID` or before.
     */
    function _releaseUnbondingStake(address _delegator, uint256 _epochID) internal {
        UnbondingQueue storage _queue = pendingUnbondingQueue[_delegator];
        uint256 _length = _queue.queue.array.length;
        uint256 _topIndex = _queue.queue.topIndex;
        uint256 _releasedStakes;

        while (_topIndex < _length) {
            UnbondingRequest storage _request = unbondingArray[_queue.queue.array[_topIndex]];
            if (_request.epochID >= _epochID) {
                break;
            }

            // release the stakes
            PoolCollection storage _pool = epochStakesPool[_request.epochID][_request.validator];
            if (_request.selfDelegation) {
                uint256 _share = StakingPoolMath.releasedStakeFromShare(
                    _request.unbondingShare,
                    _pool.validatorUnbondingPool.selfUnbondingShare,
                    _pool.delegatorUnbondingPool.releasedSelfStake
                );
                _pool.validatorUnbondingPool.selfUnbondingShare -= _request.unbondingShare;
                _pool.delegatorUnbondingPool.releasedSelfStake -= _share;
                _releasedStakes += _share;
            }
            else {
                uint256 _share = StakingPoolMath.releasedStakeFromShare(
                    _request.unbondingShare,
                    _pool.validatorUnbondingPool.unbondingShare,
                    _pool.delegatorUnbondingPool.releasedStake
                );
                _pool.validatorUnbondingPool.unbondingShare -= _request.unbondingShare;
                _pool.delegatorUnbondingPool.releasedStake -= _share;
                _releasedStakes += _share;
            }
            _topIndex++;
        }

        // TODO (tariq): consider deleting unbonding request from `unbondingArray` as they are released
        _queue.queue.dequeue(_topIndex - _queue.queue.topIndex);
        if (_queue.queue.array.length < _queue.unlockingIndex) {
            _queue.unlockingIndex = _queue.queue.array.length;
        }
        if (_releasedStakes > 0) {
            releasedStakes -= _releasedStakes;
            autonity.updateWithReleasedStake(_delegator, _releasedStakes);
        }
    }

    function _rewardsFromLiquidMinted(address _delegator, address _validator, uint256 _epochID) internal view returns (uint256) {
        UintQueue storage _queue = pendingBondingQueue[_delegator];
        uint256 _length = _queue.array.length;
        uint256 _topIndex = _queue.topIndex;
        uint256 _rewards;

        while (_topIndex < _length) {
            BondingRequest storage _request = bondingArray[_queue.array[_topIndex]];
            if (_request.epochID == _epochID) {
                break;
            }
            if (_request.validator != _validator) {
                _topIndex++;
                continue;
            }
            if (_request.selfDelegation) {
                break;
            }

            PoolCollection storage _pool = epochStakesPool[_request.epochID][_validator];
            if (!_pool.validatorBondingPool.notActive) {
                uint256 _liquidEarned = StakingPoolMath.liquidFromStake(
                    _request.amount,
                    _pool.validatorBondingPool.delegatingStake,
                    _pool.delegatorBondingPool.liquidMinted
                );
                uint256 _totalRewards = _pool.delegatorBondingPool.rewardsCollected + StakingPoolMath.computeUnrealisedRewards(
                    lastFeeFactor[_validator],
                    _pool.delegatorBondingPool.feeFactor,
                    _pool.delegatorBondingPool.liquidMinted
                );
                _rewards += StakingPoolMath.computeRewardFraction(
                    _totalRewards,
                    _liquidEarned,
                    _pool.delegatorBondingPool.liquidMinted
                );
            }
            
            _topIndex++;
        }
        return _rewards;
    }

    function _rewardsFromLiquidBurning(address _delegator, address _validator, uint256 _epochID) internal view returns (uint256) {
        UnbondingQueue storage _queue = pendingUnbondingQueue[_delegator];
        uint256 _length = _queue.queue.array.length;
        uint256 _processingIndex = _queue.unlockingIndex;
        uint256 _rewards;

        while (_processingIndex < _length) {
            UnbondingRequest storage _request = unbondingArray[_queue.queue.array[_processingIndex]];
            if (_request.epochID == _epochID) {
                break;
            }
            if (_request.validator != _validator) {
                _processingIndex++;
                continue;
            }
            if (_request.selfDelegation) {
                break;
            }

            PoolCollection storage _pool = epochStakesPool[_request.epochID][_validator];
            _rewards += StakingPoolMath.computeRewardFraction(
                _pool.delegatorUnbondingPool.rewardsCollected,
                _request.amount,
                _pool.validatorUnbondingPool.liquidBurning
            );
            _processingIndex++;
        }
        return _rewards;
    }

    /*
    ============================================================

        Calculations

    ============================================================
     */
}