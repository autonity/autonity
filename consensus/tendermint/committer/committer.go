package committer

import (
	"bytes"

	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/rlp"
)

type SignatureCache interface {
	Signatures(id []byte) [][]byte
}

type Committer struct {
	signatureCache SignatureCache
}

func (c *Committer) Commit(value, valueId []byte) error {
	// This is meant to be a block
	b := &types.Block{}
	err := rlp.Decode(bytes.NewBuffer(value), b)
	if err != nil {
		return err
	}
	// Look up the signatures for this block
	sigs := c.signatureCache.Signatures(valueId)

	// We need to have separate variable since Header() makes a copy of the
	// header.
	h := b.Header()
	// Append signatures and round into extra-data
	if err := types.WriteCommittedSeals(h, sigs); err != nil {
		return err
	}

	return nil
}
