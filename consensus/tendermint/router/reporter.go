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

	testTx, err := r.protocolContracts.Latency.Report(&bind.TransactOpts{
		From:      common.Address{},
		Nonce:     r.txOpts.Nonce,
		Signer:    r.txOpts.Signer,
		Value:     r.txOpts.Value,
		GasPrice:  r.txOpts.GasPrice,
		GasFeeCap: r.txOpts.GasFeeCap,
		GasTipCap: r.txOpts.GasTipCap,
		GasLimit:  r.txOpts.GasLimit,
		Context:   r.txOpts.Context,
		NoSend:    true,
	}, big.NewInt(int64(index)), latencyVec)
	if err != nil {
		log.Error("Reporter: failed to create test transaction", "err", err)
		return err
	}
	tx, err := r.protocolContracts.Latency.Report(
		&bind.TransactOpts{
			From:      r.txOpts.From,
			Nonce:     r.txOpts.Nonce,
			Signer:    r.txOpts.Signer,
			Value:     r.txOpts.Value,
			GasPrice:  r.txOpts.GasPrice,
			GasFeeCap: new(big.Int).Mul(testTx.GasFeeCap(), big.NewInt(10)),
			GasTipCap: new(big.Int).Mul(testTx.GasTipCap(), big.NewInt(10)),
			GasLimit:  r.txOpts.GasLimit,
			Context:   r.txOpts.Context,
			NoSend:    false,
		},
		big.NewInt(int64(index)),
		latencyVec,
	)

	if err == nil {
		log.Info("Reporter: reported latency at testTx", "tx", tx.Hash().Hex())
	} else {
		log.Info("Reporter: failed to report latency", "err", err)
	}
	return err
}
