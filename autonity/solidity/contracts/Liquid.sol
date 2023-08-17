// Copyright (c) 2015-2021 Clearmatics Technologies Ltd
//
// SPDX-License-Identifier: UNLICENSED

pragma solidity ^0.8.3;
import "./interfaces/IERC20.sol";
import {DECIMALS} from "./Autonity.sol";

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

contract Liquid is IERC20
{
    address autonityContract; //not hardcoded for testing purposes

    // TODO: Better solution to address the fractional terms in fee
    // computations?
    //
    // If fee computations are to be performed to 9 decimal places,
    // this value should be 1,000,000,000.
    uint256 public constant FEE_FACTOR_UNIT_RECIP = 1_000_000_000;
    uint256 public constant COMMISSION_RATE_PRECISION = 10_000;

    mapping(address => uint256) private balances;
    mapping(address => uint256) private lockedBalances;

    mapping(address => mapping (address => uint256)) private allowances;
    uint256 private supply;

    mapping(address => uint256) private realisedFees;
    mapping(address => uint256) private unrealisedFeeFactors;
    uint256 private lastUnrealisedFeeFactor;

    string private _name;
    string private _symbol;

    address public validator;
    address payable public treasury;
    uint256 public commissionRate;

    constructor(
        address _validator,
        address payable _treasury,
        uint256 _commissionRate,
        string memory _index)
    {
        // commissionRate <= 1.0
        require(_commissionRate <= COMMISSION_RATE_PRECISION);

        validator = _validator;
        treasury = _treasury;
        commissionRate = _commissionRate;
        _name = string.concat("LNTN-", _index);
        _symbol = string.concat("LNTN-", _index);
        autonityContract = msg.sender;
    }

    /**
    * @notice Redistribute fees, called once per epoch by the autonity contract.
    * Update lastUnrealisedFeeFactor and transfer treasury fees.
    * @dev Restricted to the autonity contract.
    */
    function redistribute() external payable onlyAutonity returns (uint256)
    {
        uint256 _reward = msg.value;
        // Step 1 : transfer entitled amount of fees to validator's
        // treasury account.
        uint256 _validatorReward =
            (_reward * commissionRate) / COMMISSION_RATE_PRECISION;
        require(_validatorReward <= _reward, "invalid validator reward");
        _reward -= _validatorReward;

        // TODO: handle failure
        treasury.call{value: _validatorReward, gas:2300}("");

        // Step 2 : perform redistribution amongst liquid stake token
        // holders for this validator.
        uint256 _feeFactorThisReward =
            ((_reward * FEE_FACTOR_UNIT_RECIP) / supply);
        lastUnrealisedFeeFactor =
            lastUnrealisedFeeFactor + _feeFactorThisReward;

        // Compute the maximum amount that can be claimed after
        // rounding.
        uint256 _maxClaimable =
            (_feeFactorThisReward * supply) / FEE_FACTOR_UNIT_RECIP;
        return _validatorReward + _maxClaimable;
    }

    /**
    * @notice Mint new tokens and transfer them to the target account.
    * @dev Restricted to the autonity contract.
    */
    function mint(address _account, uint256 _amount) external onlyAutonity
    {
        _increaseBalance(_account, _amount);
        emit Transfer(address(0), _account, _amount);
    }

    /**
    * @notice Burn tokens from the target account.
    * @dev Restricted to the autonity contract.
    */
    function burn(address _account, uint256 _amount) external onlyAutonity
    {
        _requireAndDecreaseBalance(_account, _amount);
        emit Transfer(_account, address(0), _amount);
    }

    /**
    * @notice  Returns the total claimable fees (AUT) earned by the delegator to-date.
    * @param _account the delegator account.
    */
    function unclaimedRewards(address _account) external view returns(uint256)
    {
        return realisedFees[_account] + _computeUnrealisedFees(_account);
    }

    /**
    * @notice Withdraws all fees earned so far by the caller.
    */
    function claimRewards() external
    {
        uint256 totalFees = _realiseFees(msg.sender);
        delete realisedFees[msg.sender];

        // Send the AUT
        //   solhint-disable-next-line avoid-low-level-calls
        (bool sent, ) = msg.sender.call{value: totalFees}("");
        require(sent, "Failed to send Ether");
    }

    /**
    * @notice Returns the total amount of stake token issued.
    */
    function totalSupply() public view override returns (uint256)
    {
        return supply;
    }

    /**
    * @return the number of decimals the LNTN token uses.
    * @dev ERC-20 Optional.
    */
    function decimals() public pure returns (uint8) {
        return DECIMALS;
    }

    /**
    * @notice Returns the amount of liquid newtons held by the account (ERC-20).
    */
    function balanceOf(address _delegator)
        external view override returns (uint256)
    {
        return balances[_delegator];
    }

    /**
    * @notice Returns the amount of locked liquid newtons held by the account.
    */
    function lockedBalanceOf(address _delegator)
        external view returns (uint256)
    {
        return lockedBalances[_delegator];
    }

    /**
    * @notice Returns the amount of unlocked liquid newtons held by the account.
    */
    function unlockedBalanceOf(address _delegator)
        external view returns (uint256)
    {
        return  balances[_delegator] - lockedBalances[_delegator];
    }

    /**
    * @notice Moves `_amount` LNEW tokens from the caller's account to the recipient `_to`.
    *
    * @return _success a boolean value indicating whether the operation succeeded.
    *
    * @dev Emits a {Transfer} event. Implementation of {IERC20 transfer}
    */
    function transfer(address _to, uint256 _amount)
        public override returns (bool _success)
    {
        _requireAndDecreaseBalance(msg.sender, _amount);
        _increaseBalance(_to, _amount);
        emit Transfer(msg.sender, _to, _amount);
        return true;
    }

    /**
    * @dev See {IERC20-allowance}.
    */
    function allowance(address _owner, address _spender)
        public view override returns (uint256)
    {
        return allowances[_owner][_spender];
    }

    /**
     * @dev See {IERC20-approve}.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function approve(address _spender, uint256 _amount) public override returns (bool)
    {
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
      * - `sender` and `recipient` must be allowed to hold stake.
      * - `sender` must have a balance of at least `amount`.
      * - the caller must have allowance for ``sender``'s tokens of at least
      * `amount`.
      */
    function transferFrom(address _sender, address _recipient, uint256 _amount)
        public override returns (bool _success)
    {
        uint256 currentAllowance = allowances[_sender][msg.sender];
        require(currentAllowance >= _amount, "ERC20: transfer amount exceeds allowance");
        _approve(_sender, msg.sender, currentAllowance - _amount);

        _requireAndDecreaseBalance(_sender, _amount);
        _increaseBalance(_recipient, _amount);
        emit Transfer(_sender, _recipient, _amount);
        return true;
    }


    /**
      * @notice Setter for the commission rate, restricted to the Autonity Contract.
      * @param _rate New rate.
      */
    function setCommissionRate(uint256 _rate) public onlyAutonity {
        commissionRate = _rate;
    }

    /**
      * @notice Add amount to the locked funds, restricted to the Autonity Contract.
      * @param _account address of the account to lock funds .
               _amount LNTN amount of tokens to lock.
      */
    function lock(address _account, uint256 _amount) public onlyAutonity {
        require(balances[_account] - lockedBalances[_account] >= _amount, "can't lock more funds than available");
        lockedBalances[_account] += _amount;
    }

    /**
      * @notice Unlock the locked funds, restricted to the Autonity Contract.
      * @param _account address of the account to lock funds .
               _amount LNTN amount of tokens to lock.
      */
    function unlock(address _account, uint256 _amount) public onlyAutonity {
        require(lockedBalances[_account] >= _amount, "can't unlock more funds than locked");
        lockedBalances[_account] -= _amount;
    }

    /**
    * @notice returns the name of this Liquid Newton contract
    */
    function name() external view returns (string memory){
        return _name;
    }

    /**
    * @notice returns the symbol of this Liquid Newton contract
    */
    function symbol() external view returns (string memory){
        return _symbol;
    }

    /**
    ============================================================

        Internals

    ============================================================
    */

    function _increaseBalance(address _delegator, uint256 _value) private
    {
        _realiseFees(_delegator); //always updates fee factor
        balances[_delegator] += _value;
        // when transferring, this value will just be decreased
        // again by the same amount.
        supply += _value;
    }

    function _requireAndDecreaseBalance(address _delegator, uint256 _value)
        private
    {
        _realiseFees(_delegator); // always updates fee factor
        uint256 _balance = balances[_delegator];
        require(_value <= _balance - lockedBalances[_delegator], "insufficient unlocked funds");
        balances[_delegator] = _balance - _value;

        if (_value == _balance) { // aka balances[_delegator] == 0
            // get back some gas
            delete unrealisedFeeFactors[_delegator];
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
    * @return _realisedFees that is the calculated amount of AUT that
    * the delegator is entitled to withdraw.
    */
    function _realiseFees(address _delegator) private
        returns (uint256 _realisedFees)
    {
        uint256 _unrealisedFees = _computeUnrealisedFees(_delegator);
        _realisedFees = realisedFees[_delegator] + _unrealisedFees;
        realisedFees[_delegator] = _realisedFees;
        unrealisedFeeFactors[_delegator] = lastUnrealisedFeeFactor;
    }

    function _computeUnrealisedFees(address _delegator)
        private view returns (uint256)
    {
        // TODO: save looking up the LNEW balance multiple times by passing it
        // in here.

         uint256 _stakerBalance = balances[_delegator];

         // Early out if _lnewBalance == 0
         if (_stakerBalance == 0) {
             return 0;
         }

        // If the delegator has a non-zero balance, there should
        // be a valid _unrealisedFeeFactors entry.  Currently can't
        // tell the difference between the 0 (when delegatinng from
        // the start) or a missing entry.

        // Unrealised fees are:
        //     balance x (f_{last_epoch} - f_{deposit_epoch})

        uint256 _unrealisedFeeFactor =
            lastUnrealisedFeeFactor - unrealisedFeeFactors[_delegator];

        // FEE_FACTOR_UNIT_RECIP = 10^9 won't cause overflow
        uint256 _unrealisedFee =
            (_unrealisedFeeFactor * _stakerBalance) / FEE_FACTOR_UNIT_RECIP;

        return _unrealisedFee;
    }

    /**
     * @dev Sets `amount` as the allowance of `spender` over the `owner` s tokens.
     *
     * This internal function is equivalent to `approve`, and can be used to
     * e.g. set automatic allowances for certain subsystems, etc.
     *
     * Emits an {Approval} event.
     *
     */
    function _approve(address _owner, address _spender, uint256 _amount)
        internal virtual
    {
        require(_owner != address(0), "ERC20: approve from the zero address");
        require(_spender != address(0), "ERC20: approve to the zero address");


        allowances[_owner][_spender] = _amount;
        emit Approval(_owner, _spender, _amount);
    }


    /*
    ============================================================

        Modifiers

    ============================================================
    */

    modifier onlyAutonity
    {
        require(
            msg.sender == autonityContract,
            "Call restricted to the Autonity Contract");
        _;
    }
}
