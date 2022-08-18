// SPDX-License-Identifier: MIT

pragma solidity ^0.8.3;

import "./interfaces/IERC20.sol";
import "./Liquid.sol";
import "./Upgradeable.sol";
import "./Precompiled.sol";
import "./Helpers.sol";

/** @title Proof-of-Stake Autonity Contract */
contract Autonity is IERC20, Upgradeable {

    uint256 public constant COMMISSION_RATE_PRECISION = 10_000;

    enum ValidatorState {active, paused}
    struct Validator {
        address payable treasury;
        address addr;
        string enode; //addr must match provided enode
        uint256 commissionRate;
        uint256 bondedStake;
        uint256 totalSlashed;
        Liquid liquidContract;
        uint256 liquidSupply;
        uint256 registrationBlock;
        ValidatorState state;
    }

    struct CommitteeMember {
        address addr;
        uint256 votingPower;
    }

    /* Used for epoched staking */
    struct Staking {
        address payable delegator;
        address delegatee;
        uint256 amount;
        uint256 startBlock;
    }

    /* Used to track commission rate change - See ADR-002 */
    struct CommissionRateChangeRequest {
        address validator;
        uint256 startBlock;
        uint256 rate;
    }
    // Todo: Create a FIFO structure library, integrate with Staking{}
    mapping(uint256 => CommissionRateChangeRequest) internal commissionRateChangeQueue;
    uint256 internal commissionRateChangeQueueFirst = 0;
    uint256 internal commissionRateChangeQueueLast = 0;

    /**************************************************/

    struct Config {
        address operatorAccount;
        address payable  treasuryAccount;
        uint256 treasuryFee;
        uint256 minBaseFee;
        uint256 delegationRate;
        uint256 epochPeriod;
        uint256 unbondingPeriod;
        uint256 committeeSize;
        uint256 contractVersion;
        uint256 blockPeriod;
    }

    Config public config;
    address[] internal validatorList;


    // Stake token state transitions happen every epoch.
    uint256 public epochID;
    uint256 public lastEpochBlock;
    uint256 public epochTotalBondedStake;

    CommitteeMember[] internal committee;
    uint256 public totalRedistributed;
    uint256 public epochReward;
    string[] internal committeeNodes;
    mapping(address => mapping(address => uint256)) internal allowances;

    /*
    Keep track of bonding and unbonding requests.
    */
    mapping(uint256 => Staking) internal bondingMap;
    uint256 public tailBondingID;
    uint256 public headBondingID;
    mapping(uint256 => Staking) internal unbondingMap;
    uint256 public tailUnbondingID;
    uint256 public headUnbondingID;

    /* Newton ERC-20. */
    mapping(address => uint256) internal accounts;
    mapping(address => Validator) internal validators;
    uint256 internal stakeSupply;

    /*
    We're saving the address of who is deploying the contract and we use it
    for restricting functions that could only be possibly invoked by the protocol
    itself, bypassing transaction processing and signature verification.
    In normal conditions, it is set to the zero address. We're not simply hardcoding
    it only because of testing purposes.
    */
    address public deployer;

    /* Events */
    event MintedStake(address addr, uint256 amount);
    event BurnedStake(address addr, uint256 amount);
    event CommissionRateChange(address validator, uint256 rate);
    event RegisteredValidator(address treasury, address addr, string enode, address liquidContract);
    event PausedValidator(address treasury, address addr, uint256 effectiveBlock);
    event Rewarded(address addr, uint256 amount);

    /**
     * @dev Emitted when the Minimum Gas Price was updated and set to `gasPrice`.
     * Note that `gasPrice` may be zero.
     */
    event MinimumBaseFeeUpdated(uint256 gasPrice);

    constructor(Validator[] memory _validators,
                    Config memory _config) {
        if (config.contractVersion == 0) {
            deployer = msg.sender;
            _initialize(_validators, _config);
        }
    }

    function _initialize(Validator[] memory _validators,
    Config memory _config
    ) internal {
        config = _config;

        /* We are sharing the same Validator data structure for both genesis
           initialization and runtime. It's not an ideal solution but
           it avoids us adding more complexity to the contract and running into
           stack limit issues.
        */
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _bondedStake = _validators[i].bondedStake;

            // Sanitize the validator fields for a fresh new deployment.
            _validators[i].liquidSupply = 0;
            _validators[i].liquidContract = Liquid(address(0));
            _validators[i].bondedStake = 0;
            _validators[i].registrationBlock = 0;
            _validators[i].commissionRate = config.delegationRate;
            _validators[i].state = ValidatorState.active;
            // Autonity doesn't handle voting power over 64 bits.
            require(_bondedStake < 2 ** 60, "issued Newton can't be greater than 2^60");

            _registerValidator(_validators[i]);

            accounts[_validators[i].treasury] += _bondedStake;
            stakeSupply += _bondedStake;
            epochTotalBondedStake += _bondedStake;
            _bond(_validators[i].addr, _bondedStake, payable(_validators[i].treasury));
        }
        _stakingTransitions();
        computeCommittee();
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
        return "NTN";
    }

    /**
    * @notice Register a new validator in the system.  The validator might be selected to be part of consensus.
    * This validator will have assigned to its treasury account the caller of this function.
    * A new token "Liquid Stake" is deployed at this phase.
    * @param _enode enode identifying the validator node.
    * @param _proof proof containing treasury account and signed by validator's node key.
    * @dev Emit a {RegisteredValidator} event.
    */
    function registerValidator(string memory _enode, bytes memory _proof) public {
        Validator memory _val = Validator(payable(msg.sender), //treasury
            address(0), // address
            _enode, // enode
            config.delegationRate, // validator commission rate
            0, // bonded stake
            0, // total slashed
            Liquid(address(0)), // liquid token contract
            0, // liquid token supply
            block.number,
            ValidatorState.active
        );

        _registerValidator(_val, _proof);
        emit RegisteredValidator(msg.sender, _val.addr, _enode, address(_val.liquidContract));
    }

    function bond(address _validator, uint256 _amount) public {
        require(validators[_validator].addr == _validator, "validator not registered");
        require(validators[_validator].state == ValidatorState.active, "validator need to be active");
        _bond(_validator, _amount, payable(msg.sender));
    }

    function unbond(address _validator, uint256 _amount) public {
        require(validators[_validator].addr == _validator, "validator not registered");
        _unbond(_validator, _amount, payable(msg.sender));
    }

    /**
    * @notice Pause the validator and stop it accepting delegations. See ADR-004 for more details.
    * @param _address address to be disabled.
    * @dev emit a {DisabledValidator} event.
    */
    function pauseValidator(address _address) public {
        require(validators[_address].addr == _address, "validator must be registered");
        require(validators[_address].treasury == msg.sender, "require caller to be validator admin account");
        _pauseValidator(_address);
    }

    /**
    * @notice Re-activate the specified validator. See ADR-004 for more details.
    * @param _address address to be enabled.
    */
    function activateValidator(address _address) public {
        require(validators[_address].addr == _address, "validator must be registered");
        require(validators[_address].treasury == msg.sender, "require caller to be validator admin account");
        require(validators[_address].state == ValidatorState.paused, "validator must be paused");

        validators[_address].state = ValidatorState.active;
    }

    /**
    * @notice Change commission rate for the specified validator. See ADR-002 for more details.
    * @param _validator address to be enabled.
            _rate new commission rate, ranging between 0-10000 (10000 = 100%).
    */
    function changeCommissionRate(address _validator, uint256 _rate) public {
        require(validators[_validator].addr == _validator, "validator must be registered");
        require(validators[_validator].treasury == msg.sender, "require caller to be validator admin account");
        require(_rate <= COMMISSION_RATE_PRECISION, "require correct commission rate");
        CommissionRateChangeRequest memory _newRequest = CommissionRateChangeRequest(_validator, block.number, _rate);
        commissionRateChangeQueue[commissionRateChangeQueueLast] = _newRequest;
        commissionRateChangeQueueLast += 1;
        emit CommissionRateChange(_validator, _rate);
    }

    /**
    * @notice Set the minimum gas price. Restricted to the operator account.
    * @param _price Positive integer.
    * @dev Emit a {MinimumBaseFeeUpdated} event.
    */
    function setMinimumBaseFee(uint256 _price) public onlyOperator {
        config.minBaseFee = _price;
        emit MinimumBaseFeeUpdated(_price);
    }

    /*
    * @notice Set the maximum size of the consensus committee. Restricted to the Operator account.
    * @param _size Positive integer.
    */
    function setCommitteeSize(uint256 _size) public onlyOperator {
        require(_size > 0, "committee size can't be 0");
        config.committeeSize = _size;
    }

    /*
    * @notice Set the unbonding period. Restricted to the Operator account.
    * @param _size Positive integer.
    */
    function setUnbondingPeriod(uint256 _period) public onlyOperator {
        config.unbondingPeriod = _period;
    }

    /*
    * @notice Set the epoch period. Restricted to the Operator account.
    * @param _period Positive integer.
    */
    function setEpochPeriod(uint256 _period) public onlyOperator {
        config.epochPeriod = _period;
    }

    /*
    * @notice Set the Operator account. Restricted to the Operator account.
    * @param _account the new operator account.
    */
    function setOperatorAccount(address _account) public onlyOperator {
        config.operatorAccount = _account;
    }

    /*
    Currently not supported
    * @notice Set the block period. Restricted to the Operator account.
    * @param _period Positive integer.

    function setBlockPeriod(uint256 _period) public onlyOperator {
        config.blockPeriod = _period;
    }
     */

    /*
    * @notice Set the global treasury account. Restricted to the Operator account.
    * @param _account New treasury account.
    */
    function setTreasuryAccount(address payable _account) public onlyOperator {
        config.treasuryAccount = _account;
    }

    /*
    * @notice Set the treasury fee. Restricted to the Operator account.
    * @param _treasuryFee Treasury fee. Precision TBD.
    */
    function setTreasuryFee(uint256 _treasuryFee) public onlyOperator {
        config.treasuryFee = _treasuryFee;
    }

    /*
    * @notice Mint new stake token (NTN) and add it to the recipient balance. Restricted to the Operator account.
    * @dev emit a MintStake event.
    */
    function mint(address _addr, uint256 _amount) public onlyOperator {
        // Autonity doesn't handle voting power over 64 bits.
        require(_amount < 2 ** 60, "issued Newton can't be greater than 2^60");
        accounts[_addr] += _amount;
        stakeSupply += _amount;
        emit MintedStake(_addr, _amount);
    }

    /**
    * @notice Burn the specified amount of NTN stake token from an account. Restricted to the Operator account.
    * This won't burn associated Liquid tokens.
    */
    function burn(address _addr, uint256 _amount) public onlyOperator {
        require(accounts[_addr] >= _amount, "Amount exceeds balance");
        accounts[_addr] -= _amount;
        stakeSupply -= _amount;
        emit BurnedStake(_addr, _amount);
    }

    /**
    * @notice Moves `amount` NTN stake tokens from the caller's account to `recipient`.
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
        uint256 newAllowance = allowances[sender][msg.sender] - amount;
        _approve(sender, msg.sender, newAllowance);
        return true;
    }

    /**
    * @dev See {IERC20-allowance}.
    */
    function allowance(address owner, address spender) external view override returns (uint256) {
        return allowances[owner][spender];
    }

    /** @dev finalize is the block state finalisation function. It is called
    * each block after processing every transactions within it. It must be restricted to the
    * protocol only.
    *
    * @param amount The amount of transaction fees collected for this block.
    * @return upgrade Set to true if an autonity contract upgrade is available.
    * @return committee The next block consensus committee.
    */
    function finalize(uint256 amount) external virtual onlyProtocol
    returns (bool, CommitteeMember[] memory) {
        epochReward += amount;
        if (lastEpochBlock + config.epochPeriod == block.number) {
            // - slashing should come here first -
            _performRedistribution(epochReward);
            epochReward = 0;
            _stakingTransitions();
            _applyNewCommissionRates();
            computeCommittee();
            lastEpochBlock = block.number;
            epochID += 1;
        }
        return (contractUpgradeReady, committee);
    }

    /**
    * @notice update the current committee by selecting top staking validators.
    * Restricted to the protocol.
    */
    function computeCommittee() public onlyProtocol {
        // Left public for testing purposes.
        require(validatorList.length > 0, "There must be validators");
        /*
         As opposed to storage arrays, it is not possible to resize memory arrays
         have to calculate the required size in advance
        */
        uint _len = 0;
        for (uint256 i = 0; i < validatorList.length; i++) {
            if (validators[validatorList[i]].state == ValidatorState.active &&
                validators[validatorList[i]].bondedStake > 0) {
                _len++;
            }
        }

        uint256 _committeeLength = config.committeeSize;
        if (_committeeLength >= _len) {_committeeLength = _len;}

        Validator[] memory _validatorList = new Validator[](_len);
        Validator[] memory _committeeList = new Validator[](_committeeLength);

        // since Push function does not apply to fix length array, introduce a index j to prevent the overflow,
        // not all the members in validator pool satisfy the enabled && bondedStake > 0, so the overflow happens.
        uint j = 0;
        for (uint256 i = 0; i < validatorList.length; i++) {
            if (validators[validatorList[i]].state == ValidatorState.active &&
                validators[validatorList[i]].bondedStake > 0) {
                // Perform a copy of the validator object
                Validator memory _user = validators[validatorList[i]];
                _validatorList[j] = _user;
                j++;
            }
        }

        // If there are more validators than seats in the committee
        if (_validatorList.length > config.committeeSize) {
            // sort validators by stake in ascending order
            _sortByStake(_validatorList);
            // choose the top-N (with N=maxCommitteeSize)
            // Todo: (optimisation) just pop()
            for (uint256 _j = 0; _j < config.committeeSize; _j++) {
                _committeeList[_j] = _validatorList[_j];
            }
        }
        // If all the validators fit in the committee
        else {
            _committeeList = _validatorList;
        }

        // Update committee in persistent storage
        delete committee;
        delete committeeNodes;
        epochTotalBondedStake = 0;
        for (uint256 _k = 0; _k < _committeeLength; _k++) {
            CommitteeMember memory _member = CommitteeMember(_committeeList[_k].addr, _committeeList[_k].bondedStake);
            committee.push(_member);
            committeeNodes.push(_committeeList[_k].enode);
            epochTotalBondedStake += _committeeList[_k].bondedStake;
        }
    }


    /*
    ============================================================
        Getters
    ============================================================
    */

    /**
    * @notice Returns the last epoch's end block height.
    */
    function getLastEpochBlock() external view returns (uint256) {
        return lastEpochBlock;
    }

    /**
    * @notice Returns the current contract version.
    */
    function getVersion() external view returns (uint256) {
        return config.contractVersion;
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
        return validatorList;
    }


    /**
    * @notice Returns the amount of unbonded Newton token held by the account (ERC-20).
    */
    function balanceOf(address _addr) external view override returns (uint256) {
        return accounts[_addr];
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
    function getValidator(address _addr) external view returns (Validator memory) {
        //TODO : coreturn an error if no user was found.
        return validators[_addr];
    }

    /**
    * @return Returns the maximum size of the consensus committee.
    */
    function getMaxCommitteeSize() external view returns (uint256) {
        return config.committeeSize;
    }

    /**
    * @return Returns the consensus committee enodes.
    */
    function getCommitteeEnodes() external view returns (string[] memory) {
        return committeeNodes;
    }

    /**
    * @return Returns the minimum gas price.
    * @dev Autonity transaction's gas price must be greater or equal to the minimum gas price.
    */
    function getMinimumBaseFee() external view returns (uint256) {
        return config.minBaseFee;
    }

    /**
     * @notice Returns the current operator account.
    */
    function getOperator() external view returns (address) {
        return config.operatorAccount;
    }


    /**
    * @notice getProposer returns the address of the proposer for the given height and
    * round. The proposer is selected from the committee via weighted random
    * sampling, with selection probability determined by the voting power of
    * each committee member. The selection mechanism is deterministic and will
    * always select the same address, given the same height, round and contract
    * state.
    */
    function getProposer(uint256 height, uint256 round) external view returns (address) {
        // calculate total voting power from current committee, the system does not allow validator with 0 stake/power.
        uint256 total_voting_power = 0;
        for (uint256 i = 0; i < committee.length; i++) {
            total_voting_power += committee[i].votingPower;
        }

        require(total_voting_power != 0, "The committee is not staking");

        // distribute seed into a 256bits key-space.
        uint256 key = height + round;
        uint256 value = uint256(keccak256(abi.encodePacked(key)));
        uint256 index = value % total_voting_power;

        // find the index hit which committee member which line up in the committee list.
        // we assume there is no 0 stake/power validators.
        uint256 counter = 0;
        for (uint256 i = 0; i < committee.length; i++) {
            counter += committee[i].votingPower;
            if (index <= counter - 1) {
                return committee[i].addr;
            }
        }
        revert("There is no validator left in the network");
    }

    // lastId not included
    function getBondingReq(uint256 startId, uint256 lastId) external view returns (Staking[] memory) {
        // the total length of bonding sets.
        Staking[] memory _results = new Staking[](lastId - startId);
        for (uint256 i = 0; i < lastId - startId; i++) {
            _results[i] = bondingMap[startId + i];
        }
        return _results;
    }

    function getUnbondingReq(uint256 startId, uint256 lastId) external view returns (Staking[] memory) {
        Staking[] memory _results = new Staking[](lastId - startId);
        // the total length of bonding sets.
        for (uint256 i = 0; i < lastId - startId; i++) {
            _results[i] = unbondingMap[startId + i];
        }
        return _results;
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
    modifier onlyOperator override {
        require(config.operatorAccount == msg.sender, "caller is not the operator");
        _;
    }

    /**
    * @dev Modifier that checks if the caller is not any external owned account.
    * Only the protocol itself can invoke the contract with the 0 address to the exception
    * of testing.
    */
    modifier onlyProtocol {
        require(deployer == msg.sender, "function restricted to the protocol");
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
    function _performRedistribution(uint256 _amount) internal virtual {
        require(address(this).balance >= _amount, "not enough funds to perform redistribution");
        // take treasury fee.

        uint256 _treasuryReward = (config.treasuryFee * _amount) / 10 ** 18;
        if (_treasuryReward > 0) {
            config.treasuryAccount.transfer(_treasuryReward);
            _amount -= _treasuryReward;
        }
        totalRedistributed += _amount;
        for (uint256 i = 0; i < committee.length; i++) {
            Validator storage _val = validators[committee[i].addr];
            uint256 _reward = (_val.bondedStake * _amount) / epochTotalBondedStake;
            if (_reward > 0) {
                _val.liquidContract.redistribute{value : _reward}();
            }
            emit Rewarded(_val.addr, _reward);
        }
    }

    function _transfer(address _sender, address _recipient, uint256 _amount) internal virtual {
        require(accounts[_sender] >= _amount, "amount exceeds balance");
        accounts[_sender] -= _amount;
        accounts[_recipient] += _amount;
        emit Transfer(_sender, _recipient, _amount);
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
    function _approve(address owner, address spender, uint256 amount) internal virtual {
        require(owner != address(0), "ERC20: approve from the zero address");
        require(spender != address(0), "ERC20: approve to the zero address");

        allowances[owner][spender] = amount;
        emit Approval(owner, spender, amount);
    }

    function _verifyEnode(Validator memory _validator) internal view {
    // _enode can't be empty and needs to be well-formed.
        uint _err;
        (_validator.addr, _err) = Precompiled.enodeCheck(_validator.enode);
        require(_err == 0, "enode error");
        require(validators[_validator.addr].addr == address(0), "validator already registered");
        require(_validator.commissionRate <= COMMISSION_RATE_PRECISION, "invalid commission rate");

    }

    function _deployLiquidContract(Validator memory _validator) internal {
        if (address(_validator.liquidContract) == address(0)) {
            _validator.liquidContract = new Liquid(_validator.addr,
                _validator.treasury,
                _validator.commissionRate);
        }
        validatorList.push(_validator.addr);
        validators[_validator.addr] = _validator;
    }


    function _registerValidator(Validator memory _validator) internal {
        _verifyEnode(_validator);
        // deploy liquid stake contract
        _deployLiquidContract(_validator);
    }

    function _registerValidator(Validator memory _validator, bytes memory proof) internal {
        // verify Enode
        _verifyEnode(_validator);

        // verify proof
        bytes32 r;
        bytes32 s;
        uint8 v;
        (r, s, v) = Helpers.splitProof(proof);
        bytes memory prefix = "\x19Ethereum Signed Message:\n";
        bytes memory treasury = abi.encodePacked(_validator.treasury);
        bytes memory data = abi.encodePacked(prefix, Helpers.toString(treasury.length), treasury);
        bytes32 hashed = keccak256(data);
        address signer = ecrecover(hashed, v, r, s);
        require(signer == _validator.addr, "Invalid proof provided for registration");

        // deploy liquid stake contract
       _deployLiquidContract(_validator);
    }


    /**
    * @dev Internal function pausing the specified validator. Paused validators
    * can no longer be delegated stake and can no longer be part of the consensus committe.
    * Warning: no checks are done here.
    * Emit {DisabledValidator} event.
    */
    function _pauseValidator(address _address) internal virtual {
        Validator storage val = validators[_address];
        require(val.state == ValidatorState.active, "validator must be enabled");

        val.state = ValidatorState.paused;
        //effectiveBlock may not be accurate if the epoch duration gets modified.
        emit PausedValidator(val.treasury, _address,  lastEpochBlock + config.epochPeriod);
    }


    /**
     * @dev Create a bonding object of `amount` stake token with the `_recipient` address.
     * This object will be processed at epoch finalization.
     *
     * This function assume that `_validator` is a valid validator address.
     */
    function _bond(address _validator, uint256 _amount, address payable _recipient) internal virtual{

        require(_amount > 0, "amount need to be strictly positive");
        require(accounts[_recipient] >= _amount, "insufficient Newton balance");

        accounts[_recipient] -= _amount;
        Staking memory _bonding = Staking(_recipient, _validator, _amount, block.number);
        bondingMap[headBondingID] = _bonding;
        headBondingID++;
    }

    function _applyBonding(uint256 id) internal {
        Staking storage _bonding = bondingMap[id];
        Validator storage _validator = validators[_bonding.delegatee];

        /* The conversion rate is equal to the ratio of issued liquid tokens
             over the total amount of bonded staked tokens. */
        uint256 _liquidAmount;
        if (_validator.bondedStake == 0) {
            _liquidAmount = _bonding.amount;
        } else {
            _liquidAmount = (_validator.liquidSupply * _bonding.amount) / _validator.bondedStake;
        }

        _validator.liquidContract.mint(_bonding.delegator, _liquidAmount);
        _validator.bondedStake += _bonding.amount;
        _validator.liquidSupply += _liquidAmount;
    }

    function _unbond(address _validator, uint256 _amount, address payable _recipient) internal virtual {
        uint256 liqBalance = validators[_validator].liquidContract.balanceOf(_recipient);
        require(liqBalance >= _amount, "insufficient Liquid Newton balance");

        validators[_validator].liquidContract.burn(_recipient, _amount);
        Staking memory _unbonding = Staking(_recipient, _validator, _amount, block.number);
        unbondingMap[headUnbondingID] = _unbonding;
        headUnbondingID++;
    }

    function _applyUnbonding(uint256 id) internal virtual {
        Staking storage _unbonding = unbondingMap[id];
        Validator storage validator = validators[_unbonding.delegatee];
        /* validator.liquidSupply must not be equal to zero here */
        uint256 _newtonAmount = (_unbonding.amount * validator.bondedStake) / validator.liquidSupply;

        validator.bondedStake -= _newtonAmount;
        validator.liquidSupply -= _unbonding.amount;
        accounts[_unbonding.delegator] += _newtonAmount;
    }

    function _applyNewCommissionRates() internal virtual {
        while(commissionRateChangeQueueFirst < commissionRateChangeQueueLast) {
            // check unbonding period

            CommissionRateChangeRequest storage _curRequest = commissionRateChangeQueue[commissionRateChangeQueueFirst];
            if(_curRequest.startBlock + config.unbondingPeriod > block.number){
                break;
            }

            // change commission rate for liquid staking accounts
            validators[_curRequest.validator].commissionRate = _curRequest.rate;
            validators[_curRequest.validator].liquidContract.setCommissionRate(_curRequest.rate);

            delete commissionRateChangeQueue[commissionRateChangeQueueFirst];

            commissionRateChangeQueueFirst += 1;
        }
    }

    /* Should be called at every epoch */
    function _stakingTransitions() internal virtual {
        for (uint256 i = tailBondingID; i < headBondingID; i++) {
            _applyBonding(i);
        }
        tailBondingID = headBondingID;

        uint256 _processedId = tailUnbondingID;
        for (uint256 i = tailUnbondingID; i < headUnbondingID; i++) {
            if (unbondingMap[i].startBlock + config.unbondingPeriod <= block.number) {
                _applyUnbonding(i);
                _processedId += 1;
            } else {
                break;
            }
        }
        tailUnbondingID = _processedId;
    }

    /**
    * @dev Order validators by stake
    */
    function _sortByStake(Validator[] memory _validators) internal pure {
        _structQuickSort(_validators, int(0), int(_validators.length - 1));
    }

    /**
    * @dev QuickSort algorithm sorting in ascending order by stake.
    */
    function _structQuickSort(Validator[] memory _users, int _low, int _high) internal pure {

        int _i = _low;
        int _j = _high;
        if (_i == _j) return;
        uint _pivot = _users[uint(_low + (_high - _low) / 2)].bondedStake;
        // Set the pivot element in its right sorted index in the array
        while (_i <= _j) {
            while (_users[uint(_i)].bondedStake > _pivot) _i++;
            while (_pivot > _users[uint(_j)].bondedStake) _j--;
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
