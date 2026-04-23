package route

import (
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var fetchTransfersByOrderCmd = &cobra.Command{
	Use:   "fetch-by-order <order_id>",
	Short: "Fetch transfers linked to an order",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		q := url.Values{
			"expand[]": {"transfers"},
		}
		if status, _ := cmd.Flags().GetString("status"); status != "" {
			q.Set("status", status)
		}
		data, err := client.Get("/orders/"+args[0], q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	transfersCmd.AddCommand(fetchTransfersByOrderCmd)

	fetchTransfersByOrderCmd.Flags().String("status", "", "Filter transfers by status (e.g. processing)")
}
