package strats

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

const (
	BroadcastCode = iota
	SimpleCode
	LatRandCode
	RandRandCode
	LatLatCode
	//....
)

type Strategy interface {
	Execute(data []byte) error
}

type ResultBase struct {
	Size              int
	StartTime         time.Time
	IndividualResults []IndividualDisseminateResult
}

type IndividualDisseminateResult struct {
	Sender        int
	Relay         int
	Hop           int
	ReceptionTime time.Time
	ErrorTimeout  bool
}

func (r *ResultBase) String() string {
	var builder strings.Builder
	fmt.Fprintf(&builder, "Disseminate Results \n")
	var results []IndividualDisseminateResult
	for i, res := range r.IndividualResults {
		if res.ErrorTimeout {
			continue
		}
		results = append(results, r.IndividualResults[i])
		fmt.Fprintf(&builder, "Peer %d Duration: %s Hops: %d Relay: %d\n", i, res.ReceptionTime.Sub(r.StartTime), res.Hop, res.Relay)
	}
	sort.Slice(results, func(a, b int) bool {
		return results[a].ReceptionTime.Before(results[b].ReceptionTime)
	})
	n := len(results)
	fmt.Fprintf(&builder, "min: %s, median:%s 2/3rd:%s max: %s\n", results[0].ReceptionTime.Sub(r.StartTime), results[n/2].ReceptionTime.Sub(r.StartTime), results[(2*n)/3].ReceptionTime.Sub(r.StartTime), results[n-1].ReceptionTime.Sub(r.StartTime))

	return builder.String()
}

func (r *ResultBase) CollectReports() {
	individualResults := make([]*IndividualDisseminateResult, len(p.engine.peers))
	timer := time.NewTimer(5 * time.Second)

LOOP:
	for i := 0; i < len(p.engine.peers)-1; i++ { //we're not expecting ourselves to send it back
		select {
		case report := <-p.engine.receivedReports[packetId]:
			individualResults[report.Sender] = report
		case <-timer.C:
			break LOOP
		}
	}

	for i := range individualResults {
		if individualResults[i] == nil {
			individualResults[i] = &IndividualDisseminateResult{
				Sender:        0,
				Relay:         0,
				Hop:           0,
				ReceptionTime: time.Time{},
				ErrorTimeout:  true,
			}
		}
	}
	// Wait for reports
	reply.IndividualResults = individualResults
	reply.Size = args.Size
	return nil
}
