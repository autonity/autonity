package byzantine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/consensus/tendermint/core"
	"github.com/autonity/autonity/consensus/tendermint/core/interfaces"
	"github.com/autonity/autonity/consensus/tendermint/core/message"
	e2e "github.com/autonity/autonity/e2e_test"
)

func newMalPrecommitService(c interfaces.Core) interfaces.Precommiter {
	return &malPrecommitService{c.(*core.Core), c.Precommiter()}
}

type malPrecommitService struct {
	*core.Core
	interfaces.Precommiter
}

func (c *malPrecommitService) SendPrecommit(ctx context.Context, isNil bool) {
	var precommit *message.Precommit
	self, csize := selfAndCsize(c.Core, c.Height().Uint64())
	if isNil {
		precommit = message.NewPrecommit(c.Round(), c.Height().Uint64(), common.Hash{}, c.Backend().Sign, self, csize)
	} else {
		precommit = message.NewPrecommit(c.Round(), c.Height().Uint64(), common.HexToHash("0xCAFE"), c.Backend().Sign, self, csize)
	}
	c.SetSentPrecommit()
	c.BroadcastAll(precommit)
}

func TestMaliciousPrecommitSender(t *testing.T) {
	users, err := e2e.Validators(t, 6, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	//set Malicious users
	users[0].TendermintServices = &interfaces.Services{Precommiter: newMalPrecommitService}
	// creates a network of 6 users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should be up and continue to mine blocks
	err = network.WaitToMineNBlocks(10, 120, false)
	require.NoError(t, err, "Network should be mining new blocks now, but it's not")
}

func TestMaliciousSenderDisc(t *testing.T) {
	users, err := e2e.Validators(t, 4, "10e18,v,100,0.0.0.0:%s,%s,%s,%s")
	require.NoError(t, err)

	users[0].TendermintServices = &interfaces.Services{Precommiter: newMalPrecommitService}
	users[1].TendermintServices = &interfaces.Services{Precommiter: newMalPrecommitService}

	// creates a network of users and starts all the nodes in it
	network, err := e2e.NewNetworkFromValidators(t, users, true)
	require.NoError(t, err)
	defer network.Shutdown(t)

	// network should not be able to mine blocks
	err = network.WaitToMineNBlocks(1, 120, false)
	require.Error(t, err, "Network is not supposed to be mining blocks at this point")
}
