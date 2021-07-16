pragma solidity ^0.8.3;

import "./interfaces/IERC20.sol";

contract Liquid is IERC20 {
    address autonityContract; //not hardcoded for testing purposes

    address validator;
    address payable treasury;
    uint256 public supply;
    uint256 commissionRate;

    address payable[] accountList;
    struct Account {
        uint256 balance;
        uint256 index; // index in accountList for O(1) access
    }

    mapping (address => Account) accounts;
    mapping (address  => mapping (address => uint256)) private allowances;


    constructor (address _validator, address payable _treasury, uint256 _commissionRate) {
        validator = _validator;
        treasury = _treasury;
        commissionRate = _commissionRate;
        autonityContract = msg.sender;
    }

    function mint(address payable _recipient, uint256 _amount) public onlyAutonity {
        require(_amount > 0, "amount must be strictly positive");
        if(accounts[_recipient].balance == 0) {
            accountList.push(_recipient);
            Account memory newAccount = Account(0, accountList.length-1);
            accounts[_recipient] = newAccount;
        }
        accounts[_recipient].balance += _amount;
        supply += _amount;
    }

    function burn(address payable _addr, uint256 _amount) public onlyAutonity {
        require(accounts[_addr].balance >= _amount, "address balance not sufficient");
        accounts[_addr].balance -= _amount; // test to make sure SafeMath not needed here
        supply -= _amount;
        if(accounts[_addr].balance == 0){
            _removeAddress(_addr);
        }
    }

    function redeem(uint256 _amount) public  view{
        require(accounts[msg.sender].balance >= _amount, "sender's balance has to be greater than the specified amount");
        // call the autonity Contract redeem function
    }

    function redistribute() public payable onlyAutonity  {
        uint256 _totalFees = msg.value;
        uint256 _validatorFees = (_totalFees  * commissionRate) / 100000;
        treasury.transfer(_validatorFees);
        _totalFees -= _validatorFees;
        for (uint256 i=0; i < accountList.length; i++) {
            address payable _addr = accountList[i];
            uint256 _reward = (_totalFees * accounts[_addr].balance) / supply;
            _addr.transfer(_reward);
        }
    }



    /**
    * @notice Moves `amount` NEW stake tokens from the caller's account to `recipient`.
    *
    * @return Returns a boolean value indicating whether the operation succeeded.
    *
    * @dev Emits a {Transfer} event. Implementation of {IERC20 transfer}
    */
    function transfer(address _recipient, uint256 _amount) external override returns (bool) {
        _transfer(msg.sender, _recipient, _amount);
        return true;
    }


    /**
    * @notice Returns the total amount of stake token issued.
    */
    function totalSupply() external view override returns (uint256) {
        return supply;
    }

    /**
    * @notice Returns the amount of stake token held by the account (ERC-20).
    */
    function balanceOf(address _account) external view override returns (uint256) {
        return accounts[_account].balance;
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
    function transferFrom(address sender, address recipient, uint256 amount) external override returns (bool){
        _transfer(sender, recipient, amount);
        uint256 newAllowance = allowances[sender][msg.sender] - amount;
        _approve(sender, msg.sender, newAllowance);
        return true;
    }

    /**
     * @dev See {IERC20-approve}.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function approve(address spender, uint256 amount) external override returns (bool) {
        _approve(msg.sender, spender, amount);
        return true;
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
    function _approve(address _owner, address _spender, uint256 _amount) internal virtual {
        require(_owner != address(0), "ERC20: approve from the zero address");
        require(_spender != address(0), "ERC20: approve to the zero address");

        allowances[_owner][_spender] = _amount;
        emit Approval(_owner, _spender, _amount);
    }

    function _transfer(address _sender, address _recipient, uint256 _amount) internal {
        require(accounts[_sender].balance >= _amount, "amount exceeds balance");
        accounts[_sender].balance -= _amount;
        accounts[_recipient].balance += _amount;
        emit Transfer(_sender, _recipient, _amount);
    }

    function _removeAddress(address payable _addr) internal {
        Account memory _account = accounts[_addr];
        uint256 _idx = _account.index;
        address payable _replacement = accountList[accountList.length-1];
        accountList[_idx] =  _replacement;
        accounts[_replacement].index = _idx;
        accountList.pop();
    }

    /**
    * @dev See {IERC20-allowance}.
    */
    function allowance(address _owner, address _spender) external view override returns (uint256) {
        return allowances[_owner][_spender];
    }

    modifier onlyAutonity {
        require(msg.sender == autonityContract, "Call restricted to the Autonity Contract");
        _;
    }

}
