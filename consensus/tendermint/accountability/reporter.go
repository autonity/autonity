package accountability

import (
	"errors"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/autonity/bindings"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
)

const (
	MaxEventSize = 20480 // 20KB
	NumBackups   = 8     // todo: resolve a optimal number of backups for the AFD
)

var (
	errInvalidReport = errors.New("invalid report")
	errPendingReport = errors.New("pending report")
)

// primaryIndex returns the index of the validator which is assigned with a reporting block period.
func primaryIndex(height uint64, committeeSize uint64) int {
	return int((height / reportingSlotPeriod) % committeeSize) //nolint
}

// isRuleEngineRunner check if client is a rule engine runner, as to reduce the performance cost in a large scale
// network which contains a lots of consensus message to be scanned, we select a set of nodes as the rule runner
// of a specific height, for smale scale network, all the nodes run the rule engine.
func (fd *FaultDetector) isRuleEngineRunner(height uint64) bool {
	if _, ok := fd.reportersCache.Get(height); ok {
		return true
	}

	committee, err := fd.blockchain.CommitteeByHeight(height)
	if err != nil {
		fd.logger.Error("Failed to get committee for height %d: %v", height, err)
		return false
	}

	// All members run rule engine in a small scale network.
	if committee.Len() <= NumBackups*3 {
		fd.reportersCache.Add(height, struct{}{})
		return true
	}

	// Return true if node is the primary reporter.
	primary := primaryIndex(height, uint64(committee.Len())) //nolint
	if committee.Members[primary].Address == fd.address {
		fd.reportersCache.Add(height, struct{}{})
		return true
	}

	// Return true if node is backups.
	for i := 1; i <= NumBackups; i++ {
		backup := (primary + i) % committee.Len()
		if committee.Members[backup].Address == fd.address {
			fd.reportersCache.Add(height, struct{}{})
			return true
		}
	}

	return false
}

// canReport assign the validator a dedicated time-window to submit the accountability event, if the primary fails to
// report, those backups will report once they become to primary at next dedicated time-window.
func (fd *FaultDetector) canReport(height uint64) bool {
	committee, err := fd.blockchain.CommitteeByHeight(height)
	if err != nil {
		fd.logger.Crit("Can't retrieve committee for message", "err", err, "height", height)
	}

	// each validator is assigned a reporting slot
	primary := primaryIndex(height, uint64(committee.Len())) //nolint

	// if validator is the reporter of the slot period, and if checkpoint block is the end block of the
	// slot, then it is time to report the collected events by this validator.
	if height%reportingSlotPeriod != 0 {
		return false
	}
	return committee.Members[primary].Address == fd.address
}

func (fd *FaultDetector) reportEvents(events []*bindings.IAccountabilityEvent) []*bindings.IAccountabilityEvent {
	var filtered []*bindings.IAccountabilityEvent
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

func (fd *FaultDetector) tryReport(ev *bindings.IAccountabilityEvent) error {
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
		size := len(ev.RawProof)
		if size > MaxEventSize {
			fd.logger.Warn("Ignoring too large proof reporting", "size", size)
			continue
		}
		event := bindings.IAccountabilityEvent{
			EventType:      ev.EventType,
			Rule:           ev.Rule,
			Reporter:       ev.Reporter,
			Id:             common.Big0, // not required for submission
			Block:          common.Big0, // not required for submission
			Epoch:          common.Big0, // not required for submission
			ReportingBlock: common.Big0, // not required for submission
			MessageHash:    common.Big0, // not required for submission
			Offender:       ev.Offender,
			RawProof:       ev.RawProof,
		}
		var tx *types.Transaction
		var err error
		switch event.EventType {
		case uint8(autonity.Misbehaviour):
			tx, err = fd.protocolContracts.HandleMisbehaviour(fd.txOpts, event)
		case uint8(autonity.Accusation):
			tx, err = fd.protocolContracts.HandleAccusation(fd.txOpts, event)
		case uint8(autonity.Innocence):
			tx, err = fd.protocolContracts.HandleInnocenceProof(fd.txOpts, event)
		default:
			panic("Unknown accountability event")
		}

		if err == nil {
			fd.logger.Info("Accountability transaction sent", "tx", tx.Hash(), "gas", tx.Gas(), "size", tx.Size())
		} else {
			fd.logger.Error("Cannot submit accountability transaction", "err", err)
		}
	}
}
