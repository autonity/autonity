// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.30;

import {UpgradeManager} from "../../UpgradeManager.sol";

contract UpgradeManager1 is UpgradeManager {
    mapping(bytes32 => string) internal versionHistory;

    constructor(
        address _autonity, address _operator,
        bytes32[] memory _hashes, string[] memory _versions
    ) UpgradeManager(_autonity, _operator) {
        require(_hashes.length == _versions.length, "hashes and versions have different length");
        for(uint256 i=0; i<_hashes.length; i++) {
            versionHistory[_hashes[i]] = _versions[i];
        }
    }

    /**
     * allows the operator to tag the codeHash of a contract with a semver version
     * @param _hash, the contract code hash
     * @param _version, the semver version
     */
    function setVersion(bytes32 _hash, string memory _version) external virtual onlyOperator {
        versionHistory[_hash] = _version;
    }

    /**
     * allows to fetch the semver associated to a contract hash
     * @param _hash, the contract code hash
     * @return the semver version
     */
    function getVersion(bytes32 _hash) external view virtual returns (string memory) {
        return versionHistory[_hash];
    }
}
