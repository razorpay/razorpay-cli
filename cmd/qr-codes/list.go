package qrcodes

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all QR codes",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}
		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			q.Set("count", fmt.Sprintf("%d", count))
		}
		if skip, _ := cmd.Flags().GetInt("skip"); skip > 0 {
			q.Set("skip", fmt.Sprintf("%d", skip))
		}
		if from, _ := cmd.Flags().GetInt64("from"); from > 0 {
			q.Set("from", fmt.Sprintf("%d", from))
		}
		if to, _ := cmd.Flags().GetInt64("to"); to > 0 {
			q.Set("to", fmt.Sprintf("%d", to))
		}
		if customerID, _ := cmd.Flags().GetString("customer-id"); customerID != "" {
			q.Set("customer_id", customerID)
		}
		if paymentID, _ := cmd.Flags().GetString("payment-id"); paymentID != "" {
			q.Set("payment_id", paymentID)
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

	listCmd.Flags().Int("count", 10, "Number of QR codes to fetch (max 100)")
	listCmd.Flags().Int("skip", 0, "Number of QR codes to skip")
	listCmd.Flags().Int64("from", 0, "Unix timestamp: fetch QR codes created after this time")
	listCmd.Flags().Int64("to", 0, "Unix timestamp: fetch QR codes created before this time")
	listCmd.Flags().String("customer-id", "", "Filter by customer ID")
	listCmd.Flags().String("payment-id", "", "Filter by payment ID")
}
