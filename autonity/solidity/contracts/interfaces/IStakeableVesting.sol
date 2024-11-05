// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity >=0.8.2 ;

import "../vesting/ContractBase.sol";

interface IStakeableVesting {
    function createContract(address _beneficiary, uint256 _amount, uint256 _startTime, uint256 _cliffDuration, uint256 _totalDuration) external;
    function changeContractBeneficiary(address _recipient) external;
    function setManagerContract(address _managerContract) external;
    function releaseFunds() external;
    function releaseAllNTN() external;
    function releaseAllLNTN() external;
    function releaseNTN(uint256 _amount) external;
    function releaseLNTN(address _validator, uint256 _amount) external;
    function updateFunds() external;
    function updateFundsAndGetContract() external returns (ContractBase.Contract memory);
    function bond(address _validator, uint256 _amount) external returns (uint256);
    function unbond(address _validator, uint256 _amount) external returns (uint256);
    function claimRewards(address _validator) external;
    function claimRewards() external;
    function unclaimedRewards(address _validator) external view returns (uint256);
    function unclaimedRewards() external view returns (uint256);
    function vestedFunds() external view returns (uint256);
    function withdrawableVestedFunds() external view returns (uint256);
    function contractTotalValue() external view returns (uint256);
    function getManagerContractAddress() external view returns (address);
    function getBeneficiary() external view returns (address);
    function getContract() external view returns (ContractBase.Contract memory);
    function getLinkedValidators() external view returns (address[] memory);
    function liquidBalance(address _validator) external view returns (uint256);
    function unlockedLiquidBalance(address _validator) external view returns (uint256);
    function lockedLiquidBalance(address _validator) external view returns (uint256);
}