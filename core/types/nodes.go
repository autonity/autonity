package types

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/davecgh/go-spew/spew"
)

type Nodes struct {
	List    []*enode.Node
	StrList []string
}

const (
	maxParseTries     = 300
	delayBetweenTries = time.Second
	defaultTTL        = 120
)

func NewNodes(strList []string) *Nodes {
	getEnode := getParseFunc()

	idx := new(int32)
	wg := sync.WaitGroup{}
	errCh := make(chan error, len(strList))

	n := &Nodes{
		make([]*enode.Node, len(strList)),
		make([]string, len(strList)),
	}

	for _, enodeStr := range strList {
		wg.Add(1)

		go func(enodeStr string) {
			log.Debug("performing", "node", enodeStr)
			newEnode, err := cache.Get(enodeStr, getEnode)
			if err != nil {
				errCh <- err
			}

			currentIdx := atomic.AddInt32(idx, 1) - 1
			n.List[currentIdx] = newEnode
			n.StrList[currentIdx] = enodeStr

			wg.Done()
		}(enodeStr)
	}

	wg.Wait()
	close(errCh)

	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) != 0 {
		log.Error("enodes parse errors", "errs", spew.Sdump(errs))
	}

	return filterNodes(n)
}

func filterNodes(n *Nodes) *Nodes {
	filtered := &Nodes{
		make([]*enode.Node, 0, len(n.List)),
		make([]string, 0, len(n.StrList)),
	}

	for i, node := range n.List {
		if node != nil {
			filtered.List = append(filtered.List, node)
			filtered.StrList = append(filtered.StrList, n.StrList[i])
		}
	}

	return filtered
}

func getParseFunc() func(string) (*enode.Node, error) {
	getEnode := enode.ParseV4WithResolve
	getEnode = enode.GetParseV4WithResolveMaxTry(maxParseTries, delayBetweenTries)
	return getEnode
}

var cache = &domainCache{m: make(map[string]resolvedNode)}

type domainCache struct {
	m map[string]resolvedNode
	sync.RWMutex
}

type resolvedNode struct {
	node  *enode.Node
	count int
}

func (c *domainCache) Get(enodeStr string, getter func(string) (*enode.Node, error)) (*enode.Node, error) {
	c.RLock()
	node, ok := c.m[enodeStr]
	c.RUnlock()

	if !ok {
		// could be slow, so mutex is used
		n, err := getter(enodeStr)
		if err != nil {
			return nil, err
		}
		c.Lock()
		node = resolvedNode{node: n}
		c.m[enodeStr] = node
		c.Unlock()
	}

	c.Lock()
	node.count++
	if node.count >= defaultTTL {
		// reset the cache if TTL is reached
		delete(c.m, enodeStr)
	}
	c.Unlock()
	return node.node, nil
}
