package route

import (
	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var accountPaymentsCmd = &cobra.Command{
	Use:   "payments <account_id>",
	Short: "Fetch payments received by a linked account",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client := cmdutil.NewClient()
		headers := map[string]string{
			"X-Razorpay-Account": args[0],
		}
		data, err := client.GetWithHeaders("/v1/payments", nil, headers)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	accountsCmd.AddCommand(accountPaymentsCmd)
}
