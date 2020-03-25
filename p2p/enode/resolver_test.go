package enode

import (
	"errors"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

func TestAsyncResolver(t *testing.T) {
	rs := NewResolveSet()
	rs.maxTries = 2
	rs.delayBetweenTries = time.Millisecond
	failResolved := new(int32)
	rs.resolveFunc = func(host string) (ips []net.IP, e error) {
		switch host {
		case "domainsuccess.com":
			return []net.IP{
				net.ParseIP("127.0.0.1"),
			}, nil

		case "domainfail.com":
			if atomic.LoadInt32(failResolved) != 1 {
				return nil, &net.DNSError{}
			}
			return []net.IP{
				net.ParseIP("127.0.0.1"),
			}, nil
		default:
			return nil, errors.New("resolve err")
		}

	}

	successResolveEnode := "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@domainsuccess.com"

	failResolveEnode := "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@domainfail.com"

	v, err := rs.ParseV4WithResolve(failResolveEnode)
	if err == nil {
		t.Log(v.String())
		t.Fatal()
	}

	rs.Add(successResolveEnode)
	rs.Add(failResolveEnode)

	rs.Start(time.Second)

	//wait resolve run
	time.Sleep(2 * time.Second)

	//check async reload
	rs.RLock()
	v, okFail := rs.cache[failResolveEnode]
	_, okSuccess := rs.cache[successResolveEnode]
	rs.RUnlock()

	if !okSuccess {
		t.Fatal()
	}
	if okFail {
		t.Fatal()
	}

	atomic.StoreInt32(failResolved, 1)
	//wait resolve run
	time.Sleep(2 * time.Second)

	rs.RLock()
	_, okFail = rs.cache[failResolveEnode]
	_, okSuccess = rs.cache[successResolveEnode]
	rs.RUnlock()

	if !okSuccess {
		t.Fatal()
	}
	if !okFail {
		t.Fatal()
	}

}

func TestResolve(t *testing.T) {
	rs := NewResolveSet()
	rs.maxTries = 2
	rs.delayBetweenTries = time.Millisecond

	failResolved := new(int32)
	rs.resolveFunc = func(host string) (ips []net.IP, e error) {
		switch host {
		case "domainsuccess.com":
			return []net.IP{
				net.ParseIP("127.0.0.1"),
			}, nil

		case "domainfail.com":
			if atomic.LoadInt32(failResolved) != 1 {
				t.Log("err", host)
				return nil, errors.New("resolve err")
			}
			t.Log("succ", host)
			return []net.IP{
				net.ParseIP("127.0.0.1"),
			}, nil
		default:
			return nil, errors.New("resolve err")
		}

	}

	successResolveEnode := "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@domainsuccess.com"

	failResolveEnode := "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@domainfail.com"

	_, err := rs.Get(failResolveEnode)
	if err == nil {
		t.Fatal()
	}

	_, err = rs.Get(successResolveEnode)
	if err != nil {
		t.Fatal(err)
	}

	atomic.StoreInt32(failResolved, 1)
	_, err = rs.Get(failResolveEnode)
	if err != nil {
		t.Fatal(err)
	}

	_, err = rs.Get(successResolveEnode)
	if err != nil {
		t.Fatal(err)
	}
}