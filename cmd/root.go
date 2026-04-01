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
	Short: "Razorpay CLI - interact with the Razorpay API from your terminal",
	Long: `A command-line interface for the Razorpay API.

Configure your API keys with:
  razorpay configure

Then use resource commands like:
  razorpay payments list
  razorpay orders create --param amount=5000 --param currency=INR`,
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
