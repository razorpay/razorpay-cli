package refunds

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var paymentRefundCmd = &cobra.Command{
	Use:   "payment-refund <payment_id> <refund_id>",
	Short: "Fetch a specific refund for a payment",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Get("/payments/"+args[0]+"/refunds/"+args[1], nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(paymentRefundCmd)
}
