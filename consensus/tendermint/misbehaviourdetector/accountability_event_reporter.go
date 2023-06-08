package misbehaviourdetector

import (
	"context"
	"errors"
	"fmt"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/common/hexutil"
	"github.com/autonity/autonity/core"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/internal/ethapi"
	"github.com/autonity/autonity/rpc"
	"math"
	"strings"
)

var (
	errOverSizeAccountabilityEventTX = errors.New("oversize accountability event transaction")
	errTooBigAccountabilityEvent     = errors.New("too big accountability event error")
)

const (
	AccountabilityEventHandlerMethod      = "handleAccountabilityEvents"
	DefaultChunkedAccountabilityEventSize = 32 * 1024 // 32KB the chunk size is resolved base on the testing with piccadilly's gas limit.
	MinChunkedAccountabilityEventSize     = 16 * 1024 // 16KB the chunk size is resolved base on the testing with piccadilly's gas limit.
)

func (fd *FaultDetector) sendAccountabilityTXs(accountabilityEvents []*autonity.AccountabilityEvent) {
	txs := fd.generateAccountabilityTXs(AccountabilityEventHandlerMethod, accountabilityEvents)
	for _, tx := range txs {
		e := fd.txPool.AddLocal(tx)
		if e != nil {
			fd.logger.Error("Could not add TX into TX pool", "err", e)
			continue
		}
		fd.logger.Info("Generate accountability transaction", "hash", tx.Hash())
	}
}

// generateAccountabilityTXs, it packs accountability event one by one in individual TX, for those TXs which exceed the
// block gas limit or the TX size in 512 KB, they would be chunked.
func (fd *FaultDetector) generateAccountabilityTXs(method string, accountabilityEvents []*autonity.AccountabilityEvent) (txs []*types.Transaction) {
	nonce := fd.txPool.Nonce(fd.address)

	for _, ev := range accountabilityEvents {
		tx, err := fd.genAccountabilityEventTX(nonce, method, ev)
		if err == nil {
			txs = append(txs, tx)
			nonce++
			continue
		}

		if err == errOverSizeAccountabilityEventTX {
			chunkedTXs, err := fd.genChunkedAccountabilityEventTXs(nonce, method, ev)
			if err != nil {
				fd.logger.Error("genChunkedAccountabilityEventTXs", "error", err)
				continue
			}
			txs = append(txs, chunkedTXs...)
			nonce += uint64(len(chunkedTXs))
			continue
		}

		fd.logger.Error("generateAccountabilityTXs", "error", err)
	}

	return txs
}

// resolve a proper chunk size base on the len of proof and the gas limit of current chain.
func (fd *FaultDetector) resolveAccountabilityEventChunkSize(nonce uint64, method string, ev *autonity.AccountabilityEvent) (int, error) {

	// base on the testing with piccadilly gas limit and the gas cost of chunked event handling, we shrink from the
	// default 32KB chunk size with 1 KB each time, until we got 16KB, base on the testing with piccadilly gas limit.
	for resolvedChunkSize := DefaultChunkedAccountabilityEventSize; resolvedChunkSize >= MinChunkedAccountabilityEventSize; resolvedChunkSize -= 1024 {
		if len(ev.RawProof) > resolvedChunkSize {
			var chunkedEvent = &autonity.AccountabilityEvent{
				Chunks:   uint8(2),
				ChunkID:  0,
				Type:     ev.Type,
				Rule:     ev.Rule,
				Reporter: ev.Reporter,
				Sender:   ev.Sender,
				MsgHash:  ev.MsgHash,
				RawProof: ev.RawProof[0:resolvedChunkSize],
			}
			// check if everything is good by using resolved chunked size by pre-execute the tx.
			_, err := fd.genAccountabilityEventTX(nonce, method, chunkedEvent)
			if err == nil {
				return resolvedChunkSize, nil
			}
		}
	}
	return 0, fmt.Errorf("cannot resolve afd event chunk size")
}

// genChunkedAccountabilityEventTXs, as the event is oversize due to block gas limit or the tx size limit 512KB, this
// function chunks the event into multiple TXs.
func (fd *FaultDetector) genChunkedAccountabilityEventTXs(nonce uint64, method string, ev *autonity.AccountabilityEvent) ([]*types.Transaction, error) {

	if len(ev.RawProof) == 0 {
		return nil, fmt.Errorf("no proof of accountability event")
	}

	// resolve chunk size and chunks
	resolvedChunkSize, err := fd.resolveAccountabilityEventChunkSize(nonce, method, ev)
	if err != nil {
		return nil, err
	}

	chunks := 1
	if len(ev.RawProof)%resolvedChunkSize == 0 {
		chunks = len(ev.RawProof) / resolvedChunkSize
	} else {
		if len(ev.RawProof) > resolvedChunkSize {
			chunks += len(ev.RawProof) / resolvedChunkSize
		}
	}

	if chunks == 1 {
		// this shouldn't happen base on the data collected from the test with piccadilly's gas limit.
		return nil, fmt.Errorf("an AFD event with less 16KB is treated as oversize event")
	}

	if chunks > math.MaxUint8 {
		return nil, errTooBigAccountabilityEvent
	}

	fd.logger.Info("start to chunk oversize accountability ev", "client: ", fd.address, "nonce", nonce, "ev type", ev.Type, "ev rule", ev.Rule, "ev msg hash ", ev.MsgHash, "malicious node: ", ev.Sender, "chunk size", resolvedChunkSize, "chunks", chunks)
	var txs []*types.Transaction
	start := 0
	end := 0
	for chunkID := uint8(0); chunkID < uint8(chunks); chunkID++ {
		var chunkedEvent = &autonity.AccountabilityEvent{
			Chunks:   uint8(chunks),
			ChunkID:  chunkID,
			Type:     ev.Type,
			Rule:     ev.Rule,
			Reporter: ev.Reporter,
			Sender:   ev.Sender,
			MsgHash:  ev.MsgHash,
		}

		start = int(chunkID) * resolvedChunkSize
		if chunkID+1 < uint8(chunks) {
			end = start + resolvedChunkSize
			chunkedEvent.RawProof = ev.RawProof[start : start+resolvedChunkSize]
		} else {
			end = len(ev.RawProof)
			chunkedEvent.RawProof = ev.RawProof[start:len(ev.RawProof)]
		}

		fd.logger.Info("gen chunked TX", "chunkID: ", chunkID, "byte from: ", start, "byte end: ", end, "nonce: ", nonce)
		tx, err := fd.genAccountabilityEventTX(nonce, method, chunkedEvent)
		if err != nil {
			return nil, err
		}
		nonce++
		txs = append(txs, tx)
	}

	fd.logger.Info("finish to chunk oversize accountability ev")
	return txs, nil
}

// genAccountabilityEventTX, it packs a single accountability event into a TX, and if the TX exceed the block gas
// limit or the TX size in 512KB then it returns an oversize error.
func (fd *FaultDetector) genAccountabilityEventTX(nonce uint64, method string, ev *autonity.AccountabilityEvent) (*types.Transaction, error) {
	to := autonity.ContractAddress
	abi := fd.blockchain.GetAutonityContract().ABI()

	var proofs = make([]autonity.AccountabilityEvent, 1)
	proofs[0] = *ev

	packedData, err := abi.Pack(method, proofs)
	if err != nil {
		fd.logger.Error("Cannot pack accountability transaction", "err", err)
		return nil, err
	}

	price, err := fd.ethBackend.SuggestGasTipCap(context.Background())
	if err != nil {
		fd.logger.Error("Cannot pack accountability transaction", "err", err)
		return nil, err
	}
	head := fd.blockchain.CurrentHeader()
	price.Add(price, head.BaseFee)

	callArgs := ethapi.TransactionArgs{
		From:     &fd.address,
		To:       &to,
		GasPrice: (*hexutil.Big)(price),
		Value:    new(hexutil.Big),
		Data:     (*hexutil.Bytes)(&packedData),
	}

	blockHash := fd.blockchain.CurrentBlock().Hash()
	blockNumOrHash := rpc.BlockNumberOrHash{
		BlockHash: &blockHash,
	}

	// estimate gas cost base on current chain head.
	estimatedGas, err := ethapi.DoEstimateGas(context.Background(), fd.ethBackend, callArgs, blockNumOrHash, fd.ethBackend.RPCGasCap())
	if err != nil {
		if strings.Contains(err.Error(), "gas required exceeds allowance") {
			return nil, errOverSizeAccountabilityEventTX
		}
		return nil, err
	}

	fd.logger.Info("estimated gas base on latest chain state for accountability event", "gas", uint64(estimatedGas))
	// if a single event TX exceed the current block's gas limit, return oversize error, the event will be chunked.
	if uint64(estimatedGas) > head.GasLimit {
		return nil, errOverSizeAccountabilityEventTX
	}

	// reserve more gas than the value estimated since the estimation is not always accurate.
	reservedGas := uint64(estimatedGas + 999999)
	tx, err := types.SignTx(types.NewTransaction(nonce, to, common.Big0, reservedGas, price,
		packedData), types.HomesteadSigner{}, fd.nodeKey)
	if err != nil {
		return nil, err
	}

	// if tx exceed 512KB, it should be chunked too. In some settings with big gas limit for block, it happens.
	if uint64(tx.Size()) > core.TxMaxSize {
		return nil, errOverSizeAccountabilityEventTX
	}
	return tx, nil
}

func (fd *FaultDetector) faultDetectorTXEventLoop() {
	go func() {
		for {
			select {
			case accountabilityEvents := <-fd.accountabilityTXCh:
				fd.sendAccountabilityTXs(accountabilityEvents)
			case err, ok := <-fd.accountabilityTXEventSub.Err():
				if ok {
					panic(fmt.Sprintf("accountabilityEventSub error: %v", err.Error()))
				}
				return
			}
		}
	}()
}
