// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../../interfaces/IStakableVesting.sol";
import "./StakableVestingStorage.sol";
import "./ValidatorManager.sol";

contract StakableVestingLogic is StakableVestingStorage, ContractBase, ValidatorManager, IStakableVesting {

    constructor(address payable _autonity) AccessAutonity(_autonity) {
        managerContract = StakableVestingManager(payable(msg.sender));
    }

    function createContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration
    ) virtual external onlyManager {
        require(beneficiary == address(0), "contract already created");
        require(_beneficiary != address(0), "beneficiary is not a valid address");
        beneficiary = _beneficiary;
        stakableContract = _createContract(_amount, _startTime, _cliffDuration, _totalDuration, true);
        contractValuation = ContractValuation(_amount, 0);
    }

    function setManagerContract(address _managerContract) virtual external onlyOperator {
        managerContract = StakableVestingManager(payable(_managerContract));
    }


    /**
     * @notice Used by beneficiary to transfer all unlocked NTN and LNTN of some contract to his own address.
     */
    function releaseFunds() virtual external onlyBeneficiary {
        _updateFunds();
        (uint256 _unlocked, uint256 _totalValue) = _unlockedFunds();
        // first NTN is released
        uint256 _remainingUnlocked = _releaseNTN(stakableContract, _unlocked);
        // if there still remains some unlocked funds, i.e. not enough NTN, then LNTN is released
        _remainingUnlocked = _releaseAllUnlockedLNTN(_remainingUnlocked);
        _updateWithdrawnShare(_unlocked - _remainingUnlocked, _totalValue);
        _clearValidators();
    }

    /**
     * @notice Used by beneficiary to transfer all unlocked NTN of some contract to his own address.
     */
    function releaseAllNTN() virtual external onlyBeneficiary {
        _cleanup();
        (uint256 _unlocked, uint256 _totalValue) = _unlockedFunds();
        uint256 _remainingUnlocked = _releaseNTN(stakableContract, _unlocked);
        _updateWithdrawnShare(_unlocked - _remainingUnlocked, _totalValue);
    }

    /**
     * @notice Used by beneficiary to transfer all unlocked LNTN of some contract to his own address.
     */
    function releaseAllLNTN() virtual external onlyBeneficiary {
        _updateFunds();
        (uint256 _unlocked, uint256 _totalValue) = _unlockedFunds();
        uint256 _remainingUnlocked = _releaseAllUnlockedLNTN(_unlocked);
        _updateWithdrawnShare(_unlocked - _remainingUnlocked, _totalValue);
        _clearValidators();
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address.
     * @param _amount amount to transfer
     */
    function releaseNTN(uint256 _amount) virtual external onlyBeneficiary {
        _cleanup();
        (uint256 _unlocked, uint256 _totalValue) = _unlockedFunds();
        require(_amount <= _unlocked, "not enough unlocked funds");
        uint256 _remainingUnlocked = _releaseNTN(stakableContract, _amount);
        _updateWithdrawnShare(_amount - _remainingUnlocked, _totalValue);
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of unlocked LNTN of some contract to his own address.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to transfer
     */
    function releaseLNTN(address _validator, uint256 _amount) virtual external onlyBeneficiary {
        require(_amount > 0, "require positive amount to transfer");
        _updateFunds();

        uint256 _unlockedLiquid = _unlockedLiquidBalance(_liquidStateContract(_validator));
        require(_unlockedLiquid >= _amount, "not enough unlocked LNTN");

        uint256 _value = _calculateLNTNValue(_validator, _amount);
        (uint256 _unlocked, uint256 _totalValue) = _unlockedFunds();
        require(_value <= _unlocked, "not enough unlocked funds");

        stakableContract.withdrawnValue += _value;
        _transferLNTN(_amount, _validator);
        _updateWithdrawnShare(_value, _totalValue);
        _clearValidators();
    }

    /**
     * @notice Changes the beneficiary of some contract to the recipient address. The recipient address can release and stake tokens from the contract.
     * Rewards which have been entitled to the beneficiary due to bonding from this contract are not transferred to recipient.
     * @dev Rewards earned until this point from this contract are calculated and stored in atnRewards and ntnRewards mapping so that
     * beneficiary can later claim them even though beneficiary is not entitled to this contract.
     * @param _recipient whome the contract is transferred to
     * @custom:restricted-to operator account
     */
    function changeContractBeneficiary(address _recipient) virtual external onlyManager {
        _updateFunds();
        _claimAndSendRewards();
        _clearValidators();
        beneficiary = _recipient;
    }

    /**
     * @notice In case some funds are missing due to some pending staking operation that failed,
     * this function updates the funds of some contract entitled to beneficiary by applying the pending requests.
     */
    function updateFunds() virtual external {
        _updateFunds();
    }

    /**
     * @notice Updates the funds of the contract and returns the contract.
     */
    function updateFundsAndGetContract() external returns (ContractBase.Contract memory) {
        _updateFunds();
        return stakableContract;
    }

    /**
     * @notice Used by beneficiary to bond some NTN of some contract.
     * All bondings are delegated, as vesting manager cannot own a validator.
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(address _validator, uint256 _amount) virtual external onlyBeneficiary returns (uint256) {
        require(stakableContract.start <= block.timestamp, "contract not started yet");
        uint256 _epochID = _getEpochID();
        bondingQueue.push(PendingBondingRequest(_amount, _epochID, _validator));
        _newBondingRequested(_validator, _epochID);
        return autonity.bond(_validator, _amount);
    }

    /**
     * @notice Used by beneficiary to unbond some LNTN of some contract.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(address _validator, uint256 _amount) virtual external onlyBeneficiary returns (uint256) {
        uint256 _unbondingID = autonity.unbond(_validator, _amount);
        uint256 _epochID = _getEpochID();
        unbondingQueue.push(PendingUnbondingRequest(_unbondingID, _epochID, _validator));
        _newUnbondingRequested(_validator, _epochID, _unbondingID);
        return _unbondingID;
    }

    /**
     * @notice Used by beneficiary to claim rewards from bonding some contract to validator.
     * @param _validator validator address
     */
    function claimRewards(address _validator) virtual external onlyBeneficiary {
        _claimAndSendRewards(_validator);
        _clearValidators();
    }

    /**
     * @notice Used by beneficiary to claim rewards from bonding some contract to validator.
     */
    function claimRewards() virtual external onlyBeneficiary {
        _claimAndSendRewards();
        _clearValidators();
    }

    /**
     * @dev It is not expected to fall into the fallback function. Implemeted fallback() to get a reverting message.
     */
    fallback() payable external virtual {
        revert("fallback not implemented for StakableVestingLogic");
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
     * @dev Returns equivalent amount of NTN using the current ratio.
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
     * `totalValue = currentNTN + (the value of LNTN converted to NTN using current ratio)`
     */
    function _calculateTotalValue() internal view returns (uint256) {
        address _validator;
        uint256 _balance;
        uint256 _totalValue = autonity.balanceOf(address(this));
        for (uint256 i = 0; i < bondedValidators.length; i++) {
            _validator = bondedValidators[i];
            _balance = liquidBalance(_validator);
            _totalValue += _calculateLNTNValue(_validator, _balance);
        }
        return _totalValue + _calculateNewtonUnderBonding() + _calculateNewtonUnderUnbonding();
    }

    /**
     * @dev Transfers some LNTN equivalent to beneficiary address. The amount of unlocked funds is calculated in NTN
     * and then converted to LNTN using the current ratio.
     * In case the contract has LNTN to multiple validators, we pick one validator and try to transfer
     * as much LNTN as possible. If there still remains some more uncloked funds, then we pick another validator.
     * There is no particular order in which validator should be picked first.
     */
    function _releaseAllUnlockedLNTN(
        uint256 _availableUnlockedFunds
    ) internal returns (uint256 _remaining) {
        _remaining = _availableUnlockedFunds;
        address _validator;
        uint256 _balance;
        uint256 _value;
        uint256 _liquid;
        for (uint256 i = 0; i < bondedValidators.length && _remaining > 0; i++) {
            _validator = bondedValidators[i];
            _balance = _unlockedLiquidBalance(_liquidStateContract(_validator));
            if (_balance == 0) {
                continue;
            }
            _value = _calculateLNTNValue(_validator, _balance);
            if (_remaining >= _value) {
                _remaining -= _value;
                _transferLNTN(_balance, _validator);
            }
            else {
                _liquid = _getLiquidFromNTN(_validator, _remaining);
                require(_liquid <= _balance, "conversion not working");
                _remaining = 0;
                _transferLNTN(_liquid, _validator);
            }
        }
        stakableContract.withdrawnValue += _availableUnlockedFunds - _remaining;
    }

    function _updateWithdrawnShare(uint256 _withdrawnValue, uint256 _totalValue) internal {
        if (_withdrawnValue == 0) {
            return;
        }
        uint256 _alreadyWithdrawn = contractValuation.withdrawnShare;
        uint256 _withdrawnShare = (_withdrawnValue * (contractValuation.totalShare - _alreadyWithdrawn)) / _totalValue;
        contractValuation.withdrawnShare = _alreadyWithdrawn + _withdrawnShare;
    }

    /**
     * @dev Calculates the amount of unlocked funds in NTN until last epoch time.
     */
    function _unlockedFunds() internal view returns (uint256 _unlockedValue, uint256 _totalValue) {
        uint256 _time = autonity.lastEpochTime();
        uint256 _start = stakableContract.start;
        require(_time >= _start + stakableContract.cliffDuration, "cliff period not reached yet");

        uint256 _totalDuration = stakableContract.totalDuration;
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

    /**
     * @dev Given the total value (in NTN) of the contract, calculates the amount of withdrawable tokens (in NTN).
     */
    function _calculateAvailableUnlockedFunds(
        Contract storage _contract, uint256 _totalValue, uint256 _time
    ) internal view returns (uint256) {
        require(_time >= _contract.start + _contract.cliffDuration, "cliff period not reached yet");

        uint256 _unlocked = _calculateTotalUnlockedFunds(_contract.start, _contract.totalDuration, _time, _totalValue);
        if (_unlocked > _contract.withdrawnValue) {
            return _unlocked - _contract.withdrawnValue;
        }
        return 0;
    }

    /**
     * @dev Calculates total unlocked funds while assuming cliff period has passed.
     * Check if cliff is passed before calling this function.
     */
    function _calculateTotalUnlockedFunds(
        uint256 _start, uint256 _totalDuration, uint256 _time, uint256 _totalAmount
    ) internal pure returns (uint256) {
        if (_time >= _totalDuration + _start) {
            return _totalAmount;
        }
        return (_totalAmount * (_time - _start)) / _totalDuration;
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
     * @dev Updates the funds by applying the staking requests.
     */
    function _updateFunds() internal {
        _handlePendingBondingRequest();
        _handlePendingUnbondingRequest();
        stakableContract.currentNTNAmount = autonity.balanceOf(address(this));
    }

    /**
     * @dev Updates the funds and removes any unnecessary validator from the list.
     */
    function _cleanup() internal {
        _updateFunds();
        _clearValidators();
    }

    function _deleteBondingRequest(PendingBondingRequest storage _bondingRequest) private {
        _bondingRequest.amount = 0;
        _bondingRequest.epochID = 0;
        _bondingRequest.validator = address(0);
    }

    /**
     * @dev Handles all the pending bonding requests.
     * All the requests from past epoch can be deleted as the bonding requests are
     * applied at epoch end immediately. Requests from current epoch are used to calculated `bondingNTN`.
     */
    function _handlePendingBondingRequest() internal {
        PendingBondingRequest storage _bondingRequest;
        uint256 _currentEpochID = _getEpochID();
        uint256 _length = bondingQueue.length;
        uint256 _topIndex = bondingQueueTopIndex;

        // delete all bonding requests from the past epoch
        for (uint256 i = _topIndex; i < _length; i++) {
            _bondingRequest = bondingQueue[i];
            if (_bondingRequest.epochID < _currentEpochID) {
                _bondingRequestExpired(_bondingRequest.validator, _bondingRequest.epochID);
                _deleteBondingRequest(_bondingRequest);
                _topIndex++;
            }
            else break;
        }
        bondingQueueTopIndex = _topIndex;
    }

    function _calculateNewtonUnderBonding() internal view returns (uint256) {
        uint256 _bondingNTN;
        PendingBondingRequest storage _bondingRequest;
        uint256 _currentEpochID = _getEpochID();
        uint256 _length = bondingQueue.length;

        for (uint256 i = bondingQueueTopIndex; i < _length; i++) {
            _bondingRequest = bondingQueue[i];
            if (_bondingRequest.epochID < _currentEpochID) {
                continue;
            }
            _bondingNTN += _bondingRequest.amount;
        }
        return _bondingNTN;
    }

    function _deleteUnbondingRequest(PendingUnbondingRequest storage _unbondingRequest) private {
        _unbondingRequest.unbondingID = 0;
        _unbondingRequest.epochID = 0;
        _unbondingRequest.validator = address(0);
    }

    /**
     * @dev Handles all the pending unbonding requests. All unbonding requests from past epoch are applied.
     * Unbonding request that are released in Autonity can be deleted.
     */
    function _handlePendingUnbondingRequest() internal {
        PendingUnbondingRequest storage _unbondingRequest;
        uint256 _length = unbondingQueue.length;
        uint256 _topIndex = unbondingQueueTopIndex;

        // first delete all unbonding request from queue that are released
        for (uint256 i = _topIndex; i < _length; i++) {
            _unbondingRequest = unbondingQueue[i];
            if (autonity.isUnbondingReleased(_unbondingRequest.unbondingID)) {
                _unbondingRequestExpired(_unbondingRequest.validator, _unbondingRequest.unbondingID);
                _deleteUnbondingRequest(_unbondingRequest);
                _topIndex++;
            }
            else {
                break;
            }
        }
        unbondingQueueTopIndex = _topIndex;
    }

    function _calculateNewtonUnderUnbonding() internal view returns (uint256) {
        uint256 _unbondingNTN;
        uint256 _unbondingShare;
        PendingUnbondingRequest storage _unbondingRequest;
        Autonity.Validator memory _validator;
        uint256 _currentEpochID = _getEpochID();
        uint256 _length = unbondingQueue.length;

        for (uint256 i = unbondingQueueTopIndex; i < _length; i++) {
            _unbondingRequest = unbondingQueue[i];
            if (_unbondingRequest.epochID == _currentEpochID) {
                break;
            }
            if (autonity.isUnbondingReleased(_unbondingRequest.unbondingID)) {
                continue;
            }
            _unbondingShare = autonity.getUnbondingShare(_unbondingRequest.unbondingID);
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
        for (uint256 i = 0; i < bondedValidators.length; i++) {
            _liquidStateContract(bondedValidators[i]).claimRewards();
        }
        _sendRewards(_myAddress.balance - _atnBalance, autonity.balanceOf(_myAddress) - _ntnBalance);
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns unclaimed rewards from some contract entitled to beneficiary from bonding to validator.
     * @param _validator validator address
     * @return _atnRewards unclaimed ATN rewards
     * @return _ntnRewards unclaimed NTN rewards
     */
    function unclaimedRewards(address _validator) virtual external view returns (uint256 _atnRewards, uint256 _ntnRewards) {
        (_atnRewards, _ntnRewards) = _unclaimedRewards(_validator);
    }

    /**
     * @notice Returns the amount of all unclaimed rewards due to all the bonding from contracts entitled to beneficiary.
     */
    function unclaimedRewards() virtual external view returns (uint256 _atnRewards, uint256 _ntnRewards) {
        for (uint256 i = 0; i < bondedValidators.length; i++) {
            (uint256 _atn, uint256 _ntn) = _unclaimedRewards(bondedValidators[i]);
            _atnRewards += _atn;
            _ntnRewards += _ntn;
        }
    }

    /**
     * @notice Returns the amount of released funds in NTN for some contract.
     */
    function unlockedFunds() virtual external view returns (uint256) {
        (uint256 _unlocked, ) = _unlockedFunds();
        return _unlocked;
    }

    /**
     * @notice Returns the current total value of the contract in NTN.
     */
    function contractTotalValue() external view returns (uint256) {
        return _calculateTotalValue();
    }

    function getContract() virtual external view returns (ContractBase.Contract memory) {
        return stakableContract;
    }

    function getManagerContractAddress() virtual external view returns (address) {
        return address(managerContract);
    }

    function getBeneficiary() virtual external view returns (address) {
        return beneficiary;
    }

    /**
     * @notice Returns the list of validators bonded some contract.
     */
    function getLinkedValidators() virtual external view returns (address[] memory) {
        return bondedValidators;
    }

    /**
     * @notice Returns the amount of LNTN for some contract.
     * @param _validator validator address
     */
    function liquidBalance(address _validator) virtual public view returns (uint256) {
        ILiquidLogic _liquidContract = validators[_validator].liquidStateContract;
        if (address(_liquidContract) == address(0)) {
            _liquidContract = autonity.getValidator(_validator).liquidStateContract;
        }
        return _liquidBalance(_liquidContract);
    }

    /**
     * @notice Returns the amount of unlocked LNTN for some contract.
     * @param _validator validator address
     */
    function unlockedLiquidBalance(address _validator) virtual external view returns (uint256) {
        ILiquidLogic _liquidContract = validators[_validator].liquidStateContract;
        if (address(_liquidContract) == address(0)) {
            _liquidContract = autonity.getValidator(_validator).liquidStateContract;
        }
        return _unlockedLiquidBalance(_liquidContract);
    }

    /**
     * @notice Returns the amount of locked LNTN for some contract.
     * @param _validator validator address
     */
    function lockedLiquidBalance(address _validator) virtual external view returns (uint256) {
        ILiquidLogic _liquidContract = validators[_validator].liquidStateContract;
        if (address(_liquidContract) == address(0)) {
            _liquidContract = autonity.getValidator(_validator).liquidStateContract;
        }
        return _lockedLiquidBalance(_liquidContract);
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
