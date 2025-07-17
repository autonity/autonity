// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

import {IACU} from "../asm/interfaces/IACU.sol";
import {IAuctioneer} from "../asm/interfaces/IAuctioneer.sol";
import {IStabilization} from "../asm/interfaces/IStabilization.sol";
import {ISupplyControl} from "../asm/interfaces/ISupplyControl.sol";
import {IAccountability} from "./IAccountability.sol";
import {IERC20} from "./IERC20.sol";
import {IInflationController} from "./IInflationController.sol";
import {ILiquid} from "./ILiquid.sol";
import {IOmissionAccountability} from "./IOmissionAccountability.sol";
import {IOracle} from "./IOracle.sol";
import {IScheduleController} from "./IScheduleController.sol";
import {IUpgradeManager} from "./IUpgradeManager.sol";

uint8 constant NTN_DECIMALS = 18;

/**
 * @dev Interface of the Autonity Contract.
 * Import this over Autonity.sol.
 */
interface IAutonity is IERC20, IScheduleController {

    struct Eip1559 {
        uint256 minBaseFee;
        uint256 baseFeeChangeDenominator;
        uint256 elasticityMultiplier;
        uint256 gasLimitBoundDivisor;
    }

    enum ValidatorState {active, paused, jailed, jailbound, jailedForInactivity, jailboundForInactivity}

    // any change in Validator struct must be synced with offset constants in core/vm/contracts.go
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
        ILiquid liquidStateContract;
        uint256 liquidSupply;
        uint256 registrationBlock;
        uint256 totalSlashed;
        uint256 jailReleaseBlock;
        bytes consensusKey;
        ValidatorState state;
        // NOTE: the conversionRatio is not supposed to be used for protocol computations, but rather serves
        // as a way to allow external clients to compute the performance of their delegated stake
        uint256 conversionRatio;
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

    struct UnbondingRequest {
        address payable delegator;
        address delegatee;
        uint256 amount; // NTN for self-delegation, LNTN otherwise
        uint256 unbondingShare;
        uint256 requestBlock;
        bool unlocked;
        bool released;
        bool selfDelegation;
    }

    /**************************************************/
    struct Contracts {
        IAccountability accountabilityContract;
        IOracle oracleContract;
        IACU acuContract;
        ISupplyControl supplyControlContract;
        IStabilization stabilizationContract;
        IUpgradeManager upgradeManagerContract;
        IInflationController inflationControllerContract;
        IOmissionAccountability omissionAccountabilityContract;
        IAuctioneer auctioneerContract;
    }

    // parameters that affect the economic of the system.
    struct Policy {
        uint256 treasuryFee;
        uint256 minBaseFee;
        uint256 delegationRate;
        uint256 unbondingPeriod;
        uint256 initialInflationReserve;
        uint256 withholdingThreshold;
        uint256 proposerRewardRate; // fraction of epoch fees allocated for proposer rewarding based on activity proof
        uint256 oracleRewardRate;
        address payable withheldRewardsPool; // set to the autonity global treasury at genesis, but can be changed
        address payable treasuryAccount;
        uint256 baseFeeChangeDenominator; // EIP-1559
        uint256 elasticityMultiplier; // EIP-1559
    }

    // protocol parameters unrelated to the economic of the system. Expected to rarely change.
    struct Protocol {
        address operatorAccount;
        uint256 epochPeriod;
        uint256 blockPeriod;
        uint256 committeeSize;
        uint256 maxScheduleDuration;
        uint256 gasLimit;
        uint256 gasLimitBoundDivisor;
    }

    struct Config {
        Policy policy;
        Contracts contracts;
        Protocol protocol;
        uint256 contractVersion;
    }

    /* Any change in CommitteeMember struct must be synced with:
     * 1. CommitteeSelector code to write committee in DB (see `CommitteeSelector.updateCommittee` function in core/vm/contracts.go)
     * 2. AbsenteeComputer code to read the committee from the DB (see `readCommittee` function in core/vm/contracts.go)
     */
    struct CommitteeMember {
        address addr;
        uint256 votingPower;
        bytes consensusKey;
    }

    struct EpochInfo {
        CommitteeMember[] committee;
        uint256 previousEpochBlock;
        uint256 epochBlock;
        uint256 nextEpochBlock;
        uint256 omissionDelta;
        Eip1559 eip1559;
    }

    struct Accountability {
        uint256 range;
        uint256 delta;
        uint256 gracePeriod;
    }

    // part of the config which the golang client is keeping track of
    struct ClientAwareConfig {
        uint256 epochPeriod;
        uint256 blockPeriod;
        uint256 gasLimit;
        Accountability accountability;
        Eip1559 eip1559;
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
    function registerValidator(
        string memory _enode,
        address _oracleAddress,
        bytes memory _consensusKey,
        bytes memory _signatures
    ) external;

    /**
    * @notice Update enode of a registered validator. This function updates the network connection information (IP or/and port)
    of a registered validator. you cannot change the validator's address (pubkey part of the enode)
    * @param _nodeAddress This identifies the validator you want to update
    * @param _enode new enode to be updated
    */
    function updateEnode(address _nodeAddress, string memory _enode) external;

    /**
    * @notice Create a bonding(delegation) request with the caller as delegator.
    * @param _validator address of the validator to delegate stake to.
    * @param _amount total amount of NTN to bond.
    * @return uint256 id of the bonding request in the bonding queue
    */
    function bond(address _validator, uint256 _amount) external returns (uint256);

    /**
    * @notice Create an unbonding request with the caller as delegator.
    * @param _validator address of the validator to unbond stake to.
    * @param _amount total amount of LNTN (or NTN if self delegated) to unbond.
    * @return uint256 id of the unbonding request in the unbonding queue
    */
    function unbond(address _validator, uint256 _amount) external returns (uint256);

    /**
    * @notice Create a bonding(delegation) request with the `_account` as delegator. The caller needs to have required
    * bonding-allowance (NTN) from the `_account`.
    * @param _account address of the delegator.
    * @param _validator address of the validator to delegate stake to.
    * @param _amount total amount of NTN to bond.
    * @return uint256 id of the bonding request in the bonding queue
    */
    function bondFrom(address _account, address _validator, uint256 _amount) external returns (uint256);

    /**
    * @notice Create an unbonding request with the `_account` as delegator. The caller needs to have required
    * unbonding-allowance (self-unbonding-allowance) to unbond LNTN (NTN) from the `_account`.
    * @param _account address of the delegator.
    * @param _validator address of the validator to unbond stake to.
    * @param _amount total amount of LNTN (or NTN if self delegated) to unbond.
    * @return uint256 id of the unbonding request in the unbonding queue
    */
    function unbondFrom(address _account, address _validator, uint256 _amount) external returns (uint256);

    /**
     * @notice Returns the remaining number of NTN that `_staker` will be
     * allowed to bond on behalf of `_owner` through `bondFrom`.
     * This is zero by default.
     */
    function bondingAllowance(address _owner, address _staker) external view returns (uint256);

    /**
     * @notice Sets `_amount` as the bonding-allowance (NTN) of `_staker` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits an {BondingApproval} event.
     */
    function approveBonding(address _staker, uint256 _amount) external returns (bool);

    /**
    * @dev Bonds the inflation rewards to the validator's stake at epoch finalization.
    * Restricted to protocol contracts.
    */
    function autobond(address _validator, uint256 _selfBond, uint256 _delegated) external;

    /**
    * @notice Pause the validator and stop it accepting delegations.
    * @param _address address to be disabled.
    * @dev emit a {DisabledValidator} event.
    */
    function pauseValidator(address _address) external;

    /**
    * @notice Re-activate the specified validator.
    * @param _address address to be enabled.
    */
    function activateValidator(address _address) external;

    /**
    * @notice Change commission rate for the specified validator.
    * @param _validator address to be enabled.
            _rate new commission rate, ranging between 0-10000 (10000 = 100%).
    */
    function changeCommissionRate(address _validator, uint256 _rate) external;

    /**
    * @dev jails the specified validator
    * @param _nodeAddress the node address of the validator to be jailed
    * @param _jailtime the jailing time to be assigned to the validator
    * @param _newJailedState the validator state to be applied
    * @return uint256 the block at which the validator will be released from jail
    */
    function jail(
        address _nodeAddress,
        uint256 _jailtime,
        ValidatorState _newJailedState
    ) external returns (uint256); 

    /**
    * @dev jailbounds the specified validator
    * @param _nodeAddress the node address of the validator to be jailbound
    * @param _newJailboundState the validator state to be applied
    */
    function jailbound(
        address _nodeAddress,
        ValidatorState _newJailboundState
    ) external;

    /**
    * @dev slashes the specified validator
    * @dev NOTE: 100% slash is not allowed and if attempted will cause a revert
    * @param _nodeAddress the node address of the validator to be slashed
    * @param _slashingRate the rate for the slash
    * @return slashingAmount the slashing amount
    */
    function slash(
        address _nodeAddress,
        uint256 _slashingRate
    ) external returns (uint256 slashingAmount);

    /**
      * @dev slashes and jails the specified validator
      * @param _nodeAddress the node address of the validator to be slashed
      * @param _slashingRate the rate to be used
      * @param _jailtime the jailing time to be assigned to the validator
      * @param _newJailedState the validator state to be applied for jailing
      * @param _newJailboundState the validator state to be applied in case of 100% slashing
      * @return slashingAmount the amount slashed in NTN
      * @return jailReleaseBlock the block at which the validator will be released from jail
      * @return isJailbound a flag that signals if the validator has been permanently jailed
      */
    function slashAndJail(
        address _nodeAddress,
        uint256 _slashingRate,
        uint256 _jailtime,
        ValidatorState _newJailedState,
        ValidatorState _newJailboundState
    ) external returns (
        uint256 slashingAmount,
        uint256 jailReleaseBlock,
        bool isJailbound
    );

    /**
    * @notice Returns the current operator account.
    */
    function getOperator() external view returns (address);

    /**
    * @notice Returns the current Oracle account.
    */
    function getOracle() external view returns (address);

    /**
    * @notice Returns the liquid logic contract
    */
    function getLiquidLogicContract() external view returns (address);

    /**
    * @notice Returns the current client aware config
    */
    function getClientConfig() external view returns (ClientAwareConfig memory);

    /**
    * @notice Returns the bonding request corresponding to bonding ID.
    */
    function getBondingRequestByID(uint256 _id) external view returns (BondingRequest memory);

    /**
    * @notice Returns the unbonding request corresponding to unbonding ID.
    */
    function getUnbondingRequestByID(uint256 _id) external view returns (UnbondingRequest memory);

    /**
    * @notice Returns the epoch period. If there will be an update at epoch end, the new epoch period is returned
    */
    function getEpochPeriod() external view returns (uint256);

    /**
    * @notice Returns the epoch period of the current epoch
    */
    function getCurrentEpochPeriod() external view returns (uint256);

    /**
    * @notice Returns the block period.
    */
    function getBlockPeriod() external view returns (uint256);

    /**
     * @notice Returns the un-bonding period.
     */
    function getUnbondingPeriod() external view returns (uint256);

    /**
     * @notice Returns the current epoch ID
     */
    function getEpochID() external view returns (uint256);

    /**
    * @notice Returns the last epoch's end block height.
    */
    function getLastEpochBlock() external view returns (uint256);

    /**
    * @notice Returns the last epoch's end block timestamp
    */
    function getLastEpochTime() external view returns (uint256);

    /**
    * @notice Returns the current contract config
    */
    function getConfig() external view returns (Config memory);

    /**
    * @notice Returns the current contract version.
    */
    function getVersion() external view returns (uint256);

    /**
    * @notice Returns the current inflation reserve
    */
    function getInflationReserve() external view returns (uint256);

    /**
    * @notice Returns the current epoch total bonded stake
    */
    function getEpochTotalBondedStake() external view returns (uint256);

    /**
    * @notice Returns the current epoch info of the chain.
    */
    function getEpochInfo() external view returns (EpochInfo memory);

    /**
     * @notice Returns the block committee.
     */
    function getCommittee() external view returns (CommitteeMember[] memory);

    /**
     * @notice Returns the current list of validators.
     */
    function getValidators() external view returns (address[] memory);

    /**
     * @notice Returns the current treasury account.
     */
    function getTreasuryAccount() external view returns (address);

    /**
     * @notice Returns the current treasury fee.
     */
    function getTreasuryFee() external view returns (uint256);

    /**
     * @notice Returns the next epoch block.
     */
    function getNextEpochBlock() external view returns (uint256);

    /**
     * @notice Returns the amount of tokens circulating in the network.
     */
    function circulatingSupply() external view returns (uint256);

    /**
    * @notice Returns the validator object associated with `_addr`.
    */
    function getValidator(address _addr) external view returns (Validator memory);

    /**
    * @notice Returns the state of the validator associated with `_addr`.
    */
    function getValidatorState(address _addr) external view returns (ValidatorState);

    /**
    * @notice Returns the current size of the consensus committee.
    */
    function getCurrentCommitteeSize() external view returns (uint256);

    /**
    * @notice Returns the maximum size of the consensus committee.
    */
    function getMaxCommitteeSize() external view returns (uint256);
    /**
     * @notice Returns the max allowed duration of any schedule or contract.
     */
    function getMaxScheduleDuration() external view returns (uint256);

    /**
     * @notice Returns the consensus committee enodes.
     */
    function getCommitteeEnodes() external view returns (string[] memory);

    /**
     * @notice Returns the minimum gas price.
     */
    function getMinimumBaseFee() external view returns (uint256);

    /**
    * @notice Returns the epoch info of the height.
    */
    function getEpochByHeight(uint256 _height) external view returns (EpochInfo memory);

    /**
     * @notice Returns epoch associated to the block number.
     * @param _block the input block number.
     */
    function getEpochFromBlock(uint256 _block) external view returns (uint256);

    /**
     * @notice Returns `true` if unbonding is released and `false` otherwise.
     * @param _unbondingID id of the unbonding request in unbonding queue
     */
    function isUnbondingReleased(uint256 _unbondingID) external view returns (bool);

    /**
    * @notice Returns the unbonding share after the unbonding request is applied at epoch end.
    * @param _unbondingID id of the unbonding request in unbonding queue
    */
    function getUnbondingShare(uint256 _unbondingID) external view returns (uint256);

    /**
    * @notice Emitted after updating a enode address of a validator
    * @param validator validator address
    * @param oldEnode old enode of validator
    * @param newEnode new enode of validator
    */
    event EnodeUpdate(address validator, string oldEnode, string newEnode);

    event MintedStake(address indexed addr, uint256 amount);
    event BurnedStake(address indexed addr, uint256 amount);

    /**
    * @notice Emitted after registering a commission rate change request has been submitted
    * for the validator
    * @param validator validator address
    * @param rate new rate
    */
    event CommissionRateChange(address indexed validator, uint256 rate);

    /**
    * @notice This event is emitted when a bonding request to a validator node has been registered.
    * This request will only be effective at the end of the current epoch however the stake will be
    * put in custody immediately from the delegator's account.
    * @param validator The validator node account.
    * @param delegator The caller.
    * @param selfBonded True if the validator treasury initiated the request. No LNEW will be issued.
    * @param amount The amount of NEWTON to be delegated.
    * @param headBondingID  id of the request in bonding map
    */
    event NewBondingRequest(address indexed validator, address indexed delegator,
        bool selfBonded, uint256 amount, uint256 headBondingID);

    /**
    * @notice This event is emitted when a registered bonding request to a validator is rejected
    * @param validator The validator node account.
    * @param delegator The caller.
    * @param amount The amount of NEWTON to be delegated.
    * @param state The current state of validator.
    */
    event BondingRejected(address indexed validator, address indexed delegator, uint256 amount, ValidatorState state);

    /**
    * @notice This event is emitted when an unbonding request to a validator node has been registered.
    * This request will only be effective after the unbonding period, rounded to the next epoch.
    * Please note that because of potential slashing events during this delay period, the released amount
    * may or may not be correspond to the amount requested.
    * @param validator The validator node account.
    * @param delegator The caller.
    * @param selfBonded True if the validator treasury initiated the request.
    * @param amount If self-bonded this is the requested amount of NEWTON to be unbonded.
    * @param headUnbondingID id of the request in unbonding map
    */
    event NewUnbondingRequest(address indexed validator, address indexed delegator, bool selfBonded, uint256 amount, uint256 headUnbondingID);

    /**
    * @notice emitted when a new validator is registered
    * @param treasury treasury account
    * @param addr node address of the validator.
    * @param oracleAddress address of the oracle node run by validator
    * @param enode enode of validator
    * @param liquidStateContract liquidStateContract address
    */
    event RegisteredValidator(address treasury, address addr, address oracleAddress, string enode, address liquidStateContract);

    /**
    * @notice emitted when a validator is paused, validator stops accepting delegation and becomes ineligible
    *  for committee inclusion
    * @param treasury treasury account
    * @param addr node address of the validator.
    * @param effectiveBlock the block number at which the new paused state takes effect
    */
    event PausedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock);

    /**
    * @notice emitted when a validator is activated, validator starts accepting delegation and becomes eligible
    *  for committee inclusion
    * @param treasury treasury account
    * @param addr node address of the validator.
    * @param effectiveBlock the block number at which the new activated state takes effect
    */
    event ActivatedValidator(address indexed treasury, address indexed addr, uint256 effectiveBlock);

    /**
    * @notice emitted when a validator is rewarded for taking part in block consensus
    *  for committee inclusion
    * @param addr node address of validator
    * @param atnSelfAmount validator's atn rewards
    * @param atnDelegatedAmount delegator's atn rewards (includes validator commission)
    * @param ntnSelfAmount validator's ntn rewards
    * @param ntnDelegatedAmount delegator's ntn rewards (includes validator commission)
    */
    event Rewarded(address indexed addr, uint256 atnSelfAmount, uint256 atnDelegatedAmount, uint256 ntnSelfAmount, uint256 ntnDelegatedAmount);

    /**
    * @notice emitted when epoch period is updated
    * @param period new epoch period in blocks
    * @param appliedAtBlock the block number at which the new value takes effect
    */
    event EpochPeriodUpdated(uint256 period, uint256 appliedAtBlock);

    /**
    * @notice emitted when a newEpoch begins
    * @param epoch epoch block number
    * @param inflationReserve total inflation reserves at the end of current epoch
    * @param stakeCirculating total ciculating stake  at the end of current epoch
    */
    event NewEpoch(uint256 epoch, uint256 inflationReserve, uint256 stakeCirculating);

    /**
     * @notice This event is emitted when a call to an address fails in a protocol function (like finalize()).
     * @param to address
     * @param methodSignature method signature of the call, empty in case of plain transaction
     * @param returnData low level return data
     */
    event CallFailed(address to, string methodSignature, bytes returnData);

    /**
    * @notice Event emitted after EIP-1559 parameters are updated
    */
    event Eip1559ParamsUpdate(Eip1559 oldParams, Eip1559 newParams);

    /**
     * @notice Emitted when the bonding-allowance (NTN) of a `_staker` for an `_owner` is set by
     * a call to `approveBonding`. `_value` is the new `bondingAllowance` (NTN).
     */
    event BondingApproval(address indexed _owner, address indexed _staker, uint256 _value);

}
