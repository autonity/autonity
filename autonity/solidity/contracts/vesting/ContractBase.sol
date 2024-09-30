// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

import "../AccessAutonity.sol";

abstract contract ContractBase is AccessAutonity {

    struct Contract {
        uint256 currentNTNAmount;
        uint256 withdrawnValue;
        uint256 start;
        uint256 cliffDuration;
        uint256 totalDuration;
        bool canStake;
    }

    /*
    ============================================================
         Internals
    ============================================================
     */

    function _createContract(
        uint256 _amount,
        uint256 _startTime,
        uint256 _cliffDuration,
        uint256 _totalDuration,
        bool _canStake
    ) internal pure returns (Contract memory) {

        require(_totalDuration > _cliffDuration, "end must be greater than cliff");
        return Contract(
            _amount, 0, _startTime, _cliffDuration, _totalDuration, _canStake
        );
    }

    function _releaseNTN(
        Contract storage _contract, uint256 _amount
    ) internal returns (uint256 _remaining) {
        if (_amount > _contract.currentNTNAmount) {
            _remaining = _amount - _contract.currentNTNAmount;
            _updateAndTransferNTN(_contract, msg.sender, _contract.currentNTNAmount);
        }
        else if (_amount > 0) {
            _updateAndTransferNTN(_contract, msg.sender, _amount);
        }
    }

    /**
     * @dev Given the total value (in NTN) of the contract, calculates the amount of withdrawable tokens (in NTN).
     */
    function _calculateAvailableUnlockedFunds(
        Contract storage _contract, uint256 _totalValue, uint256 _time
    ) internal view returns (uint256) {
        require(_time >= _contract.start + _contract.cliffDuration, "cliff period not reached yet");

        uint256 _unlocked = _calculateTotalUnlockedFunds(_contract.start, _contract.totalDuration, _time, _totalValue);
        if (_unlocked > _contract.withdrawnValue) {
            return _unlocked - _contract.withdrawnValue;
        }
        return 0;
    }

    /**
     * @dev Calculates total unlocked funds while assuming cliff period has passed.
     * Check if cliff is passed before calling this function.
     */
    function _calculateTotalUnlockedFunds(
        uint256 _start, uint256 _totalDuration, uint256 _time, uint256 _totalAmount
    ) internal pure returns (uint256) {
        if (_time >= _totalDuration + _start) {
            return _totalAmount;
        }
        return (_totalAmount * (_time - _start)) / _totalDuration;
    }

    /**
     * @dev Updates the contract with `contractID` and transfers NTN.
     */
    function _updateAndTransferNTN(Contract storage _contract, address _to, uint256 _amount) internal {
        _contract.currentNTNAmount -= _amount;
        _contract.withdrawnValue += _amount;
        _transferNTN(_to, _amount);
    }

    function _transferNTN(address _to, uint256 _amount) internal {
        bool _sent = autonity.transfer(_to, _amount);
        require(_sent, "NTN not transferred");
    }

    /*
    ============================================================
         Getters
    ============================================================
     */

    /**
     * @notice Returns if beneficiary can stake from his contract.
     */
    function canStake(Contract storage _contract) internal view returns (bool) {
        return _contract.canStake;
    }
}