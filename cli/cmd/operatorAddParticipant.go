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
	"github.com/clearmatics/autonity/accounts/abi/bind"
	"github.com/clearmatics/autonity/cli/autonitybindings"
	"github.com/clearmatics/autonity/core/types"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/clearmatics/autonity/p2p/enode"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	operatorKey      string
	participantEnode string
)

func addOperatorFlags(fs *pflag.FlagSet) {
	fs.StringVarP(&operatorKey, "operator-key", "k", "", "hex encoded private key of the operator, without initial '0x'")
}

// operatorAddParticipantCmd represents the operatorAddParticipant command
var operatorAddParticipantCmd = &cobra.Command{
	Use:   "operatorAddParticipant",
	Short: "Adds a the given enode as a network participant, resctricted to the operator account",
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendTx(httpRpcAddress, operatorKey, func(client *ethclient.Client, opts *bind.TransactOpts, a *autonitybindings.Autonity) (*types.Transaction, error) {

			addr, err := enode.ParseV4ToAddress(participantEnode)
			if err != nil {
				return nil, err
			}
			return a.AddParticipant(opts, addr, participantEnode)
		})
	},
}

func init() {
	rootCmd.AddCommand(operatorAddParticipantCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	addOperatorFlags(operatorAddParticipantCmd.Flags())
	//operatorAddParticipantCmd.Flags().StringVar("toggle", "t", false, "Help message for toggle")

	operatorAddParticipantCmd.Flags().StringVarP(&participantEnode, "enode", "p", "", "enode url of the participant")
	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//operatorAddParticipantCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
