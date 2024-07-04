# Solidity API

## ContractBase

### beneficiaryContracts

```solidity
mapping(address => uint256[]) beneficiaryContracts
```

_Stores the unique ids of contracts assigned to a beneficiary, but beneficiary does not need to know the id
beneficiary will number his contracts as: 0 for first contract, 1 for 2nd and so on.
We can get the unique contract id from beneficiaryContracts as follows:
`beneficiaryContracts[beneficiary][0]` is the unique id of his first contract
`beneficiaryContracts[beneficiary][1]` is the unique id of his 2nd contract and so on_

### contracts

```solidity
struct ContractBase.Contract[] contracts
```

_List of all contracts_

### _calculateAvailableUnlockedFunds

```solidity
function _calculateAvailableUnlockedFunds(uint256 _contractID, uint256 _totalValue, uint256 _time) internal view returns (uint256)
```

_Given the total value (in NTN) of the contract, calculates the amount of withdrawable tokens (in NTN)._

### _calculateTotalUnlockedFunds

```solidity
function _calculateTotalUnlockedFunds(uint256 _start, uint256 _totalDuration, uint256 _time, uint256 _totalAmount) internal pure returns (uint256)
```

_Calculates total unlocked funds while assuming cliff period has passed.
Check if cliff is passed before calling this function._

### _getUniqueContractID

```solidity
function _getUniqueContractID(address _beneficiary, uint256 _id) internal view returns (uint256)
```

_Returns a unique id for each contract._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | address of the contract holder |
| _id | uint256 | contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones) |

### _updateAndTransferNTN

```solidity
function _updateAndTransferNTN(uint256 _contractID, address _to, uint256 _amount) internal
```

_Updates the contract with `contractID` and transfers NTN._

### getContract

```solidity
function getContract(address _beneficiary, uint256 _id) external view virtual returns (struct ContractBase.Contract)
```

Returns id'th contract entitled to beneficiary.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | beneficiary address |
| _id | uint256 | contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones) |

### getContracts

```solidity
function getContracts(address _beneficiary) external view virtual returns (struct ContractBase.Contract[])
```

Returns the list of current contracts assigned to a beneficiary.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | address of the beneficiary |

### canStake

```solidity
function canStake(address _beneficiary, uint256 _id) external view virtual returns (bool)
```

Returns if beneficiary can stake from his contract.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | beneficiary address |
| _id | uint256 |  |

### totalContracts

```solidity
function totalContracts(address _beneficiary) external view virtual returns (uint256)
```

Returns the number of contracts entitled to some beneficiary.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | address of the beneficiary |

### onlyOperator

```solidity
modifier onlyOperator()
```

_Modifier that checks if the caller is the governance operator account._

## LiquidRewardManager

### RewardEvent

_This structure tracks the activity that requires to update rewards. Any activity like increase or decrease of liquid balances
or claiming rewards require updating rewards before the activity is applied. This structure tracks those activities, when such
activity is requested.

There are two types of reward event:
Pending Reward Event: A reward event that cannot be applied yet. Bonding or unbonding request creates such reward event.
Last Reward Event: A reward event that can be applied. When a pending reward event can be applied, after the epcoh is finalized,
it becomes last reward event. If there already exists a last reward event, then it is applied before being replaced by a pending reward event._

```solidity
struct RewardEvent {
  uint256 epochID;
  uint256 totalLiquid;
  uint256 stakingRequestID;
  bool isBonding;
  bool eventExist;
  bool applied;
}
```

### RewardTracker

_Tracks rewards for each validator._

```solidity
struct RewardTracker {
  uint256 atnUnclaimedRewards;
  uint256 ntnUnclaimedRewards;
  uint256 lastUpdateEpochID;
  contract Liquid liquidContract;
  struct LiquidRewardManager.RewardEvent lastRewardEvent;
  struct LiquidRewardManager.RewardEvent pendingRewardEvent;
  mapping(uint256 => uint256) atnLastUnrealisedFeeFactor;
  mapping(uint256 => uint256) ntnLastUnrealisedFeeFactor;
  mapping(uint256 => bool) unrealisedFeeFactorUpdated;
}
```

### Account

_Each account represents a bonding from some contract, `id` to some validator, `v` and stored in mapping `accounts[id][v]`.
Multiple bonding and unbonding for the same pair is aggregated, so there is at most one account for each pair._

```solidity
struct Account {
  uint256 liquidBalance;
  uint256 lockedLiquidBalance;
  uint256 atnRealisedFee;
  uint256 atnUnrealisedFeeFactor;
  uint256 ntnRealisedFee;
  uint256 ntnUnrealisedFeeFactor;
  bool newBondingRequested;
}
```

### _newBondingRequested

```solidity
function _newBondingRequested(uint256 _id, address _validator, uint256 _bondingID) internal
```

_Adds the validator in the list and inform that new bonding is requested._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 | contract id |
| _validator | address | validator address |
| _bondingID | uint256 |  |

### _burnLiquid

```solidity
function _burnLiquid(uint256 _id, address _validator, uint256 _amount, uint256 _epochID) internal
```

_Burns some liquid tokens that represents liquid bonded to some validator from some contract.
The following functions: `burnLiquid`, `mintLiquid`, `realiseFees` follow the same logic as done in Liquid.sol.
The only difference is that the liquid is not updated immediately. The updating of liquid reflects the changes after
the staking operations of `epochID` are applied in Autonity.sol._

### _mintLiquid

```solidity
function _mintLiquid(uint256 _id, address _validator, uint256 _amount, uint256 _epochID) internal
```

_Mints some liquid tokens that represents liquid bonded to some validator from some contract._

### _claimRewards

```solidity
function _claimRewards(uint256 _id) internal returns (uint256 _atnTotalFees, uint256 _ntnTotalFees)
```

_Calculates total rewards for a contract and resets `realisedFees[id][validator]` as rewards are claimed_

### _clearValidators

```solidity
function _clearValidators(uint256 _id) internal
```

_Removes all the validators that are not needed for some contract anymore, i.e. any validator
that has 0 liquid for that contract and all rewards from the validator are claimed._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 | contract id |

### getPendingRewardEvent

```solidity
function getPendingRewardEvent(address _validator) public view returns (struct LiquidRewardManager.RewardEvent)
```

Returns the last requested reward update event which is pending.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _validator | address | validator address |

### getLastRewardEvent

```solidity
function getLastRewardEvent(address _validator) public view returns (struct LiquidRewardManager.RewardEvent)
```

Returns the last requested reward update event which is not pending, i.e. can be taken into account because
the request-epoch has passed.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _validator | address | validator address |

### _newPendingRewardEvent

```solidity
function _newPendingRewardEvent(address _validator, uint256 _stakingRequestID, bool _isBonding) internal
```

_Adds a new reward event which is pending because the rewards distribution will happen at epoch end._

### _updatePendingEventLiquid

```solidity
function _updatePendingEventLiquid(address _validator) internal
```

_In case a transfer of liquid token happens, it will affect the total rewards distribution at epoch end,
because the total liquid is changed. If we already have a pending reward event to calculate reward distribution,
then this change in liquid balance needs to be considered. This function updates the total liquid in the pending event._

### _unclaimedRewards

```solidity
function _unclaimedRewards(uint256 _id, address _validator, int256 _balanceChange, uint256 _updateEpochID) internal view returns (uint256 _atnReward, uint256 _ntnReward)
```

_Calculates the rewards yet to claim for id from validator._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 | unique contract id |
| _validator | address | validator address |
| _balanceChange | int256 | change in balance after some epoch |
| _updateEpochID | uint256 | epoch ID after which balance is changed |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| _atnReward | uint256 | amount of unclaimed ATN rewards |
| _ntnReward | uint256 | amount of unclaimed NTN rewards |

### _bondedValidators

```solidity
function _bondedValidators(uint256 _id) internal view returns (address[])
```

_Returns the list of validator addresses which are bonded to some contract._

## NonStakableVesting

### totalNominal

```solidity
uint256 totalNominal
```

The total amount of funds to create new locked non-stakable schedules.
The balance is not immediately available at the vault.
Rather the unlocked amount of schedules is minted at epoch end.
The balance tells us the max size of a newly created schedule.
See `createSchedule()`

### maxAllowedDuration

```solidity
uint256 maxAllowedDuration
```

The maximum duration of any schedule or contract

### createSchedule

```solidity
function createSchedule(uint256 _amount, uint256 _startTime, uint256 _cliffDuration, uint256 _totalDuration) public virtual
```

Creates a new schedule.

_The schedule has unsubscribedAmount = amount initially. As new contracts are subscribed to the schedule, its unsubscribedAmount decreases.
At any point, `subscribedAmount of schedule = amount - unsubscribedAmount`._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _amount | uint256 | total amount of the schedule |
| _startTime | uint256 | start time |
| _cliffDuration | uint256 | cliff period, after _cliffDuration + _startTime, the schedule will have claimables |
| _totalDuration | uint256 | total duration of the schedule |

### newContract

```solidity
function newContract(address _beneficiary, uint256 _amount, uint256 _scheduleID) public virtual
```

Creates a new non-stakable contract which subscribes to some schedule.

_If the contract is created before cliff period has passed, the beneficiary is entitled to NTN as it unlocks.
Otherwise, the contract already has some unlocked NTN which is not entitled to beneficiary. However, NTN that will
be unlocked in future will be entitled to beneficiary._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | address of the beneficiary |
| _amount | uint256 | total amount of NTN to be vested |
| _scheduleID | uint256 | schedule to subscribe |

### setTotalNominal

```solidity
function setTotalNominal(uint256 _totalNominal) external virtual
```

Sets the `totalNominal` value.

### setMaxAllowedDuration

```solidity
function setMaxAllowedDuration(uint256 _newMaxDuration) external virtual
```

Sets the max allowed duration of any schedule or contract.

### releaseAllFunds

```solidity
function releaseAllFunds(uint256 _id) external virtual
```

Used by beneficiary to transfer all unlocked NTN of some contract to his own address.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 | id of the contract numbered from 0 to (n-1) where n = total contracts entitled to the beneficiary (excluding canceled ones). So any beneficiary can number their contracts from 0 to (n-1). Beneficiary does not need to know the unique global contract id. |

### releaseFund

```solidity
function releaseFund(uint256 _id, uint256 _amount) external virtual
```

Used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 |  |
| _amount | uint256 | amount of NTN to release |

### changeContractBeneficiary

```solidity
function changeContractBeneficiary(address _beneficiary, uint256 _id, address _recipient) external virtual
```

Changes the beneficiary of some contract to the recipient address. The recipient address can release tokens from the contract as it unlocks.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | beneficiary address whose contract will be canceled |
| _id | uint256 | contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones) |
| _recipient | address | whome the contract is transferred to |

### unlockTokens

```solidity
function unlockTokens() external returns (uint256 _newUnlockedSubscribed, uint256 _newUnlockedUnsubscribed)
```

Unlock tokens of all schedules upto current time.

_It calculates the newly unlocked tokens upto current time and also updates the amount
of total unlocked tokens and the time of unlock for each schedule
Autonity must mint new unlocked tokens, because this contract knows that for each schedule,
`schedule.totalUnlocked` tokens are now unlocked and available to release
`newUnlockedSubscribed` goes to the balance of address(this) and `newUnlockedUnsubscribed` goes to the treasury address.
See `finalize()` in Autonity.sol_

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| _newUnlockedSubscribed | uint256 | tokens unlocked from contract subscribed to some schedule |
| _newUnlockedUnsubscribed | uint256 | tokens unlocked from schedule.unsubscribedAmount, which is not subscribed by any contract |

### unlockedFunds

```solidity
function unlockedFunds(address _beneficiary, uint256 _id) external view virtual returns (uint256)
```

Returns the amount of withdrawable funds upto the last epoch time.

### getSchedule

```solidity
function getSchedule(uint256 _id) external view returns (struct NonStakableVesting.Schedule)
```

Returns some schedule

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 | id of some schedule numbered from 0 to (n-1), where n = total schedules created |

## StakableVesting

### totalNominal

```solidity
uint256 totalNominal
```

Sum of total amount of contracts that can be created.
Each time a new contract is created, `totalNominal` is decreased.
Address(this) should have `totalNominal` amount of NTN availabe at genesis,
otherwise withdrawing or bonding from a contract is not possible

### contractToBonding

```solidity
mapping(uint256 => uint256[]) contractToBonding
```

_We put all the bonding request id of past epoch in `contractToBonding[contractID]` array and apply them whenever needed.
All bonding requests are applied at epoch end, so we can process all of them (failed or successful) together.
See `bond` and `_handlePendingBondingRequest` for more clarity_

### contractToUnbonding

```solidity
mapping(uint256 => mapping(uint256 => uint256)) contractToUnbonding
```

_We put all the unbonding request id of past epoch in contractToUnbonding mapping. All requests from past epoch
can be applied together. But not all requests are released together at epoch end. So we need to put them in map
and use `tailPendingUnbondingID` and `headPendingUnbondingID` to keep track of contractToUnbonding.
See `unbond` and `_handlePendingUnbondingRequest` for more clarity_

### atnRewards

```solidity
mapping(address => uint256) atnRewards
```

_ATN rewards entitled to some beneficiary for bonding from some contract before it has been cancelled.
See cancelContract for more clarity._

### ntnRewards

```solidity
mapping(address => uint256) ntnRewards
```

_Same as atnRewards for NTN rewards_

### newContract

```solidity
function newContract(address _beneficiary, uint256 _amount, uint256 _startTime, uint256 _cliffDuration, uint256 _totalDuration) public virtual
```

Creates a new stakable contract.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | address of the beneficiary |
| _amount | uint256 | total amount of NTN to be vested |
| _startTime | uint256 | start time of the vesting |
| _cliffDuration | uint256 | cliff period |
| _totalDuration | uint256 | total duration of the contract |

### releaseFunds

```solidity
function releaseFunds(uint256 _id) external virtual
```

Used by beneficiary to transfer all unlocked NTN and LNTN of some contract to his own address.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 | contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones). So any beneficiary can number their contracts from 0 to (n-1). Beneficiary does not need to know the unique global contract id. |

### releaseAllNTN

```solidity
function releaseAllNTN(uint256 _id) external virtual
```

Used by beneficiary to transfer all unlocked NTN of some contract to his own address.

### releaseAllLNTN

```solidity
function releaseAllLNTN(uint256 _id) external virtual
```

Used by beneficiary to transfer all unlocked LNTN of some contract to his own address.

### releaseNTN

```solidity
function releaseNTN(uint256 _id, uint256 _amount) external virtual
```

Used by beneficiary to transfer some amount of unlocked NTN of some contract to his own address.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 |  |
| _amount | uint256 | amount to transfer |

### releaseLNTN

```solidity
function releaseLNTN(uint256 _id, address _validator, uint256 _amount) external virtual
```

Used by beneficiary to transfer some amount of unlocked LNTN of some contract to his own address.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 |  |
| _validator | address | address of the validator |
| _amount | uint256 | amount of LNTN to transfer |

### changeContractBeneficiary

```solidity
function changeContractBeneficiary(address _beneficiary, uint256 _id, address _recipient) external virtual
```

Changes the beneficiary of some contract to the recipient address. The recipient address can release and stake tokens from the contract.
Rewards which have been entitled to the beneficiary due to bonding from this contract are not transferred to recipient.

_Rewards earned until this point from this contract are calculated and stored in atnRewards and ntnRewards mapping so that
beneficiary can later claim them even though beneficiary is not entitled to this contract._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | beneficiary address whose contract will be canceled |
| _id | uint256 | contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding already canceled ones) |
| _recipient | address | whome the contract is transferred to |

### updateFunds

```solidity
function updateFunds(address _beneficiary, uint256 _id) external virtual
```

In case some funds are missing due to some pending staking operation that failed,
this function updates the funds of some contract entitled to beneficiary by applying the pending requests.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | beneficiary address |
| _id | uint256 | contract id to update |

### updateFundsAndGetContractTotalValue

```solidity
function updateFundsAndGetContractTotalValue(address _beneficiary, uint256 _id) external returns (uint256)
```

Updates the funds of the contract and returns total value of the contract

### updateFundsAndGetContract

```solidity
function updateFundsAndGetContract(address _beneficiary, uint256 _id) external returns (struct ContractBase.Contract)
```

Updates the funds of the contract and returns the contract

### setTotalNominal

```solidity
function setTotalNominal(uint256 _newTotalNominal) external virtual
```

Set the value of totalNominal
In case totalNominal is increased, the increased amount should be minted
and transferred to the address of this contract, otherwise newly created vesting
contracts will not have funds to withdraw or bond. See newContract()

### bond

```solidity
function bond(uint256 _id, address _validator, uint256 _amount) public virtual returns (uint256)
```

Used by beneficiary to bond some NTN of some contract.
All bondings are delegated, as vesting manager cannot own a validator.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 | id of the contract numbered from 0 to (n-1) where n = total contracts entitled to the beneficiary (excluding canceled ones) |
| _validator | address | address of the validator for bonding |
| _amount | uint256 | amount of NTN to bond |

### unbond

```solidity
function unbond(uint256 _id, address _validator, uint256 _amount) public virtual returns (uint256)
```

Used by beneficiary to unbond some LNTN of some contract.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 |  |
| _validator | address | address of the validator |
| _amount | uint256 | amount of LNTN to unbond |

### claimRewards

```solidity
function claimRewards(uint256 _id, address _validator) external virtual returns (uint256 _atnRewards, uint256 _ntnRewards)
```

Used by beneficiary to claim rewards from bonding some contract to validator.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _id | uint256 | contract ID |
| _validator | address | validator address |

### claimRewards

```solidity
function claimRewards(uint256 _id) external virtual returns (uint256 _atnRewards, uint256 _ntnRewards)
```

Used by beneficiary to claim rewards from bonding some contract to validator.

### claimRewards

```solidity
function claimRewards() external virtual returns (uint256 _atnRewards, uint256 _ntnRewards)
```

Used by beneficiary to claim all rewards which is entitled from bonding

_Rewards from some cancelled contracts are stored in atnRewards and ntnRewards mapping. All rewards from
contracts that are still entitled to the beneficiary need to be calculated._

### receive

```solidity
receive() external payable
```

_Receive Auton function https://solidity.readthedocs.io/en/v0.7.2/contracts.html#receive-ether-function_

### _calculateLNTNValue

```solidity
function _calculateLNTNValue(address _validator, uint256 _amount) internal view returns (uint256)
```

_Returns equivalent amount of NTN using the current ratio._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _validator | address | validator address |
| _amount | uint256 | amount of LNTN to be converted |

### _getLiquidFromNTN

```solidity
function _getLiquidFromNTN(address _validator, uint256 _amount) internal view returns (uint256)
```

_Returns equivalent amount of LNTN using the current ratio._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _validator | address | validator address |
| _amount | uint256 | amount of NTN to be converted |

### _calculateTotalValue

```solidity
function _calculateTotalValue(uint256 _contractID) internal view returns (uint256)
```

_Calculates the total value of the contract, which can vary if the contract has some LNTN.
`totalValue = currentNTN + withdrawnValue + (the value of LNTN converted to NTN using current ratio)`_

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _contractID | uint256 | unique global id of the contract |

### _releaseAllUnlockedLNTN

```solidity
function _releaseAllUnlockedLNTN(uint256 _contractID, uint256 _availableUnlockedFunds) internal returns (uint256 _remaining)
```

_Transfers some LNTN equivalent to beneficiary address. The amount of unlocked funds is calculated in NTN
and then converted to LNTN using the current ratio.
In case the contract has LNTN to multiple validators, we pick one validator and try to transfer
as much LNTN as possible. If there still remains some more uncloked funds, then we pick another validator.
There is no particular order in which validator should be picked first._

### _unlockedFunds

```solidity
function _unlockedFunds(uint256 _contractID) internal view returns (uint256)
```

_Calculates the amount of unlocked funds in NTN until last epoch time._

### _updateFunds

```solidity
function _updateFunds(uint256 _contractID) internal
```

_Updates the funds by applying the staking requests._

### _cleanup

```solidity
function _cleanup(uint256 _contractID) internal
```

_Updates the funds and removes any unnecessary validator from the list._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _contractID | uint256 | unique global contract id |

### _handlePendingBondingRequest

```solidity
function _handlePendingBondingRequest(uint256 _contractID) internal
```

_Handles all the pending bonding requests.
All the requests from past epoch can be handled as the bonding requests are
applied at epoch end immediately. Requests from current epoch are not handled._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _contractID | uint256 | unique global id of the contract |

### _handlePendingUnbondingRequest

```solidity
function _handlePendingUnbondingRequest(uint256 _contractID) internal
```

_Handles all the pending unbonding requests. All unbonding requests from past epoch are applied.
Unbonding request that are released in Autonity are released._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _contractID | uint256 | unique global id of the contract |

### _balanceChangeFromStakingRequest

```solidity
function _balanceChangeFromStakingRequest(uint256 _contractID, address _validator) internal view returns (int256 _balanceChange, uint256 _lastRequestEpoch)
```

_Calculates the balance changes for the pending staking requests to some validator.
If the requests are from current epoch, then they are not taken into account._

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _contractID | uint256 | unique contract ID |
| _validator | address | validator address |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| _balanceChange | int256 | balance change, positive if balance increases, negative otherwise |
| _lastRequestEpoch | uint256 | epochID, offset by 1, of the requests (all requests should be from the same epoch) |

### unclaimedRewards

```solidity
function unclaimedRewards(address _beneficiary, uint256 _id, address _validator) external view virtual returns (uint256 _atnRewards, uint256 _ntnRewards)
```

Returns unclaimed rewards from some contract entitled to beneficiary from bonding to validator.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | beneficiary address |
| _id | uint256 | contract ID |
| _validator | address | validator address |

#### Return Values

| Name | Type | Description |
| ---- | ---- | ----------- |
| _atnRewards | uint256 | unclaimed ATN rewards |
| _ntnRewards | uint256 | unclaimed NTN rewards |

### unclaimedRewards

```solidity
function unclaimedRewards(address _beneficiary, uint256 _id) public view virtual returns (uint256 _atnRewards, uint256 _ntnRewards)
```

Returns unclaimed rewards from some contract entitled to beneficiary from bonding.

### unclaimedRewards

```solidity
function unclaimedRewards(address _beneficiary) external view virtual returns (uint256 _atnRewards, uint256 _ntnRewards)
```

Returns the amount of all unclaimed rewards due to all the bonding from contracts entitled to beneficiary.

### liquidBalanceOf

```solidity
function liquidBalanceOf(address _beneficiary, uint256 _id, address _validator) external view virtual returns (uint256)
```

Returns the amount of LNTN for some contract.

#### Parameters

| Name | Type | Description |
| ---- | ---- | ----------- |
| _beneficiary | address | beneficiary address |
| _id | uint256 | contract id numbered from 0 to (n-1); n = total contracts entitled to the beneficiary (excluding canceled ones) |
| _validator | address | validator address |

### lockedLiquidBalanceOf

```solidity
function lockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) external view virtual returns (uint256)
```

Returns the amount of locked LNTN for some contract.

### unlockedLiquidBalanceOf

```solidity
function unlockedLiquidBalanceOf(address _beneficiary, uint256 _id, address _validator) external view virtual returns (uint256)
```

Returns the amount of unlocked LNTN for some contract.

### getBondedValidators

```solidity
function getBondedValidators(address _beneficiary, uint256 _id) external view virtual returns (address[])
```

Returns the list of validators bonded some contract.

### unlockedFunds

```solidity
function unlockedFunds(address _beneficiary, uint256 _id) external view virtual returns (uint256)
```

Returns the amount of released funds in NTN for some contract.

### contractTotalValue

```solidity
function contractTotalValue(address _beneficiary, uint256 _id) external view returns (uint256)
```

Returns the current total value of the contract in NTN.

