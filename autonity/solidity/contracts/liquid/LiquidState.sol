// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;
import "./LiquidLogic.sol";
import "./LiquidStorage.sol";

contract LiquidState is LiquidStorage {

    constructor(
        address _validator,
        address payable _treasury,
        uint256 _commissionRate,
        string memory _index,
        address _liquidLogicAddress
    ) {
        // commissionRate <= 1.0
        require(_commissionRate <= LiquidLogic(payable(_liquidLogicAddress)).COMMISSION_RATE_PRECISION());

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
        _delegate(
            _liquidLogicContract()
        );
    }

    /**
     * @dev Fallback function that delegates calls to the address returned by `_liquidLogicContract()`. Will run if call data
     * is empty.
     */
    receive() payable external {
        _delegate(
            _liquidLogicContract()
        );
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

    /**
     * @dev Delegates the current call to `_contractAddress`.
     * 
     * This function does not return to its internall call site, it will return directly to the external caller.
     */
    function _delegate(address _contractAddress) internal {
        // solhint-disable-next-line no-inline-assembly
        assembly {
            // Copy msg.data. We take full control of memory in this inline assembly
            // block because it will not return to Solidity code. We overwrite the
            // Solidity scratch pad at memory position 0.
            calldatacopy(0, 0, calldatasize())

            // Call the implementation.
            // out and outsize are 0 because we don't know the size yet.
            let result := delegatecall(gas(), _contractAddress, 0, calldatasize(), 0, 0)

            // Copy the returned data.
            returndatacopy(0, 0, returndatasize())

            if iszero(result) {
                revert(0, returndatasize())
            }
            return(0, returndatasize())
        }
    }
}
