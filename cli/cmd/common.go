package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/cli/autonitybindings"
	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/contracts/autonity"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/crypto"
	"github.com/clearmatics/autonity/ethclient"
)

// ParseUint provides support for parsing large numbers in base 10 using
// scientific notation with either e or E used as a prefix for the exponent.
// E.g 10.456e34 or 1234E5. If the resultant number cannot be parsed, is not an
// exact integer or does not fit into 256 bits an error is returned.
func ParseUint(str string) (*big.Int, error) {

	// We parse the string as a base 10 number and allow for 64 bit precision.
	// Providing the base means that the function will fail to parse anything
	// with a prefix, e.g 0x... 0b... 0X... 0B... .
	f, _, err := big.ParseFloat(str, 10, 256, big.ToZero)
	if err != nil {
		return nil, fmt.Errorf("could not parse %s: %v", str, err)

	}

	// Check to see the number is positive.
	if f.Cmp(big.NewFloat(0)) == -1 {
		return nil, fmt.Errorf("the nunber defined by %q is less than 0", str)
	}

	// Check to see if the number is an integer.
	result, accuracy := f.Int(nil)
	if accuracy != big.Exact {
		return nil, fmt.Errorf("the nunber defined by %q is not an integer", str)
	}
	return result, nil
}

func sendTx(nodeAddr string, senderKey string, action func(*ethclient.Client, *bind.TransactOpts, *autonitybindings.Autonity) (*types.Transaction, error)) error {
	client, err := ethclient.Dial("ws://" + nodeAddr)
	if err != nil {
		return err
	}
	defer client.Close()
	a, err := autonitybindings.NewAutonity(autonity.ContractAddress, client)
	if err != nil {
		return err
	}
	k, err := crypto.HexToECDSA(senderKey)
	if err != nil {
		return fmt.Errorf("failed to decode sender key: %v", err)
	}

	senderAddress := crypto.PubkeyToAddress(k.PublicKey)
	opts := bind.NewKeyedTransactor(k)
	opts.From = senderAddress
	opts.Context = context.Background()

	heads := make(chan *types.Header)
	sub, err := client.SubscribeNewHead(context.Background(), heads)
	if err != nil {
		return err
	}

	tx, err := action(client, opts, a)
	if err != nil {
		return err
	}

	txHash := tx.Hash()
	marshalled, err := json.MarshalIndent(tx, "", "  ")
	if err != nil {
		return err
	}
	fmt.Printf("Sent transaction:\n%s\n", string(marshalled))

	fmt.Print("\nAwaiting mining ...")

	for {
		select {
		case h := <-heads:
			b, err := client.BlockByHash(context.Background(), h.Hash())
			if err != nil {
				return err
			}
			for _, t := range b.Transactions() {
				if t.Hash() == txHash {
					r, err := client.TransactionReceipt(context.Background(), txHash)
					if err != nil {
						return err
					}
					marshalled, err := json.MarshalIndent(r, "", "  ")
					if err != nil {
						return err
					}
					fmt.Printf("\n\nTransaction receipt:\n%v\n", string(marshalled))
					return nil
				}
			}
			fmt.Print(".")

		case err := <-sub.Err():
			return err
		}
	}
}

func ParseAddress(s string) (common.Address, error) {
	if !common.IsHexAddress(s) {
		return common.Address{}, fmt.Errorf("string %q cannot represent a valid address", s)
	}
	return common.HexToAddress(s), nil
}
