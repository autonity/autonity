package strats

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"
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

type ChunkInfo struct {
	Partial       bool
	TotalReceived int
	SeqReceived   []bool
}

// State is a shared object amongst strategies.
type State struct {
	Id    uint64
	Peers int
	// Those need to be protected
	ReceivedPackets map[uint64]ChunkInfo
	ReceivedReports map[uint64]chan *IndividualDisseminateResult
	AverageRTT      []time.Duration
	LatencyMatrix   [][]time.Duration
	PingResults     []probing.Statistics
	InfoChannel     chan any
}

func NewState(id uint64, peers int) *State {
	state := &State{
		Id:              id,
		Peers:           peers,
		ReceivedPackets: make(map[uint64]ChunkInfo),
		ReceivedReports: make(map[uint64]chan *IndividualDisseminateResult),
		AverageRTT:      make([]time.Duration, peers),
		LatencyMatrix:   make([][]time.Duration, peers),
	}
	state.LatencyMatrix[id] = make([]time.Duration, peers)
	return state
}

func (s *State) CollectReports(packetId uint64, maxPeers int) []IndividualDisseminateResult {
	recipients := maxPeers
	individualResults := make([]IndividualDisseminateResult, recipients)
	if s.Id < uint64(recipients) {
		// local node is dismissed during propagation
		//individualResults[s.Id]
		recipients--
	}

	timer := time.NewTimer(20 * time.Second)
	fullReportReceived := 0

LOOP:
	for fullReportReceived < recipients {
		select {
		case report := <-s.ReceivedReports[packetId]:
			if report.Full {
				fullReportReceived++
				individualResults[report.Sender] = *report
			} else if individualResults[report.Sender].ReceptionTime.IsZero() {
				individualResults[report.Sender] = *report
			}
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
	HandlePacket(packetId uint64, hop uint8, originalSender uint64, receivedFrom uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error
	// ConstructGraph constructs the graph to disseminate packet for `maxPeers`
	ConstructGraph(maxPeers int) error
	// GraphReadyForPeer is called when the graph for peer `peerID` is constructed
	GraphReadyForPeer(peerID int)
	// IsGraphReadyForPeer Returns `true` if graph is constructed for all peers
	IsGraphReadyForPeer(peerID int) bool
}

type Peer interface {
	DisseminateRequest(code uint64, requestId uint64, hop uint8, originalSender uint64, maxPeers uint64, data []byte, partial bool, seqNum, total uint16) error
	RTT() time.Duration
}

type PeerGetter func(int) Peer

type IndividualDisseminateResult struct {
	Sender        int
	Relay         int
	Hop           int
	ReceptionTime time.Time
	ErrorTimeout  bool
	Full          bool
}

type BaseStrategy struct {
	Peers PeerGetter
	Code  uint64
	State *State
}
