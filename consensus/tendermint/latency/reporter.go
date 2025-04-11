package latency

import (
	"crypto/ecdsa"
	"github.com/autonity/autonity/log"
	"github.com/pkg/errors"
	"math/big"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
)

var errInvalidReporter = errors.New("Not a valid reporter")

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

func (r *Reporter) ReportLatency(committee []common.Address, latency map[common.Address]uint8) error {
	index := big.NewInt(-1)
	latencyVec := make([]uint8, len(committee))
	for i, validator := range committee {
		if validator == r.txOpts.From {
			index = big.NewInt(int64(i))
		}

		if val, ok := latency[validator]; ok {
			latencyVec[i] = val
		} else {
			latencyVec[i] = ^uint8(0)
		}
	}

	if index.Cmp(common.Big0) < 0 {
		return errInvalidReporter
	}

	tx, err := r.protocolContracts.Latency.Report(r.txOpts, index, latencyVec)
	if err == nil {
		log.Info("Latency reporting", "gas", tx.Gas())
	}

	return err
}
