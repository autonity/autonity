package accountability

import (
	"crypto/ecdsa"
	"math/big"
	"math/rand"
	"testing"

	"github.com/autonity/autonity/accounts/abi/bind/backends"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/events"
	ccore "github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/crypto/blst"
	"github.com/autonity/autonity/event"
	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/params"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestPrimaryIndex(t *testing.T) {
	tests := []struct {
		name          string
		height        uint64
		committeeSize uint64
		expected      uint64
	}{
		// Basic functionality
		{"Block 0", 0, 100, 0},
		{"Block 20", 20, 100, 1},
		{"Block 40", 40, 100, 2},
		{"Block 60", 60, 100, 3},

		// Boundary conditions
		{"Height 19", 19, 100, 0},
		{"Height 20", 20, 100, 1},
		{"Height 21", 21, 100, 1},
		{"Height 39", 39, 100, 1},
		{"Height 40", 40, 100, 2},

		// Committee size boundaries
		{"Committee size 1", 100, 1, 0},
		{"Committee size 2", 60, 2, 1},
		{"Committee size 3", 60, 3, 3 % 3}, // 60/20 = 3 → 3%3=0

		// Rollover cases
		{"Rollover 1", 200, 10, 10 % 10},
		{"Rollover 2", 220, 10, 11 % 10},
		{"Rollover 3", 240, 10, 12 % 10},

		// Large numbers
		{"Large committee", 100, 1000, (100 / 20) % 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := primaryIndex(tt.height, tt.committeeSize)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// =====================
// onDutyDetector Tests
// =====================

func TestOnDutyDetector(t *testing.T) {

	tests := []struct {
		name          string
		height        uint64
		committeeSize uint64
		clientIndex   uint64
		expected      bool
	}{
		// Small scale networks (<= SmallScaleNetSize)
		{"SmallNet-ClientIn-0", uint64(rand.Int()), 5, 0, true},
		{"SmallNet-ClientIn-5", uint64(rand.Int()), 10, 5, true},
		// Edge case: committeeSize = SmallScaleNetSize (boundary)
		{"SmallNet-ClientIn-9", uint64(rand.Int()), SmallScaleNetSize, 9, true},
		{"SmallNet-ClientIn-31", uint64(rand.Int()), SmallScaleNetSize, 31, true},

		// Primary reporter selection cases
		{"Primary-Block0", 0, 100, 0, true},
		{"Primary-Block20", reportingSlotPeriod * 1, 100, 1, true},
		{"Primary-Block40", reportingSlotPeriod * 2, 100, 2, true},
		{"Primary-Block40", reportingSlotPeriod * 4, 100, 4, true},
		// Rotations
		{"Rotations NotPrimary 99", reportingSlotPeriod * 100, 100, 99, false},
		{"Rotations Primary 0", reportingSlotPeriod * 100, 100, 0, true},

		// Non-wrapping backup cases
		{"Backup1/100", 1, 100, primaryIndex(1, 100) + 1, true},
		{"Backup10/100", 1, 100, primaryIndex(1, 100) + 10, true},
		{"Backup20/100", 1, 100, primaryIndex(1, 100) + 20, true},
		{"Backup30/100", 1, 100, primaryIndex(1, 100) + 30, true},
		{"Backup33/100", 1, 100, primaryIndex(1, 100) + 33, true}, // F == 33

		{"NoBackup34/100", 1, 100, primaryIndex(1, 100) + 34, false}, // Node 34 should not on duty.
		{"NoBackup50/100", 1, 100, primaryIndex(1, 100) + 50, false},
		{"NoBackup99/100", 1, 100, primaryIndex(1, 100) + 99, false},

		// Wrappings
		{"WrappedBackupLowBoundary", reportingSlotPeriod * 80, 100, 79, false},
		{"Non-wrapping Primary", reportingSlotPeriod * 80, 100, 80, true},
		{"Non-Wrapping Backup1", reportingSlotPeriod * 80, 100, 81, true},
		{"Non-Wrapping Backup19", reportingSlotPeriod * 80, 100, 99, true},
		{"WrappedBackup0", reportingSlotPeriod * 80, 100, 0, true},
		{"WrappedBackup13 Boundary", reportingSlotPeriod * 80, 100, 13, true},
		{"Not backup", reportingSlotPeriod * 80, 100, 14, false},

		{"Non-wrapping Primary", reportingSlotPeriod * 98, 100, 97, false},
		{"Non-wrapping Primary", reportingSlotPeriod * 98, 100, 98, true},
		{"Non-wrapping Primary", reportingSlotPeriod * 98, 100, 99, true},
		{"Non-wrapping Primary", reportingSlotPeriod * 98, 100, 0, true},

		{"Non-wrapping Primary", reportingSlotPeriod * 99, 100, 99, true},
		{"Non-wrapping Primary", reportingSlotPeriod * 99, 100, 0, true},
		{"Non-wrapping Primary", reportingSlotPeriod * 99, 100, 32, true},
		{"Non-wrapping Primary", reportingSlotPeriod * 99, 100, 33, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			fd := createFD(t, ctrl, tt.committeeSize, tt.clientIndex, tt.height)
			t.Log("height ", tt.height)
			result := fd.onDutyDetector(tt.height)
			assert.Equal(t, tt.expected, result,
				"Height: %d, Committee: %d, Client: %d",
				tt.height, tt.committeeSize, tt.clientIndex)
		})
	}
}

func createFD(t *testing.T, ctrl *gomock.Controller, committeeSize, clientIndex, height uint64) *FaultDetector {
	com := genCommittee(committeeSize)
	chainMock := NewMockChainContext(ctrl)
	chainMock.EXPECT().CommitteeByHeight(height).Return(com, nil)
	var blockSub event.Subscription
	chainMock.EXPECT().SubscribeChainEvent(gomock.Any()).AnyTimes().Return(blockSub)
	chainMock.EXPECT().Config().AnyTimes().Return(&params.ChainConfig{ChainID: common.Big1})

	fdAddr := com.Members[clientIndex].Address
	accountability, _ := bindings.NewAccountability(proposer, backends.NewSimulatedBackend(ccore.GenesisAlloc{
		fdAddr: ccore.GenesisAccount{
			Balance: big.NewInt(params.Ether),
		},
	}, 10000000))

	afdDispatchCh := make(chan events.MessageEventer, 100)
	fd := NewFaultDetector(chainMock, fdAddr, nil, core.NewMsgStore(), nil, nil, proposerNodeKey, &autonity.ProtocolContracts{Accountability: accountability}, afdDispatchCh, log.Root())
	return fd
}

func genCommittee(size uint64) *types.Committee {
	n := size
	c := new(types.Committee)
	pkeys := make([]*ecdsa.PrivateKey, n)
	consensusKeys := make([]blst.SecretKey, n)
	for i := uint64(0); i < n; i++ {
		privateKey, _ := crypto.GenerateKey()
		consensusKey, _ := blst.RandKey()
		committeeMember := types.CommitteeMember{
			Address:           crypto.PubkeyToAddress(privateKey.PublicKey),
			VotingPower:       new(big.Int).SetUint64(1),
			ConsensusKey:      consensusKey.PublicKey(),
			ConsensusKeyBytes: consensusKey.PublicKey().Marshal(),
			Index:             i,
		}
		c.Members = append(c.Members, committeeMember)
		pkeys[i] = privateKey
		consensusKeys[i] = consensusKey
	}
	return c
}
