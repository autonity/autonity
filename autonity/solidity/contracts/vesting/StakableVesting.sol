// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../interfaces/IStakeProxy.sol";
import "./LiquidRewardManager.sol";
import "./ContractBase.sol";

contract StakableVesting is IStakeProxy, ContractBase, LiquidRewardManager {
    // NTN can be here: LOCKED or UNLOCKED
    // LOCKED are tokens that can't be withdrawn yet, need to wait for the release contract
    // UNLOCKED are tokens that can be withdrawn
    uint256 public contractVersion = 1;
    uint256 private requiredGasBond = 50_000;
    uint256 private requiredGasUnbond = 50_000;

    /**
     * @notice stake reserved to create new contracts
     * each time a new contract is creasted, totalNominal is decreased
     * address(this) should have totalNominal amount of NTN availabe,
     * otherwise withdrawing or bonding from a contract is not possible
     */
    uint256 public totalNominal;


    struct PendingBondingRequest {
        uint256 amount;
        uint256 epochID;
        address validator;
        bool processed;
    }

    struct PendingUnbondingRequest {
        uint256 liquidAmount;
        uint256 epochID;
        address validator;
        bool rejected;
        bool applied;
    }

    mapping(uint256 => PendingBondingRequest) private pendingBondingRequest;
    mapping(uint256 => PendingUnbondingRequest) private pendingUnbondingRequest;

    /**
     * @dev bondingToContract is needed to handle notification from autonity when bonding is applied.
     * In case the it fails to notify the vesting contract, contractToBonding is needed to revert the failed requests.
     * All bonding requests are applied at epoch end, so we can process all of them (failed or successful) together.
     * See bond and _revertPendingBondingRequest for more clarity
     */

    mapping(uint256 => uint256) private bondingToContract;
    mapping(uint256 => uint256[]) private contractToBonding;

    /**
     * @dev unbondingToContract is needed to handle notification from autonity when unbonding is applied and released.
     * In case it fails to notify the vesting contract, contractToUnbonding is needed to revert the failed requests.
     * Not all requests are released together at epoch end, so we cannot process all the request together.
     * tailPendingUnbondingID and headPendingUnbondingID helps to keep track of contractToUnbonding.
     * See unbond and _revertPendingUnbondingRequest for more clarity
     */

    mapping(uint256 => uint256) private unbondingToContract;
    mapping(uint256 => mapping(uint256 => uint256)) private contractToUnbonding;
    mapping(uint256 => uint256) private tailPendingUnbondingID;
    mapping(uint256 => uint256) private headPendingUnbondingID;

    // AUT rewards entitled to some beneficiary for bonding from some contract before it has been cancelled
    // see cancelContract for more clarity.
    mapping(address => uint256) private atnRewards;
    mapping(address => uint256) private ntnRewards;

    constructor(
        address payable _autonity, address _operator
    ) LiquidRewardManager(_autonity) ContractBase(_autonity, _operator) {}

    /**
     * @notice creates a new stakable contract, restricted to operator only
     * @dev _amount NTN must be minted and transferred to this contract before creating the contract
     * otherwise the contract cannot be released or bonded to some validator
     * @param _beneficiary address of the beneficiary
     * @param _amount total amount of NTN to be vested
     * @param _startTime start time of the vesting
     * @param _cliffDuration cliff period
     * @param _totalDuration total duration of the contract
     */
    function newContract(
        address _beneficiary,
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration
    ) virtual onlyOperator public {
        require(_startTime + _cliffDuration > autonity.lastEpochTime(), "contract cliff duration is past");
        require(totalNominal >= _amount, "not enough stake reserved to create a new contract");

        _createContract(_beneficiary, _amount, _startTime, _cliffDuration, _totalDuration, true);
        totalNominal -= _amount;
    }


    /**
     * @notice used by beneficiary to transfer all unlocked NTN and LNTN of some contract to his own address
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones)
     * So any beneficiary can number their contracts from 0 to (n-1). Beneficiary does not need to know the 
     * unique global contract id which can be retrieved via _getUniqueContractID function
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
     * @notice used by beneficiary to transfer all unlocked NTN of some contract to his own address
     */
    function releaseAllNTN(uint256 _id) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _cleanup(_contractID);
        _releaseNTN(_contractID, _unlockedFunds(_contractID));
    }

    /**
     * @notice used by beneficiary to transfer all unlocked LNTN of some contract to his own address
     */
    function releaseAllLNTN(uint256 _id) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _updateFunds(_contractID);
        _releaseAllUnlockedLNTN(_contractID, _unlockedFunds(_contractID));
        _clearValidators(_contractID);
    }

    // do we want this method to allow beneficiary withdraw a fraction of the released amount???
    /**
     * @notice used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address
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
     * @notice used by beneficiary to transfer some amount of unlocked LNTN of some contract to his own address
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
     * @notice changes the beneficiary of some contract to the _recipient address. _recipient can release and stake tokens from the contract.
     * only operator is able to call this function.
     * rewards which have been entitled to the beneficiary due to bonding from this contract are not transferred to _recipient
     * @dev rewards earned until this point from this contract are calculated and stored in atnRewards and ntnRewards mapping so that
     * _beneficiary can later claim them even though _beneficiary is not entitled to this contract.
     * @param _beneficiary beneficiary address whose contract will be canceled
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding already canceled ones)
     * @param _recipient whome the contract is transferred to
     */
    function changeContractBeneficiary(
        address _beneficiary, uint256 _id, address _recipient
    ) virtual external onlyOperator {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        (uint256 _atnReward, uint256 _ntnReward) = _claimRewards(_contractID);
        atnRewards[_beneficiary] += _atnReward;
        ntnRewards[_beneficiary] += _ntnReward;
        _changeContractBeneficiary(_contractID, _beneficiary, _recipient);
    }

    /**
     * @notice In case some funds are missing due to some pending staking operation that failed,
     * this function updates the funds of some contract _id entitled to _beneficiary by reverting the failed requests
     */
    function updateFunds(address _beneficiary, uint256 _id) virtual external {
        _updateFunds(_getUniqueContractID(_beneficiary, _id));
    }

    /**
     * @notice Set the value of totalNominal. Restricted to operator account
     * In case totalNominal is increased, the increased amount should be minted
     * and transferred to the address of this contract, otherwise newly created vesting
     * contracts will not have funds to withdraw or bond. See newContract()
     */
    function setTotalNominal(uint256 _newTotalNominal) virtual external onlyOperator {
        totalNominal = _newTotalNominal;
    }

    /**
     * @notice Update the required gas to get notified about staking operation
     * NOTE: before updating, please check if the updated value works. It can be checked by updatting
     * the hardcoded value of requiredGasBond and then compiling the contracts and running the tests
     * in stakable_vesting_test.go
     */
    function setRequiredGasBond(uint256 _gas) external onlyOperator {
        requiredGasBond = _gas;
    }

    function setRequiredGasUnbond(uint256 _gas) external onlyOperator {
        requiredGasUnbond = _gas;
    }

    /**
     * @notice Used by beneficiary to bond some NTN of some contract _id.
     * All bondings are delegated, as vesting manager cannot own a validator
     * @param _id id of the contract numbered from 0 to (n-1) where n = total contracts entitled to the beneficiary (excluding canceled ones)
     * @param _validator address of the validator for bonding
     * @param _amount amount of NTN to bond
     */
    function bond(uint256 _id, address _validator, uint256 _amount) virtual public payable returns (uint256) {
        // TODO (tariq): do we need to wait till _contract.start before bonding??
        require(msg.value >= requiredBondingGasCost(), "not enough gas given for notification on bonding");
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _updateFunds(_contractID);
        Contract storage _contract = contracts[_contractID];
        require(_contract.start <= block.timestamp, "contract not started yet");
        require(_contract.currentNTNAmount >= _amount, "not enough tokens");

        uint256 _bondingID = autonity.bond{value: msg.value}(_validator, _amount);
        _contract.currentNTNAmount -= _amount;
        // offset by 1 to handle empty value
        contractToBonding[_contractID].push(_bondingID);
        bondingToContract[_bondingID] = _contractID+1;
        pendingBondingRequest[_bondingID] = PendingBondingRequest(_amount, _epochID(), _validator, false);
        _initiate(_contractID, _validator);
        _clearValidators(_contractID);
        return _bondingID;
    }

    /**
     * @notice Used by beneficiary to unbond some LNTN of some contract.
     * @param _validator address of the validator
     * @param _amount amount of LNTN to unbond
     */
    function unbond(uint256 _id, address _validator, uint256 _amount) virtual public payable returns (uint256) {
        require(msg.value >= requiredUnbondingGasCost(), "not enough gas given for notification on unbonding");
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        _cleanup(_contractID);
        _lock(_contractID, _validator, _amount);
        uint256 _unbondingID = autonity.unbond{value: msg.value}(_validator, _amount);
        // offset by 1 to handle empty value
        unbondingToContract[_unbondingID] = _contractID+1;
        pendingUnbondingRequest[_unbondingID] = PendingUnbondingRequest(_amount, _epochID(), _validator, false, false);
        uint256 _lastID = headPendingUnbondingID[_contractID];
        // contractToUnbonding[_contractID][_i] stores the _unbondingID of the i'th unbonding request
        contractToUnbonding[_contractID][_lastID] = _unbondingID;
        headPendingUnbondingID[_contractID] = _lastID+1;
        return _unbondingID;
    }

    /**
     * @notice used by beneficiary to claim rewards from contract _id from bonding to _validator
     * @param _id contract ID
     * @param _validator validator address
     */
    function claimRewards(uint256 _id, address _validator) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        (uint256 _atnReward, uint256 _ntnReward) = _claimRewards(_contractID, _validator);
        _sendRewards(_atnReward, _ntnReward);
    }

    /**
     * @notice used by beneficiary to claim rewards from contract _id from bonding
     */
    function claimRewards(uint256 _id) virtual external {
        uint256 _contractID = _getUniqueContractID(msg.sender, _id);
        (uint256 _atnReward, uint256 _ntnReward) = _claimRewards(_contractID);
        _sendRewards(_atnReward, _ntnReward);
    }

    /**
     * @notice used by beneficiary to claim all rewards which is entitled from bonding
     * @dev Rewards from some cancelled contracts are stored in atnRewards and ntnRewards mapping. All rewards from
     * contracts that are still entitled to the beneficiary need to be calculated via _claimRewards
     */
    function claimRewards() virtual external {
        uint256[] storage _contractIDs = beneficiaryContracts[msg.sender];
        uint256 _atnTotalFees = atnRewards[msg.sender];
        uint256 _ntnTotalFees = ntnRewards[msg.sender];
        atnRewards[msg.sender] = 0;
        ntnRewards[msg.sender] = 0;
        
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            _updateFunds(_contractIDs[i]);
            (uint256 _atnReward, uint256 _ntnReward) = _claimRewards(_contractIDs[i]);
            _atnTotalFees += _atnReward;
            _ntnTotalFees += _ntnReward;
            _clearValidators(_contractIDs[i]);
        }
        _sendRewards(_atnTotalFees, _ntnTotalFees);
    }

    /**
     * @notice can be used to send ATN to the contract
     */
    function receiveATN() external payable {
        // do nothing
    }

    /**
     * @notice callback function restricted to autonity to notify the vesting contract to update rewards (AUT) for _validators
     * @param _validators address of the validators that have staking operations and need to update their rewards (AUT)
     */
    function rewardsDistributed(address[] memory _validators) external onlyAutonity {
        _updateValidatorReward(_validators);
    }

    /**
     * @notice implements IStakeProxy.bondingApplied(), a callback function for autonity when bonding is applied
     * @param _bondingID bonding id from Autonity when bonding was requested
     * @param _liquid amount of LNTN after bonding applied successfully
     * @param _selfDelegation true if self bonded, false for delegated bonding
     * @param _rejected true if bonding request was rejected, false if applied successfully
     */
    function bondingApplied(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) external onlyAutonity {
        _applyBonding(_bondingID, _validator, _liquid, _selfDelegation, _rejected);
    }

    /**
     * @notice implements IStakeProxy.unbondingApplied(). callback function for autonity when unbonding is applied
     * @param _unbondingID unbonding id from Autonity when unbonding was requested
     * @param _rejected true if unbonding was rejected, false if applied successfully
     */
    function unbondingApplied(uint256 _unbondingID, address _validator, bool _rejected) external onlyAutonity {
        _applyUnbonding(_unbondingID, _validator, _rejected);
    }

    /**
     * @dev implements IStakeProxy.unbondingReleased(). callback function for autonity when unbonding is released
     * @param _unbondingID unbonding id from Autonity when unbonding was requested
     * @param _amount amount of NTN released
     * @param _rejected true if unbonding was rejected, false if applied and released successfully
     */
    function unbondingReleased(uint256 _unbondingID, uint256 _amount, bool _rejected) external onlyAutonity {
        _releaseUnbonding(_unbondingID, _amount, _rejected);
    }

    /**
     * @dev returns equivalent amount of NTN using the ratio.
     * @param _validator validator address
     * @param _amount amount of LNTN to be converted
     */
    function _calculateLNTNValue(address _validator, uint256 _amount) private view returns (uint256) {
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake) / _validatorInfo.liquidSupply;
    }

    /**
     * @dev returns equivalent amount of LNTN using the ratio.
     * @param _validator validator address
     * @param _amount amount of NTN to be converted
     */
    function _getLiquidFromNTN(address _validator, uint256 _amount) private view returns (uint256) {
        Autonity.Validator memory _validatorInfo = autonity.getValidator(_validator);
        return _amount * _validatorInfo.liquidSupply / (_validatorInfo.bondedStake - _validatorInfo.selfBondedStake);
    }

    /**
     * @dev calculates the total value of the contract, which can vary if the contract has some LNTN
     * total value = current_NTN + withdrawn_value + (the value of LNTN converted to NTN)
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
     * @dev transfers some LNTN equivalent to _availableUnlockedFunds NTN to beneficiary address.
     * In case the _contractID has LNTN to multiple validators, we pick one validator and try to transfer
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
     * @dev calculates the amount of unlocked funds in NTN until last epoch block time
     */
    function _unlockedFunds(uint256 _contractID) private view returns (uint256) {
        return _calculateAvailableUnlockedFunds(
            _contractID, _calculateTotalValue(_contractID), autonity.lastEpochTime()
        );
    }

    function _updateAndTransferLNTN(uint256 _contractID, address _to, uint256 _amount, address _validator) private {
        _updateUnclaimedReward(_validator);
        _decreaseLiquid(_contractID, _validator, _amount);
        _transferLNTN(_to, _amount, _validator);
    }

    function _transferLNTN(address _to, uint256 _amount, address _validator) private {
        bool _sent = _liquidContract(_validator).transfer(_to, _amount);
        require(_sent, "LNTN transfer failed");
    }

    function _updateValidatorReward(address[] memory _validators) internal {
        for (uint256 i = 0; i < _validators.length; i++) {
            _updateUnclaimedReward(_validators[i]);
        }
    }

    /**
     * @dev mimic _applyBonding from Autonity.sol
     */
    function _applyBonding(uint256 _bondingID, address _validator, uint256 _liquid, bool _selfDelegation, bool _rejected) internal {
        require(_selfDelegation == false, "bonding should be delegated");
        uint256 _contractID = bondingToContract[_bondingID] - 1;

        PendingBondingRequest storage _request = pendingBondingRequest[_bondingID];
        _request.processed = true;

        if (_rejected) {
            Contract storage _contract = contracts[_contractID];
            _contract.currentNTNAmount += _request.amount;
        }
        else {
            _increaseLiquid(_contractID, _validator, _liquid);
        }
    }

    /**
     * @dev mimic _applyUnbonding from Autonity.sol
     */
    function _applyUnbonding(uint256 _unbondingID, address _validator, bool _rejected) internal {
        uint256 _contractID = unbondingToContract[_unbondingID] - 1;
        PendingUnbondingRequest storage _unbondingRequest = pendingUnbondingRequest[_unbondingID];
        uint256 _liquid = _unbondingRequest.liquidAmount;
        _unlock(_contractID, _validator, _liquid);

        if (_rejected) {
            _unbondingRequest.rejected = true;
            return;
        }

        _unbondingRequest.applied = true;
        _decreaseLiquid(_contractID, _validator, _liquid);
    }

    function _releaseUnbonding(uint256 _unbondingID, uint256 _amount, bool _rejected) internal {
        uint256 _contractID = unbondingToContract[_unbondingID] - 1;

        if (_rejected) {
            // If released of unbonding is rejected, then it is assumed that the applying of unbonding was also rejected
            // or reverted at Autonity because Autonity could not notify us (_applyUnbonding reverted).
            // If applying of unbonding was successful, then releasing of unbonding cannot be rejected.
            // Here we assume that it was rejected at vesting manager as well, otherwise if it was reverted (due to out of gas)
            // it will revert here as well. In any case, _revertPendingUnbondingRequest will handle the reverted or rejected request
            return;
        }
        Contract storage _contract = contracts[_contractID];
        _contract.currentNTNAmount += _amount;
    }

    function _sendRewards(uint256 _atnReward, uint256 _ntnReward) private {
        // Send the AUT
        // solhint-disable-next-line avoid-low-level-calls
        (bool _sent, ) = msg.sender.call{value: _atnReward}("");
        require(_sent, "Failed to send AUT");

        _sent = autonity.transfer(msg.sender, _ntnReward);
        require(_sent, "Failed to send NTN");
    }

    function _updateFunds(uint256 _contractID) private {
        _revertPendingBondingRequest(_contractID);
        _revertPendingUnbondingRequest(_contractID);
    }

    /**
     * @dev The contract needs to be cleaned before bonding, unbonding or claiming rewards.
     * _cleanup removes any unnecessary validator from the list, removes pending bonding or unbonding requests
     * that have been rejected or reverted but vesting manager could not be notified. If the clean up is not
     * done, then liquid balance could be incorrect due to not handling the bonding or unbonding request.
     * @param _contractID unique global contract id
     */
    function _cleanup(uint256 _contractID) private {
        _revertPendingBondingRequest(_contractID);
        _revertPendingUnbondingRequest(_contractID);
        _clearValidators(_contractID);
    }

    /**
     * @dev in case some bonding request from some previous epoch was unsuccessful and vesting contract was not notified,
     * this function handles such requests. All the requests from past epoch can be handled as the bonding requests are
     * applied at epoch end immediately. Requests from current epoch are not handled.
     * @param _contractID unique global id of the contract
     */
    function _revertPendingBondingRequest(uint256 _contractID) private {
        uint256[] storage _bondingIDs = contractToBonding[_contractID];
        uint256 _length = _bondingIDs.length;
        if (_length == 0) {
            return;
        }

        uint256 _oldBondingID;
        PendingBondingRequest storage _oldBondingRequest;
        uint256 _totalAmount = 0;
        uint256 _currentEpochID = _epochID();
        for (uint256 i = 0; i < _length; i++) {
            _oldBondingID = _bondingIDs[i];
            _oldBondingRequest = pendingBondingRequest[_oldBondingID];
            // will revert request from some previous epoch, request from current epoch will not be reverted
            if (_oldBondingRequest.epochID == _currentEpochID) {
                // all the request in the array are from current epoch
                return;
            }
            _bondingRequestExpired(_contractID, _oldBondingRequest.validator);
            // if the request is not processed successfully, then we need to revert it
            if (_oldBondingRequest.processed == false) {
                _totalAmount += _oldBondingRequest.amount;
            }
            delete pendingBondingRequest[_oldBondingID];
            delete bondingToContract[_oldBondingID];
        }

        delete contractToBonding[_contractID];
        if (_totalAmount == 0) {
            return;
        }
        Contract storage _contract = contracts[_contractID];
        _contract.currentNTNAmount += _totalAmount;
    }

    /**
     * @dev in case some unbonding request from some previous epoch was unsuccessful (not applied successfully or
     * not released successfully) and vesting contract was not notified, this function handles such requests.
     * Any request that has been processed in _releaseUnbondingStake function can be handled here.
     * Other requests need to wait.
     * @param _contractID unique global id of the contract
     */
    function _revertPendingUnbondingRequest(uint256 _contractID) private {
        uint256 _unbondingID;
        PendingUnbondingRequest storage _unbondingRequest;
        mapping(uint256 => uint256) storage _unbondingIDs = contractToUnbonding[_contractID];
        uint256 _lastID = headPendingUnbondingID[_contractID];
        uint256 _processingID = tailPendingUnbondingID[_contractID];
        for (; _processingID < _lastID; _processingID++) {
            _unbondingID = _unbondingIDs[_processingID];
            Autonity.UnbondingReleaseState _releaseState = autonity.getUnbondingReleaseState(_unbondingID);
            if (_releaseState == Autonity.UnbondingReleaseState.notReleased) {
                // all the rest have not been released yet
                break;
            }
            delete _unbondingIDs[_processingID];
            delete unbondingToContract[_unbondingID];
            
            if (_releaseState == Autonity.UnbondingReleaseState.released) {
                // unbonding was released successfully
                delete pendingUnbondingRequest[_unbondingID];
                continue;
            }

            _unbondingRequest = pendingUnbondingRequest[_unbondingID];
            address _validator = _unbondingRequest.validator;

            if (_releaseState == Autonity.UnbondingReleaseState.reverted) {
                // it means the unbonding was released, but later reverted due to failing to notify us
                // that means unbonding was applied successfully
                require(_unbondingRequest.applied, "unbonding was released but not applied succesfully");
                _updateUnclaimedReward(_validator);
                _increaseLiquid(_contractID, _validator, autonity.getRevertingAmount(_unbondingID));
            }
            else if (_releaseState == Autonity.UnbondingReleaseState.rejected) {
                require(_unbondingRequest.applied == false, "unbonding was applied successfully but release was rejected");
                // _unbondingRequest.rejected = true means we already rejected it
                if (_unbondingRequest.rejected == false) {
                    _unlock(_contractID, _validator, _unbondingRequest.liquidAmount);
                }
            }
            else {
                require(false, "unknown UnbondingReleaseState, need to implement");
            }
            
            delete pendingUnbondingRequest[_unbondingID];
        }
        tailPendingUnbondingID[_contractID] = _processingID;
    }

    function _epochID() private view returns (uint256) {
        return autonity.epochID();
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice returns the gas cost required in "wei" to notify vesting manager when the bonding is applied
     */
    function requiredBondingGasCost() public view returns (uint256) {
        return requiredGasBond * autonity.stakingGasPrice();
    }

    /**
     * @notice returns the gas cost required in "wei" to notify vesting manager when the unbonding is applied and released
     */
    function requiredUnbondingGasCost() public view returns (uint256) {
        return requiredGasUnbond * autonity.stakingGasPrice();
    }

    /**
     * @notice returns unclaimed rewards from contract _id entitled to _beneficiary from bonding to _validator
     * @param _beneficiary beneficiary address
     * @param _id contract ID
     * @param _validator validator address
     */
    function unclaimedRewards(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256 _atnFee, uint256 _ntnFee) {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        (_atnFee, _ntnFee) = _unclaimedRewards(_contractID, _validator);
    }

    /**
     * @notice returns unclaimed rewards from contract _id entitled to _beneficiary from bonding
     */
    function unclaimedRewards(address _beneficiary, uint256 _id) virtual external view returns (uint256 _atnFee, uint256 _ntnFee) {
        uint256 _contractID = _getUniqueContractID(_beneficiary, _id);
        (_atnFee, _ntnFee) = _unclaimedRewards(_contractID);
    }

    /**
     * @notice returns the amount of all unclaimed rewards due to all the bonding from contracts entitled to _beneficiary
     */
    function unclaimedRewards(address _beneficiary) virtual external view returns (uint256 _atnTotalFee, uint256 _ntnTotalFee) {
        _atnTotalFee = atnRewards[_beneficiary];
        _ntnTotalFee = ntnRewards[_beneficiary];
        uint256[] storage _contractIDs = beneficiaryContracts[_beneficiary];
        for (uint256 i = 0; i < _contractIDs.length; i++) {
            (uint256 _atnFee, uint256 _ntnFee) = _unclaimedRewards(_contractIDs[i]);
            _atnTotalFee += _atnFee;
            _ntnTotalFee += _ntnFee;
        }
    }

    /**
     * @notice returns the amount of LNTN for some contract
     * @param _beneficiary beneficiary address
     * @param _id contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones)
     * @param _validator validator address
     */
    function liquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _liquidBalanceOf(_getUniqueContractID(_beneficiary, _id), _validator);
    }

    /**
     * @notice returns the amount of locked LNTN for some contract
     */
    function lockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _lockedLiquidBalanceOf(_getUniqueContractID(_beneficiary, _id), _validator);
    }

    /**
     * @notice returns the amount of unlocked LNTN for some contract
     */
    function unlockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) virtual external view returns (uint256) {
        return _unlockedLiquidBalanceOf(_getUniqueContractID(_beneficiary, _id), _validator);
    }

    /**
     * @notice returns the list of validators bonded some contract
     */
    function getBondedValidators(address _beneficiary, uint256 _id) virtual external view returns (address[] memory) {
        return _bondedValidators(_getUniqueContractID(_beneficiary, _id));
    }

    /**
     * @notice returns the amount of released funds in NTN for some contract
     */
    function unlockedFunds(address _beneficiary, uint256 _id) virtual external view returns (uint256) {
        return _unlockedFunds(_getUniqueContractID(_beneficiary, _id));
    }

    function contractTotalValue(address _beneficiary, uint256 _id) external view returns (uint256) {
        return _calculateTotalValue(_getUniqueContractID(_beneficiary, _id));
    }

}
