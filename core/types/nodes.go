package types

import (
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
	"sync"
	"time"
)

type Nodes struct {
	List    []*enode.Node
	StrList []string
}

func NewNodes(strList []string, openNetwork bool) *Nodes {
	n := &Nodes{
		[]*enode.Node{},
		[]string{},
	}

	getEnode := enode.ParseV4WithResolve
	if !openNetwork {
		getEnode = enode.GetParseV4WithResolveMaxTry(1000, time.Second)
	}

	for _, enodeStr := range strList {
		newEnode, err := cache.Get(enodeStr, getEnode)
		if err != nil {
			log.Error("Invalid whitelisted enode", "returned enode", enodeStr, "error", err.Error())

			if !openNetwork {
				panic(err)
			}
		}
		n.List = append(n.List, newEnode)
		n.StrList = append(n.StrList, enodeStr)
	}

	return n
}

var cache = &domainCache{m:make(map[string]resolvedNode)}

const defaultTTL = 20

type domainCache struct {
	m map[string]resolvedNode
	sync.RWMutex
}

type resolvedNode struct {
	node *enode.Node
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