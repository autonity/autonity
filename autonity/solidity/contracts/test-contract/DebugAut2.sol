// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.3;

import "../Autonity.sol";

/** @title Proof-of-Stake Autonity Contract */

contract DebugAut2 is Autonity {

    constructor(
        Validator[] memory _validators,
        Config memory _config
    ) Autonity(_validators, _config) {}

    function _performRedistribution(uint256 _atn, uint256 _ntn) internal override {
        // exit early if nothing to redistribute.
        if (_atn == 0 && _ntn == 0) {
            return;
        }
        // Take ATN treasury fee.
        uint256 _atnTreasuryReward = (config.policy.treasuryFee * _atn) / 10 ** 18;
        if (_atnTreasuryReward > 0) {
            // Using "call" to let the treasury contract do any kind of computation on receive.
            (bool sent,) = config.policy.treasuryAccount.call{value: _atnTreasuryReward}("");
            if (sent == true) {
                _atn -= _atnTreasuryReward;
            }
        }

        // first we need to reduce total _atn and _ntn by the oracle and proposer rewards
        uint256 _atnProposerRewards;
        uint256 _ntnProposerRewards;
        uint256 _atnOracleRewards = _atn * config.policy.oracleRewardRate / STANDARD_SCALE_FACTOR;
        uint256 _ntnOracleRewards = _ntn * config.policy.oracleRewardRate / STANDARD_SCALE_FACTOR;

        if (config.contracts.omissionAccountabilityContract.getTotalEffort() > 0) {
            // Calculate initial proposer rewards (actual distribution is done after regular rewards)
            _atnProposerRewards = (_atn * config.policy.proposerRewardRate * committee.length) / (STANDARD_SCALE_FACTOR * config.protocol.committeeSize);
            _ntnProposerRewards = (_ntn * config.policy.proposerRewardRate * committee.length) / (STANDARD_SCALE_FACTOR * config.protocol.committeeSize);
        }

        _atn -= _atnOracleRewards + _atnProposerRewards;
        _ntn -= _ntnOracleRewards + _ntnProposerRewards;

        uint256 _omissionScaleFactor = config.contracts.omissionAccountabilityContract.getScaleFactor();

        uint256[] memory _jailedValidatorLocs = new uint256[](committee.length);
        uint256 _jailedValidatorCount = 0;

        // Redistribute fees through the Liquid Newton contract
        uint256 _atnTotalWithheld = 0;
        uint256 _ntnTotalWithheld = 0;
        for (uint256 i = 0; i < committee.length; i++) {
            Validator storage _val = validators[committee[i].addr];
            // votingPower in the committee struct is the amount of bonded-stake pre-slashing event.
            uint256 _atnReward = (committee[i].votingPower * _atn) / epochTotalBondedStake;
            uint256 _ntnReward = (committee[i].votingPower * _ntn) / epochTotalBondedStake;
            if (_atnReward > 0 || _ntnReward > 0) {
                // committee members in the jailed state were just found guilty in the current epoch.
                // committee members in jailbound state are permanently jailed
                if (_val.state == ValidatorState.jailed || _val.state == ValidatorState.jailbound) {
                    _jailedValidatorLocs[_jailedValidatorCount] = i;
                    _jailedValidatorCount++;
                    continue;
                }

                // if jailed for inactivity, transfer all rewards to the withheld rewards pool
                if (_val.state == ValidatorState.jailedForInactivity || _val.state == ValidatorState.jailboundForInactivity) {
                    _atnTotalWithheld += _atnReward;
                    _ntnTotalWithheld += _ntnReward;
                    continue;
                }

                // rewards withholding based on omission accountability, only members with inactivity lower than InactivityThreshold will arrive here
                uint256 _inactivityScore = config.contracts.omissionAccountabilityContract.getInactivityScore(_val.nodeAddress);
                if (_inactivityScore > config.policy.withholdingThreshold) {
                    uint256 _atnWithheld = _atnReward * _inactivityScore / _omissionScaleFactor;
                    uint256 _ntnWithheld = _ntnReward * _inactivityScore / _omissionScaleFactor;

                    _atnTotalWithheld += _atnWithheld;
                    _ntnTotalWithheld += _ntnWithheld;

                    _atnReward -= _atnWithheld;
                    _ntnReward -= _ntnWithheld;
                }

                // non-jailed validators have a strict amount of bonded newton.
                // the distribution account for the PAS ratio post-slashing.
                uint256 _atnSelfReward = (_val.selfBondedStake * _atnReward) / _val.bondedStake;
                if (_atnSelfReward > 0) {
                    (bool _sent, bytes memory _returnData) = _val.treasury.call{value: _atnSelfReward, gas: 2300}("");
                    // if transfer doesn't go through (sneaky contract), just keep the amount at the autonity contract for future redistribution
                    // and let the treasury know that call failed
                    if (_sent == false) {
                        emit CallFailed(_val.treasury, "", _returnData);
                    }
                }
                uint256 _ntnSelfReward = (_val.selfBondedStake * _ntnReward) / _val.bondedStake;
                _transfer(address(this), address(config.contracts.stakingPool), _ntnSelfReward);
                _autobond(_val.nodeAddress, _ntnSelfReward, 0);

                uint256 _ntnDelegationReward = _ntnReward - _ntnSelfReward;
                uint256 _atnDelegationReward = _atnReward - _atnSelfReward;
                if (_atnDelegationReward > 0 || _ntnDelegationReward > 0) {
                    _transfer(address(this), address(_val.liquidStateContract), _ntnDelegationReward);
                    _applyNewCommissionRate(committee[i].addr);
                    _val.liquidStateContract.redistribute{value: _atnDelegationReward}(
                        accounts[address(_val.liquidStateContract)],
                        _val.commissionRate
                    );
                }
                // TODO: This has to be reconsidered - I feel it is too expensive
                // to emit an event per validator. But what is our recommend way to track rewards
                // from a user perspective then ?
                emit Rewarded(_val.nodeAddress, _atnReward, _ntnReward);
            }
        }

        // NOTE: the following redistribution operations could change validator bondedStake, therefore they must be done after
        // the rewards distribution to avoid inconsistencies in the rewards distribution.

        if (_jailedValidatorCount > 0) {
            for (uint256 i = 0; i < _jailedValidatorCount; i++) {
                uint256 _jailedValidatorIndex = _jailedValidatorLocs[i];
                uint256 _atnReward = (committee[_jailedValidatorIndex].votingPower * _atn) / epochTotalBondedStake;
                uint256 _ntnReward = (committee[_jailedValidatorIndex].votingPower * _ntn) / epochTotalBondedStake;
                _transfer(address(this), address(config.contracts.accountabilityContract), _ntnReward);
                config.contracts.accountabilityContract.distributeRewards{value: _atnReward}(committee[_jailedValidatorIndex].addr, _ntnReward);
            }
        }

        // proposer fees redistribution based on effort put into activity proofs
        // if the total effort is 0, just redistribute the proposer rewards based on stake
        // NOTE: reward forfeiting and withholding based on accountability and omission accountability are not applied to proposer rewards
        // e.g. a validator punished for equivocation will still receive his share of proposer rewards
        if (config.contracts.omissionAccountabilityContract.getTotalEffort() > 0) {
            address _omission = address(config.contracts.omissionAccountabilityContract);
            _transfer(address(this), _omission, _ntnProposerRewards);
            config.contracts.omissionAccountabilityContract.distributeProposerRewards{value: _atnProposerRewards}(accounts[_omission]);
        }

        _transfer(address(this), address(config.contracts.oracleContract), _ntnOracleRewards);
        require(false, "at oracle dist");
        config.contracts.oracleContract.distributeRewards{value: _atnOracleRewards}(
            accounts[address(config.contracts.oracleContract)]
        );

        // send withheld funds to the appropriate pool
        if (_atnTotalWithheld > 0) {
            // Using "call" to let the treasury contract do any kind of computation on receive.
            (bool _sent, bytes memory _returnData) = config.policy.withheldRewardsPool.call{value: _atnTotalWithheld}("");
            if (_sent == false) {
                emit CallFailed(config.policy.withheldRewardsPool, "", _returnData);
            }
        }
        if (_ntnTotalWithheld > 0) {
            _transfer(address(this), config.policy.withheldRewardsPool, _ntnTotalWithheld);
        }

        // let staking pool to collect and track its rewards
        // this process must be done before applying staking operations
        uint256 _committeeLength = committee.length;
        address[] memory _validators = new address[](_committeeLength);
        ILiquid[] memory _liquidContracts = new ILiquid[](_committeeLength);
        for (uint i = 0; i < _committeeLength; i++) {
            _validators[i] = committee[i].addr;
            _liquidContracts[i] = validators[_validators[i]].liquidStateContract;
        }
        config.contracts.stakingPool.collectRewards(_validators, _liquidContracts);
    }

    function finalize() external override onlyProtocol nonReentrant returns (
        bool,                       // contractUpgradeReady
        bool,                       // epochEnded
        CommitteeMember[] memory,   // committee
        uint256,                    // epochInfos[epochID].previousEpochBlock
        uint256,                    // epochInfos[epochID].nextEpochBlock
        uint256                     // delta
    ) {
        lastFinalizedBlock = block.number;
        blockEpochMap[block.number] = epochID;

        // use >= instead of == to facilitate tests on truffle
        bool _epochEnded = block.number >= epochInfos[epochID].nextEpochBlock;

        // finalize all auxiliary contracts
        config.contracts.accountabilityContract.finalize(_epochEnded);
        uint256 _delta = config.contracts.omissionAccountabilityContract.finalize(_epochEnded);
        // require(false, "after omissionAccountabilityContract.finalize");
        bool newRound = config.contracts.oracleContract.finalize();

        if (_epochEnded) {
            // We first calculate the new NTN injected supply for this epoch
            uint256 _inflationReward = config.contracts.inflationControllerContract.calculateSupplyDelta(
                stakeCirculating,
                inflationReserve,
                lastEpochTime,
                block.timestamp
            );
            if (inflationReserve < _inflationReward) {
                // If this code path is taken there is something deeply wrong happening in the inflation controller
                // contract.
                _inflationReward = inflationReserve;
            }
            // mint inflation NTN with the AC recipient
            // all rewards belong to the Autonity Contract before redistribution.
            _mint(address(this), _inflationReward);
            inflationReserve -= _inflationReward;
            stakeCirculating += _unlockSchedules(block.timestamp);
            // redistribute ATN tx fees and newly minted NTN inflation reward
            _performRedistribution(address(this).balance, accounts[address(this)]);
            // end of epoch here
            _stakingOperations();

            // compute the committee for new epoch
            (address[] memory _newOracles, address[] memory _newCommittee, address[] memory _newTreasuries) = computeCommittee();
            config.contracts.oracleContract.setVoters(_newOracles, _newTreasuries, _newCommittee);
            config.contracts.accountabilityContract.setCommittee(_newCommittee);
            config.contracts.omissionAccountabilityContract.setCommittee(committee, _newTreasuries);

            // apply new epoch period.
            config.protocol.epochPeriod = newEpochPeriod;

            // update epoch information
            config.contracts.omissionAccountabilityContract.setEpochBlock(block.number);
            uint256 _previousEpochBlock = epochInfos[epochID].epochBlock;
            uint256 _nextEpochBlock = block.number + config.protocol.epochPeriod;
            lastEpochTime = block.timestamp;

            // NOTE: Rewards distribution depends on the current value of epochID,
            // so we should always keep this epoch increment at the end of this block.
            epochID += 1;
            _addEpochInfo(epochID, EpochInfo(committee, _previousEpochBlock, block.number, _nextEpochBlock, _delta));
            emit NewEpoch(epochID);
        }

        if (newRound) {
            config.contracts.oracleContract.updateVotersAndSymbol();
            try config.contracts.acuContract.update() {}
            catch {}
        }

        return (contractUpgradeReady, _epochEnded, committee, epochInfos[epochID].previousEpochBlock, epochInfos[epochID].nextEpochBlock, _delta);
    }
}
