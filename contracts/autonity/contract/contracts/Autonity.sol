pragma solidity ^0.6.4;
pragma experimental ABIEncoderV2;
import "./SafeMath.sol";
import "./Precompiled.sol";

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
        //RewardDistributionData
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
        // uint256 selfStake;
        // uint256 delegatedStake;
        uint256 commissionRate; // rate must be by default 0 and must remain unchanged if not updated.
    }

    struct CommitteeMember {
        address payable addr;
        uint256 votingPower;
    }

    ////////////////////// Contract States need to be dumped for contract upgrade//////////
    address[] private usersList;
    string[] public enodesWhitelist;
    mapping (address => User) private users;
    address public deployer;
    address public operatorAccount;
    uint256 minGasPrice = 0;
    uint256 public bondingPeriod = 10*60;
    uint256 public committeeSize = 20;
    string public contractVersion = "v0.0.0";
    ///////////////////////////////////////////////////////////////////////////////////////

    ///////////////////// Contract state which can be replay from dumped states////////////
    // Below 4 meta are used in each block generation before or after, so it is more about performance consideration.
    // They are re-playable by constructor function.
    address[] private validators;
    address[] private stakeholders;
    uint256 private stakeSupply;
    CommitteeMember[] private committee;
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
    event ChangeUserType(address _address, UserType _oldType, UserType _newType);
    event SetMinimumGasPrice(uint256 _gasPrice);
    event SetCommissionRate(address _address, uint256 _value);
    event MintStake(address _address, uint256 _amount);
    event RedeemStake(address _address, uint256 _amount);
    event BlockReward(address _address, uint256 _amount);
    event Version(string version);
    // constructor get called at block #1
    // configured in the genesis file.

    constructor (address[] memory _participantAddress,
        string[] memory _participantEnode,
        uint256[] memory _participantType,
        uint256[] memory _participantStake,
        uint256[] memory _commissionRate,
        address _operatorAccount,
        uint256 _minGasPrice,
        uint256 _bondingPeriod,
        uint256 _committeeSize,
        string memory _contractVersion) public {

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
        bondingPeriod = _bondingPeriod;
        contractVersion = _contractVersion;
        committeeSize = _committeeSize;
    }

    /*
    * addUser
    * Add user to autonity contract.
    */
    function addUser(address payable _address, uint256 _stake, string memory _enode, UserType _role) public onlyOperator(msg.sender) {
        if (_role == UserType.Validator) {
            _createUser(_address, _enode, _role, _stake, 0);
            emit AddValidator(_address, _stake);
        }

        if (_role == UserType.Stakeholder) {
            _createUser(_address, _enode, _role, _stake, 0);
            emit AddStakeholder(_address, _stake);
        }

        if (_role == UserType.Participant) {
            _createUser(_address, _enode, _role, 0, 0);
            emit AddParticipant(_address, 0);
        }
    }

    /*
    * changeUserType
    * Change user status
    */
    function changeUserType( address _address , UserType newUserType ) public onlyOperator(msg.sender) {
        _changeUserType(_address, newUserType);
    }

    function _changeUserType( address _address , UserType newUserType ) internal {
        require(_address != address(0), "address must be defined");
        require(users[_address].addr != address(0), "user must exist");

        require(users[_address].userType != newUserType, "The user is already of this type.");

        // Removes the user and adds it again with the new userType
        User memory u = users[_address];
        if(newUserType == UserType.Participant){
            require(u.stake == 0);
        }
        _removeUser(u.addr);
        _createUser(u.addr, u.enode, newUserType, u.stake, u.commissionRate);

        emit ChangeUserType(u.addr , u.userType , newUserType);
    }

    /*
    * removeUser
    * remove user. function MUST be restricted to the Authority Account.
    */
    function removeUser(address _address) public onlyOperator(msg.sender) {
        _removeUser(_address);
    }

    function _removeUser(address _address) internal {
        require(_address != address(0), "address must be defined");
        require(users[_address].addr != address(0), "user must exists");
        User storage u = users[_address];

        if(u.userType == UserType.Validator || u.userType == UserType.Stakeholder){
            _removeFromArray(u.addr, stakeholders);
        }

        if(u.userType == UserType.Validator){
            require(validators.length > 1, "There must be at least 1 validator in the network");
            _removeFromArray(u.addr, validators);
        }

        if (!(bytes(u.enode).length == 0)) {
            for (uint256 i = 0; i < enodesWhitelist.length; i++) {
                if (compareStringsbyBytes(enodesWhitelist[i], u.enode)) {
                    enodesWhitelist[i] = enodesWhitelist[enodesWhitelist.length - 1];
                    enodesWhitelist.pop();
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
        _checkDowngradeValidator(_account);
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
    function setCommissionRate(uint256 _rate) public canUseStake(msg.sender) returns(bool) {
        users[msg.sender].commissionRate = _rate;
        emit SetCommissionRate(msg.sender, _rate);
        return true;
    }

    function upgradeContract(string memory _bytecode, string memory _abi, string memory _version) public onlyOperator(msg.sender) returns(bool) {
        bytecode = _bytecode;
        contractAbi = _abi;
        contractVersion = _version;
        emit Version(contractVersion);
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

    function getVersion() public view returns (string memory) {
        return contractVersion;
    }

    function retrieveState() public view
    returns (address[] memory, string[] memory, uint256[] memory, uint256[] memory, uint256[] memory, address, uint256, uint256, uint256, string memory) {

        address[] memory _addr = new address[](usersList.length);
        uint256[] memory _userType  = new uint256[](usersList.length);
        uint256[] memory _stake = new uint256[](usersList.length);
        string[] memory _enode = new string[](usersList.length);
        uint256[] memory _commissionRate = new uint256[](usersList.length);
        for(uint256 i=0; i<usersList.length; i++ ) {
            _addr[i] = users[usersList[i]].addr;
            _userType[i] = uint256(users[usersList[i]].userType);
            _stake[i] = users[usersList[i]].stake;
            _enode[i] = users[usersList[i]].enode;
            _commissionRate[i] = users[usersList[i]].commissionRate;
        }
        return (_addr, _enode, _userType, _stake, _commissionRate, operatorAccount, minGasPrice, bondingPeriod, committeeSize, contractVersion);
    }

    /*
    * getStakeholders
    *
    * Returns the macro stakeholders list
    */
    function getStakeholders() public view returns (address[] memory) {
        return stakeholders;
    }

    function getCommittee() public view returns (CommitteeMember[] memory) {
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
        return users[_account].commissionRate;
    }

    /*
    * myUserType
    * Returns sender's userType
    */
    function myUserType() public view returns(UserType)  {
        return users[msg.sender].userType;
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
    * computeCommittee
    * update the current committee by selecting top staking validators
    */
    function computeCommittee() public onlyDeployer(msg.sender) {
        require(validators.length > 0, "There must be validators");

        uint _len = validators.length;
        uint256 _committeeLength = committeeSize;
        if (_committeeLength >= _len) {_committeeLength = _len;}

        User[] memory _validatorList = new User[](_len);
        User[] memory _committeeList = new User[](_committeeLength);

        for (uint256 i = 0;i < validators.length; i++) {
            User memory _user = users[validators[i]];
            _validatorList[i] =_user;
        }

        // If there are more validators than seats in the committee
        if (_validatorList.length > committeeSize) {
            // sort validators by stake in ascending order
           _sortByStake(_validatorList);
            // choose the top-N (with N=maxCommitteeSize)
            for (uint256 _j = 0; _j < committeeSize; _j++) {
                _committeeList[_j] = _validatorList[_j];
            }
        }
        // If all the validators fit in the committee
        else {
            _committeeList = _validatorList;
        }

        // Update committee in persistent storage
        delete committee;
        for (uint256 _k =0 ; _k < _committeeLength; _k++) {
            CommitteeMember memory _member = CommitteeMember(_committeeList[_k].addr, _committeeList[_k].stake);
            committee.push(_member);
        }

    }

    /*
    * getProposer
    * getProposer returns the address of the proposer for the given height and
    * round. The proposer is selected from the committee via weighted random
    * sampling, with selection probability determined by the voting power of
    * each committee member. The selection mechanism is deterministic and will
    * always select the same address, given the same height, round and contract
    * state.
    */
    function getProposer(uint256 height, uint256 round) public view returns(address) {
        // calculate total voting power from current committee, the system does not allow validator with 0 stake/power.
        uint256 total_voting_power = 0;
        for (uint256 i = 0; i < committee.length; i++) {
            total_voting_power = total_voting_power.add(committee[i].votingPower);
        }

        require(total_voting_power != 0, "The committee is not staking");

        // distribute seed into a 256bits key-space.
        uint256 key = height.add(round);
        uint256 value = uint256(keccak256(abi.encodePacked(key)));
        uint256 index = value % total_voting_power;

        // find the index hit which committee member which line up in the committee list.
        // we assume there is no 0 stake/power validators.
        uint256 counter = 0;
        for (uint256 i = 0; i < committee.length; i++) {
            counter = counter.add(committee[i].votingPower);
            if (index <= counter - 1) {
                return committee[i].addr;
            }
        }
        revert("There is no validator left in the network");
    }

    /*
    * performRedistribution
    * return a structure contains reward distribution.
    */
    function performRedistribution(uint256 _amount) internal onlyDeployer(msg.sender)  {
        require(address(this).balance >= _amount, "not enough funds to perform redistribution");
        require(stakeholders.length > 0, "there must be stake holders");

        for (uint256 i = 0; i < stakeholders.length; i++) {
            User storage _user = users[stakeholders[i]];
            uint256 _reward = _user.stake.mul(_amount).div(stakeSupply);
            _user.addr.transfer(_reward);
            emit BlockReward(_user.addr, _reward);
        }
    }

    /*
    * finalize
    * function called once after every mined block. To avoid calling the evm multiple times we return
    * here if there is an update available and the next block committee.
    */
    function finalize(uint256 _amount) public onlyDeployer(msg.sender) returns(bool , CommitteeMember[] memory) {
        performRedistribution(_amount);
        bool _updateAvailable = bytes(bytecode).length != 0;
        computeCommittee();

       return (_updateAvailable, committee);
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
            commissionRatelist[i] = users[usersList[i]].commissionRate;
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
        require(operatorAccount == _caller, "Caller is not the operator");
        _;
    }

    modifier onlyDeployer(address _caller) {
        require(deployer == _caller, "Caller is not deployer");
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
        _checkDowngradeValidator(sender);
        emit Transfer(sender, recipient, amount);
    }

    /*
    * If the user is a validator and its stake is zero downgrade it to a
    * stakeholder. Unless the user is the only validator in which case the
    * transaction is reverted.
    *
    * This check is used to ensure that we do not end up in a situation where
    * the system has no validators.
    */
    function _checkDowngradeValidator(address _address) internal {
        User memory u = users[_address];
        if (u.stake != 0 || u.userType != UserType.Validator) {
            return;
        }

        require(validators.length > 1, "Downgrade user failed due to keep at least 1 validator in the network");
        _changeUserType(u.addr, UserType.Stakeholder);
    }

    function compareStringsbyBytes(string memory s1, string memory s2) internal pure returns(bool){
        return keccak256(abi.encodePacked(s1)) == keccak256(abi.encodePacked(s2));
    }

    function _createUser(address payable _address, string memory _enode, UserType _userType, uint256 _stake, uint256 commissionRate) internal {
        require(_address != address(0), "Addresses must be defined");
        require(Precompiled.enodeCheck(_enode)[0] != 0, "enode error");

        User memory u = User(_address, _userType, _stake, _enode, commissionRate);

        // avoid duplicated user in usersList.
        require(users[u.addr].addr == address(0), "This address is already registered");

        usersList.push(u.addr);

        users[u.addr] = u;

        if (u.userType == UserType.Stakeholder){
            stakeholders.push(u.addr);
        } else if(u.userType == UserType.Validator){
            require(u.stake != 0, "Validator with 0 stake is not permitted");
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
                _array.pop();
                break;
            }
        }
    }

    /*
    * sortByStake
    * Order validators by stake
    *
    */
    function _sortByStake(User[] memory _validators) internal pure{
        _structQuickSort(_validators, int(0), int(_validators.length - 1));
    }

    /*
    * structQuickSort
    * QuickSort algorithm sorting in ascending order by stake
    */
    function _structQuickSort(User[] memory _users, int _low, int _high) internal pure {

        int _i = _low;
        int _j = _high;
        if (_i==_j) return;
        uint _pivot = _users[uint(_low + (_high - _low) / 2)].stake;
        // Set the pivot element in its right sorted index in the array
        while (_i <= _j) {
            while (_users[uint(_i)].stake > _pivot) _i++;
            while (_pivot > _users[uint(_j)].stake) _j--;
            if (_i <= _j) {
                (_users[uint(_i)], _users[uint(_j)]) = (_users[uint(_j)], _users[uint(_i)]);
                _i++;
                _j--;
            }
        }
        // Recursion call in the left partition of the array
        if (_low < _j) {
            _structQuickSort(_users, _low, _j);
        }
        // Recursion call in the right partition
        if (_i < _high) {
            _structQuickSort(_users, _i, _high);
        }
    }

    receive() external payable {}

    fallback() external payable {}
}
