// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity ^0.8.19;

/**
 * @dev Interface of the Autonity Contract.
 * Import this over Autonity.sol.
 */
interface IAutonity {

    enum ValidatorState {active, paused, jailed, jailbound, jailedForInactivity, jailboundForInactivity}
    /**
    * @notice Returns the current operator account.
    */
    function getOperator() external view returns (address);

    /**
    * @notice Returns the current Oracle account.
    */
    function getOracle() external view returns (address);

    /**
    * @notice Emitted after updating config parameter of type uint
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    */
    event ConfigUpdateUint(string name, uint256 oldValue, uint256 newValue);

    /**
    * @notice Emitted after updating config parameter of type int
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    */
    event ConfigUpdateInt(string name, int256 oldValue, int256 newValue);

    /**
    * @notice Emitted after updating config parameter of type address
    * @param name configuration name
    * @param oldValue old value of configuration
    * @param newValue new value of configuration
    */
    event ConfigUpdateAddress(string name, address oldValue, address newValue);

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
    * @param atnAmount atnRewards
    * @param ntnAmount ntnRewards
    */
    event Rewarded(address indexed addr, uint256 atnAmount, uint256 ntnAmount);

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

}
