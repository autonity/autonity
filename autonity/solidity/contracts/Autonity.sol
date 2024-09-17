// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.19;

import "./interfaces/INonStakableVestingVault.sol";
import "./liquid/LiquidLogic.sol";
import "./liquid/LiquidState.sol";
import "./Upgradeable.sol";
import "./Precompiled.sol";
import "./Helpers.sol";
import "./UpgradeManager.sol";
import "./lib/BytesLib.sol";
import "./asm/IACU.sol";
import "./asm/ISupplyControl.sol";
import "./asm/IStabilization.sol";
import "./interfaces/IAccountability.sol";
import "./interfaces/IOracle.sol";
import "./interfaces/IAutonity.sol";
import "./interfaces/IInflationController.sol";
import "./interfaces/ILiquidLogic.sol";
import "./ReentrancyGuard.sol";

/** @title Proof-of-Stake Autonity Contract */
enum ValidatorState {active, paused, jailed, jailbound}
uint8 constant DECIMALS = 18;

contract Autonity is IAutonity, IERC20, ReentrancyGuard, Upgradeable {
    uint256 internal constant MAX_ROUND = 99;
    uint256 internal constant CONSENSUS_KEY_LEN = 48;
    uint256 internal constant BLS_PROOF_LEN = 96;
    uint256 internal constant ECDSA_SIGNATURE_LEN = 65;
    uint256 internal constant POP_LEN = 226; // Proof of possession length in bytes. (Enode, OracleNode, ValidatorNode)

    uint256 public constant COMMISSION_RATE_PRECISION = 10_000;


    struct Validator {
        // any change in Validator struct must be synced with offset constants in core/evm/contracts.go
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
        address payable liquidStateContract;
        uint256 liquidSupply;
        uint256 registrationBlock;
        uint256 totalSlashed;
        uint256 jailReleaseBlock;
        uint256 provableFaultCount;
        bytes consensusKey;
        ValidatorState state;
    }

    struct CommitteeMember {
        // any change in Validator struct must be synced with CommitteeSelector code to write committee in DB
        // see CommitteeSelector.updateCommittee function in core/evm/contracts.go
        address addr;
        uint256 votingPower;
        bytes consensusKey;
    }

    /**************************************************/
    // Todo: Create a FIFO structure library, integrate with Staking{}
    /* Used for epoched staking */
    struct BondingRequest {
        address payable delegator;
        address delegatee;
        uint256 amount;
        uint256 requestBlock;
        uint256 liquidMinted;
        uint256 atnRewardsTillBonding;
        uint256 ntnRewardsTillBonding;
        bool rejected;
    }
    mapping(uint256 => BondingRequest) internal bondingMap;
    uint256 internal tailBondingID;
    uint256 internal headBondingID;

    struct UnbondingRequest {
        address payable delegator;
        address delegatee;
        uint256 amount; // NTN for self-delegation, LNTN otherwise
        uint256 unbondingShare;
        uint256 releasedStake;
        uint256 atnRewardsTillUnbonding;
        uint256 ntnRewardsTillUnbonding;
        uint256 requestBlock;
        bool released;
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
        UpgradeManager upgradeManagerContract;
        IInflationController inflationControllerContract;
        INonStakableVestingVault nonStakableVestingContract;
    }

    struct Policy {
        uint256 treasuryFee;
        uint256 minBaseFee;
        uint256 delegationRate;
        uint256 unbondingPeriod;
        uint256 initialInflationReserve;
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

    struct EpochInfo {
        CommitteeMember[] committee;
        uint256 previousEpochBlock;
        uint256 epochBlock;
        uint256 nextEpochBlock;
    }

    Config public config;
    address[] internal validatorList;

    // Stake token state transitions happen every epoch.
    uint256 public epochID;
    mapping(uint256 => uint256) internal blockEpochMap;

    // save new epoch period on epoch period update,
    // it is applied to the protocol right after the end of current epoch.
    uint256 public epochPeriodToBeApplied;

    uint256 public lastFinalizedBlock;
    uint256 public lastEpochTime;
    uint256 public epochTotalBondedStake;

    // epochInfos, save epoch info per epoch in the history
    mapping(uint256=>EpochInfo) internal epochInfos;

    CommitteeMember[] internal committee;
    uint256 public atnTotalRedistributed;
    uint256 public epochReward;
    string[] internal committeeNodes;
    mapping(address => mapping(address => uint256)) internal allowances;


    /* Newton ERC-20. */
    mapping(address => uint256) internal accounts;
    mapping(address => Validator) internal validators;
    uint256 internal stakeSupply;
    uint256 public inflationReserve;

    /*
    We're saving the address of who is deploying the contract and we use it
    for restricting functions that could only be possibly invoked by the protocol
    itself, bypassing transaction processing and signature verification.
    In normal conditions, it is set to the zero address. We're not simply hardcoding
    it only because of testing purposes.
    */
    address public deployer;

    /**
     * @notice Address of the `LiquidLogic` contract. This contract contains all the logic for liquid newton related operations.
     * The state variables are stored in `LiquidState` contract which is different for every validator and is deployed when
     * registering a new validator. To do any operation related to liquid newton, we call `LiquidState` contract of the related
     * validator and that contract does a delegate call to `LiquidLogic` contract.
     */
    address public liquidLogicContract;

    /* Events */
    event MintedStake(address indexed addr, uint256 amount);
    event BurnedStake(address indexed addr, uint256 amount);
    event CommissionRateChange(address indexed validator, uint256 rate);

    /** @notice This event is emitted when a bonding request to a validator node has been registered.
    * This request will only be effective at the end of the current epoch however the stake will be
    * put in custody immediately from the delegator's account.
    * @param validator The validator node account.
    * @param delegator The caller.
    * @param selfBonded True if the validator treasury initiated the request. No LNEW will be issued.
    * @param amount The amount of NEWTON to be delegated.
    */
    event NewBondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount);
    event BondingRejected(address indexed validator, address indexed delegator, uint256 amount, ValidatorState state);

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

    event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidStateContract);
    event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock);
    event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock);
    event Rewarded(address indexed addr, uint256 atnAmount, uint256 ntnAmount);
    event EpochPeriodUpdated(uint256 period, uint256 toBeAppliedAtBlock);
    event NewEpoch(uint256 epoch);

    /**
     * @notice This event is emitted when a call to an address fails in a protocol function (like finalize()).
     * @param to address
     * @param methodSignature method signature of the call, empty in case of plain transaction
     * @param returnData low level return data
     */
    event CallFailed(address to, string methodSignature, bytes returnData);

    /**
     * @dev event to notify the failure in unlocking mechanism of the non-stakable schedules
     */
    event UnlockingScheduleFailed(uint256 epochTime);

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
        epochPeriodToBeApplied = _config.protocol.epochPeriod;
        inflationReserve = config.policy.initialInflationReserve;
        /* We are sharing the same Validator data structure for both genesis
           initialization and runtime. It's not an ideal solution but
           it avoids us adding more complexity to the contract and running into
           stack limit issues.
         */
        liquidLogicContract = address(new LiquidLogic());
        for (uint256 i = 0; i < _validators.length; i++) {
            uint256 _bondedStake = _validators[i].bondedStake;

            // Sanitize the validator fields for a fresh new deployment.
            _validators[i].liquidSupply = 0;
            _validators[i].liquidStateContract = payable(0);
            _validators[i].bondedStake = 0;
            _validators[i].registrationBlock = 0;
            _validators[i].commissionRate = config.policy.delegationRate;
            _validators[i].state = ValidatorState.active;
            _validators[i].selfUnbondingStakeLocked = 0;

            _verifyEnode(_validators[i]);
            _deployLiquidStateContract(_validators[i]);

            accounts[_validators[i].treasury] += _bondedStake;
            stakeSupply += _bondedStake;
            _bond(_validators[i].nodeAddress, _bondedStake, payable(_validators[i].treasury));
        }
    }

    function finalizeInitialization() onlyProtocol nonReentrant public {
        _stakingOperations();
        computeCommittee();
        lastEpochTime = block.timestamp;
        lastFinalizedBlock = block.number;
        // init the 1st epoch info for the protocol with epochID 0 and its corresponding boundary.
        blockEpochMap[block.number] = 0;
        _addEpochInfo(epochID, EpochInfo(committee, 0, block.number, config.protocol.epochPeriod));
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
    function name() external virtual pure returns (string memory) {
        return "Newton";
    }

    /**
    * @return the Stake token's symbol.
    * @dev ERC-20 Optional.
    */
    function symbol() external virtual pure returns (string memory) {
        return "NTN";
    }

    /**
    * @return the number of decimals the NTN token uses.
    * @dev ERC-20 Optional.
    */
    function decimals() public virtual pure returns (uint8) {
        return DECIMALS;
    }

    /**
    * @notice Register a new validator in the system.  The validator might be selected to be part of consensus.
    * This validator will have assigned to its treasury account the caller of this function.
    * A new token "Liquid Stake" is deployed at this phase.
    * @param _enode enode identifying the validator node.
    * @param _oracleAddress identifying the oracle server node that the validator is managing.
    * @param _consensusKey identifying the bls public key in bytes that the validator node is using.
    * @param _signatures is a combination of two ecdsa signatures, and a bls signature as the ownership proof of the
    * validator key appended sequentially. The 1st two ecdsa signatures are in below order:
        1. a message containing treasury account and signed by validator account private key .
        2. a message containing treasury account and signed by Oracle account private key .
    * @dev Emit a {RegisteredValidator} event.
    */
    function registerValidator(string memory _enode, address _oracleAddress, bytes memory _consensusKey, bytes memory _signatures) public virtual {
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
            payable(0), // liquid token contract
            0,                       // liquid token supply
            block.number,            // registration block
            0,                       // total slashed
            0,                       // jail release block
            0,                       // provable faults count
            _consensusKey,           // validator key in bytes
            ValidatorState.active    // state
        );

        _verifyAndRegisterValidator(_val, _signatures);
        emit RegisteredValidator(msg.sender, _val.nodeAddress, _oracleAddress, _enode, _val.liquidStateContract);
    }

    /**
    * @notice Update enode of a registered validator. This function updates the network connection information (IP or/and port)
    of a registered validator. you cannot change the validator's address (pubkey part of the enode)
    * @param _nodeAddress This identifies the validator you want to update
    * @param _enode new enode to be updated
    */
    function updateEnode(address _nodeAddress, string memory _enode) public virtual {
        Validator storage _val = validators[_nodeAddress];
        require(_val.nodeAddress == _nodeAddress, "validator not registered");
        require(_val.treasury == msg.sender, "require caller to be validator treasury account");
        require(!_inCommittee(_nodeAddress), "validator must not be in committee");

        uint _err;
        address _enodePubkey;
        (_enodePubkey, _err) = Precompiled.parseEnode(_enode);
        require(_err == 0, "enode error");

        require(_val.nodeAddress == _enodePubkey, "validator node address can't be updated");
        _val.enode = _enode;
    }

    /**
    * @notice Create a bonding(delegation) request with the caller as delegator. In case the caller is a contract, it needs
    * to send some gas so autonity can notify the caller about staking operations. In case autonity fails to notify
    * the caller (contract), the applied request is reverted.
    * @param _validator address of the validator to delegate stake to.
    * @param _amount total amount of NTN to bond.
    */
    function bond(address _validator, uint256 _amount) public virtual nonReentrant returns (uint256) {
        require(validators[_validator].nodeAddress == _validator, "validator not registered");
        require(validators[_validator].state == ValidatorState.active, "validator need to be active");
        return _bond(_validator, _amount, payable(msg.sender));
    }

    /**
    * @notice Create an unbonding request with the caller as delegator. In case the caller is a contract, it needs
    * to send some gas so autonity can notify the caller about staking operations. In case autonity fails to notify
    * the caller (contract), the applied request is reverted.
    * @param _validator address of the validator to unbond stake to.
    * @param _amount total amount of LNTN (or NTN if self delegated) to unbond.
    */
    function unbond(address _validator, uint256 _amount) public virtual nonReentrant returns (uint256) {
        require(validators[_validator].nodeAddress == _validator, "validator not registered");
        require(_amount > 0, "unbonding amount is 0");
        return _unbond(_validator, _amount, payable(msg.sender));
    }

    /**
    * @notice Pause the validator and stop it accepting delegations.
    * @param _address address to be disabled.
    * @dev emit a {DisabledValidator} event.
    */
    function pauseValidator(address _address) public virtual nonReentrant {
        require(validators[_address].nodeAddress == _address, "validator must be registered");
        require(validators[_address].treasury == msg.sender, "require caller to be validator admin account");
        _pauseValidator(_address);
    }

    /**
    * @notice Re-activate the specified validator.
    * @param _address address to be enabled.
    */
    function activateValidator(address _address) public virtual nonReentrant {
        require(validators[_address].nodeAddress == _address, "validator must be registered");
        Validator storage _val = validators[_address];
        require(_val.treasury == msg.sender, "require caller to be validator treasury account");
        require(_val.state != ValidatorState.active, "validator already active");
        require(!(_val.state == ValidatorState.jailed && _val.jailReleaseBlock > block.number), "validator still in jail");
        require(_val.state != ValidatorState.jailbound, "validator jailed permanently");
        _val.state = ValidatorState.active;
        emit ActivatedValidator(_val.treasury, _address, epochInfos[epochID].nextEpochBlock);
    }

    /**
    * @notice Update the validator. Only accessible to the accountability contract.
    * The difference in bondedStake will go to the treasury account.
    * @param _val Validator to be updated.
    */
    function updateValidatorAndTransferSlashedFunds(Validator calldata _val) external onlyAccountability virtual{
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
    function changeCommissionRate(address _validator, uint256 _rate) public virtual {
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
    function setMinimumBaseFee(uint256 _price) public virtual onlyOperator {
        config.policy.minBaseFee = _price;
        emit MinimumBaseFeeUpdated(_price);
    }

    /*
    * @notice Set the maximum size of the consensus committee. Restricted to the Operator account.
    * @param _size Positive integer.
    */
    function setCommitteeSize(uint256 _size) public virtual onlyOperator {
        require(_size > 0, "committee size can't be 0");
        config.protocol.committeeSize = _size;
    }

    /*
    * @notice Set the unbonding period. Restricted to the Operator account.
    * @param _size Positive integer.
    */
    function setUnbondingPeriod(uint256 _period) public virtual onlyOperator {
        config.policy.unbondingPeriod = _period;
    }

    /*
    * @notice Set the epoch period. Restricted to the Operator account.
    * @param _period Positive integer.
    */
    function setEpochPeriod(uint256 _period) public virtual onlyOperator {
        // todo: shall we need to limit a minimum epoch period?
        // the new epoch period will be activated until current epoch ends.
        epochPeriodToBeApplied = _period;
        uint256 toBeAppliedAtBlock = epochInfos[epochID].nextEpochBlock;
        emit EpochPeriodUpdated(_period, toBeAppliedAtBlock);
    }

    /*
    * @notice Set the Operator account. Restricted to the Operator account.
    * @param _account the new operator account.
    */
    function setOperatorAccount(address _account) public virtual onlyOperator {
        config.protocol.operatorAccount = _account;
        config.contracts.oracleContract.setOperator(_account);
        config.contracts.acuContract.setOperator(_account);
        config.contracts.supplyControlContract.setOperator(_account);
        config.contracts.stabilizationContract.setOperator(_account);
        config.contracts.upgradeManagerContract.setOperator(_account);
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
    function setTreasuryAccount(address payable _account) public virtual onlyOperator {
        config.policy.treasuryAccount = _account;
    }

    /*
    * @notice Set the treasury fee. Restricted to the Operator account.
    * @param _treasuryFee Treasury fee. Precision TBD.
    */
    function setTreasuryFee(uint256 _treasuryFee) public virtual onlyOperator {
        config.policy.treasuryFee = _treasuryFee;
    }

    /*
     * @notice Set the accountability contract address. Restricted to the Operator account.
     * @param _address the contract address
     */
    function setAccountabilityContract(IAccountability _address) public virtual onlyOperator {
        config.contracts.accountabilityContract = _address;
    }

    /*
    * @notice Set the oracle contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setOracleContract(address payable _address) public virtual onlyOperator {
        config.contracts.oracleContract = IOracle(_address);
        config.contracts.acuContract.setOracle(_address);
        config.contracts.stabilizationContract.setOracle(_address);
    }

    /*
    * @notice Set the ACU contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setAcuContract(IACU _address) public virtual onlyOperator {
        config.contracts.acuContract = _address;
    }

    /*
    * @notice Set the SupplyControl contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setSupplyControlContract(ISupplyControl _address) public virtual onlyOperator {
        config.contracts.supplyControlContract = _address;
    }

    /*
    * @notice Set the Stabilization contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setStabilizationContract(IStabilization _address) public virtual onlyOperator {
        config.contracts.stabilizationContract = _address;
    }

    /*
    * @notice Set the Inflation Controller contract address. Restricted to the Operator account.
    * @param _address the contract address
    */
    function setInflationControllerContract(IInflationController _address) public virtual onlyOperator {
        config.contracts.inflationControllerContract = _address;
    }

    /*
    * @notice Set the Upgrade Manager contract address. Restricted to the Operator account.
    * It is only meant to be used for internal testing purposes. Anything different than
    * 0x3C368B86AF00565Df7a3897Cfa9195B9434A59f9 will break the upgrade function live!
    * @param _address the contract address
    */
    function setUpgradeManagerContract(UpgradeManager _address) public virtual onlyOperator {
        config.contracts.upgradeManagerContract = _address;
    }

    /**
     * @notice Set the Non-stakable Vesting contract address.
     */
    function setNonStakableVestingContract(INonStakableVestingVault _address) public virtual onlyOperator {
        config.contracts.nonStakableVestingContract = _address;
    }

    /**
     * @notice Set address of the liquid logic contact.
     * @custom:restricted-to operator account
     */
    function SetLiquidLogicContract(address _contract) public virtual onlyOperator {
        require(_contract != address(0), "invalid contract address for liquid logic");
        liquidLogicContract = _contract;
    }

    /*
    * @notice Mint new stake token (NTN) and add it to the recipient balance. Restricted to the Operator account.
    * @dev emit a MintStake event.
    */
    function mint(address _addr, uint256 _amount) public virtual onlyOperator {
        _mint(_addr, _amount);
    }

    /**
    * @notice Burn the specified amount of NTN stake token from an account. Restricted to the Operator account.
    * This won't burn associated Liquid tokens.
    */
    function burn(address _addr, uint256 _amount) public virtual onlyOperator {
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
    function transfer(address _recipient, uint256 _amount) external virtual override returns (bool) {
        _transfer(msg.sender, _recipient, _amount);
        emit Transfer(msg.sender, _recipient, _amount);
        return true;
    }

    /**
     * @dev See {IERC20-approve}.
     *
     * Requirements:
     *
     * - `spender` cannot be the zero address.
     */
    function approve(address spender, uint256 amount) external virtual override returns (bool) {
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
    function transferFrom(address _sender, address _recipient, uint256 _amount) external virtual override returns (bool){
        //TODO URGENT require(allowances[_sender][msg.sender] > 0, "no allowance");
        //require(allowances[_sender][msg.sender] >= _amount, "unsufficient allowance");
        _transfer(_sender, _recipient, _amount);
        uint256 newAllowance = allowances[_sender][msg.sender] - _amount;
        _approve(_sender, msg.sender, newAllowance);
        emit Transfer(_sender, _recipient, _amount);
        return true;
    }

    /**
    * @dev See {IERC20-allowance}.
    */
    function allowance(address owner, address spender) external view virtual override returns (uint256) {
        return allowances[owner][spender];
    }

    /** @dev finalize is the block state finalisation function. It is called
    * each block after processing every transactions within it. It must be restricted to the
    * protocol only.
    *
    * @return upgrade Set to true if an autonity contract upgrade is available.
    * @return epochEnded Set to true if an epoch is ended.
    * @return committee The next epoch's consensus committee, if there is no epoch rotation, an empty set is returned.
    * @return previousEpochBlock The previous epoch block number.
    * @return nextEpochBlock The next epoch block number.
    */
    function finalize() external virtual onlyProtocol nonReentrant returns (bool, bool, CommitteeMember[] memory, uint256, uint256) {
        lastFinalizedBlock = block.number;
        blockEpochMap[block.number] = epochID;

        // Making this condition change to make the Autonity contract to have a chance to
        // finish the epoch rotation in the truffle test context, as in truffle test the
        // chain height might go beyond the lastEpochBlock+EpochPeriod making the epoch rotation
        // never happen as the new epoch period is going to be activated on epoch rotation.
        bool epochEnded = epochInfos[epochID].nextEpochBlock <= block.number;
        config.contracts.accountabilityContract.finalize(epochEnded);

        if (epochEnded) {
            // We first calculate the new NTN injected supply for this epoch
            uint256 _inflationReward = config.contracts.inflationControllerContract.calculateSupplyDelta(
                stakeSupply,
                inflationReserve,
                lastEpochTime,
                block.timestamp
            );
            if (inflationReserve < _inflationReward){
                // If this code path is taken there is something deeply wrong happening in the inflation controller
                // contract.
                _inflationReward = inflationReserve;
            }
            // mint inflation NTN with the AC recipient
            // all rewards belong to the Autonity Contract before redistribution.
            _mint(address(this), _inflationReward);
            inflationReserve -= _inflationReward;
            try config.contracts.nonStakableVestingContract.unlockTokens() returns (uint256 _newUnlockedSubscribed, uint256 _newUnlockedUnsubscribed) {
                // mint unsubsribed tokens to treasury account
                _mint(config.policy.treasuryAccount, _newUnlockedUnsubscribed);
                // and the subsribed tokens to the vault of non-stakable vesting contract
                _mint(address(config.contracts.nonStakableVestingContract), _newUnlockedSubscribed);
            } catch {
                // need immediate attention
                emit UnlockingScheduleFailed(block.timestamp);
            }
            // redistribute ATN tx fees and newly minted NTN inflation reward
            _performRedistribution(address(this).balance, _inflationReward);
            // end of epoch here
            _stakingOperations();
            _applyNewCommissionRates();

            address[] memory _voters = computeCommittee();
            config.contracts.oracleContract.setVoters(_voters);

            uint256 previousEpochBlock = epochInfos[epochID].epochBlock;
            // apply new epoch period.
            if (config.protocol.epochPeriod != epochPeriodToBeApplied && epochPeriodToBeApplied != 0) {
                config.protocol.epochPeriod = epochPeriodToBeApplied;
                config.contracts.accountabilityContract.setEpochPeriod(epochPeriodToBeApplied);
            }
            uint256 nextEpochBlock = block.number + config.protocol.epochPeriod;
            lastEpochTime = block.timestamp;
            epochID += 1;
            _addEpochInfo(epochID, EpochInfo(committee, previousEpochBlock, block.number, nextEpochBlock));
            emit NewEpoch(epochID);
        }

        bool newRound = config.contracts.oracleContract.finalize();
        if (newRound) {
            try config.contracts.acuContract.update() {}
            catch {}
        }

        return (contractUpgradeReady, epochEnded, committee, epochInfos[epochID].previousEpochBlock, epochInfos[epochID].nextEpochBlock);
    }

    /**
    * @notice update the current committee by selecting top staking validators.
    * Restricted to the protocol.
    */
    function computeCommittee() public virtual onlyProtocol returns (address[] memory){
        // Left public for testing purposes.
        require(validatorList.length > 0, "There must be validators");
        uint256[5] memory input;
        input[4] = config.protocol.committeeSize;
        assembly {
            mstore(input, validatorList.slot)
            mstore(add(input, 0x20), validators.slot)
            mstore(add(input, 0x40), committee.slot)
            mstore(add(input,0x60), epochTotalBondedStake.slot)
        }
        Precompiled.computeCommitteePrecompiled(input);
        // get oracle address of committee members
        // calculate committeeNodes
        delete committeeNodes;
        uint256 committeeSize = committee.length;
        require(committeeSize > 0, "committee is empty");
        address[] memory _voters = new address[](committeeSize);
        for (uint i = 0; i < committeeSize; i++) {
            Validator storage _member = validators[committee[i].addr];
            committeeNodes.push(_member.enode);
            _voters[i] = _member.oracleAddress;
        }
        return _voters;
    }

    /*
    ============================================================
        Getters
    ============================================================
    */

    /**
     * @notice Returns the release state of the unbonding request
     */

    function isUnbondingReleased(uint256 _unbondingID) external view returns (bool) {
        return unbondingMap[_unbondingID].released;
    }

    function getReleasedStake(uint256 _unbondingID) external view returns (uint256) {
        return unbondingMap[_unbondingID].releasedStake;
    }

    /**
    * @notice Returns the epoch period.
    */
    function getEpochPeriod() external view virtual returns (uint256) {
        // if the new epoch period haven't being applied yet, return it anyway.
        if (config.protocol.epochPeriod != epochPeriodToBeApplied) {
            return epochPeriodToBeApplied;
        }
        // otherwise we return the current applied epoch period.
        return config.protocol.epochPeriod;
    }

    /**
    * @notice Returns the block period.
    */
    function getBlockPeriod() external view virtual returns (uint256) {
        return config.protocol.blockPeriod;
    }

    /**
     * @notice Returns the un-bonding period.
     */
    function getUnbondingPeriod() external view virtual returns (uint256) {
        return config.policy.unbondingPeriod;
    }

    /**
    * @notice Returns the last epoch's end block height.
    */
    function getLastEpochBlock() external view virtual returns (uint256) {
        return epochInfos[epochID].epochBlock;
        //return lastEpochBlock;
    }

    /**
    * @notice Returns the current contract version.
    */
    function getVersion() external view virtual returns (uint256) {
        return config.contractVersion;
    }

    /**
    * @notice Returns the current epoch info of the chain.
    */
    function getEpochInfo() external view virtual returns (CommitteeMember[] memory, uint256, uint256, uint256) {
        CommitteeMember[] memory members = epochInfos[epochID].committee;
        uint256 previous = epochInfos[epochID].previousEpochBlock;
        uint256 current = epochInfos[epochID].epochBlock;
        uint256 next = epochInfos[epochID].nextEpochBlock;
        return (members, previous, current, next);
    }

    /**
     * @notice Returns the block committee.
     * @return Current block committee if called before finalize(), next block if called after.
     */
    function getCommittee() external view virtual returns (CommitteeMember[] memory) {
        return committee;
    }

    /**
     * @notice Returns the current list of validators.
     */
    function getValidators() external view virtual returns (address[] memory) {
        return validatorList;
    }

    /**
     * @notice Returns the current treasury account.
     */
    function getTreasuryAccount() external view virtual returns (address) {
        return config.policy.treasuryAccount;
    }

    /**
     * @notice Returns the current treasury fee.
     */
    function getTreasuryFee() external view virtual returns (uint256) {
        return config.policy.treasuryFee;
    }

    /**
     * @notice Returns the next epoch block.
     */
    function getNextEpochBlock() external view virtual returns(uint256) {
        return epochInfos[epochID].nextEpochBlock;
    }

    /**
    * @notice Returns the amount of unbonded Newton token held by the account (ERC-20).
    */
    function balanceOf(address _addr) external view virtual override returns (uint256) {
        return accounts[_addr];
    }

    /**
    * @notice Returns the total amount of stake token issued.
    */
    function totalSupply() external view virtual override returns (uint256) {
        return stakeSupply;
    }

    /**
    * @return Returns a user object with the `_account` parameter. The returned data
    * object might be empty if there is no user associated.
    */
    function getValidator(address _addr) external view virtual returns (Validator memory) {
        require(validators[_addr].nodeAddress == _addr, "validator not registered");
        return validators[_addr];
    }

    /**
    * @return Returns the maximum size of the consensus committee.
    */
    function getMaxCommitteeSize() external view virtual returns (uint256) {
        return config.protocol.committeeSize;
    }

    /**
    * @return Returns the consensus committee enodes.
    */
    function getCommitteeEnodes() external view virtual returns (string[] memory) {
        return committeeNodes;
    }

    /**
    * @return Returns the minimum gas price.
    * @dev Autonity transaction's gas price must be greater or equal to the minimum gas price.
    */
    function getMinimumBaseFee() external view virtual returns (uint256) {
        return config.policy.minBaseFee;
    }

    /**
     * @notice Returns the current operator account.
    */
    function getOperator() external view virtual returns (address) {
        return config.protocol.operatorAccount;
    }

    /**
    * @notice Returns the current Oracle account.
    */
    function getOracle() external view virtual returns (address) {
        return address(config.contracts.oracleContract);
    }

    /**
     * @notice Returns the committee of a specific height.
     * @param _height the input block number
     * @return committee The next epoch's consensus committee, if there is no epoch rotation, an empty set is returned.
     */
    function getCommitteeByHeight(uint256 _height) public view virtual returns (CommitteeMember[] memory) {
        require(_height <= block.number, "cannot get committee for a future height");

        // if the block was already finalized, get committee by its corresponding epoch id.
        if (_height <= lastFinalizedBlock) {
            uint256 blockEpochID = blockEpochMap[_height];
            CommitteeMember[] memory members = epochInfos[blockEpochID].committee;
            return members;
        }

        // otherwise, this _height is the latest consensus instance, return current committee.
        return committee;
    }

    /**
     * @notice Returns epoch associated to the block number.
     * @param _block the input block number.
     */
    function getEpochFromBlock(uint256 _block) external view virtual returns (uint256) {
        require(_block <= block.number, "cannot get epoch for a future block");
        if (_block <= lastFinalizedBlock) {
            return blockEpochMap[_block];
        }
        return epochID;
    }

    function isBondingRejected(uint256 _bondingID) external view returns (bool) {
        return bondingMap[_bondingID].rejected;
    }

    function getBondedLiquid(uint256 _bondingID) external view virtual returns (uint256) {
        return bondingMap[_bondingID].liquidMinted;
    }

    function getRewardsTillBonding(uint256 _bondingID) external view returns (uint256, uint256) {
        BondingRequest storage _bonding = bondingMap[_bondingID];
        return (_bonding.atnRewardsTillBonding, _bonding.ntnRewardsTillBonding);
    }

    function getRewardsTillUnbonding(uint256 _bondingID) external view returns (uint256, uint256) {
        UnbondingRequest storage _unbonding = unbondingMap[_bondingID];
        return (_unbonding.atnRewardsTillUnbonding, _unbonding.ntnRewardsTillUnbonding);
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
    * @notice Perform ATN and NTN reward distribution. The rewards fees
    * are simply re-distributed to all stake-holders, including validators,
    * pro-rata the amount of stake held.
    * @dev Emit a {Rewarded} event for every account that collected rewards.
    * @param _atn: Amount of ATN to be redistributed. The source funds will be taken from
    * this contract balance.
    * @param _ntn: Amount of NTN to be redistributed. The source funds will be minted here.
    */
    function _performRedistribution(uint256 _atn, uint256 _ntn) internal virtual {
        // exit early if nothing to redistribute.
        if (_atn == 0 && _ntn == 0) {
            return;
        }
        // Take ATN treasury fee.
        uint256 _atnTreasuryReward = (config.policy.treasuryFee * _atn) / 10 ** 18;
        if (_atnTreasuryReward > 0) {
            // Using "call" to let the treasury contract do any kind of computation on receive.
            (bool sent,) = config.policy.treasuryAccount.call{value: _atnTreasuryReward}("");
            if (sent == true) {
                _atn -= _atnTreasuryReward;
            }
        }
        // Redistribute fees through the Liquid Newton contract
        atnTotalRedistributed += _atn;
        for (uint256 i = 0; i < committee.length; i++) {
            Validator storage _val = validators[committee[i].addr];
            // votingPower in the committee struct is the amount of bonded-stake pre-slashing event.
            uint256 _atnReward = (committee[i].votingPower * _atn) / epochTotalBondedStake;
            uint256 _ntnReward = (committee[i].votingPower * _ntn) / epochTotalBondedStake;
            if (_atnReward > 0 || _ntnReward > 0) {
                // committee members in the jailed state were just found guilty in the current epoch.
                // committee members in jailbound state are permanently jailed
                if (_val.state == ValidatorState.jailed || _val.state == ValidatorState.jailbound) {
                    _transfer(address(this), address(config.contracts.accountabilityContract), _ntnReward);
                    config.contracts.accountabilityContract.distributeRewards{value: _atnReward}(committee[i].addr, _ntnReward);
                    continue;
                }
                // non-jailed validators have a strict amount of bonded newton.
                // the distribution account for the PAS ratio post-slashing.
                uint256 _atnSelfReward = (_val.selfBondedStake * _atnReward) / _val.bondedStake;
                if (_atnSelfReward > 0) {
                    (bool _sent, bytes memory _returnData) = _val.treasury.call{value: _atnSelfReward, gas: 2300}("");
                    // let the treasury know that call failed
                    if (_sent == false) {
                        emit CallFailed(_val.treasury, "", _returnData);
                    }
                }
                uint256 _ntnSelfReward = (_val.selfBondedStake * _ntnReward) / _val.bondedStake;
                if (_ntnSelfReward > 0) {
                    _transfer(address(this), _val.treasury, _ntnSelfReward);
                }
                uint256 _ntnDelegationReward = _ntnReward - _ntnSelfReward;
                uint256 _atnDelegationReward = _atnReward - _atnSelfReward;
                if (_atnDelegationReward > 0 || _ntnDelegationReward > 0) {
                    _transfer(address(this), _val.liquidStateContract, _ntnDelegationReward);
                    ILiquidLogic(_val.liquidStateContract).redistribute{value: _atnDelegationReward}(_ntnDelegationReward);
                }
                // TODO: This has to be reconsidered - I feel it is too expensive
                // to emit an event per validator. But what is our recommend way to track rewards
                // from a user perspective then ?
                emit Rewarded(_val.nodeAddress, _atnReward, _ntnReward);
            }
        }
    }

    // @dev No side effects on this function, so safe to be called in the middle of something (but may revert).
    // We may want to switch to OZ's ERC20 at one point to deal with callbacks
    // but we'll have to deal with re-entrency stuff in this case. For the time being we are conservative.
    function _transfer(address _sender, address _recipient, uint256 _amount) internal virtual {
        require(accounts[_sender] >= _amount, "amount exceeds balance");
        accounts[_sender] -= _amount;
        accounts[_recipient] += _amount;
    }

    function _mint(address _addr, uint256 _amount) internal virtual {
        accounts[_addr] += _amount;
        stakeSupply += _amount;
        emit MintedStake(_addr, _amount);
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

    function _verifyEnode(Validator memory _validator) internal virtual view {
        // _enode can't be empty and needs to be well-formed.
        uint _err;
        (_validator.nodeAddress, _err) = Precompiled.parseEnode(_validator.enode);
        require(_err == 0, "enode error");
        require(validators[_validator.nodeAddress].nodeAddress == address(0), "validator already registered");
        require(_validator.commissionRate <= COMMISSION_RATE_PRECISION, "invalid commission rate");
    }

    function _deployLiquidStateContract(Validator memory _validator) internal virtual {
        if (_validator.liquidStateContract == address(0)) {
            require(liquidLogicContract != address(0), "liquid logic contract not deployed");
            string memory stringLength = Helpers.toString(validatorList.length);
            _validator.liquidStateContract = payable(
                new LiquidState(
                    _validator.nodeAddress,
                    _validator.treasury,
                    _validator.commissionRate,
                    stringLength,
                    liquidLogicContract
                )
            );
        }
        validatorList.push(_validator.nodeAddress);
        validators[_validator.nodeAddress] = _validator;
    }

    function _verifyAndRegisterValidator(Validator memory _validator, bytes memory _signatures) internal virtual {
        require(_signatures.length == POP_LEN, "Invalid proof length");
        require(_validator.oracleAddress == address(uint160(_validator.oracleAddress)), "Invalid oracle address");
        require(_validator.consensusKey.length == CONSENSUS_KEY_LEN, "Invalid consensus key length");

        // verify enode and parse node address
        _verifyEnode(_validator);

        // verify proof of possessions.
        bytes memory prefix = "\x19Ethereum Signed Message:\n";
        bytes memory treasury = abi.encodePacked(_validator.treasury);
        bytes32 hashedData = keccak256(abi.encodePacked(prefix, Helpers.toString(treasury.length), treasury));
        address[] memory signers = new address[](2);
        bytes32 r;
        bytes32 s;
        uint8 v;
        // 1st batch bytes are signatures generated by node key and oracle node key.
        bytes memory ecdsaSignatures = BytesLib.slice(_signatures, 0, ECDSA_SIGNATURE_LEN*2);
        // 2nd batch of rest 96 bytes are the signature generated by validator BLS key.
        bytes memory blsSignature = BytesLib.slice(_signatures, ECDSA_SIGNATURE_LEN*2, BLS_PROOF_LEN);

        //start from 32th byte to skip the encoded length field from the bytes type variable
        for (uint i = 32; i < ecdsaSignatures.length; i += ECDSA_SIGNATURE_LEN) {
            (r, s, v) = Helpers.extractRSV(ecdsaSignatures, i);
            signers[i/ECDSA_SIGNATURE_LEN] = ecrecover(hashedData, v, r, s);
        }
        require(signers[0] == _validator.nodeAddress, "Invalid node key ownership proof provided");
        require(signers[1] == _validator.oracleAddress, "Invalid oracle key ownership proof provided");
        require(Precompiled.popVerification(_validator.consensusKey, blsSignature, _validator.treasury) == Precompiled.SUCCESS,
            "Invalid consensus key ownership proof for registration");

        // all good, now deploy liquidity contract.
        _deployLiquidStateContract(_validator);
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
        emit PausedValidator(val.treasury, _address, epochInfos[epochID].nextEpochBlock);
    }


    /**
     * @dev Create a bonding object of `amount` stake token with the `_recipient` address.
     * This object will be processed at epoch finalization.
     *
     * This function assume that `_validator` is a valid validator address.
     */
    function _bond(address _validator, uint256 _amount, address payable _recipient) internal virtual returns (uint256) {
        require(_amount > 0, "amount need to be strictly positive");
        require(accounts[_recipient] >= _amount, "insufficient Newton balance");

        accounts[_recipient] -= _amount;
        BondingRequest memory _bonding = BondingRequest(_recipient, _validator, _amount, block.number, 0, 0, 0, false);
        bondingMap[headBondingID] = _bonding;
        headBondingID++;

        bool _selfBonded = validators[_validator].treasury == _recipient;
        emit NewBondingRequest(_validator, _recipient, _selfBonded, _amount);
        return headBondingID-1;
    }

    function _applyBonding(uint256 id) internal virtual {
        BondingRequest storage _bonding = bondingMap[id];
        Validator storage _validator = validators[_bonding.delegatee];

        // no new bonding can be applied for jailbound or jailed or paused validator
        // in case delegator couldn't be notified about rewards distribution, we reject bonding request
        if (_validator.state != ValidatorState.active) {
            accounts[_bonding.delegator] += _bonding.amount;
            (_bonding.atnRewardsTillBonding, _bonding.ntnRewardsTillBonding) = _validator.liquidContract.realisedFees(_bonding.delegator);
            _bonding.rejected = true;
            emit BondingRejected(_bonding.delegatee, _bonding.delegator, _bonding.amount, _validator.state);
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
            ILiquidLogic(_validator.liquidStateContract).mint(_bonding.delegator, _liquidAmount);
            _validator.liquidSupply += _liquidAmount;
            _validator.bondedStake += _bonding.amount;
            _bonding.liquidMinted = _liquidAmount;
        } else {
            // Penalty Absorbing Stake : No LNTN issued if delegator is treasury
            _validator.selfBondedStake += _bonding.amount;
            _validator.bondedStake += _bonding.amount;
        }
        (_bonding.atnRewardsTillBonding, _bonding.ntnRewardsTillBonding) = _validator.liquidContract.realisedFees(_bonding.delegator);
    }

    function _unbond(address _validatorAddress, uint256 _amount, address payable _recipient) internal virtual returns (uint256) {
        Validator storage _validator = validators[_validatorAddress];
        bool selfDelegation = _recipient == _validator.treasury;
        if(!selfDelegation) {
            // Lock LNTN if it was issued (non self-delegated stake case)
            uint256 liqBalance = ILiquidLogic(_validator.liquidStateContract).unlockedBalanceOf(_recipient);
            require(liqBalance >= _amount, "insufficient unlocked Liquid Newton balance");
            ILiquidLogic(_validator.liquidStateContract).lock(_recipient, _amount);
        } else {
            require(
                _validator.selfBondedStake - _validator.selfUnbondingStakeLocked >= _amount,
                "insufficient self bonded newton balance"
            );
            _validator.selfUnbondingStakeLocked += _amount;
        }
        unbondingMap[headUnbondingID] = UnbondingRequest(
            _recipient, _validatorAddress, _amount, 0, 0, 0, 0, block.number, false, false, selfDelegation
        );
        headUnbondingID++;

        emit NewUnbondingRequest(_validatorAddress, _recipient, selfDelegation, _amount);
        return headUnbondingID-1;
    }

    function _releaseUnbondingStake(uint256 _id) internal virtual {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        _unbonding.released = true;
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
        _unbonding.releasedStake = _returnedStake;
        accounts[_unbonding.delegator] += _returnedStake;
    }

    function _applyUnbonding(uint256 _id) internal virtual {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        Validator storage _validator = validators[_unbonding.delegatee];

        uint256 _newtonAmount;
        if (!_unbonding.selfDelegation){
            // Step 1: Unlock and burn requested liquid newtons
            uint256 _liquidAmount = _unbonding.amount;
            ILiquidLogic(_validator.liquidStateContract).unlock(_unbonding.delegator, _liquidAmount);
            ILiquidLogic(_validator.liquidStateContract).burn(_unbonding.delegator, _liquidAmount);

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

        (_unbonding.atnRewardsTillUnbonding, _unbonding.ntnRewardsTillUnbonding) = _validator.liquidContract.realisedFees(_unbonding.delegator);
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

            Validator storage _validator = validators[_curRequest.validator];
            _validator.commissionRate = _curRequest.rate;
            ILiquidLogic(_validator.liquidStateContract).setCommissionRate(_curRequest.rate);

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

    function _removeFromArray(address _address, address[] storage _array) internal virtual{
        require(_array.length > 0);
        for (uint256 i = 0; i < _array.length; i++) {
            if (_array[i] == _address) {
                _array[i] = _array[_array.length - 1];
                _array.pop();
                break;
            }
        }
    }

    function _inCommittee(address _validator) internal virtual view returns (bool) {
        for (uint256 i = 0; i < committee.length; i++) {
            if (_validator == committee[i].addr) {
                return true;
            }
        }
        return false;
    }

    function _addEpochInfo(uint256 _epochID, EpochInfo memory _epoch) internal {
        EpochInfo storage epoch = epochInfos[_epochID];
        epoch.previousEpochBlock = _epoch.previousEpochBlock;
        epoch.epochBlock = _epoch.epochBlock;
        epoch.nextEpochBlock = _epoch.nextEpochBlock;
        for (uint256 i=0; i<_epoch.committee.length; i++) {
            epoch.committee.push(_epoch.committee[i]);
        }
    }
}
