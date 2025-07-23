// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity 0.8.30;

import "./lib/BytesLib.sol";

/** @title Implementation of the Autonity contract-upgrade functionality. */
abstract contract Upgradeable {

    bytes internal newContractBytecode;
    string internal newContractABI;
    bool internal contractUpgradeReady;

    /**
    * @notice Append to the contract storage buffer the new contract bytecode and abi.
    * Should be called as many times as required.
    * IMPORTANT: Should not modify config, in order to avoid invalidating the client-side cache.
    */
    function _upgradeContract(bytes memory _bytecode, string memory _abi) internal virtual {
        BytesLib.concatStorage(newContractBytecode, _bytecode);
        BytesLib.concatStorage(bytes(newContractABI), bytes(_abi));
    }

    /**
    * @notice Finalize the contract upgrade.
    * To be called once the storage buffer for the new contract are filled using {upgradeContract}
    * The protocol will then update the bytecode of the autonity contract at block finalization phase.
    */
    function _completeContractUpgrade() internal virtual {
        contractUpgradeReady = true;
    }

    /**
    * @notice Reset internal storage contract-upgrade buffers in case of issue.
    */
    function _resetContractUpgrade() internal virtual {
        delete newContractBytecode;
        delete newContractABI;
        contractUpgradeReady = false;
    }

    /**
     * @notice Getter to retrieve a new Autonity contract bytecode and ABI when an upgrade is initiated.
     * @return `bytecode` the new contract bytecode.
     * @return `contractAbi` the new contract ABI.
     */
    function _getNewContract() internal view virtual returns (bytes memory, string memory) {
        return (newContractBytecode, newContractABI);
    }
}
