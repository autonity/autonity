package ping

import "time"

type Protocol int

const (
	TCP Protocol = iota
	ICMP
)

type Pinger interface {
	Ping(Target, chan<- time.Duration)
}

type Target struct {
	IP   string
	Port int
}

func NewPinger(typ Protocol) Pinger {
	switch typ {
	case TCP:
		return NewTCPPinger()
	case ICMP:
		return NewICMPPinger()
	}
	return nil
}
