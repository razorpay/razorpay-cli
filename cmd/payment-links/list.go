package paymentlinks

import (
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all payment links",
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}
		if paymentID, _ := cmd.Flags().GetString("payment-id"); paymentID != "" {
			q.Set("payment_id", paymentID)
		}
		if referenceID, _ := cmd.Flags().GetString("reference-id"); referenceID != "" {
			q.Set("reference_id", referenceID)
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

	listCmd.Flags().String("payment-id", "", "Filter by payment ID")
	listCmd.Flags().String("reference-id", "", "Filter by reference ID set at link creation")
}
