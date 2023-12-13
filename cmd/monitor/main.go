package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/consensus/tendermint/accountability"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/ethclient"
	"github.com/autonity/autonity/rlp"
	"github.com/autonity/autonity/rpc"
)

// THIS IS VERY QUICKLY HACKED TOGETHER AND BY NO MEAN SUITABLE FOR PROD
type session struct {
	wg             sync.WaitGroup
	eth            *ethclient.Client
	accountability *autonity.Accountability
	aut            *autonity.Autonity
	sync.RWMutex
}

func (s *session) run() {
	/*
		listener, err := net.Listen("tcp", "localhost:55000")
		if err != nil {
			panic("can't bind listening socket on port 55000")
		}
		defer listener.Close()
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		buf := make([]byte, 10)
		conn.Read(buf)
		conn.Close()
		host := "ws://localhost:" + strings.Trim(string(buf), "\x00")
	*/
	host := "wss://rpc1.piccadilly.autonity.org/ws"
	connected := false
	var rpcClient *rpc.Client
	for !connected {
		var err error
		if rpcClient, err = rpc.Dial(host); err == nil {
			connected = true
		} else {
			fmt.Println("can't connect to websocket")
			time.Sleep(time.Second)
		}
	}

	fmt.Println("Connected!")
	s.eth = ethclient.NewClient(rpcClient)

	s.accountability, _ = autonity.NewAccountability(autonity.AccountabilityContractAddress, s.eth)
	s.aut, _ = autonity.NewAutonity(autonity.AutonityContractAddress, s.eth)
	faultProofsCh := make(chan *autonity.AccountabilityNewFaultProof)
	epochCh := make(chan *autonity.AutonityNewEpoch)
	slashingCh := make(chan *autonity.AccountabilitySlashingEvent)
	accusationCh := make(chan *autonity.AccountabilityNewAccusation)
	innocenceCh := make(chan *autonity.AccountabilityInnocenceProven)
	newHeadCh := make(chan *types.Header)
	starBlock := uint64(0)
	opts := &bind.FilterOpts{Start: starBlock}
	s.eth.SubscribeNewHead(context.Background(), newHeadCh)

	s.accountability.WatchNewFaultProof(nil, faultProofsCh, nil)
	s.aut.WatchNewEpoch(nil, epochCh)
	s.accountability.WatchSlashingEvent(nil, slashingCh)
	s.accountability.WatchNewAccusation(nil, accusationCh, nil)
	s.accountability.WatchInnocenceProven(nil, innocenceCh, nil)
	s.listenFaultProofs(faultProofsCh)
	s.listenNewEpoch(epochCh)
	s.listenSlashingEvent(slashingCh)
	s.listenNewAccusation(accusationCh)
	s.listenInnocence(innocenceCh)
	//s.listenNewHead(newHeadCh)
	iter, _ := s.accountability.FilterNewFaultProof(opts, nil)
	for iter.Next() {
		faultProofsCh <- iter.Event
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := s.eth.BlockNumber(ctx)
		cancel()
		if err != nil {
			break
		}
		time.Sleep(time.Second)
	}
	fmt.Println("Connection lost")
	close(faultProofsCh)
	close(epochCh)
	close(slashingCh)
	close(innocenceCh)
	close(accusationCh)
	s.wg.Wait()
}

func (s *session) listenNewEpoch(ch chan *autonity.AutonityNewEpoch) {
	s.wg.Add(1)
	defer s.wg.Done()
	go func() {
		for {
			e, ok := <-ch
			if !ok {
				return
			}
			s.Lock()
			printEvent("New epoch", "id", e.Epoch.Uint64(), "block", e.Raw.BlockNumber)
			fmt.Println("> Committee:")
			committee, _ := s.aut.GetCommittee(nil)
			for _, c := range committee {
				fmt.Println(c.Addr, " - ", c.VotingPower.Uint64())
			}
			fmt.Println()
			s.Unlock()
		}
	}()
}

func (s *session) listenSlashingEvent(ch chan *autonity.AccountabilitySlashingEvent) {
	s.wg.Add(1)
	defer s.wg.Done()
	go func() {
		for {
			e, ok := <-ch
			if !ok {
				return
			}
			s.Lock()
			printEvent(`Slashing Event`, "validator", e.Validator, "amount", e.Amount)
			s.Unlock()
		}
	}()
}

func (s *session) listenFaultProofs(ch chan *autonity.AccountabilityNewFaultProof) {

	s.wg.Add(1)
	defer s.wg.Done()
	slashingCount := 0
	go func() {
		for {
			f, ok := <-ch
			if !ok {
				return
			}
			slashingCount++
			ev, _ := s.accountability.Events(nil, f.Id)
			val, _ := s.aut.GetValidator(nil, f.Offender)
			block, _ := s.eth.BlockByNumber(context.Background(), ev.Block)
			s.Lock()
			printEvent("Accountability Event", autonity.AccountabilityEventType(ev.EventType))
			fmt.Println("> Event data:")
			fmt.Println("  Offender              ", f.Offender)
			fmt.Println("  Reporting Block       ", f.Raw.BlockNumber)
			fmt.Println("  Event ID              ", f.Id)
			fmt.Println("  Severity              ", f.Severity)
			fmt.Println("> Proof data:")
			fmt.Println("  Rule                  ", autonity.Rule(ev.Rule))
			fmt.Println("  Offense Block         ", ev.Block)
			fmt.Println("  Epoch                 ", ev.Epoch)
			fmt.Println("  Time                  ", time.Unix(int64(block.Time()), 0))
			fmt.Println("  Event Type            ", autonity.AccountabilityEventType(ev.EventType))
			fmt.Println("  Reporter              ", ev.Reporter)
			switch autonity.Rule(ev.Rule) {
			case autonity.Equivocation:
				p := new(accountability.Proof)
				if err := rlp.DecodeBytes(ev.RawProof, &p); err != nil {
					fmt.Println("PROOF NOT DECODED!!       ", err)
					break
				}
				fmt.Println(" Faulty Vote 1", p.Message)
				for i := range p.Evidences {
					fmt.Println(" Faulty Vote "+strconv.Itoa(i+2), p.Evidences[i])
				}
			}
			fmt.Println("> Validator data (pre-slashing):")
			fmt.Println("  Total bonded stake    ", val.BondedStake)
			fmt.Println("  Self bonded stake     ", val.SelfBondedStake)
			fmt.Println("  Total Slashed         ", val.TotalSlashed)
			fmt.Println("  Provable faults count ", val.ProvableFaultCount)
			fmt.Println("")
			fmt.Println("Total Chain Slashings   ", slashingCount)
			s.Unlock()
		}
	}()
}

func (s *session) listenNewAccusation(ch chan *autonity.AccountabilityNewAccusation) {
	s.wg.Add(1)
	defer s.wg.Done()
	go func() {
		for ev := range ch {
			s.Lock()
			printEvent("Accusation Event", "id", ev.Id, "offender", ev.Offender, "severity", ev.Severity)
			s.Unlock()
		}
	}()
}

func (s *session) listenInnocence(ch chan *autonity.AccountabilityInnocenceProven) {
	s.wg.Add(1)
	defer s.wg.Done()
	go func() {
		for ev := range ch {
			s.Lock()
			printEvent("Innocence Event", "id", ev.Id, "offender", ev.Offender)
			s.Unlock()
		}
	}()
}

func (s *session) listenNewHead(ch chan *types.Header) {
	s.wg.Add(1)
	defer s.wg.Done()
	go func() {
		for ev := range ch {
			s.Lock()
			printEvent("New header", "num", ev.Number)
			s.Unlock()
		}
	}()
}

func printEvent(event string, other ...any) {
	var all []any
	now := "[" + time.Now().Format("15:04:05.000") + "]"
	all = append(all, now)
	all = append(all, event+" |")
	all = append(all, other...)
	fmt.Println(all...)
}

func main() {
	fmt.Println("Starting On-Chain Accountability Events Monitor....")
	for {
		new(session).run()
	}
}
