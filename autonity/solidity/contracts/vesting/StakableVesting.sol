// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./ContractBase.sol";
import "./ValidatorManager.sol";

contract StakableVesting is ContractBase, ValidatorManager {
    
    address private beneficiary;
    address private managerContract;
    ContractBase.Contract private stakableContract;

    struct ContractValuation {
        uint256 totalShare;
        uint256 withdrawnShare;
    }

    ContractValuation private contractValuation;

    struct StakingNTN {
        uint256 bondingNTN;
        uint256 unbondingNTN;
    }

    StakingNTN private ntnUnderStaking;

    struct PendingBondingRequest {
        uint256 amount;
        uint256 epochID;
        address validator;
    }

    PendingBondingRequest[] internal bondingQueue;
    uint256 private bondingQueueTopIndex;

    struct PendingUnbondingRequest {
        uint256 unbondingID;
        uint256 epochID;
        address validator;
    }

    PendingUnbondingRequest[] internal unbondingQueue;
    uint256 private unbondingQueueTopIndex;

    constructor(
        address payable _autonity,
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration
    ) AccessAutonity(_autonity) {
        beneficiary = _beneficiary;
        managerContract = msg.sender;
        stakableContract = _createContract(_amount, _startTime, _cliffDuration, _totalDuration, true);
        contractValuation = ContractValuation(_amount, 0);
    }


    /**
     * @notice Used by beneficiary to transfer all unlocked NTN and LNTN of some contract to his own address.
     */
    function releaseFunds() virtual external onlyBeneficiary {
        _updateFunds();
        uint256 _unlocked = _unlockedFunds();
        // first NTN is released
        _unlocked = _releaseNTN(stakableContract, _unlocked);
        // if there still remains some unlocked funds, i.e. not enough NTN, then LNTN is released
        _releaseAllUnlockedLNTN(_unlocked);
        _clearValidators();
    }

    /**
     * @notice Used by beneficiary to transfer all unlocked NTN of some contract to his own address.
     */
    function releaseAllNTN() virtual external onlyBeneficiary {
        _cleanup();
        _releaseNTN(stakableContract, _unlockedFunds());
    }

    /**
     * @notice Used by beneficiary to transfer all unlocked LNTN of some contract to his own address.
     */
    function releaseAllLNTN() virtual external onlyBeneficiary {
        _updateFunds();
        _releaseAllUnlockedLNTN(_unlockedFunds());
        _clearValidators();
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address.
     * @param _amount amount to transfer
     */
    function releaseNTN(uint256 _amount) virtual external onlyBeneficiary {
        _cleanup();
        require(_amount <= _unlockedFunds(), "not enough unlocked funds");
        _releaseNTN(stakableContract, _amount);
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

        uint256 _unlockedLiquid = unlockedLiquidBalanceOf(_validator);
        require(_unlockedLiquid >= _amount, "not enough unlocked LNTN");

        uint256 _value = _calculateLNTNValue(_validator, _amount);
        require(_value <= _unlockedFunds(), "not enough unlocked funds");

        stakableContract.withdrawnValue += _value;
        _transferLNTN(_amount, _validator);
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
    function updateFunds() virtual external onlyBeneficiary {
        _updateFunds();
    }

    /**
     * @notice Updates the funds of the contract and returns total value of the contract.
     */
    function updateFundsAndGetContractTotalValue() external onlyBeneficiary returns (uint256) {
        _updateFunds();
        return _calculateTotalValue();
    }

    /**
     * @notice Updates the funds of the contract and returns the contract.
     */
    function updateFundsAndGetContract() external onlyBeneficiary returns (ContractBase.Contract memory) {
        _updateFunds();
        return stakableContract;
    }

    /**
     * @notice Used by beneficiary to bond some NTN of some contract.
     * All bondings are delegated, as vesting manager cannot own a validator.
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(address _validator, uint256 _amount) virtual public onlyBeneficiary returns (uint256) {
        require(stakableContract.start <= block.timestamp, "contract not started yet");
        bondingQueue.push(PendingBondingRequest(_amount, _getEpochID(), _validator));
        _newBondingRequested(_validator);
        return autonity.bond(_validator, _amount);
    }

    /**
     * @notice Used by beneficiary to unbond some LNTN of some contract.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(address _validator, uint256 _amount) virtual public onlyBeneficiary returns (uint256) {
        uint256 _unbondingID = autonity.unbond(_validator, _amount);
        unbondingQueue.push(PendingUnbondingRequest(_unbondingID, _getEpochID(), _validator));
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
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake) / _validatorInfo.liquidSupply;
    }

    /**
     * @dev Returns equivalent amount of LNTN using the current ratio.
     * @param _validator validator address
     * @param _amount amount of NTN to be converted
     */
    function _getLiquidFromNTN(address _validator, uint256 _amount) internal view returns (uint256) {
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * _validatorInfo.liquidSupply / (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake);
    }

    /**
     * @dev Calculates the total value of the contract, which can vary if the contract has some LNTN.
     * `totalValue = currentNTN + withdrawnValue + (the value of LNTN converted to NTN using current ratio)`
     */
    function _calculateTotalValue() internal view returns (uint256) {
        uint256 _totalValue = autonity.balanceOf(address(this)) + stakableContract.withdrawnValue;
        address _validator;
        uint256 _balance;
        for (uint256 i = 0; i < bondedValidators.length; i++) {
            _validator = bondedValidators[i];
            _balance = liquidBalanceOf(_validator);
            if (_balance == 0) {
                continue;
            }
            _totalValue += _calculateLNTNValue(_validator, _balance);
        }
        return _totalValue;
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
            _balance = unlockedLiquidBalanceOf(_validator);
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

    /**
     * @dev Calculates the amount of unlocked funds in NTN until last epoch time.
     */
    function _unlockedFunds() internal view returns (uint256) {
        return _calculateAvailableUnlockedFunds(
            stakableContract, _calculateTotalValue(), autonity.lastEpochTime()
        );
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

        // first delete all bonding requests from the queue
        for (uint256 i = _topIndex; i < _length; i++) {
            _bondingRequest = bondingQueue[i];
            if (_bondingRequest.epochID < _currentEpochID) {
                _bondingRequestExpired(_bondingRequest.validator);
                _deleteBondingRequest(_bondingRequest);
                _topIndex++;
            }
            else break;
        }

        // now calculate total `ntnUnderStaking.bondingNTN`
        // doing it in seperate loop may reduce gas cost by reducing state reading
        uint256 _bondingNTN;
        for (uint256 i = _topIndex; i < _length; i++) {
            _bondingNTN += bondingQueue[i].amount;
        }
        ntnUnderStaking.bondingNTN = _bondingNTN;
        bondingQueueTopIndex = _topIndex;
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
                _deleteUnbondingRequest(_unbondingRequest);
                _topIndex++;
            }
            else {
                break;
            }
        }

        // now calculate total `ntnUnderStaking.unbondingNTN`
        // doing it in seperate loop may reduce gas cost by reducing state reading
        uint256 _unbondingNTN;
        uint256 _unbondingShare;
        Autonity.Validator memory _validator;
        uint256 _currentEpochID = _getEpochID();
        for (uint256 i = _topIndex; i < _length; i++) {
            _unbondingRequest = unbondingQueue[i];
            if (_unbondingRequest.epochID == _currentEpochID) {
                break;
            }
            _unbondingShare = autonity.getUnbondingShare(_unbondingRequest.unbondingID);
            if (_unbondingShare == 0) {
                continue;
            }
            _validator = autonity.getValidator(_unbondingRequest.validator);
            _unbondingNTN += (_unbondingShare * _validator.unbondingStake) / _validator.unbondingShares;
        }
        ntnUnderStaking.unbondingNTN = _unbondingNTN;
        unbondingQueueTopIndex = _topIndex;
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

    function _getEpochID() internal view returns (uint256) {
        return autonity.epochID();
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
        return _unlockedFunds();
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
        require(msg.sender == managerContract, "caller is not the manager");
        _;
    }

}
