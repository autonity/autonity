// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;

import "../interfaces/ILiquid.sol";
import "./LiquidStorage.sol";

// References:
//
// - [BATOG18]
//   "Scalable Reward Distribution on the Ethereum Blockchain"
//   Bogdan Batog, Lucian Boca, Nick Johnson
//   solhint-disable-next-line max-line-length
//   https://uploads-ssl.webflow.com/5ad71ffeb79acc67c8bcdaba/5ad8d1193a40977462982470_scalable-reward-distribution-paper.pdf
//
// - [OJHA19]
//   "F1 Fee Distribution"
//   Dev Ojha, Christopher Goes
//   solhint-disable-next-line max-line-length
//   https://drops.dagstuhl.de/opus/volltexte/2020/11974/pdf/OASIcs-Tokenomics-2019-10.pdf
//
//
// Implementation Notes (see README.md for a description of the
// algorithm being demonstrated here).
//
//   The interface here closely matches `Liquid` in autonity:
//
//     https://github.com/clearmatics/autonity-internal/blob/dpos_ya/autonity/solidity/contracts/
//
//   Instead of keeping the full array of all {f_i} values, we track
//   f_{i-1} (corresponding to the previous epoch), and maintain a map
//   from delegator to the value f_{i-1} at time fees were last
//   realised.  That is, `_unrealisedFeeFactors[A]` is always the
//   value f_{a-1} above for delegator A.  In this way, when delegator
//   A's LNEW balance falls to 0, their entry in _unrealisedFeeFactors
//   can be removed, and the total state size does not increase with
//   the number of epochs.
//
//   These values f_i are referred to a "fee factors" in the
//   implementation.
//

contract LiquidLogic is ILiquid, LiquidStorage {

    // TODO: Better solution to address the fractional terms in fee computations?

    uint256 public constant FEE_FACTOR_UNIT_RECIP_DECIMALS = 9;
    uint256 public constant FEE_FACTOR_UNIT_RECIP = 10 ** FEE_FACTOR_UNIT_RECIP_DECIMALS;
    uint256 public constant COMMISSION_RATE_DECIMALS = 4;
    uint256 public constant COMMISSION_RATE_SCALE_FACTOR = 10 ** COMMISSION_RATE_DECIMALS;

    constructor() {
        autonityContract = Autonity(payable(msg.sender));
    }

    /**
     * @notice Redistribute fees, called once per epoch by the autonity contract.
     * Update lastUnrealisedFeeFactor and transfer treasury fees.
     * @custom:restricted-to the autonity contract
     */
    function redistribute(uint256 _ntnReward) external virtual payable onlyAutonity returns (uint256) {
        uint256 _atnReward = msg.value;
        // Step 1 : transfer entitled amount of fees to validator's
        // treasury account.
        uint256 _atnValidatorReward = _calculateValidatorCommission(_atnReward);
        _atnReward -= _atnValidatorReward;
        (bool _sent, ) = treasury.call{value: _atnValidatorReward, gas:2300}("");
        if (_sent == false) {
            treasuryUnclaimedATN += _atnValidatorReward;
        }

        uint256 _ntnValidatorReward = _calculateValidatorCommission(_ntnReward);
        if (_ntnReward > 0) {
            autonityContract.autobond(validator, _ntnValidatorReward, _ntnReward - _ntnValidatorReward);
        }

        // Step 2 : perform redistribution amongst liquid stake token
        // holders for this validator.
        uint256 _atnFeeFactorThisReward = (_atnReward * FEE_FACTOR_UNIT_RECIP) / supply;
        atnLastUnrealisedFeeFactor = atnLastUnrealisedFeeFactor + _atnFeeFactorThisReward;

        // Compute the maximum amount that can be claimed after
        // rounding.
        uint256 _atnMaxClaimable = (_atnFeeFactorThisReward * supply) / FEE_FACTOR_UNIT_RECIP;
        return _atnValidatorReward + _atnMaxClaimable;
    }

    /**
     * @notice Mint new tokens and transfer them to the target account.
     * @custom:restricted-to the autonity contract.
     */
    function mint(address _account, uint256 _amount) external virtual onlyAutonity {
        _increaseBalance(_account, _amount);
        emit Transfer(address(0), _account, _amount);
    }

    /**
     * @notice Burn tokens from the target account.
     * @custom:restricted-to Restricted to the autonity contract.
     */
    function burn(address _account, uint256 _amount) external virtual onlyAutonity {
        _requireAndDecreaseBalance(_account, _amount);
        emit Transfer(_account, address(0), _amount);
    }

    /**
     * @notice Send the unclaimed ATN entitled to treasury to treasury account
     */
    function claimTreasuryATN() external virtual {
        require(msg.sender == treasury, "only treasury can claim his reward");
        (bool _sent, ) = treasury.call{value: treasuryUnclaimedATN}("");
        require(_sent, "failed to send ATN");
        treasuryUnclaimedATN = 0;
    }

    /**
     * @notice Withdraws all fees earned so far by the caller.
     */
    function claimRewards() external virtual {
        uint256 _atnRealisedFees = _realiseFees(msg.sender);
        delete atnRealisedFees[msg.sender];

        //   solhint-disable-next-line avoid-low-level-calls
        (bool _sent, ) = msg.sender.call{value: _atnRealisedFees}("");
        require(_sent, "Failed to send ATN");
    }

    /**
     * @notice Moves `_amount` LNEW tokens from the caller's account to the recipient `_to`.
     *
     * @return _success a boolean value indicating whether the operation succeeded.
     *
     * @dev Emits a {Transfer} event. Implementation of {IERC20 transfer}
     */
    function transfer(address _to, uint256 _amount) external virtual returns (bool _success) {
        _requireAndDecreaseBalance(msg.sender, _amount);
        _increaseBalance(_to, _amount);
        emit Transfer(msg.sender, _to, _amount);
        return true;
    }

    /**
     * @dev See {IERC20-approve}.
     *
     * Requirements:
     *
     * - `_spender` cannot be the zero address.
     */
    function approve(address _spender, uint256 _amount) external virtual returns (bool) {
        _approve(msg.sender, _spender, _amount);
        return true;
    }

    /**
     * @dev See {IERC20-transferFrom}.
     *
     * Emits an {Approval} event indicating the updated allowance.
     *
     * Requirements:
     *
     * - `_sender` and `_recipient` must be allowed to hold stake.
     * - `_sender` must have a balance of at least `_amount`.
     * - the caller must have allowance for ``_sender``'s tokens of at least
     * `_amount`.
     */
    function transferFrom(address _sender, address _recipient, uint256 _amount) external virtual returns (bool) {
        uint256 _currentAllowance = allowances[_sender][msg.sender];
        require(_currentAllowance >= _amount, "ERC20: transfer amount exceeds allowance");
        _approve(_sender, msg.sender, _currentAllowance - _amount);

        _requireAndDecreaseBalance(_sender, _amount);
        _increaseBalance(_recipient, _amount);
        emit Transfer(_sender, _recipient, _amount);
        return true;
    }


    /**
     * @notice Setter for the commission rate, restricted to the Autonity Contract.
     * @param _rate New rate.
     */
    function setCommissionRate(uint256 _rate) external virtual onlyAutonity {
        commissionRate = _rate;
    }

    /**
     * @notice Add amount to the locked funds, restricted to the Autonity Contract.
     * @param _account address of the account to lock funds .
              _amount LNTN amount of tokens to lock.
     */
    function lock(address _account, uint256 _amount) external virtual onlyAutonity {
        require(balances[_account] - lockedBalances[_account] >= _amount, "can't lock more funds than available");
        lockedBalances[_account] += _amount;
    }

    /**
     * @notice Unlock the locked funds, restricted to the Autonity Contract.
     * @param _account address of the account to lock funds .
              _amount LNTN amount of tokens to lock.
     */
    function unlock(address _account, uint256 _amount) external virtual onlyAutonity {
        require(lockedBalances[_account] >= _amount, "can't unlock more funds than locked");
        lockedBalances[_account] -= _amount;
    }

    /**
     * @dev It is not expected to fall into the fallback function. Implemeted fallback() to get a reverting message.
     */
    fallback() payable external virtual {
        revert("fallback not implemented for LiquidLogic");
    }

    /**
     * @dev To receive ATN.
     */
    receive() payable external virtual {}

    /**
     ============================================================

        Internals

     ============================================================
     */

    function _increaseBalance(address _delegator, uint256 _value) private {
        _realiseFees(_delegator); //always updates fee factor
        balances[_delegator] += _value;
        // when transferring, this value will just be decreased
        // again by the same amount.
        supply += _value;
    }

    function _requireAndDecreaseBalance(address _delegator, uint256 _value) private {
        _realiseFees(_delegator); // always updates fee factor
        uint256 _balance = balances[_delegator];
        require(_value <= _balance - lockedBalances[_delegator], "insufficient unlocked funds");
        balances[_delegator] = _balance - _value;

        if (_value == _balance) { // aka balances[_delegator] == 0
            // get back some gas
            delete atnUnrealisedFeeFactors[_delegator];
        }
        // when transferring, this value will just be increased
        // again by the same amount.
        supply -= _value;
    }


    /**
     * @dev Compute all unrealised fees, update the fee balance and reset
     * the unrealised fee factor for the given participant.  This
     * function ALWAYS sets the unrealised fee factor for the
     * delegator, so should not be called if the delegators balance is
     * known to be zero (or the caller should handle this case itself).
     * @param _delegator, the target account to compute fees.
     * @return _atnRealisedFees that is the calculated amount of ATN that
     * the delegator is entitled to withdraw.
     */
    function _realiseFees(address _delegator) private returns (uint256 _atnRealisedFees) {
        uint256 _balance = balances[_delegator];
        uint256 _atnUnrealisedFee = _computeUnrealisedFees(_balance, atnLastUnrealisedFeeFactor, atnUnrealisedFeeFactors[_delegator]);

        _atnRealisedFees = atnRealisedFees[_delegator] + _atnUnrealisedFee;
        atnRealisedFees[_delegator] = _atnRealisedFees;
        atnUnrealisedFeeFactors[_delegator] = atnLastUnrealisedFeeFactor;
    }

    /**
     * @dev Computes atn unrealised fees.
     * @param _balance LNTN balance
     * @param _lastUnrealisedFeeFactor last unrealised fee factor for atn
     * @param _unrealisedFeeFactors unrealised fee factor for atn
     * @return uint256 atn unrealised fee.
     */
    function _computeUnrealisedFees(uint256 _balance, uint256 _lastUnrealisedFeeFactor, uint256 _unrealisedFeeFactors)
        private pure returns (uint256) {

        // Early out if _lnewBalance == 0
        if (_balance == 0) {
            return 0;
        }

        // If the delegator has a non-zero balance, there should
        // be a valid _unrealisedFeeFactors entry.  Currently can't
        // tell the difference between the 0 (when delegatinng from
        // the start) or a missing entry.

        // Unrealised fees are:
        //     balance x (f_{last_epoch} - f_{deposit_epoch})

        // FEE_FACTOR_UNIT_RECIP = 10^9 won't cause overflow
        return ((_lastUnrealisedFeeFactor - _unrealisedFeeFactors) * _balance) / FEE_FACTOR_UNIT_RECIP;
    }

    /**
     * @dev Sets `_amount` as the allowance of `_spender` over the `_owner` s tokens.
     *
     * This internal function is equivalent to `_approve`, and can be used to
     * e.g. set automatic allowances for certain subsystems, etc.
     *
     * Emits an {Approval} event.
     *
     */
    function _approve(address _owner, address _spender, uint256 _amount) internal virtual {
        require(_owner != address(0), "ERC20: approve from the zero address");
        require(_spender != address(0), "ERC20: approve to the zero address");


        allowances[_owner][_spender] = _amount;
        emit Approval(_owner, _spender, _amount);
    }

    function _calculateValidatorCommission(uint256 _reward) internal virtual view returns (uint256) {
        uint256 _commission = (_reward * commissionRate) / COMMISSION_RATE_SCALE_FACTOR;
        return _commission;
    }

    /*
     ============================================================
        Getters
     ============================================================
     */

    /**
     * @notice Calculates the total claimable fees (ATN) earned by the delegator to-date.
     * @param _account Delegator account.
     */
    function unclaimedRewards(address _account) external virtual view returns (uint256) {
        uint256 _balance = balances[_account];
        uint256 _atnUnrealisedFee = _computeUnrealisedFees(_balance, atnLastUnrealisedFeeFactor, atnUnrealisedFeeFactors[_account]);
        return atnRealisedFees[_account] + _atnUnrealisedFee;
    }

    /**
     * @notice Returns the total amount of stake token issued.
     */
    function totalSupply() external virtual view returns (uint256) {
        return supply;
    }

    /**
     * @return uint8 the number of decimals the LNTN token uses.
     * @dev ERC-20 Optional.
     */
    function decimals() external virtual pure returns (uint8) {
        return DECIMALS;
    }

    /**
     * @notice Returns the amount of liquid newtons held by the account (ERC-20).
     */
    function balanceOf(address _delegator) external virtual view returns (uint256) {
        return balances[_delegator];
    }

    /**
     * @notice Returns the amount of locked liquid newtons held by the account.
     */
    function lockedBalanceOf(address _delegator) external virtual view returns (uint256) {
        return lockedBalances[_delegator];
    }

    /**
     * @notice Returns the amount of unlocked liquid newtons held by the account.
     */
    function unlockedBalanceOf(address _delegator) external virtual view returns (uint256) {
        return  balances[_delegator] - lockedBalances[_delegator];
    }

    /**
     * @notice See {IERC20-allowance}.
     */
    function allowance(address _owner, address _spender) external virtual view returns (uint256) {
        return allowances[_owner][_spender];
    }

    function name() external virtual view returns (string memory) {
        return liquidName;
    }

    function symbol() external virtual view returns (string memory) {
        return liquidSymbol;
    }

    function getValidator() external virtual view returns (address) {
        return validator;
    }

    function getTreasury() external virtual view returns (address) {
        return treasury;
    }

    function getCommissionRate() external virtual view returns (uint256) {
        return commissionRate;
    }

    /**
     * @notice Returns the ATN amount that is yet to claim by treasury.
     * Call function `claimTreasuryATN()` to claim.
     */
    function getTreasuryUnclaimedATN() external virtual view returns (uint256) {
        return treasuryUnclaimedATN;
    }


    /*
     ============================================================

        Modifiers

     ============================================================
     */

    modifier onlyAutonity {
        require(
            msg.sender == address(autonityContract),
            "Call restricted to the Autonity Contract");
        _;
    }
}
