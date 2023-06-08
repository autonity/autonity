// SPDX-License-Identifier: MIT

pragma solidity ^0.8.3;

import "./interfaces/IERC20.sol";
import "./Liquid.sol";
import "./Upgradeable.sol";
import "./Precompiled.sol";
import "./Helpers.sol";
import "./SafeMath.sol";
import "./lib/BytesLib.sol";
import "./Oracle.sol";

uint8 constant DECIMALS = 18;

/** @title Proof-of-Stake Autonity Contract */
contract Autonity is IERC20, Upgradeable {
    using SafeMath for uint256;
    // constant settings for accountability protocol.
    uint256 constant MAX_ROUND = 99;
    uint256 constant PROVE_WINDOW = 60; // wait 60 blocks to get innocent proof of accusation before turning it slashed.
    uint256 constant LATEST_ACCOUNTABILITY_EVENTS_RANGE = 256;

    // the penalty in amount of stake token for accountability event, managed by operator account.
    uint256 faultPenalty = 1;

    uint256 public constant COMMISSION_RATE_PRECISION = 10_000;

    enum AccountabilityEventType {Misbehaviour, Accusation, Innocence}
    struct AccountabilityEvent {
        uint8 Chunks;     // Counter of number of chunks for oversize accountability event
        uint8 ChunkID;    // Chunk index to construct the oversize accountability event
        uint8 Type;       // Accountability event types: Misbehaviour, Accusation, Innocence.
        uint8 Rule;       // Rule ID defined in AFD rule engine.
        address Reporter; // The node address of the validator who report this event, for incentive protocol.
        address Sender;   // The corresponding node address of this accountability event.
        bytes32 MsgHash;  // The corresponding consensus msg's hash of this accountability event.
        bytes RawProof;   // rlp encoded bytes of Proof object.
    }

    Oracle oracleContract = Oracle(payable(0x5a443704dd4B594B382c22a083e2BD3090A6feF3));
    enum ValidatorState {active, paused}
    struct Validator {
        address payable treasury;
        address nodeAddress;
        address oracleAddress;
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

    // accountabilityEventChunks, a storage to construct oversize accountability event with chunked proof bytes.
    // Mapping(MsgHash => Mapping(Type => Mapping(Rule => Mapping(Reporter=>Mapping(ChunkID => AccountabilityEvent))))) chunkedAccountabilityEvents;
    mapping(bytes32 => mapping(uint8 => mapping(uint8 => mapping(address=>mapping(uint8 => AccountabilityEvent))))) accountabilityEventChunks;

    // accountability events against per validator.
    mapping(address => AccountabilityEvent[]) validatorMisbehaviours;
    mapping(address => AccountabilityEvent[]) validatorAccusations;

    // map accountability event with its timestamp (block height) on when it was processed.
    mapping (bytes32 => uint256) misbehaviourProcessedTS;
    mapping (bytes32 => uint256) accusationProcessedTS;

    // pending chunked accountability event.
    AccountabilityEvent[] private pendingChunkedAccountabilityEvent;

    // pending slash tasks for per epoch.
    AccountabilityEvent[] private pendingSlashTasks;

    mapping (address => uint256) stakeSlashed;

    // Stake token state transitions happen every epoch.
    uint256 public epochID;
    uint256 public lastEpochBlock;
    uint256 public epochTotalBondedStake;

    CommitteeMember[] internal committee;
    uint256 public totalRedistributed;
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
    event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract);
    event PausedValidator(address treasury, address addr, uint256 effectiveBlock);
    event Rewarded(address addr, uint256 amount);
    event MisbehaviourAdded(AccountabilityEvent ev);
    event AccusationAdded(AccountabilityEvent ev);
    event AccusationRemoved(AccountabilityEvent ev);
    event EpochPeriodUpdated(uint256 period);
    event MisbehaviourPenaltyUpdated(uint256 penalty);
    event NodeSlashed(address validator, uint256 penalty);
    event SubmitGuiltyAccusation(AccountabilityEvent ev);
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

            _registerValidator(_validators[i]);

            accounts[_validators[i].treasury] += _bondedStake;
            stakeSupply += _bondedStake;
            epochTotalBondedStake += _bondedStake;
            _bond(_validators[i].nodeAddress, _bondedStake, payable(_validators[i].treasury));
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
    * @notice handle accountability event in the Autonity Contract.
    */
    function handleAccountabilityEvents(AccountabilityEvent[] memory _events) public onlyValidator {
        for (uint256 i = 0; i < _events.length; i++) {

            if (_events[i].Reporter != msg.sender) {
                continue;
            }

            // if the event is a chunked event, store it.
            if (_events[i].Chunks != 0) {
                _storeChunkedAccountabilityEvent(_events[i]);
                continue;
            }

            if (AccountabilityEventType(_events[i].Type) == AccountabilityEventType.Misbehaviour) {
                if (misbehaviourProcessedTS[_events[i].MsgHash] == 0) {
                    _handleMisbehaviour(_events[i]);
                    continue;
                }
            }

            if (AccountabilityEventType(_events[i].Type) == AccountabilityEventType.Accusation) {
                if (accusationProcessedTS[_events[i].MsgHash] == 0) {
                    _handleAccusation(_events[i]);
                    continue;
                }
            }

            if (AccountabilityEventType(_events[i].Type) == AccountabilityEventType.Innocence) {
                if (accusationProcessedTS[_events[i].MsgHash] != 0) {
                    _handleInnocenceProof(_events[i]);
                    continue;
                }
            }
        }
    }

    /*
    * @return the number of decimals the NTN token uses.
    * @dev ERC-20 Optional.
    */
    function decimals() public pure returns (uint8) {
        return DECIMALS;
    }

    /**
    * @notice Register a new validator in the system.  The validator might be selected to be part of consensus.
    * This validator will have assigned to its treasury account the caller of this function.
    * A new token "Liquid Stake" is deployed at this phase.
    * @param _enode enode identifying the validator node.
    * @param _multisig is a combination of two signatures appended sequentially, In below order:
        1. a message containing treasury account and signed by validator account private key .
        2. a message containing treasury account and signed by Oracle account private key .
    * @dev Emit a {RegisteredValidator} event.
    */
    function registerValidator(string memory _enode, address _oracleAddress, bytes memory _multisig) public {
        Validator memory _val = Validator(payable(msg.sender), //treasury
            address(0), // address
            _oracleAddress, // voter Address //TODO: update validator registration API
            _enode, // enode
            config.delegationRate, // validator commission rate
            0, // bonded stake
            0, // total slashed
            Liquid(address(0)), // liquid token contract
            0, // liquid token supply
            block.number,
            ValidatorState.active
        );

        _registerValidator(_val, _multisig);
        emit RegisteredValidator(msg.sender, _val.nodeAddress, _oracleAddress, _enode, address(_val.liquidContract));
    }

    function bond(address _validator, uint256 _amount) public {
        require(validators[_validator].nodeAddress == _validator, "validator not registered");
        require(validators[_validator].state == ValidatorState.active, "validator need to be active");
        _bond(_validator, _amount, payable(msg.sender));
    }

    function unbond(address _validator, uint256 _amount) public {
        require(validators[_validator].nodeAddress == _validator, "validator not registered");
        _unbond(_validator, _amount, payable(msg.sender));
    }

    /**
    * @notice Pause the validator and stop it accepting delegations. See ADR-004 for more details.
    * @param _address address to be disabled.
    * @dev emit a {DisabledValidator} event.
    */
    function pauseValidator(address _address) public {
        require(validators[_address].nodeAddress == _address, "validator must be registered");
        require(validators[_address].treasury == msg.sender, "require caller to be validator admin account");
        _pauseValidator(_address);
    }

    /**
    * @notice Re-activate the specified validator. See ADR-004 for more details.
    * @param _address address to be enabled.
    */
    function activateValidator(address _address) public {
        require(validators[_address].nodeAddress == _address, "validator must be registered");
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
        require(validators[_validator].nodeAddress == _validator, "validator must be registered");
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
    * @notice Set the misbehaviour penalty. Restricted to the Operator account.
    * @param _newPenalty Positive integer.
    */
    function setMisbehaviourPenalty(uint256 _newPenalty) public onlyOperator {
        if (_newPenalty <= 0) {
            return;
        }
        faultPenalty = _newPenalty;
        emit MisbehaviourPenaltyUpdated(_newPenalty);
    }

    /*
    * @notice Set the epoch period. Restricted to the Operator account.
    * @param _period Positive integer.
    */
    function setEpochPeriod(uint256 _period) public onlyOperator {
        if (_period == config.epochPeriod) {
            return;
        }

        // to decrease the epoch period, we need to check if current chain head already exceed the window:
        // lastBlockEpoch + _newPeriod, if so, the _newPeriod cannot be applied since the finalization of current epoch
        // at finalize function will never be triggered, in such case, operator need to find better timing to do so.
        if (_period < config.epochPeriod) {
            if (block.number >= lastEpochBlock + _period) {
                revert("current chain head exceed the window: lastBlockEpoch + _newPeriod, try again latter on.");
            }
        }
        config.epochPeriod = _period;
        emit EpochPeriodUpdated(_period);
    }

    /*
    * @notice Set the Operator account. Restricted to the Operator account.
    * @param _account the new operator account.
    */
    function setOperatorAccount(address _account) public onlyOperator {
        config.operatorAccount = _account;
        oracleContract.setOperator(_account);
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
    * @return upgrade Set to true if an autonity contract upgrade is available.
    * @return committee The next block consensus committee.
    */
    function finalize() external virtual onlyProtocol
    returns (bool, CommitteeMember[] memory) {
        // on each block finalize, process those pending chunkedAccountabilityEvents
        _handleChunkedAccountabilityEvent();
        // on each block, try to promote accusations without proof of innocence to misbehaviour.
        _submitGuiltyAccusations();

        if (lastEpochBlock + config.epochPeriod == block.number) {
            // - slashing should come here first -
            _performSlashTasks();
            _performRedistribution();
            _stakingTransitions();
            _applyNewCommissionRates();
            address[] memory voters = computeCommittee();
            oracleContract.setVoters(voters);
            lastEpochBlock = block.number;
            epochID += 1;
        }
        oracleContract.finalize();
        return (contractUpgradeReady, committee);
    }

    /**
    * @dev call by validator client to check if a misbehaviour is already processed on chain.
    * @param _msgHash the msg hash of the malicious msg detected by AFD.
    */
    function misbehaviourProcessed(bytes32 _msgHash) public view returns (bool) {
        return misbehaviourProcessedTS[_msgHash] != 0;
    }

    /**
    * @dev call by validator client to check if an accusation is already processed on chain.
    * @param _msgHash the msg hash of the suspected msg detected by AFD.
    */
    function accusationProcessed(bytes32 _msgHash) public view returns (bool) {
        return accusationProcessedTS[_msgHash] != 0;
    }

    /**
    * @dev Get the slashed stake for a validator.
    * @param _addr The address a validator.
    */
    function getSlashedStake(address _addr) public view returns (uint256) {
        require(validators[_addr].nodeAddress == _addr, "validator must be registered");
        return stakeSlashed[_addr];
    }

    /**
    * @dev Dump the latest 256 misbehaviours of a validator.
    * @param _addr The address of a validator.
    */
    function getValidatorRecentMisbehaviours(address _addr) public view returns (AccountabilityEvent[] memory) {
        require(validators[_addr].nodeAddress == _addr, "validator must be registered");
        // only return latest 256 number of on-chain misbehaviours.
        if (validatorMisbehaviours[_addr].length > LATEST_ACCOUNTABILITY_EVENTS_RANGE) {
            uint256 start = validatorMisbehaviours[_addr].length - LATEST_ACCOUNTABILITY_EVENTS_RANGE;
            AccountabilityEvent[] memory _misbehaviours = new AccountabilityEvent[](LATEST_ACCOUNTABILITY_EVENTS_RANGE);
            for (uint256 i = start; i < validatorMisbehaviours[_addr].length; i++) {
                AccountabilityEvent memory _ev = validatorMisbehaviours[_addr][i];
                _misbehaviours[i-start] = _ev;
            }
            return _misbehaviours;
        }
        return validatorMisbehaviours[_addr];
    }

    function getAccountabilityEventChunk(bytes32 _msgHash, uint8 _type, uint8 _rule, address _reporter, uint8 _chunkID) public view returns (bytes memory) {
        bytes memory ret;
        if (accountabilityEventChunks[_msgHash][_type][_rule][_reporter][_chunkID].MsgHash == _msgHash) {
            return accountabilityEventChunks[_msgHash][_type][_rule][_reporter][_chunkID].RawProof;
        }
        return ret;
    }

    /**
    * @dev Dump the latest 256 accusations of a validator.
    * @param _addr The address of a validator.
    */
    function getValidatorRecentAccusations(address _addr) public view returns (AccountabilityEvent[] memory) {
        require(validators[_addr].nodeAddress == _addr, "validator must be registered");
        // only return latest 256 number of on-chain accusations.
        if (validatorAccusations[_addr].length > LATEST_ACCOUNTABILITY_EVENTS_RANGE) {
            uint256 start = validatorAccusations[_addr].length - LATEST_ACCOUNTABILITY_EVENTS_RANGE;
            AccountabilityEvent[] memory _accusations = new AccountabilityEvent[](LATEST_ACCOUNTABILITY_EVENTS_RANGE);
            for (uint256 i = start; i < validatorAccusations[_addr].length; i++) {
                AccountabilityEvent memory _ev = validatorAccusations[_addr][i];
                _accusations[i-start] = _ev;
            }
            return _accusations;
        }
        return validatorAccusations[_addr];
    }

    /**
    * @notice update the current committee by selecting top staking validators.
    * Restricted to the protocol.
    */
    function computeCommittee() public onlyProtocol returns (address[] memory){
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
        address [] memory _voterList = new address[](_committeeLength);

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
            CommitteeMember memory _member = CommitteeMember(_committeeList[_k].nodeAddress, _committeeList[_k].bondedStake);
            committee.push(_member);
            committeeNodes.push(_committeeList[_k].enode);
            _voterList[_k] = _committeeList[_k].oracleAddress;
            epochTotalBondedStake += _committeeList[_k].bondedStake;
        }
        return _voterList;
    }


    /*
    ============================================================
        Getters
    ============================================================
    */

    /**
    * @notice Returns the epoch period.
    */
    function getEpochPeriod() external view returns (uint256) {
        return config.epochPeriod;
    }

    /**
    * @notice Returns the block period.
    */
    function getBlockPeriod() external view returns (uint256) {
        return config.blockPeriod;
    }

    /**
* @notice Returns the un-bonding period.
    */
    function getUnbondingPeriod() external view returns (uint256) {
        return config.unbondingPeriod;
    }

    /**
    * @notice Returns the penalty in stake token.
    */
    function getPenalty() external view returns (uint256) {
        return faultPenalty;
    }

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
     * @notice Returns the current treasury account.
     */
    function getTreasuryAccount() external view returns (address) {
        return config.treasuryAccount;
    }

    /**
     * @notice Returns the current treasury fee.
     */
    function getTreasuryFee() external view returns (uint256) {
        return config.treasuryFee;
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
        require(validators[_addr].nodeAddress == _addr, "validator not registered");
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

    /**
    * @dev Modifier that checks if the caller is a validator from the pool.
    */
    modifier onlyValidator {
        require(validators[msg.sender].nodeAddress == msg.sender, "function restricted to the validator");
        _;
    }

    /*
    ============================================================

        Internals

    ============================================================
    */

    /**
    * @notice Take fund away from faulty node account.
    * @dev Emit a {NodeSlashed} event for account that is fined.
    */
    function _takePenalty(address addr) internal {
        // todo: take penalty by different degree for different faults.
        uint256 penalty = faultPenalty;
        // keep at least 1 stake for network liveness.
        if (validators[addr].bondedStake <= 1) {
            return;
        }

        // resolve proper penalty.
        if (validators[addr].bondedStake <= penalty) {
            penalty = validators[addr].bondedStake - 1;
        }

        // apply penalty and save fine at contract account.
        validators[addr].bondedStake -= penalty;
        accounts[address(this)] += penalty;
        validators[addr].totalSlashed += 1;
        stakeSlashed[addr] += penalty;
        emit NodeSlashed(addr, penalty);
    }

    /**
    * @notice Perform Auton reward distribution. The transaction fees
    * are simply re-distributed to all stake-holders, including validators,
    * pro-rata the amount of stake held.
    * @dev Emit a {BlockReward} event for every account that collected rewards.
    */
    function _performRedistribution() internal virtual {
        if (address(this).balance == 0) {
            return;
        }
        uint256 _amount =  address(this).balance;
        // take treasury fee.
        uint256 _treasuryReward = (config.treasuryFee * _amount) / 10 ** 18;
        if (_treasuryReward > 0) {
            //treasuryAccount.transfer(_treasuryReward);
            (bool sent, bytes memory data) = config.treasuryAccount.call{value: _treasuryReward}("");
            if (sent == true) {
                _amount -= _treasuryReward;
            } else {
                // todo: emit an event to indicate reward distribution for treasury account failed, the reward for treasury
                //  will be distributed through validators.
            }
        }
        totalRedistributed += _amount;
        // otherwise, do reward distribution.
        _rewardDistribution(_amount);
    }

    /**
    * @notice perform reward distribution with promoted members, the promotion is base on current epoch's omission counters.
    * @dev Emit a {Rewarded} event for every account that are rewarded.
    */
    function _rewardDistribution(uint256 _amount) internal {
        if (_amount <= 0) {
            return;
        }

        // count total voting powers after slashing.
        uint256 totalBondedStake = 0;
        for (uint256 j = 0; j < committee.length; j++) {
            Validator storage _val = validators[committee[j].addr];
            totalBondedStake += _val.bondedStake;
        }

        if (totalBondedStake == 0) {
            // this shouldn't happens, but to prevent finalize() from getting reverted.
            return;
        }

        uint256 totalDistributed = 0;
        for (uint256 i = 0; i < committee.length; i++) {
            Validator storage _val = validators[committee[i].addr];
            uint256 _reward = (_val.bondedStake * _amount) / totalBondedStake;
            if(_reward > 0) {
                uint256 distributed = _val.liquidContract.redistribute{value: _reward}();
                totalDistributed += distributed;
                emit Rewarded(_val.nodeAddress, distributed);
            }
        }
        // todo: to enable below dust fee handling, the corresponding testcases should be rewritten.
        /*
        // the DIV operator generates dust reward fraction, transfer dust fraction tokens to treasury account.
        uint256 dust = _amount - totalDistributed;
        if (dust > 0) {
            (bool sent, bytes memory data) = config.treasuryAccount.call{value: dust}("");
            if (sent == true) {
                emit Rewarded(config.treasuryAccount, dust);
            }
        }*/
    }

    /**
    * @notice perform slashing over faulty validators at the end of epoch. The fine in stake token are moved from
    * validator account to autonity contract account, and the corresponding slash counter as a reputation for validator
    * increase too.
    * @dev Emit a {NodeSlashed} event for every account that are slashed.
    */
    function _performSlashTasks() internal {
        // todo resolve the fine base on different accountability events.
        // slash validator of misbehaviour and accusations without innocent proof.
        for (uint256 i = 0; i < pendingSlashTasks.length; i++) {
            address addr = pendingSlashTasks[i].Sender;
            _takePenalty(addr);
        }
        // reset pending slash task queue for next epoch.
        delete pendingSlashTasks;
    }

    /**
    * @notice promote accusations without innocence proof in the prove-window into misbehaviour.
    * @dev Emit a {SubmitGuiltyAccusation} event for every accusation that are going to be slashed.
    */
    function _submitGuiltyAccusations() internal {
        // for each committee member, find no innocence proof accusations within proveWindow, promote them into
        // as misbehaviour, and push them in slashing task queue, and finally remove them from accusation list.
        for (uint256 i = 0; i < committee.length; i++) {
            uint256 len = validatorAccusations[committee[i].addr].length;
            AccountabilityEvent[] memory promoted = new AccountabilityEvent[](len);
            uint256 promotedCounter = 0;

            for (uint256 j = 0; j < len; j++) {
                AccountabilityEvent memory accusation = validatorAccusations[committee[i].addr][j];
                uint256 ttl = block.number - accusationProcessedTS[accusation.MsgHash];
                if (accusationProcessedTS[accusation.MsgHash]!= 0 && ttl > PROVE_WINDOW) {
                    // promote accusation from accusation list to misbehaviour list.
                    validatorMisbehaviours[committee[i].addr].push(accusation);
                    // push accusation for slashing.
                    pendingSlashTasks.push(accusation);
                    emit SubmitGuiltyAccusation(accusation);
                    promoted[promotedCounter] = accusation;
                    promotedCounter++;
                }
            }
            // remove accusation from accusation list.
            for (uint256 k = 0; k < promotedCounter; k++) {
                _removeAccusation(promoted[k]);
            }
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
        (_validator.nodeAddress, _err) = Precompiled.enodeCheck(_validator.enode);
        require(_err == 0, "enode error");
        require(validators[_validator.nodeAddress].nodeAddress == address(0), "validator already registered");
        require(_validator.commissionRate <= COMMISSION_RATE_PRECISION, "invalid commission rate");
    }

    function _deployLiquidContract(Validator memory _validator) internal {
        if (address(_validator.liquidContract) == address(0)) {
            string memory stringLength = Helpers.toString(validatorList.length);
            _validator.liquidContract = new Liquid(_validator.nodeAddress,
                _validator.treasury,
                _validator.commissionRate,
                stringLength);
        }
        validatorList.push(_validator.nodeAddress);
        validators[_validator.nodeAddress] = _validator;
    }


    function _registerValidator(Validator memory _validator) internal {
        _verifyEnode(_validator);
        // deploy liquid stake contract
        _deployLiquidContract(_validator);
    }

    function _removeAccusation(AccountabilityEvent memory ev) internal {
        uint256 len = validatorAccusations[ev.Sender].length;
        for (uint256 i = 0; i < len; i++) {
            if (validatorAccusations[ev.Sender][i].MsgHash == ev.MsgHash) {
                validatorAccusations[ev.Sender][i] = validatorAccusations[ev.Sender][len - 1];
                validatorAccusations[ev.Sender].pop();
                break;
            }
        }
    }

    function _registerValidator(Validator memory _validator, bytes memory _multisig) internal {
        require(_multisig.length == 130, "Invalid proof length");
        // verify Enode
        _verifyEnode(_validator);

        bytes memory prefix = "\x19Ethereum Signed Message:\n";
        bytes memory treasury = abi.encodePacked(_validator.treasury);
        bytes32 hashedData = keccak256(abi.encodePacked(prefix, Helpers.toString(treasury.length), treasury));
        address[] memory signers = new address[](2);
        bytes32 r;
        bytes32 s;
        uint8 v;
        //start from 32th byte to skip the encoded length field from the bytes type variable
        for (uint i=32; i < _multisig.length; i +=65) {
            (r, s, v) = Helpers.extractRSV(_multisig, i);
            signers[i/65] = ecrecover(hashedData, v, r, s);
        }
        require(signers[0] == _validator.nodeAddress, "Invalid node key ownership proof provided");
        require(signers[1] == _validator.oracleAddress, "Invalid oracle key ownership proof provided");

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

        uint256 _liqSupply = validators[_validator].liquidContract.totalSupply();
        require(!(_inCommittee(_validator) && _amount == _liqSupply),
            "can't have committee member without LNTN");

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

    /**
    * @dev handle misbehaviour and push the ev for slashing once it is a valid proof.
    */
    function _handleMisbehaviour(AccountabilityEvent memory _ev) internal {
        // Validate the misbehaviour proof
        (address addr, bytes32 msgHash, uint256 retCode, uint256 ruleID) =
        Precompiled.checkAccountabilityEvent(Precompiled.MISBEHAVIOUR_CONTRACT, _ev.RawProof);
        if (msgHash != _ev.MsgHash || addr != _ev.Sender || retCode == 0 || ruleID != uint256(_ev.Rule)) {
            return;
        }

        // if event is oversize, dont store duplicated raw bytes since we have copy in chunked event storage.
        if (_ev.Chunks > 0) {
            delete _ev.RawProof;
        }

        // if the misbehaviour is from validator, save proof and add slashing task.
        if (validators[_ev.Sender].nodeAddress == _ev.Sender) {
            validatorMisbehaviours[_ev.Sender].push(_ev);
            pendingSlashTasks.push(_ev);
            misbehaviourProcessedTS[_ev.MsgHash] = block.number;
            emit MisbehaviourAdded(_ev);
        }
    }
    /**
    * @dev handle accusation and push the ev into the waiting queue before it become to be slashed.
    */
    function _handleAccusation(AccountabilityEvent memory _ev) internal {
        // Validate the accusation proof
        (address addr, bytes32 msgHash, uint256 retCode, uint256 ruleID) =
        Precompiled.checkAccountabilityEvent(Precompiled.ACCUSATION_CONTRACT, _ev.RawProof);
        if (msgHash != _ev.MsgHash || addr != _ev.Sender || retCode == 0 || ruleID != uint256(_ev.Rule)) {
            return;
        }

        // if event is oversize, dont store duplicated raw bytes since we have copy in chunked event storage.
        if (_ev.Chunks > 0) {
            delete _ev.RawProof;
        }
        // accusation proof is valid, store the proof, and wait for validator to provide innocence proof.
        validatorAccusations[_ev.Sender].push(_ev);
        accusationProcessedTS[_ev.MsgHash] = block.number;
        emit AccusationAdded(_ev);
    }

    /**
    * @dev handle innocence proof and remove the corresponding accusation once it is valid proof.
    */
    function _handleInnocenceProof(AccountabilityEvent memory _ev) internal {
        // Validate the proof of innocence
        (address addr, bytes32 msgHash, uint256 retCode, uint256 ruleID) =
        Precompiled.checkAccountabilityEvent(Precompiled.INNOCENCE_CONTRACT, _ev.RawProof);
        if (msgHash != _ev.MsgHash || addr != _ev.Sender || retCode == 0 || ruleID != uint256(_ev.Rule)) {
            return;
        }
        // if event is oversize, dont store duplicated raw bytes since we have copy in chunked event storage.
        if (_ev.Chunks > 0) {
            delete _ev.RawProof;
        }
        // innocence proof is valid, remove accusation.
        _removeAccusation(_ev);
        delete accusationProcessedTS[_ev.MsgHash];
        emit AccusationRemoved(_ev);
    }

    /**
    * @dev saves chunked accountability events in the storage, once the chunks are fully collected, they will be processed
    * at the end of block finalization phase.
    */
    function _storeChunkedAccountabilityEvent(AccountabilityEvent memory _ev) internal {

        if (AccountabilityEventType(_ev.Type) == AccountabilityEventType.Misbehaviour && misbehaviourProcessed(_ev.MsgHash) == true) {
            return;
        }

        if (AccountabilityEventType(_ev.Type) == AccountabilityEventType.Accusation && accusationProcessed(_ev.MsgHash) == true) {
            return;
        }

        // save the chunk and try to construct the full event and make it ready for processing.
        accountabilityEventChunks[_ev.MsgHash][_ev.Type][_ev.Rule][_ev.Reporter][_ev.ChunkID] = _ev;

        // to save the gas from useless bytes concat, we just 1st check if all the chunks were collected.
        for (uint8 chunkID = 0; chunkID < _ev.Chunks; chunkID++) {
            if (accountabilityEventChunks[_ev.MsgHash][_ev.Type][_ev.Rule][_ev.Reporter][chunkID].Reporter != msg.sender)
            {
                return;
            }
        }

        // now, all the chunks are collected, we can start the processing. But due to the concat of chunks would cost a
        // lot of gas from the msg sender, it would limit the max size of event we can process, so we would push this
        // event into a pending list, which will be processed by the finalize() at the block finalization phase which
        // would cost none gas for the chunk concat.

        AccountabilityEvent memory pendingEvent;
        pendingEvent.Reporter = _ev.Reporter;
        pendingEvent.MsgHash = _ev.MsgHash;
        pendingEvent.Chunks = _ev.Chunks;
        pendingEvent.Rule = _ev.Rule;
        pendingEvent.ChunkID = _ev.ChunkID;
        pendingEvent.Type = _ev.Type;
        pendingEvent.Sender = _ev.Sender;

        // push the event in the pending list, it will be handled at current block finalization phase.
        pendingChunkedAccountabilityEvent.push(pendingEvent);
        // to avoid duplicated reporting from different validator, once the chunked event is pushed in the processing list,
        // we set them to be processed.
        if (AccountabilityEventType(_ev.Type) == AccountabilityEventType.Misbehaviour) {
            misbehaviourProcessedTS[_ev.MsgHash] = block.number;
        }
        if (AccountabilityEventType(_ev.Type) == AccountabilityEventType.Accusation) {
            accusationProcessedTS[_ev.MsgHash] = block.number;
        }
    }
    /**
    * @dev at block finalize phase, call this function to check if there are fully collected events, and process it.
    */
    function _handleChunkedAccountabilityEvent() internal {
        for (uint256 i = 0; i < pendingChunkedAccountabilityEvent.length; i++) {

            AccountabilityEvent memory ev;
            ev.Reporter = pendingChunkedAccountabilityEvent[i].Reporter;
            ev.MsgHash = pendingChunkedAccountabilityEvent[i].MsgHash;
            ev.Chunks = pendingChunkedAccountabilityEvent[i].Chunks;
            ev.Rule = pendingChunkedAccountabilityEvent[i].Rule;
            ev.ChunkID = pendingChunkedAccountabilityEvent[i].ChunkID;
            ev.Type = pendingChunkedAccountabilityEvent[i].Type;
            ev.Sender = pendingChunkedAccountabilityEvent[i].Sender;

            //bytes memory constructedBytes;
            for (uint8 chunkID = 0; chunkID < ev.Chunks; chunkID++) {
                ev.RawProof = BytesLib.concat(ev.RawProof, accountabilityEventChunks[ev.MsgHash][ev.Type][ev.Rule][ev.Reporter][chunkID].RawProof);
            }

            if (AccountabilityEventType(ev.Type) == AccountabilityEventType.Misbehaviour) {
                _handleMisbehaviour(ev);
                continue;
            }

            if (AccountabilityEventType(ev.Type) == AccountabilityEventType.Accusation) {
                _handleAccusation(ev);
                continue;
            }

            if (AccountabilityEventType(ev.Type) == AccountabilityEventType.Innocence) {
                _handleInnocenceProof(ev);
                continue;
            }
        }
        delete pendingChunkedAccountabilityEvent;
    }

    function _inCommittee(address _validator) internal view returns (bool) {
        for (uint256 i = 0; i < committee.length; i++) {
            if (_validator == committee[i].addr){
                return true;
            }
        }
        return false;
    }
}
