package router

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/log"
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
			latencyVec[i] = DefaultLatency
		}
	}

	index := indexOf(committee, r.txOpts.From)
	if index == -1 {
		return fmt.Errorf("validator %s not found in committee", r.txOpts.From.Hex())
	}
	if reported, err := r.protocolContracts.ClientReported(nil, big.NewInt(int64(index))); err != nil {
		log.Error("Reporter: failed to check if client reported", "err", err)
		return err
	} else if reported {
		log.Info("Reporter: client already reported latency")
		return nil
	}
	// reset values to estimate with test tx
	r.txOpts.NoSend = true
	r.txOpts.GasTipCap = nil
	r.txOpts.GasFeeCap = nil
	r.txOpts.GasLimit = 0
	testTx, err := r.protocolContracts.Latency.Report(
		r.txOpts,
		big.NewInt(int64(index)),
		latencyVec,
	)
	if err != nil {
		log.Error("Reporter: failed construct test report call", "err", err)
		return err
	}
	r.txOpts.NoSend = false
	r.txOpts.GasTipCap = new(big.Int).Mul(testTx.GasTipCap(), common.Big2)
	r.txOpts.GasFeeCap = new(big.Int).Mul(testTx.GasFeeCap(), common.Big2)
	r.txOpts.GasLimit = testTx.Gas() * 2

	if r.txOpts.GasFeeCap.Cmp(r.protocolContracts.Cache.MinimumBaseFee()) < 0 {
		r.txOpts.GasFeeCap.Set(new(big.Int).Mul(r.protocolContracts.Cache.MinimumBaseFee(), common.Big2))
	}

	tx, err := r.protocolContracts.Latency.Report(
		r.txOpts,
		big.NewInt(int64(index)),
		latencyVec,
	)

	if err == nil {
		log.Info(
			"Reporter: reported latency at tx",
			"tx", tx.Hash().Hex(),
			"from", r.txOpts.From.Hex(),
			"gasTipCap", r.txOpts.GasTipCap.String(),
			"gasFeeCap", r.txOpts.GasFeeCap.String(),
			"gasLimit", r.txOpts.GasLimit,
			"latencyLength", len(latencyVec),
			"latency", fmt.Sprintf("%v", latencyVec),
		)
	} else {
		log.Info("Reporter: failed to report latency", "err", err)
	}
	return err
}
