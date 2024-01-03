// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.19;

import "./interfaces/IERC20.sol";
import "./Liquid.sol";
import "./Upgradeable.sol";
import "./Precompiled.sol";
import "./Helpers.sol";
import "./lib/BytesLib.sol";
import "./asm/IACU.sol";
import "./asm/ISupplyControl.sol";
import "./asm/IStabilization.sol";
import "./interfaces/IAccountability.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IAutonity.sol";

/** @title Proof-of-Stake Autonity Contract */
enum ValidatorState {active, paused, jailed, jailbound}
uint8 constant DECIMALS = 18;

contract Autonity is IAutonity, IERC20, Upgradeable {
    uint256 internal constant MAX_ROUND = 99;
    uint256 public constant COMMISSION_RATE_PRECISION = 10_000;

    struct Validator {
        address payable treasury;
        address nodeAddress;
        address oracleAddress;
        string enode; //addr must match provided enode
        uint256 commissionRate;
        uint256 bondedStake;
        uint256 unbondingStake;
        uint256 unbondingShares; // not effective - used for accounting purposes
        uint256 selfBondedStake;
        // bonded stake = selfBounded stake + delegated stake
        uint256 selfUnbondingStake;
        uint256 selfUnbondingShares; // not effective - used for accounting purposes
        uint256 selfUnbondingStakeLocked;
        Liquid liquidContract;
        uint256 liquidSupply;
        uint256 registrationBlock;
        uint256 totalSlashed;
        uint256 jailReleaseBlock;
        uint256 provableFaultCount;
        ValidatorState state;
    }

    struct CommitteeMember {
        address addr;
        uint256 votingPower;
    }

    /**************************************************/
    // Todo: Create a FIFO structure library, integrate with Staking{}
    /* Used for epoched staking */
    struct BondingRequest {
        address payable delegator;
        address delegatee;
        uint256 amount;
        uint256 requestBlock;
    }
    mapping(uint256 => BondingRequest) internal bondingMap;
    uint256 internal tailBondingID;
    uint256 internal headBondingID;

    struct UnbondingRequest {
        address payable delegator;
        address delegatee;
        uint256 amount; // NTN for self-delegation, LNTN otherwise
        uint256 unbondingShare;
        uint256 requestBlock;
        bool unlocked;
        bool selfDelegation;
    }
    mapping(uint256 => UnbondingRequest) internal unbondingMap;
    uint256 internal tailUnbondingID;
    uint256 internal headUnbondingID;
    uint256 internal lastUnlockedUnbonding;

    /* Used to track commission rate change*/
    struct CommissionRateChangeRequest {
        address validator;
        uint256 startBlock;
        uint256 rate;
    }
    mapping(uint256 => CommissionRateChangeRequest) internal commissionRateChangeQueue;
    uint256 internal commissionRateChangeQueueFirst = 0;
    uint256 internal commissionRateChangeQueueLast = 0;

    /**************************************************/
    struct Contracts {
        IAccountability accountabilityContract;
        IOracle oracleContract;
        IACU acuContract;
        ISupplyControl supplyControlContract;
        IStabilization stabilizationContract;
    }

    struct Policy {
        uint256 treasuryFee;
        uint256 minBaseFee;
        uint256 delegationRate;
        uint256 unbondingPeriod;
        address payable treasuryAccount;
    }

    struct Protocol {
        address operatorAccount;
        uint256 epochPeriod;
        uint256 blockPeriod;
        uint256 committeeSize;
    }

    struct Config {
        Policy policy;
        Contracts contracts;
        Protocol protocol;
        uint256 contractVersion;
    }

    Config public config;
    address[] internal validatorList;

    // Stake token state transitions happen every epoch.
    uint256 public epochID;
    mapping(uint256 => uint256) internal blockEpochMap;
    uint256 public lastEpochBlock;
    uint256 public epochTotalBondedStake;

    CommitteeMember[] internal committee;
    uint256 public totalRedistributed;
    uint256 public epochReward;
    mapping(address => mapping(address => uint256)) internal allowances;


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
    event MintedStake(address indexed addr, uint256 amount);
    event BurnedStake(address indexed addr, uint256 amount);
    event CommissionRateChange(address indexed validator, uint256 rate);
    event BondingRejected(address delegator, address delegatee, uint256 amount);

    /** @notice This event is emitted when a bonding request to a validator node has been registered.
    * This request will only be effective at the end of the current epoch however the stake will be
    * put in custody immediately from the delegator's account.
    * @param validator The validator node account.
    * @param delegator The caller.
    * @param selfBonded True if the validator treasury initiated the request. No LNEW will be issued.
    * @param amount The amount of NEWTON to be delegated.
    */
    event NewBondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount);

    /** @notice This event is emitted when an unbonding request to a validator node has been registered.
    * This request will only be effective after the unbonding period, rounded to the next epoch.
    * Please note that because of potential slashing events during this delay period, the released amount
    * may or may not be correspond to the amount requested.
    * @param validator The validator node account.
    * @param delegator The caller.
    * @param selfBonded True if the validator treasury initiated the request.
    * @param amount If self-bonded this is the requested amount of NEWTON to be unbonded.
    * If not self-bonded, this is the amount of Liquid Newton to be unbonded.
    */
    event NewUnbondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount);

    event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract);
    event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock);
    event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock);
    event Rewarded(address indexed addr, uint256 amount);
    event EpochPeriodUpdated(uint256 period);
    event NewEpoch(uint256 epoch);

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

    function _initialize(
        Validator[] memory _validators,
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
            _validators[i].commissionRate = config.policy.delegationRate;
            _validators[i].state = ValidatorState.active;
            _validators[i].selfUnbondingStakeLocked = 0;

            _registerValidator(_validators[i]);

            accounts[_validators[i].treasury] += _bondedStake;
            stakeSupply += _bondedStake;
            _bond(_validators[i].nodeAddress, _bondedStake, payable(_validators[i].treasury));
        }
    }

    function finalizeInitialization() onlyProtocol public {
        _stakingOperations();
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
        Validator memory _val = Validator(
            payable(msg.sender),     // treasury
            address(0),              // address
            _oracleAddress,          // voter Address
            _enode,                  // enode
            config.policy.delegationRate,   // validator commission rate
            0,                       // bonded stake
            0,                       // unbonding stake
            0,                       // unbonding shares
            0,                       // self bonded stake
            0,                       // self unbonding stake
            0,                       // self unbonding shares
            0,                       // self unbonding stake locked
            Liquid(address(0)),      // liquid token contract
            0,                       // liquid token supply
            block.number,            // registration block
            0,                       // total slashed
            0,                       // jail release block
            0,                       // provable faults count
            ValidatorState.active    // state
        );

        _registerValidator(_val, _multisig);
        emit RegisteredValidator(msg.sender, _val.nodeAddress, _oracleAddress, _enode, address(_val.liquidContract));
    }

    /**
    * @notice Create a bonding(delegation) request with the caller as delegator.
    * @param _validator address of the validator to delegate stake to.
    *        _amount total amount of NTN to bond.
    */
    function bond(address _validator, uint256 _amount) public {
        require(validators[_validator].nodeAddress == _validator, "validator not registered");
        require(validators[_validator].state == ValidatorState.active, "validator need to be active");
        _bond(_validator, _amount, payable(msg.sender));
    }

    /**
    * @notice Create an unbonding request with the caller as delegator.
    * @param _validator address of the validator to unbond stake to.
    *        _amount total amount of NTN to unbond.
    */
    function unbond(address _validator, uint256 _amount) public {
        require(validators[_validator].nodeAddress == _validator, "validator not registered");
        require(_amount > 0, "unbonding amount is 0");
        _unbond(_validator, _amount, payable(msg.sender));
    }

    /**
    * @notice Pause the validator and stop it accepting delegations.
    * @param _address address to be disabled.
    * @dev emit a {DisabledValidator} event.
    */
    function pauseValidator(address _address) public {
        require(validators[_address].nodeAddress == _address, "validator must be registered");
        require(validators[_address].treasury == msg.sender, "require caller to be validator admin account");
        _pauseValidator(_address);
    }

    /**
    * @notice Re-activate the specified validator.
    * @param _address address to be enabled.
    */
    function activateValidator(address _address) public {
        require(validators[_address].nodeAddress == _address, "validator must be registered");
        Validator storage _val = validators[_address];
        require(_val.treasury == msg.sender, "require caller to be validator treasury account");
        require(_val.state != ValidatorState.active, "validator already active");
        require(!(_val.state == ValidatorState.jailed && _val.jailReleaseBlock > block.number), "validator still in jail");
        require(_val.state != ValidatorState.jailbound, "validator jailed permanently");
        _val.state = ValidatorState.active;
        emit ActivatedValidator(_val.treasury, _address, lastEpochBlock + config.protocol.epochPeriod);
    }

    /**
    * @notice Update the validator. Only accessible to the accountability contract.
    * The difference in bondedStake will go to the treasury account.
    * @param _val Validator to be updated.
    */
    function updateValidatorAndTransferSlashedFunds(Validator calldata _val) external onlyAccountability {
        uint256 _diffNewtonBalance = (validators[_val.nodeAddress].bondedStake - _val.bondedStake) +
                                     (validators[_val.nodeAddress].unbondingStake - _val.unbondingStake) +
                                     (validators[_val.nodeAddress].selfUnbondingStake - _val.selfUnbondingStake);
        accounts[config.policy.treasuryAccount] += _diffNewtonBalance;
        validators[_val.nodeAddress] = _val;
    }

    /**
    * @notice Change commission rate for the specified validator.
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
        config.policy.minBaseFee = _price;
        emit MinimumBaseFeeUpdated(_price);
    }

    /*
    * @notice Set the maximum size of the consensus committee. Restricted to the Operator account.
    * @param _size Positive integer.
    */
    function setCommitteeSize(uint256 _size) public onlyOperator {
        require(_size > 0, "committee size can't be 0");
        config.protocol.committeeSize = _size;
    }

    /*
    * @notice Set the unbonding period. Restricted to the Operator account.
    * @param _size Positive integer.
    */
    function setUnbondingPeriod(uint256 _period) public onlyOperator {
        config.policy.unbondingPeriod = _period;
    }

    /*
    * @notice Set the epoch period. Restricted to the Operator account.
    * @param _period Positive integer.
    */
    function setEpochPeriod(uint256 _period) public onlyOperator {
        // to decrease the epoch period, we need to check if current chain head already exceed the window:
        // lastBlockEpoch + _newPeriod, if so, the _newPeriod cannot be applied since the finalization of current epoch
        // at finalize function will never be triggered, in such case, operator need to find better timing to do so.
        if (_period < config.protocol.epochPeriod) {
            if (block.number >= lastEpochBlock + _period) {
                revert("current chain head exceed the window: lastBlockEpoch + _newPeriod, try again latter on.");
            }
        }
        config.protocol.epochPeriod = _period;
        config.contracts.accountabilityContract.setEpochPeriod(_period);
        emit EpochPeriodUpdated(_period);
    }

    /*
    * @notice Set the Operator account. Restricted to the Operator account.
    * @param _account the new operator account.
    */
    function setOperatorAccount(address _account) public onlyOperator {
        config.protocol.operatorAccount = _account;
        config.contracts.oracleContract.setOperator(_account);
        config.contracts.acuContract.setOperator(_account);
        config.contracts.supplyControlContract.setOperator(_account);
        config.contracts.stabilizationContract.setOperator(_account);
    }

    /*
    Currently not supported
    * @notice Set the block period. Restricted to the Operator account.
    * @param _period Positive integer.

    function setBlockPeriod(uint256 _period) public onlyOperator {
        config.protocol.blockPeriod = _period;
    }
     */

    /*
    * @notice Set the global treasury account. Restricted to the Operator account.
    * @param _account New treasury account.
    */
    function setTreasuryAccount(address payable _account) public onlyOperator {
        config.policy.treasuryAccount = _account;
    }

    /*
    * @notice Set the treasury fee. Restricted to the Operator account.
    * @param _treasuryFee Treasury fee. Precision TBD.
    */
    function setTreasuryFee(uint256 _treasuryFee) public onlyOperator {
        config.policy.treasuryFee = _treasuryFee;
    }

   /*
    * @notice Set the accountability contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setAccountabilityContract(IAccountability _address) public onlyOperator {
        config.contracts.accountabilityContract = _address;
    }

    /*
    * @notice Set the oracle contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setOracleContract(address payable _address) public onlyOperator {
        config.contracts.oracleContract = IOracle(_address);
        config.contracts.acuContract.setOracle(_address);
        config.contracts.stabilizationContract.setOracle(_address);
    }
    
    /*
    * @notice Set the ACU contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setAcuContract(IACU _address) public onlyOperator {
        config.contracts.acuContract = _address;
    }
    
    /*
    * @notice Set the SupplyControl contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setSupplyControlContract(ISupplyControl _address) public onlyOperator {
        config.contracts.supplyControlContract = _address;
    }
    
    /*
    * @notice Set the Stabilization contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setStabilizationContract(IStabilization _address) public onlyOperator {
        config.contracts.stabilizationContract = _address;
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
        blockEpochMap[block.number] = epochID;
        bool epochEnded = lastEpochBlock + config.protocol.epochPeriod == block.number;

        config.contracts.accountabilityContract.finalize(epochEnded);

        if (epochEnded) {
            _performRedistribution();
            _stakingOperations();
            _applyNewCommissionRates();
            address[] memory voters = computeCommittee();
            config.contracts.oracleContract.setVoters(voters);
            lastEpochBlock = block.number;
            epochID += 1;
            emit NewEpoch(epochID);
        }

        bool newRound = config.contracts.oracleContract.finalize();
        if (newRound) {
            try config.contracts.acuContract.update() {}
            catch {}
        }
        return (contractUpgradeReady, committee);
    }

    /**
    * @notice update the current committee by selecting top staking validators.
    * Restricted to the protocol.
    */
    function computeCommittee() public onlyProtocol returns (address[] memory){
        // Left public for testing purposes.
        require(validatorList.length > 0, "There must be validators");
        return computeCommitteePrecompiled(config.protocol.committeeSize);
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
        return config.protocol.epochPeriod;
    }

    /**
    * @notice Returns the block period.
    */
    function getBlockPeriod() external view returns (uint256) {
        return config.protocol.blockPeriod;
    }

    /**
* @notice Returns the un-bonding period.
    */
    function getUnbondingPeriod() external view returns (uint256) {
        return config.policy.unbondingPeriod;
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
        return config.policy.treasuryAccount;
    }

    /**
     * @notice Returns the current treasury fee.
     */
    function getTreasuryFee() external view returns (uint256) {
        return config.policy.treasuryFee;
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
        return config.protocol.committeeSize;
    }

    /**
    * @return Returns the consensus committee enodes.
    */
    function getCommitteeEnodes() external view returns (string[] memory) {
        uint256 len = committee.length;
        string[] memory committeeNodes = new string[](len);
        for (uint256 i = 0; i < len; i++) {
            Validator storage validator = validators[committee[i].addr];
            committeeNodes[i] = validator.enode;
        }
        return committeeNodes;
    }

    /**
    * @return Returns the minimum gas price.
    * @dev Autonity transaction's gas price must be greater or equal to the minimum gas price.
    */
    function getMinimumBaseFee() external view returns (uint256) {
        return config.policy.minBaseFee;
    }

    /**
     * @notice Returns the current operator account.
    */
    function getOperator() external view returns (address) {
        return config.protocol.operatorAccount;
    }

    /**
    * @notice Returns the current Oracle account.
    */
    function getOracle() external view returns (address) {
        return address(config.contracts.oracleContract);
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
        uint256 key = height * MAX_ROUND + round;
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

    /**
     * @notice Returns epoch associated to the block number.
     * @param _block the input block number.
    */
    function getEpochFromBlock(uint256 _block) external view returns (uint256) {
        return blockEpochMap[_block];
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
        require(config.protocol.operatorAccount == msg.sender, "caller is not the operator");
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
    * @dev Modifier that checks if the caller is the governance operator account.
    * This should be abstracted by a separate smart-contract.
    */
    modifier onlyAccountability {
        require(address(config.contracts.accountabilityContract) == msg.sender, "caller is not the slashing contract");
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
    * @dev Emit a {Rewarded} event for every account that collected rewards.
    */
    function _performRedistribution() internal virtual {
        if (address(this).balance == 0) {
            return;
        }
        uint256 _amount = address(this).balance;
        // take treasury fee.
        uint256 _treasuryReward = (config.policy.treasuryFee * _amount) / 10 ** 18;
        if (_treasuryReward > 0) {
            // Using "call" to let the treasury contract do any kind of computation on receive.
            (bool sent,) = config.policy.treasuryAccount.call{value: _treasuryReward}("");
            if (sent == true) {
                _amount -= _treasuryReward;
            }
        }
        // Redistribute fees through the Liquid Newton contract
        totalRedistributed += _amount;
        for (uint256 i = 0; i < committee.length; i++) {
            Validator storage _val = validators[committee[i].addr];
            // votingPower in the committee struct is the amount of bonded-stake pre-slashing event.
            uint256 _reward = (committee[i].votingPower * _amount) / epochTotalBondedStake;
            if (_reward > 0) {
                // committee members in the jailed state were just found guilty in the current epoch.
                // committee members in jailbound state are permanently jailed
                if (_val.state == ValidatorState.jailed || _val.state == ValidatorState.jailbound) {
                    config.contracts.accountabilityContract.distributeRewards{value: _reward}(committee[i].addr);
                    continue;
                }
                // non-jailed validators have a strict amount of bonded newton.
                // the distribution account for the PAS ratio post-slashing.
                uint256 _selfReward = (_val.selfBondedStake * _reward) / _val.bondedStake;
                uint256 _delegationReward = _reward - _selfReward;
                if (_selfReward > 0) {
                    // todo: handle failure scenario here although not critical.
                    _val.treasury.call{value: _selfReward, gas: 2300}("");
                }
                if (_delegationReward > 0) {
                    _val.liquidContract.redistribute{value: _delegationReward}();
                }
                emit Rewarded(_val.nodeAddress, _reward);
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
        (_validator.nodeAddress, _err) = Precompiled.parseEnode(_validator.enode);
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
        for (uint i = 32; i < _multisig.length; i += 65) {
            (r, s, v) = Helpers.extractRSV(_multisig, i);
            signers[i / 65] = ecrecover(hashedData, v, r, s);
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
        require(val.state == ValidatorState.active, "validator must be active");

        val.state = ValidatorState.paused;
        //effectiveBlock may not be accurate if the epoch duration gets modified.
        emit PausedValidator(val.treasury, _address, lastEpochBlock + config.protocol.epochPeriod);
    }


    /**
     * @dev Create a bonding object of `amount` stake token with the `_recipient` address.
     * This object will be processed at epoch finalization.
     *
     * This function assume that `_validator` is a valid validator address.
     */
    function _bond(address _validator, uint256 _amount, address payable _recipient) internal virtual {
        require(_amount > 0, "amount need to be strictly positive");
        require(accounts[_recipient] >= _amount, "insufficient Newton balance");

        accounts[_recipient] -= _amount;
        BondingRequest memory _bonding = BondingRequest(_recipient, _validator, _amount, block.number);
        bondingMap[headBondingID] = _bonding;
        headBondingID++;

        bool _selfBonded = validators[_validator].treasury == _recipient;
        emit NewBondingRequest(_validator, _recipient, _selfBonded, _amount);
    }

    function _applyBonding(uint256 id) internal {
        BondingRequest storage _bonding = bondingMap[id];
        Validator storage _validator = validators[_bonding.delegatee];

        // jailbound validator is jailed permanently, no new bonding can be applied for a jailbound validator
        if (_validator.state == ValidatorState.jailbound) {
            accounts[_bonding.delegator] += _bonding.amount;
            emit BondingRejected(_bonding.delegator, _bonding.delegatee, _bonding.amount);
            return;
        }

        if (_bonding.delegator != _validator.treasury) {
            /* The LNTN: NTN conversion rate is equal to the ratio of issued liquid tokens
             over the total amount of non self-delegated stake tokens. */
            uint256 _liquidAmount;
            uint256 _delegatedStake = _validator.bondedStake - _validator.selfBondedStake;
            if (_delegatedStake == 0) {
                _liquidAmount = _bonding.amount;
            } else {
                _liquidAmount = (_validator.liquidSupply * _bonding.amount) / _delegatedStake;
            }
            _validator.liquidContract.mint(_bonding.delegator, _liquidAmount);
            _validator.liquidSupply += _liquidAmount;
        } else {
            // Penalty Absorbing Stake : No LNTN issued if delegator is treasury
            _validator.selfBondedStake += _bonding.amount;
        }
        _validator.bondedStake += _bonding.amount;
    }

    function _unbond(address _validatorAddress, uint256 _amount, address payable _recipient) internal virtual {
        Validator storage _validator = validators[_validatorAddress];
        bool selfDelegation = _recipient == _validator.treasury;
        if(!selfDelegation) {
            // Lock LNTN if it was issued (non self-delegated stake case)
            uint256 liqBalance = _validator.liquidContract.unlockedBalanceOf(_recipient);
            require(liqBalance >= _amount, "insufficient unlocked Liquid Newton balance");
            _validator.liquidContract.lock(_recipient, _amount);
        } else {
            require(
                _validator.selfBondedStake - _validator.selfUnbondingStakeLocked >= _amount,
                "insufficient self bonded newton balance"
            );
            _validator.selfUnbondingStakeLocked += _amount;
        }
        unbondingMap[headUnbondingID] = UnbondingRequest(_recipient, _validatorAddress, _amount,
                                                         0, block.number, false, selfDelegation);
        headUnbondingID++;

        emit NewUnbondingRequest(_validatorAddress, _recipient, selfDelegation, _amount);
    }

    function _releaseUnbondingStake(uint256 _id) internal virtual {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        if (_unbonding.unbondingShare == 0) {
            return;
        }
        Validator storage _validator = validators[_unbonding.delegatee];
        uint256 _returnedStake;
        if(!_unbonding.selfDelegation){
            _returnedStake =  (_unbonding.unbondingShare *  _validator.unbondingStake) / _validator.unbondingShares;
            _validator.unbondingStake -= _returnedStake;
            _validator.unbondingShares -= _unbonding.unbondingShare;
        } else {
            _returnedStake =  (_unbonding.unbondingShare *  _validator.selfUnbondingStake) / _validator.selfUnbondingShares;
            _validator.selfUnbondingStake -= _returnedStake;
            _validator.selfUnbondingShares -= _unbonding.unbondingShare;
        }
        accounts[_unbonding.delegator] += _returnedStake;
    }

    function _applyUnbonding(uint256 _id) internal virtual {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        Validator storage _validator = validators[_unbonding.delegatee];

        uint256 _newtonAmount;
        if (!_unbonding.selfDelegation){
            // Step 1: Unlock and burn requested liquid newtons
            uint256 _liquidAmount = _unbonding.amount;
            _validator.liquidContract.unlock(_unbonding.delegator, _liquidAmount);
            _validator.liquidContract.burn(_unbonding.delegator, _liquidAmount);

            // Step 2: Calculate the amount of stake to reduce from the delegation pool.
            // Note: validator.liquidSupply cannot be equal to zero here
            uint256 _delegatedStake = _validator.bondedStake - _validator.selfBondedStake;
            _newtonAmount = (_liquidAmount * _delegatedStake) / _validator.liquidSupply;
           _validator.liquidSupply -= _liquidAmount;

            // Step 3: Calculate the amount of shares the staker will get in the unbonding pool.
            // Note : This accounting extra-complication is due to the possibility of slashing unbonding funds.
            if(_validator.unbondingStake == 0) {
                _unbonding.unbondingShare = _newtonAmount;
            } else {
                _unbonding.unbondingShare = (_newtonAmount * _validator.unbondingShares)/_validator.unbondingStake;
            }
            _validator.unbondingStake += _newtonAmount;
            _validator.unbondingShares +=  _unbonding.unbondingShare;
        } else {
            // self-delegated stake path, no LNTN<>NTN conversion
            _newtonAmount = _unbonding.amount;
            if (_newtonAmount > _validator.selfBondedStake) {
                _newtonAmount = _validator.selfBondedStake;
            }
            if (_validator.selfUnbondingStake == 0) {
                 _unbonding.unbondingShare = _newtonAmount;
            } else {
                _unbonding.unbondingShare = (_newtonAmount * _validator.selfUnbondingShares)/_validator.selfUnbondingStake;
            }
            _validator.selfUnbondingStake += _newtonAmount;
            _validator.selfUnbondingShares += _unbonding.unbondingShare;
            // decrease _validator.selfBondedStake for self-delegation
            _validator.selfBondedStake -= _newtonAmount;
            _validator.selfUnbondingStakeLocked -= _unbonding.amount;
        }

        _unbonding.unlocked = true;
        // Final step: Reduce amount of newton bonded
        _validator.bondedStake -= _newtonAmount;
    }

    function _applyNewCommissionRates() internal virtual {
        while (commissionRateChangeQueueFirst < commissionRateChangeQueueLast) {
            // check unbonding period
            CommissionRateChangeRequest storage _curRequest = commissionRateChangeQueue[commissionRateChangeQueueFirst];
            if (_curRequest.startBlock + config.policy.unbondingPeriod > block.number) {
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
    function _stakingOperations() internal virtual {
        // bonding operations are executed first
        for (uint256 i = tailBondingID;
                     i < headBondingID;
                     _applyBonding(i++)){}

        tailBondingID = headBondingID;
        if(tailUnbondingID == headUnbondingID) {
            // everything else already processed, return early
            return;
        }
        // Process the fresh unbonding requests, unbond NTN and burn LNTN
        for (uint256 i = lastUnlockedUnbonding;
                     i < headUnbondingID;
                      _applyUnbonding(i++)){}
        lastUnlockedUnbonding = headUnbondingID;

        // Finally we release the locked NTN tokens
        uint256 _processedId = tailUnbondingID;
        for (uint256 i = tailUnbondingID; i < headUnbondingID; i++) {
            if (unbondingMap[i].requestBlock + config.policy.unbondingPeriod <= block.number) {
                _releaseUnbondingStake(i);
                _processedId += 1;
            } else {
                break;
            }
        }
        tailUnbondingID = _processedId;
    }

    /**
     * @dev Sends necessary slots to precompiled contract.
     * Committee selection and storing the committee and writing it in persistent storage are done in precompiled contract
     */
    function computeCommitteePrecompiled(uint256 _committeeSize) internal returns (address[] memory) {
        require(_committeeSize <= 100, "hardcoded array size 102");
        address[102] memory _returnData;
        address to = Precompiled.COMPUTE_COMMITTEE_CONTRACT;
        uint256 _length = 32*5;
        uint256[5] memory input;
        input[4] = _committeeSize;
        uint _returnDataLength = 64 + _committeeSize*32;
        assembly {
            mstore(input, validatorList.slot)
            mstore(add(input, 0x20), validators.slot)
            mstore(add(input, 0x40), committee.slot)
            mstore(add(input,0x60), epochTotalBondedStake.slot)
            //staticcall(gasLimit, to, inputOffset, inputSize, outputOffset, outputSize)
            if iszero(staticcall(gas(), to, input, _length, _returnData, _returnDataLength)) {
                revert(0, 0)
            }
        }

        require(_returnData[0] == address(1), "unsuccessful call");
        // _returnData[1] has new committee size = length of voters
        require(_returnData[1] != address(0), "sorting unsuccessful");
        if (_committeeSize > uint256(uint160(_returnData[1]))) {
            _committeeSize = uint256(uint160(_returnData[1]));
        }
        address[] memory addresses = new address[](_committeeSize);
        for (uint i = 0; i < _committeeSize; i++) {
            addresses[i] = _returnData[i+2];
        }
        return addresses;
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

    function _inCommittee(address _validator) internal view returns (bool) {
        for (uint256 i = 0; i < committee.length; i++) {
            if (_validator == committee[i].addr) {
                return true;
            }
        }
        return false;
    }

}
