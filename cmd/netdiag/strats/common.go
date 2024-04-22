package strats

import (
	"time"
)

type StratCode uint64

const (
	BroadcastCode = StratCode(iota)
	SimpleCode
	// LatRandCode
	// RandRandCode
	// LatLatCode
	StrategyCount
)

func (s StratCode) String() string {
	switch s {
	case BroadcastCode:
		return "Simple Broadcast"
	case SimpleCode:
		return "Simple Tree Dissemination"
	default:
		panic("unhandled default case")
	}
}

// State is a shared object amongst strategies.
type State struct {
	Id              int
	receivedPackets map[uint64]struct{}
	receivedReports map[uint64]chan *IndividualDisseminateResult
}

func NewState(id int) *State {
	return &State{
		Id:              id,
		receivedPackets: make(map[uint64]struct{}),
		receivedReports: make(map[uint64]chan *IndividualDisseminateResult),
	}
}

func (s *State) CollectReports(packetId uint64, maxPeers int) []IndividualDisseminateResult {
	individualResults := make([]IndividualDisseminateResult, maxPeers)
	timer := time.NewTimer(5 * time.Second)

LOOP:
	for i := 0; i < maxPeers; i++ { //we're not expecting ourselves to send it back
		select {
		case report := <-s.receivedReports[packetId]:
			individualResults[report.Sender] = *report
		case <-timer.C:
			break LOOP
		}
	}

	for i := range individualResults {
		if individualResults[i].ReceptionTime.IsZero() {
			individualResults[i] = IndividualDisseminateResult{
				Sender:        0,
				Relay:         0,
				Hop:           0,
				ReceptionTime: time.Time{},
				ErrorTimeout:  true,
			}
		}
	}

	return individualResults
}

type Strategy interface {
	// Execute contains the logic for the dissemination strategy from the source sender perspective.
	Execute(packetId uint64, data []byte, maxPeers int) error
	// HandlePacket has the logic when receiving a packet.
	HandlePacket(requestId uint64, hop uint8, originalSender uint64, data any) error
}

type Peer interface {
	DisseminateRequest(code uint64, requestId uint64, hop uint8, originalSender uint64, maxPeers uint64, data any) error
}

type PeerGetter func(int) Peer

type IndividualDisseminateResult struct {
	Sender        int
	Relay         int
	Hop           int
	ReceptionTime time.Time
	ErrorTimeout  bool
}

type BaseStrategy struct {
	Peers PeerGetter
	State *State
}
