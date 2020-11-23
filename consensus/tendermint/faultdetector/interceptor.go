package faultdetector

import "github.com/clearmatics/autonity/common"

type interceptor struct {
	verifier Verifier
	msgStore Store
}

type Verifier interface {
	Verify(m *message, s Store) error
}

type message interface {
	Round() uint
	Height() uint
	Sender() common.Address
	Type() byte
	Value() common.Hash // Block hash for a proposal,
}

func (i *interceptor) Intercept(msg *message) {

	i.verifier.Verify(msg, i.msgStore)

}
