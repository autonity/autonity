// SPDX-License-Identifier: MIT

pragma solidity ^0.8.3;

import "./interfaces/IERC20.sol";
import "./Liquid.sol";
import "./Upgradeable.sol";
import "./Precompiled.sol";
import "./Autonity.sol";

/** @title Proof-of-Stake Autonity Contract */
contract AutonityUpgradeTest is Autonity {

    constructor () Autonity(new Validator[](0), Autonity.config) {
        if (Autonity.config.contractVersion == 1) {
            _initialize();
        }
    }

    function _initialize() internal onlyProtocol {
        Autonity.validators[Autonity.validatorList[1]].bondedStake /= 2;
        Autonity.validators[Autonity.validatorList[1]].selfBondedStake /= 2;
        Autonity.config.contractVersion = 2;
        Autonity.accounts[Autonity.config.operatorAccount] = 1000;
        delete Upgradeable.newContractBytecode;
        delete Upgradeable.newContractABI;
        Upgradeable.contractUpgradeReady = false;
    }

    function _transfer(address _sender, address _recipient, uint256 _amount) internal override  {
        require(Autonity.accounts[_sender] >= _amount, "amount exceeds balance");
        Autonity.accounts[_sender] -= _amount;
        Autonity.stakeSupply += _amount;
        Autonity.accounts[_recipient] += 2 * _amount;
        emit Transfer(_sender, _recipient, _amount);
    }
}
