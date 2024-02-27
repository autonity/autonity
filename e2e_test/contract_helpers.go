package e2e

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/autonity/autonity/params"

	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/ethclient"
)

// Interact is the starting point to constructing a tool to Interact with
// autonity. It accepts an rpc port to connect to a local autonity instance and
// then, through a sequence of method calls allows interaction with autonity.
// Any error encountered along the way is propagated to the end of the chain of
// method calls. All of the components of the system are designed to be single
// use only.
//
// Example usages
// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// defer cancel()
// err := Interact(url).TX(senderKey, ctx).registerValidator(address ,stake, enode)
//
// whitelist,err := Interact(rpcPort).Call(blockNumber).getWhitelist()
func Interact(url string) *Interactor {
	i := &Interactor{}
	client, err := ethclient.Dial(url)
	if err != nil {
		i.err = err
		return i
	}

	instance, err := autonity.NewAutonity(params.AutonityContractAddress, client)
	if err != nil {
		i.err = err
		return i
	}

	i.client = client
	i.instance = instance
	return i
}

// Interactor serves as the root object for interacting with autonity in a
// builder fashion. E.G. build().configure().do() any error encountered at any
// step is propagated onwards through the chain of function calls and returned
// at the end. Interactor manages an autonity instance and an ethereum client
// and provides functions to return objects that can handle sending
// transactions or making contract calls.
type Interactor struct {
	instance *autonity.Autonity
	client   *ethclient.Client
	err      error
}

// Close closes the underlying ethclient. Validators of interacotor should not
// need to Call this directly because the Call that finally interacts with the
// autonity contract will Call this.
func (i *Interactor) Close() {
	// in the case that Dial failed a.client could be nil
	if i.client != nil {
		i.client.Close()
	}
}

// TX returns a transactor through which transactions can be executed.
func (i *Interactor) TX(ctx context.Context, senderKey *ecdsa.PrivateKey) *transactor {
	t := &transactor{}
	if i.err != nil {
		t.err = i.err
		return t
	}

	sender := crypto.PubkeyToAddress(senderKey.PublicKey)
	nonce, err := i.client.PendingNonceAt(context.Background(), sender)
	if err != nil {
		t.err = err
		return t
	}

	// The test framework uses 1337 as chain ID.
	chainID := new(big.Int).SetUint64(1337)
	txOpt, err := bind.NewKeyedTransactorWithChainID(senderKey, chainID)
	if err != nil {
		t.err = err
		return t
	}
	txOpt.From = sender
	txOpt.Nonce = new(big.Int).SetUint64(nonce)
	txOpt.Context = ctx

	t.opts = txOpt
	t.i = i
	return t
}

// Call returns a Caller through which we can Call methods of the autonity
// contract.
func (i *Interactor) Call(blockNumber *big.Int) *Caller {
	c := &Caller{}
	if i.err != nil {
		c.err = i.err
		return c
	}
	c.i = i
	c.opts = &bind.CallOpts{
		Pending:     false,
		From:        common.Address{},
		BlockNumber: blockNumber,
		Context:     context.Background(),
	}
	return c
}

// The transactor provides a mechanism to send Transactions to the
// autonity contract, it is intended to be constructed by Interactor and
// provides a generic execute method which calls a callback with an autonity instance
// and transaction opts. It also can provide other methods to Call specific functions
// on the autonity contract.
type transactor struct {
	err      error
	i        *Interactor
	opts     *bind.TransactOpts
	executed []*types.Transaction
}

func (t *transactor) analyze(test *testing.T) {
	test.Log("printing executed txs receipts", "count", len(t.executed))
	for i, tx := range t.executed {
		rec, err := t.i.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			test.Log("error tx", "id", i, "err", err)
		}
		outjson, err := json.MarshalIndent(rec, "", "\t")
		require.NoError(test, err)
		fmt.Println(string(outjson))
	}
}

// execute calls the given callback
func (t *transactor) execute(action func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error)) (*types.Transaction, error) {
	if t.err != nil {
		return nil, t.err
	}
	tx, err := action(t.i.instance, t.opts)
	t.opts.Nonce.Add(t.opts.Nonce, common.Big1)
	t.executed = append(t.executed, tx)
	return tx, err
}

// public writers
func (t *transactor) registerValidator(enode string, oracleAddress common.Address, consensusKey, signatures []byte) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.RegisterValidator(opts, enode, oracleAddress, consensusKey, signatures)
	})
}

func (t *transactor) bond(validator common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.Bond(opts, validator, amount)
	})
}

func (t *transactor) unbond(validator common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.Unbond(opts, validator, amount)
	})
}

// system operator writers
func (t *transactor) upgradeContract(byteCode []byte, abi string) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		opts.GasLimit = 100000000
		return instance.UpgradeContract(opts, byteCode, abi)
	})
}

func (t *transactor) completeContractUpgrade() (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.CompleteContractUpgrade(opts)
	})
}

func (t *transactor) setMinBaseFee(fee *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.SetMinimumBaseFee(opts, fee)
	})
}

func (t *transactor) setCommitteeSize(size *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.SetCommitteeSize(opts, size)
	})
}

func (t *transactor) setUnBondingPeriod(period *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.SetUnbondingPeriod(opts, period)
	})
}

func (t *transactor) setEpochPeriod(period *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.SetEpochPeriod(opts, period)
	})
}

func (t *transactor) setOperator(address common.Address) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.SetOperatorAccount(opts, address)
	})
}

func (t *transactor) setTreasuryAccount(address common.Address) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.SetTreasuryAccount(opts, address)
	})
}

func (t *transactor) setTreasuryFee(fee *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.SetTreasuryFee(opts, fee)
	})
}

func (t *transactor) mint(address common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.Mint(opts, address, amount)
	})
}

func (t *transactor) burn(address common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.Burn(opts, address, amount)
	})
}

// ERC-20 writers
func (t *transactor) transfer(recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.Transfer(opts, recipient, amount)
	})
}

func (t *transactor) transferFrom(sender common.Address, recipient common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.TransferFrom(opts, sender, recipient, amount)
	})
}

func (t *transactor) approve(spender common.Address, amount *big.Int) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.Approve(opts, spender, amount)
	})
}

// The Caller provides a mechanism to Call the autonity contract, it is
// intended to be constructed by Interactor and provides a generic execute
// method which calls a callback with an autonity instance and Call opts. It
// also can provide other methods to Call specific functions on the autonity
// contract.
type Caller struct {
	err  error
	i    *Interactor
	opts *bind.CallOpts
}

// execute calls the given callback and ensures that its Interactor instance
// is subsequently closed.
func (c *Caller) execute(action func(instance *autonity.Autonity, opts *bind.CallOpts) error) error {
	defer c.i.Close()
	if c.err != nil {
		return c.err
	}
	return action(c.i.instance, c.opts)
}

func (c *Caller) Name() (string, error) {
	var name string
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		n, err := instance.Name(opts)
		name = n
		return err
	})
	return name, err
}

func (c *Caller) Symbol() (string, error) {
	var symbol string
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		s, err := instance.Symbol(opts)
		symbol = s
		return err
	})
	return symbol, err
}

func (c *Caller) BalanceOf(address common.Address) (*big.Int, error) {
	var newton *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		s, err := instance.BalanceOf(opts, address)
		newton = s
		return err
	})
	return newton, err
}

func (c *Caller) Allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	var al *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		a, err := instance.Allowance(opts, owner, spender)
		al = a
		return err
	})
	return al, err
}

func (c *Caller) TotalSupply() (*big.Int, error) {
	var total *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		t, err := instance.TotalSupply(opts)
		total = t
		return err
	})
	return total, err
}

func (c *Caller) GetVersion() (uint64, error) {
	var version uint64
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		v, err := instance.GetVersion(opts)
		version = v.Uint64()
		return err
	})
	return version, err
}

func (c *Caller) GetCommittee() ([]autonity.AutonityCommitteeMember, error) {
	var committee []autonity.AutonityCommitteeMember
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		cm, err := instance.GetCommittee(opts)
		committee = cm
		return err
	})
	return committee, err
}

func (c *Caller) GetValidators() ([]common.Address, error) {
	var validators []common.Address
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		vals, err := instance.GetValidators(opts)
		validators = vals
		return err
	})
	return validators, err
}

func (c *Caller) GetValidator(address common.Address) (autonity.AutonityValidator, error) {
	var val autonity.AutonityValidator
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		v, err := instance.GetValidator(opts, address)
		val = v
		return err
	})
	return val, err
}

func (c *Caller) GetMaxCommitteeSize() (*big.Int, error) {
	var size *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		s, err := instance.GetMaxCommitteeSize(opts)
		size = s
		return err
	})
	return size, err
}

func (c *Caller) GetCommitteeEnodes() ([]string, error) {
	var eNodes []string
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		es, err := instance.GetCommitteeEnodes(opts)
		eNodes = es
		return err
	})
	return eNodes, err
}

func (c *Caller) GetMinBaseFee() (*big.Int, error) {
	var fee *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		f, err := instance.GetMinimumBaseFee(opts)
		fee = f
		return err
	})
	return fee, err
}

func (c *Caller) GetUnbondingPeriod() (*big.Int, error) {
	var period *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		p, err := instance.GetUnbondingPeriod(opts)
		period = p
		return err
	})
	return period, err
}

func (c *Caller) GetEpochPeriod() (*big.Int, error) {
	var period *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		p, err := instance.GetEpochPeriod(opts)
		period = p
		return err
	})
	return period, err
}

func (c *Caller) GetTreasuryFee() (*big.Int, error) {
	var fee *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		f, err := instance.GetTreasuryFee(opts)
		fee = f
		return err
	})
	return fee, err
}

func (c *Caller) GetTreasuryAccount() (common.Address, error) {
	var treasury common.Address
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		t, err := instance.GetTreasuryAccount(opts)
		treasury = t
		return err
	})
	return treasury, err
}

func (c *Caller) GetOperator() (common.Address, error) {
	var operator common.Address
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		o, err := instance.GetOperator(opts)
		operator = o
		return err
	})
	return operator, err
}

func (c *Caller) GetNewContract() ([]byte, string, error) {
	var byteCode []byte
	var abi string
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		b, a, err := instance.GetNewContract(opts)
		byteCode = b
		abi = a
		return err
	})
	return byteCode, abi, err
}

func (n *Node) AwaitContractTXN(txFunc func(ctx context.Context, client *Interactor) (*types.Transaction, error), timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	tx, err := txFunc(ctx, n.Interactor)
	if err != nil {
		return err
	}

	return n.AwaitTransactions(ctx, tx)
}

func (n *Node) AwaitMintNTN(operatorKey *ecdsa.PrivateKey, receiver common.Address, amount *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, operatorKey).mint(receiver, amount)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitBurnNTN(optKey *ecdsa.PrivateKey, from common.Address, amount *big.Int, tm time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, optKey).burn(from, amount)
	}
	return n.AwaitContractTXN(txFunc, tm)
}

func (n *Node) AwaitNTNApprove(owner *ecdsa.PrivateKey, spender common.Address, amount *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, owner).approve(spender, amount)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitNTNTransferFrom(spenderKey *ecdsa.PrivateKey, owner common.Address, amount *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, spenderKey).transferFrom(owner, crypto.PubkeyToAddress(spenderKey.PublicKey), amount)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitTransferNTN(senderKey *ecdsa.PrivateKey, receiver common.Address, amount *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, senderKey).transfer(receiver, amount)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitRegisterValidator(validator *ecdsa.PrivateKey, enode string, oracle common.Address, consensusKey []byte, pop []byte, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, validator).registerValidator(enode, oracle, consensusKey, pop)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitBondStake(delegator *ecdsa.PrivateKey, validator common.Address, amount *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, delegator).bond(validator, amount)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitUnbondStake(delegator *ecdsa.PrivateKey, validator common.Address, amount *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, delegator).unbond(validator, amount)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitSetMinBaseFee(optKey *ecdsa.PrivateKey, newFee *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, optKey).setMinBaseFee(newFee)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitSetCommitteeSize(optKey *ecdsa.PrivateKey, newSize *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, optKey).setCommitteeSize(newSize)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitSetUnbondingPeriod(optKey *ecdsa.PrivateKey, newPeriod *big.Int, timeout time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, optKey).setUnBondingPeriod(newPeriod)
	}
	return n.AwaitContractTXN(txFunc, timeout)
}

func (n *Node) AwaitSetEpochPeriod(optKey *ecdsa.PrivateKey, newPeriod *big.Int, tm time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, optKey).setEpochPeriod(newPeriod)
	}
	return n.AwaitContractTXN(txFunc, tm)
}

func (n *Node) AwaitSetTreasuryAccount(optKey *ecdsa.PrivateKey, newAccount common.Address, tm time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, optKey).setTreasuryAccount(newAccount)
	}
	return n.AwaitContractTXN(txFunc, tm)
}

func (n *Node) AwaitSetTreasuryFee(optKey *ecdsa.PrivateKey, newFee *big.Int, tm time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, optKey).setTreasuryFee(newFee)
	}
	return n.AwaitContractTXN(txFunc, tm)
}

func (n *Node) AwaitSetOperator(optKey *ecdsa.PrivateKey, newOperator common.Address, tm time.Duration) error {
	txFunc := func(ctx context.Context, client *Interactor) (*types.Transaction, error) {
		return client.TX(ctx, optKey).setOperator(newOperator)
	}
	return n.AwaitContractTXN(txFunc, tm)
}
