package subscriptions

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var invoicesCmd = &cobra.Command{
	Use:   "invoices <subscription_id>",
	Short: "Fetch all invoices for a subscription",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{}
		q.Set("subscription_id", args[0])
		if count, _ := cmd.Flags().GetInt("count"); count > 0 {
			q.Set("count", fmt.Sprintf("%d", count))
		}
		if skip, _ := cmd.Flags().GetInt("skip"); skip > 0 {
			q.Set("skip", fmt.Sprintf("%d", skip))
		}
		data, err := client.Get("/v1/invoices", q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(invoicesCmd)

	invoicesCmd.Flags().Int("count", 10, "Number of invoices to fetch")
	invoicesCmd.Flags().Int("skip", 0, "Number of invoices to skip")
}
