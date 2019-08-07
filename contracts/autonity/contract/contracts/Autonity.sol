pragma solidity ^0.5.1;
pragma experimental ABIEncoderV2;
import "./SafeMath.sol";


contract Autonity {
    using SafeMath for uint256;

    // validators - list of validators of network
    address[] public validators;
    // enodesWhitelist - which nodes can connect to network
    string[] public enodesWhitelist;
    // owner - owner of contract
    address public owner;
    // governanceOperatorAccount - account who can manipulate enodesWhitelist
    address public governanceOperatorAccount;

    /*
    * The bonding period (BP) is specified by the Autonity System Architecture as an integer representing an interval of blocks.
    * We have identifed two differents ways to how this parameter could be used :
    * 1. Bonding/unbonding operations happening at the end of each epoch.
    * 2. BP-Delayed unbonding.
    */
    uint256 public bonding_period = 100;
    /*
    * The commission rate is set globally at the member level and is public:
    * A member can’t have multiple commission rates depending on the member
    * The commission rate MUST be by default 0 and MUST remain unchanged if not updated.
    */
    mapping (address => uint256) private comission_rate;

    /*
    * mapping of members who are able to use stacking
    */
    mapping (address => bool) private members;

    /*
    * unbonded stake token balance
    */
    mapping (address => uint256) private stake_token;
    /*
    * bonded stake token balance
    */
    mapping (address => uint256) private bonded_stake_token;
    /*
    * delegated stake token balance.
    * map[owner address][delegator address] stake
    */
    mapping(address => mapping(address => uint)) delegated_stake_token;


    struct unbondingStake {
        uint256 amount;
        uint256 block_number;
    }
    mapping (address => mapping(address => unbondingStake[])) private unbonding_stake_token;



    /*
    * Ethereum transactions gas price must be greater or equal to the minimumGasPrice, a value set by the Governance operator.
    * FM-REQ-5: The minimumGasPrice value is a Genesis file configuration, if ommitted it defaults to 0.
    */
    uint256 minGasPrice = 0;



    // constructor get called at block #1 with msg.owner equal to Soma's deployer
    // configured in the genesis file.
    constructor (address[] memory _validators, string[] memory _enodesWhitelist,  address _governanceOperatorAccount) public {
        for (uint256 i = 0; i < _validators.length; i++) {
            validators.push(_validators[i]);
        }

        for (uint256 i = 0; i < _enodesWhitelist.length; i++) {
            enodesWhitelist.push(_enodesWhitelist[i]);
        }
        owner = msg.sender;
        governanceOperatorAccount = _governanceOperatorAccount;
    }


    /*
    * AddValidator
    * Add validator to validators list. Could be
    */
    function AddValidator(address _validator) public onlyValidators(msg.sender) {
        //Need to make sure we're duplicating the entry
        validators.push(_validator);
    }


    /*
    * RemoveValidator
    * Remove validator from validators list. function MUST be restricted to the Authority Account.
    */
    function RemoveValidator(address _validator) public onlyValidators(msg.sender) {
        require(validators.length > 1);

        for (uint256 i = 0; i < validators.length-1; i++) {
            if (validators[i] == _validator){
                validators[i] = validators[validators.length - 1];
                validators.length--;
                break;
            }
        }

    }

    /*
    * AddEnode
    * add enode to permission list
    * function MUST be restricted to the Authority Account.
    */
    function AddEnode(string memory  _enode) public onlyGovernanceOperator(msg.sender) {
        //Need to make sure we're not duplicating the entry
        enodesWhitelist.push(_enode);
    }

    /*
    * RemoveEnode
    * remove enode from permission list
    * function MUST be restricted to the Authority Account.
    */
    function RemoveEnode(string memory  _enode) public onlyGovernanceOperator(msg.sender) {
        require(enodesWhitelist.length > 1);

        for (uint256 i = 0; i < enodesWhitelist.length-1; i++) {
            if (compareStringsbyBytes(enodesWhitelist[i], _enode)) {
                enodesWhitelist[i] = enodesWhitelist[enodesWhitelist.length - 1];
                enodesWhitelist.length--;
                break;
            }
        }

    }

    /*
    * setMinimumGasPrice
    * FM-REQ-4: The Autonity Contract implements the setMinimumGasPrice function that is restricted to the Governance Operator account.
    * The function takes as an argument a positive integer and modify the value of minimumGasPrice
    */
    function SetMinimumGasPrice(uint256  _value) public onlyGovernanceOperator(msg.sender) {
        minGasPrice = _value;
    }


    /*
    * MintStake
    * function capable of creating new stake token and adding it to the recipient balance
    * function MUST be restricted to theAuthority Account.
    */
    function MintStake(address _account, uint256 _amount) public onlyGovernanceOperator(msg.sender) {
        require(members[_account] == true, "Account hasn't created");
        stake_token[_account] = stake_token[_account].add(_amount);
    }

    /*
    * RedeemStake
    * Decrease unbonded stake
    * The redeemStake(amount, recipient) function MUST be restricted to the Authority Account.
    */
    function RedeemStake(address _account, uint256 _amount) public onlyGovernanceOperator(msg.sender) {
        require(members[_account] == true, "Account hasn't created");
        stake_token[_account] =  stake_token[_account].sub(_amount, "Redeem stake amount exceeds balance");
    }


    /*
    * AddNewMember
    * Add not nil account to members list
    * function MUST be restricted to the Authority Account.
    */
    function AddNewMember(address _account) public onlyGovernanceOperator(msg.sender) {
        require(_account != address(0), "Account is empty");
        require(members[_account] == false, "Account has already created");
        members[_account] = true;
    }


    /*
    * RemoveMember
    * Remove account from members list
    * function MUST be restricted to the Authority Account.
    */
    function RemoveMember(address _account) public onlyGovernanceOperator(msg.sender) {
        require(members[_account] == true, "Account hasn't created");
        members[_account] = false;
    }



    /*
    * send
    * Moves `amount` stake tokens from the caller's account to `recipient`.
    *
    * Returns a boolean value indicating whether the operation succeeded.
    *
    * Emits a {Transfer} event.
    */
    function send(address _recipient, uint256 _amount) external returns (bool) {
        require(members[msg.sender] == true, "Account hasn't created");
        require(members[_recipient] == true, "Account hasn't created");
        _transfer(msg.sender, _recipient, _amount);
        return true;
    }


    // The Autonity Contract MUST implements the bondStake(amount, recipient) function capable of delegating stake token.
    function Bonding(address _recipient, uint256 amount) public returns (bool){
        require(members[msg.sender] == true, "Account hasn't created");
        require(members[_recipient] == true, "Account hasn't created");

        stake_token[msg.sender] = stake_token[msg.sender].sub(amount);
        bonded_stake_token[_recipient] = bonded_stake_token[_recipient].add(amount);
        delegated_stake_token[msg.sender][_recipient] = delegated_stake_token[msg.sender][_recipient].add(amount);
    }




    function Unbonding(address _recipient, uint256 _amount) public returns (bool){
        require(members[msg.sender] == true, "Account hasn't created");
        require(members[_recipient] == true, "Account hasn't created");

        bonded_stake_token[_recipient] = bonded_stake_token[_recipient].sub(_amount);
        delegated_stake_token[msg.sender][_recipient] = delegated_stake_token[msg.sender][_recipient].sub(_amount);
        unbonding_stake_token[msg.sender][_recipient].push(unbondingStake(_amount,  block.number + bonding_period));
    }


    //    The Autonity Contract MUST implements the setCommissionRate(rate)
    //    function capable of fixing the caller commission rate for the next bonding period.
    function SetCommissionRate(uint256 rate) public returns(bool)  {
        require(members[msg.sender] == true, "Account hasn't created");
        comission_rate[msg.sender] = rate;
        return true;
    }



    /*
    ========================================================================================================================

        Getters - extra values we may wish to return

    ========================================================================================================================
    */

    /*
    * getValidators
    *
    * Returns the macro validator list
    */

    function GetValidators() public view returns (address[] memory) {
        return validators;
    }

    /*
    * getWhitelist
    *
    * Returns the macro participants list
    */

    function GetWhitelist() public view returns (string[] memory) {
        return enodesWhitelist;
    }


    /*
    * GetAccountStake
    *
    * Returns unbonded stake for account
    */
    function GetAccountStake(address _account) public view returns (uint256) {
        return stake_token[_account];
    }


    /*
    * CheckMember
    *
    * Returns is addres a member
    */
    function CheckMember(address _account) public view returns (bool) {
        return members[_account];
    }

    /*
    * GetStake
    *
    * Returns sender's unbonded stake
    */
    function GetStake() public view returns(uint256) {
        return stake_token[msg.sender];
    }

    /*
    * GetBondedStake
    *
    * Returns sender's ubonded stake
    */
    function GetBondedStake() public view returns(uint256) {
        return bonded_stake_token[msg.sender];
    }

    /*
    * GetDelegatedBondedStake
    *
    * Returns sender's deleagated to _account
    */
    function GetDelegatedBondedStake(address _account) public view returns(uint256) {
        return delegated_stake_token[msg.sender][_account];
    }


    function getRate(address _account) public view returns(uint256) {
        require(members[msg.sender] == true, "Account hasn't created");
        return comission_rate[msg.sender];
    }

    /*
    * GetUnbondingStake
    *
    * Returns sender's unbonding stake by account
    */
    function GetUnbondingStake(address _account) public view returns(unbondingStake[] memory ) {
        return unbonding_stake_token[msg.sender][_account];
    }




    /*
    * GetMinimumGasPrice
    * Returns minimum gas price. Ethereum transactions gas price must be greater or equal to the minimumGasPrice.
    */
    function GetMinimumGasPrice() public view returns(uint256) {
        return minGasPrice;
    }


    /*
    ========================================================================================================================

        Modifiers

    ========================================================================================================================
    */

    /*
    * onlyValidators
    *
    * Modifier that checks if the voter is an active validator
    */

    modifier onlyValidators(address _voter) {
        bool present = false;
        for (uint256 i = 0; i < validators.length; i++) {
            if(validators[i] == _voter){
                present = true;
                break;
            }
        }
        require(present, "Voter is not a validator");
        _;
    }

    /*
    * onlyGovernanceOperator
    *
    * Modifier that checks if the caller is a Governance Operator
    */
    modifier onlyGovernanceOperator(address _caller) {
        require(governanceOperatorAccount == _caller, "Caller is not a operator");
        _;
    }



    /*
    ========================================================================================================================

        Events

    ========================================================================================================================
    */

    event Transfer(address indexed from, address indexed to, uint256 value);


    /*
    ========================================================================================================================

        Internal

    ========================================================================================================================
    */



    function _transfer(address sender, address recipient, uint256 amount) internal {
        require(sender != address(0), "Transfer from the zero address");
        require(recipient != address(0), "Transfer to the zero address");

        stake_token[sender] = stake_token[sender].sub(amount, "Transfer amount exceeds balance");
        stake_token[recipient] = stake_token[recipient].add(amount);
        emit Transfer(sender, recipient, amount);
    }


    function compareStringsbyBytes(string memory s1, string memory s2) internal pure returns(bool){
        return keccak256(abi.encodePacked(s1)) == keccak256(abi.encodePacked(s2));
    }

}