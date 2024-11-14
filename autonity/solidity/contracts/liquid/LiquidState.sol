// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;

import "../lib/DelegateCaller.sol";
import "./LiquidLogic.sol";
import "./LiquidStorage.sol";

contract LiquidState is LiquidStorage {
    using DelegateCaller for address;

    constructor(
        address _validator,
        address payable _treasury,
        uint256 _commissionRate,
        string memory _index,
        address _liquidLogicAddress
    ) {
        // commissionRate <= 1.0
        require(_commissionRate <= LiquidLogic(payable(_liquidLogicAddress)).COMMISSION_RATE_SCALE_FACTOR());

        validator = _validator;
        treasury = _treasury;
        commissionRate = _commissionRate;
        liquidName = string.concat("LNTN-", _index);
        liquidSymbol = string.concat("LNTN-", _index);
        autonityContract = Autonity(payable(msg.sender));
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_liquidLogicContract()`. Will run if no other
     * function in the contract matches the call data.
     */
    fallback() payable external {
        _liquidLogicContract().delegate();
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_liquidLogicContract()`. Will run if call data
     * is empty.
     */
    receive() payable external {
        _liquidLogicContract().delegate();
    }

    /**
     ============================================================

        Internals

     ============================================================
     */

    /**
     * @dev Fetch liquid logic contract address from autonity
     */
    function _liquidLogicContract() internal view returns (address) {
        address _address = autonityContract.liquidLogicContract();
        require(_address != address(0), "liquid logic contract not set");
        return _address;
    }
}
