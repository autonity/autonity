package autonity

import "github.com/clearmatics/autonity/eth"

type AutonityAPI struct {
	eth *eth.Ethereum
}

func NewAutonityAPI(eth *eth.Ethereum) *AutonityAPI {
	return &AutonityAPI{eth}
}
