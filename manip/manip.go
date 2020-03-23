package main

import (
	"fmt"
	"os"

	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/cobra"
)

var (
	rpcEndpoint string
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
		//c, err := rpc.DialHTTP("tcp", rpcEndpoint)

		c, err := rpc.DialHTTP(rpcEndpoint)
		if err != nil {
			return err
		}

		var result []string
		err = c.Call(result, "eth_call", nil)
		if err != nil {
			return err
		}
		c.Call(result, "eth_call", `[{"to":"0xebe8efa441b9302a0d7eaecc277c09d20d684540","data":"0x45848dfc"},"latest"]`)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&rpcEndpoint, "rpc-endpoint", "", "the target rpc endpoint to send requests to")
	rootCmd.PersistentFlags().StringVar(&rpcEndpoint, "contract-address", "", "the contract address")
	// rootCmd.PersistentFlags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	// rootCmd.PersistentFlags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	// rootCmd.PersistentFlags().Bool("viper", true, "use Viper for configuration")
	// viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	// viper.BindPFlag("useViper", rootCmd.PersistentFlags().Lookup("viper"))
	// viper.SetDefault("author", "NAME HERE <EMAIL ADDRESS>")
	// viper.SetDefault("license", "apache")

	rootCmd.AddCommand(getValidators)
	// rootCmd.AddCommand(initCmd)
}
