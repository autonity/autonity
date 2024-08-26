package core

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/log"
)

func pingIcmp(address string) <-chan *probing.Statistics {
	// todo: add context, change fmt to log
	resultCh := make(chan *probing.Statistics)
	pinger, err := probing.NewPinger(address)
	if err != nil {
		panic(err)
	}
	pinger.Count = 5
	pinger.OnRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v ttl=%v (DUP!)\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt, pkt.TTL)
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
		resultCh <- stats
	}

	fmt.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	go func() {
		if err = pinger.Run(); err != nil {
			panic(err)
		}
	}()

	return resultCh
}

func PingPeers(e *Engine) []probing.Statistics {
	replyChannels := make([]<-chan *probing.Statistics, len(e.Peers))
	results := make([]probing.Statistics, len(e.Peers))
	for i, peer := range e.Peers {
		if peer == nil || !peer.Connected {
			log.Debug("Peer not connected", "peer", i)
			ch := make(chan *probing.Statistics, 1)
			ch <- &probing.Statistics{} // default result for non-connected peer to write
			replyChannels[i] = ch
			continue
		}
		replyChannels[i] = pingIcmp(peer.Ip)
	}
	for i, ch := range replyChannels {
		peerStats := <-ch
		results[i] = *peerStats
	}
	return results
}

func BroadcastLatency(e *Engine, strategy uint64, latency []probing.Statistics) error {
	errs := make([]error, len(e.Peers))
	acks := make([]bool, len(e.Peers))
	var hasError atomic.Bool
	var wg sync.WaitGroup

	for i, peer := range e.Peers {
		if i == e.Id {
			errs[i] = nil
			continue
		}
		if peer == nil || !peer.Connected {
			errs[i] = errors.New("peer not connected or nil")
			hasError.Store(true)
			log.Error("Peer not connected", "peer", i)
			continue
		}

		wg.Add(1)
		go func(id int, peer *Peer) {
			_, _, err := peer.SendLatencyArray(strategy, FilterAveRtt(latency))
			if err != nil {
				hasError.Store(true)
				errs[id] = err
			} else {
				acks[id] = true
				errs[id] = nil
			}
			wg.Done()
		}(i, peer)
	}

	wg.Wait()

	if hasError.Load() {
		return errors.Join(errs...)
	}
	return nil
}

func BroadcastGraphReady(e *Engine, strategy uint64) error {
	errs := make([]error, len(e.Peers))
	acks := make([]bool, len(e.Peers))
	var hasError atomic.Bool
	var wg sync.WaitGroup

	for i, peer := range e.Peers {
		if i == e.Id {
			errs[i] = nil
			continue
		}
		if peer == nil || !peer.Connected {
			errs[i] = errors.New("peer not connected or nil")
			hasError.Store(true)
			log.Error("Peer not connected", "peer", i)
			continue
		}

		wg.Add(1)
		go func(id int, peer *Peer) {
			err := peer.sendGraphReady(strategy)
			if err != nil {
				log.Error("sendGraphReady err:", err)
				hasError.Store(true)
				errs[id] = err
			} else {
				acks[id] = true
				errs[id] = nil
			}
			wg.Done()
		}(i, peer)
	}

	wg.Wait()

	if hasError.Load() {
		return errors.Join(errs...)
	}
	return nil
}

func TriggerLatencyBroadcast(e *Engine, strategy uint64) error {
	errs := make([]<-chan error, len(e.Peers))
	for id, peer := range e.Peers {
		ch := make(chan error, 1)
		errs[id] = ch
		if id == e.Id {
			ch <- nil
			continue
		}
		if peer == nil || !peer.Connected {
			ch <- errPeerNotConnected
		} else {
			go func(peer *Peer, ch chan error) {
				err := peer.SendTriggerRequest(strategy)
				ch <- err
			}(peer, ch)
		}
	}
	for _, ch := range errs {
		err := <-ch
		if err != nil {
			return fmt.Errorf("error in send trigger request: %s", err.Error())
		}
	}
	return nil
}

func FilterAveRtt(latency []probing.Statistics) []time.Duration {
	rtts := make([]time.Duration, len(latency))
	for i, stat := range latency {
		rtts[i] = stat.AvgRtt
	}
	return rtts
}
