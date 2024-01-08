// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;
import "./interfaces/IERC20.sol";
import {DECIMALS} from "./Autonity.sol";
import {LiquidLogic} from "./LiquidLogic.sol";

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

contract LiquidState is IERC20
{
    // storage layout - this must be compatible with LiquidLogic
    address autonityContract; //not hardcoded for testing purposes

    mapping(address => uint256) private balances;
    mapping(address => uint256) private lockedBalances;

    mapping(address => mapping (address => uint256)) private allowances;
    uint256 private supply;

    mapping(address => uint256) private realisedFees;
    mapping(address => uint256) private unrealisedFeeFactors;

    uint256 private lastUnrealisedFeeFactor;


    address payable public treasury;
    uint256 public commissionRate;
    //end of storage layout restrictions

    string private _name;
    string private _symbol;

    address public validator;

    // this must be always last, since logic is delegated to LiquidLogic and
    // they must use same storage layout
    address public liquidLogic;

    constructor(
        address _logic,
        address _validator,
        address payable _treasury,
        uint256 _commissionRate,
        string memory _index)
    {

        // this helps to DRY, but also checks if passed logic address is valid
        uint256 _commissionRatePrecision = LiquidLogic(_logic).COMMISSION_RATE_PRECISION();
//        uint256 _commissionRatePrecision = 10_000;

        // commissionRate <= 1.0
        require(_commissionRate <= _commissionRatePrecision);


        liquidLogic = _logic;
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
        (bool success, bytes memory data) = liquidLogic.delegatecall(
            abi.encodeWithSignature("redistribute()")
        );
        if (!success) {
            revert("call to logic redistribute failed");
        }

        return abi.decode(data, (uint256));
    }

    /**
    * @notice Mint new tokens and transfer them to the target account.
    * @dev Restricted to the autonity contract.
    */
    function mint(address _account, uint256 _amount) external onlyAutonity
    {
        (bool success,) = liquidLogic.delegatecall(
            abi.encodeWithSignature("mint(address,uint256)", _account, _amount)
        );
        if (!success) {
            revert("call to logic mint failed");
        }
    }

    /**
    * @notice Burn tokens from the target account.
    * @dev Restricted to the autonity contract.
    */
    function burn(address _account, uint256 _amount) external onlyAutonity
    {
        (bool success,) = liquidLogic.delegatecall(
            abi.encodeWithSignature("burn(address,uint256)", _account, _amount)
        );
        if (!success) {
            revert("call to logic burn failed");
        }
    }

    /**
    * @notice  Returns the total claimable fees (AUT) earned by the delegator to-date.
    * @param _account the delegator account.
    */
    function unclaimedRewards(address _account) external returns(uint256)
    {
        (bool success, bytes memory data) = liquidLogic.delegatecall(
            abi.encodeWithSignature("unclaimedRewards(address)", _account)
        );
        if (!success) {
            revert("call to logic unclaimedRewards failed");
        }

        return abi.decode(data, (uint256));
    }

    /**
    * @notice Withdraws all fees earned so far by the caller.
    */
    function claimRewards() external
    {
        (bool success,) = liquidLogic.delegatecall(
            abi.encodeWithSignature("claimRewards()")
        );
        if (!success) {
            revert("call to logic claimRewards failed");
        }
    }

    /**
    * @notice Returns the amount of unlocked liquid newtons held by the account.
    */
    function unlockedBalanceOf(address _delegator)
        external returns (uint256)
    {
        (bool success, bytes memory data) = liquidLogic.delegatecall(
            abi.encodeWithSignature("unlockedBalanceOf(address)", _delegator)
        );
        if (!success) {
            revert("call to logic unlockedBalanceOf failed");
        }

        return abi.decode(data, (uint256));
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
        (bool success, bytes memory data) = liquidLogic.delegatecall(
            abi.encodeWithSignature("transfer(address,uint256)", _to, _amount)
        );
        if (!success) {
            revert("call to logic transfer failed");
        }

        return abi.decode(data, (bool));
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
        (bool success, bytes memory data) = liquidLogic.delegatecall(
            abi.encodeWithSignature("approve(address,uint256)", _spender, _amount)
        );
        if (!success) {
            revert("call to logic approve failed");
        }

        return abi.decode(data, (bool));
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
        (bool success, bytes memory data) = liquidLogic.delegatecall(
            abi.encodeWithSignature("transferFrom(address,address,uint256)", _sender, _recipient, _amount)
        );
        if (!success) {
            revert("call to logic transferFrom failed");
        }

        return abi.decode(data, (bool));
    }

    /**
      * @notice Setter for the commission rate, restricted to the Autonity Contract.
      * @param _rate New rate.
      */
    function setCommissionRate(uint256 _rate) public onlyAutonity {

        (bool success,) = liquidLogic.delegatecall(
            abi.encodeWithSignature("setCommissionRate(uint256)", _rate)
        );
        if (!success) {
            revert("call to logic setCommissionRate failed");
        }
    }




    /**
      * @notice Add amount to the locked funds, restricted to the Autonity Contract.
      * @param _account address of the account to lock funds .
               _amount LNTN amount of tokens to lock.
      */
    function lock(address _account, uint256 _amount) public onlyAutonity {

        (bool success, bytes memory data) = liquidLogic.delegatecall(
            abi.encodeWithSignature("lock(address,uint256)", _account, _amount)
        );
        if (success == false) {
            // if there is a return reason string
            if (data.length > 0) {
                // bubble up any reason for revert
                assembly {
                    let returndata_size := mload(data)
                    revert(add(32, data), returndata_size)
                }
            } else {
                revert("Function call 'lock' reverted");
            }
        }
    }

    /**
      * @notice Unlock the locked funds, restricted to the Autonity Contract.
      * @param _account address of the account to lock funds .
               _amount LNTN amount of tokens to lock.
      */
    function unlock(address _account, uint256 _amount) public onlyAutonity {
        (bool success,) = liquidLogic.delegatecall(
            abi.encodeWithSignature("unlock(address,uint256)", _account, _amount)
        );
        if (!success) {
            revert("call to logic unlock failed");
        }
    }


    /*
============================================================

    Simple views, directly accessing data without any calculations

============================================================
*/


    /**
* @notice Returns the amount of locked liquid newtons held by the account.
    */
    function lockedBalanceOf(address _delegator)
    external view returns (uint256)
    {
        return lockedBalances[_delegator];
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
