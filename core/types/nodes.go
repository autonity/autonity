package types

import (
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/enode"
)

type Nodes struct {
	List    []*enode.Node
	StrList []string
}

func NewNodes(strList []string) *Nodes {
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
			log.Debug("node retrieved", "node", enodeStr)
			newEnode, err := enode.ParseV4(enodeStr)
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
		var msg string
		for _, err := range errs {
			if err != nil {
				msg += fmt.Sprintf("%v\n", err)
			}
		}
		log.Error("enodes parse errors", "errs", msg)
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
