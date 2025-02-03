package ping

import (
	"time"

	//TODO: library needs a version update(not officially released), but for POC this will suffice
	tcping "github.com/cloverstd/tcping/ping"
)

type tcpPinger struct{}

func NewTCPPinger() Pinger {
	return &tcpPinger{}
}

func (*tcpPinger) Ping(t Target, resultCh chan<- time.Duration) {
	target := &tcping.Target{Protocol: tcping.TCP, Host: t.IP, Port: t.Port, Counter: 5, Timeout: 3, Interval: 1}
	pinger := tcping.NewTCPing()
	pinger.SetTarget(target)
	go func() {
		done := pinger.Start()
		<-done
		resultCh <- pinger.Result().Avg()
	}()
}
