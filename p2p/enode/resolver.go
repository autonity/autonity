package enode

import (
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/clearmatics/autonity/log"
)

const (
	maxParseTries     = 3
	delayBetweenTries = 5 * time.Second
	resolveSetTTL     = 10 * time.Minute
)

var rs *resolveSet
var NotParsedErr = errors.New("not parsed")
func init() {
	rs = NewResolveSet()
}

func AutomaticResolveStart() {
	rs.Start(10 * time.Second)
}
func AutomaticResolveStop() {
	rs.Stop()
}

func SetResolveFunc(f func(host string) ([]net.IP, error)) {
	rs.resolveFunc = f
}

func NewResolveSet() *resolveSet {
	return &resolveSet{
		cache:             make(map[string]*Node),
		resolveSet:        make(map[string]resolveSetNode),
		started:           new(int32),
		resolveFunc:       net.LookupIP,
		maxTries:          maxParseTries,
		delayBetweenTries: delayBetweenTries,
	}
}

type resolveSet struct {
	sync.RWMutex
	cache             map[string]*Node
	resolveSet        map[string]resolveSetNode
	started           *int32
	resolveFunc       func(host string) ([]net.IP, error)
	maxTries          int
	delayBetweenTries time.Duration
}

func (rs *resolveSet) Start(resoveCycleSleepDuration time.Duration) {
	log.Warn("Async resolve started")
	swapped := atomic.CompareAndSwapInt32(rs.started, 0, 1)
	if !swapped {
		return
	}
	go func() {
		for {
			if atomic.LoadInt32(rs.started) == 0 {
				return
			}

			rs.Lock()
			currentTime := time.Now()

			for en, v := range rs.resolveSet {
				if v.resolved && currentTime.Sub(v.resolveTime) < resolveSetTTL {
					continue
				}

				node, err := rs.ParseV4WithResolveMaxTry(en, rs.maxTries, rs.delayBetweenTries)
				if err != nil {
					if err==NotParsedErr  {
						log.Warn("Node not resolved", "enode", en)
						continue
					}
					if _, ok:=errors.Unwrap(err).(*net.DNSError); ok {
						log.Warn("Node not resolved", "enode", en, "err", err)
						continue
					}

					delete(rs.resolveSet, en)
					continue
				}

				rs.cache[en] = node
				rs.resolveSet[en] = resolveSetNode{
					resolved:    true,
					resolveTime: currentTime,
				}
			}
			rs.Unlock()
			time.Sleep(resoveCycleSleepDuration)
		}
	}()
}

func (rs *resolveSet) Stop() {
	log.Warn("Async resolve stopped")

	atomic.StoreInt32(rs.started, 0)
}

func (rs *resolveSet) Add(enode string) {
	rs.Lock()
	defer rs.Unlock()
	rs.addNoLock(enode)
}

func (rs *resolveSet) addNoLock(enode string) {
	if _, ok := rs.resolveSet[enode]; !ok {
		rs.resolveSet[enode] = resolveSetNode{
			resolved: false,
		}
	}

}

func (rs *resolveSet) ParseV4WithResolveMaxTry(rawurl string, maxTry int, wait time.Duration) (*Node, error) {
	var node *Node
	var err error
	for i := 0; i < maxTry; i++ {
		node, err = rs.ParseV4WithResolve(rawurl)
		if err == nil {
			break
		}
		if _, ok:=err.(*net.DNSError); ok {
			log.Warn("trying to parse", "enode", rawurl, "attempt", i, "err", err)
		} else {
			return nil, err
		}

		time.Sleep(wait)
	}
	if node == nil {
		return nil, NotParsedErr
	}
	return node, err

}

type resolveSetNode struct {
	resolved    bool
	resolveTime time.Time
}

func (rs *resolveSet) Get(enodeStr string) (*Node, error) {
	var err error
	rs.RLock()
	node, ok := rs.cache[enodeStr]
	rs.RUnlock()

	if !ok {
		rs.Lock()
		if _, ok := rs.resolveSet[enodeStr]; !ok {
			rs.addNoLock(enodeStr)
		}
		rs.Unlock()
		node, err = rs.ParseV4WithResolveMaxTry(enodeStr, 1, rs.delayBetweenTries)
		if err != nil {
			return nil, err
		}

		rs.Lock()
		rs.cache[enodeStr] = node
		rs.resolveSet[enodeStr] = resolveSetNode{
			resolved:    true,
			resolveTime: time.Now(),
		}
		rs.Unlock()
	}

	return node, nil
}

func (rs *resolveSet) ParseV4WithResolve(rawurl string) (*Node, error) {
	return parseV4(rawurl, rs.resolveFunc)
}

func ParseV4WithResolveMaxTry(rawurl string, maxTry int, wait time.Duration) (*Node, error) {
	return rs.ParseV4WithResolveMaxTry(rawurl, maxTry, wait)
}

func ParseWithResolve(rawURL string) (*Node, error) {
	return rs.Get(rawURL)
}

func ParseV4WithResolve(rawurl string) (*Node, error) {
	return rs.ParseV4WithResolve(rawurl)
}