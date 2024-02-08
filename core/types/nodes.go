package types

import (
	"fmt"
	"sync"

	"github.com/autonity/autonity/log"
	"github.com/autonity/autonity/p2p/enode"
)

type Nodes struct {
	List    []*enode.Node
	StrList []string
}

func NewNodes(strList []string, asACN bool) *Nodes {
	wg := sync.WaitGroup{}
	errCh := make(chan error, len(strList))
	var parser func(string) (*enode.Node, error)
	if asACN {
		parser = enode.ParseACNV4
	} else {
		parser = enode.ParseV4
	}

	n := &Nodes{
		make([]*enode.Node, len(strList)),
		make([]string, len(strList)),
	}

	for i, enodeStr := range strList {
		idx := i
		wg.Add(1)

		go func(enodeStr string, idx int) {
			log.Debug("node retrieved", "node", enodeStr)

			newEnode, err := parser(enodeStr)
			if err != nil {
				errCh <- err
			}

			n.List[idx] = newEnode
			n.StrList[idx] = enodeStr

			wg.Done()
		}(enodeStr, idx)
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
