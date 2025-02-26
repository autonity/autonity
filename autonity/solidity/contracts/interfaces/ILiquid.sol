// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 ;
import "./IERC20.sol";

interface ILiquid is IERC20 {
    function mint(address _pool, uint256 _amount) external;
    function lockInPool(address _account, address _pool, uint256 _amount) external;
    function claimRewards() external;
    function claimTreasuryATN() external;
    function burn(address _pool, uint256 _amount) external;
    function redistribute(uint256 _ntnReward, uint256 _commissionRate) external payable returns (uint256);
    function transferFromPool(address _account, uint256 _amount) external;
    function unclaimedRewards(address _account) external view returns(uint256);
    function decimals() external pure returns (uint8);
    function balanceInContract(address _account) external view returns (uint256);
    function balanceInPool(address _account) external view returns (uint256);
    function lockedBalanceOf(address _delegator) external view returns (uint256);
    function unlockedBalanceOf(address _delegator) external view returns (uint256);
    function name() external view returns (string memory);
    function symbol() external view returns (string memory);
    function getValidator() external view returns (address);
    function getTreasury() external view returns (address);
    function getCommissionRate() external view returns (uint256);
    function getTreasuryUnclaimedATN() external view returns (uint256);

    event LiquidMinted(address indexed account, uint256 amount);
    event LiquidBurnt(address indexed account, uint256 amount);
}