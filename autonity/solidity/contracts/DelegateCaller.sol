// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.0;

library DelegateCaller {
    /**
     * @dev Delegates the current call to `_contractAddress`.
     * 
     * This function does not return to its internall call site, it will return directly to the external caller.
     */
    function delegate(address _contractAddress) internal {
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