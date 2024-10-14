// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 ;
import "./IERC20.sol";

interface ILiquid is IERC20 {
    function mint(address _account, uint256 _amount) external;
    function unlock(address _account, uint256 _amount) external;
    function lock(address _account, uint256 _amount) external;
    function setCommissionRate(uint256 _rate) external;
    function claimRewards() external;
    function claimTreasuryATN() external;
    function burn(address _account, uint256 _amount) external;
    function redistribute(uint256 _ntnReward) external payable returns (uint256, uint256);
    function unclaimedRewards(address _account) external view returns(uint256 _unclaimedATN, uint256 _unclaimedNTN);
    function decimals() external pure returns (uint8);
    function lockedBalanceOf(address _delegator) external view returns (uint256);
    function unlockedBalanceOf(address _delegator) external view returns (uint256);
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function getValidator() external view returns (address);
    function getTreasury() external view returns (address);
    function getCommissionRate() external view returns (uint256);
    function getTreasuryUnclaimedATN() external view returns (uint256);
}