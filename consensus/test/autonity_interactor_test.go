package test

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"strconv"

	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethclient"
)

// interact is the starting point to constructing a tool to interact with
// autonity. It accepts an rpc port to connect to a local autonity instance and
// then, through a sequence of method calls allows interaction with autonity.
// Any error encountered along the way is propagated to the end of the chain of
// method calls. All of the components of the system are designed to be single
// use only.
//
// Example usages
// err := autonityInstance(rpcPort).tx(senderKey).addValidator(address ,stake, enode)
//
// whitelist,err := autonityInstance(rpcPort).call(blockNumber).getWhitelist()
func interact(rpcPort int) *interactor {
	i := &interactor{}
	client, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(rpcPort))
	if err != nil {
		i.err = err
		return i
	}

	instance, err := NewAutonity(autonity.ContractAddress, client)
	if err != nil {
		i.err = err
		return i
	}

	i.client = client
	i.instance = instance
	return i
}

// interactor serves as the root object for interacting with autonity in a
// builder fashion. E.G. build().configure().do() any error encountered at any
// step is propagated onwards through the chain of function calls and returned
// at the end. interactor manages an autonity instance and an ethereum client
// and provides functions to return objects that can handle sending
// transactions or making contract calls.
type interactor struct {
	instance *Autonity
	client   *ethclient.Client
	err      error
}

// close closes the underlying ethclient. Users of interacotor should not
// need to call this directly because the call that finally interacts with the
// autonity contract will call this.
func (i *interactor) close() {
	// in the case that Dial failed a.client could be nil
	if i.client != nil {
		i.client.Close()
	}
}

// tx returns a transactor through which transactions can be executed.
func (i *interactor) tx(operatorKey *ecdsa.PrivateKey) *transactor {
	t := &transactor{}
	if i.err != nil {
		t.err = i.err
		return t
	}

	operatorAddress := crypto.PubkeyToAddress(operatorKey.PublicKey)
	nonce, err := i.client.PendingNonceAt(context.Background(), operatorAddress)
	if err != nil {
		t.err = err
		return t
	}
	gasPrice, err := i.client.SuggestGasPrice(context.Background())
	if err != nil {
		t.err = err
		return t
	}

	txOpt := bind.NewKeyedTransactor(operatorKey)
	txOpt.From = operatorAddress
	txOpt.Nonce = new(big.Int).SetUint64(nonce)
	txOpt.GasLimit = uint64(300000000)
	txOpt.GasPrice = gasPrice

	t.opts = txOpt
	t.i = i
	return t
}

// call returns a caller through which we can call methods of the autonity
// contract.
func (i *interactor) call(blockNumber uint64) *caller {
	c := &caller{}
	if i.err != nil {
		c.err = i.err
		return c
	}
	c.i = i
	c.opts = &bind.CallOpts{
		Pending:     false,
		From:        common.Address{},
		BlockNumber: new(big.Int).SetUint64(blockNumber),
		Context:     context.Background(),
	}
	return c
}

// The transactor provides a mechanism to send Transactions to the
// autonity contract, it is intended to be constructed by interactor and
// provides a generic execute method which calls a callback with an autonity instance
// and transaction opts. It also can provide other methods to call specific functions
// on the autonity contract.
type transactor struct {
	err  error
	i    *interactor
	opts *bind.TransactOpts
}

// execute calls the given callback and ensures that its interactor instance
// is subsequently closed.
func (t *transactor) execute(action func(instance *Autonity, opts *bind.TransactOpts) error) error {
	defer t.i.close()
	if t.err != nil {
		return t.err
	}
	return action(t.i.instance, t.opts)
}

func (t *transactor) addValidator(validatorAddress common.Address, stakeBalance *big.Int, enode string) error {
	return t.execute(func(instance *Autonity, opts *bind.TransactOpts) error {
		_, err := instance.AddValidator(opts, validatorAddress, stakeBalance, enode)
		return err
	})
}

func (t *transactor) addStakeHolder(address common.Address, stakeBalance *big.Int, enode string) error {
	return t.execute(func(instance *Autonity, opts *bind.TransactOpts) error {
		_, err := instance.AddStakeholder(opts, address, enode, stakeBalance)
		return err
	})
}

func (t *transactor) addParticipant(address common.Address, enode string) error {
	return t.execute(func(instance *Autonity, opts *bind.TransactOpts) error {
		_, err := instance.AddParticipant(opts, address, enode)
		return err
	})
}

func (t *transactor) removeUser(address common.Address) error {
	return t.execute(func(instance *Autonity, opts *bind.TransactOpts) error {
		_, err := instance.RemoveUser(opts, address)
		return err
	})
}

func (t *transactor) mintStake(address common.Address, stake *big.Int) error {
	return t.execute(func(instance *Autonity, opts *bind.TransactOpts) error {
		_, err := instance.MintStake(opts, address, stake)
		return err
	})
}

func (t *transactor) redeemStake(address common.Address, stake *big.Int) error {
	return t.execute(func(instance *Autonity, opts *bind.TransactOpts) error {
		_, err := instance.RedeemStake(opts, address, stake)
		return err
	})
}

func (t *transactor) sendStake(address common.Address, stake *big.Int) error {
	return t.execute(func(instance *Autonity, opts *bind.TransactOpts) error {
		_, err := instance.Send(opts, address, stake)
		return err
	})
}

// The caller provides a mechanism to call the autonity contract, it is
// intended to be constructed by interactor and provides a generic execute
// method which calls a callback with an autonity instance and call opts. It
// also can provide other methods to call specific functions on the autonity
// contract.
type caller struct {
	err  error
	i    *interactor
	opts *bind.CallOpts
}

// execute calls the given callback and ensures that its interactor instance
// is subsequently closed.
func (c *caller) execute(action func(instance *Autonity, opts *bind.CallOpts) error) error {
	defer c.i.close()
	if c.err != nil {
		return c.err
	}
	return action(c.i.instance, c.opts)
}

func (c *caller) getWhitelist() ([]string, error) {
	var whitelist []string
	err := c.execute(func(instance *Autonity, opts *bind.CallOpts) error {
		w, err := instance.GetWhitelist(opts)
		whitelist = w
		return err
	})
	return whitelist, err
}

func (c *caller) dumpEconomicsMetricData() (Struct1, error) {
	var metrics Struct1
	err := c.execute(func(instance *Autonity, opts *bind.CallOpts) error {
		m, err := instance.DumpEconomicsMetricData(opts)
		metrics = m
		return err
	})
	return metrics, err
}

func (c *caller) checkMember(address common.Address) (bool, error) {
	var isMember bool
	err := c.execute(func(instance *Autonity, opts *bind.CallOpts) error {
		b, err := instance.CheckMember(opts, address)
		isMember = b
		return err
	})
	return isMember, err
}
