package tests

import (
	"math/big"
	"testing"

	"go.uber.org/mock/gomock"

	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/core/vm"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"

	"github.com/stretchr/testify/require"
)

var (
	offenderKey, _ = crypto.GenerateKey()
	offender       = crypto.PubkeyToAddress(offenderKey.PublicKey)
	cm             = types.CommitteeMember{Address: offender}
	header         = &types.Header{Committee: []types.CommitteeMember{cm}}
	signer         = func(hash common.Hash) ([]byte, common.Address) {
		out, _ := crypto.Sign(hash[:], offenderKey)
		return out, offender
	}
	reporter = *params.TestAutonityContractConfig.Validators[0].NodeAddress
)

func stubVerifier(address common.Address) *types.CommitteeMember {
	return &types.CommitteeMember{
		Address:     address,
		VotingPower: common.Big1,
	}
}

func NewAccusationEvent(height uint64) AccountabilityEvent {
	prevote := message.NewPrevote(0, height, common.Hash{}, signer).MustVerify(stubVerifier)

	p := &accountability.Proof{
		Type:    autonity.Accusation,
		Rule:    autonity.PVN,
		Message: prevote,
	}
	rawProof, err := rlp.EncodeToBytes(p)
	if err != nil {
		panic(err)
	}

	return AccountabilityEvent{
		EventType:      uint8(p.Type),
		Rule:           uint8(p.Rule),
		Reporter:       reporter,
		Offender:       p.Message.Sender(),
		RawProof:       rawProof,
		Id:             common.Big0,                           // assigned contract-side
		Block:          new(big.Int).SetUint64(p.Message.H()), // assigned contract-side
		ReportingBlock: common.Big0,                           // assigned contract-side
		Epoch:          common.Big0,                           // assigned contract-side
		MessageHash:    common.Big0,                           // assigned contract-side
	}
}

func TestAccusation(t *testing.T) {
	r := setup(t, nil)

	// TODO(lorenzo) Integrate this into the `setup` function
	// if possible enable snapshotting of EXPECT() on the mocks as well
	// e.g. if I do here
	// r.chainMock.EXPECT().GetHeaderByNumber(accusationHeight - 1).Return(header)
	// it will be EXPECTed for all tests

	// load the accountability precompiles into the EVM
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chainMock := accountability.NewMockChainContext(ctrl)
	chainMock.EXPECT().GetHeaderByNumber(gomock.Any()).AnyTimes().Return(header)
	accountability.LoadPrecompiles(chainMock)

	// setup current height
	currentHeight := uint64(1024)
	r.evm.Context.BlockNumber = new(big.Int).SetUint64(currentHeight)

	r.run("submit accusation at height = currentHeight - delta (valid)", func(r *runner) {
		accusationHeight := currentHeight - accountability.DeltaBlocks

		_, err := r.accountability.HandleEvent(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight))
		require.NoError(r.t, err)
	})
	r.run("submit accusation at height = currentHeight (future)", func(r *runner) {
		accusationHeight := currentHeight

		_, err := r.accountability.HandleEvent(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = currentHeight + 5 (future)", func(r *runner) {
		accusationHeight := currentHeight + 5

		_, err := r.accountability.HandleEvent(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = currentHeight - AccountabilityHeightRange (too old)", func(r *runner) {
		accusationHeight := currentHeight - accountability.AccountabilityHeightRange

		_, err := r.accountability.HandleEvent(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = currentHeight - AccountabilityHeightRange + (AccountabilityHeightRange/4)  (too old)", func(r *runner) {
		accusationHeight := currentHeight - accountability.AccountabilityHeightRange + (accountability.AccountabilityHeightRange / 4)

		_, err := r.accountability.HandleEvent(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = currentHeight - AccountabilityHeightRange + (AccountabilityHeightRange/4) + 1  (valid)", func(r *runner) {
		accusationHeight := currentHeight - accountability.AccountabilityHeightRange + (accountability.AccountabilityHeightRange / 4) + 1

		_, err := r.accountability.HandleEvent(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight))
		require.NoError(r.t, err)
	})
}
