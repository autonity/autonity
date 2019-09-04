// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package external

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/clearmatics/autonity"
	"github.com/clearmatics/autonity/accounts"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/common/hexutil"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/event"
	"github.com/clearmatics/autonity/internal/ethapi"
	"github.com/clearmatics/autonity/log"
	"github.com/clearmatics/autonity/rpc"
	"github.com/clearmatics/autonity/signer/core"
)

type Backend struct {
	signers []accounts.Wallet
}

func (eb *Backend) Wallets() []accounts.Wallet {
	return eb.signers
}

func NewExternalBackend(endpoint string) (*Backend, error) {
	signer, err := NewExternalSigner(endpoint)
	if err != nil {
		return nil, err
	}
	return &Backend{
		signers: []accounts.Wallet{signer},
	}, nil
}

func (eb *Backend) Subscribe(sink chan<- accounts.WalletEvent) event.Subscription {
	return event.NewSubscription(func(quit <-chan struct{}) error {
		<-quit
		return nil
	})
}

// Signer provides an API to interact with an external signer (clef)
// It proxies request to the external signer while forwarding relevant
// request headers
type Signer struct {
	client   *rpc.Client
	endpoint string
	status   string
	cacheMu  sync.RWMutex
	cache    []accounts.Account
}

func NewExternalSigner(endpoint string) (*Signer, error) {
	client, err := rpc.Dial(endpoint)
	if err != nil {
		return nil, err
	}
	extsigner := &Signer{
		client:   client,
		endpoint: endpoint,
	}
	// Check if reachable
	version, err := extsigner.pingVersion()
	if err != nil {
		return nil, err
	}
	extsigner.status = fmt.Sprintf("ok [version=%v]", version)
	return extsigner, nil
}

func (api *Signer) URL() accounts.URL {
	return accounts.URL{
		Scheme: "extapi",
		Path:   api.endpoint,
	}
}

func (api *Signer) Status() (string, error) {
	return api.status, nil
}

func (api *Signer) Open(passphrase string) error {
	return fmt.Errorf("operation not supported on external signers")
}

func (api *Signer) Close() error {
	return fmt.Errorf("operation not supported on external signers")
}

func (api *Signer) Accounts() []accounts.Account {
	res, err := api.listAccounts()
	accnts := make([]accounts.Account, 0, len(res))
	if err != nil {
		log.Error("account listing failed", "error", err)
		return accnts
	}
	for _, addr := range res {
		accnts = append(accnts, accounts.Account{
			URL: accounts.URL{
				Scheme: "extapi",
				Path:   api.endpoint,
			},
			Address: addr,
		})
	}
	api.cacheMu.Lock()
	api.cache = accnts
	api.cacheMu.Unlock()
	return accnts
}

func (api *Signer) Contains(account accounts.Account) bool {
	api.cacheMu.RLock()
	defer api.cacheMu.RUnlock()
	for _, a := range api.cache {
		if a.Address == account.Address && (account.URL == (accounts.URL{}) || account.URL == api.URL()) {
			return true
		}
	}
	return false
}

func (api *Signer) Derive(path accounts.DerivationPath, pin bool) (accounts.Account, error) {
	return accounts.Account{}, fmt.Errorf("operation not supported on external signers")
}

func (api *Signer) SelfDerive(bases []accounts.DerivationPath, chain ethereum.ChainStateReader) {
	log.Error("operation SelfDerive not supported on external signers")
}

// SignData signs keccak256(data). The mimetype parameter describes the type of data being signed
func (api *Signer) SignData(account accounts.Account, mimeType string, data []byte) ([]byte, error) {
	var res hexutil.Bytes
	var signAddress = common.NewMixedcaseAddress(account.Address)
	if err := api.client.Call(&res, "account_signData",
		mimeType,
		&signAddress, // Need to use the pointer here, because of how MarshalJSON is defined
		hexutil.Encode(data)); err != nil {
		return nil, err
	}
	// If V is on 27/28-form, convert to to 0/1 for Clique
	if mimeType == accounts.MimetypeClique && (res[64] == 27 || res[64] == 28) {
		res[64] -= 27 // Transform V from 27/28 to 0/1 for Clique use
	}
	return res, nil
}

func (api *Signer) SignText(account accounts.Account, text []byte) ([]byte, error) {
	var res hexutil.Bytes
	var signAddress = common.NewMixedcaseAddress(account.Address)
	if err := api.client.Call(&res, "account_signData",
		accounts.MimetypeTextPlain,
		&signAddress, // Need to use the pointer here, because of how MarshalJSON is defined
		hexutil.Encode(text)); err != nil {
		return nil, err
	}
	return res, nil
}

func (api *Signer) SignTx(account accounts.Account, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	res := ethapi.SignTransactionResult{}
	data := hexutil.Bytes(tx.Data())
	var to *common.MixedcaseAddress
	if tx.To() != nil {
		t := common.NewMixedcaseAddress(*tx.To())
		to = &t
	}
	args := &core.SendTxArgs{
		Data:     &data,
		Nonce:    hexutil.Uint64(tx.Nonce()),
		Value:    hexutil.Big(*tx.Value()),
		Gas:      hexutil.Uint64(tx.Gas()),
		GasPrice: hexutil.Big(*tx.GasPrice()),
		To:       to,
		From:     common.NewMixedcaseAddress(account.Address),
	}
	if err := api.client.Call(&res, "account_signTransaction", args); err != nil {
		return nil, err
	}
	return res.Tx, nil
}

func (api *Signer) SignTextWithPassphrase(account accounts.Account, passphrase string, text []byte) ([]byte, error) {
	return []byte{}, fmt.Errorf("passphrase-operations not supported on external signers")
}

func (api *Signer) SignTxWithPassphrase(account accounts.Account, passphrase string, tx *types.Transaction, chainID *big.Int) (*types.Transaction, error) {
	return nil, fmt.Errorf("passphrase-operations not supported on external signers")
}
func (api *Signer) SignDataWithPassphrase(account accounts.Account, passphrase, mimeType string, data []byte) ([]byte, error) {
	return nil, fmt.Errorf("passphrase-operations not supported on external signers")
}

func (api *Signer) listAccounts() ([]common.Address, error) {
	var res []common.Address
	if err := api.client.Call(&res, "account_list"); err != nil {
		return nil, err
	}
	return res, nil
}

func (api *Signer) pingVersion() (string, error) {
	var v string
	if err := api.client.Call(&v, "account_version"); err != nil {
		return "", err
	}
	return v, nil
}
