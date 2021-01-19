package core

import (
	context "context"
	"math/big"
	"reflect"
	"testing"
	"time"

	"github.com/clearmatics/autonity/consensus"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
)

func TestStoreUnminedBlockMsg(t *testing.T) {
	t.Run("old height unminedBlock", func(t *testing.T) {
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		c := &core{
			logger:           log.New("backend", "test", "id", 0),
			messages:         messages,
			height:           big.NewInt(4),
			round:            2,
			step:             prevote,
			curRoundMessages: curRoundMessages,
		}

		unminedBlock := types.NewBlockWithHeader(&types.Header{})
		c.storeUnminedBlockMsg(context.Background(), unminedBlock)

		if s := len(c.pendingUnminedBlocks); s > 0 {
			t.Fatalf("Unmined blocks size must be 0, got %d", s)
		}
	})

	t.Run("valid block given, block is stored", func(t *testing.T) {
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		c := &core{
			logger:               log.New("backend", "test", "id", 0),
			round:                2,
			height:               big.NewInt(4),
			messages:             messages,
			curRoundMessages:     curRoundMessages,
			pendingUnminedBlocks: make(map[uint64]*types.Block),
		}

		unminedBlock := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(4)})
		c.storeUnminedBlockMsg(context.Background(), unminedBlock)

		if s := len(c.pendingUnminedBlocks); s != 1 {
			t.Fatalf("Unmined blocks size must be 1, got %d", s)
		}
	})
}

func TestUpdatePendingUnminedBlocks(t *testing.T) {
	t.Run("old pending blocks removed, new block added", func(t *testing.T) {
		anOldBlock := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		unminedBlocks := make(map[uint64]*types.Block)
		unminedBlocks[anOldBlock.NumberU64()] = anOldBlock
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		c := &core{
			round:                2,
			height:               big.NewInt(3),
			curRoundMessages:     curRoundMessages,
			pendingUnminedBlocks: unminedBlocks,
		}
		unminedBlock := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(3)})
		c.updatePendingUnminedBlocks(context.Background(), unminedBlock)

		if s := len(c.pendingUnminedBlocks); s != 1 {
			t.Fatalf("Unmined blocks size must be 1, got %d", s)
		}
	})

	t.Run("wait for unmined block, new block added", func(t *testing.T) {
		pendingUnminedBlockCh := make(chan *types.Block, 1)
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		c := &core{
			curRoundMessages:         curRoundMessages,
			messages:                 messages,
			round:                    2,
			height:                   big.NewInt(3),
			pendingUnminedBlocks:     make(map[uint64]*types.Block),
			pendingUnminedBlockCh:    pendingUnminedBlockCh,
			isWaitingForUnminedBlock: true,
		}
		unminedBlock := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(3)})

		c.updatePendingUnminedBlocks(context.Background(), unminedBlock)

		timeout := time.NewTimer(2 * time.Second)
		select {
		case block := <-pendingUnminedBlockCh:
			if block.NumberU64() != unminedBlock.NumberU64() {
				t.Errorf("block numbers mismatch: have %v, want %v", block.NumberU64(), unminedBlock.NumberU64())
			}
		case <-timeout.C:
			t.Error("unexpected timeout occurs")
		}

		if s := len(c.pendingUnminedBlocks); s != 1 {
			t.Fatalf("Unmined blocks size must be 1, got %d", s)
		}
	})
}

func TestGetUnminedBlock(t *testing.T) {
	t.Run("block exists", func(t *testing.T) {
		expectedBlock := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		messages := newMessagesMap()
		curRoundMessages := messages.getOrCreate(2)
		unminedBlocks := make(map[uint64]*types.Block)
		unminedBlocks[expectedBlock.NumberU64()] = expectedBlock
		c := &core{
			round:                1,
			height:               big.NewInt(1),
			curRoundMessages:     curRoundMessages,
			messages:             messages,
			pendingUnminedBlocks: unminedBlocks,
		}

		block := c.getUnminedBlock()
		if !reflect.DeepEqual(block, expectedBlock) {
			t.Fatalf("Want %v, got %v", expectedBlock, block)
		}
	})

	t.Run("block does not exist", func(t *testing.T) {
		c := &core{
			round:                1,
			height:               big.NewInt(1),
			pendingUnminedBlocks: make(map[uint64]*types.Block),
		}

		block := c.getUnminedBlock()
		if block != nil {
			t.Fatalf("Want <nil>. got %v", block)
		}
	})
}

func TestCheckUnminedBlockMsg(t *testing.T) {
	t.Run("valid block is given, nil returned", func(t *testing.T) {
		c := &core{
			round:  1,
			height: big.NewInt(2),
		}

		block := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)})
		err := c.checkUnminedBlockMsg(block)
		if err != nil {
			t.Fatalf("want <nil>, got %v", err)
		}
	})

	t.Run("nil block is given, error returned", func(t *testing.T) {
		c := &core{}

		err := c.checkUnminedBlockMsg(nil)
		if err != errInvalidMessage {
			t.Fatalf("want %v, got %v", errInvalidMessage, err)
		}
	})

	t.Run("old block is given, error returned", func(t *testing.T) {
		c := &core{
			round:  1,
			height: big.NewInt(2),
		}

		oldBLock := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(1)})
		err := c.checkUnminedBlockMsg(oldBLock)
		if err != errOldHeightMessage {
			t.Fatalf("want %v, got %v", errOldHeightMessage, err)
		}
	})

	t.Run("future block is given, error returned", func(t *testing.T) {
		c := &core{
			round:  1,
			height: big.NewInt(1),
		}

		futureBlock := types.NewBlockWithHeader(&types.Header{Number: big.NewInt(2)})
		err := c.checkUnminedBlockMsg(futureBlock)
		if err != consensus.ErrFutureBlock {
			t.Fatalf("want %v, got %v", consensus.ErrFutureBlock, err)
		}
	})
}
