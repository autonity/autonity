package autonity

import (
	"errors"
	"fmt"
	"github.com/autonity/autonity/autonity/bindings"
	"math/big"
	"reflect"

	"github.com/autonity/autonity/accounts/abi"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/math"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/metrics"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/params/generated"
)

var (
	finalizeGas      = metrics.NewRegisteredBufferedGauge("autonity/finalize", nil, nil)
	epochFinalizeGas = metrics.NewRegisteredBufferedGauge("autonity/epoch/finalize", nil, nil)
)

type raw []byte

func (c *evmContract) replaceAutonityBytecode(header *types.Header, statedb vm.StateDB, bytecode []byte) error {
	evm := c.evmProvider(header, params.DeployerAddress, statedb)
	_, _, _, vmerr := evm.Replace(vm.AccountRef(params.DeployerAddress), bytecode, params.AutonityContractAddress)
	if vmerr != nil {
		log.Error("replaceAutonityBytecode evm.Create", "err", vmerr)
		return vmerr
	}
	return nil
}

type errorWithRevertReason struct {
	error
	reason string
}

func (e *errorWithRevertReason) Error() string {
	return e.error.Error() + ": " + e.reason
}

func newErrorWithRevertReason(err error, ret []byte) error {
	if err == nil || !errors.Is(err, vm.ErrExecutionReverted) {
		return err
	}
	reason, unpackErr := abi.UnpackRevert(ret)
	if unpackErr != nil {
		return err // if we cannot unpack the revert reason, return the original error unmodified
	}
	return &errorWithRevertReason{err, reason}
}

// *
// Exposed package functions
// *

// AutonityContractCall calls the specified function of the autonity contract
// with the given args, and returns the output unpacked into the result
// interface.
// It returns the gas used for the call
//
//revive:disable:exported - Autonity is one of the contracts, so repetitive naming here is justified
func AutonityContractCall(autonityAbi *abi.ABI, evm *vm.EVM, function string, result any, args ...any) (uint64, error) {
	packedArgs, err := autonityAbi.Pack(function, args...)
	if err != nil {
		return 0, err
	}
	ret, usedGas, err := evm.Call(
		vm.AccountRef(params.DeployerAddress),
		params.AutonityContractAddress,
		packedArgs,
		math.MaxUint64,
		common.Big0,
	)
	if err != nil {
		return usedGas, newErrorWithRevertReason(err, ret)
	}
	// if result's type is "raw" then bypass unpacking
	if reflect.TypeOf(result) == reflect.TypeOf(&raw{}) {
		rawPtr := result.(*raw)
		*rawPtr = ret
		return usedGas, nil
	}
	if err := autonityAbi.UnpackIntoInterface(result, function, ret); err != nil {
		log.Error("Could not unpack returned value", "function", function)
		return usedGas, err
	}

	return usedGas, nil
}

func CallGetCommittee(evm *vm.EVM) (*types.Committee, error) {
	var committeeMembers []types.CommitteeMember
	if _, err := AutonityContractCall(&generated.AutonityAbi, evm, "getCommittee", &committeeMembers); err != nil {
		return nil, err
	}
	committee := &types.Committee{}
	committee.Members = committeeMembers
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error()) //nolint
	}
	return committee, nil
}

func (c *AutonityContract) CallGetCommitteeEnodes(state vm.StateDB, header *types.Header, asACN bool) (*types.Nodes, error) {
	var returnedEnodes []string
	_, err := AutonityContractCall(c.contractABI, c.evmProvider(header, params.DeployerAddress, state), "getCommitteeEnodes", &returnedEnodes)
	if err != nil {
		return nil, err
	}
	return types.NewNodes(returnedEnodes, asACN), nil
}

func (c *AutonityContract) CallConfig(state vm.StateDB, header *types.Header) (*bindings.AutonityConfig, error) {
	var config bindings.AutonityConfig
	_, err := AutonityContractCall(
		c.contractABI,
		c.evmProvider(header, params.DeployerAddress, state),
		"config",
		&config,
	)
	return &config, err
}

func (c *AutonityContract) CallEpochID(state vm.StateDB, header *types.Header) (*big.Int, error) {
	epochID := new(big.Int)
	_, err := AutonityContractCall(
		c.contractABI,
		c.evmProvider(header, params.DeployerAddress, state),
		"epochID",
		&epochID,
	)
	if err != nil {
		return nil, err
	}
	return epochID, nil
}

// CallEpochByHeight get the epoch by height.
func (c *AutonityContract) CallEpochByHeight(state vm.StateDB, header *types.Header, height *big.Int) (*types.EpochInfo, error) {
	var output raw
	if _, err := AutonityContractCall(
		c.contractABI,
		c.evmProvider(header, params.DeployerAddress, state),
		"getEpochByHeight",
		&output,
		height,
	); err != nil {
		return nil, err
	}

	data, err := c.contractABI.Unpack("getEpochByHeight", output)
	if err != nil {
		return nil, err
	}

	info := *abi.ConvertType(data[0], new(bindings.AutonityEpochInfo)).(*bindings.AutonityEpochInfo)

	committee := &types.Committee{}
	for _, member := range info.Committee {
		committee.Members = append(committee.Members, types.CommitteeMember{
			Address:           member.Addr,
			VotingPower:       member.VotingPower,
			ConsensusKeyBytes: member.ConsensusKey,
		})
	}
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error()) //nolint
	}

	epochInfo := &types.EpochInfo{
		Epoch: types.Epoch{
			Committee:          committee,
			PreviousEpochBlock: info.PreviousEpochBlock,
			NextEpochBlock:     info.NextEpochBlock,
			OmissionDelta:      info.OmissionDelta,
		},
		EpochBlock: info.EpochBlock,
	}

	return epochInfo, nil
}

// *
// Internal package functions
// *

func (c *AutonityContract) callGetMinimumBaseFee(state vm.StateDB, header *types.Header) (*big.Int, error) {
	minBaseFee := new(big.Int)
	_, err := AutonityContractCall(
		c.contractABI,
		c.evmProvider(header, params.DeployerAddress, state),
		"getMinimumBaseFee",
		&minBaseFee,
	)
	if err != nil {
		return nil, err
	}
	return minBaseFee, nil
}

func (c *AutonityContract) callGetEpochPeriod(state vm.StateDB, header *types.Header) (*big.Int, error) {
	epochPeriod := new(big.Int)
	_, err := AutonityContractCall(c.contractABI, c.evmProvider(header, params.DeployerAddress, state), "getEpochPeriod", &epochPeriod)
	if err != nil {
		return nil, err
	}
	return epochPeriod, nil
}

func (c *AutonityContract) callFinalize(state vm.StateDB, header *types.Header) (bool, *types.Epoch, error) {
	var output raw

	usedGas, err := AutonityContractCall(
		c.contractABI,
		c.evmProvider(header, params.DeployerAddress, state),
		"finalize",
		&output,
	)
	if err != nil {
		return false, nil, fmt.Errorf("call finalize failed: %w", err)
	}

	unpackedOutput, err := c.contractABI.Unpack("finalize", output)
	if err != nil {
		return false, nil, fmt.Errorf("unpacking raw finalize output failed: %w", err)
	}

	result := *abi.ConvertType(unpackedOutput[0], new(bindings.AutonityFinalizeResult)).(*bindings.AutonityFinalizeResult)

	recordFinalizeGasUsage(result.EpochEnded, header.Number.Uint64(), int64(usedGas))

	// set current config in the statedb object
	state.SetContractsConfig(&(result.Config))

	if !result.EpochEnded {
		return result.ContractUpgradeReady, nil, nil
	}

	// return with epoch info
	committee := &types.Committee{}
	committee.Members = make([]types.CommitteeMember, len(result.Committee))
	for i, member := range result.Committee {
		committee.Members[i] = types.CommitteeMember{
			Address:           member.Addr,
			VotingPower:       member.VotingPower,
			ConsensusKeyBytes: member.ConsensusKey,
		}
	}
	if err := committee.Enrich(); err != nil {
		panic("Committee member has invalid consensus key: " + err.Error())
	}

	epoch := &types.Epoch{
		PreviousEpochBlock: result.PreviousEpochBlock,
		NextEpochBlock:     result.NextEpochBlock,
		Committee:          committee,
		OmissionDelta:      result.OmissionDelta,
	}

	return result.ContractUpgradeReady, epoch, nil
}

func (c *AutonityContract) callRetrieveContract(state vm.StateDB, header *types.Header) ([]byte, string, error) {
	var bytecode []byte
	var updateAbi string
	if _, err := AutonityContractCall(
		c.contractABI,
		c.evmProvider(header, params.DeployerAddress, state),
		"getNewContract",
		&[]any{&bytecode, &updateAbi},
	); err != nil {
		return nil, "", err
	}
	return bytecode, updateAbi, nil
}

func recordFinalizeGasUsage(isEpochHeader bool, number uint64, usedGas int64) {
	if isEpochHeader {
		log.Debug("gas used to finalize epoch block", "number", number, "usedGas", usedGas)
		if metrics.Enabled {
			epochFinalizeGas.Add(usedGas)
		}
	} else {
		log.Debug("gas used to finalize block", "number", number, "usedGas", usedGas)
		if metrics.Enabled {
			finalizeGas.Add(usedGas)
		}
	}
}
