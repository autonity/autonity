// SPDX-License-Identifier: MIT

pragma solidity ^0.7.1;
pragma experimental ABIEncoderV2;
import "./interfaces/IERC20.sol";
import "./SafeMath.sol";
import "./Precompiled.sol";
import "./Accountability.sol";

/** @title Proof-of-Stake Autonity Contract */
contract Autonity is IERC20 {
    using SafeMath for uint256;

    enum UserType { Participant, Stakeholder, Validator}
    struct User {
        address payable addr;
        UserType userType;
        uint256 stake;
        string enode;
    }

    struct CommitteeMember {
        address payable addr;
        uint256 votingPower;
    }

    struct EconomicMetrics {
        address[] accounts;
        UserType[] usertypes;
        uint256[] stakes;
        uint256 mingasprice;
        uint256 stakesupply;
    }

    /* State data that needs to be dumped in-case of a contract upgrade. */
    Accountability.Proof[] public challenges;
    address[] private usersList;
    string[] private enodesWhitelist;
    mapping (address => User) private users;
    address public operatorAccount;
    uint256 private minGasPrice;
    uint256 public committeeSize;
    string private contractVersion;

    mapping (address => mapping (address => uint256)) private allowances;

    /* State data that will be recomputed during a contract upgrade. */
    address[] private validators;
    address[] private stakeholders;
    uint256 private stakeSupply;
    CommitteeMember[] private committee;

    /*
    We're saving the address of who is deploying the contract and we use it
    for restricting functions that could only be possibly invoked by the protocol
    itself, bypassing transaction processing and signature verification.
    In normal conditions, it is set to the zero address. We're not simply hardcoding
    it only because of testing purposes.
    */
    address public deployer;

    /*
     Binary code and ABI of a new contract, the default value is "" when the contract is deployed.
     If the bytecode is not empty then a contract upgrade will be triggered automatically.
    */
    string bytecode;
    string contractAbi;

    /* Events */
    event UserAdded(address _address, UserType _type, uint256 _stake);
    event RemovedUser(address _address, UserType _type);
    event ChangedUserType(address _address, UserType _oldType, UserType _newType);
    event MintedStake(address _address, uint256 _amount);
    event BurnedStake(address _address, uint256 _amount);
    event Rewarded(address _address, uint256 _amount);
    event ChallengeAdded(Accountability.Proof proof);
    event ChallengeRemoved(Accountability.Proof proof);

    /**
     * @dev Emitted when the Minimum Gas Price was updated and set to `gasPrice`.
     * Note that `gasPrice` may be zero.
     */
    event MinimumGasPriceUpdated(uint256 gasPrice);

    /**
     * @dev Emitted when the Autonity Contract was upgraded to a new version (`version`).
     */
    event ContractUpgraded(string version);

    constructor (address[] memory _participantAddress,
        string[] memory _participantEnode,
        uint256[] memory _participantType,
        uint256[] memory _participantStake,
        address _operatorAccount,
        uint256 _minGasPrice,
        uint256 _committeeSize,
        string memory _contractVersion) {

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
        operatorAccount = _operatorAccount;
        minGasPrice = _minGasPrice;
        contractVersion = _contractVersion;
        committeeSize = _committeeSize;
        deployer = msg.sender;
    }

    /**
    * @dev Receive Auton function https://solidity.readthedocs.io/en/v0.7.2/contracts.html#receive-ether-function
    *
    */
    receive() external payable {}

    /**
    * @dev Fallback function https://solidity.readthedocs.io/en/v0.7.2/contracts.html#fallback-function
    *
    */
    fallback() external payable {}


    /**
    * @return the name of the stake token.
    * @dev ERC-20 Optional.
    */
    function name() external pure returns (string memory) {
        return "Newton";
    }

    /**
    * @return the Stake token's symbol.
    * @dev ERC-20 Optional.
    */
    function symbol() external pure returns (string memory) {
        return "NEW";
    }

    /**
    * @notice Create an accountability challenge in the Autonity Contract with the specified role. Restricted to the validator account.
    */
    function addChallenge(uint256 h, uint64 r, address sender, uint8 rule, uint8 msgType, bytes memory packedProof) public onlyProtocol(msg.sender) {
        Accountability.Proof memory challenge = Accountability.Proof(h, sender, r, rule, msgType, packedProof);
        require(_isChallengeExists(challenge) == false, "Duplicated Challenge");
        require(Accountability.checkChallenge(packedProof)[0] != 0, "Not a valid challenge");

        challenges.push(challenge);
        emit ChallengeAdded(challenge);
    }

    /**
    * @notice Resolve an accountability challenge in the Autonity Contract with the specified role. Restricted to the validator account.
    */
    function resolveChallenge(uint256 h, uint64 r, address sender, uint8 rule, uint8 msgType, bytes memory packedProof) public onlyProtocol(msg.sender) {
        Accountability.Proof memory proof = Accountability.Proof(h, sender, r, rule, msgType, packedProof);
        require(_isChallengeExists(proof) == true, "Not visible challenge to be resolved");
        require(Accountability.checkInnocent(packedProof)[0] != 0, "Not a valid proof of innocent");

        _removeChallenge(proof);
        emit ChallengeRemoved(proof);
    }

    /**
    * @notice Create a user in the Autonity Contract with the specified role. Restricted to the operator account.
    */
    function addUser(address payable _address, uint256 _stake, string memory _enode, UserType _role) public onlyOperator(msg.sender) {
        require(!(_role == UserType.Participant && _stake > 0), "participant can't have stake");
        _createUser(_address, _enode, _role, _stake);
        emit UserAdded(_address, _role, _stake);
    }

    /**
    * @notice Change the user account type. Restricted to the operator account.
    */
    function changeUserType(address _address , UserType newUserType ) public onlyOperator(msg.sender) {
        _changeUserType(_address, newUserType);
    }

    /**
    * @notice Remove the user account from the contract. The account stake is burnt and the
    * associated enode(if present) is removed from the enode whitelist. Restricted to the operator account.
    * @param account address to be removed.
    * @dev emit a {RemoveUser} event.
    */
    function removeUser(address account) public onlyOperator(msg.sender) {
        _removeUser(account);
    }

    /**
    * @notice Set the minimum gas price. Restricted to the operator account.
    * @param price Positive integer.
    * @dev Emit a {MinimumGasPriceUpdated} event.
    */
    function setMinimumGasPrice(uint256 price) public onlyOperator(msg.sender) {
        minGasPrice = price;
        emit MinimumGasPriceUpdated(price);
    }

    /*
    * @notice Set the maximum size of the consensus committee. Restricted to the Operator account.
    *
    */
    function setCommitteeSize(uint256 size) public onlyOperator(msg.sender) {
        committeeSize = size;
    }

    /*
    * @notice Mint new stake token (NEW) and add it to the recipient balance. Restricted to the Operator account.
    * @dev emit a MintStake event.
    */
    function mint(address _account, uint256 _amount) public onlyOperator(msg.sender) canUseStake(_account) {
        users[_account].stake = users[_account].stake.add(_amount);
        stakeSupply = stakeSupply.add(_amount);
        emit MintedStake(_account, _amount);
    }

    /**
    * @notice Burn the specified amount of NEW stake token from an account. Restricted to the Operator account.
    */
    function burn(address _account, uint256 _amount) public onlyOperator(msg.sender) canUseStake(_account) {
        users[_account].stake = users[_account].stake.sub(_amount, "Redeem stake amount exceeds balance");
        stakeSupply = stakeSupply.sub(_amount);
        _checkDowngradeValidator(_account);
        emit BurnedStake(_account, _amount);
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
        _approve(sender, msg.sender, allowances[sender][msg.sender].sub(amount, "ERC20: transfer amount exceeds allowance"));
        return true;
    }

    /**
    * @dev See {IERC20-allowance}.
    */
    function allowance(address owner, address spender) external view override returns (uint256) {
        return allowances[owner][spender];
    }

    function upgradeContract(string memory _bytecode,
                             string memory _abi,
                             string memory _version) public onlyOperator(msg.sender) returns(bool) {
        bytecode = _bytecode;
        contractAbi = _abi;
        contractVersion = _version;
        emit ContractUpgraded(contractVersion);
        return true;
    }

    /**
    * @notice Getter to retrieve a new Autonity contract bytecode and ABI when an upgrade is initiated.
    * @return `bytecode` the new contract bytecode.
    * @return `contractAbi` the new contract ABI.
    */
    function getNewContract() external view returns(string memory, string memory) {
        return (bytecode, contractAbi);
    }


    /** @dev finalize is the block state finalisation function. It is called
    * each block after processing every transactions within it. It must be restricted to the
    * protocol only.
    *
    * @param amount The amount of transaction fees collected for this block.
    * @return upgrade Set to true if an autonity contract upgrade is available.
    * @return committee The next block consensus committee.
    */
    function finalize(uint256 amount) external onlyProtocol(msg.sender)
        returns(bool , CommitteeMember[] memory) {

        _performRedistribution(amount);
        bool _updateAvailable = bytes(bytecode).length != 0;
        computeCommittee();
        return (_updateAvailable, committee);
    }

    /**
    * @dev Dump the current internal state key elements. Called by the protocol during a contract upgrade.
    * The returned data will be passed directly to the constructor of the new contract at deployment.
    */
    function getState() external view returns(
        address[] memory _addr,
        string[] memory _enode,
        uint256[] memory _userType,
        uint256[] memory _stake,
        address _operatorAccount,
        uint256 _minGasPrice,
        uint256 _committeeSize,
        string memory _contractVersion) {

        // Exceptionally using named returns here, make things clearer.
        _addr = new address[](usersList.length);
        _userType  = new uint256[](usersList.length);
        _stake = new uint256[](usersList.length);
        _enode = new string[](usersList.length);
        for(uint256 i=0; i<usersList.length; i++ ) {
            _addr[i] = users[usersList[i]].addr;
            _userType[i] = uint256(users[usersList[i]].userType);
            _stake[i] = users[usersList[i]].stake;
            _enode[i] = users[usersList[i]].enode;
        }
        _operatorAccount = operatorAccount;
        _minGasPrice = minGasPrice;
        _committeeSize = committeeSize;
        _contractVersion = contractVersion;
    }

    /*
    ============================================================
        Getters
    ============================================================
    */

    /**
    * @notice Returns the current contract version.
    */
    function getVersion() external view returns (string memory) {
        return contractVersion;
    }

    /**
     * @notice Returns the block committee.
     * @dev Current block committee if called before finalize(), next block if called after.
     */
    function getCommittee() external view returns (CommitteeMember[] memory) {
        return committee;
    }

    /**
     * @notice Returns the current list of validators.
     */
    function getValidators() external view returns (address[] memory) {
        return validators;
    }

    /**
     * @notice Returns the current list of stakeholders.
     */
    function getStakeholders() external view returns (address[] memory) {
        return stakeholders;
    }

    /**
    * @notice Autonity Protocol function, returns the list of authorized enodes
    * able to join the network.
    */
    function getWhitelist() external view returns (string[] memory) {
        return enodesWhitelist;
    }

    /**
    * @notice Returns the amount of stake token held by the account (ERC-20).
    */
    function balanceOf(address _account) external view override returns (uint256) {
        return users[_account].stake;
    }

    /**
    * @notice Returns the total amount of stake token issued.
    */
    function totalSupply() external view override returns (uint256) {
        return stakeSupply;
    }

    /**
    * @return Returns a user object with the `_account` parameter. The returned data
    * object might be empty if there is no user associated.
    */
    function getUser(address _account) external view returns(User memory) {
        //TODO : coreturn an error if no user was found.
        return users[_account];
    }

    /**
    * @return Returns the maximum size of the consensus committee.
    */
    function getMaxCommitteeSize() external view returns(uint256) {
        return committeeSize;
    }

    /**
    * @return Returns the minimum gas price.
    * @dev Autonity transaction's gas price must be greater or equal to the minimum gas price.
    */
    function getMinimumGasPrice() external view returns(uint256) {
        return minGasPrice;
    }

    /**
    * @notice getProposer returns the address of the proposer for the given height and
    * round. The proposer is selected from the committee via weighted random
    * sampling, with selection probability determined by the voting power of
    * each committee member. The selection mechanism is deterministic and will
    * always select the same address, given the same height, round and contract
    * state.
    */
    function getProposer(uint256 height, uint256 round) external view returns(address) {
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

    /**
    * @notice Returns an object which contains all the network economics data.
    */
    function dumpEconomicMetrics() public view returns(EconomicMetrics memory) {
        uint len = usersList.length;

        address[] memory tempAddrlist = new address[](len);
        UserType[] memory tempTypelist = new UserType[](len);
        uint256[] memory tempStakelist = new uint256[](len);

        for (uint i = 0; i < len; i++) {
            tempAddrlist[i] = users[usersList[i]].addr;
            tempTypelist[i] = users[usersList[i]].userType;
            tempStakelist[i] = users[usersList[i]].stake;
        }

        EconomicMetrics memory data = EconomicMetrics(tempAddrlist, tempTypelist, tempStakelist, minGasPrice, stakeSupply);
        return data;
    }

    /**
    * @notice update the current committee by selecting top staking validators.
    * Restricted to the protocol.
    */
    function computeCommittee() public onlyProtocol(msg.sender) {
        // Left public for testing purposes.
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
    ============================================================

        Modifiers

    ============================================================
    */

    /**
    * @dev Modifier that checks if the caller is the governance operator account.
    * This should be abstracted by a separate smart-contract.
    */
    modifier onlyOperator(address _caller) {
        require(operatorAccount == _caller, "caller is not the operator");
        _;
    }

    /**
    * @dev Modifier that checks if the caller is not any external owned account.
    * Only the protocol itself can invoke the contract with the 0 address to the exception
    * of testing.
    */
    modifier onlyProtocol(address _caller) {
        require(deployer == _caller, "function restricted to the protocol");
        _;
    }

    /**
    * @dev Modifier that checks if the adress is authorized to own stake.
    */
    modifier canUseStake(address _address) {
        require(_address != address(0), "address must be defined");
        require(users[_address].userType == UserType.Stakeholder ||
        users[_address].userType ==  UserType.Validator, "address not allowed to use stake");
        require(users[_address].addr != address(0), "address must be defined");
        _;
    }


    /*
    ============================================================

        Internals

    ============================================================
    */

    /**
    * @notice Perform Auton reward distribution. The transaction fees
    * are simply re-distributed to all stake-holders, including validators,
    * pro-rata the amount of stake held.
    * @dev Emit a {BlockReward} event for every account that collected rewards.
    */
    function _performRedistribution(uint256 _amount) internal  {
        require(address(this).balance >= _amount, "not enough funds to perform redistribution");
        require(stakeholders.length > 0, "there must be stake holders");

        for (uint256 i = 0; i < stakeholders.length; i++) {
            User storage _user = users[stakeholders[i]];
            uint256 _reward = _user.stake.mul(_amount).div(stakeSupply);
            _user.addr.transfer(_reward);
            emit Rewarded(_user.addr, _reward);
        }
    }

    function _transfer(address sender, address recipient, uint256 amount) internal canUseStake(sender) canUseStake(recipient) {
        users[sender].stake = users[sender].stake.sub(amount, "Transfer amount exceeds balance");
        users[recipient].stake = users[recipient].stake.add(amount);
        _checkDowngradeValidator(sender);
        emit Transfer(sender, recipient, amount);
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
    function _approve(address owner, address spender, uint256 amount) internal canUseStake(spender) virtual {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");

        allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    /**
    * @dev If the user is a validator and its stake is zero downgrade it to a
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

    function _changeUserType(address _address , UserType newUserType ) internal {
        require(_address != address(0), "address must be defined");
        require(users[_address].addr != address(0), "user must exist");

        require(users[_address].userType != newUserType, "The user is already of this type.");

        // Removes the user and adds it again with the new userType
        User memory u = users[_address];
        if(newUserType == UserType.Participant){
            require(u.stake == 0);
        }
        _removeUser(u.addr);
        _createUser(u.addr, u.enode, newUserType, u.stake);

        emit ChangedUserType(u.addr , u.userType , newUserType);
    }

    function _isChallengeExists(Accountability.Proof memory proof) internal view returns (bool) {
        for (uint256 i = 0; i < challenges.length; i++) {
            if (challenges[i].rule == proof.rule && challenges[i].height == proof.height
            && challenges[i].round == proof.round && challenges[i].msgType == proof.msgType
            && challenges[i].sender == proof.sender) {

                return true;
            }
        }
        return false;
    }

    function _removeChallenge(Accountability.Proof memory proof) internal {
        require(challenges.length > 0);

        for (uint256 i = 0; i < challenges.length; i++) {
            if (challenges[i].rule == proof.rule && challenges[i].height == proof.height
                && challenges[i].round == proof.round && challenges[i].msgType == proof.msgType
                && challenges[i].sender == proof.sender) {

                challenges[i] = challenges[challenges.length - 1];
                challenges.pop();
                break;
            }
        }
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
                if (_compareStringsbyBytes(enodesWhitelist[i], u.enode)) {
                    enodesWhitelist[i] = enodesWhitelist[enodesWhitelist.length - 1];
                    enodesWhitelist.pop();
                    break;
                }
            }
        }
        stakeSupply = stakeSupply.sub(u.stake);
        _removeFromArray(u.addr, usersList);
        delete users[_address];
        emit RemovedUser(_address, u.userType);
    }


    function _createUser(address payable _address, string memory _enode, UserType _userType, uint256 _stake) internal {
        require(_address != address(0), "Addresses must be defined");
        require(Precompiled.enodeCheck(_enode)[0] != 0, "enode error");

        User memory u = User(_address, _userType, _stake, _enode);

        // avoid duplicated user in usersList.
        require(users[u.addr].addr == address(0), "already registered address");

        usersList.push(u.addr);

        users[u.addr] = u;

        if (u.userType == UserType.Stakeholder){
            stakeholders.push(u.addr);
        } else if(u.userType == UserType.Validator){
            require(u.stake != 0, "validator with 0 stake is not permitted");
            validators.push(u.addr);
            stakeholders.push(u.addr);
        }
        stakeSupply = stakeSupply.add(_stake);

        if(bytes(u.enode).length != 0) {
            enodesWhitelist.push(u.enode);
        }
    }



    /**
    * @dev Order validators by stake
    */
    function _sortByStake(User[] memory _validators) internal pure {
        _structQuickSort(_validators, int(0), int(_validators.length - 1));
    }

    /**
    * @dev QuickSort algorithm sorting in ascending order by stake.
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

    function _compareStringsbyBytes(string memory s1, string memory s2) internal pure returns(bool){
        return keccak256(abi.encodePacked(s1)) == keccak256(abi.encodePacked(s2));
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
}
