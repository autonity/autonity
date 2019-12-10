pragma solidity ^0.5.1;
pragma experimental ABIEncoderV2;
import "./SafeMath.sol";


contract Autonity {
    using SafeMath for uint256;

    struct ContractState {
        address[] users;
        string[] enodes;
        uint256[] types;
        uint256[] stakes;
        uint256[] commisionrates;
        address operator;
        address deployer;
        uint256 mingasprice;
        uint256 bondingperiod;
    }

    struct EconomicsMetricData {
        address[] accounts;
        UserType[] usertypes;
        uint256[] stakes;
        uint256[] commissionrates;
        uint256 mingasprice;
        uint256 stakesupply;
    }

    struct RewardDistributionData {
        bool result;
        address[] stakeholders;
        uint256[] rewardfractions;
        uint256 amount;
    }

    enum UserType { Participant, Stakeholder, Validator}

    struct User {
        address payable addr;
        UserType userType;
        uint256 stake;
        string enode;
        uint256 commission_rate; // rate must be by default 0 and must remain unchanged if not updated.
    }

    ////////////////////// Contract States need to be dumped for contract upgrade//////////
    address[] private usersList;
    string[] public enodesWhitelist;
    mapping (address => User) private users;
    address public deployer;
    address public operatorAccount;
    uint256 minGasPrice = 0;
    uint256 public bonding_period = 100;
    ///////////////////////////////////////////////////////////////////////////////////////

    ///////////////////// Contract state which can be replay from dumped states////////////
    address[] public validators;
    address[] private stakeholders;
    uint256 private stakeSupply;
    ///////////////////////////////////////////////////////////////////////////////////////

    /*
     * Binary code and ABI of new version contract, the default value is "" when contract is created.
     * if they are set by and only by operator, then contract upgrade will be triggered automatically.
    */
    string bytecode;
    string contractAbi;

    /*
    * Events
    *
    */
    event Transfer(address indexed from, address indexed to, uint256 value);
    event AddValidator(address _address, uint256 _stake);
    event AddStakeholder(address _address, uint256 _stake);
    event AddParticipant(address _address, uint256 _stake);
    event RemoveUser(address _address, UserType _type);
    event SetMinimumGasPrice(uint256 _gasPrice);
    event SetCommissionRate(address _address, uint256 _value);
    event MintStake(address _address, uint256 _amount);
    event RedeemStake(address _address, uint256 _amount);

    // constructor get called at block #1
    // configured in the genesis file.
    constructor (address[] memory _participantAddress,
        string[] memory _participantEnode,
        uint256[] memory _participantType,
        uint256[] memory _participantStake,
        uint256[] memory _commissionRate,
        address _operatorAccount,
        address _deployer,
        uint256 _minGasPrice,
        uint256 _bondingPeriod) public {


        require(_participantAddress.length == _participantEnode.length
        && _participantAddress.length == _participantType.length
        && _participantAddress.length == _participantStake.length,
            "Incorrect constructor params");


        for (uint256 i = 0; i < _participantAddress.length; i++) {
            require(_participantAddress[i] != address(0), "Addresses must be defined");
            UserType _userType = UserType(_participantType[i]);
            address payable addr = address(uint160(_participantAddress[i]));
            _createUser(addr, _participantEnode[i], _userType, _participantStake[i], _commissionRate[i]);
        }
        operatorAccount = _operatorAccount;
        deployer = msg.sender;
        minGasPrice = _minGasPrice;
        bonding_period = _bondingPeriod;
    }

    /*
    * addValidator
    * Add validator to validators list.
    */
    function addValidator(address payable _address, uint256 _stake, string memory _enode) public onlyOperator(msg.sender) {
        _createUser(_address,_enode, UserType.Validator, _stake, 0);
        emit AddValidator(_address, _stake);
    }

    function addStakeholder(address payable _address, string  memory _enode, uint256 _stake) public onlyOperator(msg.sender) {
        _createUser(_address, _enode, UserType.Stakeholder, _stake, 0);
        emit AddStakeholder(_address, _stake);
    }

    function addParticipant(address payable _address, string memory _enode) public onlyOperator(msg.sender) {
        _createUser(_address, _enode, UserType.Participant, 0, 0);
        emit AddParticipant(_address, 0);
    }

    /*
    * removeUser
    * remove user. function MUST be restricted to the Authority Account.
    */
    function removeUser(address _address) public onlyOperator(msg.sender) {
        require(_address != address(0), "address must be defined");
        require(users[_address].addr != address(0), "user must exists");
        User storage u = users[_address];

        if(u.userType == UserType.Validator || u.userType == UserType.Stakeholder){
            _removeFromArray(u.addr, stakeholders);
        }

        if(u.userType == UserType.Validator){
            _removeFromArray(u.addr, validators);
        }

        if (!(bytes(u.enode).length == 0)) {
            for (uint256 i = 0; i < enodesWhitelist.length; i++) {
                if (compareStringsbyBytes(enodesWhitelist[i], u.enode)) {
                    enodesWhitelist[i] = enodesWhitelist[enodesWhitelist.length - 1];
                    enodesWhitelist.length--;
                    break;
                }
            }
        }
        stakeSupply = stakeSupply.sub(u.stake);
        _removeFromArray(u.addr, usersList);
        delete users[_address];
        emit RemoveUser(_address, u.userType);
    }

    /*
    * setMinimumGasPrice
    * FM-REQ-4: The Autonity Contract implements the setMinimumGasPrice function that is restricted to the Governance Operator account.
    * The function takes as an argument a positive integer and modify the value of minimumGasPrice
    */
    function setMinimumGasPrice(uint256  _value) public onlyOperator(msg.sender) {
        minGasPrice = _value;
        emit SetMinimumGasPrice(_value);
    }

    /*
    * mintStake
    * function capable of creating new stake token and adding it to the recipient balance
    * function MUST be restricted to theAuthority Account.
    */
    function mintStake(address _account, uint256 _amount) public onlyOperator(msg.sender) canUseStake(_account) {
        users[_account].stake = users[_account].stake.add(_amount);
        stakeSupply = stakeSupply.add(_amount);
        emit MintStake(_account, _amount);
    }

    /*
    * redeemStake
    * Decrease unbonded stake
    * The redeemStake(amount, recipient) function MUST be restricted to the Authority Account.
    */
    function redeemStake(address _account, uint256 _amount) public onlyOperator(msg.sender) canUseStake(_account) {
        users[_account].stake = users[_account].stake.sub(_amount, "Redeem stake amount exceeds balance");
        stakeSupply = stakeSupply.sub(_amount);
        emit RedeemStake(_account, _amount);
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
        _transfer(msg.sender, _recipient, _amount);
        return true;
    }

    /*
     * TODO: msg.sender == operator address or anynode address, we might need node's address when operator perform this.
     * The Autonity Contract MUST implements the setCommissionRate(rate)
     * function capable of fixing the caller commission rate for the next bonding period.
     */
    function setCommissionRate(uint256 rate) public canUseStake(msg.sender) returns(bool) {
        users[msg.sender].commission_rate = rate;
        emit SetCommissionRate(msg.sender, rate);
        return true;
    }

    function upgradeContract(string memory _bytecode, string memory _abi) public onlyOperator(msg.sender) returns(bool) {
        bytecode = _bytecode;
        contractAbi = _abi;
        return true;
    }

    function retrieveContract() public view returns(string memory, string memory) {
        return (bytecode, contractAbi);
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

    function retrieveState() public view
    returns (address[] memory, string[] memory, uint256[] memory, uint256[] memory, uint256[] memory, address, uint256) {

        address[] memory addr = new address[](usersList.length);
        uint256[] memory userType  = new uint256[](usersList.length);
        uint256[] memory stake = new uint256[](usersList.length);
        string[] memory enode = new string[](usersList.length);
        uint256[] memory commissionRate = new uint256[](usersList.length);
        for(uint256 i=0; i<usersList.length; i++ ) {
            addr[i] = users[usersList[i]].addr;
            userType[i] = uint256(users[usersList[i]].userType);
            stake[i] = users[usersList[i]].stake;
            enode[i] = users[usersList[i]].enode;
            commissionRate[i] = users[usersList[i]].commission_rate;
        }
        return (addr, enode, userType, stake, commissionRate, operatorAccount, minGasPrice);
    }

    function retrieveStateV2() public view returns (ContractState memory state) {

        address[] memory addr = new address[](usersList.length);
        uint256[] memory userType  = new uint256[](usersList.length);
        uint256[] memory stake = new uint256[](usersList.length);
        string[] memory enode = new string[](usersList.length);
        uint256[] memory commissionRate = new uint256[](usersList.length);
        for(uint256 i=0; i<usersList.length; i++ ) {
            addr[i] = users[usersList[i]].addr;
            userType[i] = uint256(users[usersList[i]].userType);
            stake[i] = users[usersList[i]].stake;
            enode[i] = users[usersList[i]].enode;
            commissionRate[i] = users[usersList[i]].commission_rate;
        }
        ContractState memory cs = ContractState(addr, enode, userType, stake, commissionRate, operatorAccount, deployer, minGasPrice, bonding_period);
        return cs;
    }

    function getStakeholders() public view returns (address[] memory) {
        return stakeholders;
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
    * getAccountStake
    *
    * Returns unbonded stake for account
    */
    function getAccountStake(address _account) public view canUseStake(_account) returns (uint256) {
        return users[_account].stake;
    }


    /*
    * getStake
    *
    * Returns sender's unbonded stake
    */
    function getStake() public view canUseStake(msg.sender) returns(uint256)  {
        return users[msg.sender].stake;
    }

    function getRate(address _account) public view returns(uint256) {
        return users[_account].commission_rate;
    }


    /*
    * getMinimumGasPrice
    * Returns minimum gas price. Ethereum transactions gas price must be greater or equal to the minimumGasPrice.
    */
    function getMinimumGasPrice() public view returns(uint256) {
        return minGasPrice;
    }

    function checkMember(address _account) public view returns (bool) {
        return  users[_account].addr == _account;
    }

    /*
    * performRedistribution
    * return a structure contains reward distribution.
    */
    function performRedistribution(uint256 _amount) internal onlyDeployer(msg.sender) returns(RewardDistributionData memory rewarddistribution) {
        require(address(this).balance >= _amount, "not enough funds to perform redistribution");
        require(stakeholders.length > 0, "there must be stake holders");

        uint256[] memory rewardfractionlist = new uint256[](stakeholders.length);
        for (uint256 i = 0; i < stakeholders.length; i++) {
            User storage _user = users[stakeholders[i]];
            uint256 reward = _user.stake.mul(_amount).div(stakeSupply);
            _user.addr.transfer(reward);
            rewardfractionlist[i] = reward;
        }
        RewardDistributionData memory rd = RewardDistributionData(true, stakeholders, rewardfractionlist, _amount);
        return rd;
    }

    //Finalize function called once after every mined block, return if a new contract is ready for update
    function finalize(uint256 _amount) public onlyDeployer(msg.sender) returns (RewardDistributionData memory rewarddistribution) {
        RewardDistributionData memory data = performRedistribution(_amount);
        data.result = bytes(bytecode).length != 0;
        return data;
    }

    function totalSupply() public view returns (uint) {
        return stakeSupply;
    }

    /*
    * dumpEconomicsMetricData
    * Returns a struct which contains all the network economic data.
    */
    function dumpEconomicsMetricData() public view returns(EconomicsMetricData memory economics) {
        uint len = usersList.length;

        address[] memory tempAddrlist = new address[](len);
        UserType[] memory tempTypelist = new UserType[](len);
        uint256[] memory tempStakelist = new uint256[](len);
        uint256[] memory commissionRatelist = new uint256[](len);

        for (uint i = 0; i < len; i++) {
            tempAddrlist[i] = users[usersList[i]].addr;
            tempTypelist[i] = users[usersList[i]].userType;
            tempStakelist[i] = users[usersList[i]].stake;
            commissionRatelist[i] = users[usersList[i]].commission_rate;
        }

        EconomicsMetricData memory data = EconomicsMetricData(tempAddrlist, tempTypelist, tempStakelist, commissionRatelist, minGasPrice, stakeSupply);
        return data;
    }

    /*
    ========================================================================================================================

        Modifiers

    ========================================================================================================================
    */

    /*
    * onlyOperator
    *
    * Modifier that checks if the caller is a Governance Operator
    */
    modifier onlyOperator(address _caller) {
        require(operatorAccount == _caller, "Caller is not a operator");
        _;
    }

    modifier onlyDeployer(address _caller) {
        require(deployer == _caller, "Caller is not a operator");
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
        require(users[_address].addr != address(0), "address must be defined");
        _;
    }


    /*
    ========================================================================================================================

        Internal

    ========================================================================================================================
    */

    function _transfer(address sender, address recipient, uint256 amount) internal canUseStake(sender) canUseStake(recipient) {
        users[sender].stake = users[sender].stake.sub(amount, "Transfer amount exceeds balance");
        users[recipient].stake = users[recipient].stake.add(amount);
        emit Transfer(sender, recipient, amount);
    }


    function compareStringsbyBytes(string memory s1, string memory s2) internal pure returns(bool){
        return keccak256(abi.encodePacked(s1)) == keccak256(abi.encodePacked(s2));
    }

    function _createUser(address payable _address, string memory _enode, UserType _userType, uint256 _stake, uint256 commissionRate) internal {
        require(_address != address(0), "Addresses must be defined");
        User memory u = User(_address, _userType, _stake, _enode, commissionRate);

        // avoid duplicated user in usersList.
        if (users[u.addr].addr != u.addr) {
            usersList.push(u.addr);
        }

        users[u.addr] = u;

        if (u.userType == UserType.Stakeholder){
            stakeholders.push(u.addr);
        } else if(u.userType == UserType.Validator){
            validators.push(u.addr);
            stakeholders.push(u.addr);
        }
        stakeSupply = stakeSupply.add(_stake);

        if(bytes(u.enode).length != 0) {
            enodesWhitelist.push(u.enode);
        }
    }


    function _removeFromArray(address _address, address[] storage _array) internal {
        require(_array.length > 0);

        for (uint256 i = 0; i < _array.length; i++) {
            if (_array[i] == _address) {
                _array[i] = _array[_array.length - 1];
                _array.length--;
                break;
            }
        }
    }

    // @notice Will receive any eth sent to the contract
    function () external payable {
    }
}

