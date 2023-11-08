package test

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/autonity/autonity/accounts/abi/bind"
	"github.com/autonity/autonity/autonity"
	"github.com/autonity/autonity/common"
	"github.com/autonity/autonity/core/types"
	"github.com/autonity/autonity/crypto"
	"github.com/autonity/autonity/ethclient"
	"github.com/stretchr/testify/require"
)

// interact is the starting point to constructing a tool to interact with
// autonity. It accepts an rpc port to connect to a local autonity instance and
// then, through a sequence of method calls allows interaction with autonity.
// Any error encountered along the way is propagated to the end of the chain of
// method calls. All of the components of the system are designed to be single
// use only.
//
// Example usages
// err := interact(rpcPort).tx(senderKey).registerValidator(address ,stake, enode)
//
// whitelist,err := interact(rpcPort).call(blockNumber).getWhitelist()
func interact(rpcPort int) *interactor {
	i := &interactor{}
	client, err := ethclient.Dial("http://127.0.0.1:" + strconv.Itoa(rpcPort))
	if err != nil {
		i.err = err
		return i
	}

	instance, err := autonity.NewAutonity(autonity.AutonityContractAddress, client)
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
	instance *autonity.Autonity
	client   *ethclient.Client
	err      error
}

// close closes the underlying ethclient. Validators of interacotor should not
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

	chainId := new(big.Int).SetUint64(1337)
	txOpt, err := bind.NewKeyedTransactorWithChainID(operatorKey, chainId)
	if err != nil {
		t.err = err
		return t
	}
	txOpt.From = operatorAddress
	txOpt.Nonce = new(big.Int).SetUint64(nonce)

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
	err      error
	i        *interactor
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
func (t *transactor) registerValidator(enode string, oracleAddress common.Address, activityKey, signatures []byte) (*types.Transaction, error) {
	return t.execute(func(instance *autonity.Autonity, opts *bind.TransactOpts) (*types.Transaction, error) {
		return instance.RegisterValidator(opts, enode, oracleAddress, activityKey, signatures)
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
func (c *caller) execute(action func(instance *autonity.Autonity, opts *bind.CallOpts) error) error {
	defer c.i.close()
	if c.err != nil {
		return c.err
	}
	return action(c.i.instance, c.opts)
}

// erc-20 getters.
func (c *caller) name() (string, error) {
	var name string
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		n, err := instance.Name(opts)
		name = n
		return err
	})
	return name, err
}

func (c *caller) symbol() (string, error) {
	var symbol string
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		s, err := instance.Symbol(opts)
		symbol = s
		return err
	})
	return symbol, err
}

func (c *caller) balanceOf(address common.Address) (*big.Int, error) {
	var newton *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		s, err := instance.BalanceOf(opts, address)
		newton = s
		return err
	})
	return newton, err
}

func (c *caller) allowance(owner common.Address, spender common.Address) (*big.Int, error) {
	var al *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		a, err := instance.Allowance(opts, owner, spender)
		al = a
		return err
	})
	return al, err
}

func (c *caller) totalSupply() (*big.Int, error) {
	var total *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		t, err := instance.TotalSupply(opts)
		total = t
		return err
	})
	return total, err
}

// system setting / state getters
func (c *caller) getVersion() (uint64, error) {
	var version uint64
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		v, err := instance.GetVersion(opts)
		version = v.Uint64()
		return err
	})
	return version, err
}

func (c *caller) getCommittee() ([]autonity.AutonityCommitteeMember, error) {
	var committee []autonity.AutonityCommitteeMember
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		cm, err := instance.GetCommittee(opts)
		committee = cm
		return err
	})
	return committee, err
}

func (c *caller) getValidators() ([]common.Address, error) {
	var validators []common.Address
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		vals, err := instance.GetValidators(opts)
		validators = vals
		return err
	})
	return validators, err
}

func (c *caller) getValidator(address common.Address) (autonity.AutonityValidator, error) {
	var val autonity.AutonityValidator
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		v, err := instance.GetValidator(opts, address)
		val = v
		return err
	})
	return val, err
}

func (c *caller) getMaxCommitteeSize() (*big.Int, error) {
	var size *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		s, err := instance.GetMaxCommitteeSize(opts)
		size = s
		return err
	})
	return size, err
}

func (c *caller) getCommitteeEnodes() ([]string, error) {
	var eNodes []string
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		es, err := instance.GetCommitteeEnodes(opts)
		eNodes = es
		return err
	})
	return eNodes, err
}

func (c *caller) getMinBaseFee() (*big.Int, error) {
	var fee *big.Int
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		f, err := instance.GetMinimumBaseFee(opts)
		fee = f
		return err
	})
	return fee, err
}

func (c *caller) getOperator() (common.Address, error) {
	var operator common.Address
	err := c.execute(func(instance *autonity.Autonity, opts *bind.CallOpts) error {
		o, err := instance.GetOperator(opts)
		operator = o
		return err
	})
	return operator, err
}

func (c *caller) getNewContract() ([]byte, string, error) {
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
