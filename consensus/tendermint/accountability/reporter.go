package accountability

import (
	"context"
	"errors"
	"time"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
)

const (
	ChunkProofSize        = 1024 // 1KB is around 800k gas
	MaxSubmissionAttempts = 100
	SubmissionDelay       = 1 * time.Second
	MaxChunks             = 10
)

var (
	errInvalidReport = errors.New("invalid report")
	errPendingReport = errors.New("pending report")
)

func (fd *FaultDetector) reportEvents(events []*autonity.AccountabilityEvent) []*autonity.AccountabilityEvent {
	var filtered []*autonity.AccountabilityEvent
	for i, ev := range events {
		err := fd.tryReport(ev)
		switch {
		case err == nil:
			return append(filtered, events[i+1:]...)
		case errors.Is(err, errInvalidReport):
			continue
		default:
			filtered = append(filtered, ev)
		}
	}
	return filtered
}

func (fd *FaultDetector) tryReport(ev *autonity.AccountabilityEvent) error {
	// youssef: some of this logic could belong to canReport
	if ev.EventType == uint8(autonity.Misbehaviour) {
		if res, err := fd.protocolContracts.CanSlash(nil, ev.Offender, ev.Rule, ev.Block); err != nil {
			// in which scenarios err can be returned ?
			fd.logger.Debug("Accountability canSlash", "error", err)
			return errInvalidReport
		} else if !res {
			fd.logger.Info("Reporting faulty validator cancelled, already slashed")
			return errInvalidReport
		}
	} else if ev.EventType == uint8(autonity.Accusation) {
		if ret, err := fd.protocolContracts.CanAccuse(nil, ev.Offender, ev.Rule, ev.Block); err != nil {
			// again, can this really happen?
			fd.logger.Debug("Accountability canAccuse", "error", err)
			return errInvalidReport
		} else if !ret.Result && ret.Deadline.Cmp(common.Big0) == 0 {
			fd.logger.Info("Reporting accusation cancelled: already slashed")
			return errInvalidReport
		} else if !ret.Result && ret.Deadline.Cmp(common.Big0) > 0 {
			// In this scenario, there is already a pending accusation.
			delay := ret.Deadline.Int64() - fd.blockchain.CurrentHeader().Number.Int64()
			if delay <= 0 {
				fd.logger.Info("Reporting accusation cancelled: in the past")
				// this should not be possible
				return errInvalidReport
			}
			fd.logger.Info("Reporting accusation delayed", "delay", delay)
			// this accusation submission will be re-attempted at the next slot
			return errPendingReport
		}
	}
	fd.logger.Warn("Reporting faulty validator", "offender", ev.Offender, "rule", autonity.Rule(ev.Rule).String(), "block", ev.Block)
	fd.eventReporterCh <- ev
	return nil
}

func (fd *FaultDetector) eventReporter() {
	defer fd.wg.Done()
	for ev := range fd.eventReporterCh {
		chunks := len(ev.RawProof)/ChunkProofSize + 1
		if chunks > MaxChunks {
			fd.logger.Warn("Ignoring too large proof reporting", "chunks", chunks)
			continue
		}
		for i := 0; i < chunks; i++ {
			chunkedEvent := autonity.AccountabilityEvent{
				Chunks:         uint8(chunks),
				ChunkId:        uint8(i),
				EventType:      ev.EventType,
				Rule:           ev.Rule,
				Reporter:       ev.Reporter,
				Id:             common.Big0, // not required for submission
				Block:          common.Big0, // not required for submission
				Epoch:          common.Big0, // not required for submission
				ReportingBlock: common.Big0, // not required for submission
				MessageHash:    common.Big0, // not required for submission
				Offender:       ev.Offender,
				RawProof:       ev.RawProof[i*ChunkProofSize : min((i+1)*ChunkProofSize, len(ev.RawProof))],
			}
			// select the correct transaction options based on the event type
			var txOpts *bind.TransactOpts
			if autonity.AccountabilityEventType(ev.EventType) == autonity.Innocence {
				txOpts = fd.innocenceTxOpts
			} else {
				txOpts = fd.txOpts
			}
			if tx, err := fd.protocolContracts.HandleEvent(txOpts, chunkedEvent); err == nil {
				fd.logger.Warn("Accountability transaction sent", "tx", tx.Hash(), "gas", tx.Gas(), "size", tx.Size())
				// wait until it get mined before moving to the next one
				attempt := 0
				for ; attempt < MaxSubmissionAttempts; attempt++ {
					select {
					case <-fd.stopRetry:
						return
					default:
						time.Sleep(SubmissionDelay)
						_, _, blockNumber, _, _ := fd.ethBackend.GetTransaction(context.Background(), tx.Hash())
						if blockNumber != 0 {
							break
						}
					}
				}
				if attempt == MaxSubmissionAttempts {
					fd.logger.Error("Accountability transaction didn't get mined, cancelling")
					break
				}
			} else {
				fd.logger.Error("Cannot submit accountability transaction", "err", err)
			}
		}
	}
}
