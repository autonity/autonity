package strats

import (
	"time"
)

var StrategyRegistry []registeredStrategy

type registeredStrategy struct {
	Name        string
	Code        uint64
	Constructor func(peerGetter PeerGetter, state *State) Strategy
}

func registerStrategy(name string, newFn func(base BaseStrategy) Strategy) {
	id := uint64(len(StrategyRegistry))
	newRegistration := registeredStrategy{
		Name: name,
		Code: id,
		Constructor: func(peerGetter PeerGetter, state *State) Strategy {
			base := BaseStrategy{
				Peers: peerGetter,
				Code:  id,
				State: state,
			}
			return newFn(base)
		},
	}
	StrategyRegistry = append(StrategyRegistry, newRegistration)
}

// State is a shared object amongst strategies.
type State struct {
	Id uint64
	// Those need to be protected
	ReceivedPackets map[uint64]struct{}
	ReceivedReports map[uint64]chan *IndividualDisseminateResult
}

func NewState(id uint64) *State {
	return &State{
		Id:              id,
		ReceivedPackets: make(map[uint64]struct{}),
		ReceivedReports: make(map[uint64]chan *IndividualDisseminateResult),
	}
}

func (s *State) CollectReports(packetId uint64, maxPeers int) []IndividualDisseminateResult {
	individualResults := make([]IndividualDisseminateResult, maxPeers)
	timer := time.NewTimer(5 * time.Second)

LOOP:
	for i := 0; i < maxPeers; i++ { //we're not expecting ourselves to send it back
		select {
		case report := <-s.ReceivedReports[packetId]:
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
	// The strategy should disseminate packets with recipient from id = 0 to id = maxPeers-1
	Execute(packetId uint64, data []byte, maxPeers int) error
	// HandlePacket has the logic when receiving a packet.
	HandlePacket(packetId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte) error
}

type Peer interface {
	DisseminateRequest(code uint64, requestId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte) error
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
	Code  uint64
	State *State
}
