package invoices

import (
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all invoices",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}

		if invoiceType, _ := cmd.Flags().GetString("type"); invoiceType != "" {
			q.Set("type", invoiceType)
		}
		if paymentID, _ := cmd.Flags().GetString("payment-id"); paymentID != "" {
			q.Set("payment_id", paymentID)
		}
		if receipt, _ := cmd.Flags().GetString("receipt"); receipt != "" {
			q.Set("receipt", receipt)
		}
		if customerID, _ := cmd.Flags().GetString("customer-id"); customerID != "" {
			q.Set("customer_id", customerID)
		}
		data, err := client.Get(basePath, q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(listCmd)

	listCmd.Flags().String("type", "", "Filter by type: invoice, link, or ecod")
	listCmd.Flags().String("payment-id", "", "Filter by payment ID")
	listCmd.Flags().String("receipt", "", "Filter by receipt number")
	listCmd.Flags().String("customer-id", "", "Filter by customer ID")
}
