// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 ;
import "./IERC20.sol";

interface ILiquid is IERC20 {
    function mint(address _account, uint256 _amount) external;
    function unlock(address _account, uint256 _amount) external;
    function lock(address _account, uint256 _amount) external;
    function lockFrom(address _account, address _staker, uint256 _amount) external;
    function setCommissionRate(uint256 _rate) external;
    function claimRewards() external;
    function claimTreasuryATN() external;
    function burn(address _account, uint256 _amount) external;
    function redistribute(uint256 _ntnReward) external payable returns (uint256);
    function unclaimedRewards(address _account) external view returns(uint256);
    function decimals() external pure returns (uint8);
    function lockedBalanceOf(address _delegator) external view returns (uint256);
    function unlockedBalanceOf(address _delegator) external view returns (uint256);
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function getValidator() external view returns (address);
    function getTreasury() external view returns (address);
    function getCommissionRate() external view returns (uint256);
    function getTreasuryUnclaimedATN() external view returns (uint256);

    /**
     * @notice Returns the remaining number of LNTN that `_staker` will be
     * allowed to unbond on behalf of `_owner` through `unbondFrom`.
     * This is zero by default.
     */
    function unbondAllowance(address _owner, address _staker) external view returns (uint256);

    /**
     * @notice Sets `_amount` as the unbond-allowance (LNTN) of `_staker` over the caller's tokens.
     *
     * Returns a boolean value indicating whether the operation succeeded.
     *
     * Emits an {UnbondApproval} event.
     */
    function approveUnbond(address _staker, uint256 _amount) external returns (bool);

    /**
     * @notice Emitted when the unbond-allowance (LNTN) of a `_staker` for an `_owner` is set by
     * a call to `approveUnbond`. `_value` is the new unbond-allowance (LNTN).
     */
    event UnbondApproval(address indexed _owner, address indexed _staker, uint256 _value);
}