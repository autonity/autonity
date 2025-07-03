// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 ;
import "./IERC20.sol";

interface ILiquid is IERC20 {
    /**
     * @notice Issue new LNTN when some NTN are bonded.
     */
    function mint(address _account, uint256 _amount) external;

    /**
     * @notice Unlock the locked LNTN
     */
    function unlock(address _account, uint256 _amount) external;

    /**
     * @notice Lock LNTN to unbond at the epoch end.
     */
    function lock(address _account, uint256 _amount) external;

    /**
     * @notice Lock LNTN from `_account` to unbond at the epoch end. `_staker` is the caller
     * and must have enough `unbondingAllowance`.
     */
    function lockFrom(address _account, address _staker, uint256 _amount) external;

    /**
     * @notice Change validator commission rate.
     */
    function setCommissionRate(uint256 _rate) external;

    /**
     * @notice Claim accrued atn rewards.
     */
    function claimRewards() external;

    /**
     * @notice Claim the unclaimed atn rewards by the treasury.
     */
    function claimTreasuryATN() external;
    
    /**
     * @notice Burn LNTN when they are unbonded.
     */
    function burn(address _account, uint256 _amount) external;

    /**
     * @notice Distribute the atn rewards among delegators and auto-bond the ntn rewards
     * proportional to their balances
     */
    function redistribute(uint256 _ntnReward) external payable returns (uint256);

    /**
     * @notice Amount of atn rewards yet to be claimed
     */
    function unclaimedRewards(address _account) external view returns(uint256);

    function decimals() external pure returns (uint8);

    /**
     * @notice Amount of LNTN locked to be unbonded.
     */
    function lockedBalanceOf(address _delegator) external view returns (uint256);

    /**
     * @notice Amount of unlocked LNTN.
     */
    function unlockedBalanceOf(address _delegator) external view returns (uint256);

    function name() external view returns (string memory);

    function symbol() external view returns (string memory);

    function getValidator() external view returns (address);

    function getTreasury() external view returns (address);

    /**
     * @notice Returns the commission rate of the validator.
     */
    function getCommissionRate() external view returns (uint256);

    /**
     * @notice Amount of unclaimed atn rewards by the treasury.
     */
    function getTreasuryUnclaimedATN() external view returns (uint256);

    /**
     * @notice Returns the remaining number of LNTN that `_staker` will be
     * allowed to unbond on behalf of `_owner` through `unbondFrom`.
     * This is zero by default.
     */
    function unbondingAllowance(address _owner, address _staker) external view returns (uint256);

    /**
     * @notice Sets `_amount` as the unbond-allowance (LNTN) of `_staker` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits an {UnbondingApproval} event.
     */
    function approveUnbonding(address _staker, uint256 _amount) external returns (bool);

    /**
     * @notice Emitted when the unbond-allowance (LNTN) of a `_staker` for an `_owner` is set by
     * a call to `approveUnbonding`. `_value` is the new unbond-allowance (LNTN).
     */
    event UnbondingApproval(address indexed _owner, address indexed _staker, uint256 _value);
}