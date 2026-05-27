package smartcollect

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var fetchUPIPaymentCmd = &cobra.Command{
	Use:   "fetch-upi-payment <payment_id>",
	Short: "Fetch UPI transfer details for a payment (Smart Collect 2)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		data, err := client.Get("/v1/payments/"+args[0]+"/upi_transfer", nil)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(fetchUPIPaymentCmd)
}
