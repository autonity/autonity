// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "./LiquidRewardManager.sol";
import "./ContractBase.sol";

contract StakableVesting is ContractBase, LiquidRewardManager {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release contract
    // UNLOCKED are tokens that can be withdrawn
    uint256 public contractVersion = 1;

    /**
     * @notice Sum of total amount of contracts that can be created.
     * Each time a new contract is created, totalNominal is decreased.
     * Address(this) should have totalNominal amount of NTN availabe at genesis,
     * otherwise withdrawing or bonding from a contract is not possible
     */
    uint256 public totalNominal;


    struct PendingBondingRequest {
        uint256 amount;
        uint256 epochID;
        address validator;
    }

    struct PendingUnbondingRequest {
        uint256 liquidAmount;
        uint256 epochID;
        address validator;
    }

    mapping(uint256 => PendingBondingRequest) private pendingBondingRequest;
    mapping(uint256 => PendingUnbondingRequest) private pendingUnbondingRequest;

    /**
     * @dev We put all the bonding request id of past epoch in contractToBonding[contractID] array and apply them whenever needed.
     * All bonding requests are applied at epoch end, so we can process all of them (failed or successful) together.
     * See bond and _handlePendingBondingRequest for more clarity
     */

    mapping(uint256 => uint256[]) private contractToBonding;

    /**
     * @dev We put all the unbonding request id of past epoch in contractToUnbonding mapping. All requests from past epoch
     * can be applied together. But not all requests are released together at epoch end. So we need to put them in map
     * and use tailPendingUnbondingID and headPendingUnbondingID to keep track of contractToUnbonding.
     * See unbond and _handlePendingUnbondingRequest for more clarity
     */

    mapping(uint256 => mapping(uint256 => uint256)) private contractToUnbonding;
    mapping(uint256 => uint256) private appliedPendingUnbondingID;
    mapping(uint256 => uint256) private tailPendingUnbondingID;
    mapping(uint256 => uint256) private headPendingUnbondingID;

    /**
     * @dev ATN rewards entitled to some beneficiary for bonding from some contract before it has been cancelled.
     * See cancelContract for more clarity.
     */
    mapping(address => uint256) private atnRewards;

    /**
     * @dev Same as atnRewards for NTN rewards
     */
    mapping(address => uint256) private ntnRewards;

    constructor(
        address payable _autonity, address _operator
    ) LiquidRewardManager(_autonity) ContractBase(_autonity, _operator) {}

    /**
     * @notice Creates a new stakable contract.
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _startTime start time of the vesting
     * @param _cliffDuration cliff period
     * @param _totalDuration total duration of the contract
     * @custom:restricted-to operator account
     */
    function newContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration
    ) virtual onlyOperator public {
        require(_startTime + _cliffDuration >= autonity.lastEpochTime(), "contract cliff duration is past");
        require(totalNominal >= _amount, "not enough stake reserved to create a new contract");

        _createContract(_beneficiary, _amount, _startTime, _cliffDuration, _totalDuration, true);
        totalNominal -= _amount;
    }


    /**
     * @notice Used by beneficiary to transfer all unlocked NTN and LNTN of some contract to his own address.
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones).
     * So any beneficiary can number their contracts from 0 to (n-1). Beneficiary does not need to know the unique global contract id.
     */
    function releaseFunds(uint256 _id) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _updateFunds(_contractID);
        uint256 _unlocked = _unlockedFunds(_contractID);
        // first NTN is released
        _unlocked = _releaseNTN(_contractID, _unlocked);
        // if there still remains some unlocked funds, i.e. not enough NTN, then LNTN is released
        _releaseAllUnlockedLNTN(_contractID, _unlocked);
        _clearValidators(_contractID);
    }

    /**
     * @notice Used by beneficiary to transfer all unlocked NTN of some contract to his own address.
     */
    function releaseAllNTN(uint256 _id) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _cleanup(_contractID);
        _releaseNTN(_contractID, _unlockedFunds(_contractID));
    }

    /**
     * @notice Used by beneficiary to transfer all unlocked LNTN of some contract to his own address.
     */
    function releaseAllLNTN(uint256 _id) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _updateFunds(_contractID);
        _releaseAllUnlockedLNTN(_contractID, _unlockedFunds(_contractID));
        _clearValidators(_contractID);
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address.
     * @param _amount amount to transfer
     */
    function releaseNTN(uint256 _id, uint256 _amount) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _cleanup(_contractID);
        require(_amount <= _unlockedFunds(_contractID), "not enough unlocked funds");
        _releaseNTN(_contractID, _amount);
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice Used by beneficiary to transfer some amount of unlocked LNTN of some contract to his own address.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to transfer
     */
    function releaseLNTN(uint256 _id, address _validator, uint256 _amount) virtual external {
        require(_amount > 0, "require positive amount to transfer");
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _updateFunds(_contractID);

        uint256 _unlockedLiquid = _unlockedLiquidBalanceOf(_contractID, _validator);
        require(_unlockedLiquid >= _amount, "not enough unlocked LNTN");

        uint256 _value = _calculateLNTNValue(_validator, _amount);
        require(_value <= _unlockedFunds(_contractID), "not enough unlocked funds");

        Contract storage _contract = contracts[_contractID];
        _contract.withdrawnValue += _value;
        _updateAndTransferLNTN(_contractID, msg.sender, _amount, _validator);
        _clearValidators(_contractID);
    }

    /**
     * @notice Changes the beneficiary of some contract to the recipient address. The recipient address can release and stake tokens from the contract.
     * Rewards which have been entitled to the beneficiary due to bonding from this contract are not transferred to recipient.
     * @dev Rewards earned until this point from this contract are calculated and stored in atnRewards and ntnRewards mapping so that
     * beneficiary can later claim them even though beneficiary is not entitled to this contract.
     * @param _beneficiary beneficiary address whose contract will be canceled
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding already canceled ones)
     * @param _recipient whome the contract is transferred to
     * @custom:restricted-to operator account
     */
    function changeContractBeneficiary(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        _updateFunds(_contractID);
        (uint256 _atnReward, uint256 _ntnReward) = _claimRewards(_contractID);
        atnRewards[_beneficiary] += _atnReward;
        ntnRewards[_beneficiary] += _ntnReward;
        _clearValidators(_contractID);
        _changeContractBeneficiary(_contractID, _beneficiary, _recipient);
    }

    /**
     * @notice In case some funds are missing due to some pending staking operation that failed,
     * this function updates the funds of some contract entitled to beneficiary by applying the pending requests.
     * @param _beneficiary beneficiary address
     * @param _id contract id to update
     */
    function updateFunds(address _beneficiary, uint256 _id) virtual external {
        _updateFunds(_getUniqueContractID(_beneficiary, _id));
    }

    /**
     * @notice Updates the funds of the contract and returns total value of the contract
     */
    function updateFundsAndGetContractTotalValue(address _beneficiary, uint256 _id) external returns (uint256) {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        _updateFunds(_contractID);
        return _calculateTotalValue(_contractID);
    }

    /**
     * @notice Updates the funds of the contract and returns the contract
     */
    function updateFundsAndGetContract(address _beneficiary, uint256 _id) external returns (Contract memory) {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        _updateFunds(_contractID);
        return contracts[_contractID];
    }

    /**
     * @notice Set the value of totalNominal
     * In case totalNominal is increased, the increased amount should be minted
     * and transferred to the address of this contract, otherwise newly created vesting
     * contracts will not have funds to withdraw or bond. See newContract()
     * @custom:restricted-to operator account
     */
    function setTotalNominal(uint256 _newTotalNominal) virtual external onlyOperator {
        totalNominal = _newTotalNominal;
    }

    /**
     * @notice Used by beneficiary to bond some NTN of some contract.
     * All bondings are delegated, as vesting manager cannot own a validator.
     * @param _id id of the contract numbered from 0 to (n-1) where n = total contracts entitled to the beneficiary (excluding canceled ones)
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(uint256 _id, address _validator, uint256 _amount) virtual public returns (uint256) {

        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _updateFunds(_contractID);

        Contract storage _contract = contracts[_contractID];
        require(_contract.start <= block.timestamp, "contract not started yet");
        require(_contract.currentNTNAmount >= _amount, "not enough tokens");

        uint256 _bondingID = autonity.bond(_validator, _amount);
        _contract.currentNTNAmount -= _amount;

        contractToBonding[_contractID].push(_bondingID);
        pendingBondingRequest[_bondingID] = PendingBondingRequest(_amount, _getEpochID(), _validator);

        _newBondingRequested(_contractID, _validator, _bondingID);
        _clearValidators(_contractID);
        return _bondingID;
    }

    /**
     * @notice Used by beneficiary to unbond some LNTN of some contract.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public returns (uint256) {

        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _cleanup(_contractID);

        _lock(_contractID, _validator, _amount);
        uint256 _unbondingID = autonity.unbond(_validator, _amount);
        
        pendingUnbondingRequest[_unbondingID] = PendingUnbondingRequest(_amount, _getEpochID(), _validator);
        uint256 _lastID = headPendingUnbondingID[_contractID];
        // contractToUnbonding[_contractID][_i] stores the _unbondingID of the i'th unbonding request
        contractToUnbonding[_contractID][_lastID] = _unbondingID;
        headPendingUnbondingID[_contractID] = _lastID+1;

        _newPendingRewardEvent(_validator, _unbondingID, false);
        return _unbondingID;
    }

    /**
     * @notice Used by beneficiary to claim rewards from bonding some contract to validator.
     * @param _id contract ID
     * @param _validator validator address
     */
    function claimRewards(uint256 _id, address _validator) virtual external returns (uint256 _atnRewards, uint256 _ntnRewards) {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _updateFunds(_contractID);
        (_atnRewards, _ntnRewards) = _claimRewards(_contractID, _validator);
        _sendRewards(_atnRewards, _ntnRewards);
        _clearValidators(_contractID);
    }

    /**
     * @notice Used by beneficiary to claim rewards from bonding some contract to validator.
     */
    function claimRewards(uint256 _id) virtual external returns (uint256 _atnRewards, uint256 _ntnRewards) {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _updateFunds(_contractID);
        (_atnRewards, _ntnRewards) = _claimRewards(_contractID);
        _sendRewards(_atnRewards, _ntnRewards);
        _clearValidators(_contractID);
    }

    /**
     * @notice Used by beneficiary to claim all rewards which is entitled from bonding
     * @dev Rewards from some cancelled contracts are stored in atnRewards and ntnRewards mapping. All rewards from
     * contracts that are still entitled to the beneficiary need to be calculated.
     */
    function claimRewards() virtual external returns (uint256 _atnRewards, uint256 _ntnRewards) {
        uint256[] storage _contractIDs = beneficiaryContracts[msg.sender];
        _atnRewards = atnRewards[msg.sender];
        _ntnRewards = ntnRewards[msg.sender];
        atnRewards[msg.sender] = 0;
        ntnRewards[msg.sender] = 0;
        
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            _updateFunds(_contractIDs[i]);
            (uint256 _atn, uint256 _ntn) = _claimRewards(_contractIDs[i]);
            _atnRewards += _atn;
            _ntnRewards += _ntn;
            _clearValidators(_contractIDs[i]);
        }
        _sendRewards(_atnRewards, _ntnRewards);
    }

    /**
     * @dev Receive Auton function https://solidity.readthedocs.io/en/v0.7.2/contracts.html#receive-ether-function
     */
    receive() external payable {}

    /**
     * @dev Returns equivalent amount of NTN using the current ratio.
     * @param _validator validator address
     * @param _amount amount of LNTN to be converted
     */
    function _calculateLNTNValue(address _validator, uint256 _amount) private view returns (uint256) {
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake) / _validatorInfo.liquidSupply;
    }

    /**
     * @dev Returns equivalent amount of LNTN using the current ratio.
     * @param _validator validator address
     * @param _amount amount of NTN to be converted
     */
    function _getLiquidFromNTN(address _validator, uint256 _amount) private view returns (uint256) {
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * _validatorInfo.liquidSupply / (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake);
    }

    /**
     * @dev Calculates the total value of the contract, which can vary if the contract has some LNTN.
     * Total value = currentNTN + withdrawnValue + (the value of LNTN converted to NTN using current ratio)
     * @param _contractID unique global id of the contract
     */
    function _calculateTotalValue(uint256 _contractID) private view returns (uint256) {
        Contract storage _contract = contracts[_contractID];
        uint256 _totalValue = _contract.currentNTNAmount + _contract.withdrawnValue;
        address[] memory _validators = _bondedValidators(_contractID);
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _balance = _liquidBalanceOf(_contractID, _validators[i]);
            if (_balance == 0) {
                continue;
            }
            _totalValue += _calculateLNTNValue(_validators[i], _balance);
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
        uint256 _contractID, uint256 _availableUnlockedFunds
    ) private returns (uint256 _remaining) {
        _remaining = _availableUnlockedFunds;
        address[] memory _validators = _bondedValidators(_contractID);
        for (uint256 i = 0; i < _validators.length && _remaining > 0; i++) {
            uint256 _balance = _unlockedLiquidBalanceOf(_contractID, _validators[i]);
            if (_balance == 0) {
                continue;
            }
            uint256 _value = _calculateLNTNValue(_validators[i], _balance);
            if (_remaining >= _value) {
                _remaining -= _value;
                _updateAndTransferLNTN(_contractID, msg.sender, _balance, _validators[i]);
            }
            else {
                uint256 _liquid = _getLiquidFromNTN(_validators[i], _remaining);
                require(_liquid <= _balance, "conversion not working");
                _remaining = 0;
                _updateAndTransferLNTN(_contractID, msg.sender, _liquid, _validators[i]);
            }
        }
        Contract storage _contract = contracts[_contractID];
        _contract.withdrawnValue += _availableUnlockedFunds - _remaining;
    }

    /**
     * @dev Calculates the amount of unlocked funds in NTN until last epoch time.
     */
    function _unlockedFunds(uint256 _contractID) private view returns (uint256) {
        return _calculateAvailableUnlockedFunds(
            _contractID, _calculateTotalValue(_contractID), autonity.lastEpochTime()
        );
    }

    function _updateAndTransferLNTN(uint256 _contractID, address _to, uint256 _amount, address _validator) private {
        _burnLiquid(_contractID, _validator, _amount, _getEpochID()-1);
        bool _sent = _liquidContract(_validator).transfer(_to, _amount);
        require(_sent, "LNTN transfer failed");

        // this transfer decreases the total liquid balance, which will affect pending reward event
        // update pending reward event, if any
        _updatePendingEventLiquid(_validator);
    }

    function _sendRewards(uint256 _atnReward, uint256 _ntnReward) private {
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool _sent, ) = msg.sender.call{value: _atnReward}("");
        require(_sent, "failed to send ATN");

        _sent = autonity.transfer(msg.sender, _ntnReward);
        require(_sent, "failed to send NTN");
    }

    /**
     * @dev In case some staking is reverted or rejected because it was not possible to notify us, then the funds will be incorrect.
     * This function handles the rejected or reverted requests and updates the funds.
     */
    function _updateFunds(uint256 _contractID) private {
        _handlePendingBondingRequest(_contractID);
        _handlePendingUnbondingRequest(_contractID);
    }

    /**
     * @dev Updates the funds and removes any unnecessary validator from the list.
     * @param _contractID unique global contract id
     */
    function _cleanup(uint256 _contractID) private {
        _updateFunds(_contractID);
        _clearValidators(_contractID);
    }

    /**
     * @dev Handles all the pending bonding requests.
     * All the requests from past epoch can be handled as the bonding requests are
     * applied at epoch end immediately. Requests from current epoch are not handled.
     * @param _contractID unique global id of the contract
     */
    function _handlePendingBondingRequest(uint256 _contractID) private {
        uint256[] storage _bondingIDs = contractToBonding[_contractID];
        uint256 _length = _bondingIDs.length;
        if (_length == 0) {
            return;
        }

        uint256 _bondingID;
        PendingBondingRequest storage _bondingRequest;
        uint256 _totalBondingRejected;
        uint256 _currentEpochID = _getEpochID();
        for (uint256 i = 0; i < _length; i++) {
            _bondingID = _bondingIDs[i];
            _bondingRequest = pendingBondingRequest[_bondingID];
            // request is from current epoch, not applied yet
            if (_bondingRequest.epochID == _currentEpochID) {
                // all the request in the array are from current epoch
                return;
            }
            
            _bondingRequestExpired(_contractID, _bondingRequest.validator);
            if (autonity.isBondingRejected(_bondingID)) {
                _totalBondingRejected += _bondingRequest.amount;
            }
            else {
                uint256 _liquid = autonity.getBondedLiquid(_bondingID);
                _mintLiquid(_contractID, _bondingRequest.validator, _liquid, _bondingRequest.epochID);
            }
            delete pendingBondingRequest[_bondingID];
        }

        delete contractToBonding[_contractID];
        if (_totalBondingRejected == 0) {
            return;
        }
        Contract storage _contract = contracts[_contractID];
        _contract.currentNTNAmount += _totalBondingRejected;
    }

    /**
     * @dev Handles all the pending unbonding requests. All unbonding requests from past epoch are applied.
     * Unbonding request that are released in Autonity are released.
     * @param _contractID unique global id of the contract
     */
    function _handlePendingUnbondingRequest(uint256 _contractID) private {
        uint256 _unbondingID;
        PendingUnbondingRequest storage _unbondingRequest;
        mapping(uint256 => uint256) storage _unbondingIDs = contractToUnbonding[_contractID];
        uint256 _lastID = headPendingUnbondingID[_contractID];
        uint256 _currentEpochID = _getEpochID();

        // first apply all request
        uint256 _processingID = appliedPendingUnbondingID[_contractID];
        for (; _processingID < _lastID; _processingID++) {
            _unbondingID = _unbondingIDs[_processingID];
            _unbondingRequest = pendingUnbondingRequest[_unbondingID];
            if (_unbondingRequest.epochID == _currentEpochID) {
                break;
            }

            // unbonding request is always applied at Autonity
            _unlockAndBurnLiquid(_contractID, _unbondingRequest.validator, _unbondingRequest.liquidAmount, _unbondingRequest.epochID);
            delete pendingUnbondingRequest[_unbondingID];
        }
        appliedPendingUnbondingID[_contractID] = _processingID;

        // process released stake
        uint256 _totalReleasedStake;
        _lastID = appliedPendingUnbondingID[_contractID];
        _processingID = tailPendingUnbondingID[_contractID];
        for (; _processingID < _lastID; _processingID++) {
            _unbondingID = _unbondingIDs[_processingID];
            bool _released = autonity.isUnbondingReleased(_unbondingID);
            if (_released == false) {
                // all the rest have not been released yet
                break;
            }
            delete _unbondingIDs[_processingID];

            _totalReleasedStake += autonity.getReleasedStake(_unbondingID);
        }
        tailPendingUnbondingID[_contractID] = _processingID;
        if (_totalReleasedStake > 0) {
            Contract storage _contract = contracts[_contractID];
            _contract.currentNTNAmount += _totalReleasedStake;
        }
    }

    /**
     * @dev Calculates the balance changes for the pending staking requests to some validator.
     * If the requests are from current epoch, then they are not taken into account.
     * @param _contractID unique contract ID
     * @param _validator validator address
     * @return _balanceChange balance change, positive if balance increases, negative otherwise
     * @return _lastRequestEpoch epochID, offset by 1, of the requests (all requests should be from the same epoch)
     */
    function _balanceChangeFromStakingRequest(
        uint256 _contractID,
        address _validator
    ) private view returns (int256 _balanceChange, uint256 _lastRequestEpoch) {
        uint256 _currentEpochID = _getEpochID();

        // get balance increase for bonding requests
        uint256[] storage _bondingIDs = contractToBonding[_contractID];
        uint256 _length = _bondingIDs.length;
        uint256 _stakingID;
        PendingBondingRequest storage _bondingRequest;

        for (uint256 i = 0; i < _length; i++) {
            _stakingID = _bondingIDs[i];
            _bondingRequest = pendingBondingRequest[_stakingID];
            // request is from current epoch, not applied yet
            if (_bondingRequest.epochID == _currentEpochID) {
                // all the request in the array are from current epoch
                break;
            }

            if (_bondingRequest.validator != _validator) {
                continue;
            }
            
            if (autonity.isBondingRejected(_stakingID) == false) {
                _balanceChange += int256(autonity.getBondedLiquid(_stakingID));
            }
            require(
                _lastRequestEpoch == 0 || _lastRequestEpoch == _bondingRequest.epochID+1,
                "pending requests are from different epoch"
            );
            _lastRequestEpoch = _bondingRequest.epochID+1;
        }

        // get balance decrease for unbonding requests
        PendingUnbondingRequest storage _unbondingRequest;
        mapping(uint256 => uint256) storage _unbondingIDs = contractToUnbonding[_contractID];
        uint256 _lastID = headPendingUnbondingID[_contractID];
        uint256 _processingID = appliedPendingUnbondingID[_contractID];

        for (; _processingID < _lastID; _processingID++) {
            _stakingID = _unbondingIDs[_processingID];
            _unbondingRequest = pendingUnbondingRequest[_stakingID];
            if (_unbondingRequest.epochID == _currentEpochID) {
                break;
            }
            if (_unbondingRequest.validator != _validator) {
                continue;
            }

            // unbonding request is always applied at Autonity
            _balanceChange -= int256(_unbondingRequest.liquidAmount);
            require(
                _lastRequestEpoch == 0 || _lastRequestEpoch == _unbondingRequest.epochID+1,
                "pending requests are from different epoch"
            );
            _lastRequestEpoch = _unbondingRequest.epochID+1;
        }
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns unclaimed rewards from some contract entitled to beneficiary from bonding to validator.
     * @param _beneficiary beneficiary address
     * @param _id contract ID
     * @param _validator validator address
     * @return _atnReward unclaimed ATN rewards
     * @return _ntnReward unclaimed NTN rewards
     */
    function unclaimedRewards(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256 _atnRewards, uint256 _ntnRewards) {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        (int256 _balanceChange, uint256 _epochID) = _balanceChangeFromStakingRequest(_contractID, _validator);
        (_atnRewards, _ntnRewards) = _unclaimedRewards(_contractID, _validator, _balanceChange, _epochID);
    }

    /**
     * @notice Returns unclaimed rewards from some contract entitled to beneficiary from bonding.
     */
    function unclaimedRewards(address _beneficiary, uint256 _id) virtual public view returns (uint256 _atnRewards, uint256 _ntnRewards) {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        address[] memory _validators = _bondedValidators(_contractID);
        for (uint256 i = 0; i < _validators.length; i++) {
            (int256 _balanceChange, uint256 _epochID) = _balanceChangeFromStakingRequest(_contractID, _validators[i]);
            (uint256 _atn, uint256 _ntn) = _unclaimedRewards(_contractID, _validators[i], _balanceChange, _epochID);
            _atnRewards += _atn;
            _ntnRewards += _ntn;
        }
    }

    /**
     * @notice Returns the amount of all unclaimed rewards due to all the bonding from contracts entitled to beneficiary.
     */
    function unclaimedRewards(address _beneficiary) virtual external view returns (uint256 _atnRewards, uint256 _ntnRewards) {
        _atnRewards = atnRewards[_beneficiary];
        _ntnRewards = ntnRewards[_beneficiary];
        uint256 _length = beneficiaryContracts[_beneficiary].length;
        for (uint256 i = 0; i < _length; i++) {
            (uint256 _atnFee, uint256 _ntnFee) = unclaimedRewards(_beneficiary, i);
            _atnRewards += _atnFee;
            _ntnRewards += _ntnFee;
        }
    }

    /**
     * @notice Returns the amount of LNTN for some contract.
     * @param _beneficiary beneficiary address
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones)
     * @param _validator validator address
     */
    function liquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _liquidBalanceOf(_getUniqueContractID(_beneficiary, _id), _validator);
    }

    /**
     * @notice Returns the amount of locked LNTN for some contract.
     */
    function lockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _lockedLiquidBalanceOf(_getUniqueContractID(_beneficiary, _id), _validator);
    }

    /**
     * @notice Returns the amount of unlocked LNTN for some contract.
     */
    function unlockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _unlockedLiquidBalanceOf(_getUniqueContractID(_beneficiary, _id), _validator);
    }

    /**
     * @notice Returns the list of validators bonded some contract.
     */
    function getBondedValidators(address _beneficiary, uint256 _id) virtual external view returns (address[] memory) {
        return _bondedValidators(_getUniqueContractID(_beneficiary, _id));
    }

    /**
     * @notice Returns the amount of released funds in NTN for some contract.
     */
    function unlockedFunds(address _beneficiary, uint256 _id) virtual external view returns (uint256) {
        return _unlockedFunds(_getUniqueContractID(_beneficiary, _id));
    }

    /**
     * @notice Returns the current total value of the contract in NTN.
     */
    function contractTotalValue(address _beneficiary, uint256 _id) external view returns (uint256) {
        return _calculateTotalValue(_getUniqueContractID(_beneficiary, _id));
    }

}
