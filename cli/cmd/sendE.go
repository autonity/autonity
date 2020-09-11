/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	ethereum "github.com/clearmatics/autonity"
	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/cli/autonitybindings"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/spf13/cobra"
)

var (
	senderKey        string
	value            string
	recipientAddress string
)

// sendECmd represents the sendE command
var sendECmd = &cobra.Command{
	Use:   "sendE",
	Short: "Sends E from one account to another",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendTx(httpRpcAddress, senderKey, func(client *ethclient.Client, opts *bind.TransactOpts, a *autonitybindings.Autonity) (*types.Transaction, error) {

			v, err := ParseUint(value)
			if err != nil {
				return nil, err
			}
			opts.Value = v

			recipient, err := ParseAddress(recipientAddress)
			if err != nil {
				return nil, err
			}

			nonce, err := client.PendingNonceAt(opts.Context, opts.From)
			if err != nil {
				return nil, fmt.Errorf("failed to retrieve account nonce: %v", err)
			}
			// Figure out the gas allowance and gas price values
			gasPrice, err := client.SuggestGasPrice(opts.Context)
			if err != nil {
				return nil, fmt.Errorf("failed to suggest gas price: %v", err)
			}

			// If the contract surely has code (or code is not needed), estimate the transaction
			msg := ethereum.CallMsg{From: opts.From, To: &recipient, GasPrice: gasPrice, Value: v}
			gasLimit, err := client.EstimateGas(opts.Context, msg)
			if err != nil {
				return nil, fmt.Errorf("failed to estimate gas needed: %v", err)
			}
			// Create the transaction, sign it and schedule it for execution
			rawTx := types.NewTransaction(nonce, recipient, v, gasLimit, gasPrice, nil)
			signedTx, err := opts.Signer(types.HomesteadSigner{}, opts.From, rawTx)
			if err != nil {
				return nil, err
			}
			if err := client.SendTransaction(opts.Context, signedTx); err != nil {
				return nil, err
			}
			return signedTx, nil

		})
	},
}

func init() {
	rootCmd.AddCommand(sendECmd)

	sendECmd.Flags().StringVarP(&senderKey, "sender-key", "k", "", "hex encoded private key of the sender, without initial '0x'")
	sendECmd.Flags().StringVarP(&recipientAddress, "recipient", "r", "", "hex encoded address of the recipient")
	sendECmd.Flags().StringVarP(&value, "value", "v", "", "hex, scientific or decimal notation postive integer")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sendECmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sendECmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
