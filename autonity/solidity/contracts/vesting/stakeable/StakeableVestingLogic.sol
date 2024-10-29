// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../interfaces/IStakeableVesting.sol";
import {PendingStakingRequest, QueueLib} from "./QueueLib.sol";
import "./StakeableVestingStorage.sol";
import "./ValidatorManager.sol";

/**
 * @title Logic of the Stakeable Vesting Smart Contract for vesting and staking funds
 * @notice It does not support to act as a treasury account. So only delegated staking works with this.
 * @dev Only one smart contract is deployed by `StakeableVestingManager` which is used by separate accounts.
 */
contract StakeableVestingLogic is StakeableVestingStorage, ContractBase, ValidatorManager, IStakeableVesting {

    using QueueLib for StakingRequestQueue;

    constructor(address payable _autonity) AccessAutonity(_autonity) {
        managerContract = IStakeableVestingManager(payable(msg.sender));
    }

    /**
     * @notice Creates a new contract. Can be called once only.
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _startTime start time of the contract
     * @param _cliffDuration cliff duration of the contract
     * @param _totalDuration total duration of the contract
     * @custom:restricted-to StakeableVestingManager contract
     */
    function createContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration
    ) virtual external onlyManager {
        require(beneficiary == address(0), "contract already created");
        beneficiary = _beneficiary;
        stakeableContract = _createContract(_beneficiary, _amount, _startTime, _cliffDuration, _totalDuration, true);
        contractValuation = ContractValuation(_amount, 0);
    }

    /**
     * @notice Set the address of the manager contract.
     * @custom:restricted-to operator account
     */
    function setManagerContract(address _managerContract) virtual external onlyOperator {
        managerContract = IStakeableVestingManager(payable(_managerContract));
    }

    /**
     * @notice Used by beneficiary to transfer all vested NTN and LNTN to his own address.
     * When releasing funds, it tries to release everything from NTN balance first.
     * If the withdrawable vested funds is `v` NTN and NTN balance of the contract is `n`,
     * one of the following will happen
     * 
     *      1. `if (n >= v)`, all `v` NTN will be released from NTN balance and the
     *          remaining NTN balance will be `n-v` and no other asset is updated and the function exits.
     * 
     *      2. `if (n < v)`, all `n` NTN will be released from NTN balance and we update `v = v-n` and additional
     *          LNTN equivalent of `v` NTN will be released. See `releaseAllLNTN()` for how the LNTN will be released.
     * 
     * So before calling `releaseFunds()`, see the `linkedValidators` list using the function `getLinkedValidators()`.
     */
    function releaseFunds() virtual external onlyBeneficiary {
        _updateFunds();
        (uint256 _unlocked, uint256 _totalValue) = _withdrawableVestedFunds();
        // first NTN is released
        uint256 _remainingUnlocked = _releaseNTN(stakeableContract, _unlocked);
        // if there still remains some unlocked funds, i.e. not enough NTN, then LNTN is released
        _remainingUnlocked = _releaseAllVestedLNTN(_remainingUnlocked);
        _updateWithdrawnShare(_unlocked - _remainingUnlocked, _totalValue);
        _clearValidators();
    }

    /**
     * @notice Used by beneficiary to transfer all vested NTN to his own address.
     */
    function releaseAllNTN() virtual external onlyBeneficiary {
        _cleanup();
        (uint256 _unlocked, uint256 _totalValue) = _withdrawableVestedFunds();
        uint256 _remainingUnlocked = _releaseNTN(stakeableContract, _unlocked);
        _updateWithdrawnShare(_unlocked - _remainingUnlocked, _totalValue);
    }

    /**
     * @notice Used by beneficiary to transfer all vested LNTN to his own address.
     * If the withdrawable vested funds is `v` NTN, the LNTN will be released in the following order starting from `idx = 0`.
     * 
     *      1. Let `b` = unlocked LNTN balance of the contract for `linkedValidators[idx]` and `c` = equivalent NTN for `b` LNTN.
     *          `if (c >= v)`, then LNTN equivalent of `v` NTN will be released from `linkedValidators[idx]` and the function exits.
     * 
     *      2. `if (c < v)`, then all unlocked LNTN from `linkedValidators[idx]` will be released and we increase `idx = idx+1`
     *          and we update `v = v-c` and repeate from the process 1 if `idx < linkedValidators.length`.
     * 
     * So before calling `releaseAllLNTN()`, see the `linkedValidators` list using the function `getLinkedValidators()`.
     */
    function releaseAllLNTN() virtual external onlyBeneficiary {
        _updateFunds();
        (uint256 _unlocked, uint256 _totalValue) = _withdrawableVestedFunds();
        uint256 _remainingUnlocked = _releaseAllVestedLNTN(_unlocked);
        _updateWithdrawnShare(_unlocked - _remainingUnlocked, _totalValue);
        _clearValidators();
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of vested NTN to his own address.
     * @param _amount amount of NTN to transfer
     */
    function releaseNTN(uint256 _amount) virtual external onlyBeneficiary {
        _cleanup();
        (uint256 _unlocked, uint256 _totalValue) = _withdrawableVestedFunds();
        require(_amount <= _unlocked, "not enough unlocked funds");
        uint256 _remaining = _releaseNTN(stakeableContract, _amount);
        _updateWithdrawnShare(_amount - _remaining, _totalValue);
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of vested LNTN to his own address.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to transfer
     */
    function releaseLNTN(address _validator, uint256 _amount) virtual external onlyBeneficiary {
        require(_amount > 0, "require positive amount to transfer");
        _updateFunds();

        uint256 _unlockedLiquid = _unlockedLiquidBalance(_liquidStateContract(_validator));
        require(_unlockedLiquid >= _amount, "not enough unlocked LNTN");

        uint256 _value = _calculateLNTNValue(_validator, _amount);
        (uint256 _unlocked, uint256 _totalValue) = _withdrawableVestedFunds();
        require(_value <= _unlocked, "not enough unlocked funds");

        stakeableContract.withdrawnValue += _value;
        _transferLNTN(_amount, _validator);
        emit FundsReleased(beneficiary, address(_liquidStateContract(_validator)), _amount);
        _updateWithdrawnShare(_value, _totalValue);
        _clearValidators();
    }

    /**
     * @notice Changes the beneficiary of the contract to the `_recipient` address. The `_recipient` address can release and stake tokens from the contract.
     * Rewards which have been entitled to the beneficiary due to bonding from this contract are not transferred to `_recipient` address. Instead rewards
     * earned until this point from this contract are calculated and transferred to the old beneficiary address.
     * @param _recipient whome the contract is transferred to
     * @custom:restricted-to operator account
     */
    function changeContractBeneficiary(address _recipient) virtual external onlyManager {
        _claimAndSendRewards();
        _clearValidators();
        beneficiary = _recipient;
    }

    /**
     * @notice In case some funds are missing due to some pending staking operation that failed,
     * this function updates the funds by handling the pending requests.
     */
    function updateFunds() virtual external onlyBeneficiary {
        _updateFunds();
    }

    /**
     * @notice Updates the funds of the contract and returns the contract.
     */
    function updateFundsAndGetContract() external onlyBeneficiary returns (ContractBase.Contract memory) {
        _updateFunds();
        return stakeableContract;
    }

    /**
     * @notice Used by beneficiary to bond some NTN of the contract.
     * All bondings are delegated, as stakeable vesting smart contract cannot own a validator node.
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(address _validator, uint256 _amount) virtual external onlyBeneficiary returns (uint256) {
        uint256 _epochID = _getEpochID();
        bondingQueue.enqueue(
            PendingStakingRequest(_epochID, _validator, _amount, 0) // requestID not needed
        );
        _newBondingRequested(_validator, _epochID);
        return autonity.bond(_validator, _amount);
    }

    /**
     * @notice Used by beneficiary to unbond some LNTN of the contract from a validator.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(address _validator, uint256 _amount) virtual external onlyBeneficiary returns (uint256) {
        uint256 _unbondingID = autonity.unbond(_validator, _amount);
        uint256 _epochID = _getEpochID();
        unbondingQueue.enqueue(
            PendingStakingRequest(_epochID, _validator, 0, _unbondingID) // amount not needed
        );
        return _unbondingID;
    }

    /**
     * @notice Used by beneficiary to claim rewards from bonding to validator.
     * @param _validator validator address
     */
    function claimRewards(address _validator) virtual external onlyBeneficiary {
        _claimAndSendRewards(_validator);
        _clearValidators();
    }

    /**
     * @notice Used by beneficiary to claim all rewards from bonding to all the validators.
     */
    function claimRewards() virtual external onlyBeneficiary {
        _claimAndSendRewards();
        _clearValidators();
    }

    /**
     * @dev It is not expected to fall into the fallback function. Implemeted fallback() to get a proper reverting message.
     */
    fallback() payable external virtual {
        revert("fallback not implemented for StakeableVestingLogic");
    }

    /**
     * @dev Receive Auton function https://solidity.readthedocs.io/en/v0.7.2/contracts.html#receive-ether-function
     */
    receive() external payable {}

    /*
    ============================================================
         Internals
    ============================================================
     */

    /**
     * @dev Returns equivalent amount of NTN for some LNTN using the current ratio.
     * @param _validator validator address
     * @param _amount amount of LNTN to be converted
     */
    function _calculateLNTNValue(address _validator, uint256 _amount) internal view returns (uint256) {
        if (_amount == 0) {
            return 0;
        }
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake) / _validatorInfo.liquidSupply;
    }

    /**
     * @dev Returns equivalent amount of LNTN using the current ratio.
     * @param _validator validator address
     * @param _amount amount of NTN to be converted
     */
    function _getLiquidFromNTN(address _validator, uint256 _amount) internal view returns (uint256) {
        if (_amount == 0) {
            return 0;
        }
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * _validatorInfo.liquidSupply / (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake);
    }

    /**
     * @dev Calculates the total value of all the balances of the contract in NTN, which can vary if the contract has some LNTN.
     * `totalValue = currentNTN + (the value of LNTN converted to NTN using current ratio) + (newton going under bonding or unbonding)`
     */
    function _calculateTotalValue() internal view returns (uint256) {
        address _validator;
        uint256 _balance;
        uint256 _totalValue = autonity.balanceOf(address(this));
        uint256 _length = linkedValidators.length;
        for (uint256 i = 0; i < _length; i++) {
            _validator = linkedValidators[i];
            _balance = liquidBalance(_validator);
            _totalValue += _calculateLNTNValue(_validator, _balance);
        }
        return _totalValue + _calculateNewtonUnderBonding() + _calculateNewtonUnderUnbonding();
    }

    /**
     * @dev Transfers vested LNTN to beneficiary address. The amount of vested funds is calculated in NTN
     * and then converted to LNTN using the current ratio.
     * In case the contract has LNTN to multiple validators, we pick one validator and try to transfer
     * as much LNTN as possible. If there still remains some more vested funds, then we pick another validator.
     * There is no particular order in which validator should be picked first.
     */
    function _releaseAllVestedLNTN(
        uint256 _availableUnlockedFunds
    ) internal returns (uint256 _remaining) {
        _remaining = _availableUnlockedFunds;
        address _validator;
        uint256 _balance;
        uint256 _value;
        uint256 _liquid;
        uint256 _length = linkedValidators.length;
        for (uint256 i = 0; i < _length && _remaining > 0; i++) {
            _validator = linkedValidators[i];
            _balance = _unlockedLiquidBalance(_liquidStateContract(_validator));
            if (_balance == 0) {
                continue;
            }
            _value = _calculateLNTNValue(_validator, _balance);
            if (_remaining >= _value) {
                _remaining -= _value;
                _transferLNTN(_balance, _validator);
                emit FundsReleased(msg.sender, address(_liquidStateContract(_validator)), _balance);
            }
            else {
                _liquid = _getLiquidFromNTN(_validator, _remaining);
                require(_liquid <= _balance, "conversion not working");
                _remaining = 0;
                _transferLNTN(_liquid, _validator);
                emit FundsReleased(msg.sender, address(_liquidStateContract(_validator)), _liquid);
            }
        }
        stakeableContract.withdrawnValue += _availableUnlockedFunds - _remaining;
    }

    function _updateWithdrawnShare(uint256 _withdrawnValue, uint256 _totalValue) internal {
        if (_withdrawnValue == 0) {
            return;
        }
        uint256 _alreadyWithdrawn = contractValuation.withdrawnShare;
        uint256 _withdrawnShare = (_withdrawnValue * (contractValuation.totalShare - _alreadyWithdrawn)) / _totalValue;
        contractValuation.withdrawnShare = _alreadyWithdrawn + _withdrawnShare;
    }

    function _withdrawableVestedFunds() internal view returns (uint256 _unlockedValue, uint256 _totalValue) {
        if (autonity.lastEpochTime() < stakeableContract.start + stakeableContract.cliffDuration) {
            return (0, 0);
        }
        return _vestedFunds();
    }

    /**
     * @dev Calculates the amount of vested funds in NTN until last epoch time.
     */
    function _vestedFunds() internal view returns (uint256 _unlockedValue, uint256 _totalValue) {
        uint256 _time = autonity.lastEpochTime();
        uint256 _start = stakeableContract.start;
        if (_time < _start) {
            return (0, 0);
        }

        uint256 _totalDuration = stakeableContract.totalDuration;
        uint256 _totalShare = contractValuation.totalShare;
        uint256 _withdrawnShare = contractValuation.withdrawnShare;
        uint256 _unlockedShare;
        if (_start + _totalDuration <= _time) {
            _unlockedShare = _totalShare - _withdrawnShare;
        }
        else {
            _unlockedShare = (_totalShare * (_time - _start)) / _totalDuration - _withdrawnShare;
        }
        _totalValue = _calculateTotalValue();
        if (_unlockedShare > 0) {
            _unlockedValue = (_totalValue * _unlockedShare) / (_totalShare - _withdrawnShare);
        }
    }

    function _transferLNTN(uint256 _amount, address _validator) internal {
        bool _sent = _liquidStateContract(_validator).transfer(beneficiary, _amount);
        require(_sent, "LNTN transfer failed");
    }

    function _sendRewards(uint256 _atnReward, uint256 _ntnReward) internal {
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool _sent, ) = beneficiary.call{value: _atnReward}("");
        require(_sent, "failed to send ATN");

        _transferNTN(beneficiary, _ntnReward);
    }

    /**
     * @dev Updates the funds by processing the staking requests.
     */
    function _updateFunds() internal {
        _handlePendingBondingRequest();
        _handlePendingUnbondingRequest();
        stakeableContract.currentNTNAmount = autonity.balanceOf(address(this));
    }

    /**
     * @dev Updates the funds and removes any unnecessary validator from the list.
     */
    function _cleanup() internal {
        _updateFunds();
        _clearValidators();
    }

    /**
     * @dev Handles all the pending bonding requests.
     * All the requests from past epoch can be deleted as the bonding requests are
     * applied at epoch end immediately. Requests from current epoch are still pending.
     */
    function _handlePendingBondingRequest() internal {
        PendingStakingRequest storage _bondingRequest;
        uint256 _currentEpochID = _getEpochID();
        PendingStakingRequest[] storage _queue = bondingQueue.array;
        uint256 _length = _queue.length;
        uint256 _topIndex = bondingQueue.topIndex;

        // delete all bonding requests from the past epoch
        while (_topIndex < _length) {
            _bondingRequest = _queue[_topIndex];
            if (_bondingRequest.epochID < _currentEpochID) {
                _bondingRequestExpired(_bondingRequest.validator, _bondingRequest.epochID);
                _topIndex++;
            }
            else break;
        }
        bondingQueue.dequeue(_topIndex - bondingQueue.topIndex);
    }

    function _calculateNewtonUnderBonding() internal view returns (uint256) {
        uint256 _bondingNTN;
        PendingStakingRequest storage _bondingRequest;
        uint256 _currentEpochID = _getEpochID();
        PendingStakingRequest[] storage _queue = bondingQueue.array;
        uint256 _length = _queue.length;

        for (uint256 i = bondingQueue.topIndex; i < _length; i++) {
            _bondingRequest = _queue[i];
            if (_bondingRequest.epochID < _currentEpochID) {
                continue;
            }
            _bondingNTN += _bondingRequest.amount;
        }
        return _bondingNTN;
    }

    /**
     * @dev Handles all the pending unbonding requests. All unbonding requests from past epoch are applied.
     * Unbonding request that are released in Autonity can be deleted.
     */
    function _handlePendingUnbondingRequest() internal {
        PendingStakingRequest storage _unbondingRequest;
        PendingStakingRequest[] storage _queue = unbondingQueue.array;
        uint256 _length = _queue.length;
        uint256 _topIndex = unbondingQueue.topIndex;

        // first delete all unbonding request from queue that are released
        while (_topIndex < _length) {
            _unbondingRequest = _queue[_topIndex];
            if (autonity.isUnbondingReleased(_unbondingRequest.requestID)) {
                _topIndex++;
            }
            else {
                break;
            }
        }
        unbondingQueue.dequeue(_topIndex - unbondingQueue.topIndex);
    }

    function _calculateNewtonUnderUnbonding() internal view returns (uint256) {
        uint256 _unbondingNTN;
        uint256 _unbondingShare;
        PendingStakingRequest storage _unbondingRequest;
        Autonity.Validator memory _validator;
        uint256 _currentEpochID = _getEpochID();
        PendingStakingRequest[] storage _queue = unbondingQueue.array;
        uint256 _length = _queue.length;

        for (uint256 i = unbondingQueue.topIndex; i < _length; i++) {
            _unbondingRequest = _queue[i];
            if (_unbondingRequest.epochID == _currentEpochID) {
                break;
            }
            if (autonity.isUnbondingReleased(_unbondingRequest.requestID)) {
                continue;
            }
            _unbondingShare = autonity.getUnbondingShare(_unbondingRequest.requestID);
            if (_unbondingShare == 0) {
                continue;
            }
            _validator = autonity.getValidator(_unbondingRequest.validator);
            _unbondingNTN += (_unbondingShare * _validator.unbondingStake) / _validator.unbondingShares;
        }
        return _unbondingNTN;
    }

    /**
     * @dev Claims all rewards from the liquid contract of the validator.
     * @param _validator validator address
     */
    function _claimAndSendRewards(address _validator) internal {
        address _myAddress = address(this);
        uint256 _atnBalance = _myAddress.balance;
        uint256 _ntnBalance = autonity.balanceOf(_myAddress);
        _liquidStateContract(_validator).claimRewards();
        _sendRewards(_myAddress.balance - _atnBalance, autonity.balanceOf(_myAddress) - _ntnBalance);
    }

    /**
     * @dev Claims all rewards from the liquid contract from all bonded validators.
     */
    function _claimAndSendRewards() internal {
        address _myAddress = address(this);
        uint256 _atnBalance = _myAddress.balance;
        uint256 _ntnBalance = autonity.balanceOf(_myAddress);
        uint256 _length = linkedValidators.length;
        for (uint256 i = 0; i < _length; i++) {
            _liquidStateContract(linkedValidators[i]).claimRewards();
        }
        _sendRewards(_myAddress.balance - _atnBalance, autonity.balanceOf(_myAddress) - _ntnBalance);
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns unclaimed rewards from bonding to validator.
     * @param _validator validator address
     * @return _atnRewards unclaimed ATN rewards
     * @return _ntnRewards unclaimed NTN rewards
     */
    function unclaimedRewards(address _validator) virtual external view returns (uint256 _atnRewards, uint256 _ntnRewards) {
        (_atnRewards, _ntnRewards) = _unclaimedRewards(_validator);
    }

    /**
     * @notice Returns the amount of all unclaimed rewards due to all the bonding from the contract entitled to beneficiary.
     */
    function unclaimedRewards() virtual external view returns (uint256 _atnRewards, uint256 _ntnRewards) {
        for (uint256 i = 0; i < linkedValidators.length; i++) {
            (uint256 _atn, uint256 _ntn) = _unclaimedRewards(linkedValidators[i]);
            _atnRewards += _atn;
            _ntnRewards += _ntn;
        }
    }

    /**
     * @notice Returns the amount of vested funds and withdrawable in NTN.
     */
    function withdrawableVestedFunds() virtual external view returns (uint256) {
        (uint256 _unlocked, ) = _withdrawableVestedFunds();
        return _unlocked;
    }

    /**
     * @notice Returns the amount of vested funds in NTN.
     */
    function vestedFunds() virtual external view returns (uint256) {
        (uint256 _unlocked, ) = _vestedFunds();
        return _unlocked;
    }

    /**
     * @notice Returns the current total value of the contract in NTN.
     */
    function contractTotalValue() external view returns (uint256) {
        return _calculateTotalValue();
    }

    /**
     * @notice Returns the contract.
     */
    function getContract() virtual external view returns (ContractBase.Contract memory) {
        return stakeableContract;
    }

    /**
     * @notice Returns the address of the `StakeableVestingManager` smart contract.
     */
    function getManagerContractAddress() virtual external view returns (address) {
        return address(managerContract);
    }

    /**
     * @notice Returns the beneficiary address of the contract.
     */
    function getBeneficiary() virtual external view returns (address) {
        return beneficiary;
    }

    /**
     * @notice Returns the list of validators that are bonded to or have some unclaimed rewards.
     */
    function getLinkedValidators() virtual external view returns (address[] memory) {
        return linkedValidators;
    }

    /**
     * @notice Returns the amount of LNTN bonded to `_validator` from the contract.
     * @param _validator validator address
     */
    function liquidBalance(address _validator) virtual public view returns (uint256) {
        return _liquidBalance(_getLiquidStateContract(_validator));
    }

    /**
     * @notice Returns the amount of unlocked (not vested) LNTN bonded to `_validator` from the contract.
     * @param _validator validator address
     */
    function unlockedLiquidBalance(address _validator) virtual external view returns (uint256) {
        return _unlockedLiquidBalance(_getLiquidStateContract(_validator));
    }

    /**
     * @notice Returns the amount of locked (not unvested) LNTN bonded to `_validator` from the contract.
     * @param _validator validator address
     */
    function lockedLiquidBalance(address _validator) virtual external view returns (uint256) {
        return _lockedLiquidBalance(_getLiquidStateContract(_validator));
    }

    /*
    ============================================================

        Modifiers

    ============================================================
     */

    modifier onlyBeneficiary {
        require(msg.sender == beneficiary, "caller is not beneficiary of the contract");
        _;
    }

    modifier onlyManager {
        require(msg.sender == address(managerContract), "caller is not the manager");
        _;
    }

}
