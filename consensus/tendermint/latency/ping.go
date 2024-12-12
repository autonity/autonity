package latency

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/log"
)

func pingIcmp(address string) <-chan *probing.Statistics {
	resultCh := make(chan *probing.Statistics)
	pinger, err := probing.NewPinger(address)
	if err != nil {
		panic(err)
	}
	pinger.Count = 5
	pinger.Timeout = 5 * time.Second
	pinger.OnRecv = func(pkt *probing.Packet) {
		log.Debug(
			"Ping response received",
			"bytes", pkt.Nbytes,
			"from", pkt.IPAddr,
			"icmp_seq", pkt.Seq,
			"rtt", pkt.Rtt,
		)
	}

	pinger.OnDuplicateRecv = func(pkt *probing.Packet) {
		log.Debug(
			"(DUP!) Ping response received",
			"bytes", pkt.Nbytes,
			"from", pkt.IPAddr,
			"icmp_seq", pkt.Seq,
			"rtt", pkt.Rtt,
			"ttl", pkt.TTL,
		)
	}

	pinger.OnFinish = func(stats *probing.Statistics) {
		log.Debug("Peer ping completed", "stats", stats)
		resultCh <- stats
	}

	go func() {
		if err = pinger.Run(); err != nil {
			panic(err)
		}
	}()

	return resultCh
}

func PingPeers(ips []string) []probing.Statistics {
	replyChannels := make([]<-chan *probing.Statistics, len(ips))
	results := make([]probing.Statistics, len(ips))
	for i, ip := range ips {
		if ip == "" {
			ch := make(chan *probing.Statistics, 1)
			ch <- &probing.Statistics{
				// this should be a reasonable default for max RTT
				AvgRtt: time.Second * 5,
			} // default result for non-connected peer to write
			replyChannels[i] = ch
			continue
		}
		replyChannels[i] = pingIcmp(ip)
	}

	for i, ch := range replyChannels {
		peerStats := <-ch
		results[i] = *peerStats
	}
	return results
}
