package types

func (h *Header) IsGenesis() bool {
	return h.Number.Uint64() == 0
}
