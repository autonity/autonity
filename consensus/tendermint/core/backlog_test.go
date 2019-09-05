package core

import (
	"math/big"
	"testing"

	"gopkg.in/karalabe/cookiejar.v2/collections/prque"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/consensus/tendermint/validator"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/log"
)

func TestCheckMessage(t *testing.T) {
	t.Run("valid params given, nil returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
		}

		err := c.checkMessage(big.NewInt(1), big.NewInt(2))
		if err != nil {
			t.Fatalf("have %v, want nil", err)
		}
	})

	t.Run("given nil round, error returned", func(t *testing.T) {
		c := &core{}

		err := c.checkMessage(nil, big.NewInt(2))
		if err != errInvalidMessage {
			t.Fatalf("have %v, want %v", err, errInvalidMessage)
		}
	})

	t.Run("given future height, error returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		err := c.checkMessage(big.NewInt(2), big.NewInt(4))
		if err != errFutureHeightMessage {
			t.Fatalf("have %v, want %v", err, errFutureHeightMessage)
		}
	})

	t.Run("given old height, error returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		err := c.checkMessage(big.NewInt(2), big.NewInt(2))
		if err != errOldHeightMessage {
			t.Fatalf("have %v, want %v", err, errOldHeightMessage)
		}
	})

	t.Run("given future round, error returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(3)),
		}

		err := c.checkMessage(big.NewInt(3), big.NewInt(3))
		if err != errFutureRoundMessage {
			t.Fatalf("have %v, want %v", err, errFutureRoundMessage)
		}
	})

	t.Run("given old round, error returned", func(t *testing.T) {
		c := &core{
			currentRoundState: NewRoundState(big.NewInt(2), big.NewInt(2)),
		}

		err := c.checkMessage(big.NewInt(1), big.NewInt(2))
		if err != errOldRoundMessage {
			t.Fatalf("have %v, want %v", err, errOldRoundMessage)
		}
	})
}

func TestStoreBacklog(t *testing.T) {
	t.Run("backlog from self", func(t *testing.T) {
		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			address:           common.HexToAddress("0x1234567890"),
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
		}

		c.storeBacklog(nil, validator.New(common.HexToAddress("0x1234567890")))
	})

	t.Run("vote message received", func(t *testing.T) {
		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			address:           common.HexToAddress("0x1234567890"),
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
			backlogs:          make(map[validator.Validator]*prque.Prque),
		}

		vote := &Vote{
			Round:  big.NewInt(1),
			Height: big.NewInt(2),
		}

		votePayload, _ := Encode(vote)

		msg := &message{
			Code: msgPrevote,
			Msg:  votePayload,
		}

		c.storeBacklog(msg, validator.New(common.HexToAddress("0x0987654321")))
	})

	t.Run("proposal message received", func(t *testing.T) {
		c := &core{
			logger:            log.New("backend", "test", "id", 0),
			address:           common.HexToAddress("0x1234567890"),
			backlogs:          make(map[validator.Validator]*prque.Prque),
			currentRoundState: NewRoundState(big.NewInt(1), big.NewInt(2)),
		}

		proposal := &Proposal{
			Round:         big.NewInt(1),
			Height:        big.NewInt(2),
			ValidRound:    big.NewInt(1),
			ProposalBlock: types.NewBlockWithHeader(&types.Header{}),
		}

		proposalPayload, _ := Encode(proposal)

		msg := &message{
			Code: msgProposal,
			Msg:  proposalPayload,
		}

		c.storeBacklog(msg, validator.New(common.HexToAddress("0x0987654321")))
	})
}
