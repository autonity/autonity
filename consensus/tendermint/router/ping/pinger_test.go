package ping

import (
	"context"
	"errors"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/autonity/autonity/log"
)

type mockDialer struct {
	dialFunc func(ctx context.Context, network, addr string) (net.Conn, error)
}

func (m *mockDialer) DialContext(ctx context.Context, network, addr string) (net.Conn, error) {
	return m.dialFunc(ctx, network, addr)
}

type mockConn struct{}

func (m mockConn) Read(b []byte) (n int, err error) {
	return 0, nil
}

func (m mockConn) Write(b []byte) (n int, err error) {
	return 0, nil
}

func (m mockConn) Close() error {
	return nil
}

func (m mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{}
}

func (m mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{}
}

func (m mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestTCPPingerPing(t *testing.T) {
	logger := log.New()
	target := Target{IP: "127.0.0.1", Port: 8080}

	t.Run("SuccessfulPing", func(t *testing.T) {
		pinger := &TCPPinger{
			config: DefaultConfig(),
			logger: logger,
			dialer: &mockDialer{
				dialFunc: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return &mockConn{}, nil
				},
			},
		}
		result := pinger.Ping(context.Background(), target)
		assert.NoError(t, result.Err)
		assert.GreaterOrEqual(t, result.Latency, time.Duration(0))
	})

	t.Run("InvalidTargetEmptyIP", func(t *testing.T) {
		pinger := &TCPPinger{
			config: DefaultConfig(),
			logger: logger,
			dialer: &mockDialer{},
		}
		result := pinger.Ping(context.Background(), Target{Port: 8080})
		assert.Error(t, result.Err)
		assert.Equal(t, "invalid target: IP and port required", result.Err.Error())
	})

	t.Run("InvalidTargetZeroPort", func(t *testing.T) {
		pinger := &TCPPinger{
			config: DefaultConfig(),
			logger: logger,
			dialer: &mockDialer{},
		}
		result := pinger.Ping(context.Background(), Target{IP: "127.0.0.1"})
		assert.Error(t, result.Err)
		assert.Equal(t, "invalid target: IP and port required", result.Err.Error())
	})

	t.Run("PingFailure", func(t *testing.T) {
		pinger := &TCPPinger{
			config: DefaultConfig(),
			logger: logger,
			dialer: &mockDialer{
				dialFunc: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return nil, errors.New("connection refused")
				},
			},
		}
		result := pinger.Ping(context.Background(), target)
		assert.Error(t, result.Err)
		assert.Contains(t, result.Err.Error(), "ping failed after 2 retries: dial failed: connection refused")
	})

	t.Run("RetrySuccessAfterFailure", func(t *testing.T) {
		attempts := 0
		pinger := &TCPPinger{
			config: Config{
				Timeout:    1 * time.Second,
				Interval:   10 * time.Millisecond,
				Count:      1,
				MaxRetries: 2,
			},
			logger: logger,
			dialer: &mockDialer{
				dialFunc: func(ctx context.Context, network, addr string) (net.Conn, error) {
					attempts++
					if attempts < 2 {
						return nil, errors.New("temporary failure")
					}
					return &mockConn{}, nil
				},
			},
		}
		result := pinger.Ping(context.Background(), target)
		assert.NoError(t, result.Err)
		assert.Equal(t, 2, attempts)
	})

	t.Run("ContextCancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		pinger := &TCPPinger{
			config: Config{
				Timeout:    1 * time.Second,
				Interval:   1 * time.Second,
				Count:      1,
				MaxRetries: 2,
			},
			logger: logger,
			dialer: &mockDialer{
				dialFunc: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return nil, errors.New("connection refused")
				},
			},
		}
		go func() {
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()
		result := pinger.Ping(ctx, target)
		assert.Error(t, result.Err)
		assert.Equal(t, context.Canceled, result.Err)
	})

	t.Run("ConcurrentPings", func(t *testing.T) {
		pinger := &TCPPinger{
			config: DefaultConfig(),
			logger: logger,
			dialer: &mockDialer{
				dialFunc: func(ctx context.Context, network, addr string) (net.Conn, error) {
					time.Sleep(10 * time.Millisecond) // Simulate network delay
					return &mockConn{}, nil
				},
			},
		}
		var wg sync.WaitGroup
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				result := pinger.Ping(context.Background(), target)
				assert.NoError(t, result.Err)
			}()
		}
		wg.Wait()
	})
}
