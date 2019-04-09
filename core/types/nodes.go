package types

import (
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/p2p/enode"
)

type Nodes struct {
	List    []*enode.Node
	StrList []string
}

func NewNodes(strList []string, must bool) *Nodes {
	n := &Nodes{
		[]*enode.Node{},
		[]string{},
	}

	for _, enodeStr := range strList {
		newEnode, err := enode.ParseV4WithResolve(enodeStr)
		if must && err != nil {
			log.Error("Invalid whitelisted enode", "returned enode", enodeStr, "error", err.Error())
			panic(err)
		}

		n.List = append(n.List, newEnode)
		n.StrList = append(n.StrList, enodeStr)
	}

	return n
}
