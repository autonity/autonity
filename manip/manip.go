package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	// "github.com/ethereum/go-ethereum/rpc"

	"github.com/clearmatics/autonity/common"
	"github.com/clearmatics/autonity/ethclient"
	"github.com/spf13/cobra"
)

var (
	nodeAddress     string
	contractAddress string
)

func main() {
	Execute()
}

var rootCmd = &cobra.Command{
	Use:   "manip",
	Short: "Manipulate the autonity contract",
	Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at http://hugo.spf13.com`,
}

var getValidators = &cobra.Command{
	Use:   "getValidators",
	Short: "Print the validators in the autonity contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ethclient.Dial(nodeAddress)
		if err != nil {
			return err
		}
		address := common.HexToAddress(contractAddress)
		_, err = NewAutonity(address, client)
		if err != nil {
			return err
		}
		println("lsls")

		// c, err := rpc.DialHTTP("tcp", rpcEndpoint)
		// if err != nil {
		// 	return err
		// }
		// var a common.Address

		// NewAutonity(a, c)

		// var result []string
		// err = c.Call(result, "eth_call", nil)
		// if err != nil {
		// 	return err
		// }
		// c.Call(result, "eth_call", `[{"to":"0xebe8efa441b9302a0d7eaecc277c09d20d684540","data":"0x45848dfc"},"latest"]`)
		return nil
	},
}

var getFirstBlock = &cobra.Command{
	Use:   "getFirstBlock",
	Short: "Print the validators in the autonity contract",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := ethclient.Dial(nodeAddress)
		if err != nil {
			return err
		}
		b, err := client.BlockByNumber(context.Background(), big.NewInt(2))
		if err != nil {
			return err
		}
		fmt.Printf("%+v\n", b)

		// c, err := rpc.DialHTTP("tcp", rpcEndpoint)
		// if err != nil {
		// 	return err
		// }
		// var a common.Address

		// NewAutonity(a, c)

		// var result []string
		// err = c.Call(result, "eth_call", nil)
		// if err != nil {
		// 	return err
		// }
		// c.Call(result, "eth_call", `[{"to":"0xebe8efa441b9302a0d7eaecc277c09d20d684540","data":"0x45848dfc"},"latest"]`)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {

	rootCmd.AddCommand(getValidators)
	rootCmd.AddCommand(getFirstBlock)
	pf := rootCmd.PersistentFlags()

	pf.StringVar(&nodeAddress, "node-address", "", "the node address to send requests to")
	cobra.MarkFlagRequired(pf, "node-address")

	pf.StringVar(&contractAddress, "contract-address", "", "the contract address")
	// cobra.MarkFlagRequired(pf, "contract-address")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	// rootCmd.AddCommand(initCmd)
}
