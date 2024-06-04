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
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/params"
	"github.com/autonity/autonity/rlp"

	"github.com/stretchr/testify/require"
)

var (
	offenderNodeKey, _      = crypto.GenerateKey()
	offenderConsensusKey, _ = blst.RandKey()
	offender                = crypto.PubkeyToAddress(offenderNodeKey.PublicKey)
	cm                      = types.CommitteeMember{Address: offender, VotingPower: common.Big1, ConsensusKey: offenderConsensusKey.PublicKey(), ConsensusKeyBytes: offenderConsensusKey.PublicKey().Marshal(), Index: 0}
	header                  = &types.Header{Committee: []types.CommitteeMember{cm}}
	signer                  = func(hash common.Hash) blst.Signature {
		return offenderConsensusKey.Sign(hash[:])
	}
	reporter = *params.TestAutonityContractConfig.Validators[0].NodeAddress
)

func NewAccusationEvent(height uint64, value common.Hash) AccountabilityEvent {
	prevote := message.NewPrevote(0, height, value, signer, &cm, 1)

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
		Offender:       offender,
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

	// load the accountability precompiles into the EVM
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	chainMock := accountability.NewMockChainContext(ctrl)
	chainMock.EXPECT().GetHeaderByNumber(gomock.Any()).AnyTimes().Return(header)
	accountability.LoadPrecompiles(chainMock)

	// setup current height
	currentHeight := uint64(1024)
	r.evm.Context.BlockNumber = new(big.Int).SetUint64(currentHeight)
	lastCommittedHeight := currentHeight - 1

	// TODO(lorenzo) add similar tests for PVO and C1
	/*
		r.run("PVN accusation with prevote nil should revert", func(r *runner) {
			accusationHeight := lastCommittedHeight - accountability.DeltaBlocks

			chainMock.EXPECT().GetBlock(common.Hash{}, accusationHeight).Return(nil)

			_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{}, reporter))
			require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
		})
		r.run("accusation for committed value should revert", func(r *runner) {
			accusationHeight := lastCommittedHeight - accountability.DeltaBlocks

			chainMock.EXPECT().GetBlock(common.Hash{0xca, 0xfe}, accusationHeight).Return(&types.Block{})

			_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}, reporter))
			require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
		})*/
	r.run("reporting right tests", func(r *runner) {

		// reporting should be reverted since reporter is not in current committee and last committee
		accusationHeight := lastCommittedHeight - accountability.HeightRange + (accountability.HeightRange / 4) + 1
		noAccessor := common.Address{}
		_, err := r.accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
		require.Equal(r.t, "execution reverted: function restricted to a committee member", err.Error())

		// set committee with reporter
		committee, _, err := r.autonity.GetCommittee(nil)
		require.NoError(t, err)
		var newCommittee []common.Address
		for _, c := range committee {
			newCommittee = append(newCommittee, c.Addr)
		}
		newCommittee = append(newCommittee, noAccessor)

		// set the new committee that contains reporter account, then it is allowed for reporting.
		_, err = r.accountability.SetCommittee(&runOptions{origin: params.AutonityContractAddress}, newCommittee)
		require.NoError(t, err)
		chainMock.EXPECT().GetBlock(common.Hash{0xca, 0xff}, accusationHeight).Return(nil)
		_, err = r.accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}))
		require.NoError(t, err)

		// now set new committee, it will set current committee as last committee, the reporter is still allowed for reporting.
		_, err = r.accountability.SetCommittee(&runOptions{origin: params.AutonityContractAddress}, newCommittee[0:len(newCommittee)-1])
		require.NoError(t, err)
		chainMock.EXPECT().GetBlock(common.Hash{0xca, 0xff}, accusationHeight).Return(nil)
		// report same accusation should be reverted since the accusation is pending now.
		_, err = r.accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}))
		require.Equal(t, "execution reverted: already processing an accusation", err.Error())

		// nwo set new committee without having the reporter, then it is not allowed for reporting.
		_, err = r.accountability.SetCommittee(&runOptions{origin: params.AutonityContractAddress}, newCommittee[0:len(newCommittee)-1])
		require.NoError(t, err)
		// report same accusation should be reverted since the accusation is pending now.
		_, err = r.accountability.HandleAccusation(&runOptions{origin: noAccessor}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xff}))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
		require.Equal(r.t, "execution reverted: function restricted to a committee member", err.Error())
	})
}

func TestAccusationTiming(t *testing.T) {
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
	// set value to not committed for all tests, here we want to really tests only the height related checks
	chainMock.EXPECT().GetBlock(gomock.Any(), gomock.Any()).AnyTimes().Return(nil)
	accountability.LoadPrecompiles(chainMock)

	currentHeight := uint64(1024) // height of current consensus instance
	r.evm.Context.BlockNumber = new(big.Int).SetUint64(currentHeight)
	lastCommittedHeight := currentHeight - 1 // height of last committed block

	r.run("submit accusation at height = lastCommittedHeight - delta (valid)", func(r *runner) {
		accusationHeight := lastCommittedHeight - accountability.DeltaBlocks

		_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}))
		require.NoError(r.t, err)
	})
	r.run("submit accusation at height = lastCommittedHeight - delta + 1 (too recent)", func(r *runner) {
		accusationHeight := lastCommittedHeight - accountability.DeltaBlocks + 1

		_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = lastCommittedHeight (too recent)", func(r *runner) {
		accusationHeight := lastCommittedHeight

		_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = lastCommittedHeight + 5 (future)", func(r *runner) {
		accusationHeight := lastCommittedHeight + 5

		_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange (too old)", func(r *runner) {
		accusationHeight := lastCommittedHeight - accountability.HeightRange

		_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange + (AccountabilityHeightRange/4)  (too old)", func(r *runner) {
		accusationHeight := lastCommittedHeight - accountability.HeightRange + (accountability.HeightRange / 4)

		_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}))
		require.ErrorIs(r.t, err, vm.ErrExecutionReverted)
	})
	r.run("submit accusation at height = lastCommittedHeight - AccountabilityHeightRange + (AccountabilityHeightRange/4) + 1  (valid)", func(r *runner) {
		accusationHeight := lastCommittedHeight - accountability.HeightRange + (accountability.HeightRange / 4) + 1

		_, err := r.accountability.HandleAccusation(&runOptions{origin: reporter}, NewAccusationEvent(accusationHeight, common.Hash{0xca, 0xfe}))
		require.NoError(r.t, err)
	})
}
