package ping

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/autonity/autonity/log"
)

type Protocol int

const (
	ProtocolTCP Protocol = iota
)

type Target struct {
	IP   string
	Port int
}

type Result struct {
	Latency time.Duration
	Err     error
}

type Config struct {
	Timeout    time.Duration
	Interval   time.Duration
	Count      int
	MaxRetries int
}

func DefaultConfig() Config {
	return Config{
		Timeout:    3 * time.Second,
		Interval:   1 * time.Second,
		Count:      1,
		MaxRetries: 1,
	}
}

type Pinger interface {
	Ping(ctx context.Context, target Target) Result
}

type Dialer interface {
	DialContext(ctx context.Context, network, addr string) (net.Conn, error)
}

type Option func(*TCPPinger)

func WithConfig(cfg Config) Option {
	return func(p *TCPPinger) {
		p.config = cfg
	}
}

func WithLogger(logger log.Logger) Option {
	return func(p *TCPPinger) {
		p.logger = logger
	}
}

type TCPPinger struct {
	config Config
	logger log.Logger
	dialer Dialer
}

func NewPinger(protocol Protocol, logger log.Logger, opts ...Option) (Pinger, error) {
	if protocol != ProtocolTCP {
		return nil, fmt.Errorf("unsupported protocol: %d", protocol)
	}
	p := &TCPPinger{
		config: DefaultConfig(),
		logger: logger,
		dialer: &net.Dialer{},
	}
	for _, opt := range opts {
		opt(p)
	}
	if p.config.Timeout <= 0 {
		return nil, errors.New("timeout must be positive")
	}
	if p.config.Interval <= 0 {
		return nil, errors.New("interval must be positive")
	}
	if p.config.Count <= 0 {
		return nil, errors.New("count must be positive")
	}
	if p.config.MaxRetries < 0 {
		return nil, errors.New("max retries cannot be negative")
	}
	return p, nil
}

func (p *TCPPinger) Ping(ctx context.Context, target Target) Result {
	if target.IP == "" || target.Port <= 0 {
		return Result{Err: errors.New("invalid target: IP and port required")}
	}
	p.logger.Debug("starting TCP ping", "ip", target.IP, "port", target.Port, "count", p.config.Count)
	var lastErr error
	for attempt := 0; attempt <= p.config.MaxRetries; attempt++ {
		result, err := p.pingOnce(ctx, target)
		if err == nil {
			p.logger.Debug("TCP ping succeeded", "ip", target.IP, "port", target.Port, "latency_ms", result.Latency.Milliseconds())
			return result
		}
		lastErr = err
		p.logger.Warn("TCP ping attempt failed", "ip", target.IP, "port", target.Port, "attempt", attempt+1, "error", err)
		if attempt < p.config.MaxRetries {
			select {
			case <-ctx.Done():
				return Result{Err: ctx.Err()}
			case <-time.After(p.config.Interval):
			}
		}
	}
	return Result{Err: fmt.Errorf("ping failed after %d retries: %w", p.config.MaxRetries+1, lastErr)}
}

func (p *TCPPinger) pingOnce(ctx context.Context, target Target) (Result, error) {
	var totalLatency time.Duration
	successfulPings := 0
	for i := 0; i < p.config.Count; i++ {
		start := time.Now()
		conn, err := p.dialer.DialContext(ctx, "tcp", fmt.Sprintf("%s:%d", target.IP, target.Port))
		if err != nil {
			return Result{}, fmt.Errorf("dial failed: %w", err)
		}
		_ = conn.Close()
		latency := time.Since(start)
		totalLatency += latency
		successfulPings++
		if i < p.config.Count-1 {
			select {
			case <-ctx.Done():
				return Result{}, ctx.Err()
			case <-time.After(p.config.Interval):
			}
		}
	}
	if successfulPings == 0 {
		return Result{}, errors.New("no successful pings")
	}
	avgLatency := totalLatency / time.Duration(successfulPings)
	return Result{Latency: avgLatency}, nil
}
