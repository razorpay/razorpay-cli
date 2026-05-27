package smartcollect

import (
	"fmt"
	"net/url"

	"github.com/razorpay/razorpay-cli/api"
	"github.com/razorpay/razorpay-cli/cmd/cmdutil"
	"github.com/spf13/cobra"
)

var paymentsCmd = &cobra.Command{
	Use:   "payments <virtual_account_id>",
	Short: "Fetch payments made to a virtual account",
	Args:  cobra.ExactArgs(1),
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
		data, err := client.Get(basePath+"/"+args[0]+"/payments", q)
		if err != nil {
			cmdutil.HandleErr(err)
		}
		api.PrettyPrint(data)
		return nil
	},
}

func init() {
	Cmd.AddCommand(paymentsCmd)

	paymentsCmd.Flags().Int("count", 10, "Number of payments to fetch (max 100)")
	paymentsCmd.Flags().Int("skip", 0, "Number of payments to skip")
	paymentsCmd.Flags().Int64("from", 0, "Unix timestamp: fetch payments created after this time")
	paymentsCmd.Flags().Int64("to", 0, "Unix timestamp: fetch payments created before this time")
}
