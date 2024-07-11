// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.19;

import "./interfaces/IERC20.sol";
import "./interfaces/IStakeProxy.sol";
import "./interfaces/INonStakableVestingVault.sol";
import "./Liquid.sol";
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

    // TODO (tariq): review the values [already tested from stakable-vesting-contract]
    /**
     * @notice max allowed gas for notifying delegator (contract) about staking operations
     */
    uint256 public maxBondAppliedGas = 50_000;
    uint256 public maxUnbondAppliedGas = 50_000;
    uint256 public maxUnbondReleasedGas = 50_000;
    uint256 public maxRewardsDistributionGas = 10_000;


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
        Liquid liquidContract;
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
    }
    mapping(uint256 => BondingRequest) internal bondingMap;
    uint256 internal tailBondingID;
    uint256 internal headBondingID;

    enum UnbondingReleaseState {
        notReleased,
        released,
        rejected,
        reverted
    }

    struct UnbondingRequest {
        address payable delegator;
        address delegatee;
        uint256 amount; // NTN for self-delegation, LNTN otherwise
        uint256 unbondingShare;
        uint256 requestBlock;
        // new amount of NTN (self-bonding) or LNTN (delegation) in case unbonding was released, but later reverted
        uint256 revertingAmount;
        UnbondingReleaseState state;
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

    struct EpochBoundary{
        uint256 parentEpochBlock;
        uint256 epochBlock;
        uint256 nextEpochBlock;
    }

    Config public config;
    address[] internal validatorList;

    // Stake token state transitions happen every epoch.
    uint256 public epochID;
    mapping(uint256 => uint256) internal blockEpochMap;

    uint256 public parentEpochBlock;
    uint256 public nextEpochBlock;

    // save new epoch period on epoch period update,
    // it is applied to the protocol right after the end of current epoch.
    uint256 public newEpochPeriod;

    uint256 public lastEpochBlock;
    uint256 public lastEpochTime;
    uint256 public epochTotalBondedStake;

    // epochCommittees, saves committee for per epoch in the history.
    // k: epochID, value: committee members
    mapping(uint256=>CommitteeMember[]) internal epochCommittees;
    // epochBoundaries, save epoch boundaries for per epoch
    // k: epochID, value: EpochBoundary.
    mapping(uint256=>EpochBoundary) internal epochBoundaries;

    CommitteeMember[] internal committee;
    uint256 public atnTotalRedistributed;
    uint256 public epochReward;
    string[] internal committeeNodes;
    mapping(address => mapping(address => uint256)) internal allowances;


    /* For callback function at epoch end to notify about staking operations */
    mapping(address => mapping(address => uint256)) isValidatorStaked;
    mapping(address => address[]) internal validatorsStaked;
    mapping(address => uint256) internal contractSaved;
    address[] internal contractAddresses;
    mapping(address => uint256) internal stakingReverted;

    // TODO (tariq): review stakingGasPrice value
    /**
     * @notice the gas price to notify the delegator (only if contract) about the staking operation at epoch end
     */
    uint256 public stakingGasPrice = 1_000_000_000;
    /**
     * @dev stores how much gas given by delegator is left
     */
    mapping(address => uint256) internal gasLeft;

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
    event BondingReverted(address indexed validator, address indexed delegator, uint256 amount);

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
    event UnbondingRejected(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount);
    event AppliedUnbondingReverted(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount);
    event ReleasedUnbondingReverted(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount);

    event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidContract);
    event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock);
    event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock);
    event Rewarded(address indexed addr, uint256 atnAmount, uint256 ntnAmount);
    event EpochPeriodUpdated(uint256 period);
    event NewEpoch(uint256 epoch);

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
        inflationReserve = config.policy.initialInflationReserve;
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

            _verifyEnode(_validators[i]);
            _deployLiquidContract(_validators[i]);

            accounts[_validators[i].treasury] += _bondedStake;
            stakeSupply += _bondedStake;
            _bond(_validators[i].nodeAddress, _bondedStake, payable(_validators[i].treasury));
        }
    }

    function finalizeInitialization() onlyProtocol nonReentrant public {
        _stakingOperations();
        computeCommittee();
        lastEpochTime = block.timestamp;
        // init epoch info for genesis deployment of contract.
        parentEpochBlock = 0;
        nextEpochBlock = lastEpochBlock + config.protocol.epochPeriod;
        epochCommittees[epochID] = committee;
        epochBoundaries[epochID] = EpochBoundary(parentEpochBlock, block.number, nextEpochBlock);
    }

    /**
     * @notice can be used to send AUT to the contract
     */
    function receiveATN() external payable {
        // do nothing
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
            Liquid(address(0)),      // liquid token contract
            0,                       // liquid token supply
            block.number,            // registration block
            0,                       // total slashed
            0,                       // jail release block
            0,                       // provable faults count
            _consensusKey,           // validator key in bytes
            ValidatorState.active    // state
        );

        _verifyAndRegisterValidator(_val, _signatures);
        emit RegisteredValidator(msg.sender, _val.nodeAddress, _oracleAddress, _enode, address(_val.liquidContract));
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
    function bond(address _validator, uint256 _amount) public payable virtual nonReentrant returns (uint256) {
        require(validators[_validator].nodeAddress == _validator, "validator not registered");
        require(validators[_validator].state == ValidatorState.active, "validator need to be active");
        gasLeft[msg.sender] += msg.value / stakingGasPrice;
        return _bond(_validator, _amount, payable(msg.sender));
    }

    /**
    * @notice Create an unbonding request with the caller as delegator. In case the caller is a contract, it needs
    * to send some gas so autonity can notify the caller about staking operations. In case autonity fails to notify
    * the caller (contract), the applied request is reverted.
    * @param _validator address of the validator to unbond stake to.
    * @param _amount total amount of LNTN (or NTN if self delegated) to unbond.
    */
    function unbond(address _validator, uint256 _amount) public payable virtual nonReentrant returns (uint256) {
        require(validators[_validator].nodeAddress == _validator, "validator not registered");
        require(_amount > 0, "unbonding amount is 0");
        gasLeft[msg.sender] += msg.value / stakingGasPrice;
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
        emit ActivatedValidator(_val.treasury, _address, lastEpochBlock + config.protocol.epochPeriod);
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
     * @notice sets the value of max allowed gas for notifying delegator about staking operations
     * NOTE: before updating, please check if the updated value works. It can be checked by updatting
     * the hardcoded value of requiredGasBond and then compiling the contracts and running the tests
     * in stakable_vesting_test.go
     */
    function setMaxBondAppliedGas(uint256 _gas) public onlyOperator {
        maxBondAppliedGas = _gas;
    }

    function setMaxUnbondAppliedGas(uint256 _gas) public onlyOperator {
        maxUnbondAppliedGas = _gas;
    }

    function setMaxUnbondReleasedGas(uint256 _gas) public onlyOperator {
        maxUnbondReleasedGas = _gas;
    }

    function setMaxRewardsDistributionGas(uint256 _gas) public onlyOperator {
        maxRewardsDistributionGas = _gas;
    }

    /**
     * @notice Set gas price for notification on staking operation
     */
    function setStakingGasPrice(uint256 _price) public virtual onlyOperator {
        stakingGasPrice = _price;
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
        // to decrease the epoch period, we need to check if current chain head already exceed the window:
        // lastBlockEpoch + _newPeriod, if so, the _newPeriod cannot be applied since the finalization of current epoch
        // at finalize function will never be triggered, in such case, operator need to find better timing to do so.
        if (_period < config.protocol.epochPeriod) {
            if (block.number >= lastEpochBlock + _period) {
                revert("current chain head exceed the window: lastBlockEpoch + _newPeriod, try again latter on.");
            }
        }

        newEpochPeriod = _period;
        emit EpochPeriodUpdated(_period);
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
    * @return committee The next epoch's consensus committee, if there is no epoch rotation, an empty set is returned.
    * @return parentEpochBlock The parent epoch block number.
    * @return nextEpochBlock The next epoch block number.
    */
    function finalize() external virtual onlyProtocol nonReentrant returns (bool, CommitteeMember[] memory, uint256, uint256) {
        blockEpochMap[block.number] = epochID;
        // while truffle testing, we might ran into situations were currentHeight > lastEpochBlock + epochPeriod, thus
        // making it have no chance to end the epoch as we apply epoch period at the end of epoch. Thus, we change the
        // condition of epochEnded with: lastEpochBlock + config.protocol.epochPeriod <= block.number; to let the
        // endEpoch() have a chance to manipulate the finalize of epoch at the test contract of Autonity.
        // bool epochEnded = lastEpochBlock + config.protocol.epochPeriod == block.number; Change this condition does
        // not impact the protocol since blocks are finalized one by one.
        bool epochEnded = lastEpochBlock + config.protocol.epochPeriod <= block.number;

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
            _notifyRewardsDistribution();
            _stakingOperations();
            _removeContractAddresses();
            _applyNewCommissionRates();

            address[] memory _voters = computeCommittee();
            config.contracts.oracleContract.setVoters(_voters);
            parentEpochBlock = lastEpochBlock;
            lastEpochBlock = block.number;

            // apply new epoch period.
            if (config.protocol.epochPeriod != newEpochPeriod && newEpochPeriod != 0) {
                config.protocol.epochPeriod = newEpochPeriod;
                config.contracts.accountabilityContract.setEpochPeriod(newEpochPeriod);
            }

            nextEpochBlock = lastEpochBlock + config.protocol.epochPeriod;
            lastEpochTime = block.timestamp;
            epochID += 1;
            epochCommittees[epochID] = committee; // save new committee with its new epoch id.
            epochBoundaries[epochID] = EpochBoundary(parentEpochBlock, block.number, nextEpochBlock);
            emit NewEpoch(epochID);
        }

        bool newRound = config.contracts.oracleContract.finalize();
        if (newRound) {
            try config.contracts.acuContract.update() {}
            catch {}
        }

        // return newly elected committee set on epoch rotation,
        // otherwise return empty committee set for no changes.
        if (epochEnded) {
            return (contractUpgradeReady, committee, parentEpochBlock, nextEpochBlock);
        }
        CommitteeMember[] memory emptySet;
        return (contractUpgradeReady, emptySet, parentEpochBlock, nextEpochBlock);
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

    function getUnbondingReleaseState(uint256 _unbondingID) external view returns (UnbondingReleaseState) {
        return unbondingMap[_unbondingID].state;
    }

    /**
     * @notice Returns the amount of LNTN or NTN bonded when the released unbonding was reverted
     */

    function getRevertingAmount(uint256 _unbondingID) external view returns (uint256) {
        require(unbondingMap[_unbondingID].state == UnbondingReleaseState.reverted, "unbonding release not reverted");
        return unbondingMap[_unbondingID].revertingAmount;
    }

    /**
    * @notice Returns the epoch period.
    */
    function getEpochPeriod() external view virtual returns (uint256) {
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
        return lastEpochBlock;
    }

    /**
    * @notice Returns the current contract version.
    */
    function getVersion() external view virtual returns (uint256) {
        return config.contractVersion;
    }

    /**
    * @notice Returns the current epoch info: its commitee and the responding epoch boundary.
    */
    function getEpochInfo() external view virtual returns (CommitteeMember[]memory, EpochBoundary memory) {
        return (epochCommittees[epochID], epochBoundaries[epochID]);
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
        if (_height == 0) {
            return epochCommittees[0];
        }
        uint256 eID = blockEpochMap[_height];

        // if current epoch haven't been ended, then return current committee.
        CommitteeMember[] memory members = epochCommittees[eID];
        if (members.length == 0) {
            return committee;
        }
        // return finalized epoch's committee
        return members;
    }

    /**
    * @notice getProposer returns the address of the proposer for the given height and
    * round. The proposer is selected from the committee via weighted random
    * sampling, with selection probability determined by the voting power of
    * each committee member. The selection mechanism is deterministic and will
    * always select the same address, given the same height, round and contract
    * state.
    */
    function getProposer(uint256 height, uint256 round) external view virtual returns (address) {
        CommitteeMember[] memory members = getCommitteeByHeight(height);

        // calculate total voting power from the corresponding committee, the system does not allow validator with 0 stake/power.
        uint256 total_voting_power = 0;
        for (uint256 i = 0; i < members.length; i++) {
            total_voting_power += members[i].votingPower;
        }

        require(total_voting_power != 0, "The committee is not staking");

        // distribute seed into a 256bits key-space.
        uint256 key = height * MAX_ROUND + round;
        uint256 value = uint256(keccak256(abi.encodePacked(key)));
        uint256 index = value % total_voting_power;

        // find the index hit which committee member which line up in the committee list.
        // we assume there is no 0 stake/power validators.
        uint256 counter = 0;
        for (uint256 i = 0; i < members.length; i++) {
            counter += members[i].votingPower;
            if (index <= counter - 1) {
                return members[i].addr;
            }
        }
        revert("There is no validator left in the network");
    }

    /**
     * @notice Returns epoch associated to the block number.
     * @param _block the input block number.
    */
    function getEpochFromBlock(uint256 _block) external view virtual returns (uint256) {
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
                    // todo: handle failure scenario here although not critical.
                    _val.treasury.call{value: _atnSelfReward, gas: 2300}("");
                }
                uint256 _ntnSelfReward = (_val.selfBondedStake * _ntnReward) / _val.bondedStake;
                if (_ntnSelfReward > 0) {
                    _transfer(address(this), _val.treasury, _ntnSelfReward);
                }
                uint256 _ntnDelegationReward = _ntnReward - _ntnSelfReward;
                uint256 _atnDelegationReward = _atnReward - _atnSelfReward;
                if (_atnDelegationReward > 0 || _ntnDelegationReward > 0) {
                    _transfer(address(this), address(_val.liquidContract), _ntnDelegationReward);
                    _val.liquidContract.redistribute{value: _atnDelegationReward}(_ntnDelegationReward);
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

    function _deployLiquidContract(Validator memory _validator) internal virtual {
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

    function _isContract(address _to) private view returns (bool) {
        uint size;
        assembly {
            size := extcodesize(_to)
        }
        return size > 0;
    }
    
    function _storeAddress(address _delegator, address _validator) private {
        if (contractSaved[_delegator] == 0) {
            contractSaved[_delegator] = 1;
            contractAddresses.push(_delegator);
        }

        if (isValidatorStaked[_delegator][_validator] == 0) {
            isValidatorStaked[_delegator][_validator] = 1;
            validatorsStaked[_delegator].push(_validator);
        }
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
        BondingRequest memory _bonding = BondingRequest(_recipient, _validator, _amount, block.number);
        bondingMap[headBondingID] = _bonding;
        headBondingID++;

        bool _selfBonded = validators[_validator].treasury == _recipient;
        emit NewBondingRequest(_validator, _recipient, _selfBonded, _amount);
        // if the delegator is a contract, store the address to notify it about rewards distribution at epoch end
        if (_isContract(_recipient)) {
            _storeAddress(_recipient, _validator);
        }
        return headBondingID-1;
    }

    /**
     * @dev If the _delegator is a contract, then notify it if bonding request was applied or rejected.
     * Use limited gas to notify the contract. In case the call reverts, revert the staking operation as well.
     * But in case the operation is already rejected (`_rejected = true`), no need to revert anything, because
     * no operation was applied.
     */
    function _notifyBondingApplied(uint256 _id, uint256 _liquid, bool _selfDelegation, bool _rejected) private {
        BondingRequest storage _bonding = bondingMap[_id];
        address _delegator = _bonding.delegator;
        if (!_isContract(_delegator)) {
            return;
        }
        uint256 _gasAllowed = gasLeft[_delegator];
        if (_gasAllowed > maxBondAppliedGas) {
            _gasAllowed = maxBondAppliedGas;
        }
        uint256 _gasUsed = gasleft();
        try IStakeProxy(_delegator).bondingApplied{gas: _gasAllowed}(_id, _bonding.delegatee, _liquid, _selfDelegation, _rejected) {
            _gasUsed -= gasleft();
        } catch {
            _gasUsed -= gasleft();
            if (!_rejected) {
                _revertBonding(_id);
            }
        }
        if (gasLeft[_delegator] > _gasUsed) {
            gasLeft[_delegator] -= _gasUsed;
        }
        else {
            delete gasLeft[_delegator];
        }
    }

    /**
     * @dev bonding request was applied successfully, but couldn't notify the delegator.
     * so we need to revert the applied bonding
     */
    function _revertBonding(uint256 _id) internal virtual {
        BondingRequest storage _bonding = bondingMap[_id];
        accounts[_bonding.delegator] += _bonding.amount;
        Validator storage _validator = validators[_bonding.delegatee];
        // assuming that the bonding request was applied successfully, so the validator must be active
        if (_bonding.delegator != _validator.treasury) {
            // delegatedStake cannot be 0 because the bonding was applied successfully
            // calculate LNTN using current ratio of NTN:LNTN
            uint256 _liquidAmount = _validator.liquidSupply * _bonding.amount / (_validator.bondedStake - _validator.selfBondedStake);
            _validator.liquidContract.burn(_bonding.delegator, _liquidAmount);
            _validator.liquidSupply -= _liquidAmount;
        } else {
            _validator.selfBondedStake -= _bonding.amount;
        }
        _validator.bondedStake -= _bonding.amount;
        emit BondingReverted(_bonding.delegatee, _bonding.delegator, _bonding.amount);
    }

    function _applyBonding(uint256 id) internal virtual {
        BondingRequest storage _bonding = bondingMap[id];
        Validator storage _validator = validators[_bonding.delegatee];

        // no new bonding can be applied for jailbound or jailed or paused validator
        // in case delegator couldn't be notified about rewards distribution, we reject bonding request
        if (_validator.state != ValidatorState.active || stakingReverted[_bonding.delegator] == 1) {
            accounts[_bonding.delegator] += _bonding.amount;
            emit BondingRejected(_bonding.delegatee, _bonding.delegator, _bonding.amount, _validator.state);
            _notifyBondingApplied(id, 0, _bonding.delegator == _validator.treasury, true);
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
            _validator.bondedStake += _bonding.amount;
            _notifyBondingApplied(id, _liquidAmount, false, false);
        } else {
            // Penalty Absorbing Stake : No LNTN issued if delegator is treasury
            _validator.selfBondedStake += _bonding.amount;
            _validator.bondedStake += _bonding.amount;
            _notifyBondingApplied(id, 0, true, false);
        }
    }

    function _unbond(address _validatorAddress, uint256 _amount, address payable _recipient) internal virtual returns (uint256) {
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
        unbondingMap[headUnbondingID] = UnbondingRequest(
            _recipient, _validatorAddress, _amount, 0, block.number, 0, UnbondingReleaseState.notReleased, false, selfDelegation
        );
        headUnbondingID++;

        emit NewUnbondingRequest(_validatorAddress, _recipient, selfDelegation, _amount);
        // if the delegator is a contract, store the address to notify it about rewards distribution at epoch end
        if (_isContract(_recipient)) {
            _storeAddress(_recipient, _validatorAddress);
        }
        return headUnbondingID-1;
    }

    /**
     * @dev notify the delegator (only if contract) if unbonding release was successful or rejected
     */
    function _notifyUnbondingReleased(uint256 _id, uint256 _amount, bool _rejected) private {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        address _delegator = _unbonding.delegator;
        if (!_isContract(_delegator)) {
            return;
        }
        uint256 _gasAllowed = gasLeft[_delegator];
        if (_gasAllowed > maxUnbondReleasedGas) {
            _gasAllowed = maxUnbondReleasedGas;
        }
        uint256 _gasUsed = gasleft();
        try IStakeProxy(_delegator).unbondingReleased{gas: _gasAllowed}(_id, _amount, _rejected) {
            _gasUsed -= gasleft();
        } catch {
            // failed to notify
            _gasUsed -= gasleft();
            if (!_rejected) {
                // we released successfully, but failed to notify. need to revert
                _revertReleasedUnbonding(_id, _amount);
            }
        }
        if (gasLeft[_delegator] > _gasUsed) {
            gasLeft[_delegator] -= _gasUsed;
        }
        else {
            delete gasLeft[_delegator];
        }
    }

    /**
     * @dev in case the release was successful but we couldn't notify the delegator (only if contract), we revert the release.
     * We know _amount NTN was released, so we bond _amount NTN again using current NTN:LNTN ratio
     */
    function _revertReleasedUnbonding(uint256 _id, uint256 _amount) private {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        _unbonding.state = UnbondingReleaseState.reverted;
        emit ReleasedUnbondingReverted(_unbonding.delegatee, _unbonding.delegator, _unbonding.selfDelegation, _amount);
        if (_amount == 0) {
            return;
        }
        Validator storage _validator = validators[_unbonding.delegatee];
        accounts[_unbonding.delegator] -= _amount;
        if (!_unbonding.selfDelegation) {
            // calculate LNTN amount
            uint256 _liquidAmount;
            uint256 _delegatedStake = _validator.bondedStake - _validator.selfBondedStake;
            if (_delegatedStake == 0) {
                _liquidAmount = _amount;
            } else {
                _liquidAmount = (_validator.liquidSupply * _amount) / _delegatedStake;
            }
            _validator.liquidContract.mint(_unbonding.delegator, _liquidAmount);
            _validator.liquidSupply += _liquidAmount;
            _unbonding.revertingAmount = _liquidAmount;
        } else {
            _unbonding.revertingAmount = _amount;
            _validator.selfBondedStake += _amount;
        }
        _validator.bondedStake += _amount;
    }

    function _releaseUnbondingStake(uint256 _id) internal virtual {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        if (!_unbonding.unlocked) {
            // unbonding request was either rejected or reverted, in any case we reject release
            _unbonding.state = UnbondingReleaseState.rejected;
            _notifyUnbondingReleased(_id, 0, true);
            return;
        }
        _unbonding.state = UnbondingReleaseState.released;
        if (_unbonding.unbondingShare == 0) {
            _notifyUnbondingReleased(_id, 0, false);
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
        _notifyUnbondingReleased(_id, _returnedStake, false);
    }

    /**
     * @dev notify the delegator (if it is a contract) that unbonding was applied or rejected
     */
    function _notifyUnbondingApplied(uint256 _id, bool _rejected) private {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        address _delegator = _unbonding.delegator;
        if (!_isContract(_delegator)) {
            return;
        }
        uint256 _gasAllowed = gasLeft[_delegator];
        if (_gasAllowed > maxUnbondAppliedGas) {
            _gasAllowed = maxUnbondAppliedGas;
        }
        uint256 _gasUsed = gasleft();
        try IStakeProxy(_delegator).unbondingApplied{gas: _gasAllowed}(_id, _unbonding.delegatee, _rejected) {
            _gasUsed -= gasleft();
        } catch {
            // failed to notify
            _gasUsed -= gasleft();
            if (!_rejected) {
                // request was applied successfully, but failed to notify, so we need to revert it
                _revertAppliedUnbonding(_id);
            }
        }
        if (gasLeft[_delegator] > _gasUsed) {
            gasLeft[_delegator] -= _gasUsed;
        }
        else {
            delete gasLeft[_delegator];
        }
    }

    /**
     * @dev in case the unbonding request came from a contract, and we applied the request successfully but couldn't
     * notify the contract, we have to revert the applied request
     */
    function _revertAppliedUnbonding(uint256 _id) private {
        // assuming unbonding was applied successfully
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        Validator storage _validator = validators[_unbonding.delegatee];

        uint256 _newtonAmount;
        if (!_unbonding.selfDelegation){
            uint256 _liquidAmount = _unbonding.amount;
            _validator.liquidContract.mint(_unbonding.delegator, _liquidAmount);
            _validator.liquidSupply += _liquidAmount;
            // calculate newton amount from unbonding share
            _newtonAmount = _unbonding.unbondingShare *  _validator.unbondingStake / _validator.unbondingShares;
            _validator.unbondingStake -= _newtonAmount;
            _validator.unbondingShares -=  _unbonding.unbondingShare;
        } else {
            // self-delegated stake path, no LNTN<>NTN conversion
            _newtonAmount = _unbonding.unbondingShare *  _validator.selfUnbondingStake / _validator.selfUnbondingShares;
            _validator.selfUnbondingStake -= _newtonAmount;
            _validator.selfUnbondingShares -= _unbonding.unbondingShare;
            _validator.selfBondedStake += _newtonAmount;
        }

        _unbonding.unbondingShare = 0;
        _unbonding.unlocked = false;
        _validator.bondedStake += _newtonAmount;
        emit AppliedUnbondingReverted(_unbonding.delegatee, _unbonding.delegator, _unbonding.selfDelegation, _unbonding.amount);
    }

    function _applyUnbonding(uint256 _id) internal virtual {
        UnbondingRequest storage _unbonding = unbondingMap[_id];
        // in case delegator is a contract and we failed to notify the contract about rewards distribution,
        // then we cannot notify the contract about _applyUnbonding. So we reject the unbonding request
        if (stakingReverted[_unbonding.delegator] == 1) {
            emit UnbondingRejected(_unbonding.delegatee, _unbonding.delegator, _unbonding.selfDelegation, _unbonding.amount);
            _notifyUnbondingApplied(_id, true);
            return;
        }
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
        _notifyUnbondingApplied(_id, false);
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
     * @dev notify about all the contracts that made bonding or unbonding requests in current epoch
     * about rewards distribution. Because when request are applied in _applyBonding and _applyUnbonding,
     * liquid balances change, and rewads should be known before changing liquid balances.
     * NOTE: It is not necessary to notify about rewards distribution when unbonding is released,
     * because liquid balances do not change in release
     */
    function _notifyRewardsDistribution() private {
        uint256 _length = contractAddresses.length;
        for (uint256 _contractIdx = 0; _contractIdx < _length; _contractIdx++) {
            address _contract = contractAddresses[_contractIdx];
            address[] memory _validators = validatorsStaked[_contract];
            for (uint256 _validatorIdx = 0; _validatorIdx < _validators.length; _validatorIdx++) {
                delete isValidatorStaked[_contract][_validators[_validatorIdx]];
            }
            delete validatorsStaked[_contract];
            uint256 _gasAllowed = gasLeft[_contract];
            if (_gasAllowed > maxRewardsDistributionGas*_validators.length) {
                _gasAllowed = maxRewardsDistributionGas*_validators.length;
            }
            uint256 _gasUsed = gasleft();
            try IStakeProxy(_contract).rewardsDistributed{gas: _gasAllowed}(_validators) {
                _gasUsed -= gasleft();
            } catch {
                _gasUsed -= gasleft();
                stakingReverted[_contract] = 1;
            }
            if (gasLeft[_contract] > _gasUsed) {
                gasLeft[_contract] -= _gasUsed;
            }
            else {
                delete gasLeft[_contract];
            }
        }
    }

    /**
     * @dev remove saved contract addresses that made bonding or unbonding requests in current epoch
     */
    function _removeContractAddresses() private {
        uint256 _length = contractAddresses.length;
        for (uint256 i = 0; i < _length; i++) {
            address _contract = contractAddresses[i];
            delete contractSaved[_contract];
            if (stakingReverted[_contract] == 1) {
                delete stakingReverted[_contract];
            }
        }
        delete contractAddresses;
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

}
