pragma solidity ^0.5.11;
pragma experimental ABIEncoderV2;

import "./SafeMath.sol";


contract Autonity {
    using SafeMath for uint256;

    // validators - list of validators of network
    address[] public validators;
    //array of members who are able to use stacking
    address[] private stakeholders;
    // enodesWhitelist - which enodes can connect to network
    string[] public enodesWhitelist;
    // deployer - deployer's address, controlled by the native client.
    address public deployer;
    // operatorAccount - account who can access restricted functions.
    address public operatorAccount;

    enum UserType {Participant, Stakeholder, Validator}

    struct User {
        address payable addr;
        UserType userType;
        uint256 stake;
        string enode;
    }

    mapping(address => User) private users;

    uint256 private stakeSupply;


    /*
    * Ethereum transactions gas price must be greater or equal to the minimumGasPrice, a value set by the Governance operator.
    * FM-REQ-5: The minimumGasPrice value is a Genesis file configuration, if omitted it defaults to 0.
    */
    uint256 minGasPrice = 0;



    // constructor get called at block #1 with msg.owner equal to Soma's deployer
    // configured in the genesis file.
    constructor (address[] memory _userAddresses,
        string[] memory _userEnodes,
        uint256[] memory _userTypes,
        uint256[] memory _userStakes,
        address _operatorAccount,
        uint256 _minGasPrice) public {

        uint256 numUsers = _userAddresses.length;
        require(numUsers == _userEnodes.length &&
                numUsers == _userTypes.length &&
                numUsers == _userStakes.length,
            "bad user arrays length");


        for (uint256 i = 0; i < numUsers; i++) {
            UserType _userType = UserType(_userTypes[i]);
            createUser(_userAddresses[i], _userType, _userStakes[i], _userEnodes[i]);

        }

        deployer = msg.sender;
        operatorAccount = _operatorAccount;
        minGasPrice = _minGasPrice;
    }

    /*
    * addStakeholder
    * Add not nil account to members list
    * function MUST be restricted to the Authority Account.
    */
    function addStakeholder(address _address, string  memory _enode, uint256 _stake) public onlyOperator(msg.sender) {
        createUser(_address, UserType.Stakeholder, _stake, _enode);
    }

    /*
    * AddValidator
    * Add validator to validators list. Could be
    */
    function addValidator(address _address, string  memory _enode, uint256 _stake) public onlyOperator(msg.sender) {
        createUser(_address, UserType.Validator, _stake, _enode);
    }

    /*
    * AddEnode
    * add enode to permission list
    * function MUST be restricted to the Authority Account.
    */
    function addParticipant(address _address, string memory _enode) public onlyOperator(msg.sender) {
        createUser(_address, UserType.Participant, 0, _enode);
    }

    /*
    * removeUser
    * remove user from contract
    * function MUST be restricted to the Authority Account.
    */
    function removeUser(address _address) public onlyOperator(msg.sender) {
        require(_address != address(0), "address must be defined");
        require(users[_address].addr != address(0), "user must exists");
        User storage u = users[_address];
        if(u.userType == UserType.Validator || u.userType == UserType.Stakeholder){
            removeFromArray(u.addr, stakeholders);
        }
        if(u.userType == UserType.Validator){
            removeFromArray(u.addr, validators);
        }
        if (! (bytes(u.enode).length == 0)) {
            for (uint256 i = 0; i < enodesWhitelist.length - 1; i++) {
                if (compareStringsbyBytes(enodesWhitelist[i], u.enode)) {
                    enodesWhitelist[i] = enodesWhitelist[enodesWhitelist.length - 1];
                    enodesWhitelist.length--;
                    break;
                }
            }
        }
    }

    /*
    * setMinimumGasPrice
    * FM-REQ-4: The Autonity Contract implements the setMinimumGasPrice function that is restricted to the Governance Operator account.
    * The function takes as an argument a positive integer and modify the value of minimumGasPrice
    */
    function SetMinimumGasPrice(uint256 _value) public onlyOperator(msg.sender) {
        minGasPrice = _value;
    }


    /*
    * MintStake
    * function capable of creating new stake token and adding it to the recipient balance
    * function MUST be restricted to theAuthority Account.
    */
    function mintStake(address _account, uint256 _amount) public onlyOperator(msg.sender) canUseStake(_account) {
        users[_account].stake = users[_account].stake.add(_amount);
        stakeSupply = stakeSupply.add(_amount);
    }

    /*
    * RedeemStake
    * Decrease unbonded stake
    * The redeemStake(amount, recipient) function MUST be restricted to the Authority Account.
    */
    function redeemStake(address _account, uint256 _amount) public onlyOperator(msg.sender) canUseStake(_account){
        users[_account].stake = users[_account].stake.sub(_amount, "account has insufficient stake");
        stakeSupply = stakeSupply.sub(_amount);
    }


    function totalSupply() public view returns (uint) {
        return stakeSupply;
    }

    /*
    * performRedistribution
    * redistribute fee token prorata stake
    * called by the native client as part of the block finalization logic.
    */
    function performRedistribution(uint256 _amount) public onlyDeployer(msg.sender) {

        require(address(this).balance >= _amount, "not enough funds to perform redistribution");
        for (uint256 i = 0; i < stakeholders.length; i++) {
            User storage _user = users[stakeholders[i]];
            uint256 _fees = _user.stake.mul(_amount).div(stakeSupply);
            _user.addr.transfer(_fees);
        }
    }

    /*
    * send
    * Moves `amount` stake tokens from the caller's account to `recipient`.
    *
    * Returns a boolean value indicating whether the operation succeeded.
    *
    * Emits a {Transfer} event.
    */
    function send(address _recipient, uint256 _amount) external canUseStake(msg.sender) canUseStake(_recipient) returns (bool) {
        transfer(msg.sender, _recipient, _amount);
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

    function getValidators() public view returns (address[] memory) {
        return validators;
    }

    /*
    * getWhitelist
    *
    * Returns the macro participants list
    */

    function getWhitelist() public view returns (string[] memory) {
        return enodesWhitelist;
    }


    /*
    * GetAccountStake
    *
    * Returns unbonded stake for account
    */
    function getAccountStake(address _account) public view canUseStake(_account) returns (uint256) {
        return users[_account].stake;
    }


    /*
    * CheckMember
    *
    * Returns is addres a member
    */
    function checkMember(address _account) public view returns (bool) {
        return  users[_account].userType == UserType.Stakeholder ||
        users[_account].userType ==  UserType.Validator ;
    }


    /*
    * getMinimumGasPrice
    * Returns minimum gas price. Ethereum transactions gas price must be greater or equal to the minimumGasPrice.
    */
    function getMinimumGasPrice() public view returns (uint256) {
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
            if (validators[i] == _voter) {
                present = true;
                break;
            }
        }
        require(present, "caller is not a validator");
        _;
    }

    /*
    * onlyOperator
    *
    * Modifier that checks if the caller is a Governance Operator
    */
    modifier onlyOperator(address _caller) {
        require(operatorAccount == _caller, "caller is not a operator");
        _;
    }


    /*
   * onlyOperator
   *
   * Modifier that checks if the caller is the native client.
   */
    modifier onlyDeployer(address _caller) {
        require(deployer == _caller, "caller is not native client");
        _;
    }

    /*
   * canUseStake
   *
   * Modifier that checks if the adress can use stake.
   */
    modifier canUseStake(address _address) {
        require(_address != address(0), "address must be defined");
        require(users[_address].userType == UserType.Stakeholder ||
            users[_address].userType ==  UserType.Validator, "address not allowed to use stake");
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



    function transfer(address _sender, address _recipient, uint256 _amount) internal {
        users[_sender].stake = users[_sender].stake.sub(_amount, "Transfer amount exceeds balance");
        users[_recipient].stake = users[_recipient].stake.add(_amount);
        emit Transfer(_sender, _recipient, _amount);
    }


    function compareStringsbyBytes(string memory s1, string memory s2) internal pure returns (bool){
        return keccak256(abi.encodePacked(s1)) == keccak256(abi.encodePacked(s2));
    }

    function createUser(address _address, UserType _type, uint256 _stake, string memory _enode) internal returns (User storage) {
        require(_address != address(0), "user address must be defined");
        require(users[_address].addr == address(0), "user already existing");
        User memory u = User({addr : address(uint160(_address)), // casting to payable address
            enode : _enode,
            userType : _type,
            stake : _stake});

        users[u.addr] = u;

        if (u.userType == UserType.Stakeholder) {
            stakeholders.push(u.addr);
        } else if (u.userType == UserType.Validator) {
            stakeholders.push(u.addr);
            validators.push(u.addr);
        } else if (u.userType == UserType.Participant) {
            require(u.stake == 0, "stake of participant must be 0");
        }

        stakeSupply = stakeSupply.add(u.stake);

        if (!(bytes(u.enode).length == 0)) {
            enodesWhitelist.push(u.enode);
        }

        return users[u.addr];
    }

    function removeFromArray(address _address, address[] storage _array) internal {
        require(_array.length > 0);

        for (uint256 i = 0; i < _array.length - 1; i++) {
            if (_array[i] == _address) {
                _array[i] = _array[_array.length - 1];
                _array.length--;
                break;
            }
        }
    }
}