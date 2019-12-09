pragma solidity ^0.5.1;
pragma experimental ABIEncoderV2;
import "./SafeMath.sol";


contract Autonity {
    using SafeMath for uint256;

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

    address[] private usersList;

    // validators - list of validators of network
    address[] public validators;
    // enodesWhitelist - which nodes can connect to network
    string[] public enodesWhitelist;
    // deployer - deployer of contract
    address public deployer;
    // operatorAccount - account who can manipulate enodesWhitelist
    address public operatorAccount;

    uint256 private stakeSupply;

    uint256 public committeeSize = 1000;

    /*
    * The bonding period (BP) is specified by the Autonity System Architecture as an integer representing an interval of blocks.
    * We have identifed two differents ways to how this parameter could be used :
    * 1. Bonding/unbonding operations happening at the end of each epoch.
    * 2. BP-Delayed unbonding.
    */
    uint256 public bondingPeriod = 100;
    /*
    * The commission rate is set globally at the member level and is public:
    * A member can’t have multiple commission rates depending on the member
    * The commission rate MUST be by default 0 and MUST remain unchanged if not updated.
    */
    mapping (address => uint256) public commission_rate;

    //array of members who are able to use stacking
    address[] private stakeholders;

    enum UserType { Participant, Stakeholder, Validator}

    struct User {
        address payable addr;
        UserType userType;
        uint256 stake;
        string enode;
        // uint256 selfStake;
        // uint256 delegatedStake;
    }

    mapping (address => User) private users;

    User[] public committee;

    uint256 totalStake = 0;
    /*
    * Ethereum transactions gas price must be greater or equal to the minimumGasPrice, a value set by the Governance operator.
    * FM-REQ-5: The minimumGasPrice value is a Genesis file configuration, if ommitted it defaults to 0.
    */
    uint256 minGasPrice = 0;

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
        address _operatorAccount,
        uint256 _minGasPrice) public {


        require(_participantAddress.length == _participantEnode.length
        && _participantAddress.length == _participantType.length
        && _participantAddress.length == _participantStake.length,
            "Incorrect constructor params");


        for (uint256 i = 0; i < _participantAddress.length; i++) {
            require(_participantAddress[i] != address(0), "Addresses must be defined");
            UserType _userType = UserType(_participantType[i]);
            address payable addr = address(uint160(_participantAddress[i]));
            _createUser(addr, _participantEnode[i], _userType, _participantStake[i]);
        }
        deployer = msg.sender;
        operatorAccount = _operatorAccount;
        minGasPrice = _minGasPrice;
    }


    /*
    * addValidator
    * Add validator to validators list.
    */
    function addValidator(address payable _address, uint256 _stake, string memory _enode) public onlyOperator(msg.sender) {
        _createUser(_address,_enode, UserType.Validator, _stake);
        emit AddValidator(_address, _stake);
    }

    function addStakeholder(address payable _address, string  memory _enode, uint256 _stake) public onlyOperator(msg.sender) {
        _createUser(_address, _enode, UserType.Stakeholder, _stake);
        emit AddStakeholder(_address, _stake);
    }

    function addParticipant(address payable _address, string memory _enode) public onlyOperator(msg.sender) {
        _createUser(_address, _enode, UserType.Participant, 0);
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
    * setCommitteeSize
    * Set the maximum size of the commitee, restricted to the Governance Operator account
    */

    function setCommitteeSize (uint256 _size) public onlyOperator(msg.sender) {
        committeeSize = _size;
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



    //    The Autonity Contract MUST implements the setCommissionRate(rate)
    //    function capable of fixing the caller commission rate for the next bonding period.
    function setCommissionRate(uint256 rate) public canUseStake(msg.sender) returns(bool)   {
        commission_rate[msg.sender] = rate;
        emit SetCommissionRate(msg.sender, rate);
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
    * getStakeholders
    *
    * Returns the macro stakeholders list
    */

    function getStakeholders() public view returns (address[] memory) {
        return stakeholders;
    }

    function getCommittee() public view returns (User[] memory) {
        return committee;
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
        return commission_rate[_account];
    }

    /*
    * getMaxCommitteeSize
    * Returns the maximum possible size of the committee - set of validators participating in consensus.
    */
    function getMaxCommitteeSize() public view returns(uint256) {
        return committeeSize;
    }

    /*
    *getCurrentCommitteeSize
    *Returns the size of the current committee
    */
    function getCurrentCommiteeSize() public view returns(uint256) {
        return committee.length;
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
    * sortByStake
    * Order validators by stake
    *
    */
    function sortByStake(User[] memory _validators) internal pure returns(User[] memory){
        structQuickSort(_validators, int(0), int(_validators.length - 1));
        return _validators;
    }

    /*
    * structQuickSort
    * QuickSort algorithm sorting in ascending order by stake
    */
    function structQuickSort(User[] memory _users, int low, int high) internal pure {

        int i = low;
        int j = high;
        if (i==j) return;
        uint pivot = _users[uint(low + (high - low) / 2)].stake;
        // Set the pivot element in its right sorted index in the array
        while (i <= j) {
            while (_users[uint(i)].stake > pivot) i++;
            while (pivot > _users[uint(j)].stake) j--;
            if (i <= j) {
                (_users[uint(i)], _users[uint(j)]) = (_users[uint(j)], _users[uint(i)]);
                i++;
                j--;
            }
        }
        // Recursion call in the left partition of the array
        if (low < j) {
            structQuickSort(_users, low, j);
        }
        // Recursion call in the right partition
        if (i < high) {
            structQuickSort(_users, i , high);
        }
    }

    /*
    * setCommittee
    * selects the committee of validators to participate in consensus
    */
    function setCommittee() public onlyDeployer(msg.sender) returns(User[] memory){
        require(validators.length > 0, "There must be validators");

        uint len = validators.length;
        uint256 committeeLength = committeeSize;
        if (committeeLength >= len) {committeeLength = len;}

        User[] memory validatorList = new User[](len);
        User[] memory sortedValidatorList = new User[](len);
        User[] memory committeeList = new User[](committeeLength);

        for (uint256 i = 0;i < validators.length; i++) {
            User memory _user = users[validators[i]];
            validatorList[i] =_user;
        }

        // If there are more validators than seats in the committee
        if (validatorList.length > committeeSize) {
            // sort validators by stake in ascending order
            sortedValidatorList = sortByStake(validatorList);
            // choose the top-N (with N=maxCommitteeSize)
            for (uint256 j = 0; j < committeeSize; j++) {
                committeeList[j] = sortedValidatorList[j];
            }
        }
        // If all the validators fit in the committee
        else {
            committeeList = validatorList;
        }

        // Update committee in persistent storage
        delete committee;
        for (uint256 k =0 ; k < committeeLength; k++) {
            committee.push(committeeList[k]);
        }

        return committeeList;
    }

    /*
    * performRedistribution
    * return a structure contains reward distribution.
    */
    function performRedistribution(uint256 _amount) public onlyDeployer(msg.sender) returns(RewardDistributionData memory rewarddistribution) {
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
            commissionRatelist[i] = commission_rate[usersList[i]];
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

    function _createUser(address payable _address, string memory _enode, UserType _userType, uint256 _stake) internal {
        require(_address != address(0), "Addresses must be defined");
        User memory u = User(_address, _userType, _stake, _enode);
        users[u.addr] = u;
        usersList.push(u.addr);

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

