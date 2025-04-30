package ping

import (
	"time"

	probing "github.com/prometheus-community/pro-bing"

	"github.com/autonity/autonity/log"
)

type icmpPinger struct{}

func NewICMPPinger() Pinger {
	return &icmpPinger{}
}

func (*icmpPinger) Ping(t Target, resultCh chan<- time.Duration) {
	pinger, err := probing.NewPinger(t.IP)
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
		resultCh <- stats.AvgRtt
	}

	go func() {
		if err = pinger.Run(); err != nil {
			panic(err)
		}
	}()
}
