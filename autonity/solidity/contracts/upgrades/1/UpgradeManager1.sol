// SPDX-License-Identifier: LGPL-3.0-only
pragma solidity 0.8.30;

import {IConfigEvents} from "../../interfaces/IConfigEvents.sol";
import {IUpgradeManager} from "../../interfaces/IUpgradeManager.sol";
import "../../lib/Precompiled.sol";

contract UpgradeManager1 is IConfigEvents, IUpgradeManager {
    address internal autonity;
    address internal operator;

    /** @dev added in upgrade 1 */
    struct version {
        string number; // semver version number x.x.x
        uint256 block; // block at which the upgrade happened
    }

    /** @dev added in upgrade 1
     * maps runtime bytecode hash to version
     */
    mapping(bytes32 => version) internal versionHistory;

    constructor(bytes32[] memory _hashes, version[] memory _versions) {
        require(_hashes.length == _versions.length, "hashes and versions have different length");
        for(uint256 i=0; i<_hashes.length; i++) {
            versionHistory[_hashes[i]] = _versions[i];
        }
    }
    event UpgradeResult(address indexed contractAddress, bool success);

    /** @dev Call the in-protocol EVM replace mechanism. Requires specific tool to interact.
    * Restricted to the operator account.
    *  @param _target is the target contract address to be updated.
    *  @param _data is the contract creation code.
    */
    function upgrade(address _target, string memory _data) external virtual onlyOperator {
        _upgrade(_target,_data);
    }

    function _upgrade(address _target, string memory _data) internal virtual {
        address precompile = Precompiled.UPGRADER_CONTRACT;
        bytes memory _input = abi.encodePacked(_target, _data);

        bool success;
        uint256 returnSize;
        bytes memory returnData;

        assembly {
            let result := delegatecall(gas(), precompile, add(_input, 32), mload(_input), 0, 0)
            success := result
            returnSize := returndatasize()
            //load free memory pointer
            returnData := mload(0x40)
            // update free memory pointer to point to the end of the allocated space for returnData
            mstore(0x40, add(returnData, add(returnSize, 32)))
            // write the size of the returnData first
            mstore(returnData, returnSize)
            // copy actual return data after the length
            returndatacopy(add(returnData, 32), 0, returnSize)
        }
        emit UpgradeResult(_target, success);
        assembly {
            if iszero(success) {
                revert(add(returnData, 32), returnSize)
            }
            /** This was originally implemented to return the result
            *   of the precompile execution. However:
            *   1. the upgrade precompile always returns nil if execution succeeds
            *   2. the `return` yul instruction terminates the whole execution context,
            *      not the currently executing function only. Therefore it cannot be used
            *      this way now that we call _upgrade() from other functions.
            *
            *   Leaving it commented for now, as I do not see any usefulness in bubbling up
            *   this value.
            */
            // return raw data skipping the length
            //return(add(returnData, 32), returnSize)
        }
    }

    /// @return returns the autonity contract address
    function getAutonity() external virtual view returns (address) {
        return autonity;
    }

    /// @return returns the operator address
    function getOperator() external virtual view returns (address) {
        return operator;
    }

    /*
    * @notice Set the Operator account. Restricted to the Operator account.
    * @param _account the new operator account.
    */
    function setOperator(address _account) external virtual onlyAutonity {
        emit IConfigEvents.ConfigUpdateAddress("operator", operator, _account, block.number);
        operator = _account;
    }

    /**
    * @dev Modifier that checks if the caller is the operator contract.
    */
    modifier onlyOperator  {
        require(operator == msg.sender, "caller is not the operator");
        _;
    }

    /**
    * @dev Modifier that checks if the caller is the Autonity Contract
    */
    modifier onlyAutonity  {
        require(autonity == msg.sender, "caller is not the Autonity contract");
        _;
    }


    /* ----------------------------------
     * UPGRADE 1 new functions
     * ----------------------------------
     */

    /** @dev Call the in-protocol EVM replace mechanism. Requires specific tool to interact.
    * Restricted to the operator account.
    *  @param _target is the target contract address to be updated.
    *  @param _data is the contract creation code.
    *  @param _versionString, semver string of the new version.
    */
    function upgrade(address _target, string memory _data, string memory _versionString) external virtual onlyOperator {
        _upgrade(_target,_data);
        versionHistory[_target.codehash] = version(_versionString, block.number);
    }

    /** @dev Call the in-protocol EVM replace mechanism. Requires specific tool to interact.
    * Restricted to the operator account.
    *  @param _targets are the target contracts addresses to be updated.
    *  @param _bytecodes are the contracts creation code.
    */
    function upgradeMultiple(address[] memory _targets, string[] memory _bytecodes) external virtual onlyOperator {
        require(_targets.length == _bytecodes.length, "addresses and bytecodes should be in same number");
        for(uint256 i=0; i<_targets.length; i++) {
            _upgrade(_targets[i],_bytecodes[i]);
        }
    }

    /** @dev Call the in-protocol EVM replace mechanism. Requires specific tool to interact.
    * Restricted to the operator account.
    *  @param _targets are the target contracts addresses to be updated.
    *  @param _bytecodes are the contracts creation code.
    *  @param _versionStrings, semver strings of the new versions.
    */
    function upgradeMultiple(address[] memory _targets, string[] memory _bytecodes, string[] memory _versionStrings) external virtual onlyOperator {
        require(_targets.length == _bytecodes.length, "addresses and bytecodes should be in same number");
        require(_targets.length == _versionStrings.length, "addresses and version strings should be in same number");
        for(uint256 i=0; i<_targets.length; i++) {
            _upgrade(_targets[i],_bytecodes[i]);
            versionHistory[_targets[i].codehash] = version(_versionStrings[i], block.number);
        }
    }

    /**
     * allows the operator to tag the codeHash of a contract with a semver version
     * @param _hash, the contract code hash
     * @param _number, semver version number x.x.x
     * @param _block, block at which the upgrade happened
     */
    function setVersion(bytes32 _hash, string memory _number, uint256 _block) external virtual onlyOperator {
        versionHistory[_hash] = version(_number,_block);
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
