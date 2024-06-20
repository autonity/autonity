// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../Autonity.sol";
import "../LiquidState.sol";

contract LiquidRewardManager {

    uint256 public constant FEE_FACTOR_UNIT_RECIP = 1_000_000_000;

    address private operator;
    Autonity private autonity;

    // when initiated, all values are offset by 1
    uint256 private constant OFFSET = 1;

    struct LiquidInfo {
        uint256 atnLastUnrealisedFeeFactor;
        uint256 atnUnclaimedRewards;
        uint256 ntnLastUnrealisedFeeFactor;
        uint256 ntnUnclaimedRewards;
        LiquidState liquidStateContract;
    }

    struct Account {
        uint256 liquidBalance;
        uint256 lockedLiquidBalance;
        uint256 atnRealisedFee;
        uint256 atnUnrealisedFeeFactor;
        uint256 ntnRealisedFee;
        uint256 ntnUnrealisedFeeFactor;
        bool initiated;
        bool newBondingRequested;
    }

    // stores total liquid and lastUnrealisedFeeFactor for each validator
    // lastUnrealisedFeeFactor is used to calculate unrealised rewards for schedules with the same logic as done in LiquidLogic.sol
    mapping(address => LiquidInfo) private liquidInfo;

    // stores the array of validators bonded to a schedule
    mapping(uint256 => address[]) private bondedValidators;
    // validatorIdx[_id][_validator] stores the (index+1) of _validator in bondedValidators[_id] array
    mapping(uint256 => mapping(address => uint256)) private validatorIdx;

    mapping(uint256 => mapping(address => Account)) private accounts;

    constructor(address payable _autonity) {
        autonity = Autonity(_autonity);
    }

    function _bondingRequestExpired(uint256 _id, address _validator) internal {
        accounts[_id][_validator].newBondingRequested = false;
    }

    /**
     * @dev initiate a pair (_id, _validator) only once for lifetime
     * initiating the pair helps to reduces maximum allowed gas usage for notification of staking operations
     * @param _id schedule id
     * @param _validator validator address
     */
    function _initiate(uint256 _id, address _validator) internal {
        _addValidator(_id, _validator);
        
        Account storage _account = accounts[_id][_validator];
        _account.newBondingRequested = true;
        if (_account.initiated) {
            return;
        }
        accounts[_id][_validator] = Account(OFFSET, OFFSET, OFFSET, OFFSET, OFFSET, OFFSET, true, true);
    }

    function _unlock(uint256 _id, address _validator, uint256 _amount) internal {
        Account storage _account = accounts[_id][_validator];
        _account.lockedLiquidBalance -= _amount;
    }

    function _lock(uint256 _id, address _validator, uint256 _amount) internal {
        Account storage _account = accounts[_id][_validator];
        require(_account.liquidBalance - _account.lockedLiquidBalance >= _amount, "not enough unlocked liquid tokens");
        _account.lockedLiquidBalance += _amount;
    }

    /**
     * @dev _decreaseLiquid, _increaseLiquid, _realiseFees follow the same logic as done in LiquidLogic.sol
     * the only difference is we update unclaimed rewards from the validator first, because this contract does not know
     * when the epoch ends, and so cannot claim rewards at each epoch end.
     * _updateUnclaimedReward or _claimRewards must be called before calling _decreaseLiquid, _increaseLiquid functions.
     * In the process of notification from autonity for staking operation, at first vesting contract is notified about the validators
     * in rewardsDistributed function, where we call _updateUnclaimedReward. Then vesting contract is notified about staking operations
     * in bondingApplied and unbondingApplied where we call _decreaseLiquid and _increaseLiquid
     */
    function _decreaseLiquid(uint256 _id, address _validator, uint256 _amount) internal {
        _realiseFees(_id, _validator);
        Account storage _account = accounts[_id][_validator];
        _account.liquidBalance -= _amount;
    }

    function _increaseLiquid(uint256 _id, address _validator, uint256 _amount) internal {
        _realiseFees(_id, _validator);
        Account storage _account = accounts[_id][_validator];
        _account.liquidBalance += _amount;
    }

    /**
     * @dev _updateUnclaimedReward or _claimRewards must be called before this function
     */
    function _realiseFees(uint256 _id, address _validator) private returns (uint256 _atnRealisedFees, uint256 _ntnRealisedFees) {
        uint256 _atnLastUnrealisedFeeFactor = liquidInfo[_validator].atnLastUnrealisedFeeFactor;
        Account storage _account = accounts[_id][_validator];
        uint256 _balance = _account.liquidBalance-OFFSET;
        _atnRealisedFees = _account.atnRealisedFee
                        + _computeUnrealisedFees(
                            _balance, _account.atnUnrealisedFeeFactor, _atnLastUnrealisedFeeFactor
                        );
        _account.atnRealisedFee = _atnRealisedFees;
        _account.atnUnrealisedFeeFactor = _atnLastUnrealisedFeeFactor;
        // remove the offset
        _atnRealisedFees -= OFFSET;

        uint256 _ntnLastUnrealisedFeeFactor = liquidInfo[_validator].ntnLastUnrealisedFeeFactor;
        _ntnRealisedFees = _account.ntnRealisedFee
                        + _computeUnrealisedFees(
                            _balance, _account.ntnUnrealisedFeeFactor, _ntnLastUnrealisedFeeFactor
                        );
        _account.ntnRealisedFee = _ntnRealisedFees;
        _account.ntnUnrealisedFeeFactor = _ntnLastUnrealisedFeeFactor;
        // remove the offset
        _ntnRealisedFees -= OFFSET;
    }

    function _computeUnrealisedFees(
        uint256 _balance, uint256 _unrealisedFeeFactor, uint256 _lastUnrealisedFeeFactor
    ) private pure returns (uint256) {
        if (_balance == 0) {
            return 0;
        }
        return (_lastUnrealisedFeeFactor - _unrealisedFeeFactor) * _balance / FEE_FACTOR_UNIT_RECIP;
    }

    /**
     * @dev claims all rewards from the liquid contract of the _validator
     * @param _validator validator address
     */
    function _claimRewards(address _validator) private {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        _updateUnclaimedReward(_validator);
        // because all the rewards are being claimed
        // offset by 1
        _liquidInfo.atnUnclaimedRewards = OFFSET;
        _liquidInfo.ntnUnclaimedRewards = OFFSET;
        (bool _success, ) = address(_liquidInfo.liquidStateContract).call(
            abi.encodeWithSignature("claimRewards()")
        );
        require(_success, "claimRewards() from liquid state contract reverted");
    }

    /**
     * @dev calculates total rewards for a schedule and resets realisedFees[id][validator] as reward is claimed
     */ 
    function _claimRewards(uint256 _id) internal returns (uint256 _atnTotalFees, uint256 _ntnTotalFees) {
        address[] memory _validators = bondedValidators[_id];
        for (uint256 i = 0; i < _validators.length; i++) {
            (uint256 _atnFee, uint256 _ntnFee) = _claimRewards(_id, _validators[i]);
            _atnTotalFees += _atnFee;
            _ntnTotalFees += _ntnFee;
        }
    }

    function _claimRewards(uint256 _id, address _validator) internal returns (uint256 _atnFee, uint256 _ntnFee) {
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated == false) {
            return (0,0);
        }
        _claimRewards(_validator);
        (_atnFee, _ntnFee) = _realiseFees(_id, _validator);
        _account.atnRealisedFee = OFFSET;
        _account.ntnRealisedFee = OFFSET;
    }

    function _initiateValidator(address _validator) private {
        LiquidState _contract = autonity.getValidator(_validator).liquidStateContract;
        // offset by 1
        liquidInfo[_validator] = LiquidInfo(OFFSET, OFFSET, OFFSET, OFFSET, _contract);
    }

    /**
     * @dev adds _validator in bondedValidators[_id] array
     */
    function _addValidator(uint256 _id, address _validator) private {
        if (validatorIdx[_id][_validator] > 0) return;
        address[] storage _validators = bondedValidators[_id];
        _validators.push(_validator);
        // offset by 1 to handle empty value
        validatorIdx[_id][_validator] = _validators.length;
        if (address(liquidInfo[_validator].liquidStateContract) == address(0)) {
            _initiateValidator(_validator);
        }
    }

    /**
     * @dev removes all the validators that are not needed for _id anymore, i.e. any validator
     * that has 0 liquid for _id and all rewards from the validator are claimed
     * @param _id schedule id
     */
    function _clearValidators(uint256 _id) internal {
        // must take a storage pointer
        // otherwise _validators will not be updated when _removeValidator is called
        address[] storage _validators = bondedValidators[_id];
        Account storage _account;
        for (uint256 _idx = 0; _idx < _validators.length ; _idx++) {
            _account = accounts[_id][_validators[_idx]];
            // actual balance is (balance - OFFSET)
            // if liquidBalance is 0, then unrealisedFee is 0
            // if both liquidBalance and realisedFee are 0 and no new bonding is requested then the validator is not needed anymore
            while (
                _account.newBondingRequested == false && _account.liquidBalance == OFFSET
                && _account.atnRealisedFee == OFFSET && _account.ntnRealisedFee == OFFSET
            ) {
                _removeValidator(_id, _validators[_idx]);
                if (_idx >= _validators.length) {
                    break;
                }
                _account = accounts[_id][_validators[_idx]];
            }
        }
    }

    function _removeValidator(uint256 _id, address _validator) private {
        address[] storage _validators = bondedValidators[_id];
        uint256 _idxDelete = validatorIdx[_id][_validator]-1;
        if (_idxDelete+1 == _validators.length) {
            // it is the last validator in the array
            _validators.pop();
            delete validatorIdx[_id][_validator];
            return;
        }
        
        // the _validator to be deleted sits in the middle of the array
        address _lastValidator = _validators[_validators.length-1];
        // replacing the _validator in _idxDelete with _lastValidator, effectively deleting it
        _validators[_idxDelete] = _lastValidator;
        validatorIdx[_id][_lastValidator] = _idxDelete+1;
        // deleting the last one
        _validators.pop();
        delete validatorIdx[_id][_validator];
    }

    /**
     * @dev updates the unclaimed AUT from _validator and also updates lastUnrealisedFeeFactor which is used
     * to compute unrealised fees for accounts
     */
    function _updateUnclaimedReward(address _validator) internal {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        LiquidState _contract = _liquidInfo.liquidStateContract;
        uint256 _totalLiquid = _contract.balanceOf(address(this));
        if (_totalLiquid == 0) {
            return;
        }
        // reward is increased by  OFFSET
        (uint256 _atnReward, uint256 _ntnReward) = _contract.unclaimedRewards(address(this));
        _atnReward += OFFSET;
        _liquidInfo.atnLastUnrealisedFeeFactor += (_atnReward-_liquidInfo.atnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        _liquidInfo.atnUnclaimedRewards = _atnReward;

        _ntnReward += OFFSET;
        _liquidInfo.ntnLastUnrealisedFeeFactor += (_ntnReward-_liquidInfo.ntnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        _liquidInfo.ntnUnclaimedRewards = _ntnReward;
    }

    /**
     * @dev fetches the unclaimedRewards from liquid contract and calculates the changes in lastUnrealisedFeeFactor
     * but does not update any state variable. This function exists only to support _unclaimedRewards to act as view only function
     */
    function _unfetchedFeeFactor(address _validator) private view returns (uint256 _atnFeeFactor, uint256 _ntnFeeFactor) {
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        LiquidState _contract = _liquidInfo.liquidStateContract;
        uint256 _totalLiquid = _contract.balanceOf(address(this));
        if (_totalLiquid > 0) {
            // reward is increased by  OFFSET
            (uint256 _atnReward, uint256 _ntnReward) = _contract.unclaimedRewards(address(this));
            _atnFeeFactor = (_atnReward+OFFSET-_liquidInfo.atnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
            _ntnFeeFactor = (_ntnReward+OFFSET-_liquidInfo.ntnUnclaimedRewards) * FEE_FACTOR_UNIT_RECIP / _totalLiquid;
        }
    }

    /**
     * @dev calculates the rewards yet to claim for _id from _validator 
     */
    function _unclaimedRewards(uint256 _id, address _validator) internal view returns (uint256 _atnReward, uint256 _ntnReward) {
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated == false) {
            // account does not exist
            return (0,0);
        }
        (uint256 _atnLastUnrealisedFeeFactor, uint256 _ntnLastUnrealisedFeeFactor) = _unfetchedFeeFactor(_validator);
        LiquidInfo storage _liquidInfo = liquidInfo[_validator];
        _atnLastUnrealisedFeeFactor += _liquidInfo.atnLastUnrealisedFeeFactor;
        _ntnLastUnrealisedFeeFactor += _liquidInfo.ntnLastUnrealisedFeeFactor;

        uint256 _balance = _account.liquidBalance-OFFSET;
        // remove offset from realisedFee
        _atnReward = _account.atnRealisedFee - OFFSET
                    + _computeUnrealisedFees(
                        _balance, _account.atnUnrealisedFeeFactor, _atnLastUnrealisedFeeFactor
                    );

        _ntnReward = _account.ntnRealisedFee - OFFSET
                    + _computeUnrealisedFees(
                        _balance, _account.ntnUnrealisedFeeFactor, _ntnLastUnrealisedFeeFactor
                    );
    }

    /**
     * @dev calculates the rewards yet to claim for _id
     */
    function _unclaimedRewards(uint256 _id) internal view returns (uint256 _atnTotalFee, uint256 _ntnTotalFee) {
        address[] memory _validators = bondedValidators[_id];
        for (uint256 i = 0; i < _validators.length; i++) {
            (uint256 _atnReward, uint256 _ntnReward) = _unclaimedRewards(_id, _validators[i]);
            _atnTotalFee += _atnReward;
            _ntnTotalFee += _ntnReward;
        }
    }

    function _isNewBondingRequested(uint256 _id, address _validator) internal view returns (bool) {
        return accounts[_id][_validator].newBondingRequested;
    }

    function _liquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated) {
            return _account.liquidBalance - OFFSET;
        }
        return 0;
    }

    function _unlockedLiquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated) {
            return _account.liquidBalance - _account.lockedLiquidBalance;
        }
        return 0;
    }

    function _lockedLiquidBalanceOf(uint256 _id, address _validator) internal view returns (uint256) {
        Account storage _account = accounts[_id][_validator];
        if (_account.initiated) {
            return _account.lockedLiquidBalance - OFFSET;
        }
        return 0;
    }

    /**
     * @dev returns the list of validator addresses wich are bonded to schedule _id assigned to _account
     */
    function _bondedValidators(uint256 _id) internal view returns (address[] memory) {
        return bondedValidators[_id];
    }

    function _liquidStateContract(address _validator) internal view returns (LiquidState) {
        return liquidInfo[_validator].liquidStateContract;
    }

}