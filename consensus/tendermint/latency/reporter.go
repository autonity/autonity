package latency

import (
	"crypto/ecdsa"
	"errors"
	"math/big"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
)

type Reporter struct {
	txOpts            *bind.TransactOpts
	protocolContracts *autonity.ProtocolContracts
}

func NewReporter(
	chainID *big.Int,
	nodeKey *ecdsa.PrivateKey,
	contracts *autonity.ProtocolContracts,
) (*Reporter, error) {
	txOpts, err := bind.NewKeyedTransactorWithChainID(nodeKey, chainID)
	if err != nil {
		return nil, err
	}
	return &Reporter{txOpts: txOpts, protocolContracts: contracts}, nil
}

func (r *Reporter) ReportLatency(latency map[common.Address]uint8) error {
	committee, err := r.protocolContracts.Latency.GetCommittee(nil)
	if err != nil {
		return err
	}
	latencyVec := make([]uint8, len(committee))
	for i, validator := range committee {
		if val, ok := latency[validator]; ok {
			latencyVec[i] = val
		} else {
			latencyVec[i] = ^uint8(0)
		}
	}

	// check if the optimization of the clustering view is already done.
	activationHeight, err := r.protocolContracts.Latency.GetNewViewHeight(nil)
	if err != nil {
		return err
	}

	if activationHeight.Cmp(common.Big0) > 0 {
		return errors.New("clustering optimization was already done")
	}

	_, err = r.protocolContracts.Latency.Report(r.txOpts, latencyVec)
	return err
}
