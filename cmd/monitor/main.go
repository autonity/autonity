package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/ethclient"
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
	host := "wss://rpc1.piccadilly.autonity.org:8546"
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
	opts := &bind.WatchOpts{Start: &starBlock}
	s.eth.SubscribeNewHead(nil, newHeadCh)
	s.accountability.WatchNewFaultProof(opts, faultProofsCh, nil)
	s.aut.WatchNewEpoch(nil, epochCh)
	s.accountability.WatchSlashingEvent(opts, slashingCh)
	s.accountability.WatchNewAccusation(opts, accusationCh, nil)
	s.accountability.WatchInnocenceProven(opts, innocenceCh, nil)
	s.listenFaultProofs(faultProofsCh)
	s.listenNewEpoch(epochCh)
	s.listenSlashingEvent(slashingCh)
	s.listenNewAccusation(accusationCh)
	s.listenInnocence(innocenceCh)
	s.listenNewHead(newHeadCh)
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
	go func() {
		for {
			f, ok := <-ch
			if !ok {
				return
			}
			ev, _ := s.accountability.Events(nil, f.Id)
			val, _ := s.aut.GetValidator(nil, f.Offender)
			s.Lock()
			printEvent("Accountability Event", "offender:", f.Offender, "severity:", f.Severity, "id:", f.Id, "block:", f.Raw.BlockNumber)
			fmt.Println("> Proof data:")
			fmt.Println("  Rule                  ", autonity.Rule(ev.Rule))
			fmt.Println("  Offense Block         ", ev.Block)
			fmt.Println("  Epoch                 ", ev.Epoch)
			fmt.Println("  Event Type            ", autonity.AccountabilityEventType(ev.EventType))
			fmt.Println("  Reporter              ", ev.Reporter)
			fmt.Println("> Validator data (pre-slashing):")
			fmt.Println("  Total bonded stake    ", val.BondedStake)
			fmt.Println("  Self bonded stake     ", val.SelfBondedStake)
			fmt.Println("  Total Slashed         ", val.TotalSlashed)
			fmt.Println("  Provable faults count ", val.ProvableFaultCount)
			fmt.Println("")
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
