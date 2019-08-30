package core

import (
	"fmt"
	"math/big"
)

type walMessage struct {
	m       *message
	payload []byte
	height  *big.Int
	round   *big.Int
}

func (m *walMessage) Key() []byte {
	return []byte(fmt.Sprintf("message-%s-%s-%d-%s",
		m.height.String(),
		m.round.String(),
		m.m.Code,
		m.m.Address,
	))
}

func (m *walMessage) Value() []byte {
	return m.payload
}
