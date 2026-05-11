package cmd

import (
	"os"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/razorpay/razorpay-cli/cmd/customers"
	"github.com/razorpay/razorpay-cli/cmd/disputes"
	"github.com/razorpay/razorpay-cli/cmd/documents"
	"github.com/razorpay/razorpay-cli/cmd/invoices"
	"github.com/razorpay/razorpay-cli/cmd/orders"
	paymentlinks "github.com/razorpay/razorpay-cli/cmd/payment-links"
	"github.com/razorpay/razorpay-cli/cmd/payments"
	qrcodes "github.com/razorpay/razorpay-cli/cmd/qr-codes"
	"github.com/razorpay/razorpay-cli/cmd/refunds"
	"github.com/razorpay/razorpay-cli/cmd/route"
	"github.com/razorpay/razorpay-cli/cmd/settlements"
	smartcollect "github.com/razorpay/razorpay-cli/cmd/smart-collect"
	"github.com/razorpay/razorpay-cli/cmd/subscriptions"
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
  razorpay orders create --amount 5000 --currency INR`,
}

// SetVersion stamps the root command with build-time version info injected
// by goreleaser via -ldflags "-X main.version=... -X main.commit=... -X main.date=..."
func SetVersion(version, commit, date string) {
	rootCmd.Version = version
	rootCmd.Long = rootCmd.Long + "\n\nVersion: " + version + " (" + commit + ") built " + date
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// newClient and handleErr are thin wrappers around cmdutil so that the
// remaining flat cmd/*.go files (payments, customers, etc.) need no changes.
func newClient() *api.Client {
	return cmdutil.NewClient()
}

func handleErr(err error) {
	cmdutil.HandleErr(err)
}

func init() {
	rootCmd.AddCommand(configureCmd)
	rootCmd.AddCommand(payments.Cmd)
	rootCmd.AddCommand(orders.Cmd)
	rootCmd.AddCommand(customers.Cmd)
	rootCmd.AddCommand(refunds.Cmd)
	rootCmd.AddCommand(settlements.Cmd)
	rootCmd.AddCommand(disputes.Cmd)
	rootCmd.AddCommand(documents.Cmd)
	rootCmd.AddCommand(paymentlinks.Cmd)
	rootCmd.AddCommand(qrcodes.Cmd)
	rootCmd.AddCommand(invoices.Cmd)
	rootCmd.AddCommand(subscriptions.Cmd)
	rootCmd.AddCommand(route.Cmd)
	rootCmd.AddCommand(smartcollect.Cmd)
}
