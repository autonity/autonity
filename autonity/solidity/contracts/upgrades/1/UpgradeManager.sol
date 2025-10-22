// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.30;

import {UpgradeManager} from "../../UpgradeManager.sol";

contract UpgradeManager1 is UpgradeManager {
    struct version {
        string number; // semver version number x.x.x
        uint256 block; // block at which the upgrade happened
    }

    // maps runtime code hash to version
    mapping(bytes32 => version) internal versionHistory;

    constructor(
        address _autonity, address _operator,
        bytes32[] memory _hashes, version[] memory _versions
    ) UpgradeManager(_autonity, _operator) {
        require(_hashes.length == _versions.length, "hashes and versions have different length");
        for(uint256 i=0; i<_hashes.length; i++) {
            versionHistory[_hashes[i]] = _versions[i];
        }
    }

    /** @dev Call the in-protocol EVM replace mechanism. Requires specific tool to interact.
    * Restricted to the operator account.
    *  @param _target is the target contract address to be updated.
    *  @param _data is the contract creation code.
    *  @param _versionString, semver string of the new version
    */
    function upgrade(address _target, string memory _data, string memory _versionString) external virtual onlyOperator {
        this.upgrade(_target, _data);

        // TODO: test to ensure that this already gets the updated codehash
        versionHistory[_target.codehash] = version(_versionString, block.number);
    }

    /**
     * allows the operator to tag the codeHash of a contract with a semver version
     * @param _hash, the contract code hash
     * @param _version, the semver version
     */
    function setVersion(bytes32 _hash, version memory _version) external virtual onlyOperator {
        versionHistory[_hash] = _version;
    }

    /**
     * allows to fetch the semver associated to a contract hash
     * @param _hash, the contract code hash
     * @return the semver version
     */
    function getVersion(bytes32 _hash) external view virtual returns (version memory) {
        return versionHistory[_hash];
    }
}
