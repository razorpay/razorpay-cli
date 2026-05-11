package cmd

import (
	"fmt"
	"os"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "razorpay",
	Short: "Command-line interface for the Razorpay API",
	Long: `The Razorpay CLI provides command-line access to the Razorpay API.

To get started, configure your API credentials:

  razorpay configure

Then run any resource command, for example:

  razorpay payments list
  razorpay orders create --amount 50000 --currency INR

For help on a specific command, run:

  razorpay <command> --help`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func newClient() *api.Client {
	config.Init()
	return api.New(config.KeyID(), config.KeySecret())
}

func handleErr(err error) {
	fmt.Fprintln(os.Stderr, "Error:", err)
	os.Exit(1)
}

func init() {
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(paymentsCmd)
	rootCmd.AddCommand(ordersCmd)
	rootCmd.AddCommand(customersCmd)
	rootCmd.AddCommand(refundsCmd)
	rootCmd.AddCommand(settlementsCmd)
	rootCmd.AddCommand(disputesCmd)
}
