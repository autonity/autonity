// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../interfaces/IERC20.sol";

library StakingLib {
    struct Contract {
        uint256 currentNTNAmount;
        uint256 withdrawnValue;
        uint256 start;
        uint256 cliffDuration;
        uint256 totalDuration;
        bool canStake;
    }


    function createContract(
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration,
        bool _canStake
    ) internal pure returns (Contract memory) {
        require(_totalDuration > _cliffDuration, "end must be greater than cliff");
        return StakingLib.Contract(
            _amount, 0, _startTime, _cliffDuration, _totalDuration, _canStake
        );
    }

    /**
    * Struct functions
    */

    function releaseToken(
        Contract storage _contract, address _token, uint256 _amount
    ) internal returns (uint256 _remaining) {
        if (_amount > _contract.currentNTNAmount) {
            _remaining = _amount - _contract.currentNTNAmount;
            _updateAndTransferToken(_contract, _token, msg.sender, _contract.currentNTNAmount);
        }
        else if (_amount > 0) {
            _updateAndTransferToken(_contract, _token, msg.sender, _amount);
        }
    }


    function canStake(Contract storage _contract) internal view returns (bool) {
        return _contract.canStake;
    }

    /**
    * @dev Updates the contract with and transfers NTN.
    */
    function _updateAndTransferToken(Contract storage _contract, address _token, address _to, uint256 _amount) internal {
        _contract.currentNTNAmount -= _amount;
        _contract.withdrawnValue += _amount;
        transferToken(_token, _to, _amount);

    }

    /**
    * @dev Transfers a token and checks the return value.
    */
    function transferToken(address _token, address _to, uint256 _amount) internal {
        bool _sent = IERC20(_token).transfer(_to, _amount);
        require(_sent, "token not transferred");
    }
}